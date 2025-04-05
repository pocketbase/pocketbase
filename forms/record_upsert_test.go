package forms_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestRecordUpsertLoad(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	demo1Col, err := testApp.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	usersCol, err := testApp.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name                 string
		data                 map[string]any
		record               *core.Record
		managerAccessLevel   bool
		superuserAccessLevel bool
		expected             []string
		notExpected          []string
	}{
		{
			name: "base collection record",
			data: map[string]any{
				"text":         "test_text",
				"custom":       "123",                          // should be ignored
				"number":       "456",                          // should be normalized by the setter
				"select_many+": []string{"optionB", "optionC"}, // test modifier fields
				"created":      "2022-01:01 10:00:00.000Z",     // should be ignored
				// ignore special auth fields
				"oldPassword":     "123",
				"password":        "456",
				"passwordConfirm": "789",
			},
			record: core.NewRecord(demo1Col),
			expected: []string{
				`"text":"test_text"`,
				`"number":456`,
				`"select_many":["optionB","optionC"]`,
				`"created":""`,
				`"updated":""`,
				`"json":null`,
			},
			notExpected: []string{
				`"custom"`,
				`"password"`,
				`"oldPassword"`,
				`"passwordConfirm"`,
				`"select_many-"`,
				`"select_many+"`,
			},
		},
		{
			name: "auth collection record",
			data: map[string]any{
				"email": "test@example.com",
				// special auth fields
				"oldPassword":     "123",
				"password":        "456",
				"passwordConfirm": "789",
			},
			record: core.NewRecord(usersCol),
			expected: []string{
				`"email":"test@example.com"`,
				`"password":"456"`,
			},
			notExpected: []string{
				`"oldPassword"`,
				`"passwordConfirm"`,
			},
		},
		{
			name: "hidden fields (manager)",
			data: map[string]any{
				"email":    "test@example.com",
				"tokenKey": "abc", // should be ignored
				// special auth fields
				"password":        "456",
				"oldPassword":     "123",
				"passwordConfirm": "789",
			},
			managerAccessLevel: true,
			record:             core.NewRecord(usersCol),
			expected: []string{
				`"email":"test@example.com"`,
				`"tokenKey":""`,
				`"password":"456"`,
			},
			notExpected: []string{
				`"oldPassword"`,
				`"passwordConfirm"`,
			},
		},
		{
			name: "hidden fields (superuser)",
			data: map[string]any{
				"email":    "test@example.com",
				"tokenKey": "abc",
				// special auth fields
				"password":        "456",
				"oldPassword":     "123",
				"passwordConfirm": "789",
			},
			superuserAccessLevel: true,
			record:               core.NewRecord(usersCol),
			expected: []string{
				`"email":"test@example.com"`,
				`"tokenKey":"abc"`,
				`"password":"456"`,
			},
			notExpected: []string{
				`"oldPassword"`,
				`"passwordConfirm"`,
			},
		},
		{
			name: "with file field",
			data: map[string]any{
				"file_one": file,
				"url":      file, // should be ignored for non-file fields
			},
			record: core.NewRecord(demo1Col),
			expected: []string{
				`"file_one":{`,
				`"originalName":"test.txt"`,
				`"url":""`,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			form := forms.NewRecordUpsert(testApp, s.record)

			if s.managerAccessLevel {
				form.GrantManagerAccess()
			}

			if s.superuserAccessLevel {
				form.GrantSuperuserAccess()
			}

			// ensure that the form access level was updated
			if !form.HasManageAccess() && (s.superuserAccessLevel || s.managerAccessLevel) {
				t.Fatalf("Expected the form to have manage access level (manager or superuser)")
			}

			form.Load(s.data)

			loaded := map[string]any{}
			maps.Copy(loaded, s.record.FieldsData())
			maps.Copy(loaded, s.record.CustomData())

			raw, err := json.Marshal(loaded)
			if err != nil {
				t.Fatalf("Failed to serialize data: %v", err)
			}

			rawStr := string(raw)

			for _, str := range s.expected {
				if !strings.Contains(rawStr, str) {
					t.Fatalf("Couldn't find %q in \n%v", str, rawStr)
				}
			}

			for _, str := range s.notExpected {
				if strings.Contains(rawStr, str) {
					t.Fatalf("Didn't expect %q in \n%v", str, rawStr)
				}
			}
		})
	}
}

