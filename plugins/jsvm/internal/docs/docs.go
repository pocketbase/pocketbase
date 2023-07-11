package main

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/tygoja"
)

const heading = `
// -------------------------------------------------------------------
// routerBinds
// -------------------------------------------------------------------

/**
 * RouterAdd registers a new route definition.
 *
 * Example:
 *
 * ` + "```" + `js
 * routerAdd("GET", "/hello", (c) => {
 *     return c.json(200, {"message": "Hello!"})
 * }, $apis.requireAdminOrRecordAuth())
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function routerAdd(
  method: string,
  path: string,
  handler: echo.HandlerFunc,
  ...middlewares: Array<string|echo.MiddlewareFunc>,
): void;

/**
 * RouterUse registers one or more global middlewares that are executed
 * along the handler middlewares after a matching route is found.
 *
 * Example:
 *
 * ` + "```" + `js
 * routerUse((next) => {
 *     return (c) => {
 *         console.log(c.Path())
 *         return next(c)
 *     }
 * })
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function routerUse(...middlewares: Array<string|echo.MiddlewareFunc>): void;

/**
 * RouterPre registers one or more global middlewares that are executed
 * BEFORE the router processes the request. It is usually used for making
 * changes to the request properties, for example, adding or removing
 * a trailing slash or adding segments to a path so it matches a route.
 *
 * NB! Since the router will not have processed the request yet,
 * middlewares registered at this level won't have access to any path
 * related APIs from echo.Context.
 *
 * Example:
 *
 * ` + "```" + `js
 * routerPre((next) => {
 *     return (c) => {
 *         console.log(c.request().url)
 *         return next(c)
 *     }
 * })
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function routerPre(...middlewares: Array<string|echo.MiddlewareFunc>): void;

// -------------------------------------------------------------------
// baseBinds
// -------------------------------------------------------------------

// skip on* hook methods as they are registered via the global on* method
type appWithoutHooks = Omit<pocketbase.PocketBase, ` + "`on${string}`" + `>

/**
 * ` + "`$app`" + ` is the current running PocketBase instance that is globally
 * available in each .pb.js file.
 *
 * @namespace
 * @group PocketBase
 */
declare var $app: appWithoutHooks

/**
 * arrayOf creates a placeholder array of the specified models.
 * Usually used to populate DB result into an array of models.
 *
 * Example:
 *
 * ` + "```" + `js
 * const records = arrayOf(new Record)
 *
 * $app.dao().recordQuery(collection).limit(10).all(records)
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare function arrayOf<T>(model: T): Array<T>;

/**
 * DynamicModel creates a new dynamic model with fields from the provided data shape.
 *
 * Example:
 *
 * ` + "```" + `js
 * const model = new DynamicModel({
 *     name:  ""
 *     age:    0,
 *     active: false,
 *     roles:  [],
 *     meta:   {}
 * })
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class DynamicModel {
  constructor(shape?: { [key:string]: any })
}

/**
 * Record model class.
 *
 * ` + "```" + `js
 * const collection = $app.dao().findCollectionByNameOrId("article")
 *
 * const record = new Record(collection, {
 *     title: "Lorem ipsum"
 * })
 *
 * // or set field values after the initialization
 * record.set("description", "...")
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare const Record: {
  new(collection?: models.Collection, data?: { [key:string]: any }): models.Record

  // note: declare as "newable" const due to conflict with the Record TS utility type
}

interface Collection extends models.Collection{} // merge
/**
 * Collection model class.
 *
 * ` + "```" + `js
 * const collection = new Collection({
 *     name:       "article",
 *     type:       "base",
 *     listRule:   "@request.auth.id != '' || status = 'public'",
 *     viewRule:   "@request.auth.id != '' || status = 'public'",
 *     deleteRule: "@request.auth.id != ''",
 *     schema: [
 *         {
 *             name: "title",
 *             type: "text",
 *             required: true,
 *             options: { min: 6, max: 100 },
 *         },
 *         {
 *             name: "description",
 *             type: "text",
 *         },
 *     ]
 * })
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Collection implements models.Collection {
  constructor(data?: Partial<models.Collection>)
}

interface Admin extends models.Admin{} // merge
/**
 * Admin model class.
 *
 * ` + "```" + `js
 * const admin = new Admin()
 * admin.email = "test@example.com"
 * admin.setPassword(1234567890)
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Admin implements models.Admin {
  constructor(data?: Partial<models.Admin>)
}

interface Schema extends schema.Schema{} // merge
/**
 * Schema model class, usually used to define the Collection.schema field.
 *
 * @group PocketBase
 */
declare class Schema implements schema.Schema {
  constructor(data?: Partial<schema.Schema>)
}

interface SchemaField extends schema.SchemaField{} // merge
/**
 * SchemaField model class, usually used as part of the Schema model.
 *
 * @group PocketBase
 */
declare class SchemaField implements schema.SchemaField {
  constructor(data?: Partial<schema.SchemaField>)
}

interface MailerMessage extends mailer.Message{} // merge
/**
 * MailerMessage defines a single email message.
 *
 * ` + "```" + `js
 * const message = new MailerMessage({
 *     from: {
 *         address: $app.settings().meta.senderAddress,
 *         name:    $app.settings().meta.senderName,
 *     },
 *     to:      [{address: "test@example.com"}],
 *     subject: "YOUR_SUBJECT...",
 *     html:    "YOUR_HTML_BODY...",
 * })
 *
 * $app.newMailClient().send(message)
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class MailerMessage implements mailer.Message {
  constructor(message?: Partial<mailer.Message>)
}

interface Command extends cobra.Command{} // merge
/**
 * Command defines a single console command.
 *
 * ` + "```" + `js
 * const command = new Command({
 *     use: "hello",
 *     run: (cmd, args) => { console.log("Hello world!") },
 * })
 *
 * $app.rootCmd.addCommand(command);
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Command implements cobra.Command {
  constructor(cmd?: Partial<cobra.Command>)
}

interface ValidationError extends ozzo_validation.Error{} // merge
/**
 * ValidationError defines a single formatted data validation error,
 * usually used as part of a error response.
 *
 * ` + "```" + `js
 * new ValidationError("invalid_title", "Title is not valid")
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class ValidationError implements ozzo_validation.Error {
  constructor(code?: string, message?: string)
}

interface Dao extends daos.Dao{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class Dao implements daos.Dao {
  constructor(concurrentDB?: dbx.Builder, nonconcurrentDB?: dbx.Builder)
}

// -------------------------------------------------------------------
// dbxBinds
// -------------------------------------------------------------------

/**
 * ` + "`$dbx`" + ` defines common utility for working with the DB abstraction.
 * For examples and guides please check the [Database guide](https://pocketbase.io/docs/js-database).
 *
 * @group PocketBase
 */
declare namespace $dbx {
  /**
   * {@inheritDoc dbx.HashExp}
   */
  export function hashExp(pairs: { [key:string]: any }): dbx.Expression

  let _in: dbx._in
  export { _in as in }

  export let exp:        dbx.newExp
  export let not:        dbx.not
  export let and:        dbx.and
  export let or:         dbx.or
  export let notIn:      dbx.notIn
  export let like:       dbx.like
  export let orLike:     dbx.orLike
  export let notLike:    dbx.notLike
  export let orNotLike:  dbx.orNotLike
  export let exists:     dbx.exists
  export let notExists:  dbx.notExists
  export let between:    dbx.between
  export let notBetween: dbx.notBetween
}

// -------------------------------------------------------------------
// tokensBinds
// -------------------------------------------------------------------

/**
 * ` + "`" + `$tokens` + "`" + ` defines high level helpers to generate
 * various admins and auth records tokens (auth, forgotten password, etc.).
 *
 * For more control over the generated token, you can check ` + "`" + `$security` + "`" + `.
 *
 * @group PocketBase
 */
declare namespace $tokens {
  let adminAuthToken:           tokens.newAdminAuthToken
  let adminResetPasswordToken:  tokens.newAdminResetPasswordToken
  let adminFileToken:           tokens.newAdminFileToken
  let recordAuthToken:          tokens.newRecordAuthToken
  let recordVerifyToken:        tokens.newRecordVerifyToken
  let recordResetPasswordToken: tokens.newRecordResetPasswordToken
  let recordChangeEmailToken:   tokens.newRecordChangeEmailToken
  let recordFileToken:          tokens.newRecordFileToken
}

// -------------------------------------------------------------------
// securityBinds
// -------------------------------------------------------------------

/**
 * ` + "`" + `$security` + "`" + ` defines low level helpers for creating
 * and parsing JWTs, random string generation, AES encryption, etc.
 *
 * @group PocketBase
 */
declare namespace $security {
  let randomString:                   security.randomString
  let randomStringWithAlphabet:       security.randomStringWithAlphabet
  let pseudorandomString:             security.pseudorandomString
  let pseudorandomStringWithAlphabet: security.pseudorandomStringWithAlphabet
  let parseUnverifiedJWT:             security.parseUnverifiedJWT
  let parseJWT:                       security.parseJWT
  let createJWT:                      security.newJWT
  let encrypt:                        security.encrypt
  let decrypt:                        security.decrypt
}

// -------------------------------------------------------------------
// filesystemBinds
// -------------------------------------------------------------------

/**
 * ` + "`" + `$filesystem` + "`" + ` defines common helpers for working
 * with the PocketBase filesystem abstraction.
 *
 * @group PocketBase
 */
declare namespace $filesystem {
  let fileFromPath:      filesystem.newFileFromPath
  let fileFromBytes:     filesystem.newFileFromBytes
  let fileFromMultipart: filesystem.newFileFromMultipart
}

// -------------------------------------------------------------------
// formsBinds
// -------------------------------------------------------------------

interface AdminLoginForm extends forms.AdminLogin{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminLoginForm implements forms.AdminLogin {
  constructor(app: core.App)
}

interface AdminPasswordResetConfirmForm extends forms.AdminPasswordResetConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminPasswordResetConfirmForm implements forms.AdminPasswordResetConfirm {
  constructor(app: core.App)
}

interface AdminPasswordResetRequestForm extends forms.AdminPasswordResetRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminPasswordResetRequestForm implements forms.AdminPasswordResetRequest {
  constructor(app: core.App)
}

interface AdminUpsertForm extends forms.AdminUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminUpsertForm implements forms.AdminUpsert {
  constructor(app: core.App, admin: models.Admin)
}

interface AppleClientSecretCreateForm extends forms.AppleClientSecretCreate{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AppleClientSecretCreateForm implements forms.AppleClientSecretCreate {
  constructor(app: core.App)
}

interface CollectionUpsertForm extends forms.CollectionUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class CollectionUpsertForm implements forms.CollectionUpsert {
  constructor(app: core.App, collection: models.Collection)
}

interface CollectionsImportForm extends forms.CollectionsImport{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class CollectionsImportForm implements forms.CollectionsImport {
  constructor(app: core.App)
}

interface RealtimeSubscribeForm extends forms.RealtimeSubscribe{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RealtimeSubscribeForm implements forms.RealtimeSubscribe {}

interface RecordEmailChangeConfirmForm extends forms.RecordEmailChangeConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordEmailChangeConfirmForm implements forms.RecordEmailChangeConfirm {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordEmailChangeRequestForm extends forms.RecordEmailChangeRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordEmailChangeRequestForm implements forms.RecordEmailChangeRequest {
  constructor(app: core.App, record: models.Record)
}

interface RecordOAuth2LoginForm extends forms.RecordOAuth2Login{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordOAuth2LoginForm implements forms.RecordOAuth2Login {
  constructor(app: core.App, collection: models.Collection, optAuthRecord?: models.Record)
}

interface RecordPasswordLoginForm extends forms.RecordPasswordLogin{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordPasswordLoginForm implements forms.RecordPasswordLogin {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordPasswordResetConfirmForm extends forms.RecordPasswordResetConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordPasswordResetConfirmForm implements forms.RecordPasswordResetConfirm {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordPasswordResetRequestForm extends forms.RecordPasswordResetRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordPasswordResetRequestForm implements forms.RecordPasswordResetRequest {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordUpsertForm extends forms.RecordUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordUpsertForm implements forms.RecordUpsert {
  constructor(app: core.App, record: models.Record)
}

interface RecordVerificationConfirmForm extends forms.RecordVerificationConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordVerificationConfirmForm implements forms.RecordVerificationConfirm {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordVerificationRequestForm extends forms.RecordVerificationRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordVerificationRequestForm implements forms.RecordVerificationRequest {
  constructor(app: core.App, collection: models.Collection)
}

interface SettingsUpsertForm extends forms.SettingsUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class SettingsUpsertForm implements forms.SettingsUpsert {
  constructor(app: core.App)
}

interface TestEmailSendForm extends forms.TestEmailSend{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class TestEmailSendForm implements forms.TestEmailSend {
  constructor(app: core.App)
}

interface TestS3FilesystemForm extends forms.TestS3Filesystem{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class TestS3FilesystemForm implements forms.TestS3Filesystem {
  constructor(app: core.App)
}

// -------------------------------------------------------------------
// apisBinds
// -------------------------------------------------------------------

interface ApiError extends apis.ApiError{} // merge
/**
 * @inheritDoc
 *
 * @group PocketBase
 */
declare class ApiError implements apis.ApiError {
  constructor(status?: number, message?: string, data?: any)
}

interface NotFoundError extends apis.ApiError{} // merge
/**
 * NotFounderor returns 404 ApiError.
 *
 * @group PocketBase
 */
declare class NotFoundError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

interface BadRequestError extends apis.ApiError{} // merge
/**
 * BadRequestError returns 400 ApiError.
 *
 * @group PocketBase
 */
declare class BadRequestError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

interface ForbiddenError extends apis.ApiError{} // merge
/**
 * ForbiddenError returns 403 ApiError.
 *
 * @group PocketBase
 */
declare class ForbiddenError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

interface UnauthorizedError extends apis.ApiError{} // merge
/**
 * UnauthorizedError returns 401 ApiError.
 *
 * @group PocketBase
 */
declare class UnauthorizedError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

/**
 * ` + "`" + `$apis` + "`" + ` defines commonly used PocketBase api helpers and middlewares.
 *
 * @group PocketBase
 */
declare namespace $apis {
  let requireRecordAuth:         apis.requireRecordAuth
  let requireAdminAuth:          apis.requireAdminAuth
  let requireAdminAuthOnlyIfAny: apis.requireAdminAuthOnlyIfAny
  let requireAdminOrRecordAuth:  apis.requireAdminOrRecordAuth
  let requireAdminOrOwnerAuth:   apis.requireAdminOrOwnerAuth
  let activityLogger:            apis.activityLogger
  let requestData:               apis.requestData
  let recordAuthResponse:        apis.recordAuthResponse
  let enrichRecord:              apis.enrichRecord
  let enrichRecords:             apis.enrichRecords
}

// -------------------------------------------------------------------
// httpClientBinds
// -------------------------------------------------------------------

declare namespace $http {
  /**
   * Sends a single HTTP request (_currently only json and plain text requests_).
   *
   * @group PocketBase
   */
  function send(params: {
    url:     string,
    method?:  string, // default to "GET"
    data?:    { [key:string]: any },
    headers?: { [key:string]: string },
    timeout?: number // default to 120
  }): {
    statusCode: number
    raw:        string
    json:       any
  };
}

// -------------------------------------------------------------------
// migrate only
// -------------------------------------------------------------------

/**
 * Migrate defines a single migration upgrade/downgrade action.
 *
 * _Note that this method is available only in pb_migrations context._
 *
 * @group PocketBase
 */
declare function migrate(
  up: (db: dbx.Builder) => void,
  down?: (db: dbx.Builder) => void
): void;
`

