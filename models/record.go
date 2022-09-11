package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

var _ Model = (*Record)(nil)
var _ ColumnValueMapper = (*Record)(nil)
var _ FilesManager = (*Record)(nil)

type Record struct {
	BaseModel

	collection *Collection
	data       map[string]any
	expand     map[string]any
}

// NewRecord initializes a new empty Record model.
func NewRecord(collection *Collection) *Record {
	return &Record{
		collection: collection,
		data:       map[string]any{},
	}
}

// NewRecordFromNullStringMap initializes a single new Record model
// with data loaded from the provided NullStringMap.
func NewRecordFromNullStringMap(collection *Collection, data dbx.NullStringMap) *Record {
	resultMap := map[string]any{}

	for _, field := range collection.Schema.Fields() {
		var rawValue any

		nullString, ok := data[field.Name]
		if !ok || !nullString.Valid {
			rawValue = nil
		} else {
			rawValue = nullString.String
		}

		resultMap[field.Name] = rawValue
	}

	record := NewRecord(collection)

	// load base mode fields
	resultMap[schema.ReservedFieldNameId] = data[schema.ReservedFieldNameId].String
	resultMap[schema.ReservedFieldNameCreated] = data[schema.ReservedFieldNameCreated].String
	resultMap[schema.ReservedFieldNameUpdated] = data[schema.ReservedFieldNameUpdated].String

	if err := record.Load(resultMap); err != nil {
		log.Println("Failed to unmarshal record:", err)
	}

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

// GetExpand returns a shallow copy of  the optional `expand` data
// attached to the current Record model.
func (m *Record) GetExpand() map[string]any {
	return shallowCopy(m.expand)
}

// SetExpand assigns the provided data to `record.expand`.
func (m *Record) SetExpand(data map[string]any) {
	m.expand = shallowCopy(data)
}

// Data returns a shallow copy of the currently loaded record's data.
func (m *Record) Data() map[string]any {
	return shallowCopy(m.data)
}

// SetDataValue sets the provided key-value data pair for the current Record model.
//
// This method does nothing if the record doesn't have a `key` field.
func (m *Record) SetDataValue(key string, value any) {
	if m.data == nil {
		m.data = map[string]any{}
	}

	field := m.Collection().Schema.GetFieldByName(key)
	if field != nil {
		m.data[key] = field.PrepareValue(value)
	}
}

// GetDataValue returns the current record's data value for `key`.
//
// Returns nil if data value with `key` is not found or set.
func (m *Record) GetDataValue(key string) any {
	return m.data[key]
}

// GetBoolDataValue returns the data value for `key` as a bool.
func (m *Record) GetBoolDataValue(key string) bool {
	return cast.ToBool(m.GetDataValue(key))
}

// GetStringDataValue returns the data value for `key` as a string.
func (m *Record) GetStringDataValue(key string) string {
	return cast.ToString(m.GetDataValue(key))
}

// GetIntDataValue returns the data value for `key` as an int.
func (m *Record) GetIntDataValue(key string) int {
	return cast.ToInt(m.GetDataValue(key))
}

// GetFloatDataValue returns the data value for `key` as a float64.
func (m *Record) GetFloatDataValue(key string) float64 {
	return cast.ToFloat64(m.GetDataValue(key))
}

// GetTimeDataValue returns the data value for `key` as a [time.Time] instance.
func (m *Record) GetTimeDataValue(key string) time.Time {
	return cast.ToTime(m.GetDataValue(key))
}

// GetDateTimeDataValue returns the data value for `key` as a DateTime instance.
func (m *Record) GetDateTimeDataValue(key string) types.DateTime {
	d, _ := types.ParseDateTime(m.GetDataValue(key))
	return d
}

// GetStringSliceDataValue returns the data value for `key` as a slice of unique strings.
func (m *Record) GetStringSliceDataValue(key string) []string {
	return list.ToUniqueStringSlice(m.GetDataValue(key))
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
			names := m.GetStringSliceDataValue(field.Name)
			if list.ExistInSlice(filename, names) {
				return field
			}
		}
	}
	return nil
}

// Load bulk loads the provided data into the current Record model.
func (m *Record) Load(data map[string]any) error {
	if data[schema.ReservedFieldNameId] != nil {
		id, err := cast.ToStringE(data[schema.ReservedFieldNameId])
		if err != nil {
			return err
		}
		m.Id = id
	}

	if data[schema.ReservedFieldNameCreated] != nil {
		m.Created, _ = types.ParseDateTime(data[schema.ReservedFieldNameCreated])
	}

	if data[schema.ReservedFieldNameUpdated] != nil {
		m.Updated, _ = types.ParseDateTime(data[schema.ReservedFieldNameUpdated])
	}

	for k, v := range data {
		m.SetDataValue(k, v)
	}

	return nil
}

// ColumnValueMap implements [ColumnValueMapper] interface.
func (m *Record) ColumnValueMap() map[string]any {
	result := map[string]any{}
	for key := range m.data {
		result[key] = m.normalizeDataValueForDB(key)
	}

	// set base model fields
	result[schema.ReservedFieldNameId] = m.Id
	result[schema.ReservedFieldNameCreated] = m.Created
	result[schema.ReservedFieldNameUpdated] = m.Updated

	return result
}

// PublicExport exports only the record fields that are safe to be public.
//
// This method also skips the "hidden" fields, aka. fields prefixed with `#`.
func (m *Record) PublicExport() map[string]any {
	result := skipHiddenFields(m.data)

	// set base model fields
	result[schema.ReservedFieldNameId] = m.Id
	result[schema.ReservedFieldNameCreated] = m.Created
	result[schema.ReservedFieldNameUpdated] = m.Updated

	// add helper collection fields
	result["@collectionId"] = m.collection.Id
	result["@collectionName"] = m.collection.Name

	// add expand (if set)
	if m.expand != nil {
		result["@expand"] = m.expand
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

	return m.Load(result)
}

// normalizeDataValueForDB returns the `key` data value formatted for db storage.
func (m *Record) normalizeDataValueForDB(key string) any {
	val := m.GetDataValue(key)

	switch ids := val.(type) {
	case []string:
		// encode strings slice
		return append(types.JsonArray{}, list.ToInterfaceSlice(ids)...)
	case []any:
		// encode interfaces slice
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

// skipHiddenFields returns a new data map without the "#" prefixed fields.
func skipHiddenFields(data map[string]any) map[string]any {
	result := map[string]any{}

	for key, val := range data {
		// ignore "#" prefixed fields
		if strings.HasPrefix(key, "#") {
			continue
		}
		result[key] = val
	}

	return result
}