func TestRecordUpsertDrySubmitFailure(t *testing.T) {
	runTest := func(t *testing.T, testApp core.App) {
		col, err := testApp.FindCollectionByNameOrId("demo1")
		if err != nil {
			t.Fatal(err)
		}

		originalId := "imy661ixudk5izi"

		record, err := testApp.FindRecordById(col, originalId)
		if err != nil {
			t.Fatal(err)
		}

		oldRaw, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}

		file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
		if err != nil {
			t.Fatal(err)
		}

		form := forms.NewRecordUpsert(testApp, record)
		form.Load(map[string]any{
			"text":       "test_update",
			"file_one":   file,
			"select_one": "!invalid", // should be allowed even if invalid since validations are not executed
		})

		calls := ""
		testApp.OnRecordValidate(col.Name).BindFunc(func(e *core.RecordEvent) error {
			calls += "a" // shouldn't be called
			return e.Next()
		})

		result := form.DrySubmit(func(txApp core.App, drySavedRecord *core.Record) error {
			calls += "b"
			return errors.New("error...")
		})

		if result == nil {
			t.Fatal("Expected DrySubmit error, got nil")
		}

		if calls != "b" {
			t.Fatalf("Expected calls %q, got %q", "ab", calls)
		}

		// refresh the record to ensure that the changes weren't persisted
		record, err = testApp.FindRecordById(col, originalId)
		if err != nil {
			t.Fatalf("Expected record with the original id %q to exist, got\n%v", originalId, record.PublicExport())
		}

		newRaw, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(oldRaw, newRaw) {
			t.Fatalf("Expected record\n%s\ngot\n%s", oldRaw, newRaw)
		}

		testFilesCount(t, testApp, record, 0)
	}

	t.Run("without parent transaction", func(t *testing.T) {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		runTest(t, testApp)
	})

	t.Run("with parent transaction", func(t *testing.T) {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		testApp.RunInTransaction(func(txApp core.App) error {
			runTest(t, txApp)
			return nil
		})
	})
}

func TestRecordUpsertDrySubmitCreateSuccess(t *testing.T) {
	runTest := func(t *testing.T, testApp core.App) {
		col, err := testApp.FindCollectionByNameOrId("demo1")
		if err != nil {
			t.Fatal(err)
		}

		record := core.NewRecord(col)

		file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
		if err != nil {
			t.Fatal(err)
		}

		form := forms.NewRecordUpsert(testApp, record)
		form.Load(map[string]any{
			"id":         "test",
			"text":       "test_update",
			"file_one":   file,
			"select_one": "!invalid", // should be allowed even if invalid since validations are not executed
		})

		calls := ""
		testApp.OnRecordValidate(col.Name).BindFunc(func(e *core.RecordEvent) error {
			calls += "a" // shouldn't be called
			return e.Next()
		})

		result := form.DrySubmit(func(txApp core.App, drySavedRecord *core.Record) error {
			calls += "b"
			return nil
		})

		if result != nil {
			t.Fatalf("Expected DrySubmit success, got error: %v", result)
		}

		if calls != "b" {
			t.Fatalf("Expected calls %q, got %q", "ab", calls)
		}

		// refresh the record to ensure that the changes weren't persisted
		_, err = testApp.FindRecordById(col, record.Id)
		if err == nil {
			t.Fatal("Expected the created record to be deleted")
		}

		testFilesCount(t, testApp, record, 0)
	}

	t.Run("without parent transaction", func(t *testing.T) {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		runTest(t, testApp)
	})

	t.Run("with parent transaction", func(t *testing.T) {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		testApp.RunInTransaction(func(txApp core.App) error {
			runTest(t, txApp)
			return nil
		})
	})
}

