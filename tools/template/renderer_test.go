package template

import (
	"errors"
	"html/template"
	"testing"
)

func TestRendererRender(t *testing.T) {
	tpl, _ := template.New("").Parse("Hello {{.Name}}!")
	tpl.Option("missingkey=error") // enforce execute errors

	scenarios := map[string]struct {
		renderer       *Renderer
		data           any
		expectedHasErr bool
		expectedResult string
	}{
		"with nil template": {
			&Renderer{},
			nil,
			true,
			"",
		},
		"with parse error": {
			&Renderer{
				template:   tpl,
				parseError: errors.New("test"),
			},
			nil,
			true,
			"",
		},
		"with execute error": {
			&Renderer{template: tpl},
			nil,
			true,
			"",
		},
		"no error": {
			&Renderer{template: tpl},
			struct{ Name string }{"world"},
			false,
			"Hello world!",
		},
	}

	for name, s := range scenarios {
		t.Run(name, func(t *testing.T) {
			result, err := s.renderer.Render(s.data)

			hasErr := err != nil

			if s.expectedHasErr != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectedHasErr, hasErr, err)
			}

			if s.expectedResult != result {
				t.Fatalf("Expected result %v, got %v", s.expectedResult, result)
			}
		})
	}
}
