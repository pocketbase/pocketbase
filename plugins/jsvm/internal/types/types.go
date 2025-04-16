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
 * routerAdd("GET", "/hello", (e) => {
 *     return e.json(200, {"message": "Hello!"})
 * }, $apis.requireAuth())
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function routerAdd(
  method: string,
  path: string,
  handler: (e: core.RequestEvent) => void,
  ...middlewares: Array<string|((e: core.RequestEvent) => void)|Middleware>,
): void;

/**
 * RouterUse registers one or more global middlewares that are executed
 * along the handler middlewares after a matching route is found.
 *
 * Example:
 *
 * ` + "```" + `js
 * routerUse((e) => {
 *   console.log(e.request.url.path)
 *   return e.next()
 * })
 * ` + "```" + `
 *
 * _Note that this method is available only in pb_hooks context._
 *
 * @group PocketBase
 */
declare function routerUse(...middlewares: Array<string|((e: core.RequestEvent) => void)|Middleware>): void;

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
    [Property in keyof Type as Exclude<Property, ` + "`on${string}`" + `|'cron'>]: Type[Property]
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
 * This method is superseded by toString.
 *
 * @deprecated
 * @group PocketBase
 */
declare function readerToString(reader: any, maxBytes?: number): string;

/**
 * toString stringifies the specified value.
 *
 * Support optional second maxBytes argument to limit the max read bytes
 * when the value is a io.Reader (default to 32MB).
 *
 * Types that don't have explicit string representation are json serialized.
 *
 * Example:
 *
 * ` + "```" + `js
 * // io.Reader
 * const ex1 = toString(e.request.body)
 *
 * // slice of bytes ("hello")
 * const ex2 = toString([104 101 108 108 111])
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare function toString(val: any, maxBytes?: number): string;

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
 * $app.recordQuery("articles").limit(10).all(records)
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare function arrayOf<T>(model: T): Array<T>;

/**
 * DynamicModel creates a new dynamic model with fields from the provided data shape.
 *
 * Note that in order to use 0 as double/float initialization number you have to use negative zero (` + "`-0`" + `).
 *
 * Example:
 *
 * ` + "```" + `js
 * const model = new DynamicModel({
 *     name:       ""
 *     age:        0,  // int64
 *     totalSpent: -0, // float64
 *     active:     false,
 *     roles:      [],
 *     meta:       {}
 * })
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class DynamicModel {
  constructor(shape?: { [key:string]: any })
}

interface Context extends context.Context{} // merge
/**
 * Context creates a new empty Go context.Context.
 *
 * This is usually used as part of some Go transitive bindings.
 *
 * Example:
 *
 * ` + "```" + `js
 * const blank = new Context()
 *
 * // with single key-value pair
 * const base = new Context(null, "a", 123)
 * console.log(base.value("a")) // 123
 *
 * // extend with additional key-value pair
 * const sub = new Context(base, "b", 456)
 * console.log(sub.value("a")) // 123
 * console.log(sub.value("b")) // 456
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Context implements context.Context {
  constructor(parentCtx?: Context, key?: any, value?: any)
}

/**
 * Record model class.
 *
 * ` + "```" + `js
 * const collection = $app.findCollectionByNameOrId("article")
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
  new(collection?: core.Collection, data?: { [key:string]: any }): core.Record

  // note: declare as "newable" const due to conflict with the Record TS utility type
}

interface Collection extends core.Collection{
  type: "base" | "view" | "auth"
} // merge
/**
 * Collection model class.
 *
 * ` + "```" + `js
 * const collection = new Collection({
 *     type:       "base",
 *     name:       "article",
 *     listRule:   "@request.auth.id != '' || status = 'public'",
 *     viewRule:   "@request.auth.id != '' || status = 'public'",
 *     deleteRule: "@request.auth.id != ''",
 *     fields: [
 *         {
 *             name: "title",
 *             type: "text",
 *             required: true,
 *             min: 6,
 *             max: 100,
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
declare class Collection implements core.Collection {
  constructor(data?: Partial<Collection>)
}

interface FieldsList extends core.FieldsList{} // merge
/**
 * FieldsList model class, usually used to define the Collection.fields.
 *
 * @group PocketBase
 */
declare class FieldsList implements core.FieldsList {
  constructor(data?: Partial<core.FieldsList>)
}

interface Field extends core.Field{} // merge
/**
 * Field model class, usually used as part of the FieldsList model.
 *
 * @group PocketBase
 */
declare class Field implements core.Field {
  constructor(data?: Partial<core.Field>)
}

interface NumberField extends core.NumberField{} // merge
/**
 * {@inheritDoc core.NumberField}
 *
 * @group PocketBase
 */
