package apis_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/router"
)

func TestBatchRequest(t *testing.T) {
	t.Parallel()

	formData, mp, err := tests.MockMultipartData(
		map[string]string{
			router.JSONPayloadKey: `{
				"requests":[
					{"method":"POST", "url":"/api/collections/demo3/records", "body": {"title": "batch1"}},
					{"method":"POST", "url":"/api/collections/demo3/records", "body": {"title": "batch2"}},
					{"method":"POST", "url":"/api/collections/demo3/records", "body": {"title": "batch3"}},
					{"method":"PATCH", "url":"/api/collections/demo3/records/lcl9d87w22ml6jy", "body": {"files-": "test_FLurQTgrY8.txt"}}
				]
			}`,
		},
		"requests.0.files",
		"requests.0.files",
		"requests.0.files",
		"requests[2].files",
	)
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:   "disabled batch requets",
			Method: http.MethodPost,
			URL:    "/api/batch",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().Batch.Enabled = false
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "max request limits reached",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
					{"method":"GET", "url":"/test1"},
					{"method":"GET", "url":"/test2"},
					{"method":"GET", "url":"/test3"}
				]
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().Batch.Enabled = true
				app.Settings().Batch.MaxRequests = 2
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"requests":{"code":"validation_length_too_long"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "trigger requests validations",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
					{},
					{"method":"GET", "url":"/valid"},
					{"method":"invalid", "url":"/valid"},
					{"method":"POST", "url":"` + strings.Repeat("a", 2001) + `"}
				]
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().Batch.Enabled = true
				app.Settings().Batch.MaxRequests = 100
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"requests":{`,
				`"0":{"method":{"code":"validation_required"`,
				`"2":{"method":{"code":"validation_in_invalid"`,
				`"3":{"url":{"code":"validation_length_too_long"`,
			},
			NotExpectedContent: []string{
				`"1":`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "unknown batch request action",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
					{"method":"GET", "url":"/api/health"}
				]
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"requests":{`,
				`0":{"code":"batch_request_failed"`,
				`"response":{`,
			},
			ExpectedEvents: map[string]int{
				"*":              0,
				"OnBatchRequest": 1,
			},
		},
		{
			Name:   "base 2 successful and 1 failed (public collection)",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
                    {"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch1"}},
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch2"}},
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": ""}}
				]
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"response":{`,
				`"2":{"code":"batch_request_failed"`,
				`"response":{"data":{"title":{"code":"validation_required"`,
			},
			NotExpectedContent: []string{
				`"0":`,
				`"1":`,
			},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnBatchRequest":           1,
				"OnRecordCreateRequest":    3,
				"OnModelCreate":            3,
				"OnModelCreateExecute":     2,
				"OnModelAfterCreateError":  3,
				"OnModelValidate":          3,
				"OnRecordCreate":           3,
				"OnRecordCreateExecute":    2,
				"OnRecordAfterCreateError": 3,
				"OnRecordValidate":         3,
				"OnRecordEnrich":           2,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter("demo2", `title~"batch"`, "", 0, 0)
				if err != nil {
					t.Fatal(err)
				}

				if len(records) != 0 {
					t.Fatalf("Expected no batch records to be persisted, got %d", len(records))
				}
			},
		},
		{
			Name:   "base 4 successful (public collection)",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch1"}},
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch2"}},
					{"method":"PUT", "url":"/api/collections/demo2/records", "body": {"title": "batch3"}},
					{"method":"PUT", "url":"/api/collections/demo2/records?fields=*,id:excerpt(4,true)", "body": {"id":"achvryl401bhse3","title": "batch4"}}
				]
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"title":"batch1"`,
				`"title":"batch2"`,
				`"title":"batch3"`,
				`"title":"batch4"`,
				`"id":"achv..."`,
				`"active":false`,
				`"active":true`,
				`"status":200`,
				`"body":{`,
			},
			ExpectedEvents: map[string]int{
				"*":                0,
				"OnBatchRequest":   1,
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
				"OnRecordEnrich":   4,

				"OnRecordCreateRequest":      3,
				"OnModelCreate":              3,
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,

				"OnRecordUpdateRequest":      1,
				"OnModelUpdate":              1,
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter("demo2", `title~"batch"`, "", 0, 0)
				if err != nil {
					t.Fatal(err)
				}

				if len(records) != 4 {
					t.Fatalf("Expected %d batch records to be persisted, got %d", 3, len(records))
				}
			},
		},
		{
			Name:   "mixed create/update/delete (rules failure)",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch_create"}},
					{"method":"DELETE", "url":"/api/collections/demo2/records/achvryl401bhse3"},
					{"method":"PATCH", "url":"/api/collections/demo3/records/1tmknxy2868d869", "body": {"title": "batch_update"}}
				]
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"requests":{`,
				`"2":{"code":"batch_request_failed"`,
				`"response":{`,
			},
			NotExpectedContent: []string{
				// only demo3 requires authentication
				`"0":`,
				`"1":`,
			},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnBatchRequest":           1,
				"OnModelCreate":            1,
				"OnModelCreateExecute":     1,
				"OnModelAfterCreateError":  1,
				"OnModelDelete":            1,
				"OnModelDeleteExecute":     1,
				"OnModelAfterDeleteError":  1,
				"OnModelValidate":          1,
				"OnRecordCreateRequest":    1,
				"OnRecordCreate":           1,
				"OnRecordCreateExecute":    1,
				"OnRecordAfterCreateError": 1,
				"OnRecordDeleteRequest":    1,
				"OnRecordDelete":           1,
				"OnRecordDeleteExecute":    1,
				"OnRecordAfterDeleteError": 1,
				"OnRecordEnrich":           1,
				"OnRecordValidate":         1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				_, err := app.FindFirstRecordByFilter("demo2", `title="batch_create"`)
				if err == nil {
					t.Fatal("Expected record to not be created")
				}

				_, err = app.FindFirstRecordByFilter("demo3", `title="batch_update"`)
				if err == nil {
					t.Fatal("Expected record to not be updated")
				}

				_, err = app.FindRecordById("demo2", "achvryl401bhse3")
				if err != nil {
					t.Fatal("Expected record to not be deleted")
				}
			},
		},
		{
			Name:   "mixed create/update/delete (rules success)",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Headers: map[string]string{
				// test@example.com, clients
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.0ONnm_BsvPRZyDNT31GN1CKUB6uQRxvVvQ-Wc9AZfG0",
			},
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch_create"}, "headers": {"Authorization": "ignored"}},
					{"method":"DELETE", "url":"/api/collections/demo2/records/achvryl401bhse3", "headers": {"Authorization": "ignored"}},
					{"method":"PATCH", "url":"/api/collections/demo3/records/1tmknxy2868d869", "body": {"title": "batch_update"}, "headers": {"Authorization": "ignored"}}
				]
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"title":"batch_create"`,
				`"title":"batch_update"`,
				`"status":200`,
				`"status":204`,
				`"body":{`,
				`"body":null`,
			},
			ExpectedEvents: map[string]int{
				"*":              0,
				"OnBatchRequest": 1,
				// ---
				"OnModelCreate":             1,
				"OnModelCreateExecute":      1,
				"OnModelAfterCreateSuccess": 1,
				"OnModelDelete":             1,
				"OnModelDeleteExecute":      1,
				"OnModelAfterDeleteSuccess": 1,
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				"OnModelValidate":           2,
				// ---
				"OnRecordCreateRequest":      1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordDeleteRequest":      1,
				"OnRecordDelete":             1,
				"OnRecordDeleteExecute":      1,
				"OnRecordAfterDeleteSuccess": 1,
				"OnRecordUpdateRequest":      1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				"OnRecordValidate":           2,
				"OnRecordEnrich":             2,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				_, err := app.FindFirstRecordByFilter("demo2", `title="batch_create"`)
				if err != nil {
					t.Fatal(err)
				}

				_, err = app.FindFirstRecordByFilter("demo3", `title="batch_update"`)
				if err != nil {
					t.Fatal(err)
				}

				_, err = app.FindRecordById("demo2", "achvryl401bhse3")
				if err == nil {
					t.Fatal("Expected record to be deleted")
				}
			},
		},
		{
			Name:   "mixed create/update/delete (superuser auth)",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Headers: map[string]string{
				// test@example.com, _superusers
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch_create"}},
					{"method":"DELETE", "url":"/api/collections/demo2/records/achvryl401bhse3"},
					{"method":"PATCH", "url":"/api/collections/demo3/records/1tmknxy2868d869", "body": {"title": "batch_update"}}
				]
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"title":"batch_create"`,
				`"title":"batch_update"`,
				`"status":200`,
				`"status":204`,
				`"body":{`,
				`"body":null`,
			},
			ExpectedEvents: map[string]int{
				"*":              0,
				"OnBatchRequest": 1,
				// ---
				"OnModelCreate":             1,
				"OnModelCreateExecute":      1,
				"OnModelAfterCreateSuccess": 1,
				"OnModelDelete":             1,
				"OnModelDeleteExecute":      1,
				"OnModelAfterDeleteSuccess": 1,
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				"OnModelValidate":           2,
				// ---
				"OnRecordCreateRequest":      1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordDeleteRequest":      1,
				"OnRecordDelete":             1,
				"OnRecordDeleteExecute":      1,
				"OnRecordAfterDeleteSuccess": 1,
				"OnRecordUpdateRequest":      1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				"OnRecordValidate":           2,
				"OnRecordEnrich":             2,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				_, err := app.FindFirstRecordByFilter("demo2", `title="batch_create"`)
				if err != nil {
					t.Fatal(err)
				}

				_, err = app.FindFirstRecordByFilter("demo3", `title="batch_update"`)
				if err != nil {
					t.Fatal(err)
				}

				_, err = app.FindRecordById("demo2", "achvryl401bhse3")
				if err == nil {
					t.Fatal("Expected record to be deleted")
				}
			},
		},
		{
			Name:   "cascade delete/update",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Headers: map[string]string{
				// test@example.com, _superusers
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{
				"requests": [
					{"method":"DELETE", "url":"/api/collections/demo3/records/1tmknxy2868d869"},
					{"method":"DELETE", "url":"/api/collections/demo3/records/mk5fmymtx4wsprk"}
				]
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":204`,
				`"body":null`,
			},
			NotExpectedContent: []string{
				`"status":200`,
				`"body":{`,
			},
			ExpectedEvents: map[string]int{
				"*":              0,
				"OnBatchRequest": 1,
				// ---
				"OnModelDelete":             3, // 2 batch + 1 cascade delete
				"OnModelDeleteExecute":      3,
				"OnModelAfterDeleteSuccess": 3,
				"OnModelUpdate":             5, // 5 cascade update
				"OnModelUpdateExecute":      5,
				"OnModelAfterUpdateSuccess": 5,
				// ---
				"OnRecordDeleteRequest":      2,
				"OnRecordDelete":             3,
				"OnRecordDeleteExecute":      3,
				"OnRecordAfterDeleteSuccess": 3,
				"OnRecordUpdate":             5,
				"OnRecordUpdateExecute":      5,
				"OnRecordAfterUpdateSuccess": 5,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				ids := []string{
					"1tmknxy2868d869",
					"mk5fmymtx4wsprk",
					"qzaqccwrmva4o1n",
				}

				for _, id := range ids {
					_, err := app.FindRecordById("demo2", id)
					if err == nil {
						t.Fatalf("Expected record %q to be deleted", id)
					}
				}
			},
		},
		{
			Name:   "transaction timeout",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch1"}},
					{"method":"POST", "url":"/api/collections/demo2/records", "body": {"title": "batch2"}}
				]
			}`),
			Headers: map[string]string{
				// test@example.com, _superusers
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().Batch.Timeout = 1
				app.OnRecordCreateRequest("demo2").BindFunc(func(e *core.RecordRequestEvent) error {
					time.Sleep(600 * time.Millisecond) // < 1s so that the first request can succeed
					return e.Next()
				})
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnBatchRequest":           1,
				"OnRecordCreateRequest":    2,
				"OnModelCreate":            1,
				"OnModelCreateExecute":     1,
				"OnModelAfterCreateError":  1,
				"OnModelValidate":          1,
				"OnRecordCreate":           1,
				"OnRecordCreateExecute":    1,
				"OnRecordAfterCreateError": 1,
				"OnRecordEnrich":           1,
				"OnRecordValidate":         1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				records, err := app.FindRecordsByFilter("demo2", `title~"batch"`, "", 0, 0)
				if err != nil {
					t.Fatal(err)
				}

				if len(records) != 0 {
					t.Fatalf("Expected %d batch records to be persisted, got %d", 0, len(records))
				}
			},
		},
		{
			Name:   "multipart/form-data + file upload",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Body:   formData,
			Headers: map[string]string{
				// test@example.com, clients
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.0ONnm_BsvPRZyDNT31GN1CKUB6uQRxvVvQ-Wc9AZfG0",
				"Content-Type":  mp.FormDataContentType(),
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"title":"batch1"`,
				`"title":"batch2"`,
				`"title":"batch3"`,
				`"id":"lcl9d87w22ml6jy"`,
				`"files":["300_UhLKX91HVb.png"]`,
				`"tmpfile_`,
				`"status":200`,
				`"body":{`,
			},
			ExpectedEvents: map[string]int{
				"*":              0,
				"OnBatchRequest": 1,
				// ---
				"OnModelCreate":             3,
				"OnModelCreateExecute":      3,
				"OnModelAfterCreateSuccess": 3,
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				"OnModelValidate":           4,
				// ---
				"OnRecordCreateRequest":      3,
				"OnRecordUpdateRequest":      1,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				"OnRecordValidate":           4,
				"OnRecordEnrich":             4,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				batch1, err := app.FindFirstRecordByFilter("demo3", `title="batch1"`)
				if err != nil {
					t.Fatalf("missing batch1: %v", err)
				}
				batch1Files := batch1.GetStringSlice("files")
				if len(batch1Files) != 3 {
					t.Fatalf("Expected %d batch1 file(s), got %d", 3, len(batch1Files))
				}

				batch2, err := app.FindFirstRecordByFilter("demo3", `title="batch2"`)
				if err != nil {
					t.Fatalf("missing batch2: %v", err)
				}
				batch2Files := batch2.GetStringSlice("files")
				if len(batch2Files) != 0 {
					t.Fatalf("Expected %d batch2 file(s), got %d", 0, len(batch2Files))
				}

				batch3, err := app.FindFirstRecordByFilter("demo3", `title="batch3"`)
				if err != nil {
					t.Fatalf("missing batch3: %v", err)
				}
				batch3Files := batch3.GetStringSlice("files")
				if len(batch3Files) != 1 {
					t.Fatalf("Expected %d batch3 file(s), got %d", 1, len(batch3Files))
				}

				batch4, err := app.FindRecordById("demo3", "lcl9d87w22ml6jy")
				if err != nil {
					t.Fatalf("missing batch4: %v", err)
				}
				batch4Files := batch4.GetStringSlice("files")
				if len(batch4Files) != 1 {
					t.Fatalf("Expected %d batch4 file(s), got %d", 1, len(batch4Files))
				}
			},
		},
		{
			Name:   "create/update with expand query params",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Headers: map[string]string{
				// test@example.com, _superusers
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo5/records?expand=rel_one", "body": {"total": 9, "rel_one":"qzaqccwrmva4o1n"}},
					{"method":"PATCH", "url":"/api/collections/demo5/records/qjeql998mtp1azp?expand=rel_many", "body": {"total": 10}}
				]
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"body":{`,
				`"id":"qjeql998mtp1azp"`,
				`"id":"qzaqccwrmva4o1n"`,
				`"id":"i9naidtvr6qsgb4"`,
				`"expand":{"rel_one"`,
				`"expand":{"rel_many"`,
			},
			ExpectedEvents: map[string]int{
				"*":              0,
				"OnBatchRequest": 1,
				// ---
				"OnModelCreate":             1,
				"OnModelCreateExecute":      1,
				"OnModelAfterCreateSuccess": 1,
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				"OnModelValidate":           2,
				// ---
				"OnRecordCreateRequest":      1,
				"OnRecordUpdateRequest":      1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				"OnRecordValidate":           2,
				"OnRecordEnrich":             5,
			},
		},
		{
			Name:   "check body limit middleware",
			Method: http.MethodPost,
			URL:    "/api/batch",
			Headers: map[string]string{
				// test@example.com, _superusers
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{
				"requests": [
					{"method":"POST", "url":"/api/collections/demo5/records?expand=rel_one", "body": {"total": 9, "rel_one":"qzaqccwrmva4o1n"}},
					{"method":"PATCH", "url":"/api/collections/demo5/records/qjeql998mtp1azp?expand=rel_many", "body": {"total": 10}}
				]
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().Batch.MaxBodySize = 10
			},
			ExpectedStatus:  413,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
