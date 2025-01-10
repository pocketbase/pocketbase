package core

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeFile] = func() Field {
		return &FileField{}
	}
}

const FieldTypeFile = "file"

const DefaultFileFieldMaxSize int64 = 5 << 20

var looseFilenameRegex = regexp.MustCompile(`^[^\./\\][^/\\]+$`)

const (
	deletedFilesPrefix  = internalCustomFieldKeyPrefix + "_deletedFilesPrefix_"
	uploadedFilesPrefix = internalCustomFieldKeyPrefix + "_uploadedFilesPrefix_"
)

var (
	_ Field                 = (*FileField)(nil)
	_ MultiValuer           = (*FileField)(nil)
	_ DriverValuer          = (*FileField)(nil)
	_ GetterFinder          = (*FileField)(nil)
	_ SetterFinder          = (*FileField)(nil)
	_ RecordInterceptor     = (*FileField)(nil)
	_ MaxBodySizeCalculator = (*FileField)(nil)
)

// FileField defines "file" type field for managing record file(s).
//
// Only the file name is stored as part of the record value.
// New files (aka. files to upload) are expected to be of *filesytem.File.
//
// If MaxSelect is not set or <= 1, then the field value is expected to be a single record id.
//
// If MaxSelect is > 1, then the field value is expected to be a slice of record ids.
//
// The respective zero record field value is either empty string (single) or empty string slice (multiple).
//
// ---
//
// The following additional setter keys are available:
//
//   - "fieldName+" - append one or more files to the existing record one. For example:
//
//     // []string{"old1.txt", "old2.txt", "new1_ajkvass.txt", "new2_klhfnwd.txt"}
//     record.Set("documents+", []*filesystem.File{new1, new2})
//
//   - "+fieldName" - prepend one or more files to the existing record one. For example:
//
//     // []string{"new1_ajkvass.txt", "new2_klhfnwd.txt", "old1.txt", "old2.txt",}
//     record.Set("+documents", []*filesystem.File{new1, new2})
//
//   - "fieldName-" - subtract/delete one or more files from the existing record one. For example:
//
//     // []string{"old2.txt",}
//     record.Set("documents-", "old1.txt")
type FileField struct {
	// Name (required) is the unique name of the field.
	Name string `form:"name" json:"name"`

	// Id is the unique stable field identifier.
	//
	// It is automatically generated from the name when adding to a collection FieldsList.
	Id string `form:"id" json:"id"`

	// System prevents the renaming and removal of the field.
	System bool `form:"system" json:"system"`

	// Hidden hides the field from the API response.
	Hidden bool `form:"hidden" json:"hidden"`

	// Presentable hints the Dashboard UI to use the underlying
	// field record value in the relation preview label.
	Presentable bool `form:"presentable" json:"presentable"`

	// ---

	// MaxSize specifies the maximum size of a single uploaded file (in bytes and up to 2^53-1).
	//
	// If zero, a default limit of 5MB is applied.
	MaxSize int64 `form:"maxSize" json:"maxSize"`

	// MaxSelect specifies the max allowed files.
	//
	// For multiple files the value must be > 1, otherwise fallbacks to single (default).
	MaxSelect int `form:"maxSelect" json:"maxSelect"`

	// MimeTypes specifies an optional list of the allowed file mime types.
	//
	// Leave it empty to disable the validator.
	MimeTypes []string `form:"mimeTypes" json:"mimeTypes"`

	// Thumbs specifies an optional list of the supported thumbs for image based files.
	//
	// Each entry must be in one of the following formats:
	//
	//   - WxH  (eg. 100x300) - crop to WxH viewbox (from center)
	//   - WxHt (eg. 100x300t) - crop to WxH viewbox (from top)
	//   - WxHb (eg. 100x300b) - crop to WxH viewbox (from bottom)
	//   - WxHf (eg. 100x300f) - fit inside a WxH viewbox (without cropping)
	//   - 0xH  (eg. 0x300)    - resize to H height preserving the aspect ratio
	//   - Wx0  (eg. 100x0)    - resize to W width preserving the aspect ratio
	Thumbs []string `form:"thumbs" json:"thumbs"`

	// Protected will require the users to provide a special file token to access the file.
	//
	// Note that by default all files are publicly accessible.
	//
	// For the majority of the cases this is fine because by default
	// all file names have random part appended to their name which
	// need to be known by the user before accessing the file.
	Protected bool `form:"protected" json:"protected"`

	// Required will require the field value to have at least one file.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *FileField) Type() string {
	return FieldTypeFile
}

