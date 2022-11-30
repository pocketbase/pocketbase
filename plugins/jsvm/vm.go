package jsvm

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode"

	"github.com/dop251/goja"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func NewBaseVM(app core.App) *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(fieldMapper{})
	vm.Set("$app", app)

	baseBind(vm)
	dbxBind(vm)

	return vm
}

func baseBind(vm *goja.Runtime) {
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

	vm.Set("Collection", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Collection{}
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())
		return instanceValue
	})

	vm.Set("Record", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Record{}
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())
		return instanceValue
	})

	vm.Set("Admin", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Admin{}
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())
		return instanceValue
	})

	vm.Set("Schema", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.Schema{}
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())
		return instanceValue
	})

	vm.Set("SchemaField", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.SchemaField{}
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())
		return instanceValue
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

func dbxBind(vm *goja.Runtime) {
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

// fieldMapper provides custom mapping between Go and JavaScript property names.
//
// It is similar to the builtin "uncapFieldNameMapper" but also converts
// all uppercase identifiers to their lowercase equivalent (eg. "GET" -> "get").
type fieldMapper struct {
}

// FieldName implements the [FieldNameMapper.FieldName] interface method.
func (u fieldMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	return convertGoToJSName(f.Name)
}

// MethodName implements the [FieldNameMapper.MethodName] interface method.
func (u fieldMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return convertGoToJSName(m.Name)
}

func convertGoToJSName(name string) string {
	allUppercase := true
	for _, c := range name {
		if !unicode.IsUpper(c) {
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
