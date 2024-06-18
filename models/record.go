package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

var (
	_ Model             = (*Record)(nil)
	_ ColumnValueMapper = (*Record)(nil)
	_ FilesManager      = (*Record)(nil)
)

type Record struct {
	BaseModel

	collection *Collection

	exportUnknown         bool // whether to export unknown fields
	ignoreEmailVisibility bool // whether to ignore the emailVisibility flag for auth collections
	loaded                bool
	originalData          map[string]any    // the original (aka. first loaded) model data
	expand                *store.Store[any] // expanded relations
	data                  *store.Store[any] // any custom data in addition to the base model fields
}

// NewRecord initializes a new empty Record model.
func NewRecord(collection *Collection) *Record {
	return &Record{
		collection: collection,
		data:       store.New[any](nil),
	}
}

// nullStringMapValue returns the raw string value if it exist and
// its not NULL, otherwise - nil.
func nullStringMapValue(data dbx.NullStringMap, key string) any {
	nullString, ok := data[key]

	if ok && nullString.Valid {
		return nullString.String
	}

	return nil
}

// NewRecordFromNullStringMap initializes a single new Record model
// with data loaded from the provided NullStringMap.
//
// Note that this method is intended to load and Scan data from a database
// result and calls PostScan() which marks the record as "not new".
func NewRecordFromNullStringMap(collection *Collection, data dbx.NullStringMap) *Record {
	resultMap := make(map[string]any, len(data))

	// load schema fields
	for _, field := range collection.Schema.Fields() {
		resultMap[field.Name] = nullStringMapValue(data, field.Name)
	}

	// load base model fields
	for _, name := range schema.BaseModelFieldNames() {
		resultMap[name] = nullStringMapValue(data, name)
	}

	// load auth fields
	if collection.IsAuth() {
		for _, name := range schema.AuthFieldNames() {
			resultMap[name] = nullStringMapValue(data, name)
		}
	}

	record := NewRecord(collection)

	record.Load(resultMap)
	record.PostScan()

	return record
}

// NewRecordsFromNullStringMaps initializes a new Record model for
// each row in the provided NullStringMap slice.
//
// Note that this method is intended to load and Scan data from a database
// result and calls PostScan() for each record marking them as "not new".
func NewRecordsFromNullStringMaps(collection *Collection, rows []dbx.NullStringMap) []*Record {
	result := make([]*Record, len(rows))

	for i, row := range rows {
		result[i] = NewRecordFromNullStringMap(collection, row)
	}

	return result
}

// TableName returns the table name associated to the current Record model.
func (m *Record) TableName() string {
	return m.collection.Name
}

// Collection returns the Collection model associated to the current Record model.
func (m *Record) Collection() *Collection {
	return m.collection
}

// OriginalCopy returns a copy of the current record model populated
// with its ORIGINAL data state (aka. the initially loaded) and
// everything else reset to the defaults.
func (m *Record) OriginalCopy() *Record {
	newRecord := NewRecord(m.collection)
	newRecord.Load(m.originalData)

	if m.IsNew() {
		newRecord.MarkAsNew()
	} else {
		newRecord.MarkAsNotNew()
	}

	return newRecord
}

// CleanCopy returns a copy of the current record model populated only
// with its LATEST data state and everything else reset to the defaults.
func (m *Record) CleanCopy() *Record {
	newRecord := NewRecord(m.collection)
	newRecord.Load(m.data.GetAll())
	newRecord.Id = m.Id
	newRecord.Created = m.Created
	newRecord.Updated = m.Updated

	if m.IsNew() {
		newRecord.MarkAsNew()
	} else {
		newRecord.MarkAsNotNew()
	}

	return newRecord
}

// Expand returns a shallow copy of the current Record model expand data.
func (m *Record) Expand() map[string]any {
	if m.expand == nil {
		m.expand = store.New[any](nil)
	}

	return m.expand.GetAll()
}

// SetExpand shallow copies the provided data to the current Record model's expand.
func (m *Record) SetExpand(expand map[string]any) {
	if m.expand == nil {
		m.expand = store.New[any](nil)
	}

	m.expand.Reset(expand)
}

