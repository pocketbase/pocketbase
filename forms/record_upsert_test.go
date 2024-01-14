package forms_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/types"
)

func hasRecordFile(app core.App, record *models.Record, filename string) bool {
	fs, _ := app.NewFilesystem()
	defer fs.Close()

	fileKey := filepath.Join(
		record.Collection().Id,
		record.Id,
		filename,
	)

	exists, _ := fs.Exists(fileKey)

	return exists
}

func TestNewRecordUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo2")
	record := models.NewRecord(collection)
	record.Set("title", "test_value")

	form := forms.NewRecordUpsert(app, record)

	val := form.Data()["title"]
	if val != "test_value" {
		t.Errorf("Expected record data to be loaded, got %v", form.Data())
	}
}

func TestRecordUpsertLoadRequestUnsupported(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.Dao().FindRecordById("demo2", "0yxhwia2amd8gec")
	if err != nil {
		t.Fatal(err)
	}

	testData := "title=test123"

	form := forms.NewRecordUpsert(app, record)
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(testData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if err := form.LoadRequest(req, ""); err == nil {
		t.Fatal("Expected LoadRequest to fail, got nil")
	}
}

func TestRecordUpsertLoadRequestJson(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.Dao().FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	testData := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"id":      "test_id",
				"text":    "test123",
				"unknown": "test456",
				// file fields unset/delete
				"file_one":                     nil,
				"file_many.0":                  "",                                                     // delete by index
				"file_many-":                   []string{"test_MaWC6mWyrP.txt", "test_tC1Yc87DfC.txt"}, // multiple delete with modifier
				"file_many.300_WlbFWSGmW9.png": nil,                                                    // delete by filename
				"file_many.2":                  "test.png",                                             // should be ignored
			},
		},
	}

	form := forms.NewRecordUpsert(app, record)
	jsonBody, _ := json.Marshal(testData)
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loadErr := form.LoadRequest(req, "a.b")
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	if form.Id != "test_id" {
		t.Fatalf("Expect id field to be %q, got %q", "test_id", form.Id)
	}

	if v, ok := form.Data()["text"]; !ok || v != "test123" {
		t.Fatalf("Expect title field to be %q, got %q", "test123", v)
	}

	if v, ok := form.Data()["unknown"]; ok {
		t.Fatalf("Didn't expect unknown field to be set, got %v", v)
	}

	fileOne, ok := form.Data()["file_one"]
	if !ok {
		t.Fatal("Expect file_one field to be set")
	}
	if fileOne != "" {
		t.Fatalf("Expect file_one field to be empty string, got %v", fileOne)
	}

	fileMany, ok := form.Data()["file_many"]
	if !ok || fileMany == nil {
		t.Fatal("Expect file_many field to be set")
	}
	manyfilesRemains := len(list.ToUniqueStringSlice(fileMany))
	if manyfilesRemains != 1 {
		t.Fatalf("Expect only 1 file_many to remain, got \n%v", fileMany)
	}
}

