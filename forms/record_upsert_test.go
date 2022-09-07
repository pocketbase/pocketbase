package forms_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestRecordUpsertPanic1(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewRecordUpsert(nil, nil)
}

func TestRecordUpsertPanic2(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewRecordUpsert(app, nil)
}

func TestNewRecordUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")
	record := models.NewRecord(collection)
	record.SetDataValue("title", "test_value")

	form := forms.NewRecordUpsert(app, record)

	val := form.Data["title"]
	if val != "test_value" {
		t.Errorf("Expected record data to be loaded, got %v", form.Data)
	}
}

func TestRecordUpsertLoadDataUnsupported(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	record, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	testData := "title=test123"

	form := forms.NewRecordUpsert(app, record)
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(testData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	if err := form.LoadData(req); err == nil {
		t.Fatal("Expected LoadData to fail, got nil")
	}
}

func TestRecordUpsertLoadDataJson(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	record, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	testData := map[string]any{
		"id":      "test_id",
		"title":   "test123",
		"unknown": "test456",
		// file fields unset/delete
		"onefile":     nil,
		"manyfiles.0": "",
		"manyfiles.1": "test.png", // should be ignored
		"onlyimages":  nil,
	}

	form := forms.NewRecordUpsert(app, record)
	jsonBody, _ := json.Marshal(testData)
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loadErr := form.LoadData(req)
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	if form.Id != "test_id" {
		t.Fatalf("Expect id field to be %q, got %q", "test_id", form.Id)
	}

	if v, ok := form.Data["title"]; !ok || v != "test123" {
		t.Fatalf("Expect title field to be %q, got %q", "test123", v)
	}

	if v, ok := form.Data["unknown"]; ok {
		t.Fatalf("Didn't expect unknown field to be set, got %v", v)
	}

	onefile, ok := form.Data["onefile"]
	if !ok {
		t.Fatal("Expect onefile field to be set")
	}
	if onefile != "" {
		t.Fatalf("Expect onefile field to be empty string, got %v", onefile)
	}

	manyfiles, ok := form.Data["manyfiles"]
	if !ok || manyfiles == nil {
		t.Fatal("Expect manyfiles field to be set")
	}
	manyfilesRemains := len(list.ToUniqueStringSlice(manyfiles))
	if manyfilesRemains != 1 {
		t.Fatalf("Expect only 1 manyfiles to remain, got \n%v", manyfiles)
	}

	onlyimages := form.Data["onlyimages"]
	if len(list.ToUniqueStringSlice(onlyimages)) != 0 {
		t.Fatalf("Expect onlyimages field to be deleted, got \n%v", onlyimages)
	}
}

func TestRecordUpsertLoadDataMultipart(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	record, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"id":      "test_id",
		"title":   "test123",
		"unknown": "test456",
		// file fields unset/delete
		"onefile":     "",
		"manyfiles.0": "", // delete by index
		"manyfiles.b635c395-6837-49e5-8535-b0a6ebfbdbf3.png": "", // delete by name
		"manyfiles.1": "test.png", // should be ignored
		"onlyimages":  "",
	}, "onlyimages")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, record)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	loadErr := form.LoadData(req)
	if loadErr != nil {
		t.Fatal(loadErr)
	}

	if form.Id != "test_id" {
		t.Fatalf("Expect id field to be %q, got %q", "test_id", form.Id)
	}

	if v, ok := form.Data["title"]; !ok || v != "test123" {
		t.Fatalf("Expect title field to be %q, got %q", "test123", v)
	}

	if v, ok := form.Data["unknown"]; ok {
		t.Fatalf("Didn't expect unknown field to be set, got %v", v)
	}

	onefile, ok := form.Data["onefile"]
	if !ok {
		t.Fatal("Expect onefile field to be set")
	}
	if onefile != "" {
		t.Fatalf("Expect onefile field to be empty string, got %v", onefile)
	}

	manyfiles, ok := form.Data["manyfiles"]
	if !ok || manyfiles == nil {
		t.Fatal("Expect manyfiles field to be set")
	}
	manyfilesRemains := len(list.ToUniqueStringSlice(manyfiles))
	if manyfilesRemains != 0 {
		t.Fatalf("Expect 0 manyfiles to remain, got %v", manyfiles)
	}

	onlyimages, ok := form.Data["onlyimages"]
	if !ok || onlyimages == nil {
		t.Fatal("Expect onlyimages field to be set")
	}
	onlyimagesRemains := len(list.ToUniqueStringSlice(onlyimages))
	expectedRemains := 1 // -2 removed + 1 new upload
	if onlyimagesRemains != expectedRemains {
		t.Fatalf("Expect onlyimages to be %d, got %d (%v)", expectedRemains, onlyimagesRemains, onlyimages)
	}
}

func TestRecordUpsertValidateFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	record, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	// try with invalid test data to check whether the RecordDataValidator is triggered
	formData, mp, err := tests.MockMultipartData(map[string]string{
		"id":      "",
		"unknown": "test456", // should be ignored
		"title":   "a",
		"onerel":  "00000000-84ab-4057-a592-4604a731f78f",
	}, "manyfiles", "manyfiles")
	if err != nil {
		t.Fatal(err)
	}

	expectedErrors := []string{"title", "onerel", "manyfiles"}

	form := forms.NewRecordUpsert(app, record)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

	result := form.Validate()

	// parse errors
	errs, ok := result.(validation.Errors)
	if !ok && result != nil {
		t.Fatalf("Failed to parse errors %v", result)
	}

	// check errors
	if len(errs) > len(expectedErrors) {
		t.Fatalf("Expected error keys %v, got %v", expectedErrors, errs)
	}
	for _, k := range expectedErrors {
		if _, ok := errs[k]; !ok {
			t.Errorf("Missing expected error key %q in %v", k, errs)
		}
	}
}