declare class NumberField implements core.NumberField {
  constructor(data?: Partial<core.NumberField>)
}

interface BoolField extends core.BoolField{} // merge
/**
 * {@inheritDoc core.BoolField}
 *
 * @group PocketBase
 */
declare class BoolField implements core.BoolField {
  constructor(data?: Partial<core.BoolField>)
}

interface TextField extends core.TextField{} // merge
/**
 * {@inheritDoc core.TextField}
 *
 * @group PocketBase
 */
declare class TextField implements core.TextField {
  constructor(data?: Partial<core.TextField>)
}

interface URLField extends core.URLField{} // merge
/**
 * {@inheritDoc core.URLField}
 *
 * @group PocketBase
 */
declare class URLField implements core.URLField {
  constructor(data?: Partial<core.URLField>)
}

interface EmailField extends core.EmailField{} // merge
/**
 * {@inheritDoc core.EmailField}
 *
 * @group PocketBase
 */
declare class EmailField implements core.EmailField {
  constructor(data?: Partial<core.EmailField>)
}

interface EditorField extends core.EditorField{} // merge
/**
 * {@inheritDoc core.EditorField}
 *
 * @group PocketBase
 */
declare class EditorField implements core.EditorField {
  constructor(data?: Partial<core.EditorField>)
}

interface PasswordField extends core.PasswordField{} // merge
/**
 * {@inheritDoc core.PasswordField}
 *
 * @group PocketBase
 */
declare class PasswordField implements core.PasswordField {
  constructor(data?: Partial<core.PasswordField>)
}

interface DateField extends core.DateField{} // merge
/**
 * {@inheritDoc core.DateField}
 *
 * @group PocketBase
 */
declare class DateField implements core.DateField {
  constructor(data?: Partial<core.DateField>)
}

interface AutodateField extends core.AutodateField{} // merge
/**
 * {@inheritDoc core.AutodateField}
 *
 * @group PocketBase
 */
declare class AutodateField implements core.AutodateField {
  constructor(data?: Partial<core.AutodateField>)
}

interface JSONField extends core.JSONField{} // merge
/**
 * {@inheritDoc core.JSONField}
 *
 * @group PocketBase
 */
declare class JSONField implements core.JSONField {
  constructor(data?: Partial<core.JSONField>)
}

interface RelationField extends core.RelationField{} // merge
/**
 * {@inheritDoc core.RelationField}
 *
 * @group PocketBase
 */
declare class RelationField implements core.RelationField {
  constructor(data?: Partial<core.RelationField>)
}

interface SelectField extends core.SelectField{} // merge
/**
 * {@inheritDoc core.SelectField}
 *
 * @group PocketBase
 */
declare class SelectField implements core.SelectField {
  constructor(data?: Partial<core.SelectField>)
}

interface FileField extends core.FileField{} // merge
/**
 * {@inheritDoc core.FileField}
 *
 * @group PocketBase
 */
declare class FileField implements core.FileField {
  constructor(data?: Partial<core.FileField>)
}

interface GeoPointField extends core.GeoPointField{} // merge
/**
 * {@inheritDoc core.GeoPointField}
 *
 * @group PocketBase
 */
