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
	"strings"
	"time"

	"github.com/dop251/goja"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/security"
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
	excludeHooks := []string{"OnBeforeServe"}

	for i := 0; i < totalMethods; i++ {
		method := appType.Method(i)
		if !strings.HasPrefix(method.Name, "On") || list.ExistInSlice(method.Name, excludeHooks) {
			continue // not a hook or excluded
		}

		jsName := fm.MethodName(appType, method)

		// register the hook to the loader
		loader.Set(jsName, func(callback string, tags ...string) {
			pr := goja.MustCompile("", "{("+callback+").apply(undefined, __args)}", true)

			tagsAsValues := make([]reflect.Value, len(tags))
			for i, tag := range tags {
				tagsAsValues[i] = reflect.ValueOf(tag)
			}

			hookInstance := appValue.MethodByName(method.Name).Call(tagsAsValues)[0]
			addFunc := hookInstance.MethodByName("Add")

			handlerType := addFunc.Type().In(0)

			handler := reflect.MakeFunc(handlerType, func(args []reflect.Value) (results []reflect.Value) {
				handlerArgs := make([]any, len(args))
				for i, arg := range args {
					handlerArgs[i] = arg.Interface()
				}

				err := executors.run(func(executor *goja.Runtime) error {
					executor.Set("__args", handlerArgs)
					res, err := executor.RunProgram(pr)
					executor.Set("__args", goja.Undefined())

					// check for returned error or false
					if res != nil {
						switch v := res.Export().(type) {
						case error:
							return v
						case bool:
							if !v {
								return hook.StopPropagation
							}
						}
					}

					return err
				})

				return []reflect.Value{reflect.ValueOf(&err).Elem()}
			})

			// register the wrapped hook handler
			addFunc.Call([]reflect.Value{handler})
		})
	}
}

func cronBinds(app core.App, loader *goja.Runtime, executors *vmsPool) {
	scheduler := cron.New()

	var wasServeTriggered bool

	loader.Set("cronAdd", func(jobId, cronExpr, handler string) {
		pr := goja.MustCompile("", "{("+handler+").apply(undefined)}", true)

		err := scheduler.Add(jobId, cronExpr, func() {
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

		// start the ticker (if not already)
		if wasServeTriggered && scheduler.Total() > 0 && !scheduler.HasStarted() {
			scheduler.Start()
		}
	})

	loader.Set("cronRemove", func(jobId string) {
		scheduler.Remove(jobId)

		// stop the ticker if there are no other jobs
		if scheduler.Total() == 0 {
			scheduler.Stop()
		}
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// start the ticker (if not already)
		if scheduler.Total() > 0 && !scheduler.HasStarted() {
			scheduler.Start()
		}

		wasServeTriggered = true

		return nil
	})
}

func routerBinds(app core.App, loader *goja.Runtime, executors *vmsPool) {
	loader.Set("routerAdd", func(method string, path string, handler goja.Value, middlewares ...goja.Value) {
		wrappedMiddlewares, err := wrapMiddlewares(executors, middlewares...)
		if err != nil {
			panic("[routerAdd] failed to wrap middlewares: " + err.Error())
		}

		wrappedHandler, err := wrapHandler(executors, handler)
		if err != nil {
			panic("[routerAdd] failed to wrap handler: " + err.Error())
		}

		app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			e.Router.Add(strings.ToUpper(method), path, wrappedHandler, wrappedMiddlewares...)

			return nil
		})
	})

	loader.Set("routerUse", func(middlewares ...goja.Value) {
		wrappedMiddlewares, err := wrapMiddlewares(executors, middlewares...)
		if err != nil {
			panic("[routerUse] failed to wrap middlewares: " + err.Error())
		}

		app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			e.Router.Use(wrappedMiddlewares...)
			return nil
		})
	})

	loader.Set("routerPre", func(middlewares ...goja.Value) {
		wrappedMiddlewares, err := wrapMiddlewares(executors, middlewares...)
		if err != nil {
			panic("[routerPre] failed to wrap middlewares: " + err.Error())
		}

		app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			e.Router.Pre(wrappedMiddlewares...)
			return nil
		})
	})
}

func wrapHandler(executors *vmsPool, handler goja.Value) (echo.HandlerFunc, error) {
	if handler == nil {
		return nil, errors.New("handler must be non-nil")
	}

	switch h := handler.Export().(type) {
	case echo.HandlerFunc:
		// "native" handler - no need to wrap
		return h, nil
	case func(goja.FunctionCall) goja.Value, string:
		pr := goja.MustCompile("", "{("+handler.String()+").apply(undefined, __args)}", true)

		wrappedHandler := func(c echo.Context) error {
			return executors.run(func(executor *goja.Runtime) error {
				executor.Set("__args", []any{c})
				res, err := executor.RunProgram(pr)
				executor.Set("__args", goja.Undefined())

				// check for returned error
				if res != nil {
					if v, ok := res.Export().(error); ok {
						return v
					}
				}

				return err
			})
		}

		return wrappedHandler, nil
	default:
		return nil, errors.New("unsupported goja handler type")
	}
}

