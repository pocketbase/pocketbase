package forms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

// RecordUpsert specifies a [models.Record] upsert (create/update) form.
type RecordUpsert struct {
	config RecordUpsertConfig
	record *models.Record

	filesToDelete []string // names list
	filesToUpload []*rest.UploadedFile

	Id   string         `form:"id" json:"id"`
	Data map[string]any `json:"data"`
}

// RecordUpsertConfig is the [RecordUpsert] factory initializer config.
//
// NB! App is required struct member.
type RecordUpsertConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewRecordUpsert creates a new [RecordUpsert] form with initializer
// config created from the provided [core.App] and [models.Record] instances
// (for create you could pass a pointer to an empty Record - `models.NewRecord(collection)`).
//
// If you want to submit the form as part of another transaction, use
// [NewRecordUpsertWithConfig] with explicitly set Dao.
func NewRecordUpsert(app core.App, record *models.Record) *RecordUpsert {
	return NewRecordUpsertWithConfig(RecordUpsertConfig{
		App: app,
	}, record)
}

// NewRecordUpsertWithConfig creates a new [RecordUpsert] form
// with the provided config and [models.Record] instance or panics on invalid configuration
// (for create you could pass a pointer to an empty Record - `models.NewRecord(collection)`).
func NewRecordUpsertWithConfig(config RecordUpsertConfig, record *models.Record) *RecordUpsert {
	form := &RecordUpsert{
		config:        config,
		record:        record,
		filesToDelete: []string{},
		filesToUpload: []*rest.UploadedFile{},
	}

	if form.config.App == nil || form.record == nil {
		panic("Invalid initializer config or nil upsert model.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	form.Id = record.Id

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

	if id, ok := requestData["id"]; ok {
		form.Id = cast.ToString(id)
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

		if field.Type != schema.FieldTypeFile {
			form.Data[key] = value
			continue
		}

		options, _ := field.Options.(*schema.FileOptions)
		oldNames := list.ToUniqueStringSlice(form.Data[key])

		// -----------------------------------------------------------
		// Delete previously uploaded file(s)
		// -----------------------------------------------------------

		// if empty value was set, mark all previously uploaded files for deletion
		if len(list.ToUniqueStringSlice(value)) == 0 && len(oldNames) > 0 {
			form.filesToDelete = append(form.filesToDelete, oldNames...)
			form.Data[key] = []string{}
		} else if len(oldNames) > 0 {
			indexesToDelete := make([]int, 0, len(extendedData))

			// search for individual file name to delete (eg. "file.test.png = null")
			for i, name := range oldNames {
				if v, ok := extendedData[key+"."+name]; ok && cast.ToString(v) == "" {
					indexesToDelete = append(indexesToDelete, i)
				}
			}

			// search for individual file index to delete (eg. "file.0 = null")
			keyExp, _ := regexp.Compile(`^` + regexp.QuoteMeta(key) + `\.\d+$`)
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
			nonDeleted := make([]string, 0, len(oldNames))
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

		// -----------------------------------------------------------
		// Check for new uploaded file
		// -----------------------------------------------------------

		if form.getContentType(r) != "multipart/form-data" {
			continue // file upload is supported only via multipart/form-data
		}

		files, err := rest.FindUploadedFiles(r, key)
		if err != nil {
			if form.config.App.IsDebug() {
				log.Printf("%q uploaded file error: %v\n", key, err)
			}

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
	}

	return form.normalizeData()
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordUpsert) Validate() error {
	// base form fields validator
	baseFieldsErrors := validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.record.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
			).Else(validation.In(form.record.Id)),
		),
	)
	if baseFieldsErrors != nil {
		return baseFieldsErrors
	}

	// record data validator
	dataValidator := validators.NewRecordDataValidator(
		form.config.Dao,
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

	isNew := form.record.IsNew()

	// custom insertion id can be set only on create
	if isNew && form.Id != "" {
		form.record.MarkAsNew()
		form.record.SetId(form.Id)
	}

	// bulk load form data
	if err := form.record.Load(form.Data); err != nil {
		return err
	}

	return form.config.Dao.RunInTransaction(func(txDao *daos.Dao) error {
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

		// restore record isNew state
		if isNew {
			form.record.MarkAsNew()
		}

		if callback != nil {
			return callback(txDao)
		}

		return nil
	})
}

// Submit validates the form and upserts the form Record model.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *RecordUpsert) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}

	// custom insertion id can be set only on create
	if form.record.IsNew() && form.Id != "" {
		form.record.MarkAsNew()
		form.record.SetId(form.Id)
	}

	// bulk load form data
	if err := form.record.Load(form.Data); err != nil {
		return err
	}

	return runInterceptors(func() error {
		return form.config.Dao.RunInTransaction(func(txDao *daos.Dao) error {
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
	}, interceptors...)
}

func (form *RecordUpsert) processFilesToUpload() error {
	if len(form.filesToUpload) == 0 {
		return nil // nothing to upload
	}

	if !form.record.HasId() {
		return errors.New("The record is not persisted yet.")
	}

	fs, err := form.config.App.NewFilesystem()
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

	fs, err := form.config.App.NewFilesystem()
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
