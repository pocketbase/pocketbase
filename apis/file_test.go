package apis_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestFileToken(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/files/token",
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnFileBeforeTokenRequest": 1,
			},
		},
		{
			Name:   "unauthorized with model and token via hook",
			Method: http.MethodPost,
			Url:    "/api/files/token",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnFileBeforeTokenRequest().Add(func(e *core.FileTokenEvent) error {
					record, _ := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
					e.Model = record
					e.Token = "test"
					return nil
				})
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"test"`,
			},
			ExpectedEvents: map[string]int{
				"OnFileBeforeTokenRequest": 1,
				"OnFileAfterTokenRequest":  1,
			},
		},
		{
			Name:   "auth record",
			Method: http.MethodPost,
			Url:    "/api/files/token",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"`,
			},
			ExpectedEvents: map[string]int{
				"OnFileBeforeTokenRequest": 1,
				"OnFileAfterTokenRequest":  1,
			},
		},
		{
			Name:   "admin",
			Method: http.MethodPost,
			Url:    "/api/files/token",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"`,
			},
			ExpectedEvents: map[string]int{
				"OnFileBeforeTokenRequest": 1,
				"OnFileAfterTokenRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestFileDownload(t *testing.T) {
	t.Parallel()

	_, currentFile, _, _ := runtime.Caller(0)
	dataDirRelPath := "../tests/data/"

	testFilePath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/oap640cot4yru2s/test_kfd2wYLxkz.txt")
	testImgPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png")
	testThumbCropCenterPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50_300_1SEi6Q6U72.png")
	testThumbCropTopPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50t_300_1SEi6Q6U72.png")
	testThumbCropBottomPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50b_300_1SEi6Q6U72.png")
	testThumbFitPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50f_300_1SEi6Q6U72.png")
	testThumbZeroWidthPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/0x50_300_1SEi6Q6U72.png")
	testThumbZeroHeightPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x0_300_1SEi6Q6U72.png")

	testFile, fileErr := os.ReadFile(testFilePath)
	if fileErr != nil {
		t.Fatal(fileErr)
	}

	testImg, imgErr := os.ReadFile(testImgPath)
	if imgErr != nil {
		t.Fatal(imgErr)
	}

	testThumbCropCenter, thumbErr := os.ReadFile(testThumbCropCenterPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbCropTop, thumbErr := os.ReadFile(testThumbCropTopPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbCropBottom, thumbErr := os.ReadFile(testThumbCropBottomPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbFit, thumbErr := os.ReadFile(testThumbFitPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbZeroWidth, thumbErr := os.ReadFile(testThumbZeroWidthPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbZeroHeight, thumbErr := os.ReadFile(testThumbZeroHeightPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/files/missing/4q1xlclmfloku33/300_1SEi6Q6U72.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/missing/300_1SEi6Q6U72.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing file",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/missing.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "existing image",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testImg)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - missing thumb (should fallback to the original)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=999x999",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testImg)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop center)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropCenter)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop top)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50t",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropTop)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop bottom)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50b",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropBottom)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (fit)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50f",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbFit)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (zero width)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=0x50",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbZeroWidth)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (zero height)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x0",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbZeroHeight)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing non image file - thumb parameter should be ignored",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/oap640cot4yru2s/test_kfd2wYLxkz.txt?thumb=100x100",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testFile)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},

		// protected file access checks
		{
			Name:            "protected file - expired token",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/oap640cot4yru2s/test_kfd2wYLxkz.txt?thumb=100x100",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testFile)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "protected file - admin with expired file token",
			Method:          http.MethodGet,
			Url:             "/api/files/demo1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsImV4cCI6MTY0MDk5MTY2MSwidHlwZSI6ImFkbWluIn0.g7Q_3UX6H--JWJ7yt1Hoe-1ugTX1KpbKzdt0zjGSe-E",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "protected file - admin with valid file token",
			Method:          http.MethodGet,
			Url:             "/api/files/demo1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsImV4cCI6MTg5MzQ1MjQ2MSwidHlwZSI6ImFkbWluIn0.LyAMpSfaHVsuUqIlqqEbhDQSdFzoPz_EIDcb2VJMBsU",
			ExpectedStatus:  200,
			ExpectedContent: []string{"PNG"},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "protected file - guest without view access",
			Method:          http.MethodGet,
			Url:             "/api/files/demo1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "protected file - guest with view access",
			Method: http.MethodGet,
			Url:    "/api/files/demo1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				dao := daos.New(app.Dao().DB())

				// mock public view access
				c, err := dao.FindCollectionByNameOrId("demo1")
				if err != nil {
					t.Fatalf("Failed to fetch mock collection: %v", err)
				}
				c.ViewRule = types.Pointer("")
				if err := dao.SaveCollection(c); err != nil {
					t.Fatalf("Failed to update mock collection: %v", err)
				}
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"PNG"},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:   "protected file - auth record without view access",
			Method: http.MethodGet,
			Url:    "/api/files/demo1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MTg5MzQ1MjQ2MSwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwidHlwZSI6ImF1dGhSZWNvcmQifQ.0d_0EO6kfn9ijZIQWAqgRi8Bo1z7MKcg1LQpXhQsEPk",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				dao := daos.New(app.Dao().DB())

				// mock restricted user view access
				c, err := dao.FindCollectionByNameOrId("demo1")
				if err != nil {
					t.Fatalf("Failed to fetch mock collection: %v", err)
				}
				c.ViewRule = types.Pointer("@request.auth.verified = true")
				if err := dao.SaveCollection(c); err != nil {
					t.Fatalf("Failed to update mock collection: %v", err)
				}
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "protected file - auth record with view access",
			Method: http.MethodGet,
			Url:    "/api/files/demo1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MTg5MzQ1MjQ2MSwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwidHlwZSI6ImF1dGhSZWNvcmQifQ.0d_0EO6kfn9ijZIQWAqgRi8Bo1z7MKcg1LQpXhQsEPk",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				dao := daos.New(app.Dao().DB())

				// mock user view access
				c, err := dao.FindCollectionByNameOrId("demo1")
				if err != nil {
					t.Fatalf("Failed to fetch mock collection: %v", err)
				}
				c.ViewRule = types.Pointer("@request.auth.verified = false")
				if err := dao.SaveCollection(c); err != nil {
					t.Fatalf("Failed to update mock collection: %v", err)
				}
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"PNG"},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "protected file in view (view's View API rule failure)",
			Method:          http.MethodGet,
			Url:             "/api/files/view1/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MTg5MzQ1MjQ2MSwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwidHlwZSI6ImF1dGhSZWNvcmQifQ.0d_0EO6kfn9ijZIQWAqgRi8Bo1z7MKcg1LQpXhQsEPk",
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "protected file in view (view's View API rule success)",
			Method:          http.MethodGet,
			Url:             "/api/files/view1/84nmscqy84lsi1t/test_d61b33QdDU.txt?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MTg5MzQ1MjQ2MSwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwidHlwZSI6ImF1dGhSZWNvcmQifQ.0d_0EO6kfn9ijZIQWAqgRi8Bo1z7MKcg1LQpXhQsEPk",
			ExpectedStatus:  200,
			ExpectedContent: []string{"test"},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		// clone for the HEAD test (the same as the original scenario but without body)
		head := scenario
		head.Method = http.MethodHead
		head.Name = ("(HEAD) " + scenario.Name)
		head.ExpectedContent = nil
		head.Test(t)

		// regular request test
		scenario.Test(t)
	}
}

