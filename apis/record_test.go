package apis_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordsList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/collections/missing/records",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "unauthorized trying to access nil rule collection (aka. need admin auth)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo/records",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "public collection but with admin only filter/sort (aka. @collection)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo3/records?filter=@collection.demo.title='test'",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "public collection but with ENCODED admin only filter/sort (aka. @collection)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo3/records?filter=%40collection.demo.title%3D%27test%27",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"848a1dea-5ddd-42d6-a00d-030547bffcfe"`,
				`"id":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"id":"b5c2ffc2-bafd-48f7-b8b7-090638afe209"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "public collection",
			Method:         http.MethodGet,
			Url:            "/api/collections/demo3/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"2c542824-9de1-42fe-8924-e57c86267760"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "using the collection id as identifier",
			Method:         http.MethodGet,
			Url:            "/api/collections/3cd6fe92-70dc-4819-8542-4d036faabd89/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"2c542824-9de1-42fe-8924-e57c86267760"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "valid query params",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records?filter=title%7E%27test%27&sort=-title",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"848a1dea-5ddd-42d6-a00d-030547bffcfe"`,
				`"id":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "invalid filter",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records?filter=invalid~'test'",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expand",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records?expand=manyrels,onerel&perPage=2&sort=created",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":2`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"id":"848a1dea-5ddd-42d6-a00d-030547bffcfe"`,
				`"manyrels":[{`,
				`"manyrels":[]`,
				`"rel_cascade":"`,
				`"onerel":{"@collectionId":"3f2888f8-075d-49fe-9d09-ea7e951000dc","@collectionName":"demo",`,
				`"json":[1,2,3]`,
				`"select":["a","b"]`,
				`"select":[]`,
				`"user":""`,
				`"bool":true`,
				`"number":456`,
				`"user":"97cc3d3d-6ba2-383f-b42a-7bc84d27410c"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "authorized as user that DOESN'T match the collection list rule",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records",
			RequestHeaders: map[string]string{
				// test@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":0`,
				`"items":[]`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "authorized as user that matches the collection list rule",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records",
			RequestHeaders: map[string]string{
				// test3@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"63c2ab80-84ab-4057-a592-4604a731f78f"`,
				`"id":"94568ca2-0bee-49d7-b749-06cb97956fd9"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordView(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/collections/missing/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record (unauthorized)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo/records/00000000-bafd-48f7-b8b7-090638afe209",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "invalid record id (authorized)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records/invalid",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "missing record (authorized)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records/00000000-bafd-48f7-b8b7-090638afe209",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "mismatched collection-record pair (unauthorized)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "mismatched collection-record pair (authorized)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "unauthorized trying to access nil rule collection (aka. need admin auth)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "access record as admin",
			Method: http.MethodGet,
			Url:    "/api/collections/demo/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"@collectionId":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"@collectionName":"demo"`,
				`"id":"b5c2ffc2-bafd-48f7-b8b7-090638afe209"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "access record as admin (using the collection id as identifier)",
			Method: http.MethodGet,
			Url:    "/api/collections/3f2888f8-075d-49fe-9d09-ea7e951000dc/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"@collectionId":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"@collectionName":"demo"`,
				`"id":"b5c2ffc2-bafd-48f7-b8b7-090638afe209"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "access record as admin (test rule skipping)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records/94568ca2-0bee-49d7-b749-06cb97956fd9",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"@collectionId":"2c1010aa-b8fe-41d9-a980-99534ca8a167"`,
				`"@collectionName":"demo2"`,
				`"id":"94568ca2-0bee-49d7-b749-06cb97956fd9"`,
				`"manyrels":[]`,
				`"onerel":"b5c2ffc2-bafd-48f7-b8b7-090638afe209"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "access record as user (filter mismatch)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records/94568ca2-0bee-49d7-b749-06cb97956fd9",
			RequestHeaders: map[string]string{
				// test3@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "access record as user (filter match)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			RequestHeaders: map[string]string{
				// test3@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"@collectionId":"2c1010aa-b8fe-41d9-a980-99534ca8a167"`,
				`"@collectionName":"demo2"`,
				`"id":"63c2ab80-84ab-4057-a592-4604a731f78f"`,
				`"manyrels":["848a1dea-5ddd-42d6-a00d-030547bffcfe","577bd676-aacb-4072-b7da-99d00ee210a4"]`,
				`"onerel":"848a1dea-5ddd-42d6-a00d-030547bffcfe"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "expand relations",
			Method: http.MethodGet,
			Url:    "/api/collections/demo2/records/63c2ab80-84ab-4057-a592-4604a731f78f?expand=manyrels,onerel",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"@collectionId":"2c1010aa-b8fe-41d9-a980-99534ca8a167"`,
				`"@collectionName":"demo2"`,
				`"id":"63c2ab80-84ab-4057-a592-4604a731f78f"`,
				`"manyrels":[{`,
				`"onerel":{`,
				`"@collectionId":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"@collectionName":"demo"`,
				`"id":"848a1dea-5ddd-42d6-a00d-030547bffcfe"`,
				`"id":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordDelete(t *testing.T) {
	ensureDeletedFiles := func(app *tests.TestApp, collectionId string, recordId string) {
		storageDir := filepath.Join(app.DataDir(), "storage", collectionId, recordId)

		entries, _ := os.ReadDir(storageDir)
		if len(entries) != 0 {
			t.Errorf("Expected empty/deleted dir, found %d", len(entries))
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodDelete,
			Url:             "/api/collections/missing/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record (unauthorized)",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo/records/00000000-bafd-48f7-b8b7-090638afe209",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "missing record (authorized)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo/records/00000000-bafd-48f7-b8b7-090638afe209",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "mismatched collection-record pair (unauthorized)",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "mismatched collection-record pair (authorized)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "unauthorized trying to access nil rule collection (aka. need admin auth)",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo/records/577bd676-aacb-4072-b7da-99d00ee210a4",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo/records/577bd676-aacb-4072-b7da-99d00ee210a4",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "access record as admin",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo/records/577bd676-aacb-4072-b7da-99d00ee210a4",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
				"OnModelAfterUpdate":          1, // nullify related record
				"OnModelBeforeUpdate":         1, // nullify related record
				"OnModelBeforeDelete":         3, // +1 cascade delete related record
				"OnModelAfterDelete":          3, // +1 cascade delete related record
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				ensureDeletedFiles(app, "3f2888f8-075d-49fe-9d09-ea7e951000dc", "577bd676-aacb-4072-b7da-99d00ee210a4")
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "94568ca2-0bee-49d7-b749-06cb97956fd9")
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "63c2ab80-84ab-4057-a592-4604a731f78f")
			},
		},
		{
			Name:   "access record as admin (using the collection id as identifier)",
			Method: http.MethodDelete,
			Url:    "/api/collections/3f2888f8-075d-49fe-9d09-ea7e951000dc/records/577bd676-aacb-4072-b7da-99d00ee210a4",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
				"OnModelAfterUpdate":          1, // nullify related record
				"OnModelBeforeUpdate":         1, // nullify related record
				"OnModelBeforeDelete":         3, // +1 cascade delete related record
				"OnModelAfterDelete":          3, // +1 cascade delete related record
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				ensureDeletedFiles(app, "3f2888f8-075d-49fe-9d09-ea7e951000dc", "577bd676-aacb-4072-b7da-99d00ee210a4")
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "94568ca2-0bee-49d7-b749-06cb97956fd9")
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "63c2ab80-84ab-4057-a592-4604a731f78f")
			},
		},
		{
			Name:   "deleting record as admin (test rule skipping)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo2/records/94568ca2-0bee-49d7-b749-06cb97956fd9",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
				"OnModelBeforeDelete":         1,
				"OnModelAfterDelete":          1,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "94568ca2-0bee-49d7-b749-06cb97956fd9")
			},
		},
		{
			Name:   "deleting record as user (filter mismatch)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo2/records/94568ca2-0bee-49d7-b749-06cb97956fd9",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "deleting record as user (filter match)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo2/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
				"OnModelBeforeDelete":         1,
				"OnModelAfterDelete":          1,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "63c2ab80-84ab-4057-a592-4604a731f78f")
			},
		},
		{
			Name:   "trying to delete record while being part of a non-cascade required relation",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo/records/848a1dea-5ddd-42d6-a00d-030547bffcfe",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "cascade delete referenced records",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo/records/577bd676-aacb-4072-b7da-99d00ee210a4",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
				"OnModelBeforeDelete":         3,
				"OnModelAfterDelete":          3,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				recId := "63c2ab80-84ab-4057-a592-4604a731f78f"
				col, _ := app.Dao().FindCollectionByNameOrId("demo2")
				rec, _ := app.Dao().FindRecordById(col, recId, nil)
				if rec != nil {
					t.Errorf("Expected record %s to be cascade deleted", recId)
				}
				ensureDeletedFiles(app, "3f2888f8-075d-49fe-9d09-ea7e951000dc", "577bd676-aacb-4072-b7da-99d00ee210a4")
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "94568ca2-0bee-49d7-b749-06cb97956fd9")
				ensureDeletedFiles(app, "2c1010aa-b8fe-41d9-a980-99534ca8a167", "63c2ab80-84ab-4057-a592-4604a731f78f")
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordCreate(t *testing.T) {
	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title": "new",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodPost,
			Url:             "/api/collections/missing/records",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "guest trying to access nil-rule collection",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo/records",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "user trying to access nil-rule collection",
			Method: http.MethodPost,
			Url:    "/api/collections/demo/records",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit invalid format",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo3/records",
			Body:            strings.NewReader(`{"`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit nil body",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo3/records",
			Body:            nil,
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "guest submit in public collection",
			Method:         http.MethodPost,
			Url:            "/api/collections/demo3/records",
			Body:           strings.NewReader(`{"title":"new"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"title":"new"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "user submit in restricted collection (rule failure check)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo2/records",
			Body: strings.NewReader(`{
				"rel_cascade": "577bd676-aacb-4072-b7da-99d00ee210a4",
				"onerel": "577bd676-aacb-4072-b7da-99d00ee210a4",
				"manyrels": ["577bd676-aacb-4072-b7da-99d00ee210a4"],
				"text": "test123",
				"bool": "false",
				"number": 1
			}`),
			RequestHeaders: map[string]string{
				// test@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "user submit in restricted collection (rule pass check)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo2/records",
			Body: strings.NewReader(`{
				"rel_cascade":"577bd676-aacb-4072-b7da-99d00ee210a4",
				"onerel":"577bd676-aacb-4072-b7da-99d00ee210a4",
				"manyrels":["577bd676-aacb-4072-b7da-99d00ee210a4"],
				"text":"test123",
				"bool":true,
				"number":1
			}`),
			RequestHeaders: map[string]string{
				// test3@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"rel_cascade":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"onerel":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"manyrels":["577bd676-aacb-4072-b7da-99d00ee210a4"]`,
				`"text":"test123"`,
				`"bool":true`,
				`"number":1`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "admin submit in restricted collection (rule skip check)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo2/records",
			Body: strings.NewReader(`{
				"rel_cascade": "577bd676-aacb-4072-b7da-99d00ee210a4",
				"onerel": "577bd676-aacb-4072-b7da-99d00ee210a4",
				"manyrels" :["577bd676-aacb-4072-b7da-99d00ee210a4"],
				"text": "test123",
				"bool": false,
				"number": 1
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"rel_cascade":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"onerel":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"manyrels":["577bd676-aacb-4072-b7da-99d00ee210a4"]`,
				`"text":"test123"`,
				`"bool":false`,
				`"number":1`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "submit via multipart form data",
			Method: http.MethodPost,
			Url:    "/api/collections/demo/records",
			Body:   formData,
			RequestHeaders: map[string]string{
				"Content-Type":  mp.FormDataContentType(),
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"`,
				`"title":"new"`,
				`"file":"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordUpdate(t *testing.T) {
	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title": "new",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodPatch,
			Url:             "/api/collections/missing/records/2c542824-9de1-42fe-8924-e57c86267760",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo3/records/00000000-9de1-42fe-8924-e57c86267760",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "guest trying to edit nil-rule collection record",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "user trying to edit nil-rule collection record",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit invalid format",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo3/records/2c542824-9de1-42fe-8924-e57c86267760",
			Body:            strings.NewReader(`{"`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit nil body",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo3/records/2c542824-9de1-42fe-8924-e57c86267760",
			Body:            nil,
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "guest submit in public collection",
			Method:         http.MethodPatch,
			Url:            "/api/collections/demo3/records/2c542824-9de1-42fe-8924-e57c86267760",
			Body:           strings.NewReader(`{"title":"new"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"2c542824-9de1-42fe-8924-e57c86267760"`,
				`"title":"new"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "user submit in restricted collection (rule failure check)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2/records/94568ca2-0bee-49d7-b749-06cb97956fd9",
			Body:   strings.NewReader(`{"text": "test_new"}`),
			RequestHeaders: map[string]string{
				// test@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "user submit in restricted collection (rule pass check)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			Body: strings.NewReader(`{
				"text":"test_new",
				"bool":false
			}`),
			RequestHeaders: map[string]string{
				// test3@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"63c2ab80-84ab-4057-a592-4604a731f78f"`,
				`"rel_cascade":"577bd676-aacb-4072-b7da-99d00ee210a4"`,
				`"onerel":"848a1dea-5ddd-42d6-a00d-030547bffcfe"`,
				`"manyrels":["848a1dea-5ddd-42d6-a00d-030547bffcfe","577bd676-aacb-4072-b7da-99d00ee210a4"]`,
				`"bool":false`,
				`"text":"test_new"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "admin submit in restricted collection (rule skip check)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2/records/63c2ab80-84ab-4057-a592-4604a731f78f",
			Body: strings.NewReader(`{
				"text":"test_new",
				"number":1
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"63c2ab80-84ab-4057-a592-4604a731f78f"`,
				`"text":"test_new"`,
				`"number":1`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "submit via multipart form data",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo/records/b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			Body:   formData,
			RequestHeaders: map[string]string{
				"Content-Type":  mp.FormDataContentType(),
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"b5c2ffc2-bafd-48f7-b8b7-090638afe209"`,
				`"title":"new"`,
				`"file":"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
