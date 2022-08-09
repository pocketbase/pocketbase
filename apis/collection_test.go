package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestCollectionsList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/collections",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/collections",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/collections",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":5`,
				`"items":[{`,
				`"id":"abe78266-fd4d-4aea-962d-8c0138ac522b"`,
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"id":"2c1010aa-b8fe-41d9-a980-99534ca8a167"`,
				`"id":"3cd6fe92-70dc-4819-8542-4d036faabd89"`,
				`"id":"f12f3eb6-b980-4bf6-b1e4-36de0450c8be"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + paging and sorting",
			Method: http.MethodGet,
			Url:    "/api/collections?page=2&perPage=2&sort=-created",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":2`,
				`"perPage":2`,
				`"totalItems":5`,
				`"items":[{`,
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"id":"2c1010aa-b8fe-41d9-a980-99534ca8a167"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + invalid filter",
			Method: http.MethodGet,
			Url:    "/api/collections?filter=invalidfield~'demo2'",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + valid filter",
			Method: http.MethodGet,
			Url:    "/api/collections?filter=name~'demo2'",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"2c1010aa-b8fe-41d9-a980-99534ca8a167"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionView(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/collections/demo",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting collection identifier",
			Method: http.MethodGet,
			Url:    "/api/collections/missing",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + using the collection name",
			Method: http.MethodGet,
			Url:    "/api/collections/demo",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionViewRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + using the collection id",
			Method: http.MethodGet,
			Url:    "/api/collections/3f2888f8-075d-49fe-9d09-ea7e951000dc",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionViewRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionDelete(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo3",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo3",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting collection identifier",
			Method: http.MethodDelete,
			Url:    "/api/collections/b97ccf83-34a2-4d01-a26b-3d77bc842d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + using the collection name",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo3",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":             1,
				"OnModelAfterDelete":              1,
				"OnCollectionBeforeDeleteRequest": 1,
				"OnCollectionAfterDeleteRequest":  1,
			},
		},
		{
			Name:   "authorized as admin + using the collection id",
			Method: http.MethodDelete,
			Url:    "/api/collections/3cd6fe92-70dc-4819-8542-4d036faabd89",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":             1,
				"OnModelAfterDelete":              1,
				"OnCollectionBeforeDeleteRequest": 1,
				"OnCollectionAfterDeleteRequest":  1,
			},
		},
		{
			Name:   "authorized as admin + trying to delete a system collection",
			Method: http.MethodDelete,
			Url:    "/api/collections/profiles",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnCollectionBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + trying to delete a referenced collection",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnCollectionBeforeDeleteRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionCreate(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/collections",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/collections",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_required"`,
				`"schema":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin + invalid data (eg. existing name)",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(`{"name":"demo","schema":[{"type":"text","name":""}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_collection_name_exists"`,
				`"schema":{"0":{"name":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","schema":[{"type":"text","id":"12345789","name":"test"}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
				`"system":false`,
				`"schema":[{"system":false,"id":"12345789","name":"test","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}]`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionUpdate(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "authorized as admin + invalid data (eg. existing name)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo",
			Body:   strings.NewReader(`{"name":"demo2"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_collection_name_exists"`,
			},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo",
			Body:   strings.NewReader(`{"name":"new"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"name":"new"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "authorized as admin + valid data and id as identifier",
			Method: http.MethodPatch,
			Url:    "/api/collections/3f2888f8-075d-49fe-9d09-ea7e951000dc",
			Body:   strings.NewReader(`{"name":"new"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc"`,
				`"name":"new"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionImport(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPut,
			Url:             "/api/collections/import",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty collections",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body:   strings.NewReader(`{"collections":[]}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"validation_required"`,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				if len(collections) != 5 {
					t.Fatalf("Expected %d collections, got %d", 5, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + trying to delete system collections",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body:   strings.NewReader(`{"deleteMissing": true, "collections":[{"name": "test123"}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"collections_import_failure"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				if len(collections) != 5 {
					t.Fatalf("Expected %d collections, got %d", 5, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + collections validator failure",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{"name": "import2"}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"collections_import_failure"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
				"OnModelBeforeCreate":              2,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				if len(collections) != 5 {
					t.Fatalf("Expected %d collections, got %d", 5, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + successful collections save",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{
						"name": "import2",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
				"OnCollectionsAfterImportRequest":  1,
				"OnModelBeforeCreate":              2,
				"OnModelAfterCreate":               2,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				if len(collections) != 7 {
					t.Fatalf("Expected %d collections, got %d", 7, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + successful collections save and old non-system collections deletion",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"deleteMissing": true,
				"collections":[
					{
						"id":"abe78266-fd4d-4aea-962d-8c0138ac522b",
						"name":"profiles",
						"system":true,
						"listRule":"userId = @request.user.id",
						"viewRule":"created > 'test_change'",
						"createRule":"userId = @request.user.id",
						"updateRule":"userId = @request.user.id",
						"deleteRule":"userId = @request.user.id",
						"schema":[
							{
								"id":"koih1lqx",
								"name":"userId",
								"type":"user",
								"system":true,
								"required":true,
								"unique":true,
								"options":{
									"maxSelect":1,
									"cascadeDelete":true
								}
							},
							{
								"id":"69ycbg3q",
								"name":"rel",
								"type":"relation",
								"system":false,
								"required":false,
								"unique":false,
								"options":{
									"maxSelect":2,
									"collectionId":"abe78266-fd4d-4aea-962d-8c0138ac522b",
									"cascadeDelete":false
								}
							}
						]
					},
					{
						"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
						"name":"demo",
						"schema":[
							{
								"id":"_2hlxbmp",
								"name":"title",
								"type":"text",
								"system":false,
								"required":true,
								"unique":false,
								"options":{
									"min":3,
									"max":null,
									"pattern":""
								}
							}
						]
					},
					{
						"name": "new_import",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnCollectionsAfterImportRequest":  1,
				"OnCollectionsBeforeImportRequest": 1,
				"OnModelBeforeDelete":              3,
				"OnModelAfterDelete":               3,
				"OnModelBeforeUpdate":              2,
				"OnModelAfterUpdate":               2,
				"OnModelBeforeCreate":              1,
				"OnModelAfterCreate":               1,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				if len(collections) != 3 {
					t.Fatalf("Expected %d collections, got %d", 3, len(collections))
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
