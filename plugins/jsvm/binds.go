package jsvm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/dop251/goja"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// hooksBinds adds wrapped "on*" hook methods by reflecting on core.App.
func hooksBinds(app core.App, loader *goja.Runtime, executors *vmsPool) {
	fm := FieldMapper{}

	appType := reflect.TypeOf(app)
	appValue := reflect.ValueOf(app)
	totalMethods := appType.NumMethod()
	excludeHooks := []string{"OnServe"}

	for i := 0; i < totalMethods; i++ {
		method := appType.Method(i)
		if !strings.HasPrefix(method.Name, "On") || slices.Contains(excludeHooks, method.Name) {
			continue // not a hook or excluded
		}

		jsName := fm.MethodName(appType, method)

		// register the hook to the loader
		loader.Set(jsName, func(callback string, tags ...string) {
			// overwrite the global $app with the hook scoped instance
			callback = `function(e) { $app = e.app; return (` + callback + `).call(undefined, e) }`
			pr := goja.MustCompile("", "{("+callback+").apply(undefined, __args)}", true)

			tagsAsValues := make([]reflect.Value, len(tags))
			for i, tag := range tags {
				tagsAsValues[i] = reflect.ValueOf(tag)
			}

			hookInstance := appValue.MethodByName(method.Name).Call(tagsAsValues)[0]
			hookBindFunc := hookInstance.MethodByName("BindFunc")

			handlerType := hookBindFunc.Type().In(0)

			handler := reflect.MakeFunc(handlerType, func(args []reflect.Value) (results []reflect.Value) {
				handlerArgs := make([]any, len(args))
				for i, arg := range args {
					handlerArgs[i] = arg.Interface()
				}

				err := executors.run(func(executor *goja.Runtime) error {
					executor.Set("$app", goja.Undefined())
					executor.Set("__args", handlerArgs)
					res, err := executor.RunProgram(pr)
					executor.Set("__args", goja.Undefined())

					// (legacy) check for returned Go error value
					if res != nil {
						if resErr, ok := res.Export().(error); ok {
							return resErr
						}
					}

					return normalizeException(err)
				})

				return []reflect.Value{reflect.ValueOf(&err).Elem()}
			})

			// register the wrapped hook handler
			hookBindFunc.Call([]reflect.Value{handler})
		})
	}
}

func cronBinds(app core.App, loader *goja.Runtime, executors *vmsPool) {
	loader.Set("cronAdd", func(jobId, cronExpr, handler string) {
		pr := goja.MustCompile("", "{("+handler+").apply(undefined)}", true)

		err := app.Cron().Add(jobId, cronExpr, func() {
			err := executors.run(func(executor *goja.Runtime) error {
				_, err := executor.RunProgram(pr)
				return err
			})

			if err != nil {
				app.Logger().Error(
					"[cronAdd] failed to execute cron job",
					slog.String("jobId", jobId),
					slog.String("error", err.Error()),
				)
			}
		})
		if err != nil {
			panic("[cronAdd] failed to register cron job " + jobId + ": " + err.Error())
		}
	})

	// note: it is not necessary needed but it is here for consistency
	loader.Set("cronRemove", func(jobId string) {
		app.Cron().Remove(jobId)
	})

	// register the removal helper also in the executors to allow removing cron jobs from everywhere
	oldFactory := executors.factory
	executors.factory = func() *goja.Runtime {
		vm := oldFactory()

		vm.Set("cronRemove", func(jobId string) {
			app.Cron().Remove(jobId)
		})

		return vm
	}
	for _, item := range executors.items {
		item.vm.Set("cronRemove", func(jobId string) {
			app.Cron().Remove(jobId)
		})
	}
}