func TestRecordUpsertLoadRequestMultipart(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.Dao().FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"a.b.id":                       "test_id",
		"a.b.text":                     "test123",
		"a.b.unknown":                  "test456",
		"a.b." + rest.MultipartJsonKey: `{"json":["a","b"],"email":"test3@example.com"}`,
		// file fields unset/delete
		"a.b.file_one-":                    "test_d61b33QdDU.txt", // delete with modifier
		"a.b.file_many.0":                  "",                    // delete by index
		"a.b.file_many-":                   "test_tC1Yc87DfC.txt", // delete with modifier
		"a.b.file_many.300_WlbFWSGmW9.png": "",                    // delete by filename
		"a.b.file_many.2":                  "test.png",            // should be ignored
	}, "a.b.file_many")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, record)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	loadErr := form.LoadRequest(req, "a.b")
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	if form.Id != "test_id" {
		t.Fatalf("Expect id field to be %q, got %q", "test_id", form.Id)
	}

	if v, ok := form.Data()["text"]; !ok || v != "test123" {
		t.Fatalf("Expect text field to be %q, got %q", "test123", v)
	}

	if v, ok := form.Data()["unknown"]; ok {
		t.Fatalf("Didn't expect unknown field to be set, got %v", v)
	}

	if v, ok := form.Data()["email"]; !ok || v != "test3@example.com" {
		t.Fatalf("Expect email field to be %q, got %q", "test3@example.com", v)
	}

	rawJsonValue, ok := form.Data()["json"].(types.JsonRaw)
	if !ok {
		t.Fatal("Expect json field to be set")
	}
	expectedJsonValue := `["a","b"]`
	if rawJsonValue.String() != expectedJsonValue {
		t.Fatalf("Expect json field %v, got %v", expectedJsonValue, rawJsonValue)
	}

	fileOne, ok := form.Data()["file_one"]
	if !ok {
		t.Fatal("Expect file_one field to be set")
	}
	if fileOne != "" {
		t.Fatalf("Expect file_one field to be empty string, got %v", fileOne)
	}

	fileMany, ok := form.Data()["file_many"]
	if !ok || fileMany == nil {
		t.Fatal("Expect file_many field to be set")
	}
	manyfilesRemains := len(list.ToUniqueStringSlice(fileMany))
	expectedRemains := 3 // 5 old; 3 deleted and 1 new uploaded
	if manyfilesRemains != expectedRemains {
		t.Fatalf("Expect file_many to be %d, got %d (%v)", expectedRemains, manyfilesRemains, fileMany)
	}
}

func TestRecordUpsertLoadData(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.Dao().FindRecordById("demo2", "llvuca81nly1qls")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, record)

	loadErr := form.LoadData(map[string]any{
		"title":  "test_new",
		"active": true,
	})
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	if v, ok := form.Data()["title"]; !ok || v != "test_new" {
		t.Fatalf("Expect title field to be %v, got %v", "test_new", v)
	}

	if v, ok := form.Data()["active"]; !ok || v != true {
		t.Fatalf("Expect active field to be %v, got %v", true, v)
	}
}

func TestRecordUpsertDrySubmitFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo1")
	recordBefore, err := app.Dao().FindRecordById(collection.Id, "al1h9ijdeojtsjy")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title":   "abc",
		"rel_one": "missing",
	})
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadRequest(req, "")

	callbackCalls := 0

	// ensure that validate is triggered
	// ---
	result := form.DrySubmit(func(txDao *daos.Dao) error {
		callbackCalls++
		return nil
	})
	if result == nil {
		t.Fatal("Expected error, got nil")
	}
	if callbackCalls != 0 {
		t.Fatalf("Expected callbackCalls to be 0, got %d", callbackCalls)
	}

	// ensure that the record changes weren't persisted
	// ---
	recordAfter, err := app.Dao().FindRecordById(collection.Id, recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if recordAfter.GetString("title") == "abc" {
		t.Fatalf("Expected record.title to be %v, got %v", recordAfter.GetString("title"), "abc")
	}

	if recordAfter.GetString("rel_one") == "missing" {
		t.Fatalf("Expected record.rel_one to be %s, got %s", recordBefore.GetString("rel_one"), "missing")
	}
}

func TestRecordUpsertDrySubmitSuccess(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo1")
	recordBefore, err := app.Dao().FindRecordById(collection.Id, "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title":    "dry_test",
		"file_one": "",
	}, "file_many")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadRequest(req, "")

	callbackCalls := 0

	result := form.DrySubmit(func(txDao *daos.Dao) error {
		callbackCalls++
		return nil
	})
	if result != nil {
		t.Fatalf("Expected nil, got error %v", result)
	}

	// ensure callback was called
	if callbackCalls != 1 {
		t.Fatalf("Expected callbackCalls to be 1, got %d", callbackCalls)
	}

	// ensure that the record changes weren't persisted
	recordAfter, err := app.Dao().FindRecordById(collection.Id, recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}
	if recordAfter.GetString("title") == "dry_test" {
		t.Fatalf("Expected record.title to be %v, got %v", recordAfter.GetString("title"), "dry_test")
	}
	if recordAfter.GetString("file_one") == "" {
		t.Fatal("Expected record.file_one to not be changed, got empty string")
	}

	// file wasn't removed
	if !hasRecordFile(app, recordAfter, recordAfter.GetString("file_one")) {
		t.Fatal("file_one file should not have been deleted")
	}
}