declare class GeoPointField implements core.GeoPointField {
  constructor(data?: Partial<core.GeoPointField>)
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

/**
 * RequestInfo defines a single core.RequestInfo instance, usually used
 * as part of various filter checks.
 *
 * Example:
 *
 * ` + "```" + `js
 * const authRecord = $app.findAuthRecordByEmail("users", "test@example.com")
 *
 * const info = new RequestInfo({
 *     auth:    authRecord,
 *     body:    {"name": 123},
 *     headers: {"x-token": "..."},
 * })
 *
 * const record = $app.findFirstRecordByData("articles", "slug", "hello")
 *
 * const canAccess = $app.canAccessRecord(record, info, "@request.auth.id != '' && @request.body.name = 123")
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare const RequestInfo: {
  new(info?: Partial<core.RequestInfo>): core.RequestInfo

 // note: declare as "newable" const due to conflict with the RequestInfo TS node type
}

/**
 * Middleware defines a single request middleware handler.
 *
 * This class is usually used when you want to explicitly specify a priority to your custom route middleware.
 *
 * Example:
 *
 * ` + "```" + `js
 * routerUse(new Middleware((e) => {
 *   console.log(e.request.url.path)
 *   return e.next()
 * }, -10))
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Middleware {
  constructor(
    func: string|((e: core.RequestEvent) => void),
    priority?: number,
    id?: string,
  )
}

interface Timezone extends time.Location{} // merge
/**
 * Timezone returns the timezone location with the given name.
 *
 * The name is expected to be a location name corresponding to a file
 * in the IANA Time Zone database, such as "America/New_York".
 *
 * If the name is "Local", LoadLocation returns Local.
 *
 * If the name is "", invalid or "UTC", returns UTC.
 *
 * The constructor is equivalent to calling the Go ` + "`" + `time.LoadLocation(name)` + "`" + ` method.
 *
 * Example:
 *
 * ` + "```" + `js
 * const zone = new Timezone("America/New_York")
 * $app.cron().setTimezone(zone)
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class Timezone implements time.Location {
  constructor(name?: string)
}

interface DateTime extends types.DateTime{} // merge
/**
 * DateTime defines a single DateTime type instance.
 * The returned date is always represented in UTC.
 *
 * Example:
 *
 * ` + "```" + `js
 * const dt0 = new DateTime() // now
 *
 * // full datetime string
 * const dt1 = new DateTime('2023-07-01 00:00:00.000Z')
 *
 * // datetime string with default "parse in" timezone location
 * //
 * // similar to new DateTime('2023-07-01 00:00:00 +01:00') or new DateTime('2023-07-01 00:00:00 +02:00')
 * // but accounts for the daylight saving time (DST)
 * const dt2 = new DateTime('2023-07-01 00:00:00', 'Europe/Amsterdam')
 * ` + "```" + `
 *
 * @group PocketBase
 */
declare class DateTime implements types.DateTime {
  constructor(date?: string, defaultParseInLocation?: string)
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
// mailsBinds
// -------------------------------------------------------------------

/**
 * ` + "`" + `$mails` + "`" + ` defines helpers to send common
 * auth records emails like verification, password reset, etc.
 *
 * @group PocketBase
 */
declare namespace $mails {
  let sendRecordPasswordReset: mails.sendRecordPasswordReset
  let sendRecordVerification:  mails.sendRecordVerification
  let sendRecordChangeEmail:   mails.sendRecordChangeEmail
  let sendRecordOTP:           mails.sendRecordOTP
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
  let randomStringByRegex:            security.randomStringByRegex
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

  /**
   * {@inheritDoc security.newJWT}
   */
  export function createJWT(payload: { [key:string]: any }, signingKey: string, secDuration: number): string

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
   * fileFromURL creates a new File from the provided url by
   * downloading the resource and creating a BytesReader.
   *
   * Example:
   *
   * ` + "```" + `js
   * // with default max timeout of 120sec
   * const file1 = $filesystem.fileFromURL("https://...")
   *
   * // with custom timeout of 15sec
   * const file2 = $filesystem.fileFromURL("https://...", 15)
   * ` + "```" + `
   */
  export function fileFromURL(url: string, secTimeout?: number): filesystem.File
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
   * const output = toString(cmd.output());
   * ` + "```" + `
   */
  export let cmd: exec.command

  /**
   * Args hold the command-line arguments, starting with the program name.
   */
  export let args: Array<string>

  export let exit:      os.exit
  export let getenv:    os.getenv
  export let dirFS:     os.dirFS
  export let readFile:  os.readFile
  export let writeFile: os.writeFile
  export let stat:      os.stat
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

interface AppleClientSecretCreateForm extends forms.AppleClientSecretCreate{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class AppleClientSecretCreateForm implements forms.AppleClientSecretCreate {
  constructor(app: CoreApp)
}

interface RecordUpsertForm extends forms.RecordUpsert{} // merge
/**
 * @inheritDoc
 * @group PocketBase
 */
declare class RecordUpsertForm implements forms.RecordUpsert {
  constructor(app: CoreApp, record: core.Record)
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

interface ApiError extends router.ApiError{} // merge
/**
 * @inheritDoc
 *
 * @group PocketBase
 */
declare class ApiError implements router.ApiError {
  constructor(status?: number, message?: string, data?: any)
}

interface NotFoundError extends router.ApiError{} // merge
/**
 * NotFounderor returns 404 ApiError.
 *
 * @group PocketBase
 */
declare class NotFoundError implements router.ApiError {
  constructor(message?: string, data?: any)
}

interface BadRequestError extends router.ApiError{} // merge
/**
 * BadRequestError returns 400 ApiError.
 *
 * @group PocketBase
 */
declare class BadRequestError implements router.ApiError {
  constructor(message?: string, data?: any)
}

interface ForbiddenError extends router.ApiError{} // merge
/**
 * ForbiddenError returns 403 ApiError.
 *
 * @group PocketBase
 */
declare class ForbiddenError implements router.ApiError {
  constructor(message?: string, data?: any)
}

interface UnauthorizedError extends router.ApiError{} // merge
/**
 * UnauthorizedError returns 401 ApiError.
 *
 * @group PocketBase
 */
declare class UnauthorizedError implements router.ApiError {
  constructor(message?: string, data?: any)
}

interface TooManyRequestsError extends router.ApiError{} // merge
/**
 * TooManyRequestsError returns 429 ApiError.
 *
 * @group PocketBase
 */
declare class TooManyRequestsError implements router.ApiError {
  constructor(message?: string, data?: any)
}

interface InternalServerError extends router.ApiError{} // merge
/**
 * InternalServerError returns 429 ApiError.
 *
 * @group PocketBase
 */
declare class InternalServerError implements router.ApiError {
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
  export function static(dir: string, indexFallback: boolean): (e: core.RequestEvent) => void