func TestRecordUpsertValidateSuccess(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	record, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"id":      record.Id,
		"unknown": "test456", // should be ignored
		"title":   "abc",
		"onerel":  "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2",
	}, "manyfiles", "onefile")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, record)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

	result := form.Validate()
	if result != nil {
		t.Fatal(result)
	}
}

func TestRecordUpsertDrySubmitFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	recordBefore, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title":  "a",
		"onerel": "00000000-84ab-4057-a592-4604a731f78f",
	})
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

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
	recordAfter, err := app.Dao().FindFirstRecordByData(collection, "id", recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if recordAfter.GetStringDataValue("title") == "a" {
		t.Fatalf("Expected record.title to be %v, got %v", recordAfter.GetStringDataValue("title"), "a")
	}

	if recordAfter.GetStringDataValue("onerel") == "00000000-84ab-4057-a592-4604a731f78f" {
		t.Fatalf("Expected record.onerel to be %s, got %s", recordBefore.GetStringDataValue("onerel"), recordAfter.GetStringDataValue("onerel"))
	}
}

func TestRecordUpsertDrySubmitSuccess(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	recordBefore, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title":   "dry_test",
		"onefile": "",
	}, "manyfiles")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

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
	// ---
	recordAfter, err := app.Dao().FindFirstRecordByData(collection, "id", recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if recordAfter.GetStringDataValue("title") == "dry_test" {
		t.Fatalf("Expected record.title to be %v, got %v", recordAfter.GetStringDataValue("title"), "dry_test")
	}
	if recordAfter.GetStringDataValue("onefile") == "" {
		t.Fatal("Expected record.onefile to be set, got empty string")
	}

	// file wasn't removed
	if !hasRecordFile(app, recordAfter, recordAfter.GetStringDataValue("onefile")) {
		t.Fatal("onefile file should not have been deleted")
	}
}

func TestRecordUpsertSubmitFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	recordBefore, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title":   "a",
		"onefile": "",
	})
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

	interceptorCalls := 0
	interceptor := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			interceptorCalls++
			return next()
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
	recordAfter, err := app.Dao().FindFirstRecordByData(collection, "id", recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if recordAfter.GetStringDataValue("title") == "a" {
		t.Fatalf("Expected record.title to be %v, got %v", recordAfter.GetStringDataValue("title"), "a")
	}

	if recordAfter.GetStringDataValue("onefile") == "" {
		t.Fatal("Expected record.onefile to be set, got empty string")
	}

	// file wasn't removed
	if !hasRecordFile(app, recordAfter, recordAfter.GetStringDataValue("onefile")) {
		t.Fatal("onefile file should not have been deleted")
	}
}

func TestRecordUpsertSubmitSuccess(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	recordBefore, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	formData, mp, err := tests.MockMultipartData(map[string]string{
		"title":      "test_save",
		"onefile":    "",
		"onlyimages": "",
	}, "manyfiles.1", "manyfiles") // replace + new file
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

	interceptorCalls := 0
	interceptor := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			interceptorCalls++
			return next()
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
	recordAfter, err := app.Dao().FindFirstRecordByData(collection, "id", recordBefore.Id)
	if err != nil {
		t.Fatal(err)
	}

	if recordAfter.GetStringDataValue("title") != "test_save" {
		t.Fatalf("Expected record.title to be %v, got %v", recordAfter.GetStringDataValue("title"), "test_save")
	}

	if hasRecordFile(app, recordAfter, recordAfter.GetStringDataValue("onefile")) {
		t.Fatal("Expected record.onefile to be deleted")
	}

	onlyimages := (recordAfter.GetStringSliceDataValue("onlyimages"))
	if len(onlyimages) != 0 {
		t.Fatalf("Expected all onlyimages files to be deleted, got %d (%v)", len(onlyimages), onlyimages)
	}

	manyfiles := (recordAfter.GetStringSliceDataValue("manyfiles"))
	if len(manyfiles) != 3 {
		t.Fatalf("Expected 3 manyfiles, got %d (%v)", len(manyfiles), manyfiles)
	}
	for _, f := range manyfiles {
		if !hasRecordFile(app, recordAfter, f) {
			t.Fatalf("Expected file %q to exist", f)
		}
	}
}

func TestRecordUpsertSubmitInterceptors(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")
	record, err := app.Dao().FindFirstRecordByData(collection, "id", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, record)
	form.Data["title"] = "test_new"

	testErr := errors.New("test_error")
	interceptorRecordTitle := ""

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			interceptor1Called = true
			return next()
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			interceptorRecordTitle = record.GetStringDataValue("title") // to check if the record was filled
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

	if interceptorRecordTitle != form.Data["title"].(string) {
		t.Fatalf("Expected the form model to be filled before calling the interceptors")
	}
}

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

func TestRecordUpsertWithCustomId(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo3")
	existingRecord, err := app.Dao().FindFirstRecordByData(collection, "id", "2c542824-9de1-42fe-8924-e57c86267760")
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
		form.LoadData(req)

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
			_, err := app.Dao().FindRecordById(collection, id, nil)
			if err != nil {
				t.Errorf("[%s] Expected to find record with id %s, got %v", scenario.name, id, err)
			}
		}
	}
}
