// Package jsvm implements optional utilities for binding a JS goja runtime
// to the PocketBase instance (loading migrations, attaching to app hooks, etc.).
//
// Currently it provides the following plugins:
//
// 1. JS Migrations loader:
//
//	jsvm.MustRegisterMigrations(app, &jsvm.MigrationsOptions{
//		Dir: "custom_js_migrations_dir_path", // default to "pb_data/../pb_migrations"
//	})
package jsvm

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode"

	"github.com/dop251/goja"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func NewBaseVM() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(FieldMapper{})

	baseBinds(vm)
	dbxBinds(vm)

	return vm
}

func baseBinds(vm *goja.Runtime) {
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
		return defaultConstructor(vm, call, instance)
	})

	vm.Set("Admin", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Admin{}
		return defaultConstructor(vm, call, instance)
	})

	vm.Set("Schema", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.Schema{}
		return defaultConstructor(vm, call, instance)
	})

	vm.Set("SchemaField", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.SchemaField{}
		return defaultConstructor(vm, call, instance)
	})

	vm.Set("Dao", func(call goja.ConstructorCall) *goja.Object {
		db, ok := call.Argument(0).Export().(dbx.Builder)
		if !ok || db == nil {
			panic("missing required Dao(db) argument")
		}

		instance := daos.New(db)
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})
}

func defaultConstructor(vm *goja.Runtime, call goja.ConstructorCall, instance any) *goja.Object {
	if data := call.Argument(0).Export(); data != nil {
		if raw, err := json.Marshal(data); err == nil {
			json.Unmarshal(raw, instance)
		}
	}

	instanceValue := vm.ToValue(instance).(*goja.Object)
	instanceValue.SetPrototype(call.This.Prototype())

	return instanceValue
}

func dbxBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$dbx", obj)

	obj.Set("exp", dbx.NewExp)
	obj.Set("hashExp", func(data map[string]any) dbx.HashExp {
		exp := dbx.HashExp{}
		for k, v := range data {
			exp[k] = v
		}
		return exp
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

func apisBind(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$apis", obj)

	// middlewares
	obj.Set("requireRecordAuth", apis.RequireRecordAuth)
	obj.Set("requireRecordAuth", apis.RequireRecordAuth)
	obj.Set("requireSameContextRecordAuth", apis.RequireSameContextRecordAuth)
	obj.Set("requireAdminAuth", apis.RequireAdminAuth)
	obj.Set("requireAdminAuthOnlyIfAny", apis.RequireAdminAuthOnlyIfAny)
	obj.Set("requireAdminOrRecordAuth", apis.RequireAdminOrRecordAuth)
	obj.Set("requireAdminOrOwnerAuth", apis.RequireAdminOrOwnerAuth)
	obj.Set("activityLogger", apis.ActivityLogger)

	// api errors
	obj.Set("notFoundError", apis.NewNotFoundError)
	obj.Set("badRequestError", apis.NewBadRequestError)
	obj.Set("forbiddenError", apis.NewForbiddenError)
	obj.Set("unauthorizedError", apis.NewUnauthorizedError)

	// record helpers
	obj.Set("requestData", apis.RequestData)
	obj.Set("enrichRecord", apis.EnrichRecord)
	obj.Set("enrichRecords", apis.EnrichRecords)
}

// FieldMapper provides custom mapping between Go and JavaScript property names.
//
// It is similar to the builtin "uncapFieldNameMapper" but also converts
// all uppercase identifiers to their lowercase equivalent (eg. "GET" -> "get").
type FieldMapper struct {
}

// FieldName implements the [FieldNameMapper.FieldName] interface method.
func (u FieldMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	return convertGoToJSName(f.Name)
}

// MethodName implements the [FieldNameMapper.MethodName] interface method.
func (u FieldMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return convertGoToJSName(m.Name)
}

func convertGoToJSName(name string) string {
	allUppercase := true
	for _, c := range name {
		if c != '_' && !unicode.IsUpper(c) {
			allUppercase = false
			break
		}
	}

	// eg. "JSON" -> "json"
	if allUppercase {
		return strings.ToLower(name)
	}

	// eg. "GetField" -> "getField"
	return strings.ToLower(name[0:1]) + name[1:]
}