func routerBinds(app core.App, loader *goja.Runtime, executors *vmsPool) {
	loader.Set("routerAdd", func(method string, path string, handler goja.Value, middlewares ...goja.Value) {
		wrappedMiddlewares, err := wrapMiddlewares(executors, middlewares...)
		if err != nil {
			panic("[routerAdd] failed to wrap middlewares: " + err.Error())
		}

		wrappedHandler, err := wrapHandlerFunc(executors, handler)
		if err != nil {
			panic("[routerAdd] failed to wrap handler: " + err.Error())
		}

		app.OnServe().BindFunc(func(e *core.ServeEvent) error {
			e.Router.Route(strings.ToUpper(method), path, wrappedHandler).Bind(wrappedMiddlewares...)

			return e.Next()
		})
	})

	loader.Set("routerUse", func(middlewares ...goja.Value) {
		wrappedMiddlewares, err := wrapMiddlewares(executors, middlewares...)
		if err != nil {
			panic("[routerUse] failed to wrap middlewares: " + err.Error())
		}

		app.OnServe().BindFunc(func(e *core.ServeEvent) error {
			e.Router.Bind(wrappedMiddlewares...)
			return e.Next()
		})
	})
}

func wrapHandlerFunc(executors *vmsPool, handler goja.Value) (func(*core.RequestEvent) error, error) {
	if handler == nil {
		return nil, errors.New("handler must be non-nil")
	}

	switch h := handler.Export().(type) {
	case func(*core.RequestEvent) error:
		// "native" handler func - no need to wrap
		return h, nil
	case func(goja.FunctionCall) goja.Value, string:
		pr := goja.MustCompile("", "{("+handler.String()+").apply(undefined, __args)}", true)

		wrappedHandler := func(e *core.RequestEvent) error {
			return executors.run(func(executor *goja.Runtime) error {
				executor.Set("$app", e.App) // overwrite the global $app with the hook scoped instance
				executor.Set("__args", []any{e})
				res, err := executor.RunProgram(pr)
				executor.Set("__args", goja.Undefined())

				// (legacy) check for returned Go error value
				if res != nil {
					if v, ok := res.Export().(error); ok {
						return v
					}
				}

				return normalizeException(err)
			})
		}

		return wrappedHandler, nil
	default:
		return nil, errors.New("unsupported goja handler type")
	}
}

type gojaHookHandler struct {
	id             string
	serializedFunc string
	priority       int
}

func wrapMiddlewares(executors *vmsPool, rawMiddlewares ...goja.Value) ([]*hook.Handler[*core.RequestEvent], error) {
	wrappedMiddlewares := make([]*hook.Handler[*core.RequestEvent], len(rawMiddlewares))

	for i, m := range rawMiddlewares {
		if m == nil {
			return nil, errors.New("middleware must be non-nil")
		}

		switch v := m.Export().(type) {
		case *hook.Handler[*core.RequestEvent]:
			// "native" middleware handler - no need to wrap
			wrappedMiddlewares[i] = v
		case func(*core.RequestEvent) error:
			// "native" middleware func - wrap as handler
			wrappedMiddlewares[i] = &hook.Handler[*core.RequestEvent]{
				Func: v,
			}
		case *gojaHookHandler:
			if v.serializedFunc == "" {
				return nil, errors.New("missing or invalid Middleware function")
			}

			pr := goja.MustCompile("", "{("+v.serializedFunc+").apply(undefined, __args)}", true)

			wrappedMiddlewares[i] = &hook.Handler[*core.RequestEvent]{
				Id:       v.id,
				Priority: v.priority,
				Func: func(e *core.RequestEvent) error {
					return executors.run(func(executor *goja.Runtime) error {
						executor.Set("$app", e.App) // overwrite the global $app with the hook scoped instance
						executor.Set("__args", []any{e})
						res, err := executor.RunProgram(pr)
						executor.Set("__args", goja.Undefined())

						// (legacy) check for returned Go error value
						if res != nil {
							if v, ok := res.Export().(error); ok {
								return v
							}
						}

						return normalizeException(err)
					})
				},
			}
		case func(goja.FunctionCall) goja.Value, string:
			pr := goja.MustCompile("", "{("+m.String()+").apply(undefined, __args)}", true)

			wrappedMiddlewares[i] = &hook.Handler[*core.RequestEvent]{
				Func: func(e *core.RequestEvent) error {
					return executors.run(func(executor *goja.Runtime) error {
						executor.Set("$app", e.App) // overwrite the global $app with the hook scoped instance
						executor.Set("__args", []any{e})
						res, err := executor.RunProgram(pr)
						executor.Set("__args", goja.Undefined())

						// (legacy) check for returned Go error value
						if res != nil {
							if v, ok := res.Export().(error); ok {
								return v
							}
						}

						return normalizeException(err)
					})
				},
			}
		default:
			return nil, errors.New("unsupported goja middleware type")
		}
	}

	return wrappedMiddlewares, nil
}