func TestRecordUpsertDrySubmitWithNestedTx(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo1")
	recordBefore, err := app.Dao().FindRecordById(collection.Id, "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title": "dry_test",
	})
	if err != nil {
		t.Fatal(err)
	}

	txErr := app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		form := forms.NewRecordUpsert(app, recordBefore)
		form.SetDao(txDao)
		req := httptest.NewRequest(http.MethodGet, "/", formData)
		req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
		form.LoadRequest(req, "")

		callbackCalls := 0

		result := form.DrySubmit(func(innerTxDao *daos.Dao) error {
			callbackCalls++
			return nil
		})
		if result != nil {
			t.Fatalf("Expected nil, got error %v", result)
		}

		// ensure callback was called
		if callbackCalls != 1 {
			t.Fatalf("Expected callbackCalls to be 1, got %d", callbackCalls)
		}

		// ensure that the original txDao can still be used after the DrySubmit rollback
		if _, err := txDao.FindRecordById(collection.Id, recordBefore.Id); err != nil {
			t.Fatalf("Expected the dry submit rollback to not affect the outer tx context, got %v", err)
		}

		// ensure that the record changes weren't persisted
		recordAfter, err := app.Dao().FindRecordById(collection.Id, recordBefore.Id)
		if err != nil {
			t.Fatal(err)
		}
		if recordAfter.GetString("title") == "dry_test" {
			t.Fatalf("Expected record.title to be %v, got %v", recordBefore.GetString("title"), "dry_test")
		}

		return nil
	})
	if txErr != nil {
		t.Fatalf("Nested transactions failure: %v", txErr)
	}
}

func TestRecordUpsertSubmitFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	recordBefore, err := app.Dao().FindRecordById(collection.Id, "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"text":       "abc",
		"bool":       "false",
		"select_one": "invalid",
		"file_many":  "invalid",
		"email":      "invalid",
	}, "file_one")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadRequest(req, "")

	interceptorCalls := 0
	interceptor := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(r *models.Record) error {
			interceptorCalls++
			return next(r)
		}
	}

	// ensure that validate is triggered
	// ---
	result := form.Submit(interceptor)
	if result == nil {
		t.Fatal("Expected error, got nil")
	}

	// check interceptor calls
	// ---
	if interceptorCalls != 0 {
		t.Fatalf("Expected interceptor to be called 0 times, got %d", interceptorCalls)
	}

	// ensure that the record changes weren't persisted
	// ---
	recordAfter, err := app.Dao().FindRecordById(collection.Id, recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if v := recordAfter.Get("text"); v == "abc" {
		t.Fatalf("Expected record.text not to change, got %v", v)
	}
	if v := recordAfter.Get("bool"); v == false {
		t.Fatalf("Expected record.bool not to change, got %v", v)
	}
	if v := recordAfter.Get("select_one"); v == "invalid" {
		t.Fatalf("Expected record.select_one not to change, got %v", v)
	}
	if v := recordAfter.Get("email"); v == "invalid" {
		t.Fatalf("Expected record.email not to change, got %v", v)
	}
	if v := recordAfter.GetStringSlice("file_many"); len(v) != 5 {
		t.Fatalf("Expected record.file_many not to change, got %v", v)
	}

	// ensure the files weren't removed
	for _, f := range recordAfter.GetStringSlice("file_many") {
		if !hasRecordFile(app, recordAfter, f) {
			t.Fatal("file_many file should not have been deleted")
		}
	}
}

