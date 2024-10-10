package picker_test

import (
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/tools/picker"
	"github.com/pocketbase/pocketbase/tools/search"
)

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
			result, err := picker.Pick(s.data, s.fields)

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
