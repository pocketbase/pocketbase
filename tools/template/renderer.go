package template

import (
	"bytes"
	"errors"
	"html/template"
)

// Renderer defines a single parsed template.
type Renderer struct {
	template   *template.Template
	parseError error
}

// Render executes the template with the specified data as the dot object
// and returns the result as plain string.
func (r *Renderer) Render(data any) (string, error) {
	if r.parseError != nil {
		return "", r.parseError
	}

	if r.template == nil {
		return "", errors.New("invalid or nil template")
	}

	buf := new(bytes.Buffer)

	if err := r.template.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
