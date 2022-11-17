package models

// FilterRequestData defines a HTTP request data struct, usually used
// as part of the `@request.*` filter resolver.
type FilterRequestData struct {
	Method     string         `json:"method"`
	Query      map[string]any `json:"query"`
	Data       map[string]any `json:"data"`
	AuthRecord *Record        `json:"authRecord"`
	Admin      *Admin         `json:"admin"`
}