func TestRecordUpsertSubmitSuccess(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo1")
	recordBefore, err := app.Dao().FindRecordById(collection.Id, "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"text":       "test_save",
		"bool":       "true",
		"select_one": "optionA",
		"file_one":   "",
	}, "file_many.1", "file_many") // replace + new file
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadRequest(req, "")

	interceptorCalls := 0
	interceptor := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(r *models.Record) error {
			interceptorCalls++
			return next(r)
		}
	}

	result := form.Submit(interceptor)
	if result != nil {
		t.Fatalf("Expected nil, got error %v", result)
	}

	// check interceptor calls
	// ---
	if interceptorCalls != 1 {
		t.Fatalf("Expected interceptor to be called 1 time, got %d", interceptorCalls)
	}

	// ensure that the record changes were persisted
	// ---
	recordAfter, err := app.Dao().FindRecordById(collection.Id, recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if v := recordAfter.GetString("text"); v != "test_save" {
		t.Fatalf("Expected record.text to be %v, got %v", v, "test_save")
	}

	if hasRecordFile(app, recordAfter, recordAfter.GetString("file_one")) {
		t.Fatal("Expected record.file_one to be deleted")
	}

	fileMany := (recordAfter.GetStringSlice("file_many"))
	if len(fileMany) != 6 { // 1 replace + 1 new
		t.Fatalf("Expected 6 record.file_many, got %d (%v)", len(fileMany), fileMany)
	}
	for _, f := range fileMany {
		if !hasRecordFile(app, recordAfter, f) {
			t.Fatalf("Expected file %q to exist", f)
		}
	}
}

func TestRecordUpsertSubmitInterceptors(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo3")
	record, err := app.Dao().FindRecordById(collection.Id, "mk5fmymtx4wsprk")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, record)
	form.Data()["title"] = "test_new"

	testErr := errors.New("test_error")
	interceptorRecordTitle := ""

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(r *models.Record) error {
			interceptor1Called = true
			return next(r)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(r *models.Record) error {
			interceptorRecordTitle = record.GetString("title") // to check if the record was filled
			interceptor2Called = true
			return testErr
		}
	}

	submitErr := form.Submit(interceptor1, interceptor2)
	if submitErr != testErr {
		t.Fatalf("Expected submitError %v, got %v", testErr, submitErr)
	}

	if !interceptor1Called {
		t.Fatalf("Expected interceptor1 to be called")
	}

	if !interceptor2Called {
		t.Fatalf("Expected interceptor2 to be called")
	}

	if interceptorRecordTitle != form.Data()["title"].(string) {
		t.Fatalf("Expected the form model to be filled before calling the interceptors")
	}
}

func TestRecordUpsertWithCustomId(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}

	existingRecord, err := app.Dao().FindRecordById(collection.Id, "mk5fmymtx4wsprk")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		data        map[string]string
		record      *models.Record
		expectError bool
	}{
		{
			"empty data",
			map[string]string{},
			models.NewRecord(collection),
			false,
		},
		{
			"empty id",
			map[string]string{"id": ""},
			models.NewRecord(collection),
			false,
		},
		{
			"id < 15 chars",
			map[string]string{"id": "a23"},
			models.NewRecord(collection),
			true,
		},
		{
			"id > 15 chars",
			map[string]string{"id": "a234567890123456"},
			models.NewRecord(collection),
			true,
		},
		{
			"id = 15 chars (invalid chars)",
			map[string]string{"id": "a@3456789012345"},
			models.NewRecord(collection),
			true,
		},
		{
			"id = 15 chars (valid chars)",
			map[string]string{"id": "a23456789012345"},
			models.NewRecord(collection),
			false,
		},
		{
			"changing the id of an existing record",
			map[string]string{"id": "b23456789012345"},
			existingRecord,
			true,
		},
		{
			"using the same existing record id",
			map[string]string{"id": existingRecord.Id},
			existingRecord,
			false,
		},
		{
			"skipping the id for existing record",
			map[string]string{},
			existingRecord,
			false,
		},
	}

	for _, scenario := range scenarios {
		formData, mp, err := tests.MockMultipartData(scenario.data)
		if err != nil {
			t.Fatal(err)
		}

		form := forms.NewRecordUpsert(app, scenario.record)
		req := httptest.NewRequest(http.MethodGet, "/", formData)
		req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
		form.LoadRequest(req, "")

		dryErr := form.DrySubmit(nil)
		hasDryErr := dryErr != nil

		submitErr := form.Submit()
		hasSubmitErr := submitErr != nil

		if hasDryErr != hasSubmitErr {
			t.Errorf("[%s] Expected hasDryErr and hasSubmitErr to have the same value, got %v vs %v", scenario.name, hasDryErr, hasSubmitErr)
		}

		if hasSubmitErr != scenario.expectError {
			t.Errorf("[%s] Expected hasSubmitErr to be %v, got %v (%v)", scenario.name, scenario.expectError, hasSubmitErr, submitErr)
		}

		if id, ok := scenario.data["id"]; ok && id != "" && !hasSubmitErr {
			_, err := app.Dao().FindRecordById(collection.Id, id)
			if err != nil {
				t.Errorf("[%s] Expected to find record with id %s, got %v", scenario.name, id, err)
			}
		}
	}
}

