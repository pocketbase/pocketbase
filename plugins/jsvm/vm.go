// Package jsvm implements optional utilities for binding a JS goja runtime
// to the PocketBase instance (loading migrations, attaching to app hooks, etc.).
//
// Currently it provides the following plugins:
//
// 1. JS Migrations loader:
//
//	jsvm.MustRegisterMigrations(app, jsvm.MigrationsConfig{
//		Dir: "/custom/js/migrations/dir", // default to "pb_data/../pb_migrations"
//	})
//
// 2. JS app hooks:
//
//	jsvm.MustRegisterHooks(app, jsvm.HooksConfig{
//		Dir: "/custom/js/hooks/dir", // default to "pb_data/../pb_hooks"
//	})
package jsvm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/dop251/goja"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/security"
)

func baseBinds(vm *goja.Runtime) {
	vm.SetFieldNameMapper(FieldMapper{})

	// override primitive class constructors to return pointers
	// (this is useful when unmarshaling or scaning a db result)
	vm.Set("_numberPointer", func(arg float64) *float64 {
		return &arg
	})
	vm.Set("_stringPointer", func(arg string) *string {
		return &arg
	})
	vm.Set("_boolPointer", func(arg bool) *bool {
		return &arg
	})
	vm.RunString(`
		this.Number = function(arg) {
			return _numberPointer(arg)
		}
		this.String = function(arg) {
			return _stringPointer(arg)
		}
		this.Boolean = function(arg) {
			return _boolPointer(arg)
		}
	`)

	vm.Set("unmarshal", func(src map[string]any, dest any) (any, error) {
		raw, err := json.Marshal(src)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(raw, &dest); err != nil {
			return nil, err
		}

		return dest, nil
	})

	vm.Set("Record", func(call goja.ConstructorCall) *goja.Object {
		var instance *models.Record

		collection, ok := call.Argument(0).Export().(*models.Collection)
		if ok {
			instance = models.NewRecord(collection)
			data, ok := call.Argument(1).Export().(map[string]any)
			if ok {
				if raw, err := json.Marshal(data); err == nil {
					json.Unmarshal(raw, instance)
				}
			}
		} else {
			instance = &models.Record{}
		}

		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	vm.Set("Collection", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Collection{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("Admin", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Admin{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("Schema", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.Schema{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("SchemaField", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.SchemaField{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("Mail", func(call goja.ConstructorCall) *goja.Object {
		instance := &mailer.Message{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("ValidationError", func(call goja.ConstructorCall) *goja.Object {
		code, _ := call.Argument(0).Export().(string)
		message, _ := call.Argument(1).Export().(string)

		instance := validation.NewError(code, message)
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	vm.Set("Dao", func(call goja.ConstructorCall) *goja.Object {
		concurrentDB, _ := call.Argument(0).Export().(dbx.Builder)
		if concurrentDB == nil {
			panic("missing required Dao(concurrentDB, [nonconcurrentDB]) argument")
		}

		nonConcurrentDB, _ := call.Argument(1).Export().(dbx.Builder)
		if nonConcurrentDB == nil {
			nonConcurrentDB = concurrentDB
		}

		instance := daos.NewMultiDB(concurrentDB, nonConcurrentDB)
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})
}

func dbxBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$dbx", obj)

	obj.Set("exp", dbx.NewExp)
	obj.Set("hashExp", func(data map[string]any) dbx.HashExp {
		return dbx.HashExp(data)
	})
	obj.Set("not", dbx.Not)
	obj.Set("and", dbx.And)
	obj.Set("or", dbx.Or)
	obj.Set("in", dbx.In)
	obj.Set("notIn", dbx.NotIn)
	obj.Set("like", dbx.Like)
	obj.Set("orLike", dbx.OrLike)
	obj.Set("notLike", dbx.NotLike)
	obj.Set("orNotLike", dbx.OrNotLike)
	obj.Set("exists", dbx.Exists)
	obj.Set("notExists", dbx.NotExists)
	obj.Set("between", dbx.Between)
	obj.Set("notBetween", dbx.NotBetween)
}

func tokensBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$tokens", obj)

	// admin
	obj.Set("adminAuthToken", tokens.NewAdminAuthToken)
	obj.Set("adminResetPasswordToken", tokens.NewAdminResetPasswordToken)
	obj.Set("adminFileToken", tokens.NewAdminFileToken)

	// record
	obj.Set("recordAuthToken", tokens.NewRecordAuthToken)
	obj.Set("recordVerifyToken", tokens.NewRecordVerifyToken)
	obj.Set("recordResetPasswordToken", tokens.NewRecordResetPasswordToken)
	obj.Set("recordChangeEmailToken", tokens.NewRecordChangeEmailToken)
	obj.Set("recordFileToken", tokens.NewRecordFileToken)
}

func securityBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$security", obj)

	// random
	obj.Set("randomString", security.RandomString)
	obj.Set("randomStringWithAlphabet", security.RandomStringWithAlphabet)
	obj.Set("pseudorandomString", security.PseudorandomString)
	obj.Set("pseudorandomStringWithAlphabet", security.PseudorandomStringWithAlphabet)

	// jwt
	obj.Set("parseUnverifiedToken", security.ParseUnverifiedJWT)
	obj.Set("parseToken", security.ParseJWT)
	obj.Set("createToken", security.NewToken)
}

func filesystemBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$filesystem", obj)

	obj.Set("fileFromPath", filesystem.NewFileFromPath)
	obj.Set("fileFromBytes", filesystem.NewFileFromBytes)
	obj.Set("fileFromMultipart", filesystem.NewFileFromMultipart)
}

func formsBinds(vm *goja.Runtime) {
	registerFactoryAsConstructor(vm, "AdminLoginForm", forms.NewAdminLogin)
	registerFactoryAsConstructor(vm, "AdminPasswordResetConfirmForm", forms.NewAdminPasswordResetConfirm)
	registerFactoryAsConstructor(vm, "AdminPasswordResetRequestForm", forms.NewAdminPasswordResetRequest)
	registerFactoryAsConstructor(vm, "AdminUpsertForm", forms.NewAdminUpsert)
	registerFactoryAsConstructor(vm, "AppleClientSecretCreateForm", forms.NewAppleClientSecretCreate)
	registerFactoryAsConstructor(vm, "CollectionUpsertForm", forms.NewCollectionUpsert)
	registerFactoryAsConstructor(vm, "CollectionsImportForm", forms.NewCollectionsImport)
	registerFactoryAsConstructor(vm, "RealtimeSubscribeForm", forms.NewRealtimeSubscribe)
	registerFactoryAsConstructor(vm, "RecordEmailChangeConfirmForm", forms.NewRecordEmailChangeConfirm)
	registerFactoryAsConstructor(vm, "RecordEmailChangeRequestForm", forms.NewRecordEmailChangeRequest)
	registerFactoryAsConstructor(vm, "RecordOAuth2LoginForm", forms.NewRecordOAuth2Login)
	registerFactoryAsConstructor(vm, "RecordPasswordLoginForm", forms.NewRecordPasswordLogin)
	registerFactoryAsConstructor(vm, "RecordPasswordResetConfirmForm", forms.NewRecordPasswordResetConfirm)
	registerFactoryAsConstructor(vm, "RecordPasswordResetRequestForm", forms.NewRecordPasswordResetRequest)
	registerFactoryAsConstructor(vm, "RecordUpsertForm", forms.NewRecordUpsert)
	registerFactoryAsConstructor(vm, "RecordVerificationConfirmForm", forms.NewRecordVerificationConfirm)
	registerFactoryAsConstructor(vm, "RecordVerificationRequestForm", forms.NewRecordVerificationRequest)
	registerFactoryAsConstructor(vm, "SettingsUpsertForm", forms.NewSettingsUpsert)
	registerFactoryAsConstructor(vm, "TestEmailSendForm", forms.NewTestEmailSend)
	registerFactoryAsConstructor(vm, "TestS3FilesystemForm", forms.NewTestS3Filesystem)
}

func apisBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$apis", obj)

	// middlewares
	obj.Set("requireRecordAuth", apis.RequireRecordAuth)
	obj.Set("requireSameContextRecordAuth", apis.RequireSameContextRecordAuth)
	obj.Set("requireAdminAuth", apis.RequireAdminAuth)
	obj.Set("requireAdminAuthOnlyIfAny", apis.RequireAdminAuthOnlyIfAny)
	obj.Set("requireAdminOrRecordAuth", apis.RequireAdminOrRecordAuth)
	obj.Set("requireAdminOrOwnerAuth", apis.RequireAdminOrOwnerAuth)
	obj.Set("activityLogger", apis.ActivityLogger)

	// record helpers
	obj.Set("requestData", apis.RequestData)
	obj.Set("recordAuthResponse", apis.RecordAuthResponse)
	obj.Set("enrichRecord", apis.EnrichRecord)
	obj.Set("enrichRecords", apis.EnrichRecords)

	// api errors
	vm.Set("ApiError", func(call goja.ConstructorCall) *goja.Object {
		status, _ := call.Argument(0).Export().(int64)
		message, _ := call.Argument(1).Export().(string)
		data := call.Argument(2).Export()

		instance := apis.NewApiError(int(status), message, data)
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})
	obj.Set("notFoundError", apis.NewNotFoundError)
	obj.Set("badRequestError", apis.NewBadRequestError)
	obj.Set("forbiddenError", apis.NewForbiddenError)
	obj.Set("unauthorizedError", apis.NewUnauthorizedError)

	vm.Set("Route", func(call goja.ConstructorCall) *goja.Object {
		instance := echo.Route{}
		return structConstructor(vm, call, &instance)
	})
}

// -------------------------------------------------------------------

// registerFactoryAsConstructor registers the factory function as native JS constructor.
func registerFactoryAsConstructor(vm *goja.Runtime, constructorName string, factoryFunc any) {
	vm.Set(constructorName, func(call goja.ConstructorCall) *goja.Object {
		f := reflect.ValueOf(factoryFunc)

		args := []reflect.Value{}

		for _, v := range call.Arguments {
			args = append(args, reflect.ValueOf(v.Export()))
		}

		result := f.Call(args)

		if len(result) != 1 {
			panic("the factory function should return only 1 item")
		}

		value := vm.ToValue(result[0].Interface()).(*goja.Object)
		value.SetPrototype(call.This.Prototype())

		return value
	})
}

// structConstructor wraps the provided struct with a native JS constructor.
func structConstructor(vm *goja.Runtime, call goja.ConstructorCall, instance any) *goja.Object {
	if data := call.Argument(0).Export(); data != nil {
		if raw, err := json.Marshal(data); err == nil {
			json.Unmarshal(raw, instance)
		}
	}

	instanceValue := vm.ToValue(instance).(*goja.Object)
	instanceValue.SetPrototype(call.This.Prototype())

	return instanceValue
}

// filesContent returns a map with all direct files within the specified dir and their content.
//
// If directory with dirPath is missing or no files matching the pattern were found,
// it returns an empty map and no error.
//
// If pattern is empty string it matches all root files.
func filesContent(dirPath string, pattern string) (map[string][]byte, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]byte{}, nil
		}
		return nil, err
	}

	var exp *regexp.Regexp
	if pattern != "" {
		var err error
		if exp, err = regexp.Compile(pattern); err != nil {
			return nil, err
		}
	}

	result := map[string][]byte{}

	for _, f := range files {
		if f.IsDir() || (exp != nil && !exp.MatchString(f.Name())) {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, err
		}

		result[f.Name()] = raw
	}

	return result, nil
}
