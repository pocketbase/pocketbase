package forms

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/spf13/cast"
)

// RecordUpsert defines a Record upsert form.
type RecordUpsert struct {
	app    core.App
	record *models.Record

	isCreate      bool
	filesToDelete []string // names list
	filesToUpload []*rest.UploadedFile

	Data map[string]any `json:"data"`
}

// NewRecordUpsert creates a new Record upsert form.
// (pass a new Record model instance (`models.NewRecord(...)`) for create).
func NewRecordUpsert(app core.App, record *models.Record) *RecordUpsert {
	form := &RecordUpsert{
		app:           app,
		record:        record,
		isCreate:      !record.HasId(),
		filesToDelete: []string{},
		filesToUpload: []*rest.UploadedFile{},
	}

	form.Data = map[string]any{}
	for _, field := range record.Collection().Schema.Fields() {
		form.Data[field.Name] = record.GetDataValue(field.Name)
	}

	return form
}

func (form *RecordUpsert) getContentType(r *http.Request) string {
	t := r.Header.Get("Content-Type")
	for i, c := range t {
		if c == ' ' || c == ';' {
			return t[:i]
		}
	}
	return t
}

func (form *RecordUpsert) extractRequestData(r *http.Request) (map[string]any, error) {
	switch form.getContentType(r) {
	case "application/json":
		return form.extractJsonData(r)
	case "multipart/form-data":
		return form.extractMultipartFormData(r)
	default:
		return nil, errors.New("Unsupported request Content-Type.")
	}
}

func (form *RecordUpsert) extractJsonData(r *http.Request) (map[string]any, error) {
	result := map[string]any{}

	err := rest.ReadJsonBodyCopy(r, &result)

	return result, err
}

func (form *RecordUpsert) extractMultipartFormData(r *http.Request) (map[string]any, error) {
	result := map[string]any{}

	// parse form data (if not already)
	if err := r.ParseMultipartForm(rest.DefaultMaxMemory); err != nil {
		return result, err
	}

	arrayValueSupportTypes := schema.ArraybleFieldTypes()

	for key, values := range r.PostForm {
		if len(values) == 0 {
			result[key] = nil
			continue
		}

		field := form.record.Collection().Schema.GetFieldByName(key)
		if field != nil && list.ExistInSlice(field.Type, arrayValueSupportTypes) {
			result[key] = values
		} else {
			result[key] = values[0]
		}
	}

	return result, nil
}

func (form *RecordUpsert) normalizeData() error {
	for _, field := range form.record.Collection().Schema.Fields() {
		if v, ok := form.Data[field.Name]; ok {
			form.Data[field.Name] = field.PrepareValue(v)
		}
	}

	return nil
}

// LoadData loads and normalizes json OR multipart/form-data request data.
//
// File upload is supported only via multipart/form-data.
//
// To REPLACE previously uploaded file(s) you can suffix the field name
// with the file index (eg. `myfile.0`) and set the new value.
// For single file upload fields, you can skip the index and directly
// assign the file value to the field name (eg. `myfile`).
//
// To DELETE previously uploaded file(s) you can suffix the field name
// with the file index (eg. `myfile.0`) and set it to null or empty string.
// For single file upload fields, you can skip the index and directly
// reset the field using its field name (eg. `myfile`).
func (form *RecordUpsert) LoadData(r *http.Request) error {
	requestData, err := form.extractRequestData(r)
	if err != nil {
		return err
	}

	// extend base data with the extracted one
	extendedData := form.record.Data()
	rawData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(rawData, &extendedData); err != nil {
		return err
	}

	for _, field := range form.record.Collection().Schema.Fields() {
		key := field.Name
		value := extendedData[key]
		value = field.PrepareValue(value)

		if field.Type == schema.FieldTypeFile {
			options, _ := field.Options.(*schema.FileOptions)
			oldNames := list.ToUniqueStringSlice(form.Data[key])

			// delete previously uploaded file(s)
			if options.MaxSelect == 1 {
				// search for unset zero indexed key as a fallback
				indexedKeyValue, hasIndexedKey := extendedData[key+".0"]

				if cast.ToString(value) == "" || (hasIndexedKey && cast.ToString(indexedKeyValue) == "") {
					if len(oldNames) > 0 {
						form.filesToDelete = append(form.filesToDelete, oldNames...)
					}
					form.Data[key] = nil
				}
			} else if options.MaxSelect > 1 {
				// search for individual file index to delete (eg. "file.0")
				keyExp, _ := regexp.Compile(`^` + regexp.QuoteMeta(key) + `\.\d+$`)
				indexesToDelete := []int{}
				for indexedKey := range extendedData {
					if keyExp.MatchString(indexedKey) && cast.ToString(extendedData[indexedKey]) == "" {
						index, indexErr := strconv.Atoi(indexedKey[len(key)+1:])
						if indexErr != nil || index >= len(oldNames) {
							continue
						}
						indexesToDelete = append(indexesToDelete, index)
					}
				}

				// slice to fill only with the non-deleted indexes
				nonDeleted := []string{}
				for i, name := range oldNames {
					// not marked for deletion
					if !list.ExistInSlice(i, indexesToDelete) {
						nonDeleted = append(nonDeleted, name)
						continue
					}

					// store the id to actually delete the file later
					form.filesToDelete = append(form.filesToDelete, name)
				}
				form.Data[key] = nonDeleted
			}

			// check if there are any new uploaded form files
			files, err := rest.FindUploadedFiles(r, key)
			if err != nil {
				continue // skip invalid or missing file(s)
			}

			// refresh oldNames list
			oldNames = list.ToUniqueStringSlice(form.Data[key])

			if options.MaxSelect == 1 {
				// delete previous file(s) before replacing
				if len(oldNames) > 0 {
					form.filesToDelete = list.ToUniqueStringSlice(append(form.filesToDelete, oldNames...))
				}
				form.filesToUpload = append(form.filesToUpload, files[0])
				form.Data[key] = files[0].Name()
			} else if options.MaxSelect > 1 {
				// append the id of each uploaded file instance
				form.filesToUpload = append(form.filesToUpload, files...)
				for _, file := range files {
					oldNames = append(oldNames, file.Name())
				}
				form.Data[key] = oldNames
			}
		} else {
			form.Data[key] = value
		}
	}

	return form.normalizeData()
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordUpsert) Validate() error {
	dataValidator := validators.NewRecordDataValidator(
		form.app.Dao(),
		form.record,
		form.filesToUpload,
	)

	return dataValidator.Validate(form.Data)
}

