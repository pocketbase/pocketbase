package picker

import (
	"encoding/json"
	"strings"

	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/tokenizer"
)

// Pick converts data into a []any, map[string]any, etc. (using json marshal->unmarshal)
// containing only the fields from the parsed rawFields expression.
//
// rawFields is a comma separated string of the fields to include.
// Nested fields should be listed with dot-notation.
// Fields value modifiers are also supported using the `:modifier(args)` format (see Modifiers).
//
// Example:
//
//	data := map[string]any{"a": 1, "b": 2, "c": map[string]any{"c1": 11, "c2": 22}}
//	Pick(data, "a,c.c1") // map[string]any{"a": 1, "c": map[string]any{"c1": 11}}
func Pick(data any, rawFields string) (any, error) {
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
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var decoded any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, err
	}
	// ---

	// special cases to preserve the same fields format when used with single item or search results data.
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

func parseFields(rawFields string) (map[string]Modifier, error) {
	t := tokenizer.NewFromString(rawFields)

	fields, err := t.ScanAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]Modifier, len(fields))

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

func pickParsedFields(data any, fields map[string]Modifier) error {
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

func pickMapFields(data map[string]any, fields map[string]Modifier) error {
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
		matchingFields := make(map[string]Modifier, len(fields))
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
