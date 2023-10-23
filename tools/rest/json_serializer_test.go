package rest_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/search"
)

func TestSerialize(t *testing.T) {
	scenarios := []struct {
		name       string
		serializer rest.Serializer
		statusCode int
		data       any
		query      string
		expected   string
	}{
		{
			"empty query",
			rest.Serializer{},
			200,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"",
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"empty fields",
			rest.Serializer{},
			200,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=",
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"missing fields",
			rest.Serializer{},
			200,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=missing",
			`{}`,
		},
		{
			">299 response",
			rest.Serializer{},
			300,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=missing",
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"<200 response",
			rest.Serializer{},
			199,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=missing",
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"non map response",
			rest.Serializer{},
			200,
			"test",
			"fields=a,b,test",
			`"test"`,
		},
		{
			"non slice of map response",
			rest.Serializer{},
			200,
			[]any{"a", "b", "test"},
			"fields=a,test",
			`["a","b","test"]`,
		},
		{
			"map with no matching field",
			rest.Serializer{},
			200,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=missing", // test individual fields trim
			`{}`,
		},
		{
			"map with existing and missing fields",
			rest.Serializer{},
			200,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=a,  c  ,missing", // test individual fields trim
			`{"a":1,"c":"test"}`,
		},
		{
			"custom fields param",
			rest.Serializer{FieldsParam: "custom"},
			200,
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"custom=a,  c  ,missing", // test individual fields trim
			`{"a":1,"c":"test"}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.URL.RawQuery = s.query
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.Response().Status = s.statusCode

			if err := s.serializer.Serialize(c, s.data, ""); err != nil {
				t.Fatalf("Serialize failure: %v", err)
			}

			rawBody, err := io.ReadAll(rec.Result().Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			if v := strings.TrimSpace(string(rawBody)); v != s.expected {
				t.Fatalf("Expected body\n%v \ngot \n%v", s.expected, v)
			}
		})
	}
}

func TestPickFields(t *testing.T) {
	scenarios := []struct {
		name        string
		data        any
		fields      string
		expectError bool
		result      string
	}{
		{
			"empty fields",
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"",
			false,
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"missing fields",
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"missing",
			false,
			`{}`,
		},
		{
			"non map data",
			"test",
			"a,b,test",
			false,
			`"test"`,
		},
		{
			"non slice of map data",
			[]any{"a", "b", "test"},
			"a,test",
			false,
			`["a","b","test"]`,
		},
		{
			"map with no matching field",
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"missing", // test individual fields trim
			false,
			`{}`,
		},
		{
			"map with existing and missing fields",
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"a,  c  ,missing", // test individual fields trim
			false,
			`{"a":1,"c":"test"}`,
		},
		{
			"slice of maps with existing and missing fields",
			[]any{
				map[string]any{"a": 11, "b": 11, "c": "test1"},
				map[string]any{"a": 22, "b": 22, "c": "test2"},
			},
			"a,  c  ,missing", // test individual fields trim
			false,
			`[{"a":11,"c":"test1"},{"a":22,"c":"test2"}]`,
		},
		{
			"nested fields with mixed map and any slices",
			map[string]any{
				"a": 1,
				"b": 2,
				"c": "test",
				"anySlice": []any{
					map[string]any{
						"A": []int{1, 2, 3},
						"B": []any{"1", "2", 3},
						"C": "test",
						"D": map[string]any{
							"DA": 1,
							"DB": 2,
						},
					},
					map[string]any{
						"A": "test",
					},
				},
				"mapSlice": []map[string]any{
					{
						"A": []int{1, 2, 3},
						"B": []any{"1", "2", 3},
						"C": "test",
						"D": []any{
							map[string]any{"DA": 1},
						},
					},
					{
						"B": []any{"1", "2", 3},
						"D": []any{
							map[string]any{"DA": 2},
							map[string]any{"DA": 3},
							map[string]any{"DB": 4}, // will result to empty since there is no DA
						},
					},
				},
				"fullMap": []map[string]any{
					{
						"A": []int{1, 2, 3},
						"B": []any{"1", "2", 3},
						"C": "test",
					},
					{
						"B": []any{"1", "2", 3},
						"D": []any{
							map[string]any{"DA": 2},
							map[string]any{"DA": 3}, // will result to empty since there is no DA
						},
					},
				},
			},
			"a, c, anySlice.A, mapSlice.C, mapSlice.D.DA, anySlice.D,fullMap",
			false,
			`{"a":1,"anySlice":[{"A":[1,2,3],"D":{"DA":1,"DB":2}},{"A":"test"}],"c":"test","fullMap":[{"A":[1,2,3],"B":["1","2",3],"C":"test"},{"B":["1","2",3],"D":[{"DA":2},{"DA":3}]}],"mapSlice":[{"C":"test","D":[{"DA":1}]},{"D":[{"DA":2},{"DA":3},{}]}]}`,
		},
		{
			"SearchResult",
			search.Result{
				Page:       1,
				PerPage:    10,
				TotalItems: 20,
				TotalPages: 30,
				Items: []any{
					map[string]any{"a": 11, "b": 11, "c": "test1"},
					map[string]any{"a": 22, "b": 22, "c": "test2"},
				},
			},
			"a,c,missing",
			false,
			`{"items":[{"a":11,"c":"test1"},{"a":22,"c":"test2"}],"page":1,"perPage":10,"totalItems":20,"totalPages":30}`,
		},
		{
			"*SearchResult",
			&search.Result{
				Page:       1,
				PerPage:    10,
				TotalItems: 20,
				TotalPages: 30,
				Items: []any{
					map[string]any{"a": 11, "b": 11, "c": "test1"},
					map[string]any{"a": 22, "b": 22, "c": "test2"},
				},
			},
			"a,c",
			false,
			`{"items":[{"a":11,"c":"test1"},{"a":22,"c":"test2"}],"page":1,"perPage":10,"totalItems":20,"totalPages":30}`,
		},
		{
			"root wildcard",
			&search.Result{
				Page:       1,
				PerPage:    10,
				TotalItems: 20,
				TotalPages: 30,
				Items: []any{
					map[string]any{"a": 11, "b": 11, "c": "test1"},
					map[string]any{"a": 22, "b": 22, "c": "test2"},
				},
			},
			"*",
			false,
			`{"items":[{"a":11,"b":11,"c":"test1"},{"a":22,"b":22,"c":"test2"}],"page":1,"perPage":10,"totalItems":20,"totalPages":30}`,
		},
		{
			"root wildcard with nested exception",
			map[string]any{
				"id":    "123",
				"title": "lorem",
				"rel": map[string]any{
					"id":    "456",
					"title": "rel_title",
				},
			},
			"*,rel.id",
			false,
			`{"id":"123","rel":{"id":"456"},"title":"lorem"}`,
		},
		{
			"sub wildcard",
			map[string]any{
				"id":    "123",
				"title": "lorem",
				"rel": map[string]any{
					"id":    "456",
					"title": "rel_title",
					"sub": map[string]any{
						"id":    "789",
						"title": "sub_title",
					},
				},
			},
			"id,rel.*",
			false,
			`{"id":"123","rel":{"id":"456","sub":{"id":"789","title":"sub_title"},"title":"rel_title"}}`,
		},
		{
			"sub wildcard with nested exception",
			map[string]any{
				"id":    "123",
				"title": "lorem",
				"rel": map[string]any{
					"id":    "456",
					"title": "rel_title",
					"sub": map[string]any{
						"id":    "789",
						"title": "sub_title",
					},
				},
			},
			"id,rel.*,rel.sub.id",
			false,
			`{"id":"123","rel":{"id":"456","sub":{"id":"789"},"title":"rel_title"}}`,
		},
		{
			"invalid excerpt modifier",
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"*:excerpt",
			true,
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"valid excerpt modifier",
			map[string]any{
				"id":    "123",
				"title": "lorem",
				"rel": map[string]any{
					"id":    "456",
					"title": "<p>rel_title</p>",
					"sub": map[string]any{
						"id":    "789",
						"title": "sub_title",
					},
				},
			},
			"*:excerpt(2),rel.title:excerpt(3, true)",
			false,
			`{"id":"12","rel":{"title":"rel..."},"title":"lo"}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result, err := rest.PickFields(s.data, s.fields)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			serialized, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			if v := string(serialized); v != s.result {
				t.Fatalf("Expected body\n%s \ngot \n%s", s.result, v)
			}
		})
	}
}