// DrySubmit performs a form submit within a transaction and reverts it.
// For actual record persistence, check the `form.Submit()` method.
//
// This method doesn't handle file uploads/deletes or trigger any app events!
func (form *RecordUpsert) DrySubmit(callback func(txDao *daos.Dao) error) error {
	if err := form.Validate(); err != nil {
		return err
	}

	// bulk load form data
	if err := form.record.Load(form.Data); err != nil {
		return err
	}

	return form.app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		tx, ok := txDao.DB().(*dbx.Tx)
		if !ok {
			return errors.New("failed to get transaction db")
		}
		defer tx.Rollback()
		txDao.BeforeCreateFunc = nil
		txDao.AfterCreateFunc = nil
		txDao.BeforeUpdateFunc = nil
		txDao.AfterUpdateFunc = nil

		if err := txDao.SaveRecord(form.record); err != nil {
			return err
		}

		return callback(txDao)
	})
}

// Submit validates the form and upserts the form Record model.
func (form *RecordUpsert) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	// bulk load form data
	if err := form.record.Load(form.Data); err != nil {
		return err
	}

	return form.app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		// persist record model
		if err := txDao.SaveRecord(form.record); err != nil {
			return err
		}

		// upload new files (if any)
		if err := form.processFilesToUpload(); err != nil {
			return err
		}

		// delete old files (if any)
		if err := form.processFilesToDelete(); err != nil { //nolint:staticcheck
			// for now fail silently to avoid reupload when `form.Submit()`
			// is called manually (aka. not from an api request)...
		}

		return nil
	})
}

func (form *RecordUpsert) processFilesToUpload() error {
	if len(form.filesToUpload) == 0 {
		return nil // nothing to upload
	}

	if !form.record.HasId() {
		return errors.New("The record is not persisted yet.")
	}

	fs, err := form.app.NewFilesystem()
	if err != nil {
		return err
	}
	defer fs.Close()

	var uploadErrors []error
	for i := len(form.filesToUpload) - 1; i >= 0; i-- {
		file := form.filesToUpload[i]
		path := form.record.BaseFilesPath() + "/" + file.Name()

		if err := fs.Upload(file.Bytes(), path); err == nil {
			// remove the uploaded file from the list
			form.filesToUpload = append(form.filesToUpload[:i], form.filesToUpload[i+1:]...)
		} else {
			// store the upload error
			uploadErrors = append(uploadErrors, fmt.Errorf("File %d: %v", i, err))
		}
	}

	if len(uploadErrors) > 0 {
		return fmt.Errorf("Failed to upload all files: %v", uploadErrors)
	}

	return nil
}

func (form *RecordUpsert) processFilesToDelete() error {
	if len(form.filesToDelete) == 0 {
		return nil // nothing to delete
	}

	if !form.record.HasId() {
		return errors.New("The record is not persisted yet.")
	}

	fs, err := form.app.NewFilesystem()
	if err != nil {
		return err
	}
	defer fs.Close()

	var deleteErrors []error
	for i := len(form.filesToDelete) - 1; i >= 0; i-- {
		filename := form.filesToDelete[i]
		path := form.record.BaseFilesPath() + "/" + filename

		if err := fs.Delete(path); err == nil {
			// remove the deleted file from the list
			form.filesToDelete = append(form.filesToDelete[:i], form.filesToDelete[i+1:]...)
		} else {
			// store the delete error
			deleteErrors = append(deleteErrors, fmt.Errorf("File %d: %v", i, err))
		}

		// try to delete the related file thumbs (if any)
		fs.DeletePrefix(form.record.BaseFilesPath() + "/thumbs_" + filename + "/")
	}

	if len(deleteErrors) > 0 {
		return fmt.Errorf("Failed to delete all files: %v", deleteErrors)
	}

	return nil
}
