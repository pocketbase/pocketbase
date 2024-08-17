package forms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// username value regex pattern
var usernameRegex = regexp.MustCompile(`^\w[\w.\-]*$`)

// ---------------------------
// Types
// ---------------------------

// RecordUpsert is a [models.Record] upsert (create/update) form.
type RecordUpsert struct {
	app          core.App
	dao          *daos.Dao
	manageAccess bool
	record       *models.Record

	filesToUpload map[string][]*filesystem.File
	filesToDelete []string

	Id              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	OldPassword     string `json:"oldPassword"`

	data map[string]any
}

// ---------------------------
// Exported API
// ---------------------------

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
		filesToUpload: map[string][]*filesystem.File{},
	}
	form.loadFormDefaults()
	return form
}

// Data returns the loaded form's data.
func (form *RecordUpsert) Data() map[string]any {
	return form.data
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

// LoadRequest extracts the json or multipart/form-data request data
// and loads it into the form.
//
// File upload is supported only via multipart/form-data.
func (form *RecordUpsert) LoadRequest(r *http.Request, keyPrefix string) error {
	requestData, uploadedFiles, err := form.extractRequestData(r, keyPrefix)
	if err != nil {
		return err
	}

	if err := form.LoadData(requestData); err != nil {
		return err
	}

	for key, files := range uploadedFiles {
		err := form.AddFiles(key, files...)
		if err != nil {
			return err
		}
	}

	return nil
}

// FilesToUpload returns the parsed request files ready for upload.
func (form *RecordUpsert) FilesToUpload() map[string][]*filesystem.File {
	return form.filesToUpload
}

// FilesToDelete returns the parsed request filenames ready to be deleted.
func (form *RecordUpsert) FilesToDelete() []string {
	return form.filesToDelete
}

// AddFiles adds the provided file(s) to the specified file field.
//
// If the file field is a SINGLE-value file field (aka. "Max Select = 1"),
// then the newly added file will REPLACE the existing one.
// In this case if you pass more than 1 files only the first one will be assigned.
//
// If the file field is a MULTI-value file field (aka. "Max Select > 1"),
// then the newly added file(s) will be APPENDED to the existing one(s).
//
// Example
//
//	f1, _ := filesystem.NewFileFromPath("/path/to/file1.txt")
//	f2, _ := filesystem.NewFileFromPath("/path/to/file2.txt")
//	form.AddFiles("documents", f1, f2)
func (form *RecordUpsert) AddFiles(key string, files ...*filesystem.File) error {
	field := form.record.Collection().Schema.GetFieldByName(key)
	if field == nil || field.Type != schema.FieldTypeFile {
		return errors.New("invalid field key")
	}

	options, ok := field.Options.(*schema.FileOptions)
	if !ok {
		return errors.New("failed to initialize field options")
	}

	if len(files) == 0 {
		return nil // nothing to upload
	}

	if form.filesToUpload == nil {
		form.filesToUpload = map[string][]*filesystem.File{}
	}

	oldNames := list.ToUniqueStringSlice(form.data[key])

	if options.MaxSelect == 1 {
		// mark previous file(s) for deletion before replacing
		if len(oldNames) > 0 {
			form.filesToDelete = list.ToUniqueStringSlice(append(form.filesToDelete, oldNames...))
		}

		// replace
		form.filesToUpload[key] = []*filesystem.File{files[0]}
		form.data[key] = field.PrepareValue(files[0].Name)
	} else {
		// append
		form.filesToUpload[key] = append(form.filesToUpload[key], files...)
		for _, f := range files {
			oldNames = append(oldNames, f.Name)
		}
		form.data[key] = field.PrepareValue(oldNames)
	}

	return nil
}

// RemoveFiles removes a single or multiple file from the specified file field.
//
// NB! If filesToDelete is not set it will remove all existing files
// assigned to the file field (including those assigned with AddFiles)!
//
// Example
//
//	// mark only 2 files for removal
//	form.RemoveFiles("documents", "file1_aw4bdrvws6.txt", "file2_xwbs36bafv.txt")
//
//	// mark all "documents" files for removal
//	form.RemoveFiles("documents")
func (form *RecordUpsert) RemoveFiles(key string, toDelete ...string) error {
	field := form.record.Collection().Schema.GetFieldByName(key)
	if field == nil || field.Type != schema.FieldTypeFile {
		return errors.New("invalid field key")
	}

	existing := list.ToUniqueStringSlice(form.data[key])

	// mark all files for deletion
	if len(toDelete) == 0 {
		toDelete = make([]string, len(existing))
		copy(toDelete, existing)
	}

	form.removeFilesFromLists(existing, key, toDelete)

	return nil
}

// LoadData loads and normalizes the provided regular record data fields into the form.
func (form *RecordUpsert) LoadData(requestData map[string]any) error {
	form.loadBaseFields(requestData)

	// replace modifiers (if any)
	requestData = form.record.ReplaceModifers(requestData)

	extendedData, err := form.mergeRequestData(requestData)
	if err != nil {
		return err
	}

	for _, field := range form.record.Collection().Schema.Fields() {
		key := field.Name
		value := field.PrepareValue(extendedData[key])

		if field.Type != schema.FieldTypeFile {
			form.data[key] = value
			continue
		}

		// -----------------------------------------------------------
		// Delete previously uploaded file(s)
		// -----------------------------------------------------------
		form.handleFileField(key, value)
	}

	return nil
}

func (form *RecordUpsert) Submit(interceptors ...InterceptorFunc[*models.Record]) error {
	if err := form.ValidateAndFill(); err != nil {
		return err
	}

	return form.submitInternal(interceptors...)
}

func (form *RecordUpsert) SubmitWithoutValidation(interceptors ...InterceptorFunc[*models.Record]) error {
	if err := form.FillWithoutValidation(); err != nil {
		return err
	}

	return form.submitInternal(interceptors...)
}

// ValidateAndFill -> Validates the form and fills its record fields.
func (form *RecordUpsert) ValidateAndFill() error {
	return form.validateAndFill(true)
}

// FillWithoutValidation -> Fills the form's record fields without performing any validation.
func (form *RecordUpsert) FillWithoutValidation() error {
	return form.validateAndFill(false)
}

// DrySubmit performs a form submit within a transaction and reverts it.
// For actual record persistence, check the form.Submit() method.
//
// This method doesn't handle file uploads/deletes or trigger any app events!
func (form *RecordUpsert) DrySubmit(callback func(txDao *daos.Dao) error) error {
	isNew := form.record.IsNew()

	if err := form.ValidateAndFill(); err != nil {
		return err
	}

	var dryDao *daos.Dao
	if form.dao.ConcurrentDB() == form.dao.NonconcurrentDB() {
		// it is already in a transaction and therefore use the app concurrent db pool
		// to prevent "transaction has already been committed or rolled back" error
		dryDao = daos.New(form.app.Dao().ConcurrentDB())
	} else {
		// otherwise use the form non-concurrent dao db pool
		dryDao = daos.New(form.dao.NonconcurrentDB())
	}

	return dryDao.RunInTransaction(func(txDao *daos.Dao) error {
		tx, ok := txDao.DB().(*dbx.Tx)
		if !ok {
			return errors.New("failed to get transaction db")
		}
		defer func(tx *dbx.Tx) {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}(tx)

		if err := txDao.SaveRecord(form.record); err != nil {
			return form.prepareError(err)
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

// ---------------------------
// Internal Helper Functions
// ---------------------------

func (form *RecordUpsert) getContentType(r *http.Request) string {
	t := r.Header.Get("Content-Type")
	for i, c := range t {
		if c == ' ' || c == ';' {
			return t[:i]
		}
	}
	return t
}

func (form *RecordUpsert) loadFormDefaults() {
	form.Id = form.record.Id

	if form.record.Collection().IsAuth() {
		form.Username = form.record.Username()
		form.Email = form.record.Email()
		form.EmailVisibility = form.record.EmailVisibility()
		form.Verified = form.record.Verified()
	}

	form.data = map[string]any{}
	for _, field := range form.record.Collection().Schema.Fields() {
		form.data[field.Name] = form.record.Get(field.Name)
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordUpsert) Validate() error {
	// base form fields validator
	if err := form.validateBaseFields(); err != nil {
		return err
	}

	return validators.NewRecordDataValidator(
		form.dao,
		form.record,
		form.filesToUpload,
	).Validate(form.data)
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

func (form *RecordUpsert) validateAndFill(validate bool) error {
	if validate {
		if err := form.Validate(); err != nil {
			return err
		}
	}

	return form.fillRecordFields()
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

func (form *RecordUpsert) extractRequestData(r *http.Request, keyPrefix string) (map[string]any, map[string][]*filesystem.File, error) {
	switch form.getContentType(r) {
	case "application/json":
		return form.extractJsonData(r, keyPrefix)
	case "multipart/form-data":
		return form.extractMultipartFormData(r, keyPrefix)
	default:
		return nil, nil, errors.New("unsupported request content-type")
	}
}

func (form *RecordUpsert) extractJsonData(r *http.Request, keyPrefix string) (map[string]any, map[string][]*filesystem.File, error) {
	data := map[string]any{}
	err := rest.CopyJsonBody(r, &data)
	if err != nil {
		return nil, nil, err
	}

	if keyPrefix != "" {
		parts := strings.Split(keyPrefix, ".")
		for _, part := range parts {
			if data[part] == nil {
				break
			}
			if v, ok := data[part].(map[string]any); ok {
				data = v
			}
		}
	}

	return data, nil, err
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

func (form *RecordUpsert) extractMultipartFormData(r *http.Request, keyPrefix string) (map[string]any, map[string][]*filesystem.File, error) {
	if err := r.ParseMultipartForm(rest.DefaultMaxMemory); err != nil {
		return nil, nil, err
	}

	data := map[string]any{}
	filesToUpload := map[string][]*filesystem.File{}
	arraybleFieldTypes := schema.ArraybleFieldTypes()

	for fullKey, values := range r.PostForm {
		key := trimKeyPrefix(fullKey, keyPrefix)
		handleMultipartField(form, data, key, values, arraybleFieldTypes)
	}

	form.loadUploadedFiles(r, keyPrefix, filesToUpload)

	return data, filesToUpload, nil
}

func trimKeyPrefix(key, prefix string) string {
	if prefix != "" {
		return strings.TrimPrefix(key, prefix+".")
	}
	return key
}

func handleMultipartField(form *RecordUpsert, data map[string]any, key string, values []string, arraybleFieldTypes []string) {
	if len(values) == 0 {
		data[key] = nil
		return
	}

	// special case for multipart json encoded fields
	if key == rest.MultipartJsonKey {
		for _, v := range values {
			if err := json.Unmarshal([]byte(v), &data); err != nil {
				form.app.Logger().Debug("Failed to decode @json value into the data map", "error", err, "value", v)
			}
		}
		return
	}

	field := form.record.Collection().Schema.GetFieldByName(key)
	if field != nil && list.ExistInSlice(field.Type, arraybleFieldTypes) {
		data[key] = values
	} else {
		data[key] = values[0]
	}
}

func (form *RecordUpsert) loadUploadedFiles(r *http.Request, keyPrefix string, filesToUpload map[string][]*filesystem.File) {
	// load uploaded files (if any)
	for _, field := range form.record.Collection().Schema.Fields() {
		if field.Type != schema.FieldTypeFile {
			continue
		}

		key := field.Name
		fullKey := key
		if keyPrefix != "" {
			fullKey = keyPrefix + "." + key
		}

		files, err := rest.FindUploadedFiles(r, fullKey)
		if err != nil || len(files) == 0 {
			if err != nil && !errors.Is(err, http.ErrMissingFile) {
				form.app.Logger().Debug("Uploaded file error", slog.String("key", fullKey), slog.String("error", err.Error()))
			}
			// skip invalid or missing file(s)
			continue
		}

		filesToUpload[key] = append(filesToUpload[key], files...)
	}
}

func (form *RecordUpsert) removeFilesFromLists(existing []string, key string, toDelete []string) {
	// check for existing files
	for i := len(existing) - 1; i >= 0; i-- {
		if list.ExistInSlice(existing[i], toDelete) {
			form.filesToDelete = append(form.filesToDelete, existing[i])
			existing = append(existing[:i], existing[i+1:]...)
		}
	}

	// check for newly uploaded files
	for i := len(form.filesToUpload[key]) - 1; i >= 0; i-- {
		f := form.filesToUpload[key][i]
		if list.ExistInSlice(f.Name, toDelete) {
			form.filesToUpload[key] = append(form.filesToUpload[key][:i], form.filesToUpload[key][i+1:]...)
		}
	}

	form.data[key] = existing
}

func (form *RecordUpsert) loadBaseFields(requestData map[string]any) {
	// load base system fields
	if v, ok := requestData[schema.FieldNameId]; ok {
		form.Id = cast.ToString(v)
	}
	// load auth system fields
	if form.record.Collection().IsAuth() {
		form.loadAuthFields(requestData)
	}
}

func (form *RecordUpsert) loadAuthFields(requestData map[string]any) {
	if v, ok := requestData[schema.FieldNameUsername]; ok {
		form.Username = cast.ToString(v)
	}
	if v, ok := requestData[schema.FieldNameEmail]; ok {
		form.Email = cast.ToString(v)
	}
	if v, ok := requestData[schema.FieldNameEmailVisibility]; ok {
		form.EmailVisibility = cast.ToBool(v)
	}
	if v, ok := requestData[schema.FieldNameVerified]; ok {
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

func (form *RecordUpsert) mergeRequestData(requestData map[string]any) (map[string]any, error) {
	// create a shallow copy of form.data
	var extendedData = make(map[string]any, len(form.data))
	for k, v := range form.data {
		extendedData[k] = v
	}
	// extend form.data with the request data
	rawData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawData, &extendedData); err != nil {
		return nil, err
	}

	return extendedData, nil
}

func (form *RecordUpsert) handleFileField(key string, value any) {
	oldNames := form.record.GetStringSlice(key)
	submittedNames := list.ToUniqueStringSlice(value)

	// if empty value was set, mark all previously uploaded files for deletion
	// otherwise check for "deleted" (aka. unsubmitted) file names
	if len(submittedNames) > len(oldNames) || len(list.SubtractSlice(submittedNames, oldNames)) != 0 {
		err := form.RemoveFiles(key)
		if err != nil {
			return
		}
		return
	}
	// allow file key reassignments for file names sorting
	// (only if all submitted values already exists)
	if len(submittedNames) > 0 && len(list.SubtractSlice(submittedNames, oldNames)) == 0 {
		form.data[key] = submittedNames
	}
}

func (form *RecordUpsert) validateBaseFields() error {
	baseFieldsRules := []*validation.FieldRules{
		validation.Field(
			&form.Id,
			validation.When(
				form.record.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
				validation.By(validators.UniqueId(form.dao, form.record.TableName())),
			).Else(validation.In(form.record.Id)),
		),
	}

	if form.record.Collection().IsAuth() {
		baseFieldsRules = append(baseFieldsRules, form.getAuthFieldRules()...)
	}

	return validation.ValidateStruct(form, baseFieldsRules...)
}

func (form *RecordUpsert) getAuthFieldRules() []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(
			&form.Username,
			validation.When(!form.record.IsNew(), validation.Required),
			validation.Length(3, 150),
			validation.Match(usernameRegex),
			validation.By(form.checkUniqueUsername),
		),
		validation.Field(
			&form.Email,
			validation.When(
				form.record.Collection().AuthOptions().RequireEmail,
				validation.Required,
			),
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
			validation.When(
				!form.manageAccess,
				validation.In(form.record.Verified()),
			),
		),
		validation.Field(
			&form.Password,
			validation.When(
				form.record.IsNew() || form.PasswordConfirm != "" || form.OldPassword != "",
				validation.Required,
			),
			validation.Length(form.record.Collection().AuthOptions().MinPasswordLength, 72),
		),
		validation.Field(
			&form.PasswordConfirm,
			validation.When(
				form.record.IsNew() || form.Password != "" || form.OldPassword != "",
				validation.Required,
			),
			validation.By(validators.Compare(form.Password)),
		),
		validation.Field(
			&form.OldPassword,
			validation.When(
				!form.record.IsNew() && !form.manageAccess && (form.Password != "" || form.PasswordConfirm != ""),
				validation.Required,
				validation.By(form.checkOldPassword),
			),
		),
	}
}

func (form *RecordUpsert) fillRecordFields() error {
	isNew := form.record.IsNew()

	if isNew && form.Id != "" {
		form.record.SetId(form.Id)
		form.record.MarkAsNew()
	}

	if form.record.Collection().IsAuth() {
		err := form.fillAuthFields(isNew)
		if err != nil {
			return err
		}
	}

	form.record.Load(form.data)

	return nil
}

func (form *RecordUpsert) fillAuthFields(isNew bool) error {
	if isNew && form.Username == "" {
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

	if form.Password != "" && form.Password == form.PasswordConfirm {
		if err := form.record.SetPassword(form.Password); err != nil {
			return err
		}
	}

	return nil
}

func (form *RecordUpsert) submitInternal(interceptors ...InterceptorFunc[*models.Record]) error {
	return runInterceptors(form.record, func(record *models.Record) error {
		form.record = record

		if !form.record.HasId() {
			form.record.RefreshId()
			form.record.MarkAsNew()
		}

		dao := form.dao.Clone()

		dao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
			return form.beforeSave(eventDao, m, action, form.dao.BeforeCreateFunc)
		}

		dao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
			return form.beforeSave(eventDao, m, action, form.dao.BeforeUpdateFunc)
		}

		if err := dao.SaveRecord(form.record); err != nil {
			return form.prepareError(err)
		}

		if err := form.processFilesToDelete(); err != nil {
			form.app.Logger().Debug("Failed to delete old files", slog.String("error", err.Error()))
		}

		return nil
	}, interceptors...)
}

func (form *RecordUpsert) beforeSave(eventDao *daos.Dao, m models.Model, action func() error, originalFunc func(eventDao *daos.Dao, m models.Model, action func() error) error) error {
	newAction := func() error {
		if m.TableName() == form.record.TableName() && m.GetId() == form.record.GetId() {
			if err := form.processFilesToUpload(); err != nil {
				return err
			}
		}

		return action()
	}

	if originalFunc != nil {
		return originalFunc(eventDao, m, newAction)
	}

	return newAction()
}

func (form *RecordUpsert) processFilesToUpload() error {
	if len(form.filesToUpload) == 0 {
		return nil
	}

	if !form.record.HasId() {
		return errors.New("the record doesn't have an id")
	}

	fs, err := form.app.NewFilesystem()
	if err != nil {
		return err
	}
	defer func(fs *filesystem.System) {
		err := fs.Close()
		if err != nil {
			return
		}
	}(fs)

	var uploadErrors []error
	var uploaded []string

	for _, files := range form.filesToUpload {
		for i, file := range files {
			path := form.record.BaseFilesPath() + "/" + file.Name
			if err := fs.UploadFile(file, path); err == nil {
				uploaded = append(uploaded, path)
			} else {
				uploadErrors = append(uploadErrors, fmt.Errorf("file %d: %v", i, err))
			}
		}
	}

	if len(uploadErrors) > 0 {
		_, err := form.deleteFilesByNamesList(uploaded)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to upload all files: %v", uploadErrors)
	}

	return nil
}

func (form *RecordUpsert) processFilesToDelete() error {
	files, err := form.deleteFilesByNamesList(form.filesToDelete)
	form.filesToDelete = files
	return err
}

func (form *RecordUpsert) deleteFilesByNamesList(filenames []string) ([]string, error) {
	if len(filenames) == 0 {
		return filenames, nil
	}

	if !form.record.HasId() {
		return filenames, errors.New("the record doesn't have an id")
	}

	fs, err := form.app.NewFilesystem()
	if err != nil {
		return filenames, err
	}
	defer func(fs *filesystem.System) {
		err := fs.Close()
		if err != nil {
			return
		}
	}(fs)

	var deleteErrors []error

	for i := len(filenames) - 1; i >= 0; i-- {
		filename := filenames[i]
		path := form.record.BaseFilesPath() + "/" + filename

		if err := fs.Delete(path); err == nil {
			filenames = append(filenames[:i], filenames[i+1:]...)
			fs.DeletePrefix(form.record.BaseFilesPath() + "/thumbs_" + filename + "/")
		} else {
			deleteErrors = append(deleteErrors, fmt.Errorf("file %d: %v", i, err))
		}
	}

	if len(deleteErrors) > 0 {
		return filenames, fmt.Errorf("failed to delete all files: %v", deleteErrors)
	}

	return filenames, nil
}

func (form *RecordUpsert) prepareError(err error) error {
	msg := strings.ToLower(err.Error())

	validationErrs := validation.Errors{}

	if strings.Contains(msg, "unique constraint failed") {
		msg = strings.ReplaceAll(strings.TrimSpace(msg), ",", " ")

		c := form.record.Collection()
		for _, f := range c.Schema.Fields() {
			if strings.Contains(msg+" ", strings.ToLower(c.Name+"."+f.Name)) {
				validationErrs[f.Name] = validation.NewError("validation_not_unique", "Value must be unique")
			}
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return err
}
