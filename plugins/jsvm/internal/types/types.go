package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/tygoja"
)

const heading = `
// -------------------------------------------------------------------
// cronBinds
// -------------------------------------------------------------------

/**
 * CronAdd registers a new cron job.
 *
 * If a cron job with the specified name already exist, it will be
 * replaced with the new one.
 *
 * Example:
 *
 * ` + "```" + `js
 * // prints "Hello world!" on every 30 minutes
 * cronAdd("hello", "*\/30 * * * *", () => {
 *     console.log("Hello world!")
 * })
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function cronAdd(
  jobId:    string,
  cronExpr: string,
  handler:  () => void,
): void;

/**
 * CronRemove removes a single registered cron job by its name.
 *
 * Example:
 *
 * ` + "```" + `js
 * cronRemove("hello")
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function cronRemove(jobId: string): void;

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
 *         console.log(c.path())
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

/**
 * Global helper variable that contains the absolute path to the app pb_hooks directory.
 *
 * @group PocketBase
 */
declare var __hooks: string

// Utility type to exclude the on* hook methods from a type
// (hooks are separately generated as global methods).
//
// See https://www.typescriptlang.org/docs/handbook/2/mapped-types.html#key-remapping-via-as
type excludeHooks<Type> = {
    [Property in keyof Type as Exclude<Property, ` + "`on${string}`" + `>]: Type[Property]
};

// core.App without the on* hook methods
type CoreApp = excludeHooks<ORIGINAL_CORE_APP>

// pocketbase.PocketBase without the on* hook methods
type PocketBase = excludeHooks<ORIGINAL_POCKETBASE>

/**
 * ` + "`$app`" + ` is the current running PocketBase instance that is globally
 * available in each .pb.js file.
 *
 * _Note that this variable is available only in pb_hooks context._
 *
 * @namespace
 * @group PocketBase
 */
declare var $app: PocketBase

/**
 * ` + "`$template`" + ` is a global helper to load and cache HTML templates on the fly.
 *
 * The templates uses the standard Go [html/template](https://pkg.go.dev/html/template)
 * and [text/template](https://pkg.go.dev/text/template) package syntax.
 *
 * Example:
 *
 * ` + "```" + `js
 * const html = $template.loadFiles(
 *     "views/layout.html",
 *     "views/content.html",
 * ).render({"name": "John"})
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @namespace
 * @group PocketBase
 */
declare var $template: template.Registry

/**
 * readerToString reads the content of the specified io.Reader until
 * EOF or maxBytes are reached.
 *
 * If maxBytes is not specified it will read up to 32MB.
 *
 * Note that after this call the reader can't be used anymore.
 *
 * Example:
 *
 * ` + "```" + `js
 * const rawBody = readerToString(c.request().body)
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare function readerToString(reader: any, maxBytes?: number): string;

/**
 * sleep pauses the current goroutine for at least the specified user duration (in ms).
 * A zero or negative duration returns immediately.
 *
 * Example:
 *
 * ` + "```" + `js
 * sleep(250) // sleeps for 250ms
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare function sleep(milliseconds: number): void;

/**
 * arrayOf creates a placeholder array of the specified models.
 * Usually used to populate DB result into an array of models.
 *
 * Example:
 *
 * ` + "```" + `js
 * const records = arrayOf(new Record)
 *
 * $app.dao().recordQuery("articles").limit(10).all(records)
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
 * Example:
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

interface RequestInfo extends models.RequestInfo{} // merge
/**
 * RequestInfo defines a single models.RequestInfo instance, usually used
 * as part of various filter checks.
 *
 * Example:
 *
 * ` + "```" + `js
 * const authRecord = $app.dao().findAuthRecordByEmail("users", "test@example.com")
 *
 * const info = new RequestInfo({
 *     authRecord: authRecord,
 *     data:       {"name": 123},
 *     headers:    {"x-token": "..."},
 * })
 *
 * const record = $app.dao().findFirstRecordByData("articles", "slug", "hello")
 *
 * const canAccess = $app.dao().canAccessRecord(record, info, "@request.auth.id != '' && @request.data.name = 123")
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class RequestInfo implements models.RequestInfo {
  constructor(date?: Partial<models.RequestInfo>)
}

interface DateTime extends types.DateTime{} // merge
/**
 * DateTime defines a single DateTime type instance.
 *
 * Example:
 *
 * ` + "```" + `js
 * const dt0 = new DateTime() // now
 *
 * const dt1 = new DateTime('2023-07-01 00:00:00.000Z')
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class DateTime implements types.DateTime {
  constructor(date?: string)
}

interface ValidationError extends ozzo_validation.Error{} // merge
/**
 * ValidationError defines a single formatted data validation error,
 * usually used as part of an error response.
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

interface Cookie extends http.Cookie{} // merge
/**
 * A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an
 * HTTP response.
 *
 * Example:
 *
 * ` + "```" + `js
 * routerAdd("POST", "/example", (c) => {
 *     c.setCookie(new Cookie({
 *         name:     "example_name",
 *         value:    "example_value",
 *         path:     "/",
 *         domain:   "example.com",
 *         maxAge:   10,
 *         secure:   true,
 *         httpOnly: true,
 *         sameSite: 3,
 *     }))
 *
 *     return c.redirect(200, "/");
 * })
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Cookie implements http.Cookie {
  constructor(options?: Partial<http.Cookie>)
}

interface SubscriptionMessage extends subscriptions.Message{} // merge
/**
 * SubscriptionMessage defines a realtime subscription payload.
 *
 * Example:
 *
 * ` + "```" + `js
 * onRealtimeConnectRequest((e) => {
 *     e.client.send(new SubscriptionMessage({
 *         name: "example",
 *         data: '{"greeting": "Hello world"}'
 *     }))
 * })
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class SubscriptionMessage implements subscriptions.Message {
  constructor(options?: Partial<subscriptions.Message>)
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
// mailsBinds
// -------------------------------------------------------------------

/**
 * ` + "`" + `$mails` + "`" + ` defines helpers to send common
 * admins and auth records emails like verification, password reset, etc.
 *
 * @group PocketBase
 */
declare namespace $mails {
  let sendAdminPasswordReset:  mails.sendAdminPasswordReset
  let sendRecordPasswordReset: mails.sendRecordPasswordReset
  let sendRecordVerification:  mails.sendRecordVerification
  let sendRecordChangeEmail:   mails.sendRecordChangeEmail
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
  let encrypt:                        security.encrypt
  let decrypt:                        security.decrypt
  let hs256:                          security.hs256
  let hs512:                          security.hs512
  let equal:                          security.equal
  let md5:                            security.md5
  let sha256:                         security.sha256
  let sha512:                         security.sha512
  let createJWT:                      security.newJWT

  /**
   * {@inheritDoc security.parseUnverifiedJWT}
   */
  export function parseUnverifiedJWT(token: string): _TygojaDict

  /**
   * {@inheritDoc security.parseJWT}
   */
  export function parseJWT(token: string, verificationKey: string): _TygojaDict
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

  /**
   * fileFromUrl creates a new File from the provided url by
   * downloading the resource and creating a BytesReader.
   *
   * Example:
   *
   * ` + "```" + `js
   * // with default max timeout of 120sec
   * const file1 = $filesystem.fileFromUrl("https://...")
   *
   * // with custom timeout of 15sec
   * const file2 = $filesystem.fileFromUrl("https://...", 15)
   * ` + "```" + `
   */
  export function fileFromUrl(url: string, secTimeout?: number): filesystem.File
}

// -------------------------------------------------------------------
// filepathBinds
// -------------------------------------------------------------------

/**
 * ` + "`$filepath`" + ` defines common helpers for manipulating filename
 * paths in a way compatible with the target operating system-defined file paths.
 *
 * @group PocketBase
 */
declare namespace $filepath {
  export let base:      filepath.base
  export let clean:     filepath.clean
  export let dir:       filepath.dir
  export let ext:       filepath.ext
  export let fromSlash: filepath.fromSlash
  export let glob:      filepath.glob
  export let isAbs:     filepath.isAbs
  export let join:      filepath.join
  export let match:     filepath.match
  export let rel:       filepath.rel
  export let split:     filepath.split
  export let splitList: filepath.splitList
  export let toSlash:   filepath.toSlash
  export let walk:      filepath.walk
  export let walkDir:   filepath.walkDir
}

// -------------------------------------------------------------------
// osBinds
// -------------------------------------------------------------------

/**
 * ` + "`$os`" + ` defines common helpers for working with the OS level primitives
 * (eg. deleting directories, executing shell commands, etc.).
 *
 * @group PocketBase
 */
declare namespace $os {
  /**
   * Legacy alias for $os.cmd().
   */
  export let exec: exec.command

  /**
   * Prepares an external OS command.
   *
   * Example:
   *
   * ` + "```" + `js
   * // prepare the command to execute
   * const cmd = $os.cmd('ls', '-sl')
   *
   * // execute the command and return its standard output as string
   * const output = String.fromCharCode(...cmd.output());
   * ` + "```" + `
   */
  export let cmd: exec.command

  export let args:      os.args
  export let exit:      os.exit
  export let getenv:    os.getenv
  export let dirFS:     os.dirFS
  export let readFile:  os.readFile
  export let writeFile: os.writeFile
  export let readDir:   os.readDir
  export let tempDir:   os.tempDir
  export let truncate:  os.truncate
  export let getwd:     os.getwd
  export let mkdir:     os.mkdir
  export let mkdirAll:  os.mkdirAll
  export let rename:    os.rename
  export let remove:    os.remove
  export let removeAll: os.removeAll
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
  constructor(app: CoreApp)
}

interface AdminPasswordResetConfirmForm extends forms.AdminPasswordResetConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminPasswordResetConfirmForm implements forms.AdminPasswordResetConfirm {
  constructor(app: CoreApp)
}

interface AdminPasswordResetRequestForm extends forms.AdminPasswordResetRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminPasswordResetRequestForm implements forms.AdminPasswordResetRequest {
  constructor(app: CoreApp)
}

interface AdminUpsertForm extends forms.AdminUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AdminUpsertForm implements forms.AdminUpsert {
  constructor(app: CoreApp, admin: models.Admin)
}

interface AppleClientSecretCreateForm extends forms.AppleClientSecretCreate{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AppleClientSecretCreateForm implements forms.AppleClientSecretCreate {
  constructor(app: CoreApp)
}

interface CollectionUpsertForm extends forms.CollectionUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class CollectionUpsertForm implements forms.CollectionUpsert {
  constructor(app: CoreApp, collection: models.Collection)
}

interface CollectionsImportForm extends forms.CollectionsImport{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class CollectionsImportForm implements forms.CollectionsImport {
  constructor(app: CoreApp)
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
  constructor(app: CoreApp, collection: models.Collection)
}

interface RecordEmailChangeRequestForm extends forms.RecordEmailChangeRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordEmailChangeRequestForm implements forms.RecordEmailChangeRequest {
  constructor(app: CoreApp, record: models.Record)
}

interface RecordOAuth2LoginForm extends forms.RecordOAuth2Login{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordOAuth2LoginForm implements forms.RecordOAuth2Login {
  constructor(app: CoreApp, collection: models.Collection, optAuthRecord?: models.Record)
}

interface RecordPasswordLoginForm extends forms.RecordPasswordLogin{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordPasswordLoginForm implements forms.RecordPasswordLogin {
  constructor(app: CoreApp, collection: models.Collection)
}

interface RecordPasswordResetConfirmForm extends forms.RecordPasswordResetConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordPasswordResetConfirmForm implements forms.RecordPasswordResetConfirm {
  constructor(app: CoreApp, collection: models.Collection)
}

interface RecordPasswordResetRequestForm extends forms.RecordPasswordResetRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordPasswordResetRequestForm implements forms.RecordPasswordResetRequest {
  constructor(app: CoreApp, collection: models.Collection)
}

interface RecordUpsertForm extends forms.RecordUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordUpsertForm implements forms.RecordUpsert {
  constructor(app: CoreApp, record: models.Record)
}

interface RecordVerificationConfirmForm extends forms.RecordVerificationConfirm{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordVerificationConfirmForm implements forms.RecordVerificationConfirm {
  constructor(app: CoreApp, collection: models.Collection)
}

interface RecordVerificationRequestForm extends forms.RecordVerificationRequest{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordVerificationRequestForm implements forms.RecordVerificationRequest {
  constructor(app: CoreApp, collection: models.Collection)
}

interface SettingsUpsertForm extends forms.SettingsUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class SettingsUpsertForm implements forms.SettingsUpsert {
  constructor(app: CoreApp)
}

interface TestEmailSendForm extends forms.TestEmailSend{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class TestEmailSendForm implements forms.TestEmailSend {
  constructor(app: CoreApp)
}

interface TestS3FilesystemForm extends forms.TestS3Filesystem{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class TestS3FilesystemForm implements forms.TestS3Filesystem {
  constructor(app: CoreApp)
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
  /**
   * Route handler to serve static directory content (html, js, css, etc.).
   *
   * If a file resource is missing and indexFallback is set, the request
   * will be forwarded to the base index.html (useful for SPA).
   */
  export function staticDirectoryHandler(dir: string, indexFallback: boolean): echo.HandlerFunc

  let requireGuestOnly:          apis.requireGuestOnly
  let requireRecordAuth:         apis.requireRecordAuth
  let requireAdminAuth:          apis.requireAdminAuth
  let requireAdminAuthOnlyIfAny: apis.requireAdminAuthOnlyIfAny
  let requireAdminOrRecordAuth:  apis.requireAdminOrRecordAuth
  let requireAdminOrOwnerAuth:   apis.requireAdminOrOwnerAuth
  let activityLogger:            apis.activityLogger
  let requestInfo:               apis.requestInfo
  let recordAuthResponse:        apis.recordAuthResponse
  let gzip:                      middleware.gzip
  let bodyLimit:                 middleware.bodyLimit
  let enrichRecord:              apis.enrichRecord
  let enrichRecords:             apis.enrichRecords
}

// -------------------------------------------------------------------
// httpClientBinds
// -------------------------------------------------------------------

// extra FormData overload to prevent TS warnings when used with non File/Blob value.
interface FormData {
  append(key:string, value:any): void
  set(key:string, value:any): void
}

/**
 * ` + "`" + `$http` + "`" + ` defines common methods for working with HTTP requests.
 *
 * @group PocketBase
 */
declare namespace $http {
  /**
   * Sends a single HTTP request.
   *
   * Example:
   *
   * ` + "```" + `js
   * const res = $http.send({
   *     url:    "https://example.com",
   *     body:   JSON.stringify({"title": "test"})
   *     method: "post",
   * })
   *
   * console.log(res.statusCode) // the response HTTP status code
   * console.log(res.headers)    // the response headers (eg. res.headers['X-Custom'][0])
   * console.log(res.cookies)    // the response cookies (eg. res.cookies.sessionId.value)
   * console.log(res.raw)        // the response body as plain text
   * console.log(res.json)       // the response body as parsed json array or map
   * ` + "```" + `
   */
  function send(config: {
    url:      string,
    body?:    string|FormData,
    method?:  string, // default to "GET"
    headers?: { [key:string]: string },
    timeout?: number, // default to 120

    // deprecated, please use body instead
    data?: { [key:string]: any },
  }): {
    statusCode: number,
    headers:    { [key:string]: Array<string> },
    cookies:    { [key:string]: http.Cookie },
    raw:        string,
    json:       any,
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
			"github.com/labstack/echo/v5/middleware":            {"Gzip", "BodyLimit"},
			"github.com/go-ozzo/ozzo-validation/v4":             {"Error"},
			"github.com/pocketbase/dbx":                         {"*"},
			"github.com/pocketbase/pocketbase/tools/security":   {"*"},
			"github.com/pocketbase/pocketbase/tools/filesystem": {"*"},
			"github.com/pocketbase/pocketbase/tools/template":   {"*"},
			"github.com/pocketbase/pocketbase/tokens":           {"*"},
			"github.com/pocketbase/pocketbase/mails":            {"*"},
			"github.com/pocketbase/pocketbase/apis":             {"*"},
			"github.com/pocketbase/pocketbase/forms":            {"*"},
			"github.com/pocketbase/pocketbase":                  {"*"},
			"path/filepath":                                     {"*"},
			"os":                                                {"*"},
			"os/exec":                                           {"Command"},
		},
		FieldNameFormatter: func(s string) string {
			return mapper.FieldName(nil, reflect.StructField{Name: s})
		},
		MethodNameFormatter: func(s string) string {
			return mapper.MethodName(nil, reflect.Method{Name: s})
		},
		TypeMappings: map[string]string{
			"crypto.*":    "any",
			"acme.*":      "any",
			"autocert.*":  "any",
			"driver.*":    "any",
			"reflect.*":   "any",
			"fmt.*":       "any",
			"rand.*":      "any",
			"tls.*":       "any",
			"asn1.*":      "any",
			"pkix.*":      "any",
			"x509.*":      "any",
			"pflag.*":     "any",
			"flag.*":      "any",
			"log.*":       "any",
			"http.Client": "any",
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

	// replace the original app interfaces with their non-"on*"" hooks equivalents
	result = strings.ReplaceAll(result, "core.App", "CoreApp")
	result = strings.ReplaceAll(result, "pocketbase.PocketBase", "PocketBase")
	result = strings.ReplaceAll(result, "ORIGINAL_CORE_APP", "core.App")
	result = strings.ReplaceAll(result, "ORIGINAL_POCKETBASE", "pocketbase.PocketBase")

	// prepend a timestamp with the generation time
	// so that it can be compared without reading the entire file
	result = fmt.Sprintf("// %d\n%s", time.Now().Unix(), result)

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