func TestConcurrentThumbsGeneration(t *testing.T) {
	t.Parallel()

	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatal(err)
	}
	defer app.Cleanup()

	fsys, err := app.NewFilesystem()
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	// create a dummy file field collection
	demo1, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}
	fileField := demo1.Schema.GetFieldByName("file_one")
	fileField.Options = &schema.FileOptions{
		Protected: false,
		MaxSelect: 1,
		MaxSize:   999999,
		// new thumbs
		Thumbs: []string{"111x111", "111x222", "111x333"},
	}
	demo1.Schema.AddField(fileField)
	if err := app.Dao().SaveCollection(demo1); err != nil {
		t.Fatal(err)
	}

	fileKey := "wsmn24bux7wo113/al1h9ijdeojtsjy/300_Jsjq7RdBgA.png"

	e, err := apis.InitApi(app)
	if err != nil {
		t.Fatal(err)
	}

	urls := []string{
		"/api/files/" + fileKey + "?thumb=111x111",
		"/api/files/" + fileKey + "?thumb=111x111", // should still result in single thumb
		"/api/files/" + fileKey + "?thumb=111x222",
		"/api/files/" + fileKey + "?thumb=111x333",
	}

	var wg sync.WaitGroup

	wg.Add(len(urls))

	for _, url := range urls {
		url := url
		go func() {
			defer wg.Done()

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest("GET", url, nil)

			e.ServeHTTP(recorder, req)
		}()
	}

	wg.Wait()

	// ensure that all new requested thumbs were created
	thumbKeys := []string{
		"wsmn24bux7wo113/al1h9ijdeojtsjy/thumbs_300_Jsjq7RdBgA.png/111x111_" + filepath.Base(fileKey),
		"wsmn24bux7wo113/al1h9ijdeojtsjy/thumbs_300_Jsjq7RdBgA.png/111x222_" + filepath.Base(fileKey),
		"wsmn24bux7wo113/al1h9ijdeojtsjy/thumbs_300_Jsjq7RdBgA.png/111x333_" + filepath.Base(fileKey),
	}
	for _, k := range thumbKeys {
		if exists, _ := fsys.Exists(k); !exists {
			t.Fatalf("Missing thumb %q: %v", k, err)
		}
	}
}