var mapper = &jsvm.FieldMapper{}

func main() {
	declarations := heading + hooksDeclarations()

	gen := tygoja.New(tygoja.Config{
		Packages: map[string][]string{
			"github.com/go-ozzo/ozzo-validation/v4":             {"Error"},
			"github.com/pocketbase/dbx":                         {"*"},
			"github.com/pocketbase/pocketbase/tools/security":   {"*"},
			"github.com/pocketbase/pocketbase/tools/filesystem": {"*"},
			"github.com/pocketbase/pocketbase/tokens":           {"*"},
			"github.com/pocketbase/pocketbase/apis":             {"*"},
			"github.com/pocketbase/pocketbase/forms":            {"*"},
			"github.com/pocketbase/pocketbase":                  {"*"},
		},
		FieldNameFormatter: func(s string) string {
			return mapper.FieldName(nil, reflect.StructField{Name: s})
		},
		MethodNameFormatter: func(s string) string {
			return mapper.MethodName(nil, reflect.Method{Name: s})
		},
		Indent:               " ", // use only a single space to reduce slight the size
		WithPackageFunctions: true,
		Heading:              declarations,
	})

	result, err := gen.Generate()
	if err != nil {
		log.Fatal(err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get the current docs directory")
	}

	parentDir := filepath.Dir(filename)
	typesFile := filepath.Join(parentDir, "generated", "types.d.ts")

	if err := os.WriteFile(typesFile, []byte(result), 0644); err != nil {
		log.Fatal(err)
	}
}

func hooksDeclarations() string {
	var result strings.Builder

	excluded := []string{"OnBeforeServe"}
	appType := reflect.TypeOf(struct{ core.App }{})
	totalMethods := appType.NumMethod()

	for i := 0; i < totalMethods; i++ {
		method := appType.Method(i)
		if !strings.HasPrefix(method.Name, "On") || list.ExistInSlice(method.Name, excluded) {
			continue // not a hook or excluded
		}

		hookType := method.Type.Out(0)

		withTags := strings.HasPrefix(hookType.String(), "*hook.TaggedHook")

		addMethod, ok := hookType.MethodByName("Add")
		if !ok {
			continue
		}

		addHanlder := addMethod.Type.In(1)
		eventTypeName := strings.TrimPrefix(addHanlder.In(0).String(), "*")

		jsName := mapper.MethodName(appType, method)
		result.WriteString("/** @group PocketBase */")
		result.WriteString("declare function ")
		result.WriteString(jsName)
		result.WriteString("(handler: (e: ")
		result.WriteString(eventTypeName)
		result.WriteString(") => void")
		if withTags {
			result.WriteString(", ...tags: string[]")
		}
		result.WriteString("): void")
		result.WriteString("\n")
	}

	return result.String()
}
