package search

import (
	"strings"

	"github.com/pocketbase/pocketbase/models/schema"
)

func getExcludeSelects() map[string]bool {
	excludeSelectList := []string{}
	// exclude special filter literals
	excludeSelectList = append(excludeSelectList, "null", "true", "false")
	// exclude system literals
	excludeSelectList = append(excludeSelectList, schema.SystemFieldNames()...)

	excludeSelectMap := make(map[string]bool)
	for _, field := range excludeSelectList {
		excludeSelectMap[field] = true
	}
	return excludeSelectMap
}

// ParseSelectsFromString parses the provided string expression
// into a slice of strings.
//
// Example:
//
//	fields := search.ParseSelectsFromString("email,name")
func ParseSelectsFromString(str string) (fields []string) {
	excludeSelects := getExcludeSelects()
	data := strings.Split(str, ",")

	for _, field := range data {
		field = strings.TrimSpace(field)
		if _, ok := excludeSelects[field]; !ok {
			fields = append(fields, field)
		}
	}

	return
}
