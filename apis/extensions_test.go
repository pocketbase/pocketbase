package apis_test

import (
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/ui"
)

// note: don't run in parallel to avoid conflicts with the ui.DistDirFS nil test
func TestUIExtensions_Mainjs(t *testing.T) {
	successAfterTestFunc := func(t testing.TB, app *tests.TestApp, res *http.Response) {
		expected := "text/javascript"
		if ct := res.Header.Get("content-type"); ct != expected {
			t.Fatalf("Expected response Content-Type %q, got %q", expected, ct)
		}
	}

	oldDistDirFS := ui.DistDirFS

	scenarios := []tests.ApiScenario{
		{
			Name:   "disabled UI",
			Method: http.MethodGet,
			URL:    "/_/extensions.js",
			TestAppFactory: func(t testing.TB) *tests.TestApp {
				app, err := tests.NewTestApp()
				if err != nil {
					t.Fatal(err)
				}

				// simulate no_ui tag (needs to be cleared before the router is initialized)
				ui.DistDirFS = nil

				return app
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				ui.DistDirFS = oldDistDirFS
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "no extensions",
			Method:          http.MethodGet,
			URL:             "/_/extensions.js",
			AfterTestFunc:   successAfterTestFunc,
			ExpectedStatus:  200,
			ExpectedContent: []string{},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "with extensions",
			Method: http.MethodGet,
			URL:    "/_/extensions.js",
			TestAppFactory: func(t testing.TB) *tests.TestApp {
				app, err := tests.NewTestApp()
				if err != nil {
					t.Fatal(err)
				}

				app.OnServe().BindFunc(func(e *core.ServeEvent) error {
					e.UIExtensions = createTestExtensions()
					return e.Next()
				})

				return app
			},
			AfterTestFunc:   successAfterTestFunc,
			ExpectedStatus:  200,
			ExpectedContent: []string{"(function(){ext1_main})();(function(){ext3_main})();"},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// note: don't run in parallel to avoid conflicts with the ui.DistDirFS nil test
func TestUIExtensions_Files(t *testing.T) {
	testAppFactory := func(t testing.TB) *tests.TestApp {
		app, err := tests.NewTestApp()
		if err != nil {
			t.Fatal(err)
		}

		app.OnServe().BindFunc(func(e *core.ServeEvent) error {
			e.UIExtensions = createTestExtensions()
			return e.Next()
		})

		return app
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "no extensions",
			Method:          http.MethodGet,
			URL:             "/_/extensions/ext1/test.txt",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "with missing extension file",
			Method:          http.MethodGet,
			URL:             "/_/extensions/ext1/missing",
			TestAppFactory:  testAppFactory,
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "with existing extension file (ext1)",
			Method:          http.MethodGet,
			URL:             "/_/extensions/ext1/test.txt",
			TestAppFactory:  testAppFactory,
			ExpectedStatus:  200,
			ExpectedContent: []string{"ext1_txt"},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "with existing extension file (extension name escape)",
			Method:          http.MethodGet,
			URL:             "/_/extensions/ext3%20with%20spaces/test.txt",
			TestAppFactory:  testAppFactory,
			ExpectedStatus:  200,
			ExpectedContent: []string{"ext3_txt"},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func createTestExtensions() []core.UIExtension {
	return []core.UIExtension{
		{
			Name: "ext1",
			FS: fstest.MapFS{
				"main.js": &fstest.MapFile{
					Data: []byte("ext1_main"),
				},
				"test.txt": &fstest.MapFile{
					Data: []byte("ext1_txt"),
				},
			},
		},
		{
			Name: "ext2",
			FS: fstest.MapFS{
				"test.txt": &fstest.MapFile{
					Data: []byte("ext2_txt"),
				},
			},
		},
		{
			Name: "ext3 with spaces",
			FS: fstest.MapFS{
				"main.js": &fstest.MapFile{
					Data: []byte("ext3_main"),
				},
				"test.txt": &fstest.MapFile{
					Data: []byte("ext3_txt"),
				},
			},
		},
	}
}
