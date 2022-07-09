package forms_test

import (
	"bytes"
	"encoding/json"
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

func TestNewRecordUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")
	record := models.NewRecord(collection)
	record.SetDataValue("title", "test_value")

	form := forms.NewRecordUpsert(app, record)

	val := form.Data["title"]
	if val != "test_value" {
		t.Errorf("Expected record data to be load, got %v", form.Data)
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
		"title":   "test123",
		"unknown": "test456",
		// file fields unset/delete
		"onefile":     nil,
		"manyfiles.0": "",
		"manyfiles.1": "test.png", // should be ignored
		"onlyimages":  nil,        // should be ignored
	}

	form := forms.NewRecordUpsert(app, record)
	jsonBody, _ := json.Marshal(testData)
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loadErr := form.LoadData(req)
	if loadErr != nil {
		t.Fatal(loadErr)
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
	if onefile != nil {
		t.Fatalf("Expect onefile field to be nil, got %v", onefile)
	}

	manyfiles, ok := form.Data["manyfiles"]
	if !ok || manyfiles == nil {
		t.Fatal("Expect manyfiles field to be set")
	}
	manyfilesRemains := len(list.ToUniqueStringSlice(manyfiles))
	if manyfilesRemains != 1 {
		t.Fatalf("Expect only 1 manyfiles to remain, got %v", manyfiles)
	}

	// cannot reset multiple file upload field with just using the field name
	onlyimages, ok := form.Data["onlyimages"]
	if !ok || onlyimages == nil {
		t.Fatal("Expect onlyimages field to be set and not be altered")
	}
	onlyimagesRemains := len(list.ToUniqueStringSlice(onlyimages))
	expectedRemains := 2 // 2 existing
	if onlyimagesRemains != expectedRemains {
		t.Fatalf("Expect onlyimages to be %d, got %d (%v)", expectedRemains, onlyimagesRemains, onlyimages)
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
		"title":   "test123",
		"unknown": "test456",
		// file fields unset/delete
		"onefile":     "",
		"manyfiles.0": "",
		"manyfiles.1": "test.png", // should be ignored
		"onlyimages":  "",         // should be ignored
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
	if onefile != nil {
		t.Fatalf("Expect onefile field to be nil, got %v", onefile)
	}

	manyfiles, ok := form.Data["manyfiles"]
	if !ok || manyfiles == nil {
		t.Fatal("Expect manyfiles field to be set")
	}
	manyfilesRemains := len(list.ToUniqueStringSlice(manyfiles))
	if manyfilesRemains != 1 {
		t.Fatalf("Expect only 1 manyfiles to remain, got %v", manyfiles)
	}

	onlyimages, ok := form.Data["onlyimages"]
	if !ok || onlyimages == nil {
		t.Fatal("Expect onlyimages field to be set and not be altered")
	}
	onlyimagesRemains := len(list.ToUniqueStringSlice(onlyimages))
	expectedRemains := 3 // 2 existing + 1 new upload
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

	// ensure that validate is triggered
	// ---
	result := form.Submit()
	if result == nil {
		t.Fatal("Expected error, got nil")
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
		"title":   "test_save",
		"onefile": "",
	}, "manyfiles.1", "manyfiles") // replace + new file
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(app, recordBefore)
	req := httptest.NewRequest(http.MethodGet, "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	form.LoadData(req)

	result := form.Submit()
	if result != nil {
		t.Fatalf("Expected nil, got error %v", result)
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