var cachedArrayOfTypes = store.New[reflect.Type, reflect.Type](nil)

func baseBinds(vm *goja.Runtime) {
	vm.SetFieldNameMapper(FieldMapper{})

	// deprecated: use toString
	vm.Set("readerToString", func(r io.Reader, maxBytes int) (string, error) {
		if maxBytes == 0 {
			maxBytes = router.DefaultMaxMemory
		}

		limitReader := io.LimitReader(r, int64(maxBytes))

		bodyBytes, readErr := io.ReadAll(limitReader)
		if readErr != nil {
			return "", readErr
		}

		return string(bodyBytes), nil
	})

	vm.Set("toString", func(raw any, maxReaderBytes int) (string, error) {
		switch v := raw.(type) {
		case io.Reader:
			if maxReaderBytes == 0 {
				maxReaderBytes = router.DefaultMaxMemory
			}

			limitReader := io.LimitReader(v, int64(maxReaderBytes))

			bodyBytes, readErr := io.ReadAll(limitReader)
			if readErr != nil {
				return "", readErr
			}

			return string(bodyBytes), nil
		default:
			str, err := cast.ToStringE(v)
			if err == nil {
				return str, nil
			}

			// as a last attempt try to json encode the value
			rawBytes, _ := json.Marshal(raw)

			return string(rawBytes), nil
		}
	})

	vm.Set("sleep", func(milliseconds int64) {
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
	})

	vm.Set("arrayOf", func(model any) any {
		mt := reflect.TypeOf(model)
		st := cachedArrayOfTypes.GetOrSet(mt, func() reflect.Type {
			return reflect.SliceOf(mt)
		})

		return reflect.New(st).Elem().Addr().Interface()
	})

	vm.Set("unmarshal", func(data, dst any) error {
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return json.Unmarshal(raw, &dst)
	})

	vm.Set("Context", func(call goja.ConstructorCall) *goja.Object {
		var instance context.Context

		oldCtx, ok := call.Argument(0).Export().(context.Context)
		if ok {
			instance = oldCtx
		} else {
			instance = context.Background()
		}

		key := call.Argument(1).Export()
		if key != nil {
			instance = context.WithValue(instance, key, call.Argument(2).Export())
		}

		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	vm.Set("DynamicModel", func(call goja.ConstructorCall) *goja.Object {
		shape, ok := call.Argument(0).Export().(map[string]any)
		if !ok || len(shape) == 0 {
			panic("[DynamicModel] missing shape data")
		}

		instance := newDynamicModel(shape)
		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	vm.Set("Record", func(call goja.ConstructorCall) *goja.Object {
		var instance *core.Record

		collection, ok := call.Argument(0).Export().(*core.Collection)
		if ok {
			instance = core.NewRecord(collection)
			data, ok := call.Argument(1).Export().(map[string]any)
			if ok {
				instance.Load(data)
			}
		} else {
			instance = &core.Record{}
		}

		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	vm.Set("Collection", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.Collection{}
		return structConstructorUnmarshal(vm, call, instance)
	})

	vm.Set("FieldsList", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.FieldsList{}
		return structConstructorUnmarshal(vm, call, instance)
	})

	// fields
	// ---
	vm.Set("Field", func(call goja.ConstructorCall) *goja.Object {
		data, _ := call.Argument(0).Export().(map[string]any)
		rawDataSlice, _ := json.Marshal([]any{data})

		fieldsList := core.NewFieldsList()
		_ = fieldsList.UnmarshalJSON(rawDataSlice)

		if len(fieldsList) == 0 {
			return nil
		}

		field := fieldsList[0]

		fieldValue := vm.ToValue(field).(*goja.Object)
		fieldValue.SetPrototype(call.This.Prototype())

		return fieldValue
	})
	vm.Set("NumberField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.NumberField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("BoolField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.BoolField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("TextField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.TextField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("URLField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.URLField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("EmailField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.EmailField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("EditorField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.EditorField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("PasswordField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.PasswordField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("DateField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.DateField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("AutodateField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.AutodateField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("JSONField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.JSONField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("RelationField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.RelationField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("SelectField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.SelectField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	vm.Set("FileField", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.FileField{}
		return structConstructorUnmarshal(vm, call, instance)
	})
	// ---

	vm.Set("MailerMessage", func(call goja.ConstructorCall) *goja.Object {
		instance := &mailer.Message{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("Command", func(call goja.ConstructorCall) *goja.Object {
		instance := &cobra.Command{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("RequestInfo", func(call goja.ConstructorCall) *goja.Object {
		instance := &core.RequestInfo{Context: core.RequestInfoContextDefault}
		return structConstructor(vm, call, instance)
	})

	// ```js
	// new Middleware((e) => {
	//    return e.next()
	// }, 100, "example_middleware")
	// ```
	vm.Set("Middleware", func(call goja.ConstructorCall) *goja.Object {
		instance := &gojaHookHandler{}

		instance.serializedFunc = call.Argument(0).String()
		instance.priority = cast.ToInt(call.Argument(1).Export())
		instance.id = cast.ToString(call.Argument(2).Export())

		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	vm.Set("DateTime", func(call goja.ConstructorCall) *goja.Object {
		instance := types.NowDateTime()

		val, _ := call.Argument(0).Export().(string)
		if val != "" {
			instance, _ = types.ParseDateTime(val)
		}

		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

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

	vm.Set("Cookie", func(call goja.ConstructorCall) *goja.Object {
		instance := &http.Cookie{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("SubscriptionMessage", func(call goja.ConstructorCall) *goja.Object {
		instance := &subscriptions.Message{}
		return structConstructor(vm, call, instance)
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

func mailsBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$mails", obj)

	obj.Set("sendRecordPasswordReset", mails.SendRecordPasswordReset)
	obj.Set("sendRecordVerification", mails.SendRecordVerification)
	obj.Set("sendRecordChangeEmail", mails.SendRecordChangeEmail)
	obj.Set("sendRecordOTP", mails.SendRecordOTP)
}

func securityBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$security", obj)

	// crypto
	obj.Set("md5", security.MD5)
	obj.Set("sha256", security.SHA256)
	obj.Set("sha512", security.SHA512)
	obj.Set("hs256", security.HS256)
	obj.Set("hs512", security.HS512)
	obj.Set("equal", security.Equal)

	// random
	obj.Set("randomString", security.RandomString)
	obj.Set("randomStringByRegex", security.RandomStringByRegex)
	obj.Set("randomStringWithAlphabet", security.RandomStringWithAlphabet)
	obj.Set("pseudorandomString", security.PseudorandomString)
	obj.Set("pseudorandomStringWithAlphabet", security.PseudorandomStringWithAlphabet)

	// jwt
	obj.Set("parseUnverifiedJWT", func(token string) (map[string]any, error) {
		return security.ParseUnverifiedJWT(token)
	})
	obj.Set("parseJWT", func(token string, verificationKey string) (map[string]any, error) {
		return security.ParseJWT(token, verificationKey)
	})
	obj.Set("createJWT", func(payload jwt.MapClaims, signingKey string, secDuration int) (string, error) {
		return security.NewJWT(payload, signingKey, time.Duration(secDuration)*time.Second)
	})

	// encryption
	obj.Set("encrypt", security.Encrypt)
	obj.Set("decrypt", func(cipherText, key string) (string, error) {
		result, err := security.Decrypt(cipherText, key)

		if err != nil {
			return "", err
		}

		return string(result), err
	})
}

func filesystemBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$filesystem", obj)

	obj.Set("fileFromPath", filesystem.NewFileFromPath)
	obj.Set("fileFromBytes", filesystem.NewFileFromBytes)
	obj.Set("fileFromMultipart", filesystem.NewFileFromMultipart)
	obj.Set("fileFromURL", func(url string, secTimeout int) (*filesystem.File, error) {
		if secTimeout == 0 {
			secTimeout = 120
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(secTimeout)*time.Second)
		defer cancel()

		return filesystem.NewFileFromURL(ctx, url)
	})
}

func filepathBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$filepath", obj)

	obj.Set("base", filepath.Base)
	obj.Set("clean", filepath.Clean)
	obj.Set("dir", filepath.Dir)
	obj.Set("ext", filepath.Ext)
	obj.Set("fromSlash", filepath.FromSlash)
	obj.Set("glob", filepath.Glob)
	obj.Set("isAbs", filepath.IsAbs)
	obj.Set("join", filepath.Join)
	obj.Set("match", filepath.Match)
	obj.Set("rel", filepath.Rel)
	obj.Set("split", filepath.Split)
	obj.Set("splitList", filepath.SplitList)
	obj.Set("toSlash", filepath.ToSlash)
	obj.Set("walk", filepath.Walk)
	obj.Set("walkDir", filepath.WalkDir)
}

func osBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$os", obj)

	obj.Set("args", os.Args)
	obj.Set("exec", exec.Command) // @deprecated
	obj.Set("cmd", exec.Command)
	obj.Set("exit", os.Exit)
	obj.Set("getenv", os.Getenv)
	obj.Set("dirFS", os.DirFS)
	obj.Set("readFile", os.ReadFile)
	obj.Set("writeFile", os.WriteFile)
	obj.Set("readDir", os.ReadDir)
	obj.Set("tempDir", os.TempDir)
	obj.Set("truncate", os.Truncate)
	obj.Set("getwd", os.Getwd)
	obj.Set("mkdir", os.Mkdir)
	obj.Set("mkdirAll", os.MkdirAll)
	obj.Set("rename", os.Rename)
	obj.Set("remove", os.Remove)
	obj.Set("removeAll", os.RemoveAll)
}

func formsBinds(vm *goja.Runtime) {
	registerFactoryAsConstructor(vm, "AppleClientSecretCreateForm", forms.NewAppleClientSecretCreate)
	registerFactoryAsConstructor(vm, "RecordUpsertForm", forms.NewRecordUpsert)
	registerFactoryAsConstructor(vm, "TestEmailSendForm", forms.NewTestEmailSend)
	registerFactoryAsConstructor(vm, "TestS3FilesystemForm", forms.NewTestS3Filesystem)
}

func apisBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$apis", obj)

	obj.Set("static", func(dir string, indexFallback bool) func(*core.RequestEvent) error {
		return apis.Static(os.DirFS(dir), indexFallback)
	})

	// middlewares
	obj.Set("requireGuestOnly", apis.RequireGuestOnly)
	obj.Set("requireAuth", apis.RequireAuth)
	obj.Set("requireSuperuserAuth", apis.RequireSuperuserAuth)
	obj.Set("requireSuperuserOrOwnerAuth", apis.RequireSuperuserOrOwnerAuth)
	obj.Set("skipSuccessActivityLog", apis.SkipSuccessActivityLog)
	obj.Set("gzip", apis.Gzip)
	obj.Set("bodyLimit", apis.BodyLimit)

	// record helpers
	obj.Set("recordAuthResponse", apis.RecordAuthResponse)
	obj.Set("enrichRecord", apis.EnrichRecord)
	obj.Set("enrichRecords", apis.EnrichRecords)

	// api errors
	registerFactoryAsConstructor(vm, "ApiError", router.NewApiError)
	registerFactoryAsConstructor(vm, "NotFoundError", router.NewNotFoundError)
	registerFactoryAsConstructor(vm, "BadRequestError", router.NewBadRequestError)
	registerFactoryAsConstructor(vm, "ForbiddenError", router.NewForbiddenError)
	registerFactoryAsConstructor(vm, "UnauthorizedError", router.NewUnauthorizedError)
	registerFactoryAsConstructor(vm, "TooManyRequestsError", router.NewTooManyRequestsError)
	registerFactoryAsConstructor(vm, "InternalServerError", router.NewInternalServerError)
}

func httpClientBinds(vm *goja.Runtime) {
	obj := vm.NewObject()
	vm.Set("$http", obj)

	vm.Set("FormData", func(call goja.ConstructorCall) *goja.Object {
		instance := FormData{}

		instanceValue := vm.ToValue(instance).(*goja.Object)
		instanceValue.SetPrototype(call.This.Prototype())

		return instanceValue
	})

	type sendResult struct {
		JSON       any                     `json:"json"`
		Headers    map[string][]string     `json:"headers"`
		Cookies    map[string]*http.Cookie `json:"cookies"`
		Raw        string                  `json:"raw"`
		StatusCode int                     `json:"statusCode"`
	}

	type sendConfig struct {
		// Deprecated: consider using Body instead
		Data map[string]any

		Body    any // raw string or FormData
		Headers map[string]string
		Method  string
		Url     string
		Timeout int // seconds (default to 120)
	}

	obj.Set("send", func(params map[string]any) (*sendResult, error) {
		config := sendConfig{
			Method: "GET",
		}

		if v, ok := params["data"]; ok {
			config.Data = cast.ToStringMap(v)
		}

		if v, ok := params["body"]; ok {
			config.Body = v
		}

		if v, ok := params["headers"]; ok {
			config.Headers = cast.ToStringMapString(v)
		}

		if v, ok := params["method"]; ok {
			config.Method = cast.ToString(v)
		}

		if v, ok := params["url"]; ok {
			config.Url = cast.ToString(v)
		}

		if v, ok := params["timeout"]; ok {
			config.Timeout = cast.ToInt(v)
		}

		if config.Timeout <= 0 {
			config.Timeout = 120
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
		defer cancel()

		var reqBody io.Reader
		var contentType string

		// legacy json body data
		if len(config.Data) != 0 {
			encoded, err := json.Marshal(config.Data)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewReader(encoded)
		} else {
			switch v := config.Body.(type) {
			case io.Reader:
				reqBody = v
			case FormData:
				body, mp, err := v.toMultipart()
				if err != nil {
					return nil, err
				}

				reqBody = body
				contentType = mp.FormDataContentType()
			default:
				reqBody = strings.NewReader(cast.ToString(config.Body))
			}
		}

		req, err := http.NewRequestWithContext(ctx, strings.ToUpper(config.Method), config.Url, reqBody)
		if err != nil {
			return nil, err
		}

		for k, v := range config.Headers {
			req.Header.Add(k, v)
		}

		// set the explicit content type
		// (overwriting the user provided header value if any)
		if contentType != "" {
			req.Header.Set("content-type", contentType)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		bodyRaw, _ := io.ReadAll(res.Body)

		result := &sendResult{
			StatusCode: res.StatusCode,
			Headers:    map[string][]string{},
			Cookies:    map[string]*http.Cookie{},
			Raw:        string(bodyRaw),
		}

		for k, v := range res.Header {
			result.Headers[k] = v
		}

		for _, v := range res.Cookies() {
			result.Cookies[v.Name] = v
		}

		if len(result.Raw) != 0 {
			// try as map
			result.JSON = map[string]any{}
			if err := json.Unmarshal(bodyRaw, &result.JSON); err != nil {
				// try as slice
				result.JSON = []any{}
				if err := json.Unmarshal(bodyRaw, &result.JSON); err != nil {
					result.JSON = nil
				}
			}
		}

		return result, nil
	})
}

// -------------------------------------------------------------------

// normalizeException checks if the provided error is a goja.Exception
// and attempts to return its underlying Go error.
//
// note: using just goja.Exception.Unwrap() is insufficient and may falsely result in nil.
func normalizeException(err error) error {
	if err == nil {
		return nil
	}

	jsException, ok := err.(*goja.Exception)
	if !ok {
		return err // no exception
	}

	switch v := jsException.Value().Export().(type) {
	case error:
		err = v
	case map[string]any: // goja.GoError
		if vErr, ok := v["value"].(error); ok {
			err = vErr
		}
	}

	return err
}

var cachedFactoryFuncTypes = store.New[string, reflect.Type](nil)

// registerFactoryAsConstructor registers the factory function as native JS constructor.
//
// If there is missing or nil arguments, their type zero value is used.
func registerFactoryAsConstructor(vm *goja.Runtime, constructorName string, factoryFunc any) {
	rv := reflect.ValueOf(factoryFunc)
	rt := cachedFactoryFuncTypes.GetOrSet(constructorName, func() reflect.Type {
		return reflect.TypeOf(factoryFunc)
	})
	totalArgs := rt.NumIn()

	vm.Set(constructorName, func(call goja.ConstructorCall) *goja.Object {
		args := make([]reflect.Value, totalArgs)

		for i := 0; i < totalArgs; i++ {
			v := call.Argument(i).Export()

			// use the arg type zero value
			if v == nil {
				args[i] = reflect.New(rt.In(i)).Elem()
			} else if number, ok := v.(int64); ok {
				// goja uses int64 for "int"-like numbers but we rarely do that and use int most of the times
				// (at later stage we can use reflection on the arguments to validate the types in case this is not sufficient anymore)
				args[i] = reflect.ValueOf(int(number))
			} else {
				args[i] = reflect.ValueOf(v)
			}
		}

		result := rv.Call(args)

		if len(result) != 1 {
			panic("the factory function should return only 1 item")
		}

		value := vm.ToValue(result[0].Interface()).(*goja.Object)
		value.SetPrototype(call.This.Prototype())

		return value
	})
}

// structConstructor wraps the provided struct with a native JS constructor.
//
// If the constructor argument is a map, each entry of the map will be loaded into the wrapped goja.Object.
func structConstructor(vm *goja.Runtime, call goja.ConstructorCall, instance any) *goja.Object {
	data, _ := call.Argument(0).Export().(map[string]any)

	instanceValue := vm.ToValue(instance).(*goja.Object)
	for k, v := range data {
		instanceValue.Set(k, v)
	}

	instanceValue.SetPrototype(call.This.Prototype())

	return instanceValue
}

// structConstructorUnmarshal wraps the provided struct with a native JS constructor.
//
// The constructor first argument will be loaded via json.Unmarshal into the instance.
func structConstructorUnmarshal(vm *goja.Runtime, call goja.ConstructorCall, instance any) *goja.Object {
	if data := call.Argument(0).Export(); data != nil {
		if raw, err := json.Marshal(data); err == nil {
			_ = json.Unmarshal(raw, instance)
		}
	}

	instanceValue := vm.ToValue(instance).(*goja.Object)
	instanceValue.SetPrototype(call.This.Prototype())

	return instanceValue
}

var cachedDynamicModels = store.New[string, *dynamicModelType](nil)

// newDynamicModel creates a new dynamic struct with fields based
// on the specified "shape".
//
// Example:
//
//	m := newDynamicModel(map[string]any{
//		"title": "",
//		"total": 0,
//	})
func newDynamicModel(shape map[string]any) any {
	var modelType *dynamicModelType

	shapeRaw, err := json.Marshal(shape)
	if err != nil {
		modelType = getDynamicModelStruct(shape)
	} else {
		modelType = cachedDynamicModels.GetOrSet(string(shapeRaw), func() *dynamicModelType {
			return getDynamicModelStruct(shape)
		})
	}

	rvShapeValues := make([]reflect.Value, len(modelType.shapeValues))
	for i, v := range modelType.shapeValues {
		rvShapeValues[i] = reflect.ValueOf(v)
	}

	elem := reflect.New(modelType.structType).Elem()

	for i, v := range rvShapeValues {
		elem.Field(i).Set(v)
	}

	return elem.Addr().Interface()
}

type dynamicModelType struct {
	structType  reflect.Type
	shapeValues []any
}

func getDynamicModelStruct(shape map[string]any) *dynamicModelType {
	result := new(dynamicModelType)
	result.shapeValues = make([]any, 0, len(shape))

	structFields := make([]reflect.StructField, 0, len(shape))

	for k, v := range shape {
		vt := reflect.TypeOf(v)

		switch kind := vt.Kind(); kind {
		case reflect.Map:
			raw, _ := json.Marshal(v)
			newV := types.JSONMap[any]{}
			newV.Scan(raw)
			v = newV
			vt = reflect.TypeOf(v)
		case reflect.Slice, reflect.Array:
			raw, _ := json.Marshal(v)
			newV := types.JSONArray[any]{}
			newV.Scan(raw)
			v = newV
			vt = reflect.TypeOf(newV)
		}

		result.shapeValues = append(result.shapeValues, v)

		structFields = append(structFields, reflect.StructField{
			Name: inflector.UcFirst(k), // ensures that the field is exportable
			Type: vt,
			Tag:  reflect.StructTag(`db:"` + k + `" json:"` + k + `" form:"` + k + `"`),
		})
	}

	result.structType = reflect.StructOf(structFields)

	return result
}
