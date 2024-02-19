package dbutils

import (
	"fmt"
	"strings"
)

// JsonEach returns JSON_EACH SQLite string expression with
// some normalizations for non-json columns.
func JsonEach(column string) string {
	return fmt.Sprintf(
		`json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END)`,
		column, column, column,
	)
}

// JsonArrayLength returns JSON_ARRAY_LENGTH SQLite string expression
// with some normalizations for non-json columns.
//
// It works with both json and non-json column values.
//
// Returns 0 for empty string or NULL column values.
func JsonArrayLength(column string) string {
	return fmt.Sprintf(
		`json_array_length(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE (CASE WHEN [[%s]] = '' OR [[%s]] IS NULL THEN json_array() ELSE json_array([[%s]]) END) END)`,
		column, column, column, column, column,
	)
}

// JsonExtract returns a JSON_EXTRACT SQLite string expression with
// some normalizations for non-json columns.
func JsonExtract(column string, path string) string {
	// prefix the path with dot if it is not starting with array notation
	if path != "" && !strings.HasPrefix(path, "[") {
		path = "." + path
	}

	return fmt.Sprintf(
		// note: the extra object wrapping is needed to workaround the cases where a json_extract is used with non-json columns.
		"(CASE WHEN json_valid([[%s]]) THEN JSON_EXTRACT([[%s]], '$%s') ELSE JSON_EXTRACT(json_object('pb', [[%s]]), '$.pb%s') END)",
		column,
		column,
		path,
		column,
		path,
	)
}
