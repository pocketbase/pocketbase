package rest_test

import (
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
		data       any
		query      string
		expected   string
	}{
		{
			"empty query",
			rest.Serializer{},
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"",
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"empty fields",
			rest.Serializer{},
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=",
			`{"a":1,"b":2,"c":"test"}`,
		},
		{
			"missing fields",
			rest.Serializer{},
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=missing",
			`{}`,
		},
		{
			"non map response",
			rest.Serializer{},
			"test",
			"fields=a,b,test",
			`"test"`,
		},
		{
			"non slice of map response",
			rest.Serializer{},
			[]any{"a", "b", "test"},
			"fields=a,test",
			`["a","b","test"]`,
		},
		{
			"map with existing and missing fields",
			rest.Serializer{},
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"fields=a,  c  ,missing", // test individual fields trim
			`{"a":1,"c":"test"}`,
		},
		{
			"custom fields param",
			rest.Serializer{FieldsParam: "custom"},
			map[string]any{"a": 1, "b": 2, "c": "test"},
			"custom=a,  c  ,missing", // test individual fields trim
			`{"a":1,"c":"test"}`,
		},
		{
			"slice of maps with existing and missing fields",
			rest.Serializer{},
			[]any{
				map[string]any{"a": 11, "b": 11, "c": "test1"},
				map[string]any{"a": 22, "b": 22, "c": "test2"},
			},
			"fields=a,  c  ,missing", // test individual fields trim
			`[{"a":11,"c":"test1"},{"a":22,"c":"test2"}]`,
		},
		{
			"nested fields with mixed map and any slices",
			rest.Serializer{},
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
			"fields=a, c, anySlice.A, mapSlice.C, mapSlice.D.DA, anySlice.D,fullMap",
			`{"a":1,"anySlice":[{"A":[1,2,3],"D":{"DA":1,"DB":2}},{"A":"test"}],"c":"test","fullMap":[{"A":[1,2,3],"B":["1","2",3],"C":"test"},{"B":["1","2",3],"D":[{"DA":2},{"DA":3}]}],"mapSlice":[{"C":"test","D":[{"DA":1}]},{"D":[{"DA":2},{"DA":3},{}]}]}`,
		},
		{
			"SearchResult",
			rest.Serializer{},
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
			"fields=a,c,missing",
			`{"items":[{"a":11,"c":"test1"},{"a":22,"c":"test2"}],"page":1,"perPage":10,"totalItems":20,"totalPages":30}`,
		},
		{
			"*SearchResult",
			rest.Serializer{},
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
			"fields=a,c",
			`{"items":[{"a":11,"c":"test1"},{"a":22,"c":"test2"}],"page":1,"perPage":10,"totalItems":20,"totalPages":30}`,
		},
	}

	for _, s := range scenarios {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.URL.RawQuery = s.query
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := s.serializer.Serialize(c, s.data, ""); err != nil {
			t.Errorf("[%s] Serialize failure: %v", s.name, err)
			continue
		}

		rawBody, err := io.ReadAll(rec.Result().Body)
		if err != nil {
			t.Errorf("[%s] Failed to read request body: %v", s.name, err)
			continue
		}

		if v := strings.TrimSpace(string(rawBody)); v != s.expected {
			t.Fatalf("[%s] Expected body\n%v \ngot: \n%v", s.name, s.expected, v)
		}
	}
}