func wrapMiddlewares(executors *vmsPool, rawMiddlewares ...goja.Value) ([]echo.MiddlewareFunc, error) {
	wrappedMiddlewares := make([]echo.MiddlewareFunc, len(rawMiddlewares))

	for i, m := range rawMiddlewares {
		if m == nil {
			return nil, errors.New("middleware func must be non-nil")
		}

		switch v := m.Export().(type) {
		case echo.MiddlewareFunc:
			// "native" middleware - no need to wrap
			wrappedMiddlewares[i] = v
		case func(goja.FunctionCall) goja.Value, string:
			pr := goja.MustCompile("", "{(("+m.String()+").apply(undefined, __args)).apply(undefined, __args2)}", true)

			wrappedMiddlewares[i] = func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return executors.run(func(executor *goja.Runtime) error {
						executor.Set("__args", []any{next})
						executor.Set("__args2", []any{c})
						res, err := executor.RunProgram(pr)
						executor.Set("__args", goja.Undefined())
						executor.Set("__args2", goja.Undefined())

						// check for returned error
						if res != nil {
							if v, ok := res.Export().(error); ok {
								return v
							}
						}

						return err
					})
				}
			}
		default:
			return nil, errors.New("unsupported goja middleware type")
		}
	}

	return wrappedMiddlewares, nil
}

