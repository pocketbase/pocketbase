package models

import (
	"strings"

	"github.com/pocketbase/pocketbase/models/schema"
)

// RequestData defines a HTTP request data struct, usually used
// as part of the `@request.*` filter resolver.
type RequestData struct {
	Method     string         `json:"method"`
	Query      map[string]any `json:"query"`
	Data       map[string]any `json:"data"`
	Headers    map[string]any `json:"headers"`
	AuthRecord *Record        `json:"authRecord"`
	Admin      *Admin         `json:"admin"`
}

// HasModifierDataKeys loosely checks if the current struct has any modifier Data keys.
func (r *RequestData) HasModifierDataKeys() bool {
	allModifiers := schema.FieldValueModifiers()

	for key := range r.Data {
		for _, m := range allModifiers {
			if strings.HasSuffix(key, m) {
				return true
			}
		}
	}

	return false
}