func TestRecordUpsertDrySubmitUpdateSuccess(t *testing.T) {
	runTest := func(t *testing.T, testApp core.App) {
		col, err := testApp.FindCollectionByNameOrId("demo1")
		if err != nil {
			t.Fatal(err)
		}

		record, err := testApp.FindRecordById(col, "imy661ixudk5izi")
		if err != nil {
			t.Fatal(err)
		}

		oldRaw, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}

		file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
		if err != nil {
			t.Fatal(err)
		}

		form := forms.NewRecordUpsert(testApp, record)
		form.Load(map[string]any{
			"text":     "test_update",
			"file_one": file,
		})

		calls := ""
		testApp.OnRecordValidate(col.Name).BindFunc(func(e *core.RecordEvent) error {
			calls += "a" // shouldn't be called
			return e.Next()
		})

		result := form.DrySubmit(func(txApp core.App, drySavedRecord *core.Record) error {
			calls += "b"
			return nil
		})

		if result != nil {
			t.Fatalf("Expected DrySubmit success, got error: %v", result)
		}

		if calls != "b" {
			t.Fatalf("Expected calls %q, got %q", "ab", calls)
		}

		// refresh the record to ensure that the changes weren't persisted
		record, err = testApp.FindRecordById(col, record.Id)
		if err != nil {
			t.Fatal(err)
		}

		newRaw, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(oldRaw, newRaw) {
			t.Fatalf("Expected record\n%s\ngot\n%s", oldRaw, newRaw)
		}

		testFilesCount(t, testApp, record, 0)
	}

	t.Run("without parent transaction", func(t *testing.T) {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		runTest(t, testApp)
	})

	t.Run("with parent transaction", func(t *testing.T) {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		testApp.RunInTransaction(func(txApp core.App) error {
			runTest(t, txApp)
			return nil
		})
	})
}

