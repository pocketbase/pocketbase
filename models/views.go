package models

import (
	"encoding/json"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models/schema"
)

var _ Model = (*View)(nil)

type View struct {
	BaseModel

	Sql      string        `db:"sql" json:"sql"`
	Name     string        `db:"name" json:"name"`
	ListRule *string       `db:"listRule" json:"listRule"`
	Schema   schema.Schema `db:"schema" json:"schema"`
}

func (m *View) TableName() string {
	return "_views"
}

func recordFromViewSchema(view *View, data dbx.NullStringMap) map[string]any {
	resultMap := map[string]any{}

	for _, field := range view.Schema.Fields() {
		var rawValue any = nil

		nullString, ok := data[field.Name]
		if ok && nullString.Valid {
			rawValue = parseNullStringToViewRecord(nullString.String, field.Type)
		}
		resultMap[field.Name] = rawValue
	}
	return resultMap
}

func NewRecordsFromViewSchema(view *View, rows []dbx.NullStringMap) []map[string]any {
	result := make([]map[string]any, len(rows))

	for i, row := range rows {
		// newRow := replaceColonInColumns(&row)
		result[i] = recordFromViewSchema(view, row)
	}

	return result
}

func parseNullStringToViewRecord(data string, fieldType string) any {
	var rawValue any = data
	switch fieldType {
	case schema.FieldTypeBool:
		rawValue, _ = strconv.ParseBool(data)
	case schema.FieldTypeNumber:
		rawValue, _ = strconv.Atoi(data)
	case schema.FieldTypeJson:
		json.Unmarshal([]byte(data), &rawValue)
	}
	return rawValue
}
