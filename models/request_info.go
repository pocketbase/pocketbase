package models

import (
	"strings"

	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	RequestInfoContextDefault       = "default"
	RequestInfoContextRealtime      = "realtime"
	RequestInfoContextProtectedFile = "protectedFile"
	RequestInfoContextOAuth2        = "oauth2"
)

// RequestInfo defines a HTTP request data struct, usually used
// as part of the `@request.*` filter resolver.
type RequestInfo struct {
	Context    string         `json:"context"`
	Query      map[string]any `json:"query"`
	Data       map[string]any `json:"data"`
	Headers    map[string]any `json:"headers"`
	AuthRecord *Record        `json:"authRecord"`
	Admin      *Admin         `json:"admin"`
	Method     string         `json:"method"`
}

// HasModifierDataKeys loosely checks if the current struct has any modifier Data keys.
func (r *RequestInfo) HasModifierDataKeys() bool {
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
