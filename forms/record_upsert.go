package forms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// username value regex pattern
var usernameRegex = regexp.MustCompile(`^[\w][\w\.]*$`)

// RecordUpsert is a [models.Record] upsert (create/update) form.
type RecordUpsert struct {
	app          core.App
	dao          *daos.Dao
	manageAccess bool
	record       *models.Record

	filesToUpload map[string][]*rest.UploadedFile
	filesToDelete []string // names list

	// base model fields
	Id string `json:"id"`

	// auth collection fields
	// ---
	Username        string `json:"username"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	OldPassword     string `json:"oldPassword"`
	// ---

	Data map[string]any `json:"data"`
}

// NewRecordUpsert creates a new [RecordUpsert] form with initializer
// config created from the provided [core.App] and [models.Record] instances
// (for create you could pass a pointer to an empty Record - models.NewRecord(collection)).
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordUpsert(app core.App, record *models.Record) *RecordUpsert {
	form := &RecordUpsert{
		app:           app,
		dao:           app.Dao(),
		record:        record,
		filesToDelete: []string{},
		filesToUpload: map[string][]*rest.UploadedFile{},
	}

	form.loadFormDefaults()

	return form
}

// SetFullManageAccess sets the manageAccess bool flag of the current
// form to enable/disable directly changing some system record fields
// (often used with auth collection records).
func (form *RecordUpsert) SetFullManageAccess(fullManageAccess bool) {
	form.manageAccess = fullManageAccess
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordUpsert) SetDao(dao *daos.Dao) {
	form.dao = dao
}

func (form *RecordUpsert) loadFormDefaults() {
	form.Id = form.record.Id

	if form.record.Collection().IsAuth() {
		form.Username = form.record.Username()
		form.Email = form.record.Email()
		form.EmailVisibility = form.record.EmailVisibility()
		form.Verified = form.record.Verified()
	}

	form.Data = map[string]any{}
	for _, field := range form.record.Collection().Schema.Fields() {
		form.Data[field.Name] = form.record.Get(field.Name)
	}
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

func (form *RecordUpsert) extractRequestData(r *http.Request, keyPrefix string) (map[string]any, error) {
	switch form.getContentType(r) {
	case "application/json":
		return form.extractJsonData(r, keyPrefix)
	case "multipart/form-data":
		return form.extractMultipartFormData(r, keyPrefix)
	default:
		return nil, errors.New("Unsupported request Content-Type.")
	}
}

func (form *RecordUpsert) extractJsonData(r *http.Request, keyPrefix string) (map[string]any, error) {
	result := map[string]any{}

	err := rest.CopyJsonBody(r, &result)

	if keyPrefix != "" {
		parts := strings.Split(keyPrefix, ".")
		for _, part := range parts {
			if result[part] == nil {
				break
			}
			if v, ok := result[part].(map[string]any); ok {
				result = v
			}
		}
	}

	return result, err
}

func (form *RecordUpsert) extractMultipartFormData(r *http.Request, keyPrefix string) (map[string]any, error) {
	result := map[string]any{}

	// parse form data (if not already)
	if err := r.ParseMultipartForm(rest.DefaultMaxMemory); err != nil {
		return result, err
	}

	arrayValueSupportTypes := schema.ArraybleFieldTypes()

	form.filesToUpload = map[string][]*rest.UploadedFile{}

	for fullKey, values := range r.PostForm {
		key := fullKey
		if keyPrefix != "" {
			key = strings.TrimPrefix(key, keyPrefix+".")
		}

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

	// load uploaded files (if any)
	for _, field := range form.record.Collection().Schema.Fields() {
		if field.Type != schema.FieldTypeFile {
			continue // not a file field
		}

		key := field.Name
		fullKey := key
		if keyPrefix != "" {
			fullKey = keyPrefix + "." + key
		}

		files, err := rest.FindUploadedFiles(r, fullKey)
		if err != nil || len(files) == 0 {
			if err != nil && err != http.ErrMissingFile && form.app.IsDebug() {
				log.Printf("%q uploaded file error: %v\n", fullKey, err)
			}

			// skip invalid or missing file(s)
			continue
		}

		options, ok := field.Options.(*schema.FileOptions)
		if !ok {
			continue
		}

		if form.filesToUpload[key] == nil {
			form.filesToUpload[key] = []*rest.UploadedFile{}
		}

		if options.MaxSelect == 1 {
			form.filesToUpload[key] = append(form.filesToUpload[key], files[0])
		} else if options.MaxSelect > 1 {
			form.filesToUpload[key] = append(form.filesToUpload[key], files...)
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

// LoadRequest extracts the json or multipart/form-data request data
// and lods it into the form.
//
// File upload is supported only via multipart/form-data.
//
// To DELETE previously uploaded file(s) you can suffix the field name
// with the file index or filename (eg. `myfile.0`) and set it to null or empty string.
// For single file upload fields, you can skip the index and directly
// reset the field using its field name (eg. `myfile = null`).
func (form *RecordUpsert) LoadRequest(r *http.Request, keyPrefix string) error {
	requestData, err := form.extractRequestData(r, keyPrefix)
	if err != nil {
		return err
	}

	return form.LoadData(requestData)
}

// LoadData loads and normalizes the provided data into the form.
//
// To DELETE previously uploaded file(s) you can suffix the field name
// with the file index or filename (eg. `myfile.0`) and set it to null or empty string.
// For single file upload fields, you can skip the index and directly
// reset the field using its field name (eg. `myfile = null`).
func (form *RecordUpsert) LoadData(requestData map[string]any) error {
	// load base system fields
	if v, ok := requestData["id"]; ok {
		form.Id = cast.ToString(v)
	}

	// load auth system fields
	if form.record.Collection().IsAuth() {
		if v, ok := requestData["username"]; ok {
			form.Username = cast.ToString(v)
		}
		if v, ok := requestData["email"]; ok {
			form.Email = cast.ToString(v)
		}
		if v, ok := requestData["emailVisibility"]; ok {
			form.EmailVisibility = cast.ToBool(v)
		}
		if v, ok := requestData["verified"]; ok {
			form.Verified = cast.ToBool(v)
		}
		if v, ok := requestData["password"]; ok {
			form.Password = cast.ToString(v)
		}
		if v, ok := requestData["passwordConfirm"]; ok {
			form.PasswordConfirm = cast.ToString(v)
		}
		if v, ok := requestData["oldPassword"]; ok {
			form.OldPassword = cast.ToString(v)
		}
	}

	// extend the record schema data with the request data
	extendedData := form.record.SchemaData()
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

		if len(form.filesToUpload[key]) == 0 {
			continue
		}

		// refresh oldNames list
		oldNames = list.ToUniqueStringSlice(form.Data[key])

		if options.MaxSelect == 1 {
			// delete previous file(s) before replacing
			if len(oldNames) > 0 {
				form.filesToDelete = list.ToUniqueStringSlice(append(form.filesToDelete, oldNames...))
			}
			form.Data[key] = form.filesToUpload[key][0].Name()
		} else if options.MaxSelect > 1 {
			// append the id of each uploaded file instance
			for _, file := range form.filesToUpload[key] {
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
	baseFieldsRules := []*validation.FieldRules{
		validation.Field(
			&form.Id,
			validation.When(
				form.record.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
			).Else(validation.In(form.record.Id)),
		),
	}

	// auth fields validators
	if form.record.Collection().IsAuth() {
		baseFieldsRules = append(baseFieldsRules,
			validation.Field(
				&form.Username,
				// require only on update, because on create we fallback to auto generated username
				validation.When(!form.record.IsNew(), validation.Required),
				validation.Length(3, 100),
				validation.Match(usernameRegex),
				validation.By(form.checkUniqueUsername),
			),
			validation.Field(
				&form.Email,
				validation.When(
					form.record.Collection().AuthOptions().RequireEmail,
					validation.Required,
				),
				// don't allow direct email change (or unset) if the form doesn't have manage access permissions
				// (aka. allow only admin or authorized auth models to directly update the field)
				validation.When(
					!form.record.IsNew() && !form.manageAccess,
					validation.In(form.record.Email()),
				),
				validation.Length(1, 255),
				is.EmailFormat,
				validation.By(form.checkEmailDomain),
				validation.By(form.checkUniqueEmail),
			),
			validation.Field(
				&form.Verified,
				// don't allow changing verified if the form doesn't have manage access permissions
				// (aka. allow only admin or authorized auth models to directly change the field)
				validation.When(
					!form.manageAccess,
					validation.In(form.record.Verified()),
				),
			),
			validation.Field(
				&form.Password,
				validation.When(form.record.IsNew(), validation.Required),
				validation.Length(form.record.Collection().AuthOptions().MinPasswordLength, 72),
			),
			validation.Field(
				&form.PasswordConfirm,
				validation.When(
					(form.record.IsNew() || form.Password != ""),
					validation.Required,
				),
				validation.By(validators.Compare(form.Password)),
			),
			validation.Field(
				&form.OldPassword,
				// require old password only on update when:
				// - form.manageAccess is not set
				// - changing the existing password
				validation.When(
					!form.record.IsNew() && !form.manageAccess && form.Password != "",
					validation.Required,
					validation.By(form.checkOldPassword),
				),
			),
		)
	}

	if err := validation.ValidateStruct(form, baseFieldsRules...); err != nil {
		return err
	}

	// record data validator
	return validators.NewRecordDataValidator(
		form.dao,
		form.record,
		form.filesToUpload,
	).Validate(form.Data)
}

func (form *RecordUpsert) checkUniqueUsername(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	isUnique := form.dao.IsRecordValueUnique(
		form.record.Collection().Id,
		schema.FieldNameUsername,
		v,
		form.record.Id,
	)
	if !isUnique {
		return validation.NewError("validation_invalid_username", "The username is invalid or already in use.")
	}

	return nil
}

func (form *RecordUpsert) checkUniqueEmail(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	isUnique := form.dao.IsRecordValueUnique(
		form.record.Collection().Id,
		schema.FieldNameEmail,
		v,
		form.record.Id,
	)
	if !isUnique {
		return validation.NewError("validation_invalid_email", "The email is invalid or already in use.")
	}

	return nil
}

func (form *RecordUpsert) checkEmailDomain(value any) error {
	val, _ := value.(string)
	if val == "" {
		return nil // nothing to check
	}

	domain := val[strings.LastIndex(val, "@")+1:]
	only := form.record.Collection().AuthOptions().OnlyEmailDomains
	except := form.record.Collection().AuthOptions().ExceptEmailDomains

	// only domains check
	if len(only) > 0 && !list.ExistInSlice(domain, only) {
		return validation.NewError("validation_email_domain_not_allowed", "Email domain is not allowed.")
	}

	// except domains check
	if len(except) > 0 && list.ExistInSlice(domain, except) {
		return validation.NewError("validation_email_domain_not_allowed", "Email domain is not allowed.")
	}

	return nil
}

func (form *RecordUpsert) checkOldPassword(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	if !form.record.ValidatePassword(v) {
		return validation.NewError("validation_invalid_old_password", "Missing or invalid old password.")
	}

	return nil
}

func (form *RecordUpsert) ValidateAndFill() error {
	if err := form.Validate(); err != nil {
		return err
	}

	isNew := form.record.IsNew()

	// custom insertion id can be set only on create
	if isNew && form.Id != "" {
		form.record.SetId(form.Id)
		form.record.MarkAsNew()
	}

	// set auth fields
	if form.record.Collection().IsAuth() {
		// generate a default username during create (if missing)
		if form.record.IsNew() && form.Username == "" {
			baseUsername := form.record.Collection().Name + security.RandomStringWithAlphabet(5, "123456789")
			form.Username = form.dao.SuggestUniqueAuthRecordUsername(form.record.Collection().Id, baseUsername)
		}

		if form.Username != "" {
			if err := form.record.SetUsername(form.Username); err != nil {
				return err
			}
		}

		if isNew || form.manageAccess {
			if err := form.record.SetEmail(form.Email); err != nil {
				return err
			}
		}

		if err := form.record.SetEmailVisibility(form.EmailVisibility); err != nil {
			return err
		}

		if form.manageAccess {
			if err := form.record.SetVerified(form.Verified); err != nil {
				return err
			}
		}

		if form.Password != "" {
			if err := form.record.SetPassword(form.Password); err != nil {
				return err
			}
		}
	}

	// bulk load the remaining form data
	form.record.Load(form.Data)

	return nil
}

// DrySubmit performs a form submit within a transaction and reverts it.
// For actual record persistence, check the `form.Submit()` method.
//
// This method doesn't handle file uploads/deletes or trigger any app events!
func (form *RecordUpsert) DrySubmit(callback func(txDao *daos.Dao) error) error {
	isNew := form.record.IsNew()

	if err := form.ValidateAndFill(); err != nil {
		return err
	}

	// use the default app.Dao to prevent changing the transaction form.Dao
	// and causing "transaction has already been committed or rolled back" error
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
	if err := form.ValidateAndFill(); err != nil {
		return err
	}

	return runInterceptors(func() error {
		if !form.record.HasId() {
			form.record.RefreshId()
			form.record.MarkAsNew()
		}

		// upload new files (if any)
		if err := form.processFilesToUpload(); err != nil {
			return fmt.Errorf("failed to process the uploaded files: %w", err)
		}

		// persist the record model
		if saveErr := form.dao.SaveRecord(form.record); saveErr != nil {
			// try to cleanup the successfully uploaded files
			if _, err := form.deleteFilesByNamesList(form.getFilesToUploadNames()); err != nil && form.app.IsDebug() {
				log.Println(err)
			}

			return fmt.Errorf("failed to save the record: %w", saveErr)
		}

		// delete old files (if any)
		//
		// for now fail silently to avoid reupload when `form.Submit()`
		// is called manually (aka. not from an api request)...
		if err := form.processFilesToDelete(); err != nil && form.app.IsDebug() {
			log.Println(err)
		}

		return nil
	}, interceptors...)
}

func (form *RecordUpsert) getFilesToUploadNames() []string {
	names := []string{}

	for fieldKey := range form.filesToUpload {
		for _, file := range form.filesToUpload[fieldKey] {
			names = append(names, file.Name())
		}
	}

	return names
}

func (form *RecordUpsert) processFilesToUpload() error {
	if len(form.filesToUpload) == 0 {
		return nil // no parsed file fields
	}

	if !form.record.HasId() {
		return errors.New("the record is not persisted yet")
	}

	fs, err := form.app.NewFilesystem()
	if err != nil {
		return err
	}
	defer fs.Close()

	var uploadErrors []error // list of upload errors
	var uploaded []string    // list of uploaded file paths

	for fieldKey := range form.filesToUpload {
		for i, file := range form.filesToUpload[fieldKey] {
			path := form.record.BaseFilesPath() + "/" + file.Name()
			if err := fs.UploadMultipart(file.Header(), path); err == nil {
				// keep track of the already uploaded file
				uploaded = append(uploaded, path)
			} else {
				// store the upload error
				uploadErrors = append(uploadErrors, fmt.Errorf("file %d: %v", i, err))
			}
		}
	}

	if len(uploadErrors) > 0 {
		// cleanup - try to delete the successfully uploaded files (if any)
		form.deleteFilesByNamesList(uploaded)

		return fmt.Errorf("failed to upload all files: %v", uploadErrors)
	}

	return nil
}

func (form *RecordUpsert) processFilesToDelete() (err error) {
	form.filesToDelete, err = form.deleteFilesByNamesList(form.filesToDelete)
	return
}

// deleteFiles deletes a list of record files by their names.
// Returns the failed/remaining files.
func (form *RecordUpsert) deleteFilesByNamesList(filenames []string) ([]string, error) {
	if len(filenames) == 0 {
		return filenames, nil // nothing to delete
	}

	if !form.record.HasId() {
		return filenames, errors.New("the record doesn't have a unique ID")
	}

	fs, err := form.app.NewFilesystem()
	if err != nil {
		return filenames, err
	}
	defer fs.Close()

	var deleteErrors []error

	for i := len(filenames) - 1; i >= 0; i-- {
		filename := filenames[i]
		path := form.record.BaseFilesPath() + "/" + filename

		if err := fs.Delete(path); err == nil {
			// remove the deleted file from the list
			filenames = append(filenames[:i], filenames[i+1:]...)

			// try to delete the related file thumbs (if any)
			fs.DeletePrefix(form.record.BaseFilesPath() + "/thumbs_" + filename + "/")
		} else {
			// store the delete error
			deleteErrors = append(deleteErrors, fmt.Errorf("file %d: %v", i, err))
		}
	}

	if len(deleteErrors) > 0 {
		return filenames, fmt.Errorf("failed to delete all files: %v", deleteErrors)
	}

	return filenames, nil
}
