package daos_test

import (
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestParamQuery(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_params}}.* FROM `_params`"

	sql := app.Dao().ParamQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindParamByKey(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		key         string
		expectError bool
	}{
		{"", true},
		{"missing", true},
		{models.ParamAppSettings, false},
	}

	for i, scenario := range scenarios {
		param, err := app.Dao().FindParamByKey(scenario.key)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if param != nil && param.Key != scenario.key {
			t.Errorf("(%d) Expected param with identifier %s, got %v", i, scenario.key, param.Key)
		}
	}
}

func TestSaveParam(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		key   string
		value any
	}{
		{"", "demo"},
		{"test", nil},
		{"test", ""},
		{"test", 1},
		{"test", 123},
		{models.ParamAppSettings, map[string]any{"test": 123}},
	}

	for i, scenario := range scenarios {
		err := app.Dao().SaveParam(scenario.key, scenario.value)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		jsonRaw := types.JsonRaw{}
		jsonRaw.Scan(scenario.value)
		encodedScenarioValue, err := jsonRaw.MarshalJSON()
		if err != nil {
			t.Errorf("(%d) Encoded error %v", i, err)
		}

		// check if the param was really saved
		param, _ := app.Dao().FindParamByKey(scenario.key)
		encodedParamValue, err := param.Value.MarshalJSON()
		if err != nil {
			t.Errorf("(%d) Encoded error %v", i, err)
		}

		if string(encodedParamValue) != string(encodedScenarioValue) {
			t.Errorf("(%d) Expected the two values to be equal, got %v vs %v", i, string(encodedParamValue), string(encodedScenarioValue))
		}
	}
}

func TestSaveParamEncrypted(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	encryptionKey := security.RandomString(32)
	data := map[string]int{"test": 123}
	expected := map[string]int{}

	err := app.Dao().SaveParam("test", data, encryptionKey)
	if err != nil {
		t.Fatal(err)
	}

	// check if the param was really saved
	param, _ := app.Dao().FindParamByKey("test")

	// decrypt
	decrypted, decryptErr := security.Decrypt(string(param.Value), encryptionKey)
	if decryptErr != nil {
		t.Fatal(decryptErr)
	}

	// decode
	decryptedDecodeErr := json.Unmarshal(decrypted, &expected)
	if decryptedDecodeErr != nil {
		t.Fatal(decryptedDecodeErr)
	}

	// check if the decoded value is correct
	if len(expected) != len(data) || expected["test"] != data["test"] {
		t.Fatalf("Expected %v, got %v", expected, data)
	}
}

func TestDeleteParam(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// unsaved param
	err1 := app.Dao().DeleteParam(&models.Param{})
	if err1 == nil {
		t.Fatal("Expected error, got nil")
	}

	// existing param
	param, _ := app.Dao().FindParamByKey(models.ParamAppSettings)
	err2 := app.Dao().DeleteParam(param)
	if err2 != nil {
		t.Fatalf("Expected nil, got error %v", err2)
	}

	// check if it was really deleted
	paramCheck, _ := app.Dao().FindParamByKey(models.ParamAppSettings)
	if paramCheck != nil {
		t.Fatalf("Expected param to be deleted, got %v", paramCheck)
	}
}