// MergeExpand merges recursively the provided expand data into
// the current model's expand (if any).
//
// Note that if an expanded prop with the same key is a slice (old or new expand)
// then both old and new records will be merged into a new slice (aka. a :merge: [b,c] => [a,b,c]).
// Otherwise the "old" expanded record will be replace with the "new" one (aka. a :merge: aNew => aNew).
func (m *Record) MergeExpand(expand map[string]any) {
	// nothing to merge
	if len(expand) == 0 {
		return
	}

	// no old expand
	if m.expand == nil {
		m.expand = store.New(expand)
		return
	}

	oldExpand := m.expand.GetAll()

	for key, new := range expand {
		old, ok := oldExpand[key]
		if !ok {
			oldExpand[key] = new
			continue
		}

		var wasOldSlice bool
		var oldSlice []*Record
		switch v := old.(type) {
		case *Record:
			oldSlice = []*Record{v}
		case []*Record:
			wasOldSlice = true
			oldSlice = v
		default:
			// invalid old expand data -> assign directly the new
			// (no matter whether new is valid or not)
			oldExpand[key] = new
			continue
		}

		var wasNewSlice bool
		var newSlice []*Record
		switch v := new.(type) {
		case *Record:
			newSlice = []*Record{v}
		case []*Record:
			wasNewSlice = true
			newSlice = v
		default:
			// invalid new expand data -> skip
			continue
		}

		oldIndexed := make(map[string]*Record, len(oldSlice))
		for _, oldRecord := range oldSlice {
			oldIndexed[oldRecord.Id] = oldRecord
		}

		for _, newRecord := range newSlice {
			oldRecord := oldIndexed[newRecord.Id]
			if oldRecord != nil {
				// note: there is no need to update oldSlice since oldRecord is a reference
				oldRecord.MergeExpand(newRecord.Expand())
			} else {
				// missing new entry
				oldSlice = append(oldSlice, newRecord)
			}
		}

		if wasOldSlice || wasNewSlice || len(oldSlice) == 0 {
			oldExpand[key] = oldSlice
		} else {
			oldExpand[key] = oldSlice[0]
		}
	}

	m.expand.Reset(oldExpand)
}

// SchemaData returns a shallow copy ONLY of the defined record schema fields data.
func (m *Record) SchemaData() map[string]any {
	result := make(map[string]any, len(m.collection.Schema.Fields()))

	data := m.data.GetAll()

	for _, field := range m.collection.Schema.Fields() {
		if v, ok := data[field.Name]; ok {
			result[field.Name] = v
		}
	}

	return result
}

// UnknownData returns a shallow copy ONLY of the unknown record fields data,
// aka. fields that are neither one of the base and special system ones,
// nor defined by the collection schema.
func (m *Record) UnknownData() map[string]any {
	if m.data == nil {
		return nil
	}

	return m.extractUnknownData(m.data.GetAll())
}

// IgnoreEmailVisibility toggles the flag to ignore the auth record email visibility check.
func (m *Record) IgnoreEmailVisibility(state bool) {
	m.ignoreEmailVisibility = state
}

// WithUnknownData toggles the export/serialization of unknown data fields
// (false by default).
func (m *Record) WithUnknownData(state bool) {
	m.exportUnknown = state
}

// Set sets the provided key-value data pair for the current Record model.
//
// If the record collection has field with name matching the provided "key",
// the value will be further normalized according to the field rules.
func (m *Record) Set(key string, value any) {
	switch key {
	case schema.FieldNameId:
		m.Id = cast.ToString(value)
	case schema.FieldNameCreated:
		m.Created, _ = types.ParseDateTime(value)
	case schema.FieldNameUpdated:
		m.Updated, _ = types.ParseDateTime(value)
	case schema.FieldNameExpand:
		m.SetExpand(cast.ToStringMap(value))
	default:
		var v = value

		if field := m.Collection().Schema.GetFieldByName(key); field != nil {
			v = field.PrepareValue(value)
		} else if m.collection.IsAuth() {
			// normalize auth fields
			switch key {
			case schema.FieldNameEmailVisibility, schema.FieldNameVerified:
				v = cast.ToBool(value)
			case schema.FieldNameLastResetSentAt, schema.FieldNameLastVerificationSentAt, schema.FieldNameLastLoginAlertSentAt:
				v, _ = types.ParseDateTime(value)
			case schema.FieldNameUsername, schema.FieldNameEmail, schema.FieldNameTokenKey, schema.FieldNamePasswordHash:
				v = cast.ToString(value)
			}
		}

		if m.data == nil {
			m.data = store.New[any](nil)
		}

		m.data.Set(key, v)
	}
}

