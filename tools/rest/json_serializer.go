package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	// Experimental!
	//
	// Need more tests before replacing encoding/json entirely.
	// Test also encoding/json/v2 once released (see https://github.com/golang/go/discussions/63397)
	goccy "github.com/goccy/go-json"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/tokenizer"
)

type FieldModifier interface {
	// Modify executes the modifier and returns a new modified value.
	Modify(value any) (any, error)
}

// Serializer represents custom REST JSON serializer based on echo.DefaultJSONSerializer,
// with support for additional generic response data transformation (eg. fields picker).
type Serializer struct {
	echo.DefaultJSONSerializer

	FieldsParam string
}

// Serialize converts an interface into a json and writes it to the response.
//
// It also provides a generic response data fields picker via the FieldsParam query parameter (default to "fields").
//
// Note: for the places where it is safe, the std encoding/json is replaced
// with goccy due to its slightly better Unmarshal/Marshal performance.
func (s *Serializer) Serialize(c echo.Context, i any, indent string) error {
	fieldsParam := s.FieldsParam
	if fieldsParam == "" {
		fieldsParam = "fields"
	}

	statusCode := c.Response().Status

	rawFields := c.QueryParam(fieldsParam)
	if rawFields == "" || statusCode < 200 || statusCode > 299 {
		return s.DefaultJSONSerializer.Serialize(c, i, indent)
	}

	decoded, err := PickFields(i, rawFields)
	if err != nil {
		return err
	}

	enc := goccy.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(decoded)
}

// PickFields parses the provided fields string expression and
// returns a new subset of data with only the requested fields.
//
// Fields transformations with modifiers are also supported (see initModifer()).
//
// Example:
//
//	data := map[string]any{"a": 1, "b": 2, "c": map[string]any{"c1": 11, "c2": 22}}
//	PickFields(data, "a,c.c1") // map[string]any{"a": 1, "c": map[string]any{"c1": 11}}
func PickFields(data any, rawFields string) (any, error) {
	parsedFields, err := parseFields(rawFields)
	if err != nil {
		return nil, err
	}

	// marshalize the provided data to ensure that the related json.Marshaler
	// implementations are invoked, and then convert it back to a plain
	// json value that we can further operate on.
	//
	// @todo research other approaches to avoid the double serialization
	// ---
	encoded, err := json.Marshal(data) // use the std json since goccy has several bugs reported with struct marshaling and it is not safe
	if err != nil {
		return nil, err
	}

	var decoded any
	if err := goccy.Unmarshal(encoded, &decoded); err != nil {
		return nil, err
	}
	// ---

	// special cases to preserve the same fields format when used with single item or array data.
	var isSearchResult bool
	switch data.(type) {
	case search.Result, *search.Result:
		isSearchResult = true
	}

	if isSearchResult {
		if decodedMap, ok := decoded.(map[string]any); ok {
			pickParsedFields(decodedMap["items"], parsedFields)
		}
	} else {
		pickParsedFields(decoded, parsedFields)
	}

	return decoded, nil
}

func parseFields(rawFields string) (map[string]FieldModifier, error) {
	t := tokenizer.NewFromString(rawFields)

	fields, err := t.ScanAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]FieldModifier, len(fields))

	for _, f := range fields {
		parts := strings.SplitN(strings.TrimSpace(f), ":", 2)

		if len(parts) > 1 {
			m, err := initModifer(parts[1])
			if err != nil {
				return nil, err
			}
			result[parts[0]] = m
		} else {
			result[parts[0]] = nil
		}
	}

	return result, nil
}

func initModifer(rawModifier string) (FieldModifier, error) {
	t := tokenizer.NewFromString(rawModifier)
	t.Separators('(', ')', ',', ' ')
	t.IgnoreParenthesis(true)

	parts, err := t.ScanAll()
	if err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid or empty modifier expression %q", rawModifier)
	}

	name := parts[0]
	args := parts[1:]

	switch name {
	case "excerpt":
		m, err := newExcerptModifier(args...)
		if err != nil {
			return nil, fmt.Errorf("invalid excerpt modifier: %w", err)
		}
		return m, nil
	}

	return nil, fmt.Errorf("missing or invalid modifier %q", name)
}

func pickParsedFields(data any, fields map[string]FieldModifier) error {
	switch v := data.(type) {
	case map[string]any:
		pickMapFields(v, fields)
	case []map[string]any:
		for _, item := range v {
			if err := pickMapFields(item, fields); err != nil {
				return err
			}
		}
	case []any:
		if len(v) == 0 {
			return nil // nothing to pick
		}

		if _, ok := v[0].(map[string]any); !ok {
			return nil // for now ignore non-map values
		}

		for _, item := range v {
			if err := pickMapFields(item.(map[string]any), fields); err != nil {
				return nil
			}
		}
	}

	return nil
}

func pickMapFields(data map[string]any, fields map[string]FieldModifier) error {
	if len(fields) == 0 {
		return nil // nothing to pick
	}

	if m, ok := fields["*"]; ok {
		// append all missing root level data keys
		for k := range data {
			var exists bool

			for f := range fields {
				if strings.HasPrefix(f+".", k+".") {
					exists = true
					break
				}
			}

			if !exists {
				fields[k] = m
			}
		}
	}

DataLoop:
	for k := range data {
		matchingFields := make(map[string]FieldModifier, len(fields))
		for f, m := range fields {
			if strings.HasPrefix(f+".", k+".") {
				matchingFields[f] = m
				continue
			}
		}

		if len(matchingFields) == 0 {
			delete(data, k)
			continue DataLoop
		}

		// remove the current key from the matching fields path
		for f, m := range matchingFields {
			remains := strings.TrimSuffix(strings.TrimPrefix(f+".", k+"."), ".")

			// final key
			if remains == "" {
				if m != nil {
					var err error
					data[k], err = m.Modify(data[k])
					if err != nil {
						return err
					}
				}
				continue DataLoop
			}

			// cleanup the old field key and continue with the rest of the field path
			delete(matchingFields, f)
			matchingFields[remains] = m
		}

		if err := pickParsedFields(data[k], matchingFields); err != nil {
			return err
		}
	}

	return nil
}
