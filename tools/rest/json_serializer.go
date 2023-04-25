package rest

import (
	"encoding/json"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/search"
)

// Serializer represents custom REST JSON serializer based on echo.DefaultJSONSerializer,
// with support for additional generic response data transformation (eg. fields picker).
type Serializer struct {
	echo.DefaultJSONSerializer

	FieldsParam string
}

// Serialize converts an interface into a json and writes it to the response.
//
// It also provides a generic response data fields picker via the FieldsParam query parameter (default to "fields").
func (s *Serializer) Serialize(c echo.Context, i any, indent string) error {
	fieldsParam := s.FieldsParam
	if fieldsParam == "" {
		fieldsParam = "fields"
	}

	param := c.QueryParam(fieldsParam)
	if param == "" {
		return s.DefaultJSONSerializer.Serialize(c, i, indent)
	}

	fields := strings.Split(param, ",")
	for i, f := range fields {
		fields[i] = strings.TrimSpace(f)
	}

	encoded, err := json.Marshal(i)
	if err != nil {
		return err
	}

	var decoded any

	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return err
	}

	var isSearchResult bool

	switch i.(type) {
	case search.Result, *search.Result:
		isSearchResult = true
	}

	if isSearchResult {
		if decodedMap, ok := decoded.(map[string]any); ok {
			pickFields(decodedMap["items"], fields)
		}
	} else {
		pickFields(decoded, fields)
	}

	return s.DefaultJSONSerializer.Serialize(c, decoded, indent)
}

func pickFields(data any, fields []string) {
	switch v := data.(type) {
	case map[string]any:
		pickMapFields(v, fields)
	case []map[string]any:
		for _, item := range v {
			pickMapFields(item, fields)
		}
	case []any:
		if len(v) == 0 {
			return // nothing to pick
		}

		if _, ok := v[0].(map[string]any); !ok {
			return // for now ignore non-map values
		}

		for _, item := range v {
			pickMapFields(item.(map[string]any), fields)
		}
	}
}

func pickMapFields(data map[string]any, fields []string) {
	if len(fields) == 0 {
		return // nothing to pick
	}

DataLoop:
	for k := range data {
		matchingFields := make([]string, 0, len(fields))
		for _, f := range fields {
			if strings.HasPrefix(f+".", k+".") {
				matchingFields = append(matchingFields, f)
				continue
			}
		}

		if len(matchingFields) == 0 {
			delete(data, k)
			continue DataLoop
		}

		// trim the key from the fields
		for i, v := range matchingFields {
			trimmed := strings.TrimSuffix(strings.TrimPrefix(v+".", k+"."), ".")
			if trimmed == "" {
				continue DataLoop
			}
			matchingFields[i] = trimmed
		}

		pickFields(data[k], matchingFields)
	}
}