// Get returns a normalized single record model data value for "key".
func (m *Record) Get(key string) any {
	switch key {
	case schema.FieldNameId:
		return m.Id
	case schema.FieldNameCreated:
		return m.Created
	case schema.FieldNameUpdated:
		return m.Updated
	default:
		var v any
		if m.data != nil {
			v = m.data.Get(key)
		}

		// normalize the field value in case it is missing or an incorrect type was set
		// to ensure that the DB will always have normalized columns value.
		if field := m.Collection().Schema.GetFieldByName(key); field != nil {
			v = field.PrepareValue(v)
		} else if m.collection.IsAuth() {
			switch key {
			case schema.FieldNameEmailVisibility, schema.FieldNameVerified:
				v = cast.ToBool(v)
			case schema.FieldNameLastResetSentAt, schema.FieldNameLastVerificationSentAt, schema.FieldNameLastLoginAlertSentAt:
				v, _ = types.ParseDateTime(v)
			case schema.FieldNameUsername, schema.FieldNameEmail, schema.FieldNameTokenKey, schema.FieldNamePasswordHash:
				v = cast.ToString(v)
			}
		}

		return v
	}
}

// GetBool returns the data value for "key" as a bool.
func (m *Record) GetBool(key string) bool {
	return cast.ToBool(m.Get(key))
}

// GetString returns the data value for "key" as a string.
func (m *Record) GetString(key string) string {
	return cast.ToString(m.Get(key))
}

// GetInt returns the data value for "key" as an int.
func (m *Record) GetInt(key string) int {
	return cast.ToInt(m.Get(key))
}

// GetFloat returns the data value for "key" as a float64.
func (m *Record) GetFloat(key string) float64 {
	return cast.ToFloat64(m.Get(key))
}

// GetTime returns the data value for "key" as a [time.Time] instance.
func (m *Record) GetTime(key string) time.Time {
	return cast.ToTime(m.Get(key))
}

// GetDateTime returns the data value for "key" as a DateTime instance.
func (m *Record) GetDateTime(key string) types.DateTime {
	d, _ := types.ParseDateTime(m.Get(key))
	return d
}

// GetStringSlice returns the data value for "key" as a slice of unique strings.
func (m *Record) GetStringSlice(key string) []string {
	return list.ToUniqueStringSlice(m.Get(key))
}

// ExpandedOne retrieves a single relation Record from the already
// loaded expand data of the current model.
//
// If the requested expand relation is multiple, this method returns
// only first available Record from the expanded relation.
//
// Returns nil if there is no such expand relation loaded.
func (m *Record) ExpandedOne(relField string) *Record {
	if m.expand == nil {
		return nil
	}

	rel := m.expand.Get(relField)

	switch v := rel.(type) {
	case *Record:
		return v
	case []*Record:
		if len(v) > 0 {
			return v[0]
		}
	}

	return nil
}

// ExpandedAll retrieves a slice of relation Records from the already
// loaded expand data of the current model.
//
// If the requested expand relation is single, this method normalizes
// the return result and will wrap the single model as a slice.
//
// Returns nil slice if there is no such expand relation loaded.
func (m *Record) ExpandedAll(relField string) []*Record {
	if m.expand == nil {
		return nil
	}

	rel := m.expand.Get(relField)

	switch v := rel.(type) {
	case *Record:
		return []*Record{v}
	case []*Record:
		return v
	}

	return nil
}

// Retrieves the "key" json field value and unmarshals it into "result".
//
// Example
//
//	result := struct {
//	    FirstName string `json:"first_name"`
//	}{}
//	err := m.UnmarshalJSONField("my_field_name", &result)
func (m *Record) UnmarshalJSONField(key string, result any) error {
	return json.Unmarshal([]byte(m.GetString(key)), &result)
}

// BaseFilesPath returns the storage dir path used by the record.
func (m *Record) BaseFilesPath() string {
	return fmt.Sprintf("%s/%s", m.Collection().BaseFilesPath(), m.Id)
}