  let requireGuestOnly:              apis.requireGuestOnly
  let requireAuth:                   apis.requireAuth
  let requireSuperuserAuth:          apis.requireSuperuserAuth
  let requireSuperuserOrOwnerAuth:   apis.requireSuperuserOrOwnerAuth
  let skipSuccessActivityLog:        apis.skipSuccessActivityLog
  let gzip:                          apis.gzip
  let bodyLimit:                     apis.bodyLimit
  let enrichRecord:                  apis.enrichRecord
  let enrichRecords:                 apis.enrichRecords

  /**
   * RecordAuthResponse writes standardized json record auth response
   * into the specified request event.
   *
   * The authMethod argument specify the name of the current authentication method (eg. password, oauth2, etc.)
   * that it is used primarily as an auth identifier during MFA and for login alerts.
   *
   * Set authMethod to empty string if you want to ignore the MFA checks and the login alerts
   * (can be also adjusted additionally via the onRecordAuthRequest hook).
   */
  export function recordAuthResponse(e: core.RequestEvent, authRecord: core.Record, authMethod: string, meta?: any): void
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
   *     method:  "POST",
   *     url:     "https://example.com",
   *     body:    JSON.stringify({"title": "test"}),
   *     headers: { 'Content-Type': 'application/json' }
   * })
   *
   * console.log(res.statusCode) // the response HTTP status code
   * console.log(res.headers)    // the response headers (eg. res.headers['X-Custom'][0])
   * console.log(res.cookies)    // the response cookies (eg. res.cookies.sessionId.value)
   * console.log(res.body)       // the response body as raw bytes slice
   * console.log(res.json)       // the response body as parsed json array or map
   * ` + "```" + `
   */
  function send(config: {
    url:      string,
    body?:    string|FormData,
    method?:  string, // default to "GET"
    headers?: { [key:string]: string },
    timeout?: number, // default to 120

    // @deprecated please use body instead
    data?: { [key:string]: any },
  }): {
    statusCode: number,
    headers:    { [key:string]: Array<string> },
    cookies:    { [key:string]: http.Cookie },
    json:       any,
    body:       Array<number>,

    // @deprecated please use toString(result.body) instead
    raw: string,
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
  up: (txApp: CoreApp) => void,
  down?: (txApp: CoreApp) => void
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
			"github.com/pocketbase/pocketbase/tools/template":   {"*"},
			"github.com/pocketbase/pocketbase/mails":            {"*"},
			"github.com/pocketbase/pocketbase/apis":             {"*"},
			"github.com/pocketbase/pocketbase/core":             {"*"},
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
			"crypto.*":     "any",
			"acme.*":       "any",
			"autocert.*":   "any",
			"driver.*":     "any",
			"reflect.*":    "any",
			"fmt.*":        "any",
			"rand.*":       "any",
			"tls.*":        "any",
			"asn1.*":       "any",
			"pkix.*":       "any",
			"x509.*":       "any",
			"pflag.*":      "any",
			"flag.*":       "any",
			"log.*":        "any",
			"http.Client":  "any",
			"mail.Address": "{ address: string; name?: string; }", // prevents the LSP to complain in case no name is provided
		},
		Indent:               " ", // use only a single space to reduce slightly the size
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

	excluded := []string{"OnServe"}
	appType := reflect.TypeOf(struct{ core.App }{})
	totalMethods := appType.NumMethod()

	for i := 0; i < totalMethods; i++ {
		method := appType.Method(i)
		if !strings.HasPrefix(method.Name, "On") || list.ExistInSlice(method.Name, excluded) {
			continue // not a hook or excluded
		}

		hookType := method.Type.Out(0)

		withTags := strings.HasPrefix(hookType.String(), "*hook.TaggedHook")

		addMethod, ok := hookType.MethodByName("BindFunc")
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
