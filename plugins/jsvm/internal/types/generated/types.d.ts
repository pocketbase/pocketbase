// 1762248160
// GENERATED CODE - DO NOT MODIFY BY HAND

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
 * ```js
 * // prints "Hello world!" on every 30 minutes
 * cronAdd("hello", "*\/30 * * * *", () => {
 *     console.log("Hello world!")
 * })
 * ```
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
 * ```js
 * cronRemove("hello")
 * ```
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
 * ```js
 * routerAdd("GET", "/hello", (e) => {
 *     return e.json(200, {"message": "Hello!"})
 * }, $apis.requireAuth())
 * ```
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
 * ```js
 * routerUse((e) => {
 *   console.log(e.request.url.path)
 *   return e.next()
 * })
 * ```
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
    [Property in keyof Type as Exclude<Property, `on${string}`|'cron'>]: Type[Property]
};

// CoreApp without the on* hook methods
type CoreApp = excludeHooks<core.App>

// PocketBase without the on* hook methods
type PocketBase = excludeHooks<pocketbase.PocketBase>

/**
 * `$app` is the current running PocketBase instance that is globally
 * available in each .pb.js file.
 *
 * _Note that this variable is available only in pb_hooks context._
 *
 * @namespace
 * @group PocketBase
 */
declare var $app: PocketBase

/**
 * `$template` is a global helper to load and cache HTML templates on the fly.
 *
 * The templates uses the standard Go [html/template](https://pkg.go.dev/html/template)
 * and [text/template](https://pkg.go.dev/text/template) package syntax.
 *
 * Example:
 *
 * ```js
 * const html = $template.loadFiles(
 *     "views/layout.html",
 *     "views/content.html",
 * ).render({"name": "John"})
 * ```
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
 * ```js
 * // io.Reader
 * const ex1 = toString(e.request.body)
 *
 * // slice of bytes
 * const ex2 = toString([104 101 108 108 111]) // "hello"
 *
 * // null
 * const ex3 = toString(null) // ""
 * ```
 *
 * @group PocketBase
 */
declare function toString(val: any, maxBytes?: number): string;

/**
 * toBytes converts the specified value into a bytes slice.
 *
 * Support optional second maxBytes argument to limit the max read bytes
 * when the value is a io.Reader (default to 32MB).
 *
 * Types that don't have Go slice representation (bool, objects, etc.)
 * are serialized to UTF8 string and its bytes slice is returned.
 *
 * Example:
 *
 * ```js
 * // io.Reader
 * const ex1 = toBytes(e.request.body)
 *
 * // string
 * const ex2 = toBytes("hello") // [104 101 108 108 111]
 *
 * // object (the same as the string '{"test":1}')
 * const ex3 = toBytes({"test":1}) // [123 34 116 101 115 116 34 58 49 125]
 *
 * // null
 * const ex4 = toBytes(null) // []
 * ```
 *
 * @group PocketBase
 */
declare function toBytes(val: any, maxBytes?: number): Array<number>;

/**
 * sleep pauses the current goroutine for at least the specified user duration (in ms).
 * A zero or negative duration returns immediately.
 *
 * Example:
 *
 * ```js
 * sleep(250) // sleeps for 250ms
 * ```
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
 * ```js
 * const records = arrayOf(new Record)
 *
 * $app.recordQuery("articles").limit(10).all(records)
 * ```
 *
 * @group PocketBase
 */
declare function arrayOf<T>(model: T): Array<T>;

/**
 * DynamicModel creates a new dynamic model with fields from the provided data shape.
 *
 * Caveats:
 * - In order to use 0 as double/float initialization number you have to negate it (`-0`).
 * - You need to use lowerCamelCase when accessing the model fields (e.g. `model.roles` and not `model.Roles`).
 *
 * Example:
 *
 * ```js
 * const model = new DynamicModel({
 *     name:       ""
 *     age:        0,  // int64
 *     totalSpent: -0, // float64
 *     active:     false,
 *     Roles:      [], // maps to "Roles" in the DB/JSON but the prop would be accessible via "model.roles"
 *     meta:       {}
 * })
 * ```
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
 * ```js
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
 * ```
 *
 * @group PocketBase
 */
declare class Context implements context.Context {
  constructor(parentCtx?: Context, key?: any, value?: any)
}

/**
 * Record model class.
 *
 * ```js
 * const collection = $app.findCollectionByNameOrId("article")
 *
 * const record = new Record(collection, {
 *     title: "Lorem ipsum"
 * })
 *
 * // or set field values after the initialization
 * record.set("description", "...")
 * ```
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
 * ```js
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
 * ```
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
 * ```js
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
 * ```
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
 * ```js
 * const command = new Command({
 *     use: "hello",
 *     run: (cmd, args) => { console.log("Hello world!") },
 * })
 *
 * $app.rootCmd.addCommand(command);
 * ```
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
 * ```js
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
 * ```
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
 * ```js
 * routerUse(new Middleware((e) => {
 *   console.log(e.request.url.path)
 *   return e.next()
 * }, -10))
 * ```
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
 * The constructor is equivalent to calling the Go `time.LoadLocation(name)` method.
 *
 * Example:
 *
 * ```js
 * const zone = new Timezone("America/New_York")
 * $app.cron().setTimezone(zone)
 * ```
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
 * ```js
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
 * ```
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
 * ```js
 * new ValidationError("invalid_title", "Title is not valid")
 * ```
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
 * ```js
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
 * ```
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
 * ```js
 * onRealtimeConnectRequest((e) => {
 *     e.client.send(new SubscriptionMessage({
 *         name: "example",
 *         data: '{"greeting": "Hello world"}'
 *     }))
 * })
 * ```
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
 * `$dbx` defines common utility for working with the DB abstraction.
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
 * `$mails` defines helpers to send common
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
 * `$security` defines low level helpers for creating
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
 * `$filesystem` defines common helpers for working
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
   * ```js
   * // with default max timeout of 120sec
   * const file1 = $filesystem.fileFromURL("https://...")
   *
   * // with custom timeout of 15sec
   * const file2 = $filesystem.fileFromURL("https://...", 15)
   * ```
   */
  export function fileFromURL(url: string, secTimeout?: number): filesystem.File
}

// -------------------------------------------------------------------
// filepathBinds
// -------------------------------------------------------------------

/**
 * `$filepath` defines common helpers for manipulating filename
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
 * `$os` defines common helpers for working with the OS level primitives
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
   * ```js
   * // prepare the command to execute
   * const cmd = $os.cmd('ls', '-sl')
   *
   * // execute the command and return its standard output as string
   * const output = toString(cmd.output());
   * ```
   */
  export let cmd: exec.command

  /**
   * Args hold the command-line arguments, starting with the program name.
   */
  export let args: Array<string>

  export let exit:       os.exit
  export let getenv:     os.getenv
  export let dirFS:      os.dirFS
  export let readFile:   os.readFile
  export let writeFile:  os.writeFile
  export let stat:       os.stat
  export let readDir:    os.readDir
  export let tempDir:    os.tempDir
  export let truncate:   os.truncate
  export let getwd:      os.getwd
  export let mkdir:      os.mkdir
  export let mkdirAll:   os.mkdirAll
  export let rename:     os.rename
  export let remove:     os.remove
  export let removeAll:  os.removeAll
  export let openRoot:   os.openRoot
  export let openInRoot: os.openInRoot
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
 * `$apis` defines commonly used PocketBase api helpers and middlewares.
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
 * `$http` defines common methods for working with HTTP requests.
 *
 * @group PocketBase
 */
declare namespace $http {
  /**
   * Sends a single HTTP request.
   *
   * Example:
   *
   * ```js
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
   * ```
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
/** @group PocketBase */declare function onBackupCreate(handler: (e: core.BackupEvent) => void): void
/** @group PocketBase */declare function onBackupRestore(handler: (e: core.BackupEvent) => void): void
/** @group PocketBase */declare function onBatchRequest(handler: (e: core.BatchRequestEvent) => void): void
/** @group PocketBase */declare function onBootstrap(handler: (e: core.BootstrapEvent) => void): void
/** @group PocketBase */declare function onCollectionAfterCreateError(handler: (e: core.CollectionErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionAfterCreateSuccess(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionAfterDeleteError(handler: (e: core.CollectionErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionAfterDeleteSuccess(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionAfterUpdateError(handler: (e: core.CollectionErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionAfterUpdateSuccess(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionCreate(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionCreateExecute(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionCreateRequest(handler: (e: core.CollectionRequestEvent) => void): void
/** @group PocketBase */declare function onCollectionDelete(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionDeleteExecute(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionDeleteRequest(handler: (e: core.CollectionRequestEvent) => void): void
/** @group PocketBase */declare function onCollectionUpdate(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionUpdateExecute(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionUpdateRequest(handler: (e: core.CollectionRequestEvent) => void): void
/** @group PocketBase */declare function onCollectionValidate(handler: (e: core.CollectionEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onCollectionViewRequest(handler: (e: core.CollectionRequestEvent) => void): void
/** @group PocketBase */declare function onCollectionsImportRequest(handler: (e: core.CollectionsImportRequestEvent) => void): void
/** @group PocketBase */declare function onCollectionsListRequest(handler: (e: core.CollectionsListRequestEvent) => void): void
/** @group PocketBase */declare function onFileDownloadRequest(handler: (e: core.FileDownloadRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onFileTokenRequest(handler: (e: core.FileTokenRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerRecordAuthAlertSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerRecordEmailChangeSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerRecordOTPSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerRecordPasswordResetSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerRecordVerificationSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerSend(handler: (e: core.MailerEvent) => void): void
/** @group PocketBase */declare function onModelAfterCreateError(handler: (e: core.ModelErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterCreateSuccess(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterDeleteError(handler: (e: core.ModelErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterDeleteSuccess(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterUpdateError(handler: (e: core.ModelErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterUpdateSuccess(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelCreate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelCreateExecute(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelDelete(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelDeleteExecute(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelUpdate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelUpdateExecute(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelValidate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRealtimeConnectRequest(handler: (e: core.RealtimeConnectRequestEvent) => void): void
/** @group PocketBase */declare function onRealtimeMessageSend(handler: (e: core.RealtimeMessageEvent) => void): void
/** @group PocketBase */declare function onRealtimeSubscribeRequest(handler: (e: core.RealtimeSubscribeRequestEvent) => void): void
/** @group PocketBase */declare function onRecordAfterCreateError(handler: (e: core.RecordErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterCreateSuccess(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterDeleteError(handler: (e: core.RecordErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterDeleteSuccess(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterUpdateError(handler: (e: core.RecordErrorEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterUpdateSuccess(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAuthRefreshRequest(handler: (e: core.RecordAuthRefreshRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAuthRequest(handler: (e: core.RecordAuthRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAuthWithOAuth2Request(handler: (e: core.RecordAuthWithOAuth2RequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAuthWithOTPRequest(handler: (e: core.RecordAuthWithOTPRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAuthWithPasswordRequest(handler: (e: core.RecordAuthWithPasswordRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordConfirmEmailChangeRequest(handler: (e: core.RecordConfirmEmailChangeRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordConfirmPasswordResetRequest(handler: (e: core.RecordConfirmPasswordResetRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordConfirmVerificationRequest(handler: (e: core.RecordConfirmVerificationRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordCreate(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordCreateExecute(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordCreateRequest(handler: (e: core.RecordRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordDelete(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordDeleteExecute(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordDeleteRequest(handler: (e: core.RecordRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordEnrich(handler: (e: core.RecordEnrichEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordRequestEmailChangeRequest(handler: (e: core.RecordRequestEmailChangeRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordRequestOTPRequest(handler: (e: core.RecordCreateOTPRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordRequestPasswordResetRequest(handler: (e: core.RecordRequestPasswordResetRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordRequestVerificationRequest(handler: (e: core.RecordRequestVerificationRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordUpdate(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordUpdateExecute(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordUpdateRequest(handler: (e: core.RecordRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordValidate(handler: (e: core.RecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordViewRequest(handler: (e: core.RecordRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordsListRequest(handler: (e: core.RecordsListRequestEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onSettingsListRequest(handler: (e: core.SettingsListRequestEvent) => void): void
/** @group PocketBase */declare function onSettingsReload(handler: (e: core.SettingsReloadEvent) => void): void
/** @group PocketBase */declare function onSettingsUpdateRequest(handler: (e: core.SettingsUpdateRequestEvent) => void): void
/** @group PocketBase */declare function onTerminate(handler: (e: core.TerminateEvent) => void): void
type _TygojaDict = { [key:string | number | symbol]: any; }
type _TygojaAny = any

/**
 * Package os provides a platform-independent interface to operating system
 * functionality. The design is Unix-like, although the error handling is
 * Go-like; failing calls return values of type error rather than error numbers.
 * Often, more information is available within the error. For example,
 * if a call that takes a file name fails, such as [Open] or [Stat], the error
 * will include the failing file name when printed and will be of type
 * [*PathError], which may be unpacked for more information.
 * 
 * The os interface is intended to be uniform across all operating systems.
 * Features not generally available appear in the system-specific package syscall.
 * 
 * Here is a simple example, opening a file and reading some of it.
 * 
 * ```
 * 	file, err := os.Open("file.go") // For read access.
 * 	if err != nil {
 * 		log.Fatal(err)
 * 	}
 * ```
 * 
 * If the open fails, the error string will be self-explanatory, like
 * 
 * ```
 * 	open file.go: no such file or directory
 * ```
 * 
 * The file's data can then be read into a slice of bytes. Read and
 * Write take their byte counts from the length of the argument slice.
 * 
 * ```
 * 	data := make([]byte, 100)
 * 	count, err := file.Read(data)
 * 	if err != nil {
 * 		log.Fatal(err)
 * 	}
 * 	fmt.Printf("read %d bytes: %q\n", count, data[:count])
 * ```
 * 
 * # Concurrency
 * 
 * The methods of [File] correspond to file system operations. All are
 * safe for concurrent use. The maximum number of concurrent
 * operations on a File may be limited by the OS or the system. The
 * number should be high, but exceeding it may degrade performance or
 * cause other issues.
 */
namespace os {
 interface readdirMode extends Number{}
 interface File {
  /**
   * Readdir reads the contents of the directory associated with file and
   * returns a slice of up to n [FileInfo] values, as would be returned
   * by [Lstat], in directory order. Subsequent calls on the same file will yield
   * further FileInfos.
   * 
   * If n > 0, Readdir returns at most n FileInfo structures. In this case, if
   * Readdir returns an empty slice, it will return a non-nil error
   * explaining why. At the end of a directory, the error is [io.EOF].
   * 
   * If n <= 0, Readdir returns all the FileInfo from the directory in
   * a single slice. In this case, if Readdir succeeds (reads all
   * the way to the end of the directory), it returns the slice and a
   * nil error. If it encounters an error before the end of the
   * directory, Readdir returns the FileInfo read until that point
   * and a non-nil error.
   * 
   * Most clients are better served by the more efficient ReadDir method.
   */
  readdir(n: number): Array<FileInfo>
 }
 interface File {
  /**
   * Readdirnames reads the contents of the directory associated with file
   * and returns a slice of up to n names of files in the directory,
   * in directory order. Subsequent calls on the same file will yield
   * further names.
   * 
   * If n > 0, Readdirnames returns at most n names. In this case, if
   * Readdirnames returns an empty slice, it will return a non-nil error
   * explaining why. At the end of a directory, the error is [io.EOF].
   * 
   * If n <= 0, Readdirnames returns all the names from the directory in
   * a single slice. In this case, if Readdirnames succeeds (reads all
   * the way to the end of the directory), it returns the slice and a
   * nil error. If it encounters an error before the end of the
   * directory, Readdirnames returns the names read until that point and
   * a non-nil error.
   */
  readdirnames(n: number): Array<string>
 }
 /**
  * A DirEntry is an entry read from a directory
  * (using the [ReadDir] function or a [File.ReadDir] method).
  */
 interface DirEntry extends fs.DirEntry{}
 interface File {
  /**
   * ReadDir reads the contents of the directory associated with the file f
   * and returns a slice of [DirEntry] values in directory order.
   * Subsequent calls on the same file will yield later DirEntry records in the directory.
   * 
   * If n > 0, ReadDir returns at most n DirEntry records.
   * In this case, if ReadDir returns an empty slice, it will return an error explaining why.
   * At the end of a directory, the error is [io.EOF].
   * 
   * If n <= 0, ReadDir returns all the DirEntry records remaining in the directory.
   * When it succeeds, it returns a nil error (not io.EOF).
   */
  readDir(n: number): Array<DirEntry>
 }
 interface readDir {
  /**
   * ReadDir reads the named directory,
   * returning all its directory entries sorted by filename.
   * If an error occurs reading the directory,
   * ReadDir returns the entries it was able to read before the error,
   * along with the error.
   */
  (name: string): Array<DirEntry>
 }
 interface copyFS {
  /**
   * CopyFS copies the file system fsys into the directory dir,
   * creating dir if necessary.
   * 
   * Files are created with mode 0o666 plus any execute permissions
   * from the source, and directories are created with mode 0o777
   * (before umask).
   * 
   * CopyFS will not overwrite existing files. If a file name in fsys
   * already exists in the destination, CopyFS will return an error
   * such that errors.Is(err, fs.ErrExist) will be true.
   * 
   * Symbolic links in fsys are not supported. A *PathError with Err set
   * to ErrInvalid is returned when copying from a symbolic link.
   * 
   * Symbolic links in dir are followed.
   * 
   * New files added to fsys (including if dir is a subdirectory of fsys)
   * while CopyFS is running are not guaranteed to be copied.
   * 
   * Copying stops at and returns the first error encountered.
   */
  (dir: string, fsys: fs.FS): void
 }
 /**
  * Auxiliary information if the File describes a directory
  */
 interface dirInfo {
 }
 interface expand {
  /**
   * Expand replaces ${var} or $var in the string based on the mapping function.
   * For example, [os.ExpandEnv](s) is equivalent to [os.Expand](s, [os.Getenv]).
   */
  (s: string, mapping: (_arg0: string) => string): string
 }
 interface expandEnv {
  /**
   * ExpandEnv replaces ${var} or $var in the string according to the values
   * of the current environment variables. References to undefined
   * variables are replaced by the empty string.
   */
  (s: string): string
 }
 interface getenv {
  /**
   * Getenv retrieves the value of the environment variable named by the key.
   * It returns the value, which will be empty if the variable is not present.
   * To distinguish between an empty value and an unset value, use [LookupEnv].
   */
  (key: string): string
 }
 interface lookupEnv {
  /**
   * LookupEnv retrieves the value of the environment variable named
   * by the key. If the variable is present in the environment the
   * value (which may be empty) is returned and the boolean is true.
   * Otherwise the returned value will be empty and the boolean will
   * be false.
   */
  (key: string): [string, boolean]
 }
 interface setenv {
  /**
   * Setenv sets the value of the environment variable named by the key.
   * It returns an error, if any.
   */
  (key: string, value: string): void
 }
 interface unsetenv {
  /**
   * Unsetenv unsets a single environment variable.
   */
  (key: string): void
 }
 interface clearenv {
  /**
   * Clearenv deletes all environment variables.
   */
  (): void
 }
 interface environ {
  /**
   * Environ returns a copy of strings representing the environment,
   * in the form "key=value".
   */
  (): Array<string>
 }
 interface timeout {
  [key:string]: any;
  timeout(): boolean
 }
 /**
  * PathError records an error and the operation and file path that caused it.
  */
 interface PathError extends fs.PathError{}
 /**
  * SyscallError records an error from a specific system call.
  */
 interface SyscallError {
  syscall: string
  err: Error
 }
 interface SyscallError {
  error(): string
 }
 interface SyscallError {
  unwrap(): void
 }
 interface SyscallError {
  /**
   * Timeout reports whether this error represents a timeout.
   */
  timeout(): boolean
 }
 interface newSyscallError {
  /**
   * NewSyscallError returns, as an error, a new [SyscallError]
   * with the given system call name and error details.
   * As a convenience, if err is nil, NewSyscallError returns nil.
   */
  (syscall: string, err: Error): void
 }
 interface isExist {
  /**
   * IsExist returns a boolean indicating whether its argument is known to report
   * that a file or directory already exists. It is satisfied by [ErrExist] as
   * well as some syscall errors.
   * 
   * This function predates [errors.Is]. It only supports errors returned by
   * the os package. New code should use errors.Is(err, fs.ErrExist).
   */
  (err: Error): boolean
 }
 interface isNotExist {
  /**
   * IsNotExist returns a boolean indicating whether its argument is known to
   * report that a file or directory does not exist. It is satisfied by
   * [ErrNotExist] as well as some syscall errors.
   * 
   * This function predates [errors.Is]. It only supports errors returned by
   * the os package. New code should use errors.Is(err, fs.ErrNotExist).
   */
  (err: Error): boolean
 }
 interface isPermission {
  /**
   * IsPermission returns a boolean indicating whether its argument is known to
   * report that permission is denied. It is satisfied by [ErrPermission] as well
   * as some syscall errors.
   * 
   * This function predates [errors.Is]. It only supports errors returned by
   * the os package. New code should use errors.Is(err, fs.ErrPermission).
   */
  (err: Error): boolean
 }
 interface isTimeout {
  /**
   * IsTimeout returns a boolean indicating whether its argument is known
   * to report that a timeout occurred.
   * 
   * This function predates [errors.Is], and the notion of whether an
   * error indicates a timeout can be ambiguous. For example, the Unix
   * error EWOULDBLOCK sometimes indicates a timeout and sometimes does not.
   * New code should use errors.Is with a value appropriate to the call
   * returning the error, such as [os.ErrDeadlineExceeded].
   */
  (err: Error): boolean
 }
 interface syscallErrorType extends syscall.Errno{}
 interface processMode extends Number{}
 interface processStatus extends Number{}
 /**
  * Process stores the information about a process created by [StartProcess].
  */
 interface Process {
  pid: number
 }
 /**
  * ProcAttr holds the attributes that will be applied to a new process
  * started by StartProcess.
  */
 interface ProcAttr {
  /**
   * If Dir is non-empty, the child changes into the directory before
   * creating the process.
   */
  dir: string
  /**
   * If Env is non-nil, it gives the environment variables for the
   * new process in the form returned by Environ.
   * If it is nil, the result of Environ will be used.
   */
  env: Array<string>
  /**
   * Files specifies the open files inherited by the new process. The
   * first three entries correspond to standard input, standard output, and
   * standard error. An implementation may support additional entries,
   * depending on the underlying operating system. A nil entry corresponds
   * to that file being closed when the process starts.
   * On Unix systems, StartProcess will change these File values
   * to blocking mode, which means that SetDeadline will stop working
   * and calling Close will not interrupt a Read or Write.
   */
  files: Array<(File | undefined)>
  /**
   * Operating system-specific process creation attributes.
   * Note that setting this field means that your program
   * may not execute properly or even compile on some
   * operating systems.
   */
  sys?: syscall.SysProcAttr
 }
 /**
  * A Signal represents an operating system signal.
  * The usual underlying implementation is operating system-dependent:
  * on Unix it is syscall.Signal.
  */
 interface Signal {
  [key:string]: any;
  string(): string
  signal(): void // to distinguish from other Stringers
 }
 interface getpid {
  /**
   * Getpid returns the process id of the caller.
   */
  (): number
 }
 interface getppid {
  /**
   * Getppid returns the process id of the caller's parent.
   */
  (): number
 }
 interface findProcess {
  /**
   * FindProcess looks for a running process by its pid.
   * 
   * The [Process] it returns can be used to obtain information
   * about the underlying operating system process.
   * 
   * On Unix systems, FindProcess always succeeds and returns a Process
   * for the given pid, regardless of whether the process exists. To test whether
   * the process actually exists, see whether p.Signal(syscall.Signal(0)) reports
   * an error.
   */
  (pid: number): (Process)
 }
 interface startProcess {
  /**
   * StartProcess starts a new process with the program, arguments and attributes
   * specified by name, argv and attr. The argv slice will become [os.Args] in the
   * new process, so it normally starts with the program name.
   * 
   * If the calling goroutine has locked the operating system thread
   * with [runtime.LockOSThread] and modified any inheritable OS-level
   * thread state (for example, Linux or Plan 9 name spaces), the new
   * process will inherit the caller's thread state.
   * 
   * StartProcess is a low-level interface. The [os/exec] package provides
   * higher-level interfaces.
   * 
   * If there is an error, it will be of type [*PathError].
   */
  (name: string, argv: Array<string>, attr: ProcAttr): (Process)
 }
 interface Process {
  /**
   * Release releases any resources associated with the [Process] p,
   * rendering it unusable in the future.
   * Release only needs to be called if [Process.Wait] is not.
   */
  release(): void
 }
 interface Process {
  /**
   * Kill causes the [Process] to exit immediately. Kill does not wait until
   * the Process has actually exited. This only kills the Process itself,
   * not any other processes it may have started.
   */
  kill(): void
 }
 interface Process {
  /**
   * Wait waits for the [Process] to exit, and then returns a
   * ProcessState describing its status and an error, if any.
   * Wait releases any resources associated with the Process.
   * On most operating systems, the Process must be a child
   * of the current process or an error will be returned.
   */
  wait(): (ProcessState)
 }
 interface Process {
  /**
   * Signal sends a signal to the [Process].
   * Sending [Interrupt] on Windows is not implemented.
   */
  signal(sig: Signal): void
 }
 interface ProcessState {
  /**
   * UserTime returns the user CPU time of the exited process and its children.
   */
  userTime(): time.Duration
 }
 interface ProcessState {
  /**
   * SystemTime returns the system CPU time of the exited process and its children.
   */
  systemTime(): time.Duration
 }
 interface ProcessState {
  /**
   * Exited reports whether the program has exited.
   * On Unix systems this reports true if the program exited due to calling exit,
   * but false if the program terminated due to a signal.
   */
  exited(): boolean
 }
 interface ProcessState {
  /**
   * Success reports whether the program exited successfully,
   * such as with exit status 0 on Unix.
   */
  success(): boolean
 }
 interface ProcessState {
  /**
   * Sys returns system-dependent exit information about
   * the process. Convert it to the appropriate underlying
   * type, such as [syscall.WaitStatus] on Unix, to access its contents.
   */
  sys(): any
 }
 interface ProcessState {
  /**
   * SysUsage returns system-dependent resource usage information about
   * the exited process. Convert it to the appropriate underlying
   * type, such as [*syscall.Rusage] on Unix, to access its contents.
   * (On Unix, *syscall.Rusage matches struct rusage as defined in the
   * getrusage(2) manual page.)
   */
  sysUsage(): any
 }
 /**
  * ProcessState stores information about a process, as reported by Wait.
  */
 interface ProcessState {
 }
 interface ProcessState {
  /**
   * Pid returns the process id of the exited process.
   */
  pid(): number
 }
 interface ProcessState {
  string(): string
 }
 interface ProcessState {
  /**
   * ExitCode returns the exit code of the exited process, or -1
   * if the process hasn't exited or was terminated by a signal.
   */
  exitCode(): number
 }
 interface executable {
  /**
   * Executable returns the path name for the executable that started
   * the current process. There is no guarantee that the path is still
   * pointing to the correct executable. If a symlink was used to start
   * the process, depending on the operating system, the result might
   * be the symlink or the path it pointed to. If a stable result is
   * needed, [path/filepath.EvalSymlinks] might help.
   * 
   * Executable returns an absolute path unless an error occurred.
   * 
   * The main use case is finding resources located relative to an
   * executable.
   */
  (): string
 }
 interface File {
  /**
   * Name returns the name of the file as presented to Open.
   * 
   * It is safe to call Name after [Close].
   */
  name(): string
 }
 /**
  * LinkError records an error during a link or symlink or rename
  * system call and the paths that caused it.
  */
 interface LinkError {
  op: string
  old: string
  new: string
  err: Error
 }
 interface LinkError {
  error(): string
 }
 interface LinkError {
  unwrap(): void
 }
 interface File {
  /**
   * Read reads up to len(b) bytes from the File and stores them in b.
   * It returns the number of bytes read and any error encountered.
   * At end of file, Read returns 0, io.EOF.
   */
  read(b: string|Array<number>): number
 }
 interface File {
  /**
   * ReadAt reads len(b) bytes from the File starting at byte offset off.
   * It returns the number of bytes read and the error, if any.
   * ReadAt always returns a non-nil error when n < len(b).
   * At end of file, that error is io.EOF.
   */
  readAt(b: string|Array<number>, off: number): number
 }
 interface File {
  /**
   * ReadFrom implements io.ReaderFrom.
   */
  readFrom(r: io.Reader): number
 }
 /**
  * noReadFrom can be embedded alongside another type to
  * hide the ReadFrom method of that other type.
  */
 interface noReadFrom {
 }
 interface noReadFrom {
  /**
   * ReadFrom hides another ReadFrom method.
   * It should never be called.
   */
  readFrom(_arg0: io.Reader): number
 }
 /**
  * fileWithoutReadFrom implements all the methods of *File other
  * than ReadFrom. This is used to permit ReadFrom to call io.Copy
  * without leading to a recursive call to ReadFrom.
  */
 type _sMDmKqA = noReadFrom&File
 interface fileWithoutReadFrom extends _sMDmKqA {
 }
 interface File {
  /**
   * Write writes len(b) bytes from b to the File.
   * It returns the number of bytes written and an error, if any.
   * Write returns a non-nil error when n != len(b).
   */
  write(b: string|Array<number>): number
 }
 interface File {
  /**
   * WriteAt writes len(b) bytes to the File starting at byte offset off.
   * It returns the number of bytes written and an error, if any.
   * WriteAt returns a non-nil error when n != len(b).
   * 
   * If file was opened with the O_APPEND flag, WriteAt returns an error.
   */
  writeAt(b: string|Array<number>, off: number): number
 }
 interface File {
  /**
   * WriteTo implements io.WriterTo.
   */
  writeTo(w: io.Writer): number
 }
 /**
  * noWriteTo can be embedded alongside another type to
  * hide the WriteTo method of that other type.
  */
 interface noWriteTo {
 }
 interface noWriteTo {
  /**
   * WriteTo hides another WriteTo method.
   * It should never be called.
   */
  writeTo(_arg0: io.Writer): number
 }
 /**
  * fileWithoutWriteTo implements all the methods of *File other
  * than WriteTo. This is used to permit WriteTo to call io.Copy
  * without leading to a recursive call to WriteTo.
  */
 type _sKepPvx = noWriteTo&File
 interface fileWithoutWriteTo extends _sKepPvx {
 }
 interface File {
  /**
   * Seek sets the offset for the next Read or Write on file to offset, interpreted
   * according to whence: 0 means relative to the origin of the file, 1 means
   * relative to the current offset, and 2 means relative to the end.
   * It returns the new offset and an error, if any.
   * The behavior of Seek on a file opened with O_APPEND is not specified.
   */
  seek(offset: number, whence: number): number
 }
 interface File {
  /**
   * WriteString is like Write, but writes the contents of string s rather than
   * a slice of bytes.
   */
  writeString(s: string): number
 }
 interface mkdir {
  /**
   * Mkdir creates a new directory with the specified name and permission
   * bits (before umask).
   * If there is an error, it will be of type *PathError.
   */
  (name: string, perm: FileMode): void
 }
 interface chdir {
  /**
   * Chdir changes the current working directory to the named directory.
   * If there is an error, it will be of type *PathError.
   */
  (dir: string): void
 }
 interface open {
  /**
   * Open opens the named file for reading. If successful, methods on
   * the returned file can be used for reading; the associated file
   * descriptor has mode O_RDONLY.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): (File)
 }
 interface create {
  /**
   * Create creates or truncates the named file. If the file already exists,
   * it is truncated. If the file does not exist, it is created with mode 0o666
   * (before umask). If successful, methods on the returned File can
   * be used for I/O; the associated file descriptor has mode O_RDWR.
   * The directory containing the file must already exist.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): (File)
 }
 interface openFile {
  /**
   * OpenFile is the generalized open call; most users will use Open
   * or Create instead. It opens the named file with specified flag
   * (O_RDONLY etc.). If the file does not exist, and the O_CREATE flag
   * is passed, it is created with mode perm (before umask);
   * the containing directory must exist. If successful,
   * methods on the returned File can be used for I/O.
   * If there is an error, it will be of type *PathError.
   */
  (name: string, flag: number, perm: FileMode): (File)
 }
 interface rename {
  /**
   * Rename renames (moves) oldpath to newpath.
   * If newpath already exists and is not a directory, Rename replaces it.
   * If newpath already exists and is a directory, Rename returns an error.
   * OS-specific restrictions may apply when oldpath and newpath are in different directories.
   * Even within the same directory, on non-Unix platforms Rename is not an atomic operation.
   * If there is an error, it will be of type *LinkError.
   */
  (oldpath: string, newpath: string): void
 }
 interface readlink {
  /**
   * Readlink returns the destination of the named symbolic link.
   * If there is an error, it will be of type *PathError.
   * 
   * If the link destination is relative, Readlink returns the relative path
   * without resolving it to an absolute one.
   */
  (name: string): string
 }
 interface tempDir {
  /**
   * TempDir returns the default directory to use for temporary files.
   * 
   * On Unix systems, it returns $TMPDIR if non-empty, else /tmp.
   * On Windows, it uses GetTempPath, returning the first non-empty
   * value from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.
   * On Plan 9, it returns /tmp.
   * 
   * The directory is neither guaranteed to exist nor have accessible
   * permissions.
   */
  (): string
 }
 interface userCacheDir {
  /**
   * UserCacheDir returns the default root directory to use for user-specific
   * cached data. Users should create their own application-specific subdirectory
   * within this one and use that.
   * 
   * On Unix systems, it returns $XDG_CACHE_HOME as specified by
   * https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
   * non-empty, else $HOME/.cache.
   * On Darwin, it returns $HOME/Library/Caches.
   * On Windows, it returns %LocalAppData%.
   * On Plan 9, it returns $home/lib/cache.
   * 
   * If the location cannot be determined (for example, $HOME is not defined) or
   * the path in $XDG_CACHE_HOME is relative, then it will return an error.
   */
  (): string
 }
 interface userConfigDir {
  /**
   * UserConfigDir returns the default root directory to use for user-specific
   * configuration data. Users should create their own application-specific
   * subdirectory within this one and use that.
   * 
   * On Unix systems, it returns $XDG_CONFIG_HOME as specified by
   * https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
   * non-empty, else $HOME/.config.
   * On Darwin, it returns $HOME/Library/Application Support.
   * On Windows, it returns %AppData%.
   * On Plan 9, it returns $home/lib.
   * 
   * If the location cannot be determined (for example, $HOME is not defined) or
   * the path in $XDG_CONFIG_HOME is relative, then it will return an error.
   */
  (): string
 }
 interface userHomeDir {
  /**
   * UserHomeDir returns the current user's home directory.
   * 
   * On Unix, including macOS, it returns the $HOME environment variable.
   * On Windows, it returns %USERPROFILE%.
   * On Plan 9, it returns the $home environment variable.
   * 
   * If the expected variable is not set in the environment, UserHomeDir
   * returns either a platform-specific default value or a non-nil error.
   */
  (): string
 }
 interface chmod {
  /**
   * Chmod changes the mode of the named file to mode.
   * If the file is a symbolic link, it changes the mode of the link's target.
   * If there is an error, it will be of type *PathError.
   * 
   * A different subset of the mode bits are used, depending on the
   * operating system.
   * 
   * On Unix, the mode's permission bits, ModeSetuid, ModeSetgid, and
   * ModeSticky are used.
   * 
   * On Windows, only the 0o200 bit (owner writable) of mode is used; it
   * controls whether the file's read-only attribute is set or cleared.
   * The other bits are currently unused. For compatibility with Go 1.12
   * and earlier, use a non-zero mode. Use mode 0o400 for a read-only
   * file and 0o600 for a readable+writable file.
   * 
   * On Plan 9, the mode's permission bits, ModeAppend, ModeExclusive,
   * and ModeTemporary are used.
   */
  (name: string, mode: FileMode): void
 }
 interface File {
  /**
   * Chmod changes the mode of the file to mode.
   * If there is an error, it will be of type *PathError.
   */
  chmod(mode: FileMode): void
 }
 interface File {
  /**
   * SetDeadline sets the read and write deadlines for a File.
   * It is equivalent to calling both SetReadDeadline and SetWriteDeadline.
   * 
   * Only some kinds of files support setting a deadline. Calls to SetDeadline
   * for files that do not support deadlines will return ErrNoDeadline.
   * On most systems ordinary files do not support deadlines, but pipes do.
   * 
   * A deadline is an absolute time after which I/O operations fail with an
   * error instead of blocking. The deadline applies to all future and pending
   * I/O, not just the immediately following call to Read or Write.
   * After a deadline has been exceeded, the connection can be refreshed
   * by setting a deadline in the future.
   * 
   * If the deadline is exceeded a call to Read or Write or to other I/O
   * methods will return an error that wraps ErrDeadlineExceeded.
   * This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
   * That error implements the Timeout method, and calling the Timeout
   * method will return true, but there are other possible errors for which
   * the Timeout will return true even if the deadline has not been exceeded.
   * 
   * An idle timeout can be implemented by repeatedly extending
   * the deadline after successful Read or Write calls.
   * 
   * A zero value for t means I/O operations will not time out.
   */
  setDeadline(t: time.Time): void
 }
 interface File {
  /**
   * SetReadDeadline sets the deadline for future Read calls and any
   * currently-blocked Read call.
   * A zero value for t means Read will not time out.
   * Not all files support setting deadlines; see SetDeadline.
   */
  setReadDeadline(t: time.Time): void
 }
 interface File {
  /**
   * SetWriteDeadline sets the deadline for any future Write calls and any
   * currently-blocked Write call.
   * Even if Write times out, it may return n > 0, indicating that
   * some of the data was successfully written.
   * A zero value for t means Write will not time out.
   * Not all files support setting deadlines; see SetDeadline.
   */
  setWriteDeadline(t: time.Time): void
 }
 interface File {
  /**
   * SyscallConn returns a raw file.
   * This implements the syscall.Conn interface.
   */
  syscallConn(): syscall.RawConn
 }
 interface dirFS {
  /**
   * DirFS returns a file system (an fs.FS) for the tree of files rooted at the directory dir.
   * 
   * Note that DirFS("/prefix") only guarantees that the Open calls it makes to the
   * operating system will begin with "/prefix": DirFS("/prefix").Open("file") is the
   * same as os.Open("/prefix/file"). So if /prefix/file is a symbolic link pointing outside
   * the /prefix tree, then using DirFS does not stop the access any more than using
   * os.Open does. Additionally, the root of the fs.FS returned for a relative path,
   * DirFS("prefix"), will be affected by later calls to Chdir. DirFS is therefore not
   * a general substitute for a chroot-style security mechanism when the directory tree
   * contains arbitrary content.
   * 
   * Use [Root.FS] to obtain a fs.FS that prevents escapes from the tree via symbolic links.
   * 
   * The directory dir must not be "".
   * 
   * The result implements [io/fs.StatFS], [io/fs.ReadFileFS] and
   * [io/fs.ReadDirFS].
   */
  (dir: string): fs.FS
 }
 interface dirFS extends String{}
 interface dirFS {
  open(name: string): fs.File
 }
 interface dirFS {
  /**
   * The ReadFile method calls the [ReadFile] function for the file
   * with the given name in the directory. The function provides
   * robust handling for small files and special file systems.
   * Through this method, dirFS implements [io/fs.ReadFileFS].
   */
  readFile(name: string): string|Array<number>
 }
 interface dirFS {
  /**
   * ReadDir reads the named directory, returning all its directory entries sorted
   * by filename. Through this method, dirFS implements [io/fs.ReadDirFS].
   */
  readDir(name: string): Array<DirEntry>
 }
 interface dirFS {
  stat(name: string): fs.FileInfo
 }
 interface readFile {
  /**
   * ReadFile reads the named file and returns the contents.
   * A successful call returns err == nil, not err == EOF.
   * Because ReadFile reads the whole file, it does not treat an EOF from Read
   * as an error to be reported.
   */
  (name: string): string|Array<number>
 }
 interface writeFile {
  /**
   * WriteFile writes data to the named file, creating it if necessary.
   * If the file does not exist, WriteFile creates it with permissions perm (before umask);
   * otherwise WriteFile truncates it before writing, without changing permissions.
   * Since WriteFile requires multiple system calls to complete, a failure mid-operation
   * can leave the file in a partially written state.
   */
  (name: string, data: string|Array<number>, perm: FileMode): void
 }
 interface File {
  /**
   * Close closes the [File], rendering it unusable for I/O.
   * On files that support [File.SetDeadline], any pending I/O operations will
   * be canceled and return immediately with an [ErrClosed] error.
   * Close will return an error if it has already been called.
   */
  close(): void
 }
 interface chown {
  /**
   * Chown changes the numeric uid and gid of the named file.
   * If the file is a symbolic link, it changes the uid and gid of the link's target.
   * A uid or gid of -1 means to not change that value.
   * If there is an error, it will be of type [*PathError].
   * 
   * On Windows or Plan 9, Chown always returns the [syscall.EWINDOWS] or
   * EPLAN9 error, wrapped in *PathError.
   */
  (name: string, uid: number, gid: number): void
 }
 interface lchown {
  /**
   * Lchown changes the numeric uid and gid of the named file.
   * If the file is a symbolic link, it changes the uid and gid of the link itself.
   * If there is an error, it will be of type [*PathError].
   * 
   * On Windows, it always returns the [syscall.EWINDOWS] error, wrapped
   * in *PathError.
   */
  (name: string, uid: number, gid: number): void
 }
 interface File {
  /**
   * Chown changes the numeric uid and gid of the named file.
   * If there is an error, it will be of type [*PathError].
   * 
   * On Windows, it always returns the [syscall.EWINDOWS] error, wrapped
   * in *PathError.
   */
  chown(uid: number, gid: number): void
 }
 interface File {
  /**
   * Truncate changes the size of the file.
   * It does not change the I/O offset.
   * If there is an error, it will be of type [*PathError].
   */
  truncate(size: number): void
 }
 interface File {
  /**
   * Sync commits the current contents of the file to stable storage.
   * Typically, this means flushing the file system's in-memory copy
   * of recently written data to disk.
   */
  sync(): void
 }
 interface chtimes {
  /**
   * Chtimes changes the access and modification times of the named
   * file, similar to the Unix utime() or utimes() functions.
   * A zero [time.Time] value will leave the corresponding file time unchanged.
   * 
   * The underlying filesystem may truncate or round the values to a
   * less precise time unit.
   * If there is an error, it will be of type [*PathError].
   */
  (name: string, atime: time.Time, mtime: time.Time): void
 }
 interface File {
  /**
   * Chdir changes the current working directory to the file,
   * which must be a directory.
   * If there is an error, it will be of type [*PathError].
   */
  chdir(): void
 }
 /**
  * file is the real representation of *File.
  * The extra level of indirection ensures that no clients of os
  * can overwrite this data, which could cause the finalizer
  * to close the wrong file descriptor.
  */
 interface file {
 }
 interface File {
  /**
   * Fd returns the integer Unix file descriptor referencing the open file.
   * If f is closed, the file descriptor becomes invalid.
   * If f is garbage collected, a finalizer may close the file descriptor,
   * making it invalid; see [runtime.SetFinalizer] for more information on when
   * a finalizer might be run. On Unix systems this will cause the [File.SetDeadline]
   * methods to stop working.
   * Because file descriptors can be reused, the returned file descriptor may
   * only be closed through the [File.Close] method of f, or by its finalizer during
   * garbage collection. Otherwise, during garbage collection the finalizer
   * may close an unrelated file descriptor with the same (reused) number.
   * 
   * As an alternative, see the f.SyscallConn method.
   */
  fd(): number
 }
 interface newFile {
  /**
   * NewFile returns a new File with the given file descriptor and
   * name. The returned value will be nil if fd is not a valid file
   * descriptor. On Unix systems, if the file descriptor is in
   * non-blocking mode, NewFile will attempt to return a pollable File
   * (one for which the SetDeadline methods work).
   * 
   * After passing it to NewFile, fd may become invalid under the same
   * conditions described in the comments of the Fd method, and the same
   * constraints apply.
   */
  (fd: number, name: string): (File)
 }
 /**
  * newFileKind describes the kind of file to newFile.
  */
 interface newFileKind extends Number{}
 interface truncate {
  /**
   * Truncate changes the size of the named file.
   * If the file is a symbolic link, it changes the size of the link's target.
   * If there is an error, it will be of type *PathError.
   */
  (name: string, size: number): void
 }
 interface remove {
  /**
   * Remove removes the named file or (empty) directory.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): void
 }
 interface link {
  /**
   * Link creates newname as a hard link to the oldname file.
   * If there is an error, it will be of type *LinkError.
   */
  (oldname: string, newname: string): void
 }
 interface symlink {
  /**
   * Symlink creates newname as a symbolic link to oldname.
   * On Windows, a symlink to a non-existent oldname creates a file symlink;
   * if oldname is later created as a directory the symlink will not work.
   * If there is an error, it will be of type *LinkError.
   */
  (oldname: string, newname: string): void
 }
 interface unixDirent {
 }
 interface unixDirent {
  name(): string
 }
 interface unixDirent {
  isDir(): boolean
 }
 interface unixDirent {
  type(): FileMode
 }
 interface unixDirent {
  info(): FileInfo
 }
 interface unixDirent {
  string(): string
 }
 interface getwd {
  /**
   * Getwd returns an absolute path name corresponding to the
   * current directory. If the current directory can be
   * reached via multiple paths (due to symbolic links),
   * Getwd may return any one of them.
   * 
   * On Unix platforms, if the environment variable PWD
   * provides an absolute name, and it is a name of the
   * current directory, it is returned.
   */
  (): string
 }
 interface mkdirAll {
  /**
   * MkdirAll creates a directory named path,
   * along with any necessary parents, and returns nil,
   * or else returns an error.
   * The permission bits perm (before umask) are used for all
   * directories that MkdirAll creates.
   * If path is already a directory, MkdirAll does nothing
   * and returns nil.
   */
  (path: string, perm: FileMode): void
 }
 interface removeAll {
  /**
   * RemoveAll removes path and any children it contains.
   * It removes everything it can but returns the first error
   * it encounters. If the path does not exist, RemoveAll
   * returns nil (no error).
   * If there is an error, it will be of type [*PathError].
   */
  (path: string): void
 }
 interface isPathSeparator {
  /**
   * IsPathSeparator reports whether c is a directory separator character.
   */
  (c: number): boolean
 }
 interface pipe {
  /**
   * Pipe returns a connected pair of Files; reads from r return bytes written to w.
   * It returns the files and an error, if any.
   */
  (): [(File), (File)]
 }
 interface getuid {
  /**
   * Getuid returns the numeric user id of the caller.
   * 
   * On Windows, it returns -1.
   */
  (): number
 }
 interface geteuid {
  /**
   * Geteuid returns the numeric effective user id of the caller.
   * 
   * On Windows, it returns -1.
   */
  (): number
 }
 interface getgid {
  /**
   * Getgid returns the numeric group id of the caller.
   * 
   * On Windows, it returns -1.
   */
  (): number
 }
 interface getegid {
  /**
   * Getegid returns the numeric effective group id of the caller.
   * 
   * On Windows, it returns -1.
   */
  (): number
 }
 interface getgroups {
  /**
   * Getgroups returns a list of the numeric ids of groups that the caller belongs to.
   * 
   * On Windows, it returns [syscall.EWINDOWS]. See the [os/user] package
   * for a possible alternative.
   */
  (): Array<number>
 }
 interface exit {
  /**
   * Exit causes the current program to exit with the given status code.
   * Conventionally, code zero indicates success, non-zero an error.
   * The program terminates immediately; deferred functions are not run.
   * 
   * For portability, the status code should be in the range [0, 125].
   */
  (code: number): void
 }
 /**
  * rawConn implements syscall.RawConn.
  */
 interface rawConn {
 }
 interface rawConn {
  control(f: (_arg0: number) => void): void
 }
 interface rawConn {
  read(f: (_arg0: number) => boolean): void
 }
 interface rawConn {
  write(f: (_arg0: number) => boolean): void
 }
 interface openInRoot {
  /**
   * OpenInRoot opens the file name in the directory dir.
   * It is equivalent to OpenRoot(dir) followed by opening the file in the root.
   * 
   * OpenInRoot returns an error if any component of the name
   * references a location outside of dir.
   * 
   * See [Root] for details and limitations.
   */
  (dir: string, name: string): (File)
 }
 /**
  * Root may be used to only access files within a single directory tree.
  * 
  * Methods on Root can only access files and directories beneath a root directory.
  * If any component of a file name passed to a method of Root references a location
  * outside the root, the method returns an error.
  * File names may reference the directory itself (.).
  * 
  * Methods on Root will follow symbolic links, but symbolic links may not
  * reference a location outside the root.
  * Symbolic links must not be absolute.
  * 
  * Methods on Root do not prohibit traversal of filesystem boundaries,
  * Linux bind mounts, /proc special files, or access to Unix device files.
  * 
  * Methods on Root are safe to be used from multiple goroutines simultaneously.
  * 
  * On most platforms, creating a Root opens a file descriptor or handle referencing
  * the directory. If the directory is moved, methods on Root reference the original
  * directory in its new location.
  * 
  * Root's behavior differs on some platforms:
  * 
  * ```
  *   - When GOOS=windows, file names may not reference Windows reserved device names
  *     such as NUL and COM1.
  *   - When GOOS=js, Root is vulnerable to TOCTOU (time-of-check-time-of-use)
  *     attacks in symlink validation, and cannot ensure that operations will not
  *     escape the root.
  *   - When GOOS=plan9 or GOOS=js, Root does not track directories across renames.
  *     On these platforms, a Root references a directory name, not a file descriptor.
  * ```
  */
 interface Root {
 }
 interface openRoot {
  /**
   * OpenRoot opens the named directory.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): (Root)
 }
 interface Root {
  /**
   * Name returns the name of the directory presented to OpenRoot.
   * 
   * It is safe to call Name after [Close].
   */
  name(): string
 }
 interface Root {
  /**
   * Close closes the Root.
   * After Close is called, methods on Root return errors.
   */
  close(): void
 }
 interface Root {
  /**
   * Open opens the named file in the root for reading.
   * See [Open] for more details.
   */
  open(name: string): (File)
 }
 interface Root {
  /**
   * Create creates or truncates the named file in the root.
   * See [Create] for more details.
   */
  create(name: string): (File)
 }
 interface Root {
  /**
   * OpenFile opens the named file in the root.
   * See [OpenFile] for more details.
   * 
   * If perm contains bits other than the nine least-significant bits (0o777),
   * OpenFile returns an error.
   */
  openFile(name: string, flag: number, perm: FileMode): (File)
 }
 interface Root {
  /**
   * OpenRoot opens the named directory in the root.
   * If there is an error, it will be of type *PathError.
   */
  openRoot(name: string): (Root)
 }
 interface Root {
  /**
   * Mkdir creates a new directory in the root
   * with the specified name and permission bits (before umask).
   * See [Mkdir] for more details.
   * 
   * If perm contains bits other than the nine least-significant bits (0o777),
   * OpenFile returns an error.
   */
  mkdir(name: string, perm: FileMode): void
 }
 interface Root {
  /**
   * Remove removes the named file or (empty) directory in the root.
   * See [Remove] for more details.
   */
  remove(name: string): void
 }
 interface Root {
  /**
   * Stat returns a [FileInfo] describing the named file in the root.
   * See [Stat] for more details.
   */
  stat(name: string): FileInfo
 }
 interface Root {
  /**
   * Lstat returns a [FileInfo] describing the named file in the root.
   * If the file is a symbolic link, the returned FileInfo
   * describes the symbolic link.
   * See [Lstat] for more details.
   */
  lstat(name: string): FileInfo
 }
 interface Root {
  /**
   * FS returns a file system (an fs.FS) for the tree of files in the root.
   * 
   * The result implements [io/fs.StatFS], [io/fs.ReadFileFS] and
   * [io/fs.ReadDirFS].
   */
  fs(): fs.FS
 }
 interface rootFS extends Root{}
 interface rootFS {
  open(name: string): fs.File
 }
 interface rootFS {
  readDir(name: string): Array<DirEntry>
 }
 interface rootFS {
  readFile(name: string): string|Array<number>
 }
 interface rootFS {
  stat(name: string): FileInfo
 }
 /**
  * root implementation for platforms with a function to open a file
  * relative to a directory.
  */
 interface root {
 }
 interface root {
  close(): void
 }
 interface root {
  name(): string
 }
 /**
  * errSymlink reports that a file being operated on is actually a symlink,
  * and the target of that symlink.
  */
 interface errSymlink extends String{}
 interface errSymlink {
  error(): string
 }
 interface sysfdType extends Number{}
 interface stat {
  /**
   * Stat returns a [FileInfo] describing the named file.
   * If there is an error, it will be of type [*PathError].
   */
  (name: string): FileInfo
 }
 interface lstat {
  /**
   * Lstat returns a [FileInfo] describing the named file.
   * If the file is a symbolic link, the returned FileInfo
   * describes the symbolic link. Lstat makes no attempt to follow the link.
   * If there is an error, it will be of type [*PathError].
   * 
   * On Windows, if the file is a reparse point that is a surrogate for another
   * named entity (such as a symbolic link or mounted folder), the returned
   * FileInfo describes the reparse point, and makes no attempt to resolve it.
   */
  (name: string): FileInfo
 }
 interface File {
  /**
   * Stat returns the [FileInfo] structure describing file.
   * If there is an error, it will be of type [*PathError].
   */
  stat(): FileInfo
 }
 interface hostname {
  /**
   * Hostname returns the host name reported by the kernel.
   */
  (): string
 }
 interface createTemp {
  /**
   * CreateTemp creates a new temporary file in the directory dir,
   * opens the file for reading and writing, and returns the resulting file.
   * The filename is generated by taking pattern and adding a random string to the end.
   * If pattern includes a "*", the random string replaces the last "*".
   * The file is created with mode 0o600 (before umask).
   * If dir is the empty string, CreateTemp uses the default directory for temporary files, as returned by [TempDir].
   * Multiple programs or goroutines calling CreateTemp simultaneously will not choose the same file.
   * The caller can use the file's Name method to find the pathname of the file.
   * It is the caller's responsibility to remove the file when it is no longer needed.
   */
  (dir: string, pattern: string): (File)
 }
 interface mkdirTemp {
  /**
   * MkdirTemp creates a new temporary directory in the directory dir
   * and returns the pathname of the new directory.
   * The new directory's name is generated by adding a random string to the end of pattern.
   * If pattern includes a "*", the random string replaces the last "*" instead.
   * The directory is created with mode 0o700 (before umask).
   * If dir is the empty string, MkdirTemp uses the default directory for temporary files, as returned by TempDir.
   * Multiple programs or goroutines calling MkdirTemp simultaneously will not choose the same directory.
   * It is the caller's responsibility to remove the directory when it is no longer needed.
   */
  (dir: string, pattern: string): string
 }
 interface getpagesize {
  /**
   * Getpagesize returns the underlying system's memory page size.
   */
  (): number
 }
 /**
  * File represents an open file descriptor.
  * 
  * The methods of File are safe for concurrent use.
  */
 type _sZngstS = file
 interface File extends _sZngstS {
 }
 /**
  * A FileInfo describes a file and is returned by [Stat] and [Lstat].
  */
 interface FileInfo extends fs.FileInfo{}
 /**
  * A FileMode represents a file's mode and permission bits.
  * The bits have the same definition on all systems, so that
  * information about files can be moved from one system
  * to another portably. Not all bits apply to all systems.
  * The only required bit is [ModeDir] for directories.
  */
 interface FileMode extends fs.FileMode{}
 interface fileStat {
  name(): string
 }
 interface fileStat {
  isDir(): boolean
 }
 interface sameFile {
  /**
   * SameFile reports whether fi1 and fi2 describe the same file.
   * For example, on Unix this means that the device and inode fields
   * of the two underlying structures are identical; on other systems
   * the decision may be based on the path names.
   * SameFile only applies to results returned by this package's [Stat].
   * It returns false in other cases.
   */
  (fi1: FileInfo, fi2: FileInfo): boolean
 }
 /**
  * A fileStat is the implementation of FileInfo returned by Stat and Lstat.
  */
 interface fileStat {
 }
 interface fileStat {
  size(): number
 }
 interface fileStat {
  mode(): FileMode
 }
 interface fileStat {
  modTime(): time.Time
 }
 interface fileStat {
  sys(): any
 }
}

/**
 * Package filepath implements utility routines for manipulating filename paths
 * in a way compatible with the target operating system-defined file paths.
 * 
 * The filepath package uses either forward slashes or backslashes,
 * depending on the operating system. To process paths such as URLs
 * that always use forward slashes regardless of the operating
 * system, see the [path] package.
 */
namespace filepath {
 interface match {
  /**
   * Match reports whether name matches the shell file name pattern.
   * The pattern syntax is:
   * 
   * ```
   * 	pattern:
   * 		{ term }
   * 	term:
   * 		'*'         matches any sequence of non-Separator characters
   * 		'?'         matches any single non-Separator character
   * 		'[' [ '^' ] { character-range } ']'
   * 		            character class (must be non-empty)
   * 		c           matches character c (c != '*', '?', '\\', '[')
   * 		'\\' c      matches character c
   * 
   * 	character-range:
   * 		c           matches character c (c != '\\', '-', ']')
   * 		'\\' c      matches character c
   * 		lo '-' hi   matches character c for lo <= c <= hi
   * ```
   * 
   * Match requires pattern to match all of name, not just a substring.
   * The only possible returned error is [ErrBadPattern], when pattern
   * is malformed.
   * 
   * On Windows, escaping is disabled. Instead, '\\' is treated as
   * path separator.
   */
  (pattern: string, name: string): boolean
 }
 interface glob {
  /**
   * Glob returns the names of all files matching pattern or nil
   * if there is no matching file. The syntax of patterns is the same
   * as in [Match]. The pattern may describe hierarchical names such as
   * /usr/*\/bin/ed (assuming the [Separator] is '/').
   * 
   * Glob ignores file system errors such as I/O errors reading directories.
   * The only possible returned error is [ErrBadPattern], when pattern
   * is malformed.
   */
  (pattern: string): Array<string>
 }
 interface clean {
  /**
   * Clean returns the shortest path name equivalent to path
   * by purely lexical processing. It applies the following rules
   * iteratively until no further processing can be done:
   * 
   *  1. Replace multiple [Separator] elements with a single one.
   *  2. Eliminate each . path name element (the current directory).
   *  3. Eliminate each inner .. path name element (the parent directory)
   * ```
   *     along with the non-.. element that precedes it.
   * ```
   *  4. Eliminate .. elements that begin a rooted path:
   * ```
   *     that is, replace "/.." by "/" at the beginning of a path,
   *     assuming Separator is '/'.
   * ```
   * 
   * The returned path ends in a slash only if it represents a root directory,
   * such as "/" on Unix or `C:\` on Windows.
   * 
   * Finally, any occurrences of slash are replaced by Separator.
   * 
   * If the result of this process is an empty string, Clean
   * returns the string ".".
   * 
   * On Windows, Clean does not modify the volume name other than to replace
   * occurrences of "/" with `\`.
   * For example, Clean("//host/share/../x") returns `\\host\share\x`.
   * 
   * See also Rob Pike, “Lexical File Names in Plan 9 or
   * Getting Dot-Dot Right,”
   * https://9p.io/sys/doc/lexnames.html
   */
  (path: string): string
 }
 interface isLocal {
  /**
   * IsLocal reports whether path, using lexical analysis only, has all of these properties:
   * 
   * ```
   *   - is within the subtree rooted at the directory in which path is evaluated
   *   - is not an absolute path
   *   - is not empty
   *   - on Windows, is not a reserved name such as "NUL"
   * ```
   * 
   * If IsLocal(path) returns true, then
   * Join(base, path) will always produce a path contained within base and
   * Clean(path) will always produce an unrooted path with no ".." path elements.
   * 
   * IsLocal is a purely lexical operation.
   * In particular, it does not account for the effect of any symbolic links
   * that may exist in the filesystem.
   */
  (path: string): boolean
 }
 interface localize {
  /**
   * Localize converts a slash-separated path into an operating system path.
   * The input path must be a valid path as reported by [io/fs.ValidPath].
   * 
   * Localize returns an error if the path cannot be represented by the operating system.
   * For example, the path a\b is rejected on Windows, on which \ is a separator
   * character and cannot be part of a filename.
   * 
   * The path returned by Localize will always be local, as reported by IsLocal.
   */
  (path: string): string
 }
 interface toSlash {
  /**
   * ToSlash returns the result of replacing each separator character
   * in path with a slash ('/') character. Multiple separators are
   * replaced by multiple slashes.
   */
  (path: string): string
 }
 interface fromSlash {
  /**
   * FromSlash returns the result of replacing each slash ('/') character
   * in path with a separator character. Multiple slashes are replaced
   * by multiple separators.
   * 
   * See also the Localize function, which converts a slash-separated path
   * as used by the io/fs package to an operating system path.
   */
  (path: string): string
 }
 interface splitList {
  /**
   * SplitList splits a list of paths joined by the OS-specific [ListSeparator],
   * usually found in PATH or GOPATH environment variables.
   * Unlike strings.Split, SplitList returns an empty slice when passed an empty
   * string.
   */
  (path: string): Array<string>
 }
 interface split {
  /**
   * Split splits path immediately following the final [Separator],
   * separating it into a directory and file name component.
   * If there is no Separator in path, Split returns an empty dir
   * and file set to path.
   * The returned values have the property that path = dir+file.
   */
  (path: string): [string, string]
 }
 interface join {
  /**
   * Join joins any number of path elements into a single path,
   * separating them with an OS specific [Separator]. Empty elements
   * are ignored. The result is Cleaned. However, if the argument
   * list is empty or all its elements are empty, Join returns
   * an empty string.
   * On Windows, the result will only be a UNC path if the first
   * non-empty element is a UNC path.
   */
  (...elem: string[]): string
 }
 interface ext {
  /**
   * Ext returns the file name extension used by path.
   * The extension is the suffix beginning at the final dot
   * in the final element of path; it is empty if there is
   * no dot.
   */
  (path: string): string
 }
 interface evalSymlinks {
  /**
   * EvalSymlinks returns the path name after the evaluation of any symbolic
   * links.
   * If path is relative the result will be relative to the current directory,
   * unless one of the components is an absolute symbolic link.
   * EvalSymlinks calls [Clean] on the result.
   */
  (path: string): string
 }
 interface isAbs {
  /**
   * IsAbs reports whether the path is absolute.
   */
  (path: string): boolean
 }
 interface abs {
  /**
   * Abs returns an absolute representation of path.
   * If the path is not absolute it will be joined with the current
   * working directory to turn it into an absolute path. The absolute
   * path name for a given file is not guaranteed to be unique.
   * Abs calls [Clean] on the result.
   */
  (path: string): string
 }
 interface rel {
  /**
   * Rel returns a relative path that is lexically equivalent to targpath when
   * joined to basepath with an intervening separator. That is,
   * [Join](basepath, Rel(basepath, targpath)) is equivalent to targpath itself.
   * On success, the returned path will always be relative to basepath,
   * even if basepath and targpath share no elements.
   * An error is returned if targpath can't be made relative to basepath or if
   * knowing the current working directory would be necessary to compute it.
   * Rel calls [Clean] on the result.
   */
  (basepath: string, targpath: string): string
 }
 /**
  * WalkFunc is the type of the function called by [Walk] to visit each
  * file or directory.
  * 
  * The path argument contains the argument to Walk as a prefix.
  * That is, if Walk is called with root argument "dir" and finds a file
  * named "a" in that directory, the walk function will be called with
  * argument "dir/a".
  * 
  * The directory and file are joined with Join, which may clean the
  * directory name: if Walk is called with the root argument "x/../dir"
  * and finds a file named "a" in that directory, the walk function will
  * be called with argument "dir/a", not "x/../dir/a".
  * 
  * The info argument is the fs.FileInfo for the named path.
  * 
  * The error result returned by the function controls how Walk continues.
  * If the function returns the special value [SkipDir], Walk skips the
  * current directory (path if info.IsDir() is true, otherwise path's
  * parent directory). If the function returns the special value [SkipAll],
  * Walk skips all remaining files and directories. Otherwise, if the function
  * returns a non-nil error, Walk stops entirely and returns that error.
  * 
  * The err argument reports an error related to path, signaling that Walk
  * will not walk into that directory. The function can decide how to
  * handle that error; as described earlier, returning the error will
  * cause Walk to stop walking the entire tree.
  * 
  * Walk calls the function with a non-nil err argument in two cases.
  * 
  * First, if an [os.Lstat] on the root directory or any directory or file
  * in the tree fails, Walk calls the function with path set to that
  * directory or file's path, info set to nil, and err set to the error
  * from os.Lstat.
  * 
  * Second, if a directory's Readdirnames method fails, Walk calls the
  * function with path set to the directory's path, info, set to an
  * [fs.FileInfo] describing the directory, and err set to the error from
  * Readdirnames.
  */
 interface WalkFunc {(path: string, info: fs.FileInfo, err: Error): void }
 interface walkDir {
  /**
   * WalkDir walks the file tree rooted at root, calling fn for each file or
   * directory in the tree, including root.
   * 
   * All errors that arise visiting files and directories are filtered by fn:
   * see the [fs.WalkDirFunc] documentation for details.
   * 
   * The files are walked in lexical order, which makes the output deterministic
   * but requires WalkDir to read an entire directory into memory before proceeding
   * to walk that directory.
   * 
   * WalkDir does not follow symbolic links.
   * 
   * WalkDir calls fn with paths that use the separator character appropriate
   * for the operating system. This is unlike [io/fs.WalkDir], which always
   * uses slash separated paths.
   */
  (root: string, fn: fs.WalkDirFunc): void
 }
 interface walk {
  /**
   * Walk walks the file tree rooted at root, calling fn for each file or
   * directory in the tree, including root.
   * 
   * All errors that arise visiting files and directories are filtered by fn:
   * see the [WalkFunc] documentation for details.
   * 
   * The files are walked in lexical order, which makes the output deterministic
   * but requires Walk to read an entire directory into memory before proceeding
   * to walk that directory.
   * 
   * Walk does not follow symbolic links.
   * 
   * Walk is less efficient than [WalkDir], introduced in Go 1.16,
   * which avoids calling os.Lstat on every visited file or directory.
   */
  (root: string, fn: WalkFunc): void
 }
 interface base {
  /**
   * Base returns the last element of path.
   * Trailing path separators are removed before extracting the last element.
   * If the path is empty, Base returns ".".
   * If the path consists entirely of separators, Base returns a single separator.
   */
  (path: string): string
 }
 interface dir {
  /**
   * Dir returns all but the last element of path, typically the path's directory.
   * After dropping the final element, Dir calls [Clean] on the path and trailing
   * slashes are removed.
   * If the path is empty, Dir returns ".".
   * If the path consists entirely of separators, Dir returns a single separator.
   * The returned path does not end in a separator unless it is the root directory.
   */
  (path: string): string
 }
 interface volumeName {
  /**
   * VolumeName returns leading volume name.
   * Given "C:\foo\bar" it returns "C:" on Windows.
   * Given "\\host\share\foo" it returns "\\host\share".
   * On other platforms it returns "".
   */
  (path: string): string
 }
 interface hasPrefix {
  /**
   * HasPrefix exists for historical compatibility and should not be used.
   * 
   * Deprecated: HasPrefix does not respect path boundaries and
   * does not ignore case when required.
   */
  (p: string, prefix: string): boolean
 }
}

/**
 * Package validation provides configurable and extensible rules for validating data of various types.
 */
namespace ozzo_validation {
 /**
  * Error interface represents an validation error
  */
 interface Error {
  [key:string]: any;
  error(): string
  code(): string
  message(): string
  setMessage(_arg0: string): Error
  params(): _TygojaDict
  setParams(_arg0: _TygojaDict): Error
 }
}

/**
 * Package dbx provides a set of DB-agnostic and easy-to-use query building methods for relational databases.
 */
namespace dbx {
 /**
  * Builder supports building SQL statements in a DB-agnostic way.
  * Builder mainly provides two sets of query building methods: those building SELECT statements
  * and those manipulating DB data or schema (e.g. INSERT statements, CREATE TABLE statements).
  */
 interface Builder {
  [key:string]: any;
  /**
   * NewQuery creates a new Query object with the given SQL statement.
   * The SQL statement may contain parameter placeholders which can be bound with actual parameter
   * values before the statement is executed.
   */
  newQuery(_arg0: string): (Query)
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(..._arg0: string[]): (SelectQuery)
  /**
   * ModelQuery returns a new ModelQuery object that can be used to perform model insertion, update, and deletion.
   * The parameter to this method should be a pointer to the model struct that needs to be inserted, updated, or deleted.
   */
  model(_arg0: {
  }): (ModelQuery)
  /**
   * GeneratePlaceholder generates an anonymous parameter placeholder with the given parameter ID.
   */
  generatePlaceholder(_arg0: number): string
  /**
   * Quote quotes a string so that it can be embedded in a SQL statement as a string value.
   */
  quote(_arg0: string): string
  /**
   * QuoteSimpleTableName quotes a simple table name.
   * A simple table name does not contain any schema prefix.
   */
  quoteSimpleTableName(_arg0: string): string
  /**
   * QuoteSimpleColumnName quotes a simple column name.
   * A simple column name does not contain any table prefix.
   */
  quoteSimpleColumnName(_arg0: string): string
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
  /**
   * Insert creates a Query that represents an INSERT SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  insert(table: string, cols: Params): (Query)
  /**
   * Upsert creates a Query that represents an UPSERT SQL statement.
   * Upsert inserts a row into the table if the primary key or unique index is not found.
   * Otherwise it will update the row with the new values.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  upsert(table: string, cols: Params, ...constraints: string[]): (Query)
  /**
   * Update creates a Query that represents an UPDATE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding new column
   * values. If the "where" expression is nil, the UPDATE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will update ALL rows in the table).
   */
  update(table: string, cols: Params, where: Expression): (Query)
  /**
   * Delete creates a Query that represents a DELETE SQL statement.
   * If the "where" expression is nil, the DELETE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will delete ALL rows in the table).
   */
  delete(table: string, where: Expression): (Query)
  /**
   * CreateTable creates a Query that represents a CREATE TABLE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding column types.
   * The optional "options" parameters will be appended to the generated SQL statement.
   */
  createTable(table: string, cols: _TygojaDict, ...options: string[]): (Query)
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string, newName: string): (Query)
  /**
   * DropTable creates a Query that can be used to drop a table.
   */
  dropTable(table: string): (Query)
  /**
   * TruncateTable creates a Query that can be used to truncate a table.
   */
  truncateTable(table: string): (Query)
  /**
   * AddColumn creates a Query that can be used to add a column to a table.
   */
  addColumn(table: string, col: string, typ: string): (Query)
  /**
   * DropColumn creates a Query that can be used to drop a column from a table.
   */
  dropColumn(table: string, col: string): (Query)
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string, oldName: string, newName: string): (Query)
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string, col: string, typ: string): (Query)
  /**
   * AddPrimaryKey creates a Query that can be used to specify primary key(s) for a table.
   * The "name" parameter specifies the name of the primary key constraint.
   */
  addPrimaryKey(table: string, name: string, ...cols: string[]): (Query)
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string, name: string): (Query)
  /**
   * AddForeignKey creates a Query that can be used to add a foreign key constraint to a table.
   * The length of cols and refCols must be the same as they refer to the primary and referential columns.
   * The optional "options" parameters will be appended to the SQL statement. They can be used to
   * specify options such as "ON DELETE CASCADE".
   */
  addForeignKey(table: string, name: string, cols: Array<string>, refCols: Array<string>, refTable: string, ...options: string[]): (Query)
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string, name: string): (Query)
  /**
   * CreateIndex creates a Query that can be used to create an index for a table.
   */
  createIndex(table: string, name: string, ...cols: string[]): (Query)
  /**
   * CreateUniqueIndex creates a Query that can be used to create a unique index for a table.
   */
  createUniqueIndex(table: string, name: string, ...cols: string[]): (Query)
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string, name: string): (Query)
 }
 /**
  * BaseBuilder provides a basic implementation of the Builder interface.
  */
 interface BaseBuilder {
 }
 interface newBaseBuilder {
  /**
   * NewBaseBuilder creates a new BaseBuilder instance.
   */
  (db: DB, executor: Executor): (BaseBuilder)
 }
 interface BaseBuilder {
  /**
   * DB returns the DB instance that this builder is associated with.
   */
  db(): (DB)
 }
 interface BaseBuilder {
  /**
   * Executor returns the executor object (a DB instance or a transaction) for executing SQL statements.
   */
  executor(): Executor
 }
 interface BaseBuilder {
  /**
   * NewQuery creates a new Query object with the given SQL statement.
   * The SQL statement may contain parameter placeholders which can be bound with actual parameter
   * values before the statement is executed.
   */
  newQuery(sql: string): (Query)
 }
 interface BaseBuilder {
  /**
   * GeneratePlaceholder generates an anonymous parameter placeholder with the given parameter ID.
   */
  generatePlaceholder(_arg0: number): string
 }
 interface BaseBuilder {
  /**
   * Quote quotes a string so that it can be embedded in a SQL statement as a string value.
   */
  quote(s: string): string
 }
 interface BaseBuilder {
  /**
   * QuoteSimpleTableName quotes a simple table name.
   * A simple table name does not contain any schema prefix.
   */
  quoteSimpleTableName(s: string): string
 }
 interface BaseBuilder {
  /**
   * QuoteSimpleColumnName quotes a simple column name.
   * A simple column name does not contain any table prefix.
   */
  quoteSimpleColumnName(s: string): string
 }
 interface BaseBuilder {
  /**
   * Insert creates a Query that represents an INSERT SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  insert(table: string, cols: Params): (Query)
 }
 interface BaseBuilder {
  /**
   * Upsert creates a Query that represents an UPSERT SQL statement.
   * Upsert inserts a row into the table if the primary key or unique index is not found.
   * Otherwise it will update the row with the new values.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  upsert(table: string, cols: Params, ...constraints: string[]): (Query)
 }
 interface BaseBuilder {
  /**
   * Update creates a Query that represents an UPDATE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding new column
   * values. If the "where" expression is nil, the UPDATE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will update ALL rows in the table).
   */
  update(table: string, cols: Params, where: Expression): (Query)
 }
 interface BaseBuilder {
  /**
   * Delete creates a Query that represents a DELETE SQL statement.
   * If the "where" expression is nil, the DELETE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will delete ALL rows in the table).
   */
  delete(table: string, where: Expression): (Query)
 }
 interface BaseBuilder {
  /**
   * CreateTable creates a Query that represents a CREATE TABLE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding column types.
   * The optional "options" parameters will be appended to the generated SQL statement.
   */
  createTable(table: string, cols: _TygojaDict, ...options: string[]): (Query)
 }
 interface BaseBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string, newName: string): (Query)
 }
 interface BaseBuilder {
  /**
   * DropTable creates a Query that can be used to drop a table.
   */
  dropTable(table: string): (Query)
 }
 interface BaseBuilder {
  /**
   * TruncateTable creates a Query that can be used to truncate a table.
   */
  truncateTable(table: string): (Query)
 }
 interface BaseBuilder {
  /**
   * AddColumn creates a Query that can be used to add a column to a table.
   */
  addColumn(table: string, col: string, typ: string): (Query)
 }
 interface BaseBuilder {
  /**
   * DropColumn creates a Query that can be used to drop a column from a table.
   */
  dropColumn(table: string, col: string): (Query)
 }
 interface BaseBuilder {
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string, oldName: string, newName: string): (Query)
 }
 interface BaseBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string, col: string, typ: string): (Query)
 }
 interface BaseBuilder {
  /**
   * AddPrimaryKey creates a Query that can be used to specify primary key(s) for a table.
   * The "name" parameter specifies the name of the primary key constraint.
   */
  addPrimaryKey(table: string, name: string, ...cols: string[]): (Query)
 }
 interface BaseBuilder {
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string, name: string): (Query)
 }
 interface BaseBuilder {
  /**
   * AddForeignKey creates a Query that can be used to add a foreign key constraint to a table.
   * The length of cols and refCols must be the same as they refer to the primary and referential columns.
   * The optional "options" parameters will be appended to the SQL statement. They can be used to
   * specify options such as "ON DELETE CASCADE".
   */
  addForeignKey(table: string, name: string, cols: Array<string>, refCols: Array<string>, refTable: string, ...options: string[]): (Query)
 }
 interface BaseBuilder {
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string, name: string): (Query)
 }
 interface BaseBuilder {
  /**
   * CreateIndex creates a Query that can be used to create an index for a table.
   */
  createIndex(table: string, name: string, ...cols: string[]): (Query)
 }
 interface BaseBuilder {
  /**
   * CreateUniqueIndex creates a Query that can be used to create a unique index for a table.
   */
  createUniqueIndex(table: string, name: string, ...cols: string[]): (Query)
 }
 interface BaseBuilder {
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string, name: string): (Query)
 }
 /**
  * MssqlBuilder is the builder for SQL Server databases.
  */
 type _snSipwC = BaseBuilder
 interface MssqlBuilder extends _snSipwC {
 }
 /**
  * MssqlQueryBuilder is the query builder for SQL Server databases.
  */
 type _sarHbbt = BaseQueryBuilder
 interface MssqlQueryBuilder extends _sarHbbt {
 }
 interface newMssqlBuilder {
  /**
   * NewMssqlBuilder creates a new MssqlBuilder instance.
   */
  (db: DB, executor: Executor): Builder
 }
 interface MssqlBuilder {
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
 }
 interface MssqlBuilder {
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface MssqlBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery)
 }
 interface MssqlBuilder {
  /**
   * QuoteSimpleTableName quotes a simple table name.
   * A simple table name does not contain any schema prefix.
   */
  quoteSimpleTableName(s: string): string
 }
 interface MssqlBuilder {
  /**
   * QuoteSimpleColumnName quotes a simple column name.
   * A simple column name does not contain any table prefix.
   */
  quoteSimpleColumnName(s: string): string
 }
 interface MssqlBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string, newName: string): (Query)
 }
 interface MssqlBuilder {
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string, oldName: string, newName: string): (Query)
 }
 interface MssqlBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string, col: string, typ: string): (Query)
 }
 interface MssqlQueryBuilder {
  /**
   * BuildOrderByAndLimit generates the ORDER BY and LIMIT clauses.
   */
  buildOrderByAndLimit(sql: string, cols: Array<string>, limit: number, offset: number): string
 }
 /**
  * MysqlBuilder is the builder for MySQL databases.
  */
 type _sVuRlWe = BaseBuilder
 interface MysqlBuilder extends _sVuRlWe {
 }
 interface newMysqlBuilder {
  /**
   * NewMysqlBuilder creates a new MysqlBuilder instance.
   */
  (db: DB, executor: Executor): Builder
 }
 interface MysqlBuilder {
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
 }
 interface MysqlBuilder {
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface MysqlBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery)
 }
 interface MysqlBuilder {
  /**
   * QuoteSimpleTableName quotes a simple table name.
   * A simple table name does not contain any schema prefix.
   */
  quoteSimpleTableName(s: string): string
 }
 interface MysqlBuilder {
  /**
   * QuoteSimpleColumnName quotes a simple column name.
   * A simple column name does not contain any table prefix.
   */
  quoteSimpleColumnName(s: string): string
 }
 interface MysqlBuilder {
  /**
   * Upsert creates a Query that represents an UPSERT SQL statement.
   * Upsert inserts a row into the table if the primary key or unique index is not found.
   * Otherwise it will update the row with the new values.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  upsert(table: string, cols: Params, ...constraints: string[]): (Query)
 }
 interface MysqlBuilder {
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string, oldName: string, newName: string): (Query)
 }
 interface MysqlBuilder {
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string, name: string): (Query)
 }
 interface MysqlBuilder {
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string, name: string): (Query)
 }
 /**
  * OciBuilder is the builder for Oracle databases.
  */
 type _syrzuTA = BaseBuilder
 interface OciBuilder extends _syrzuTA {
 }
 /**
  * OciQueryBuilder is the query builder for Oracle databases.
  */
 type _sXjWjty = BaseQueryBuilder
 interface OciQueryBuilder extends _sXjWjty {
 }
 interface newOciBuilder {
  /**
   * NewOciBuilder creates a new OciBuilder instance.
   */
  (db: DB, executor: Executor): Builder
 }
 interface OciBuilder {
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface OciBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery)
 }
 interface OciBuilder {
  /**
   * GeneratePlaceholder generates an anonymous parameter placeholder with the given parameter ID.
   */
  generatePlaceholder(i: number): string
 }
 interface OciBuilder {
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
 }
 interface OciBuilder {
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string, name: string): (Query)
 }
 interface OciBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string, newName: string): (Query)
 }
 interface OciBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string, col: string, typ: string): (Query)
 }
 interface OciQueryBuilder {
  /**
   * BuildOrderByAndLimit generates the ORDER BY and LIMIT clauses.
   */
  buildOrderByAndLimit(sql: string, cols: Array<string>, limit: number, offset: number): string
 }
 /**
  * PgsqlBuilder is the builder for PostgreSQL databases.
  */
 type _sSFIxNF = BaseBuilder
 interface PgsqlBuilder extends _sSFIxNF {
 }
 interface newPgsqlBuilder {
  /**
   * NewPgsqlBuilder creates a new PgsqlBuilder instance.
   */
  (db: DB, executor: Executor): Builder
 }
 interface PgsqlBuilder {
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface PgsqlBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery)
 }
 interface PgsqlBuilder {
  /**
   * GeneratePlaceholder generates an anonymous parameter placeholder with the given parameter ID.
   */
  generatePlaceholder(i: number): string
 }
 interface PgsqlBuilder {
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
 }
 interface PgsqlBuilder {
  /**
   * Upsert creates a Query that represents an UPSERT SQL statement.
   * Upsert inserts a row into the table if the primary key or unique index is not found.
   * Otherwise it will update the row with the new values.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  upsert(table: string, cols: Params, ...constraints: string[]): (Query)
 }
 interface PgsqlBuilder {
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string, name: string): (Query)
 }
 interface PgsqlBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string, newName: string): (Query)
 }
 interface PgsqlBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string, col: string, typ: string): (Query)
 }
 /**
  * SqliteBuilder is the builder for SQLite databases.
  */
 type _sTHEwza = BaseBuilder
 interface SqliteBuilder extends _sTHEwza {
 }
 interface newSqliteBuilder {
  /**
   * NewSqliteBuilder creates a new SqliteBuilder instance.
   */
  (db: DB, executor: Executor): Builder
 }
 interface SqliteBuilder {
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
 }
 interface SqliteBuilder {
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface SqliteBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery)
 }
 interface SqliteBuilder {
  /**
   * QuoteSimpleTableName quotes a simple table name.
   * A simple table name does not contain any schema prefix.
   */
  quoteSimpleTableName(s: string): string
 }
 interface SqliteBuilder {
  /**
   * QuoteSimpleColumnName quotes a simple column name.
   * A simple column name does not contain any table prefix.
   */
  quoteSimpleColumnName(s: string): string
 }
 interface SqliteBuilder {
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string, name: string): (Query)
 }
 interface SqliteBuilder {
  /**
   * TruncateTable creates a Query that can be used to truncate a table.
   */
  truncateTable(table: string): (Query)
 }
 interface SqliteBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string, newName: string): (Query)
 }
 interface SqliteBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string, col: string, typ: string): (Query)
 }
 interface SqliteBuilder {
  /**
   * AddPrimaryKey creates a Query that can be used to specify primary key(s) for a table.
   * The "name" parameter specifies the name of the primary key constraint.
   */
  addPrimaryKey(table: string, name: string, ...cols: string[]): (Query)
 }
 interface SqliteBuilder {
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string, name: string): (Query)
 }
 interface SqliteBuilder {
  /**
   * AddForeignKey creates a Query that can be used to add a foreign key constraint to a table.
   * The length of cols and refCols must be the same as they refer to the primary and referential columns.
   * The optional "options" parameters will be appended to the SQL statement. They can be used to
   * specify options such as "ON DELETE CASCADE".
   */
  addForeignKey(table: string, name: string, cols: Array<string>, refCols: Array<string>, refTable: string, ...options: string[]): (Query)
 }
 interface SqliteBuilder {
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string, name: string): (Query)
 }
 /**
  * StandardBuilder is the builder that is used by DB for an unknown driver.
  */
 type _sdrMQZN = BaseBuilder
 interface StandardBuilder extends _sdrMQZN {
 }
 interface newStandardBuilder {
  /**
   * NewStandardBuilder creates a new StandardBuilder instance.
   */
  (db: DB, executor: Executor): Builder
 }
 interface StandardBuilder {
  /**
   * QueryBuilder returns the query builder supporting the current DB.
   */
  queryBuilder(): QueryBuilder
 }
 interface StandardBuilder {
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface StandardBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery)
 }
 /**
  * LogFunc logs a message for each SQL statement being executed.
  * This method takes one or multiple parameters. If a single parameter
  * is provided, it will be treated as the log message. If multiple parameters
  * are provided, they will be passed to fmt.Sprintf() to generate the log message.
  */
 interface LogFunc {(format: string, ...a: {
  }[]): void }
 /**
  * PerfFunc is called when a query finishes execution.
  * The query execution time is passed to this function so that the DB performance
  * can be profiled. The "ns" parameter gives the number of nanoseconds that the
  * SQL statement takes to execute, while the "execute" parameter indicates whether
  * the SQL statement is executed or queried (usually SELECT statements).
  */
 interface PerfFunc {(ns: number, sql: string, execute: boolean): void }
 /**
  * QueryLogFunc is called each time when performing a SQL query.
  * The "t" parameter gives the time that the SQL statement takes to execute,
  * while rows and err are the result of the query.
  */
 interface QueryLogFunc {(ctx: context.Context, t: time.Duration, sql: string, rows: sql.Rows, err: Error): void }
 /**
  * ExecLogFunc is called each time when a SQL statement is executed.
  * The "t" parameter gives the time that the SQL statement takes to execute,
  * while result and err refer to the result of the execution.
  */
 interface ExecLogFunc {(ctx: context.Context, t: time.Duration, sql: string, result: sql.Result, err: Error): void }
 /**
  * BuilderFunc creates a Builder instance using the given DB instance and Executor.
  */
 interface BuilderFunc {(_arg0: DB, _arg1: Executor): Builder }
 /**
  * DB enhances sql.DB by providing a set of DB-agnostic query building methods.
  * DB allows easier query building and population of data into Go variables.
  */
 type _sibLwPB = Builder
 interface DB extends _sibLwPB {
  /**
   * FieldMapper maps struct fields to DB columns. Defaults to DefaultFieldMapFunc.
   */
  fieldMapper: FieldMapFunc
  /**
   * TableMapper maps structs to table names. Defaults to GetTableName.
   */
  tableMapper: TableMapFunc
  /**
   * LogFunc logs the SQL statements being executed. Defaults to nil, meaning no logging.
   */
  logFunc: LogFunc
  /**
   * PerfFunc logs the SQL execution time. Defaults to nil, meaning no performance profiling.
   * Deprecated: Please use QueryLogFunc and ExecLogFunc instead.
   */
  perfFunc: PerfFunc
  /**
   * QueryLogFunc is called each time when performing a SQL query that returns data.
   */
  queryLogFunc: QueryLogFunc
  /**
   * ExecLogFunc is called each time when a SQL statement is executed.
   */
  execLogFunc: ExecLogFunc
 }
 /**
  * Errors represents a list of errors.
  */
 interface Errors extends Array<Error>{}
 interface newFromDB {
  /**
   * NewFromDB encapsulates an existing database connection.
   */
  (sqlDB: sql.DB, driverName: string): (DB)
 }
 interface open {
  /**
   * Open opens a database specified by a driver name and data source name (DSN).
   * Note that Open does not check if DSN is specified correctly. It doesn't try to establish a DB connection either.
   * Please refer to sql.Open() for more information.
   */
  (driverName: string, dsn: string): (DB)
 }
 interface mustOpen {
  /**
   * MustOpen opens a database and establishes a connection to it.
   * Please refer to sql.Open() and sql.Ping() for more information.
   */
  (driverName: string, dsn: string): (DB)
 }
 interface DB {
  /**
   * Clone makes a shallow copy of DB.
   */
  clone(): (DB)
 }
 interface DB {
  /**
   * WithContext returns a new instance of DB associated with the given context.
   */
  withContext(ctx: context.Context): (DB)
 }
 interface DB {
  /**
   * Context returns the context associated with the DB instance.
   * It returns nil if no context is associated.
   */
  context(): context.Context
 }
 interface DB {
  /**
   * DB returns the sql.DB instance encapsulated by dbx.DB.
   */
  db(): (sql.DB)
 }
 interface DB {
  /**
   * Close closes the database, releasing any open resources.
   * It is rare to Close a DB, as the DB handle is meant to be
   * long-lived and shared between many goroutines.
   */
  close(): void
 }
 interface DB {
  /**
   * Begin starts a transaction.
   */
  begin(): (Tx)
 }
 interface DB {
  /**
   * BeginTx starts a transaction with the given context and transaction options.
   */
  beginTx(ctx: context.Context, opts: sql.TxOptions): (Tx)
 }
 interface DB {
  /**
   * Wrap encapsulates an existing transaction.
   */
  wrap(sqlTx: sql.Tx): (Tx)
 }
 interface DB {
  /**
   * Transactional starts a transaction and executes the given function.
   * If the function returns an error, the transaction will be rolled back.
   * Otherwise, the transaction will be committed.
   */
  transactional(f: (_arg0: Tx) => void): void
 }
 interface DB {
  /**
   * TransactionalContext starts a transaction and executes the given function with the given context and transaction options.
   * If the function returns an error, the transaction will be rolled back.
   * Otherwise, the transaction will be committed.
   */
  transactionalContext(ctx: context.Context, opts: sql.TxOptions, f: (_arg0: Tx) => void): void
 }
 interface DB {
  /**
   * DriverName returns the name of the DB driver.
   */
  driverName(): string
 }
 interface DB {
  /**
   * QuoteTableName quotes the given table name appropriately.
   * If the table name contains DB schema prefix, it will be handled accordingly.
   * This method will do nothing if the table name is already quoted or if it contains parenthesis.
   */
  quoteTableName(s: string): string
 }
 interface DB {
  /**
   * QuoteColumnName quotes the given column name appropriately.
   * If the table name contains table name prefix, it will be handled accordingly.
   * This method will do nothing if the column name is already quoted or if it contains parenthesis.
   */
  quoteColumnName(s: string): string
 }
 interface Errors {
  /**
   * Error returns the error string of Errors.
   */
  error(): string
 }
 /**
  * Expression represents a DB expression that can be embedded in a SQL statement.
  */
 interface Expression {
  [key:string]: any;
  /**
   * Build converts an expression into a SQL fragment.
   * If the expression contains binding parameters, they will be added to the given Params.
   */
  build(_arg0: DB, _arg1: Params): string
 }
 /**
  * HashExp represents a hash expression.
  * 
  * A hash expression is a map whose keys are DB column names which need to be filtered according
  * to the corresponding values. For example, HashExp{"level": 2, "dept": 10} will generate
  * the SQL: "level"=2 AND "dept"=10.
  * 
  * HashExp also handles nil values and slice values. For example, HashExp{"level": []interface{}{1, 2}, "dept": nil}
  * will generate: "level" IN (1, 2) AND "dept" IS NULL.
  */
 interface HashExp extends _TygojaDict{}
 interface newExp {
  /**
   * NewExp generates an expression with the specified SQL fragment and the optional binding parameters.
   */
  (e: string, ...params: Params[]): Expression
 }
 interface not {
  /**
   * Not generates a NOT expression which prefixes "NOT" to the specified expression.
   */
  (e: Expression): Expression
 }
 interface and {
  /**
   * And generates an AND expression which concatenates the given expressions with "AND".
   */
  (...exps: Expression[]): Expression
 }
 interface or {
  /**
   * Or generates an OR expression which concatenates the given expressions with "OR".
   */
  (...exps: Expression[]): Expression
 }
 interface _in {
  /**
   * In generates an IN expression for the specified column and the list of allowed values.
   * If values is empty, a SQL "0=1" will be generated which represents a false expression.
   */
  (col: string, ...values: {
   }[]): Expression
 }
 interface notIn {
  /**
   * NotIn generates an NOT IN expression for the specified column and the list of disallowed values.
   * If values is empty, an empty string will be returned indicating a true expression.
   */
  (col: string, ...values: {
   }[]): Expression
 }
 interface like {
  /**
   * Like generates a LIKE expression for the specified column and the possible strings that the column should be like.
   * If multiple values are present, the column should be like *all* of them. For example, Like("name", "key", "word")
   * will generate a SQL expression: "name" LIKE "%key%" AND "name" LIKE "%word%".
   * 
   * By default, each value will be surrounded by "%" to enable partial matching. If a value contains special characters
   * such as "%", "\", "_", they will also be properly escaped.
   * 
   * You may call Escape() and/or Match() to change the default behavior. For example, Like("name", "key").Match(false, true)
   * generates "name" LIKE "key%".
   */
  (col: string, ...values: string[]): (LikeExp)
 }
 interface notLike {
  /**
   * NotLike generates a NOT LIKE expression.
   * For example, NotLike("name", "key", "word") will generate a SQL expression:
   * "name" NOT LIKE "%key%" AND "name" NOT LIKE "%word%". Please see Like() for more details.
   */
  (col: string, ...values: string[]): (LikeExp)
 }
 interface orLike {
  /**
   * OrLike generates an OR LIKE expression.
   * This is similar to Like() except that the column should be like one of the possible values.
   * For example, OrLike("name", "key", "word") will generate a SQL expression:
   * "name" LIKE "%key%" OR "name" LIKE "%word%". Please see Like() for more details.
   */
  (col: string, ...values: string[]): (LikeExp)
 }
 interface orNotLike {
  /**
   * OrNotLike generates an OR NOT LIKE expression.
   * For example, OrNotLike("name", "key", "word") will generate a SQL expression:
   * "name" NOT LIKE "%key%" OR "name" NOT LIKE "%word%". Please see Like() for more details.
   */
  (col: string, ...values: string[]): (LikeExp)
 }
 interface exists {
  /**
   * Exists generates an EXISTS expression by prefixing "EXISTS" to the given expression.
   */
  (exp: Expression): Expression
 }
 interface notExists {
  /**
   * NotExists generates an EXISTS expression by prefixing "NOT EXISTS" to the given expression.
   */
  (exp: Expression): Expression
 }
 interface between {
  /**
   * Between generates a BETWEEN expression.
   * For example, Between("age", 10, 30) generates: "age" BETWEEN 10 AND 30
   */
  (col: string, from: {
   }, to: {
   }): Expression
 }
 interface notBetween {
  /**
   * NotBetween generates a NOT BETWEEN expression.
   * For example, NotBetween("age", 10, 30) generates: "age" NOT BETWEEN 10 AND 30
   */
  (col: string, from: {
   }, to: {
   }): Expression
 }
 /**
  * Exp represents an expression with a SQL fragment and a list of optional binding parameters.
  */
 interface Exp {
 }
 interface Exp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 interface HashExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * NotExp represents an expression that should prefix "NOT" to a specified expression.
  */
 interface NotExp {
 }
 interface NotExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * AndOrExp represents an expression that concatenates multiple expressions using either "AND" or "OR".
  */
 interface AndOrExp {
 }
 interface AndOrExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * InExp represents an "IN" or "NOT IN" expression.
  */
 interface InExp {
 }
 interface InExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * LikeExp represents a variant of LIKE expressions.
  */
 interface LikeExp {
  /**
   * Like stores the LIKE operator. It can be "LIKE", "NOT LIKE".
   * It may also be customized as something like "ILIKE".
   */
  like: string
 }
 interface LikeExp {
  /**
   * Escape specifies how a LIKE expression should be escaped.
   * Each string at position 2i represents a special character and the string at position 2i+1 is
   * the corresponding escaped version.
   */
  escape(...chars: string[]): (LikeExp)
 }
 interface LikeExp {
  /**
   * Match specifies whether to do wildcard matching on the left and/or right of given strings.
   */
  match(left: boolean, right: boolean): (LikeExp)
 }
 interface LikeExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * ExistsExp represents an EXISTS or NOT EXISTS expression.
  */
 interface ExistsExp {
 }
 interface ExistsExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * BetweenExp represents a BETWEEN or a NOT BETWEEN expression.
  */
 interface BetweenExp {
 }
 interface BetweenExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 interface enclose {
  /**
   * Enclose surrounds the provided nonempty expression with parenthesis "()".
   */
  (exp: Expression): Expression
 }
 /**
  * EncloseExp represents a parenthesis enclosed expression.
  */
 interface EncloseExp {
 }
 interface EncloseExp {
  /**
   * Build converts an expression into a SQL fragment.
   */
  build(db: DB, params: Params): string
 }
 /**
  * TableModel is the interface that should be implemented by models which have unconventional table names.
  */
 interface TableModel {
  [key:string]: any;
  tableName(): string
 }
 /**
  * ModelQuery represents a query associated with a struct model.
  */
 interface ModelQuery {
 }
 interface newModelQuery {
  (model: {
   }, fieldMapFunc: FieldMapFunc, db: DB, builder: Builder): (ModelQuery)
 }
 interface ModelQuery {
  /**
   * Context returns the context associated with the query.
   */
  context(): context.Context
 }
 interface ModelQuery {
  /**
   * WithContext associates a context with the query.
   */
  withContext(ctx: context.Context): (ModelQuery)
 }
 interface ModelQuery {
  /**
   * Exclude excludes the specified struct fields from being inserted/updated into the DB table.
   */
  exclude(...attrs: string[]): (ModelQuery)
 }
 interface ModelQuery {
  /**
   * Insert inserts a row in the table using the struct model associated with this query.
   * 
   * By default, it inserts *all* public fields into the table, including those nil or empty ones.
   * You may pass a list of the fields to this method to indicate that only those fields should be inserted.
   * You may also call Exclude to exclude some fields from being inserted.
   * 
   * If a model has an empty primary key, it is considered auto-incremental and the corresponding struct
   * field will be filled with the generated primary key value after a successful insertion.
   */
  insert(...attrs: string[]): void
 }
 interface ModelQuery {
  /**
   * Update updates a row in the table using the struct model associated with this query.
   * The row being updated has the same primary key as specified by the model.
   * 
   * By default, it updates *all* public fields in the table, including those nil or empty ones.
   * You may pass a list of the fields to this method to indicate that only those fields should be updated.
   * You may also call Exclude to exclude some fields from being updated.
   */
  update(...attrs: string[]): void
 }
 interface ModelQuery {
  /**
   * Delete deletes a row in the table using the primary key specified by the struct model associated with this query.
   */
  delete(): void
 }
 /**
  * ExecHookFunc executes before op allowing custom handling like auto fail/retry.
  */
 interface ExecHookFunc {(q: Query, op: () => void): void }
 /**
  * OneHookFunc executes right before the query populate the row result from One() call (aka. op).
  */
 interface OneHookFunc {(q: Query, a: {
  }, op: (b: {
  }) => void): void }
 /**
  * AllHookFunc executes right before the query populate the row result from All() call (aka. op).
  */
 interface AllHookFunc {(q: Query, sliceA: {
  }, op: (sliceB: {
  }) => void): void }
 /**
  * Params represents a list of parameter values to be bound to a SQL statement.
  * The map keys are the parameter names while the map values are the corresponding parameter values.
  */
 interface Params extends _TygojaDict{}
 /**
  * Executor prepares, executes, or queries a SQL statement.
  */
 interface Executor {
  [key:string]: any;
  /**
   * Exec executes a SQL statement
   */
  exec(query: string, ...args: {
  }[]): sql.Result
  /**
   * ExecContext executes a SQL statement with the given context
   */
  execContext(ctx: context.Context, query: string, ...args: {
  }[]): sql.Result
  /**
   * Query queries a SQL statement
   */
  query(query: string, ...args: {
  }[]): (sql.Rows)
  /**
   * QueryContext queries a SQL statement with the given context
   */
  queryContext(ctx: context.Context, query: string, ...args: {
  }[]): (sql.Rows)
  /**
   * Prepare creates a prepared statement
   */
  prepare(query: string): (sql.Stmt)
 }
 /**
  * Query represents a SQL statement to be executed.
  */
 interface Query {
  /**
   * FieldMapper maps struct field names to DB column names.
   */
  fieldMapper: FieldMapFunc
  /**
   * LastError contains the last error (if any) of the query.
   * LastError is cleared by Execute(), Row(), Rows(), One(), and All().
   */
  lastError: Error
  /**
   * LogFunc is used to log the SQL statement being executed.
   */
  logFunc: LogFunc
  /**
   * PerfFunc is used to log the SQL execution time. It is ignored if nil.
   * Deprecated: Please use QueryLogFunc and ExecLogFunc instead.
   */
  perfFunc: PerfFunc
  /**
   * QueryLogFunc is called each time when performing a SQL query that returns data.
   */
  queryLogFunc: QueryLogFunc
  /**
   * ExecLogFunc is called each time when a SQL statement is executed.
   */
  execLogFunc: ExecLogFunc
 }
 interface newQuery {
  /**
   * NewQuery creates a new Query with the given SQL statement.
   */
  (db: DB, executor: Executor, sql: string): (Query)
 }
 interface Query {
  /**
   * SQL returns the original SQL used to create the query.
   * The actual SQL (RawSQL) being executed is obtained by replacing the named
   * parameter placeholders with anonymous ones.
   */
  sql(): string
 }
 interface Query {
  /**
   * Context returns the context associated with the query.
   */
  context(): context.Context
 }
 interface Query {
  /**
   * WithContext associates a context with the query.
   */
  withContext(ctx: context.Context): (Query)
 }
 interface Query {
  /**
   * WithExecHook associates the provided exec hook function with the query.
   * 
   * It is called for every Query resolver (Execute(), One(), All(), Row(), Column()),
   * allowing you to implement auto fail/retry or any other additional handling.
   */
  withExecHook(fn: ExecHookFunc): (Query)
 }
 interface Query {
  /**
   * WithOneHook associates the provided hook function with the query,
   * called on q.One(), allowing you to implement custom struct scan based
   * on the One() argument and/or result.
   */
  withOneHook(fn: OneHookFunc): (Query)
 }
 interface Query {
  /**
   * WithOneHook associates the provided hook function with the query,
   * called on q.All(), allowing you to implement custom slice scan based
   * on the All() argument and/or result.
   */
  withAllHook(fn: AllHookFunc): (Query)
 }
 interface Query {
  /**
   * Params returns the parameters to be bound to the SQL statement represented by this query.
   */
  params(): Params
 }
 interface Query {
  /**
   * Prepare creates a prepared statement for later queries or executions.
   * Close() should be called after finishing all queries.
   */
  prepare(): (Query)
 }
 interface Query {
  /**
   * Close closes the underlying prepared statement.
   * Close does nothing if the query has not been prepared before.
   */
  close(): void
 }
 interface Query {
  /**
   * Bind sets the parameters that should be bound to the SQL statement.
   * The parameter placeholders in the SQL statement are in the format of "{:ParamName}".
   */
  bind(params: Params): (Query)
 }
 interface Query {
  /**
   * Execute executes the SQL statement without retrieving data.
   */
  execute(): sql.Result
 }
 interface Query {
  /**
   * One executes the SQL statement and populates the first row of the result into a struct or NullStringMap.
   * Refer to Rows.ScanStruct() and Rows.ScanMap() for more details on how to specify
   * the variable to be populated.
   * Note that when the query has no rows in the result set, an sql.ErrNoRows will be returned.
   */
  one(a: {
   }): void
 }
 interface Query {
  /**
   * All executes the SQL statement and populates all the resulting rows into a slice of struct or NullStringMap.
   * The slice must be given as a pointer. Each slice element must be either a struct or a NullStringMap.
   * Refer to Rows.ScanStruct() and Rows.ScanMap() for more details on how each slice element can be.
   * If the query returns no row, the slice will be an empty slice (not nil).
   */
  all(slice: {
   }): void
 }
 interface Query {
  /**
   * Row executes the SQL statement and populates the first row of the result into a list of variables.
   * Note that the number of the variables should match to that of the columns in the query result.
   * Note that when the query has no rows in the result set, an sql.ErrNoRows will be returned.
   */
  row(...a: {
   }[]): void
 }
 interface Query {
  /**
   * Column executes the SQL statement and populates the first column of the result into a slice.
   * Note that the parameter must be a pointer to a slice.
   */
  column(a: {
   }): void
 }
 interface Query {
  /**
   * Rows executes the SQL statement and returns a Rows object to allow retrieving data row by row.
   */
  rows(): (Rows)
 }
 /**
  * QueryBuilder builds different clauses for a SELECT SQL statement.
  */
 interface QueryBuilder {
  [key:string]: any;
  /**
   * BuildSelect generates a SELECT clause from the given selected column names.
   */
  buildSelect(cols: Array<string>, distinct: boolean, option: string): string
  /**
   * BuildFrom generates a FROM clause from the given tables.
   */
  buildFrom(tables: Array<string>): string
  /**
   * BuildGroupBy generates a GROUP BY clause from the given group-by columns.
   */
  buildGroupBy(cols: Array<string>): string
  /**
   * BuildJoin generates a JOIN clause from the given join information.
   */
  buildJoin(_arg0: Array<JoinInfo>, _arg1: Params): string
  /**
   * BuildWhere generates a WHERE clause from the given expression.
   */
  buildWhere(_arg0: Expression, _arg1: Params): string
  /**
   * BuildHaving generates a HAVING clause from the given expression.
   */
  buildHaving(_arg0: Expression, _arg1: Params): string
  /**
   * BuildOrderByAndLimit generates the ORDER BY and LIMIT clauses.
   */
  buildOrderByAndLimit(_arg0: string, _arg1: Array<string>, _arg2: number, _arg3: number): string
  /**
   * BuildUnion generates a UNION clause from the given union information.
   */
  buildUnion(_arg0: Array<UnionInfo>, _arg1: Params): string
 }
 /**
  * BaseQueryBuilder provides a basic implementation of QueryBuilder.
  */
 interface BaseQueryBuilder {
 }
 interface newBaseQueryBuilder {
  /**
   * NewBaseQueryBuilder creates a new BaseQueryBuilder instance.
   */
  (db: DB): (BaseQueryBuilder)
 }
 interface BaseQueryBuilder {
  /**
   * DB returns the DB instance associated with the query builder.
   */
  db(): (DB)
 }
 interface BaseQueryBuilder {
  /**
   * BuildSelect generates a SELECT clause from the given selected column names.
   */
  buildSelect(cols: Array<string>, distinct: boolean, option: string): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildFrom generates a FROM clause from the given tables.
   */
  buildFrom(tables: Array<string>): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildJoin generates a JOIN clause from the given join information.
   */
  buildJoin(joins: Array<JoinInfo>, params: Params): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildWhere generates a WHERE clause from the given expression.
   */
  buildWhere(e: Expression, params: Params): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildHaving generates a HAVING clause from the given expression.
   */
  buildHaving(e: Expression, params: Params): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildGroupBy generates a GROUP BY clause from the given group-by columns.
   */
  buildGroupBy(cols: Array<string>): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildOrderByAndLimit generates the ORDER BY and LIMIT clauses.
   */
  buildOrderByAndLimit(sql: string, cols: Array<string>, limit: number, offset: number): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildUnion generates a UNION clause from the given union information.
   */
  buildUnion(unions: Array<UnionInfo>, params: Params): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildOrderBy generates the ORDER BY clause.
   */
  buildOrderBy(cols: Array<string>): string
 }
 interface BaseQueryBuilder {
  /**
   * BuildLimit generates the LIMIT clause.
   */
  buildLimit(limit: number, offset: number): string
 }
 /**
  * VarTypeError indicates a variable type error when trying to populating a variable with DB result.
  */
 interface VarTypeError extends String{}
 interface VarTypeError {
  /**
   * Error returns the error message.
   */
  error(): string
 }
 /**
  * NullStringMap is a map of sql.NullString that can be used to hold DB query result.
  * The map keys correspond to the DB column names, while the map values are their corresponding column values.
  */
 interface NullStringMap extends _TygojaDict{}
 /**
  * Rows enhances sql.Rows by providing additional data query methods.
  * Rows can be obtained by calling Query.Rows(). It is mainly used to populate data row by row.
  */
 type _sLdXoMB = sql.Rows
 interface Rows extends _sLdXoMB {
 }
 interface Rows {
  /**
   * ScanMap populates the current row of data into a NullStringMap.
   * Note that the NullStringMap must not be nil, or it will panic.
   * The NullStringMap will be populated using column names as keys and their values as
   * the corresponding element values.
   */
  scanMap(a: NullStringMap): void
 }
 interface Rows {
  /**
   * ScanStruct populates the current row of data into a struct.
   * The struct must be given as a pointer.
   * 
   * ScanStruct associates struct fields with DB table columns through a field mapping function.
   * It populates a struct field with the data of its associated column.
   * Note that only exported struct fields will be populated.
   * 
   * By default, DefaultFieldMapFunc() is used to map struct fields to table columns.
   * This function separates each word in a field name with a underscore and turns every letter into lower case.
   * For example, "LastName" is mapped to "last_name", "MyID" is mapped to "my_id", and so on.
   * To change the default behavior, set DB.FieldMapper with your custom mapping function.
   * You may also set Query.FieldMapper to change the behavior for particular queries.
   */
  scanStruct(a: {
   }): void
 }
 /**
  * BuildHookFunc defines a callback function that is executed on Query creation.
  */
 interface BuildHookFunc {(q: Query): void }
 /**
  * SelectQuery represents a DB-agnostic SELECT query.
  * It can be built into a DB-specific query by calling the Build() method.
  */
 interface SelectQuery {
  /**
   * FieldMapper maps struct field names to DB column names.
   */
  fieldMapper: FieldMapFunc
  /**
   * TableMapper maps structs to DB table names.
   */
  tableMapper: TableMapFunc
 }
 /**
  * JoinInfo contains the specification for a JOIN clause.
  */
 interface JoinInfo {
  join: string
  table: string
  on: Expression
 }
 /**
  * UnionInfo contains the specification for a UNION clause.
  */
 interface UnionInfo {
  all: boolean
  query?: Query
 }
 interface newSelectQuery {
  /**
   * NewSelectQuery creates a new SelectQuery instance.
   */
  (builder: Builder, db: DB): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * WithBuildHook runs the provided hook function with the query created on Build().
   */
  withBuildHook(fn: BuildHookFunc): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Context returns the context associated with the query.
   */
  context(): context.Context
 }
 interface SelectQuery {
  /**
   * WithContext associates a context with the query.
   */
  withContext(ctx: context.Context): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * PreFragment sets SQL fragment that should be prepended before the select query (e.g. WITH clause).
   */
  preFragment(fragment: string): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * PostFragment sets SQL fragment that should be appended at the end of the select query.
   */
  postFragment(fragment: string): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Select specifies the columns to be selected.
   * Column names will be automatically quoted.
   */
  select(...cols: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * AndSelect adds additional columns to be selected.
   * Column names will be automatically quoted.
   */
  andSelect(...cols: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Distinct specifies whether to select columns distinctively.
   * By default, distinct is false.
   */
  distinct(v: boolean): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * SelectOption specifies additional option that should be append to "SELECT".
   */
  selectOption(option: string): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * From specifies which tables to select from.
   * Table names will be automatically quoted.
   */
  from(...tables: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Where specifies the WHERE condition.
   */
  where(e: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * AndWhere concatenates a new WHERE condition with the existing one (if any) using "AND".
   */
  andWhere(e: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * OrWhere concatenates a new WHERE condition with the existing one (if any) using "OR".
   */
  orWhere(e: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Join specifies a JOIN clause.
   * The "typ" parameter specifies the JOIN type (e.g. "INNER JOIN", "LEFT JOIN").
   */
  join(typ: string, table: string, on: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * InnerJoin specifies an INNER JOIN clause.
   * This is a shortcut method for Join.
   */
  innerJoin(table: string, on: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * LeftJoin specifies a LEFT JOIN clause.
   * This is a shortcut method for Join.
   */
  leftJoin(table: string, on: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * RightJoin specifies a RIGHT JOIN clause.
   * This is a shortcut method for Join.
   */
  rightJoin(table: string, on: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * OrderBy specifies the ORDER BY clause.
   * Column names will be properly quoted. A column name can contain "ASC" or "DESC" to indicate its ordering direction.
   */
  orderBy(...cols: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * AndOrderBy appends additional columns to the existing ORDER BY clause.
   * Column names will be properly quoted. A column name can contain "ASC" or "DESC" to indicate its ordering direction.
   */
  andOrderBy(...cols: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * GroupBy specifies the GROUP BY clause.
   * Column names will be properly quoted.
   */
  groupBy(...cols: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * AndGroupBy appends additional columns to the existing GROUP BY clause.
   * Column names will be properly quoted.
   */
  andGroupBy(...cols: string[]): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Having specifies the HAVING clause.
   */
  having(e: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * AndHaving concatenates a new HAVING condition with the existing one (if any) using "AND".
   */
  andHaving(e: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * OrHaving concatenates a new HAVING condition with the existing one (if any) using "OR".
   */
  orHaving(e: Expression): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Union specifies a UNION clause.
   */
  union(q: Query): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * UnionAll specifies a UNION ALL clause.
   */
  unionAll(q: Query): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Limit specifies the LIMIT clause.
   * A negative limit means no limit.
   */
  limit(limit: number): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Offset specifies the OFFSET clause.
   * A negative offset means no offset.
   */
  offset(offset: number): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Bind specifies the parameter values to be bound to the query.
   */
  bind(params: Params): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * AndBind appends additional parameters to be bound to the query.
   */
  andBind(params: Params): (SelectQuery)
 }
 interface SelectQuery {
  /**
   * Build builds the SELECT query and returns an executable Query object.
   */
  build(): (Query)
 }
 interface SelectQuery {
  /**
   * One executes the SELECT query and populates the first row of the result into the specified variable.
   * 
   * If the query does not specify a "from" clause, the method will try to infer the name of the table
   * to be selected from by calling getTableName() which will return either the variable type name
   * or the TableName() method if the variable implements the TableModel interface.
   * 
   * Note that when the query has no rows in the result set, an sql.ErrNoRows will be returned.
   */
  one(a: {
   }): void
 }
 interface SelectQuery {
  /**
   * Model selects the row with the specified primary key and populates the model with the row data.
   * 
   * The model variable should be a pointer to a struct. If the query does not specify a "from" clause,
   * it will use the model struct to determine which table to select data from. It will also use the model
   * to infer the name of the primary key column. Only simple primary key is supported. For composite primary keys,
   * please use Where() to specify the filtering condition.
   */
  model(pk: {
   }, model: {
   }): void
 }
 interface SelectQuery {
  /**
   * All executes the SELECT query and populates all rows of the result into a slice.
   * 
   * Note that the slice must be passed in as a pointer.
   * 
   * If the query does not specify a "from" clause, the method will try to infer the name of the table
   * to be selected from by calling getTableName() which will return either the type name of the slice elements
   * or the TableName() method if the slice element implements the TableModel interface.
   */
  all(slice: {
   }): void
 }
 interface SelectQuery {
  /**
   * Rows builds and executes the SELECT query and returns a Rows object for data retrieval purpose.
   * This is a shortcut to SelectQuery.Build().Rows()
   */
  rows(): (Rows)
 }
 interface SelectQuery {
  /**
   * Row builds and executes the SELECT query and populates the first row of the result into the specified variables.
   * This is a shortcut to SelectQuery.Build().Row()
   */
  row(...a: {
   }[]): void
 }
 interface SelectQuery {
  /**
   * Column builds and executes the SELECT statement and populates the first column of the result into a slice.
   * Note that the parameter must be a pointer to a slice.
   * This is a shortcut to SelectQuery.Build().Column()
   */
  column(a: {
   }): void
 }
 /**
  * QueryInfo represents a debug/info struct with exported SelectQuery fields.
  */
 interface QueryInfo {
  preFragment: string
  postFragment: string
  builder: Builder
  selects: Array<string>
  distinct: boolean
  selectOption: string
  from: Array<string>
  where: Expression
  join: Array<JoinInfo>
  orderBy: Array<string>
  groupBy: Array<string>
  having: Expression
  union: Array<UnionInfo>
  limit: number
  offset: number
  params: Params
  context: context.Context
  buildHook: BuildHookFunc
 }
 interface SelectQuery {
  /**
   * Info exports common SelectQuery fields allowing to inspect the
   * current select query options.
   */
  info(): (QueryInfo)
 }
 /**
  * FieldMapFunc converts a struct field name into a DB column name.
  */
 interface FieldMapFunc {(_arg0: string): string }
 /**
  * TableMapFunc converts a sample struct into a DB table name.
  */
 interface TableMapFunc {(a: {
  }): string }
 interface structInfo {
 }
 type _sFCCfIW = structInfo
 interface structValue extends _sFCCfIW {
 }
 interface fieldInfo {
 }
 interface structInfoMapKey {
 }
 /**
  * PostScanner is an optional interface used by ScanStruct.
  */
 interface PostScanner {
  [key:string]: any;
  /**
   * PostScan executes right after the struct has been populated
   * with the DB values, allowing you to further normalize or validate
   * the loaded data.
   */
  postScan(): void
 }
 interface defaultFieldMapFunc {
  /**
   * DefaultFieldMapFunc maps a field name to a DB column name.
   * The mapping rule set by this method is that words in a field name will be separated by underscores
   * and the name will be turned into lower case. For example, "FirstName" maps to "first_name", and "MyID" becomes "my_id".
   * See DB.FieldMapper for more details.
   */
  (f: string): string
 }
 interface getTableName {
  /**
   * GetTableName implements the default way of determining the table name corresponding to the given model struct
   * or slice of structs. To get the actual table name for a model, you should use DB.TableMapFunc() instead.
   * Do not call this method in a model's TableName() method because it will cause infinite loop.
   */
  (a: {
   }): string
 }
 /**
  * Tx enhances sql.Tx with additional querying methods.
  */
 type _szqtPqt = Builder
 interface Tx extends _szqtPqt {
 }
 interface Tx {
  /**
   * Commit commits the transaction.
   */
  commit(): void
 }
 interface Tx {
  /**
   * Rollback aborts the transaction.
   */
  rollback(): void
 }
}

namespace security {
 interface s256Challenge {
  /**
   * S256Challenge creates base64 encoded sha256 challenge string derived from code.
   * The padding of the result base64 string is stripped per [RFC 7636].
   * 
   * [RFC 7636]: https://datatracker.ietf.org/doc/html/rfc7636#section-4.2
   */
  (code: string): string
 }
 interface md5 {
  /**
   * MD5 creates md5 hash from the provided plain text.
   */
  (text: string): string
 }
 interface sha256 {
  /**
   * SHA256 creates sha256 hash as defined in FIPS 180-4 from the provided text.
   */
  (text: string): string
 }
 interface sha512 {
  /**
   * SHA512 creates sha512 hash as defined in FIPS 180-4 from the provided text.
   */
  (text: string): string
 }
 interface hs256 {
  /**
   * HS256 creates a HMAC hash with sha256 digest algorithm.
   */
  (text: string, secret: string): string
 }
 interface hs512 {
  /**
   * HS512 creates a HMAC hash with sha512 digest algorithm.
   */
  (text: string, secret: string): string
 }
 interface equal {
  /**
   * Equal compares two hash strings for equality without leaking timing information.
   */
  (hash1: string, hash2: string): boolean
 }
 // @ts-ignore
 import crand = rand
 interface encrypt {
  /**
   * Encrypt encrypts "data" with the specified "key" (must be valid 32 char AES key).
   * 
   * This method uses AES-256-GCM block cypher mode.
   */
  (data: string|Array<number>, key: string): string
 }
 interface decrypt {
  /**
   * Decrypt decrypts encrypted text with key (must be valid 32 chars AES key).
   * 
   * This method uses AES-256-GCM block cypher mode.
   */
  (cipherText: string, key: string): string|Array<number>
 }
 interface parseUnverifiedJWT {
  /**
   * ParseUnverifiedJWT parses JWT and returns its claims
   * but DOES NOT verify the signature.
   * 
   * It verifies only the exp, iat and nbf claims.
   */
  (token: string): jwt.MapClaims
 }
 interface parseJWT {
  /**
   * ParseJWT verifies and parses JWT and returns its claims.
   */
  (token: string, verificationKey: string): jwt.MapClaims
 }
 interface newJWT {
  /**
   * NewJWT generates and returns new HS256 signed JWT.
   */
  (payload: jwt.MapClaims, signingKey: string, duration: time.Duration): string
 }
 // @ts-ignore
 import cryptoRand = rand
 // @ts-ignore
 import mathRand = rand
 interface randomString {
  /**
   * RandomString generates a cryptographically random string with the specified length.
   * 
   * The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
   */
  (length: number): string
 }
 interface randomStringWithAlphabet {
  /**
   * RandomStringWithAlphabet generates a cryptographically random string
   * with the specified length and characters set.
   * 
   * It panics if for some reason rand.Int returns a non-nil error.
   */
  (length: number, alphabet: string): string
 }
 interface pseudorandomString {
  /**
   * PseudorandomString generates a pseudorandom string with the specified length.
   * 
   * The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
   * 
   * For a cryptographically random string (but a little bit slower) use RandomString instead.
   */
  (length: number): string
 }
 interface pseudorandomStringWithAlphabet {
  /**
   * PseudorandomStringWithAlphabet generates a pseudorandom string
   * with the specified length and characters set.
   * 
   * For a cryptographically random (but a little bit slower) use RandomStringWithAlphabet instead.
   */
  (length: number, alphabet: string): string
 }
 interface randomStringByRegex {
  /**
   * RandomStringByRegex generates a random string matching the regex pattern.
   * If optFlags is not set, fallbacks to [syntax.Perl].
   * 
   * NB! While the source of the randomness comes from [crypto/rand] this method
   * is not recommended to be used on its own in critical secure contexts because
   * the generated length could vary too much on the used pattern and may not be
   * as secure as simply calling [security.RandomString].
   * If you still insist on using it for such purposes, consider at least
   * a large enough minimum length for the generated string, e.g. `[a-z0-9]{30}`.
   * 
   * This function is inspired by github.com/pipe01/revregexp, github.com/lucasjones/reggen and other similar packages.
   */
  (pattern: string, ...optFlags: syntax.Flags[]): string
 }
}

namespace filesystem {
 /**
  * FileReader defines an interface for a file resource reader.
  */
 interface FileReader {
  [key:string]: any;
  open(): io.ReadSeekCloser
 }
 /**
  * File defines a single file [io.ReadSeekCloser] resource.
  * 
  * The file could be from a local path, multipart/form-data header, etc.
  */
 interface File {
  reader: FileReader
  name: string
  originalName: string
  size: number
 }
 interface File {
  /**
   * AsMap implements [core.mapExtractor] and returns a value suitable
   * to be used in an API rule expression.
   */
  asMap(): _TygojaDict
 }
 interface newFileFromPath {
  /**
   * NewFileFromPath creates a new File instance from the provided local file path.
   */
  (path: string): (File)
 }
 interface newFileFromBytes {
  /**
   * NewFileFromBytes creates a new File instance from the provided byte slice.
   */
  (b: string|Array<number>, name: string): (File)
 }
 interface newFileFromMultipart {
  /**
   * NewFileFromMultipart creates a new File from the provided multipart header.
   */
  (mh: multipart.FileHeader): (File)
 }
 interface newFileFromURL {
  /**
   * NewFileFromURL creates a new File from the provided url by
   * downloading the resource and load it as BytesReader.
   * 
   * Example
   * 
   * ```
   * 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   * 	defer cancel()
   * 
   * 	file, err := filesystem.NewFileFromURL(ctx, "https://example.com/image.png")
   * ```
   */
  (ctx: context.Context, url: string): (File)
 }
 /**
  * MultipartReader defines a FileReader from [multipart.FileHeader].
  */
 interface MultipartReader {
  header?: multipart.FileHeader
 }
 interface MultipartReader {
  /**
   * Open implements the [filesystem.FileReader] interface.
   */
  open(): io.ReadSeekCloser
 }
 /**
  * PathReader defines a FileReader from a local file path.
  */
 interface PathReader {
  path: string
 }
 interface PathReader {
  /**
   * Open implements the [filesystem.FileReader] interface.
   */
  open(): io.ReadSeekCloser
 }
 /**
  * BytesReader defines a FileReader from bytes content.
  */
 interface BytesReader {
  bytes: string|Array<number>
 }
 interface BytesReader {
  /**
   * Open implements the [filesystem.FileReader] interface.
   */
  open(): io.ReadSeekCloser
 }
 type _sPUPdtB = bytes.Reader
 interface bytesReadSeekCloser extends _sPUPdtB {
 }
 interface bytesReadSeekCloser {
  /**
   * Close implements the [io.ReadSeekCloser] interface.
   */
  close(): void
 }
 /**
  * openFuncAsReader defines a FileReader from a bare Open function.
  */
 interface openFuncAsReader {(): io.ReadSeekCloser }
 interface openFuncAsReader {
  /**
   * Open implements the [filesystem.FileReader] interface.
   */
  open(): io.ReadSeekCloser
 }
 interface System {
 }
 interface newS3 {
  /**
   * NewS3 initializes an S3 filesystem instance.
   * 
   * NB! Make sure to call `Close()` after you are done working with it.
   */
  (bucketName: string, region: string, endpoint: string, accessKey: string, secretKey: string, s3ForcePathStyle: boolean): (System)
 }
 interface newLocal {
  /**
   * NewLocal initializes a new local filesystem instance.
   * 
   * NB! Make sure to call `Close()` after you are done working with it.
   */
  (dirPath: string): (System)
 }
 interface System {
  /**
   * SetContext assigns the specified context to the current filesystem.
   */
  setContext(ctx: context.Context): void
 }
 interface System {
  /**
   * Close releases any resources used for the related filesystem.
   */
  close(): void
 }
 interface System {
  /**
   * Exists checks if file with fileKey path exists or not.
   */
  exists(fileKey: string): boolean
 }
 interface System {
  /**
   * Attributes returns the attributes for the file with fileKey path.
   * 
   * If the file doesn't exist it returns ErrNotFound.
   */
  attributes(fileKey: string): (blob.Attributes)
 }
 interface System {
  /**
   * GetReader returns a file content reader for the given fileKey.
   * 
   * NB! Make sure to call Close() on the file after you are done working with it.
   * 
   * If the file doesn't exist returns ErrNotFound.
   */
  getReader(fileKey: string): (blob.Reader)
 }
 interface System {
  /**
   * Deprecated: Please use GetReader(fileKey) instead.
   */
  getFile(fileKey: string): (blob.Reader)
 }
 interface System {
  /**
   * GetReuploadableFile constructs a new reuploadable File value
   * from the associated fileKey blob.Reader.
   * 
   * If preserveName is false then the returned File.Name will have
   * a new randomly generated suffix, otherwise it will reuse the original one.
   * 
   * This method could be useful in case you want to clone an existing
   * Record file and assign it to a new Record (e.g. in a Record duplicate action).
   * 
   * If you simply want to copy an existing file to a new location you
   * could check the Copy(srcKey, dstKey) method.
   */
  getReuploadableFile(fileKey: string, preserveName: boolean): (File)
 }
 interface System {
  /**
   * Copy copies the file stored at srcKey to dstKey.
   * 
   * If srcKey file doesn't exist, it returns ErrNotFound.
   * 
   * If dstKey file already exists, it is overwritten.
   */
  copy(srcKey: string, dstKey: string): void
 }
 interface System {
  /**
   * List returns a flat list with info for all files under the specified prefix.
   */
  list(prefix: string): Array<(blob.ListObject | undefined)>
 }
 interface System {
  /**
   * Upload writes content into the fileKey location.
   */
  upload(content: string|Array<number>, fileKey: string): void
 }
 interface System {
  /**
   * UploadFile uploads the provided File to the fileKey location.
   */
  uploadFile(file: File, fileKey: string): void
 }
 interface System {
  /**
   * UploadMultipart uploads the provided multipart file to the fileKey location.
   */
  uploadMultipart(fh: multipart.FileHeader, fileKey: string): void
 }
 interface System {
  /**
   * Delete deletes stored file at fileKey location.
   * 
   * If the file doesn't exist returns ErrNotFound.
   */
  delete(fileKey: string): void
 }
 interface System {
  /**
   * DeletePrefix deletes everything starting with the specified prefix.
   * 
   * The prefix could be subpath (ex. "/a/b/") or filename prefix (ex. "/a/b/file_").
   */
  deletePrefix(prefix: string): Array<Error>
 }
 interface System {
  /**
   * Checks if the provided dir prefix doesn't have any files.
   * 
   * A trailing slash will be appended to a non-empty dir string argument
   * to ensure that the checked prefix is a "directory".
   * 
   * Returns "false" in case the has at least one file, otherwise - "true".
   */
  isEmptyDir(dir: string): boolean
 }
 interface System {
  /**
   * Serve serves the file at fileKey location to an HTTP response.
   * 
   * If the `download` query parameter is used the file will be always served for
   * download no matter of its type (aka. with "Content-Disposition: attachment").
   * 
   * Internally this method uses [http.ServeContent] so Range requests,
   * If-Match, If-Unmodified-Since, etc. headers are handled transparently.
   */
  serve(res: http.ResponseWriter, req: http.Request, fileKey: string, name: string): void
 }
 interface System {
  /**
   * CreateThumb creates a new thumb image for the file at originalKey location.
   * The new thumb file is stored at thumbKey location.
   * 
   * thumbSize is in the format:
   * - 0xH  (eg. 0x100)    - resize to H height preserving the aspect ratio
   * - Wx0  (eg. 300x0)    - resize to W width preserving the aspect ratio
   * - WxH  (eg. 300x100)  - resize and crop to WxH viewbox (from center)
   * - WxHt (eg. 300x100t) - resize and crop to WxH viewbox (from top)
   * - WxHb (eg. 300x100b) - resize and crop to WxH viewbox (from bottom)
   * - WxHf (eg. 300x100f) - fit inside a WxH viewbox (without cropping)
   */
  createThumb(originalKey: string, thumbKey: string, thumbSize: string): void
 }
}

/**
 * Package exec runs external commands. It wraps os.StartProcess to make it
 * easier to remap stdin and stdout, connect I/O with pipes, and do other
 * adjustments.
 * 
 * Unlike the "system" library call from C and other languages, the
 * os/exec package intentionally does not invoke the system shell and
 * does not expand any glob patterns or handle other expansions,
 * pipelines, or redirections typically done by shells. The package
 * behaves more like C's "exec" family of functions. To expand glob
 * patterns, either call the shell directly, taking care to escape any
 * dangerous input, or use the [path/filepath] package's Glob function.
 * To expand environment variables, use package os's ExpandEnv.
 * 
 * Note that the examples in this package assume a Unix system.
 * They may not run on Windows, and they do not run in the Go Playground
 * used by golang.org and godoc.org.
 * 
 * # Executables in the current directory
 * 
 * The functions [Command] and [LookPath] look for a program
 * in the directories listed in the current path, following the
 * conventions of the host operating system.
 * Operating systems have for decades included the current
 * directory in this search, sometimes implicitly and sometimes
 * configured explicitly that way by default.
 * Modern practice is that including the current directory
 * is usually unexpected and often leads to security problems.
 * 
 * To avoid those security problems, as of Go 1.19, this package will not resolve a program
 * using an implicit or explicit path entry relative to the current directory.
 * That is, if you run [LookPath]("go"), it will not successfully return
 * ./go on Unix nor .\go.exe on Windows, no matter how the path is configured.
 * Instead, if the usual path algorithms would result in that answer,
 * these functions return an error err satisfying [errors.Is](err, [ErrDot]).
 * 
 * For example, consider these two program snippets:
 * 
 * ```
 * 	path, err := exec.LookPath("prog")
 * 	if err != nil {
 * 		log.Fatal(err)
 * 	}
 * 	use(path)
 * ```
 * 
 * and
 * 
 * ```
 * 	cmd := exec.Command("prog")
 * 	if err := cmd.Run(); err != nil {
 * 		log.Fatal(err)
 * 	}
 * ```
 * 
 * These will not find and run ./prog or .\prog.exe,
 * no matter how the current path is configured.
 * 
 * Code that always wants to run a program from the current directory
 * can be rewritten to say "./prog" instead of "prog".
 * 
 * Code that insists on including results from relative path entries
 * can instead override the error using an errors.Is check:
 * 
 * ```
 * 	path, err := exec.LookPath("prog")
 * 	if errors.Is(err, exec.ErrDot) {
 * 		err = nil
 * 	}
 * 	if err != nil {
 * 		log.Fatal(err)
 * 	}
 * 	use(path)
 * ```
 * 
 * and
 * 
 * ```
 * 	cmd := exec.Command("prog")
 * 	if errors.Is(cmd.Err, exec.ErrDot) {
 * 		cmd.Err = nil
 * 	}
 * 	if err := cmd.Run(); err != nil {
 * 		log.Fatal(err)
 * 	}
 * ```
 * 
 * Setting the environment variable GODEBUG=execerrdot=0
 * disables generation of ErrDot entirely, temporarily restoring the pre-Go 1.19
 * behavior for programs that are unable to apply more targeted fixes.
 * A future version of Go may remove support for this variable.
 * 
 * Before adding such overrides, make sure you understand the
 * security implications of doing so.
 * See https://go.dev/blog/path-security for more information.
 */
namespace exec {
 interface command {
  /**
   * Command returns the [Cmd] struct to execute the named program with
   * the given arguments.
   * 
   * It sets only the Path and Args in the returned structure.
   * 
   * If name contains no path separators, Command uses [LookPath] to
   * resolve name to a complete path if possible. Otherwise it uses name
   * directly as Path.
   * 
   * The returned Cmd's Args field is constructed from the command name
   * followed by the elements of arg, so arg should not include the
   * command name itself. For example, Command("echo", "hello").
   * Args[0] is always name, not the possibly resolved Path.
   * 
   * On Windows, processes receive the whole command line as a single string
   * and do their own parsing. Command combines and quotes Args into a command
   * line string with an algorithm compatible with applications using
   * CommandLineToArgvW (which is the most common way). Notable exceptions are
   * msiexec.exe and cmd.exe (and thus, all batch files), which have a different
   * unquoting algorithm. In these or other similar cases, you can do the
   * quoting yourself and provide the full command line in SysProcAttr.CmdLine,
   * leaving Args empty.
   */
  (name: string, ...arg: string[]): (Cmd)
 }
}

/**
 * Package core is the backbone of PocketBase.
 * 
 * It defines the main PocketBase App interface and its base implementation.
 */
namespace core {
 /**
  * App defines the main PocketBase app interface.
  * 
  * Note that the interface is not intended to be implemented manually by users
  * and instead they should use core.BaseApp (either directly or as embedded field in a custom struct).
  * 
  * This interface exists to make testing easier and to allow users to
  * create common and pluggable helpers and methods that doesn't rely
  * on a specific wrapped app struct (hence the large interface size).
  */
 interface App {
  [key:string]: any;
  /**
   * UnsafeWithoutHooks returns a shallow copy of the current app WITHOUT any registered hooks.
   * 
   * NB! Note that using the returned app instance may cause data integrity errors
   * since the Record validations and data normalizations (including files uploads)
   * rely on the app hooks to work.
   */
  unsafeWithoutHooks(): App
  /**
   * Logger returns the default app logger.
   * 
   * If the application is not bootstrapped yet, fallbacks to slog.Default().
   */
  logger(): (slog.Logger)
  /**
   * IsBootstrapped checks if the application was initialized
   * (aka. whether Bootstrap() was called).
   */
  isBootstrapped(): boolean
  /**
   * IsTransactional checks if the current app instance is part of a transaction.
   */
  isTransactional(): boolean
  /**
   * TxInfo returns the transaction associated with the current app instance (if any).
   * 
   * Could be used if you want to execute indirectly a function after
   * the related app transaction completes using `app.TxInfo().OnAfterFunc(callback)`.
   */
  txInfo(): (TxAppInfo)
  /**
   * Bootstrap initializes the application
   * (aka. create data dir, open db connections, load settings, etc.).
   * 
   * It will call ResetBootstrapState() if the application was already bootstrapped.
   */
  bootstrap(): void
  /**
   * ResetBootstrapState releases the initialized core app resources
   * (closing db connections, stopping cron ticker, etc.).
   */
  resetBootstrapState(): void
  /**
   * DataDir returns the app data directory path.
   */
  dataDir(): string
  /**
   * EncryptionEnv returns the name of the app secret env key
   * (currently used primarily for optional settings encryption but this may change in the future).
   */
  encryptionEnv(): string
  /**
   * IsDev returns whether the app is in dev mode.
   * 
   * When enabled logs, executed sql statements, etc. are printed to the stderr.
   */
  isDev(): boolean
  /**
   * Settings returns the loaded app settings.
   */
  settings(): (Settings)
  /**
   * Store returns the app runtime store.
   */
  store(): (store.Store<string, any>)
  /**
   * Cron returns the app cron instance.
   */
  cron(): (cron.Cron)
  /**
   * SubscriptionsBroker returns the app realtime subscriptions broker instance.
   */
  subscriptionsBroker(): (subscriptions.Broker)
  /**
   * NewMailClient creates and returns a new SMTP or Sendmail client
   * based on the current app settings.
   */
  newMailClient(): mailer.Mailer
  /**
   * NewFilesystem creates a new local or S3 filesystem instance
   * for managing regular app files (ex. record uploads)
   * based on the current app settings.
   * 
   * NB! Make sure to call Close() on the returned result
   * after you are done working with it.
   */
  newFilesystem(): (filesystem.System)
  /**
   * NewBackupsFilesystem creates a new local or S3 filesystem instance
   * for managing app backups based on the current app settings.
   * 
   * NB! Make sure to call Close() on the returned result
   * after you are done working with it.
   */
  newBackupsFilesystem(): (filesystem.System)
  /**
   * ReloadSettings reinitializes and reloads the stored application settings.
   */
  reloadSettings(): void
  /**
   * CreateBackup creates a new backup of the current app pb_data directory.
   * 
   * Backups can be stored on S3 if it is configured in app.Settings().Backups.
   * 
   * Please refer to the godoc of the specific CoreApp implementation
   * for details on the backup procedures.
   */
  createBackup(ctx: context.Context, name: string): void
  /**
   * RestoreBackup restores the backup with the specified name and restarts
   * the current running application process.
   * 
   * The safely perform the restore it is recommended to have free disk space
   * for at least 2x the size of the restored pb_data backup.
   * 
   * Please refer to the godoc of the specific CoreApp implementation
   * for details on the restore procedures.
   * 
   * NB! This feature is experimental and currently is expected to work only on UNIX based systems.
   */
  restoreBackup(ctx: context.Context, name: string): void
  /**
   * Restart restarts (aka. replaces) the current running application process.
   * 
   * NB! It relies on execve which is supported only on UNIX based systems.
   */
  restart(): void
  /**
   * RunSystemMigrations applies all new migrations registered in the [core.SystemMigrations] list.
   */
  runSystemMigrations(): void
  /**
   * RunAppMigrations applies all new migrations registered in the [CoreAppMigrations] list.
   */
  runAppMigrations(): void
  /**
   * RunAllMigrations applies all system and app migrations
   * (aka. from both [core.SystemMigrations] and [CoreAppMigrations]).
   */
  runAllMigrations(): void
  /**
   * DB returns the default app data.db builder instance.
   * 
   * To minimize SQLITE_BUSY errors, it automatically routes the
   * SELECT queries to the underlying concurrent db pool and everything else
   * to the nonconcurrent one.
   * 
   * For more finer control over the used connections pools you can
   * call directly ConcurrentDB() or NonconcurrentDB().
   */
  db(): dbx.Builder
  /**
   * ConcurrentDB returns the concurrent app data.db builder instance.
   * 
   * This method is used mainly internally for executing db read
   * operations in a concurrent/non-blocking manner.
   * 
   * Most users should use simply DB() as it will automatically
   * route the query execution to ConcurrentDB() or NonconcurrentDB().
   * 
   * In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
   */
  concurrentDB(): dbx.Builder
  /**
   * NonconcurrentDB returns the nonconcurrent app data.db builder instance.
   * 
   * The returned db instance is limited only to a single open connection,
   * meaning that it can process only 1 db operation at a time (other queries queue up).
   * 
   * This method is used mainly internally and in the tests to execute write
   * (save/delete) db operations as it helps with minimizing the SQLITE_BUSY errors.
   * 
   * Most users should use simply DB() as it will automatically
   * route the query execution to ConcurrentDB() or NonconcurrentDB().
   * 
   * In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
   */
  nonconcurrentDB(): dbx.Builder
  /**
   * AuxDB returns the app auxiliary.db builder instance.
   * 
   * To minimize SQLITE_BUSY errors, it automatically routes the
   * SELECT queries to the underlying concurrent db pool and everything else
   * to the nonconcurrent one.
   * 
   * For more finer control over the used connections pools you can
   * call directly AuxConcurrentDB() or AuxNonconcurrentDB().
   */
  auxDB(): dbx.Builder
  /**
   * AuxConcurrentDB returns the concurrent app auxiliary.db builder instance.
   * 
   * This method is used mainly internally for executing db read
   * operations in a concurrent/non-blocking manner.
   * 
   * Most users should use simply AuxDB() as it will automatically
   * route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
   * 
   * In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
   */
  auxConcurrentDB(): dbx.Builder
  /**
   * AuxNonconcurrentDB returns the nonconcurrent app auxiliary.db builder instance.
   * 
   * The returned db instance is limited only to a single open connection,
   * meaning that it can process only 1 db operation at a time (other queries queue up).
   * 
   * This method is used mainly internally and in the tests to execute write
   * (save/delete) db operations as it helps with minimizing the SQLITE_BUSY errors.
   * 
   * Most users should use simply AuxDB() as it will automatically
   * route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
   * 
   * In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
   */
  auxNonconcurrentDB(): dbx.Builder
  /**
   * HasTable checks if a table (or view) with the provided name exists (case insensitive).
   * in the data.db.
   */
  hasTable(tableName: string): boolean
  /**
   * AuxHasTable checks if a table (or view) with the provided name exists (case insensitive)
   * in the auxiliary.db.
   */
  auxHasTable(tableName: string): boolean
  /**
   * TableColumns returns all column names of a single table by its name.
   */
  tableColumns(tableName: string): Array<string>
  /**
   * TableInfo returns the "table_info" pragma result for the specified table.
   */
  tableInfo(tableName: string): Array<(TableInfoRow | undefined)>
  /**
   * TableIndexes returns a name grouped map with all non empty index of the specified table.
   * 
   * Note: This method doesn't return an error on nonexisting table.
   */
  tableIndexes(tableName: string): _TygojaDict
  /**
   * DeleteTable drops the specified table.
   * 
   * This method is a no-op if a table with the provided name doesn't exist.
   * 
   * NB! Be aware that this method is vulnerable to SQL injection and the
   * "tableName" argument must come only from trusted input!
   */
  deleteTable(tableName: string): void
  /**
   * DeleteView drops the specified view name.
   * 
   * This method is a no-op if a view with the provided name doesn't exist.
   * 
   * NB! Be aware that this method is vulnerable to SQL injection and the
   * "name" argument must come only from trusted input!
   */
  deleteView(name: string): void
  /**
   * SaveView creates (or updates already existing) persistent SQL view.
   * 
   * NB! Be aware that this method is vulnerable to SQL injection and the
   * "selectQuery" argument must come only from trusted input!
   */
  saveView(name: string, selectQuery: string): void
  /**
   * CreateViewFields creates a new FieldsList from the provided select query.
   * 
   * There are some caveats:
   * - The select query must have an "id" column.
   * - Wildcard ("*") columns are not supported to avoid accidentally leaking sensitive data.
   */
  createViewFields(selectQuery: string): FieldsList
  /**
   * FindRecordByViewFile returns the original Record of the provided view collection file.
   */
  findRecordByViewFile(viewCollectionModelOrIdentifier: any, fileFieldName: string, filename: string): (Record)
  /**
   * Vacuum executes VACUUM on the data.db in order to reclaim unused data db disk space.
   */
  vacuum(): void
  /**
   * AuxVacuum executes VACUUM on the auxiliary.db in order to reclaim unused auxiliary db disk space.
   */
  auxVacuum(): void
  /**
   * ModelQuery creates a new preconfigured select data.db query with preset
   * SELECT, FROM and other common fields based on the provided model.
   */
  modelQuery(model: Model): (dbx.SelectQuery)
  /**
   * AuxModelQuery creates a new preconfigured select auxiliary.db query with preset
   * SELECT, FROM and other common fields based on the provided model.
   */
  auxModelQuery(model: Model): (dbx.SelectQuery)
  /**
   * Delete deletes the specified model from the regular app database.
   */
  delete(model: Model): void
  /**
   * Delete deletes the specified model from the regular app database
   * (the context could be used to limit the query execution).
   */
  deleteWithContext(ctx: context.Context, model: Model): void
  /**
   * AuxDelete deletes the specified model from the auxiliary database.
   */
  auxDelete(model: Model): void
  /**
   * AuxDeleteWithContext deletes the specified model from the auxiliary database
   * (the context could be used to limit the query execution).
   */
  auxDeleteWithContext(ctx: context.Context, model: Model): void
  /**
   * Save validates and saves the specified model into the regular app database.
   * 
   * If you don't want to run validations, use [App.SaveNoValidate()].
   */
  save(model: Model): void
  /**
   * SaveWithContext is the same as [App.Save()] but allows specifying a context to limit the db execution.
   * 
   * If you don't want to run validations, use [App.SaveNoValidateWithContext()].
   */
  saveWithContext(ctx: context.Context, model: Model): void
  /**
   * SaveNoValidate saves the specified model into the regular app database without performing validations.
   * 
   * If you want to also run validations before persisting, use [App.Save()].
   */
  saveNoValidate(model: Model): void
  /**
   * SaveNoValidateWithContext is the same as [App.SaveNoValidate()]
   * but allows specifying a context to limit the db execution.
   * 
   * If you want to also run validations before persisting, use [App.SaveWithContext()].
   */
  saveNoValidateWithContext(ctx: context.Context, model: Model): void
  /**
   * AuxSave validates and saves the specified model into the auxiliary app database.
   * 
   * If you don't want to run validations, use [App.AuxSaveNoValidate()].
   */
  auxSave(model: Model): void
  /**
   * AuxSaveWithContext is the same as [App.AuxSave()] but allows specifying a context to limit the db execution.
   * 
   * If you don't want to run validations, use [App.AuxSaveNoValidateWithContext()].
   */
  auxSaveWithContext(ctx: context.Context, model: Model): void
  /**
   * AuxSaveNoValidate saves the specified model into the auxiliary app database without performing validations.
   * 
   * If you want to also run validations before persisting, use [App.AuxSave()].
   */
  auxSaveNoValidate(model: Model): void
  /**
   * AuxSaveNoValidateWithContext is the same as [App.AuxSaveNoValidate()]
   * but allows specifying a context to limit the db execution.
   * 
   * If you want to also run validations before persisting, use [App.AuxSaveWithContext()].
   */
  auxSaveNoValidateWithContext(ctx: context.Context, model: Model): void
  /**
   * Validate triggers the OnModelValidate hook for the specified model.
   */
  validate(model: Model): void
  /**
   * ValidateWithContext is the same as Validate but allows specifying the ModelEvent context.
   */
  validateWithContext(ctx: context.Context, model: Model): void
  /**
   * RunInTransaction wraps fn into a transaction for the regular app database.
   * 
   * It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
   */
  runInTransaction(fn: (txApp: App) => void): void
  /**
   * AuxRunInTransaction wraps fn into a transaction for the auxiliary app database.
   * 
   * It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
   */
  auxRunInTransaction(fn: (txApp: App) => void): void
  /**
   * LogQuery returns a new Log select query.
   */
  logQuery(): (dbx.SelectQuery)
  /**
   * FindLogById finds a single Log entry by its id.
   */
  findLogById(id: string): (Log)
  /**
   * LogsStatsItem returns hourly grouped logs statistics.
   */
  logsStats(expr: dbx.Expression): Array<(LogsStatsItem | undefined)>
  /**
   * DeleteOldLogs delete all logs that are created before createdBefore.
   */
  deleteOldLogs(createdBefore: time.Time): void
  /**
   * CollectionQuery returns a new Collection select query.
   */
  collectionQuery(): (dbx.SelectQuery)
  /**
   * FindCollections finds all collections by the given type(s).
   * 
   * If collectionTypes is not set, it returns all collections.
   * 
   * Example:
   * 
   * ```
   * 	app.FindAllCollections() // all collections
   * 	app.FindAllCollections("auth", "view") // only auth and view collections
   * ```
   */
  findAllCollections(...collectionTypes: string[]): Array<(Collection | undefined)>
  /**
   * ReloadCachedCollections fetches all collections and caches them into the app store.
   */
  reloadCachedCollections(): void
  /**
   * FindCollectionByNameOrId finds a single collection by its name (case insensitive) or id.s
   */
  findCollectionByNameOrId(nameOrId: string): (Collection)
  /**
   * FindCachedCollectionByNameOrId is similar to [App.FindCollectionByNameOrId]
   * but retrieves the Collection from the app cache instead of making a db call.
   * 
   * NB! This method is suitable for read-only Collection operations.
   * 
   * Returns [sql.ErrNoRows] if no Collection is found for consistency
   * with the [App.FindCollectionByNameOrId] method.
   * 
   * If you plan making changes to the returned Collection model,
   * use [App.FindCollectionByNameOrId] instead.
   * 
   * Caveats:
   * 
   * ```
   *   - The returned Collection should be used only for read-only operations.
   *     Avoid directly modifying the returned cached Collection as it will affect
   *     the global cached value even if you don't persist the changes in the database!
   *   - If you are updating a Collection in a transaction and then call this method before commit,
   *     it'll return the cached Collection state and not the one from the uncommitted transaction.
   *   - The cache is automatically updated on collections db change (create/update/delete).
   *     To manually reload the cache you can call [App.ReloadCachedCollections]
   * ```
   */
  findCachedCollectionByNameOrId(nameOrId: string): (Collection)
  /**
   * FindCollectionReferences returns information for all relation
   * fields referencing the provided collection.
   * 
   * If the provided collection has reference to itself then it will be
   * also included in the result. To exclude it, pass the collection id
   * as the excludeIds argument.
   */
  findCollectionReferences(collection: Collection, ...excludeIds: string[]): _TygojaDict
  /**
   * FindCachedCollectionReferences is similar to [App.FindCollectionReferences]
   * but retrieves the Collection from the app cache instead of making a db call.
   * 
   * NB! This method is suitable for read-only Collection operations.
   * 
   * If you plan making changes to the returned Collection model,
   * use [App.FindCollectionReferences] instead.
   * 
   * Caveats:
   * 
   * ```
   *   - The returned Collection should be used only for read-only operations.
   *     Avoid directly modifying the returned cached Collection as it will affect
   *     the global cached value even if you don't persist the changes in the database!
   *   - If you are updating a Collection in a transaction and then call this method before commit,
   *     it'll return the cached Collection state and not the one from the uncommitted transaction.
   *   - The cache is automatically updated on collections db change (create/update/delete).
   *     To manually reload the cache you can call [App.ReloadCachedCollections].
   * ```
   */
  findCachedCollectionReferences(collection: Collection, ...excludeIds: string[]): _TygojaDict
  /**
   * IsCollectionNameUnique checks that there is no existing collection
   * with the provided name (case insensitive!).
   * 
   * Note: case insensitive check because the name is used also as
   * table name for the records.
   */
  isCollectionNameUnique(name: string, ...excludeIds: string[]): boolean
  /**
   * TruncateCollection deletes all records associated with the provided collection.
   * 
   * The truncate operation is executed in a single transaction,
   * aka. either everything is deleted or none.
   * 
   * Note that this method will also trigger the records related
   * cascade and file delete actions.
   */
  truncateCollection(collection: Collection): void
  /**
   * ImportCollections imports the provided collections data in a single transaction.
   * 
   * For existing matching collections, the imported data is unmarshaled on top of the existing model.
   * 
   * NB! If deleteMissing is true, ALL NON-SYSTEM COLLECTIONS AND SCHEMA FIELDS,
   * that are not present in the imported configuration, WILL BE DELETED
   * (this includes their related records data).
   */
  importCollections(toImport: Array<_TygojaDict>, deleteMissing: boolean): void
  /**
   * ImportCollectionsByMarshaledJSON is the same as [ImportCollections]
   * but accept marshaled json array as import data (usually used for the autogenerated snapshots).
   */
  importCollectionsByMarshaledJSON(rawSliceOfMaps: string|Array<number>, deleteMissing: boolean): void
  /**
   * SyncRecordTableSchema compares the two provided collections
   * and applies the necessary related record table changes.
   * 
   * If oldCollection is null, then only newCollection is used to create the record table.
   * 
   * This method is automatically invoked as part of a collection create/update/delete operation.
   */
  syncRecordTableSchema(newCollection: Collection, oldCollection: Collection): void
  /**
   * FindAllExternalAuthsByRecord returns all ExternalAuth models
   * linked to the provided auth record.
   */
  findAllExternalAuthsByRecord(authRecord: Record): Array<(ExternalAuth | undefined)>
  /**
   * FindAllExternalAuthsByCollection returns all ExternalAuth models
   * linked to the provided auth collection.
   */
  findAllExternalAuthsByCollection(collection: Collection): Array<(ExternalAuth | undefined)>
  /**
   * FindFirstExternalAuthByExpr returns the first available (the most recent created)
   * ExternalAuth model that satisfies the non-nil expression.
   */
  findFirstExternalAuthByExpr(expr: dbx.Expression): (ExternalAuth)
  /**
   * FindAllMFAsByRecord returns all MFA models linked to the provided auth record.
   */
  findAllMFAsByRecord(authRecord: Record): Array<(MFA | undefined)>
  /**
   * FindAllMFAsByCollection returns all MFA models linked to the provided collection.
   */
  findAllMFAsByCollection(collection: Collection): Array<(MFA | undefined)>
  /**
   * FindMFAById returns a single MFA model by its id.
   */
  findMFAById(id: string): (MFA)
  /**
   * DeleteAllMFAsByRecord deletes all MFA models associated with the provided record.
   * 
   * Returns a combined error with the failed deletes.
   */
  deleteAllMFAsByRecord(authRecord: Record): void
  /**
   * DeleteExpiredMFAs deletes the expired MFAs for all auth collections.
   */
  deleteExpiredMFAs(): void
  /**
   * FindAllOTPsByRecord returns all OTP models linked to the provided auth record.
   */
  findAllOTPsByRecord(authRecord: Record): Array<(OTP | undefined)>
  /**
   * FindAllOTPsByCollection returns all OTP models linked to the provided collection.
   */
  findAllOTPsByCollection(collection: Collection): Array<(OTP | undefined)>
  /**
   * FindOTPById returns a single OTP model by its id.
   */
  findOTPById(id: string): (OTP)
  /**
   * DeleteAllOTPsByRecord deletes all OTP models associated with the provided record.
   * 
   * Returns a combined error with the failed deletes.
   */
  deleteAllOTPsByRecord(authRecord: Record): void
  /**
   * DeleteExpiredOTPs deletes the expired OTPs for all auth collections.
   */
  deleteExpiredOTPs(): void
  /**
   * FindAllAuthOriginsByRecord returns all AuthOrigin models linked to the provided auth record (in DESC order).
   */
  findAllAuthOriginsByRecord(authRecord: Record): Array<(AuthOrigin | undefined)>
  /**
   * FindAllAuthOriginsByCollection returns all AuthOrigin models linked to the provided collection (in DESC order).
   */
  findAllAuthOriginsByCollection(collection: Collection): Array<(AuthOrigin | undefined)>
  /**
   * FindAuthOriginById returns a single AuthOrigin model by its id.
   */
  findAuthOriginById(id: string): (AuthOrigin)
  /**
   * FindAuthOriginByRecordAndFingerprint returns a single AuthOrigin model
   * by its authRecord relation and fingerprint.
   */
  findAuthOriginByRecordAndFingerprint(authRecord: Record, fingerprint: string): (AuthOrigin)
  /**
   * DeleteAllAuthOriginsByRecord deletes all AuthOrigin models associated with the provided record.
   * 
   * Returns a combined error with the failed deletes.
   */
  deleteAllAuthOriginsByRecord(authRecord: Record): void
  /**
   * RecordQuery returns a new Record select query from a collection model, id or name.
   * 
   * In case a collection id or name is provided and that collection doesn't
   * actually exists, the generated query will be created with a cancelled context
   * and will fail once an executor (Row(), One(), All(), etc.) is called.
   */
  recordQuery(collectionModelOrIdentifier: any): (dbx.SelectQuery)
  /**
   * FindRecordById finds the Record model by its id.
   */
  findRecordById(collectionModelOrIdentifier: any, recordId: string, ...optFilters: ((q: dbx.SelectQuery) => void)[]): (Record)
  /**
   * FindRecordsByIds finds all records by the specified ids.
   * If no records are found, returns an empty slice.
   */
  findRecordsByIds(collectionModelOrIdentifier: any, recordIds: Array<string>, ...optFilters: ((q: dbx.SelectQuery) => void)[]): Array<(Record | undefined)>
  /**
   * FindAllRecords finds all records matching specified db expressions.
   * 
   * Returns all collection records if no expression is provided.
   * 
   * Returns an empty slice if no records are found.
   * 
   * Example:
   * 
   * ```
   * 	// no extra expressions
   * 	app.FindAllRecords("example")
   * 
   * 	// with extra expressions
   * 	expr1 := dbx.HashExp{"email": "test@example.com"}
   * 	expr2 := dbx.NewExp("LOWER(username) = {:username}", dbx.Params{"username": "test"})
   * 	app.FindAllRecords("example", expr1, expr2)
   * ```
   */
  findAllRecords(collectionModelOrIdentifier: any, ...exprs: dbx.Expression[]): Array<(Record | undefined)>
  /**
   * FindFirstRecordByData returns the first found record matching
   * the provided key-value pair.
   */
  findFirstRecordByData(collectionModelOrIdentifier: any, key: string, value: any): (Record)
  /**
   * FindRecordsByFilter returns limit number of records matching the
   * provided string filter.
   * 
   * NB! Use the last "params" argument to bind untrusted user variables!
   * 
   * The filter argument is optional and can be empty string to target
   * all available records.
   * 
   * The sort argument is optional and can be empty string OR the same format
   * used in the web APIs, ex. "-created,title".
   * 
   * If the limit argument is <= 0, no limit is applied to the query and
   * all matching records are returned.
   * 
   * Returns an empty slice if no records are found.
   * 
   * Example:
   * 
   * ```
   * 	app.FindRecordsByFilter(
   * 		"posts",
   * 		"title ~ {:title} && visible = {:visible}",
   * 		"-created",
   * 		10,
   * 		0,
   * 		dbx.Params{"title": "lorem ipsum", "visible": true}
   * 	)
   * ```
   */
  findRecordsByFilter(collectionModelOrIdentifier: any, filter: string, sort: string, limit: number, offset: number, ...params: dbx.Params[]): Array<(Record | undefined)>
  /**
   * FindFirstRecordByFilter returns the first available record matching the provided filter (if any).
   * 
   * NB! Use the last params argument to bind untrusted user variables!
   * 
   * Returns sql.ErrNoRows if no record is found.
   * 
   * Example:
   * 
   * ```
   * 	app.FindFirstRecordByFilter("posts", "")
   * 	app.FindFirstRecordByFilter("posts", "slug={:slug} && status='public'", dbx.Params{"slug": "test"})
   * ```
   */
  findFirstRecordByFilter(collectionModelOrIdentifier: any, filter: string, ...params: dbx.Params[]): (Record)
  /**
   * CountRecords returns the total number of records in a collection.
   */
  countRecords(collectionModelOrIdentifier: any, ...exprs: dbx.Expression[]): number
  /**
   * FindAuthRecordByToken finds the auth record associated with the provided JWT
   * (auth, file, verifyEmail, changeEmail, passwordReset types).
   * 
   * Optionally specify a list of validTypes to check tokens only from those types.
   * 
   * Returns an error if the JWT is invalid, expired or not associated to an auth collection record.
   */
  findAuthRecordByToken(token: string, ...validTypes: string[]): (Record)
  /**
   * FindAuthRecordByEmail finds the auth record associated with the provided email.
   * 
   * Returns an error if it is not an auth collection or the record is not found.
   */
  findAuthRecordByEmail(collectionModelOrIdentifier: any, email: string): (Record)
  /**
   * CanAccessRecord checks if a record is allowed to be accessed by the
   * specified requestInfo and accessRule.
   * 
   * Rule and db checks are ignored in case requestInfo.Auth is a superuser.
   * 
   * The returned error indicate that something unexpected happened during
   * the check (eg. invalid rule or db query error).
   * 
   * The method always return false on invalid rule or db query error.
   * 
   * Example:
   * 
   * ```
   * 	requestInfo, _ := e.RequestInfo()
   * 	record, _ := app.FindRecordById("example", "RECORD_ID")
   * 	rule := types.Pointer("@request.auth.id != '' || status = 'public'")
   * 	// ... or use one of the record collection's rule, eg. record.Collection().ViewRule
   * 
   * 	if ok, _ := app.CanAccessRecord(record, requestInfo, rule); ok { ... }
   * ```
   */
  canAccessRecord(record: Record, requestInfo: RequestInfo, accessRule: string): boolean
  /**
   * ExpandRecord expands the relations of a single Record model.
   * 
   * If optFetchFunc is not set, then a default function will be used
   * that returns all relation records.
   * 
   * Returns a map with the failed expand parameters and their errors.
   */
  expandRecord(record: Record, expands: Array<string>, optFetchFunc: ExpandFetchFunc): _TygojaDict
  /**
   * ExpandRecords expands the relations of the provided Record models list.
   * 
   * If optFetchFunc is not set, then a default function will be used
   * that returns all relation records.
   * 
   * Returns a map with the failed expand parameters and their errors.
   */
  expandRecords(records: Array<(Record | undefined)>, expands: Array<string>, optFetchFunc: ExpandFetchFunc): _TygojaDict
  /**
   * OnBootstrap hook is triggered when initializing the main application
   * resources (db, app settings, etc).
   */
  onBootstrap(): (hook.Hook<BootstrapEvent | undefined>)
  /**
   * OnServe hook is triggered when the app web server is started
   * (after starting the TCP listener but before initializing the blocking serve task),
   * allowing you to adjust its options and attach new routes or middlewares.
   */
  onServe(): (hook.Hook<ServeEvent | undefined>)
  /**
   * OnTerminate hook is triggered when the app is in the process
   * of being terminated (ex. on SIGTERM signal).
   * 
   * Note that the app could be terminated abruptly without awaiting the hook completion.
   */
  onTerminate(): (hook.Hook<TerminateEvent | undefined>)
  /**
   * OnBackupCreate hook is triggered on each [App.CreateBackup] call.
   */
  onBackupCreate(): (hook.Hook<BackupEvent | undefined>)
  /**
   * OnBackupRestore hook is triggered before app backup restore (aka. [App.RestoreBackup] call).
   * 
   * Note that by default on success the application is restarted and the after state of the hook is ignored.
   */
  onBackupRestore(): (hook.Hook<BackupEvent | undefined>)
  /**
   * OnModelValidate is triggered every time when a model is being validated
   * (e.g. triggered by App.Validate() or App.Save()).
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelValidate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelCreate is triggered every time when a new model is being created
   * (e.g. triggered by App.Save()).
   * 
   * Operations BEFORE the e.Next() execute before the model validation
   * and the INSERT DB statement.
   * 
   * Operations AFTER the e.Next() execute after the model validation
   * and the INSERT DB statement.
   * 
   * Note that successful execution doesn't guarantee that the model
   * is persisted in the database since its wrapping transaction may
   * not have been committed yet.
   * If you want to listen to only the actual persisted events, you can
   * bind to [OnModelAfterCreateSuccess] or [OnModelAfterCreateError] hooks.
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelCreate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelCreateExecute is triggered after successful Model validation
   * and right before the model INSERT DB statement execution.
   * 
   * Usually it is triggered as part of the App.Save() in the following firing order:
   * OnModelCreate {
   * ```
   *    -> OnModelValidate (skipped with App.SaveNoValidate())
   *    -> OnModelCreateExecute
   * ```
   * }
   * 
   * Note that successful execution doesn't guarantee that the model
   * is persisted in the database since its wrapping transaction may have been
   * committed yet.
   * If you want to listen to only the actual persisted events,
   * you can bind to [OnModelAfterCreateSuccess] or [OnModelAfterCreateError] hooks.
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelCreateExecute(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelAfterCreateSuccess is triggered after each successful
   * Model DB create persistence.
   * 
   * Note that when a Model is persisted as part of a transaction,
   * this hook is delayed and executed only AFTER the transaction has been committed.
   * This hook is NOT triggered in case the transaction rollbacks
   * (aka. when the model wasn't persisted).
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelAfterCreateSuccess(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelAfterCreateError is triggered after each failed
   * Model DB create persistence.
   * 
   * Note that the execution of this hook is either immediate or delayed
   * depending on the error:
   * ```
   *   - "immediate" on App.Save() failure
   *   - "delayed" on transaction rollback
   * ```
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelAfterCreateError(...tags: string[]): (hook.TaggedHook<ModelErrorEvent | undefined>)
  /**
   * OnModelUpdate is triggered every time when a new model is being updated
   * (e.g. triggered by App.Save()).
   * 
   * Operations BEFORE the e.Next() execute before the model validation
   * and the UPDATE DB statement.
   * 
   * Operations AFTER the e.Next() execute after the model validation
   * and the UPDATE DB statement.
   * 
   * Note that successful execution doesn't guarantee that the model
   * is persisted in the database since its wrapping transaction may
   * not have been committed yet.
   * If you want to listen to only the actual persisted events, you can
   * bind to [OnModelAfterUpdateSuccess] or [OnModelAfterUpdateError] hooks.
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelUpdate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelUpdateExecute is triggered after successful Model validation
   * and right before the model UPDATE DB statement execution.
   * 
   * Usually it is triggered as part of the App.Save() in the following firing order:
   * OnModelUpdate {
   * ```
   *    -> OnModelValidate (skipped with App.SaveNoValidate())
   *    -> OnModelUpdateExecute
   * ```
   * }
   * 
   * Note that successful execution doesn't guarantee that the model
   * is persisted in the database since its wrapping transaction may have been
   * committed yet.
   * If you want to listen to only the actual persisted events,
   * you can bind to [OnModelAfterUpdateSuccess] or [OnModelAfterUpdateError] hooks.
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelUpdateExecute(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelAfterUpdateSuccess is triggered after each successful
   * Model DB update persistence.
   * 
   * Note that when a Model is persisted as part of a transaction,
   * this hook is delayed and executed only AFTER the transaction has been committed.
   * This hook is NOT triggered in case the transaction rollbacks
   * (aka. when the model changes weren't persisted).
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelAfterUpdateSuccess(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelAfterUpdateError is triggered after each failed
   * Model DB update persistence.
   * 
   * Note that the execution of this hook is either immediate or delayed
   * depending on the error:
   * ```
   *   - "immediate" on App.Save() failure
   *   - "delayed" on transaction rollback
   * ```
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelAfterUpdateError(...tags: string[]): (hook.TaggedHook<ModelErrorEvent | undefined>)
  /**
   * OnModelDelete is triggered every time when a new model is being deleted
   * (e.g. triggered by App.Delete()).
   * 
   * Note that successful execution doesn't guarantee that the model
   * is deleted from the database since its wrapping transaction may
   * not have been committed yet.
   * If you want to listen to only the actual persisted deleted events, you can
   * bind to [OnModelAfterDeleteSuccess] or [OnModelAfterDeleteError] hooks.
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelDelete(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelUpdateExecute is triggered right before the model
   * DELETE DB statement execution.
   * 
   * Usually it is triggered as part of the App.Delete() in the following firing order:
   * OnModelDelete {
   * ```
   *    -> (internal delete checks)
   *    -> OnModelDeleteExecute
   * ```
   * }
   * 
   * Note that successful execution doesn't guarantee that the model
   * is deleted from the database since its wrapping transaction may
   * not have been committed yet.
   * If you want to listen to only the actual persisted deleted events, you can
   * bind to [OnModelAfterDeleteSuccess] or [OnModelAfterDeleteError] hooks.
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelDeleteExecute(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelAfterDeleteSuccess is triggered after each successful
   * Model DB delete persistence.
   * 
   * Note that when a Model is deleted as part of a transaction,
   * this hook is delayed and executed only AFTER the transaction has been committed.
   * This hook is NOT triggered in case the transaction rollbacks
   * (aka. when the model delete wasn't persisted).
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelAfterDeleteSuccess(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
  /**
   * OnModelAfterDeleteError is triggered after each failed
   * Model DB delete persistence.
   * 
   * Note that the execution of this hook is either immediate or delayed
   * depending on the error:
   * ```
   *   - "immediate" on App.Delete() failure
   *   - "delayed" on transaction rollback
   * ```
   * 
   * For convenience, if you want to listen to only the Record models
   * events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
   * 
   * If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onModelAfterDeleteError(...tags: string[]): (hook.TaggedHook<ModelErrorEvent | undefined>)
  /**
   * OnRecordEnrich is triggered every time when a record is enriched
   * (as part of the builtin Record responses, during realtime message seriazation, or when [apis.EnrichRecord] is invoked).
   * 
   * It could be used for example to redact/hide or add computed temporary
   * Record model props only for the specific request info. For example:
   * 
   *  app.OnRecordEnrich("posts").BindFunc(func(e core.*RecordEnrichEvent) {
   * ```
   *      // hide one or more fields
   *      e.Record.Hide("role")
   * 
   *      // add new custom field for registered users
   *      if e.RequestInfo.Auth != nil && e.RequestInfo.Auth.Collection().Name == "users" {
   *          e.Record.WithCustomData(true) // for security requires explicitly allowing it
   *          e.Record.Set("computedScore", e.Record.GetInt("score") * e.RequestInfo.Auth.GetInt("baseScore"))
   *      }
   * 
   *      return e.Next()
   * ```
   *  })
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordEnrich(...tags: string[]): (hook.TaggedHook<RecordEnrichEvent | undefined>)
  /**
   * OnRecordValidate is a Record proxy model hook of [OnModelValidate].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordValidate(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordCreate is a Record proxy model hook of [OnModelCreate].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordCreate(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordCreateExecute is a Record proxy model hook of [OnModelCreateExecute].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordCreateExecute(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordAfterCreateSuccess is a Record proxy model hook of [OnModelAfterCreateSuccess].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterCreateSuccess(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordAfterCreateError is a Record proxy model hook of [OnModelAfterCreateError].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterCreateError(...tags: string[]): (hook.TaggedHook<RecordErrorEvent | undefined>)
  /**
   * OnRecordUpdate is a Record proxy model hook of [OnModelUpdate].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordUpdate(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordUpdateExecute is a Record proxy model hook of [OnModelUpdateExecute].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordUpdateExecute(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordAfterUpdateSuccess is a Record proxy model hook of [OnModelAfterUpdateSuccess].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterUpdateSuccess(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordAfterUpdateError is a Record proxy model hook of [OnModelAfterUpdateError].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterUpdateError(...tags: string[]): (hook.TaggedHook<RecordErrorEvent | undefined>)
  /**
   * OnRecordDelete is a Record proxy model hook of [OnModelDelete].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordDelete(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordDeleteExecute is a Record proxy model hook of [OnModelDeleteExecute].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordDeleteExecute(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordAfterDeleteSuccess is a Record proxy model hook of [OnModelAfterDeleteSuccess].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterDeleteSuccess(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
  /**
   * OnRecordAfterDeleteError is a Record proxy model hook of [OnModelAfterDeleteError].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterDeleteError(...tags: string[]): (hook.TaggedHook<RecordErrorEvent | undefined>)
  /**
   * OnCollectionValidate is a Collection proxy model hook of [OnModelValidate].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionValidate(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionCreate is a Collection proxy model hook of [OnModelCreate].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionCreate(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionCreateExecute is a Collection proxy model hook of [OnModelCreateExecute].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionCreateExecute(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionAfterCreateSuccess is a Collection proxy model hook of [OnModelAfterCreateSuccess].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionAfterCreateSuccess(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionAfterCreateError is a Collection proxy model hook of [OnModelAfterCreateError].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionAfterCreateError(...tags: string[]): (hook.TaggedHook<CollectionErrorEvent | undefined>)
  /**
   * OnCollectionUpdate is a Collection proxy model hook of [OnModelUpdate].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionUpdate(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionUpdateExecute is a Collection proxy model hook of [OnModelUpdateExecute].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionUpdateExecute(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionAfterUpdateSuccess is a Collection proxy model hook of [OnModelAfterUpdateSuccess].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionAfterUpdateSuccess(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionAfterUpdateError is a Collection proxy model hook of [OnModelAfterUpdateError].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionAfterUpdateError(...tags: string[]): (hook.TaggedHook<CollectionErrorEvent | undefined>)
  /**
   * OnCollectionDelete is a Collection proxy model hook of [OnModelDelete].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionDelete(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionDeleteExecute is a Collection proxy model hook of [OnModelDeleteExecute].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionDeleteExecute(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionAfterDeleteSuccess is a Collection proxy model hook of [OnModelAfterDeleteSuccess].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionAfterDeleteSuccess(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
  /**
   * OnCollectionAfterDeleteError is a Collection proxy model hook of [OnModelAfterDeleteError].
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onCollectionAfterDeleteError(...tags: string[]): (hook.TaggedHook<CollectionErrorEvent | undefined>)
  /**
   * OnMailerSend hook is triggered every time when a new email is
   * being sent using the [App.NewMailClient()] instance.
   * 
   * It allows intercepting the email message or to use a custom mailer client.
   */
  onMailerSend(): (hook.Hook<MailerEvent | undefined>)
  /**
   * OnMailerRecordAuthAlertSend hook is triggered when
   * sending a new device login auth alert email, allowing you to
   * intercept and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerRecordAuthAlertSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
  /**
   * OnMailerBeforeRecordResetPasswordSend hook is triggered when
   * sending a password reset email to an auth record, allowing
   * you to intercept and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerRecordPasswordResetSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
  /**
   * OnMailerBeforeRecordVerificationSend hook is triggered when
   * sending a verification email to an auth record, allowing
   * you to intercept and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerRecordVerificationSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
  /**
   * OnMailerRecordEmailChangeSend hook is triggered when sending a
   * confirmation new address email to an auth record, allowing
   * you to intercept and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerRecordEmailChangeSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
  /**
   * OnMailerRecordOTPSend hook is triggered when sending an OTP email
   * to an auth record, allowing you to intercept and customize the
   * email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerRecordOTPSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
  /**
   * OnRealtimeConnectRequest hook is triggered when establishing the SSE client connection.
   * 
   * Any execution after e.Next() of a hook handler happens after the client disconnects.
   */
  onRealtimeConnectRequest(): (hook.Hook<RealtimeConnectRequestEvent | undefined>)
  /**
   * OnRealtimeMessageSend hook is triggered when sending an SSE message to a client.
   */
  onRealtimeMessageSend(): (hook.Hook<RealtimeMessageEvent | undefined>)
  /**
   * OnRealtimeSubscribeRequest hook is triggered when updating the
   * client subscriptions, allowing you to further validate and
   * modify the submitted change.
   */
  onRealtimeSubscribeRequest(): (hook.Hook<RealtimeSubscribeRequestEvent | undefined>)
  /**
   * OnSettingsListRequest hook is triggered on each API Settings list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onSettingsListRequest(): (hook.Hook<SettingsListRequestEvent | undefined>)
  /**
   * OnSettingsUpdateRequest hook is triggered on each API Settings update request.
   * 
   * Could be used to additionally validate the request data or
   * implement completely different persistence behavior.
   */
  onSettingsUpdateRequest(): (hook.Hook<SettingsUpdateRequestEvent | undefined>)
  /**
   * OnSettingsReload hook is triggered every time when the App.Settings()
   * is being replaced with a new state.
   * 
   * Calling App.Settings() after e.Next() returns the new state.
   */
  onSettingsReload(): (hook.Hook<SettingsReloadEvent | undefined>)
  /**
   * OnFileDownloadRequest hook is triggered before each API File download request.
   * 
   * Could be used to validate or modify the file response before
   * returning it to the client.
   */
  onFileDownloadRequest(...tags: string[]): (hook.TaggedHook<FileDownloadRequestEvent | undefined>)
  /**
   * OnFileBeforeTokenRequest hook is triggered on each auth file token API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onFileTokenRequest(...tags: string[]): (hook.TaggedHook<FileTokenRequestEvent | undefined>)
  /**
   * OnRecordAuthRequest hook is triggered on each successful API
   * record authentication request (sign-in, token refresh, etc.).
   * 
   * Could be used to additionally validate or modify the authenticated
   * record data and token.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAuthRequest(...tags: string[]): (hook.TaggedHook<RecordAuthRequestEvent | undefined>)
  /**
   * OnRecordAuthWithPasswordRequest hook is triggered on each
   * Record auth with password API request.
   * 
   * [RecordAuthWithPasswordRequestEvent.Record] could be nil if no matching identity is found, allowing
   * you to manually locate a different Record model (by reassigning [RecordAuthWithPasswordRequestEvent.Record]).
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAuthWithPasswordRequest(...tags: string[]): (hook.TaggedHook<RecordAuthWithPasswordRequestEvent | undefined>)
  /**
   * OnRecordAuthWithOAuth2Request hook is triggered on each Record
   * OAuth2 sign-in/sign-up API request (after token exchange and before external provider linking).
   * 
   * If [RecordAuthWithOAuth2RequestEvent.Record] is not set, then the OAuth2
   * request will try to create a new auth Record.
   * 
   * To assign or link a different existing record model you can
   * change the [RecordAuthWithOAuth2RequestEvent.Record] field.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAuthWithOAuth2Request(...tags: string[]): (hook.TaggedHook<RecordAuthWithOAuth2RequestEvent | undefined>)
  /**
   * OnRecordAuthRefreshRequest hook is triggered on each Record
   * auth refresh API request (right before generating a new auth token).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different auth refresh behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAuthRefreshRequest(...tags: string[]): (hook.TaggedHook<RecordAuthRefreshRequestEvent | undefined>)
  /**
   * OnRecordRequestPasswordResetRequest hook is triggered on
   * each Record request password reset API request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different password reset behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordRequestPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordRequestPasswordResetRequestEvent | undefined>)
  /**
   * OnRecordConfirmPasswordResetRequest hook is triggered on
   * each Record confirm password reset API request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordConfirmPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmPasswordResetRequestEvent | undefined>)
  /**
   * OnRecordRequestVerificationRequest hook is triggered on
   * each Record request verification API request.
   * 
   * Could be used to additionally validate the loaded request data or implement
   * completely different verification behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordRequestVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordRequestVerificationRequestEvent | undefined>)
  /**
   * OnRecordConfirmVerificationRequest hook is triggered on each
   * Record confirm verification API request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordConfirmVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmVerificationRequestEvent | undefined>)
  /**
   * OnRecordRequestEmailChangeRequest hook is triggered on each
   * Record request email change API request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different request email change behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordRequestEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEmailChangeRequestEvent | undefined>)
  /**
   * OnRecordConfirmEmailChangeRequest hook is triggered on each
   * Record confirm email change API request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordConfirmEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmEmailChangeRequestEvent | undefined>)
  /**
   * OnRecordRequestOTPRequest hook is triggered on each Record
   * request OTP API request.
   * 
   * [RecordCreateOTPRequestEvent.Record] could be nil if no matching identity is found, allowing
   * you to manually create or locate a different Record model (by reassigning [RecordCreateOTPRequestEvent.Record]).
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordRequestOTPRequest(...tags: string[]): (hook.TaggedHook<RecordCreateOTPRequestEvent | undefined>)
  /**
   * OnRecordAuthWithOTPRequest hook is triggered on each Record
   * auth with OTP API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAuthWithOTPRequest(...tags: string[]): (hook.TaggedHook<RecordAuthWithOTPRequestEvent | undefined>)
  /**
   * OnRecordsListRequest hook is triggered on each API Records list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordsListRequest(...tags: string[]): (hook.TaggedHook<RecordsListRequestEvent | undefined>)
  /**
   * OnRecordViewRequest hook is triggered on each API Record view request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordViewRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
  /**
   * OnRecordCreateRequest hook is triggered on each API Record create request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordCreateRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
  /**
   * OnRecordUpdateRequest hook is triggered on each API Record update request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordUpdateRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
  /**
   * OnRecordDeleteRequest hook is triggered on each API Record delete request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different delete behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordDeleteRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
  /**
   * OnCollectionsListRequest hook is triggered on each API Collections list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onCollectionsListRequest(): (hook.Hook<CollectionsListRequestEvent | undefined>)
  /**
   * OnCollectionViewRequest hook is triggered on each API Collection view request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onCollectionViewRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
  /**
   * OnCollectionCreateRequest hook is triggered on each API Collection create request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onCollectionCreateRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
  /**
   * OnCollectionUpdateRequest hook is triggered on each API Collection update request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onCollectionUpdateRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
  /**
   * OnCollectionDeleteRequest hook is triggered on each API Collection delete request.
   * 
   * Could be used to additionally validate the request data or implement
   * completely different delete behavior.
   */
  onCollectionDeleteRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
  /**
   * OnCollectionsBeforeImportRequest hook is triggered on each API
   * collections import request.
   * 
   * Could be used to additionally validate the imported collections or
   * to implement completely different import behavior.
   */
  onCollectionsImportRequest(): (hook.Hook<CollectionsImportRequestEvent | undefined>)
  /**
   * OnBatchRequest hook is triggered on each API batch request.
   * 
   * Could be used to additionally validate or modify the submitted batch requests.
   */
  onBatchRequest(): (hook.Hook<BatchRequestEvent | undefined>)
 }
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * AuthOrigin defines a Record proxy for working with the authOrigins collection.
  */
 type _sdlZlAF = Record
 interface AuthOrigin extends _sdlZlAF {
 }
 interface newAuthOrigin {
  /**
   * NewAuthOrigin instantiates and returns a new blank *AuthOrigin model.
   * 
   * Example usage:
   * 
   * ```
   * 	origin := core.NewOrigin(app)
   * 	origin.SetRecordRef(user.Id)
   * 	origin.SetCollectionRef(user.Collection().Id)
   * 	origin.SetFingerprint("...")
   * 	app.Save(origin)
   * ```
   */
  (app: App): (AuthOrigin)
 }
 interface AuthOrigin {
  /**
   * PreValidate implements the [PreValidator] interface and checks
   * whether the proxy is properly loaded.
   */
  preValidate(ctx: context.Context, app: App): void
 }
 interface AuthOrigin {
  /**
   * ProxyRecord returns the proxied Record model.
   */
  proxyRecord(): (Record)
 }
 interface AuthOrigin {
  /**
   * SetProxyRecord loads the specified record model into the current proxy.
   */
  setProxyRecord(record: Record): void
 }
 interface AuthOrigin {
  /**
   * CollectionRef returns the "collectionRef" field value.
   */
  collectionRef(): string
 }
 interface AuthOrigin {
  /**
   * SetCollectionRef updates the "collectionRef" record field value.
   */
  setCollectionRef(collectionId: string): void
 }
 interface AuthOrigin {
  /**
   * RecordRef returns the "recordRef" record field value.
   */
  recordRef(): string
 }
 interface AuthOrigin {
  /**
   * SetRecordRef updates the "recordRef" record field value.
   */
  setRecordRef(recordId: string): void
 }
 interface AuthOrigin {
  /**
   * Fingerprint returns the "fingerprint" record field value.
   */
  fingerprint(): string
 }
 interface AuthOrigin {
  /**
   * SetFingerprint updates the "fingerprint" record field value.
   */
  setFingerprint(fingerprint: string): void
 }
 interface AuthOrigin {
  /**
   * Created returns the "created" record field value.
   */
  created(): types.DateTime
 }
 interface AuthOrigin {
  /**
   * Updated returns the "updated" record field value.
   */
  updated(): types.DateTime
 }
 interface BaseApp {
  /**
   * FindAllAuthOriginsByRecord returns all AuthOrigin models linked to the provided auth record (in DESC order).
   */
  findAllAuthOriginsByRecord(authRecord: Record): Array<(AuthOrigin | undefined)>
 }
 interface BaseApp {
  /**
   * FindAllAuthOriginsByCollection returns all AuthOrigin models linked to the provided collection (in DESC order).
   */
  findAllAuthOriginsByCollection(collection: Collection): Array<(AuthOrigin | undefined)>
 }
 interface BaseApp {
  /**
   * FindAuthOriginById returns a single AuthOrigin model by its id.
   */
  findAuthOriginById(id: string): (AuthOrigin)
 }
 interface BaseApp {
  /**
   * FindAuthOriginByRecordAndFingerprint returns a single AuthOrigin model
   * by its authRecord relation and fingerprint.
   */
  findAuthOriginByRecordAndFingerprint(authRecord: Record, fingerprint: string): (AuthOrigin)
 }
 interface BaseApp {
  /**
   * DeleteAllAuthOriginsByRecord deletes all AuthOrigin models associated with the provided record.
   * 
   * Returns a combined error with the failed deletes.
   */
  deleteAllAuthOriginsByRecord(authRecord: Record): void
 }
 /**
  * FilesManager defines an interface with common methods that files manager models should implement.
  */
 interface FilesManager {
  [key:string]: any;
  /**
   * BaseFilesPath returns the storage dir path used by the interface instance.
   */
  baseFilesPath(): string
 }
 /**
  * DBConnectFunc defines a database connection initialization function.
  */
 interface DBConnectFunc {(dbPath: string): (dbx.DB) }
 /**
  * BaseAppConfig defines a BaseApp configuration option
  */
 interface BaseAppConfig {
  dbConnect: DBConnectFunc
  dataDir: string
  encryptionEnv: string
  queryTimeout: time.Duration
  dataMaxOpenConns: number
  dataMaxIdleConns: number
  auxMaxOpenConns: number
  auxMaxIdleConns: number
  isDev: boolean
 }
 /**
  * BaseApp implements CoreApp and defines the base PocketBase app structure.
  */
 interface BaseApp {
 }
 interface newBaseApp {
  /**
   * NewBaseApp creates and returns a new BaseApp instance
   * configured with the provided arguments.
   * 
   * To initialize the app, you need to call `app.Bootstrap()`.
   */
  (config: BaseAppConfig): (BaseApp)
 }
 interface BaseApp {
  /**
   * UnsafeWithoutHooks returns a shallow copy of the current app WITHOUT any registered hooks.
   * 
   * NB! Note that using the returned app instance may cause data integrity errors
   * since the Record validations and data normalizations (including files uploads)
   * rely on the app hooks to work.
   */
  unsafeWithoutHooks(): App
 }
 interface BaseApp {
  /**
   * Logger returns the default app logger.
   * 
   * If the application is not bootstrapped yet, fallbacks to slog.Default().
   */
  logger(): (slog.Logger)
 }
 interface BaseApp {
  /**
   * TxInfo returns the transaction associated with the current app instance (if any).
   * 
   * Could be used if you want to execute indirectly a function after
   * the related app transaction completes using `app.TxInfo().OnAfterFunc(callback)`.
   */
  txInfo(): (TxAppInfo)
 }
 interface BaseApp {
  /**
   * IsTransactional checks if the current app instance is part of a transaction.
   */
  isTransactional(): boolean
 }
 interface BaseApp {
  /**
   * IsBootstrapped checks if the application was initialized
   * (aka. whether Bootstrap() was called).
   */
  isBootstrapped(): boolean
 }
 interface BaseApp {
  /**
   * Bootstrap initializes the application
   * (aka. create data dir, open db connections, load settings, etc.).
   * 
   * It will call ResetBootstrapState() if the application was already bootstrapped.
   */
  bootstrap(): void
 }
 interface closer {
  [key:string]: any;
  close(): void
 }
 interface BaseApp {
  /**
   * ResetBootstrapState releases the initialized core app resources
   * (closing db connections, stopping cron ticker, etc.).
   */
  resetBootstrapState(): void
 }
 interface BaseApp {
  /**
   * DB returns the default app data.db builder instance.
   * 
   * To minimize SQLITE_BUSY errors, it automatically routes the
   * SELECT queries to the underlying concurrent db pool and everything
   * else to the nonconcurrent one.
   * 
   * For more finer control over the used connections pools you can
   * call directly ConcurrentDB() or NonconcurrentDB().
   */
  db(): dbx.Builder
 }
 interface BaseApp {
  /**
   * ConcurrentDB returns the concurrent app data.db builder instance.
   * 
   * This method is used mainly internally for executing db read
   * operations in a concurrent/non-blocking manner.
   * 
   * Most users should use simply DB() as it will automatically
   * route the query execution to ConcurrentDB() or NonconcurrentDB().
   * 
   * In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
   */
  concurrentDB(): dbx.Builder
 }
 interface BaseApp {
  /**
   * NonconcurrentDB returns the nonconcurrent app data.db builder instance.
   * 
   * The returned db instance is limited only to a single open connection,
   * meaning that it can process only 1 db operation at a time (other queries queue up).
   * 
   * This method is used mainly internally and in the tests to execute write
   * (save/delete) db operations as it helps with minimizing the SQLITE_BUSY errors.
   * 
   * Most users should use simply DB() as it will automatically
   * route the query execution to ConcurrentDB() or NonconcurrentDB().
   * 
   * In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
   */
  nonconcurrentDB(): dbx.Builder
 }
 interface BaseApp {
  /**
   * AuxDB returns the app auxiliary.db builder instance.
   * 
   * To minimize SQLITE_BUSY errors, it automatically routes the
   * SELECT queries to the underlying concurrent db pool and everything
   * else to the nonconcurrent one.
   * 
   * For more finer control over the used connections pools you can
   * call directly AuxConcurrentDB() or AuxNonconcurrentDB().
   */
  auxDB(): dbx.Builder
 }
 interface BaseApp {
  /**
   * AuxConcurrentDB returns the concurrent app auxiliary.db builder instance.
   * 
   * This method is used mainly internally for executing db read
   * operations in a concurrent/non-blocking manner.
   * 
   * Most users should use simply AuxDB() as it will automatically
   * route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
   * 
   * In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
   */
  auxConcurrentDB(): dbx.Builder
 }
 interface BaseApp {
  /**
   * AuxNonconcurrentDB returns the nonconcurrent app auxiliary.db builder instance.
   * 
   * The returned db instance is limited only to a single open connection,
   * meaning that it can process only 1 db operation at a time (other queries queue up).
   * 
   * This method is used mainly internally and in the tests to execute write
   * (save/delete) db operations as it helps with minimizing the SQLITE_BUSY errors.
   * 
   * Most users should use simply AuxDB() as it will automatically
   * route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
   * 
   * In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
   */
  auxNonconcurrentDB(): dbx.Builder
 }
 interface BaseApp {
  /**
   * DataDir returns the app data directory path.
   */
  dataDir(): string
 }
 interface BaseApp {
  /**
   * EncryptionEnv returns the name of the app secret env key
   * (currently used primarily for optional settings encryption but this may change in the future).
   */
  encryptionEnv(): string
 }
 interface BaseApp {
  /**
   * IsDev returns whether the app is in dev mode.
   * 
   * When enabled logs, executed sql statements, etc. are printed to the stderr.
   */
  isDev(): boolean
 }
 interface BaseApp {
  /**
   * Settings returns the loaded app settings.
   */
  settings(): (Settings)
 }
 interface BaseApp {
  /**
   * Store returns the app runtime store.
   */
  store(): (store.Store<string, any>)
 }
 interface BaseApp {
  /**
   * Cron returns the app cron instance.
   */
  cron(): (cron.Cron)
 }
 interface BaseApp {
  /**
   * SubscriptionsBroker returns the app realtime subscriptions broker instance.
   */
  subscriptionsBroker(): (subscriptions.Broker)
 }
 interface BaseApp {
  /**
   * NewMailClient creates and returns a new SMTP or Sendmail client
   * based on the current app settings.
   */
  newMailClient(): mailer.Mailer
 }
 interface BaseApp {
  /**
   * NewFilesystem creates a new local or S3 filesystem instance
   * for managing regular app files (ex. record uploads)
   * based on the current app settings.
   * 
   * NB! Make sure to call Close() on the returned result
   * after you are done working with it.
   */
  newFilesystem(): (filesystem.System)
 }
 interface BaseApp {
  /**
   * NewBackupsFilesystem creates a new local or S3 filesystem instance
   * for managing app backups based on the current app settings.
   * 
   * NB! Make sure to call Close() on the returned result
   * after you are done working with it.
   */
  newBackupsFilesystem(): (filesystem.System)
 }
 interface BaseApp {
  /**
   * Restart restarts (aka. replaces) the current running application process.
   * 
   * NB! It relies on execve which is supported only on UNIX based systems.
   */
  restart(): void
 }
 interface BaseApp {
  /**
   * RunSystemMigrations applies all new migrations registered in the [core.SystemMigrations] list.
   */
  runSystemMigrations(): void
 }
 interface BaseApp {
  /**
   * RunAppMigrations applies all new migrations registered in the [CoreAppMigrations] list.
   */
  runAppMigrations(): void
 }
 interface BaseApp {
  /**
   * RunAllMigrations applies all system and app migrations
   * (aka. from both [core.SystemMigrations] and [CoreAppMigrations]).
   */
  runAllMigrations(): void
 }
 interface BaseApp {
  onBootstrap(): (hook.Hook<BootstrapEvent | undefined>)
 }
 interface BaseApp {
  onServe(): (hook.Hook<ServeEvent | undefined>)
 }
 interface BaseApp {
  onTerminate(): (hook.Hook<TerminateEvent | undefined>)
 }
 interface BaseApp {
  onBackupCreate(): (hook.Hook<BackupEvent | undefined>)
 }
 interface BaseApp {
  onBackupRestore(): (hook.Hook<BackupEvent | undefined>)
 }
 interface BaseApp {
  onModelCreate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelCreateExecute(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelAfterCreateSuccess(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelAfterCreateError(...tags: string[]): (hook.TaggedHook<ModelErrorEvent | undefined>)
 }
 interface BaseApp {
  onModelUpdate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelUpdateExecute(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelAfterUpdateSuccess(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelAfterUpdateError(...tags: string[]): (hook.TaggedHook<ModelErrorEvent | undefined>)
 }
 interface BaseApp {
  onModelValidate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelDelete(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelDeleteExecute(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelAfterDeleteSuccess(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined>)
 }
 interface BaseApp {
  onModelAfterDeleteError(...tags: string[]): (hook.TaggedHook<ModelErrorEvent | undefined>)
 }
 interface BaseApp {
  onRecordEnrich(...tags: string[]): (hook.TaggedHook<RecordEnrichEvent | undefined>)
 }
 interface BaseApp {
  onRecordValidate(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordCreate(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordCreateExecute(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordAfterCreateSuccess(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordAfterCreateError(...tags: string[]): (hook.TaggedHook<RecordErrorEvent | undefined>)
 }
 interface BaseApp {
  onRecordUpdate(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordUpdateExecute(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordAfterUpdateSuccess(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordAfterUpdateError(...tags: string[]): (hook.TaggedHook<RecordErrorEvent | undefined>)
 }
 interface BaseApp {
  onRecordDelete(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordDeleteExecute(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordAfterDeleteSuccess(...tags: string[]): (hook.TaggedHook<RecordEvent | undefined>)
 }
 interface BaseApp {
  onRecordAfterDeleteError(...tags: string[]): (hook.TaggedHook<RecordErrorEvent | undefined>)
 }
 interface BaseApp {
  onCollectionValidate(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionCreate(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionCreateExecute(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionAfterCreateSuccess(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionAfterCreateError(...tags: string[]): (hook.TaggedHook<CollectionErrorEvent | undefined>)
 }
 interface BaseApp {
  onCollectionUpdate(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionUpdateExecute(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionAfterUpdateSuccess(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionAfterUpdateError(...tags: string[]): (hook.TaggedHook<CollectionErrorEvent | undefined>)
 }
 interface BaseApp {
  onCollectionDelete(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionDeleteExecute(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionAfterDeleteSuccess(...tags: string[]): (hook.TaggedHook<CollectionEvent | undefined>)
 }
 interface BaseApp {
  onCollectionAfterDeleteError(...tags: string[]): (hook.TaggedHook<CollectionErrorEvent | undefined>)
 }
 interface BaseApp {
  onMailerSend(): (hook.Hook<MailerEvent | undefined>)
 }
 interface BaseApp {
  onMailerRecordPasswordResetSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
 }
 interface BaseApp {
  onMailerRecordVerificationSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
 }
 interface BaseApp {
  onMailerRecordEmailChangeSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
 }
 interface BaseApp {
  onMailerRecordOTPSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
 }
 interface BaseApp {
  onMailerRecordAuthAlertSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined>)
 }
 interface BaseApp {
  onRealtimeConnectRequest(): (hook.Hook<RealtimeConnectRequestEvent | undefined>)
 }
 interface BaseApp {
  onRealtimeMessageSend(): (hook.Hook<RealtimeMessageEvent | undefined>)
 }
 interface BaseApp {
  onRealtimeSubscribeRequest(): (hook.Hook<RealtimeSubscribeRequestEvent | undefined>)
 }
 interface BaseApp {
  onSettingsListRequest(): (hook.Hook<SettingsListRequestEvent | undefined>)
 }
 interface BaseApp {
  onSettingsUpdateRequest(): (hook.Hook<SettingsUpdateRequestEvent | undefined>)
 }
 interface BaseApp {
  onSettingsReload(): (hook.Hook<SettingsReloadEvent | undefined>)
 }
 interface BaseApp {
  onFileDownloadRequest(...tags: string[]): (hook.TaggedHook<FileDownloadRequestEvent | undefined>)
 }
 interface BaseApp {
  onFileTokenRequest(...tags: string[]): (hook.TaggedHook<FileTokenRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordAuthRequest(...tags: string[]): (hook.TaggedHook<RecordAuthRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordAuthWithPasswordRequest(...tags: string[]): (hook.TaggedHook<RecordAuthWithPasswordRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordAuthWithOAuth2Request(...tags: string[]): (hook.TaggedHook<RecordAuthWithOAuth2RequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordAuthRefreshRequest(...tags: string[]): (hook.TaggedHook<RecordAuthRefreshRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordRequestPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordRequestPasswordResetRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordConfirmPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmPasswordResetRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordRequestVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordRequestVerificationRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordConfirmVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmVerificationRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordRequestEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEmailChangeRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordConfirmEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmEmailChangeRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordRequestOTPRequest(...tags: string[]): (hook.TaggedHook<RecordCreateOTPRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordAuthWithOTPRequest(...tags: string[]): (hook.TaggedHook<RecordAuthWithOTPRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordsListRequest(...tags: string[]): (hook.TaggedHook<RecordsListRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordViewRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordCreateRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordUpdateRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
 }
 interface BaseApp {
  onRecordDeleteRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEvent | undefined>)
 }
 interface BaseApp {
  onCollectionsListRequest(): (hook.Hook<CollectionsListRequestEvent | undefined>)
 }
 interface BaseApp {
  onCollectionViewRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
 }
 interface BaseApp {
  onCollectionCreateRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
 }
 interface BaseApp {
  onCollectionUpdateRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
 }
 interface BaseApp {
  onCollectionDeleteRequest(): (hook.Hook<CollectionRequestEvent | undefined>)
 }
 interface BaseApp {
  onCollectionsImportRequest(): (hook.Hook<CollectionsImportRequestEvent | undefined>)
 }
 interface BaseApp {
  onBatchRequest(): (hook.Hook<BatchRequestEvent | undefined>)
 }
 interface BaseApp {
  /**
   * CreateBackup creates a new backup of the current app pb_data directory.
   * 
   * If name is empty, it will be autogenerated.
   * If backup with the same name exists, the new backup file will replace it.
   * 
   * The backup is executed within a transaction, meaning that new writes
   * will be temporary "blocked" until the backup file is generated.
   * 
   * To safely perform the backup, it is recommended to have free disk space
   * for at least 2x the size of the pb_data directory.
   * 
   * By default backups are stored in pb_data/backups
   * (the backups directory itself is excluded from the generated backup).
   * 
   * When using S3 storage for the uploaded collection files, you have to
   * take care manually to backup those since they are not part of the pb_data.
   * 
   * Backups can be stored on S3 if it is configured in app.Settings().Backups.
   */
  createBackup(ctx: context.Context, name: string): void
 }
 interface BaseApp {
  /**
   * RestoreBackup restores the backup with the specified name and restarts
   * the current running application process.
   * 
   * NB! This feature is experimental and currently is expected to work only on UNIX based systems.
   * 
   * To safely perform the restore it is recommended to have free disk space
   * for at least 2x the size of the restored pb_data backup.
   * 
   * The performed steps are:
   * 
   *  1. Download the backup with the specified name in a temp location
   * ```
   *     (this is in case of S3; otherwise it creates a temp copy of the zip)
   * ```
   * 
   *  2. Extract the backup in a temp directory inside the app "pb_data"
   * ```
   *     (eg. "pb_data/.pb_temp_to_delete/pb_restore").
   * ```
   * 
   *  3. Move the current app "pb_data" content (excluding the local backups and the special temp dir)
   * ```
   *     under another temp sub dir that will be deleted on the next app start up
   *     (eg. "pb_data/.pb_temp_to_delete/old_pb_data").
   *     This is because on some environments it may not be allowed
   *     to delete the currently open "pb_data" files.
   * ```
   * 
   *  4. Move the extracted dir content to the app "pb_data".
   * 
   *  5. Restart the app (on successful app bootstap it will also remove the old pb_data).
   * 
   * If a failure occure during the restore process the dir changes are reverted.
   * If for whatever reason the revert is not possible, it panics.
   * 
   * Note that if your pb_data has custom network mounts as subdirectories, then
   * it is possible the restore to fail during the `os.Rename` operations
   * (see https://github.com/pocketbase/pocketbase/issues/4647).
   */
  restoreBackup(ctx: context.Context, name: string): void
 }
 interface BaseApp {
  /**
   * ImportCollectionsByMarshaledJSON is the same as [ImportCollections]
   * but accept marshaled json array as import data (usually used for the autogenerated snapshots).
   */
  importCollectionsByMarshaledJSON(rawSliceOfMaps: string|Array<number>, deleteMissing: boolean): void
 }
 interface BaseApp {
  /**
   * ImportCollections imports the provided collections data in a single transaction.
   * 
   * For existing matching collections, the imported data is unmarshaled on top of the existing model.
   * 
   * NB! If deleteMissing is true, ALL NON-SYSTEM COLLECTIONS AND SCHEMA FIELDS,
   * that are not present in the imported configuration, WILL BE DELETED
   * (this includes their related records data).
   */
  importCollections(toImport: Array<_TygojaDict>, deleteMissing: boolean): void
 }
 /**
  * @todo experiment eventually replacing the rules *string with a struct?
  */
 type _sLoVRbR = BaseModel
 interface baseCollection extends _sLoVRbR {
  listRule?: string
  viewRule?: string
  createRule?: string
  updateRule?: string
  deleteRule?: string
  /**
   * RawOptions represents the raw serialized collection option loaded from the DB.
   * NB! This field shouldn't be modified manually. It is automatically updated
   * with the collection type specific option before save.
   */
  rawOptions: types.JSONRaw
  name: string
  type: string
  fields: FieldsList
  indexes: types.JSONArray<string>
  created: types.DateTime
  updated: types.DateTime
  /**
   * System prevents the collection rename, deletion and rules change.
   * It is used primarily for internal purposes for collections like "_superusers", "_externalAuths", etc.
   */
  system: boolean
 }
 /**
  * Collection defines the table, fields and various options related to a set of records.
  */
 type _sCiyZlG = baseCollection&collectionAuthOptions&collectionViewOptions
 interface Collection extends _sCiyZlG {
 }
 interface newCollection {
  /**
   * NewCollection initializes and returns a new Collection model with the specified type and name.
   * 
   * It also loads the minimal default configuration for the collection
   * (eg. system fields, indexes, type specific options, etc.).
   */
  (typ: string, name: string, ...optId: string[]): (Collection)
 }
 interface newBaseCollection {
  /**
   * NewBaseCollection initializes and returns a new "base" Collection model.
   * 
   * It also loads the minimal default configuration for the collection
   * (eg. system fields, indexes, type specific options, etc.).
   */
  (name: string, ...optId: string[]): (Collection)
 }
 interface newViewCollection {
  /**
   * NewViewCollection initializes and returns a new "view" Collection model.
   * 
   * It also loads the minimal default configuration for the collection
   * (eg. system fields, indexes, type specific options, etc.).
   */
  (name: string, ...optId: string[]): (Collection)
 }
 interface newAuthCollection {
  /**
   * NewAuthCollection initializes and returns a new "auth" Collection model.
   * 
   * It also loads the minimal default configuration for the collection
   * (eg. system fields, indexes, type specific options, etc.).
   */
  (name: string, ...optId: string[]): (Collection)
 }
 interface Collection {
  /**
   * TableName returns the Collection model SQL table name.
   */
  tableName(): string
 }
 interface Collection {
  /**
   * BaseFilesPath returns the storage dir path used by the collection.
   */
  baseFilesPath(): string
 }
 interface Collection {
  /**
   * IsBase checks if the current collection has "base" type.
   */
  isBase(): boolean
 }
 interface Collection {
  /**
   * IsAuth checks if the current collection has "auth" type.
   */
  isAuth(): boolean
 }
 interface Collection {
  /**
   * IsView checks if the current collection has "view" type.
   */
  isView(): boolean
 }
 interface Collection {
  /**
   * IntegrityChecks toggles the current collection integrity checks (ex. checking references on delete).
   */
  integrityChecks(enable: boolean): void
 }
 interface Collection {
  /**
   * PostScan implements the [dbx.PostScanner] interface to auto unmarshal
   * the raw serialized options into the concrete type specific fields.
   */
  postScan(): void
 }
 interface Collection {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   * 
   * For new/"blank" Collection models it replaces the model with a factory
   * instance and then unmarshal the provided data one on top of it.
   */
  unmarshalJSON(b: string|Array<number>): void
 }
 interface Collection {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   * 
   * Note that non-type related fields are ignored from the serialization
   * (ex. for "view" colections the "auth" fields are skipped).
   */
  marshalJSON(): string|Array<number>
 }
 interface Collection {
  /**
   * String returns a string representation of the current collection.
   */
  string(): string
 }
 interface Collection {
  /**
   * DBExport prepares and exports the current collection data for db persistence.
   */
  dbExport(app: App): _TygojaDict
 }
 interface Collection {
  /**
   * GetIndex returns s single Collection index expression by its name.
   */
  getIndex(name: string): string
 }
 interface Collection {
  /**
   * AddIndex adds a new index into the current collection.
   * 
   * If the collection has an existing index matching the new name it will be replaced with the new one.
   */
  addIndex(name: string, unique: boolean, columnsExpr: string, optWhereExpr: string): void
 }
 interface Collection {
  /**
   * RemoveIndex removes a single index with the specified name from the current collection.
   */
  removeIndex(name: string): void
 }
 /**
  * collectionAuthOptions defines the options for the "auth" type collection.
  */
 interface collectionAuthOptions {
  /**
   * AuthRule could be used to specify additional record constraints
   * applied after record authentication and right before returning the
   * auth token response to the client.
   * 
   * For example, to allow only verified users you could set it to
   * "verified = true".
   * 
   * Set it to empty string to allow any Auth collection record to authenticate.
   * 
   * Set it to nil to disallow authentication altogether for the collection
   * (that includes password, OAuth2, etc.).
   */
  authRule?: string
  /**
   * ManageRule gives admin-like permissions to allow fully managing
   * the auth record(s), eg. changing the password without requiring
   * to enter the old one, directly updating the verified state and email, etc.
   * 
   * This rule is executed in addition to the Create and Update API rules.
   */
  manageRule?: string
  /**
   * AuthAlert defines options related to the auth alerts on new device login.
   */
  authAlert: AuthAlertConfig
  /**
   * OAuth2 specifies whether OAuth2 auth is enabled for the collection
   * and which OAuth2 providers are allowed.
   */
  oauth2: OAuth2Config
  /**
   * PasswordAuth defines options related to the collection password authentication.
   */
  passwordAuth: PasswordAuthConfig
  /**
   * MFA defines options related to the Multi-factor authentication (MFA).
   */
  mfa: MFAConfig
  /**
   * OTP defines options related to the One-time password authentication (OTP).
   */
  otp: OTPConfig
  /**
   * Various token configurations
   * ---
   */
  authToken: TokenConfig
  passwordResetToken: TokenConfig
  emailChangeToken: TokenConfig
  verificationToken: TokenConfig
  fileToken: TokenConfig
  /**
   * Default email templates
   * ---
   */
  verificationTemplate: EmailTemplate
  resetPasswordTemplate: EmailTemplate
  confirmEmailChangeTemplate: EmailTemplate
 }
 interface EmailTemplate {
  subject: string
  body: string
 }
 interface EmailTemplate {
  /**
   * Validate makes EmailTemplate validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface EmailTemplate {
  /**
   * Resolve replaces the placeholder parameters in the current email
   * template and returns its components as ready-to-use strings.
   */
  resolve(placeholders: _TygojaDict): [string, string]
 }
 interface AuthAlertConfig {
  enabled: boolean
  emailTemplate: EmailTemplate
 }
 interface AuthAlertConfig {
  /**
   * Validate makes AuthAlertConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface TokenConfig {
  secret: string
  /**
   * Duration specifies how long an issued token to be valid (in seconds)
   */
  duration: number
 }
 interface TokenConfig {
  /**
   * Validate makes TokenConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface TokenConfig {
  /**
   * DurationTime returns the current Duration as [time.Duration].
   */
  durationTime(): time.Duration
 }
 interface OTPConfig {
  enabled: boolean
  /**
   * Duration specifies how long the OTP to be valid (in seconds)
   */
  duration: number
  /**
   * Length specifies the auto generated password length.
   */
  length: number
  /**
   * EmailTemplate is the default OTP email template that will be send to the auth record.
   * 
   * In addition to the system placeholders you can also make use of
   * [core.EmailPlaceholderOTPId] and [core.EmailPlaceholderOTP].
   */
  emailTemplate: EmailTemplate
 }
 interface OTPConfig {
  /**
   * Validate makes OTPConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface OTPConfig {
  /**
   * DurationTime returns the current Duration as [time.Duration].
   */
  durationTime(): time.Duration
 }
 interface MFAConfig {
  enabled: boolean
  /**
   * Duration specifies how long an issued MFA to be valid (in seconds)
   */
  duration: number
  /**
   * Rule is an optional field to restrict MFA only for the records that satisfy the rule.
   * 
   * Leave it empty to enable MFA for everyone.
   */
  rule: string
 }
 interface MFAConfig {
  /**
   * Validate makes MFAConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface MFAConfig {
  /**
   * DurationTime returns the current Duration as [time.Duration].
   */
  durationTime(): time.Duration
 }
 interface PasswordAuthConfig {
  enabled: boolean
  /**
   * IdentityFields is a list of field names that could be used as
   * identity during password authentication.
   * 
   * Usually only fields that has single column UNIQUE index are accepted as values.
   */
  identityFields: Array<string>
 }
 interface PasswordAuthConfig {
  /**
   * Validate makes PasswordAuthConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface OAuth2KnownFields {
  id: string
  name: string
  username: string
  avatarURL: string
 }
 interface OAuth2Config {
  providers: Array<OAuth2ProviderConfig>
  mappedFields: OAuth2KnownFields
  enabled: boolean
 }
 interface OAuth2Config {
  /**
   * GetProviderConfig returns the first OAuth2ProviderConfig that matches the specified name.
   * 
   * Returns false and zero config if no such provider is available in c.Providers.
   */
  getProviderConfig(name: string): [OAuth2ProviderConfig, boolean]
 }
 interface OAuth2Config {
  /**
   * Validate makes OAuth2Config validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface OAuth2ProviderConfig {
  /**
   * PKCE overwrites the default provider PKCE config option.
   * 
   * This usually shouldn't be needed but some OAuth2 vendors, like the LinkedIn OIDC,
   * may require manual adjustment due to returning error if extra parameters are added to the request
   * (https://github.com/pocketbase/pocketbase/discussions/3799#discussioncomment-7640312)
   */
  pkce?: boolean
  name: string
  clientId: string
  clientSecret: string
  authURL: string
  tokenURL: string
  userInfoURL: string
  displayName: string
  extra: _TygojaDict
 }
 interface OAuth2ProviderConfig {
  /**
   * Validate makes OAuth2ProviderConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface OAuth2ProviderConfig {
  /**
   * InitProvider returns a new auth.Provider instance loaded with the current OAuth2ProviderConfig options.
   */
  initProvider(): auth.Provider
 }
 /**
  * collectionBaseOptions defines the options for the "base" type collection.
  */
 interface collectionBaseOptions {
 }
 /**
  * collectionViewOptions defines the options for the "view" type collection.
  */
 interface collectionViewOptions {
  viewQuery: string
 }
 interface BaseApp {
  /**
   * CollectionQuery returns a new Collection select query.
   */
  collectionQuery(): (dbx.SelectQuery)
 }
 interface BaseApp {
  /**
   * FindCollections finds all collections by the given type(s).
   * 
   * If collectionTypes is not set, it returns all collections.
   * 
   * Example:
   * 
   * ```
   * 	app.FindAllCollections() // all collections
   * 	app.FindAllCollections("auth", "view") // only auth and view collections
   * ```
   */
  findAllCollections(...collectionTypes: string[]): Array<(Collection | undefined)>
 }
 interface BaseApp {
  /**
   * ReloadCachedCollections fetches all collections and caches them into the app store.
   */
  reloadCachedCollections(): void
 }
 interface BaseApp {
  /**
   * FindCollectionByNameOrId finds a single collection by its name (case insensitive) or id.
   */
  findCollectionByNameOrId(nameOrId: string): (Collection)
 }
 interface BaseApp {
  /**
   * FindCachedCollectionByNameOrId is similar to [BaseApp.FindCollectionByNameOrId]
   * but retrieves the Collection from the app cache instead of making a db call.
   * 
   * NB! This method is suitable for read-only Collection operations.
   * 
   * Returns [sql.ErrNoRows] if no Collection is found for consistency
   * with the [BaseApp.FindCollectionByNameOrId] method.
   * 
   * If you plan making changes to the returned Collection model,
   * use [BaseApp.FindCollectionByNameOrId] instead.
   * 
   * Caveats:
   * 
   * ```
   *   - The returned Collection should be used only for read-only operations.
   *     Avoid directly modifying the returned cached Collection as it will affect
   *     the global cached value even if you don't persist the changes in the database!
   *   - If you are updating a Collection in a transaction and then call this method before commit,
   *     it'll return the cached Collection state and not the one from the uncommitted transaction.
   *   - The cache is automatically updated on collections db change (create/update/delete).
   *     To manually reload the cache you can call [BaseApp.ReloadCachedCollections].
   * ```
   */
  findCachedCollectionByNameOrId(nameOrId: string): (Collection)
 }
 interface BaseApp {
  /**
   * FindCollectionReferences returns information for all relation fields
   * referencing the provided collection.
   * 
   * If the provided collection has reference to itself then it will be
   * also included in the result. To exclude it, pass the collection id
   * as the excludeIds argument.
   */
  findCollectionReferences(collection: Collection, ...excludeIds: string[]): _TygojaDict
 }
 interface BaseApp {
  /**
   * FindCachedCollectionReferences is similar to [BaseApp.FindCollectionReferences]
   * but retrieves the Collection from the app cache instead of making a db call.
   * 
   * NB! This method is suitable for read-only Collection operations.
   * 
   * If you plan making changes to the returned Collection model,
   * use [BaseApp.FindCollectionReferences] instead.
   * 
   * Caveats:
   * 
   * ```
   *   - The returned Collection should be used only for read-only operations.
   *     Avoid directly modifying the returned cached Collection as it will affect
   *     the global cached value even if you don't persist the changes in the database!
   *   - If you are updating a Collection in a transaction and then call this method before commit,
   *     it'll return the cached Collection state and not the one from the uncommitted transaction.
   *   - The cache is automatically updated on collections db change (create/update/delete).
   *     To manually reload the cache you can call [BaseApp.ReloadCachedCollections].
   * ```
   */
  findCachedCollectionReferences(collection: Collection, ...excludeIds: string[]): _TygojaDict
 }
 interface BaseApp {
  /**
   * IsCollectionNameUnique checks that there is no existing collection
   * with the provided name (case insensitive!).
   * 
   * Note: case insensitive check because the name is used also as
   * table name for the records.
   */
  isCollectionNameUnique(name: string, ...excludeIds: string[]): boolean
 }
 interface BaseApp {
  /**
   * TruncateCollection deletes all records associated with the provided collection.
   * 
   * The truncate operation is executed in a single transaction,
   * aka. either everything is deleted or none.
   * 
   * Note that this method will also trigger the records related
   * cascade and file delete actions.
   */
  truncateCollection(collection: Collection): void
 }
 interface BaseApp {
  /**
   * SyncRecordTableSchema compares the two provided collections
   * and applies the necessary related record table changes.
   * 
   * If oldCollection is null, then only newCollection is used to create the record table.
   * 
   * This method is automatically invoked as part of a collection create/update/delete operation.
   */
  syncRecordTableSchema(newCollection: Collection, oldCollection: Collection): void
 }
 interface collectionValidator {
 }
 interface optionsValidator {
  [key:string]: any;
 }
 /**
  * DBExporter defines an interface for custom DB data export.
  * Usually used as part of [App.Save].
  */
 interface DBExporter {
  [key:string]: any;
  /**
   * DBExport returns a key-value map with the data to be used when saving the struct in the database.
   */
  dbExport(app: App): _TygojaDict
 }
 /**
  * PreValidator defines an optional model interface for registering a
  * function that will run BEFORE firing the validation hooks (see [App.ValidateWithContext]).
  */
 interface PreValidator {
  [key:string]: any;
  /**
   * PreValidate defines a function that runs BEFORE the validation hooks.
   */
  preValidate(ctx: context.Context, app: App): void
 }
 /**
  * PostValidator defines an optional model interface for registering a
  * function that will run AFTER executing the validation hooks (see [App.ValidateWithContext]).
  */
 interface PostValidator {
  [key:string]: any;
  /**
   * PostValidate defines a function that runs AFTER the successful
   * execution of the validation hooks.
   */
  postValidate(ctx: context.Context, app: App): void
 }
 interface generateDefaultRandomId {
  /**
   * GenerateDefaultRandomId generates a default random id string
   * (note: the generated random string is not intended for security purposes).
   */
  (): string
 }
 interface BaseApp {
  /**
   * ModelQuery creates a new preconfigured select data.db query with preset
   * SELECT, FROM and other common fields based on the provided model.
   */
  modelQuery(m: Model): (dbx.SelectQuery)
 }
 interface BaseApp {
  /**
   * AuxModelQuery creates a new preconfigured select auxiliary.db query with preset
   * SELECT, FROM and other common fields based on the provided model.
   */
  auxModelQuery(m: Model): (dbx.SelectQuery)
 }
 interface BaseApp {
  /**
   * Delete deletes the specified model from the regular app database.
   */
  delete(model: Model): void
 }
 interface BaseApp {
  /**
   * Delete deletes the specified model from the regular app database
   * (the context could be used to limit the query execution).
   */
  deleteWithContext(ctx: context.Context, model: Model): void
 }
 interface BaseApp {
  /**
   * AuxDelete deletes the specified model from the auxiliary database.
   */
  auxDelete(model: Model): void
 }
 interface BaseApp {
  /**
   * AuxDeleteWithContext deletes the specified model from the auxiliary database
   * (the context could be used to limit the query execution).
   */
  auxDeleteWithContext(ctx: context.Context, model: Model): void
 }
 interface BaseApp {
  /**
   * Save validates and saves the specified model into the regular app database.
   * 
   * If you don't want to run validations, use [App.SaveNoValidate()].
   */
  save(model: Model): void
 }
 interface BaseApp {
  /**
   * SaveWithContext is the same as [App.Save()] but allows specifying a context to limit the db execution.
   * 
   * If you don't want to run validations, use [App.SaveNoValidateWithContext()].
   */
  saveWithContext(ctx: context.Context, model: Model): void
 }
 interface BaseApp {
  /**
   * SaveNoValidate saves the specified model into the regular app database without performing validations.
   * 
   * If you want to also run validations before persisting, use [App.Save()].
   */
  saveNoValidate(model: Model): void
 }
 interface BaseApp {
  /**
   * SaveNoValidateWithContext is the same as [App.SaveNoValidate()]
   * but allows specifying a context to limit the db execution.
   * 
   * If you want to also run validations before persisting, use [App.SaveWithContext()].
   */
  saveNoValidateWithContext(ctx: context.Context, model: Model): void
 }
 interface BaseApp {
  /**
   * AuxSave validates and saves the specified model into the auxiliary app database.
   * 
   * If you don't want to run validations, use [App.AuxSaveNoValidate()].
   */
  auxSave(model: Model): void
 }
 interface BaseApp {
  /**
   * AuxSaveWithContext is the same as [App.AuxSave()] but allows specifying a context to limit the db execution.
   * 
   * If you don't want to run validations, use [App.AuxSaveNoValidateWithContext()].
   */
  auxSaveWithContext(ctx: context.Context, model: Model): void
 }
 interface BaseApp {
  /**
   * AuxSaveNoValidate saves the specified model into the auxiliary app database without performing validations.
   * 
   * If you want to also run validations before persisting, use [App.AuxSave()].
   */
  auxSaveNoValidate(model: Model): void
 }
 interface BaseApp {
  /**
   * AuxSaveNoValidateWithContext is the same as [App.AuxSaveNoValidate()]
   * but allows specifying a context to limit the db execution.
   * 
   * If you want to also run validations before persisting, use [App.AuxSaveWithContext()].
   */
  auxSaveNoValidateWithContext(ctx: context.Context, model: Model): void
 }
 interface BaseApp {
  /**
   * Validate triggers the OnModelValidate hook for the specified model.
   */
  validate(model: Model): void
 }
 interface BaseApp {
  /**
   * ValidateWithContext is the same as Validate but allows specifying the ModelEvent context.
   */
  validateWithContext(ctx: context.Context, model: Model): void
 }
 /**
  * note: expects both builder to use the same driver
  */
 interface dualDBBuilder {
 }
 interface dualDBBuilder {
  /**
   * Select implements the [dbx.Builder.Select] interface method.
   */
  select(...cols: string[]): (dbx.SelectQuery)
 }
 interface dualDBBuilder {
  /**
   * Model implements the [dbx.Builder.Model] interface method.
   */
  model(data: {
   }): (dbx.ModelQuery)
 }
 interface dualDBBuilder {
  /**
   * GeneratePlaceholder implements the [dbx.Builder.GeneratePlaceholder] interface method.
   */
  generatePlaceholder(i: number): string
 }
 interface dualDBBuilder {
  /**
   * Quote implements the [dbx.Builder.Quote] interface method.
   */
  quote(str: string): string
 }
 interface dualDBBuilder {
  /**
   * QuoteSimpleTableName implements the [dbx.Builder.QuoteSimpleTableName] interface method.
   */
  quoteSimpleTableName(table: string): string
 }
 interface dualDBBuilder {
  /**
   * QuoteSimpleColumnName implements the [dbx.Builder.QuoteSimpleColumnName] interface method.
   */
  quoteSimpleColumnName(col: string): string
 }
 interface dualDBBuilder {
  /**
   * QueryBuilder implements the [dbx.Builder.QueryBuilder] interface method.
   */
  queryBuilder(): dbx.QueryBuilder
 }
 interface dualDBBuilder {
  /**
   * Insert implements the [dbx.Builder.Insert] interface method.
   */
  insert(table: string, cols: dbx.Params): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * Upsert implements the [dbx.Builder.Upsert] interface method.
   */
  upsert(table: string, cols: dbx.Params, ...constraints: string[]): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * Update implements the [dbx.Builder.Update] interface method.
   */
  update(table: string, cols: dbx.Params, where: dbx.Expression): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * Delete implements the [dbx.Builder.Delete] interface method.
   */
  delete(table: string, where: dbx.Expression): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * CreateTable implements the [dbx.Builder.CreateTable] interface method.
   */
  createTable(table: string, cols: _TygojaDict, ...options: string[]): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * RenameTable implements the [dbx.Builder.RenameTable] interface method.
   */
  renameTable(oldName: string, newName: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * DropTable implements the [dbx.Builder.DropTable] interface method.
   */
  dropTable(table: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * TruncateTable implements the [dbx.Builder.TruncateTable] interface method.
   */
  truncateTable(table: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * AddColumn implements the [dbx.Builder.AddColumn] interface method.
   */
  addColumn(table: string, col: string, typ: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * DropColumn implements the [dbx.Builder.DropColumn] interface method.
   */
  dropColumn(table: string, col: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * RenameColumn implements the [dbx.Builder.RenameColumn] interface method.
   */
  renameColumn(table: string, oldName: string, newName: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * AlterColumn implements the [dbx.Builder.AlterColumn] interface method.
   */
  alterColumn(table: string, col: string, typ: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * AddPrimaryKey implements the [dbx.Builder.AddPrimaryKey] interface method.
   */
  addPrimaryKey(table: string, name: string, ...cols: string[]): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * DropPrimaryKey implements the [dbx.Builder.DropPrimaryKey] interface method.
   */
  dropPrimaryKey(table: string, name: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * AddForeignKey implements the [dbx.Builder.AddForeignKey] interface method.
   */
  addForeignKey(table: string, name: string, cols: Array<string>, refCols: Array<string>, refTable: string, ...options: string[]): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * DropForeignKey implements the [dbx.Builder.DropForeignKey] interface method.
   */
  dropForeignKey(table: string, name: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * CreateIndex implements the [dbx.Builder.CreateIndex] interface method.
   */
  createIndex(table: string, name: string, ...cols: string[]): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * CreateUniqueIndex implements the [dbx.Builder.CreateUniqueIndex] interface method.
   */
  createUniqueIndex(table: string, name: string, ...cols: string[]): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * DropIndex implements the [dbx.Builder.DropIndex] interface method.
   */
  dropIndex(table: string, name: string): (dbx.Query)
 }
 interface dualDBBuilder {
  /**
   * NewQuery implements the [dbx.Builder.NewQuery] interface method by
   * routing the SELECT queries to the concurrent builder instance.
   */
  newQuery(str: string): (dbx.Query)
 }
 interface defaultDBConnect {
  (dbPath: string): (dbx.DB)
 }
 /**
  * Model defines an interface with common methods that all db models should have.
  * 
  * Note: for simplicity composite pk are not supported.
  */
 interface Model {
  [key:string]: any;
  tableName(): string
  pk(): any
  lastSavedPK(): any
  isNew(): boolean
  markAsNew(): void
  markAsNotNew(): void
 }
 /**
  * BaseModel defines a base struct that is intended to be embedded into other custom models.
  */
 interface BaseModel {
  /**
   * Id is the primary key of the model.
   * It is usually autogenerated by the parent model implementation.
   */
  id: string
 }
 interface BaseModel {
  /**
   * LastSavedPK returns the last saved primary key of the model.
   * 
   * Its value is updated to the latest PK value after MarkAsNotNew() or PostScan() calls.
   */
  lastSavedPK(): any
 }
 interface BaseModel {
  pk(): any
 }
 interface BaseModel {
  /**
   * IsNew indicates what type of db query (insert or update)
   * should be used with the model instance.
   */
  isNew(): boolean
 }
 interface BaseModel {
  /**
   * MarkAsNew clears the pk field and marks the current model as "new"
   * (aka. forces m.IsNew() to be true).
   */
  markAsNew(): void
 }
 interface BaseModel {
  /**
   * MarkAsNew set the pk field to the Id value and marks the current model
   * as NOT "new" (aka. forces m.IsNew() to be false).
   */
  markAsNotNew(): void
 }
 interface BaseModel {
  /**
   * PostScan implements the [dbx.PostScanner] interface.
   * 
   * It is usually executed right after the model is populated with the db row values.
   */
  postScan(): void
 }
 interface BaseApp {
  /**
   * TableColumns returns all column names of a single table by its name.
   */
  tableColumns(tableName: string): Array<string>
 }
 interface TableInfoRow {
  /**
   * the `db:"pk"` tag has special semantic so we cannot rename
   * the original field without specifying a custom mapper
   */
  pk: number
  index: number
  name: string
  type: string
  notNull: boolean
  defaultValue: sql.NullString
 }
 interface BaseApp {
  /**
   * TableInfo returns the "table_info" pragma result for the specified table.
   */
  tableInfo(tableName: string): Array<(TableInfoRow | undefined)>
 }
 interface BaseApp {
  /**
   * TableIndexes returns a name grouped map with all non empty index of the specified table.
   * 
   * Note: This method doesn't return an error on nonexisting table.
   */
  tableIndexes(tableName: string): _TygojaDict
 }
 interface BaseApp {
  /**
   * DeleteTable drops the specified table.
   * 
   * This method is a no-op if a table with the provided name doesn't exist.
   * 
   * NB! Be aware that this method is vulnerable to SQL injection and the
   * "tableName" argument must come only from trusted input!
   */
  deleteTable(tableName: string): void
 }
 interface BaseApp {
  /**
   * HasTable checks if a table (or view) with the provided name exists (case insensitive).
   * in the data.db.
   */
  hasTable(tableName: string): boolean
 }
 interface BaseApp {
  /**
   * AuxHasTable checks if a table (or view) with the provided name exists (case insensitive)
   * in the auixiliary.db.
   */
  auxHasTable(tableName: string): boolean
 }
 interface BaseApp {
  /**
   * Vacuum executes VACUUM on the data.db in order to reclaim unused data db disk space.
   */
  vacuum(): void
 }
 interface BaseApp {
  /**
   * AuxVacuum executes VACUUM on the auxiliary.db in order to reclaim unused auxiliary db disk space.
   */
  auxVacuum(): void
 }
 interface BaseApp {
  /**
   * RunInTransaction wraps fn into a transaction for the regular app database.
   * 
   * It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
   */
  runInTransaction(fn: (txApp: App) => void): void
 }
 interface BaseApp {
  /**
   * AuxRunInTransaction wraps fn into a transaction for the auxiliary app database.
   * 
   * It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
   */
  auxRunInTransaction(fn: (txApp: App) => void): void
 }
 /**
  * TxAppInfo represents an active transaction context associated to an existing app instance.
  */
 interface TxAppInfo {
 }
 interface TxAppInfo {
  /**
   * OnComplete registers the provided callback that will be invoked
   * once the related transaction ends (either completes successfully or rollbacked with an error).
   * 
   * The callback receives the transaction error (if any) as its argument.
   * Any additional errors returned by the OnComplete callbacks will be
   * joined together with txErr when returning the final transaction result.
   */
  onComplete(fn: (txErr: Error) => void): void
 }
 /**
  * RequestEvent defines the PocketBase router handler event.
  */
 type _snSurdV = router.Event
 interface RequestEvent extends _snSurdV {
  app: App
  auth?: Record
 }
 interface RequestEvent {
  /**
   * RealIP returns the "real" IP address from the configured trusted proxy headers.
   * 
   * If Settings.TrustedProxy is not configured or the found IP is empty,
   * it fallbacks to e.RemoteIP().
   * 
   * NB!
   * Be careful when used in a security critical context as it relies on
   * the trusted proxy to be properly configured and your app to be accessible only through it.
   * If you are not sure, use e.RemoteIP().
   */
  realIP(): string
 }
 interface RequestEvent {
  /**
   * HasSuperuserAuth checks whether the current RequestEvent has superuser authentication loaded.
   */
  hasSuperuserAuth(): boolean
 }
 interface RequestEvent {
  /**
   * RequestInfo parses the current request into RequestInfo instance.
   * 
   * Note that the returned result is cached to avoid copying the request data multiple times
   * but the auth state and other common store items are always refreshed in case they were changed by another handler.
   */
  requestInfo(): (RequestInfo)
 }
 /**
  * RequestInfo defines a HTTP request data struct, usually used
  * as part of the `@request.*` filter resolver.
  * 
  * The Query and Headers fields contains only the first value for each found entry.
  */
 interface RequestInfo {
  query: _TygojaDict
  headers: _TygojaDict
  body: _TygojaDict
  auth?: Record
  method: string
  context: string
 }
 interface RequestInfo {
  /**
   * HasSuperuserAuth checks whether the current RequestInfo instance
   * has superuser authentication loaded.
   */
  hasSuperuserAuth(): boolean
 }
 interface RequestInfo {
  /**
   * Clone creates a new shallow copy of the current RequestInfo and its Auth record (if any).
   */
  clone(): (RequestInfo)
 }
 type _sXmDOQC = hook.Event&RequestEvent
 interface BatchRequestEvent extends _sXmDOQC {
  batch: Array<(InternalRequest | undefined)>
 }
 interface InternalRequest {
  /**
   * note: for uploading files the value must be either *filesystem.File or []*filesystem.File
   */
  body: _TygojaDict
  headers: _TygojaDict
  method: string
  url: string
 }
 interface InternalRequest {
  validate(): void
 }
 interface HookTagger {
  [key:string]: any;
  hookTags(): Array<string>
 }
 interface baseModelEventData {
  model: Model
 }
 interface baseModelEventData {
  tags(): Array<string>
 }
 interface baseRecordEventData {
  record?: Record
 }
 interface baseRecordEventData {
  tags(): Array<string>
 }
 interface baseCollectionEventData {
  collection?: Collection
 }
 interface baseCollectionEventData {
  tags(): Array<string>
 }
 type _sDtSitQ = hook.Event
 interface BootstrapEvent extends _sDtSitQ {
  app: App
 }
 type _szBGPIi = hook.Event
 interface TerminateEvent extends _szBGPIi {
  app: App
  isRestart: boolean
 }
 type _swKcgFZ = hook.Event
 interface BackupEvent extends _swKcgFZ {
  app: App
  context: context.Context
  name: string // the name of the backup to create/restore.
  exclude: Array<string> // list of dir entries to exclude from the backup create/restore.
 }
 type _szZVUoM = hook.Event
 interface ServeEvent extends _szZVUoM {
  app: App
  router?: router.Router<RequestEvent | undefined>
  server?: http.Server
  certManager?: any
  /**
   * Listener allow specifying a custom network listener.
   * 
   * Leave it nil to use the default net.Listen("tcp", e.Server.Addr).
   */
  listener: net.Listener
  /**
   * InstallerFunc is the "installer" function that is called after
   * successful server tcp bind but only if there is no explicit
   * superuser record created yet.
   * 
   * It runs in a separate goroutine and its default value is [apis.DefaultInstallerFunc].
   * 
   * It receives a system superuser record as argument that you can use to generate
   * a short-lived auth token (e.g. systemSuperuser.NewStaticAuthToken(30 * time.Minute))
   * and concatenate it as query param for your installer page
   * (if you are using the client-side SDKs, you can then load the
   * token with pb.authStore.save(token) and perform any Web API request
   * e.g. creating a new superuser).
   * 
   * Set it to nil if you want to skip the installer.
   */
  installerFunc: (app: App, systemSuperuser: Record, baseURL: string) => void
 }
 type _sTUpSQl = hook.Event&RequestEvent
 interface SettingsListRequestEvent extends _sTUpSQl {
  settings?: Settings
 }
 type _sFHUfjZ = hook.Event&RequestEvent
 interface SettingsUpdateRequestEvent extends _sFHUfjZ {
  oldSettings?: Settings
  newSettings?: Settings
 }
 type _sgiSZoY = hook.Event
 interface SettingsReloadEvent extends _sgiSZoY {
  app: App
 }
 type _slVrcYm = hook.Event
 interface MailerEvent extends _slVrcYm {
  app: App
  mailer: mailer.Mailer
  message?: mailer.Message
 }
 type _sKleLwL = MailerEvent&baseRecordEventData
 interface MailerRecordEvent extends _sKleLwL {
  meta: _TygojaDict
 }
 type _suVBtNg = hook.Event&baseModelEventData
 interface ModelEvent extends _suVBtNg {
  app: App
  context: context.Context
  /**
   * Could be any of the ModelEventType* constants, like:
   * - create
   * - update
   * - delete
   * - validate
   */
  type: string
 }
 type _sgCNkEC = ModelEvent
 interface ModelErrorEvent extends _sgCNkEC {
  error: Error
 }
 type _sHIZkgL = hook.Event&baseRecordEventData
 interface RecordEvent extends _sHIZkgL {
  app: App
  context: context.Context
  /**
   * Could be any of the ModelEventType* constants, like:
   * - create
   * - update
   * - delete
   * - validate
   */
  type: string
 }
 type _sBXRmvD = RecordEvent
 interface RecordErrorEvent extends _sBXRmvD {
  error: Error
 }
 type _sSWqrgB = hook.Event&baseCollectionEventData
 interface CollectionEvent extends _sSWqrgB {
  app: App
  context: context.Context
  /**
   * Could be any of the ModelEventType* constants, like:
   * - create
   * - update
   * - delete
   * - validate
   */
  type: string
 }
 type _stCPMoB = CollectionEvent
 interface CollectionErrorEvent extends _stCPMoB {
  error: Error
 }
 type _sbaGHtE = hook.Event&RequestEvent&baseRecordEventData
 interface FileTokenRequestEvent extends _sbaGHtE {
  token: string
 }
 type _scJJkRS = hook.Event&RequestEvent&baseCollectionEventData
 interface FileDownloadRequestEvent extends _scJJkRS {
  record?: Record
  fileField?: FileField
  servedPath: string
  servedName: string
  /**
   * ThumbError indicates the a thumb wasn't able to be generated
   * (e.g. because it didn't satisfy the support image formats or it timed out).
   * 
   * Note that PocketBase fallbacks to the original file in case of a thumb error,
   * but developers can check the field and provide their own custom thumb generation if necessary.
   */
  thumbError: Error
 }
 type _sGJOfcO = hook.Event&RequestEvent
 interface CollectionsListRequestEvent extends _sGJOfcO {
  collections: Array<(Collection | undefined)>
  result?: search.Result
 }
 type _sPlLzRU = hook.Event&RequestEvent
 interface CollectionsImportRequestEvent extends _sPlLzRU {
  collectionsData: Array<_TygojaDict>
  deleteMissing: boolean
 }
 type _sYEEfkx = hook.Event&RequestEvent&baseCollectionEventData
 interface CollectionRequestEvent extends _sYEEfkx {
 }
 type _sjQDOje = hook.Event&RequestEvent
 interface RealtimeConnectRequestEvent extends _sjQDOje {
  client: subscriptions.Client
  /**
   * note: modifying it after the connect has no effect
   */
  idleTimeout: time.Duration
 }
 type _shhfygx = hook.Event&RequestEvent
 interface RealtimeMessageEvent extends _shhfygx {
  client: subscriptions.Client
  message?: subscriptions.Message
 }
 type _sSbfVWu = hook.Event&RequestEvent
 interface RealtimeSubscribeRequestEvent extends _sSbfVWu {
  client: subscriptions.Client
  subscriptions: Array<string>
 }
 type _svZuxxU = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordsListRequestEvent extends _svZuxxU {
  /**
   * @todo consider removing and maybe add as generic to the search.Result?
   */
  records: Array<(Record | undefined)>
  result?: search.Result
 }
 type _snfWzTi = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordRequestEvent extends _snfWzTi {
  record?: Record
 }
 type _sieomkT = hook.Event&baseRecordEventData
 interface RecordEnrichEvent extends _sieomkT {
  app: App
  requestInfo?: RequestInfo
 }
 type _sbxczGT = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordCreateOTPRequestEvent extends _sbxczGT {
  record?: Record
  password: string
 }
 type _sPeVvjB = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordAuthWithOTPRequestEvent extends _sPeVvjB {
  record?: Record
  otp?: OTP
 }
 type _sHZzLzi = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordAuthRequestEvent extends _sHZzLzi {
  record?: Record
  token: string
  meta: any
  authMethod: string
 }
 type _shTKHaF = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordAuthWithPasswordRequestEvent extends _shTKHaF {
  record?: Record
  identity: string
  identityField: string
  password: string
 }
 type _sXSLTFx = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordAuthWithOAuth2RequestEvent extends _sXSLTFx {
  providerName: string
  providerClient: auth.Provider
  record?: Record
  oAuth2User?: auth.AuthUser
  createData: _TygojaDict
  isNewRecord: boolean
 }
 type _sbXnYFq = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordAuthRefreshRequestEvent extends _sbXnYFq {
  record?: Record
 }
 type _sJSWTdf = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordRequestPasswordResetRequestEvent extends _sJSWTdf {
  record?: Record
 }
 type _snngaqX = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordConfirmPasswordResetRequestEvent extends _snngaqX {
  record?: Record
 }
 type _sfYeStQ = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordRequestVerificationRequestEvent extends _sfYeStQ {
  record?: Record
 }
 type _sCSuVjQ = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordConfirmVerificationRequestEvent extends _sCSuVjQ {
  record?: Record
 }
 type _sTwJvQM = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordRequestEmailChangeRequestEvent extends _sTwJvQM {
  record?: Record
  newEmail: string
 }
 type _sMSEdOI = hook.Event&RequestEvent&baseCollectionEventData
 interface RecordConfirmEmailChangeRequestEvent extends _sMSEdOI {
  record?: Record
  newEmail: string
 }
 /**
  * ExternalAuth defines a Record proxy for working with the externalAuths collection.
  */
 type _sDwFWea = Record
 interface ExternalAuth extends _sDwFWea {
 }
 interface newExternalAuth {
  /**
   * NewExternalAuth instantiates and returns a new blank *ExternalAuth model.
   * 
   * Example usage:
   * 
   * ```
   * 	ea := core.NewExternalAuth(app)
   * 	ea.SetRecordRef(user.Id)
   * 	ea.SetCollectionRef(user.Collection().Id)
   * 	ea.SetProvider("google")
   * 	ea.SetProviderId("...")
   * 	app.Save(ea)
   * ```
   */
  (app: App): (ExternalAuth)
 }
 interface ExternalAuth {
  /**
   * PreValidate implements the [PreValidator] interface and checks
   * whether the proxy is properly loaded.
   */
  preValidate(ctx: context.Context, app: App): void
 }
 interface ExternalAuth {
  /**
   * ProxyRecord returns the proxied Record model.
   */
  proxyRecord(): (Record)
 }
 interface ExternalAuth {
  /**
   * SetProxyRecord loads the specified record model into the current proxy.
   */
  setProxyRecord(record: Record): void
 }
 interface ExternalAuth {
  /**
   * CollectionRef returns the "collectionRef" field value.
   */
  collectionRef(): string
 }
 interface ExternalAuth {
  /**
   * SetCollectionRef updates the "collectionRef" record field value.
   */
  setCollectionRef(collectionId: string): void
 }
 interface ExternalAuth {
  /**
   * RecordRef returns the "recordRef" record field value.
   */
  recordRef(): string
 }
 interface ExternalAuth {
  /**
   * SetRecordRef updates the "recordRef" record field value.
   */
  setRecordRef(recordId: string): void
 }
 interface ExternalAuth {
  /**
   * Provider returns the "provider" record field value.
   */
  provider(): string
 }
 interface ExternalAuth {
  /**
   * SetProvider updates the "provider" record field value.
   */
  setProvider(provider: string): void
 }
 interface ExternalAuth {
  /**
   * Provider returns the "providerId" record field value.
   */
  providerId(): string
 }
 interface ExternalAuth {
  /**
   * SetProvider updates the "providerId" record field value.
   */
  setProviderId(providerId: string): void
 }
 interface ExternalAuth {
  /**
   * Created returns the "created" record field value.
   */
  created(): types.DateTime
 }
 interface ExternalAuth {
  /**
   * Updated returns the "updated" record field value.
   */
  updated(): types.DateTime
 }
 interface BaseApp {
  /**
   * FindAllExternalAuthsByRecord returns all ExternalAuth models
   * linked to the provided auth record.
   */
  findAllExternalAuthsByRecord(authRecord: Record): Array<(ExternalAuth | undefined)>
 }
 interface BaseApp {
  /**
   * FindAllExternalAuthsByCollection returns all ExternalAuth models
   * linked to the provided auth collection.
   */
  findAllExternalAuthsByCollection(collection: Collection): Array<(ExternalAuth | undefined)>
 }
 interface BaseApp {
  /**
   * FindFirstExternalAuthByExpr returns the first available (the most recent created)
   * ExternalAuth model that satisfies the non-nil expression.
   */
  findFirstExternalAuthByExpr(expr: dbx.Expression): (ExternalAuth)
 }
 /**
  * FieldFactoryFunc defines a simple function to construct a specific Field instance.
  */
 interface FieldFactoryFunc {(): Field }
 /**
  * Field defines a common interface that all Collection fields should implement.
  */
 interface Field {
  [key:string]: any;
  /**
   * GetId returns the field id.
   */
  getId(): string
  /**
   * SetId changes the field id.
   */
  setId(id: string): void
  /**
   * GetName returns the field name.
   */
  getName(): string
  /**
   * SetName changes the field name.
   */
  setName(name: string): void
  /**
   * GetSystem returns the field system flag state.
   */
  getSystem(): boolean
  /**
   * SetSystem changes the field system flag state.
   */
  setSystem(system: boolean): void
  /**
   * GetHidden returns the field hidden flag state.
   */
  getHidden(): boolean
  /**
   * SetHidden changes the field hidden flag state.
   */
  setHidden(hidden: boolean): void
  /**
   * Type returns the unique type of the field.
   */
  type(): string
  /**
   * ColumnType returns the DB column definition of the field.
   */
  columnType(app: App): string
  /**
   * PrepareValue returns a properly formatted field value based on the provided raw one.
   * 
   * This method is also called on record construction to initialize its default field value.
   */
  prepareValue(record: Record, raw: any): any
  /**
   * ValidateSettings validates the current field value associated with the provided record.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
  /**
   * ValidateSettings validates the current field settings.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 /**
  * MaxBodySizeCalculator defines an optional field interface for
  * specifying the max size of a field value.
  */
 interface MaxBodySizeCalculator {
  [key:string]: any;
  /**
   * CalculateMaxBodySize returns the approximate max body size of a field value.
   */
  calculateMaxBodySize(): number
 }
 interface SetterFunc {(record: Record, raw: any): void }
 /**
  * SetterFinder defines a field interface for registering custom field value setters.
  */
 interface SetterFinder {
  [key:string]: any;
  /**
   * FindSetter returns a single field value setter function
   * by performing pattern-like field matching using the specified key.
   * 
   * The key is usually just the field name but it could also
   * contains "modifier" characters based on which you can perform custom set operations
   * (ex. "users+" could be mapped to a function that will append new user to the existing field value).
   * 
   * Return nil if you want to fallback to the default field value setter.
   */
  findSetter(key: string): SetterFunc
 }
 interface GetterFunc {(record: Record): any }
 /**
  * GetterFinder defines a field interface for registering custom field value getters.
  */
 interface GetterFinder {
  [key:string]: any;
  /**
   * FindGetter returns a single field value getter function
   * by performing pattern-like field matching using the specified key.
   * 
   * The key is usually just the field name but it could also
   * contains "modifier" characters based on which you can perform custom get operations
   * (ex. "description:excerpt" could be mapped to a function that will return an excerpt of the current field value).
   * 
   * Return nil if you want to fallback to the default field value setter.
   */
  findGetter(key: string): GetterFunc
 }
 /**
  * DriverValuer defines a Field interface for exporting and formatting
  * a field value for the database.
  */
 interface DriverValuer {
  [key:string]: any;
  /**
   * DriverValue exports a single field value for persistence in the database.
   */
  driverValue(record: Record): any
 }
 /**
  * MultiValuer defines a field interface that every multi-valued (eg. with MaxSelect) field has.
  */
 interface MultiValuer {
  [key:string]: any;
  /**
   * IsMultiple checks whether the field is configured to support multiple or single values.
   */
  isMultiple(): boolean
 }
 /**
  * RecordInterceptor defines a field interface for reacting to various
  * Record related operations (create, delete, validate, etc.).
  */
 interface RecordInterceptor {
  [key:string]: any;
  /**
   * Interceptor is invoked when a specific record action occurs
   * allowing you to perform extra validations and normalization
   * (ex. uploading or deleting files).
   * 
   * Note that users must call actionFunc() manually if they want to
   * execute the specific record action.
   */
  intercept(ctx: context.Context, app: App, record: Record, actionName: string, actionFunc: () => void): void
 }
 interface defaultFieldIdValidationRule {
  /**
   * DefaultFieldIdValidationRule performs base validation on a field id value.
   */
  (value: any): void
 }
 interface defaultFieldNameValidationRule {
  /**
   * DefaultFieldIdValidationRule performs base validation on a field name value.
   */
  (value: any): void
 }
 /**
  * AutodateField defines an "autodate" type field, aka.
  * field which datetime value could be auto set on record create/update.
  * 
  * This field is usually used for defining timestamp fields like "created" and "updated".
  * 
  * Requires either both or at least one of the OnCreate or OnUpdate options to be set.
  */
 interface AutodateField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * OnCreate auto sets the current datetime as field value on record create.
   */
  onCreate: boolean
  /**
   * OnUpdate auto sets the current datetime as field value on record update.
   */
  onUpdate: boolean
 }
 interface AutodateField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface AutodateField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface AutodateField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface AutodateField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface AutodateField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface AutodateField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface AutodateField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface AutodateField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface AutodateField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface AutodateField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface AutodateField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface AutodateField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface AutodateField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface AutodateField {
  /**
   * FindSetter implements the [SetterFinder] interface.
   */
  findSetter(key: string): SetterFunc
 }
 interface AutodateField {
  /**
   * Intercept implements the [RecordInterceptor] interface.
   */
  intercept(ctx: context.Context, app: App, record: Record, actionName: string, actionFunc: () => void): void
 }
 /**
  * BoolField defines "bool" type field to store a single true/false value.
  * 
  * The respective zero record field value is false.
  */
 interface BoolField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Required will require the field value to be always "true".
   */
  required: boolean
 }
 interface BoolField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface BoolField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface BoolField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface BoolField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface BoolField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface BoolField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface BoolField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface BoolField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface BoolField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface BoolField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface BoolField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface BoolField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface BoolField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 /**
  * DateField defines "date" type field to store a single [types.DateTime] value.
  * 
  * The respective zero record field value is the zero [types.DateTime].
  */
 interface DateField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Min specifies the min allowed field value.
   * 
   * Leave it empty to skip the validator.
   */
  min: types.DateTime
  /**
   * Max specifies the max allowed field value.
   * 
   * Leave it empty to skip the validator.
   */
  max: types.DateTime
  /**
   * Required will require the field value to be non-zero [types.DateTime].
   */
  required: boolean
 }
 interface DateField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface DateField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface DateField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface DateField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface DateField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface DateField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface DateField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface DateField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface DateField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface DateField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface DateField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface DateField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface DateField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 /**
  * EditorField defines "editor" type field to store HTML formatted text.
  * 
  * The respective zero record field value is empty string.
  */
 interface EditorField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * MaxSize specifies the maximum size of the allowed field value (in bytes and up to 2^53-1).
   * 
   * If zero, a default limit of ~5MB is applied.
   */
  maxSize: number
  /**
   * ConvertURLs is usually used to instruct the editor whether to
   * apply url conversion (eg. stripping the domain name in case the
   * urls are using the same domain as the one where the editor is loaded).
   * 
   * (see also https://www.tiny.cloud/docs/tinymce/6/url-handling/#convert_urls)
   */
  convertURLs: boolean
  /**
   * Required will require the field value to be non-empty string.
   */
  required: boolean
 }
 interface EditorField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface EditorField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface EditorField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface EditorField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface EditorField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface EditorField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface EditorField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface EditorField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface EditorField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface EditorField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface EditorField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface EditorField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface EditorField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface EditorField {
  /**
   * CalculateMaxBodySize implements the [MaxBodySizeCalculator] interface.
   */
  calculateMaxBodySize(): number
 }
 /**
  * EmailField defines "email" type field for storing a single email string address.
  * 
  * The respective zero record field value is empty string.
  */
 interface EmailField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * ExceptDomains will require the email domain to NOT be included in the listed ones.
   * 
   * This validator can be set only if OnlyDomains is empty.
   */
  exceptDomains: Array<string>
  /**
   * OnlyDomains will require the email domain to be included in the listed ones.
   * 
   * This validator can be set only if ExceptDomains is empty.
   */
  onlyDomains: Array<string>
  /**
   * Required will require the field value to be non-empty email string.
   */
  required: boolean
 }
 interface EmailField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface EmailField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface EmailField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface EmailField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface EmailField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface EmailField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface EmailField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface EmailField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface EmailField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface EmailField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface EmailField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface EmailField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface EmailField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 /**
  * FileField defines "file" type field for managing record file(s).
  * 
  * Only the file name is stored as part of the record value.
  * New files (aka. files to upload) are expected to be of *filesytem.File.
  * 
  * If MaxSelect is not set or <= 1, then the field value is expected to be a single record id.
  * 
  * If MaxSelect is > 1, then the field value is expected to be a slice of record ids.
  * 
  * The respective zero record field value is either empty string (single) or empty string slice (multiple).
  * 
  * ---
  * 
  * The following additional setter keys are available:
  * 
  * ```
  *   - "fieldName+" - append one or more files to the existing record one. For example:
  * 
  *     // []string{"old1.txt", "old2.txt", "new1_ajkvass.txt", "new2_klhfnwd.txt"}
  *     record.Set("documents+", []*filesystem.File{new1, new2})
  * 
  *   - "+fieldName" - prepend one or more files to the existing record one. For example:
  * 
  *     // []string{"new1_ajkvass.txt", "new2_klhfnwd.txt", "old1.txt", "old2.txt",}
  *     record.Set("+documents", []*filesystem.File{new1, new2})
  * 
  *   - "fieldName-" - subtract/delete one or more files from the existing record one. For example:
  * 
  *     // []string{"old2.txt",}
  *     record.Set("documents-", "old1.txt")
  * ```
  */
 interface FileField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * MaxSize specifies the maximum size of a single uploaded file (in bytes and up to 2^53-1).
   * 
   * If zero, a default limit of 5MB is applied.
   */
  maxSize: number
  /**
   * MaxSelect specifies the max allowed files.
   * 
   * For multiple files the value must be > 1, otherwise fallbacks to single (default).
   */
  maxSelect: number
  /**
   * MimeTypes specifies an optional list of the allowed file mime types.
   * 
   * Leave it empty to disable the validator.
   */
  mimeTypes: Array<string>
  /**
   * Thumbs specifies an optional list of the supported thumbs for image based files.
   * 
   * Each entry must be in one of the following formats:
   * 
   * ```
   *   - WxH  (eg. 100x300) - crop to WxH viewbox (from center)
   *   - WxHt (eg. 100x300t) - crop to WxH viewbox (from top)
   *   - WxHb (eg. 100x300b) - crop to WxH viewbox (from bottom)
   *   - WxHf (eg. 100x300f) - fit inside a WxH viewbox (without cropping)
   *   - 0xH  (eg. 0x300)    - resize to H height preserving the aspect ratio
   *   - Wx0  (eg. 100x0)    - resize to W width preserving the aspect ratio
   * ```
   */
  thumbs: Array<string>
  /**
   * Protected will require the users to provide a special file token to access the file.
   * 
   * Note that by default all files are publicly accessible.
   * 
   * For the majority of the cases this is fine because by default
   * all file names have random part appended to their name which
   * need to be known by the user before accessing the file.
   */
  protected: boolean
  /**
   * Required will require the field value to have at least one file.
   */
  required: boolean
 }
 interface FileField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface FileField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface FileField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface FileField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface FileField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface FileField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface FileField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface FileField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface FileField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface FileField {
  /**
   * IsMultiple implements MultiValuer interface and checks whether the
   * current field options support multiple values.
   */
  isMultiple(): boolean
 }
 interface FileField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface FileField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface FileField {
  /**
   * DriverValue implements the [DriverValuer] interface.
   */
  driverValue(record: Record): any
 }
 interface FileField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface FileField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface FileField {
  /**
   * CalculateMaxBodySize implements the [MaxBodySizeCalculator] interface.
   */
  calculateMaxBodySize(): number
 }
 interface FileField {
  /**
   * Intercept implements the [RecordInterceptor] interface.
   * 
   * note: files delete after records deletion is handled globally by the app FileManager hook
   */
  intercept(ctx: context.Context, app: App, record: Record, actionName: string, actionFunc: () => void): void
 }
 interface FileField {
  /**
   * FindGetter implements the [GetterFinder] interface.
   */
  findGetter(key: string): GetterFunc
 }
 interface FileField {
  /**
   * FindSetter implements the [SetterFinder] interface.
   */
  findSetter(key: string): SetterFunc
 }
 /**
  * GeoPointField defines "geoPoint" type field for storing latitude and longitude GPS coordinates.
  * 
  * You can set the record field value as [types.GeoPoint], map or serialized json object with lat-lon props.
  * The stored value is always converted to [types.GeoPoint].
  * Nil, empty map, empty bytes slice, etc. results in zero [types.GeoPoint].
  * 
  * Examples of updating a record's GeoPointField value programmatically:
  * 
  * ```
  * 	record.Set("location", types.GeoPoint{Lat: 123, Lon: 456})
  * 	record.Set("location", map[string]any{"lat":123, "lon":456})
  * 	record.Set("location", []byte(`{"lat":123, "lon":456}`)
  * ```
  */
 interface GeoPointField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Required will require the field coordinates to be non-zero (aka. not "Null Island").
   */
  required: boolean
 }
 interface GeoPointField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface GeoPointField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface GeoPointField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface GeoPointField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface GeoPointField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface GeoPointField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface GeoPointField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface GeoPointField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface GeoPointField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface GeoPointField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface GeoPointField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface GeoPointField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface GeoPointField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 /**
  * JSONField defines "json" type field for storing any serialized JSON value.
  * 
  * The respective zero record field value is the zero [types.JSONRaw].
  */
 interface JSONField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * MaxSize specifies the maximum size of the allowed field value (in bytes and up to 2^53-1).
   * 
   * If zero, a default limit of 1MB is applied.
   */
  maxSize: number
  /**
   * Required will require the field value to be non-empty JSON value
   * (aka. not "null", `""`, "[]", "{}").
   */
  required: boolean
 }
 interface JSONField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface JSONField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface JSONField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface JSONField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface JSONField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface JSONField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface JSONField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface JSONField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface JSONField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface JSONField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface JSONField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface JSONField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface JSONField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface JSONField {
  /**
   * CalculateMaxBodySize implements the [MaxBodySizeCalculator] interface.
   */
  calculateMaxBodySize(): number
 }
 /**
  * NumberField defines "number" type field for storing numeric (float64) value.
  * 
  * The respective zero record field value is 0.
  * 
  * The following additional setter keys are available:
  * 
  * ```
  *   - "fieldName+" - appends to the existing record value. For example:
  *     record.Set("total+", 5)
  *   - "fieldName-" - subtracts from the existing record value. For example:
  *     record.Set("total-", 5)
  * ```
  */
 interface NumberField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Min specifies the min allowed field value.
   * 
   * Leave it nil to skip the validator.
   */
  min?: number
  /**
   * Max specifies the max allowed field value.
   * 
   * Leave it nil to skip the validator.
   */
  max?: number
  /**
   * OnlyInt will require the field value to be integer.
   */
  onlyInt: boolean
  /**
   * Required will require the field value to be non-zero.
   */
  required: boolean
 }
 interface NumberField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface NumberField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface NumberField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface NumberField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface NumberField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface NumberField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface NumberField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface NumberField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface NumberField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface NumberField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface NumberField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface NumberField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface NumberField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface NumberField {
  /**
   * FindSetter implements the [SetterFinder] interface.
   */
  findSetter(key: string): SetterFunc
 }
 /**
  * PasswordField defines "password" type field for storing bcrypt hashed strings
  * (usually used only internally for the "password" auth collection system field).
  * 
  * If you want to set a direct bcrypt hash as record field value you can use the SetRaw method, for example:
  * 
  * ```
  * 	// generates a bcrypt hash of "123456" and set it as field value
  * 	// (record.GetString("password") returns the plain password until persisted, otherwise empty string)
  * 	record.Set("password", "123456")
  * 
  * 	// set directly a bcrypt hash of "123456" as field value
  * 	// (record.GetString("password") returns empty string)
  * 	record.SetRaw("password", "$2a$10$.5Elh8fgxypNUWhpUUr/xOa2sZm0VIaE0qWuGGl9otUfobb46T1Pq")
  * ```
  * 
  * The following additional getter keys are available:
  * 
  * ```
  *   - "fieldName:hash" - returns the bcrypt hash string of the record field value (if any). For example:
  *     record.GetString("password:hash")
  * ```
  */
 interface PasswordField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Pattern specifies an optional regex pattern to match against the field value.
   * 
   * Leave it empty to skip the pattern check.
   */
  pattern: string
  /**
   * Min specifies an optional required field string length.
   */
  min: number
  /**
   * Max specifies an optional required field string length.
   * 
   * If zero, fallback to max 71 bytes.
   */
  max: number
  /**
   * Cost specifies the cost/weight/iteration/etc. bcrypt factor.
   * 
   * If zero, fallback to [bcrypt.DefaultCost].
   * 
   * If explicitly set, must be between [bcrypt.MinCost] and [bcrypt.MaxCost].
   */
  cost: number
  /**
   * Required will require the field value to be non-empty string.
   */
  required: boolean
 }
 interface PasswordField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface PasswordField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface PasswordField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface PasswordField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface PasswordField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface PasswordField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface PasswordField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface PasswordField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface PasswordField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface PasswordField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface PasswordField {
  /**
   * DriverValue implements the [DriverValuer] interface.
   */
  driverValue(record: Record): any
 }
 interface PasswordField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface PasswordField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface PasswordField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface PasswordField {
  /**
   * Intercept implements the [RecordInterceptor] interface.
   */
  intercept(ctx: context.Context, app: App, record: Record, actionName: string, actionFunc: () => void): void
 }
 interface PasswordField {
  /**
   * FindGetter implements the [GetterFinder] interface.
   */
  findGetter(key: string): GetterFunc
 }
 interface PasswordField {
  /**
   * FindSetter implements the [SetterFinder] interface.
   */
  findSetter(key: string): SetterFunc
 }
 interface PasswordFieldValue {
  lastError: Error
  hash: string
  plain: string
 }
 interface PasswordFieldValue {
  validate(pass: string): boolean
 }
 /**
  * RelationField defines "relation" type field for storing single or
  * multiple collection record references.
  * 
  * Requires the CollectionId option to be set.
  * 
  * If MaxSelect is not set or <= 1, then the field value is expected to be a single record id.
  * 
  * If MaxSelect is > 1, then the field value is expected to be a slice of record ids.
  * 
  * The respective zero record field value is either empty string (single) or empty string slice (multiple).
  * 
  * ---
  * 
  * The following additional setter keys are available:
  * 
  * ```
  *   - "fieldName+" - append one or more values to the existing record one. For example:
  * 
  *     record.Set("categories+", []string{"new1", "new2"}) // []string{"old1", "old2", "new1", "new2"}
  * 
  *   - "+fieldName" - prepend one or more values to the existing record one. For example:
  * 
  *     record.Set("+categories", []string{"new1", "new2"}) // []string{"new1", "new2", "old1", "old2"}
  * 
  *   - "fieldName-" - subtract one or more values from the existing record one. For example:
  * 
  *     record.Set("categories-", "old1") // []string{"old2"}
  * ```
  */
 interface RelationField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * CollectionId is the id of the related collection.
   */
  collectionId: string
  /**
   * CascadeDelete indicates whether the root model should be deleted
   * in case of delete of all linked relations.
   */
  cascadeDelete: boolean
  /**
   * MinSelect indicates the min number of allowed relation records
   * that could be linked to the main model.
   * 
   * No min limit is applied if it is zero or negative value.
   */
  minSelect: number
  /**
   * MaxSelect indicates the max number of allowed relation records
   * that could be linked to the main model.
   * 
   * For multiple select the value must be > 1, otherwise fallbacks to single (default).
   * 
   * If MinSelect is set, MaxSelect must be at least >= MinSelect.
   */
  maxSelect: number
  /**
   * Required will require the field value to be non-empty.
   */
  required: boolean
 }
 interface RelationField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface RelationField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface RelationField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface RelationField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface RelationField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface RelationField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface RelationField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface RelationField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface RelationField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface RelationField {
  /**
   * IsMultiple implements [MultiValuer] interface and checks whether the
   * current field options support multiple values.
   */
  isMultiple(): boolean
 }
 interface RelationField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface RelationField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface RelationField {
  /**
   * DriverValue implements the [DriverValuer] interface.
   */
  driverValue(record: Record): any
 }
 interface RelationField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface RelationField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface RelationField {
  /**
   * FindSetter implements [SetterFinder] interface method.
   */
  findSetter(key: string): SetterFunc
 }
 /**
  * SelectField defines "select" type field for storing single or
  * multiple string values from a predefined list.
  * 
  * Requires the Values option to be set.
  * 
  * If MaxSelect is not set or <= 1, then the field value is expected to be a single Values element.
  * 
  * If MaxSelect is > 1, then the field value is expected to be a subset of Values slice.
  * 
  * The respective zero record field value is either empty string (single) or empty string slice (multiple).
  * 
  * ---
  * 
  * The following additional setter keys are available:
  * 
  * ```
  *   - "fieldName+" - append one or more values to the existing record one. For example:
  * 
  *     record.Set("roles+", []string{"new1", "new2"}) // []string{"old1", "old2", "new1", "new2"}
  * 
  *   - "+fieldName" - prepend one or more values to the existing record one. For example:
  * 
  *     record.Set("+roles", []string{"new1", "new2"}) // []string{"new1", "new2", "old1", "old2"}
  * 
  *   - "fieldName-" - subtract one or more values from the existing record one. For example:
  * 
  *     record.Set("roles-", "old1") // []string{"old2"}
  * ```
  */
 interface SelectField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Values specifies the list of accepted values.
   */
  values: Array<string>
  /**
   * MaxSelect specifies the max allowed selected values.
   * 
   * For multiple select the value must be > 1, otherwise fallbacks to single (default).
   */
  maxSelect: number
  /**
   * Required will require the field value to be non-empty.
   */
  required: boolean
 }
 interface SelectField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface SelectField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface SelectField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface SelectField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface SelectField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface SelectField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface SelectField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface SelectField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface SelectField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface SelectField {
  /**
   * IsMultiple implements [MultiValuer] interface and checks whether the
   * current field options support multiple values.
   */
  isMultiple(): boolean
 }
 interface SelectField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface SelectField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface SelectField {
  /**
   * DriverValue implements the [DriverValuer] interface.
   */
  driverValue(record: Record): any
 }
 interface SelectField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface SelectField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface SelectField {
  /**
   * FindSetter implements the [SetterFinder] interface.
   */
  findSetter(key: string): SetterFunc
 }
 /**
  * TextField defines "text" type field for storing any string value.
  * 
  * The respective zero record field value is empty string.
  * 
  * The following additional setter keys are available:
  * 
  * - "fieldName:autogenerate" - autogenerate field value if AutogeneratePattern is set. For example:
  * 
  * ```
  * 	record.Set("slug:autogenerate", "") // [random value]
  * 	record.Set("slug:autogenerate", "abc-") // abc-[random value]
  * ```
  */
 interface TextField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * Min specifies the minimum required string characters.
   * 
   * if zero value, no min limit is applied.
   */
  min: number
  /**
   * Max specifies the maximum allowed string characters.
   * 
   * If zero, a default limit of 5000 is applied.
   */
  max: number
  /**
   * Pattern specifies an optional regex pattern to match against the field value.
   * 
   * Leave it empty to skip the pattern check.
   */
  pattern: string
  /**
   * AutogeneratePattern specifies an optional regex pattern that could
   * be used to generate random string from it and set it automatically
   * on record create if no explicit value is set or when the `:autogenerate` modifier is used.
   * 
   * Note: the generated value still needs to satisfy min, max, pattern (if set)
   */
  autogeneratePattern: string
  /**
   * Required will require the field value to be non-empty string.
   */
  required: boolean
  /**
   * PrimaryKey will mark the field as primary key.
   * 
   * A single collection can have only 1 field marked as primary key.
   */
  primaryKey: boolean
 }
 interface TextField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface TextField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface TextField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface TextField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface TextField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface TextField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface TextField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface TextField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface TextField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface TextField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface TextField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface TextField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface TextField {
  /**
   * ValidatePlainValue validates the provided string against the field options.
   */
  validatePlainValue(value: string): void
 }
 interface TextField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface TextField {
  /**
   * Intercept implements the [RecordInterceptor] interface.
   */
  intercept(ctx: context.Context, app: App, record: Record, actionName: string, actionFunc: () => void): void
 }
 interface TextField {
  /**
   * FindSetter implements the [SetterFinder] interface.
   */
  findSetter(key: string): SetterFunc
 }
 /**
  * URLField defines "url" type field for storing a single URL string value.
  * 
  * The respective zero record field value is empty string.
  */
 interface URLField {
  /**
   * Name (required) is the unique name of the field.
   */
  name: string
  /**
   * Id is the unique stable field identifier.
   * 
   * It is automatically generated from the name when adding to a collection FieldsList.
   */
  id: string
  /**
   * System prevents the renaming and removal of the field.
   */
  system: boolean
  /**
   * Hidden hides the field from the API response.
   */
  hidden: boolean
  /**
   * Presentable hints the Dashboard UI to use the underlying
   * field record value in the relation preview label.
   */
  presentable: boolean
  /**
   * ExceptDomains will require the URL domain to NOT be included in the listed ones.
   * 
   * This validator can be set only if OnlyDomains is empty.
   */
  exceptDomains: Array<string>
  /**
   * OnlyDomains will require the URL domain to be included in the listed ones.
   * 
   * This validator can be set only if ExceptDomains is empty.
   */
  onlyDomains: Array<string>
  /**
   * Required will require the field value to be non-empty URL string.
   */
  required: boolean
 }
 interface URLField {
  /**
   * Type implements [Field.Type] interface method.
   */
  type(): string
 }
 interface URLField {
  /**
   * GetId implements [Field.GetId] interface method.
   */
  getId(): string
 }
 interface URLField {
  /**
   * SetId implements [Field.SetId] interface method.
   */
  setId(id: string): void
 }
 interface URLField {
  /**
   * GetName implements [Field.GetName] interface method.
   */
  getName(): string
 }
 interface URLField {
  /**
   * SetName implements [Field.SetName] interface method.
   */
  setName(name: string): void
 }
 interface URLField {
  /**
   * GetSystem implements [Field.GetSystem] interface method.
   */
  getSystem(): boolean
 }
 interface URLField {
  /**
   * SetSystem implements [Field.SetSystem] interface method.
   */
  setSystem(system: boolean): void
 }
 interface URLField {
  /**
   * GetHidden implements [Field.GetHidden] interface method.
   */
  getHidden(): boolean
 }
 interface URLField {
  /**
   * SetHidden implements [Field.SetHidden] interface method.
   */
  setHidden(hidden: boolean): void
 }
 interface URLField {
  /**
   * ColumnType implements [Field.ColumnType] interface method.
   */
  columnType(app: App): string
 }
 interface URLField {
  /**
   * PrepareValue implements [Field.PrepareValue] interface method.
   */
  prepareValue(record: Record, raw: any): any
 }
 interface URLField {
  /**
   * ValidateValue implements [Field.ValidateValue] interface method.
   */
  validateValue(ctx: context.Context, app: App, record: Record): void
 }
 interface URLField {
  /**
   * ValidateSettings implements [Field.ValidateSettings] interface method.
   */
  validateSettings(ctx: context.Context, app: App, collection: Collection): void
 }
 interface newFieldsList {
  /**
   * NewFieldsList creates a new FieldsList instance with the provided fields.
   */
  (...fields: Field[]): FieldsList
 }
 /**
  * FieldsList defines a Collection slice of fields.
  */
 interface FieldsList extends Array<Field>{}
 interface FieldsList {
  /**
   * Clone creates a deep clone of the current list.
   */
  clone(): FieldsList
 }
 interface FieldsList {
  /**
   * FieldNames returns a slice with the name of all list fields.
   */
  fieldNames(): Array<string>
 }
 interface FieldsList {
  /**
   * AsMap returns a map with all registered list field.
   * The returned map is indexed with each field name.
   */
  asMap(): _TygojaDict
 }
 interface FieldsList {
  /**
   * GetById returns a single field by its id.
   */
  getById(fieldId: string): Field
 }
 interface FieldsList {
  /**
   * GetByName returns a single field by its name.
   */
  getByName(fieldName: string): Field
 }
 interface FieldsList {
  /**
   * RemoveById removes a single field by its id.
   * 
   * This method does nothing if field with the specified id doesn't exist.
   */
  removeById(fieldId: string): void
 }
 interface FieldsList {
  /**
   * RemoveByName removes a single field by its name.
   * 
   * This method does nothing if field with the specified name doesn't exist.
   */
  removeByName(fieldName: string): void
 }
 interface FieldsList {
  /**
   * Add adds one or more fields to the current list.
   * 
   * By default this method will try to REPLACE existing fields with
   * the new ones by their id or by their name if the new field doesn't have an explicit id.
   * 
   * If no matching existing field is found, it will APPEND the field to the end of the list.
   * 
   * In all cases, if any of the new fields don't have an explicit id it will auto generate a default one for them
   * (the id value doesn't really matter and it is mostly used as a stable identifier in case of a field rename).
   */
  add(...fields: Field[]): void
 }
 interface FieldsList {
  /**
   * AddAt is the same as Add but insert/move the fields at the specific position.
   * 
   * If pos < 0, then this method acts the same as calling Add.
   * 
   * If pos > FieldsList total items, then the specified fields are inserted/moved at the end of the list.
   */
  addAt(pos: number, ...fields: Field[]): void
 }
 interface FieldsList {
  /**
   * AddMarshaledJSON parses the provided raw json data and adds the
   * found fields into the current list (following the same rule as the Add method).
   * 
   * The rawJSON argument could be one of:
   * ```
   *   - serialized array of field objects
   *   - single field object.
   * ```
   * 
   * Example:
   * 
   * ```
   * 	l.AddMarshaledJSON([]byte{`{"type":"text", name: "test"}`})
   * 	l.AddMarshaledJSON([]byte{`[{"type":"text", name: "test1"}, {"type":"text", name: "test2"}]`})
   * ```
   */
  addMarshaledJSON(rawJSON: string|Array<number>): void
 }
 interface FieldsList {
  /**
   * AddMarshaledJSONAt is the same as AddMarshaledJSON but insert/move the fields at the specific position.
   * 
   * If pos < 0, then this method acts the same as calling AddMarshaledJSON.
   * 
   * If pos > FieldsList total items, then the specified fields are inserted/moved at the end of the list.
   */
  addMarshaledJSONAt(pos: number, rawJSON: string|Array<number>): void
 }
 interface FieldsList {
  /**
   * String returns the string representation of the current list.
   */
  string(): string
 }
 interface onlyFieldType {
  type: string
 }
 type _sdaAYmS = Field
 interface fieldWithType extends _sdaAYmS {
  type: string
 }
 interface fieldWithType {
  unmarshalJSON(data: string|Array<number>): void
 }
 interface FieldsList {
  /**
   * UnmarshalJSON implements [json.Unmarshaler] and
   * loads the provided json data into the current FieldsList.
   */
  unmarshalJSON(data: string|Array<number>): void
 }
 interface FieldsList {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface FieldsList {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface FieldsList {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current FieldsList instance.
   */
  scan(value: any): void
 }
 type _sgJNjtq = BaseModel
 interface Log extends _sgJNjtq {
  created: types.DateTime
  data: types.JSONMap<any>
  message: string
  level: number
 }
 interface Log {
  tableName(): string
 }
 interface BaseApp {
  /**
   * LogQuery returns a new Log select query.
   */
  logQuery(): (dbx.SelectQuery)
 }
 interface BaseApp {
  /**
   * FindLogById finds a single Log entry by its id.
   */
  findLogById(id: string): (Log)
 }
 /**
  * LogsStatsItem defines the total number of logs for a specific time period.
  */
 interface LogsStatsItem {
  date: types.DateTime
  total: number
 }
 interface BaseApp {
  /**
   * LogsStats returns hourly grouped logs statistics.
   */
  logsStats(expr: dbx.Expression): Array<(LogsStatsItem | undefined)>
 }
 interface BaseApp {
  /**
   * DeleteOldLogs delete all logs that are created before createdBefore.
   * 
   * For better performance the logs delete is executed as plain SQL statement,
   * aka. no delete model hook events will be fired.
   */
  deleteOldLogs(createdBefore: time.Time): void
 }
 /**
  * MFA defines a Record proxy for working with the mfas collection.
  */
 type _sPqobze = Record
 interface MFA extends _sPqobze {
 }
 interface newMFA {
  /**
   * NewMFA instantiates and returns a new blank *MFA model.
   * 
   * Example usage:
   * 
   * ```
   * 	mfa := core.NewMFA(app)
   * 	mfa.SetRecordRef(user.Id)
   * 	mfa.SetCollectionRef(user.Collection().Id)
   * 	mfa.SetMethod(core.MFAMethodPassword)
   * 	app.Save(mfa)
   * ```
   */
  (app: App): (MFA)
 }
 interface MFA {
  /**
   * PreValidate implements the [PreValidator] interface and checks
   * whether the proxy is properly loaded.
   */
  preValidate(ctx: context.Context, app: App): void
 }
 interface MFA {
  /**
   * ProxyRecord returns the proxied Record model.
   */
  proxyRecord(): (Record)
 }
 interface MFA {
  /**
   * SetProxyRecord loads the specified record model into the current proxy.
   */
  setProxyRecord(record: Record): void
 }
 interface MFA {
  /**
   * CollectionRef returns the "collectionRef" field value.
   */
  collectionRef(): string
 }
 interface MFA {
  /**
   * SetCollectionRef updates the "collectionRef" record field value.
   */
  setCollectionRef(collectionId: string): void
 }
 interface MFA {
  /**
   * RecordRef returns the "recordRef" record field value.
   */
  recordRef(): string
 }
 interface MFA {
  /**
   * SetRecordRef updates the "recordRef" record field value.
   */
  setRecordRef(recordId: string): void
 }
 interface MFA {
  /**
   * Method returns the "method" record field value.
   */
  method(): string
 }
 interface MFA {
  /**
   * SetMethod updates the "method" record field value.
   */
  setMethod(method: string): void
 }
 interface MFA {
  /**
   * Created returns the "created" record field value.
   */
  created(): types.DateTime
 }
 interface MFA {
  /**
   * Updated returns the "updated" record field value.
   */
  updated(): types.DateTime
 }
 interface MFA {
  /**
   * HasExpired checks if the mfa is expired, aka. whether it has been
   * more than maxElapsed time since its creation.
   */
  hasExpired(maxElapsed: time.Duration): boolean
 }
 interface BaseApp {
  /**
   * FindAllMFAsByRecord returns all MFA models linked to the provided auth record.
   */
  findAllMFAsByRecord(authRecord: Record): Array<(MFA | undefined)>
 }
 interface BaseApp {
  /**
   * FindAllMFAsByCollection returns all MFA models linked to the provided collection.
   */
  findAllMFAsByCollection(collection: Collection): Array<(MFA | undefined)>
 }
 interface BaseApp {
  /**
   * FindMFAById returns a single MFA model by its id.
   */
  findMFAById(id: string): (MFA)
 }
 interface BaseApp {
  /**
   * DeleteAllMFAsByRecord deletes all MFA models associated with the provided record.
   * 
   * Returns a combined error with the failed deletes.
   */
  deleteAllMFAsByRecord(authRecord: Record): void
 }
 interface BaseApp {
  /**
   * DeleteExpiredMFAs deletes the expired MFAs for all auth collections.
   */
  deleteExpiredMFAs(): void
 }
 interface Migration {
  up: (txApp: App) => void
  down: (txApp: App) => void
  file: string
  reapplyCondition: (txApp: App, runner: MigrationsRunner, fileName: string) => boolean
 }
 /**
  * MigrationsList defines a list with migration definitions
  */
 interface MigrationsList {
 }
 interface MigrationsList {
  /**
   * Item returns a single migration from the list by its index.
   */
  item(index: number): (Migration)
 }
 interface MigrationsList {
  /**
   * Items returns the internal migrations list slice.
   */
  items(): Array<(Migration | undefined)>
 }
 interface MigrationsList {
  /**
   * Copy copies all provided list migrations into the current one.
   */
  copy(list: MigrationsList): void
 }
 interface MigrationsList {
  /**
   * Add adds adds an existing migration definition to the list.
   * 
   * If m.File is not provided, it will try to get the name from its .go file.
   * 
   * The list will be sorted automatically based on the migrations file name.
   */
  add(m: Migration): void
 }
 interface MigrationsList {
  /**
   * Register adds new migration definition to the list.
   * 
   * If optFilename is not provided, it will try to get the name from its .go file.
   * 
   * The list will be sorted automatically based on the migrations file name.
   */
  register(up: (txApp: App) => void, down: (txApp: App) => void, ...optFilename: string[]): void
 }
 /**
  * MigrationsRunner defines a simple struct for managing the execution of db migrations.
  */
 interface MigrationsRunner {
 }
 interface newMigrationsRunner {
  /**
   * NewMigrationsRunner creates and initializes a new db migrations MigrationsRunner instance.
   */
  (app: App, migrationsList: MigrationsList): (MigrationsRunner)
 }
 interface MigrationsRunner {
  /**
   * Run interactively executes the current runner with the provided args.
   * 
   * The following commands are supported:
   * - up           - applies all migrations
   * - down [n]     - reverts the last n (default 1) applied migrations
   * - history-sync - syncs the migrations table with the runner's migrations list
   */
  run(...args: string[]): void
 }
 interface MigrationsRunner {
  /**
   * Up executes all unapplied migrations for the provided runner.
   * 
   * On success returns list with the applied migrations file names.
   */
  up(): Array<string>
 }
 interface MigrationsRunner {
  /**
   * Down reverts the last `toRevertCount` applied migrations
   * (in the order they were applied).
   * 
   * On success returns list with the reverted migrations file names.
   */
  down(toRevertCount: number): Array<string>
 }
 interface MigrationsRunner {
  /**
   * RemoveMissingAppliedMigrations removes the db entries of all applied migrations
   * that are not listed in the runner's migrations list.
   */
  removeMissingAppliedMigrations(): void
 }
 /**
  * OTP defines a Record proxy for working with the otps collection.
  */
 type _sZgdqnw = Record
 interface OTP extends _sZgdqnw {
 }
 interface newOTP {
  /**
   * NewOTP instantiates and returns a new blank *OTP model.
   * 
   * Example usage:
   * 
   * ```
   * 	otp := core.NewOTP(app)
   * 	otp.SetRecordRef(user.Id)
   * 	otp.SetCollectionRef(user.Collection().Id)
   * 	otp.SetPassword(security.RandomStringWithAlphabet(6, "1234567890"))
   * 	app.Save(otp)
   * ```
   */
  (app: App): (OTP)
 }
 interface OTP {
  /**
   * PreValidate implements the [PreValidator] interface and checks
   * whether the proxy is properly loaded.
   */
  preValidate(ctx: context.Context, app: App): void
 }
 interface OTP {
  /**
   * ProxyRecord returns the proxied Record model.
   */
  proxyRecord(): (Record)
 }
 interface OTP {
  /**
   * SetProxyRecord loads the specified record model into the current proxy.
   */
  setProxyRecord(record: Record): void
 }
 interface OTP {
  /**
   * CollectionRef returns the "collectionRef" field value.
   */
  collectionRef(): string
 }
 interface OTP {
  /**
   * SetCollectionRef updates the "collectionRef" record field value.
   */
  setCollectionRef(collectionId: string): void
 }
 interface OTP {
  /**
   * RecordRef returns the "recordRef" record field value.
   */
  recordRef(): string
 }
 interface OTP {
  /**
   * SetRecordRef updates the "recordRef" record field value.
   */
  setRecordRef(recordId: string): void
 }
 interface OTP {
  /**
   * SentTo returns the "sentTo" record field value.
   * 
   * It could be any string value (email, phone, message app id, etc.)
   * and usually is used as part of the auth flow to update the verified
   * user state in case for example the sentTo value matches with the user record email.
   */
  sentTo(): string
 }
 interface OTP {
  /**
   * SetSentTo updates the "sentTo" record field value.
   */
  setSentTo(val: string): void
 }
 interface OTP {
  /**
   * Created returns the "created" record field value.
   */
  created(): types.DateTime
 }
 interface OTP {
  /**
   * Updated returns the "updated" record field value.
   */
  updated(): types.DateTime
 }
 interface OTP {
  /**
   * HasExpired checks if the otp is expired, aka. whether it has been
   * more than maxElapsed time since its creation.
   */
  hasExpired(maxElapsed: time.Duration): boolean
 }
 interface BaseApp {
  /**
   * FindAllOTPsByRecord returns all OTP models linked to the provided auth record.
   */
  findAllOTPsByRecord(authRecord: Record): Array<(OTP | undefined)>
 }
 interface BaseApp {
  /**
   * FindAllOTPsByCollection returns all OTP models linked to the provided collection.
   */
  findAllOTPsByCollection(collection: Collection): Array<(OTP | undefined)>
 }
 interface BaseApp {
  /**
   * FindOTPById returns a single OTP model by its id.
   */
  findOTPById(id: string): (OTP)
 }
 interface BaseApp {
  /**
   * DeleteAllOTPsByRecord deletes all OTP models associated with the provided record.
   * 
   * Returns a combined error with the failed deletes.
   */
  deleteAllOTPsByRecord(authRecord: Record): void
 }
 interface BaseApp {
  /**
   * DeleteExpiredOTPs deletes the expired OTPs for all auth collections.
   */
  deleteExpiredOTPs(): void
 }
 /**
  * RecordFieldResolver defines a custom search resolver struct for
  * managing Record model search fields.
  * 
  * Usually used together with `search.Provider`.
  * Example:
  * 
  * ```
  * 	resolver := resolvers.NewRecordFieldResolver(
  * 	    app,
  * 	    myCollection,
  * 	    &models.RequestInfo{...},
  * 	    true,
  * 	)
  * 	provider := search.NewProvider(resolver)
  * 	...
  * ```
  */
 interface RecordFieldResolver {
 }
 interface RecordFieldResolver {
  /**
   * AllowedFields returns a copy of the resolver's allowed fields.
   */
  allowedFields(): Array<string>
 }
 interface RecordFieldResolver {
  /**
   * SetAllowedFields replaces the resolver's allowed fields with the new ones.
   */
  setAllowedFields(newAllowedFields: Array<string>): void
 }
 interface RecordFieldResolver {
  /**
   * AllowHiddenFields returns whether the current resolver allows filtering hidden fields.
   */
  allowHiddenFields(): boolean
 }
 interface RecordFieldResolver {
  /**
   * SetAllowHiddenFields enables or disables hidden fields filtering.
   */
  setAllowHiddenFields(allowHiddenFields: boolean): void
 }
 interface newRecordFieldResolver {
  /**
   * NewRecordFieldResolver creates and initializes a new `RecordFieldResolver`.
   */
  (app: App, baseCollection: Collection, requestInfo: RequestInfo, allowHiddenFields: boolean): (RecordFieldResolver)
 }
 interface RecordFieldResolver {
  /**
   * @todo think of a better a way how to call it automatically after BuildExpr
   * 
   * UpdateQuery implements `search.FieldResolver` interface.
   * 
   * Conditionally updates the provided search query based on the
   * resolved fields (eg. dynamically joining relations).
   */
  updateQuery(query: dbx.SelectQuery): void
 }
 interface RecordFieldResolver {
  /**
   * Resolve implements `search.FieldResolver` interface.
   * 
   * Example of some resolvable fieldName formats:
   * 
   * ```
   * 	id
   * 	someSelect.each
   * 	project.screen.status
   * 	screen.project_via_prototype.name
   * 	@request.context
   * 	@request.method
   * 	@request.query.filter
   * 	@request.headers.x_token
   * 	@request.auth.someRelation.name
   * 	@request.body.someRelation.name
   * 	@request.body.someField
   * 	@request.body.someSelect:each
   * 	@request.body.someField:isset
   * 	@collection.product.name
   * ```
   */
  resolve(fieldName: string): (search.ResolverResult)
 }
 interface mapExtractor {
  [key:string]: any;
  asMap(): _TygojaDict
 }
 /**
  * join defines the specification for a single SQL JOIN clause.
  */
 interface join {
 }
 /**
  * multiMatchSubquery defines a record multi-match subquery expression.
  */
 interface multiMatchSubquery {
 }
 interface multiMatchSubquery {
  /**
   * Build converts the expression into a SQL fragment.
   * 
   * Implements [dbx.Expression] interface.
   */
  build(db: dbx.DB, params: dbx.Params): string
 }
 interface runner {
 }
 type _skFHppE = BaseModel
 interface Record extends _skFHppE {
 }
 interface newRecord {
  /**
   * NewRecord initializes a new empty Record model.
   */
  (collection: Collection): (Record)
 }
 interface Record {
  /**
   * Collection returns the Collection model associated with the current Record model.
   * 
   * NB! The returned collection is only for read purposes and it shouldn't be modified
   * because it could have unintended side-effects on other Record models from the same collection.
   */
  collection(): (Collection)
 }
 interface Record {
  /**
   * TableName returns the table name associated with the current Record model.
   */
  tableName(): string
 }
 interface Record {
  /**
   * PostScan implements the [dbx.PostScanner] interface.
   * 
   * It essentially refreshes/updates the current Record original state
   * as if the model was fetched from the databases for the first time.
   * 
   * Or in other words, it means that m.Original().FieldsData() will have
   * the same values as m.Record().FieldsData().
   */
  postScan(): void
 }
 interface Record {
  /**
   * HookTags returns the hook tags associated with the current record.
   */
  hookTags(): Array<string>
 }
 interface Record {
  /**
   * BaseFilesPath returns the storage dir path used by the record.
   */
  baseFilesPath(): string
 }
 interface Record {
  /**
   * Original returns a shallow copy of the current record model populated
   * with its ORIGINAL db data state (aka. right after PostScan())
   * and everything else reset to the defaults.
   * 
   * If record was created using NewRecord() the original will be always
   * a blank record (until PostScan() is invoked).
   */
  original(): (Record)
 }
 interface Record {
  /**
   * Fresh returns a shallow copy of the current record model populated
   * with its LATEST data state and everything else reset to the defaults
   * (aka. no expand, no unknown fields and with default visibility flags).
   */
  fresh(): (Record)
 }
 interface Record {
  /**
   * Clone returns a shallow copy of the current record model with all of
   * its collection and unknown fields data, expand and flags copied.
   * 
   * use [Record.Fresh()] instead if you want a copy with only the latest
   * collection fields data and everything else reset to the defaults.
   */
  clone(): (Record)
 }
 interface Record {
  /**
   * Expand returns a shallow copy of the current Record model expand data (if any).
   */
  expand(): _TygojaDict
 }
 interface Record {
  /**
   * SetExpand replaces the current Record's expand with the provided expand arg data (shallow copied).
   */
  setExpand(expand: _TygojaDict): void
 }
 interface Record {
  /**
   * MergeExpand merges recursively the provided expand data into
   * the current model's expand (if any).
   * 
   * Note that if an expanded prop with the same key is a slice (old or new expand)
   * then both old and new records will be merged into a new slice (aka. a :merge: [b,c] => [a,b,c]).
   * Otherwise the "old" expanded record will be replace with the "new" one (aka. a :merge: aNew => aNew).
   */
  mergeExpand(expand: _TygojaDict): void
 }
 interface Record {
  /**
   * FieldsData returns a shallow copy ONLY of the collection's fields record's data.
   */
  fieldsData(): _TygojaDict
 }
 interface Record {
  /**
   * CustomData returns a shallow copy ONLY of the custom record fields data,
   * aka. fields that are neither defined by the collection, nor special system ones.
   * 
   * Note that custom fields prefixed with "@pbInternal" are always skipped.
   */
  customData(): _TygojaDict
 }
 interface Record {
  /**
   * WithCustomData toggles the export/serialization of custom data fields
   * (false by default).
   */
  withCustomData(state: boolean): (Record)
 }
 interface Record {
  /**
   * IgnoreEmailVisibility toggles the flag to ignore the auth record email visibility check.
   */
  ignoreEmailVisibility(state: boolean): (Record)
 }
 interface Record {
  /**
   * IgnoreUnchangedFields toggles the flag to ignore the unchanged fields
   * from the DB export for the UPDATE SQL query.
   * 
   * This could be used if you want to save only the record fields that you've changed
   * without overwrite other untouched fields in case of concurrent update.
   * 
   * Note that the fields change comparison is based on the current fields against m.Original()
   * (aka. if you have performed save on the same Record instance multiple times you may have to refetch it,
   * so that m.Original() could reflect the last saved change).
   */
  ignoreUnchangedFields(state: boolean): (Record)
 }
 interface Record {
  /**
   * Set sets the provided key-value data pair into the current Record
   * model directly as it is WITHOUT NORMALIZATIONS.
   * 
   * See also [Record.Set].
   */
  setRaw(key: string, value: any): void
 }
 interface Record {
  /**
   * SetIfFieldExists sets the provided key-value data pair into the current Record model
   * ONLY if key is existing Collection field name/modifier.
   * 
   * This method does nothing if key is not a known Collection field name/modifier.
   * 
   * On success returns the matched Field, otherwise - nil.
   * 
   * To set any key-value, including custom/unknown fields, use the [Record.Set] method.
   */
  setIfFieldExists(key: string, value: any): Field
 }
 interface Record {
  /**
   * Set sets the provided key-value data pair into the current Record model.
   * 
   * If the record collection has field with name matching the provided "key",
   * the value will be further normalized according to the field setter(s).
   */
  set(key: string, value: any): void
 }
 interface Record {
  getRaw(key: string): any
 }
 interface Record {
  /**
   * Get returns a normalized single record model data value for "key".
   */
  get(key: string): any
 }
 interface Record {
  /**
   * Load bulk loads the provided data into the current Record model.
   */
  load(data: _TygojaDict): void
 }
 interface Record {
  /**
   * GetBool returns the data value for "key" as a bool.
   */
  getBool(key: string): boolean
 }
 interface Record {
  /**
   * GetString returns the data value for "key" as a string.
   */
  getString(key: string): string
 }
 interface Record {
  /**
   * GetInt returns the data value for "key" as an int.
   */
  getInt(key: string): number
 }
 interface Record {
  /**
   * GetFloat returns the data value for "key" as a float64.
   */
  getFloat(key: string): number
 }
 interface Record {
  /**
   * GetDateTime returns the data value for "key" as a DateTime instance.
   */
  getDateTime(key: string): types.DateTime
 }
 interface Record {
  /**
   * GetGeoPoint returns the data value for "key" as a GeoPoint instance.
   */
  getGeoPoint(key: string): types.GeoPoint
 }
 interface Record {
  /**
   * GetStringSlice returns the data value for "key" as a slice of non-zero unique strings.
   */
  getStringSlice(key: string): Array<string>
 }
 interface Record {
  /**
   * GetUnsavedFiles returns the uploaded files for the provided "file" field key,
   * (aka. the current [*filesytem.File] values) so that you can apply further
   * validations or modifications (including changing the file name or content before persisting).
   * 
   * Example:
   * 
   * ```
   * 	files := record.GetUnsavedFiles("documents")
   * 	for _, f := range files {
   * 	    f.Name = "doc_" + f.Name // add a prefix to each file name
   * 	}
   * 	app.Save(record) // the files are pointers so the applied changes will transparently reflect on the record value
   * ```
   */
  getUnsavedFiles(key: string): Array<(filesystem.File | undefined)>
 }
 interface Record {
  /**
   * Deprecated: replaced with GetUnsavedFiles.
   */
  getUploadedFiles(key: string): Array<(filesystem.File | undefined)>
 }
 interface Record {
  /**
   * Retrieves the "key" json field value and unmarshals it into "result".
   * 
   * Example
   * 
   * ```
   * 	result := struct {
   * 	    FirstName string `json:"first_name"`
   * 	}{}
   * 	err := m.UnmarshalJSONField("my_field_name", &result)
   * ```
   */
  unmarshalJSONField(key: string, result: any): void
 }
 interface Record {
  /**
   * ExpandedOne retrieves a single relation Record from the already
   * loaded expand data of the current model.
   * 
   * If the requested expand relation is multiple, this method returns
   * only first available Record from the expanded relation.
   * 
   * Returns nil if there is no such expand relation loaded.
   */
  expandedOne(relField: string): (Record)
 }
 interface Record {
  /**
   * ExpandedAll retrieves a slice of relation Records from the already
   * loaded expand data of the current model.
   * 
   * If the requested expand relation is single, this method normalizes
   * the return result and will wrap the single model as a slice.
   * 
   * Returns nil slice if there is no such expand relation loaded.
   */
  expandedAll(relField: string): Array<(Record | undefined)>
 }
 interface Record {
  /**
   * FindFileFieldByFile returns the first file type field for which
   * any of the record's data contains the provided filename.
   */
  findFileFieldByFile(filename: string): (FileField)
 }
 interface Record {
  /**
   * DBExport implements the [DBExporter] interface and returns a key-value
   * map with the data to be persisted when saving the Record in the database.
   */
  dbExport(app: App): _TygojaDict
 }
 interface Record {
  /**
   * Hide hides the specified fields from the public safe serialization of the record.
   */
  hide(...fieldNames: string[]): (Record)
 }
 interface Record {
  /**
   * Unhide forces to unhide the specified fields from the public safe serialization
   * of the record (even when the collection field itself is marked as hidden).
   */
  unhide(...fieldNames: string[]): (Record)
 }
 interface Record {
  /**
   * PublicExport exports only the record fields that are safe to be public.
   * 
   * To export unknown data fields you need to set record.WithCustomData(true).
   * 
   * For auth records, to force the export of the email field you need to set
   * record.IgnoreEmailVisibility(true).
   */
  publicExport(): _TygojaDict
 }
 interface Record {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   * 
   * Only the data exported by `PublicExport()` will be serialized.
   */
  marshalJSON(): string|Array<number>
 }
 interface Record {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   */
  unmarshalJSON(data: string|Array<number>): void
 }
 interface Record {
  /**
   * ReplaceModifiers returns a new map with applied modifier
   * values based on the current record and the specified data.
   * 
   * The resolved modifier keys will be removed.
   * 
   * Multiple modifiers will be applied one after another,
   * while reusing the previous base key value result (ex. 1; -5; +2 => -2).
   * 
   * Note that because Go doesn't guaranteed the iteration order of maps,
   * we would explicitly apply shorter keys first for a more consistent and reproducible behavior.
   * 
   * Example usage:
   * 
   * ```
   * 	 newData := record.ReplaceModifiers(data)
   * 		// record: {"field": 10}
   * 		// data:   {"field+": 5}
   * 		// result: {"field": 15}
   * ```
   */
  replaceModifiers(data: _TygojaDict): _TygojaDict
 }
 interface Record {
  /**
   * Email returns the "email" record field value (usually available with Auth collections).
   */
  email(): string
 }
 interface Record {
  /**
   * SetEmail sets the "email" record field value (usually available with Auth collections).
   */
  setEmail(email: string): void
 }
 interface Record {
  /**
   * Verified returns the "emailVisibility" record field value (usually available with Auth collections).
   */
  emailVisibility(): boolean
 }
 interface Record {
  /**
   * SetEmailVisibility sets the "emailVisibility" record field value (usually available with Auth collections).
   */
  setEmailVisibility(visible: boolean): void
 }
 interface Record {
  /**
   * Verified returns the "verified" record field value (usually available with Auth collections).
   */
  verified(): boolean
 }
 interface Record {
  /**
   * SetVerified sets the "verified" record field value (usually available with Auth collections).
   */
  setVerified(verified: boolean): void
 }
 interface Record {
  /**
   * TokenKey returns the "tokenKey" record field value (usually available with Auth collections).
   */
  tokenKey(): string
 }
 interface Record {
  /**
   * SetTokenKey sets the "tokenKey" record field value (usually available with Auth collections).
   */
  setTokenKey(key: string): void
 }
 interface Record {
  /**
   * RefreshTokenKey generates and sets a new random auth record "tokenKey".
   */
  refreshTokenKey(): void
 }
 interface Record {
  /**
   * SetPassword sets the "password" record field value (usually available with Auth collections).
   */
  setPassword(password: string): void
 }
 interface Record {
  /**
   * SetRandomPassword sets the "password" auth record field to a random autogenerated value.
   * 
   * The autogenerated password is ~30 characters and it is set directly as hash,
   * aka. the field plain password value validators (length, pattern, etc.) are ignored
   * (this is usually used as part of the auto created OTP or OAuth2 user flows).
   */
  setRandomPassword(): string
 }
 interface Record {
  /**
   * ValidatePassword validates a plain password against the "password" record field.
   * 
   * Returns false if the password is incorrect.
   */
  validatePassword(password: string): boolean
 }
 interface Record {
  /**
   * IsSuperuser returns whether the current record is a superuser, aka.
   * whether the record is from the _superusers collection.
   */
  isSuperuser(): boolean
 }
 /**
  * RecordProxy defines an interface for a Record proxy/project model,
  * aka. custom model struct that acts on behalve the proxied Record to
  * allow for example typed getter/setters for the Record fields.
  * 
  * To implement the interface it is usually enough to embed the [BaseRecordProxy] struct.
  */
 interface RecordProxy {
  [key:string]: any;
  /**
   * ProxyRecord returns the proxied Record model.
   */
  proxyRecord(): (Record)
  /**
   * SetProxyRecord loads the specified record model into the current proxy.
   */
  setProxyRecord(record: Record): void
 }
 /**
  * BaseRecordProxy implements the [RecordProxy] interface and it is intended
  * to be used as embed to custom user provided Record proxy structs.
  */
 type _sjwWfBO = Record
 interface BaseRecordProxy extends _sjwWfBO {
 }
 interface BaseRecordProxy {
  /**
   * ProxyRecord returns the proxied Record model.
   */
  proxyRecord(): (Record)
 }
 interface BaseRecordProxy {
  /**
   * SetProxyRecord loads the specified record model into the current proxy.
   */
  setProxyRecord(record: Record): void
 }
 interface BaseApp {
  /**
   * RecordQuery returns a new Record select query from a collection model, id or name.
   * 
   * In case a collection id or name is provided and that collection doesn't
   * actually exists, the generated query will be created with a cancelled context
   * and will fail once an executor (Row(), One(), All(), etc.) is called.
   */
  recordQuery(collectionModelOrIdentifier: any): (dbx.SelectQuery)
 }
 interface BaseApp {
  /**
   * FindRecordById finds the Record model by its id.
   */
  findRecordById(collectionModelOrIdentifier: any, recordId: string, ...optFilters: ((q: dbx.SelectQuery) => void)[]): (Record)
 }
 interface BaseApp {
  /**
   * FindRecordsByIds finds all records by the specified ids.
   * If no records are found, returns an empty slice.
   */
  findRecordsByIds(collectionModelOrIdentifier: any, recordIds: Array<string>, ...optFilters: ((q: dbx.SelectQuery) => void)[]): Array<(Record | undefined)>
 }
 interface BaseApp {
  /**
   * FindAllRecords finds all records matching specified db expressions.
   * 
   * Returns all collection records if no expression is provided.
   * 
   * Returns an empty slice if no records are found.
   * 
   * Example:
   * 
   * ```
   * 	// no extra expressions
   * 	app.FindAllRecords("example")
   * 
   * 	// with extra expressions
   * 	expr1 := dbx.HashExp{"email": "test@example.com"}
   * 	expr2 := dbx.NewExp("LOWER(username) = {:username}", dbx.Params{"username": "test"})
   * 	app.FindAllRecords("example", expr1, expr2)
   * ```
   */
  findAllRecords(collectionModelOrIdentifier: any, ...exprs: dbx.Expression[]): Array<(Record | undefined)>
 }
 interface BaseApp {
  /**
   * FindFirstRecordByData returns the first found record matching
   * the provided key-value pair.
   */
  findFirstRecordByData(collectionModelOrIdentifier: any, key: string, value: any): (Record)
 }
 interface BaseApp {
  /**
   * FindRecordsByFilter returns limit number of records matching the
   * provided string filter.
   * 
   * NB! Use the last "params" argument to bind untrusted user variables!
   * 
   * The filter argument is optional and can be empty string to target
   * all available records.
   * 
   * The sort argument is optional and can be empty string OR the same format
   * used in the web APIs, ex. "-created,title".
   * 
   * If the limit argument is <= 0, no limit is applied to the query and
   * all matching records are returned.
   * 
   * Returns an empty slice if no records are found.
   * 
   * Example:
   * 
   * ```
   * 	app.FindRecordsByFilter(
   * 		"posts",
   * 		"title ~ {:title} && visible = {:visible}",
   * 		"-created",
   * 		10,
   * 		0,
   * 		dbx.Params{"title": "lorem ipsum", "visible": true}
   * 	)
   * ```
   */
  findRecordsByFilter(collectionModelOrIdentifier: any, filter: string, sort: string, limit: number, offset: number, ...params: dbx.Params[]): Array<(Record | undefined)>
 }
 interface BaseApp {
  /**
   * FindFirstRecordByFilter returns the first available record matching the provided filter (if any).
   * 
   * NB! Use the last params argument to bind untrusted user variables!
   * 
   * Returns sql.ErrNoRows if no record is found.
   * 
   * Example:
   * 
   * ```
   * 	app.FindFirstRecordByFilter("posts", "")
   * 	app.FindFirstRecordByFilter("posts", "slug={:slug} && status='public'", dbx.Params{"slug": "test"})
   * ```
   */
  findFirstRecordByFilter(collectionModelOrIdentifier: any, filter: string, ...params: dbx.Params[]): (Record)
 }
 interface BaseApp {
  /**
   * CountRecords returns the total number of records in a collection.
   */
  countRecords(collectionModelOrIdentifier: any, ...exprs: dbx.Expression[]): number
 }
 interface BaseApp {
  /**
   * FindAuthRecordByToken finds the auth record associated with the provided JWT
   * (auth, file, verifyEmail, changeEmail, passwordReset types).
   * 
   * Optionally specify a list of validTypes to check tokens only from those types.
   * 
   * Returns an error if the JWT is invalid, expired or not associated to an auth collection record.
   */
  findAuthRecordByToken(token: string, ...validTypes: string[]): (Record)
 }
 interface BaseApp {
  /**
   * FindAuthRecordByEmail finds the auth record associated with the provided email.
   * 
   * The email check would be case-insensitive if the related collection
   * email unique index has COLLATE NOCASE specified for the email column.
   * 
   * Returns an error if it is not an auth collection or the record is not found.
   */
  findAuthRecordByEmail(collectionModelOrIdentifier: any, email: string): (Record)
 }
 interface BaseApp {
  /**
   * CanAccessRecord checks if a record is allowed to be accessed by the
   * specified requestInfo and accessRule.
   * 
   * Rule and db checks are ignored in case requestInfo.Auth is a superuser.
   * 
   * The returned error indicate that something unexpected happened during
   * the check (eg. invalid rule or db query error).
   * 
   * The method always return false on invalid rule or db query error.
   * 
   * Example:
   * 
   * ```
   * 	requestInfo, _ := e.RequestInfo()
   * 	record, _ := app.FindRecordById("example", "RECORD_ID")
   * 	rule := types.Pointer("@request.auth.id != '' || status = 'public'")
   * 	// ... or use one of the record collection's rule, eg. record.Collection().ViewRule
   * 
   * 	if ok, _ := app.CanAccessRecord(record, requestInfo, rule); ok { ... }
   * ```
   */
  canAccessRecord(record: Record, requestInfo: RequestInfo, accessRule: string): boolean
 }
 /**
  * ExpandFetchFunc defines the function that is used to fetch the expanded relation records.
  */
 interface ExpandFetchFunc {(relCollection: Collection, relIds: Array<string>): Array<(Record | undefined)> }
 interface BaseApp {
  /**
   * ExpandRecord expands the relations of a single Record model.
   * 
   * If optFetchFunc is not set, then a default function will be used
   * that returns all relation records.
   * 
   * Returns a map with the failed expand parameters and their errors.
   */
  expandRecord(record: Record, expands: Array<string>, optFetchFunc: ExpandFetchFunc): _TygojaDict
 }
 interface BaseApp {
  /**
   * ExpandRecords expands the relations of the provided Record models list.
   * 
   * If optFetchFunc is not set, then a default function will be used
   * that returns all relation records.
   * 
   * Returns a map with the failed expand parameters and their errors.
   */
  expandRecords(records: Array<(Record | undefined)>, expands: Array<string>, optFetchFunc: ExpandFetchFunc): _TygojaDict
 }
 interface Record {
  /**
   * NewStaticAuthToken generates and returns a new static record authentication token.
   * 
   * Static auth tokens are similar to the regular auth tokens, but are
   * non-refreshable and support custom duration.
   * 
   * Zero or negative duration will fallback to the duration from the auth collection settings.
   */
  newStaticAuthToken(duration: time.Duration): string
 }
 interface Record {
  /**
   * NewAuthToken generates and returns a new record authentication token.
   */
  newAuthToken(): string
 }
 interface Record {
  /**
   * NewVerificationToken generates and returns a new record verification token.
   */
  newVerificationToken(): string
 }
 interface Record {
  /**
   * NewPasswordResetToken generates and returns a new auth record password reset request token.
   */
  newPasswordResetToken(): string
 }
 interface Record {
  /**
   * NewEmailChangeToken generates and returns a new auth record change email request token.
   */
  newEmailChangeToken(newEmail: string): string
 }
 interface Record {
  /**
   * NewFileToken generates and returns a new record private file access token.
   */
  newFileToken(): string
 }
 interface settings {
  smtp: SMTPConfig
  backups: BackupsConfig
  s3: S3Config
  meta: MetaConfig
  rateLimits: RateLimitsConfig
  trustedProxy: TrustedProxyConfig
  batch: BatchConfig
  logs: LogsConfig
 }
 /**
  * Settings defines the PocketBase app settings.
  */
 type _sHSTpUg = settings
 interface Settings extends _sHSTpUg {
 }
 interface Settings {
  /**
   * TableName implements [Model.TableName] interface method.
   */
  tableName(): string
 }
 interface Settings {
  /**
   * PK implements [Model.LastSavedPK] interface method.
   */
  lastSavedPK(): any
 }
 interface Settings {
  /**
   * PK implements [Model.PK] interface method.
   */
  pk(): any
 }
 interface Settings {
  /**
   * IsNew implements [Model.IsNew] interface method.
   */
  isNew(): boolean
 }
 interface Settings {
  /**
   * MarkAsNew implements [Model.MarkAsNew] interface method.
   */
  markAsNew(): void
 }
 interface Settings {
  /**
   * MarkAsNew implements [Model.MarkAsNotNew] interface method.
   */
  markAsNotNew(): void
 }
 interface Settings {
  /**
   * PostScan implements [Model.PostScan] interface method.
   */
  postScan(): void
 }
 interface Settings {
  /**
   * String returns a serialized string representation of the current settings.
   */
  string(): string
 }
 interface Settings {
  /**
   * DBExport prepares and exports the current settings for db persistence.
   */
  dbExport(app: App): _TygojaDict
 }
 interface Settings {
  /**
   * PostValidate implements the [PostValidator] interface and defines
   * the Settings model validations.
   */
  postValidate(ctx: context.Context, app: App): void
 }
 interface Settings {
  /**
   * Merge merges the "other" settings into the current one.
   */
  merge(other: Settings): void
 }
 interface Settings {
  /**
   * Clone creates a new deep copy of the current settings.
   */
  clone(): (Settings)
 }
 interface Settings {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   * 
   * Note that sensitive fields (S3 secret, SMTP password, etc.) are excluded.
   */
  marshalJSON(): string|Array<number>
 }
 interface SMTPConfig {
  enabled: boolean
  port: number
  host: string
  username: string
  password: string
  /**
   * SMTP AUTH - PLAIN (default) or LOGIN
   */
  authMethod: string
  /**
   * Whether to enforce TLS encryption for the mail server connection.
   * 
   * When set to false StartTLS command is send, leaving the server
   * to decide whether to upgrade the connection or not.
   */
  tls: boolean
  /**
   * LocalName is optional domain name or IP address used for the
   * EHLO/HELO exchange (if not explicitly set, defaults to "localhost").
   * 
   * This is required only by some SMTP servers, such as Gmail SMTP-relay.
   */
  localName: string
 }
 interface SMTPConfig {
  /**
   * Validate makes SMTPConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface S3Config {
  enabled: boolean
  bucket: string
  region: string
  endpoint: string
  accessKey: string
  secret: string
  forcePathStyle: boolean
 }
 interface S3Config {
  /**
   * Validate makes S3Config validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface BatchConfig {
  enabled: boolean
  /**
   * MaxRequests is the maximum allowed batch request to execute.
   */
  maxRequests: number
  /**
   * Timeout is the the max duration in seconds to wait before cancelling the batch transaction.
   */
  timeout: number
  /**
   * MaxBodySize is the maximum allowed batch request body size in bytes.
   * 
   * If not set, fallbacks to max ~128MB.
   */
  maxBodySize: number
 }
 interface BatchConfig {
  /**
   * Validate makes BatchConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface BackupsConfig {
  /**
   * Cron is a cron expression to schedule auto backups, eg. "* * * * *".
   * 
   * Leave it empty to disable the auto backups functionality.
   */
  cron: string
  /**
   * CronMaxKeep is the the max number of cron generated backups to
   * keep before removing older entries.
   * 
   * This field works only when the cron config has valid cron expression.
   */
  cronMaxKeep: number
  /**
   * S3 is an optional S3 storage config specifying where to store the app backups.
   */
  s3: S3Config
 }
 interface BackupsConfig {
  /**
   * Validate makes BackupsConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface MetaConfig {
  appName: string
  appURL: string
  senderName: string
  senderAddress: string
  hideControls: boolean
 }
 interface MetaConfig {
  /**
   * Validate makes MetaConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface LogsConfig {
  maxDays: number
  minLevel: number
  logIP: boolean
  logAuthId: boolean
 }
 interface LogsConfig {
  /**
   * Validate makes LogsConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface TrustedProxyConfig {
  /**
   * Headers is a list of explicit trusted header(s) to check.
   */
  headers: Array<string>
  /**
   * UseLeftmostIP specifies to use the left-mostish IP from the trusted headers.
   * 
   * Note that this could be insecure when used with X-Forwarded-For header
   * because some proxies like AWS ELB allow users to prepend their own header value
   * before appending the trusted ones.
   */
  useLeftmostIP: boolean
 }
 interface TrustedProxyConfig {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface TrustedProxyConfig {
  /**
   * Validate makes RateLimitRule validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RateLimitsConfig {
  rules: Array<RateLimitRule>
  enabled: boolean
 }
 interface RateLimitsConfig {
  /**
   * FindRateLimitRule returns the first matching rule based on the provided labels.
   * 
   * Optionally you can further specify a list of valid RateLimitRule.Audience values to further filter the matching rule
   * (aka. the rule Audience will have to exist in one of the specified options).
   */
  findRateLimitRule(searchLabels: Array<string>, ...optOnlyAudience: string[]): [RateLimitRule, boolean]
 }
 interface RateLimitsConfig {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface RateLimitsConfig {
  /**
   * Validate makes RateLimitsConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RateLimitRule {
  /**
   * Label is the identifier of the current rule.
   * 
   * It could be a tag, complete path or path prerefix (when ends with `/`).
   * 
   * Example supported labels:
   * ```
   *   - test_a (plain text "tag")
   *   - users:create
   *   - *:create
   *   - /
   *   - /api
   *   - POST /api/collections/
   * ```
   */
  label: string
  /**
   * Audience specifies the auth group the rule should apply for:
   * ```
   *   - ""      - both guests and authenticated users (default)
   *   - "@guest" - only for guests
   *   - "@auth"  - only for authenticated users
   * ```
   */
  audience: string
  /**
   * Duration specifies the interval (in seconds) per which to reset
   * the counted/accumulated rate limiter tokens.
   */
  duration: number
  /**
   * MaxRequests is the max allowed number of requests per Duration.
   */
  maxRequests: number
 }
 interface RateLimitRule {
  /**
   * Validate makes RateLimitRule validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RateLimitRule {
  /**
   * DurationTime returns the tag's Duration as [time.Duration].
   */
  durationTime(): time.Duration
 }
 interface RateLimitRule {
  /**
   * String returns a string representation of the rule.
   */
  string(): string
 }
 type _sPqlRfa = BaseModel
 interface Param extends _sPqlRfa {
  created: types.DateTime
  updated: types.DateTime
  value: types.JSONRaw
 }
 interface Param {
  tableName(): string
 }
 interface BaseApp {
  /**
   * ReloadSettings initializes and reloads the stored application settings.
   * 
   * If no settings were stored it will persist the current app ones.
   */
  reloadSettings(): void
 }
 interface BaseApp {
  /**
   * DeleteView drops the specified view name.
   * 
   * This method is a no-op if a view with the provided name doesn't exist.
   * 
   * NB! Be aware that this method is vulnerable to SQL injection and the
   * "name" argument must come only from trusted input!
   */
  deleteView(name: string): void
 }
 interface BaseApp {
  /**
   * SaveView creates (or updates already existing) persistent SQL view.
   * 
   * NB! Be aware that this method is vulnerable to SQL injection and the
   * "selectQuery" argument must come only from trusted input!
   */
  saveView(name: string, selectQuery: string): void
 }
 interface BaseApp {
  /**
   * CreateViewFields creates a new FieldsList from the provided select query.
   * 
   * There are some caveats:
   * - The select query must have an "id" column.
   * - Wildcard ("*") columns are not supported to avoid accidentally leaking sensitive data.
   */
  createViewFields(selectQuery: string): FieldsList
 }
 interface BaseApp {
  /**
   * FindRecordByViewFile returns the original Record of the provided view collection file.
   */
  findRecordByViewFile(viewCollectionModelOrIdentifier: any, fileFieldName: string, filename: string): (Record)
 }
 interface queryField {
 }
 interface identifier {
 }
 interface identifiersParser {
 }
}

/**
 * Package mails implements various helper methods for sending common
 * emails like forgotten password, verification, etc.
 */
namespace mails {
 interface sendRecordAuthAlert {
  /**
   * SendRecordAuthAlert sends a new device login alert to the specified auth record.
   */
  (app: CoreApp, authRecord: core.Record): void
 }
 interface sendRecordOTP {
  /**
   * SendRecordOTP sends OTP email to the specified auth record.
   * 
   * This method will also update the "sentTo" field of the related OTP record to the mail sent To address (if the OTP exists and not already assigned).
   */
  (app: CoreApp, authRecord: core.Record, otpId: string, pass: string): void
 }
 interface sendRecordPasswordReset {
  /**
   * SendRecordPasswordReset sends a password reset request email to the specified auth record.
   */
  (app: CoreApp, authRecord: core.Record): void
 }
 interface sendRecordVerification {
  /**
   * SendRecordVerification sends a verification request email to the specified auth record.
   */
  (app: CoreApp, authRecord: core.Record): void
 }
 interface sendRecordChangeEmail {
  /**
   * SendRecordChangeEmail sends a change email confirmation email to the specified auth record.
   */
  (app: CoreApp, authRecord: core.Record, newEmail: string): void
 }
}

namespace forms {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * AppleClientSecretCreate is a form struct to generate a new Apple Client Secret.
  * 
  * Reference: https://developer.apple.com/documentation/sign_in_with_apple/generate_and_validate_tokens
  */
 interface AppleClientSecretCreate {
  /**
   * ClientId is the identifier of your app (aka. Service ID).
   */
  clientId: string
  /**
   * TeamId is a 10-character string associated with your developer account
   * (usually could be found next to your name in the Apple Developer site).
   */
  teamId: string
  /**
   * KeyId is a 10-character key identifier generated for the "Sign in with Apple"
   * private key associated with your developer account.
   */
  keyId: string
  /**
   * PrivateKey is the private key associated to your app.
   * Usually wrapped within -----BEGIN PRIVATE KEY----- X -----END PRIVATE KEY-----.
   */
  privateKey: string
  /**
   * Duration specifies how long the generated JWT should be considered valid.
   * The specified value must be in seconds and max 15777000 (~6months).
   */
  duration: number
 }
 interface newAppleClientSecretCreate {
  /**
   * NewAppleClientSecretCreate creates a new [AppleClientSecretCreate] form with initializer
   * config created from the provided [CoreApp] instances.
   */
  (app: CoreApp): (AppleClientSecretCreate)
 }
 interface AppleClientSecretCreate {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface AppleClientSecretCreate {
  /**
   * Submit validates the form and returns a new Apple Client Secret JWT.
   */
  submit(): string
 }
 interface RecordUpsert {
 }
 interface newRecordUpsert {
  /**
   * NewRecordUpsert creates a new [RecordUpsert] form from the provided [CoreApp] and [core.Record] instances
   * (for create you could pass a pointer to an empty Record - core.NewRecord(collection)).
   */
  (app: CoreApp, record: core.Record): (RecordUpsert)
 }
 interface RecordUpsert {
  /**
   * SetContext assigns ctx as context of the current form.
   */
  setContext(ctx: context.Context): void
 }
 interface RecordUpsert {
  /**
   * SetApp replaces the current form app instance.
   * 
   * This could be used for example if you want to change at later stage
   * before submission to change from regular -> transactional app instance.
   */
  setApp(app: CoreApp): void
 }
 interface RecordUpsert {
  /**
   * SetRecord replaces the current form record instance.
   */
  setRecord(record: core.Record): void
 }
 interface RecordUpsert {
  /**
   * ResetAccess resets the form access level to the accessLevelDefault.
   */
  resetAccess(): void
 }
 interface RecordUpsert {
  /**
   * GrantManagerAccess updates the form access level to "manager" allowing
   * directly changing some system record fields (often used with auth collection records).
   */
  grantManagerAccess(): void
 }
 interface RecordUpsert {
  /**
   * GrantSuperuserAccess updates the form access level to "superuser" allowing
   * directly changing all system record fields, including those marked as "Hidden".
   */
  grantSuperuserAccess(): void
 }
 interface RecordUpsert {
  /**
   * HasManageAccess reports whether the form has "manager" or "superuser" level access.
   */
  hasManageAccess(): boolean
 }
 interface RecordUpsert {
  /**
   * Load loads the provided data into the form and the related record.
   */
  load(data: _TygojaDict): void
 }
 interface RecordUpsert {
  /**
   * Deprecated: It was previously used as part of the record create action but it is not needed anymore and will be removed in the future.
   * 
   * DrySubmit performs a temp form submit within a transaction and reverts it at the end.
   * For actual record persistence, check the [RecordUpsert.Submit()] method.
   * 
   * This method doesn't perform validations, handle file uploads/deletes or trigger app save events!
   */
  drySubmit(callback: (txApp: CoreApp, drySavedRecord: core.Record) => void): void
 }
 interface RecordUpsert {
  /**
   * Submit validates the form specific validations and attempts to save the form record.
   */
  submit(): void
 }
 /**
  * TestEmailSend is a email template test request form.
  */
 interface TestEmailSend {
  email: string
  template: string
  collection: string // optional, fallbacks to _superusers
 }
 interface newTestEmailSend {
  /**
   * NewTestEmailSend creates and initializes new TestEmailSend form.
   */
  (app: CoreApp): (TestEmailSend)
 }
 interface TestEmailSend {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface TestEmailSend {
  /**
   * Submit validates and sends a test email to the form.Email address.
   */
  submit(): void
 }
 /**
  * TestS3Filesystem defines a S3 filesystem connection test.
  */
 interface TestS3Filesystem {
  /**
   * The name of the filesystem - storage or backups
   */
  filesystem: string
 }
 interface newTestS3Filesystem {
  /**
   * NewTestS3Filesystem creates and initializes new TestS3Filesystem form.
   */
  (app: CoreApp): (TestS3Filesystem)
 }
 interface TestS3Filesystem {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface TestS3Filesystem {
  /**
   * Submit validates and performs a S3 filesystem connection test.
   */
  submit(): void
 }
}

namespace apis {
 interface toApiError {
  /**
   * ToApiError wraps err into ApiError instance (if not already).
   */
  (err: Error): (router.ApiError)
 }
 interface newApiError {
  /**
   * NewApiError is an alias for [router.NewApiError].
   */
  (status: number, message: string, errData: any): (router.ApiError)
 }
 interface newBadRequestError {
  /**
   * NewBadRequestError is an alias for [router.NewBadRequestError].
   */
  (message: string, errData: any): (router.ApiError)
 }
 interface newNotFoundError {
  /**
   * NewNotFoundError is an alias for [router.NewNotFoundError].
   */
  (message: string, errData: any): (router.ApiError)
 }
 interface newForbiddenError {
  /**
   * NewForbiddenError is an alias for [router.NewForbiddenError].
   */
  (message: string, errData: any): (router.ApiError)
 }
 interface newUnauthorizedError {
  /**
   * NewUnauthorizedError is an alias for [router.NewUnauthorizedError].
   */
  (message: string, errData: any): (router.ApiError)
 }
 interface newTooManyRequestsError {
  /**
   * NewTooManyRequestsError is an alias for [router.NewTooManyRequestsError].
   */
  (message: string, errData: any): (router.ApiError)
 }
 interface newInternalServerError {
  /**
   * NewInternalServerError is an alias for [router.NewInternalServerError].
   */
  (message: string, errData: any): (router.ApiError)
 }
 interface backupFileInfo {
  modified: types.DateTime
  key: string
  size: number
 }
 // @ts-ignore
 import validation = ozzo_validation
 interface backupCreateForm {
  name: string
 }
 interface backupUploadForm {
  file?: filesystem.File
 }
 interface newRouter {
  /**
   * NewRouter returns a new router instance loaded with the default app middlewares and api routes.
   */
  (app: CoreApp): (router.Router<core.RequestEvent | undefined>)
 }
 interface wrapStdHandler {
  /**
   * WrapStdHandler wraps Go [http.Handler] into a PocketBase handler func.
   */
  (h: http.Handler): (_arg0: core.RequestEvent) => void
 }
 interface wrapStdMiddleware {
  /**
   * WrapStdMiddleware wraps Go [func(http.Handler) http.Handle] into a PocketBase middleware func.
   */
  (m: (_arg0: http.Handler) => http.Handler): (_arg0: core.RequestEvent) => void
 }
 interface mustSubFS {
  /**
   * MustSubFS returns an [fs.FS] corresponding to the subtree rooted at fsys's dir.
   * 
   * This is similar to [fs.Sub] but panics on failure.
   */
  (fsys: fs.FS, dir: string): fs.FS
 }
 interface _static {
  /**
   * Static is a handler function to serve static directory content from fsys.
   * 
   * If a file resource is missing and indexFallback is set, the request
   * will be forwarded to the base index.html (useful for SPA with pretty urls).
   * 
   * NB! Expects the route to have a "{path...}" wildcard parameter.
   * 
   * Special redirects:
   * ```
   *   - if "path" is a file that ends in index.html, it is redirected to its non-index.html version (eg. /test/index.html -> /test/)
   *   - if "path" is a directory that has index.html, the index.html file is rendered,
   *     otherwise if missing - returns 404 or fallback to the root index.html if indexFallback is set
   * ```
   * 
   * Example:
   * 
   * ```
   * 	fsys := os.DirFS("./pb_public")
   * 	router.GET("/files/{path...}", apis.Static(fsys, false))
   * ```
   */
  (fsys: fs.FS, indexFallback: boolean): (_arg0: core.RequestEvent) => void
 }
 interface HandleFunc {(e: core.RequestEvent): void }
 interface BatchActionHandlerFunc {(app: CoreApp, ir: core.InternalRequest, params: _TygojaDict, next: (data: any) => void): HandleFunc }
 interface BatchRequestResult {
  body: any
  status: number
 }
 interface batchRequestsForm {
  requests: Array<(core.InternalRequest | undefined)>
 }
 interface batchProcessor {
 }
 interface batchProcessor {
  process(batch: Array<(core.InternalRequest | undefined)>, timeout: time.Duration): void
 }
 interface BatchResponseError {
 }
 interface BatchResponseError {
  error(): string
 }
 interface BatchResponseError {
  code(): string
 }
 interface BatchResponseError {
  resolve(errData: _TygojaDict): any
 }
 interface BatchResponseError {
  marshalJSON(): string|Array<number>
 }
 interface collectionsImportForm {
  collections: Array<_TygojaDict>
  deleteMissing: boolean
 }
 interface fileApi {
 }
 interface defaultInstallerFunc {
  /**
   * DefaultInstallerFunc is the default PocketBase installer function.
   * 
   * It will attempt to open a link in the browser (with a short-lived auth
   * token for the systemSuperuser) to the installer UI so that users can
   * create their own custom superuser record.
   * 
   * See https://github.com/pocketbase/pocketbase/discussions/5814.
   */
  (app: CoreApp, systemSuperuser: core.Record, baseURL: string): void
 }
 interface requireGuestOnly {
  /**
   * RequireGuestOnly middleware requires a request to NOT have a valid
   * Authorization header.
   * 
   * This middleware is the opposite of [apis.RequireAuth()].
   */
  (): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface requireAuth {
  /**
   * RequireAuth middleware requires a request to have a valid record Authorization header.
   * 
   * The auth record could be from any collection.
   * You can further filter the allowed record auth collections by specifying their names.
   * 
   * Example:
   * 
   * ```
   * 	apis.RequireAuth()                      // any auth collection
   * 	apis.RequireAuth("_superusers", "users") // only the listed auth collections
   * ```
   */
  (...optCollectionNames: string[]): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface requireSuperuserAuth {
  /**
   * RequireSuperuserAuth middleware requires a request to have
   * a valid superuser Authorization header.
   */
  (): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface requireSuperuserOrOwnerAuth {
  /**
   * RequireSuperuserOrOwnerAuth middleware requires a request to have
   * a valid superuser or regular record owner Authorization header set.
   * 
   * This middleware is similar to [apis.RequireAuth()] but
   * for the auth record token expects to have the same id as the path
   * parameter ownerIdPathParam (default to "id" if empty).
   */
  (ownerIdPathParam: string): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface requireSameCollectionContextAuth {
  /**
   * RequireSameCollectionContextAuth middleware requires a request to have
   * a valid record Authorization header and the auth record's collection to
   * match the one from the route path parameter (default to "collection" if collectionParam is empty).
   */
  (collectionPathParam: string): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface skipSuccessActivityLog {
  /**
   * SkipSuccessActivityLog is a helper middleware that instructs the global
   * activity logger to log only requests that have failed/returned an error.
   */
  (): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface bodyLimit {
  /**
   * BodyLimit returns a middleware handler that changes the default request body size limit.
   * 
   * If limitBytes <= 0, no limit is applied.
   * 
   * Otherwise, if the request body size exceeds the configured limitBytes,
   * it sends 413 error response.
   */
  (limitBytes: number): (hook.Handler<core.RequestEvent | undefined>)
 }
 type _sTIPhiP = io.ReadCloser
 interface limitedReader extends _sTIPhiP {
 }
 interface limitedReader {
  read(b: string|Array<number>): number
 }
 interface limitedReader {
  reread(): void
 }
 /**
  * CORSConfig defines the config for CORS middleware.
  */
 interface CORSConfig {
  /**
   * AllowOrigins determines the value of the Access-Control-Allow-Origin
   * response header.  This header defines a list of origins that may access the
   * resource.  The wildcard characters '*' and '?' are supported and are
   * converted to regex fragments '.*' and '.' accordingly.
   * 
   * Security: use extreme caution when handling the origin, and carefully
   * validate any logic. Remember that attackers may register hostile domain names.
   * See https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
   * 
   * Optional. Default value []string{"*"}.
   * 
   * See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
   */
  allowOrigins: Array<string>
  /**
   * AllowOriginFunc is a custom function to validate the origin. It takes the
   * origin as an argument and returns true if allowed or false otherwise. If
   * an error is returned, it is returned by the handler. If this option is
   * set, AllowOrigins is ignored.
   * 
   * Security: use extreme caution when handling the origin, and carefully
   * validate any logic. Remember that attackers may register hostile domain names.
   * See https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
   * 
   * Optional.
   */
  allowOriginFunc: (origin: string) => boolean
  /**
   * AllowMethods determines the value of the Access-Control-Allow-Methods
   * response header.  This header specified the list of methods allowed when
   * accessing the resource.  This is used in response to a preflight request.
   * 
   * Optional. Default value DefaultCORSConfig.AllowMethods.
   * 
   * See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
   */
  allowMethods: Array<string>
  /**
   * AllowHeaders determines the value of the Access-Control-Allow-Headers
   * response header.  This header is used in response to a preflight request to
   * indicate which HTTP headers can be used when making the actual request.
   * 
   * Optional. Default value []string{}.
   * 
   * See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
   */
  allowHeaders: Array<string>
  /**
   * AllowCredentials determines the value of the
   * Access-Control-Allow-Credentials response header.  This header indicates
   * whether or not the response to the request can be exposed when the
   * credentials mode (Request.credentials) is true. When used as part of a
   * response to a preflight request, this indicates whether or not the actual
   * request can be made using credentials.  See also
   * [MDN: Access-Control-Allow-Credentials].
   * 
   * Optional. Default value false, in which case the header is not set.
   * 
   * Security: avoid using `AllowCredentials = true` with `AllowOrigins = *`.
   * See "Exploiting CORS misconfigurations for Bitcoins and bounties",
   * https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
   * 
   * See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
   */
  allowCredentials: boolean
  /**
   * UnsafeWildcardOriginWithAllowCredentials UNSAFE/INSECURE: allows wildcard '*' origin to be used with AllowCredentials
   * flag. In that case we consider any origin allowed and send it back to the client with `Access-Control-Allow-Origin` header.
   * 
   * This is INSECURE and potentially leads to [cross-origin](https://portswigger.net/research/exploiting-cors-misconfigurations-for-bitcoins-and-bounties)
   * attacks. See: https://github.com/labstack/echo/issues/2400 for discussion on the subject.
   * 
   * Optional. Default value is false.
   */
  unsafeWildcardOriginWithAllowCredentials: boolean
  /**
   * ExposeHeaders determines the value of Access-Control-Expose-Headers, which
   * defines a list of headers that clients are allowed to access.
   * 
   * Optional. Default value []string{}, in which case the header is not set.
   * 
   * See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Header
   */
  exposeHeaders: Array<string>
  /**
   * MaxAge determines the value of the Access-Control-Max-Age response header.
   * This header indicates how long (in seconds) the results of a preflight
   * request can be cached.
   * The header is set only if MaxAge != 0, negative value sends "0" which instructs browsers not to cache that response.
   * 
   * Optional. Default value 0 - meaning header is not sent.
   * 
   * See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
   */
  maxAge: number
 }
 interface cors {
  /**
   * CORS returns a CORS middleware.
   */
  (config: CORSConfig): (hook.Handler<core.RequestEvent | undefined>)
 }
 /**
  * GzipConfig defines the config for Gzip middleware.
  */
 interface GzipConfig {
  /**
   * Gzip compression level.
   * Optional. Default value -1.
   */
  level: number
  /**
   * Length threshold before gzip compression is applied.
   * Optional. Default value 0.
   * 
   * Most of the time you will not need to change the default. Compressing
   * a short response might increase the transmitted data because of the
   * gzip format overhead. Compressing the response will also consume CPU
   * and time on the server and the client (for decompressing). Depending on
   * your use case such a threshold might be useful.
   * 
   * See also:
   * https://webmasters.stackexchange.com/questions/31750/what-is-recommended-minimum-object-size-for-gzip-performance-benefits
   */
  minLength: number
 }
 interface gzip {
  /**
   * Gzip returns a middleware which compresses HTTP response using Gzip compression scheme.
   */
  (): (hook.Handler<core.RequestEvent | undefined>)
 }
 interface gzipWithConfig {
  /**
   * GzipWithConfig returns a middleware which compresses HTTP response using gzip compression scheme.
   */
  (config: GzipConfig): (hook.Handler<core.RequestEvent | undefined>)
 }
 type _sTutvJD = http.ResponseWriter&io.Writer
 interface gzipResponseWriter extends _sTutvJD {
 }
 interface gzipResponseWriter {
  writeHeader(code: number): void
 }
 interface gzipResponseWriter {
  write(b: string|Array<number>): number
 }
 interface gzipResponseWriter {
  flush(): void
 }
 interface gzipResponseWriter {
  hijack(): [net.Conn, (bufio.ReadWriter)]
 }
 interface gzipResponseWriter {
  push(target: string, opts: http.PushOptions): void
 }
 interface gzipResponseWriter {
  unwrap(): http.ResponseWriter
 }
 type _sZdHxEJ = sync.RWMutex
 interface rateLimiter extends _sZdHxEJ {
 }
 type _sIwZCwA = sync.Mutex
 interface fixedWindow extends _sIwZCwA {
 }
 interface realtimeSubscribeForm {
  clientId: string
  subscriptions: Array<string>
 }
 /**
  * recordData represents the broadcasted record subscrition message data.
  */
 interface recordData {
  record: any //  map or core.Record
  action: string
 }
 interface EmailChangeConfirmForm {
  token: string
  password: string
 }
 interface emailChangeRequestForm {
  newEmail: string
 }
 interface impersonateForm {
  /**
   * Duration is the optional custom token duration in seconds.
   */
  duration: number
 }
 interface otpResponse {
  enabled: boolean
  duration: number // in seconds
 }
 interface mfaResponse {
  enabled: boolean
  duration: number // in seconds
 }
 interface passwordResponse {
  identityFields: Array<string>
  enabled: boolean
 }
 interface oauth2Response {
  providers: Array<providerInfo>
  enabled: boolean
 }
 interface providerInfo {
  name: string
  displayName: string
  state: string
  authURL: string
  /**
   * @todo
   * deprecated: use AuthURL instead
   * AuthUrl will be removed after dropping v0.22 support
   */
  authUrl: string
  /**
   * technically could be omitted if the provider doesn't support PKCE,
   * but to avoid breaking existing typed clients we'll return them as empty string
   */
  codeVerifier: string
  codeChallenge: string
  codeChallengeMethod: string
 }
 interface authMethodsResponse {
  password: passwordResponse
  oauth2: oauth2Response
  mfa: mfaResponse
  otp: otpResponse
  /**
   * legacy fields
   * @todo remove after dropping v0.22 support
   */
  authProviders: Array<providerInfo>
  usernamePassword: boolean
  emailPassword: boolean
 }
 interface createOTPForm {
  email: string
 }
 interface recordConfirmPasswordResetForm {
  token: string
  password: string
  passwordConfirm: string
 }
 interface recordRequestPasswordResetForm {
  email: string
 }
 interface recordConfirmVerificationForm {
  token: string
 }
 interface recordRequestVerificationForm {
  email: string
 }
 interface recordOAuth2LoginForm {
  /**
   * Additional data that will be used for creating a new auth record
   * if an existing OAuth2 account doesn't exist.
   */
  createData: _TygojaDict
  /**
   * The name of the OAuth2 client provider (eg. "google")
   */
  provider: string
  /**
   * The authorization code returned from the initial request.
   */
  code: string
  /**
   * The optional PKCE code verifier as part of the code_challenge sent with the initial request.
   */
  codeVerifier: string
  /**
   * The redirect url sent with the initial request.
   */
  redirectURL: string
  /**
   * @todo
   * deprecated: use RedirectURL instead
   * RedirectUrl will be removed after dropping v0.22 support
   */
  redirectUrl: string
 }
 interface oauth2RedirectData {
  state: string
  code: string
  error: string
  /**
   * returned by Apple only
   */
  appleUser: string
 }
 interface authWithOTPForm {
  otpId: string
  password: string
 }
 interface authWithPasswordForm {
  identity: string
  password: string
  /**
   * IdentityField specifies the field to use to search for the identity
   * (leave it empty for "auto" detection).
   */
  identityField: string
 }
 // @ts-ignore
 import cryptoRand = rand
 interface recordAuthResponse {
  /**
   * RecordAuthResponse writes standardized json record auth response
   * into the specified request context.
   * 
   * The authMethod argument specify the name of the current authentication method (eg. password, oauth2, etc.)
   * that it is used primarily as an auth identifier during MFA and for login alerts.
   * 
   * Set authMethod to empty string if you want to ignore the MFA checks and the login alerts
   * (can be also adjusted additionally via the OnRecordAuthRequest hook).
   */
  (e: core.RequestEvent, authRecord: core.Record, authMethod: string, meta: any): void
 }
 interface enrichRecord {
  /**
   * EnrichRecord parses the request context and enrich the provided record:
   * ```
   *   - expands relations (if defaultExpands and/or ?expand query param is set)
   *   - ensures that the emails of the auth record and its expanded auth relations
   *     are visible only for the current logged superuser, record owner or record with manage access
   * ```
   */
  (e: core.RequestEvent, record: core.Record, ...defaultExpands: string[]): void
 }
 interface enrichRecords {
  /**
   * EnrichRecords parses the request context and enriches the provided records:
   * ```
   *   - expands relations (if defaultExpands and/or ?expand query param is set)
   *   - ensures that the emails of the auth records and their expanded auth relations
   *     are visible only for the current logged superuser, record owner or record with manage access
   * ```
   * 
   * Note: Expects all records to be from the same collection!
   */
  (e: core.RequestEvent, records: Array<(core.Record | undefined)>, ...defaultExpands: string[]): void
 }
 interface iterator<T> {
 }
 /**
  * ServeConfig defines a configuration struct for apis.Serve().
  */
 interface ServeConfig {
  /**
   * ShowStartBanner indicates whether to show or hide the server start console message.
   */
  showStartBanner: boolean
  /**
   * HttpAddr is the TCP address to listen for the HTTP server (eg. "127.0.0.1:80").
   */
  httpAddr: string
  /**
   * HttpsAddr is the TCP address to listen for the HTTPS server (eg. "127.0.0.1:443").
   */
  httpsAddr: string
  /**
   * Optional domains list to use when issuing the TLS certificate.
   * 
   * If not set, the host from the bound server address will be used.
   * 
   * For convenience, for each "non-www" domain a "www" entry and
   * redirect will be automatically added.
   */
  certificateDomains: Array<string>
  /**
   * AllowedOrigins is an optional list of CORS origins (default to "*").
   */
  allowedOrigins: Array<string>
 }
 interface serve {
  /**
   * Serve starts a new app web server.
   * 
   * NB! The app should be bootstrapped before starting the web server.
   * 
   * Example:
   * 
   * ```
   * 	app.Bootstrap()
   * 	apis.Serve(app, apis.ServeConfig{
   * 		HttpAddr:        "127.0.0.1:8080",
   * 		ShowStartBanner: false,
   * 	})
   * ```
   */
  (app: CoreApp, config: ServeConfig): void
 }
 interface serverErrorLogWriter {
 }
 interface serverErrorLogWriter {
  write(p: string|Array<number>): number
 }
}

namespace pocketbase {
 /**
  * PocketBase defines a PocketBase app launcher.
  * 
  * It implements [CoreApp] via embedding and all of the app interface methods
  * could be accessed directly through the instance (eg. PocketBase.DataDir()).
  */
 type _sBtBPkV = CoreApp
 interface PocketBase extends _sBtBPkV {
  /**
   * RootCmd is the main console command
   */
  rootCmd?: cobra.Command
 }
 /**
  * Config is the PocketBase initialization config struct.
  */
 interface Config {
  /**
   * hide the default console server info on app startup
   */
  hideStartBanner: boolean
  /**
   * optional default values for the console flags
   */
  defaultDev: boolean
  defaultDataDir: string // if not set, it will fallback to "./pb_data"
  defaultEncryptionEnv: string
  defaultQueryTimeout: time.Duration // default to core.DefaultQueryTimeout (in seconds)
  /**
   * optional DB configurations
   */
  dataMaxOpenConns: number // default to core.DefaultDataMaxOpenConns
  dataMaxIdleConns: number // default to core.DefaultDataMaxIdleConns
  auxMaxOpenConns: number // default to core.DefaultAuxMaxOpenConns
  auxMaxIdleConns: number // default to core.DefaultAuxMaxIdleConns
  dbConnect: core.DBConnectFunc // default to core.dbConnect
 }
 interface _new {
  /**
   * New creates a new PocketBase instance with the default configuration.
   * Use [NewWithConfig] if you want to provide a custom configuration.
   * 
   * Note that the application will not be initialized/bootstrapped yet,
   * aka. DB connections, migrations, app settings, etc. will not be accessible.
   * Everything will be initialized when [PocketBase.Start] is executed.
   * If you want to initialize the application before calling [PocketBase.Start],
   * then you'll have to manually call [PocketBase.Bootstrap].
   */
  (): (PocketBase)
 }
 interface newWithConfig {
  /**
   * NewWithConfig creates a new PocketBase instance with the provided config.
   * 
   * Note that the application will not be initialized/bootstrapped yet,
   * aka. DB connections, migrations, app settings, etc. will not be accessible.
   * Everything will be initialized when [PocketBase.Start] is executed.
   * If you want to initialize the application before calling [PocketBase.Start],
   * then you'll have to manually call [PocketBase.Bootstrap].
   */
  (config: Config): (PocketBase)
 }
 interface PocketBase {
  /**
   * Start starts the application, aka. registers the default system
   * commands (serve, superuser, version) and executes pb.RootCmd.
   */
  start(): void
 }
 interface PocketBase {
  /**
   * Execute initializes the application (if not already) and executes
   * the pb.RootCmd with graceful shutdown support.
   * 
   * This method differs from pb.Start() by not registering the default
   * system commands!
   */
  execute(): void
 }
 /**
  * coloredWriter is a small wrapper struct to construct a [color.Color] writter.
  */
 interface coloredWriter {
 }
 interface coloredWriter {
  /**
   * Write writes the p bytes using the colored writer.
   */
  write(p: string|Array<number>): number
 }
}

/**
 * Package template is a thin wrapper around the standard html/template
 * and text/template packages that implements a convenient registry to
 * load and cache templates on the fly concurrently.
 * 
 * It was created to assist the JSVM plugin HTML rendering, but could be used in other Go code.
 * 
 * Example:
 * 
 * ```
 * 	registry := template.NewRegistry()
 * 
 * 	html1, err := registry.LoadFiles(
 * 		// the files set wil be parsed only once and then cached
 * 		"layout.html",
 * 		"content.html",
 * 	).Render(map[string]any{"name": "John"})
 * 
 * 	html2, err := registry.LoadFiles(
 * 		// reuse the already parsed and cached files set
 * 		"layout.html",
 * 		"content.html",
 * 	).Render(map[string]any{"name": "Jane"})
 * ```
 */
namespace template {
 interface newRegistry {
  /**
   * NewRegistry creates and initializes a new templates registry with
   * some defaults (eg. global "raw" template function for unescaped HTML).
   * 
   * Use the Registry.Load* methods to load templates into the registry.
   */
  (): (Registry)
 }
 /**
  * Registry defines a templates registry that is safe to be used by multiple goroutines.
  * 
  * Use the Registry.Load* methods to load templates into the registry.
  */
 interface Registry {
 }
 interface Registry {
  /**
   * AddFuncs registers new global template functions.
   * 
   * The key of each map entry is the function name that will be used in the templates.
   * If a function with the map entry name already exists it will be replaced with the new one.
   * 
   * The value of each map entry is a function that must have either a
   * single return value, or two return values of which the second has type error.
   * 
   * Example:
   * 
   * ```
   * 	r.AddFuncs(map[string]any{
   * 	  "toUpper": func(str string) string {
   * 	      return strings.ToUppser(str)
   * 	  },
   * 	  ...
   * 	})
   * ```
   */
  addFuncs(funcs: _TygojaDict): (Registry)
 }
 interface Registry {
  /**
   * LoadFiles caches (if not already) the specified filenames set as a
   * single template and returns a ready to use Renderer instance.
   * 
   * There must be at least 1 filename specified.
   */
  loadFiles(...filenames: string[]): (Renderer)
 }
 interface Registry {
  /**
   * LoadString caches (if not already) the specified inline string as a
   * single template and returns a ready to use Renderer instance.
   */
  loadString(text: string): (Renderer)
 }
 interface Registry {
  /**
   * LoadFS caches (if not already) the specified fs and globPatterns
   * pair as single template and returns a ready to use Renderer instance.
   * 
   * There must be at least 1 file matching the provided globPattern(s)
   * (note that most file names serves as glob patterns matching themselves).
   */
  loadFS(fsys: fs.FS, ...globPatterns: string[]): (Renderer)
 }
 /**
  * Renderer defines a single parsed template.
  */
 interface Renderer {
 }
 interface Renderer {
  /**
   * Render executes the template with the specified data as the dot object
   * and returns the result as plain string.
   */
  render(data: any): string
 }
}

/**
 * Package sync provides basic synchronization primitives such as mutual
 * exclusion locks. Other than the [Once] and [WaitGroup] types, most are intended
 * for use by low-level library routines. Higher-level synchronization is
 * better done via channels and communication.
 * 
 * Values containing the types defined in this package should not be copied.
 */
namespace sync {
 // @ts-ignore
 import isync = sync
 /**
  * A Mutex is a mutual exclusion lock.
  * The zero value for a Mutex is an unlocked mutex.
  * 
  * A Mutex must not be copied after first use.
  * 
  * In the terminology of [the Go memory model],
  * the n'th call to [Mutex.Unlock] “synchronizes before” the m'th call to [Mutex.Lock]
  * for any n < m.
  * A successful call to [Mutex.TryLock] is equivalent to a call to Lock.
  * A failed call to TryLock does not establish any “synchronizes before”
  * relation at all.
  * 
  * [the Go memory model]: https://go.dev/ref/mem
  */
 interface Mutex {
 }
 interface Mutex {
  /**
   * Lock locks m.
   * If the lock is already in use, the calling goroutine
   * blocks until the mutex is available.
   */
  lock(): void
 }
 interface Mutex {
  /**
   * TryLock tries to lock m and reports whether it succeeded.
   * 
   * Note that while correct uses of TryLock do exist, they are rare,
   * and use of TryLock is often a sign of a deeper problem
   * in a particular use of mutexes.
   */
  tryLock(): boolean
 }
 interface Mutex {
  /**
   * Unlock unlocks m.
   * It is a run-time error if m is not locked on entry to Unlock.
   * 
   * A locked [Mutex] is not associated with a particular goroutine.
   * It is allowed for one goroutine to lock a Mutex and then
   * arrange for another goroutine to unlock it.
   */
  unlock(): void
 }
 /**
  * A RWMutex is a reader/writer mutual exclusion lock.
  * The lock can be held by an arbitrary number of readers or a single writer.
  * The zero value for a RWMutex is an unlocked mutex.
  * 
  * A RWMutex must not be copied after first use.
  * 
  * If any goroutine calls [RWMutex.Lock] while the lock is already held by
  * one or more readers, concurrent calls to [RWMutex.RLock] will block until
  * the writer has acquired (and released) the lock, to ensure that
  * the lock eventually becomes available to the writer.
  * Note that this prohibits recursive read-locking.
  * A [RWMutex.RLock] cannot be upgraded into a [RWMutex.Lock],
  * nor can a [RWMutex.Lock] be downgraded into a [RWMutex.RLock].
  * 
  * In the terminology of [the Go memory model],
  * the n'th call to [RWMutex.Unlock] “synchronizes before” the m'th call to Lock
  * for any n < m, just as for [Mutex].
  * For any call to RLock, there exists an n such that
  * the n'th call to Unlock “synchronizes before” that call to RLock,
  * and the corresponding call to [RWMutex.RUnlock] “synchronizes before”
  * the n+1'th call to Lock.
  * 
  * [the Go memory model]: https://go.dev/ref/mem
  */
 interface RWMutex {
 }
 interface RWMutex {
  /**
   * RLock locks rw for reading.
   * 
   * It should not be used for recursive read locking; a blocked Lock
   * call excludes new readers from acquiring the lock. See the
   * documentation on the [RWMutex] type.
   */
  rLock(): void
 }
 interface RWMutex {
  /**
   * TryRLock tries to lock rw for reading and reports whether it succeeded.
   * 
   * Note that while correct uses of TryRLock do exist, they are rare,
   * and use of TryRLock is often a sign of a deeper problem
   * in a particular use of mutexes.
   */
  tryRLock(): boolean
 }
 interface RWMutex {
  /**
   * RUnlock undoes a single [RWMutex.RLock] call;
   * it does not affect other simultaneous readers.
   * It is a run-time error if rw is not locked for reading
   * on entry to RUnlock.
   */
  rUnlock(): void
 }
 interface RWMutex {
  /**
   * Lock locks rw for writing.
   * If the lock is already locked for reading or writing,
   * Lock blocks until the lock is available.
   */
  lock(): void
 }
 interface RWMutex {
  /**
   * TryLock tries to lock rw for writing and reports whether it succeeded.
   * 
   * Note that while correct uses of TryLock do exist, they are rare,
   * and use of TryLock is often a sign of a deeper problem
   * in a particular use of mutexes.
   */
  tryLock(): boolean
 }
 interface RWMutex {
  /**
   * Unlock unlocks rw for writing. It is a run-time error if rw is
   * not locked for writing on entry to Unlock.
   * 
   * As with Mutexes, a locked [RWMutex] is not associated with a particular
   * goroutine. One goroutine may [RWMutex.RLock] ([RWMutex.Lock]) a RWMutex and then
   * arrange for another goroutine to [RWMutex.RUnlock] ([RWMutex.Unlock]) it.
   */
  unlock(): void
 }
 interface RWMutex {
  /**
   * RLocker returns a [Locker] interface that implements
   * the [Locker.Lock] and [Locker.Unlock] methods by calling rw.RLock and rw.RUnlock.
   */
  rLocker(): Locker
 }
}

/**
 * Package syscall contains an interface to the low-level operating system
 * primitives. The details vary depending on the underlying system, and
 * by default, godoc will display the syscall documentation for the current
 * system. If you want godoc to display syscall documentation for another
 * system, set $GOOS and $GOARCH to the desired system. For example, if
 * you want to view documentation for freebsd/arm on linux/amd64, set $GOOS
 * to freebsd and $GOARCH to arm.
 * The primary use of syscall is inside other packages that provide a more
 * portable interface to the system, such as "os", "time" and "net".  Use
 * those packages rather than this one if you can.
 * For details of the functions and data types in this package consult
 * the manuals for the appropriate operating system.
 * These calls return err == nil to indicate success; otherwise
 * err is an operating system error describing the failure.
 * On most systems, that error has type [Errno].
 * 
 * NOTE: Most of the functions, types, and constants defined in
 * this package are also available in the [golang.org/x/sys] package.
 * That package has more system call support than this one,
 * and most new code should prefer that package where possible.
 * See https://golang.org/s/go1.4-syscall for more information.
 */
namespace syscall {
 // @ts-ignore
 import errpkg = errors
 interface SysProcAttr {
  chroot: string // Chroot.
  credential?: Credential // Credential.
  /**
   * Ptrace tells the child to call ptrace(PTRACE_TRACEME).
   * Call runtime.LockOSThread before starting a process with this set,
   * and don't call UnlockOSThread until done with PtraceSyscall calls.
   */
  ptrace: boolean
  setsid: boolean // Create session.
  /**
   * Setpgid sets the process group ID of the child to Pgid,
   * or, if Pgid == 0, to the new child's process ID.
   */
  setpgid: boolean
  /**
   * Setctty sets the controlling terminal of the child to
   * file descriptor Ctty. Ctty must be a descriptor number
   * in the child process: an index into ProcAttr.Files.
   * This is only meaningful if Setsid is true.
   */
  setctty: boolean
  noctty: boolean // Detach fd 0 from controlling terminal.
  ctty: number // Controlling TTY fd.
  /**
   * Foreground places the child process group in the foreground.
   * This implies Setpgid. The Ctty field must be set to
   * the descriptor of the controlling TTY.
   * Unlike Setctty, in this case Ctty must be a descriptor
   * number in the parent process.
   */
  foreground: boolean
  pgid: number // Child's process group ID if Setpgid.
  /**
   * Pdeathsig, if non-zero, is a signal that the kernel will send to
   * the child process when the creating thread dies. Note that the signal
   * is sent on thread termination, which may happen before process termination.
   * There are more details at https://go.dev/issue/27505.
   */
  pdeathsig: Signal
  cloneflags: number // Flags for clone calls.
  unshareflags: number // Flags for unshare calls.
  uidMappings: Array<SysProcIDMap> // User ID mappings for user namespaces.
  gidMappings: Array<SysProcIDMap> // Group ID mappings for user namespaces.
  /**
   * GidMappingsEnableSetgroups enabling setgroups syscall.
   * If false, then setgroups syscall will be disabled for the child process.
   * This parameter is no-op if GidMappings == nil. Otherwise for unprivileged
   * users this should be set to false for mappings work.
   */
  gidMappingsEnableSetgroups: boolean
  ambientCaps: Array<number> // Ambient capabilities.
  useCgroupFD: boolean // Whether to make use of the CgroupFD field.
  cgroupFD: number // File descriptor of a cgroup to put the new process into.
  /**
   * PidFD, if not nil, is used to store the pidfd of a child, if the
   * functionality is supported by the kernel, or -1. Note *PidFD is
   * changed only if the process starts successfully.
   */
  pidFD?: number
 }
 // @ts-ignore
 import errorspkg = errors
 /**
  * A RawConn is a raw network connection.
  */
 interface RawConn {
  [key:string]: any;
  /**
   * Control invokes f on the underlying connection's file
   * descriptor or handle.
   * The file descriptor fd is guaranteed to remain valid while
   * f executes but not after f returns.
   */
  control(f: (fd: number) => void): void
  /**
   * Read invokes f on the underlying connection's file
   * descriptor or handle; f is expected to try to read from the
   * file descriptor.
   * If f returns true, Read returns. Otherwise Read blocks
   * waiting for the connection to be ready for reading and
   * tries again repeatedly.
   * The file descriptor is guaranteed to remain valid while f
   * executes but not after f returns.
   */
  read(f: (fd: number) => boolean): void
  /**
   * Write is like Read but for writing.
   */
  write(f: (fd: number) => boolean): void
 }
 // @ts-ignore
 import runtimesyscall = syscall
 /**
  * An Errno is an unsigned number describing an error condition.
  * It implements the error interface. The zero Errno is by convention
  * a non-error, so code to convert from Errno to error should use:
  * 
  * ```
  * 	err = nil
  * 	if errno != 0 {
  * 		err = errno
  * 	}
  * ```
  * 
  * Errno values can be tested against error values using [errors.Is].
  * For example:
  * 
  * ```
  * 	_, _, err := syscall.Syscall(...)
  * 	if errors.Is(err, fs.ErrNotExist) ...
  * ```
  */
 interface Errno extends Number{}
 interface Errno {
  error(): string
 }
 interface Errno {
  is(target: Error): boolean
 }
 interface Errno {
  temporary(): boolean
 }
 interface Errno {
  timeout(): boolean
 }
}

/**
 * Package io provides basic interfaces to I/O primitives.
 * Its primary job is to wrap existing implementations of such primitives,
 * such as those in package os, into shared public interfaces that
 * abstract the functionality, plus some other related primitives.
 * 
 * Because these interfaces and primitives wrap lower-level operations with
 * various implementations, unless otherwise informed clients should not
 * assume they are safe for parallel execution.
 */
namespace io {
 /**
  * Reader is the interface that wraps the basic Read method.
  * 
  * Read reads up to len(p) bytes into p. It returns the number of bytes
  * read (0 <= n <= len(p)) and any error encountered. Even if Read
  * returns n < len(p), it may use all of p as scratch space during the call.
  * If some data is available but not len(p) bytes, Read conventionally
  * returns what is available instead of waiting for more.
  * 
  * When Read encounters an error or end-of-file condition after
  * successfully reading n > 0 bytes, it returns the number of
  * bytes read. It may return the (non-nil) error from the same call
  * or return the error (and n == 0) from a subsequent call.
  * An instance of this general case is that a Reader returning
  * a non-zero number of bytes at the end of the input stream may
  * return either err == EOF or err == nil. The next Read should
  * return 0, EOF.
  * 
  * Callers should always process the n > 0 bytes returned before
  * considering the error err. Doing so correctly handles I/O errors
  * that happen after reading some bytes and also both of the
  * allowed EOF behaviors.
  * 
  * If len(p) == 0, Read should always return n == 0. It may return a
  * non-nil error if some error condition is known, such as EOF.
  * 
  * Implementations of Read are discouraged from returning a
  * zero byte count with a nil error, except when len(p) == 0.
  * Callers should treat a return of 0 and nil as indicating that
  * nothing happened; in particular it does not indicate EOF.
  * 
  * Implementations must not retain p.
  */
 interface Reader {
  [key:string]: any;
  read(p: string|Array<number>): number
 }
 /**
  * Writer is the interface that wraps the basic Write method.
  * 
  * Write writes len(p) bytes from p to the underlying data stream.
  * It returns the number of bytes written from p (0 <= n <= len(p))
  * and any error encountered that caused the write to stop early.
  * Write must return a non-nil error if it returns n < len(p).
  * Write must not modify the slice data, even temporarily.
  * 
  * Implementations must not retain p.
  */
 interface Writer {
  [key:string]: any;
  write(p: string|Array<number>): number
 }
 /**
  * ReadCloser is the interface that groups the basic Read and Close methods.
  */
 interface ReadCloser {
  [key:string]: any;
 }
 /**
  * ReadSeekCloser is the interface that groups the basic Read, Seek and Close
  * methods.
  */
 interface ReadSeekCloser {
  [key:string]: any;
 }
}

/**
 * Package bytes implements functions for the manipulation of byte slices.
 * It is analogous to the facilities of the [strings] package.
 */
namespace bytes {
 /**
  * A Reader implements the [io.Reader], [io.ReaderAt], [io.WriterTo], [io.Seeker],
  * [io.ByteScanner], and [io.RuneScanner] interfaces by reading from
  * a byte slice.
  * Unlike a [Buffer], a Reader is read-only and supports seeking.
  * The zero value for Reader operates like a Reader of an empty slice.
  */
 interface Reader {
 }
 interface Reader {
  /**
   * Len returns the number of bytes of the unread portion of the
   * slice.
   */
  len(): number
 }
 interface Reader {
  /**
   * Size returns the original length of the underlying byte slice.
   * Size is the number of bytes available for reading via [Reader.ReadAt].
   * The result is unaffected by any method calls except [Reader.Reset].
   */
  size(): number
 }
 interface Reader {
  /**
   * Read implements the [io.Reader] interface.
   */
  read(b: string|Array<number>): number
 }
 interface Reader {
  /**
   * ReadAt implements the [io.ReaderAt] interface.
   */
  readAt(b: string|Array<number>, off: number): number
 }
 interface Reader {
  /**
   * ReadByte implements the [io.ByteReader] interface.
   */
  readByte(): number
 }
 interface Reader {
  /**
   * UnreadByte complements [Reader.ReadByte] in implementing the [io.ByteScanner] interface.
   */
  unreadByte(): void
 }
 interface Reader {
  /**
   * ReadRune implements the [io.RuneReader] interface.
   */
  readRune(): [number, number]
 }
 interface Reader {
  /**
   * UnreadRune complements [Reader.ReadRune] in implementing the [io.RuneScanner] interface.
   */
  unreadRune(): void
 }
 interface Reader {
  /**
   * Seek implements the [io.Seeker] interface.
   */
  seek(offset: number, whence: number): number
 }
 interface Reader {
  /**
   * WriteTo implements the [io.WriterTo] interface.
   */
  writeTo(w: io.Writer): number
 }
 interface Reader {
  /**
   * Reset resets the [Reader] to be reading from b.
   */
  reset(b: string|Array<number>): void
 }
}

/**
 * Package time provides functionality for measuring and displaying time.
 * 
 * The calendrical calculations always assume a Gregorian calendar, with
 * no leap seconds.
 * 
 * # Monotonic Clocks
 * 
 * Operating systems provide both a “wall clock,” which is subject to
 * changes for clock synchronization, and a “monotonic clock,” which is
 * not. The general rule is that the wall clock is for telling time and
 * the monotonic clock is for measuring time. Rather than split the API,
 * in this package the Time returned by [time.Now] contains both a wall
 * clock reading and a monotonic clock reading; later time-telling
 * operations use the wall clock reading, but later time-measuring
 * operations, specifically comparisons and subtractions, use the
 * monotonic clock reading.
 * 
 * For example, this code always computes a positive elapsed time of
 * approximately 20 milliseconds, even if the wall clock is changed during
 * the operation being timed:
 * 
 * ```
 * 	start := time.Now()
 * 	... operation that takes 20 milliseconds ...
 * 	t := time.Now()
 * 	elapsed := t.Sub(start)
 * ```
 * 
 * Other idioms, such as [time.Since](start), [time.Until](deadline), and
 * time.Now().Before(deadline), are similarly robust against wall clock
 * resets.
 * 
 * The rest of this section gives the precise details of how operations
 * use monotonic clocks, but understanding those details is not required
 * to use this package.
 * 
 * The Time returned by time.Now contains a monotonic clock reading.
 * If Time t has a monotonic clock reading, t.Add adds the same duration to
 * both the wall clock and monotonic clock readings to compute the result.
 * Because t.AddDate(y, m, d), t.Round(d), and t.Truncate(d) are wall time
 * computations, they always strip any monotonic clock reading from their results.
 * Because t.In, t.Local, and t.UTC are used for their effect on the interpretation
 * of the wall time, they also strip any monotonic clock reading from their results.
 * The canonical way to strip a monotonic clock reading is to use t = t.Round(0).
 * 
 * If Times t and u both contain monotonic clock readings, the operations
 * t.After(u), t.Before(u), t.Equal(u), t.Compare(u), and t.Sub(u) are carried out
 * using the monotonic clock readings alone, ignoring the wall clock
 * readings. If either t or u contains no monotonic clock reading, these
 * operations fall back to using the wall clock readings.
 * 
 * On some systems the monotonic clock will stop if the computer goes to sleep.
 * On such a system, t.Sub(u) may not accurately reflect the actual
 * time that passed between t and u. The same applies to other functions and
 * methods that subtract times, such as [Since], [Until], [Time.Before], [Time.After],
 * [Time.Add], [Time.Equal] and [Time.Compare]. In some cases, you may need to strip
 * the monotonic clock to get accurate results.
 * 
 * Because the monotonic clock reading has no meaning outside
 * the current process, the serialized forms generated by t.GobEncode,
 * t.MarshalBinary, t.MarshalJSON, and t.MarshalText omit the monotonic
 * clock reading, and t.Format provides no format for it. Similarly, the
 * constructors [time.Date], [time.Parse], [time.ParseInLocation], and [time.Unix],
 * as well as the unmarshalers t.GobDecode, t.UnmarshalBinary.
 * t.UnmarshalJSON, and t.UnmarshalText always create times with
 * no monotonic clock reading.
 * 
 * The monotonic clock reading exists only in [Time] values. It is not
 * a part of [Duration] values or the Unix times returned by t.Unix and
 * friends.
 * 
 * Note that the Go == operator compares not just the time instant but
 * also the [Location] and the monotonic clock reading. See the
 * documentation for the Time type for a discussion of equality
 * testing for Time values.
 * 
 * For debugging, the result of t.String does include the monotonic
 * clock reading if present. If t != u because of different monotonic clock readings,
 * that difference will be visible when printing t.String() and u.String().
 * 
 * # Timer Resolution
 * 
 * [Timer] resolution varies depending on the Go runtime, the operating system
 * and the underlying hardware.
 * On Unix, the resolution is ~1ms.
 * On Windows version 1803 and newer, the resolution is ~0.5ms.
 * On older Windows versions, the default resolution is ~16ms, but
 * a higher resolution may be requested using [golang.org/x/sys/windows.TimeBeginPeriod].
 */
namespace time {
 interface Time {
  /**
   * String returns the time formatted using the format string
   * 
   * ```
   * 	"2006-01-02 15:04:05.999999999 -0700 MST"
   * ```
   * 
   * If the time has a monotonic clock reading, the returned string
   * includes a final field "m=±<value>", where value is the monotonic
   * clock reading formatted as a decimal number of seconds.
   * 
   * The returned string is meant for debugging; for a stable serialized
   * representation, use t.MarshalText, t.MarshalBinary, or t.Format
   * with an explicit format string.
   */
  string(): string
 }
 interface Time {
  /**
   * GoString implements [fmt.GoStringer] and formats t to be printed in Go source
   * code.
   */
  goString(): string
 }
 interface Time {
  /**
   * Format returns a textual representation of the time value formatted according
   * to the layout defined by the argument. See the documentation for the
   * constant called [Layout] to see how to represent the layout format.
   * 
   * The executable example for [Time.Format] demonstrates the working
   * of the layout string in detail and is a good reference.
   */
  format(layout: string): string
 }
 interface Time {
  /**
   * AppendFormat is like [Time.Format] but appends the textual
   * representation to b and returns the extended buffer.
   */
  appendFormat(b: string|Array<number>, layout: string): string|Array<number>
 }
 /**
  * A Time represents an instant in time with nanosecond precision.
  * 
  * Programs using times should typically store and pass them as values,
  * not pointers. That is, time variables and struct fields should be of
  * type [time.Time], not *time.Time.
  * 
  * A Time value can be used by multiple goroutines simultaneously except
  * that the methods [Time.GobDecode], [Time.UnmarshalBinary], [Time.UnmarshalJSON] and
  * [Time.UnmarshalText] are not concurrency-safe.
  * 
  * Time instants can be compared using the [Time.Before], [Time.After], and [Time.Equal] methods.
  * The [Time.Sub] method subtracts two instants, producing a [Duration].
  * The [Time.Add] method adds a Time and a Duration, producing a Time.
  * 
  * The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC.
  * As this time is unlikely to come up in practice, the [Time.IsZero] method gives
  * a simple way of detecting a time that has not been initialized explicitly.
  * 
  * Each time has an associated [Location]. The methods [Time.Local], [Time.UTC], and Time.In return a
  * Time with a specific Location. Changing the Location of a Time value with
  * these methods does not change the actual instant it represents, only the time
  * zone in which to interpret it.
  * 
  * Representations of a Time value saved by the [Time.GobEncode], [Time.MarshalBinary], [Time.AppendBinary],
  * [Time.MarshalJSON], [Time.MarshalText] and [Time.AppendText] methods store the [Time.Location]'s offset,
  * but not the location name. They therefore lose information about Daylight Saving Time.
  * 
  * In addition to the required “wall clock” reading, a Time may contain an optional
  * reading of the current process's monotonic clock, to provide additional precision
  * for comparison or subtraction.
  * See the “Monotonic Clocks” section in the package documentation for details.
  * 
  * Note that the Go == operator compares not just the time instant but also the
  * Location and the monotonic clock reading. Therefore, Time values should not
  * be used as map or database keys without first guaranteeing that the
  * identical Location has been set for all values, which can be achieved
  * through use of the UTC or Local method, and that the monotonic clock reading
  * has been stripped by setting t = t.Round(0). In general, prefer t.Equal(u)
  * to t == u, since t.Equal uses the most accurate comparison available and
  * correctly handles the case when only one of its arguments has a monotonic
  * clock reading.
  */
 interface Time {
 }
 interface Time {
  /**
   * IsZero reports whether t represents the zero time instant,
   * January 1, year 1, 00:00:00 UTC.
   */
  isZero(): boolean
 }
 interface Time {
  /**
   * After reports whether the time instant t is after u.
   */
  after(u: Time): boolean
 }
 interface Time {
  /**
   * Before reports whether the time instant t is before u.
   */
  before(u: Time): boolean
 }
 interface Time {
  /**
   * Compare compares the time instant t with u. If t is before u, it returns -1;
   * if t is after u, it returns +1; if they're the same, it returns 0.
   */
  compare(u: Time): number
 }
 interface Time {
  /**
   * Equal reports whether t and u represent the same time instant.
   * Two times can be equal even if they are in different locations.
   * For example, 6:00 +0200 and 4:00 UTC are Equal.
   * See the documentation on the Time type for the pitfalls of using == with
   * Time values; most code should use Equal instead.
   */
  equal(u: Time): boolean
 }
 interface Time {
  /**
   * Date returns the year, month, and day in which t occurs.
   */
  date(): [number, Month, number]
 }
 interface Time {
  /**
   * Year returns the year in which t occurs.
   */
  year(): number
 }
 interface Time {
  /**
   * Month returns the month of the year specified by t.
   */
  month(): Month
 }
 interface Time {
  /**
   * Day returns the day of the month specified by t.
   */
  day(): number
 }
 interface Time {
  /**
   * Weekday returns the day of the week specified by t.
   */
  weekday(): Weekday
 }
 interface Time {
  /**
   * ISOWeek returns the ISO 8601 year and week number in which t occurs.
   * Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
   * week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
   * of year n+1.
   */
  isoWeek(): [number, number]
 }
 interface Time {
  /**
   * Clock returns the hour, minute, and second within the day specified by t.
   */
  clock(): [number, number, number]
 }
 interface Time {
  /**
   * Hour returns the hour within the day specified by t, in the range [0, 23].
   */
  hour(): number
 }
 interface Time {
  /**
   * Minute returns the minute offset within the hour specified by t, in the range [0, 59].
   */
  minute(): number
 }
 interface Time {
  /**
   * Second returns the second offset within the minute specified by t, in the range [0, 59].
   */
  second(): number
 }
 interface Time {
  /**
   * Nanosecond returns the nanosecond offset within the second specified by t,
   * in the range [0, 999999999].
   */
  nanosecond(): number
 }
 interface Time {
  /**
   * YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years,
   * and [1,366] in leap years.
   */
  yearDay(): number
 }
 /**
  * A Duration represents the elapsed time between two instants
  * as an int64 nanosecond count. The representation limits the
  * largest representable duration to approximately 290 years.
  */
 interface Duration extends Number{}
 interface Duration {
  /**
   * String returns a string representing the duration in the form "72h3m0.5s".
   * Leading zero units are omitted. As a special case, durations less than one
   * second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
   * that the leading digit is non-zero. The zero duration formats as 0s.
   */
  string(): string
 }
 interface Duration {
  /**
   * Nanoseconds returns the duration as an integer nanosecond count.
   */
  nanoseconds(): number
 }
 interface Duration {
  /**
   * Microseconds returns the duration as an integer microsecond count.
   */
  microseconds(): number
 }
 interface Duration {
  /**
   * Milliseconds returns the duration as an integer millisecond count.
   */
  milliseconds(): number
 }
 interface Duration {
  /**
   * Seconds returns the duration as a floating point number of seconds.
   */
  seconds(): number
 }
 interface Duration {
  /**
   * Minutes returns the duration as a floating point number of minutes.
   */
  minutes(): number
 }
 interface Duration {
  /**
   * Hours returns the duration as a floating point number of hours.
   */
  hours(): number
 }
 interface Duration {
  /**
   * Truncate returns the result of rounding d toward zero to a multiple of m.
   * If m <= 0, Truncate returns d unchanged.
   */
  truncate(m: Duration): Duration
 }
 interface Duration {
  /**
   * Round returns the result of rounding d to the nearest multiple of m.
   * The rounding behavior for halfway values is to round away from zero.
   * If the result exceeds the maximum (or minimum)
   * value that can be stored in a [Duration],
   * Round returns the maximum (or minimum) duration.
   * If m <= 0, Round returns d unchanged.
   */
  round(m: Duration): Duration
 }
 interface Duration {
  /**
   * Abs returns the absolute value of d.
   * As a special case, Duration([math.MinInt64]) is converted to Duration([math.MaxInt64]),
   * reducing its magnitude by 1 nanosecond.
   */
  abs(): Duration
 }
 interface Time {
  /**
   * Add returns the time t+d.
   */
  add(d: Duration): Time
 }
 interface Time {
  /**
   * Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
   * value that can be stored in a [Duration], the maximum (or minimum) duration
   * will be returned.
   * To compute t-d for a duration d, use t.Add(-d).
   */
  sub(u: Time): Duration
 }
 interface Time {
  /**
   * AddDate returns the time corresponding to adding the
   * given number of years, months, and days to t.
   * For example, AddDate(-1, 2, 3) applied to January 1, 2011
   * returns March 4, 2010.
   * 
   * Note that dates are fundamentally coupled to timezones, and calendrical
   * periods like days don't have fixed durations. AddDate uses the Location of
   * the Time value to determine these durations. That means that the same
   * AddDate arguments can produce a different shift in absolute time depending on
   * the base Time value and its Location. For example, AddDate(0, 0, 1) applied
   * to 12:00 on March 27 always returns 12:00 on March 28. At some locations and
   * in some years this is a 24 hour shift. In others it's a 23 hour shift due to
   * daylight savings time transitions.
   * 
   * AddDate normalizes its result in the same way that Date does,
   * so, for example, adding one month to October 31 yields
   * December 1, the normalized form for November 31.
   */
  addDate(years: number, months: number, days: number): Time
 }
 interface Time {
  /**
   * UTC returns t with the location set to UTC.
   */
  utc(): Time
 }
 interface Time {
  /**
   * Local returns t with the location set to local time.
   */
  local(): Time
 }
 interface Time {
  /**
   * In returns a copy of t representing the same time instant, but
   * with the copy's location information set to loc for display
   * purposes.
   * 
   * In panics if loc is nil.
   */
  in(loc: Location): Time
 }
 interface Time {
  /**
   * Location returns the time zone information associated with t.
   */
  location(): (Location)
 }
 interface Time {
  /**
   * Zone computes the time zone in effect at time t, returning the abbreviated
   * name of the zone (such as "CET") and its offset in seconds east of UTC.
   */
  zone(): [string, number]
 }
 interface Time {
  /**
   * ZoneBounds returns the bounds of the time zone in effect at time t.
   * The zone begins at start and the next zone begins at end.
   * If the zone begins at the beginning of time, start will be returned as a zero Time.
   * If the zone goes on forever, end will be returned as a zero Time.
   * The Location of the returned times will be the same as t.
   */
  zoneBounds(): [Time, Time]
 }
 interface Time {
  /**
   * Unix returns t as a Unix time, the number of seconds elapsed
   * since January 1, 1970 UTC. The result does not depend on the
   * location associated with t.
   * Unix-like operating systems often record time as a 32-bit
   * count of seconds, but since the method here returns a 64-bit
   * value it is valid for billions of years into the past or future.
   */
  unix(): number
 }
 interface Time {
  /**
   * UnixMilli returns t as a Unix time, the number of milliseconds elapsed since
   * January 1, 1970 UTC. The result is undefined if the Unix time in
   * milliseconds cannot be represented by an int64 (a date more than 292 million
   * years before or after 1970). The result does not depend on the
   * location associated with t.
   */
  unixMilli(): number
 }
 interface Time {
  /**
   * UnixMicro returns t as a Unix time, the number of microseconds elapsed since
   * January 1, 1970 UTC. The result is undefined if the Unix time in
   * microseconds cannot be represented by an int64 (a date before year -290307 or
   * after year 294246). The result does not depend on the location associated
   * with t.
   */
  unixMicro(): number
 }
 interface Time {
  /**
   * UnixNano returns t as a Unix time, the number of nanoseconds elapsed
   * since January 1, 1970 UTC. The result is undefined if the Unix time
   * in nanoseconds cannot be represented by an int64 (a date before the year
   * 1678 or after 2262). Note that this means the result of calling UnixNano
   * on the zero Time is undefined. The result does not depend on the
   * location associated with t.
   */
  unixNano(): number
 }
 interface Time {
  /**
   * AppendBinary implements the [encoding.BinaryAppender] interface.
   */
  appendBinary(b: string|Array<number>): string|Array<number>
 }
 interface Time {
  /**
   * MarshalBinary implements the [encoding.BinaryMarshaler] interface.
   */
  marshalBinary(): string|Array<number>
 }
 interface Time {
  /**
   * UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
   */
  unmarshalBinary(data: string|Array<number>): void
 }
 interface Time {
  /**
   * GobEncode implements the gob.GobEncoder interface.
   */
  gobEncode(): string|Array<number>
 }
 interface Time {
  /**
   * GobDecode implements the gob.GobDecoder interface.
   */
  gobDecode(data: string|Array<number>): void
 }
 interface Time {
  /**
   * MarshalJSON implements the [encoding/json.Marshaler] interface.
   * The time is a quoted string in the RFC 3339 format with sub-second precision.
   * If the timestamp cannot be represented as valid RFC 3339
   * (e.g., the year is out of range), then an error is reported.
   */
  marshalJSON(): string|Array<number>
 }
 interface Time {
  /**
   * UnmarshalJSON implements the [encoding/json.Unmarshaler] interface.
   * The time must be a quoted string in the RFC 3339 format.
   */
  unmarshalJSON(data: string|Array<number>): void
 }
 interface Time {
  /**
   * AppendText implements the [encoding.TextAppender] interface.
   * The time is formatted in RFC 3339 format with sub-second precision.
   * If the timestamp cannot be represented as valid RFC 3339
   * (e.g., the year is out of range), then an error is returned.
   */
  appendText(b: string|Array<number>): string|Array<number>
 }
 interface Time {
  /**
   * MarshalText implements the [encoding.TextMarshaler] interface. The output
   * matches that of calling the [Time.AppendText] method.
   * 
   * See [Time.AppendText] for more information.
   */
  marshalText(): string|Array<number>
 }
 interface Time {
  /**
   * UnmarshalText implements the [encoding.TextUnmarshaler] interface.
   * The time must be in the RFC 3339 format.
   */
  unmarshalText(data: string|Array<number>): void
 }
 interface Time {
  /**
   * IsDST reports whether the time in the configured location is in Daylight Savings Time.
   */
  isDST(): boolean
 }
 interface Time {
  /**
   * Truncate returns the result of rounding t down to a multiple of d (since the zero time).
   * If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
   * 
   * Truncate operates on the time as an absolute duration since the
   * zero time; it does not operate on the presentation form of the
   * time. Thus, Truncate(Hour) may return a time with a non-zero
   * minute, depending on the time's Location.
   */
  truncate(d: Duration): Time
 }
 interface Time {
  /**
   * Round returns the result of rounding t to the nearest multiple of d (since the zero time).
   * The rounding behavior for halfway values is to round up.
   * If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
   * 
   * Round operates on the time as an absolute duration since the
   * zero time; it does not operate on the presentation form of the
   * time. Thus, Round(Hour) may return a time with a non-zero
   * minute, depending on the time's Location.
   */
  round(d: Duration): Time
 }
}

/**
 * Package fs defines basic interfaces to a file system.
 * A file system can be provided by the host operating system
 * but also by other packages.
 * 
 * See the [testing/fstest] package for support with testing
 * implementations of file systems.
 */
namespace fs {
 /**
  * An FS provides access to a hierarchical file system.
  * 
  * The FS interface is the minimum implementation required of the file system.
  * A file system may implement additional interfaces,
  * such as [ReadFileFS], to provide additional or optimized functionality.
  * 
  * [testing/fstest.TestFS] may be used to test implementations of an FS for
  * correctness.
  */
 interface FS {
  [key:string]: any;
  /**
   * Open opens the named file.
   * [File.Close] must be called to release any associated resources.
   * 
   * When Open returns an error, it should be of type *PathError
   * with the Op field set to "open", the Path field set to name,
   * and the Err field describing the problem.
   * 
   * Open should reject attempts to open names that do not satisfy
   * ValidPath(name), returning a *PathError with Err set to
   * ErrInvalid or ErrNotExist.
   */
  open(name: string): File
 }
 /**
  * A File provides access to a single file.
  * The File interface is the minimum implementation required of the file.
  * Directory files should also implement [ReadDirFile].
  * A file may implement [io.ReaderAt] or [io.Seeker] as optimizations.
  */
 interface File {
  [key:string]: any;
  stat(): FileInfo
  read(_arg0: string|Array<number>): number
  close(): void
 }
 /**
  * A DirEntry is an entry read from a directory
  * (using the [ReadDir] function or a [ReadDirFile]'s ReadDir method).
  */
 interface DirEntry {
  [key:string]: any;
  /**
   * Name returns the name of the file (or subdirectory) described by the entry.
   * This name is only the final element of the path (the base name), not the entire path.
   * For example, Name would return "hello.go" not "home/gopher/hello.go".
   */
  name(): string
  /**
   * IsDir reports whether the entry describes a directory.
   */
  isDir(): boolean
  /**
   * Type returns the type bits for the entry.
   * The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
   */
  type(): FileMode
  /**
   * Info returns the FileInfo for the file or subdirectory described by the entry.
   * The returned FileInfo may be from the time of the original directory read
   * or from the time of the call to Info. If the file has been removed or renamed
   * since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
   * If the entry denotes a symbolic link, Info reports the information about the link itself,
   * not the link's target.
   */
  info(): FileInfo
 }
 /**
  * A FileInfo describes a file and is returned by [Stat].
  */
 interface FileInfo {
  [key:string]: any;
  name(): string // base name of the file
  size(): number // length in bytes for regular files; system-dependent for others
  mode(): FileMode // file mode bits
  modTime(): time.Time // modification time
  isDir(): boolean // abbreviation for Mode().IsDir()
  sys(): any // underlying data source (can return nil)
 }
 /**
  * A FileMode represents a file's mode and permission bits.
  * The bits have the same definition on all systems, so that
  * information about files can be moved from one system
  * to another portably. Not all bits apply to all systems.
  * The only required bit is [ModeDir] for directories.
  */
 interface FileMode extends Number{}
 interface FileMode {
  string(): string
 }
 interface FileMode {
  /**
   * IsDir reports whether m describes a directory.
   * That is, it tests for the [ModeDir] bit being set in m.
   */
  isDir(): boolean
 }
 interface FileMode {
  /**
   * IsRegular reports whether m describes a regular file.
   * That is, it tests that no mode type bits are set.
   */
  isRegular(): boolean
 }
 interface FileMode {
  /**
   * Perm returns the Unix permission bits in m (m & [ModePerm]).
   */
  perm(): FileMode
 }
 interface FileMode {
  /**
   * Type returns type bits in m (m & [ModeType]).
   */
  type(): FileMode
 }
 /**
  * PathError records an error and the operation and file path that caused it.
  */
 interface PathError {
  op: string
  path: string
  err: Error
 }
 interface PathError {
  error(): string
 }
 interface PathError {
  unwrap(): void
 }
 interface PathError {
  /**
   * Timeout reports whether this error represents a timeout.
   */
  timeout(): boolean
 }
 /**
  * WalkDirFunc is the type of the function called by [WalkDir] to visit
  * each file or directory.
  * 
  * The path argument contains the argument to [WalkDir] as a prefix.
  * That is, if WalkDir is called with root argument "dir" and finds a file
  * named "a" in that directory, the walk function will be called with
  * argument "dir/a".
  * 
  * The d argument is the [DirEntry] for the named path.
  * 
  * The error result returned by the function controls how [WalkDir]
  * continues. If the function returns the special value [SkipDir], WalkDir
  * skips the current directory (path if d.IsDir() is true, otherwise
  * path's parent directory). If the function returns the special value
  * [SkipAll], WalkDir skips all remaining files and directories. Otherwise,
  * if the function returns a non-nil error, WalkDir stops entirely and
  * returns that error.
  * 
  * The err argument reports an error related to path, signaling that
  * [WalkDir] will not walk into that directory. The function can decide how
  * to handle that error; as described earlier, returning the error will
  * cause WalkDir to stop walking the entire tree.
  * 
  * [WalkDir] calls the function with a non-nil err argument in two cases.
  * 
  * First, if the initial [Stat] on the root directory fails, WalkDir
  * calls the function with path set to root, d set to nil, and err set to
  * the error from [fs.Stat].
  * 
  * Second, if a directory's ReadDir method (see [ReadDirFile]) fails, WalkDir calls the
  * function with path set to the directory's path, d set to an
  * [DirEntry] describing the directory, and err set to the error from
  * ReadDir. In this second case, the function is called twice with the
  * path of the directory: the first call is before the directory read is
  * attempted and has err set to nil, giving the function a chance to
  * return [SkipDir] or [SkipAll] and avoid the ReadDir entirely. The second call
  * is after a failed ReadDir and reports the error from ReadDir.
  * (If ReadDir succeeds, there is no second call.)
  * 
  * The differences between WalkDirFunc compared to [path/filepath.WalkFunc] are:
  * 
  * ```
  *   - The second argument has type [DirEntry] instead of [FileInfo].
  *   - The function is called before reading a directory, to allow [SkipDir]
  *     or [SkipAll] to bypass the directory read entirely or skip all remaining
  *     files and directories respectively.
  *   - If a directory read fails, the function is called a second time
  *     for that directory to report the error.
  * ```
  */
 interface WalkDirFunc {(path: string, d: DirEntry, err: Error): void }
}

/**
 * Package context defines the Context type, which carries deadlines,
 * cancellation signals, and other request-scoped values across API boundaries
 * and between processes.
 * 
 * Incoming requests to a server should create a [Context], and outgoing
 * calls to servers should accept a Context. The chain of function
 * calls between them must propagate the Context, optionally replacing
 * it with a derived Context created using [WithCancel], [WithDeadline],
 * [WithTimeout], or [WithValue].
 * 
 * A Context may be canceled to indicate that work done on its behalf should stop.
 * A Context with a deadline is canceled after the deadline passes.
 * When a Context is canceled, all Contexts derived from it are also canceled.
 * 
 * The [WithCancel], [WithDeadline], and [WithTimeout] functions take a
 * Context (the parent) and return a derived Context (the child) and a
 * [CancelFunc]. Calling the CancelFunc directly cancels the child and its
 * children, removes the parent's reference to the child, and stops
 * any associated timers. Failing to call the CancelFunc leaks the
 * child and its children until the parent is canceled. The go vet tool
 * checks that CancelFuncs are used on all control-flow paths.
 * 
 * The [WithCancelCause], [WithDeadlineCause], and [WithTimeoutCause] functions
 * return a [CancelCauseFunc], which takes an error and records it as
 * the cancellation cause. Calling [Cause] on the canceled context
 * or any of its children retrieves the cause. If no cause is specified,
 * Cause(ctx) returns the same value as ctx.Err().
 * 
 * Programs that use Contexts should follow these rules to keep interfaces
 * consistent across packages and enable static analysis tools to check context
 * propagation:
 * 
 * Do not store Contexts inside a struct type; instead, pass a Context
 * explicitly to each function that needs it. This is discussed further in
 * https://go.dev/blog/context-and-structs. The Context should be the first
 * parameter, typically named ctx:
 * 
 * ```
 * 	func DoSomething(ctx context.Context, arg Arg) error {
 * 		// ... use ctx ...
 * 	}
 * ```
 * 
 * Do not pass a nil [Context], even if a function permits it. Pass [context.TODO]
 * if you are unsure about which Context to use.
 * 
 * Use context Values only for request-scoped data that transits processes and
 * APIs, not for passing optional parameters to functions.
 * 
 * The same Context may be passed to functions running in different goroutines;
 * Contexts are safe for simultaneous use by multiple goroutines.
 * 
 * See https://go.dev/blog/context for example code for a server that uses
 * Contexts.
 */
namespace context {
 /**
  * A Context carries a deadline, a cancellation signal, and other values across
  * API boundaries.
  * 
  * Context's methods may be called by multiple goroutines simultaneously.
  */
 interface Context {
  [key:string]: any;
  /**
   * Deadline returns the time when work done on behalf of this context
   * should be canceled. Deadline returns ok==false when no deadline is
   * set. Successive calls to Deadline return the same results.
   */
  deadline(): [time.Time, boolean]
  /**
   * Done returns a channel that's closed when work done on behalf of this
   * context should be canceled. Done may return nil if this context can
   * never be canceled. Successive calls to Done return the same value.
   * The close of the Done channel may happen asynchronously,
   * after the cancel function returns.
   * 
   * WithCancel arranges for Done to be closed when cancel is called;
   * WithDeadline arranges for Done to be closed when the deadline
   * expires; WithTimeout arranges for Done to be closed when the timeout
   * elapses.
   * 
   * Done is provided for use in select statements:
   * 
   *  // Stream generates values with DoSomething and sends them to out
   *  // until DoSomething returns an error or ctx.Done is closed.
   *  func Stream(ctx context.Context, out chan<- Value) error {
   *  	for {
   *  		v, err := DoSomething(ctx)
   *  		if err != nil {
   *  			return err
   *  		}
   *  		select {
   *  		case <-ctx.Done():
   *  			return ctx.Err()
   *  		case out <- v:
   *  		}
   *  	}
   *  }
   * 
   * See https://blog.golang.org/pipelines for more examples of how to use
   * a Done channel for cancellation.
   */
  done(): undefined
  /**
   * If Done is not yet closed, Err returns nil.
   * If Done is closed, Err returns a non-nil error explaining why:
   * DeadlineExceeded if the context's deadline passed,
   * or Canceled if the context was canceled for some other reason.
   * After Err returns a non-nil error, successive calls to Err return the same error.
   */
  err(): void
  /**
   * Value returns the value associated with this context for key, or nil
   * if no value is associated with key. Successive calls to Value with
   * the same key returns the same result.
   * 
   * Use context values only for request-scoped data that transits
   * processes and API boundaries, not for passing optional parameters to
   * functions.
   * 
   * A key identifies a specific value in a Context. Functions that wish
   * to store values in Context typically allocate a key in a global
   * variable then use that key as the argument to context.WithValue and
   * Context.Value. A key can be any type that supports equality;
   * packages should define keys as an unexported type to avoid
   * collisions.
   * 
   * Packages that define a Context key should provide type-safe accessors
   * for the values stored using that key:
   * 
   * ```
   * 	// Package user defines a User type that's stored in Contexts.
   * 	package user
   * 
   * 	import "context"
   * 
   * 	// User is the type of value stored in the Contexts.
   * 	type User struct {...}
   * 
   * 	// key is an unexported type for keys defined in this package.
   * 	// This prevents collisions with keys defined in other packages.
   * 	type key int
   * 
   * 	// userKey is the key for user.User values in Contexts. It is
   * 	// unexported; clients use user.NewContext and user.FromContext
   * 	// instead of using this key directly.
   * 	var userKey key
   * 
   * 	// NewContext returns a new Context that carries value u.
   * 	func NewContext(ctx context.Context, u *User) context.Context {
   * 		return context.WithValue(ctx, userKey, u)
   * 	}
   * 
   * 	// FromContext returns the User value stored in ctx, if any.
   * 	func FromContext(ctx context.Context) (*User, bool) {
   * 		u, ok := ctx.Value(userKey).(*User)
   * 		return u, ok
   * 	}
   * ```
   */
  value(key: any): any
 }
}

/**
 * Package net provides a portable interface for network I/O, including
 * TCP/IP, UDP, domain name resolution, and Unix domain sockets.
 * 
 * Although the package provides access to low-level networking
 * primitives, most clients will need only the basic interface provided
 * by the [Dial], [Listen], and Accept functions and the associated
 * [Conn] and [Listener] interfaces. The crypto/tls package uses
 * the same interfaces and similar Dial and Listen functions.
 * 
 * The Dial function connects to a server:
 * 
 * ```
 * 	conn, err := net.Dial("tcp", "golang.org:80")
 * 	if err != nil {
 * 		// handle error
 * 	}
 * 	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
 * 	status, err := bufio.NewReader(conn).ReadString('\n')
 * 	// ...
 * ```
 * 
 * The Listen function creates servers:
 * 
 * ```
 * 	ln, err := net.Listen("tcp", ":8080")
 * 	if err != nil {
 * 		// handle error
 * 	}
 * 	for {
 * 		conn, err := ln.Accept()
 * 		if err != nil {
 * 			// handle error
 * 		}
 * 		go handleConnection(conn)
 * 	}
 * ```
 * 
 * # Name Resolution
 * 
 * The method for resolving domain names, whether indirectly with functions like Dial
 * or directly with functions like [LookupHost] and [LookupAddr], varies by operating system.
 * 
 * On Unix systems, the resolver has two options for resolving names.
 * It can use a pure Go resolver that sends DNS requests directly to the servers
 * listed in /etc/resolv.conf, or it can use a cgo-based resolver that calls C
 * library routines such as getaddrinfo and getnameinfo.
 * 
 * On Unix the pure Go resolver is preferred over the cgo resolver, because a blocked DNS
 * request consumes only a goroutine, while a blocked C call consumes an operating system thread.
 * When cgo is available, the cgo-based resolver is used instead under a variety of
 * conditions: on systems that do not let programs make direct DNS requests (OS X),
 * when the LOCALDOMAIN environment variable is present (even if empty),
 * when the RES_OPTIONS or HOSTALIASES environment variable is non-empty,
 * when the ASR_CONFIG environment variable is non-empty (OpenBSD only),
 * when /etc/resolv.conf or /etc/nsswitch.conf specify the use of features that the
 * Go resolver does not implement.
 * 
 * On all systems (except Plan 9), when the cgo resolver is being used
 * this package applies a concurrent cgo lookup limit to prevent the system
 * from running out of system threads. Currently, it is limited to 500 concurrent lookups.
 * 
 * The resolver decision can be overridden by setting the netdns value of the
 * GODEBUG environment variable (see package runtime) to go or cgo, as in:
 * 
 * ```
 * 	export GODEBUG=netdns=go    # force pure Go resolver
 * 	export GODEBUG=netdns=cgo   # force native resolver (cgo, win32)
 * ```
 * 
 * The decision can also be forced while building the Go source tree
 * by setting the netgo or netcgo build tag.
 * The netgo build tag disables entirely the use of the native (CGO) resolver,
 * meaning the Go resolver is the only one that can be used.
 * With the netcgo build tag the native and the pure Go resolver are compiled into the binary,
 * but the native (CGO) resolver is preferred over the Go resolver.
 * With netcgo, the Go resolver can still be forced at runtime with GODEBUG=netdns=go.
 * 
 * A numeric netdns setting, as in GODEBUG=netdns=1, causes the resolver
 * to print debugging information about its decisions.
 * To force a particular resolver while also printing debugging information,
 * join the two settings by a plus sign, as in GODEBUG=netdns=go+1.
 * 
 * The Go resolver will send an EDNS0 additional header with a DNS request,
 * to signal a willingness to accept a larger DNS packet size.
 * This can reportedly cause sporadic failures with the DNS server run
 * by some modems and routers. Setting GODEBUG=netedns0=0 will disable
 * sending the additional header.
 * 
 * On macOS, if Go code that uses the net package is built with
 * -buildmode=c-archive, linking the resulting archive into a C program
 * requires passing -lresolv when linking the C code.
 * 
 * On Plan 9, the resolver always accesses /net/cs and /net/dns.
 * 
 * On Windows, in Go 1.18.x and earlier, the resolver always used C
 * library functions, such as GetAddrInfo and DnsQuery.
 */
namespace net {
 /**
  * Conn is a generic stream-oriented network connection.
  * 
  * Multiple goroutines may invoke methods on a Conn simultaneously.
  */
 interface Conn {
  [key:string]: any;
  /**
   * Read reads data from the connection.
   * Read can be made to time out and return an error after a fixed
   * time limit; see SetDeadline and SetReadDeadline.
   */
  read(b: string|Array<number>): number
  /**
   * Write writes data to the connection.
   * Write can be made to time out and return an error after a fixed
   * time limit; see SetDeadline and SetWriteDeadline.
   */
  write(b: string|Array<number>): number
  /**
   * Close closes the connection.
   * Any blocked Read or Write operations will be unblocked and return errors.
   */
  close(): void
  /**
   * LocalAddr returns the local network address, if known.
   */
  localAddr(): Addr
  /**
   * RemoteAddr returns the remote network address, if known.
   */
  remoteAddr(): Addr
  /**
   * SetDeadline sets the read and write deadlines associated
   * with the connection. It is equivalent to calling both
   * SetReadDeadline and SetWriteDeadline.
   * 
   * A deadline is an absolute time after which I/O operations
   * fail instead of blocking. The deadline applies to all future
   * and pending I/O, not just the immediately following call to
   * Read or Write. After a deadline has been exceeded, the
   * connection can be refreshed by setting a deadline in the future.
   * 
   * If the deadline is exceeded a call to Read or Write or to other
   * I/O methods will return an error that wraps os.ErrDeadlineExceeded.
   * This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
   * The error's Timeout method will return true, but note that there
   * are other possible errors for which the Timeout method will
   * return true even if the deadline has not been exceeded.
   * 
   * An idle timeout can be implemented by repeatedly extending
   * the deadline after successful Read or Write calls.
   * 
   * A zero value for t means I/O operations will not time out.
   */
  setDeadline(t: time.Time): void
  /**
   * SetReadDeadline sets the deadline for future Read calls
   * and any currently-blocked Read call.
   * A zero value for t means Read will not time out.
   */
  setReadDeadline(t: time.Time): void
  /**
   * SetWriteDeadline sets the deadline for future Write calls
   * and any currently-blocked Write call.
   * Even if write times out, it may return n > 0, indicating that
   * some of the data was successfully written.
   * A zero value for t means Write will not time out.
   */
  setWriteDeadline(t: time.Time): void
 }
 /**
  * A Listener is a generic network listener for stream-oriented protocols.
  * 
  * Multiple goroutines may invoke methods on a Listener simultaneously.
  */
 interface Listener {
  [key:string]: any;
  /**
   * Accept waits for and returns the next connection to the listener.
   */
  accept(): Conn
  /**
   * Close closes the listener.
   * Any blocked Accept operations will be unblocked and return errors.
   */
  close(): void
  /**
   * Addr returns the listener's network address.
   */
  addr(): Addr
 }
}

/**
 * Package jwt is a Go implementation of JSON Web Tokens: http://self-issued.info/docs/draft-jones-json-web-token.html
 * 
 * See README.md for more info.
 */
namespace jwt {
 /**
  * MapClaims is a claims type that uses the map[string]any for JSON
  * decoding. This is the default claims type if you don't supply one
  */
 interface MapClaims extends _TygojaDict{}
 interface MapClaims {
  /**
   * GetExpirationTime implements the Claims interface.
   */
  getExpirationTime(): (NumericDate)
 }
 interface MapClaims {
  /**
   * GetNotBefore implements the Claims interface.
   */
  getNotBefore(): (NumericDate)
 }
 interface MapClaims {
  /**
   * GetIssuedAt implements the Claims interface.
   */
  getIssuedAt(): (NumericDate)
 }
 interface MapClaims {
  /**
   * GetAudience implements the Claims interface.
   */
  getAudience(): ClaimStrings
 }
 interface MapClaims {
  /**
   * GetIssuer implements the Claims interface.
   */
  getIssuer(): string
 }
 interface MapClaims {
  /**
   * GetSubject implements the Claims interface.
   */
  getSubject(): string
 }
}

namespace store {
 /**
  * Store defines a concurrent safe in memory key-value data store.
  */
 interface Store<K,T> {
 }
 interface Store<K, T> {
  /**
   * Reset clears the store and replaces the store data with a
   * shallow copy of the provided newData.
   */
  reset(newData: _TygojaDict): void
 }
 interface Store<K, T> {
  /**
   * Length returns the current number of elements in the store.
   */
  length(): number
 }
 interface Store<K, T> {
  /**
   * RemoveAll removes all the existing store entries.
   */
  removeAll(): void
 }
 interface Store<K, T> {
  /**
   * Remove removes a single entry from the store.
   * 
   * Remove does nothing if key doesn't exist in the store.
   */
  remove(key: K): void
 }
 interface Store<K, T> {
  /**
   * Has checks if element with the specified key exist or not.
   */
  has(key: K): boolean
 }
 interface Store<K, T> {
  /**
   * Get returns a single element value from the store.
   * 
   * If key is not set, the zero T value is returned.
   */
  get(key: K): T
 }
 interface Store<K, T> {
  /**
   * GetOk is similar to Get but returns also a boolean indicating whether the key exists or not.
   */
  getOk(key: K): [T, boolean]
 }
 interface Store<K, T> {
  /**
   * GetAll returns a shallow copy of the current store data.
   */
  getAll(): _TygojaDict
 }
 interface Store<K, T> {
  /**
   * Values returns a slice with all of the current store values.
   */
  values(): Array<T>
 }
 interface Store<K, T> {
  /**
   * Set sets (or overwrite if already exists) a new value for key.
   */
  set(key: K, value: T): void
 }
 interface Store<K, T> {
  /**
   * SetFunc sets (or overwrite if already exists) a new value resolved
   * from the function callback for the provided key.
   * 
   * The function callback receives as argument the old store element value (if exists).
   * If there is no old store element, the argument will be the T zero value.
   * 
   * Example:
   * 
   * ```
   * 	s := store.New[string, int](nil)
   * 	s.SetFunc("count", func(old int) int {
   * 	    return old + 1
   * 	})
   * ```
   */
  setFunc(key: K, fn: (old: T) => T): void
 }
 interface Store<K, T> {
  /**
   * GetOrSet retrieves a single existing value for the provided key
   * or stores a new one if it doesn't exist.
   */
  getOrSet(key: K, setFunc: () => T): T
 }
 interface Store<K, T> {
  /**
   * SetIfLessThanLimit sets (or overwrite if already exist) a new value for key.
   * 
   * This method is similar to Set() but **it will skip adding new elements**
   * to the store if the store length has reached the specified limit.
   * false is returned if maxAllowedElements limit is reached.
   */
  setIfLessThanLimit(key: K, value: T, maxAllowedElements: number): boolean
 }
 interface Store<K, T> {
  /**
   * UnmarshalJSON implements [json.Unmarshaler] and imports the
   * provided JSON data into the store.
   * 
   * The store entries that match with the ones from the data will be overwritten with the new value.
   */
  unmarshalJSON(data: string|Array<number>): void
 }
 interface Store<K, T> {
  /**
   * MarshalJSON implements [json.Marshaler] and export the current
   * store data into valid JSON.
   */
  marshalJSON(): string|Array<number>
 }
}

/**
 * Package syntax parses regular expressions into parse trees and compiles
 * parse trees into programs. Most clients of regular expressions will use the
 * facilities of package [regexp] (such as [regexp.Compile] and [regexp.Match]) instead of this package.
 * 
 * # Syntax
 * 
 * The regular expression syntax understood by this package when parsing with the [Perl] flag is as follows.
 * Parts of the syntax can be disabled by passing alternate flags to [Parse].
 * 
 * Single characters:
 * 
 * ```
 * 	.              any character, possibly including newline (flag s=true)
 * 	[xyz]          character class
 * 	[^xyz]         negated character class
 * 	\d             Perl character class
 * 	\D             negated Perl character class
 * 	[[:alpha:]]    ASCII character class
 * 	[[:^alpha:]]   negated ASCII character class
 * 	\pN            Unicode character class (one-letter name)
 * 	\p{Greek}      Unicode character class
 * 	\PN            negated Unicode character class (one-letter name)
 * 	\P{Greek}      negated Unicode character class
 * ```
 * 
 * Composites:
 * 
 * ```
 * 	xy             x followed by y
 * 	x|y            x or y (prefer x)
 * ```
 * 
 * Repetitions:
 * 
 * ```
 * 	x*             zero or more x, prefer more
 * 	x+             one or more x, prefer more
 * 	x?             zero or one x, prefer one
 * 	x{n,m}         n or n+1 or ... or m x, prefer more
 * 	x{n,}          n or more x, prefer more
 * 	x{n}           exactly n x
 * 	x*?            zero or more x, prefer fewer
 * 	x+?            one or more x, prefer fewer
 * 	x??            zero or one x, prefer zero
 * 	x{n,m}?        n or n+1 or ... or m x, prefer fewer
 * 	x{n,}?         n or more x, prefer fewer
 * 	x{n}?          exactly n x
 * ```
 * 
 * Implementation restriction: The counting forms x{n,m}, x{n,}, and x{n}
 * reject forms that create a minimum or maximum repetition count above 1000.
 * Unlimited repetitions are not subject to this restriction.
 * 
 * Grouping:
 * 
 * ```
 * 	(re)           numbered capturing group (submatch)
 * 	(?P<name>re)   named & numbered capturing group (submatch)
 * 	(?<name>re)    named & numbered capturing group (submatch)
 * 	(?:re)         non-capturing group
 * 	(?flags)       set flags within current group; non-capturing
 * 	(?flags:re)    set flags during re; non-capturing
 * 
 * 	Flag syntax is xyz (set) or -xyz (clear) or xy-z (set xy, clear z). The flags are:
 * 
 * 	i              case-insensitive (default false)
 * 	m              multi-line mode: ^ and $ match begin/end line in addition to begin/end text (default false)
 * 	s              let . match \n (default false)
 * 	U              ungreedy: swap meaning of x* and x*?, x+ and x+?, etc (default false)
 * ```
 * 
 * Empty strings:
 * 
 * ```
 * 	^              at beginning of text or line (flag m=true)
 * 	$              at end of text (like \z not \Z) or line (flag m=true)
 * 	\A             at beginning of text
 * 	\b             at ASCII word boundary (\w on one side and \W, \A, or \z on the other)
 * 	\B             not at ASCII word boundary
 * 	\z             at end of text
 * ```
 * 
 * Escape sequences:
 * 
 * ```
 * 	\a             bell (== \007)
 * 	\f             form feed (== \014)
 * 	\t             horizontal tab (== \011)
 * 	\n             newline (== \012)
 * 	\r             carriage return (== \015)
 * 	\v             vertical tab character (== \013)
 * 	\*             literal *, for any punctuation character *
 * 	\123           octal character code (up to three digits)
 * 	\x7F           hex character code (exactly two digits)
 * 	\x{10FFFF}     hex character code
 * 	\Q...\E        literal text ... even if ... has punctuation
 * ```
 * 
 * Character class elements:
 * 
 * ```
 * 	x              single character
 * 	A-Z            character range (inclusive)
 * 	\d             Perl character class
 * 	[:foo:]        ASCII character class foo
 * 	\p{Foo}        Unicode character class Foo
 * 	\pF            Unicode character class F (one-letter name)
 * ```
 * 
 * Named character classes as character class elements:
 * 
 * ```
 * 	[\d]           digits (== \d)
 * 	[^\d]          not digits (== \D)
 * 	[\D]           not digits (== \D)
 * 	[^\D]          not not digits (== \d)
 * 	[[:name:]]     named ASCII class inside character class (== [:name:])
 * 	[^[:name:]]    named ASCII class inside negated character class (== [:^name:])
 * 	[\p{Name}]     named Unicode property inside character class (== \p{Name})
 * 	[^\p{Name}]    named Unicode property inside negated character class (== \P{Name})
 * ```
 * 
 * Perl character classes (all ASCII-only):
 * 
 * ```
 * 	\d             digits (== [0-9])
 * 	\D             not digits (== [^0-9])
 * 	\s             whitespace (== [\t\n\f\r ])
 * 	\S             not whitespace (== [^\t\n\f\r ])
 * 	\w             word characters (== [0-9A-Za-z_])
 * 	\W             not word characters (== [^0-9A-Za-z_])
 * ```
 * 
 * ASCII character classes:
 * 
 * ```
 * 	[[:alnum:]]    alphanumeric (== [0-9A-Za-z])
 * 	[[:alpha:]]    alphabetic (== [A-Za-z])
 * 	[[:ascii:]]    ASCII (== [\x00-\x7F])
 * 	[[:blank:]]    blank (== [\t ])
 * 	[[:cntrl:]]    control (== [\x00-\x1F\x7F])
 * 	[[:digit:]]    digits (== [0-9])
 * 	[[:graph:]]    graphical (== [!-~] == [A-Za-z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~])
 * 	[[:lower:]]    lower case (== [a-z])
 * 	[[:print:]]    printable (== [ -~] == [ [:graph:]])
 * 	[[:punct:]]    punctuation (== [!-/:-@[-`{-~])
 * 	[[:space:]]    whitespace (== [\t\n\v\f\r ])
 * 	[[:upper:]]    upper case (== [A-Z])
 * 	[[:word:]]     word characters (== [0-9A-Za-z_])
 * 	[[:xdigit:]]   hex digit (== [0-9A-Fa-f])
 * ```
 * 
 * Unicode character classes are those in [unicode.Categories] and [unicode.Scripts].
 */
namespace syntax {
 /**
  * Flags control the behavior of the parser and record information about regexp context.
  */
 interface Flags extends Number{}
}

namespace hook {
 /**
  * Event implements [Resolver] and it is intended to be used as a base
  * Hook event that you can embed in your custom typed event structs.
  * 
  * Example:
  * 
  * ```
  * 	type CustomEvent struct {
  * 		hook.Event
  * 
  * 		SomeField int
  * 	}
  * ```
  */
 interface Event {
 }
 interface Event {
  /**
   * Next calls the next hook handler.
   */
  next(): void
 }
 /**
  * Handler defines a single Hook handler.
  * Multiple handlers can share the same id.
  * If Id is not explicitly set it will be autogenerated by Hook.Add and Hook.AddHandler.
  */
 interface Handler<T> {
  /**
   * Func defines the handler function to execute.
   * 
   * Note that users need to call e.Next() in order to proceed with
   * the execution of the hook chain.
   */
  func: (_arg0: T) => void
  /**
   * Id is the unique identifier of the handler.
   * 
   * It could be used later to remove the handler from a hook via [Hook.Remove].
   * 
   * If missing, an autogenerated value will be assigned when adding
   * the handler to a hook.
   */
  id: string
  /**
   * Priority allows changing the default exec priority of the handler within a hook.
   * 
   * If 0, the handler will be executed in the same order it was registered.
   */
  priority: number
 }
 /**
  * Hook defines a generic concurrent safe structure for managing event hooks.
  * 
  * When using custom event it must embed the base [hook.Event].
  * 
  * Example:
  * 
  * ```
  * 	type CustomEvent struct {
  * 		hook.Event
  * 		SomeField int
  * 	}
  * 
  * 	h := Hook[*CustomEvent]{}
  * 
  * 	h.BindFunc(func(e *CustomEvent) error {
  * 		println(e.SomeField)
  * 
  * 		return e.Next()
  * 	})
  * 
  * 	h.Trigger(&CustomEvent{ SomeField: 123 })
  * ```
  */
 interface Hook<T> {
 }
 interface Hook<T> {
  /**
   * Bind registers the provided handler to the current hooks queue.
   * 
   * If handler.Id is empty it is updated with autogenerated value.
   * 
   * If a handler from the current hook list has Id matching handler.Id
   * then the old handler is replaced with the new one.
   */
  bind(handler: Handler<T>): string
 }
 interface Hook<T> {
  /**
   * BindFunc is similar to Bind but registers a new handler from just the provided function.
   * 
   * The registered handler is added with a default 0 priority and the id will be autogenerated.
   * 
   * If you want to register a handler with custom priority or id use the [Hook.Bind] method.
   */
  bindFunc(fn: (e: T) => void): string
 }
 interface Hook<T> {
  /**
   * Unbind removes one or many hook handler by their id.
   */
  unbind(...idsToRemove: string[]): void
 }
 interface Hook<T> {
  /**
   * UnbindAll removes all registered handlers.
   */
  unbindAll(): void
 }
 interface Hook<T> {
  /**
   * Length returns to total number of registered hook handlers.
   */
  length(): number
 }
 interface Hook<T> {
  /**
   * Trigger executes all registered hook handlers one by one
   * with the specified event as an argument.
   * 
   * Optionally, this method allows also to register additional one off
   * handler funcs that will be temporary appended to the handlers queue.
   * 
   * NB! Each hook handler must call event.Next() in order the hook chain to proceed.
   */
  trigger(event: T, ...oneOffHandlerFuncs: ((_arg0: T) => void)[]): void
 }
 /**
  * TaggedHook defines a proxy hook which register handlers that are triggered only
  * if the TaggedHook.tags are empty or includes at least one of the event data tag(s).
  */
 type _sdlSYRj<T> = mainHook<T>
 interface TaggedHook<T> extends _sdlSYRj<T> {
 }
 interface TaggedHook<T> {
  /**
   * CanTriggerOn checks if the current TaggedHook can be triggered with
   * the provided event data tags.
   * 
   * It returns always true if the hook doens't have any tags.
   */
  canTriggerOn(tagsToCheck: Array<string>): boolean
 }
 interface TaggedHook<T> {
  /**
   * Bind registers the provided handler to the current hooks queue.
   * 
   * It is similar to [Hook.Bind] with the difference that the handler
   * function is invoked only if the event data tags satisfy h.CanTriggerOn.
   */
  bind(handler: Handler<T>): string
 }
 interface TaggedHook<T> {
  /**
   * BindFunc registers a new handler with the specified function.
   * 
   * It is similar to [Hook.Bind] with the difference that the handler
   * function is invoked only if the event data tags satisfy h.CanTriggerOn.
   */
  bindFunc(fn: (e: T) => void): string
 }
}

/**
 * Package cron implements a crontab-like service to execute and schedule
 * repeative tasks/jobs.
 * 
 * Example:
 * 
 * ```
 * 	c := cron.New()
 * 	c.MustAdd("dailyReport", "0 0 * * *", func() { ... })
 * 	c.Start()
 * ```
 */
namespace cron {
 /**
  * Cron is a crontab-like struct for tasks/jobs scheduling.
  */
 interface Cron {
 }
 interface Cron {
  /**
   * SetInterval changes the current cron tick interval
   * (it usually should be >= 1 minute).
   */
  setInterval(d: time.Duration): void
 }
 interface Cron {
  /**
   * SetTimezone changes the current cron tick timezone.
   */
  setTimezone(l: time.Location): void
 }
 interface Cron {
  /**
   * MustAdd is similar to Add() but panic on failure.
   */
  mustAdd(jobId: string, cronExpr: string, run: () => void): void
 }
 interface Cron {
  /**
   * Add registers a single cron job.
   * 
   * If there is already a job with the provided id, then the old job
   * will be replaced with the new one.
   * 
   * cronExpr is a regular cron expression, eg. "0 *\/3 * * *" (aka. at minute 0 past every 3rd hour).
   * Check cron.NewSchedule() for the supported tokens.
   */
  add(jobId: string, cronExpr: string, fn: () => void): void
 }
 interface Cron {
  /**
   * Remove removes a single cron job by its id.
   */
  remove(jobId: string): void
 }
 interface Cron {
  /**
   * RemoveAll removes all registered cron jobs.
   */
  removeAll(): void
 }
 interface Cron {
  /**
   * Total returns the current total number of registered cron jobs.
   */
  total(): number
 }
 interface Cron {
  /**
   * Jobs returns a shallow copy of the currently registered cron jobs.
   */
  jobs(): Array<(Job | undefined)>
 }
 interface Cron {
  /**
   * Stop stops the current cron ticker (if not already).
   * 
   * You can resume the ticker by calling Start().
   */
  stop(): void
 }
 interface Cron {
  /**
   * Start starts the cron ticker.
   * 
   * Calling Start() on already started cron will restart the ticker.
   */
  start(): void
 }
 interface Cron {
  /**
   * HasStarted checks whether the current Cron ticker has been started.
   */
  hasStarted(): boolean
 }
}

/**
 * Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer
 * object, creating another object (Reader or Writer) that also implements
 * the interface but provides buffering and some help for textual I/O.
 */
namespace bufio {
 /**
  * ReadWriter stores pointers to a [Reader] and a [Writer].
  * It implements [io.ReadWriter].
  */
 type _sIymuZL = Reader&Writer
 interface ReadWriter extends _sIymuZL {
 }
}

/**
 * Package sql provides a generic interface around SQL (or SQL-like)
 * databases.
 * 
 * The sql package must be used in conjunction with a database driver.
 * See https://golang.org/s/sqldrivers for a list of drivers.
 * 
 * Drivers that do not support context cancellation will not return until
 * after the query is completed.
 * 
 * For usage examples, see the wiki page at
 * https://golang.org/s/sqlwiki.
 */
namespace sql {
 /**
  * TxOptions holds the transaction options to be used in [DB.BeginTx].
  */
 interface TxOptions {
  /**
   * Isolation is the transaction isolation level.
   * If zero, the driver or database's default level is used.
   */
  isolation: IsolationLevel
  readOnly: boolean
 }
 /**
  * NullString represents a string that may be null.
  * NullString implements the [Scanner] interface so
  * it can be used as a scan destination:
  * 
  * ```
  * 	var s NullString
  * 	err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
  * 	...
  * 	if s.Valid {
  * 	   // use s.String
  * 	} else {
  * 	   // NULL value
  * 	}
  * ```
  */
 interface NullString {
  string: string
  valid: boolean // Valid is true if String is not NULL
 }
 interface NullString {
  /**
   * Scan implements the [Scanner] interface.
   */
  scan(value: any): void
 }
 interface NullString {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 /**
  * DB is a database handle representing a pool of zero or more
  * underlying connections. It's safe for concurrent use by multiple
  * goroutines.
  * 
  * The sql package creates and frees connections automatically; it
  * also maintains a free pool of idle connections. If the database has
  * a concept of per-connection state, such state can be reliably observed
  * within a transaction ([Tx]) or connection ([Conn]). Once [DB.Begin] is called, the
  * returned [Tx] is bound to a single connection. Once [Tx.Commit] or
  * [Tx.Rollback] is called on the transaction, that transaction's
  * connection is returned to [DB]'s idle connection pool. The pool size
  * can be controlled with [DB.SetMaxIdleConns].
  */
 interface DB {
 }
 interface DB {
  /**
   * PingContext verifies a connection to the database is still alive,
   * establishing a connection if necessary.
   */
  pingContext(ctx: context.Context): void
 }
 interface DB {
  /**
   * Ping verifies a connection to the database is still alive,
   * establishing a connection if necessary.
   * 
   * Ping uses [context.Background] internally; to specify the context, use
   * [DB.PingContext].
   */
  ping(): void
 }
 interface DB {
  /**
   * Close closes the database and prevents new queries from starting.
   * Close then waits for all queries that have started processing on the server
   * to finish.
   * 
   * It is rare to Close a [DB], as the [DB] handle is meant to be
   * long-lived and shared between many goroutines.
   */
  close(): void
 }
 interface DB {
  /**
   * SetMaxIdleConns sets the maximum number of connections in the idle
   * connection pool.
   * 
   * If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
   * then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
   * 
   * If n <= 0, no idle connections are retained.
   * 
   * The default max idle connections is currently 2. This may change in
   * a future release.
   */
  setMaxIdleConns(n: number): void
 }
 interface DB {
  /**
   * SetMaxOpenConns sets the maximum number of open connections to the database.
   * 
   * If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
   * MaxIdleConns, then MaxIdleConns will be reduced to match the new
   * MaxOpenConns limit.
   * 
   * If n <= 0, then there is no limit on the number of open connections.
   * The default is 0 (unlimited).
   */
  setMaxOpenConns(n: number): void
 }
 interface DB {
  /**
   * SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
   * 
   * Expired connections may be closed lazily before reuse.
   * 
   * If d <= 0, connections are not closed due to a connection's age.
   */
  setConnMaxLifetime(d: time.Duration): void
 }
 interface DB {
  /**
   * SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
   * 
   * Expired connections may be closed lazily before reuse.
   * 
   * If d <= 0, connections are not closed due to a connection's idle time.
   */
  setConnMaxIdleTime(d: time.Duration): void
 }
 interface DB {
  /**
   * Stats returns database statistics.
   */
  stats(): DBStats
 }
 interface DB {
  /**
   * PrepareContext creates a prepared statement for later queries or executions.
   * Multiple queries or executions may be run concurrently from the
   * returned statement.
   * The caller must call the statement's [*Stmt.Close] method
   * when the statement is no longer needed.
   * 
   * The provided context is used for the preparation of the statement, not for the
   * execution of the statement.
   */
  prepareContext(ctx: context.Context, query: string): (Stmt)
 }
 interface DB {
  /**
   * Prepare creates a prepared statement for later queries or executions.
   * Multiple queries or executions may be run concurrently from the
   * returned statement.
   * The caller must call the statement's [*Stmt.Close] method
   * when the statement is no longer needed.
   * 
   * Prepare uses [context.Background] internally; to specify the context, use
   * [DB.PrepareContext].
   */
  prepare(query: string): (Stmt)
 }
 interface DB {
  /**
   * ExecContext executes a query without returning any rows.
   * The args are for any placeholder parameters in the query.
   */
  execContext(ctx: context.Context, query: string, ...args: any[]): Result
 }
 interface DB {
  /**
   * Exec executes a query without returning any rows.
   * The args are for any placeholder parameters in the query.
   * 
   * Exec uses [context.Background] internally; to specify the context, use
   * [DB.ExecContext].
   */
  exec(query: string, ...args: any[]): Result
 }
 interface DB {
  /**
   * QueryContext executes a query that returns rows, typically a SELECT.
   * The args are for any placeholder parameters in the query.
   */
  queryContext(ctx: context.Context, query: string, ...args: any[]): (Rows)
 }
 interface DB {
  /**
   * Query executes a query that returns rows, typically a SELECT.
   * The args are for any placeholder parameters in the query.
   * 
   * Query uses [context.Background] internally; to specify the context, use
   * [DB.QueryContext].
   */
  query(query: string, ...args: any[]): (Rows)
 }
 interface DB {
  /**
   * QueryRowContext executes a query that is expected to return at most one row.
   * QueryRowContext always returns a non-nil value. Errors are deferred until
   * [Row]'s Scan method is called.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, [*Row.Scan] scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, query: string, ...args: any[]): (Row)
 }
 interface DB {
  /**
   * QueryRow executes a query that is expected to return at most one row.
   * QueryRow always returns a non-nil value. Errors are deferred until
   * [Row]'s Scan method is called.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, [*Row.Scan] scans the first selected row and discards
   * the rest.
   * 
   * QueryRow uses [context.Background] internally; to specify the context, use
   * [DB.QueryRowContext].
   */
  queryRow(query: string, ...args: any[]): (Row)
 }
 interface DB {
  /**
   * BeginTx starts a transaction.
   * 
   * The provided context is used until the transaction is committed or rolled back.
   * If the context is canceled, the sql package will roll back
   * the transaction. [Tx.Commit] will return an error if the context provided to
   * BeginTx is canceled.
   * 
   * The provided [TxOptions] is optional and may be nil if defaults should be used.
   * If a non-default isolation level is used that the driver doesn't support,
   * an error will be returned.
   */
  beginTx(ctx: context.Context, opts: TxOptions): (Tx)
 }
 interface DB {
  /**
   * Begin starts a transaction. The default isolation level is dependent on
   * the driver.
   * 
   * Begin uses [context.Background] internally; to specify the context, use
   * [DB.BeginTx].
   */
  begin(): (Tx)
 }
 interface DB {
  /**
   * Driver returns the database's underlying driver.
   */
  driver(): any
 }
 interface DB {
  /**
   * Conn returns a single connection by either opening a new connection
   * or returning an existing connection from the connection pool. Conn will
   * block until either a connection is returned or ctx is canceled.
   * Queries run on the same Conn will be run in the same database session.
   * 
   * Every Conn must be returned to the database pool after use by
   * calling [Conn.Close].
   */
  conn(ctx: context.Context): (Conn)
 }
 /**
  * Tx is an in-progress database transaction.
  * 
  * A transaction must end with a call to [Tx.Commit] or [Tx.Rollback].
  * 
  * After a call to [Tx.Commit] or [Tx.Rollback], all operations on the
  * transaction fail with [ErrTxDone].
  * 
  * The statements prepared for a transaction by calling
  * the transaction's [Tx.Prepare] or [Tx.Stmt] methods are closed
  * by the call to [Tx.Commit] or [Tx.Rollback].
  */
 interface Tx {
 }
 interface Tx {
  /**
   * Commit commits the transaction.
   */
  commit(): void
 }
 interface Tx {
  /**
   * Rollback aborts the transaction.
   */
  rollback(): void
 }
 interface Tx {
  /**
   * PrepareContext creates a prepared statement for use within a transaction.
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   * 
   * To use an existing prepared statement on this transaction, see [Tx.Stmt].
   * 
   * The provided context will be used for the preparation of the context, not
   * for the execution of the returned statement. The returned statement
   * will run in the transaction context.
   */
  prepareContext(ctx: context.Context, query: string): (Stmt)
 }
 interface Tx {
  /**
   * Prepare creates a prepared statement for use within a transaction.
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   * 
   * To use an existing prepared statement on this transaction, see [Tx.Stmt].
   * 
   * Prepare uses [context.Background] internally; to specify the context, use
   * [Tx.PrepareContext].
   */
  prepare(query: string): (Stmt)
 }
 interface Tx {
  /**
   * StmtContext returns a transaction-specific prepared statement from
   * an existing statement.
   * 
   * Example:
   * 
   * ```
   * 	updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
   * 	...
   * 	tx, err := db.Begin()
   * 	...
   * 	res, err := tx.StmtContext(ctx, updateMoney).Exec(123.45, 98293203)
   * ```
   * 
   * The provided context is used for the preparation of the statement, not for the
   * execution of the statement.
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   */
  stmtContext(ctx: context.Context, stmt: Stmt): (Stmt)
 }
 interface Tx {
  /**
   * Stmt returns a transaction-specific prepared statement from
   * an existing statement.
   * 
   * Example:
   * 
   * ```
   * 	updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
   * 	...
   * 	tx, err := db.Begin()
   * 	...
   * 	res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
   * ```
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   * 
   * Stmt uses [context.Background] internally; to specify the context, use
   * [Tx.StmtContext].
   */
  stmt(stmt: Stmt): (Stmt)
 }
 interface Tx {
  /**
   * ExecContext executes a query that doesn't return rows.
   * For example: an INSERT and UPDATE.
   */
  execContext(ctx: context.Context, query: string, ...args: any[]): Result
 }
 interface Tx {
  /**
   * Exec executes a query that doesn't return rows.
   * For example: an INSERT and UPDATE.
   * 
   * Exec uses [context.Background] internally; to specify the context, use
   * [Tx.ExecContext].
   */
  exec(query: string, ...args: any[]): Result
 }
 interface Tx {
  /**
   * QueryContext executes a query that returns rows, typically a SELECT.
   */
  queryContext(ctx: context.Context, query: string, ...args: any[]): (Rows)
 }
 interface Tx {
  /**
   * Query executes a query that returns rows, typically a SELECT.
   * 
   * Query uses [context.Background] internally; to specify the context, use
   * [Tx.QueryContext].
   */
  query(query: string, ...args: any[]): (Rows)
 }
 interface Tx {
  /**
   * QueryRowContext executes a query that is expected to return at most one row.
   * QueryRowContext always returns a non-nil value. Errors are deferred until
   * [Row]'s Scan method is called.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, the [*Row.Scan] scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, query: string, ...args: any[]): (Row)
 }
 interface Tx {
  /**
   * QueryRow executes a query that is expected to return at most one row.
   * QueryRow always returns a non-nil value. Errors are deferred until
   * [Row]'s Scan method is called.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, the [*Row.Scan] scans the first selected row and discards
   * the rest.
   * 
   * QueryRow uses [context.Background] internally; to specify the context, use
   * [Tx.QueryRowContext].
   */
  queryRow(query: string, ...args: any[]): (Row)
 }
 /**
  * Stmt is a prepared statement.
  * A Stmt is safe for concurrent use by multiple goroutines.
  * 
  * If a Stmt is prepared on a [Tx] or [Conn], it will be bound to a single
  * underlying connection forever. If the [Tx] or [Conn] closes, the Stmt will
  * become unusable and all operations will return an error.
  * If a Stmt is prepared on a [DB], it will remain usable for the lifetime of the
  * [DB]. When the Stmt needs to execute on a new underlying connection, it will
  * prepare itself on the new connection automatically.
  */
 interface Stmt {
 }
 interface Stmt {
  /**
   * ExecContext executes a prepared statement with the given arguments and
   * returns a [Result] summarizing the effect of the statement.
   */
  execContext(ctx: context.Context, ...args: any[]): Result
 }
 interface Stmt {
  /**
   * Exec executes a prepared statement with the given arguments and
   * returns a [Result] summarizing the effect of the statement.
   * 
   * Exec uses [context.Background] internally; to specify the context, use
   * [Stmt.ExecContext].
   */
  exec(...args: any[]): Result
 }
 interface Stmt {
  /**
   * QueryContext executes a prepared query statement with the given arguments
   * and returns the query results as a [*Rows].
   */
  queryContext(ctx: context.Context, ...args: any[]): (Rows)
 }
 interface Stmt {
  /**
   * Query executes a prepared query statement with the given arguments
   * and returns the query results as a *Rows.
   * 
   * Query uses [context.Background] internally; to specify the context, use
   * [Stmt.QueryContext].
   */
  query(...args: any[]): (Rows)
 }
 interface Stmt {
  /**
   * QueryRowContext executes a prepared query statement with the given arguments.
   * If an error occurs during the execution of the statement, that error will
   * be returned by a call to Scan on the returned [*Row], which is always non-nil.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, the [*Row.Scan] scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, ...args: any[]): (Row)
 }
 interface Stmt {
  /**
   * QueryRow executes a prepared query statement with the given arguments.
   * If an error occurs during the execution of the statement, that error will
   * be returned by a call to Scan on the returned [*Row], which is always non-nil.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, the [*Row.Scan] scans the first selected row and discards
   * the rest.
   * 
   * Example usage:
   * 
   * ```
   * 	var name string
   * 	err := nameByUseridStmt.QueryRow(id).Scan(&name)
   * ```
   * 
   * QueryRow uses [context.Background] internally; to specify the context, use
   * [Stmt.QueryRowContext].
   */
  queryRow(...args: any[]): (Row)
 }
 interface Stmt {
  /**
   * Close closes the statement.
   */
  close(): void
 }
 /**
  * Rows is the result of a query. Its cursor starts before the first row
  * of the result set. Use [Rows.Next] to advance from row to row.
  */
 interface Rows {
 }
 interface Rows {
  /**
   * Next prepares the next result row for reading with the [Rows.Scan] method. It
   * returns true on success, or false if there is no next result row or an error
   * happened while preparing it. [Rows.Err] should be consulted to distinguish between
   * the two cases.
   * 
   * Every call to [Rows.Scan], even the first one, must be preceded by a call to [Rows.Next].
   */
  next(): boolean
 }
 interface Rows {
  /**
   * NextResultSet prepares the next result set for reading. It reports whether
   * there is further result sets, or false if there is no further result set
   * or if there is an error advancing to it. The [Rows.Err] method should be consulted
   * to distinguish between the two cases.
   * 
   * After calling NextResultSet, the [Rows.Next] method should always be called before
   * scanning. If there are further result sets they may not have rows in the result
   * set.
   */
  nextResultSet(): boolean
 }
 interface Rows {
  /**
   * Err returns the error, if any, that was encountered during iteration.
   * Err may be called after an explicit or implicit [Rows.Close].
   */
  err(): void
 }
 interface Rows {
  /**
   * Columns returns the column names.
   * Columns returns an error if the rows are closed.
   */
  columns(): Array<string>
 }
 interface Rows {
  /**
   * ColumnTypes returns column information such as column type, length,
   * and nullable. Some information may not be available from some drivers.
   */
  columnTypes(): Array<(ColumnType | undefined)>
 }
 interface Rows {
  /**
   * Scan copies the columns in the current row into the values pointed
   * at by dest. The number of values in dest must be the same as the
   * number of columns in [Rows].
   * 
   * Scan converts columns read from the database into the following
   * common Go types and special types provided by the sql package:
   * 
   * ```
   * 	*string
   * 	*[]byte
   * 	*int, *int8, *int16, *int32, *int64
   * 	*uint, *uint8, *uint16, *uint32, *uint64
   * 	*bool
   * 	*float32, *float64
   * 	*interface{}
   * 	*RawBytes
   * 	*Rows (cursor value)
   * 	any type implementing Scanner (see Scanner docs)
   * ```
   * 
   * In the most simple case, if the type of the value from the source
   * column is an integer, bool or string type T and dest is of type *T,
   * Scan simply assigns the value through the pointer.
   * 
   * Scan also converts between string and numeric types, as long as no
   * information would be lost. While Scan stringifies all numbers
   * scanned from numeric database columns into *string, scans into
   * numeric types are checked for overflow. For example, a float64 with
   * value 300 or a string with value "300" can scan into a uint16, but
   * not into a uint8, though float64(255) or "255" can scan into a
   * uint8. One exception is that scans of some float64 numbers to
   * strings may lose information when stringifying. In general, scan
   * floating point columns into *float64.
   * 
   * If a dest argument has type *[]byte, Scan saves in that argument a
   * copy of the corresponding data. The copy is owned by the caller and
   * can be modified and held indefinitely. The copy can be avoided by
   * using an argument of type [*RawBytes] instead; see the documentation
   * for [RawBytes] for restrictions on its use.
   * 
   * If an argument has type *interface{}, Scan copies the value
   * provided by the underlying driver without conversion. When scanning
   * from a source value of type []byte to *interface{}, a copy of the
   * slice is made and the caller owns the result.
   * 
   * Source values of type [time.Time] may be scanned into values of type
   * *time.Time, *interface{}, *string, or *[]byte. When converting to
   * the latter two, [time.RFC3339Nano] is used.
   * 
   * Source values of type bool may be scanned into types *bool,
   * *interface{}, *string, *[]byte, or [*RawBytes].
   * 
   * For scanning into *bool, the source may be true, false, 1, 0, or
   * string inputs parseable by [strconv.ParseBool].
   * 
   * Scan can also convert a cursor returned from a query, such as
   * "select cursor(select * from my_table) from dual", into a
   * [*Rows] value that can itself be scanned from. The parent
   * select query will close any cursor [*Rows] if the parent [*Rows] is closed.
   * 
   * If any of the first arguments implementing [Scanner] returns an error,
   * that error will be wrapped in the returned error.
   */
  scan(...dest: any[]): void
 }
 interface Rows {
  /**
   * Close closes the [Rows], preventing further enumeration. If [Rows.Next] is called
   * and returns false and there are no further result sets,
   * the [Rows] are closed automatically and it will suffice to check the
   * result of [Rows.Err]. Close is idempotent and does not affect the result of [Rows.Err].
   */
  close(): void
 }
 /**
  * A Result summarizes an executed SQL command.
  */
 interface Result {
  [key:string]: any;
  /**
   * LastInsertId returns the integer generated by the database
   * in response to a command. Typically this will be from an
   * "auto increment" column when inserting a new row. Not all
   * databases support this feature, and the syntax of such
   * statements varies.
   */
  lastInsertId(): number
  /**
   * RowsAffected returns the number of rows affected by an
   * update, insert, or delete. Not every database or database
   * driver may support this.
   */
  rowsAffected(): number
 }
}

/**
 * Package multipart implements MIME multipart parsing, as defined in RFC
 * 2046.
 * 
 * The implementation is sufficient for HTTP (RFC 2388) and the multipart
 * bodies generated by popular browsers.
 * 
 * # Limits
 * 
 * To protect against malicious inputs, this package sets limits on the size
 * of the MIME data it processes.
 * 
 * [Reader.NextPart] and [Reader.NextRawPart] limit the number of headers in a
 * part to 10000 and [Reader.ReadForm] limits the total number of headers in all
 * FileHeaders to 10000.
 * These limits may be adjusted with the GODEBUG=multipartmaxheaders=<values>
 * setting.
 * 
 * Reader.ReadForm further limits the number of parts in a form to 1000.
 * This limit may be adjusted with the GODEBUG=multipartmaxparts=<value>
 * setting.
 */
namespace multipart {
 /**
  * A FileHeader describes a file part of a multipart request.
  */
 interface FileHeader {
  filename: string
  header: textproto.MIMEHeader
  size: number
 }
 interface FileHeader {
  /**
   * Open opens and returns the [FileHeader]'s associated File.
   */
  open(): File
 }
}

/**
 * Package http provides HTTP client and server implementations.
 * 
 * [Get], [Head], [Post], and [PostForm] make HTTP (or HTTPS) requests:
 * 
 * ```
 * 	resp, err := http.Get("http://example.com/")
 * 	...
 * 	resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
 * 	...
 * 	resp, err := http.PostForm("http://example.com/form",
 * 		url.Values{"key": {"Value"}, "id": {"123"}})
 * ```
 * 
 * The caller must close the response body when finished with it:
 * 
 * ```
 * 	resp, err := http.Get("http://example.com/")
 * 	if err != nil {
 * 		// handle error
 * 	}
 * 	defer resp.Body.Close()
 * 	body, err := io.ReadAll(resp.Body)
 * 	// ...
 * ```
 * 
 * # Clients and Transports
 * 
 * For control over HTTP client headers, redirect policy, and other
 * settings, create a [Client]:
 * 
 * ```
 * 	client := &http.Client{
 * 		CheckRedirect: redirectPolicyFunc,
 * 	}
 * 
 * 	resp, err := client.Get("http://example.com")
 * 	// ...
 * 
 * 	req, err := http.NewRequest("GET", "http://example.com", nil)
 * 	// ...
 * 	req.Header.Add("If-None-Match", `W/"wyzzy"`)
 * 	resp, err := client.Do(req)
 * 	// ...
 * ```
 * 
 * For control over proxies, TLS configuration, keep-alives,
 * compression, and other settings, create a [Transport]:
 * 
 * ```
 * 	tr := &http.Transport{
 * 		MaxIdleConns:       10,
 * 		IdleConnTimeout:    30 * time.Second,
 * 		DisableCompression: true,
 * 	}
 * 	client := &http.Client{Transport: tr}
 * 	resp, err := client.Get("https://example.com")
 * ```
 * 
 * Clients and Transports are safe for concurrent use by multiple
 * goroutines and for efficiency should only be created once and re-used.
 * 
 * # Servers
 * 
 * ListenAndServe starts an HTTP server with a given address and handler.
 * The handler is usually nil, which means to use [DefaultServeMux].
 * [Handle] and [HandleFunc] add handlers to [DefaultServeMux]:
 * 
 * ```
 * 	http.Handle("/foo", fooHandler)
 * 
 * 	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
 * 		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
 * 	})
 * 
 * 	log.Fatal(http.ListenAndServe(":8080", nil))
 * ```
 * 
 * More control over the server's behavior is available by creating a
 * custom Server:
 * 
 * ```
 * 	s := &http.Server{
 * 		Addr:           ":8080",
 * 		Handler:        myHandler,
 * 		ReadTimeout:    10 * time.Second,
 * 		WriteTimeout:   10 * time.Second,
 * 		MaxHeaderBytes: 1 << 20,
 * 	}
 * 	log.Fatal(s.ListenAndServe())
 * ```
 * 
 * # HTTP/2
 * 
 * Starting with Go 1.6, the http package has transparent support for the
 * HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2
 * can do so by setting [Transport.TLSNextProto] (for clients) or
 * [Server.TLSNextProto] (for servers) to a non-nil, empty
 * map. Alternatively, the following GODEBUG settings are
 * currently supported:
 * 
 * ```
 * 	GODEBUG=http2client=0  # disable HTTP/2 client support
 * 	GODEBUG=http2server=0  # disable HTTP/2 server support
 * 	GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
 * 	GODEBUG=http2debug=2   # ... even more verbose, with frame dumps
 * ```
 * 
 * Please report any issues before disabling HTTP/2 support: https://golang.org/s/http2bug
 * 
 * The http package's [Transport] and [Server] both automatically enable
 * HTTP/2 support for simple configurations. To enable HTTP/2 for more
 * complex configurations, to use lower-level HTTP/2 features, or to use
 * a newer version of Go's http2 package, import "golang.org/x/net/http2"
 * directly and use its ConfigureTransport and/or ConfigureServer
 * functions. Manually configuring HTTP/2 via the golang.org/x/net/http2
 * package takes precedence over the net/http package's built-in HTTP/2
 * support.
 */
namespace http {
 // @ts-ignore
 import mathrand = rand
 /**
  * PushOptions describes options for [Pusher.Push].
  */
 interface PushOptions {
  /**
   * Method specifies the HTTP method for the promised request.
   * If set, it must be "GET" or "HEAD". Empty means "GET".
   */
  method: string
  /**
   * Header specifies additional promised request headers. This cannot
   * include HTTP/2 pseudo header fields like ":path" and ":scheme",
   * which will be added automatically.
   */
  header: Header
 }
 // @ts-ignore
 import urlpkg = url
 /**
  * A Request represents an HTTP request received by a server
  * or to be sent by a client.
  * 
  * The field semantics differ slightly between client and server
  * usage. In addition to the notes on the fields below, see the
  * documentation for [Request.Write] and [RoundTripper].
  */
 interface Request {
  /**
   * Method specifies the HTTP method (GET, POST, PUT, etc.).
   * For client requests, an empty string means GET.
   */
  method: string
  /**
   * URL specifies either the URI being requested (for server
   * requests) or the URL to access (for client requests).
   * 
   * For server requests, the URL is parsed from the URI
   * supplied on the Request-Line as stored in RequestURI.  For
   * most requests, fields other than Path and RawQuery will be
   * empty. (See RFC 7230, Section 5.3)
   * 
   * For client requests, the URL's Host specifies the server to
   * connect to, while the Request's Host field optionally
   * specifies the Host header value to send in the HTTP
   * request.
   */
  url?: url.URL
  /**
   * The protocol version for incoming server requests.
   * 
   * For client requests, these fields are ignored. The HTTP
   * client code always uses either HTTP/1.1 or HTTP/2.
   * See the docs on Transport for details.
   */
  proto: string // "HTTP/1.0"
  protoMajor: number // 1
  protoMinor: number // 0
  /**
   * Header contains the request header fields either received
   * by the server or to be sent by the client.
   * 
   * If a server received a request with header lines,
   * 
   * ```
   * 	Host: example.com
   * 	accept-encoding: gzip, deflate
   * 	Accept-Language: en-us
   * 	fOO: Bar
   * 	foo: two
   * ```
   * 
   * then
   * 
   * ```
   * 	Header = map[string][]string{
   * 		"Accept-Encoding": {"gzip, deflate"},
   * 		"Accept-Language": {"en-us"},
   * 		"Foo": {"Bar", "two"},
   * 	}
   * ```
   * 
   * For incoming requests, the Host header is promoted to the
   * Request.Host field and removed from the Header map.
   * 
   * HTTP defines that header names are case-insensitive. The
   * request parser implements this by using CanonicalHeaderKey,
   * making the first character and any characters following a
   * hyphen uppercase and the rest lowercase.
   * 
   * For client requests, certain headers such as Content-Length
   * and Connection are automatically written when needed and
   * values in Header may be ignored. See the documentation
   * for the Request.Write method.
   */
  header: Header
  /**
   * Body is the request's body.
   * 
   * For client requests, a nil body means the request has no
   * body, such as a GET request. The HTTP Client's Transport
   * is responsible for calling the Close method.
   * 
   * For server requests, the Request Body is always non-nil
   * but will return EOF immediately when no body is present.
   * The Server will close the request body. The ServeHTTP
   * Handler does not need to.
   * 
   * Body must allow Read to be called concurrently with Close.
   * In particular, calling Close should unblock a Read waiting
   * for input.
   */
  body: io.ReadCloser
  /**
   * GetBody defines an optional func to return a new copy of
   * Body. It is used for client requests when a redirect requires
   * reading the body more than once. Use of GetBody still
   * requires setting Body.
   * 
   * For server requests, it is unused.
   */
  getBody: () => io.ReadCloser
  /**
   * ContentLength records the length of the associated content.
   * The value -1 indicates that the length is unknown.
   * Values >= 0 indicate that the given number of bytes may
   * be read from Body.
   * 
   * For client requests, a value of 0 with a non-nil Body is
   * also treated as unknown.
   */
  contentLength: number
  /**
   * TransferEncoding lists the transfer encodings from outermost to
   * innermost. An empty list denotes the "identity" encoding.
   * TransferEncoding can usually be ignored; chunked encoding is
   * automatically added and removed as necessary when sending and
   * receiving requests.
   */
  transferEncoding: Array<string>
  /**
   * Close indicates whether to close the connection after
   * replying to this request (for servers) or after sending this
   * request and reading its response (for clients).
   * 
   * For server requests, the HTTP server handles this automatically
   * and this field is not needed by Handlers.
   * 
   * For client requests, setting this field prevents re-use of
   * TCP connections between requests to the same hosts, as if
   * Transport.DisableKeepAlives were set.
   */
  close: boolean
  /**
   * For server requests, Host specifies the host on which the
   * URL is sought. For HTTP/1 (per RFC 7230, section 5.4), this
   * is either the value of the "Host" header or the host name
   * given in the URL itself. For HTTP/2, it is the value of the
   * ":authority" pseudo-header field.
   * It may be of the form "host:port". For international domain
   * names, Host may be in Punycode or Unicode form. Use
   * golang.org/x/net/idna to convert it to either format if
   * needed.
   * To prevent DNS rebinding attacks, server Handlers should
   * validate that the Host header has a value for which the
   * Handler considers itself authoritative. The included
   * ServeMux supports patterns registered to particular host
   * names and thus protects its registered Handlers.
   * 
   * For client requests, Host optionally overrides the Host
   * header to send. If empty, the Request.Write method uses
   * the value of URL.Host. Host may contain an international
   * domain name.
   */
  host: string
  /**
   * Form contains the parsed form data, including both the URL
   * field's query parameters and the PATCH, POST, or PUT form data.
   * This field is only available after ParseForm is called.
   * The HTTP client ignores Form and uses Body instead.
   */
  form: url.Values
  /**
   * PostForm contains the parsed form data from PATCH, POST
   * or PUT body parameters.
   * 
   * This field is only available after ParseForm is called.
   * The HTTP client ignores PostForm and uses Body instead.
   */
  postForm: url.Values
  /**
   * MultipartForm is the parsed multipart form, including file uploads.
   * This field is only available after ParseMultipartForm is called.
   * The HTTP client ignores MultipartForm and uses Body instead.
   */
  multipartForm?: multipart.Form
  /**
   * Trailer specifies additional headers that are sent after the request
   * body.
   * 
   * For server requests, the Trailer map initially contains only the
   * trailer keys, with nil values. (The client declares which trailers it
   * will later send.)  While the handler is reading from Body, it must
   * not reference Trailer. After reading from Body returns EOF, Trailer
   * can be read again and will contain non-nil values, if they were sent
   * by the client.
   * 
   * For client requests, Trailer must be initialized to a map containing
   * the trailer keys to later send. The values may be nil or their final
   * values. The ContentLength must be 0 or -1, to send a chunked request.
   * After the HTTP request is sent the map values can be updated while
   * the request body is read. Once the body returns EOF, the caller must
   * not mutate Trailer.
   * 
   * Few HTTP clients, servers, or proxies support HTTP trailers.
   */
  trailer: Header
  /**
   * RemoteAddr allows HTTP servers and other software to record
   * the network address that sent the request, usually for
   * logging. This field is not filled in by ReadRequest and
   * has no defined format. The HTTP server in this package
   * sets RemoteAddr to an "IP:port" address before invoking a
   * handler.
   * This field is ignored by the HTTP client.
   */
  remoteAddr: string
  /**
   * RequestURI is the unmodified request-target of the
   * Request-Line (RFC 7230, Section 3.1.1) as sent by the client
   * to a server. Usually the URL field should be used instead.
   * It is an error to set this field in an HTTP client request.
   */
  requestURI: string
  /**
   * TLS allows HTTP servers and other software to record
   * information about the TLS connection on which the request
   * was received. This field is not filled in by ReadRequest.
   * The HTTP server in this package sets the field for
   * TLS-enabled connections before invoking a handler;
   * otherwise it leaves the field nil.
   * This field is ignored by the HTTP client.
   */
  tls?: any
  /**
   * Cancel is an optional channel whose closure indicates that the client
   * request should be regarded as canceled. Not all implementations of
   * RoundTripper may support Cancel.
   * 
   * For server requests, this field is not applicable.
   * 
   * Deprecated: Set the Request's context with NewRequestWithContext
   * instead. If a Request's Cancel field and context are both
   * set, it is undefined whether Cancel is respected.
   */
  cancel: undefined
  /**
   * Response is the redirect response which caused this request
   * to be created. This field is only populated during client
   * redirects.
   */
  response?: Response
  /**
   * Pattern is the [ServeMux] pattern that matched the request.
   * It is empty if the request was not matched against a pattern.
   */
  pattern: string
 }
 interface Request {
  /**
   * Context returns the request's context. To change the context, use
   * [Request.Clone] or [Request.WithContext].
   * 
   * The returned context is always non-nil; it defaults to the
   * background context.
   * 
   * For outgoing client requests, the context controls cancellation.
   * 
   * For incoming server requests, the context is canceled when the
   * client's connection closes, the request is canceled (with HTTP/2),
   * or when the ServeHTTP method returns.
   */
  context(): context.Context
 }
 interface Request {
  /**
   * WithContext returns a shallow copy of r with its context changed
   * to ctx. The provided ctx must be non-nil.
   * 
   * For outgoing client request, the context controls the entire
   * lifetime of a request and its response: obtaining a connection,
   * sending the request, and reading the response headers and body.
   * 
   * To create a new request with a context, use [NewRequestWithContext].
   * To make a deep copy of a request with a new context, use [Request.Clone].
   */
  withContext(ctx: context.Context): (Request)
 }
 interface Request {
  /**
   * Clone returns a deep copy of r with its context changed to ctx.
   * The provided ctx must be non-nil.
   * 
   * Clone only makes a shallow copy of the Body field.
   * 
   * For an outgoing client request, the context controls the entire
   * lifetime of a request and its response: obtaining a connection,
   * sending the request, and reading the response headers and body.
   */
  clone(ctx: context.Context): (Request)
 }
 interface Request {
  /**
   * ProtoAtLeast reports whether the HTTP protocol used
   * in the request is at least major.minor.
   */
  protoAtLeast(major: number, minor: number): boolean
 }
 interface Request {
  /**
   * UserAgent returns the client's User-Agent, if sent in the request.
   */
  userAgent(): string
 }
 interface Request {
  /**
   * Cookies parses and returns the HTTP cookies sent with the request.
   */
  cookies(): Array<(Cookie | undefined)>
 }
 interface Request {
  /**
   * CookiesNamed parses and returns the named HTTP cookies sent with the request
   * or an empty slice if none matched.
   */
  cookiesNamed(name: string): Array<(Cookie | undefined)>
 }
 interface Request {
  /**
   * Cookie returns the named cookie provided in the request or
   * [ErrNoCookie] if not found.
   * If multiple cookies match the given name, only one cookie will
   * be returned.
   */
  cookie(name: string): (Cookie)
 }
 interface Request {
  /**
   * AddCookie adds a cookie to the request. Per RFC 6265 section 5.4,
   * AddCookie does not attach more than one [Cookie] header field. That
   * means all cookies, if any, are written into the same line,
   * separated by semicolon.
   * AddCookie only sanitizes c's name and value, and does not sanitize
   * a Cookie header already present in the request.
   */
  addCookie(c: Cookie): void
 }
 interface Request {
  /**
   * Referer returns the referring URL, if sent in the request.
   * 
   * Referer is misspelled as in the request itself, a mistake from the
   * earliest days of HTTP.  This value can also be fetched from the
   * [Header] map as Header["Referer"]; the benefit of making it available
   * as a method is that the compiler can diagnose programs that use the
   * alternate (correct English) spelling req.Referrer() but cannot
   * diagnose programs that use Header["Referrer"].
   */
  referer(): string
 }
 interface Request {
  /**
   * MultipartReader returns a MIME multipart reader if this is a
   * multipart/form-data or a multipart/mixed POST request, else returns nil and an error.
   * Use this function instead of [Request.ParseMultipartForm] to
   * process the request body as a stream.
   */
  multipartReader(): (multipart.Reader)
 }
 interface Request {
  /**
   * Write writes an HTTP/1.1 request, which is the header and body, in wire format.
   * This method consults the following fields of the request:
   * 
   * ```
   * 	Host
   * 	URL
   * 	Method (defaults to "GET")
   * 	Header
   * 	ContentLength
   * 	TransferEncoding
   * 	Body
   * ```
   * 
   * If Body is present, Content-Length is <= 0 and [Request.TransferEncoding]
   * hasn't been set to "identity", Write adds "Transfer-Encoding:
   * chunked" to the header. Body is closed after it is sent.
   */
  write(w: io.Writer): void
 }
 interface Request {
  /**
   * WriteProxy is like [Request.Write] but writes the request in the form
   * expected by an HTTP proxy. In particular, [Request.WriteProxy] writes the
   * initial Request-URI line of the request with an absolute URI, per
   * section 5.3 of RFC 7230, including the scheme and host.
   * In either case, WriteProxy also writes a Host header, using
   * either r.Host or r.URL.Host.
   */
  writeProxy(w: io.Writer): void
 }
 interface Request {
  /**
   * BasicAuth returns the username and password provided in the request's
   * Authorization header, if the request uses HTTP Basic Authentication.
   * See RFC 2617, Section 2.
   */
  basicAuth(): [string, string, boolean]
 }
 interface Request {
  /**
   * SetBasicAuth sets the request's Authorization header to use HTTP
   * Basic Authentication with the provided username and password.
   * 
   * With HTTP Basic Authentication the provided username and password
   * are not encrypted. It should generally only be used in an HTTPS
   * request.
   * 
   * The username may not contain a colon. Some protocols may impose
   * additional requirements on pre-escaping the username and
   * password. For instance, when used with OAuth2, both arguments must
   * be URL encoded first with [url.QueryEscape].
   */
  setBasicAuth(username: string, password: string): void
 }
 interface Request {
  /**
   * ParseForm populates r.Form and r.PostForm.
   * 
   * For all requests, ParseForm parses the raw query from the URL and updates
   * r.Form.
   * 
   * For POST, PUT, and PATCH requests, it also reads the request body, parses it
   * as a form and puts the results into both r.PostForm and r.Form. Request body
   * parameters take precedence over URL query string values in r.Form.
   * 
   * If the request Body's size has not already been limited by [MaxBytesReader],
   * the size is capped at 10MB.
   * 
   * For other HTTP methods, or when the Content-Type is not
   * application/x-www-form-urlencoded, the request Body is not read, and
   * r.PostForm is initialized to a non-nil, empty value.
   * 
   * [Request.ParseMultipartForm] calls ParseForm automatically.
   * ParseForm is idempotent.
   */
  parseForm(): void
 }
 interface Request {
  /**
   * ParseMultipartForm parses a request body as multipart/form-data.
   * The whole request body is parsed and up to a total of maxMemory bytes of
   * its file parts are stored in memory, with the remainder stored on
   * disk in temporary files.
   * ParseMultipartForm calls [Request.ParseForm] if necessary.
   * If ParseForm returns an error, ParseMultipartForm returns it but also
   * continues parsing the request body.
   * After one call to ParseMultipartForm, subsequent calls have no effect.
   */
  parseMultipartForm(maxMemory: number): void
 }
 interface Request {
  /**
   * FormValue returns the first value for the named component of the query.
   * The precedence order:
   *  1. application/x-www-form-urlencoded form body (POST, PUT, PATCH only)
   *  2. query parameters (always)
   *  3. multipart/form-data form body (always)
   * 
   * FormValue calls [Request.ParseMultipartForm] and [Request.ParseForm]
   * if necessary and ignores any errors returned by these functions.
   * If key is not present, FormValue returns the empty string.
   * To access multiple values of the same key, call ParseForm and
   * then inspect [Request.Form] directly.
   */
  formValue(key: string): string
 }
 interface Request {
  /**
   * PostFormValue returns the first value for the named component of the POST,
   * PUT, or PATCH request body. URL query parameters are ignored.
   * PostFormValue calls [Request.ParseMultipartForm] and [Request.ParseForm] if necessary and ignores
   * any errors returned by these functions.
   * If key is not present, PostFormValue returns the empty string.
   */
  postFormValue(key: string): string
 }
 interface Request {
  /**
   * FormFile returns the first file for the provided form key.
   * FormFile calls [Request.ParseMultipartForm] and [Request.ParseForm] if necessary.
   */
  formFile(key: string): [multipart.File, (multipart.FileHeader)]
 }
 interface Request {
  /**
   * PathValue returns the value for the named path wildcard in the [ServeMux] pattern
   * that matched the request.
   * It returns the empty string if the request was not matched against a pattern
   * or there is no such wildcard in the pattern.
   */
  pathValue(name: string): string
 }
 interface Request {
  /**
   * SetPathValue sets name to value, so that subsequent calls to r.PathValue(name)
   * return value.
   */
  setPathValue(name: string, value: string): void
 }
 /**
  * A Handler responds to an HTTP request.
  * 
  * [Handler.ServeHTTP] should write reply headers and data to the [ResponseWriter]
  * and then return. Returning signals that the request is finished; it
  * is not valid to use the [ResponseWriter] or read from the
  * [Request.Body] after or concurrently with the completion of the
  * ServeHTTP call.
  * 
  * Depending on the HTTP client software, HTTP protocol version, and
  * any intermediaries between the client and the Go server, it may not
  * be possible to read from the [Request.Body] after writing to the
  * [ResponseWriter]. Cautious handlers should read the [Request.Body]
  * first, and then reply.
  * 
  * Except for reading the body, handlers should not modify the
  * provided Request.
  * 
  * If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
  * that the effect of the panic was isolated to the active request.
  * It recovers the panic, logs a stack trace to the server error log,
  * and either closes the network connection or sends an HTTP/2
  * RST_STREAM, depending on the HTTP protocol. To abort a handler so
  * the client sees an interrupted response but the server doesn't log
  * an error, panic with the value [ErrAbortHandler].
  */
 interface Handler {
  [key:string]: any;
  serveHTTP(_arg0: ResponseWriter, _arg1: Request): void
 }
 /**
  * A ResponseWriter interface is used by an HTTP handler to
  * construct an HTTP response.
  * 
  * A ResponseWriter may not be used after [Handler.ServeHTTP] has returned.
  */
 interface ResponseWriter {
  [key:string]: any;
  /**
   * Header returns the header map that will be sent by
   * [ResponseWriter.WriteHeader]. The [Header] map also is the mechanism with which
   * [Handler] implementations can set HTTP trailers.
   * 
   * Changing the header map after a call to [ResponseWriter.WriteHeader] (or
   * [ResponseWriter.Write]) has no effect unless the HTTP status code was of the
   * 1xx class or the modified headers are trailers.
   * 
   * There are two ways to set Trailers. The preferred way is to
   * predeclare in the headers which trailers you will later
   * send by setting the "Trailer" header to the names of the
   * trailer keys which will come later. In this case, those
   * keys of the Header map are treated as if they were
   * trailers. See the example. The second way, for trailer
   * keys not known to the [Handler] until after the first [ResponseWriter.Write],
   * is to prefix the [Header] map keys with the [TrailerPrefix]
   * constant value.
   * 
   * To suppress automatic response headers (such as "Date"), set
   * their value to nil.
   */
  header(): Header
  /**
   * Write writes the data to the connection as part of an HTTP reply.
   * 
   * If [ResponseWriter.WriteHeader] has not yet been called, Write calls
   * WriteHeader(http.StatusOK) before writing the data. If the Header
   * does not contain a Content-Type line, Write adds a Content-Type set
   * to the result of passing the initial 512 bytes of written data to
   * [DetectContentType]. Additionally, if the total size of all written
   * data is under a few KB and there are no Flush calls, the
   * Content-Length header is added automatically.
   * 
   * Depending on the HTTP protocol version and the client, calling
   * Write or WriteHeader may prevent future reads on the
   * Request.Body. For HTTP/1.x requests, handlers should read any
   * needed request body data before writing the response. Once the
   * headers have been flushed (due to either an explicit Flusher.Flush
   * call or writing enough data to trigger a flush), the request body
   * may be unavailable. For HTTP/2 requests, the Go HTTP server permits
   * handlers to continue to read the request body while concurrently
   * writing the response. However, such behavior may not be supported
   * by all HTTP/2 clients. Handlers should read before writing if
   * possible to maximize compatibility.
   */
  write(_arg0: string|Array<number>): number
  /**
   * WriteHeader sends an HTTP response header with the provided
   * status code.
   * 
   * If WriteHeader is not called explicitly, the first call to Write
   * will trigger an implicit WriteHeader(http.StatusOK).
   * Thus explicit calls to WriteHeader are mainly used to
   * send error codes or 1xx informational responses.
   * 
   * The provided code must be a valid HTTP 1xx-5xx status code.
   * Any number of 1xx headers may be written, followed by at most
   * one 2xx-5xx header. 1xx headers are sent immediately, but 2xx-5xx
   * headers may be buffered. Use the Flusher interface to send
   * buffered data. The header map is cleared when 2xx-5xx headers are
   * sent, but not with 1xx headers.
   * 
   * The server will automatically send a 100 (Continue) header
   * on the first read from the request body if the request has
   * an "Expect: 100-continue" header.
   */
  writeHeader(statusCode: number): void
 }
 /**
  * A Server defines parameters for running an HTTP server.
  * The zero value for Server is a valid configuration.
  */
 interface Server {
  /**
   * Addr optionally specifies the TCP address for the server to listen on,
   * in the form "host:port". If empty, ":http" (port 80) is used.
   * The service names are defined in RFC 6335 and assigned by IANA.
   * See net.Dial for details of the address format.
   */
  addr: string
  handler: Handler // handler to invoke, http.DefaultServeMux if nil
  /**
   * DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler,
   * otherwise responds with 200 OK and Content-Length: 0.
   */
  disableGeneralOptionsHandler: boolean
  /**
   * TLSConfig optionally provides a TLS configuration for use
   * by ServeTLS and ListenAndServeTLS. Note that this value is
   * cloned by ServeTLS and ListenAndServeTLS, so it's not
   * possible to modify the configuration with methods like
   * tls.Config.SetSessionTicketKeys. To use
   * SetSessionTicketKeys, use Server.Serve with a TLS Listener
   * instead.
   */
  tlsConfig?: any
  /**
   * ReadTimeout is the maximum duration for reading the entire
   * request, including the body. A zero or negative value means
   * there will be no timeout.
   * 
   * Because ReadTimeout does not let Handlers make per-request
   * decisions on each request body's acceptable deadline or
   * upload rate, most users will prefer to use
   * ReadHeaderTimeout. It is valid to use them both.
   */
  readTimeout: time.Duration
  /**
   * ReadHeaderTimeout is the amount of time allowed to read
   * request headers. The connection's read deadline is reset
   * after reading the headers and the Handler can decide what
   * is considered too slow for the body. If zero, the value of
   * ReadTimeout is used. If negative, or if zero and ReadTimeout
   * is zero or negative, there is no timeout.
   */
  readHeaderTimeout: time.Duration
  /**
   * WriteTimeout is the maximum duration before timing out
   * writes of the response. It is reset whenever a new
   * request's header is read. Like ReadTimeout, it does not
   * let Handlers make decisions on a per-request basis.
   * A zero or negative value means there will be no timeout.
   */
  writeTimeout: time.Duration
  /**
   * IdleTimeout is the maximum amount of time to wait for the
   * next request when keep-alives are enabled. If zero, the value
   * of ReadTimeout is used. If negative, or if zero and ReadTimeout
   * is zero or negative, there is no timeout.
   */
  idleTimeout: time.Duration
  /**
   * MaxHeaderBytes controls the maximum number of bytes the
   * server will read parsing the request header's keys and
   * values, including the request line. It does not limit the
   * size of the request body.
   * If zero, DefaultMaxHeaderBytes is used.
   */
  maxHeaderBytes: number
  /**
   * TLSNextProto optionally specifies a function to take over
   * ownership of the provided TLS connection when an ALPN
   * protocol upgrade has occurred. The map key is the protocol
   * name negotiated. The Handler argument should be used to
   * handle HTTP requests and will initialize the Request's TLS
   * and RemoteAddr if not already set. The connection is
   * automatically closed when the function returns.
   * If TLSNextProto is not nil, HTTP/2 support is not enabled
   * automatically.
   */
  tlsNextProto: _TygojaDict
  /**
   * ConnState specifies an optional callback function that is
   * called when a client connection changes state. See the
   * ConnState type and associated constants for details.
   */
  connState: (_arg0: net.Conn, _arg1: ConnState) => void
  /**
   * ErrorLog specifies an optional logger for errors accepting
   * connections, unexpected behavior from handlers, and
   * underlying FileSystem errors.
   * If nil, logging is done via the log package's standard logger.
   */
  errorLog?: any
  /**
   * BaseContext optionally specifies a function that returns
   * the base context for incoming requests on this server.
   * The provided Listener is the specific Listener that's
   * about to start accepting requests.
   * If BaseContext is nil, the default is context.Background().
   * If non-nil, it must return a non-nil context.
   */
  baseContext: (_arg0: net.Listener) => context.Context
  /**
   * ConnContext optionally specifies a function that modifies
   * the context used for a new connection c. The provided ctx
   * is derived from the base context and has a ServerContextKey
   * value.
   */
  connContext: (ctx: context.Context, c: net.Conn) => context.Context
  /**
   * HTTP2 configures HTTP/2 connections.
   * 
   * This field does not yet have any effect.
   * See https://go.dev/issue/67813.
   */
  http2?: HTTP2Config
  /**
   * Protocols is the set of protocols accepted by the server.
   * 
   * If Protocols includes UnencryptedHTTP2, the server will accept
   * unencrypted HTTP/2 connections. The server can serve both
   * HTTP/1 and unencrypted HTTP/2 on the same address and port.
   * 
   * If Protocols is nil, the default is usually HTTP/1 and HTTP/2.
   * If TLSNextProto is non-nil and does not contain an "h2" entry,
   * the default is HTTP/1 only.
   */
  protocols?: Protocols
 }
 interface Server {
  /**
   * Close immediately closes all active net.Listeners and any
   * connections in state [StateNew], [StateActive], or [StateIdle]. For a
   * graceful shutdown, use [Server.Shutdown].
   * 
   * Close does not attempt to close (and does not even know about)
   * any hijacked connections, such as WebSockets.
   * 
   * Close returns any error returned from closing the [Server]'s
   * underlying Listener(s).
   */
  close(): void
 }
 interface Server {
  /**
   * Shutdown gracefully shuts down the server without interrupting any
   * active connections. Shutdown works by first closing all open
   * listeners, then closing all idle connections, and then waiting
   * indefinitely for connections to return to idle and then shut down.
   * If the provided context expires before the shutdown is complete,
   * Shutdown returns the context's error, otherwise it returns any
   * error returned from closing the [Server]'s underlying Listener(s).
   * 
   * When Shutdown is called, [Serve], [ListenAndServe], and
   * [ListenAndServeTLS] immediately return [ErrServerClosed]. Make sure the
   * program doesn't exit and waits instead for Shutdown to return.
   * 
   * Shutdown does not attempt to close nor wait for hijacked
   * connections such as WebSockets. The caller of Shutdown should
   * separately notify such long-lived connections of shutdown and wait
   * for them to close, if desired. See [Server.RegisterOnShutdown] for a way to
   * register shutdown notification functions.
   * 
   * Once Shutdown has been called on a server, it may not be reused;
   * future calls to methods such as Serve will return ErrServerClosed.
   */
  shutdown(ctx: context.Context): void
 }
 interface Server {
  /**
   * RegisterOnShutdown registers a function to call on [Server.Shutdown].
   * This can be used to gracefully shutdown connections that have
   * undergone ALPN protocol upgrade or that have been hijacked.
   * This function should start protocol-specific graceful shutdown,
   * but should not wait for shutdown to complete.
   */
  registerOnShutdown(f: () => void): void
 }
 interface Server {
  /**
   * ListenAndServe listens on the TCP network address s.Addr and then
   * calls [Serve] to handle requests on incoming connections.
   * Accepted connections are configured to enable TCP keep-alives.
   * 
   * If s.Addr is blank, ":http" is used.
   * 
   * ListenAndServe always returns a non-nil error. After [Server.Shutdown] or [Server.Close],
   * the returned error is [ErrServerClosed].
   */
  listenAndServe(): void
 }
 interface Server {
  /**
   * Serve accepts incoming connections on the Listener l, creating a
   * new service goroutine for each. The service goroutines read requests and
   * then call s.Handler to reply to them.
   * 
   * HTTP/2 support is only enabled if the Listener returns [*tls.Conn]
   * connections and they were configured with "h2" in the TLS
   * Config.NextProtos.
   * 
   * Serve always returns a non-nil error and closes l.
   * After [Server.Shutdown] or [Server.Close], the returned error is [ErrServerClosed].
   */
  serve(l: net.Listener): void
 }
 interface Server {
  /**
   * ServeTLS accepts incoming connections on the Listener l, creating a
   * new service goroutine for each. The service goroutines perform TLS
   * setup and then read requests, calling s.Handler to reply to them.
   * 
   * Files containing a certificate and matching private key for the
   * server must be provided if neither the [Server]'s
   * TLSConfig.Certificates, TLSConfig.GetCertificate nor
   * config.GetConfigForClient are populated.
   * If the certificate is signed by a certificate authority, the
   * certFile should be the concatenation of the server's certificate,
   * any intermediates, and the CA's certificate.
   * 
   * ServeTLS always returns a non-nil error. After [Server.Shutdown] or [Server.Close], the
   * returned error is [ErrServerClosed].
   */
  serveTLS(l: net.Listener, certFile: string, keyFile: string): void
 }
 interface Server {
  /**
   * SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled.
   * By default, keep-alives are always enabled. Only very
   * resource-constrained environments or servers in the process of
   * shutting down should disable them.
   */
  setKeepAlivesEnabled(v: boolean): void
 }
 interface Server {
  /**
   * ListenAndServeTLS listens on the TCP network address s.Addr and
   * then calls [ServeTLS] to handle requests on incoming TLS connections.
   * Accepted connections are configured to enable TCP keep-alives.
   * 
   * Filenames containing a certificate and matching private key for the
   * server must be provided if neither the [Server]'s TLSConfig.Certificates
   * nor TLSConfig.GetCertificate are populated. If the certificate is
   * signed by a certificate authority, the certFile should be the
   * concatenation of the server's certificate, any intermediates, and
   * the CA's certificate.
   * 
   * If s.Addr is blank, ":https" is used.
   * 
   * ListenAndServeTLS always returns a non-nil error. After [Server.Shutdown] or
   * [Server.Close], the returned error is [ErrServerClosed].
   */
  listenAndServeTLS(certFile: string, keyFile: string): void
 }
}

namespace exec {
 /**
  * Cmd represents an external command being prepared or run.
  * 
  * A Cmd cannot be reused after calling its [Cmd.Run], [Cmd.Output] or [Cmd.CombinedOutput]
  * methods.
  */
 interface Cmd {
  /**
   * Path is the path of the command to run.
   * 
   * This is the only field that must be set to a non-zero
   * value. If Path is relative, it is evaluated relative
   * to Dir.
   */
  path: string
  /**
   * Args holds command line arguments, including the command as Args[0].
   * If the Args field is empty or nil, Run uses {Path}.
   * 
   * In typical use, both Path and Args are set by calling Command.
   */
  args: Array<string>
  /**
   * Env specifies the environment of the process.
   * Each entry is of the form "key=value".
   * If Env is nil, the new process uses the current process's
   * environment.
   * If Env contains duplicate environment keys, only the last
   * value in the slice for each duplicate key is used.
   * As a special case on Windows, SYSTEMROOT is always added if
   * missing and not explicitly set to the empty string.
   * 
   * See also the Dir field, which may set PWD in the environment.
   */
  env: Array<string>
  /**
   * Dir specifies the working directory of the command.
   * If Dir is the empty string, Run runs the command in the
   * calling process's current directory.
   * 
   * On Unix systems, the value of Dir also determines the
   * child process's PWD environment variable if not otherwise
   * specified. A Unix process represents its working directory
   * not by name but as an implicit reference to a node in the
   * file tree. So, if the child process obtains its working
   * directory by calling a function such as C's getcwd, which
   * computes the canonical name by walking up the file tree, it
   * will not recover the original value of Dir if that value
   * was an alias involving symbolic links. However, if the
   * child process calls Go's [os.Getwd] or GNU C's
   * get_current_dir_name, and the value of PWD is an alias for
   * the current directory, those functions will return the
   * value of PWD, which matches the value of Dir.
   */
  dir: string
  /**
   * Stdin specifies the process's standard input.
   * 
   * If Stdin is nil, the process reads from the null device (os.DevNull).
   * 
   * If Stdin is an *os.File, the process's standard input is connected
   * directly to that file.
   * 
   * Otherwise, during the execution of the command a separate
   * goroutine reads from Stdin and delivers that data to the command
   * over a pipe. In this case, Wait does not complete until the goroutine
   * stops copying, either because it has reached the end of Stdin
   * (EOF or a read error), or because writing to the pipe returned an error,
   * or because a nonzero WaitDelay was set and expired.
   */
  stdin: io.Reader
  /**
   * Stdout and Stderr specify the process's standard output and error.
   * 
   * If either is nil, Run connects the corresponding file descriptor
   * to the null device (os.DevNull).
   * 
   * If either is an *os.File, the corresponding output from the process
   * is connected directly to that file.
   * 
   * Otherwise, during the execution of the command a separate goroutine
   * reads from the process over a pipe and delivers that data to the
   * corresponding Writer. In this case, Wait does not complete until the
   * goroutine reaches EOF or encounters an error or a nonzero WaitDelay
   * expires.
   * 
   * If Stdout and Stderr are the same writer, and have a type that can
   * be compared with ==, at most one goroutine at a time will call Write.
   */
  stdout: io.Writer
  stderr: io.Writer
  /**
   * ExtraFiles specifies additional open files to be inherited by the
   * new process. It does not include standard input, standard output, or
   * standard error. If non-nil, entry i becomes file descriptor 3+i.
   * 
   * ExtraFiles is not supported on Windows.
   */
  extraFiles: Array<(os.File | undefined)>
  /**
   * SysProcAttr holds optional, operating system-specific attributes.
   * Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
   */
  sysProcAttr?: syscall.SysProcAttr
  /**
   * Process is the underlying process, once started.
   */
  process?: os.Process
  /**
   * ProcessState contains information about an exited process.
   * If the process was started successfully, Wait or Run will
   * populate its ProcessState when the command completes.
   */
  processState?: os.ProcessState
  err: Error // LookPath error, if any.
  /**
   * If Cancel is non-nil, the command must have been created with
   * CommandContext and Cancel will be called when the command's
   * Context is done. By default, CommandContext sets Cancel to
   * call the Kill method on the command's Process.
   * 
   * Typically a custom Cancel will send a signal to the command's
   * Process, but it may instead take other actions to initiate cancellation,
   * such as closing a stdin or stdout pipe or sending a shutdown request on a
   * network socket.
   * 
   * If the command exits with a success status after Cancel is
   * called, and Cancel does not return an error equivalent to
   * os.ErrProcessDone, then Wait and similar methods will return a non-nil
   * error: either an error wrapping the one returned by Cancel,
   * or the error from the Context.
   * (If the command exits with a non-success status, or Cancel
   * returns an error that wraps os.ErrProcessDone, Wait and similar methods
   * continue to return the command's usual exit status.)
   * 
   * If Cancel is set to nil, nothing will happen immediately when the command's
   * Context is done, but a nonzero WaitDelay will still take effect. That may
   * be useful, for example, to work around deadlocks in commands that do not
   * support shutdown signals but are expected to always finish quickly.
   * 
   * Cancel will not be called if Start returns a non-nil error.
   */
  cancel: () => void
  /**
   * If WaitDelay is non-zero, it bounds the time spent waiting on two sources
   * of unexpected delay in Wait: a child process that fails to exit after the
   * associated Context is canceled, and a child process that exits but leaves
   * its I/O pipes unclosed.
   * 
   * The WaitDelay timer starts when either the associated Context is done or a
   * call to Wait observes that the child process has exited, whichever occurs
   * first. When the delay has elapsed, the command shuts down the child process
   * and/or its I/O pipes.
   * 
   * If the child process has failed to exit — perhaps because it ignored or
   * failed to receive a shutdown signal from a Cancel function, or because no
   * Cancel function was set — then it will be terminated using os.Process.Kill.
   * 
   * Then, if the I/O pipes communicating with the child process are still open,
   * those pipes are closed in order to unblock any goroutines currently blocked
   * on Read or Write calls.
   * 
   * If pipes are closed due to WaitDelay, no Cancel call has occurred,
   * and the command has otherwise exited with a successful status, Wait and
   * similar methods will return ErrWaitDelay instead of nil.
   * 
   * If WaitDelay is zero (the default), I/O pipes will be read until EOF,
   * which might not occur until orphaned subprocesses of the command have
   * also closed their descriptors for the pipes.
   */
  waitDelay: time.Duration
 }
 interface Cmd {
  /**
   * String returns a human-readable description of c.
   * It is intended only for debugging.
   * In particular, it is not suitable for use as input to a shell.
   * The output of String may vary across Go releases.
   */
  string(): string
 }
 interface Cmd {
  /**
   * Run starts the specified command and waits for it to complete.
   * 
   * The returned error is nil if the command runs, has no problems
   * copying stdin, stdout, and stderr, and exits with a zero exit
   * status.
   * 
   * If the command starts but does not complete successfully, the error is of
   * type [*ExitError]. Other error types may be returned for other situations.
   * 
   * If the calling goroutine has locked the operating system thread
   * with [runtime.LockOSThread] and modified any inheritable OS-level
   * thread state (for example, Linux or Plan 9 name spaces), the new
   * process will inherit the caller's thread state.
   */
  run(): void
 }
 interface Cmd {
  /**
   * Start starts the specified command but does not wait for it to complete.
   * 
   * If Start returns successfully, the c.Process field will be set.
   * 
   * After a successful call to Start the [Cmd.Wait] method must be called in
   * order to release associated system resources.
   */
  start(): void
 }
 interface Cmd {
  /**
   * Wait waits for the command to exit and waits for any copying to
   * stdin or copying from stdout or stderr to complete.
   * 
   * The command must have been started by [Cmd.Start].
   * 
   * The returned error is nil if the command runs, has no problems
   * copying stdin, stdout, and stderr, and exits with a zero exit
   * status.
   * 
   * If the command fails to run or doesn't complete successfully, the
   * error is of type [*ExitError]. Other error types may be
   * returned for I/O problems.
   * 
   * If any of c.Stdin, c.Stdout or c.Stderr are not an [*os.File], Wait also waits
   * for the respective I/O loop copying to or from the process to complete.
   * 
   * Wait releases any resources associated with the [Cmd].
   */
  wait(): void
 }
 interface Cmd {
  /**
   * Output runs the command and returns its standard output.
   * Any returned error will usually be of type [*ExitError].
   * If c.Stderr was nil and the returned error is of type
   * [*ExitError], Output populates the Stderr field of the
   * returned error.
   */
  output(): string|Array<number>
 }
 interface Cmd {
  /**
   * CombinedOutput runs the command and returns its combined standard
   * output and standard error.
   */
  combinedOutput(): string|Array<number>
 }
 interface Cmd {
  /**
   * StdinPipe returns a pipe that will be connected to the command's
   * standard input when the command starts.
   * The pipe will be closed automatically after [Cmd.Wait] sees the command exit.
   * A caller need only call Close to force the pipe to close sooner.
   * For example, if the command being run will not exit until standard input
   * is closed, the caller must close the pipe.
   */
  stdinPipe(): io.WriteCloser
 }
 interface Cmd {
  /**
   * StdoutPipe returns a pipe that will be connected to the command's
   * standard output when the command starts.
   * 
   * [Cmd.Wait] will close the pipe after seeing the command exit, so most callers
   * need not close the pipe themselves. It is thus incorrect to call Wait
   * before all reads from the pipe have completed.
   * For the same reason, it is incorrect to call [Cmd.Run] when using StdoutPipe.
   * See the example for idiomatic usage.
   */
  stdoutPipe(): io.ReadCloser
 }
 interface Cmd {
  /**
   * StderrPipe returns a pipe that will be connected to the command's
   * standard error when the command starts.
   * 
   * [Cmd.Wait] will close the pipe after seeing the command exit, so most callers
   * need not close the pipe themselves. It is thus incorrect to call Wait
   * before all reads from the pipe have completed.
   * For the same reason, it is incorrect to use [Cmd.Run] when using StderrPipe.
   * See the StdoutPipe example for idiomatic usage.
   */
  stderrPipe(): io.ReadCloser
 }
 interface Cmd {
  /**
   * Environ returns a copy of the environment in which the command would be run
   * as it is currently configured.
   */
  environ(): Array<string>
 }
}

namespace mailer {
 /**
  * Message defines a generic email message struct.
  */
 interface Message {
  from: { address: string; name?: string; }
  to: Array<{ address: string; name?: string; }>
  bcc: Array<{ address: string; name?: string; }>
  cc: Array<{ address: string; name?: string; }>
  subject: string
  html: string
  text: string
  headers: _TygojaDict
  attachments: _TygojaDict
  inlineAttachments: _TygojaDict
 }
 /**
  * Mailer defines a base mail client interface.
  */
 interface Mailer {
  [key:string]: any;
  /**
   * Send sends an email with the provided Message.
   */
  send(message: Message): void
 }
}

/**
 * Package blob defines a lightweight abstration for interacting with
 * various storage services (local filesystem, S3, etc.).
 * 
 * NB!
 * For compatibility with earlier PocketBase versions and to prevent
 * unnecessary breaking changes, this package is based and implemented
 * as a minimal, stripped down version of the previously used gocloud.dev/blob.
 * While there is no promise that it won't diverge in the future to accommodate
 * better some PocketBase specific use cases, currently it copies and
 * tries to follow as close as possible the same implementations,
 * conventions and rules for the key escaping/unescaping, blob read/write
 * interfaces and struct options as gocloud.dev/blob, therefore the
 * credits goes to the original Go Cloud Development Kit Authors.
 */
namespace blob {
 /**
  * ListObject represents a single blob returned from List.
  */
 interface ListObject {
  /**
   * Key is the key for this blob.
   */
  key: string
  /**
   * ModTime is the time the blob was last modified.
   */
  modTime: time.Time
  /**
   * Size is the size of the blob's content in bytes.
   */
  size: number
  /**
   * MD5 is an MD5 hash of the blob contents or nil if not available.
   */
  md5: string|Array<number>
  /**
   * IsDir indicates that this result represents a "directory" in the
   * hierarchical namespace, ending in ListOptions.Delimiter. Key can be
   * passed as ListOptions.Prefix to list items in the "directory".
   * Fields other than Key and IsDir will not be set if IsDir is true.
   */
  isDir: boolean
 }
 /**
  * Attributes contains attributes about a blob.
  */
 interface Attributes {
  /**
   * CacheControl specifies caching attributes that services may use
   * when serving the blob.
   * https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
   */
  cacheControl: string
  /**
   * ContentDisposition specifies whether the blob content is expected to be
   * displayed inline or as an attachment.
   * https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
   */
  contentDisposition: string
  /**
   * ContentEncoding specifies the encoding used for the blob's content, if any.
   * https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
   */
  contentEncoding: string
  /**
   * ContentLanguage specifies the language used in the blob's content, if any.
   * https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Language
   */
  contentLanguage: string
  /**
   * ContentType is the MIME type of the blob. It will not be empty.
   * https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
   */
  contentType: string
  /**
   * Metadata holds key/value pairs associated with the blob.
   * Keys are guaranteed to be in lowercase, even if the backend service
   * has case-sensitive keys (although note that Metadata written via
   * this package will always be lowercased). If there are duplicate
   * case-insensitive keys (e.g., "foo" and "FOO"), only one value
   * will be kept, and it is undefined which one.
   */
  metadata: _TygojaDict
  /**
   * CreateTime is the time the blob was created, if available. If not available,
   * CreateTime will be the zero time.
   */
  createTime: time.Time
  /**
   * ModTime is the time the blob was last modified.
   */
  modTime: time.Time
  /**
   * Size is the size of the blob's content in bytes.
   */
  size: number
  /**
   * MD5 is an MD5 hash of the blob contents or nil if not available.
   */
  md5: string|Array<number>
  /**
   * ETag for the blob; see https://en.wikipedia.org/wiki/HTTP_ETag.
   */
  eTag: string
 }
 /**
  * Reader reads bytes from a blob.
  * It implements io.ReadSeekCloser, and must be closed after reads are finished.
  */
 interface Reader {
 }
 interface Reader {
  /**
   * Read implements io.Reader (https://golang.org/pkg/io/#Reader).
   */
  read(p: string|Array<number>): number
 }
 interface Reader {
  /**
   * Seek implements io.Seeker (https://golang.org/pkg/io/#Seeker).
   */
  seek(offset: number, whence: number): number
 }
 interface Reader {
  /**
   * Close implements io.Closer (https://golang.org/pkg/io/#Closer).
   */
  close(): void
 }
 interface Reader {
  /**
   * ContentType returns the MIME type of the blob.
   */
  contentType(): string
 }
 interface Reader {
  /**
   * ModTime returns the time the blob was last modified.
   */
  modTime(): time.Time
 }
 interface Reader {
  /**
   * Size returns the size of the blob content in bytes.
   */
  size(): number
 }
 interface Reader {
  /**
   * WriteTo reads from r and writes to w until there's no more data or
   * an error occurs.
   * The return value is the number of bytes written to w.
   * 
   * It implements the io.WriterTo interface.
   */
  writeTo(w: io.Writer): number
 }
}

/**
 * Package types implements some commonly used db serializable types
 * like datetime, json, etc.
 */
namespace types {
 /**
  * DateTime represents a [time.Time] instance in UTC that is wrapped
  * and serialized using the app default date layout.
  */
 interface DateTime {
 }
 interface DateTime {
  /**
   * Time returns the internal [time.Time] instance.
   */
  time(): time.Time
 }
 interface DateTime {
  /**
   * Add returns a new DateTime based on the current DateTime + the specified duration.
   */
  add(duration: time.Duration): DateTime
 }
 interface DateTime {
  /**
   * Sub returns a [time.Duration] by subtracting the specified DateTime from the current one.
   * 
   * If the result exceeds the maximum (or minimum) value that can be stored in a [time.Duration],
   * the maximum (or minimum) duration will be returned.
   */
  sub(u: DateTime): time.Duration
 }
 interface DateTime {
  /**
   * AddDate returns a new DateTime based on the current one + duration.
   * 
   * It follows the same rules as [time.AddDate].
   */
  addDate(years: number, months: number, days: number): DateTime
 }
 interface DateTime {
  /**
   * After reports whether the current DateTime instance is after u.
   */
  after(u: DateTime): boolean
 }
 interface DateTime {
  /**
   * Before reports whether the current DateTime instance is before u.
   */
  before(u: DateTime): boolean
 }
 interface DateTime {
  /**
   * Compare compares the current DateTime instance with u.
   * If the current instance is before u, it returns -1.
   * If the current instance is after u, it returns +1.
   * If they're the same, it returns 0.
   */
  compare(u: DateTime): number
 }
 interface DateTime {
  /**
   * Equal reports whether the current DateTime and u represent the same time instant.
   * Two DateTime can be equal even if they are in different locations.
   * For example, 6:00 +0200 and 4:00 UTC are Equal.
   */
  equal(u: DateTime): boolean
 }
 interface DateTime {
  /**
   * Unix returns the current DateTime as a Unix time, aka.
   * the number of seconds elapsed since January 1, 1970 UTC.
   */
  unix(): number
 }
 interface DateTime {
  /**
   * IsZero checks whether the current DateTime instance has zero time value.
   */
  isZero(): boolean
 }
 interface DateTime {
  /**
   * String serializes the current DateTime instance into a formatted
   * UTC date string.
   * 
   * The zero value is serialized to an empty string.
   */
  string(): string
 }
 interface DateTime {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface DateTime {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   */
  unmarshalJSON(b: string|Array<number>): void
 }
 interface DateTime {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface DateTime {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current DateTime instance.
   */
  scan(value: any): void
 }
 /**
  * GeoPoint defines a struct for storing geo coordinates as serialized json object
  * (e.g. {lon:0,lat:0}).
  * 
  * Note: using object notation and not a plain array to avoid the confusion
  * as there doesn't seem to be a fixed standard for the coordinates order.
  */
 interface GeoPoint {
  lon: number
  lat: number
 }
 interface GeoPoint {
  /**
   * String returns the string representation of the current GeoPoint instance.
   */
  string(): string
 }
 interface GeoPoint {
  /**
   * AsMap implements [core.mapExtractor] and returns a value suitable
   * to be used in an API rule expression.
   */
  asMap(): _TygojaDict
 }
 interface GeoPoint {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface GeoPoint {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current GeoPoint instance.
   * 
   * The value argument could be nil (no-op), another GeoPoint instance,
   * map or serialized json object with lat-lon props.
   */
  scan(value: any): void
 }
 /**
  * JSONArray defines a slice that is safe for json and db read/write.
  */
 interface JSONArray<T> extends Array<T>{}
 interface JSONArray<T> {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface JSONArray<T> {
  /**
   * String returns the string representation of the current json array.
   */
  string(): string
 }
 interface JSONArray<T> {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface JSONArray<T> {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current JSONArray[T] instance.
   */
  scan(value: any): void
 }
 /**
  * JSONMap defines a map that is safe for json and db read/write.
  */
 interface JSONMap<T> extends _TygojaDict{}
 interface JSONMap<T> {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface JSONMap<T> {
  /**
   * String returns the string representation of the current json map.
   */
  string(): string
 }
 interface JSONMap<T> {
  /**
   * Get retrieves a single value from the current JSONMap[T].
   * 
   * This helper was added primarily to assist the goja integration since custom map types
   * don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
   */
  get(key: string): T
 }
 interface JSONMap<T> {
  /**
   * Set sets a single value in the current JSONMap[T].
   * 
   * This helper was added primarily to assist the goja integration since custom map types
   * don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
   */
  set(key: string, value: T): void
 }
 interface JSONMap<T> {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface JSONMap<T> {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current JSONMap[T] instance.
   */
  scan(value: any): void
 }
 /**
  * JSONRaw defines a json value type that is safe for db read/write.
  */
 interface JSONRaw extends Array<number>{}
 interface JSONRaw {
  /**
   * String returns the current JSONRaw instance as a json encoded string.
   */
  string(): string
 }
 interface JSONRaw {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string|Array<number>
 }
 interface JSONRaw {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   */
  unmarshalJSON(b: string|Array<number>): void
 }
 interface JSONRaw {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface JSONRaw {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current JSONRaw instance.
   */
  scan(value: any): void
 }
}

namespace search {
 /**
  * Result defines the returned search result structure.
  */
 interface Result {
  items: any
  page: number
  perPage: number
  totalItems: number
  totalPages: number
 }
 /**
  * ResolverResult defines a single FieldResolver.Resolve() successfully parsed result.
  */
 interface ResolverResult {
  /**
   * Identifier is the plain SQL identifier/column that will be used
   * in the final db expression as left or right operand.
   */
  identifier: string
  /**
   * NoCoalesce instructs to not use COALESCE or NULL fallbacks
   * when building the identifier expression.
   */
  noCoalesce: boolean
  /**
   * Params is a map with db placeholder->value pairs that will be added
   * to the query when building both resolved operands/sides in a single expression.
   */
  params: dbx.Params
  /**
   * MultiMatchSubQuery is an optional sub query expression that will be added
   * in addition to the combined ResolverResult expression during build.
   */
  multiMatchSubQuery: dbx.Expression
  /**
   * AfterBuild is an optional function that will be called after building
   * and combining the result of both resolved operands/sides in a single expression.
   */
  afterBuild: (expr: dbx.Expression) => dbx.Expression
 }
}

namespace router {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * ApiError defines the struct for a basic api error response.
  */
 interface ApiError {
  data: _TygojaDict
  message: string
  status: number
 }
 interface ApiError {
  /**
   * Error makes it compatible with the `error` interface.
   */
  error(): string
 }
 interface ApiError {
  /**
   * RawData returns the unformatted error data (could be an internal error, text, etc.)
   */
  rawData(): any
 }
 interface ApiError {
  /**
   * Is reports whether the current ApiError wraps the target.
   */
  is(target: Error): boolean
 }
 /**
  * Event specifies based Route handler event that is usually intended
  * to be embedded as part of a custom event struct.
  * 
  * NB! It is expected that the Response and Request fields are always set.
  */
 type _scPkqIP = hook.Event
 interface Event extends _scPkqIP {
  response: http.ResponseWriter
  request?: http.Request
 }
 interface Event {
  /**
   * Written reports whether the current response has already been written.
   * 
   * This method always returns false if e.ResponseWritter doesn't implement the WriteTracker interface
   * (all router package handlers receives a ResponseWritter that implements it unless explicitly replaced with a custom one).
   */
  written(): boolean
 }
 interface Event {
  /**
   * Status reports the status code of the current response.
   * 
   * This method always returns 0 if e.Response doesn't implement the StatusTracker interface
   * (all router package handlers receives a ResponseWritter that implements it unless explicitly replaced with a custom one).
   */
  status(): number
 }
 interface Event {
  /**
   * Flush flushes buffered data to the current response.
   * 
   * Returns [http.ErrNotSupported] if e.Response doesn't implement the [http.Flusher] interface
   * (all router package handlers receives a ResponseWritter that implements it unless explicitly replaced with a custom one).
   */
  flush(): void
 }
 interface Event {
  /**
   * IsTLS reports whether the connection on which the request was received is TLS.
   */
  isTLS(): boolean
 }
 interface Event {
  /**
   * SetCookie is an alias for [http.SetCookie].
   * 
   * SetCookie adds a Set-Cookie header to the current response's headers.
   * The provided cookie must have a valid Name.
   * Invalid cookies may be silently dropped.
   */
  setCookie(cookie: http.Cookie): void
 }
 interface Event {
  /**
   * RemoteIP returns the IP address of the client that sent the request.
   * 
   * IPv6 addresses are returned expanded.
   * For example, "2001:db8::1" becomes "2001:0db8:0000:0000:0000:0000:0000:0001".
   * 
   * Note that if you are behind reverse proxy(ies), this method returns
   * the IP of the last connecting proxy.
   */
  remoteIP(): string
 }
 interface Event {
  /**
   * FindUploadedFiles extracts all form files of "key" from a http request
   * and returns a slice with filesystem.File instances (if any).
   */
  findUploadedFiles(key: string): Array<(filesystem.File | undefined)>
 }
 interface Event {
  /**
   * Get retrieves single value from the current event data store.
   */
  get(key: string): any
 }
 interface Event {
  /**
   * GetAll returns a copy of the current event data store.
   */
  getAll(): _TygojaDict
 }
 interface Event {
  /**
   * Set saves single value into the current event data store.
   */
  set(key: string, value: any): void
 }
 interface Event {
  /**
   * SetAll saves all items from m into the current event data store.
   */
  setAll(m: _TygojaDict): void
 }
 interface Event {
  /**
   * String writes a plain string response.
   */
  string(status: number, data: string): void
 }
 interface Event {
  /**
   * HTML writes an HTML response.
   */
  html(status: number, data: string): void
 }
 interface Event {
  /**
   * JSON writes a JSON response.
   * 
   * It also provides a generic response data fields picker if the "fields" query parameter is set.
   * For example, if you are requesting `?fields=a,b` for `e.JSON(200, map[string]int{ "a":1, "b":2, "c":3 })`,
   * it should result in a JSON response like: `{"a":1, "b": 2}`.
   */
  json(status: number, data: any): void
 }
 interface Event {
  /**
   * XML writes an XML response.
   * It automatically prepends the generic [xml.Header] string to the response.
   */
  xml(status: number, data: any): void
 }
 interface Event {
  /**
   * Stream streams the specified reader into the response.
   */
  stream(status: number, contentType: string, reader: io.Reader): void
 }
 interface Event {
  /**
   * Blob writes a blob (bytes slice) response.
   */
  blob(status: number, contentType: string, b: string|Array<number>): void
 }
 interface Event {
  /**
   * FileFS serves the specified filename from fsys.
   * 
   * It is similar to [echo.FileFS] for consistency with earlier versions.
   */
  fileFS(fsys: fs.FS, filename: string): void
 }
 interface Event {
  /**
   * NoContent writes a response with no body (ex. 204).
   */
  noContent(status: number): void
 }
 interface Event {
  /**
   * Redirect writes a redirect response to the specified url.
   * The status code must be in between 300 – 399 range.
   */
  redirect(status: number, url: string): void
 }
 interface Event {
  error(status: number, message: string, errData: any): (ApiError)
 }
 interface Event {
  badRequestError(message: string, errData: any): (ApiError)
 }
 interface Event {
  notFoundError(message: string, errData: any): (ApiError)
 }
 interface Event {
  forbiddenError(message: string, errData: any): (ApiError)
 }
 interface Event {
  unauthorizedError(message: string, errData: any): (ApiError)
 }
 interface Event {
  tooManyRequestsError(message: string, errData: any): (ApiError)
 }
 interface Event {
  internalServerError(message: string, errData: any): (ApiError)
 }
 interface Event {
  /**
   * BindBody unmarshal the request body into the provided dst.
   * 
   * dst must be either a struct pointer or map[string]any.
   * 
   * The rules how the body will be scanned depends on the request Content-Type.
   * 
   * Currently the following Content-Types are supported:
   * ```
   *   - application/json
   *   - text/xml, application/xml
   *   - multipart/form-data, application/x-www-form-urlencoded
   * ```
   * 
   * Respectively the following struct tags are supported (again, which one will be used depends on the Content-Type):
   * ```
   *   - "json" (json body)- uses the builtin Go json package for unmarshaling.
   *   - "xml" (xml body) - uses the builtin Go xml package for unmarshaling.
   *   - "form" (form data) - utilizes the custom [router.UnmarshalRequestData] method.
   * ```
   * 
   * NB! When dst is a struct make sure that it doesn't have public fields
   * that shouldn't be bindable and it is advisible such fields to be unexported
   * or have a separate struct just for the binding. For example:
   * 
   * ```
   * 	data := struct{
   * 	   somethingPrivate string
   * 
   * 	   Title string `json:"title" form:"title"`
   * 	   Total int    `json:"total" form:"total"`
   * 	}
   * 	err := e.BindBody(&data)
   * ```
   */
  bindBody(dst: any): void
 }
 /**
  * Router defines a thin wrapper around the standard Go [http.ServeMux] by
  * adding support for routing sub-groups, middlewares and other common utils.
  * 
  * Example:
  * 
  * ```
  * 	r := NewRouter[*MyEvent](eventFactory)
  * 
  * 	// middlewares
  * 	r.BindFunc(m1, m2)
  * 
  * 	// routes
  * 	r.GET("/test", handler1)
  * 
  * 	// sub-routers/groups
  * 	api := r.Group("/api")
  * 	api.GET("/admins", handler2)
  * 
  * 	// generate a http.ServeMux instance based on the router configurations
  * 	mux, _ := r.BuildMux()
  * 
  * 	http.ListenAndServe("localhost:8090", mux)
  * ```
  */
 type _sMoYDMq<T> = RouterGroup<T>
 interface Router<T> extends _sMoYDMq<T> {
 }
 interface Router<T> {
  /**
   * BuildMux constructs a new mux [http.Handler] instance from the current router configurations.
   */
  buildMux(): http.Handler
 }
}

/**
 * Package slog provides structured logging,
 * in which log records include a message,
 * a severity level, and various other attributes
 * expressed as key-value pairs.
 * 
 * It defines a type, [Logger],
 * which provides several methods (such as [Logger.Info] and [Logger.Error])
 * for reporting events of interest.
 * 
 * Each Logger is associated with a [Handler].
 * A Logger output method creates a [Record] from the method arguments
 * and passes it to the Handler, which decides how to handle it.
 * There is a default Logger accessible through top-level functions
 * (such as [Info] and [Error]) that call the corresponding Logger methods.
 * 
 * A log record consists of a time, a level, a message, and a set of key-value
 * pairs, where the keys are strings and the values may be of any type.
 * As an example,
 * 
 * ```
 * 	slog.Info("hello", "count", 3)
 * ```
 * 
 * creates a record containing the time of the call,
 * a level of Info, the message "hello", and a single
 * pair with key "count" and value 3.
 * 
 * The [Info] top-level function calls the [Logger.Info] method on the default Logger.
 * In addition to [Logger.Info], there are methods for Debug, Warn and Error levels.
 * Besides these convenience methods for common levels,
 * there is also a [Logger.Log] method which takes the level as an argument.
 * Each of these methods has a corresponding top-level function that uses the
 * default logger.
 * 
 * The default handler formats the log record's message, time, level, and attributes
 * as a string and passes it to the [log] package.
 * 
 * ```
 * 	2022/11/08 15:28:26 INFO hello count=3
 * ```
 * 
 * For more control over the output format, create a logger with a different handler.
 * This statement uses [New] to create a new logger with a [TextHandler]
 * that writes structured records in text form to standard error:
 * 
 * ```
 * 	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
 * ```
 * 
 * [TextHandler] output is a sequence of key=value pairs, easily and unambiguously
 * parsed by machine. This statement:
 * 
 * ```
 * 	logger.Info("hello", "count", 3)
 * ```
 * 
 * produces this output:
 * 
 * ```
 * 	time=2022-11-08T15:28:26.000-05:00 level=INFO msg=hello count=3
 * ```
 * 
 * The package also provides [JSONHandler], whose output is line-delimited JSON:
 * 
 * ```
 * 	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
 * 	logger.Info("hello", "count", 3)
 * ```
 * 
 * produces this output:
 * 
 * ```
 * 	{"time":"2022-11-08T15:28:26.000000000-05:00","level":"INFO","msg":"hello","count":3}
 * ```
 * 
 * Both [TextHandler] and [JSONHandler] can be configured with [HandlerOptions].
 * There are options for setting the minimum level (see Levels, below),
 * displaying the source file and line of the log call, and
 * modifying attributes before they are logged.
 * 
 * Setting a logger as the default with
 * 
 * ```
 * 	slog.SetDefault(logger)
 * ```
 * 
 * will cause the top-level functions like [Info] to use it.
 * [SetDefault] also updates the default logger used by the [log] package,
 * so that existing applications that use [log.Printf] and related functions
 * will send log records to the logger's handler without needing to be rewritten.
 * 
 * Some attributes are common to many log calls.
 * For example, you may wish to include the URL or trace identifier of a server request
 * with all log events arising from the request.
 * Rather than repeat the attribute with every log call, you can use [Logger.With]
 * to construct a new Logger containing the attributes:
 * 
 * ```
 * 	logger2 := logger.With("url", r.URL)
 * ```
 * 
 * The arguments to With are the same key-value pairs used in [Logger.Info].
 * The result is a new Logger with the same handler as the original, but additional
 * attributes that will appear in the output of every call.
 * 
 * # Levels
 * 
 * A [Level] is an integer representing the importance or severity of a log event.
 * The higher the level, the more severe the event.
 * This package defines constants for the most common levels,
 * but any int can be used as a level.
 * 
 * In an application, you may wish to log messages only at a certain level or greater.
 * One common configuration is to log messages at Info or higher levels,
 * suppressing debug logging until it is needed.
 * The built-in handlers can be configured with the minimum level to output by
 * setting [HandlerOptions.Level].
 * The program's `main` function typically does this.
 * The default value is LevelInfo.
 * 
 * Setting the [HandlerOptions.Level] field to a [Level] value
 * fixes the handler's minimum level throughout its lifetime.
 * Setting it to a [LevelVar] allows the level to be varied dynamically.
 * A LevelVar holds a Level and is safe to read or write from multiple
 * goroutines.
 * To vary the level dynamically for an entire program, first initialize
 * a global LevelVar:
 * 
 * ```
 * 	var programLevel = new(slog.LevelVar) // Info by default
 * ```
 * 
 * Then use the LevelVar to construct a handler, and make it the default:
 * 
 * ```
 * 	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
 * 	slog.SetDefault(slog.New(h))
 * ```
 * 
 * Now the program can change its logging level with a single statement:
 * 
 * ```
 * 	programLevel.Set(slog.LevelDebug)
 * ```
 * 
 * # Groups
 * 
 * Attributes can be collected into groups.
 * A group has a name that is used to qualify the names of its attributes.
 * How this qualification is displayed depends on the handler.
 * [TextHandler] separates the group and attribute names with a dot.
 * [JSONHandler] treats each group as a separate JSON object, with the group name as the key.
 * 
 * Use [Group] to create a Group attribute from a name and a list of key-value pairs:
 * 
 * ```
 * 	slog.Group("request",
 * 	    "method", r.Method,
 * 	    "url", r.URL)
 * ```
 * 
 * TextHandler would display this group as
 * 
 * ```
 * 	request.method=GET request.url=http://example.com
 * ```
 * 
 * JSONHandler would display it as
 * 
 * ```
 * 	"request":{"method":"GET","url":"http://example.com"}
 * ```
 * 
 * Use [Logger.WithGroup] to qualify all of a Logger's output
 * with a group name. Calling WithGroup on a Logger results in a
 * new Logger with the same Handler as the original, but with all
 * its attributes qualified by the group name.
 * 
 * This can help prevent duplicate attribute keys in large systems,
 * where subsystems might use the same keys.
 * Pass each subsystem a different Logger with its own group name so that
 * potential duplicates are qualified:
 * 
 * ```
 * 	logger := slog.Default().With("id", systemID)
 * 	parserLogger := logger.WithGroup("parser")
 * 	parseInput(input, parserLogger)
 * ```
 * 
 * When parseInput logs with parserLogger, its keys will be qualified with "parser",
 * so even if it uses the common key "id", the log line will have distinct keys.
 * 
 * # Contexts
 * 
 * Some handlers may wish to include information from the [context.Context] that is
 * available at the call site. One example of such information
 * is the identifier for the current span when tracing is enabled.
 * 
 * The [Logger.Log] and [Logger.LogAttrs] methods take a context as a first
 * argument, as do their corresponding top-level functions.
 * 
 * Although the convenience methods on Logger (Info and so on) and the
 * corresponding top-level functions do not take a context, the alternatives ending
 * in "Context" do. For example,
 * 
 * ```
 * 	slog.InfoContext(ctx, "message")
 * ```
 * 
 * It is recommended to pass a context to an output method if one is available.
 * 
 * # Attrs and Values
 * 
 * An [Attr] is a key-value pair. The Logger output methods accept Attrs as well as
 * alternating keys and values. The statement
 * 
 * ```
 * 	slog.Info("hello", slog.Int("count", 3))
 * ```
 * 
 * behaves the same as
 * 
 * ```
 * 	slog.Info("hello", "count", 3)
 * ```
 * 
 * There are convenience constructors for [Attr] such as [Int], [String], and [Bool]
 * for common types, as well as the function [Any] for constructing Attrs of any
 * type.
 * 
 * The value part of an Attr is a type called [Value].
 * Like an [any], a Value can hold any Go value,
 * but it can represent typical values, including all numbers and strings,
 * without an allocation.
 * 
 * For the most efficient log output, use [Logger.LogAttrs].
 * It is similar to [Logger.Log] but accepts only Attrs, not alternating
 * keys and values; this allows it, too, to avoid allocation.
 * 
 * The call
 * 
 * ```
 * 	logger.LogAttrs(ctx, slog.LevelInfo, "hello", slog.Int("count", 3))
 * ```
 * 
 * is the most efficient way to achieve the same output as
 * 
 * ```
 * 	slog.InfoContext(ctx, "hello", "count", 3)
 * ```
 * 
 * # Customizing a type's logging behavior
 * 
 * If a type implements the [LogValuer] interface, the [Value] returned from its LogValue
 * method is used for logging. You can use this to control how values of the type
 * appear in logs. For example, you can redact secret information like passwords,
 * or gather a struct's fields in a Group. See the examples under [LogValuer] for
 * details.
 * 
 * A LogValue method may return a Value that itself implements [LogValuer]. The [Value.Resolve]
 * method handles these cases carefully, avoiding infinite loops and unbounded recursion.
 * Handler authors and others may wish to use [Value.Resolve] instead of calling LogValue directly.
 * 
 * # Wrapping output methods
 * 
 * The logger functions use reflection over the call stack to find the file name
 * and line number of the logging call within the application. This can produce
 * incorrect source information for functions that wrap slog. For instance, if you
 * define this function in file mylog.go:
 * 
 * ```
 * 	func Infof(logger *slog.Logger, format string, args ...any) {
 * 	    logger.Info(fmt.Sprintf(format, args...))
 * 	}
 * ```
 * 
 * and you call it like this in main.go:
 * 
 * ```
 * 	Infof(slog.Default(), "hello, %s", "world")
 * ```
 * 
 * then slog will report the source file as mylog.go, not main.go.
 * 
 * A correct implementation of Infof will obtain the source location
 * (pc) and pass it to NewRecord.
 * The Infof function in the package-level example called "wrapping"
 * demonstrates how to do this.
 * 
 * # Working with Records
 * 
 * Sometimes a Handler will need to modify a Record
 * before passing it on to another Handler or backend.
 * A Record contains a mixture of simple public fields (e.g. Time, Level, Message)
 * and hidden fields that refer to state (such as attributes) indirectly. This
 * means that modifying a simple copy of a Record (e.g. by calling
 * [Record.Add] or [Record.AddAttrs] to add attributes)
 * may have unexpected effects on the original.
 * Before modifying a Record, use [Record.Clone] to
 * create a copy that shares no state with the original,
 * or create a new Record with [NewRecord]
 * and build up its Attrs by traversing the old ones with [Record.Attrs].
 * 
 * # Performance considerations
 * 
 * If profiling your application demonstrates that logging is taking significant time,
 * the following suggestions may help.
 * 
 * If many log lines have a common attribute, use [Logger.With] to create a Logger with
 * that attribute. The built-in handlers will format that attribute only once, at the
 * call to [Logger.With]. The [Handler] interface is designed to allow that optimization,
 * and a well-written Handler should take advantage of it.
 * 
 * The arguments to a log call are always evaluated, even if the log event is discarded.
 * If possible, defer computation so that it happens only if the value is actually logged.
 * For example, consider the call
 * 
 * ```
 * 	slog.Info("starting request", "url", r.URL.String())  // may compute String unnecessarily
 * ```
 * 
 * The URL.String method will be called even if the logger discards Info-level events.
 * Instead, pass the URL directly:
 * 
 * ```
 * 	slog.Info("starting request", "url", &r.URL) // calls URL.String only if needed
 * ```
 * 
 * The built-in [TextHandler] will call its String method, but only
 * if the log event is enabled.
 * Avoiding the call to String also preserves the structure of the underlying value.
 * For example [JSONHandler] emits the components of the parsed URL as a JSON object.
 * If you want to avoid eagerly paying the cost of the String call
 * without causing the handler to potentially inspect the structure of the value,
 * wrap the value in a fmt.Stringer implementation that hides its Marshal methods.
 * 
 * You can also use the [LogValuer] interface to avoid unnecessary work in disabled log
 * calls. Say you need to log some expensive value:
 * 
 * ```
 * 	slog.Debug("frobbing", "value", computeExpensiveValue(arg))
 * ```
 * 
 * Even if this line is disabled, computeExpensiveValue will be called.
 * To avoid that, define a type implementing LogValuer:
 * 
 * ```
 * 	type expensive struct { arg int }
 * 
 * 	func (e expensive) LogValue() slog.Value {
 * 	    return slog.AnyValue(computeExpensiveValue(e.arg))
 * 	}
 * ```
 * 
 * Then use a value of that type in log calls:
 * 
 * ```
 * 	slog.Debug("frobbing", "value", expensive{arg})
 * ```
 * 
 * Now computeExpensiveValue will only be called when the line is enabled.
 * 
 * The built-in handlers acquire a lock before calling [io.Writer.Write]
 * to ensure that exactly one [Record] is written at a time in its entirety.
 * Although each log record has a timestamp,
 * the built-in handlers do not use that time to sort the written records.
 * User-defined handlers are responsible for their own locking and sorting.
 * 
 * # Writing a handler
 * 
 * For a guide to writing a custom handler, see https://golang.org/s/slog-handler-guide.
 */
namespace slog {
 // @ts-ignore
 import loginternal = internal
 /**
  * A Logger records structured information about each call to its
  * Log, Debug, Info, Warn, and Error methods.
  * For each call, it creates a [Record] and passes it to a [Handler].
  * 
  * To create a new Logger, call [New] or a Logger method
  * that begins "With".
  */
 interface Logger {
 }
 interface Logger {
  /**
   * Handler returns l's Handler.
   */
  handler(): Handler
 }
 interface Logger {
  /**
   * With returns a Logger that includes the given attributes
   * in each output operation. Arguments are converted to
   * attributes as if by [Logger.Log].
   */
  with(...args: any[]): (Logger)
 }
 interface Logger {
  /**
   * WithGroup returns a Logger that starts a group, if name is non-empty.
   * The keys of all attributes added to the Logger will be qualified by the given
   * name. (How that qualification happens depends on the [Handler.WithGroup]
   * method of the Logger's Handler.)
   * 
   * If name is empty, WithGroup returns the receiver.
   */
  withGroup(name: string): (Logger)
 }
 interface Logger {
  /**
   * Enabled reports whether l emits log records at the given context and level.
   */
  enabled(ctx: context.Context, level: Level): boolean
 }
 interface Logger {
  /**
   * Log emits a log record with the current time and the given level and message.
   * The Record's Attrs consist of the Logger's attributes followed by
   * the Attrs specified by args.
   * 
   * The attribute arguments are processed as follows:
   * ```
   *   - If an argument is an Attr, it is used as is.
   *   - If an argument is a string and this is not the last argument,
   *     the following argument is treated as the value and the two are combined
   *     into an Attr.
   *   - Otherwise, the argument is treated as a value with key "!BADKEY".
   * ```
   */
  log(ctx: context.Context, level: Level, msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
   */
  logAttrs(ctx: context.Context, level: Level, msg: string, ...attrs: Attr[]): void
 }
 interface Logger {
  /**
   * Debug logs at [LevelDebug].
   */
  debug(msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * DebugContext logs at [LevelDebug] with the given context.
   */
  debugContext(ctx: context.Context, msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * Info logs at [LevelInfo].
   */
  info(msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * InfoContext logs at [LevelInfo] with the given context.
   */
  infoContext(ctx: context.Context, msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * Warn logs at [LevelWarn].
   */
  warn(msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * WarnContext logs at [LevelWarn] with the given context.
   */
  warnContext(ctx: context.Context, msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * Error logs at [LevelError].
   */
  error(msg: string, ...args: any[]): void
 }
 interface Logger {
  /**
   * ErrorContext logs at [LevelError] with the given context.
   */
  errorContext(ctx: context.Context, msg: string, ...args: any[]): void
 }
}

namespace subscriptions {
 /**
  * Broker defines a struct for managing subscriptions clients.
  */
 interface Broker {
 }
 interface Broker {
  /**
   * Clients returns a shallow copy of all registered clients indexed
   * with their connection id.
   */
  clients(): _TygojaDict
 }
 interface Broker {
  /**
   * ChunkedClients splits the current clients into a chunked slice.
   */
  chunkedClients(chunkSize: number): Array<Array<Client>>
 }
 interface Broker {
  /**
   * TotalClients returns the total number of registered clients.
   */
  totalClients(): number
 }
 interface Broker {
  /**
   * ClientById finds a registered client by its id.
   * 
   * Returns non-nil error when client with clientId is not registered.
   */
  clientById(clientId: string): Client
 }
 interface Broker {
  /**
   * Register adds a new client to the broker instance.
   */
  register(client: Client): void
 }
 interface Broker {
  /**
   * Unregister removes a single client by its id and marks it as discarded.
   * 
   * If client with clientId doesn't exist, this method does nothing.
   */
  unregister(clientId: string): void
 }
 /**
  * Client is an interface for a generic subscription client.
  */
 interface Client {
  [key:string]: any;
  /**
   * Id Returns the unique id of the client.
   */
  id(): string
  /**
   * Channel returns the client's communication channel.
   * 
   * NB! The channel shouldn't be used after calling Discard().
   */
  channel(): undefined
  /**
   * Subscriptions returns a shallow copy of the client subscriptions matching the prefixes.
   * If no prefix is specified, returns all subscriptions.
   */
  subscriptions(...prefixes: string[]): _TygojaDict
  /**
   * Subscribe subscribes the client to the provided subscriptions list.
   * 
   * Each subscription can also have "options" (json serialized SubscriptionOptions) as query parameter.
   * 
   * Example:
   * 
   * ```
   * 	Subscribe(
   * 	    "subscriptionA",
   * 	    `subscriptionB?options={"query":{"a":1},"headers":{"x_token":"abc"}}`,
   * 	)
   * ```
   */
  subscribe(...subs: string[]): void
  /**
   * Unsubscribe unsubscribes the client from the provided subscriptions list.
   */
  unsubscribe(...subs: string[]): void
  /**
   * HasSubscription checks if the client is subscribed to `sub`.
   */
  hasSubscription(sub: string): boolean
  /**
   * Set stores any value to the client's context.
   */
  set(key: string, value: any): void
  /**
   * Unset removes a single value from the client's context.
   */
  unset(key: string): void
  /**
   * Get retrieves the key value from the client's context.
   */
  get(key: string): any
  /**
   * Discard marks the client as "discarded" (and closes its channel),
   * meaning that it shouldn't be used anymore for sending new messages.
   * 
   * It is safe to call Discard() multiple times.
   */
  discard(): void
  /**
   * IsDiscarded indicates whether the client has been "discarded"
   * and should no longer be used.
   */
  isDiscarded(): boolean
  /**
   * Send sends the specified message to the client's channel (if not discarded).
   */
  send(m: Message): void
 }
 /**
  * Message defines a client's channel data.
  */
 interface Message {
  name: string
  data: string|Array<number>
 }
 interface Message {
  /**
   * WriteSSE writes the current message in a SSE format into the provided writer.
   * 
   * For example, writing to a router.Event:
   * 
   * ```
   * 	m := Message{Name: "users/create", Data: []byte{...}}
   * 	m.WriteSSE(e.Response, "yourEventId")
   * 	e.Flush()
   * ```
   */
  writeSSE(w: io.Writer, eventId: string): void
 }
}

namespace auth {
 /**
  * Provider defines a common interface for an OAuth2 client.
  */
 interface Provider {
  [key:string]: any;
  /**
   * Context returns the context associated with the provider (if any).
   */
  context(): context.Context
  /**
   * SetContext assigns the specified context to the current provider.
   */
  setContext(ctx: context.Context): void
  /**
   * PKCE indicates whether the provider can use the PKCE flow.
   */
  pkce(): boolean
  /**
   * SetPKCE toggles the state whether the provider can use the PKCE flow or not.
   */
  setPKCE(enable: boolean): void
  /**
   * DisplayName usually returns provider name as it is officially written
   * and it could be used directly in the UI.
   */
  displayName(): string
  /**
   * SetDisplayName sets the provider's display name.
   */
  setDisplayName(displayName: string): void
  /**
   * Scopes returns the provider access permissions that will be requested.
   */
  scopes(): Array<string>
  /**
   * SetScopes sets the provider access permissions that will be requested later.
   */
  setScopes(scopes: Array<string>): void
  /**
   * ClientId returns the provider client's app ID.
   */
  clientId(): string
  /**
   * SetClientId sets the provider client's ID.
   */
  setClientId(clientId: string): void
  /**
   * ClientSecret returns the provider client's app secret.
   */
  clientSecret(): string
  /**
   * SetClientSecret sets the provider client's app secret.
   */
  setClientSecret(secret: string): void
  /**
   * RedirectURL returns the end address to redirect the user
   * going through the OAuth flow.
   */
  redirectURL(): string
  /**
   * SetRedirectURL sets the provider's RedirectURL.
   */
  setRedirectURL(url: string): void
  /**
   * AuthURL returns the provider's authorization service url.
   */
  authURL(): string
  /**
   * SetAuthURL sets the provider's AuthURL.
   */
  setAuthURL(url: string): void
  /**
   * TokenURL returns the provider's token exchange service url.
   */
  tokenURL(): string
  /**
   * SetTokenURL sets the provider's TokenURL.
   */
  setTokenURL(url: string): void
  /**
   * UserInfoURL returns the provider's user info api url.
   */
  userInfoURL(): string
  /**
   * SetUserInfoURL sets the provider's UserInfoURL.
   */
  setUserInfoURL(url: string): void
  /**
   * Extra returns a shallow copy of any custom config data
   * that the provider may be need.
   */
  extra(): _TygojaDict
  /**
   * SetExtra updates the provider's custom config data.
   */
  setExtra(data: _TygojaDict): void
  /**
   * Client returns an http client using the provided token.
   */
  client(token: oauth2.Token): (any)
  /**
   * BuildAuthURL returns a URL to the provider's consent page
   * that asks for permissions for the required scopes explicitly.
   */
  buildAuthURL(state: string, ...opts: oauth2.AuthCodeOption[]): string
  /**
   * FetchToken converts an authorization code to token.
   */
  fetchToken(code: string, ...opts: oauth2.AuthCodeOption[]): (oauth2.Token)
  /**
   * FetchRawUserInfo requests and marshalizes into `result` the
   * the OAuth user api response.
   */
  fetchRawUserInfo(token: oauth2.Token): string|Array<number>
  /**
   * FetchAuthUser is similar to FetchRawUserInfo, but normalizes and
   * marshalizes the user api response into a standardized AuthUser struct.
   */
  fetchAuthUser(token: oauth2.Token): (AuthUser)
 }
 /**
  * AuthUser defines a standardized OAuth2 user data structure.
  */
 interface AuthUser {
  expiry: types.DateTime
  rawUser: _TygojaDict
  id: string
  name: string
  username: string
  email: string
  avatarURL: string
  accessToken: string
  refreshToken: string
  /**
   * @todo
   * deprecated: use AvatarURL instead
   * AvatarUrl will be removed after dropping v0.22 support
   */
  avatarUrl: string
 }
 interface AuthUser {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   * 
   * @todo remove after dropping v0.22 support
   */
  marshalJSON(): string|Array<number>
 }
}

/**
 * Package cobra is a commander providing a simple interface to create powerful modern CLI interfaces.
 * In addition to providing an interface, Cobra simultaneously provides a controller to organize your application code.
 */
namespace cobra {
 interface Command {
  /**
   * GenBashCompletion generates bash completion file and writes to the passed writer.
   */
  genBashCompletion(w: io.Writer): void
 }
 interface Command {
  /**
   * GenBashCompletionFile generates bash completion file.
   */
  genBashCompletionFile(filename: string): void
 }
 interface Command {
  /**
   * GenBashCompletionFileV2 generates Bash completion version 2.
   */
  genBashCompletionFileV2(filename: string, includeDesc: boolean): void
 }
 interface Command {
  /**
   * GenBashCompletionV2 generates Bash completion file version 2
   * and writes it to the passed writer.
   */
  genBashCompletionV2(w: io.Writer, includeDesc: boolean): void
 }
 // @ts-ignore
 import flag = pflag
 /**
  * Command is just that, a command for your application.
  * E.g.  'go run ...' - 'run' is the command. Cobra requires
  * you to define the usage and description as part of your command
  * definition to ensure usability.
  */
 interface Command {
  /**
   * Use is the one-line usage message.
   * Recommended syntax is as follows:
   * ```
   *   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
   *   ... indicates that you can specify multiple values for the previous argument.
   *   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
   *       argument to the right of the separator. You cannot use both arguments in a single use of the command.
   *   { } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
   *       optional, they are enclosed in brackets ([ ]).
   * ```
   * Example: add [-F file | -D dir]... [-f format] profile
   */
  use: string
  /**
   * Aliases is an array of aliases that can be used instead of the first word in Use.
   */
  aliases: Array<string>
  /**
   * SuggestFor is an array of command names for which this command will be suggested -
   * similar to aliases but only suggests.
   */
  suggestFor: Array<string>
  /**
   * Short is the short description shown in the 'help' output.
   */
  short: string
  /**
   * The group id under which this subcommand is grouped in the 'help' output of its parent.
   */
  groupID: string
  /**
   * Long is the long message shown in the 'help <this-command>' output.
   */
  long: string
  /**
   * Example is examples of how to use the command.
   */
  example: string
  /**
   * ValidArgs is list of all valid non-flag arguments that are accepted in shell completions
   */
  validArgs: Array<Completion>
  /**
   * ValidArgsFunction is an optional function that provides valid non-flag arguments for shell completion.
   * It is a dynamic version of using ValidArgs.
   * Only one of ValidArgs and ValidArgsFunction can be used for a command.
   */
  validArgsFunction: CompletionFunc
  /**
   * Expected arguments
   */
  args: PositionalArgs
  /**
   * ArgAliases is List of aliases for ValidArgs.
   * These are not suggested to the user in the shell completion,
   * but accepted if entered manually.
   */
  argAliases: Array<string>
  /**
   * BashCompletionFunction is custom bash functions used by the legacy bash autocompletion generator.
   * For portability with other shells, it is recommended to instead use ValidArgsFunction
   */
  bashCompletionFunction: string
  /**
   * Deprecated defines, if this command is deprecated and should print this string when used.
   */
  deprecated: string
  /**
   * Annotations are key/value pairs that can be used by applications to identify or
   * group commands or set special options.
   */
  annotations: _TygojaDict
  /**
   * Version defines the version for this command. If this value is non-empty and the command does not
   * define a "version" flag, a "version" boolean flag will be added to the command and, if specified,
   * will print content of the "Version" variable. A shorthand "v" flag will also be added if the
   * command does not define one.
   */
  version: string
  /**
   * The *Run functions are executed in the following order:
   * ```
   *   * PersistentPreRun()
   *   * PreRun()
   *   * Run()
   *   * PostRun()
   *   * PersistentPostRun()
   * ```
   * All functions get the same args, the arguments after the command name.
   * The *PreRun and *PostRun functions will only be executed if the Run function of the current
   * command has been declared.
   * 
   * PersistentPreRun: children of this command will inherit and execute.
   */
  persistentPreRun: (cmd: Command, args: Array<string>) => void
  /**
   * PersistentPreRunE: PersistentPreRun but returns an error.
   */
  persistentPreRunE: (cmd: Command, args: Array<string>) => void
  /**
   * PreRun: children of this command will not inherit.
   */
  preRun: (cmd: Command, args: Array<string>) => void
  /**
   * PreRunE: PreRun but returns an error.
   */
  preRunE: (cmd: Command, args: Array<string>) => void
  /**
   * Run: Typically the actual work function. Most commands will only implement this.
   */
  run: (cmd: Command, args: Array<string>) => void
  /**
   * RunE: Run but returns an error.
   */
  runE: (cmd: Command, args: Array<string>) => void
  /**
   * PostRun: run after the Run command.
   */
  postRun: (cmd: Command, args: Array<string>) => void
  /**
   * PostRunE: PostRun but returns an error.
   */
  postRunE: (cmd: Command, args: Array<string>) => void
  /**
   * PersistentPostRun: children of this command will inherit and execute after PostRun.
   */
  persistentPostRun: (cmd: Command, args: Array<string>) => void
  /**
   * PersistentPostRunE: PersistentPostRun but returns an error.
   */
  persistentPostRunE: (cmd: Command, args: Array<string>) => void
  /**
   * FParseErrWhitelist flag parse errors to be ignored
   */
  fParseErrWhitelist: FParseErrWhitelist
  /**
   * CompletionOptions is a set of options to control the handling of shell completion
   */
  completionOptions: CompletionOptions
  /**
   * TraverseChildren parses flags on all parents before executing child command.
   */
  traverseChildren: boolean
  /**
   * Hidden defines, if this command is hidden and should NOT show up in the list of available commands.
   */
  hidden: boolean
  /**
   * SilenceErrors is an option to quiet errors down stream.
   */
  silenceErrors: boolean
  /**
   * SilenceUsage is an option to silence usage when an error occurs.
   */
  silenceUsage: boolean
  /**
   * DisableFlagParsing disables the flag parsing.
   * If this is true all flags will be passed to the command as arguments.
   */
  disableFlagParsing: boolean
  /**
   * DisableAutoGenTag defines, if gen tag ("Auto generated by spf13/cobra...")
   * will be printed by generating docs for this command.
   */
  disableAutoGenTag: boolean
  /**
   * DisableFlagsInUseLine will disable the addition of [flags] to the usage
   * line of a command when printing help or generating docs
   */
  disableFlagsInUseLine: boolean
  /**
   * DisableSuggestions disables the suggestions based on Levenshtein distance
   * that go along with 'unknown command' messages.
   */
  disableSuggestions: boolean
  /**
   * SuggestionsMinimumDistance defines minimum levenshtein distance to display suggestions.
   * Must be > 0.
   */
  suggestionsMinimumDistance: number
 }
 interface Command {
  /**
   * Context returns underlying command context. If command was executed
   * with ExecuteContext or the context was set with SetContext, the
   * previously set context will be returned. Otherwise, nil is returned.
   * 
   * Notice that a call to Execute and ExecuteC will replace a nil context of
   * a command with a context.Background, so a background context will be
   * returned by Context after one of these functions has been called.
   */
  context(): context.Context
 }
 interface Command {
  /**
   * SetContext sets context for the command. This context will be overwritten by
   * Command.ExecuteContext or Command.ExecuteContextC.
   */
  setContext(ctx: context.Context): void
 }
 interface Command {
  /**
   * SetArgs sets arguments for the command. It is set to os.Args[1:] by default, if desired, can be overridden
   * particularly useful when testing.
   */
  setArgs(a: Array<string>): void
 }
 interface Command {
  /**
   * SetOutput sets the destination for usage and error messages.
   * If output is nil, os.Stderr is used.
   * 
   * Deprecated: Use SetOut and/or SetErr instead
   */
  setOutput(output: io.Writer): void
 }
 interface Command {
  /**
   * SetOut sets the destination for usage messages.
   * If newOut is nil, os.Stdout is used.
   */
  setOut(newOut: io.Writer): void
 }
 interface Command {
  /**
   * SetErr sets the destination for error messages.
   * If newErr is nil, os.Stderr is used.
   */
  setErr(newErr: io.Writer): void
 }
 interface Command {
  /**
   * SetIn sets the source for input data
   * If newIn is nil, os.Stdin is used.
   */
  setIn(newIn: io.Reader): void
 }
 interface Command {
  /**
   * SetUsageFunc sets usage function. Usage can be defined by application.
   */
  setUsageFunc(f: (_arg0: Command) => void): void
 }
 interface Command {
  /**
   * SetUsageTemplate sets usage template. Can be defined by Application.
   */
  setUsageTemplate(s: string): void
 }
 interface Command {
  /**
   * SetFlagErrorFunc sets a function to generate an error when flag parsing
   * fails.
   */
  setFlagErrorFunc(f: (_arg0: Command, _arg1: Error) => void): void
 }
 interface Command {
  /**
   * SetHelpFunc sets help function. Can be defined by Application.
   */
  setHelpFunc(f: (_arg0: Command, _arg1: Array<string>) => void): void
 }
 interface Command {
  /**
   * SetHelpCommand sets help command.
   */
  setHelpCommand(cmd: Command): void
 }
 interface Command {
  /**
   * SetHelpCommandGroupID sets the group id of the help command.
   */
  setHelpCommandGroupID(groupID: string): void
 }
 interface Command {
  /**
   * SetCompletionCommandGroupID sets the group id of the completion command.
   */
  setCompletionCommandGroupID(groupID: string): void
 }
 interface Command {
  /**
   * SetHelpTemplate sets help template to be used. Application can use it to set custom template.
   */
  setHelpTemplate(s: string): void
 }
 interface Command {
  /**
   * SetVersionTemplate sets version template to be used. Application can use it to set custom template.
   */
  setVersionTemplate(s: string): void
 }
 interface Command {
  /**
   * SetErrPrefix sets error message prefix to be used. Application can use it to set custom prefix.
   */
  setErrPrefix(s: string): void
 }
 interface Command {
  /**
   * SetGlobalNormalizationFunc sets a normalization function to all flag sets and also to child commands.
   * The user should not have a cyclic dependency on commands.
   */
  setGlobalNormalizationFunc(n: (f: any, name: string) => any): void
 }
 interface Command {
  /**
   * OutOrStdout returns output to stdout.
   */
  outOrStdout(): io.Writer
 }
 interface Command {
  /**
   * OutOrStderr returns output to stderr
   */
  outOrStderr(): io.Writer
 }
 interface Command {
  /**
   * ErrOrStderr returns output to stderr
   */
  errOrStderr(): io.Writer
 }
 interface Command {
  /**
   * InOrStdin returns input to stdin
   */
  inOrStdin(): io.Reader
 }
 interface Command {
  /**
   * UsageFunc returns either the function set by SetUsageFunc for this command
   * or a parent, or it returns a default usage function.
   */
  usageFunc(): (_arg0: Command) => void
 }
 interface Command {
  /**
   * Usage puts out the usage for the command.
   * Used when a user provides invalid input.
   * Can be defined by user by overriding UsageFunc.
   */
  usage(): void
 }
 interface Command {
  /**
   * HelpFunc returns either the function set by SetHelpFunc for this command
   * or a parent, or it returns a function with default help behavior.
   */
  helpFunc(): (_arg0: Command, _arg1: Array<string>) => void
 }
 interface Command {
  /**
   * Help puts out the help for the command.
   * Used when a user calls help [command].
   * Can be defined by user by overriding HelpFunc.
   */
  help(): void
 }
 interface Command {
  /**
   * UsageString returns usage string.
   */
  usageString(): string
 }
 interface Command {
  /**
   * FlagErrorFunc returns either the function set by SetFlagErrorFunc for this
   * command or a parent, or it returns a function which returns the original
   * error.
   */
  flagErrorFunc(): (_arg0: Command, _arg1: Error) => void
 }
 interface Command {
  /**
   * UsagePadding return padding for the usage.
   */
  usagePadding(): number
 }
 interface Command {
  /**
   * CommandPathPadding return padding for the command path.
   */
  commandPathPadding(): number
 }
 interface Command {
  /**
   * NamePadding returns padding for the name.
   */
  namePadding(): number
 }
 interface Command {
  /**
   * UsageTemplate returns usage template for the command.
   * This function is kept for backwards-compatibility reasons.
   */
  usageTemplate(): string
 }
 interface Command {
  /**
   * HelpTemplate return help template for the command.
   * This function is kept for backwards-compatibility reasons.
   */
  helpTemplate(): string
 }
 interface Command {
  /**
   * VersionTemplate return version template for the command.
   * This function is kept for backwards-compatibility reasons.
   */
  versionTemplate(): string
 }
 interface Command {
  /**
   * ErrPrefix return error message prefix for the command
   */
  errPrefix(): string
 }
 interface Command {
  /**
   * Find the target command given the args and command tree
   * Meant to be run on the highest node. Only searches down.
   */
  find(args: Array<string>): [(Command), Array<string>]
 }
 interface Command {
  /**
   * Traverse the command tree to find the command, and parse args for
   * each parent.
   */
  traverse(args: Array<string>): [(Command), Array<string>]
 }
 interface Command {
  /**
   * SuggestionsFor provides suggestions for the typedName.
   */
  suggestionsFor(typedName: string): Array<string>
 }
 interface Command {
  /**
   * VisitParents visits all parents of the command and invokes fn on each parent.
   */
  visitParents(fn: (_arg0: Command) => void): void
 }
 interface Command {
  /**
   * Root finds root command.
   */
  root(): (Command)
 }
 interface Command {
  /**
   * ArgsLenAtDash will return the length of c.Flags().Args at the moment
   * when a -- was found during args parsing.
   */
  argsLenAtDash(): number
 }
 interface Command {
  /**
   * ExecuteContext is the same as Execute(), but sets the ctx on the command.
   * Retrieve ctx by calling cmd.Context() inside your *Run lifecycle or ValidArgs
   * functions.
   */
  executeContext(ctx: context.Context): void
 }
 interface Command {
  /**
   * Execute uses the args (os.Args[1:] by default)
   * and run through the command tree finding appropriate matches
   * for commands and then corresponding flags.
   */
  execute(): void
 }
 interface Command {
  /**
   * ExecuteContextC is the same as ExecuteC(), but sets the ctx on the command.
   * Retrieve ctx by calling cmd.Context() inside your *Run lifecycle or ValidArgs
   * functions.
   */
  executeContextC(ctx: context.Context): (Command)
 }
 interface Command {
  /**
   * ExecuteC executes the command.
   */
  executeC(): (Command)
 }
 interface Command {
  validateArgs(args: Array<string>): void
 }
 interface Command {
  /**
   * ValidateRequiredFlags validates all required flags are present and returns an error otherwise
   */
  validateRequiredFlags(): void
 }
 interface Command {
  /**
   * InitDefaultHelpFlag adds default help flag to c.
   * It is called automatically by executing the c or by calling help and usage.
   * If c already has help flag, it will do nothing.
   */
  initDefaultHelpFlag(): void
 }
 interface Command {
  /**
   * InitDefaultVersionFlag adds default version flag to c.
   * It is called automatically by executing the c.
   * If c already has a version flag, it will do nothing.
   * If c.Version is empty, it will do nothing.
   */
  initDefaultVersionFlag(): void
 }
 interface Command {
  /**
   * InitDefaultHelpCmd adds default help command to c.
   * It is called automatically by executing the c or by calling help and usage.
   * If c already has help command or c has no subcommands, it will do nothing.
   */
  initDefaultHelpCmd(): void
 }
 interface Command {
  /**
   * ResetCommands delete parent, subcommand and help command from c.
   */
  resetCommands(): void
 }
 interface Command {
  /**
   * Commands returns a sorted slice of child commands.
   */
  commands(): Array<(Command | undefined)>
 }
 interface Command {
  /**
   * AddCommand adds one or more commands to this parent command.
   */
  addCommand(...cmds: (Command | undefined)[]): void
 }
 interface Command {
  /**
   * Groups returns a slice of child command groups.
   */
  groups(): Array<(Group | undefined)>
 }
 interface Command {
  /**
   * AllChildCommandsHaveGroup returns if all subcommands are assigned to a group
   */
  allChildCommandsHaveGroup(): boolean
 }
 interface Command {
  /**
   * ContainsGroup return if groupID exists in the list of command groups.
   */
  containsGroup(groupID: string): boolean
 }
 interface Command {
  /**
   * AddGroup adds one or more command groups to this parent command.
   */
  addGroup(...groups: (Group | undefined)[]): void
 }
 interface Command {
  /**
   * RemoveCommand removes one or more commands from a parent command.
   */
  removeCommand(...cmds: (Command | undefined)[]): void
 }
 interface Command {
  /**
   * Print is a convenience method to Print to the defined output, fallback to Stderr if not set.
   */
  print(...i: {
   }[]): void
 }
 interface Command {
  /**
   * Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
   */
  println(...i: {
   }[]): void
 }
 interface Command {
  /**
   * Printf is a convenience method to Printf to the defined output, fallback to Stderr if not set.
   */
  printf(format: string, ...i: {
   }[]): void
 }
 interface Command {
  /**
   * PrintErr is a convenience method to Print to the defined Err output, fallback to Stderr if not set.
   */
  printErr(...i: {
   }[]): void
 }
 interface Command {
  /**
   * PrintErrln is a convenience method to Println to the defined Err output, fallback to Stderr if not set.
   */
  printErrln(...i: {
   }[]): void
 }
 interface Command {
  /**
   * PrintErrf is a convenience method to Printf to the defined Err output, fallback to Stderr if not set.
   */
  printErrf(format: string, ...i: {
   }[]): void
 }
 interface Command {
  /**
   * CommandPath returns the full path to this command.
   */
  commandPath(): string
 }
 interface Command {
  /**
   * DisplayName returns the name to display in help text. Returns command Name()
   * If CommandDisplayNameAnnoation is not set
   */
  displayName(): string
 }
 interface Command {
  /**
   * UseLine puts out the full usage for a given command (including parents).
   */
  useLine(): string
 }
 interface Command {
  /**
   * DebugFlags used to determine which flags have been assigned to which commands
   * and which persist.
   */
  debugFlags(): void
 }
 interface Command {
  /**
   * Name returns the command's name: the first word in the use line.
   */
  name(): string
 }
 interface Command {
  /**
   * HasAlias determines if a given string is an alias of the command.
   */
  hasAlias(s: string): boolean
 }
 interface Command {
  /**
   * CalledAs returns the command name or alias that was used to invoke
   * this command or an empty string if the command has not been called.
   */
  calledAs(): string
 }
 interface Command {
  /**
   * NameAndAliases returns a list of the command name and all aliases
   */
  nameAndAliases(): string
 }
 interface Command {
  /**
   * HasExample determines if the command has example.
   */
  hasExample(): boolean
 }
 interface Command {
  /**
   * Runnable determines if the command is itself runnable.
   */
  runnable(): boolean
 }
 interface Command {
  /**
   * HasSubCommands determines if the command has children commands.
   */
  hasSubCommands(): boolean
 }
 interface Command {
  /**
   * IsAvailableCommand determines if a command is available as a non-help command
   * (this includes all non deprecated/hidden commands).
   */
  isAvailableCommand(): boolean
 }
 interface Command {
  /**
   * IsAdditionalHelpTopicCommand determines if a command is an additional
   * help topic command; additional help topic command is determined by the
   * fact that it is NOT runnable/hidden/deprecated, and has no sub commands that
   * are runnable/hidden/deprecated.
   * Concrete example: https://github.com/spf13/cobra/issues/393#issuecomment-282741924.
   */
  isAdditionalHelpTopicCommand(): boolean
 }
 interface Command {
  /**
   * HasHelpSubCommands determines if a command has any available 'help' sub commands
   * that need to be shown in the usage/help default template under 'additional help
   * topics'.
   */
  hasHelpSubCommands(): boolean
 }
 interface Command {
  /**
   * HasAvailableSubCommands determines if a command has available sub commands that
   * need to be shown in the usage/help default template under 'available commands'.
   */
  hasAvailableSubCommands(): boolean
 }
 interface Command {
  /**
   * HasParent determines if the command is a child command.
   */
  hasParent(): boolean
 }
 interface Command {
  /**
   * GlobalNormalizationFunc returns the global normalization function or nil if it doesn't exist.
   */
  globalNormalizationFunc(): (f: any, name: string) => any
 }
 interface Command {
  /**
   * Flags returns the complete FlagSet that applies
   * to this command (local and persistent declared here and by all parents).
   */
  flags(): (any)
 }
 interface Command {
  /**
   * LocalNonPersistentFlags are flags specific to this command which will NOT persist to subcommands.
   * This function does not modify the flags of the current command, it's purpose is to return the current state.
   */
  localNonPersistentFlags(): (any)
 }
 interface Command {
  /**
   * LocalFlags returns the local FlagSet specifically set in the current command.
   * This function does not modify the flags of the current command, it's purpose is to return the current state.
   */
  localFlags(): (any)
 }
 interface Command {
  /**
   * InheritedFlags returns all flags which were inherited from parent commands.
   * This function does not modify the flags of the current command, it's purpose is to return the current state.
   */
  inheritedFlags(): (any)
 }
 interface Command {
  /**
   * NonInheritedFlags returns all flags which were not inherited from parent commands.
   * This function does not modify the flags of the current command, it's purpose is to return the current state.
   */
  nonInheritedFlags(): (any)
 }
 interface Command {
  /**
   * PersistentFlags returns the persistent FlagSet specifically set in the current command.
   */
  persistentFlags(): (any)
 }
 interface Command {
  /**
   * ResetFlags deletes all flags from command.
   */
  resetFlags(): void
 }
 interface Command {
  /**
   * HasFlags checks if the command contains any flags (local plus persistent from the entire structure).
   */
  hasFlags(): boolean
 }
 interface Command {
  /**
   * HasPersistentFlags checks if the command contains persistent flags.
   */
  hasPersistentFlags(): boolean
 }
 interface Command {
  /**
   * HasLocalFlags checks if the command has flags specifically declared locally.
   */
  hasLocalFlags(): boolean
 }
 interface Command {
  /**
   * HasInheritedFlags checks if the command has flags inherited from its parent command.
   */
  hasInheritedFlags(): boolean
 }
 interface Command {
  /**
   * HasAvailableFlags checks if the command contains any flags (local plus persistent from the entire
   * structure) which are not hidden or deprecated.
   */
  hasAvailableFlags(): boolean
 }
 interface Command {
  /**
   * HasAvailablePersistentFlags checks if the command contains persistent flags which are not hidden or deprecated.
   */
  hasAvailablePersistentFlags(): boolean
 }
 interface Command {
  /**
   * HasAvailableLocalFlags checks if the command has flags specifically declared locally which are not hidden
   * or deprecated.
   */
  hasAvailableLocalFlags(): boolean
 }
 interface Command {
  /**
   * HasAvailableInheritedFlags checks if the command has flags inherited from its parent command which are
   * not hidden or deprecated.
   */
  hasAvailableInheritedFlags(): boolean
 }
 interface Command {
  /**
   * Flag climbs up the command tree looking for matching flag.
   */
  flag(name: string): (any)
 }
 interface Command {
  /**
   * ParseFlags parses persistent flag tree and local flags.
   */
  parseFlags(args: Array<string>): void
 }
 interface Command {
  /**
   * Parent returns a commands parent command.
   */
  parent(): (Command)
 }
 interface Command {
  /**
   * RegisterFlagCompletionFunc should be called to register a function to provide completion for a flag.
   * 
   * You can use pre-defined completion functions such as [FixedCompletions] or [NoFileCompletions],
   * or you can define your own.
   */
  registerFlagCompletionFunc(flagName: string, f: CompletionFunc): void
 }
 interface Command {
  /**
   * GetFlagCompletionFunc returns the completion function for the given flag of the command, if available.
   */
  getFlagCompletionFunc(flagName: string): [CompletionFunc, boolean]
 }
 interface Command {
  /**
   * InitDefaultCompletionCmd adds a default 'completion' command to c.
   * This function will do nothing if any of the following is true:
   * 1- the feature has been explicitly disabled by the program,
   * 2- c has no subcommands (to avoid creating one),
   * 3- c already has a 'completion' command provided by the program.
   */
  initDefaultCompletionCmd(...args: string[]): void
 }
 interface Command {
  /**
   * GenFishCompletion generates fish completion file and writes to the passed writer.
   */
  genFishCompletion(w: io.Writer, includeDesc: boolean): void
 }
 interface Command {
  /**
   * GenFishCompletionFile generates fish completion file.
   */
  genFishCompletionFile(filename: string, includeDesc: boolean): void
 }
 interface Command {
  /**
   * MarkFlagsRequiredTogether marks the given flags with annotations so that Cobra errors
   * if the command is invoked with a subset (but not all) of the given flags.
   */
  markFlagsRequiredTogether(...flagNames: string[]): void
 }
 interface Command {
  /**
   * MarkFlagsOneRequired marks the given flags with annotations so that Cobra errors
   * if the command is invoked without at least one flag from the given set of flags.
   */
  markFlagsOneRequired(...flagNames: string[]): void
 }
 interface Command {
  /**
   * MarkFlagsMutuallyExclusive marks the given flags with annotations so that Cobra errors
   * if the command is invoked with more than one flag from the given set of flags.
   */
  markFlagsMutuallyExclusive(...flagNames: string[]): void
 }
 interface Command {
  /**
   * ValidateFlagGroups validates the mutuallyExclusive/oneRequired/requiredAsGroup logic and returns the
   * first error encountered.
   */
  validateFlagGroups(): void
 }
 interface Command {
  /**
   * GenPowerShellCompletionFile generates powershell completion file without descriptions.
   */
  genPowerShellCompletionFile(filename: string): void
 }
 interface Command {
  /**
   * GenPowerShellCompletion generates powershell completion file without descriptions
   * and writes it to the passed writer.
   */
  genPowerShellCompletion(w: io.Writer): void
 }
 interface Command {
  /**
   * GenPowerShellCompletionFileWithDesc generates powershell completion file with descriptions.
   */
  genPowerShellCompletionFileWithDesc(filename: string): void
 }
 interface Command {
  /**
   * GenPowerShellCompletionWithDesc generates powershell completion file with descriptions
   * and writes it to the passed writer.
   */
  genPowerShellCompletionWithDesc(w: io.Writer): void
 }
 interface Command {
  /**
   * MarkFlagRequired instructs the various shell completion implementations to
   * prioritize the named flag when performing completion,
   * and causes your command to report an error if invoked without the flag.
   */
  markFlagRequired(name: string): void
 }
 interface Command {
  /**
   * MarkPersistentFlagRequired instructs the various shell completion implementations to
   * prioritize the named persistent flag when performing completion,
   * and causes your command to report an error if invoked without the flag.
   */
  markPersistentFlagRequired(name: string): void
 }
 interface Command {
  /**
   * MarkFlagFilename instructs the various shell completion implementations to
   * limit completions for the named flag to the specified file extensions.
   */
  markFlagFilename(name: string, ...extensions: string[]): void
 }
 interface Command {
  /**
   * MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
   * The bash completion script will call the bash function f for the flag.
   * 
   * This will only work for bash completion.
   * It is recommended to instead use c.RegisterFlagCompletionFunc(...) which allows
   * to register a Go function which will work across all shells.
   */
  markFlagCustom(name: string, f: string): void
 }
 interface Command {
  /**
   * MarkPersistentFlagFilename instructs the various shell completion
   * implementations to limit completions for the named persistent flag to the
   * specified file extensions.
   */
  markPersistentFlagFilename(name: string, ...extensions: string[]): void
 }
 interface Command {
  /**
   * MarkFlagDirname instructs the various shell completion implementations to
   * limit completions for the named flag to directory names.
   */
  markFlagDirname(name: string): void
 }
 interface Command {
  /**
   * MarkPersistentFlagDirname instructs the various shell completion
   * implementations to limit completions for the named persistent flag to
   * directory names.
   */
  markPersistentFlagDirname(name: string): void
 }
 interface Command {
  /**
   * GenZshCompletionFile generates zsh completion file including descriptions.
   */
  genZshCompletionFile(filename: string): void
 }
 interface Command {
  /**
   * GenZshCompletion generates zsh completion file including descriptions
   * and writes it to the passed writer.
   */
  genZshCompletion(w: io.Writer): void
 }
 interface Command {
  /**
   * GenZshCompletionFileNoDesc generates zsh completion file without descriptions.
   */
  genZshCompletionFileNoDesc(filename: string): void
 }
 interface Command {
  /**
   * GenZshCompletionNoDesc generates zsh completion file without descriptions
   * and writes it to the passed writer.
   */
  genZshCompletionNoDesc(w: io.Writer): void
 }
 interface Command {
  /**
   * MarkZshCompPositionalArgumentFile only worked for zsh and its behavior was
   * not consistent with Bash completion. It has therefore been disabled.
   * Instead, when no other completion is specified, file completion is done by
   * default for every argument. One can disable file completion on a per-argument
   * basis by using ValidArgsFunction and ShellCompDirectiveNoFileComp.
   * To achieve file extension filtering, one can use ValidArgsFunction and
   * ShellCompDirectiveFilterFileExt.
   * 
   * Deprecated
   */
  markZshCompPositionalArgumentFile(argPosition: number, ...patterns: string[]): void
 }
 interface Command {
  /**
   * MarkZshCompPositionalArgumentWords only worked for zsh. It has therefore
   * been disabled.
   * To achieve the same behavior across all shells, one can use
   * ValidArgs (for the first argument only) or ValidArgsFunction for
   * any argument (can include the first one also).
   * 
   * Deprecated
   */
  markZshCompPositionalArgumentWords(argPosition: number, ...words: string[]): void
 }
}

namespace sync {
 // @ts-ignore
 import isync = sync
 /**
  * A Locker represents an object that can be locked and unlocked.
  */
 interface Locker {
  [key:string]: any;
  lock(): void
  unlock(): void
 }
}

namespace io {
 /**
  * WriteCloser is the interface that groups the basic Write and Close methods.
  */
 interface WriteCloser {
  [key:string]: any;
 }
}

namespace syscall {
 // @ts-ignore
 import errpkg = errors
 /**
  * SysProcIDMap holds Container ID to Host ID mappings used for User Namespaces in Linux.
  * See user_namespaces(7).
  * 
  * Note that User Namespaces are not available on a number of popular Linux
  * versions (due to security issues), or are available but subject to AppArmor
  * restrictions like in Ubuntu 24.04.
  */
 interface SysProcIDMap {
  containerID: number // Container ID.
  hostID: number // Host ID.
  size: number // Size.
 }
 // @ts-ignore
 import errorspkg = errors
 /**
  * Credential holds user and group identities to be assumed
  * by a child process started by [StartProcess].
  */
 interface Credential {
  uid: number // User ID.
  gid: number // Group ID.
  groups: Array<number> // Supplementary group IDs.
  noSetGroups: boolean // If true, don't set supplementary groups
 }
 // @ts-ignore
 import runtimesyscall = syscall
 /**
  * A Signal is a number describing a process signal.
  * It implements the [os.Signal] interface.
  */
 interface Signal extends Number{}
 interface Signal {
  signal(): void
 }
 interface Signal {
  string(): string
 }
}

namespace time {
 /**
  * A Month specifies a month of the year (January = 1, ...).
  */
 interface Month extends Number{}
 interface Month {
  /**
   * String returns the English name of the month ("January", "February", ...).
   */
  string(): string
 }
 /**
  * A Weekday specifies a day of the week (Sunday = 0, ...).
  */
 interface Weekday extends Number{}
 interface Weekday {
  /**
   * String returns the English name of the day ("Sunday", "Monday", ...).
   */
  string(): string
 }
 /**
  * A Location maps time instants to the zone in use at that time.
  * Typically, the Location represents the collection of time offsets
  * in use in a geographical area. For many Locations the time offset varies
  * depending on whether daylight savings time is in use at the time instant.
  * 
  * Location is used to provide a time zone in a printed Time value and for
  * calculations involving intervals that may cross daylight savings time
  * boundaries.
  */
 interface Location {
 }
 interface Location {
  /**
   * String returns a descriptive name for the time zone information,
   * corresponding to the name argument to [LoadLocation] or [FixedZone].
   */
  string(): string
 }
}

namespace fs {
}

namespace store {
}

/**
 * Package url parses URLs and implements query escaping.
 */
namespace url {
 /**
  * A URL represents a parsed URL (technically, a URI reference).
  * 
  * The general form represented is:
  * 
  * ```
  * 	[scheme:][//[userinfo@]host][/]path[?query][#fragment]
  * ```
  * 
  * URLs that do not start with a slash after the scheme are interpreted as:
  * 
  * ```
  * 	scheme:opaque[?query][#fragment]
  * ```
  * 
  * The Host field contains the host and port subcomponents of the URL.
  * When the port is present, it is separated from the host with a colon.
  * When the host is an IPv6 address, it must be enclosed in square brackets:
  * "[fe80::1]:80". The [net.JoinHostPort] function combines a host and port
  * into a string suitable for the Host field, adding square brackets to
  * the host when necessary.
  * 
  * Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/.
  * A consequence is that it is impossible to tell which slashes in the Path were
  * slashes in the raw URL and which were %2f. This distinction is rarely important,
  * but when it is, the code should use the [URL.EscapedPath] method, which preserves
  * the original encoding of Path.
  * 
  * The RawPath field is an optional field which is only set when the default
  * encoding of Path is different from the escaped path. See the EscapedPath method
  * for more details.
  * 
  * URL's String method uses the EscapedPath method to obtain the path.
  */
 interface URL {
  scheme: string
  opaque: string // encoded opaque data
  user?: Userinfo // username and password information
  host: string // host or host:port (see Hostname and Port methods)
  path: string // path (relative paths may omit leading slash)
  rawPath: string // encoded path hint (see EscapedPath method)
  omitHost: boolean // do not emit empty host (authority)
  forceQuery: boolean // append a query ('?') even if RawQuery is empty
  rawQuery: string // encoded query values, without '?'
  fragment: string // fragment for references, without '#'
  rawFragment: string // encoded fragment hint (see EscapedFragment method)
 }
 interface URL {
  /**
   * EscapedPath returns the escaped form of u.Path.
   * In general there are multiple possible escaped forms of any path.
   * EscapedPath returns u.RawPath when it is a valid escaping of u.Path.
   * Otherwise EscapedPath ignores u.RawPath and computes an escaped
   * form on its own.
   * The [URL.String] and [URL.RequestURI] methods use EscapedPath to construct
   * their results.
   * In general, code should call EscapedPath instead of
   * reading u.RawPath directly.
   */
  escapedPath(): string
 }
 interface URL {
  /**
   * EscapedFragment returns the escaped form of u.Fragment.
   * In general there are multiple possible escaped forms of any fragment.
   * EscapedFragment returns u.RawFragment when it is a valid escaping of u.Fragment.
   * Otherwise EscapedFragment ignores u.RawFragment and computes an escaped
   * form on its own.
   * The [URL.String] method uses EscapedFragment to construct its result.
   * In general, code should call EscapedFragment instead of
   * reading u.RawFragment directly.
   */
  escapedFragment(): string
 }
 interface URL {
  /**
   * String reassembles the [URL] into a valid URL string.
   * The general form of the result is one of:
   * 
   * ```
   * 	scheme:opaque?query#fragment
   * 	scheme://userinfo@host/path?query#fragment
   * ```
   * 
   * If u.Opaque is non-empty, String uses the first form;
   * otherwise it uses the second form.
   * Any non-ASCII characters in host are escaped.
   * To obtain the path, String uses u.EscapedPath().
   * 
   * In the second form, the following rules apply:
   * ```
   *   - if u.Scheme is empty, scheme: is omitted.
   *   - if u.User is nil, userinfo@ is omitted.
   *   - if u.Host is empty, host/ is omitted.
   *   - if u.Scheme and u.Host are empty and u.User is nil,
   *     the entire scheme://userinfo@host/ is omitted.
   *   - if u.Host is non-empty and u.Path begins with a /,
   *     the form host/path does not add its own /.
   *   - if u.RawQuery is empty, ?query is omitted.
   *   - if u.Fragment is empty, #fragment is omitted.
   * ```
   */
  string(): string
 }
 interface URL {
  /**
   * Redacted is like [URL.String] but replaces any password with "xxxxx".
   * Only the password in u.User is redacted.
   */
  redacted(): string
 }
 /**
  * Values maps a string key to a list of values.
  * It is typically used for query parameters and form values.
  * Unlike in the http.Header map, the keys in a Values map
  * are case-sensitive.
  */
 interface Values extends _TygojaDict{}
 interface Values {
  /**
   * Get gets the first value associated with the given key.
   * If there are no values associated with the key, Get returns
   * the empty string. To access multiple values, use the map
   * directly.
   */
  get(key: string): string
 }
 interface Values {
  /**
   * Set sets the key to value. It replaces any existing
   * values.
   */
  set(key: string, value: string): void
 }
 interface Values {
  /**
   * Add adds the value to key. It appends to any existing
   * values associated with key.
   */
  add(key: string, value: string): void
 }
 interface Values {
  /**
   * Del deletes the values associated with key.
   */
  del(key: string): void
 }
 interface Values {
  /**
   * Has checks whether a given key is set.
   */
  has(key: string): boolean
 }
 interface Values {
  /**
   * Encode encodes the values into “URL encoded” form
   * ("bar=baz&foo=quux") sorted by key.
   */
  encode(): string
 }
 interface URL {
  /**
   * IsAbs reports whether the [URL] is absolute.
   * Absolute means that it has a non-empty scheme.
   */
  isAbs(): boolean
 }
 interface URL {
  /**
   * Parse parses a [URL] in the context of the receiver. The provided URL
   * may be relative or absolute. Parse returns nil, err on parse
   * failure, otherwise its return value is the same as [URL.ResolveReference].
   */
  parse(ref: string): (URL)
 }
 interface URL {
  /**
   * ResolveReference resolves a URI reference to an absolute URI from
   * an absolute base URI u, per RFC 3986 Section 5.2. The URI reference
   * may be relative or absolute. ResolveReference always returns a new
   * [URL] instance, even if the returned URL is identical to either the
   * base or reference. If ref is an absolute URL, then ResolveReference
   * ignores base and returns a copy of ref.
   */
  resolveReference(ref: URL): (URL)
 }
 interface URL {
  /**
   * Query parses RawQuery and returns the corresponding values.
   * It silently discards malformed value pairs.
   * To check errors use [ParseQuery].
   */
  query(): Values
 }
 interface URL {
  /**
   * RequestURI returns the encoded path?query or opaque?query
   * string that would be used in an HTTP request for u.
   */
  requestURI(): string
 }
 interface URL {
  /**
   * Hostname returns u.Host, stripping any valid port number if present.
   * 
   * If the result is enclosed in square brackets, as literal IPv6 addresses are,
   * the square brackets are removed from the result.
   */
  hostname(): string
 }
 interface URL {
  /**
   * Port returns the port part of u.Host, without the leading colon.
   * 
   * If u.Host doesn't contain a valid numeric port, Port returns an empty string.
   */
  port(): string
 }
 interface URL {
  marshalBinary(): string|Array<number>
 }
 interface URL {
  appendBinary(b: string|Array<number>): string|Array<number>
 }
 interface URL {
  unmarshalBinary(text: string|Array<number>): void
 }
 interface URL {
  /**
   * JoinPath returns a new [URL] with the provided path elements joined to
   * any existing path and the resulting path cleaned of any ./ or ../ elements.
   * Any sequences of multiple / characters will be reduced to a single /.
   */
  joinPath(...elem: string[]): (URL)
 }
}

namespace context {
}

namespace net {
 /**
  * Addr represents a network end point address.
  * 
  * The two methods [Addr.Network] and [Addr.String] conventionally return strings
  * that can be passed as the arguments to [Dial], but the exact form
  * and meaning of the strings is up to the implementation.
  */
 interface Addr {
  [key:string]: any;
  network(): string // name of the network (for example, "tcp", "udp")
  string(): string // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
 }
}

namespace jwt {
 /**
  * NumericDate represents a JSON numeric date value, as referenced at
  * https://datatracker.ietf.org/doc/html/rfc7519#section-2.
  */
 type _sPgnyHi = time.Time
 interface NumericDate extends _sPgnyHi {
 }
 interface NumericDate {
  /**
   * MarshalJSON is an implementation of the json.RawMessage interface and serializes the UNIX epoch
   * represented in NumericDate to a byte array, using the precision specified in TimePrecision.
   */
  marshalJSON(): string|Array<number>
 }
 interface NumericDate {
  /**
   * UnmarshalJSON is an implementation of the json.RawMessage interface and
   * deserializes a [NumericDate] from a JSON representation, i.e. a
   * [json.Number]. This number represents an UNIX epoch with either integer or
   * non-integer seconds.
   */
  unmarshalJSON(b: string|Array<number>): void
 }
 /**
  * ClaimStrings is basically just a slice of strings, but it can be either
  * serialized from a string array or just a string. This type is necessary,
  * since the "aud" claim can either be a single string or an array.
  */
 interface ClaimStrings extends Array<string>{}
 interface ClaimStrings {
  unmarshalJSON(data: string|Array<number>): void
 }
 interface ClaimStrings {
  marshalJSON(): string|Array<number>
 }
}

namespace hook {
 /**
  * wrapped local Hook embedded struct to limit the public API surface.
  */
 type _spkyOMF<T> = Hook<T>
 interface mainHook<T> extends _spkyOMF<T> {
 }
}

namespace bufio {
 /**
  * Reader implements buffering for an io.Reader object.
  * A new Reader is created by calling [NewReader] or [NewReaderSize];
  * alternatively the zero value of a Reader may be used after calling [Reset]
  * on it.
  */
 interface Reader {
 }
 interface Reader {
  /**
   * Size returns the size of the underlying buffer in bytes.
   */
  size(): number
 }
 interface Reader {
  /**
   * Reset discards any buffered data, resets all state, and switches
   * the buffered reader to read from r.
   * Calling Reset on the zero value of [Reader] initializes the internal buffer
   * to the default size.
   * Calling b.Reset(b) (that is, resetting a [Reader] to itself) does nothing.
   */
  reset(r: io.Reader): void
 }
 interface Reader {
  /**
   * Peek returns the next n bytes without advancing the reader. The bytes stop
   * being valid at the next read call. If necessary, Peek will read more bytes
   * into the buffer in order to make n bytes available. If Peek returns fewer
   * than n bytes, it also returns an error explaining why the read is short.
   * The error is [ErrBufferFull] if n is larger than b's buffer size.
   * 
   * Calling Peek prevents a [Reader.UnreadByte] or [Reader.UnreadRune] call from succeeding
   * until the next read operation.
   */
  peek(n: number): string|Array<number>
 }
 interface Reader {
  /**
   * Discard skips the next n bytes, returning the number of bytes discarded.
   * 
   * If Discard skips fewer than n bytes, it also returns an error.
   * If 0 <= n <= b.Buffered(), Discard is guaranteed to succeed without
   * reading from the underlying io.Reader.
   */
  discard(n: number): number
 }
 interface Reader {
  /**
   * Read reads data into p.
   * It returns the number of bytes read into p.
   * The bytes are taken from at most one Read on the underlying [Reader],
   * hence n may be less than len(p).
   * To read exactly len(p) bytes, use io.ReadFull(b, p).
   * If the underlying [Reader] can return a non-zero count with io.EOF,
   * then this Read method can do so as well; see the [io.Reader] docs.
   */
  read(p: string|Array<number>): number
 }
 interface Reader {
  /**
   * ReadByte reads and returns a single byte.
   * If no byte is available, returns an error.
   */
  readByte(): number
 }
 interface Reader {
  /**
   * UnreadByte unreads the last byte. Only the most recently read byte can be unread.
   * 
   * UnreadByte returns an error if the most recent method called on the
   * [Reader] was not a read operation. Notably, [Reader.Peek], [Reader.Discard], and [Reader.WriteTo] are not
   * considered read operations.
   */
  unreadByte(): void
 }
 interface Reader {
  /**
   * ReadRune reads a single UTF-8 encoded Unicode character and returns the
   * rune and its size in bytes. If the encoded rune is invalid, it consumes one byte
   * and returns unicode.ReplacementChar (U+FFFD) with a size of 1.
   */
  readRune(): [number, number]
 }
 interface Reader {
  /**
   * UnreadRune unreads the last rune. If the most recent method called on
   * the [Reader] was not a [Reader.ReadRune], [Reader.UnreadRune] returns an error. (In this
   * regard it is stricter than [Reader.UnreadByte], which will unread the last byte
   * from any read operation.)
   */
  unreadRune(): void
 }
 interface Reader {
  /**
   * Buffered returns the number of bytes that can be read from the current buffer.
   */
  buffered(): number
 }
 interface Reader {
  /**
   * ReadSlice reads until the first occurrence of delim in the input,
   * returning a slice pointing at the bytes in the buffer.
   * The bytes stop being valid at the next read.
   * If ReadSlice encounters an error before finding a delimiter,
   * it returns all the data in the buffer and the error itself (often io.EOF).
   * ReadSlice fails with error [ErrBufferFull] if the buffer fills without a delim.
   * Because the data returned from ReadSlice will be overwritten
   * by the next I/O operation, most clients should use
   * [Reader.ReadBytes] or ReadString instead.
   * ReadSlice returns err != nil if and only if line does not end in delim.
   */
  readSlice(delim: number): string|Array<number>
 }
 interface Reader {
  /**
   * ReadLine is a low-level line-reading primitive. Most callers should use
   * [Reader.ReadBytes]('\n') or [Reader.ReadString]('\n') instead or use a [Scanner].
   * 
   * ReadLine tries to return a single line, not including the end-of-line bytes.
   * If the line was too long for the buffer then isPrefix is set and the
   * beginning of the line is returned. The rest of the line will be returned
   * from future calls. isPrefix will be false when returning the last fragment
   * of the line. The returned buffer is only valid until the next call to
   * ReadLine. ReadLine either returns a non-nil line or it returns an error,
   * never both.
   * 
   * The text returned from ReadLine does not include the line end ("\r\n" or "\n").
   * No indication or error is given if the input ends without a final line end.
   * Calling [Reader.UnreadByte] after ReadLine will always unread the last byte read
   * (possibly a character belonging to the line end) even if that byte is not
   * part of the line returned by ReadLine.
   */
  readLine(): [string|Array<number>, boolean]
 }
 interface Reader {
  /**
   * ReadBytes reads until the first occurrence of delim in the input,
   * returning a slice containing the data up to and including the delimiter.
   * If ReadBytes encounters an error before finding a delimiter,
   * it returns the data read before the error and the error itself (often io.EOF).
   * ReadBytes returns err != nil if and only if the returned data does not end in
   * delim.
   * For simple uses, a Scanner may be more convenient.
   */
  readBytes(delim: number): string|Array<number>
 }
 interface Reader {
  /**
   * ReadString reads until the first occurrence of delim in the input,
   * returning a string containing the data up to and including the delimiter.
   * If ReadString encounters an error before finding a delimiter,
   * it returns the data read before the error and the error itself (often io.EOF).
   * ReadString returns err != nil if and only if the returned data does not end in
   * delim.
   * For simple uses, a Scanner may be more convenient.
   */
  readString(delim: number): string
 }
 interface Reader {
  /**
   * WriteTo implements io.WriterTo.
   * This may make multiple calls to the [Reader.Read] method of the underlying [Reader].
   * If the underlying reader supports the [Reader.WriteTo] method,
   * this calls the underlying [Reader.WriteTo] without buffering.
   */
  writeTo(w: io.Writer): number
 }
 /**
  * Writer implements buffering for an [io.Writer] object.
  * If an error occurs writing to a [Writer], no more data will be
  * accepted and all subsequent writes, and [Writer.Flush], will return the error.
  * After all data has been written, the client should call the
  * [Writer.Flush] method to guarantee all data has been forwarded to
  * the underlying [io.Writer].
  */
 interface Writer {
 }
 interface Writer {
  /**
   * Size returns the size of the underlying buffer in bytes.
   */
  size(): number
 }
 interface Writer {
  /**
   * Reset discards any unflushed buffered data, clears any error, and
   * resets b to write its output to w.
   * Calling Reset on the zero value of [Writer] initializes the internal buffer
   * to the default size.
   * Calling w.Reset(w) (that is, resetting a [Writer] to itself) does nothing.
   */
  reset(w: io.Writer): void
 }
 interface Writer {
  /**
   * Flush writes any buffered data to the underlying [io.Writer].
   */
  flush(): void
 }
 interface Writer {
  /**
   * Available returns how many bytes are unused in the buffer.
   */
  available(): number
 }
 interface Writer {
  /**
   * AvailableBuffer returns an empty buffer with b.Available() capacity.
   * This buffer is intended to be appended to and
   * passed to an immediately succeeding [Writer.Write] call.
   * The buffer is only valid until the next write operation on b.
   */
  availableBuffer(): string|Array<number>
 }
 interface Writer {
  /**
   * Buffered returns the number of bytes that have been written into the current buffer.
   */
  buffered(): number
 }
 interface Writer {
  /**
   * Write writes the contents of p into the buffer.
   * It returns the number of bytes written.
   * If nn < len(p), it also returns an error explaining
   * why the write is short.
   */
  write(p: string|Array<number>): number
 }
 interface Writer {
  /**
   * WriteByte writes a single byte.
   */
  writeByte(c: number): void
 }
 interface Writer {
  /**
   * WriteRune writes a single Unicode code point, returning
   * the number of bytes written and any error.
   */
  writeRune(r: number): number
 }
 interface Writer {
  /**
   * WriteString writes a string.
   * It returns the number of bytes written.
   * If the count is less than len(s), it also returns an error explaining
   * why the write is short.
   */
  writeString(s: string): number
 }
 interface Writer {
  /**
   * ReadFrom implements [io.ReaderFrom]. If the underlying writer
   * supports the ReadFrom method, this calls the underlying ReadFrom.
   * If there is buffered data and an underlying ReadFrom, this fills
   * the buffer and writes it before calling ReadFrom.
   */
  readFrom(r: io.Reader): number
 }
}

namespace cron {
 /**
  * Job defines a single registered cron job.
  */
 interface Job {
 }
 interface Job {
  /**
   * Id returns the cron job id.
   */
  id(): string
 }
 interface Job {
  /**
   * Expression returns the plain cron job schedule expression.
   */
  expression(): string
 }
 interface Job {
  /**
   * Run runs the cron job function.
   */
  run(): void
 }
 interface Job {
  /**
   * MarshalJSON implements [json.Marshaler] and export the current
   * jobs data into valid JSON.
   */
  marshalJSON(): string|Array<number>
 }
}

namespace sql {
 /**
  * IsolationLevel is the transaction isolation level used in [TxOptions].
  */
 interface IsolationLevel extends Number{}
 interface IsolationLevel {
  /**
   * String returns the name of the transaction isolation level.
   */
  string(): string
 }
 /**
  * DBStats contains database statistics.
  */
 interface DBStats {
  maxOpenConnections: number // Maximum number of open connections to the database.
  /**
   * Pool Status
   */
  openConnections: number // The number of established connections both in use and idle.
  inUse: number // The number of connections currently in use.
  idle: number // The number of idle connections.
  /**
   * Counters
   */
  waitCount: number // The total number of connections waited for.
  waitDuration: time.Duration // The total time blocked waiting for a new connection.
  maxIdleClosed: number // The total number of connections closed due to SetMaxIdleConns.
  maxIdleTimeClosed: number // The total number of connections closed due to SetConnMaxIdleTime.
  maxLifetimeClosed: number // The total number of connections closed due to SetConnMaxLifetime.
 }
 /**
  * Conn represents a single database connection rather than a pool of database
  * connections. Prefer running queries from [DB] unless there is a specific
  * need for a continuous single database connection.
  * 
  * A Conn must call [Conn.Close] to return the connection to the database pool
  * and may do so concurrently with a running query.
  * 
  * After a call to [Conn.Close], all operations on the
  * connection fail with [ErrConnDone].
  */
 interface Conn {
 }
 interface Conn {
  /**
   * PingContext verifies the connection to the database is still alive.
   */
  pingContext(ctx: context.Context): void
 }
 interface Conn {
  /**
   * ExecContext executes a query without returning any rows.
   * The args are for any placeholder parameters in the query.
   */
  execContext(ctx: context.Context, query: string, ...args: any[]): Result
 }
 interface Conn {
  /**
   * QueryContext executes a query that returns rows, typically a SELECT.
   * The args are for any placeholder parameters in the query.
   */
  queryContext(ctx: context.Context, query: string, ...args: any[]): (Rows)
 }
 interface Conn {
  /**
   * QueryRowContext executes a query that is expected to return at most one row.
   * QueryRowContext always returns a non-nil value. Errors are deferred until
   * the [*Row.Scan] method is called.
   * If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
   * Otherwise, the [*Row.Scan] scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, query: string, ...args: any[]): (Row)
 }
 interface Conn {
  /**
   * PrepareContext creates a prepared statement for later queries or executions.
   * Multiple queries or executions may be run concurrently from the
   * returned statement.
   * The caller must call the statement's [*Stmt.Close] method
   * when the statement is no longer needed.
   * 
   * The provided context is used for the preparation of the statement, not for the
   * execution of the statement.
   */
  prepareContext(ctx: context.Context, query: string): (Stmt)
 }
 interface Conn {
  /**
   * Raw executes f exposing the underlying driver connection for the
   * duration of f. The driverConn must not be used outside of f.
   * 
   * Once f returns and err is not [driver.ErrBadConn], the [Conn] will continue to be usable
   * until [Conn.Close] is called.
   */
  raw(f: (driverConn: any) => void): void
 }
 interface Conn {
  /**
   * BeginTx starts a transaction.
   * 
   * The provided context is used until the transaction is committed or rolled back.
   * If the context is canceled, the sql package will roll back
   * the transaction. [Tx.Commit] will return an error if the context provided to
   * BeginTx is canceled.
   * 
   * The provided [TxOptions] is optional and may be nil if defaults should be used.
   * If a non-default isolation level is used that the driver doesn't support,
   * an error will be returned.
   */
  beginTx(ctx: context.Context, opts: TxOptions): (Tx)
 }
 interface Conn {
  /**
   * Close returns the connection to the connection pool.
   * All operations after a Close will return with [ErrConnDone].
   * Close is safe to call concurrently with other operations and will
   * block until all other operations finish. It may be useful to first
   * cancel any used context and then call close directly after.
   */
  close(): void
 }
 /**
  * ColumnType contains the name and type of a column.
  */
 interface ColumnType {
 }
 interface ColumnType {
  /**
   * Name returns the name or alias of the column.
   */
  name(): string
 }
 interface ColumnType {
  /**
   * Length returns the column type length for variable length column types such
   * as text and binary field types. If the type length is unbounded the value will
   * be [math.MaxInt64] (any database limits will still apply).
   * If the column type is not variable length, such as an int, or if not supported
   * by the driver ok is false.
   */
  length(): [number, boolean]
 }
 interface ColumnType {
  /**
   * DecimalSize returns the scale and precision of a decimal type.
   * If not applicable or if not supported ok is false.
   */
  decimalSize(): [number, number, boolean]
 }
 interface ColumnType {
  /**
   * ScanType returns a Go type suitable for scanning into using [Rows.Scan].
   * If a driver does not support this property ScanType will return
   * the type of an empty interface.
   */
  scanType(): any
 }
 interface ColumnType {
  /**
   * Nullable reports whether the column may be null.
   * If a driver does not support this property ok will be false.
   */
  nullable(): [boolean, boolean]
 }
 interface ColumnType {
  /**
   * DatabaseTypeName returns the database system name of the column type. If an empty
   * string is returned, then the driver type name is not supported.
   * Consult your driver documentation for a list of driver data types. [ColumnType.Length] specifiers
   * are not included.
   * Common type names include "VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL",
   * "INT", and "BIGINT".
   */
  databaseTypeName(): string
 }
 /**
  * Row is the result of calling [DB.QueryRow] to select a single row.
  */
 interface Row {
 }
 interface Row {
  /**
   * Scan copies the columns from the matched row into the values
   * pointed at by dest. See the documentation on [Rows.Scan] for details.
   * If more than one row matches the query,
   * Scan uses the first row and discards the rest. If no row matches
   * the query, Scan returns [ErrNoRows].
   */
  scan(...dest: any[]): void
 }
 interface Row {
  /**
   * Err provides a way for wrapping packages to check for
   * query errors without calling [Row.Scan].
   * Err returns the error, if any, that was encountered while running the query.
   * If this error is not nil, this error will also be returned from [Row.Scan].
   */
  err(): void
 }
}

/**
 * Package textproto implements generic support for text-based request/response
 * protocols in the style of HTTP, NNTP, and SMTP.
 * 
 * The package provides:
 * 
 * [Error], which represents a numeric error response from
 * a server.
 * 
 * [Pipeline], to manage pipelined requests and responses
 * in a client.
 * 
 * [Reader], to read numeric response code lines,
 * key: value headers, lines wrapped with leading spaces
 * on continuation lines, and whole text blocks ending
 * with a dot on a line by itself.
 * 
 * [Writer], to write dot-encoded text blocks.
 * 
 * [Conn], a convenient packaging of [Reader], [Writer], and [Pipeline] for use
 * with a single network connection.
 */
namespace textproto {
 /**
  * A MIMEHeader represents a MIME-style header mapping
  * keys to sets of values.
  */
 interface MIMEHeader extends _TygojaDict{}
 interface MIMEHeader {
  /**
   * Add adds the key, value pair to the header.
   * It appends to any existing values associated with key.
   */
  add(key: string, value: string): void
 }
 interface MIMEHeader {
  /**
   * Set sets the header entries associated with key to
   * the single element value. It replaces any existing
   * values associated with key.
   */
  set(key: string, value: string): void
 }
 interface MIMEHeader {
  /**
   * Get gets the first value associated with the given key.
   * It is case insensitive; [CanonicalMIMEHeaderKey] is used
   * to canonicalize the provided key.
   * If there are no values associated with the key, Get returns "".
   * To use non-canonical keys, access the map directly.
   */
  get(key: string): string
 }
 interface MIMEHeader {
  /**
   * Values returns all values associated with the given key.
   * It is case insensitive; [CanonicalMIMEHeaderKey] is
   * used to canonicalize the provided key. To use non-canonical
   * keys, access the map directly.
   * The returned slice is not a copy.
   */
  values(key: string): Array<string>
 }
 interface MIMEHeader {
  /**
   * Del deletes the values associated with key.
   */
  del(key: string): void
 }
}

namespace multipart {
 interface Reader {
  /**
   * ReadForm parses an entire multipart message whose parts have
   * a Content-Disposition of "form-data".
   * It stores up to maxMemory bytes + 10MB (reserved for non-file parts)
   * in memory. File parts which can't be stored in memory will be stored on
   * disk in temporary files.
   * It returns [ErrMessageTooLarge] if all non-file parts can't be stored in
   * memory.
   */
  readForm(maxMemory: number): (Form)
 }
 /**
  * Form is a parsed multipart form.
  * Its File parts are stored either in memory or on disk,
  * and are accessible via the [*FileHeader]'s Open method.
  * Its Value parts are stored as strings.
  * Both are keyed by field name.
  */
 interface Form {
  value: _TygojaDict
  file: _TygojaDict
 }
 interface Form {
  /**
   * RemoveAll removes any temporary files associated with a [Form].
   */
  removeAll(): void
 }
 /**
  * File is an interface to access the file part of a multipart message.
  * Its contents may be either stored in memory or on disk.
  * If stored on disk, the File's underlying concrete type will be an *os.File.
  */
 interface File {
  [key:string]: any;
 }
 /**
  * Reader is an iterator over parts in a MIME multipart body.
  * Reader's underlying parser consumes its input as needed. Seeking
  * isn't supported.
  */
 interface Reader {
 }
 interface Reader {
  /**
   * NextPart returns the next part in the multipart or an error.
   * When there are no more parts, the error [io.EOF] is returned.
   * 
   * As a special case, if the "Content-Transfer-Encoding" header
   * has a value of "quoted-printable", that header is instead
   * hidden and the body is transparently decoded during Read calls.
   */
  nextPart(): (Part)
 }
 interface Reader {
  /**
   * NextRawPart returns the next part in the multipart or an error.
   * When there are no more parts, the error [io.EOF] is returned.
   * 
   * Unlike [Reader.NextPart], it does not have special handling for
   * "Content-Transfer-Encoding: quoted-printable".
   */
  nextRawPart(): (Part)
 }
}

namespace http {
 /**
  * A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an
  * HTTP response or the Cookie header of an HTTP request.
  * 
  * See https://tools.ietf.org/html/rfc6265 for details.
  */
 interface Cookie {
  name: string
  value: string
  quoted: boolean // indicates whether the Value was originally quoted
  path: string // optional
  domain: string // optional
  expires: time.Time // optional
  rawExpires: string // for reading cookies only
  /**
   * MaxAge=0 means no 'Max-Age' attribute specified.
   * MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
   * MaxAge>0 means Max-Age attribute present and given in seconds
   */
  maxAge: number
  secure: boolean
  httpOnly: boolean
  sameSite: SameSite
  partitioned: boolean
  raw: string
  unparsed: Array<string> // Raw text of unparsed attribute-value pairs
 }
 interface Cookie {
  /**
   * String returns the serialization of the cookie for use in a [Cookie]
   * header (if only Name and Value are set) or a Set-Cookie response
   * header (if other fields are set).
   * If c is nil or c.Name is invalid, the empty string is returned.
   */
  string(): string
 }
 interface Cookie {
  /**
   * Valid reports whether the cookie is valid.
   */
  valid(): void
 }
 // @ts-ignore
 import mathrand = rand
 /**
  * A Header represents the key-value pairs in an HTTP header.
  * 
  * The keys should be in canonical form, as returned by
  * [CanonicalHeaderKey].
  */
 interface Header extends _TygojaDict{}
 interface Header {
  /**
   * Add adds the key, value pair to the header.
   * It appends to any existing values associated with key.
   * The key is case insensitive; it is canonicalized by
   * [CanonicalHeaderKey].
   */
  add(key: string, value: string): void
 }
 interface Header {
  /**
   * Set sets the header entries associated with key to the
   * single element value. It replaces any existing values
   * associated with key. The key is case insensitive; it is
   * canonicalized by [textproto.CanonicalMIMEHeaderKey].
   * To use non-canonical keys, assign to the map directly.
   */
  set(key: string, value: string): void
 }
 interface Header {
  /**
   * Get gets the first value associated with the given key. If
   * there are no values associated with the key, Get returns "".
   * It is case insensitive; [textproto.CanonicalMIMEHeaderKey] is
   * used to canonicalize the provided key. Get assumes that all
   * keys are stored in canonical form. To use non-canonical keys,
   * access the map directly.
   */
  get(key: string): string
 }
 interface Header {
  /**
   * Values returns all values associated with the given key.
   * It is case insensitive; [textproto.CanonicalMIMEHeaderKey] is
   * used to canonicalize the provided key. To use non-canonical
   * keys, access the map directly.
   * The returned slice is not a copy.
   */
  values(key: string): Array<string>
 }
 interface Header {
  /**
   * Del deletes the values associated with key.
   * The key is case insensitive; it is canonicalized by
   * [CanonicalHeaderKey].
   */
  del(key: string): void
 }
 interface Header {
  /**
   * Write writes a header in wire format.
   */
  write(w: io.Writer): void
 }
 interface Header {
  /**
   * Clone returns a copy of h or nil if h is nil.
   */
  clone(): Header
 }
 interface Header {
  /**
   * WriteSubset writes a header in wire format.
   * If exclude is not nil, keys where exclude[key] == true are not written.
   * Keys are not canonicalized before checking the exclude map.
   */
  writeSubset(w: io.Writer, exclude: _TygojaDict): void
 }
 /**
  * Protocols is a set of HTTP protocols.
  * The zero value is an empty set of protocols.
  * 
  * The supported protocols are:
  * 
  * ```
  *   - HTTP1 is the HTTP/1.0 and HTTP/1.1 protocols.
  *     HTTP1 is supported on both unsecured TCP and secured TLS connections.
  * 
  *   - HTTP2 is the HTTP/2 protcol over a TLS connection.
  * 
  *   - UnencryptedHTTP2 is the HTTP/2 protocol over an unsecured TCP connection.
  * ```
  */
 interface Protocols {
 }
 interface Protocols {
  /**
   * HTTP1 reports whether p includes HTTP/1.
   */
  http1(): boolean
 }
 interface Protocols {
  /**
   * SetHTTP1 adds or removes HTTP/1 from p.
   */
  setHTTP1(ok: boolean): void
 }
 interface Protocols {
  /**
   * HTTP2 reports whether p includes HTTP/2.
   */
  http2(): boolean
 }
 interface Protocols {
  /**
   * SetHTTP2 adds or removes HTTP/2 from p.
   */
  setHTTP2(ok: boolean): void
 }
 interface Protocols {
  /**
   * UnencryptedHTTP2 reports whether p includes unencrypted HTTP/2.
   */
  unencryptedHTTP2(): boolean
 }
 interface Protocols {
  /**
   * SetUnencryptedHTTP2 adds or removes unencrypted HTTP/2 from p.
   */
  setUnencryptedHTTP2(ok: boolean): void
 }
 interface Protocols {
  string(): string
 }
 /**
  * HTTP2Config defines HTTP/2 configuration parameters common to
  * both [Transport] and [Server].
  */
 interface HTTP2Config {
  /**
   * MaxConcurrentStreams optionally specifies the number of
   * concurrent streams that a peer may have open at a time.
   * If zero, MaxConcurrentStreams defaults to at least 100.
   */
  maxConcurrentStreams: number
  /**
   * MaxDecoderHeaderTableSize optionally specifies an upper limit for the
   * size of the header compression table used for decoding headers sent
   * by the peer.
   * A valid value is less than 4MiB.
   * If zero or invalid, a default value is used.
   */
  maxDecoderHeaderTableSize: number
  /**
   * MaxEncoderHeaderTableSize optionally specifies an upper limit for the
   * header compression table used for sending headers to the peer.
   * A valid value is less than 4MiB.
   * If zero or invalid, a default value is used.
   */
  maxEncoderHeaderTableSize: number
  /**
   * MaxReadFrameSize optionally specifies the largest frame
   * this endpoint is willing to read.
   * A valid value is between 16KiB and 16MiB, inclusive.
   * If zero or invalid, a default value is used.
   */
  maxReadFrameSize: number
  /**
   * MaxReceiveBufferPerConnection is the maximum size of the
   * flow control window for data received on a connection.
   * A valid value is at least 64KiB and less than 4MiB.
   * If invalid, a default value is used.
   */
  maxReceiveBufferPerConnection: number
  /**
   * MaxReceiveBufferPerStream is the maximum size of
   * the flow control window for data received on a stream (request).
   * A valid value is less than 4MiB.
   * If zero or invalid, a default value is used.
   */
  maxReceiveBufferPerStream: number
  /**
   * SendPingTimeout is the timeout after which a health check using a ping
   * frame will be carried out if no frame is received on a connection.
   * If zero, no health check is performed.
   */
  sendPingTimeout: time.Duration
  /**
   * PingTimeout is the timeout after which a connection will be closed
   * if a response to a ping is not received.
   * If zero, a default of 15 seconds is used.
   */
  pingTimeout: time.Duration
  /**
   * WriteByteTimeout is the timeout after which a connection will be
   * closed if no data can be written to it. The timeout begins when data is
   * available to write, and is extended whenever any bytes are written.
   */
  writeByteTimeout: time.Duration
  /**
   * PermitProhibitedCipherSuites, if true, permits the use of
   * cipher suites prohibited by the HTTP/2 spec.
   */
  permitProhibitedCipherSuites: boolean
  /**
   * CountError, if non-nil, is called on HTTP/2 errors.
   * It is intended to increment a metric for monitoring.
   * The errType contains only lowercase letters, digits, and underscores
   * (a-z, 0-9, _).
   */
  countError: (errType: string) => void
 }
 // @ts-ignore
 import urlpkg = url
 /**
  * Response represents the response from an HTTP request.
  * 
  * The [Client] and [Transport] return Responses from servers once
  * the response headers have been received. The response body
  * is streamed on demand as the Body field is read.
  */
 interface Response {
  status: string // e.g. "200 OK"
  statusCode: number // e.g. 200
  proto: string // e.g. "HTTP/1.0"
  protoMajor: number // e.g. 1
  protoMinor: number // e.g. 0
  /**
   * Header maps header keys to values. If the response had multiple
   * headers with the same key, they may be concatenated, with comma
   * delimiters.  (RFC 7230, section 3.2.2 requires that multiple headers
   * be semantically equivalent to a comma-delimited sequence.) When
   * Header values are duplicated by other fields in this struct (e.g.,
   * ContentLength, TransferEncoding, Trailer), the field values are
   * authoritative.
   * 
   * Keys in the map are canonicalized (see CanonicalHeaderKey).
   */
  header: Header
  /**
   * Body represents the response body.
   * 
   * The response body is streamed on demand as the Body field
   * is read. If the network connection fails or the server
   * terminates the response, Body.Read calls return an error.
   * 
   * The http Client and Transport guarantee that Body is always
   * non-nil, even on responses without a body or responses with
   * a zero-length body. It is the caller's responsibility to
   * close Body. The default HTTP client's Transport may not
   * reuse HTTP/1.x "keep-alive" TCP connections if the Body is
   * not read to completion and closed.
   * 
   * The Body is automatically dechunked if the server replied
   * with a "chunked" Transfer-Encoding.
   * 
   * As of Go 1.12, the Body will also implement io.Writer
   * on a successful "101 Switching Protocols" response,
   * as used by WebSockets and HTTP/2's "h2c" mode.
   */
  body: io.ReadCloser
  /**
   * ContentLength records the length of the associated content. The
   * value -1 indicates that the length is unknown. Unless Request.Method
   * is "HEAD", values >= 0 indicate that the given number of bytes may
   * be read from Body.
   */
  contentLength: number
  /**
   * Contains transfer encodings from outer-most to inner-most. Value is
   * nil, means that "identity" encoding is used.
   */
  transferEncoding: Array<string>
  /**
   * Close records whether the header directed that the connection be
   * closed after reading Body. The value is advice for clients: neither
   * ReadResponse nor Response.Write ever closes a connection.
   */
  close: boolean
  /**
   * Uncompressed reports whether the response was sent compressed but
   * was decompressed by the http package. When true, reading from
   * Body yields the uncompressed content instead of the compressed
   * content actually set from the server, ContentLength is set to -1,
   * and the "Content-Length" and "Content-Encoding" fields are deleted
   * from the responseHeader. To get the original response from
   * the server, set Transport.DisableCompression to true.
   */
  uncompressed: boolean
  /**
   * Trailer maps trailer keys to values in the same
   * format as Header.
   * 
   * The Trailer initially contains only nil values, one for
   * each key specified in the server's "Trailer" header
   * value. Those values are not added to Header.
   * 
   * Trailer must not be accessed concurrently with Read calls
   * on the Body.
   * 
   * After Body.Read has returned io.EOF, Trailer will contain
   * any trailer values sent by the server.
   */
  trailer: Header
  /**
   * Request is the request that was sent to obtain this Response.
   * Request's Body is nil (having already been consumed).
   * This is only populated for Client requests.
   */
  request?: Request
  /**
   * TLS contains information about the TLS connection on which the
   * response was received. It is nil for unencrypted responses.
   * The pointer is shared between responses and should not be
   * modified.
   */
  tls?: any
 }
 interface Response {
  /**
   * Cookies parses and returns the cookies set in the Set-Cookie headers.
   */
  cookies(): Array<(Cookie | undefined)>
 }
 interface Response {
  /**
   * Location returns the URL of the response's "Location" header,
   * if present. Relative redirects are resolved relative to
   * [Response.Request]. [ErrNoLocation] is returned if no
   * Location header is present.
   */
  location(): (url.URL)
 }
 interface Response {
  /**
   * ProtoAtLeast reports whether the HTTP protocol used
   * in the response is at least major.minor.
   */
  protoAtLeast(major: number, minor: number): boolean
 }
 interface Response {
  /**
   * Write writes r to w in the HTTP/1.x server response format,
   * including the status line, headers, body, and optional trailer.
   * 
   * This method consults the following fields of the response r:
   * 
   * ```
   * 	StatusCode
   * 	ProtoMajor
   * 	ProtoMinor
   * 	Request.Method
   * 	TransferEncoding
   * 	Trailer
   * 	Body
   * 	ContentLength
   * 	Header, values for non-canonical keys will have unpredictable behavior
   * ```
   * 
   * The Response Body is closed after it is sent.
   */
  write(w: io.Writer): void
 }
 /**
  * A ConnState represents the state of a client connection to a server.
  * It's used by the optional [Server.ConnState] hook.
  */
 interface ConnState extends Number{}
 interface ConnState {
  string(): string
 }
}

namespace types {
}

namespace search {
}

namespace router {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * RouterGroup represents a collection of routes and other sub groups
  * that share common pattern prefix and middlewares.
  */
 interface RouterGroup<T> {
  prefix: string
  middlewares: Array<(hook.Handler<T> | undefined)>
 }
 interface RouterGroup<T> {
  /**
   * Group creates and register a new child Group into the current one
   * with the specified prefix.
   * 
   * The prefix follows the standard Go net/http ServeMux pattern format ("[HOST]/[PATH]")
   * and will be concatenated recursively into the final route path, meaning that
   * only the root level group could have HOST as part of the prefix.
   * 
   * Returns the newly created group to allow chaining and registering
   * sub-routes and group specific middlewares.
   */
  group(prefix: string): (RouterGroup<T>)
 }
 interface RouterGroup<T> {
  /**
   * BindFunc registers one or multiple middleware functions to the current group.
   * 
   * The registered middleware functions are "anonymous" and with default priority,
   * aka. executes in the order they were registered.
   * 
   * If you need to specify a named middleware (ex. so that it can be removed)
   * or middleware with custom exec prirority, use [RouterGroup.Bind] method.
   */
  bindFunc(...middlewareFuncs: ((e: T) => void)[]): (RouterGroup<T>)
 }
 interface RouterGroup<T> {
  /**
   * Bind registers one or multiple middleware handlers to the current group.
   */
  bind(...middlewares: (hook.Handler<T> | undefined)[]): (RouterGroup<T>)
 }
 interface RouterGroup<T> {
  /**
   * Unbind removes one or more middlewares with the specified id(s)
   * from the current group and its children (if any).
   * 
   * Anonymous middlewares are not removable, aka. this method does nothing
   * if the middleware id is an empty string.
   */
  unbind(...middlewareIds: string[]): (RouterGroup<T>)
 }
 interface RouterGroup<T> {
  /**
   * Route registers a single route into the current group.
   * 
   * Note that the final route path will be the concatenation of all parent groups prefixes + the route path.
   * The path follows the standard Go net/http ServeMux format ("[HOST]/[PATH]"),
   * meaning that only a top level group route could have HOST as part of the prefix.
   * 
   * Returns the newly created route to allow attaching route-only middlewares.
   */
  route(method: string, path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * Any is a shorthand for [RouterGroup.AddRoute] with "" as route method (aka. matches any method).
   */
  any(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * GET is a shorthand for [RouterGroup.AddRoute] with GET as route method.
   */
  get(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * SEARCH is a shorthand for [RouterGroup.AddRoute] with SEARCH as route method.
   */
  search(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * POST is a shorthand for [RouterGroup.AddRoute] with POST as route method.
   */
  post(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * DELETE is a shorthand for [RouterGroup.AddRoute] with DELETE as route method.
   */
  delete(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * PATCH is a shorthand for [RouterGroup.AddRoute] with PATCH as route method.
   */
  patch(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * PUT is a shorthand for [RouterGroup.AddRoute] with PUT as route method.
   */
  put(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * HEAD is a shorthand for [RouterGroup.AddRoute] with HEAD as route method.
   */
  head(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * OPTIONS is a shorthand for [RouterGroup.AddRoute] with OPTIONS as route method.
   */
  options(path: string, action: (e: T) => void): (Route<T>)
 }
 interface RouterGroup<T> {
  /**
   * HasRoute checks whether the specified route pattern (method + path)
   * is registered in the current group or its children.
   * 
   * This could be useful to conditionally register and checks for routes
   * in order prevent panic on duplicated routes.
   * 
   * Note that routes with anonymous and named wildcard placeholder are treated as equal,
   * aka. "GET /abc/" is considered the same as "GET /abc/{something...}".
   */
  hasRoute(method: string, path: string): boolean
 }
}

namespace cobra {
 interface PositionalArgs {(cmd: Command, args: Array<string>): void }
 // @ts-ignore
 import flag = pflag
 /**
  * FParseErrWhitelist configures Flag parse errors to be ignored
  */
 interface FParseErrWhitelist extends _TygojaAny{}
 /**
  * Group Structure to manage groups for commands
  */
 interface Group {
  id: string
  title: string
 }
 /**
  * CompletionOptions are the options to control shell completion
  */
 interface CompletionOptions {
  /**
   * DisableDefaultCmd prevents Cobra from creating a default 'completion' command
   */
  disableDefaultCmd: boolean
  /**
   * DisableNoDescFlag prevents Cobra from creating the '--no-descriptions' flag
   * for shells that support completion descriptions
   */
  disableNoDescFlag: boolean
  /**
   * DisableDescriptions turns off all completion descriptions for shells
   * that support them
   */
  disableDescriptions: boolean
  /**
   * HiddenDefaultCmd makes the default 'completion' command hidden
   */
  hiddenDefaultCmd: boolean
  /**
   * DefaultShellCompDirective sets the ShellCompDirective that is returned
   * if no special directive can be determined
   */
  defaultShellCompDirective?: ShellCompDirective
 }
 interface CompletionOptions {
  setDefaultShellCompDirective(directive: ShellCompDirective): void
 }
 /**
  * Completion is a string that can be used for completions
  * 
  * two formats are supported:
  * ```
  *   - the completion choice
  *   - the completion choice with a textual description (separated by a TAB).
  * ```
  * 
  * [CompletionWithDesc] can be used to create a completion string with a textual description.
  * 
  * Note: Go type alias is used to provide a more descriptive name in the documentation, but any string can be used.
  */
 interface Completion extends String{}
 /**
  * CompletionFunc is a function that provides completion results.
  */
 interface CompletionFunc {(cmd: Command, args: Array<string>, toComplete: string): [Array<Completion>, ShellCompDirective] }
}

namespace slog {
 /**
  * An Attr is a key-value pair.
  */
 interface Attr {
  key: string
  value: Value
 }
 interface Attr {
  /**
   * Equal reports whether a and b have equal keys and values.
   */
  equal(b: Attr): boolean
 }
 interface Attr {
  string(): string
 }
 /**
  * A Handler handles log records produced by a Logger.
  * 
  * A typical handler may print log records to standard error,
  * or write them to a file or database, or perhaps augment them
  * with additional attributes and pass them on to another handler.
  * 
  * Any of the Handler's methods may be called concurrently with itself
  * or with other methods. It is the responsibility of the Handler to
  * manage this concurrency.
  * 
  * Users of the slog package should not invoke Handler methods directly.
  * They should use the methods of [Logger] instead.
  */
 interface Handler {
  [key:string]: any;
  /**
   * Enabled reports whether the handler handles records at the given level.
   * The handler ignores records whose level is lower.
   * It is called early, before any arguments are processed,
   * to save effort if the log event should be discarded.
   * If called from a Logger method, the first argument is the context
   * passed to that method, or context.Background() if nil was passed
   * or the method does not take a context.
   * The context is passed so Enabled can use its values
   * to make a decision.
   */
  enabled(_arg0: context.Context, _arg1: Level): boolean
  /**
   * Handle handles the Record.
   * It will only be called when Enabled returns true.
   * The Context argument is as for Enabled.
   * It is present solely to provide Handlers access to the context's values.
   * Canceling the context should not affect record processing.
   * (Among other things, log messages may be necessary to debug a
   * cancellation-related problem.)
   * 
   * Handle methods that produce output should observe the following rules:
   * ```
   *   - If r.Time is the zero time, ignore the time.
   *   - If r.PC is zero, ignore it.
   *   - Attr's values should be resolved.
   *   - If an Attr's key and value are both the zero value, ignore the Attr.
   *     This can be tested with attr.Equal(Attr{}).
   *   - If a group's key is empty, inline the group's Attrs.
   *   - If a group has no Attrs (even if it has a non-empty key),
   *     ignore it.
   * ```
   */
  handle(_arg0: context.Context, _arg1: Record): void
  /**
   * WithAttrs returns a new Handler whose attributes consist of
   * both the receiver's attributes and the arguments.
   * The Handler owns the slice: it may retain, modify or discard it.
   */
  withAttrs(attrs: Array<Attr>): Handler
  /**
   * WithGroup returns a new Handler with the given group appended to
   * the receiver's existing groups.
   * The keys of all subsequent attributes, whether added by With or in a
   * Record, should be qualified by the sequence of group names.
   * 
   * How this qualification happens is up to the Handler, so long as
   * this Handler's attribute keys differ from those of another Handler
   * with a different sequence of group names.
   * 
   * A Handler should treat WithGroup as starting a Group of Attrs that ends
   * at the end of the log event. That is,
   * 
   * ```
   *     logger.WithGroup("s").LogAttrs(ctx, level, msg, slog.Int("a", 1), slog.Int("b", 2))
   * ```
   * 
   * should behave like
   * 
   * ```
   *     logger.LogAttrs(ctx, level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
   * ```
   * 
   * If the name is empty, WithGroup returns the receiver.
   */
  withGroup(name: string): Handler
 }
 /**
  * A Level is the importance or severity of a log event.
  * The higher the level, the more important or severe the event.
  */
 interface Level extends Number{}
 interface Level {
  /**
   * String returns a name for the level.
   * If the level has a name, then that name
   * in uppercase is returned.
   * If the level is between named values, then
   * an integer is appended to the uppercased name.
   * Examples:
   * 
   * ```
   * 	LevelWarn.String() => "WARN"
   * 	(LevelInfo+2).String() => "INFO+2"
   * ```
   */
  string(): string
 }
 interface Level {
  /**
   * MarshalJSON implements [encoding/json.Marshaler]
   * by quoting the output of [Level.String].
   */
  marshalJSON(): string|Array<number>
 }
 interface Level {
  /**
   * UnmarshalJSON implements [encoding/json.Unmarshaler]
   * It accepts any string produced by [Level.MarshalJSON],
   * ignoring case.
   * It also accepts numeric offsets that would result in a different string on
   * output. For example, "Error-8" would marshal as "INFO".
   */
  unmarshalJSON(data: string|Array<number>): void
 }
 interface Level {
  /**
   * AppendText implements [encoding.TextAppender]
   * by calling [Level.String].
   */
  appendText(b: string|Array<number>): string|Array<number>
 }
 interface Level {
  /**
   * MarshalText implements [encoding.TextMarshaler]
   * by calling [Level.AppendText].
   */
  marshalText(): string|Array<number>
 }
 interface Level {
  /**
   * UnmarshalText implements [encoding.TextUnmarshaler].
   * It accepts any string produced by [Level.MarshalText],
   * ignoring case.
   * It also accepts numeric offsets that would result in a different string on
   * output. For example, "Error-8" would marshal as "INFO".
   */
  unmarshalText(data: string|Array<number>): void
 }
 interface Level {
  /**
   * Level returns the receiver.
   * It implements [Leveler].
   */
  level(): Level
 }
 // @ts-ignore
 import loginternal = internal
}

/**
 * Package oauth2 provides support for making
 * OAuth2 authorized and authenticated HTTP requests,
 * as specified in RFC 6749.
 * It can additionally grant authorization with Bearer JWT.
 */
namespace oauth2 {
 /**
  * An AuthCodeOption is passed to Config.AuthCodeURL.
  */
 interface AuthCodeOption {
  [key:string]: any;
 }
 /**
  * Token represents the credentials used to authorize
  * the requests to access protected resources on the OAuth 2.0
  * provider's backend.
  * 
  * Most users of this package should not access fields of Token
  * directly. They're exported mostly for use by related packages
  * implementing derivative OAuth2 flows.
  */
 interface Token {
  /**
   * AccessToken is the token that authorizes and authenticates
   * the requests.
   */
  accessToken: string
  /**
   * TokenType is the type of token.
   * The Type method returns either this or "Bearer", the default.
   */
  tokenType: string
  /**
   * RefreshToken is a token that's used by the application
   * (as opposed to the user) to refresh the access token
   * if it expires.
   */
  refreshToken: string
  /**
   * Expiry is the optional expiration time of the access token.
   * 
   * If zero, [TokenSource] implementations will reuse the same
   * token forever and RefreshToken or equivalent
   * mechanisms for that TokenSource will not be used.
   */
  expiry: time.Time
  /**
   * ExpiresIn is the OAuth2 wire format "expires_in" field,
   * which specifies how many seconds later the token expires,
   * relative to an unknown time base approximately around "now".
   * It is the application's responsibility to populate
   * `Expiry` from `ExpiresIn` when required.
   */
  expiresIn: number
 }
 interface Token {
  /**
   * Type returns t.TokenType if non-empty, else "Bearer".
   */
  type(): string
 }
 interface Token {
  /**
   * SetAuthHeader sets the Authorization header to r using the access
   * token in t.
   * 
   * This method is unnecessary when using [Transport] or an HTTP Client
   * returned by this package.
   */
  setAuthHeader(r: http.Request): void
 }
 interface Token {
  /**
   * WithExtra returns a new [Token] that's a clone of t, but using the
   * provided raw extra map. This is only intended for use by packages
   * implementing derivative OAuth2 flows.
   */
  withExtra(extra: any): (Token)
 }
 interface Token {
  /**
   * Extra returns an extra field.
   * Extra fields are key-value pairs returned by the server as a
   * part of the token retrieval response.
   */
  extra(key: string): any
 }
 interface Token {
  /**
   * Valid reports whether t is non-nil, has an AccessToken, and is not expired.
   */
  valid(): boolean
 }
}

namespace subscriptions {
}

namespace url {
 /**
  * The Userinfo type is an immutable encapsulation of username and
  * password details for a [URL]. An existing Userinfo value is guaranteed
  * to have a username set (potentially empty, as allowed by RFC 2396),
  * and optionally a password.
  */
 interface Userinfo {
 }
 interface Userinfo {
  /**
   * Username returns the username.
   */
  username(): string
 }
 interface Userinfo {
  /**
   * Password returns the password in case it is set, and whether it is set.
   */
  password(): [string, boolean]
 }
 interface Userinfo {
  /**
   * String returns the encoded userinfo information in the standard form
   * of "username[:password]".
   */
  string(): string
 }
}

namespace multipart {
 /**
  * A Part represents a single part in a multipart body.
  */
 interface Part {
  /**
   * The headers of the body, if any, with the keys canonicalized
   * in the same fashion that the Go http.Request headers are.
   * For example, "foo-bar" changes case to "Foo-Bar"
   */
  header: textproto.MIMEHeader
 }
 interface Part {
  /**
   * FormName returns the name parameter if p has a Content-Disposition
   * of type "form-data".  Otherwise it returns the empty string.
   */
  formName(): string
 }
 interface Part {
  /**
   * FileName returns the filename parameter of the [Part]'s Content-Disposition
   * header. If not empty, the filename is passed through filepath.Base (which is
   * platform dependent) before being returned.
   */
  fileName(): string
 }
 interface Part {
  /**
   * Read reads the body of a part, after its headers and before the
   * next part (if any) begins.
   */
  read(d: string|Array<number>): number
 }
 interface Part {
  close(): void
 }
}

namespace http {
 /**
  * SameSite allows a server to define a cookie attribute making it impossible for
  * the browser to send this cookie along with cross-site requests. The main
  * goal is to mitigate the risk of cross-origin information leakage, and provide
  * some protection against cross-site request forgery attacks.
  * 
  * See https://tools.ietf.org/html/draft-ietf-httpbis-cookie-same-site-00 for details.
  */
 interface SameSite extends Number{}
 // @ts-ignore
 import mathrand = rand
 // @ts-ignore
 import urlpkg = url
}

namespace router {
 // @ts-ignore
 import validation = ozzo_validation
 interface Route<T> {
  action: (e: T) => void
  method: string
  path: string
  middlewares: Array<(hook.Handler<T> | undefined)>
 }
 interface Route<T> {
  /**
   * BindFunc registers one or multiple middleware functions to the current route.
   * 
   * The registered middleware functions are "anonymous" and with default priority,
   * aka. executes in the order they were registered.
   * 
   * If you need to specify a named middleware (ex. so that it can be removed)
   * or middleware with custom exec prirority, use the [Route.Bind] method.
   */
  bindFunc(...middlewareFuncs: ((e: T) => void)[]): (Route<T>)
 }
 interface Route<T> {
  /**
   * Bind registers one or multiple middleware handlers to the current route.
   */
  bind(...middlewares: (hook.Handler<T> | undefined)[]): (Route<T>)
 }
 interface Route<T> {
  /**
   * Unbind removes one or more middlewares with the specified id(s) from the current route.
   * 
   * It also adds the removed middleware ids to an exclude list so that they could be skipped from
   * the execution chain in case the middleware is registered in a parent group.
   * 
   * Anonymous middlewares are considered non-removable, aka. this method
   * does nothing if the middleware id is an empty string.
   */
  unbind(...middlewareIds: string[]): (Route<T>)
 }
}

namespace cobra {
 // @ts-ignore
 import flag = pflag
 /**
  * ShellCompDirective is a bit map representing the different behaviors the shell
  * can be instructed to have once completions have been provided.
  */
 interface ShellCompDirective extends Number{}
}

namespace slog {
 // @ts-ignore
 import loginternal = internal
 /**
  * A Record holds information about a log event.
  * Copies of a Record share state.
  * Do not modify a Record after handing out a copy to it.
  * Call [NewRecord] to create a new Record.
  * Use [Record.Clone] to create a copy with no shared state.
  */
 interface Record {
  /**
   * The time at which the output method (Log, Info, etc.) was called.
   */
  time: time.Time
  /**
   * The log message.
   */
  message: string
  /**
   * The level of the event.
   */
  level: Level
  /**
   * The program counter at the time the record was constructed, as determined
   * by runtime.Callers. If zero, no program counter is available.
   * 
   * The only valid use for this value is as an argument to
   * [runtime.CallersFrames]. In particular, it must not be passed to
   * [runtime.FuncForPC].
   */
  pc: number
 }
 interface Record {
  /**
   * Clone returns a copy of the record with no shared state.
   * The original record and the clone can both be modified
   * without interfering with each other.
   */
  clone(): Record
 }
 interface Record {
  /**
   * NumAttrs returns the number of attributes in the [Record].
   */
  numAttrs(): number
 }
 interface Record {
  /**
   * Attrs calls f on each Attr in the [Record].
   * Iteration stops if f returns false.
   */
  attrs(f: (_arg0: Attr) => boolean): void
 }
 interface Record {
  /**
   * AddAttrs appends the given Attrs to the [Record]'s list of Attrs.
   * It omits empty groups.
   */
  addAttrs(...attrs: Attr[]): void
 }
 interface Record {
  /**
   * Add converts the args to Attrs as described in [Logger.Log],
   * then appends the Attrs to the [Record]'s list of Attrs.
   * It omits empty groups.
   */
  add(...args: any[]): void
 }
 /**
  * A Value can represent any Go value, but unlike type any,
  * it can represent most small values without an allocation.
  * The zero Value corresponds to nil.
  */
 interface Value {
 }
 interface Value {
  /**
   * Kind returns v's Kind.
   */
  kind(): Kind
 }
 interface Value {
  /**
   * Any returns v's value as an any.
   */
  any(): any
 }
 interface Value {
  /**
   * String returns Value's value as a string, formatted like [fmt.Sprint]. Unlike
   * the methods Int64, Float64, and so on, which panic if v is of the
   * wrong kind, String never panics.
   */
  string(): string
 }
 interface Value {
  /**
   * Int64 returns v's value as an int64. It panics
   * if v is not a signed integer.
   */
  int64(): number
 }
 interface Value {
  /**
   * Uint64 returns v's value as a uint64. It panics
   * if v is not an unsigned integer.
   */
  uint64(): number
 }
 interface Value {
  /**
   * Bool returns v's value as a bool. It panics
   * if v is not a bool.
   */
  bool(): boolean
 }
 interface Value {
  /**
   * Duration returns v's value as a [time.Duration]. It panics
   * if v is not a time.Duration.
   */
  duration(): time.Duration
 }
 interface Value {
  /**
   * Float64 returns v's value as a float64. It panics
   * if v is not a float64.
   */
  float64(): number
 }
 interface Value {
  /**
   * Time returns v's value as a [time.Time]. It panics
   * if v is not a time.Time.
   */
  time(): time.Time
 }
 interface Value {
  /**
   * LogValuer returns v's value as a LogValuer. It panics
   * if v is not a LogValuer.
   */
  logValuer(): LogValuer
 }
 interface Value {
  /**
   * Group returns v's value as a []Attr.
   * It panics if v's [Kind] is not [KindGroup].
   */
  group(): Array<Attr>
 }
 interface Value {
  /**
   * Equal reports whether v and w represent the same Go value.
   */
  equal(w: Value): boolean
 }
 interface Value {
  /**
   * Resolve repeatedly calls LogValue on v while it implements [LogValuer],
   * and returns the result.
   * If v resolves to a group, the group's attributes' values are not recursively
   * resolved.
   * If the number of LogValue calls exceeds a threshold, a Value containing an
   * error is returned.
   * Resolve's return value is guaranteed not to be of Kind [KindLogValuer].
   */
  resolve(): Value
 }
}

namespace oauth2 {
}

namespace router {
 // @ts-ignore
 import validation = ozzo_validation
}

namespace slog {
 // @ts-ignore
 import loginternal = internal
 /**
  * Kind is the kind of a [Value].
  */
 interface Kind extends Number{}
 interface Kind {
  string(): string
 }
 /**
  * A LogValuer is any Go value that can convert itself into a Value for logging.
  * 
  * This mechanism may be used to defer expensive operations until they are
  * needed, or to expand a single value into a sequence of components.
  */
 interface LogValuer {
  [key:string]: any;
  logValue(): Value
 }
}