// FindFileFieldByFile returns the first file type field for which
// any of the record's data contains the provided filename.
func (m *Record) FindFileFieldByFile(filename string) *schema.SchemaField {
	for _, field := range m.Collection().Schema.Fields() {
		if field.Type == schema.FieldTypeFile {
			names := m.GetStringSlice(field.Name)
			if list.ExistInSlice(filename, names) {
				return field
			}
		}
	}
	return nil
}

// Load bulk loads the provided data into the current Record model.
func (m *Record) Load(data map[string]any) {
	if !m.loaded {
		m.loaded = true
		m.originalData = data
	}

	for k, v := range data {
		m.Set(k, v)
	}
}

// ColumnValueMap implements [ColumnValueMapper] interface.
func (m *Record) ColumnValueMap() map[string]any {
	result := make(map[string]any, len(m.collection.Schema.Fields())+3)

	// export schema field values
	for _, field := range m.collection.Schema.Fields() {
		result[field.Name] = m.getNormalizeDataValueForDB(field.Name)
	}

	// export auth collection fields
	if m.collection.IsAuth() {
		for _, name := range schema.AuthFieldNames() {
			result[name] = m.getNormalizeDataValueForDB(name)
		}
	}

	// export base model fields
	result[schema.FieldNameId] = m.getNormalizeDataValueForDB(schema.FieldNameId)
	result[schema.FieldNameCreated] = m.getNormalizeDataValueForDB(schema.FieldNameCreated)
	result[schema.FieldNameUpdated] = m.getNormalizeDataValueForDB(schema.FieldNameUpdated)

	return result
}

// PublicExport exports only the record fields that are safe to be public.
//
// For auth records, to force the export of the email field you need to set
// `m.IgnoreEmailVisibility(true)`.
func (m *Record) PublicExport() map[string]any {
	result := make(map[string]any, len(m.collection.Schema.Fields())+5)

	// export unknown data fields if allowed
	if m.exportUnknown {
		for k, v := range m.UnknownData() {
			result[k] = v
		}
	}

	// export schema field values
	for _, field := range m.collection.Schema.Fields() {
		result[field.Name] = m.Get(field.Name)
	}

	// export some of the safe auth collection fields
	if m.collection.IsAuth() {
		result[schema.FieldNameVerified] = m.Verified()
		result[schema.FieldNameUsername] = m.Username()
		result[schema.FieldNameEmailVisibility] = m.EmailVisibility()
		if m.ignoreEmailVisibility || m.EmailVisibility() {
			result[schema.FieldNameEmail] = m.Email()
		}
	}

	// export base model fields
	result[schema.FieldNameId] = m.GetId()
	if created := m.GetCreated(); !m.Collection().IsView() || !created.IsZero() {
		result[schema.FieldNameCreated] = created
	}
	if updated := m.GetUpdated(); !m.Collection().IsView() || !updated.IsZero() {
		result[schema.FieldNameUpdated] = updated
	}

	// add helper collection reference fields
	result[schema.FieldNameCollectionId] = m.collection.Id
	result[schema.FieldNameCollectionName] = m.collection.Name

	// add expand (if set)
	if m.expand != nil && m.expand.Length() > 0 {
		result[schema.FieldNameExpand] = m.expand.GetAll()
	}

	return result
}