func TestRecordUpsertAuthRecord(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		existingId   string
		data         map[string]any
		manageAccess bool
		expectError  bool
	}{
		{
			"empty create data",
			"",
			map[string]any{},
			false,
			true,
		},
		{
			"empty update data",
			"4q1xlclmfloku33",
			map[string]any{},
			false,
			false,
		},
		{
			"minimum valid create data",
			"",
			map[string]any{
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			false,
			false,
		},
		{
			"create with all allowed auth fields",
			"",
			map[string]any{
				"username":        "test_new-a.b",
				"email":           "test_new@example.com",
				"emailVisibility": true,
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			false,
			false,
		},

		// username
		{
			"invalid username characters",
			"",
			map[string]any{
				"username":        "test abc!@#",
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			false,
			true,
		},
		{
			"invalid username length (less than 3)",
			"",
			map[string]any{
				"username":        "ab",
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			false,
			true,
		},
		{
			"invalid username length (more than 150)",
			"",
			map[string]any{
				"username":        strings.Repeat("a", 151),
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			false,
			true,
		},

		// verified
		{
			"try to set verified without managed access",
			"",
			map[string]any{
				"verified":        true,
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			false,
			true,
		},
		{
			"try to update verified without managed access",
			"4q1xlclmfloku33",
			map[string]any{
				"verified": true,
			},
			false,
			true,
		},
		{
			"set verified with managed access",
			"",
			map[string]any{
				"verified":        true,
				"password":        "12345678",
				"passwordConfirm": "12345678",
			},
			true,
			false,
		},
		{
			"update verified with managed access",
			"4q1xlclmfloku33",
			map[string]any{
				"verified": true,
			},
			true,
			false,
		},

		// email
		{
			"try to update email without managed access",
			"4q1xlclmfloku33",
			map[string]any{
				"email": "test_update@example.com",
			},
			false,
			true,
		},
		{
			"update email with managed access",
			"4q1xlclmfloku33",
			map[string]any{
				"email": "test_update@example.com",
			},
			true,
			false,
		},

		// password
		{
			"trigger the password validations if only oldPassword is set",
			"4q1xlclmfloku33",
			map[string]any{
				"oldPassword": "1234567890",
			},
			false,
			true,
		},
		{
			"trigger the password validations if only passwordConfirm is set",
			"4q1xlclmfloku33",
			map[string]any{
				"passwordConfirm": "1234567890",
			},
			false,
			true,
		},
		{
			"try to update password without managed access",
			"4q1xlclmfloku33",
			map[string]any{
				"password":        "1234567890",
				"passwordConfirm": "1234567890",
			},
			false,
			true,
		},
		{
			"update password without managed access but with oldPassword",
			"4q1xlclmfloku33",
			map[string]any{
				"oldPassword":     "1234567890",
				"password":        "1234567890",
				"passwordConfirm": "1234567890",
			},
			false,
			false,
		},
		{
			"update email with managed access (without oldPassword)",
			"4q1xlclmfloku33",
			map[string]any{
				"password":        "1234567890",
				"passwordConfirm": "1234567890",
			},
			true,
			false,
		},
	}

	for _, s := range scenarios {
		collection, err := app.Dao().FindCollectionByNameOrId("users")
		if err != nil {
			t.Fatal(err)
		}

		record := models.NewRecord(collection)
		if s.existingId != "" {
			var err error
			record, err = app.Dao().FindRecordById(collection.Id, s.existingId)
			if err != nil {
				t.Errorf("[%s] Failed to fetch auth record with id %s", s.name, s.existingId)
				continue
			}
		}

		form := forms.NewRecordUpsert(app, record)
		form.SetFullManageAccess(s.manageAccess)
		if err := form.LoadData(s.data); err != nil {
			t.Errorf("[%s] Failed to load form data", s.name)
			continue
		}

		submitErr := form.Submit()

		hasErr := submitErr != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, submitErr)
		}

		if !hasErr && record.Username() == "" {
			t.Errorf("[%s] Expected username to be set, got empty string: \n%v", s.name, record)
		}
	}
}

func TestRecordUpsertUniqueValidator(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create a dummy collection
	collection := &models.Collection{
		Name: "test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Type: "text",
				Name: "fieldA",
			},
			&schema.SchemaField{
				Type: "text",
				Name: "fieldB",
			},
			&schema.SchemaField{
				Type: "text",
				Name: "fieldC",
			},
		),
		Indexes: types.JsonArray[string]{
			// the field case shouldn't matter
			"create unique index unique_single_idx on test (fielda)",
			"create unique index unique_combined_idx on test (fieldb, FIELDC)",
		},
	}
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	dummyRecord := models.NewRecord(collection)
	dummyRecord.Set("fieldA", "a")
	dummyRecord.Set("fieldB", "b")
	dummyRecord.Set("fieldC", "c")
	if err := app.Dao().SaveRecord(dummyRecord); err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name           string
		data           map[string]any
		expectedErrors []string
	}{
		{
			"duplicated unique value",
			map[string]any{
				"fieldA": "a",
			},
			[]string{"fieldA"},
		},
		{
			"duplicated combined unique value",
			map[string]any{
				"fieldB": "b",
				"fieldC": "c",
			},
			[]string{"fieldB", "fieldC"},
		},
		{
			"non-duplicated unique value",
			map[string]any{
				"fieldA": "a2",
			},
			nil,
		},
		{
			"non-duplicated combined unique value",
			map[string]any{
				"fieldB": "b",
				"fieldC": "d",
			},
			nil,
		},
	}

	for _, s := range scenarios {
		record := models.NewRecord(collection)

		form := forms.NewRecordUpsert(app, record)
		if err := form.LoadData(s.data); err != nil {
			t.Errorf("[%s] Failed to load form data", s.name)
			continue
		}

		result := form.Submit()

		// parse errors
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("[%s] Failed to parse errors %v", s.name, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", s.name, s.expectedErrors, errs)
			continue
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", s.name, k, errs)
				continue
			}
		}
	}
}