// GetId implements [Field.GetId] interface method.
func (f *FileField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *FileField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *FileField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *FileField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *FileField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *FileField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *FileField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *FileField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// IsMultiple implements MultiValuer interface and checks whether the
// current field options support multiple values.
func (f *FileField) IsMultiple() bool {
	return f.MaxSelect > 1
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *FileField) ColumnType(app App) string {
	if f.IsMultiple() {
		return "JSON DEFAULT '[]' NOT NULL"
	}

	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *FileField) PrepareValue(record *Record, raw any) (any, error) {
	return f.normalizeValue(raw), nil
}

// DriverValue implements the [DriverValuer] interface.
func (f *FileField) DriverValue(record *Record) (driver.Value, error) {
	files := f.toSliceValue(record.GetRaw(f.Name))

	if f.IsMultiple() {
		ja := make(types.JSONArray[string], len(files))
		for i, v := range files {
			ja[i] = f.getFileName(v)
		}
		return ja, nil
	}

	if len(files) == 0 {
		return "", nil
	}

	return f.getFileName(files[len(files)-1]), nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *FileField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.MaxSelect, validation.Min(0), validation.Max(maxSafeJSONInt)),
		validation.Field(&f.MaxSize, validation.Min(0), validation.Max(maxSafeJSONInt)),
		validation.Field(&f.Thumbs, validation.Each(
			validation.NotIn("0x0", "0x0t", "0x0b", "0x0f"),
			validation.Match(filesystem.ThumbSizeRegex),
		)),
	)
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *FileField) ValidateValue(ctx context.Context, app App, record *Record) error {
	files := f.toSliceValue(record.GetRaw(f.Name))
	if len(files) == 0 {
		if f.Required {
			return validation.ErrRequired
		}
		return nil // nothing to check
	}

	// validate existing and disallow new plain string filenames submission
	// (new files must be *filesystem.File)
	// ---
	oldExistingStrings := f.toSliceValue(f.getLatestOldValue(app, record))
	existingStrings := list.ToInterfaceSlice(f.extractPlainStrings(files))
	addedStrings := f.excludeFiles(existingStrings, oldExistingStrings)

	if len(addedStrings) > 0 {
		invalidFiles := make([]string, len(addedStrings))
		for i, invalid := range addedStrings {
			invalidStr := cast.ToString(invalid)
			if len(invalidStr) > 250 {
				invalidStr = invalidStr[:250]
			}
			invalidFiles[i] = invalidStr
		}

		return validation.NewError("validation_invalid_file", "Invalid new files: {{.invalidFiles}}.").
			SetParams(map[string]any{"invalidFiles": invalidFiles})
	}

	maxSelect := f.maxSelect()
	if len(files) > maxSelect {
		return validation.NewError("validation_too_many_files", "The maximum allowed files is {{.maxSelect}}").
			SetParams(map[string]any{"maxSelect": maxSelect})
	}

	// validate uploaded
	// ---
	uploads := f.extractUploadableFiles(files)
	for _, upload := range uploads {
		// loosely check the filename just in case it was manually changed after the normalization
		err := validation.Length(1, 150).Validate(upload.Name)
		if err != nil {
			return err
		}
		err = validation.Match(looseFilenameRegex).Validate(upload.Name)
		if err != nil {
			return err
		}

		// check size
		err = validators.UploadedFileSize(f.maxSize())(upload)
		if err != nil {
			return err
		}

		// check type
		if len(f.MimeTypes) > 0 {
			err = validators.UploadedFileMimeType(f.MimeTypes)(upload)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *FileField) maxSize() int64 {
	if f.MaxSize <= 0 {
		return DefaultFileFieldMaxSize
	}

	return f.MaxSize
}

func (f *FileField) maxSelect() int {
	if f.MaxSelect <= 1 {
		return 1
	}

	return f.MaxSelect
}

// CalculateMaxBodySize implements the [MaxBodySizeCalculator] interface.
func (f *FileField) CalculateMaxBodySize() int64 {
	return f.maxSize() * int64(f.maxSelect())
}

// Interceptors
// -------------------------------------------------------------------

// Intercept implements the [RecordInterceptor] interface.
//
// note: files delete after records deletion is handled globally by the app FileManager hook
func (f *FileField) Intercept(
	ctx context.Context,
	app App,
	record *Record,
	actionName string,
	actionFunc func() error,
) error {
	switch actionName {
	case InterceptorActionCreateExecute, InterceptorActionUpdateExecute:
		oldValue := f.getLatestOldValue(app, record)

		err := f.processFilesToUpload(ctx, app, record)
		if err != nil {
			return err
		}

		err = actionFunc()
		if err != nil {
			return errors.Join(err, f.afterRecordExecuteFailure(newContextIfInvalid(ctx), app, record))
		}

		f.rememberFilesToDelete(app, record, oldValue)

		f.afterRecordExecuteSuccess(newContextIfInvalid(ctx), app, record)

		return nil
	case InterceptorActionAfterCreateError, InterceptorActionAfterUpdateError:
		// when in transaction we assume that the error was handled by afterRecordExecuteFailure
		_, insideTransaction := app.DB().(*dbx.Tx)
		if insideTransaction {
			return actionFunc()
		}

		failedToDelete, deleteErr := f.deleteNewlyUploadedFiles(newContextIfInvalid(ctx), app, record)
		if deleteErr != nil {
			app.Logger().Warn(
				"Failed to cleanup all new files after record commit failure",
				"error", deleteErr,
				"failedToDelete", failedToDelete,
			)
		}

		record.SetRaw(deletedFilesPrefix+f.Name, nil)

		if record.IsNew() {
			// try to delete the record directory if there are no other files
			//
			// note: executed only on create failure to avoid accidentally
			// deleting a concurrently updating directory due to the
			// eventual consistent nature of some storage providers
			err := f.deleteEmptyRecordDir(newContextIfInvalid(ctx), app, record)
			if err != nil {
				app.Logger().Warn("Failed to delete empty dir after new record commit failure", "error", err)
			}
		}

		return actionFunc()
	case InterceptorActionAfterCreate, InterceptorActionAfterUpdate:
		record.SetRaw(uploadedFilesPrefix+f.Name, nil)

		err := f.processFilesToDelete(ctx, app, record)
		if err != nil {
			return err
		}

		return actionFunc()
	default:
		return actionFunc()
	}
}
func (f *FileField) getLatestOldValue(app App, record *Record) any {
	if !record.IsNew() {
		latestOriginal, err := app.FindRecordById(record.Collection(), cast.ToString(record.LastSavedPK()))
		if err == nil {
			return latestOriginal.GetRaw(f.Name)
		}
	}

	return record.Original().GetRaw(f.Name)
}

func (f *FileField) afterRecordExecuteSuccess(ctx context.Context, app App, record *Record) {
	uploaded, _ := record.GetRaw(uploadedFilesPrefix + f.Name).([]*filesystem.File)

	// replace the uploaded file objects with their plain string names
	newValue := f.toSliceValue(record.GetRaw(f.Name))
	for i, v := range newValue {
		if file, ok := v.(*filesystem.File); ok {
			uploaded = append(uploaded, file)
			newValue[i] = file.Name
		}
	}
	f.setValue(record, newValue)

	record.SetRaw(uploadedFilesPrefix+f.Name, uploaded)
}

func (f *FileField) afterRecordExecuteFailure(ctx context.Context, app App, record *Record) error {
	uploaded := f.extractUploadableFiles(f.toSliceValue(record.GetRaw(f.Name)))

	toDelete := make([]string, len(uploaded))
	for i, file := range uploaded {
		toDelete[i] = file.Name
	}

	// delete previously uploaded files
	failedToDelete, deleteErr := f.deleteFilesByNamesList(ctx, app, record, list.ToUniqueStringSlice(toDelete))

	if len(failedToDelete) > 0 {
		app.Logger().Warn(
			"Failed to cleanup the new uploaded file after record db write failure",
			"error", deleteErr,
			"failedToDelete", failedToDelete,
		)
	}

	return deleteErr
}

func (f *FileField) deleteEmptyRecordDir(ctx context.Context, app App, record *Record) error {
	fsys, err := app.NewFilesystem()
	if err != nil {
		return err
	}
	defer fsys.Close()
	fsys.SetContext(newContextIfInvalid(ctx))

	dir := record.BaseFilesPath()

	if !fsys.IsEmptyDir(dir) {
		return nil // no-op
	}

	err = fsys.Delete(dir)
	if err != nil && !errors.Is(err, filesystem.ErrNotFound) {
		return err
	}

	return nil
}

func (f *FileField) processFilesToDelete(ctx context.Context, app App, record *Record) error {
	markedForDelete, _ := record.GetRaw(deletedFilesPrefix + f.Name).([]string)
	if len(markedForDelete) == 0 {
		return nil
	}

	old := list.ToInterfaceSlice(markedForDelete)
	new := list.ToInterfaceSlice(f.extractPlainStrings(f.toSliceValue(record.GetRaw(f.Name))))
	diff := f.excludeFiles(old, new)

	toDelete := make([]string, len(diff))
	for i, del := range diff {
		toDelete[i] = f.getFileName(del)
	}

	failedToDelete, err := f.deleteFilesByNamesList(ctx, app, record, list.ToUniqueStringSlice(toDelete))

	record.SetRaw(deletedFilesPrefix+f.Name, failedToDelete)

	return err
}

func (f *FileField) rememberFilesToDelete(app App, record *Record, oldValue any) {
	old := list.ToInterfaceSlice(f.extractPlainStrings(f.toSliceValue(oldValue)))
	new := list.ToInterfaceSlice(f.extractPlainStrings(f.toSliceValue(record.GetRaw(f.Name))))
	diff := f.excludeFiles(old, new)

	toDelete, _ := record.GetRaw(deletedFilesPrefix + f.Name).([]string)

	for _, del := range diff {
		toDelete = append(toDelete, f.getFileName(del))
	}

	record.SetRaw(deletedFilesPrefix+f.Name, toDelete)
}

func (f *FileField) processFilesToUpload(ctx context.Context, app App, record *Record) error {
	uploads := f.extractUploadableFiles(f.toSliceValue(record.GetRaw(f.Name)))
	if len(uploads) == 0 {
		return nil
	}

	if record.Id == "" {
		return errors.New("uploading files requires the record to have a valid nonempty id")
	}

	fsys, err := app.NewFilesystem()
	if err != nil {
		return err
	}
	defer fsys.Close()
	fsys.SetContext(ctx)

	var failed []error     // list of upload errors
	var succeeded []string // list of uploaded file names

	for _, upload := range uploads {
		path := record.BaseFilesPath() + "/" + upload.Name
		if err := fsys.UploadFile(upload, path); err == nil {
			succeeded = append(succeeded, upload.Name)
		} else {
			failed = append(failed, fmt.Errorf("%q: %w", upload.Name, err))
			break // for now stop on the first error since we currently don't allow partial uploads
		}
	}

	if len(failed) > 0 {
		// cleanup - try to delete the successfully uploaded files (if any)
		_, cleanupErr := f.deleteFilesByNamesList(newContextIfInvalid(ctx), app, record, succeeded)

		failed = append(failed, cleanupErr)

		return fmt.Errorf("failed to upload all files: %w", errors.Join(failed...))
	}

	return nil
}

func (f *FileField) deleteNewlyUploadedFiles(ctx context.Context, app App, record *Record) ([]string, error) {
	uploaded, _ := record.GetRaw(uploadedFilesPrefix + f.Name).([]*filesystem.File)
	if len(uploaded) == 0 {
		return nil, nil
	}

	names := make([]string, len(uploaded))
	for i, file := range uploaded {
		names[i] = file.Name
	}

	failed, err := f.deleteFilesByNamesList(ctx, app, record, list.ToUniqueStringSlice(names))
	if err != nil {
		return failed, err
	}

	record.SetRaw(uploadedFilesPrefix+f.Name, nil)

	return nil, nil
}

// deleteFiles deletes a list of record files by their names.
// Returns the failed/remaining files.
func (f *FileField) deleteFilesByNamesList(ctx context.Context, app App, record *Record, filenames []string) ([]string, error) {
	if len(filenames) == 0 {
		return nil, nil // nothing to delete
	}

	if record.Id == "" {
		return filenames, errors.New("the record doesn't have an id")
	}

	fsys, err := app.NewFilesystem()
	if err != nil {
		return filenames, err
	}
	defer fsys.Close()
	fsys.SetContext(ctx)

	var failures []error

	for i := len(filenames) - 1; i >= 0; i-- {
		filename := filenames[i]
		if filename == "" || strings.ContainsAny(filename, "/\\") {
			continue // empty or not a plain filename
		}

		path := record.BaseFilesPath() + "/" + filename

		err := fsys.Delete(path)
		if err != nil && !errors.Is(err, filesystem.ErrNotFound) {
			// store the delete error
			failures = append(failures, fmt.Errorf("file %d (%q): %w", i, filename, err))
		} else {
			// remove the deleted file from the list
			filenames = append(filenames[:i], filenames[i+1:]...)

			// try to delete the related file thumbs (if any)
			thumbsErr := fsys.DeletePrefix(record.BaseFilesPath() + "/thumbs_" + filename + "/")
			if len(thumbsErr) > 0 {
				app.Logger().Warn("Failed to delete file thumbs", "error", errors.Join(thumbsErr...))
			}
		}
	}

	if len(failures) > 0 {
		return filenames, fmt.Errorf("failed to delete all files: %w", errors.Join(failures...))
	}

	return nil, nil
}

// newContextIfInvalid returns a new Background context if the provided one was cancelled.
func newContextIfInvalid(ctx context.Context) context.Context {
	if ctx.Err() == nil {
		return ctx
	}

	return context.Background()
}

// -------------------------------------------------------------------

// FindGetter implements the [GetterFinder] interface.
func (f *FileField) FindGetter(key string) GetterFunc {
	switch key {
	case f.Name:
		return func(record *Record) any {
			return record.GetRaw(f.Name)
		}
	case f.Name + ":unsaved":
		return func(record *Record) any {
			return f.extractUploadableFiles(f.toSliceValue(record.GetRaw(f.Name)))
		}
	case f.Name + ":uploaded":
		// deprecated
		log.Println("[file field getter] please replace :uploaded with :unsaved")
		return func(record *Record) any {
			return f.extractUploadableFiles(f.toSliceValue(record.GetRaw(f.Name)))
		}
	default:
		return nil
	}
}

// -------------------------------------------------------------------

// FindSetter implements the [SetterFinder] interface.
func (f *FileField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		return f.setValue
	case "+" + f.Name:
		return f.prependValue
	case f.Name + "+":
		return f.appendValue
	case f.Name + "-":
		return f.subtractValue
	default:
		return nil
	}
}

func (f *FileField) setValue(record *Record, raw any) {
	val := f.normalizeValue(raw)

	record.SetRaw(f.Name, val)
}

func (f *FileField) prependValue(record *Record, toPrepend any) {
	files := f.toSliceValue(record.GetRaw(f.Name))
	prepends := f.toSliceValue(toPrepend)

	if len(prepends) > 0 {
		files = append(prepends, files...)
	}

	f.setValue(record, files)
}

func (f *FileField) appendValue(record *Record, toAppend any) {
	files := f.toSliceValue(record.GetRaw(f.Name))
	appends := f.toSliceValue(toAppend)

	if len(appends) > 0 {
		files = append(files, appends...)
	}

	f.setValue(record, files)
}

func (f *FileField) subtractValue(record *Record, toRemove any) {
	files := f.excludeFiles(
		f.toSliceValue(record.GetRaw(f.Name)),
		f.toSliceValue(toRemove),
	)

	f.setValue(record, files)
}

func (f *FileField) normalizeValue(raw any) any {
	files := f.toSliceValue(raw)

	if f.IsMultiple() {
		return files
	}

	if len(files) > 0 {
		return files[len(files)-1] // the last selected
	}

	return ""
}

func (f *FileField) toSliceValue(raw any) []any {
	var result []any

	switch value := raw.(type) {
	case nil:
		// nothing to cast
	case *filesystem.File:
		result = append(result, value)
	case filesystem.File:
		result = append(result, &value)
	case []*filesystem.File:
		for _, v := range value {
			result = append(result, v)
		}
	case []filesystem.File:
		for _, v := range value {
			result = append(result, &v)
		}
	case []any:
		for _, v := range value {
			casted := f.toSliceValue(v)
			if len(casted) == 1 {
				result = append(result, casted[0])
			}
		}
	default:
		result = list.ToInterfaceSlice(list.ToUniqueStringSlice(value))
	}

	return f.uniqueFiles(result)
}

func (f *FileField) uniqueFiles(files []any) []any {
	found := make(map[string]struct{}, len(files))
	result := make([]any, 0, len(files))

	for _, fv := range files {
		name := f.getFileName(fv)
		if _, ok := found[name]; !ok {
			result = append(result, fv)
			found[name] = struct{}{}
		}
	}

	return result
}

func (f *FileField) extractPlainStrings(files []any) []string {
	result := []string{}

	for _, raw := range files {
		if f, ok := raw.(string); ok {
			result = append(result, f)
		}
	}

	return result
}

func (f *FileField) extractUploadableFiles(files []any) []*filesystem.File {
	result := []*filesystem.File{}

	for _, raw := range files {
		if upload, ok := raw.(*filesystem.File); ok {
			result = append(result, upload)
		}
	}

	return result
}

func (f *FileField) excludeFiles(base []any, toExclude []any) []any {
	result := make([]any, 0, len(base))

SUBTRACT_LOOP:
	for _, fv := range base {
		for _, exclude := range toExclude {
			if f.getFileName(exclude) == f.getFileName(fv) {
				continue SUBTRACT_LOOP // skip
			}
		}

		result = append(result, fv)
	}

	return result
}

func (f *FileField) getFileName(file any) string {
	switch v := file.(type) {
	case string:
		return v
	case *filesystem.File:
		return v.Name
	default:
		return ""
	}
}