// MarshalJSON implements the [json.Marshaler] interface.
//
// Only the data exported by `PublicExport()` will be serialized.
func (m Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.PublicExport())
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (m *Record) UnmarshalJSON(data []byte) error {
	result := map[string]any{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	m.Load(result)

	return nil
}

// ReplaceModifers returns a new map with applied modifier
// values based on the current record and the specified data.
//
// The resolved modifier keys will be removed.
//
// Multiple modifiers will be applied one after another,
// while reusing the previous base key value result (eg. 1; -5; +2 => -2).
//
// Example usage:
//
//	 newData := record.ReplaceModifers(data)
//		// record:  {"field": 10}
//		// data:    {"field+": 5}
//		// newData: {"field": 15}
func (m *Record) ReplaceModifers(data map[string]any) map[string]any {
	var clone = shallowCopy(data)
	if len(clone) == 0 {
		return clone
	}

	var recordDataCache map[string]any

	// export recordData lazily
	recordData := func() map[string]any {
		if recordDataCache == nil {
			recordDataCache = m.SchemaData()
		}
		return recordDataCache
	}

	modifiers := schema.FieldValueModifiers()

	for _, field := range m.Collection().Schema.Fields() {
		key := field.Name

		for _, m := range modifiers {
			if mv, mOk := clone[key+m]; mOk {
				if _, ok := clone[key]; !ok {
					// get base value from the merged data
					clone[key] = recordData()[key]
				}

				clone[key] = field.PrepareValueWithModifier(clone[key], m, mv)
				delete(clone, key+m)
			}
		}

		if field.Type != schema.FieldTypeFile {
			continue
		}

		// -----------------------------------------------------------
		// legacy file field modifiers (kept for backward compatibility)
		// -----------------------------------------------------------

		var oldNames []string
		var toDelete []string
		if _, ok := clone[key]; ok {
			oldNames = list.ToUniqueStringSlice(clone[key])
		} else {
			// get oldNames from the model
			oldNames = list.ToUniqueStringSlice(recordData()[key])
		}

		// search for individual file name to delete (eg. "file.test.png = null")
		for _, name := range oldNames {
			suffixedKey := key + "." + name
			if v, ok := clone[suffixedKey]; ok && cast.ToString(v) == "" {
				toDelete = append(toDelete, name)
				delete(clone, suffixedKey)
				continue
			}
		}

		// search for individual file index to delete (eg. "file.0 = null")
		keyExp, _ := regexp.Compile(`^` + regexp.QuoteMeta(key) + `\.\d+$`)
		for indexedKey := range clone {
			if keyExp.MatchString(indexedKey) && cast.ToString(clone[indexedKey]) == "" {
				index, indexErr := strconv.Atoi(indexedKey[len(key)+1:])
				if indexErr != nil || index < 0 || index >= len(oldNames) {
					continue
				}
				toDelete = append(toDelete, oldNames[index])
				delete(clone, indexedKey)
			}
		}

		if toDelete != nil {
			clone[key] = field.PrepareValue(list.SubtractSlice(oldNames, toDelete))
		}
	}

	return clone
}

// getNormalizeDataValueForDB returns the "key" data value formatted for db storage.
func (m *Record) getNormalizeDataValueForDB(key string) any {
	var val any

	// normalize auth fields
	if m.collection.IsAuth() {
		switch key {
		case schema.FieldNameEmailVisibility, schema.FieldNameVerified:
			return m.GetBool(key)
		case schema.FieldNameLastResetSentAt, schema.FieldNameLastVerificationSentAt, schema.FieldNameLastLoginAlertSentAt:
			return m.GetDateTime(key)
		case schema.FieldNameUsername, schema.FieldNameEmail, schema.FieldNameTokenKey, schema.FieldNamePasswordHash:
			return m.GetString(key)
		}
	}

	val = m.Get(key)

	switch ids := val.(type) {
	case []string:
		// encode string slice
		return append(types.JsonArray[string]{}, ids...)
	case []int:
		// encode int slice
		return append(types.JsonArray[int]{}, ids...)
	case []float64:
		// encode float64 slice
		return append(types.JsonArray[float64]{}, ids...)
	case []any:
		// encode interface slice
		return append(types.JsonArray[any]{}, ids...)
	default:
		// no changes
		return val
	}
}

// shallowCopy shallow copy data into a new map.
func shallowCopy(data map[string]any) map[string]any {
	result := make(map[string]any, len(data))

	for k, v := range data {
		result[k] = v
	}

	return result
}

func (m *Record) extractUnknownData(data map[string]any) map[string]any {
	knownFields := map[string]struct{}{}

	for _, name := range schema.SystemFieldNames() {
		knownFields[name] = struct{}{}
	}
	for _, name := range schema.BaseModelFieldNames() {
		knownFields[name] = struct{}{}
	}

	for _, f := range m.collection.Schema.Fields() {
		knownFields[f.Name] = struct{}{}
	}

	if m.collection.IsAuth() {
		for _, name := range schema.AuthFieldNames() {
			knownFields[name] = struct{}{}
		}
	}

	result := map[string]any{}

	for k, v := range data {
		if _, ok := knownFields[k]; !ok {
			result[k] = v
		}
	}

	return result
}

// -------------------------------------------------------------------
// Auth helpers
// -------------------------------------------------------------------

var notAuthRecordErr = errors.New("Not an auth collection record.")

// Username returns the "username" auth record data value.
func (m *Record) Username() string {
	return m.GetString(schema.FieldNameUsername)
}

// SetUsername sets the "username" auth record data value.
//
// This method doesn't check whether the provided value is a valid username.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetUsername(username string) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameUsername, username)

	return nil
}

