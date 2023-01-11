package forms_test

import (
	"encoding/json"
	"github.com/pocketbase/pocketbase/tools/auth"
	"reflect"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestUserTelegramLoginValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		testName       string
		collectionName string
		jsonData       string
		expectedErrors []string
	}{
		{
			"empty payload",
			"users",
			"{}",
			[]string{"data"},
		},
		{
			"empty data",
			"users",
			`{"data":""}`,
			[]string{"data"},
		},
		{
			"invalid data",
			"users",
			`{"data":"invalid"}`,
			[]string{"data"},
		},
		{
			"disabled tg auth",
			"clients",
			`{"data":"query_id=AAGSTRQLAAAAAJJNFAsbizs2&user=%7B%22id%22%3A185879954%2C%22first_name%22%3A%22Ilya%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22beer13%22%2C%22language_code%22%3A%22ru%22%7D&auth_date=1673317539&hash=74e1b67c230d2343f5d317a4d77841e9c673cae1bde28606a40825a98c7be638"}`,
			[]string{"data"},
		},
		{
			"valid data web auth",
			"users",
			`{"data": "query_id=AAGSTRQLAAAAAJJNFAsbizs2&user=%7B%22id%22%3A185879954%2C%22first_name%22%3A%22Ilya%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22beer13%22%2C%22language_code%22%3A%22ru%22%7D&auth_date=1673317539&hash=74e1b67c230d2343f5d317a4d77841e9c673cae1bde28606a40825a98c7be638"}`,
			[]string{},
		},
		{
			"valid data login widget",
			"users",
			`{"data":"id=185879954&first_name=Ilya&last_name=&username=beer13&language_code=ru&hash=bf8e28bc7dfed2415ef50b70b2ed64759a94cd4ff647fec150cbc721f988066a"}`,
			[]string{},
		},
	}

	for _, s := range scenarios {
		authCollection, _ := app.Dao().FindCollectionByNameOrId(s.collectionName)
		if authCollection == nil {
			t.Errorf("[%s] Failed to fetch auth collection", s.testName)
		}

		form := forms.NewRecordTelegramLogin(app, authCollection, nil)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", s.testName, loadErr)
			continue
		}

		err := form.Validate()

		// parse errors
		errs, ok := err.(validation.Errors)
		if !ok && err != nil {
			t.Errorf("[%s] Failed to parse errors %v", s.testName, err)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", s.testName, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", s.testName, k, errs)
			}
		}
	}
}

func TestUserTelegramGetDataParsed(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		testName       string
		collectionName string
		jsonData       string
		expectedOutput auth.AuthUser
	}{
		{
			"web app data",
			"users",
			`{"data":"query_id=AAGSTRQLAAAAAJJNFAsbizs2&user=%7B%22id%22%3A185879954%2C%22first_name%22%3A%22Ilya%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22beer13%22%2C%22language_code%22%3A%22ru%22%7D&auth_date=1673317539&hash=74e1b67c230d2343f5d317a4d77841e9c673cae1bde28606a40825a98c7be638"}`,
			auth.AuthUser{
				Id:       "185879954",
				Name:     "Ilya",
				Username: "beer13",
				RawUser: map[string]any{
					"auth_date": "1673317539",
					"hash":      "74e1b67c230d2343f5d317a4d77841e9c673cae1bde28606a40825a98c7be638",
					"query_id":  "AAGSTRQLAAAAAJJNFAsbizs2",
					"user":      "{\"id\":185879954,\"first_name\":\"Ilya\",\"last_name\":\"\",\"username\":\"beer13\",\"language_code\":\"ru\"}",
				},
			},
		},
		{
			"login widget data",
			"users",
			`{"data":"id=185879954&first_name=Ilya&last_name=&username=beer13&language_code=ru&hash=bf8e28bc7dfed2415ef50b70b2ed64759a94cd4ff647fec150cbc721f988066a"}`,
			auth.AuthUser{
				Id:       "185879954",
				Name:     "Ilya",
				Username: "beer13",
				RawUser: map[string]any{
					"first_name":    "Ilya",
					"id":            "185879954",
					"language_code": "ru",
					"last_name":     "",
					"username":      "beer13",
					"hash":          "bf8e28bc7dfed2415ef50b70b2ed64759a94cd4ff647fec150cbc721f988066a",
				},
			},
		},
	}

	for _, s := range scenarios {
		authCollection, _ := app.Dao().FindCollectionByNameOrId(s.collectionName)
		if authCollection == nil {
			t.Errorf("[%s] Failed to fetch auth collection", s.testName)
		}

		form := forms.NewRecordTelegramLogin(app, authCollection, nil)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", s.testName, loadErr)
			continue
		}

		authData, parseDataErr := form.GetAuthUserFromData()
		if parseDataErr != nil {
			t.Errorf("[%s] Failed to parse form data: %v", s.testName, parseDataErr)
			continue
		}

		if !reflect.DeepEqual(authData, &s.expectedOutput) {
			t.Errorf("[%s] Auth data not equal. Expected\n %#v\n got\n %#v", s.testName, s.expectedOutput, authData)
		}

	}
}