func TestRecordUpsertAddAndRemoveFiles(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	recordBefore, err := app.Dao().FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	// create test temp files
	tempDir := filepath.Join(app.DataDir(), "temp")
	if err := os.MkdirAll(app.DataDir(), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	tmpFile, _ := os.CreateTemp(os.TempDir(), "tmpfile1-*.txt")
	tmpFile.Close()

	form := forms.NewRecordUpsert(app, recordBefore)

	f1, err := filesystem.NewFileFromPath(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	f2, err := filesystem.NewFileFromPath(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	f3, err := filesystem.NewFileFromPath(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	removed0 := "test_d61b33QdDU.txt" // replaced
	removed1 := "300_WlbFWSGmW9.png"
	removed2 := "logo_vcfJJG5TAh.svg"

	form.AddFiles("file_one", f1) // should replace the existin file

	form.AddFiles("file_many", f2, f3) // should append

	form.RemoveFiles("file_many", removed1, removed2) // should remove

	filesToUpload := form.FilesToUpload()
	if v, ok := filesToUpload["file_one"]; !ok || len(v) != 1 {
		t.Fatalf("Expected filesToUpload[file_one] to have exactly 1 file, got %v", v)
	}
	if v, ok := filesToUpload["file_many"]; !ok || len(v) != 2 {
		t.Fatalf("Expected filesToUpload[file_many] to have exactly 2 file, got %v", v)
	}

	filesToDelete := form.FilesToDelete()
	if len(filesToDelete) != 3 {
		t.Fatalf("Expected exactly 2 file to delete, got %v", filesToDelete)
	}
	for _, f := range []string{removed0, removed1, removed2} {
		if !list.ExistInSlice(f, filesToDelete) {
			t.Fatalf("Missing file %q from filesToDelete %v", f, filesToDelete)
		}
	}

	if err := form.Submit(); err != nil {
		t.Fatalf("Failed to submit the RecordUpsert form, got %v", err)
	}

	recordAfter, err := app.Dao().FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	// ensure files deletion
	if hasRecordFile(app, recordAfter, removed0) {
		t.Fatalf("Expected the old file_one file to be deleted")
	}
	if hasRecordFile(app, recordAfter, removed1) {
		t.Fatalf("Expected %s to be deleted", removed1)
	}
	if hasRecordFile(app, recordAfter, removed2) {
		t.Fatalf("Expected %s to be deleted", removed2)
	}

	fileOne := recordAfter.GetStringSlice("file_one")
	if len(fileOne) == 0 {
		t.Fatalf("Expected new file_one file to be uploaded")
	}

	fileMany := recordAfter.GetStringSlice("file_many")
	if len(fileMany) != 5 {
		t.Fatalf("Expected file_many to be 5, got %v", fileMany)
	}
}

func TestRecordUpsertUploadFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}

	testDaos := []*daos.Dao{
		app.Dao(),                // with hooks
		daos.New(app.Dao().DB()), // without hooks
	}

	for i, dao := range testDaos {
		// create with invalid file
		{
			prefix := fmt.Sprintf("%d-create", i)

			new := models.NewRecord(collection)
			new.Id = "123456789012341"

			form := forms.NewRecordUpsert(app, new)
			form.SetDao(dao)
			form.LoadData(map[string]any{"title": "new_test"})
			form.AddFiles("files", &filesystem.File{Reader: &filesystem.PathReader{Path: "/tmp/__missing__"}})

			if err := form.Submit(); err == nil {
				t.Fatalf("[%s] Expected error, got nil", prefix)
			}

			if r, err := app.Dao().FindRecordById(collection.Id, new.Id); err == nil {
				t.Fatalf("[%s] Expected the inserted record to be deleted, found \n%v", prefix, r.PublicExport())
			}
		}

		// update with invalid file
		{
			prefix := fmt.Sprintf("%d-update", i)

			record, err := app.Dao().FindRecordById(collection.Id, "1tmknxy2868d869")
			if err != nil {
				t.Fatal(err)
			}

			form := forms.NewRecordUpsert(app, record)
			form.SetDao(dao)
			form.LoadData(map[string]any{"title": "update_test"})
			form.AddFiles("files", &filesystem.File{Reader: &filesystem.PathReader{Path: "/tmp/__missing__"}})

			if err := form.Submit(); err == nil {
				t.Fatalf("[%s] Expected error, got nil", prefix)
			}

			if r, _ := app.Dao().FindRecordById(collection.Id, record.Id); r == nil || r.GetString("title") == "update_test" {
				t.Fatalf("[%s] Expected the record changes to be reverted, got \n%v", prefix, r.PublicExport())
			}
		}
	}
}