// Email returns the "email" auth record data value.
func (m *Record) Email() string {
	return m.GetString(schema.FieldNameEmail)
}

// SetEmail sets the "email" auth record data value.
//
// This method doesn't check whether the provided value is a valid email.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetEmail(email string) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameEmail, email)

	return nil
}

// Verified returns the "emailVisibility" auth record data value.
func (m *Record) EmailVisibility() bool {
	return m.GetBool(schema.FieldNameEmailVisibility)
}

// SetEmailVisibility sets the "emailVisibility" auth record data value.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetEmailVisibility(visible bool) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameEmailVisibility, visible)

	return nil
}

// Verified returns the "verified" auth record data value.
func (m *Record) Verified() bool {
	return m.GetBool(schema.FieldNameVerified)
}

// SetVerified sets the "verified" auth record data value.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetVerified(verified bool) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameVerified, verified)

	return nil
}

// TokenKey returns the "tokenKey" auth record data value.
func (m *Record) TokenKey() string {
	return m.GetString(schema.FieldNameTokenKey)
}

// SetTokenKey sets the "tokenKey" auth record data value.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetTokenKey(key string) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameTokenKey, key)

	return nil
}

// RefreshTokenKey generates and sets new random auth record "tokenKey".
//
// Returns an error if the record is not from an auth collection.
func (m *Record) RefreshTokenKey() error {
	return m.SetTokenKey(security.RandomString(50))
}

// LastResetSentAt returns the "lastResentSentAt" auth record data value.
func (m *Record) LastResetSentAt() types.DateTime {
	return m.GetDateTime(schema.FieldNameLastResetSentAt)
}

// SetLastResetSentAt sets the "lastResentSentAt" auth record data value.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetLastResetSentAt(dateTime types.DateTime) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameLastResetSentAt, dateTime)

	return nil
}

// LastVerificationSentAt returns the "lastVerificationSentAt" auth record data value.
func (m *Record) LastVerificationSentAt() types.DateTime {
	return m.GetDateTime(schema.FieldNameLastVerificationSentAt)
}

// SetLastVerificationSentAt sets an "lastVerificationSentAt" auth record data value.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetLastVerificationSentAt(dateTime types.DateTime) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameLastVerificationSentAt, dateTime)

	return nil
}

// LastLoginAlertSentAt returns the "lastLoginAlertSentAt" auth record data value.
func (m *Record) LastLoginAlertSentAt() types.DateTime {
	return m.GetDateTime(schema.FieldNameLastLoginAlertSentAt)
}

// SetLastLoginAlertSentAt sets an "lastLoginAlertSentAt" auth record data value.
//
// Returns an error if the record is not from an auth collection.
func (m *Record) SetLastLoginAlertSentAt(dateTime types.DateTime) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	m.Set(schema.FieldNameLastLoginAlertSentAt, dateTime)

	return nil
}

// PasswordHash returns the "passwordHash" auth record data value.
func (m *Record) PasswordHash() string {
	return m.GetString(schema.FieldNamePasswordHash)
}

// ValidatePassword validates a plain password against the auth record password.
//
// Returns false if the password is incorrect or record is not from an auth collection.
func (m *Record) ValidatePassword(password string) bool {
	if !m.collection.IsAuth() {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(m.PasswordHash()), []byte(password))

	return err == nil
}

// SetPassword sets cryptographically secure string to the auth record "password" field.
// This method also resets the "lastResetSentAt" and the "tokenKey" fields.
//
// Returns an error if the record is not from an auth collection or
// an empty password is provided.
func (m *Record) SetPassword(password string) error {
	if !m.collection.IsAuth() {
		return notAuthRecordErr
	}

	if password == "" {
		return errors.New("The provided plain password is empty")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	m.Set(schema.FieldNamePasswordHash, string(hashedPassword))
	m.Set(schema.FieldNameLastResetSentAt, types.DateTime{})

	// invalidate previously issued tokens
	return m.RefreshTokenKey()
}