func TestRecordUpsertSubmitValidations(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo2Col, err := app.FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}

	demo2Rec, err := app.FindRecordById(demo2Col, "llvuca81nly1qls")
	if err != nil {
		t.Fatal(err)
	}

	usersCol, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	userRec, err := app.FindRecordById(usersCol, "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name           string
		record         *core.Record
		data           map[string]any
		managerAccess  bool
		expectedErrors []string
	}{
		// base
		{
			name:           "new base collection record with empty data",
			record:         core.NewRecord(demo2Col),
			data:           map[string]any{},
			expectedErrors: []string{"title"},
		},
		{
			name:   "new base collection record with invalid data",
			record: core.NewRecord(demo2Col),
			data: map[string]any{
				"title": "",
				// should be ignored
				"custom":          "abc",
				"oldPassword":     "123",
				"password":        "456",
				"passwordConfirm": "789",
			},
			expectedErrors: []string{"title"},
		},
		{
			name:   "new base collection record with valid data",
			record: core.NewRecord(demo2Col),
			data: map[string]any{
				"title": "abc",
				// should be ignored
				"custom":          "abc",
				"oldPassword":     "123",
				"password":        "456",
				"passwordConfirm": "789",
			},
			expectedErrors: []string{},
		},
		{
			name:           "existing base collection record with empty data",
			record:         demo2Rec,
			data:           map[string]any{},
			expectedErrors: []string{},
		},
		{
			name:   "existing base collection record with invalid data",
			record: demo2Rec,
			data: map[string]any{
				"title": "",
			},
			expectedErrors: []string{"title"},
		},
		{
			name:   "existing base collection record with valid data",
			record: demo2Rec,
			data: map[string]any{
				"title": "abc",
			},
			expectedErrors: []string{},
		},

		// auth
		{
			name:           "new auth collection record with empty data",
			record:         core.NewRecord(usersCol),
			data:           map[string]any{},
			expectedErrors: []string{"password", "passwordConfirm"},
		},
		{
			name:   "new auth collection record with invalid record and invalid form data (without manager acess)",
			record: core.NewRecord(usersCol),
			data: map[string]any{
				"verified":        true,
				"emailVisibility": true,
				"email":           "test@example.com",
				"password":        "456",
				"passwordConfirm": "789",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey":    strings.Repeat("a", 2),
				"custom":      "abc",
				"oldPassword": "123",
			},
			// fail the form validator
			expectedErrors: []string{"verified", "passwordConfirm"},
		},
		{
			name:   "new auth collection record with invalid record and valid form data (without manager acess)",
			record: core.NewRecord(usersCol),
			data: map[string]any{
				"verified":        false,
				"emailVisibility": true,
				"email":           "test@example.com",
				"password":        "456",
				"passwordConfirm": "456",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey":    strings.Repeat("a", 2),
				"custom":      "abc",
				"oldPassword": "123",
			},
			// fail the record fields validator
			expectedErrors: []string{"password", "username"},
		},
		{
			name:          "new auth collection record with invalid record and invalid form data (with manager acess)",
			record:        core.NewRecord(usersCol),
			managerAccess: true,
			data: map[string]any{
				"verified":        true,
				"emailVisibility": true,
				"email":           "test@example.com",
				"password":        "456",
				"passwordConfirm": "789",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey":    strings.Repeat("a", 2),
				"custom":      "abc",
				"oldPassword": "123",
			},
			// fail the form validator
			expectedErrors: []string{"passwordConfirm"},
		},
		{
			name:          "new auth collection record with invalid record and valid form data (with manager acess)",
			record:        core.NewRecord(usersCol),
			managerAccess: true,
			data: map[string]any{
				"verified":        true,
				"emailVisibility": true,
				"email":           "test@example.com",
				"password":        "456",
				"passwordConfirm": "456",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey":    strings.Repeat("a", 2),
				"custom":      "abc",
				"oldPassword": "123",
			},
			// fail the record fields validator
			expectedErrors: []string{"password", "username"},
		},
		{
			name:   "new auth collection record with valid data",
			record: core.NewRecord(usersCol),
			data: map[string]any{
				"emailVisibility": true,
				"email":           "test_new@example.com",
				"password":        "1234567890",
				"passwordConfirm": "1234567890",
				// should be ignored (custom or hidden fields)
				"tokenKey":    strings.Repeat("a", 2),
				"custom":      "abc",
				"oldPassword": "123",
			},
			expectedErrors: []string{},
		},
		{
			name:   "new auth collection record with valid data and duplicated email",
			record: core.NewRecord(usersCol),
			data: map[string]any{
				"email":           "test@example.com",
				"password":        "1234567890",
				"passwordConfirm": "1234567890",
				// should be ignored (custom or hidden fields)
				"tokenKey":    strings.Repeat("a", 2),
				"custom":      "abc",
				"oldPassword": "123",
			},
			// fail the unique db validator
			expectedErrors: []string{"email"},
		},
		{
			name:           "existing auth collection record with empty data",
			record:         userRec,
			data:           map[string]any{},
			expectedErrors: []string{},
		},
		{
			name:   "existing auth collection record with invalid record data and invalid form data (without manager access)",
			record: userRec,
			data: map[string]any{
				"verified":        true,
				"email":           "test_new@example.com", // not allowed to change
				"oldPassword":     "123",
				"password":        "456",
				"passwordConfirm": "789",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey": strings.Repeat("a", 2),
				"custom":   "abc",
			},
			// fail form validator
			expectedErrors: []string{"verified", "email", "oldPassword", "passwordConfirm"},
		},
		{
			name:   "existing auth collection record with invalid record data and valid form data (without manager access)",
			record: userRec,
			data: map[string]any{
				"oldPassword":     "1234567890",
				"password":        "12345678901",
				"passwordConfirm": "12345678901",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey": strings.Repeat("a", 2),
				"custom":   "abc",
			},
			// fail record fields validator
			expectedErrors: []string{"username"},
		},
		{
			name:          "existing auth collection record with invalid record data and invalid form data (with manager access)",
			record:        userRec,
			managerAccess: true,
			data: map[string]any{
				"verified":        true,
				"email":           "test_new@example.com",
				"oldPassword":     "123", // should be ignored
				"password":        "456",
				"passwordConfirm": "789",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey": strings.Repeat("a", 2),
				"custom":   "abc",
			},
			// fail form validator
			expectedErrors: []string{"passwordConfirm"},
		},
		{
			name:          "existing auth collection record with invalid record data and valid form data (with manager access)",
			record:        userRec,
			managerAccess: true,
			data: map[string]any{
				"verified":        true,
				"email":           "test_new@example.com",
				"oldPassword":     "1234567890",
				"password":        "12345678901",
				"passwordConfirm": "12345678901",
				"username":        "!invalid",
				// should be ignored (custom or hidden fields)
				"tokenKey": strings.Repeat("a", 2),
				"custom":   "abc",
			},
			// fail record fields validator
			expectedErrors: []string{"username"},
		},
		{
			name:   "existing auth collection record with base valid data",
			record: userRec,
			data: map[string]any{
				"name": "test",
			},
			expectedErrors: []string{},
		},
		{
			name:   "existing auth collection record with valid password and invalid oldPassword data",
			record: userRec,
			data: map[string]any{
				"name":            "test",
				"oldPassword":     "invalid",
				"password":        "1234567890",
				"passwordConfirm": "1234567890",
			},
			expectedErrors: []string{"oldPassword"},
		},
		{
			name:   "existing auth collection record with valid password data",
			record: userRec,
			data: map[string]any{
				"name":            "test",
				"oldPassword":     "1234567890",
				"password":        "0987654321",
				"passwordConfirm": "0987654321",
			},
			expectedErrors: []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testApp, _ := tests.NewTestApp()
			defer testApp.Cleanup()

			form := forms.NewRecordUpsert(testApp, s.record.Original())
			if s.managerAccess {
				form.GrantManagerAccess()
			}
			form.Load(s.data)

			result := form.Submit()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestRecordUpsertSubmitFailure(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	col, err := testApp.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	record, err := testApp.FindRecordById(col, "imy661ixudk5izi")
	if err != nil {
		t.Fatal(err)
	}

	file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(testApp, record)
	form.Load(map[string]any{
		"text":       "test_update",
		"file_one":   file,
		"select_one": "invalid",
	})

	validateCalls := 0
	testApp.OnRecordValidate(col.Name).BindFunc(func(e *core.RecordEvent) error {
		validateCalls++
		return e.Next()
	})

	result := form.Submit()

	if result == nil {
		t.Fatal("Expected Submit error, got nil")
	}

	if validateCalls != 1 {
		t.Fatalf("Expected validateCalls %d, got %d", 1, validateCalls)
	}

	// refresh the record to ensure that the changes weren't persisted
	record, err = testApp.FindRecordById(col, record.Id)
	if err != nil {
		t.Fatal(err)
	}

	if v := record.GetString("text"); v == "test_update" {
		t.Fatalf("Expected record.text to remain the same, got %q", v)
	}

	if v := record.GetString("select_one"); v != "" {
		t.Fatalf("Expected record.select_one to remain the same, got %q", v)
	}

	if v := record.GetString("file_one"); v != "" {
		t.Fatalf("Expected record.file_one to remain the same, got %q", v)
	}

	testFilesCount(t, testApp, record, 0)
}

func TestRecordUpsertSubmitSuccess(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	col, err := testApp.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	record, err := testApp.FindRecordById(col, "imy661ixudk5izi")
	if err != nil {
		t.Fatal(err)
	}

	file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordUpsert(testApp, record)
	form.Load(map[string]any{
		"text":       "test_update",
		"file_one":   file,
		"select_one": "optionC",
	})

	validateCalls := 0
	testApp.OnRecordValidate(col.Name).BindFunc(func(e *core.RecordEvent) error {
		validateCalls++
		return e.Next()
	})

	result := form.Submit()

	if result != nil {
		t.Fatalf("Expected Submit success, got error: %v", result)
	}

	if validateCalls != 1 {
		t.Fatalf("Expected validateCalls %d, got %d", 1, validateCalls)
	}

	// refresh the record to ensure that the changes were persisted
	record, err = testApp.FindRecordById(col, record.Id)
	if err != nil {
		t.Fatal(err)
	}

	if v := record.GetString("text"); v != "test_update" {
		t.Fatalf("Expected record.text %q, got %q", "test_update", v)
	}

	if v := record.GetString("select_one"); v != "optionC" {
		t.Fatalf("Expected record.select_one %q, got %q", "optionC", v)
	}

	if v := record.GetString("file_one"); v != file.Name {
		t.Fatalf("Expected record.file_one %q, got %q", file.Name, v)
	}

	testFilesCount(t, testApp, record, 2) // the file + attrs
}

func TestRecordUpsertPasswordsSync(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	users, err := testApp.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("new user without password", func(t *testing.T) {
		record := core.NewRecord(users)

		form := forms.NewRecordUpsert(testApp, record)

		err := form.Submit()

		tests.TestValidationErrors(t, err, []string{"password", "passwordConfirm"})
	})

	t.Run("new user with manual password", func(t *testing.T) {
		record := core.NewRecord(users)

		form := forms.NewRecordUpsert(testApp, record)

		record.SetPassword("1234567890")

		err := form.Submit()
		if err != nil {
			t.Fatalf("Expected no errors, got %v", err)
		}
	})

	t.Run("new user with random password", func(t *testing.T) {
		record := core.NewRecord(users)

		form := forms.NewRecordUpsert(testApp, record)

		record.SetRandomPassword()

		err := form.Submit()
		if err != nil {
			t.Fatalf("Expected no errors, got %v", err)
		}
	})

	t.Run("update user with no password change", func(t *testing.T) {
		record, err := testApp.FindAuthRecordByEmail(users, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		oldHash := record.GetString("password:hash")

		form := forms.NewRecordUpsert(testApp, record)

		err = form.Submit()
		if err != nil {
			t.Fatalf("Expected no errors, got %v", err)
		}

		newHash := record.GetString("password:hash")
		if newHash == "" || newHash != oldHash {
			t.Fatal("Expected no password change")
		}
	})

	t.Run("update user with manual password change", func(t *testing.T) {
		record, err := testApp.FindAuthRecordByEmail(users, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		oldHash := record.GetString("password:hash")

		form := forms.NewRecordUpsert(testApp, record)

		record.SetPassword("1234567890")

		err = form.Submit()
		if err != nil {
			t.Fatalf("Expected no errors, got %v", err)
		}

		newHash := record.GetString("password:hash")
		if newHash == "" || newHash == oldHash {
			t.Fatal("Expected password change")
		}
	})

	t.Run("update user with random password change", func(t *testing.T) {
		record, err := testApp.FindAuthRecordByEmail(users, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		oldHash := record.GetString("password:hash")

		form := forms.NewRecordUpsert(testApp, record)

		record.SetRandomPassword()

		err = form.Submit()
		if err != nil {
			t.Fatalf("Expected no errors, got %v", err)
		}

		newHash := record.GetString("password:hash")
		if newHash == "" || newHash == oldHash {
			t.Fatal("Expected password change")
		}
	})
}

// -------------------------------------------------------------------

func testFilesCount(t *testing.T, app core.App, record *core.Record, count int) {
	storageDir := filepath.Join(app.DataDir(), "storage", record.Collection().Id, record.Id)

	entries, _ := os.ReadDir(storageDir)
	if len(entries) != count {
		t.Errorf("Expected %d entries, got %d\n%v", count, len(entries), entries)
	}
}