func baseBinds(vm *goja.Runtime) {
	vm.SetFieldNameMapper(FieldMapper{})

	vm.Set("readerToString", func(r io.Reader, maxBytes int) (string, error) {
		if maxBytes == 0 {
			maxBytes = rest.DefaultMaxMemory
		}

		limitReader := io.LimitReader(r, int64(maxBytes))

		bodyBytes, readErr := io.ReadAll(limitReader)
		if readErr != nil {
			return "", readErr
		}

		return string(bodyBytes), nil
	})

	vm.Set("sleep", func(milliseconds int64) {
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
	})

	vm.Set("arrayOf", func(model any) any {
		mt := reflect.TypeOf(model)
		st := reflect.SliceOf(mt)
		elem := reflect.New(st).Elem()

		return elem.Addr().Interface()
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
		var instance *models.Record

		collection, ok := call.Argument(0).Export().(*models.Collection)
		if ok {
			instance = models.NewRecord(collection)
			data, ok := call.Argument(1).Export().(map[string]any)
			if ok {
				instance.Load(data)
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
		return structConstructorUnmarshal(vm, call, instance)
	})

	vm.Set("Admin", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.Admin{}
		return structConstructorUnmarshal(vm, call, instance)
	})

	vm.Set("Schema", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.Schema{}
		return structConstructorUnmarshal(vm, call, instance)
	})

	vm.Set("SchemaField", func(call goja.ConstructorCall) *goja.Object {
		instance := &schema.SchemaField{}
		return structConstructorUnmarshal(vm, call, instance)
	})

	vm.Set("MailerMessage", func(call goja.ConstructorCall) *goja.Object {
		instance := &mailer.Message{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("Command", func(call goja.ConstructorCall) *goja.Object {
		instance := &cobra.Command{}
		return structConstructor(vm, call, instance)
	})

	vm.Set("RequestInfo", func(call goja.ConstructorCall) *goja.Object {
		instance := &models.RequestInfo{Context: models.RequestInfoContextDefault}
		return structConstructor(vm, call, instance)
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

	vm.Set("Dao", func(call goja.ConstructorCall) *goja.Object {
		concurrentDB, _ := call.Argument(0).Export().(dbx.Builder)
		if concurrentDB == nil {
			panic("[Dao] missing required Dao(concurrentDB, [nonconcurrentDB]) argument")
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

	// admin
	obj.Set("sendAdminPasswordReset", mails.SendAdminPasswordReset)

	// record
	obj.Set("sendRecordPasswordReset", mails.SendRecordPasswordReset)
	obj.Set("sendRecordVerification", mails.SendRecordVerification)
	obj.Set("sendRecordChangeEmail", mails.SendRecordChangeEmail)
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

	// crypto
	obj.Set("md5", security.MD5)
	obj.Set("sha256", security.SHA256)
	obj.Set("sha512", security.SHA512)
	obj.Set("hs256", security.HS256)
	obj.Set("hs512", security.HS512)
	obj.Set("equal", security.Equal)

	// random
	obj.Set("randomString", security.RandomString)
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
	obj.Set("createJWT", security.NewJWT)

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
	obj.Set("fileFromUrl", func(url string, secTimeout int) (*filesystem.File, error) {
		if secTimeout == 0 {
			secTimeout = 120
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(secTimeout)*time.Second)
		defer cancel()

		return filesystem.NewFileFromUrl(ctx, url)
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

	obj.Set("staticDirectoryHandler", func(dir string, indexFallback bool) echo.HandlerFunc {
		return apis.StaticDirectoryHandler(os.DirFS(dir), indexFallback)
	})

	// middlewares
	obj.Set("requireGuestOnly", apis.RequireGuestOnly)
	obj.Set("requireRecordAuth", apis.RequireRecordAuth)
	obj.Set("requireAdminAuth", apis.RequireAdminAuth)
	obj.Set("requireAdminAuthOnlyIfAny", apis.RequireAdminAuthOnlyIfAny)
	obj.Set("requireAdminOrRecordAuth", apis.RequireAdminOrRecordAuth)
	obj.Set("requireAdminOrOwnerAuth", apis.RequireAdminOrOwnerAuth)
	obj.Set("activityLogger", apis.ActivityLogger)
	obj.Set("gzip", middleware.Gzip)
	obj.Set("bodyLimit", middleware.BodyLimit)

	// record helpers
	obj.Set("requestInfo", apis.RequestInfo)
	obj.Set("recordAuthResponse", apis.RecordAuthResponse)
	obj.Set("enrichRecord", apis.EnrichRecord)
	obj.Set("enrichRecords", apis.EnrichRecords)

	// api errors
	registerFactoryAsConstructor(vm, "ApiError", apis.NewApiError)
	registerFactoryAsConstructor(vm, "NotFoundError", apis.NewNotFoundError)
	registerFactoryAsConstructor(vm, "BadRequestError", apis.NewBadRequestError)
	registerFactoryAsConstructor(vm, "ForbiddenError", apis.NewForbiddenError)
	registerFactoryAsConstructor(vm, "UnauthorizedError", apis.NewUnauthorizedError)
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
		Json       any                     `json:"json"`
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

		// @todo consider removing during the refactoring
		//
		// fallback to json content-type
		if req.Header.Get("content-type") == "" {
			req.Header.Set("content-type", "application/json")
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
			result.Json = map[string]any{}
			if err := json.Unmarshal(bodyRaw, &result.Json); err != nil {
				// try as slice
				result.Json = []any{}
				if err := json.Unmarshal(bodyRaw, &result.Json); err != nil {
					result.Json = nil
				}
			}
		}

		return result, nil
	})
}

// -------------------------------------------------------------------

// registerFactoryAsConstructor registers the factory function as native JS constructor.
//
// If there is missing or nil arguments, their type zero value is used.
func registerFactoryAsConstructor(vm *goja.Runtime, constructorName string, factoryFunc any) {
	rv := reflect.ValueOf(factoryFunc)
	rt := reflect.TypeOf(factoryFunc)
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
			json.Unmarshal(raw, instance)
		}
	}

	instanceValue := vm.ToValue(instance).(*goja.Object)
	instanceValue.SetPrototype(call.This.Prototype())

	return instanceValue
}

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
	shapeValues := make([]reflect.Value, 0, len(shape))
	structFields := make([]reflect.StructField, 0, len(shape))

	for k, v := range shape {
		vt := reflect.TypeOf(v)

		switch kind := vt.Kind(); kind {
		case reflect.Map:
			raw, _ := json.Marshal(v)
			newV := types.JsonMap{}
			newV.Scan(raw)
			v = newV
			vt = reflect.TypeOf(v)
		case reflect.Slice, reflect.Array:
			raw, _ := json.Marshal(v)
			newV := types.JsonArray[any]{}
			newV.Scan(raw)
			v = newV
			vt = reflect.TypeOf(newV)
		}

		shapeValues = append(shapeValues, reflect.ValueOf(v))

		structFields = append(structFields, reflect.StructField{
			Name: inflector.UcFirst(k), // ensures that the field is exportable
			Type: vt,
			Tag:  reflect.StructTag(`db:"` + k + `" json:"` + k + `" form:"` + k + `"`),
		})
	}

	st := reflect.StructOf(structFields)
	elem := reflect.New(st).Elem()

	for i, v := range shapeValues {
		elem.Field(i).Set(v)
	}

	return elem.Addr().Interface()
}
