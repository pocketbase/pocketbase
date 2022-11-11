package apis_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordCreateFail(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()
	users, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Error(err)
	}
	// Failed to create record.
	_, err = apis.RecordCreate(testApp.BaseApp, users, nil, nil)
	if err == nil {
		t.Fail()
	}
	// Failed to create record.
	fieldsMap := map[string]interface{}{}
	fieldsMap["name"] = "username"
	_, err = apis.RecordCreate(testApp.BaseApp, users, nil, fieldsMap)
	if err == nil {
		t.Fail()
	}
	// password: cannot be blank; passwordConfirm: cannot be blank.
	admin := models.Admin{} // create as admin
	fieldsMap["password"] = "password"
	_, err = apis.RecordCreate(testApp.BaseApp, users, &admin, fieldsMap)
	if err == nil {
		t.Fail()
	}
}

func TestRecordCreateOk(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	users, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Error(err)
	}
	admin := models.Admin{}

	fieldsMap := map[string]interface{}{}
	fieldsMap["name"] = "username"
	fieldsMap["password"] = "password"
	fieldsMap["passwordConfirm"] = "password"
	user, err := apis.RecordCreate(testApp.BaseApp, users, &admin, fieldsMap)
	if err != nil {
		t.Error(err)
	}

	if user.GetString("name") != "username" {
		t.Error("wrong username:" + user.GetString("name"))
	}
}
