package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
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

	exportUnknown         bool           // whether to export unknown fields
	ignoreEmailVisibility bool           // whether to ignore the emailVisibility flag for auth collections
	data                  map[string]any // any custom data in addition to the base model fields
	expand                map[string]any // expanded relations
	loaded                bool
	originalData          map[string]any // the original (aka. first loaded) model data
}

// NewRecord initializes a new empty Record model.
func NewRecord(collection *Collection) *Record {
	return &Record{
		collection: collection,
		data:       map[string]any{},
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
func NewRecordFromNullStringMap(collection *Collection, data dbx.NullStringMap) *Record {
	resultMap := map[string]any{}

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
// with its original (aka. the initially loaded) data state.
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

// Expand returns a shallow copy of  the record.expand data
// attached to the current Record model.
func (m *Record) Expand() map[string]any {
	return shallowCopy(m.expand)
}

// SetExpand assigns the provided data to record.expand.
func (m *Record) SetExpand(expand map[string]any) {
	m.expand = shallowCopy(expand)
}

// SchemaData returns a shallow copy ONLY of the defined record schema fields data.
func (m *Record) SchemaData() map[string]any {
	result := map[string]any{}

	for _, field := range m.collection.Schema.Fields() {
		if v, ok := m.data[field.Name]; ok {
			result[field.Name] = v
		}
	}

	return result
}

// UnknownData returns a shallow copy ONLY of the unknown record fields data,
// aka. fields that are neither one of the base and special system ones,
// nor defined by the collection schema.
func (m *Record) UnknownData() map[string]any {
	return m.extractUnknownData(m.data)
}

// IgnoreEmailVisibility toggles the flag to ignore the auth record email visibility check.
func (m *Record) IgnoreEmailVisibility(state bool) {
	m.ignoreEmailVisibility = state
}

// WithUnkownData toggles the export/serialization of unknown data fields
// (false by default).
func (m *Record) WithUnkownData(state bool) {
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
			case schema.FieldNameLastResetSentAt, schema.FieldNameLastVerificationSentAt:
				v, _ = types.ParseDateTime(value)
			case schema.FieldNameUsername, schema.FieldNameEmail, schema.FieldNameTokenKey, schema.FieldNamePasswordHash:
				v = cast.ToString(value)
			}
		}

		if m.data == nil {
			m.data = map[string]any{}
		}

		m.data[key] = v
	}
}

// Get returns a single record model data value for "key".
func (m *Record) Get(key string) any {
	switch key {
	case schema.FieldNameId:
		return m.Id
	case schema.FieldNameCreated:
		return m.Created
	case schema.FieldNameUpdated:
		return m.Updated
	default:
		if v, ok := m.data[key]; ok {
			return v
		}

		return nil
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

// Retrieves the "key" json field value and unmarshals it into "result".
//
// Example
//  result := struct {
//      FirstName string `json:"first_name"`
//  }{}
//  err := m.UnmarshalJSONField("my_field_name", &result)
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
	result := map[string]any{}

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
// Fields marked as hidden will be exported only if `m.IgnoreEmailVisibility(true)` is set.
func (m *Record) PublicExport() map[string]any {
	result := map[string]any{}

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
	result[schema.FieldNameCreated] = m.GetCreated()
	result[schema.FieldNameUpdated] = m.GetUpdated()

	// add helper collection reference fields
	result[schema.FieldNameCollectionId] = m.collection.Id
	result[schema.FieldNameCollectionName] = m.collection.Name

	// add expand (if set)
	if m.expand != nil {
		result[schema.FieldNameExpand] = m.expand
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

// getNormalizeDataValueForDB returns the "key" data value formatted for db storage.
func (m *Record) getNormalizeDataValueForDB(key string) any {
	var val any

	// normalize auth fields
	if m.collection.IsAuth() {
		switch key {
		case schema.FieldNameEmailVisibility, schema.FieldNameVerified:
			return m.GetBool(key)
		case schema.FieldNameLastResetSentAt, schema.FieldNameLastVerificationSentAt:
			return m.GetDateTime(key)
		case schema.FieldNameUsername, schema.FieldNameEmail, schema.FieldNameTokenKey, schema.FieldNamePasswordHash:
			return m.GetString(key)
		}
	}

	val = m.Get(key)

	switch ids := val.(type) {
	case []string:
		// encode string slice
		return append(types.JsonArray{}, list.ToInterfaceSlice(ids)...)
	case []int:
		// encode int slice
		return append(types.JsonArray{}, list.ToInterfaceSlice(ids)...)
	case []float64:
		// encode float64 slice
		return append(types.JsonArray{}, list.ToInterfaceSlice(ids)...)
	case []any:
		// encode interface slice
		return append(types.JsonArray{}, ids...)
	default:
		// no changes
		return val
	}
}

// shallowCopy shallow copy data into a new map.
func shallowCopy(data map[string]any) map[string]any {
	result := map[string]any{}

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

	for k, v := range m.data {
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return err
	}

	m.Set(schema.FieldNamePasswordHash, string(hashedPassword))
	m.Set(schema.FieldNameLastResetSentAt, types.DateTime{})

	// invalidate previously issued tokens
	return m.RefreshTokenKey()
}
