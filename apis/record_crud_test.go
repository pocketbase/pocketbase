package apis_test

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordCrudList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/collections/missing/records",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "unauthenticated trying to access nil rule collection (aka. need admin auth)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo1/records",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authenticated record trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "public collection but with admin only filter param (aka. @collection, @request, etc.)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo2/records?filter=%40collection.demo2.title='test1'",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "public collection but with admin only sort param (aka. @collection, @request, etc.)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo2/records?sort=@request.auth.title",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "public collection but with ENCODED admin only filter/sort (aka. @collection)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo2/records?filter=%40collection.demo2.title%3D%27test1%27",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "public collection",
			Method:         http.MethodGet,
			Url:            "/api/collections/demo2/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"0yxhwia2amd8gec"`,
				`"id":"achvryl401bhse3"`,
				`"id":"llvuca81nly1qls"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "public collection (using the collection id)",
			Method:         http.MethodGet,
			Url:            "/api/collections/sz5l5z67tg7gku0/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"0yxhwia2amd8gec"`,
				`"id":"achvryl401bhse3"`,
				`"id":"llvuca81nly1qls"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "authorized as admin trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"al1h9ijdeojtsjy"`,
				`"id":"84nmscqy84lsi1t"`,
				`"id":"imy661ixudk5izi"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "valid query params",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records?filter=text~'test'&sort=-bool",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"al1h9ijdeojtsjy"`,
				`"id":"84nmscqy84lsi1t"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "invalid filter",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records?filter=invalid~'test'",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expand relations",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records?expand=rel_one,rel_many.rel,missing&perPage=2&sort=created",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":2`,
				`"totalPages":2`,
				`"totalItems":3`,
				`"items":[{`,
				`"collectionName":"demo1"`,
				`"id":"84nmscqy84lsi1t"`,
				`"id":"al1h9ijdeojtsjy"`,
				`"expand":{`,
				`"rel_one":""`,
				`"rel_one":{"`,
				`"rel_many":[{`,
				`"rel":{`,
				`"rel":""`,
				`"json":[1,2,3]`,
				`"select_many":["optionB","optionC"]`,
				`"select_many":["optionB"]`,
				// subrel items
				`"id":"0yxhwia2amd8gec"`,
				`"id":"llvuca81nly1qls"`,
				// email visibility should be ignored for admins even in expanded rels
				`"email":"test@example.com"`,
				`"email":"test2@example.com"`,
				`"email":"test3@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "authenticated record model that DOESN'T match the collection list rule",
			Method: http.MethodGet,
			Url:    "/api/collections/demo3/records",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
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
			Name:   "authenticated record that matches the collection list rule",
			Method: http.MethodGet,
			Url:    "/api/collections/demo3/records",
			RequestHeaders: map[string]string{
				// clients, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.q34IWXrRWsjLvbbVNRfAs_J4SoTHloNBfdGEiLmy-D8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":4`,
				`"items":[{`,
				`"id":"1tmknxy2868d869"`,
				`"id":"lcl9d87w22ml6jy"`,
				`"id":"7nwo8tuiatetxdm"`,
				`"id":"mk5fmymtx4wsprk"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           ":rule modifer",
			Method:         http.MethodGet,
			Url:            "/api/collections/demo5/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"qjeql998mtp1azp"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "multi-match - at least one of",
			Method:         http.MethodGet,
			Url:            "/api/collections/demo4/records?filter=" + url.QueryEscape("rel_many_no_cascade_required.files:length?=2"),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"qzaqccwrmva4o1n"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "multi-match - all",
			Method:         http.MethodGet,
			Url:            "/api/collections/demo4/records?filter=" + url.QueryEscape("rel_many_no_cascade_required.files:length=2"),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":0`,
				`"totalItems":0`,
				`"items":[]`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},

		// auth collection
		// -----------------------------------------------------------
		{
			Name:           "check email visibility as guest",
			Method:         http.MethodGet,
			Url:            "/api/collections/nologin/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"phhq3wr65cap535"`,
				`"id":"dc49k6jgejn40h3"`,
				`"id":"oos036e9xvqeexy"`,
				`"email":"test2@example.com"`,
				`"emailVisibility":true`,
				`"emailVisibility":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
				`"email":"test@example.com"`,
				`"email":"test3@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "check email visibility as any authenticated record",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records",
			RequestHeaders: map[string]string{
				// clients, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.q34IWXrRWsjLvbbVNRfAs_J4SoTHloNBfdGEiLmy-D8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"phhq3wr65cap535"`,
				`"id":"dc49k6jgejn40h3"`,
				`"id":"oos036e9xvqeexy"`,
				`"email":"test2@example.com"`,
				`"emailVisibility":true`,
				`"emailVisibility":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
				`"email":"test@example.com"`,
				`"email":"test3@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "check email visibility as manage auth record",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"phhq3wr65cap535"`,
				`"id":"dc49k6jgejn40h3"`,
				`"id":"oos036e9xvqeexy"`,
				`"email":"test@example.com"`,
				`"email":"test2@example.com"`,
				`"email":"test3@example.com"`,
				`"emailVisibility":true`,
				`"emailVisibility":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "check email visibility as admin",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"phhq3wr65cap535"`,
				`"id":"dc49k6jgejn40h3"`,
				`"id":"oos036e9xvqeexy"`,
				`"email":"test@example.com"`,
				`"email":"test2@example.com"`,
				`"email":"test3@example.com"`,
				`"emailVisibility":true`,
				`"emailVisibility":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "check self email visibility resolver",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records",
			RequestHeaders: map[string]string{
				// nologin, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImRjNDlrNmpnZWpuNDBoMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoia3B2NzA5c2sybHFicWs4IiwiZXhwIjoyMjA4OTg1MjYxfQ.DOYSon3x1-C0hJbwjEU6dp2-6oLeEa8bOlkyP1CinyM",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"phhq3wr65cap535"`,
				`"id":"dc49k6jgejn40h3"`,
				`"id":"oos036e9xvqeexy"`,
				`"email":"test2@example.com"`,
				`"email":"test@example.com"`,
				`"emailVisibility":true`,
				`"emailVisibility":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
				`"email":"test3@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},

		// view collection
		// -----------------------------------------------------------
		{
			Name:           "public view records",
			Method:         http.MethodGet,
			Url:            "/api/collections/view2/records?filter=state=false",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"al1h9ijdeojtsjy"`,
				`"id":"imy661ixudk5izi"`,
			},
			NotExpectedContent: []string{
				`"created"`,
				`"updated"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "guest that doesn't match the view collection list rule",
			Method:         http.MethodGet,
			Url:            "/api/collections/view1/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":0`,
				`"totalItems":0`,
				`"items":[]`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:   "authenticated record that matches the view collection list rule",
			Method: http.MethodGet,
			Url:    "/api/collections/view1/records",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"84nmscqy84lsi1t"`,
				`"bool":true`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
		{
			Name:           "view collection with numeric ids",
			Method:         http.MethodGet,
			Url:            "/api/collections/numeric_id_view/records",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalPages":1`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"1"`,
				`"id":"2"`,
			},
			ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordCrudView(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/collections/missing/records/0yxhwia2amd8gec",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo2/records/missing",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "unauthenticated trying to access nil rule collection (aka. need admin auth)",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo1/records/imy661ixudk5izi",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authenticated record trying to access nil rule collection (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records/imy661ixudk5izi",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authenticated record that doesn't match the collection view rule",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/bgs820n361vj1qd",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "public collection view",
			Method:         http.MethodGet,
			Url:            "/api/collections/demo2/records/0yxhwia2amd8gec",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"0yxhwia2amd8gec"`,
				`"collectionName":"demo2"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:           "public collection view (using the collection id)",
			Method:         http.MethodGet,
			Url:            "/api/collections/sz5l5z67tg7gku0/records/0yxhwia2amd8gec",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"0yxhwia2amd8gec"`,
				`"collectionName":"demo2"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "authorized as admin trying to access nil rule collection view (aka. need admin auth)",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records/imy661ixudk5izi",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"imy661ixudk5izi"`,
				`"collectionName":"demo1"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "authenticated record that does match the collection view rule",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"4q1xlclmfloku33"`,
				`"collectionName":"users"`,
				// owners can always view their email
				`"emailVisibility":false`,
				`"email":"test@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "expand relations",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1/records/al1h9ijdeojtsjy?expand=rel_one,rel_many.rel,missing&perPage=2&sort=created",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"al1h9ijdeojtsjy"`,
				`"collectionName":"demo1"`,
				`"rel_many":[{`,
				`"rel_one":{`,
				`"collectionName":"users"`,
				`"id":"bgs820n361vj1qd"`,
				`"expand":{"rel":{`,
				`"id":"0yxhwia2amd8gec"`,
				`"collectionName":"demo2"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},

		// auth collection
		// -----------------------------------------------------------
		{
			Name:           "check email visibility as guest",
			Method:         http.MethodGet,
			Url:            "/api/collections/nologin/records/oos036e9xvqeexy",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"oos036e9xvqeexy"`,
				`"emailVisibility":false`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
				`"email":"test3@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "check email visibility as any authenticated record",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records/oos036e9xvqeexy",
			RequestHeaders: map[string]string{
				// clients, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.q34IWXrRWsjLvbbVNRfAs_J4SoTHloNBfdGEiLmy-D8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"oos036e9xvqeexy"`,
				`"emailVisibility":false`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
				`"email":"test3@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "check email visibility as manage auth record",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records/oos036e9xvqeexy",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"oos036e9xvqeexy"`,
				`"emailVisibility":false`,
				`"email":"test3@example.com"`,
				`"verified":true`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "check email visibility as admin",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records/oos036e9xvqeexy",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"oos036e9xvqeexy"`,
				`"emailVisibility":false`,
				`"email":"test3@example.com"`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:   "check self email visibility resolver",
			Method: http.MethodGet,
			Url:    "/api/collections/nologin/records/dc49k6jgejn40h3",
			RequestHeaders: map[string]string{
				// nologin, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImRjNDlrNmpnZWpuNDBoMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoia3B2NzA5c2sybHFicWs4IiwiZXhwIjoyMjA4OTg1MjYxfQ.DOYSon3x1-C0hJbwjEU6dp2-6oLeEa8bOlkyP1CinyM",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"dc49k6jgejn40h3"`,
				`"email":"test@example.com"`,
				`"emailVisibility":false`,
				`"verified":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},

		// view collection
		// -----------------------------------------------------------
		{
			Name:           "public view record",
			Method:         http.MethodGet,
			Url:            "/api/collections/view2/records/84nmscqy84lsi1t",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"84nmscqy84lsi1t"`,
				`"state":true`,
				`"file_many":["`,
				`"rel_many":["`,
			},
			NotExpectedContent: []string{
				`"created"`,
				`"updated"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:            "guest that doesn't match the view collection view rule",
			Method:          http.MethodGet,
			Url:             "/api/collections/view1/records/84nmscqy84lsi1t",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authenticated record that matches the view collection view rule",
			Method: http.MethodGet,
			Url:    "/api/collections/view1/records/84nmscqy84lsi1t",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"84nmscqy84lsi1t"`,
				`"bool":true`,
				`"text":"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
		{
			Name:           "view record with numeric id",
			Method:         http.MethodGet,
			Url:            "/api/collections/numeric_id_view/records/1",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"1"`,
			},
			ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordCrudDelete(t *testing.T) {
	t.Parallel()

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
			Url:             "/api/collections/missing/records/0yxhwia2amd8gec",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo2/records/missing",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "unauthenticated trying to delete nil rule collection (aka. need admin auth)",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo1/records/imy661ixudk5izi",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authenticated record trying to delete nil rule collection (aka. need admin auth)",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo1/records/imy661ixudk5izi",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authenticated record that doesn't match the collection delete rule",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/bgs820n361vj1qd",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "trying to delete a view collection record",
			Method:          http.MethodDelete,
			Url:             "/api/collections/view1/records/imy661ixudk5izi",
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "public collection record delete",
			Method:         http.MethodDelete,
			Url:            "/api/collections/nologin/records/dc49k6jgejn40h3",
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":          1,
				"OnModelBeforeDelete":         1,
				"OnRecordAfterDeleteRequest":  1,
				"OnRecordBeforeDeleteRequest": 1,
			},
		},
		{
			Name:           "public collection record delete (using the collection id as identifier)",
			Method:         http.MethodDelete,
			Url:            "/api/collections/kpv709sk2lqbqk8/records/dc49k6jgejn40h3",
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":          1,
				"OnModelBeforeDelete":         1,
				"OnRecordAfterDeleteRequest":  1,
				"OnRecordBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "authorized as admin trying to delete nil rule collection view (aka. need admin auth)",
			Method: http.MethodDelete,
			Url:    "/api/collections/clients/records/o1y0dd0spd786md",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":          1,
				"OnModelBeforeDelete":         1,
				"OnRecordAfterDeleteRequest":  1,
				"OnRecordBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "OnRecordAfterDeleteRequest error response",
			Method: http.MethodDelete,
			Url:    "/api/collections/clients/records/o1y0dd0spd786md",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":          1,
				"OnModelBeforeDelete":         1,
				"OnRecordAfterDeleteRequest":  1,
				"OnRecordBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "authenticated record that match the collection delete rule",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":          3, // +2 because of the external auths
				"OnModelBeforeDelete":         3, // +2 because of the external auths
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterDeleteRequest":  1,
				"OnRecordBeforeDeleteRequest": 1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				ensureDeletedFiles(app, "_pb_users_auth_", "4q1xlclmfloku33")

				// check if all the external auths records were deleted
				collection, _ := app.Dao().FindCollectionByNameOrId("users")
				record := models.NewRecord(collection)
				record.Id = "4q1xlclmfloku33"
				externalAuths, err := app.Dao().FindAllExternalAuthsByRecord(record)
				if err != nil {
					t.Errorf("Failed to fetch external auths: %v", err)
				}
				if len(externalAuths) > 0 {
					t.Errorf("Expected the linked external auths to be deleted, got %d", len(externalAuths))
				}
			},
		},
		{
			Name:            "@request :isset (rule failure check)",
			Method:          http.MethodDelete,
			Url:             "/api/collections/demo5/records/la4y2w4o98acwuj",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "@request :isset (rule pass check)",
			Method:         http.MethodDelete,
			Url:            "/api/collections/demo5/records/la4y2w4o98acwuj?test=1",
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":          1,
				"OnModelBeforeDelete":         1,
				"OnRecordAfterDeleteRequest":  1,
				"OnRecordBeforeDeleteRequest": 1,
			},
		},

		// cascade delete checks
		// -----------------------------------------------------------
		{
			Name:   "trying to delete a record while being part of a non-cascade required relation",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo3/records/7nwo8tuiatetxdm",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnModelBeforeUpdate":         2, // self_rel_many update of test1 record + rel_one_cascade demo4 cascaded in demo5
				"OnModelBeforeDelete":         2, // the record itself + rel_one_cascade of test1 record
			},
		},
		{
			Name:   "delete a record with non-cascade references",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo3/records/1tmknxy2868d869",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":         1,
				"OnModelAfterDelete":          1,
				"OnModelBeforeUpdate":         2,
				"OnModelAfterUpdate":          2,
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
			},
		},
		{
			Name:   "delete a record with cascade references",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/oap640cot4yru2s",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":         2,
				"OnModelAfterDelete":          2,
				"OnModelBeforeUpdate":         2,
				"OnModelAfterUpdate":          2,
				"OnRecordBeforeDeleteRequest": 1,
				"OnRecordAfterDeleteRequest":  1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				recId := "84nmscqy84lsi1t"
				rec, _ := app.Dao().FindRecordById("demo1", recId, nil)
				if rec != nil {
					t.Errorf("Expected record %s to be cascade deleted", recId)
				}
				ensureDeletedFiles(app, "wsmn24bux7wo113", recId)
				ensureDeletedFiles(app, "_pb_users_auth_", "oap640cot4yru2s")
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordCrudCreate(t *testing.T) {
	t.Parallel()

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title": "title_test",
	}, "files")
	if err != nil {
		t.Fatal(err)
	}

	formData2, mp2, err2 := tests.MockMultipartData(map[string]string{
		rest.MultipartJsonKey: `{"title": "title_test2", "testPayload": 123}`,
	}, "files")
	if err2 != nil {
		t.Fatal(err2)
	}

	formData3, mp3, err3 := tests.MockMultipartData(map[string]string{
		rest.MultipartJsonKey: `{"title": "title_test3", "testPayload": 123}`,
	}, "files")
	if err3 != nil {
		t.Fatal(err3)
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
			Url:             "/api/collections/demo1/records",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record trying to access nil-rule collection",
			Method: http.MethodPost,
			Url:    "/api/collections/demo1/records",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "trying to create a new view collection record",
			Method:          http.MethodPost,
			Url:             "/api/collections/view1/records",
			Body:            strings.NewReader(`{"text":"new"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit nil body",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo2/records",
			Body:            nil,
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit invalid format",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo2/records",
			Body:            strings.NewReader(`{"`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "submit empty json body",
			Method:         http.MethodPost,
			Url:            "/api/collections/nologin/records",
			Body:           strings.NewReader(`{}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{"code":"validation_required"`,
				`"password":{"code":"validation_required"`,
				`"passwordConfirm":{"code":"validation_required"`,
			},
		},
		{
			Name:           "guest submit in public collection",
			Method:         http.MethodPost,
			Url:            "/api/collections/demo2/records",
			Body:           strings.NewReader(`{"title":"new"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"title":"new"`,
				`"active":false`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:            "guest trying to submit in restricted collection",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo3/records",
			Body:            strings.NewReader(`{"title":"test123"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record submit in restricted collection (rule failure check)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body:   strings.NewReader(`{"title":"test123"}`),
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record submit in restricted collection (rule pass check) + expand relations",
			Method: http.MethodPost,
			Url:    "/api/collections/demo4/records?expand=missing,rel_one_no_cascade,rel_many_no_cascade_required",
			Body: strings.NewReader(`{
				"title":"test123",
				"rel_one_no_cascade":"mk5fmymtx4wsprk",
				"rel_one_no_cascade_required":"7nwo8tuiatetxdm",
				"rel_one_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"],
				"rel_many_cascade":"lcl9d87w22ml6jy"
			}`),
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"title":"test123"`,
				`"rel_one_no_cascade":"mk5fmymtx4wsprk"`,
				`"rel_one_no_cascade_required":"7nwo8tuiatetxdm"`,
				`"rel_one_cascade":"mk5fmymtx4wsprk"`,
				`"rel_many_no_cascade":["mk5fmymtx4wsprk"]`,
				`"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"]`,
				`"rel_many_cascade":["lcl9d87w22ml6jy"]`,
			},
			NotExpectedContent: []string{
				// the users auth records don't have access to view the demo3 expands
				`"expand":{`,
				`"missing"`,
				`"id":"mk5fmymtx4wsprk"`,
				`"id":"7nwo8tuiatetxdm"`,
				`"id":"lcl9d87w22ml6jy"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "admin submit in restricted collection (rule skip check) + expand relations",
			Method: http.MethodPost,
			Url:    "/api/collections/demo4/records?expand=missing,rel_one_no_cascade,rel_many_no_cascade_required",
			Body: strings.NewReader(`{
				"title":"test123",
				"rel_one_no_cascade":"mk5fmymtx4wsprk",
				"rel_one_no_cascade_required":"7nwo8tuiatetxdm",
				"rel_one_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"],
				"rel_many_cascade":"lcl9d87w22ml6jy"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"title":"test123"`,
				`"rel_one_no_cascade":"mk5fmymtx4wsprk"`,
				`"rel_one_no_cascade_required":"7nwo8tuiatetxdm"`,
				`"rel_one_cascade":"mk5fmymtx4wsprk"`,
				`"rel_many_no_cascade":["mk5fmymtx4wsprk"]`,
				`"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"]`,
				`"rel_many_cascade":["lcl9d87w22ml6jy"]`,
				`"expand":{`,
				`"id":"mk5fmymtx4wsprk"`,
				`"id":"7nwo8tuiatetxdm"`,
				`"id":"lcl9d87w22ml6jy"`,
			},
			NotExpectedContent: []string{
				`"missing"`,
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
			Url:    "/api/collections/demo3/records",
			Body:   formData,
			RequestHeaders: map[string]string{
				"Content-Type":  mp.FormDataContentType(),
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"`,
				`"title":"title_test"`,
				`"files":["`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "submit via multipart form data with @jsonPayload key and unsatisfied @request.data rule",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body:   formData2,
			RequestHeaders: map[string]string{
				"Content-Type": mp2.FormDataContentType(),
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collection, err := app.Dao().FindCollectionByNameOrId("demo3")
				if err != nil {
					t.Fatalf("failed to find demo3 collection: %v", err)
				}
				collection.CreateRule = types.Pointer("@request.data.testPayload != 123")
				if err := app.Dao().WithoutHooks().SaveCollection(collection); err != nil {
					t.Fatalf("failed to update demo3 collection create rule: %v", err)
				}
				core.ReloadCachedCollections(app)
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "submit via multipart form data with @jsonPayload key and satisfied @request.data rule",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body:   formData3,
			RequestHeaders: map[string]string{
				"Content-Type": mp3.FormDataContentType(),
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collection, err := app.Dao().FindCollectionByNameOrId("demo3")
				if err != nil {
					t.Fatalf("failed to find demo3 collection: %v", err)
				}
				collection.CreateRule = types.Pointer("@request.data.testPayload = 123")
				if err := app.Dao().WithoutHooks().SaveCollection(collection); err != nil {
					t.Fatalf("failed to update demo3 collection create rule: %v", err)
				}
				core.ReloadCachedCollections(app)
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"`,
				`"title":"title_test3"`,
				`"files":["`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "unique field error check",
			Method: http.MethodPost,
			Url:    "/api/collections/demo2/records",
			Body: strings.NewReader(`{
				"title":"test2"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"title":{`,
				`"code":"validation_not_unique"`,
			},
		},
		{
			Name:   "OnRecordAfterCreateRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections/demo2/records",
			Body:   strings.NewReader(`{"title":"new"}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},

		// ID checks
		// -----------------------------------------------------------
		{
			Name:   "invalid custom insertion id (less than 15 chars)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body: strings.NewReader(`{
				"id": "12345678901234",
				"title": "test"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"id":{"code":"validation_length_invalid"`,
			},
		},
		{
			Name:   "invalid custom insertion id (more than 15 chars)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body: strings.NewReader(`{
				"id": "1234567890123456",
				"title": "test"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"id":{"code":"validation_length_invalid"`,
			},
		},
		{
			Name:   "valid custom insertion id (exactly 15 chars)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body: strings.NewReader(`{
				"id": "123456789012345",
				"title": "test"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"123456789012345"`,
				`"title":"test"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "valid custom insertion id existing in another non-auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/demo3/records",
			Body: strings.NewReader(`{
				"id": "0yxhwia2amd8gec",
				"title": "test"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"0yxhwia2amd8gec"`,
				`"title":"test"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
				"OnRecordAfterCreateRequest":  1,
				"OnModelBeforeCreate":         1,
				"OnModelAfterCreate":          1,
			},
		},
		{
			Name:   "valid custom insertion auth id duplicating in another auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/users/records",
			Body: strings.NewReader(`{
				"id":"o1y0dd0spd786md",
				"title":"test",
				"password":"1234567890",
				"passwordConfirm":"1234567890"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeCreateRequest": 1,
			},
		},

		// fields modifier checks
		// -----------------------------------------------------------
		{
			Name:   "trying to delete a record while being part of a non-cascade required relation",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo3/records/7nwo8tuiatetxdm",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeDeleteRequest": 1,
				"OnModelBeforeUpdate":         2, // self_rel_many update of test1 record + rel_one_cascade demo4 cascaded in demo5
				"OnModelBeforeDelete":         2, // the record itself + rel_one_cascade of test1 record
			},
		},

		// check whether if @request.data modifer fields are properly resolved
		// -----------------------------------------------------------
		{
			Name:   "@request.data.field with compute modifers (rule failure check)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo5/records",
			Body: strings.NewReader(`{
				"total":1,
				"total+":4,
				"total-":1
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
		},
		{
			Name:   "@request.data.field with compute modifers (rule pass check)",
			Method: http.MethodPost,
			Url:    "/api/collections/demo5/records",
			Body: strings.NewReader(`{
				"total":1,
				"total+":3,
				"total-":1
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"`,
				`"collectionName":"demo5"`,
				`"total":3`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterCreate":          1,
				"OnModelBeforeCreate":         1,
				"OnRecordAfterCreateRequest":  1,
				"OnRecordBeforeCreateRequest": 1,
			},
		},

		// auth records
		// -----------------------------------------------------------
		{
			Name:   "auth record with invalid data",
			Method: http.MethodPost,
			Url:    "/api/collections/users/records",
			Body: strings.NewReader(`{
				"id":"o1y0pd786mq",
				"username":"Users75657",
				"email":"invalid",
				"password":"1234567",
				"passwordConfirm":"1234560"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"id":{"code":"validation_length_invalid"`,
				`"username":{"code":"validation_invalid_username"`, // for duplicated case-insensitive username
				`"email":{"code":"validation_is_email"`,
				`"password":{"code":"validation_length_out_of_range"`,
				`"passwordConfirm":{"code":"validation_values_mismatch"`,
			},
			NotExpectedContent: []string{
				// schema fields are not checked if the base fields has errors
				`"rel":{"code":`,
			},
		},
		{
			Name:   "auth record with valid base fields but invalid schema data",
			Method: http.MethodPost,
			Url:    "/api/collections/users/records",
			Body: strings.NewReader(`{
				"password":"12345678",
				"passwordConfirm":"12345678",
				"rel":"invalid"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"rel":{"code":`,
			},
		},
		{
			Name:   "auth record with valid data and explicitly verified state by guest",
			Method: http.MethodPost,
			Url:    "/api/collections/users/records",
			Body: strings.NewReader(`{
				"password":"12345678",
				"passwordConfirm":"12345678",
				"verified":true
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"verified":{"code":`,
			},
		},
		{
			Name:   "auth record with valid data and explicitly verified state by random user",
			Method: http.MethodPost,
			Url:    "/api/collections/users/records",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			Body: strings.NewReader(`{
				"password":"12345678",
				"passwordConfirm":"12345678",
				"emailVisibility":true,
				"verified":true
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"verified":{"code":`,
			},
			NotExpectedContent: []string{
				`"emailVisibility":{"code":`,
			},
		},
		{
			Name:   "auth record with valid data by admin",
			Method: http.MethodPost,
			Url:    "/api/collections/users/records",
			Body: strings.NewReader(`{
				"id":"o1o1y0pd78686mq",
				"username":"test.valid",
				"email":"new@example.com",
				"password":"12345678",
				"passwordConfirm":"12345678",
				"rel":"achvryl401bhse3",
				"emailVisibility":true,
				"verified":true
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"o1o1y0pd78686mq"`,
				`"username":"test.valid"`,
				`"email":"new@example.com"`,
				`"rel":"achvryl401bhse3"`,
				`"emailVisibility":true`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"password"`,
				`"passwordConfirm"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterCreate":          1,
				"OnModelBeforeCreate":         1,
				"OnRecordAfterCreateRequest":  1,
				"OnRecordBeforeCreateRequest": 1,
			},
		},
		{
			Name:   "auth record with valid data by auth record with manage access",
			Method: http.MethodPost,
			Url:    "/api/collections/nologin/records",
			Body: strings.NewReader(`{
				"email":"new@example.com",
				"password":"12345678",
				"passwordConfirm":"12345678",
				"name":"test_name",
				"emailVisibility":true,
				"verified":true
			}`),
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"`,
				`"username":"`,
				`"email":"new@example.com"`,
				`"name":"test_name"`,
				`"emailVisibility":true`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"password"`,
				`"passwordConfirm"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterCreate":          1,
				"OnModelBeforeCreate":         1,
				"OnRecordAfterCreateRequest":  1,
				"OnRecordBeforeCreateRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordCrudUpdate(t *testing.T) {
	t.Parallel()

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title": "title_test",
	}, "files")
	if err != nil {
		t.Fatal(err)
	}

	formData2, mp2, err2 := tests.MockMultipartData(map[string]string{
		rest.MultipartJsonKey: `{"title": "title_test2", "testPayload": 123}`,
	}, "files")
	if err2 != nil {
		t.Fatal(err2)
	}

	formData3, mp3, err3 := tests.MockMultipartData(map[string]string{
		rest.MultipartJsonKey: `{"title": "title_test3", "testPayload": 123}`,
	}, "files")
	if err3 != nil {
		t.Fatal(err3)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodPatch,
			Url:             "/api/collections/missing/records/0yxhwia2amd8gec",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "guest trying to access nil-rule collection record",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo1/records/imy661ixudk5izi",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record trying to access nil-rule collection",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1/records/imy661ixudk5izi",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit invalid body",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo2/records/0yxhwia2amd8gec",
			Body:            strings.NewReader(`{"`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "trying to update a view collection record",
			Method:          http.MethodPatch,
			Url:             "/api/collections/view1/records/imy661ixudk5izi",
			Body:            strings.NewReader(`{"text":"new"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "submit nil body",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo2/records/0yxhwia2amd8gec",
			Body:            nil,
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "submit empty body (aka. no fields change)",
			Method:         http.MethodPatch,
			Url:            "/api/collections/demo2/records/0yxhwia2amd8gec",
			Body:           strings.NewReader(`{}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"collectionName":"demo2"`,
				`"id":"0yxhwia2amd8gec"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterUpdateRequest":  1,
				"OnRecordBeforeUpdateRequest": 1,
			},
		},
		{
			Name:           "trigger field validation",
			Method:         http.MethodPatch,
			Url:            "/api/collections/demo2/records/0yxhwia2amd8gec",
			Body:           strings.NewReader(`{"title":"a"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`data":{`,
				`"title":{"code":"validation_min_text_constraint"`,
			},
		},
		{
			Name:           "guest submit in public collection",
			Method:         http.MethodPatch,
			Url:            "/api/collections/demo2/records/0yxhwia2amd8gec",
			Body:           strings.NewReader(`{"title":"new"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"0yxhwia2amd8gec"`,
				`"title":"new"`,
				`"active":true`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:            "guest trying to submit in restricted collection",
			Method:          http.MethodPatch,
			Url:             "/api/collections/demo3/records/mk5fmymtx4wsprk",
			Body:            strings.NewReader(`{"title":"new"}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record submit in restricted collection (rule failure check)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo3/records/mk5fmymtx4wsprk",
			Body:   strings.NewReader(`{"title":"new"}`),
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record submit in restricted collection (rule pass check) + expand relations",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo4/records/i9naidtvr6qsgb4?expand=missing,rel_one_no_cascade,rel_many_no_cascade_required",
			Body: strings.NewReader(`{
				"title":"test123",
				"rel_one_no_cascade":"mk5fmymtx4wsprk",
				"rel_one_no_cascade_required":"7nwo8tuiatetxdm",
				"rel_one_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"],
				"rel_many_cascade":"lcl9d87w22ml6jy"
			}`),
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"i9naidtvr6qsgb4"`,
				`"title":"test123"`,
				`"rel_one_no_cascade":"mk5fmymtx4wsprk"`,
				`"rel_one_no_cascade_required":"7nwo8tuiatetxdm"`,
				`"rel_one_cascade":"mk5fmymtx4wsprk"`,
				`"rel_many_no_cascade":["mk5fmymtx4wsprk"]`,
				`"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"]`,
				`"rel_many_cascade":["lcl9d87w22ml6jy"]`,
			},
			NotExpectedContent: []string{
				// the users auth records don't have access to view the demo3 expands
				`"expand":{`,
				`"missing"`,
				`"id":"mk5fmymtx4wsprk"`,
				`"id":"7nwo8tuiatetxdm"`,
				`"id":"lcl9d87w22ml6jy"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "admin submit in restricted collection (rule skip check) + expand relations",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo4/records/i9naidtvr6qsgb4?expand=missing,rel_one_no_cascade,rel_many_no_cascade_required",
			Body: strings.NewReader(`{
				"title":"test123",
				"rel_one_no_cascade":"mk5fmymtx4wsprk",
				"rel_one_no_cascade_required":"7nwo8tuiatetxdm",
				"rel_one_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade":"mk5fmymtx4wsprk",
				"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"],
				"rel_many_cascade":"lcl9d87w22ml6jy"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"i9naidtvr6qsgb4"`,
				`"title":"test123"`,
				`"rel_one_no_cascade":"mk5fmymtx4wsprk"`,
				`"rel_one_no_cascade_required":"7nwo8tuiatetxdm"`,
				`"rel_one_cascade":"mk5fmymtx4wsprk"`,
				`"rel_many_no_cascade":["mk5fmymtx4wsprk"]`,
				`"rel_many_no_cascade_required":["7nwo8tuiatetxdm","lcl9d87w22ml6jy"]`,
				`"rel_many_cascade":["lcl9d87w22ml6jy"]`,
				`"expand":{`,
				`"id":"mk5fmymtx4wsprk"`,
				`"id":"7nwo8tuiatetxdm"`,
				`"id":"lcl9d87w22ml6jy"`,
			},
			NotExpectedContent: []string{
				`"missing"`,
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
			Url:    "/api/collections/demo3/records/mk5fmymtx4wsprk",
			Body:   formData,
			RequestHeaders: map[string]string{
				"Content-Type":  mp.FormDataContentType(),
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"mk5fmymtx4wsprk"`,
				`"title":"title_test"`,
				`"files":["`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "submit via multipart form data with @jsonPayload key and unsatisfied @request.data rule",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo3/records/mk5fmymtx4wsprk",
			Body:   formData2,
			RequestHeaders: map[string]string{
				"Content-Type": mp2.FormDataContentType(),
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collection, err := app.Dao().FindCollectionByNameOrId("demo3")
				if err != nil {
					t.Fatalf("failed to find demo3 collection: %v", err)
				}
				collection.UpdateRule = types.Pointer("@request.data.testPayload != 123")
				if err := app.Dao().WithoutHooks().SaveCollection(collection); err != nil {
					t.Fatalf("failed to update demo3 collection update rule: %v", err)
				}
				core.ReloadCachedCollections(app)
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "submit via multipart form data with @jsonPayload key and satisfied @request.data rule",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo3/records/mk5fmymtx4wsprk",
			Body:   formData3,
			RequestHeaders: map[string]string{
				"Content-Type": mp3.FormDataContentType(),
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				collection, err := app.Dao().FindCollectionByNameOrId("demo3")
				if err != nil {
					t.Fatalf("failed to find demo3 collection: %v", err)
				}
				collection.UpdateRule = types.Pointer("@request.data.testPayload = 123")
				if err := app.Dao().WithoutHooks().SaveCollection(collection); err != nil {
					t.Fatalf("failed to update demo3 collection update rule: %v", err)
				}
				core.ReloadCachedCollections(app)
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"mk5fmymtx4wsprk"`,
				`"title":"title_test3"`,
				`"files":["`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "OnRecordAfterUpdateRequest error response",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2/records/0yxhwia2amd8gec",
			Body:   strings.NewReader(`{"title":"new"}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnRecordAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":         1,
				"OnModelAfterUpdate":          1,
			},
		},
		{
			Name:   "try to change the id of an existing record",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo3/records/mk5fmymtx4wsprk",
			Body: strings.NewReader(`{
				"id": "mk5fmymtx4wspra"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"id":{"code":"validation_in_invalid"`,
			},
		},
		{
			Name:   "unique field error check",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2/records/llvuca81nly1qls",
			Body: strings.NewReader(`{
				"title":"test2"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"title":{`,
				`"code":"validation_not_unique"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeUpdateRequest": 1,
				"OnModelBeforeUpdate":         1,
			},
		},

		// check whether if @request.data modifer fields are properly resolved
		// -----------------------------------------------------------
		{
			Name:   "@request.data.field with compute modifers (rule failure check)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo5/records/la4y2w4o98acwuj",
			Body: strings.NewReader(`{
				"total+":3,
				"total-":1
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"data":{}`,
			},
		},
		{
			Name:   "@request.data.field with compute modifers (rule pass check)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo5/records/la4y2w4o98acwuj",
			Body: strings.NewReader(`{
				"total+":2,
				"total-":1
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"la4y2w4o98acwuj"`,
				`"collectionName":"demo5"`,
				`"total":3`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterUpdateRequest":  1,
				"OnRecordBeforeUpdateRequest": 1,
			},
		},

		// auth records
		// -----------------------------------------------------------
		{
			Name:   "auth record with invalid data",
			Method: http.MethodPatch,
			Url:    "/api/collections/users/records/bgs820n361vj1qd",
			Body: strings.NewReader(`{
				"username":"Users75657",
				"email":"invalid",
				"password":"1234567",
				"passwordConfirm":"1234560",
				"verified":false
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"username":{"code":"validation_invalid_username"`, // for duplicated case-insensitive username
				`"email":{"code":"validation_is_email"`,
				`"password":{"code":"validation_length_out_of_range"`,
				`"passwordConfirm":{"code":"validation_values_mismatch"`,
			},
			NotExpectedContent: []string{
				// admins are allowed to change the verified state
				`"verified"`,
				// schema fields are not checked if the base fields has errors
				`"rel":{"code":`,
			},
		},
		{
			Name:   "auth record with valid base fields but invalid schema data",
			Method: http.MethodPatch,
			Url:    "/api/collections/users/records/bgs820n361vj1qd",
			Body: strings.NewReader(`{
				"password":"12345678",
				"passwordConfirm":"12345678",
				"rel":"invalid"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"rel":{"code":`,
			},
		},
		{
			Name:   "try to change account managing fields by guest",
			Method: http.MethodPatch,
			Url:    "/api/collections/nologin/records/phhq3wr65cap535",
			Body: strings.NewReader(`{
				"password":"12345678",
				"passwordConfirm":"12345678",
				"emailVisibility":true,
				"verified":true
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"verified":{"code":`,
				`"oldPassword":{"code":`,
			},
			NotExpectedContent: []string{
				`"emailVisibility":{"code":`,
			},
		},
		{
			Name:   "try to change account managing fields by auth record (owner)",
			Method: http.MethodPatch,
			Url:    "/api/collections/users/records/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			Body: strings.NewReader(`{
				"password":"12345678",
				"passwordConfirm":"12345678",
				"emailVisibility":true,
				"verified":true
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"verified":{"code":`,
				`"oldPassword":{"code":`,
			},
			NotExpectedContent: []string{
				`"emailVisibility":{"code":`,
			},
		},
		{
			Name:   "try to change account managing fields by auth record with managing rights",
			Method: http.MethodPatch,
			Url:    "/api/collections/nologin/records/phhq3wr65cap535",
			Body: strings.NewReader(`{
				"email":"new@example.com",
				"password":"12345678",
				"passwordConfirm":"12345678",
				"name":"test_name",
				"emailVisibility":true,
				"verified":true
			}`),
			RequestHeaders: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"new@example.com"`,
				`"name":"test_name"`,
				`"emailVisibility":true`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"password"`,
				`"passwordConfirm"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterUpdateRequest":  1,
				"OnRecordBeforeUpdateRequest": 1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				record, _ := app.Dao().FindRecordById("nologin", "phhq3wr65cap535")
				if !record.ValidatePassword("12345678") {
					t.Fatal("Password update failed.")
				}
			},
		},
		{
			Name:   "update auth record with valid data by admin",
			Method: http.MethodPatch,
			Url:    "/api/collections/users/records/oap640cot4yru2s",
			Body: strings.NewReader(`{
				"username":"test.valid",
				"email":"new@example.com",
				"password":"12345678",
				"passwordConfirm":"12345678",
				"rel":"achvryl401bhse3",
				"emailVisibility":true,
				"verified":false
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"username":"test.valid"`,
				`"email":"new@example.com"`,
				`"rel":"achvryl401bhse3"`,
				`"emailVisibility":true`,
				`"verified":false`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"password"`,
				`"passwordConfirm"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterUpdateRequest":  1,
				"OnRecordBeforeUpdateRequest": 1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				record, _ := app.Dao().FindRecordById("users", "oap640cot4yru2s")
				if !record.ValidatePassword("12345678") {
					t.Fatal("Password update failed.")
				}
			},
		},
		{
			Name:   "update auth record with valid data by guest (empty update filter)",
			Method: http.MethodPatch,
			Url:    "/api/collections/nologin/records/dc49k6jgejn40h3",
			Body: strings.NewReader(`{
				"username":"test_new",
				"emailVisibility":true,
				"name":"test"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"username":"test_new"`,
				`"email":"test@example.com"`, // the email should be visible since we updated the emailVisibility
				`"emailVisibility":true`,
				`"verified":false`,
				`"name":"test"`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"password"`,
				`"passwordConfirm"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterUpdateRequest":  1,
				"OnRecordBeforeUpdateRequest": 1,
			},
		},
		{
			Name:   "success password change with oldPassword",
			Method: http.MethodPatch,
			Url:    "/api/collections/nologin/records/dc49k6jgejn40h3",
			Body: strings.NewReader(`{
				"password":"123456789",
				"passwordConfirm":"123456789",
				"oldPassword":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"dc49k6jgejn40h3"`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"password"`,
				`"passwordConfirm"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":          1,
				"OnModelBeforeUpdate":         1,
				"OnRecordAfterUpdateRequest":  1,
				"OnRecordBeforeUpdateRequest": 1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				record, _ := app.Dao().FindRecordById("nologin", "dc49k6jgejn40h3")
				if !record.ValidatePassword("123456789") {
					t.Fatal("Password update failed.")
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
