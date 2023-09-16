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
 * routerAdd("GET", "/hello", (c) => {
 *     return c.json(200, {"message": "Hello!"})
 * }, $apis.requireAdminOrRecordAuth())
 * ```
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
 * ```js
 * routerUse((next) => {
 *     return (c) => {
 *         console.log(c.Path())
 *         return next(c)
 *     }
 * })
 * ```
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
 * ```js
 * routerPre((next) => {
 *     return (c) => {
 *         console.log(c.request().url)
 *         return next(c)
 *     }
 * })
 * ```
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

// skip on* hook methods as they are registered via the global on* method
type appWithoutHooks = Omit<pocketbase.PocketBase, `on${string}`>

/**
 * `$app` is the current running PocketBase instance that is globally
 * available in each .pb.js file.
 *
 * _Note that this variable is available only in pb_hooks context._
 *
 * @namespace
 * @group PocketBase
 */
declare var $app: appWithoutHooks

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
 * ```js
 * const rawBody = readerToString(c.request().body)
 * ```
 *
 * @group PocketBase
 */
declare function readerToString(reader: any, maxBytes?: number): string;

/**
 * arrayOf creates a placeholder array of the specified models.
 * Usually used to populate DB result into an array of models.
 *
 * Example:
 *
 * ```js
 * const records = arrayOf(new Record)
 *
 * $app.dao().recordQuery("articles").limit(10).all(records)
 * ```
 *
 * @group PocketBase
 */
declare function arrayOf<T>(model: T): Array<T>;

/**
 * DynamicModel creates a new dynamic model with fields from the provided data shape.
 *
 * Example:
 *
 * ```js
 * const model = new DynamicModel({
 *     name:  ""
 *     age:    0,
 *     active: false,
 *     roles:  [],
 *     meta:   {}
 * })
 * ```
 *
 * @group PocketBase
 */
declare class DynamicModel {
  constructor(shape?: { [key:string]: any })
}

/**
 * Record model class.
 *
 * ```js
 * const collection = $app.dao().findCollectionByNameOrId("article")
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
  new(collection?: models.Collection, data?: { [key:string]: any }): models.Record

  // note: declare as "newable" const due to conflict with the Record TS utility type
}

interface Collection extends models.Collection{} // merge
/**
 * Collection model class.
 *
 * ```js
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
 * ```
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
 * ```js
 * const admin = new Admin()
 * admin.email = "test@example.com"
 * admin.setPassword(1234567890)
 * ```
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

interface RequestInfo extends models.RequestInfo{} // merge
/**
 * RequestInfo defines a single models.RequestInfo instance, usually used
 * as part of various filter checks.
 *
 * Example:
 *
 * ```js
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
 * ```
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
 * ```js
 * const dt0 = new DateTime() // now
 *
 * const dt1 = new DateTime('2023-07-01 00:00:00.000Z')
 * ```
 *
 * @group PocketBase
 */
declare class DateTime implements types.DateTime {
  constructor(date?: string)
}

interface ValidationError extends ozzo_validation.Error{} // merge
/**
 * ValidationError defines a single formatted data validation error,
 * usually used as part of a error response.
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
// tokensBinds
// -------------------------------------------------------------------

/**
 * `$tokens` defines high level helpers to generate
 * various admins and auth records tokens (auth, forgotten password, etc.).
 *
 * For more control over the generated token, you can check `$security`.
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
 * `$security` defines low level helpers for creating
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
  let hs256:                          security.hs256
  let hs512:                          security.hs512
  let equal:                          security.equal
  let md5:                            security.md5
  let sha256:                         security.sha256
  let sha512:                         security.sha512
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
  export let exec:      exec.command
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
  export function staticDirectoryHandler(dir: string, indexFallback: boolean): echo.HandlerFunc

  let requireRecordAuth:         apis.requireRecordAuth
  let requireAdminAuth:          apis.requireAdminAuth
  let requireAdminAuthOnlyIfAny: apis.requireAdminAuthOnlyIfAny
  let requireAdminOrRecordAuth:  apis.requireAdminOrRecordAuth
  let requireAdminOrOwnerAuth:   apis.requireAdminOrOwnerAuth
  let activityLogger:            apis.activityLogger
  let requestInfo:               apis.requestInfo
  let recordAuthResponse:        apis.recordAuthResponse
  let enrichRecord:              apis.enrichRecord
  let enrichRecords:             apis.enrichRecords
}

// -------------------------------------------------------------------
// httpClientBinds
// -------------------------------------------------------------------

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
   *     url:    "https://example.com",
   *     data:   {"title": "test"}
   *     method: "post",
   * })
   *
   * console.log(res.statusCode) // the response HTTP status code
   * console.log(res.headers)    // the response headers (eg. res.headers['X-Custom'][0])
   * console.log(res.cookies)    // the response cookies (eg. res.cookies.sessionId.value)
   * console.log(res.raw)        // the response body as plain text
   * console.log(res.json)       // the response body as parsed json array or map
   * ```
   */
  function send(config: {
    url:      string,
    body?:    string,
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
/** @group PocketBase */declare function onAdminAfterAuthRefreshRequest(handler: (e: core.AdminAuthRefreshEvent) => void): void
/** @group PocketBase */declare function onAdminAfterAuthWithPasswordRequest(handler: (e: core.AdminAuthWithPasswordEvent) => void): void
/** @group PocketBase */declare function onAdminAfterConfirmPasswordResetRequest(handler: (e: core.AdminConfirmPasswordResetEvent) => void): void
/** @group PocketBase */declare function onAdminAfterCreateRequest(handler: (e: core.AdminCreateEvent) => void): void
/** @group PocketBase */declare function onAdminAfterDeleteRequest(handler: (e: core.AdminDeleteEvent) => void): void
/** @group PocketBase */declare function onAdminAfterRequestPasswordResetRequest(handler: (e: core.AdminRequestPasswordResetEvent) => void): void
/** @group PocketBase */declare function onAdminAfterUpdateRequest(handler: (e: core.AdminUpdateEvent) => void): void
/** @group PocketBase */declare function onAdminAuthRequest(handler: (e: core.AdminAuthEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeAuthRefreshRequest(handler: (e: core.AdminAuthRefreshEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeAuthWithPasswordRequest(handler: (e: core.AdminAuthWithPasswordEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeConfirmPasswordResetRequest(handler: (e: core.AdminConfirmPasswordResetEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeCreateRequest(handler: (e: core.AdminCreateEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeDeleteRequest(handler: (e: core.AdminDeleteEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeRequestPasswordResetRequest(handler: (e: core.AdminRequestPasswordResetEvent) => void): void
/** @group PocketBase */declare function onAdminBeforeUpdateRequest(handler: (e: core.AdminUpdateEvent) => void): void
/** @group PocketBase */declare function onAdminViewRequest(handler: (e: core.AdminViewEvent) => void): void
/** @group PocketBase */declare function onAdminsListRequest(handler: (e: core.AdminsListEvent) => void): void
/** @group PocketBase */declare function onAfterApiError(handler: (e: core.ApiErrorEvent) => void): void
/** @group PocketBase */declare function onAfterBootstrap(handler: (e: core.BootstrapEvent) => void): void
/** @group PocketBase */declare function onBeforeApiError(handler: (e: core.ApiErrorEvent) => void): void
/** @group PocketBase */declare function onBeforeBootstrap(handler: (e: core.BootstrapEvent) => void): void
/** @group PocketBase */declare function onCollectionAfterCreateRequest(handler: (e: core.CollectionCreateEvent) => void): void
/** @group PocketBase */declare function onCollectionAfterDeleteRequest(handler: (e: core.CollectionDeleteEvent) => void): void
/** @group PocketBase */declare function onCollectionAfterUpdateRequest(handler: (e: core.CollectionUpdateEvent) => void): void
/** @group PocketBase */declare function onCollectionBeforeCreateRequest(handler: (e: core.CollectionCreateEvent) => void): void
/** @group PocketBase */declare function onCollectionBeforeDeleteRequest(handler: (e: core.CollectionDeleteEvent) => void): void
/** @group PocketBase */declare function onCollectionBeforeUpdateRequest(handler: (e: core.CollectionUpdateEvent) => void): void
/** @group PocketBase */declare function onCollectionViewRequest(handler: (e: core.CollectionViewEvent) => void): void
/** @group PocketBase */declare function onCollectionsAfterImportRequest(handler: (e: core.CollectionsImportEvent) => void): void
/** @group PocketBase */declare function onCollectionsBeforeImportRequest(handler: (e: core.CollectionsImportEvent) => void): void
/** @group PocketBase */declare function onCollectionsListRequest(handler: (e: core.CollectionsListEvent) => void): void
/** @group PocketBase */declare function onFileAfterTokenRequest(handler: (e: core.FileTokenEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onFileBeforeTokenRequest(handler: (e: core.FileTokenEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onFileDownloadRequest(handler: (e: core.FileDownloadEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerAfterAdminResetPasswordSend(handler: (e: core.MailerAdminEvent) => void): void
/** @group PocketBase */declare function onMailerAfterRecordChangeEmailSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerAfterRecordResetPasswordSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerAfterRecordVerificationSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerBeforeAdminResetPasswordSend(handler: (e: core.MailerAdminEvent) => void): void
/** @group PocketBase */declare function onMailerBeforeRecordChangeEmailSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerBeforeRecordResetPasswordSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onMailerBeforeRecordVerificationSend(handler: (e: core.MailerRecordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterCreate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterDelete(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelAfterUpdate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelBeforeCreate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelBeforeDelete(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onModelBeforeUpdate(handler: (e: core.ModelEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRealtimeAfterMessageSend(handler: (e: core.RealtimeMessageEvent) => void): void
/** @group PocketBase */declare function onRealtimeAfterSubscribeRequest(handler: (e: core.RealtimeSubscribeEvent) => void): void
/** @group PocketBase */declare function onRealtimeBeforeMessageSend(handler: (e: core.RealtimeMessageEvent) => void): void
/** @group PocketBase */declare function onRealtimeBeforeSubscribeRequest(handler: (e: core.RealtimeSubscribeEvent) => void): void
/** @group PocketBase */declare function onRealtimeConnectRequest(handler: (e: core.RealtimeConnectEvent) => void): void
/** @group PocketBase */declare function onRealtimeDisconnectRequest(handler: (e: core.RealtimeDisconnectEvent) => void): void
/** @group PocketBase */declare function onRecordAfterAuthRefreshRequest(handler: (e: core.RecordAuthRefreshEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterAuthWithOAuth2Request(handler: (e: core.RecordAuthWithOAuth2Event) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterAuthWithPasswordRequest(handler: (e: core.RecordAuthWithPasswordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterConfirmEmailChangeRequest(handler: (e: core.RecordConfirmEmailChangeEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterConfirmPasswordResetRequest(handler: (e: core.RecordConfirmPasswordResetEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterConfirmVerificationRequest(handler: (e: core.RecordConfirmVerificationEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterCreateRequest(handler: (e: core.RecordCreateEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterDeleteRequest(handler: (e: core.RecordDeleteEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterRequestEmailChangeRequest(handler: (e: core.RecordRequestEmailChangeEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterRequestPasswordResetRequest(handler: (e: core.RecordRequestPasswordResetEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterRequestVerificationRequest(handler: (e: core.RecordRequestVerificationEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterUnlinkExternalAuthRequest(handler: (e: core.RecordUnlinkExternalAuthEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAfterUpdateRequest(handler: (e: core.RecordUpdateEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordAuthRequest(handler: (e: core.RecordAuthEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeAuthRefreshRequest(handler: (e: core.RecordAuthRefreshEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeAuthWithOAuth2Request(handler: (e: core.RecordAuthWithOAuth2Event) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeAuthWithPasswordRequest(handler: (e: core.RecordAuthWithPasswordEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeConfirmEmailChangeRequest(handler: (e: core.RecordConfirmEmailChangeEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeConfirmPasswordResetRequest(handler: (e: core.RecordConfirmPasswordResetEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeConfirmVerificationRequest(handler: (e: core.RecordConfirmVerificationEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeCreateRequest(handler: (e: core.RecordCreateEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeDeleteRequest(handler: (e: core.RecordDeleteEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeRequestEmailChangeRequest(handler: (e: core.RecordRequestEmailChangeEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeRequestPasswordResetRequest(handler: (e: core.RecordRequestPasswordResetEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeRequestVerificationRequest(handler: (e: core.RecordRequestVerificationEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeUnlinkExternalAuthRequest(handler: (e: core.RecordUnlinkExternalAuthEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordBeforeUpdateRequest(handler: (e: core.RecordUpdateEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordListExternalAuthsRequest(handler: (e: core.RecordListExternalAuthsEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordViewRequest(handler: (e: core.RecordViewEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onRecordsListRequest(handler: (e: core.RecordsListEvent) => void, ...tags: string[]): void
/** @group PocketBase */declare function onSettingsAfterUpdateRequest(handler: (e: core.SettingsUpdateEvent) => void): void
/** @group PocketBase */declare function onSettingsBeforeUpdateRequest(handler: (e: core.SettingsUpdateEvent) => void): void
/** @group PocketBase */declare function onSettingsListRequest(handler: (e: core.SettingsListEvent) => void): void
/** @group PocketBase */declare function onTerminate(handler: (e: core.TerminateEvent) => void): void
type _TygojaDict = { [key:string | number | symbol]: any; }
type _TygojaAny = any

/**
 * Package os provides a platform-independent interface to operating system
 * functionality. The design is Unix-like, although the error handling is
 * Go-like; failing calls return values of type error rather than error numbers.
 * Often, more information is available within the error. For example,
 * if a call that takes a file name fails, such as Open or Stat, the error
 * will include the failing file name when printed and will be of type
 * *PathError, which may be unpacked for more information.
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
 * Note: The maximum number of concurrent operations on a File may be limited by
 * the OS or the system. The number should be high, but exceeding it may degrade
 * performance or cause other issues.
 */
namespace os {
 interface readdirMode extends Number{}
 interface File {
  /**
   * Readdir reads the contents of the directory associated with file and
   * returns a slice of up to n FileInfo values, as would be returned
   * by Lstat, in directory order. Subsequent calls on the same file will yield
   * further FileInfos.
   * 
   * If n > 0, Readdir returns at most n FileInfo structures. In this case, if
   * Readdir returns an empty slice, it will return a non-nil error
   * explaining why. At the end of a directory, the error is io.EOF.
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
   * explaining why. At the end of a directory, the error is io.EOF.
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
  * (using the ReadDir function or a File's ReadDir method).
  */
 interface DirEntry extends fs.DirEntry{}
 interface File {
  /**
   * ReadDir reads the contents of the directory associated with the file f
   * and returns a slice of DirEntry values in directory order.
   * Subsequent calls on the same file will yield later DirEntry records in the directory.
   * 
   * If n > 0, ReadDir returns at most n DirEntry records.
   * In this case, if ReadDir returns an empty slice, it will return an error explaining why.
   * At the end of a directory, the error is io.EOF.
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
 /**
  * Auxiliary information if the File describes a directory
  */
 interface dirInfo {
 }
 interface expand {
  /**
   * Expand replaces ${var} or $var in the string based on the mapping function.
   * For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).
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
   * To distinguish between an empty value and an unset value, use LookupEnv.
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
  (key: string): void
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
   * NewSyscallError returns, as an error, a new SyscallError
   * with the given system call name and error details.
   * As a convenience, if err is nil, NewSyscallError returns nil.
   */
  (syscall: string, err: Error): void
 }
 interface isExist {
  /**
   * IsExist returns a boolean indicating whether the error is known to report
   * that a file or directory already exists. It is satisfied by ErrExist as
   * well as some syscall errors.
   * 
   * This function predates errors.Is. It only supports errors returned by
   * the os package. New code should use errors.Is(err, fs.ErrExist).
   */
  (err: Error): boolean
 }
 interface isNotExist {
  /**
   * IsNotExist returns a boolean indicating whether the error is known to
   * report that a file or directory does not exist. It is satisfied by
   * ErrNotExist as well as some syscall errors.
   * 
   * This function predates errors.Is. It only supports errors returned by
   * the os package. New code should use errors.Is(err, fs.ErrNotExist).
   */
  (err: Error): boolean
 }
 interface isPermission {
  /**
   * IsPermission returns a boolean indicating whether the error is known to
   * report that permission is denied. It is satisfied by ErrPermission as well
   * as some syscall errors.
   * 
   * This function predates errors.Is. It only supports errors returned by
   * the os package. New code should use errors.Is(err, fs.ErrPermission).
   */
  (err: Error): boolean
 }
 interface isTimeout {
  /**
   * IsTimeout returns a boolean indicating whether the error is known
   * to report that a timeout occurred.
   * 
   * This function predates errors.Is, and the notion of whether an
   * error indicates a timeout can be ambiguous. For example, the Unix
   * error EWOULDBLOCK sometimes indicates a timeout and sometimes does not.
   * New code should use errors.Is with a value appropriate to the call
   * returning the error, such as os.ErrDeadlineExceeded.
   */
  (err: Error): boolean
 }
 interface syscallErrorType extends syscall.Errno{}
 /**
  * Process stores the information about a process created by StartProcess.
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
   * The Process it returns can be used to obtain information
   * about the underlying operating system process.
   * 
   * On Unix systems, FindProcess always succeeds and returns a Process
   * for the given pid, regardless of whether the process exists.
   */
  (pid: number): (Process | undefined)
 }
 interface startProcess {
  /**
   * StartProcess starts a new process with the program, arguments and attributes
   * specified by name, argv and attr. The argv slice will become os.Args in the
   * new process, so it normally starts with the program name.
   * 
   * If the calling goroutine has locked the operating system thread
   * with runtime.LockOSThread and modified any inheritable OS-level
   * thread state (for example, Linux or Plan 9 name spaces), the new
   * process will inherit the caller's thread state.
   * 
   * StartProcess is a low-level interface. The os/exec package provides
   * higher-level interfaces.
   * 
   * If there is an error, it will be of type *PathError.
   */
  (name: string, argv: Array<string>, attr: ProcAttr): (Process | undefined)
 }
 interface Process {
  /**
   * Release releases any resources associated with the Process p,
   * rendering it unusable in the future.
   * Release only needs to be called if Wait is not.
   */
  release(): void
 }
 interface Process {
  /**
   * Kill causes the Process to exit immediately. Kill does not wait until
   * the Process has actually exited. This only kills the Process itself,
   * not any other processes it may have started.
   */
  kill(): void
 }
 interface Process {
  /**
   * Wait waits for the Process to exit, and then returns a
   * ProcessState describing its status and an error, if any.
   * Wait releases any resources associated with the Process.
   * On most operating systems, the Process must be a child
   * of the current process or an error will be returned.
   */
  wait(): (ProcessState | undefined)
 }
 interface Process {
  /**
   * Signal sends a signal to the Process.
   * Sending Interrupt on Windows is not implemented.
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
   * type, such as syscall.WaitStatus on Unix, to access its contents.
   */
  sys(): any
 }
 interface ProcessState {
  /**
   * SysUsage returns system-dependent resource usage information about
   * the exited process. Convert it to the appropriate underlying
   * type, such as *syscall.Rusage on Unix, to access its contents.
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
   * needed, path/filepath.EvalSymlinks might help.
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
  read(b: string): number
 }
 interface File {
  /**
   * ReadAt reads len(b) bytes from the File starting at byte offset off.
   * It returns the number of bytes read and the error, if any.
   * ReadAt always returns a non-nil error when n < len(b).
   * At end of file, that error is io.EOF.
   */
  readAt(b: string, off: number): number
 }
 interface File {
  /**
   * ReadFrom implements io.ReaderFrom.
   */
  readFrom(r: io.Reader): number
 }
 type _subFPBxR = io.Writer
 interface onlyWriter extends _subFPBxR {
 }
 interface File {
  /**
   * Write writes len(b) bytes from b to the File.
   * It returns the number of bytes written and an error, if any.
   * Write returns a non-nil error when n != len(b).
   */
  write(b: string): number
 }
 interface File {
  /**
   * WriteAt writes len(b) bytes to the File starting at byte offset off.
   * It returns the number of bytes written and an error, if any.
   * WriteAt returns a non-nil error when n != len(b).
   * 
   * If file was opened with the O_APPEND flag, WriteAt returns an error.
   */
  writeAt(b: string, off: number): number
 }
 interface File {
  /**
   * Seek sets the offset for the next Read or Write on file to offset, interpreted
   * according to whence: 0 means relative to the origin of the file, 1 means
   * relative to the current offset, and 2 means relative to the end.
   * It returns the new offset and an error, if any.
   * The behavior of Seek on a file opened with O_APPEND is not specified.
   * 
   * If f is a directory, the behavior of Seek varies by operating
   * system; you can seek to the beginning of the directory on Unix-like
   * operating systems, but not on Windows.
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
  (name: string): (File | undefined)
 }
 interface create {
  /**
   * Create creates or truncates the named file. If the file already exists,
   * it is truncated. If the file does not exist, it is created with mode 0666
   * (before umask). If successful, methods on the returned File can
   * be used for I/O; the associated file descriptor has mode O_RDWR.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): (File | undefined)
 }
 interface openFile {
  /**
   * OpenFile is the generalized open call; most users will use Open
   * or Create instead. It opens the named file with specified flag
   * (O_RDONLY etc.). If the file does not exist, and the O_CREATE flag
   * is passed, it is created with mode perm (before umask). If successful,
   * methods on the returned File can be used for I/O.
   * If there is an error, it will be of type *PathError.
   */
  (name: string, flag: number, perm: FileMode): (File | undefined)
 }
 interface rename {
  /**
   * Rename renames (moves) oldpath to newpath.
   * If newpath already exists and is not a directory, Rename replaces it.
   * OS-specific restrictions may apply when oldpath and newpath are in different directories.
   * If there is an error, it will be of type *LinkError.
   */
  (oldpath: string): void
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
   * If the location cannot be determined (for example, $HOME is not defined),
   * then it will return an error.
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
   * If the location cannot be determined (for example, $HOME is not defined),
   * then it will return an error.
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
   * On Windows, only the 0200 bit (owner writable) of mode is used; it
   * controls whether the file's read-only attribute is set or cleared.
   * The other bits are currently unused. For compatibility with Go 1.12
   * and earlier, use a non-zero mode. Use mode 0400 for a read-only
   * file and 0600 for a readable+writable file.
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
   * os.Open does. DirFS is therefore not a general substitute for a chroot-style security
   * mechanism when the directory tree contains arbitrary content.
   * 
   * The directory dir must not be "".
   */
  (dir: string): fs.FS
 }
 interface dirFS extends String{}
 interface dirFS {
  open(name: string): fs.File
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
  (name: string): string
 }
 interface writeFile {
  /**
   * WriteFile writes data to the named file, creating it if necessary.
   * If the file does not exist, WriteFile creates it with permissions perm (before umask);
   * otherwise WriteFile truncates it before writing, without changing permissions.
   */
  (name: string, data: string, perm: FileMode): void
 }
 interface File {
  /**
   * Close closes the File, rendering it unusable for I/O.
   * On files that support SetDeadline, any pending I/O operations will
   * be canceled and return immediately with an ErrClosed error.
   * Close will return an error if it has already been called.
   */
  close(): void
 }
 interface chown {
  /**
   * Chown changes the numeric uid and gid of the named file.
   * If the file is a symbolic link, it changes the uid and gid of the link's target.
   * A uid or gid of -1 means to not change that value.
   * If there is an error, it will be of type *PathError.
   * 
   * On Windows or Plan 9, Chown always returns the syscall.EWINDOWS or
   * EPLAN9 error, wrapped in *PathError.
   */
  (name: string, uid: number): void
 }
 interface lchown {
  /**
   * Lchown changes the numeric uid and gid of the named file.
   * If the file is a symbolic link, it changes the uid and gid of the link itself.
   * If there is an error, it will be of type *PathError.
   * 
   * On Windows, it always returns the syscall.EWINDOWS error, wrapped
   * in *PathError.
   */
  (name: string, uid: number): void
 }
 interface File {
  /**
   * Chown changes the numeric uid and gid of the named file.
   * If there is an error, it will be of type *PathError.
   * 
   * On Windows, it always returns the syscall.EWINDOWS error, wrapped
   * in *PathError.
   */
  chown(uid: number): void
 }
 interface File {
  /**
   * Truncate changes the size of the file.
   * It does not change the I/O offset.
   * If there is an error, it will be of type *PathError.
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
   * 
   * The underlying filesystem may truncate or round the values to a
   * less precise time unit.
   * If there is an error, it will be of type *PathError.
   */
  (name: string, atime: time.Time, mtime: time.Time): void
 }
 interface File {
  /**
   * Chdir changes the current working directory to the file,
   * which must be a directory.
   * If there is an error, it will be of type *PathError.
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
   * making it invalid; see runtime.SetFinalizer for more information on when
   * a finalizer might be run. On Unix systems this will cause the SetDeadline
   * methods to stop working.
   * Because file descriptors can be reused, the returned file descriptor may
   * only be closed through the Close method of f, or by its finalizer during
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
  (fd: number, name: string): (File | undefined)
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
  (oldname: string): void
 }
 interface symlink {
  /**
   * Symlink creates newname as a symbolic link to oldname.
   * On Windows, a symlink to a non-existent oldname creates a file symlink;
   * if oldname is later created as a directory the symlink will not work.
   * If there is an error, it will be of type *LinkError.
   */
  (oldname: string): void
 }
 interface readlink {
  /**
   * Readlink returns the destination of the named symbolic link.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): string
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
 interface getwd {
  /**
   * Getwd returns a rooted path name corresponding to the
   * current directory. If the current directory can be
   * reached via multiple paths (due to symbolic links),
   * Getwd may return any one of them.
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
   * If there is an error, it will be of type *PathError.
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
  (): [(File | undefined), (File | undefined)]
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
   * On Windows, it returns syscall.EWINDOWS. See the os/user package
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
 interface stat {
  /**
   * Stat returns a FileInfo describing the named file.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): FileInfo
 }
 interface lstat {
  /**
   * Lstat returns a FileInfo describing the named file.
   * If the file is a symbolic link, the returned FileInfo
   * describes the symbolic link. Lstat makes no attempt to follow the link.
   * If there is an error, it will be of type *PathError.
   */
  (name: string): FileInfo
 }
 interface File {
  /**
   * Stat returns the FileInfo structure describing file.
   * If there is an error, it will be of type *PathError.
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
   * If dir is the empty string, CreateTemp uses the default directory for temporary files, as returned by TempDir.
   * Multiple programs or goroutines calling CreateTemp simultaneously will not choose the same file.
   * The caller can use the file's Name method to find the pathname of the file.
   * It is the caller's responsibility to remove the file when it is no longer needed.
   */
  (dir: string): (File | undefined)
 }
 interface mkdirTemp {
  /**
   * MkdirTemp creates a new temporary directory in the directory dir
   * and returns the pathname of the new directory.
   * The new directory's name is generated by adding a random string to the end of pattern.
   * If pattern includes a "*", the random string replaces the last "*" instead.
   * If dir is the empty string, MkdirTemp uses the default directory for temporary files, as returned by TempDir.
   * Multiple programs or goroutines calling MkdirTemp simultaneously will not choose the same directory.
   * It is the caller's responsibility to remove the directory when it is no longer needed.
   */
  (dir: string): string
 }
 interface getpagesize {
  /**
   * Getpagesize returns the underlying system's memory page size.
   */
  (): number
 }
 /**
  * File represents an open file descriptor.
  */
 type _subtNQMD = file
 interface File extends _subtNQMD {
 }
 /**
  * A FileInfo describes a file and is returned by Stat and Lstat.
  */
 interface FileInfo extends fs.FileInfo{}
 /**
  * A FileMode represents a file's mode and permission bits.
  * The bits have the same definition on all systems, so that
  * information about files can be moved from one system
  * to another portably. Not all bits apply to all systems.
  * The only required bit is ModeDir for directories.
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
   * SameFile only applies to results returned by this package's Stat.
   * It returns false in other cases.
   */
  (fi1: FileInfo): boolean
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
 * system, see the path package.
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
   * The only possible returned error is ErrBadPattern, when pattern
   * is malformed.
   * 
   * On Windows, escaping is disabled. Instead, '\\' is treated as
   * path separator.
   */
  (pattern: string): boolean
 }
 interface glob {
  /**
   * Glob returns the names of all files matching pattern or nil
   * if there is no matching file. The syntax of patterns is the same
   * as in Match. The pattern may describe hierarchical names such as
   * /usr/*\/bin/ed (assuming the Separator is '/').
   * 
   * Glob ignores file system errors such as I/O errors reading directories.
   * The only possible returned error is ErrBadPattern, when pattern
   * is malformed.
   */
  (pattern: string): Array<string>
 }
 /**
  * A lazybuf is a lazily constructed path buffer.
  * It supports append, reading previously appended bytes,
  * and retrieving the final string. It does not allocate a buffer
  * to hold the output until that output diverges from s.
  */
 interface lazybuf {
 }
 interface clean {
  /**
   * Clean returns the shortest path name equivalent to path
   * by purely lexical processing. It applies the following rules
   * iteratively until no further processing can be done:
   * 
   * ```
   * 	1. Replace multiple Separator elements with a single one.
   * 	2. Eliminate each . path name element (the current directory).
   * 	3. Eliminate each inner .. path name element (the parent directory)
   * 	   along with the non-.. element that precedes it.
   * 	4. Eliminate .. elements that begin a rooted path:
   * 	   that is, replace "/.." by "/" at the beginning of a path,
   * 	   assuming Separator is '/'.
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
   * See also Rob Pike, ``Lexical File Names in Plan 9 or
   * Getting Dot-Dot Right,''
   * https://9p.io/sys/doc/lexnames.html
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
   */
  (path: string): string
 }
 interface splitList {
  /**
   * SplitList splits a list of paths joined by the OS-specific ListSeparator,
   * usually found in PATH or GOPATH environment variables.
   * Unlike strings.Split, SplitList returns an empty slice when passed an empty
   * string.
   */
  (path: string): Array<string>
 }
 interface split {
  /**
   * Split splits path immediately following the final Separator,
   * separating it into a directory and file name component.
   * If there is no Separator in path, Split returns an empty dir
   * and file set to path.
   * The returned values have the property that path = dir+file.
   */
  (path: string): string
 }
 interface join {
  /**
   * Join joins any number of path elements into a single path,
   * separating them with an OS specific Separator. Empty elements
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
   * EvalSymlinks calls Clean on the result.
   */
  (path: string): string
 }
 interface abs {
  /**
   * Abs returns an absolute representation of path.
   * If the path is not absolute it will be joined with the current
   * working directory to turn it into an absolute path. The absolute
   * path name for a given file is not guaranteed to be unique.
   * Abs calls Clean on the result.
   */
  (path: string): string
 }
 interface rel {
  /**
   * Rel returns a relative path that is lexically equivalent to targpath when
   * joined to basepath with an intervening separator. That is,
   * Join(basepath, Rel(basepath, targpath)) is equivalent to targpath itself.
   * On success, the returned path will always be relative to basepath,
   * even if basepath and targpath share no elements.
   * An error is returned if targpath can't be made relative to basepath or if
   * knowing the current working directory would be necessary to compute it.
   * Rel calls Clean on the result.
   */
  (basepath: string): string
 }
 /**
  * WalkFunc is the type of the function called by Walk to visit each
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
  * If the function returns the special value SkipDir, Walk skips the
  * current directory (path if info.IsDir() is true, otherwise path's
  * parent directory). Otherwise, if the function returns a non-nil error,
  * Walk stops entirely and returns that error.
  * 
  * The err argument reports an error related to path, signaling that Walk
  * will not walk into that directory. The function can decide how to
  * handle that error; as described earlier, returning the error will
  * cause Walk to stop walking the entire tree.
  * 
  * Walk calls the function with a non-nil err argument in two cases.
  * 
  * First, if an os.Lstat on the root directory or any directory or file
  * in the tree fails, Walk calls the function with path set to that
  * directory or file's path, info set to nil, and err set to the error
  * from os.Lstat.
  * 
  * Second, if a directory's Readdirnames method fails, Walk calls the
  * function with path set to the directory's path, info, set to an
  * fs.FileInfo describing the directory, and err set to the error from
  * Readdirnames.
  */
 interface WalkFunc {(path: string, info: fs.FileInfo, err: Error): void }
 interface walkDir {
  /**
   * WalkDir walks the file tree rooted at root, calling fn for each file or
   * directory in the tree, including root.
   * 
   * All errors that arise visiting files and directories are filtered by fn:
   * see the fs.WalkDirFunc documentation for details.
   * 
   * The files are walked in lexical order, which makes the output deterministic
   * but requires WalkDir to read an entire directory into memory before proceeding
   * to walk that directory.
   * 
   * WalkDir does not follow symbolic links.
   */
  (root: string, fn: fs.WalkDirFunc): void
 }
 interface statDirEntry {
 }
 interface statDirEntry {
  name(): string
 }
 interface statDirEntry {
  isDir(): boolean
 }
 interface statDirEntry {
  type(): fs.FileMode
 }
 interface statDirEntry {
  info(): fs.FileInfo
 }
 interface walk {
  /**
   * Walk walks the file tree rooted at root, calling fn for each file or
   * directory in the tree, including root.
   * 
   * All errors that arise visiting files and directories are filtered by fn:
   * see the WalkFunc documentation for details.
   * 
   * The files are walked in lexical order, which makes the output deterministic
   * but requires Walk to read an entire directory into memory before proceeding
   * to walk that directory.
   * 
   * Walk does not follow symbolic links.
   * 
   * Walk is less efficient than WalkDir, introduced in Go 1.16,
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
   * After dropping the final element, Dir calls Clean on the path and trailing
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
 interface isAbs {
  /**
   * IsAbs reports whether the path is absolute.
   */
  (path: string): boolean
 }
 interface hasPrefix {
  /**
   * HasPrefix exists for historical compatibility and should not be used.
   * 
   * Deprecated: HasPrefix does not respect path boundaries and
   * does not ignore case when required.
   */
  (p: string): boolean
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
  error(): string
  code(): string
  message(): string
  setMessage(_arg0: string): Error
  params(): _TygojaDict
  setParams(_arg0: _TygojaDict): Error
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
   * Encrypt encrypts data with key (must be valid 32 char aes key).
   */
  (data: string, key: string): string
 }
 interface decrypt {
  /**
   * Decrypt decrypts encrypted text with key (must be valid 32 chars aes key).
   */
  (cipherText: string, key: string): string
 }
 interface parseUnverifiedJWT {
  /**
   * ParseUnverifiedJWT parses JWT token and returns its claims
   * but DOES NOT verify the signature.
   * 
   * It verifies only the exp, iat and nbf claims.
   */
  (token: string): jwt.MapClaims
 }
 interface parseJWT {
  /**
   * ParseJWT verifies and parses JWT token and returns its claims.
   */
  (token: string, verificationKey: string): jwt.MapClaims
 }
 interface newJWT {
  /**
   * NewJWT generates and returns new HS256 signed JWT token.
   */
  (payload: jwt.MapClaims, signingKey: string, secondsDuration: number): string
 }
 interface newToken {
  /**
   * Deprecated:
   * Consider replacing with NewJWT().
   * 
   * NewToken is a legacy alias for NewJWT that generates a HS256 signed JWT token.
   */
  (payload: jwt.MapClaims, signingKey: string, secondsDuration: number): string
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
  /**
   * NewQuery creates a new Query object with the given SQL statement.
   * The SQL statement may contain parameter placeholders which can be bound with actual parameter
   * values before the statement is executed.
   */
  newQuery(_arg0: string): (Query | undefined)
  /**
   * Select returns a new SelectQuery object that can be used to build a SELECT statement.
   * The parameters to this method should be the list column names to be selected.
   * A column name may have an optional alias name. For example, Select("id", "my_name AS name").
   */
  select(..._arg0: string[]): (SelectQuery | undefined)
  /**
   * ModelQuery returns a new ModelQuery object that can be used to perform model insertion, update, and deletion.
   * The parameter to this method should be a pointer to the model struct that needs to be inserted, updated, or deleted.
   */
  model(_arg0: {
  }): (ModelQuery | undefined)
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
  insert(table: string, cols: Params): (Query | undefined)
  /**
   * Upsert creates a Query that represents an UPSERT SQL statement.
   * Upsert inserts a row into the table if the primary key or unique index is not found.
   * Otherwise it will update the row with the new values.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  upsert(table: string, cols: Params, ...constraints: string[]): (Query | undefined)
  /**
   * Update creates a Query that represents an UPDATE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding new column
   * values. If the "where" expression is nil, the UPDATE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will update ALL rows in the table).
   */
  update(table: string, cols: Params, where: Expression): (Query | undefined)
  /**
   * Delete creates a Query that represents a DELETE SQL statement.
   * If the "where" expression is nil, the DELETE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will delete ALL rows in the table).
   */
  delete(table: string, where: Expression): (Query | undefined)
  /**
   * CreateTable creates a Query that represents a CREATE TABLE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding column types.
   * The optional "options" parameters will be appended to the generated SQL statement.
   */
  createTable(table: string, cols: _TygojaDict, ...options: string[]): (Query | undefined)
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string): (Query | undefined)
  /**
   * DropTable creates a Query that can be used to drop a table.
   */
  dropTable(table: string): (Query | undefined)
  /**
   * TruncateTable creates a Query that can be used to truncate a table.
   */
  truncateTable(table: string): (Query | undefined)
  /**
   * AddColumn creates a Query that can be used to add a column to a table.
   */
  addColumn(table: string): (Query | undefined)
  /**
   * DropColumn creates a Query that can be used to drop a column from a table.
   */
  dropColumn(table: string): (Query | undefined)
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string): (Query | undefined)
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string): (Query | undefined)
  /**
   * AddPrimaryKey creates a Query that can be used to specify primary key(s) for a table.
   * The "name" parameter specifies the name of the primary key constraint.
   */
  addPrimaryKey(table: string, ...cols: string[]): (Query | undefined)
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string): (Query | undefined)
  /**
   * AddForeignKey creates a Query that can be used to add a foreign key constraint to a table.
   * The length of cols and refCols must be the same as they refer to the primary and referential columns.
   * The optional "options" parameters will be appended to the SQL statement. They can be used to
   * specify options such as "ON DELETE CASCADE".
   */
  addForeignKey(table: string, cols: Array<string>, refTable: string, ...options: string[]): (Query | undefined)
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string): (Query | undefined)
  /**
   * CreateIndex creates a Query that can be used to create an index for a table.
   */
  createIndex(table: string, ...cols: string[]): (Query | undefined)
  /**
   * CreateUniqueIndex creates a Query that can be used to create a unique index for a table.
   */
  createUniqueIndex(table: string, ...cols: string[]): (Query | undefined)
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string): (Query | undefined)
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
  (db: DB, executor: Executor): (BaseBuilder | undefined)
 }
 interface BaseBuilder {
  /**
   * DB returns the DB instance that this builder is associated with.
   */
  db(): (DB | undefined)
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
  newQuery(sql: string): (Query | undefined)
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
  insert(table: string, cols: Params): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * Upsert creates a Query that represents an UPSERT SQL statement.
   * Upsert inserts a row into the table if the primary key or unique index is not found.
   * Otherwise it will update the row with the new values.
   * The keys of cols are the column names, while the values of cols are the corresponding column
   * values to be inserted.
   */
  upsert(table: string, cols: Params, ...constraints: string[]): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * Update creates a Query that represents an UPDATE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding new column
   * values. If the "where" expression is nil, the UPDATE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will update ALL rows in the table).
   */
  update(table: string, cols: Params, where: Expression): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * Delete creates a Query that represents a DELETE SQL statement.
   * If the "where" expression is nil, the DELETE SQL statement will have no WHERE clause
   * (be careful in this case as the SQL statement will delete ALL rows in the table).
   */
  delete(table: string, where: Expression): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * CreateTable creates a Query that represents a CREATE TABLE SQL statement.
   * The keys of cols are the column names, while the values of cols are the corresponding column types.
   * The optional "options" parameters will be appended to the generated SQL statement.
   */
  createTable(table: string, cols: _TygojaDict, ...options: string[]): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * DropTable creates a Query that can be used to drop a table.
   */
  dropTable(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * TruncateTable creates a Query that can be used to truncate a table.
   */
  truncateTable(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * AddColumn creates a Query that can be used to add a column to a table.
   */
  addColumn(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * DropColumn creates a Query that can be used to drop a column from a table.
   */
  dropColumn(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * AddPrimaryKey creates a Query that can be used to specify primary key(s) for a table.
   * The "name" parameter specifies the name of the primary key constraint.
   */
  addPrimaryKey(table: string, ...cols: string[]): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * AddForeignKey creates a Query that can be used to add a foreign key constraint to a table.
   * The length of cols and refCols must be the same as they refer to the primary and referential columns.
   * The optional "options" parameters will be appended to the SQL statement. They can be used to
   * specify options such as "ON DELETE CASCADE".
   */
  addForeignKey(table: string, cols: Array<string>, refTable: string, ...options: string[]): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * CreateIndex creates a Query that can be used to create an index for a table.
   */
  createIndex(table: string, ...cols: string[]): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * CreateUniqueIndex creates a Query that can be used to create a unique index for a table.
   */
  createUniqueIndex(table: string, ...cols: string[]): (Query | undefined)
 }
 interface BaseBuilder {
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string): (Query | undefined)
 }
 /**
  * MssqlBuilder is the builder for SQL Server databases.
  */
 type _subklAen = BaseBuilder
 interface MssqlBuilder extends _subklAen {
 }
 /**
  * MssqlQueryBuilder is the query builder for SQL Server databases.
  */
 type _subdEVgQ = BaseQueryBuilder
 interface MssqlQueryBuilder extends _subdEVgQ {
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
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface MssqlBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery | undefined)
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
  renameTable(oldName: string): (Query | undefined)
 }
 interface MssqlBuilder {
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string): (Query | undefined)
 }
 interface MssqlBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string): (Query | undefined)
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
 type _subXuKqX = BaseBuilder
 interface MysqlBuilder extends _subXuKqX {
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
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface MysqlBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery | undefined)
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
  upsert(table: string, cols: Params, ...constraints: string[]): (Query | undefined)
 }
 interface MysqlBuilder {
  /**
   * RenameColumn creates a Query that can be used to rename a column in a table.
   */
  renameColumn(table: string): (Query | undefined)
 }
 interface MysqlBuilder {
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string): (Query | undefined)
 }
 interface MysqlBuilder {
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string): (Query | undefined)
 }
 /**
  * OciBuilder is the builder for Oracle databases.
  */
 type _subMOhJY = BaseBuilder
 interface OciBuilder extends _subMOhJY {
 }
 /**
  * OciQueryBuilder is the query builder for Oracle databases.
  */
 type _subNPnto = BaseQueryBuilder
 interface OciQueryBuilder extends _subNPnto {
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
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface OciBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery | undefined)
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
  dropIndex(table: string): (Query | undefined)
 }
 interface OciBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string): (Query | undefined)
 }
 interface OciBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string): (Query | undefined)
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
 type _submlkgA = BaseBuilder
 interface PgsqlBuilder extends _submlkgA {
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
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface PgsqlBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery | undefined)
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
  upsert(table: string, cols: Params, ...constraints: string[]): (Query | undefined)
 }
 interface PgsqlBuilder {
  /**
   * DropIndex creates a Query that can be used to remove the named index from a table.
   */
  dropIndex(table: string): (Query | undefined)
 }
 interface PgsqlBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string): (Query | undefined)
 }
 interface PgsqlBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string): (Query | undefined)
 }
 /**
  * SqliteBuilder is the builder for SQLite databases.
  */
 type _subJtnhw = BaseBuilder
 interface SqliteBuilder extends _subJtnhw {
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
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface SqliteBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery | undefined)
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
  dropIndex(table: string): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * TruncateTable creates a Query that can be used to truncate a table.
   */
  truncateTable(table: string): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * RenameTable creates a Query that can be used to rename a table.
   */
  renameTable(oldName: string): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * AlterColumn creates a Query that can be used to change the definition of a table column.
   */
  alterColumn(table: string): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * AddPrimaryKey creates a Query that can be used to specify primary key(s) for a table.
   * The "name" parameter specifies the name of the primary key constraint.
   */
  addPrimaryKey(table: string, ...cols: string[]): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * DropPrimaryKey creates a Query that can be used to remove the named primary key constraint from a table.
   */
  dropPrimaryKey(table: string): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * AddForeignKey creates a Query that can be used to add a foreign key constraint to a table.
   * The length of cols and refCols must be the same as they refer to the primary and referential columns.
   * The optional "options" parameters will be appended to the SQL statement. They can be used to
   * specify options such as "ON DELETE CASCADE".
   */
  addForeignKey(table: string, cols: Array<string>, refTable: string, ...options: string[]): (Query | undefined)
 }
 interface SqliteBuilder {
  /**
   * DropForeignKey creates a Query that can be used to remove the named foreign key constraint from a table.
   */
  dropForeignKey(table: string): (Query | undefined)
 }
 /**
  * StandardBuilder is the builder that is used by DB for an unknown driver.
  */
 type _subLFjYA = BaseBuilder
 interface StandardBuilder extends _subLFjYA {
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
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface StandardBuilder {
  /**
   * Model returns a new ModelQuery object that can be used to perform model-based DB operations.
   * The model passed to this method should be a pointer to a model struct.
   */
  model(model: {
   }): (ModelQuery | undefined)
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
 type _subvKgBu = Builder
 interface DB extends _subvKgBu {
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
  (sqlDB: sql.DB, driverName: string): (DB | undefined)
 }
 interface open {
  /**
   * Open opens a database specified by a driver name and data source name (DSN).
   * Note that Open does not check if DSN is specified correctly. It doesn't try to establish a DB connection either.
   * Please refer to sql.Open() for more information.
   */
  (driverName: string): (DB | undefined)
 }
 interface mustOpen {
  /**
   * MustOpen opens a database and establishes a connection to it.
   * Please refer to sql.Open() and sql.Ping() for more information.
   */
  (driverName: string): (DB | undefined)
 }
 interface DB {
  /**
   * Clone makes a shallow copy of DB.
   */
  clone(): (DB | undefined)
 }
 interface DB {
  /**
   * WithContext returns a new instance of DB associated with the given context.
   */
  withContext(ctx: context.Context): (DB | undefined)
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
  db(): (sql.DB | undefined)
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
  begin(): (Tx | undefined)
 }
 interface DB {
  /**
   * BeginTx starts a transaction with the given context and transaction options.
   */
  beginTx(ctx: context.Context, opts: sql.TxOptions): (Tx | undefined)
 }
 interface DB {
  /**
   * Wrap encapsulates an existing transaction.
   */
  wrap(sqlTx: sql.Tx): (Tx | undefined)
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
  (col: string, ...values: string[]): (LikeExp | undefined)
 }
 interface notLike {
  /**
   * NotLike generates a NOT LIKE expression.
   * For example, NotLike("name", "key", "word") will generate a SQL expression:
   * "name" NOT LIKE "%key%" AND "name" NOT LIKE "%word%". Please see Like() for more details.
   */
  (col: string, ...values: string[]): (LikeExp | undefined)
 }
 interface orLike {
  /**
   * OrLike generates an OR LIKE expression.
   * This is similar to Like() except that the column should be like one of the possible values.
   * For example, OrLike("name", "key", "word") will generate a SQL expression:
   * "name" LIKE "%key%" OR "name" LIKE "%word%". Please see Like() for more details.
   */
  (col: string, ...values: string[]): (LikeExp | undefined)
 }
 interface orNotLike {
  /**
   * OrNotLike generates an OR NOT LIKE expression.
   * For example, OrNotLike("name", "key", "word") will generate a SQL expression:
   * "name" NOT LIKE "%key%" OR "name" NOT LIKE "%word%". Please see Like() for more details.
   */
  (col: string, ...values: string[]): (LikeExp | undefined)
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
   }): Expression
 }
 interface notBetween {
  /**
   * NotBetween generates a NOT BETWEEN expression.
   * For example, NotBetween("age", 10, 30) generates: "age" NOT BETWEEN 10 AND 30
   */
  (col: string, from: {
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
  escape(...chars: string[]): (LikeExp | undefined)
 }
 interface LikeExp {
  /**
   * Match specifies whether to do wildcard matching on the left and/or right of given strings.
   */
  match(left: boolean): (LikeExp | undefined)
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
  tableName(): string
 }
 /**
  * ModelQuery represents a query associated with a struct model.
  */
 interface ModelQuery {
 }
 interface newModelQuery {
  (model: {
   }, fieldMapFunc: FieldMapFunc, db: DB, builder: Builder): (ModelQuery | undefined)
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
  withContext(ctx: context.Context): (ModelQuery | undefined)
 }
 interface ModelQuery {
  /**
   * Exclude excludes the specified struct fields from being inserted/updated into the DB table.
   */
  exclude(...attrs: string[]): (ModelQuery | undefined)
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
  }[]): (sql.Rows | undefined)
  /**
   * QueryContext queries a SQL statement with the given context
   */
  queryContext(ctx: context.Context, query: string, ...args: {
  }[]): (sql.Rows | undefined)
  /**
   * Prepare creates a prepared statement
   */
  prepare(query: string): (sql.Stmt | undefined)
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
  (db: DB, executor: Executor, sql: string): (Query | undefined)
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
  withContext(ctx: context.Context): (Query | undefined)
 }
 interface Query {
  /**
   * WithExecHook associates the provided exec hook function with the query.
   * 
   * It is called for every Query resolver (Execute(), One(), All(), Row(), Column()),
   * allowing you to implement auto fail/retry or any other additional handling.
   */
  withExecHook(fn: ExecHookFunc): (Query | undefined)
 }
 interface Query {
  /**
   * WithOneHook associates the provided hook function with the query,
   * called on q.One(), allowing you to implement custom struct scan based
   * on the One() argument and/or result.
   */
  withOneHook(fn: OneHookFunc): (Query | undefined)
 }
 interface Query {
  /**
   * WithOneHook associates the provided hook function with the query,
   * called on q.All(), allowing you to implement custom slice scan based
   * on the All() argument and/or result.
   */
  withAllHook(fn: AllHookFunc): (Query | undefined)
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
  prepare(): (Query | undefined)
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
  bind(params: Params): (Query | undefined)
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
  rows(): (Rows | undefined)
 }
 /**
  * QueryBuilder builds different clauses for a SELECT SQL statement.
  */
 interface QueryBuilder {
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
  (db: DB): (BaseQueryBuilder | undefined)
 }
 interface BaseQueryBuilder {
  /**
   * DB returns the DB instance associated with the query builder.
   */
  db(): (DB | undefined)
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
 type _subpryha = sql.Rows
 interface Rows extends _subpryha {
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
  (builder: Builder, db: DB): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * WithBuildHook runs the provided hook function with the query created on Build().
   */
  withBuildHook(fn: BuildHookFunc): (SelectQuery | undefined)
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
  withContext(ctx: context.Context): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Select specifies the columns to be selected.
   * Column names will be automatically quoted.
   */
  select(...cols: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * AndSelect adds additional columns to be selected.
   * Column names will be automatically quoted.
   */
  andSelect(...cols: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Distinct specifies whether to select columns distinctively.
   * By default, distinct is false.
   */
  distinct(v: boolean): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * SelectOption specifies additional option that should be append to "SELECT".
   */
  selectOption(option: string): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * From specifies which tables to select from.
   * Table names will be automatically quoted.
   */
  from(...tables: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Where specifies the WHERE condition.
   */
  where(e: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * AndWhere concatenates a new WHERE condition with the existing one (if any) using "AND".
   */
  andWhere(e: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * OrWhere concatenates a new WHERE condition with the existing one (if any) using "OR".
   */
  orWhere(e: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Join specifies a JOIN clause.
   * The "typ" parameter specifies the JOIN type (e.g. "INNER JOIN", "LEFT JOIN").
   */
  join(typ: string, table: string, on: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * InnerJoin specifies an INNER JOIN clause.
   * This is a shortcut method for Join.
   */
  innerJoin(table: string, on: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * LeftJoin specifies a LEFT JOIN clause.
   * This is a shortcut method for Join.
   */
  leftJoin(table: string, on: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * RightJoin specifies a RIGHT JOIN clause.
   * This is a shortcut method for Join.
   */
  rightJoin(table: string, on: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * OrderBy specifies the ORDER BY clause.
   * Column names will be properly quoted. A column name can contain "ASC" or "DESC" to indicate its ordering direction.
   */
  orderBy(...cols: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * AndOrderBy appends additional columns to the existing ORDER BY clause.
   * Column names will be properly quoted. A column name can contain "ASC" or "DESC" to indicate its ordering direction.
   */
  andOrderBy(...cols: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * GroupBy specifies the GROUP BY clause.
   * Column names will be properly quoted.
   */
  groupBy(...cols: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * AndGroupBy appends additional columns to the existing GROUP BY clause.
   * Column names will be properly quoted.
   */
  andGroupBy(...cols: string[]): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Having specifies the HAVING clause.
   */
  having(e: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * AndHaving concatenates a new HAVING condition with the existing one (if any) using "AND".
   */
  andHaving(e: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * OrHaving concatenates a new HAVING condition with the existing one (if any) using "OR".
   */
  orHaving(e: Expression): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Union specifies a UNION clause.
   */
  union(q: Query): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * UnionAll specifies a UNION ALL clause.
   */
  unionAll(q: Query): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Limit specifies the LIMIT clause.
   * A negative limit means no limit.
   */
  limit(limit: number): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Offset specifies the OFFSET clause.
   * A negative offset means no offset.
   */
  offset(offset: number): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Bind specifies the parameter values to be bound to the query.
   */
  bind(params: Params): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * AndBind appends additional parameters to be bound to the query.
   */
  andBind(params: Params): (SelectQuery | undefined)
 }
 interface SelectQuery {
  /**
   * Build builds the SELECT query and returns an executable Query object.
   */
  build(): (Query | undefined)
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
  rows(): (Rows | undefined)
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
  info(): (QueryInfo | undefined)
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
 type _subYkrPN = structInfo
 interface structValue extends _subYkrPN {
 }
 interface fieldInfo {
 }
 interface structInfoMapKey {
 }
 /**
  * PostScanner is an optional interface used by ScanStruct.
  */
 interface PostScanner {
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
 type _subbdnmI = Builder
 interface Tx extends _subbdnmI {
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
 * dangerous input, or use the path/filepath package's Glob function.
 * To expand environment variables, use package os's ExpandEnv.
 * 
 * Note that the examples in this package assume a Unix system.
 * They may not run on Windows, and they do not run in the Go Playground
 * used by golang.org and godoc.org.
 */
namespace exec {
 interface command {
  /**
   * Command returns the Cmd struct to execute the named program with
   * the given arguments.
   * 
   * It sets only the Path and Args in the returned structure.
   * 
   * If name contains no path separators, Command uses LookPath to
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
  (name: string, ...arg: string[]): (Cmd | undefined)
 }
}

namespace filesystem {
 /**
  * FileReader defines an interface for a file resource reader.
  */
 interface FileReader {
  open(): io.ReadSeekCloser
 }
 /**
  * File defines a single file [io.ReadSeekCloser] resource.
  * 
  * The file could be from a local path, multipipart/formdata header, etc.
  */
 interface File {
  name: string
  originalName: string
  size: number
  reader: FileReader
 }
 interface newFileFromPath {
  /**
   * NewFileFromPath creates a new File instance from the provided local file path.
   */
  (path: string): (File | undefined)
 }
 interface newFileFromBytes {
  /**
   * NewFileFromBytes creates a new File instance from the provided byte slice.
   */
  (b: string, name: string): (File | undefined)
 }
 interface newFileFromMultipart {
  /**
   * NewFileFromMultipart creates a new File instace from the provided multipart header.
   */
  (mh: multipart.FileHeader): (File | undefined)
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
  bytes: string
 }
 interface BytesReader {
  /**
   * Open implements the [filesystem.FileReader] interface.
   */
  open(): io.ReadSeekCloser
 }
 type _subYsCWp = bytes.Reader
 interface bytesReadSeekCloser extends _subYsCWp {
 }
 interface bytesReadSeekCloser {
  /**
   * Close implements the [io.ReadSeekCloser] interface.
   */
  close(): void
 }
 interface System {
 }
 interface newS3 {
  /**
   * NewS3 initializes an S3 filesystem instance.
   * 
   * NB! Make sure to call `Close()` after you are done working with it.
   */
  (bucketName: string, region: string, endpoint: string, accessKey: string, secretKey: string, s3ForcePathStyle: boolean): (System | undefined)
 }
 interface newLocal {
  /**
   * NewLocal initializes a new local filesystem instance.
   * 
   * NB! Make sure to call `Close()` after you are done working with it.
   */
  (dirPath: string): (System | undefined)
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
   */
  attributes(fileKey: string): (blob.Attributes | undefined)
 }
 interface System {
  /**
   * GetFile returns a file content reader for the given fileKey.
   * 
   * NB! Make sure to call `Close()` after you are done working with it.
   */
  getFile(fileKey: string): (blob.Reader | undefined)
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
  upload(content: string, fileKey: string): void
 }
 interface System {
  /**
   * UploadFile uploads the provided multipart file to the fileKey location.
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
   */
  delete(fileKey: string): void
 }
 interface System {
  /**
   * DeletePrefix deletes everything starting with the specified prefix.
   */
  deletePrefix(prefix: string): Array<Error>
 }
 interface System {
  /**
   * Serve serves the file at fileKey location to an HTTP response.
   * 
   * If the `download` query parameter is used the file will be always served for
   * download no matter of its type (aka. with "Content-Disposition: attachment").
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
  createThumb(originalKey: string, thumbKey: string): void
 }
}

/**
 * Package tokens implements various user and admin tokens generation methods.
 */
namespace tokens {
 interface newAdminAuthToken {
  /**
   * NewAdminAuthToken generates and returns a new admin authentication token.
   */
  (app: core.App, admin: models.Admin): string
 }
 interface newAdminResetPasswordToken {
  /**
   * NewAdminResetPasswordToken generates and returns a new admin password reset request token.
   */
  (app: core.App, admin: models.Admin): string
 }
 interface newAdminFileToken {
  /**
   * NewAdminFileToken generates and returns a new admin private file access token.
   */
  (app: core.App, admin: models.Admin): string
 }
 interface newRecordAuthToken {
  /**
   * NewRecordAuthToken generates and returns a new auth record authentication token.
   */
  (app: core.App, record: models.Record): string
 }
 interface newRecordVerifyToken {
  /**
   * NewRecordVerifyToken generates and returns a new record verification token.
   */
  (app: core.App, record: models.Record): string
 }
 interface newRecordResetPasswordToken {
  /**
   * NewRecordResetPasswordToken generates and returns a new auth record password reset request token.
   */
  (app: core.App, record: models.Record): string
 }
 interface newRecordChangeEmailToken {
  /**
   * NewRecordChangeEmailToken generates and returns a new auth record change email request token.
   */
  (app: core.App, record: models.Record, newEmail: string): string
 }
 interface newRecordFileToken {
  /**
   * NewRecordFileToken generates and returns a new record private file access token.
   */
  (app: core.App, record: models.Record): string
 }
}

/**
 * Package models implements various services used for request data
 * validation and applying changes to existing DB models through the app Dao.
 */
namespace forms {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * AdminLogin is an admin email/pass login form.
  */
 interface AdminLogin {
  identity: string
  password: string
 }
 interface newAdminLogin {
  /**
   * NewAdminLogin creates a new [AdminLogin] form initialized with
   * the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App): (AdminLogin | undefined)
 }
 interface AdminLogin {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface AdminLogin {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface AdminLogin {
  /**
   * Submit validates and submits the admin form.
   * On success returns the authorized admin model.
   * 
   * You can optionally provide a list of InterceptorFunc to
   * further modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Admin | undefined>[]): (models.Admin | undefined)
 }
 /**
  * AdminPasswordResetConfirm is an admin password reset confirmation form.
  */
 interface AdminPasswordResetConfirm {
  token: string
  password: string
  passwordConfirm: string
 }
 interface newAdminPasswordResetConfirm {
  /**
   * NewAdminPasswordResetConfirm creates a new [AdminPasswordResetConfirm]
   * form initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App): (AdminPasswordResetConfirm | undefined)
 }
 interface AdminPasswordResetConfirm {
  /**
   * SetDao replaces the form Dao instance with the provided one.
   * 
   * This is useful if you want to use a specific transaction Dao instance
   * instead of the default app.Dao().
   */
  setDao(dao: daos.Dao): void
 }
 interface AdminPasswordResetConfirm {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface AdminPasswordResetConfirm {
  /**
   * Submit validates and submits the admin password reset confirmation form.
   * On success returns the updated admin model associated to `form.Token`.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Admin | undefined>[]): (models.Admin | undefined)
 }
 /**
  * AdminPasswordResetRequest is an admin password reset request form.
  */
 interface AdminPasswordResetRequest {
  email: string
 }
 interface newAdminPasswordResetRequest {
  /**
   * NewAdminPasswordResetRequest creates a new [AdminPasswordResetRequest]
   * form initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App): (AdminPasswordResetRequest | undefined)
 }
 interface AdminPasswordResetRequest {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface AdminPasswordResetRequest {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   * 
   * This method doesn't verify that admin with `form.Email` exists (this is done on Submit).
   */
  validate(): void
 }
 interface AdminPasswordResetRequest {
  /**
   * Submit validates and submits the form.
   * On success sends a password reset email to the `form.Email` admin.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Admin | undefined>[]): void
 }
 /**
  * AdminUpsert is a [models.Admin] upsert (create/update) form.
  */
 interface AdminUpsert {
  id: string
  avatar: number
  email: string
  password: string
  passwordConfirm: string
 }
 interface newAdminUpsert {
  /**
   * NewAdminUpsert creates a new [AdminUpsert] form with initializer
   * config created from the provided [core.App] and [models.Admin] instances
   * (for create you could pass a pointer to an empty Admin - `&models.Admin{}`).
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, admin: models.Admin): (AdminUpsert | undefined)
 }
 interface AdminUpsert {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface AdminUpsert {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface AdminUpsert {
  /**
   * Submit validates the form and upserts the form admin model.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Admin | undefined>[]): void
 }
 /**
  * AppleClientSecretCreate is a [models.Admin] upsert (create/update) form.
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
   * Duration specifies how long the generated JWT token should be considered valid.
   * The specified value must be in seconds and max 15777000 (~6months).
   */
  duration: number
 }
 interface newAppleClientSecretCreate {
  /**
   * NewAppleClientSecretCreate creates a new [AppleClientSecretCreate] form with initializer
   * config created from the provided [core.App] instances.
   */
  (app: core.App): (AppleClientSecretCreate | undefined)
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
 /**
  * BackupCreate is a request form for creating a new app backup.
  */
 interface BackupCreate {
  name: string
 }
 interface newBackupCreate {
  /**
   * NewBackupCreate creates new BackupCreate request form.
   */
  (app: core.App): (BackupCreate | undefined)
 }
 interface BackupCreate {
  /**
   * SetContext replaces the default form context with the provided one.
   */
  setContext(ctx: context.Context): void
 }
 interface BackupCreate {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface BackupCreate {
  /**
   * Submit validates the form and creates the app backup.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before creating the backup.
   */
  submit(...interceptors: InterceptorFunc<string>[]): void
 }
 /**
  * BackupUpload is a request form for uploading a new app backup.
  */
 interface BackupUpload {
  file?: filesystem.File
 }
 interface newBackupUpload {
  /**
   * NewBackupUpload creates new BackupUpload request form.
   */
  (app: core.App): (BackupUpload | undefined)
 }
 interface BackupUpload {
  /**
   * SetContext replaces the default form upload context with the provided one.
   */
  setContext(ctx: context.Context): void
 }
 interface BackupUpload {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface BackupUpload {
  /**
   * Submit validates the form and upload the backup file.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before uploading the backup.
   */
  submit(...interceptors: InterceptorFunc<filesystem.File | undefined>[]): void
 }
 /**
  * InterceptorNextFunc is a interceptor handler function.
  * Usually used in combination with InterceptorFunc.
  */
 interface InterceptorNextFunc<T> {(t: T): void }
 /**
  * InterceptorFunc defines a single interceptor function that
  * will execute the provided next func handler.
  */
 interface InterceptorFunc<T> {(next: InterceptorNextFunc<T>): InterceptorNextFunc<T> }
 /**
  * CollectionUpsert is a [models.Collection] upsert (create/update) form.
  */
 interface CollectionUpsert {
  id: string
  type: string
  name: string
  system: boolean
  schema: schema.Schema
  indexes: types.JsonArray<string>
  listRule?: string
  viewRule?: string
  createRule?: string
  updateRule?: string
  deleteRule?: string
  options: types.JsonMap
 }
 interface newCollectionUpsert {
  /**
   * NewCollectionUpsert creates a new [CollectionUpsert] form with initializer
   * config created from the provided [core.App] and [models.Collection] instances
   * (for create you could pass a pointer to an empty Collection - `&models.Collection{}`).
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (CollectionUpsert | undefined)
 }
 interface CollectionUpsert {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface CollectionUpsert {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface CollectionUpsert {
  /**
   * Submit validates the form and upserts the form's Collection model.
   * 
   * On success the related record table schema will be auto updated.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Collection | undefined>[]): void
 }
 /**
  * CollectionsImport is a form model to bulk import
  * (create, replace and delete) collections from a user provided list.
  */
 interface CollectionsImport {
  collections: Array<(models.Collection | undefined)>
  deleteMissing: boolean
 }
 interface newCollectionsImport {
  /**
   * NewCollectionsImport creates a new [CollectionsImport] form with
   * initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App): (CollectionsImport | undefined)
 }
 interface CollectionsImport {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface CollectionsImport {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface CollectionsImport {
  /**
   * Submit applies the import, aka.:
   * - imports the form collections (create or replace)
   * - sync the collection changes with their related records table
   * - ensures the integrity of the imported structure (aka. run validations for each collection)
   * - if [form.DeleteMissing] is set, deletes all local collections that are not found in the imports list
   * 
   * All operations are wrapped in a single transaction that are
   * rollbacked on the first encountered error.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<Array<(models.Collection | undefined)>>[]): void
 }
 /**
  * RealtimeSubscribe is a realtime subscriptions request form.
  */
 interface RealtimeSubscribe {
  clientId: string
  subscriptions: Array<string>
 }
 interface newRealtimeSubscribe {
  /**
   * NewRealtimeSubscribe creates new RealtimeSubscribe request form.
   */
  (): (RealtimeSubscribe | undefined)
 }
 interface RealtimeSubscribe {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 /**
  * RecordEmailChangeConfirm is an auth record email change confirmation form.
  */
 interface RecordEmailChangeConfirm {
  token: string
  password: string
 }
 interface newRecordEmailChangeConfirm {
  /**
   * NewRecordEmailChangeConfirm creates a new [RecordEmailChangeConfirm] form
   * initialized with from the provided [core.App] and [models.Collection] instances.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (RecordEmailChangeConfirm | undefined)
 }
 interface RecordEmailChangeConfirm {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordEmailChangeConfirm {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordEmailChangeConfirm {
  /**
   * Submit validates and submits the auth record email change confirmation form.
   * On success returns the updated auth record associated to `form.Token`.
   * 
   * You can optionally provide a list of InterceptorFunc to
   * further modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): (models.Record | undefined)
 }
 /**
  * RecordEmailChangeRequest is an auth record email change request form.
  */
 interface RecordEmailChangeRequest {
  newEmail: string
 }
 interface newRecordEmailChangeRequest {
  /**
   * NewRecordEmailChangeRequest creates a new [RecordEmailChangeRequest] form
   * initialized with from the provided [core.App] and [models.Record] instances.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, record: models.Record): (RecordEmailChangeRequest | undefined)
 }
 interface RecordEmailChangeRequest {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordEmailChangeRequest {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordEmailChangeRequest {
  /**
   * Submit validates and sends the change email request.
   * 
   * You can optionally provide a list of InterceptorFunc to
   * further modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): void
 }
 /**
  * RecordOAuth2LoginData defines the OA
  */
 interface RecordOAuth2LoginData {
  externalAuth?: models.ExternalAuth
  record?: models.Record
  oAuth2User?: auth.AuthUser
  providerClient: auth.Provider
 }
 /**
  * BeforeOAuth2RecordCreateFunc defines a callback function that will
  * be called before OAuth2 new Record creation.
  */
 interface BeforeOAuth2RecordCreateFunc {(createForm: RecordUpsert, authRecord: models.Record, authUser: auth.AuthUser): void }
 /**
  * RecordOAuth2Login is an auth record OAuth2 login form.
  */
 interface RecordOAuth2Login {
  /**
   * The name of the OAuth2 client provider (eg. "google")
   */
  provider: string
  /**
   * The authorization code returned from the initial request.
   */
  code: string
  /**
   * The code verifier sent with the initial request as part of the code_challenge.
   */
  codeVerifier: string
  /**
   * The redirect url sent with the initial request.
   */
  redirectUrl: string
  /**
   * Additional data that will be used for creating a new auth record
   * if an existing OAuth2 account doesn't exist.
   */
  createData: _TygojaDict
 }
 interface newRecordOAuth2Login {
  /**
   * NewRecordOAuth2Login creates a new [RecordOAuth2Login] form with
   * initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection, optAuthRecord: models.Record): (RecordOAuth2Login | undefined)
 }
 interface RecordOAuth2Login {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordOAuth2Login {
  /**
   * SetBeforeNewRecordCreateFunc sets a before OAuth2 record create callback handler.
   */
  setBeforeNewRecordCreateFunc(f: BeforeOAuth2RecordCreateFunc): void
 }
 interface RecordOAuth2Login {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordOAuth2Login {
  /**
   * Submit validates and submits the form.
   * 
   * If an auth record doesn't exist, it will make an attempt to create it
   * based on the fetched OAuth2 profile data via a local [RecordUpsert] form.
   * You can intercept/modify the Record create form with [form.SetBeforeNewRecordCreateFunc()].
   * 
   * You can also optionally provide a list of InterceptorFunc to
   * further modify the form behavior before persisting it.
   * 
   * On success returns the authorized record model and the fetched provider's data.
   */
  submit(...interceptors: InterceptorFunc<RecordOAuth2LoginData | undefined>[]): [(models.Record | undefined), (auth.AuthUser | undefined)]
 }
 /**
  * RecordPasswordLogin is record username/email + password login form.
  */
 interface RecordPasswordLogin {
  identity: string
  password: string
 }
 interface newRecordPasswordLogin {
  /**
   * NewRecordPasswordLogin creates a new [RecordPasswordLogin] form initialized
   * with from the provided [core.App] and [models.Collection] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (RecordPasswordLogin | undefined)
 }
 interface RecordPasswordLogin {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordPasswordLogin {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordPasswordLogin {
  /**
   * Submit validates and submits the form.
   * On success returns the authorized record model.
   * 
   * You can optionally provide a list of InterceptorFunc to
   * further modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): (models.Record | undefined)
 }
 /**
  * RecordPasswordResetConfirm is an auth record password reset confirmation form.
  */
 interface RecordPasswordResetConfirm {
  token: string
  password: string
  passwordConfirm: string
 }
 interface newRecordPasswordResetConfirm {
  /**
   * NewRecordPasswordResetConfirm creates a new [RecordPasswordResetConfirm]
   * form initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (RecordPasswordResetConfirm | undefined)
 }
 interface RecordPasswordResetConfirm {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordPasswordResetConfirm {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordPasswordResetConfirm {
  /**
   * Submit validates and submits the form.
   * On success returns the updated auth record associated to `form.Token`.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): (models.Record | undefined)
 }
 /**
  * RecordPasswordResetRequest is an auth record reset password request form.
  */
 interface RecordPasswordResetRequest {
  email: string
 }
 interface newRecordPasswordResetRequest {
  /**
   * NewRecordPasswordResetRequest creates a new [RecordPasswordResetRequest]
   * form initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (RecordPasswordResetRequest | undefined)
 }
 interface RecordPasswordResetRequest {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordPasswordResetRequest {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   * 
   * This method doesn't checks whether auth record with `form.Email` exists (this is done on Submit).
   */
  validate(): void
 }
 interface RecordPasswordResetRequest {
  /**
   * Submit validates and submits the form.
   * On success, sends a password reset email to the `form.Email` auth record.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): void
 }
 /**
  * RecordUpsert is a [models.Record] upsert (create/update) form.
  */
 interface RecordUpsert {
  /**
   * base model fields
   */
  id: string
  /**
   * auth collection fields
   * ---
   */
  username: string
  email: string
  emailVisibility: boolean
  verified: boolean
  password: string
  passwordConfirm: string
  oldPassword: string
 }
 interface newRecordUpsert {
  /**
   * NewRecordUpsert creates a new [RecordUpsert] form with initializer
   * config created from the provided [core.App] and [models.Record] instances
   * (for create you could pass a pointer to an empty Record - models.NewRecord(collection)).
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, record: models.Record): (RecordUpsert | undefined)
 }
 interface RecordUpsert {
  /**
   * Data returns the loaded form's data.
   */
  data(): _TygojaDict
 }
 interface RecordUpsert {
  /**
   * SetFullManageAccess sets the manageAccess bool flag of the current
   * form to enable/disable directly changing some system record fields
   * (often used with auth collection records).
   */
  setFullManageAccess(fullManageAccess: boolean): void
 }
 interface RecordUpsert {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordUpsert {
  /**
   * LoadRequest extracts the json or multipart/form-data request data
   * and lods it into the form.
   * 
   * File upload is supported only via multipart/form-data.
   */
  loadRequest(r: http.Request, keyPrefix: string): void
 }
 interface RecordUpsert {
  /**
   * FilesToUpload returns the parsed request files ready for upload.
   */
  filesToUpload(): _TygojaDict
 }
 interface RecordUpsert {
  /**
   * FilesToUpload returns the parsed request filenames ready to be deleted.
   */
  filesToDelete(): Array<string>
 }
 interface RecordUpsert {
  /**
   * AddFiles adds the provided file(s) to the specified file field.
   * 
   * If the file field is a SINGLE-value file field (aka. "Max Select = 1"),
   * then the newly added file will REPLACE the existing one.
   * In this case if you pass more than 1 files only the first one will be assigned.
   * 
   * If the file field is a MULTI-value file field (aka. "Max Select > 1"),
   * then the newly added file(s) will be APPENDED to the existing one(s).
   * 
   * Example
   * 
   * ```
   * 	f1, _ := filesystem.NewFileFromPath("/path/to/file1.txt")
   * 	f2, _ := filesystem.NewFileFromPath("/path/to/file2.txt")
   * 	form.AddFiles("documents", f1, f2)
   * ```
   */
  addFiles(key: string, ...files: (filesystem.File | undefined)[]): void
 }
 interface RecordUpsert {
  /**
   * RemoveFiles removes a single or multiple file from the specified file field.
   * 
   * NB! If filesToDelete is not set it will remove all existing files
   * assigned to the file field (including those assigned with AddFiles)!
   * 
   * Example
   * 
   * ```
   * 	// mark only only 2 files for removal
   * 	form.AddFiles("documents", "file1_aw4bdrvws6.txt", "file2_xwbs36bafv.txt")
   * 
   * 	// mark all "documents" files for removal
   * 	form.AddFiles("documents")
   * ```
   */
  removeFiles(key: string, ...toDelete: string[]): void
 }
 interface RecordUpsert {
  /**
   * LoadData loads and normalizes the provided regular record data fields into the form.
   */
  loadData(requestInfo: _TygojaDict): void
 }
 interface RecordUpsert {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordUpsert {
  validateAndFill(): void
 }
 interface RecordUpsert {
  /**
   * DrySubmit performs a form submit within a transaction and reverts it.
   * For actual record persistence, check the `form.Submit()` method.
   * 
   * This method doesn't handle file uploads/deletes or trigger any app events!
   */
  drySubmit(callback: (txDao: daos.Dao) => void): void
 }
 interface RecordUpsert {
  /**
   * Submit validates the form and upserts the form Record model.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): void
 }
 /**
  * RecordVerificationConfirm is an auth record email verification confirmation form.
  */
 interface RecordVerificationConfirm {
  token: string
 }
 interface newRecordVerificationConfirm {
  /**
   * NewRecordVerificationConfirm creates a new [RecordVerificationConfirm]
   * form initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (RecordVerificationConfirm | undefined)
 }
 interface RecordVerificationConfirm {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordVerificationConfirm {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface RecordVerificationConfirm {
  /**
   * Submit validates and submits the form.
   * On success returns the verified auth record associated to `form.Token`.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): (models.Record | undefined)
 }
 /**
  * RecordVerificationRequest is an auth record email verification request form.
  */
 interface RecordVerificationRequest {
  email: string
 }
 interface newRecordVerificationRequest {
  /**
   * NewRecordVerificationRequest creates a new [RecordVerificationRequest]
   * form initialized with from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App, collection: models.Collection): (RecordVerificationRequest | undefined)
 }
 interface RecordVerificationRequest {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface RecordVerificationRequest {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   * 
   * // This method doesn't verify that auth record with `form.Email` exists (this is done on Submit).
   */
  validate(): void
 }
 interface RecordVerificationRequest {
  /**
   * Submit validates and sends a verification request email
   * to the `form.Email` auth record.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<models.Record | undefined>[]): void
 }
 /**
  * SettingsUpsert is a [settings.Settings] upsert (create/update) form.
  */
 type _subBpyXp = settings.Settings
 interface SettingsUpsert extends _subBpyXp {
 }
 interface newSettingsUpsert {
  /**
   * NewSettingsUpsert creates a new [SettingsUpsert] form with initializer
   * config created from the provided [core.App] instance.
   * 
   * If you want to submit the form as part of a transaction,
   * you can change the default Dao via [SetDao()].
   */
  (app: core.App): (SettingsUpsert | undefined)
 }
 interface SettingsUpsert {
  /**
   * SetDao replaces the default form Dao instance with the provided one.
   */
  setDao(dao: daos.Dao): void
 }
 interface SettingsUpsert {
  /**
   * Validate makes the form validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface SettingsUpsert {
  /**
   * Submit validates the form and upserts the loaded settings.
   * 
   * On success the app settings will be refreshed with the form ones.
   * 
   * You can optionally provide a list of InterceptorFunc to further
   * modify the form behavior before persisting it.
   */
  submit(...interceptors: InterceptorFunc<settings.Settings | undefined>[]): void
 }
 /**
  * TestEmailSend is a email template test request form.
  */
 interface TestEmailSend {
  template: string
  email: string
 }
 interface newTestEmailSend {
  /**
   * NewTestEmailSend creates and initializes new TestEmailSend form.
   */
  (app: core.App): (TestEmailSend | undefined)
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
  (app: core.App): (TestS3Filesystem | undefined)
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

/**
 * Package apis implements the default PocketBase api services and middlewares.
 */
namespace apis {
 interface adminApi {
 }
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * ApiError defines the struct for a basic api error response.
  */
 interface ApiError {
  code: number
  message: string
  data: _TygojaDict
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
 interface newNotFoundError {
  /**
   * NewNotFoundError creates and returns 404 `ApiError`.
   */
  (message: string, data: any): (ApiError | undefined)
 }
 interface newBadRequestError {
  /**
   * NewBadRequestError creates and returns 400 `ApiError`.
   */
  (message: string, data: any): (ApiError | undefined)
 }
 interface newForbiddenError {
  /**
   * NewForbiddenError creates and returns 403 `ApiError`.
   */
  (message: string, data: any): (ApiError | undefined)
 }
 interface newUnauthorizedError {
  /**
   * NewUnauthorizedError creates and returns 401 `ApiError`.
   */
  (message: string, data: any): (ApiError | undefined)
 }
 interface newApiError {
  /**
   * NewApiError creates and returns new normalized `ApiError` instance.
   */
  (status: number, message: string, data: any): (ApiError | undefined)
 }
 interface backupApi {
 }
 interface initApi {
  /**
   * InitApi creates a configured echo instance with registered
   * system and app specific routes and middlewares.
   */
  (app: core.App): (echo.Echo | undefined)
 }
 interface staticDirectoryHandler {
  /**
   * StaticDirectoryHandler is similar to `echo.StaticDirectoryHandler`
   * but without the directory redirect which conflicts with RemoveTrailingSlash middleware.
   * 
   * If a file resource is missing and indexFallback is set, the request
   * will be forwarded to the base index.html (useful also for SPA).
   * 
   * @see https://github.com/labstack/echo/issues/2211
   */
  (fileSystem: fs.FS, indexFallback: boolean): echo.HandlerFunc
 }
 interface collectionApi {
 }
 interface fileApi {
 }
 interface healthApi {
 }
 interface healthCheckResponse {
  code: number
  message: string
  data: {
   canBackup: boolean
  }
 }
 interface logsApi {
 }
 interface requireGuestOnly {
  /**
   * RequireGuestOnly middleware requires a request to NOT have a valid
   * Authorization header.
   * 
   * This middleware is the opposite of [apis.RequireAdminOrRecordAuth()].
   */
  (): echo.MiddlewareFunc
 }
 interface requireRecordAuth {
  /**
   * RequireRecordAuth middleware requires a request to have
   * a valid record auth Authorization header.
   * 
   * The auth record could be from any collection.
   * 
   * You can further filter the allowed record auth collections by
   * specifying their names.
   * 
   * Example:
   * 
   * ```
   * 	apis.RequireRecordAuth()
   * ```
   * 
   * Or:
   * 
   * ```
   * 	apis.RequireRecordAuth("users", "supervisors")
   * ```
   * 
   * To restrict the auth record only to the loaded context collection,
   * use [apis.RequireSameContextRecordAuth()] instead.
   */
  (...optCollectionNames: string[]): echo.MiddlewareFunc
 }
 interface requireSameContextRecordAuth {
  /**
   * RequireSameContextRecordAuth middleware requires a request to have
   * a valid record Authorization header.
   * 
   * The auth record must be from the same collection already loaded in the context.
   */
  (): echo.MiddlewareFunc
 }
 interface requireAdminAuth {
  /**
   * RequireAdminAuth middleware requires a request to have
   * a valid admin Authorization header.
   */
  (): echo.MiddlewareFunc
 }
 interface requireAdminAuthOnlyIfAny {
  /**
   * RequireAdminAuthOnlyIfAny middleware requires a request to have
   * a valid admin Authorization header ONLY if the application has
   * at least 1 existing Admin model.
   */
  (app: core.App): echo.MiddlewareFunc
 }
 interface requireAdminOrRecordAuth {
  /**
   * RequireAdminOrRecordAuth middleware requires a request to have
   * a valid admin or record Authorization header set.
   * 
   * You can further filter the allowed auth record collections by providing their names.
   * 
   * This middleware is the opposite of [apis.RequireGuestOnly()].
   */
  (...optCollectionNames: string[]): echo.MiddlewareFunc
 }
 interface requireAdminOrOwnerAuth {
  /**
   * RequireAdminOrOwnerAuth middleware requires a request to have
   * a valid admin or auth record owner Authorization header set.
   * 
   * This middleware is similar to [apis.RequireAdminOrRecordAuth()] but
   * for the auth record token expects to have the same id as the path
   * parameter ownerIdParam (default to "id" if empty).
   */
  (ownerIdParam: string): echo.MiddlewareFunc
 }
 interface loadAuthContext {
  /**
   * LoadAuthContext middleware reads the Authorization request header
   * and loads the token related record or admin instance into the
   * request's context.
   * 
   * This middleware is expected to be already registered by default for all routes.
   */
  (app: core.App): echo.MiddlewareFunc
 }
 interface loadCollectionContext {
  /**
   * LoadCollectionContext middleware finds the collection with related
   * path identifier and loads it into the request context.
   * 
   * Set optCollectionTypes to further filter the found collection by its type.
   */
  (app: core.App, ...optCollectionTypes: string[]): echo.MiddlewareFunc
 }
 interface activityLogger {
  /**
   * ActivityLogger middleware takes care to save the request information
   * into the logs database.
   * 
   * The middleware does nothing if the app logs retention period is zero
   * (aka. app.Settings().Logs.MaxDays = 0).
   */
  (app: core.App): echo.MiddlewareFunc
 }
 interface realtimeApi {
 }
 interface recordData {
  action: string
  record?: models.Record
 }
 interface getter {
  get(_arg0: string): any
 }
 interface recordAuthApi {
 }
 interface providerInfo {
  name: string
  state: string
  codeVerifier: string
  codeChallenge: string
  codeChallengeMethod: string
  authUrl: string
 }
 interface recordApi {
 }
 interface requestData {
  /**
   * Deprecated: Use RequestInfo instead.
   */
  (c: echo.Context): (models.RequestInfo | undefined)
 }
 interface requestInfo {
  /**
   * RequestInfo exports cached common request data fields
   * (query, body, logged auth state, etc.) from the provided context.
   */
  (c: echo.Context): (models.RequestInfo | undefined)
 }
 interface recordAuthResponse {
  /**
   * RecordAuthResponse writes standardised json record auth response
   * into the specified request context.
   */
  (app: core.App, c: echo.Context, authRecord: models.Record, meta: any, ...finalizers: ((token: string) => void)[]): void
 }
 interface enrichRecord {
  /**
   * EnrichRecord parses the request context and enrich the provided record:
   * ```
   *   - expands relations (if defaultExpands and/or ?expand query param is set)
   *   - ensures that the emails of the auth record and its expanded auth relations
   *     are visibe only for the current logged admin, record owner or record with manage access
   * ```
   */
  (c: echo.Context, dao: daos.Dao, record: models.Record, ...defaultExpands: string[]): void
 }
 interface enrichRecords {
  /**
   * EnrichRecords parses the request context and enriches the provided records:
   * ```
   *   - expands relations (if defaultExpands and/or ?expand query param is set)
   *   - ensures that the emails of the auth records and their expanded auth relations
   *     are visibe only for the current logged admin, record owner or record with manage access
   * ```
   */
  (c: echo.Context, dao: daos.Dao, records: Array<(models.Record | undefined)>, ...defaultExpands: string[]): void
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
   * HttpAddr is the TCP address to listen for the HTTP server (eg. `127.0.0.1:80`).
   */
  httpAddr: string
  /**
   * HttpsAddr is the TCP address to listen for the HTTPS server (eg. `127.0.0.1:443`).
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
  (app: core.App, config: ServeConfig): (http.Server | undefined)
 }
 interface migrationsConnection {
  db?: dbx.DB
  migrationsList: migrate.MigrationsList
 }
 interface settingsApi {
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
   * NewRegistry creates and initializes a new blank templates registry.
   * 
   * Use the Registry.Load* methods to load templates into the registry.
   */
  (): (Registry | undefined)
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
   * LoadFiles caches (if not already) the specified filenames set as a
   * single template and returns a ready to use Renderer instance.
   * 
   * There must be at least 1 filename specified.
   */
  loadFiles(...filenames: string[]): (Renderer | undefined)
 }
 interface Registry {
  /**
   * LoadString caches (if not already) the specified inline string as a
   * single template and returns a ready to use Renderer instance.
   */
  loadString(text: string): (Renderer | undefined)
 }
 interface Registry {
  /**
   * LoadString caches (if not already) the specified fs and globPatterns
   * pair as single template and returns a ready to use Renderer instance.
   * 
   * There must be at least 1 file matching the provided globPattern(s)
   * (note that most file names serves as glob patterns matching themselves).
   */
  loadFS(fs: fs.FS, ...globPatterns: string[]): (Renderer | undefined)
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

namespace pocketbase {
 /**
  * appWrapper serves as a private core.App instance wrapper.
  */
 type _subSsZSM = core.App
 interface appWrapper extends _subSsZSM {
 }
 /**
  * PocketBase defines a PocketBase app launcher.
  * 
  * It implements [core.App] via embedding and all of the app interface methods
  * could be accessed directly through the instance (eg. PocketBase.DataDir()).
  */
 type _subtNSbs = appWrapper
 interface PocketBase extends _subtNSbs {
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
   * optional default values for the console flags
   */
  defaultDebug: boolean
  defaultDataDir: string // if not set, it will fallback to "./pb_data"
  defaultEncryptionEnv: string
  /**
   * hide the default console server info on app startup
   */
  hideStartBanner: boolean
  /**
   * optional DB configurations
   */
  dataMaxOpenConns: number // default to core.DefaultDataMaxOpenConns
  dataMaxIdleConns: number // default to core.DefaultDataMaxIdleConns
  logsMaxOpenConns: number // default to core.DefaultLogsMaxOpenConns
  logsMaxIdleConns: number // default to core.DefaultLogsMaxIdleConns
 }
 interface _new {
  /**
   * New creates a new PocketBase instance with the default configuration.
   * Use [NewWithConfig()] if you want to provide a custom configuration.
   * 
   * Note that the application will not be initialized/bootstrapped yet,
   * aka. DB connections, migrations, app settings, etc. will not be accessible.
   * Everything will be initialized when [Start()] is executed.
   * If you want to initialize the application before calling [Start()],
   * then you'll have to manually call [Bootstrap()].
   */
  (): (PocketBase | undefined)
 }
 interface newWithConfig {
  /**
   * NewWithConfig creates a new PocketBase instance with the provided config.
   * 
   * Note that the application will not be initialized/bootstrapped yet,
   * aka. DB connections, migrations, app settings, etc. will not be accessible.
   * Everything will be initialized when [Start()] is executed.
   * If you want to initialize the application before calling [Start()],
   * then you'll have to manually call [Bootstrap()].
   */
  (config: Config): (PocketBase | undefined)
 }
 interface PocketBase {
  /**
   * Start starts the application, aka. registers the default system
   * commands (serve, migrate, version) and executes pb.RootCmd.
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
  * Implementations of Read are discouraged from returning a
  * zero byte count with a nil error, except when len(p) == 0.
  * Callers should treat a return of 0 and nil as indicating that
  * nothing happened; in particular it does not indicate EOF.
  * 
  * Implementations must not retain p.
  */
 interface Reader {
  read(p: string): number
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
  write(p: string): number
 }
 /**
  * ReadSeekCloser is the interface that groups the basic Read, Seek and Close
  * methods.
  */
 interface ReadSeekCloser {
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
 * On most systems, that error has type syscall.Errno.
 * 
 * Deprecated: this package is locked down. Callers should use the
 * corresponding package in the golang.org/x/sys repository instead.
 * That is also where updates required by new systems or versions
 * should be applied. See https://golang.org/s/go1.4-syscall for more
 * information.
 */
namespace syscall {
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
  noctty: boolean // Detach fd 0 from controlling terminal
  ctty: number // Controlling TTY fd
  /**
   * Foreground places the child process group in the foreground.
   * This implies Setpgid. The Ctty field must be set to
   * the descriptor of the controlling TTY.
   * Unlike Setctty, in this case Ctty must be a descriptor
   * number in the parent process.
   */
  foreground: boolean
  pgid: number // Child's process group ID if Setpgid.
  pdeathsig: Signal // Signal that the process will get when its parent dies (Linux and FreeBSD only)
  cloneflags: number // Flags for clone calls (Linux only)
  unshareflags: number // Flags for unshare calls (Linux only)
  uidMappings: Array<SysProcIDMap> // User ID mappings for user namespaces.
  gidMappings: Array<SysProcIDMap> // Group ID mappings for user namespaces.
  /**
   * GidMappingsEnableSetgroups enabling setgroups syscall.
   * If false, then setgroups syscall will be disabled for the child process.
   * This parameter is no-op if GidMappings == nil. Otherwise for unprivileged
   * users this should be set to false for mappings work.
   */
  gidMappingsEnableSetgroups: boolean
  ambientCaps: Array<number> // Ambient capabilities (Linux only)
 }
 // @ts-ignore
 import errorspkg = errors
 /**
  * A RawConn is a raw network connection.
  */
 interface RawConn {
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
 /**
  * An Errno is an unsigned number describing an error condition.
  * It implements the error interface. The zero Errno is by convention
  * a non-error, so code to convert from Errno to error should use:
  * ```
  * 	err = nil
  * 	if errno != 0 {
  * 		err = errno
  * 	}
  * ```
  * 
  * Errno values can be tested against error values from the os package
  * using errors.Is. For example:
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
 * Package time provides functionality for measuring and displaying time.
 * 
 * The calendrical calculations always assume a Gregorian calendar, with
 * no leap seconds.
 * 
 * Monotonic Clocks
 * 
 * Operating systems provide both a “wall clock,” which is subject to
 * changes for clock synchronization, and a “monotonic clock,” which is
 * not. The general rule is that the wall clock is for telling time and
 * the monotonic clock is for measuring time. Rather than split the API,
 * in this package the Time returned by time.Now contains both a wall
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
 * Other idioms, such as time.Since(start), time.Until(deadline), and
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
 * t.After(u), t.Before(u), t.Equal(u), and t.Sub(u) are carried out
 * using the monotonic clock readings alone, ignoring the wall clock
 * readings. If either t or u contains no monotonic clock reading, these
 * operations fall back to using the wall clock readings.
 * 
 * On some systems the monotonic clock will stop if the computer goes to sleep.
 * On such a system, t.Sub(u) may not accurately reflect the actual
 * time that passed between t and u.
 * 
 * Because the monotonic clock reading has no meaning outside
 * the current process, the serialized forms generated by t.GobEncode,
 * t.MarshalBinary, t.MarshalJSON, and t.MarshalText omit the monotonic
 * clock reading, and t.Format provides no format for it. Similarly, the
 * constructors time.Date, time.Parse, time.ParseInLocation, and time.Unix,
 * as well as the unmarshalers t.GobDecode, t.UnmarshalBinary.
 * t.UnmarshalJSON, and t.UnmarshalText always create times with
 * no monotonic clock reading.
 * 
 * Note that the Go == operator compares not just the time instant but
 * also the Location and the monotonic clock reading. See the
 * documentation for the Time type for a discussion of equality
 * testing for Time values.
 * 
 * For debugging, the result of t.String does include the monotonic
 * clock reading if present. If t != u because of different monotonic clock readings,
 * that difference will be visible when printing t.String() and u.String().
 */
namespace time {
 interface Time {
  /**
   * String returns the time formatted using the format string
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
   * GoString implements fmt.GoStringer and formats t to be printed in Go source
   * code.
   */
  goString(): string
 }
 interface Time {
  /**
   * Format returns a textual representation of the time value formatted according
   * to the layout defined by the argument. See the documentation for the
   * constant called Layout to see how to represent the layout format.
   * 
   * The executable example for Time.Format demonstrates the working
   * of the layout string in detail and is a good reference.
   */
  format(layout: string): string
 }
 interface Time {
  /**
   * AppendFormat is like Format but appends the textual
   * representation to b and returns the extended buffer.
   */
  appendFormat(b: string, layout: string): string
 }
 /**
  * A Time represents an instant in time with nanosecond precision.
  * 
  * Programs using times should typically store and pass them as values,
  * not pointers. That is, time variables and struct fields should be of
  * type time.Time, not *time.Time.
  * 
  * A Time value can be used by multiple goroutines simultaneously except
  * that the methods GobDecode, UnmarshalBinary, UnmarshalJSON and
  * UnmarshalText are not concurrency-safe.
  * 
  * Time instants can be compared using the Before, After, and Equal methods.
  * The Sub method subtracts two instants, producing a Duration.
  * The Add method adds a Time and a Duration, producing a Time.
  * 
  * The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC.
  * As this time is unlikely to come up in practice, the IsZero method gives
  * a simple way of detecting a time that has not been initialized explicitly.
  * 
  * Each Time has associated with it a Location, consulted when computing the
  * presentation form of the time, such as in the Format, Hour, and Year methods.
  * The methods Local, UTC, and In return a Time with a specific location.
  * Changing the location in this way changes only the presentation; it does not
  * change the instant in time being denoted and therefore does not affect the
  * computations described in earlier paragraphs.
  * 
  * Representations of a Time value saved by the GobEncode, MarshalBinary,
  * MarshalJSON, and MarshalText methods store the Time.Location's offset, but not
  * the location name. They therefore lose information about Daylight Saving Time.
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
   * IsZero reports whether t represents the zero time instant,
   * January 1, year 1, 00:00:00 UTC.
   */
  isZero(): boolean
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
  isoWeek(): number
 }
 interface Time {
  /**
   * Clock returns the hour, minute, and second within the day specified by t.
   */
  clock(): number
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
   * value that can be stored in a Duration,
   * Round returns the maximum (or minimum) duration.
   * If m <= 0, Round returns d unchanged.
   */
  round(m: Duration): Duration
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
   * value that can be stored in a Duration, the maximum (or minimum) duration
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
  location(): (Location | undefined)
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
   * MarshalBinary implements the encoding.BinaryMarshaler interface.
   */
  marshalBinary(): string
 }
 interface Time {
  /**
   * UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
   */
  unmarshalBinary(data: string): void
 }
 interface Time {
  /**
   * GobEncode implements the gob.GobEncoder interface.
   */
  gobEncode(): string
 }
 interface Time {
  /**
   * GobDecode implements the gob.GobDecoder interface.
   */
  gobDecode(data: string): void
 }
 interface Time {
  /**
   * MarshalJSON implements the json.Marshaler interface.
   * The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
   */
  marshalJSON(): string
 }
 interface Time {
  /**
   * UnmarshalJSON implements the json.Unmarshaler interface.
   * The time is expected to be a quoted string in RFC 3339 format.
   */
  unmarshalJSON(data: string): void
 }
 interface Time {
  /**
   * MarshalText implements the encoding.TextMarshaler interface.
   * The time is formatted in RFC 3339 format, with sub-second precision added if present.
   */
  marshalText(): string
 }
 interface Time {
  /**
   * UnmarshalText implements the encoding.TextUnmarshaler interface.
   * The time is expected to be in RFC 3339 format.
   */
  unmarshalText(data: string): void
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
 */
namespace fs {
 /**
  * An FS provides access to a hierarchical file system.
  * 
  * The FS interface is the minimum implementation required of the file system.
  * A file system may implement additional interfaces,
  * such as ReadFileFS, to provide additional or optimized functionality.
  */
 interface FS {
  /**
   * Open opens the named file.
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
  * Directory files should also implement ReadDirFile.
  * A file may implement io.ReaderAt or io.Seeker as optimizations.
  */
 interface File {
  stat(): FileInfo
  read(_arg0: string): number
  close(): void
 }
 /**
  * A DirEntry is an entry read from a directory
  * (using the ReadDir function or a ReadDirFile's ReadDir method).
  */
 interface DirEntry {
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
  * A FileInfo describes a file and is returned by Stat.
  */
 interface FileInfo {
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
  * The only required bit is ModeDir for directories.
  */
 interface FileMode extends Number{}
 interface FileMode {
  string(): string
 }
 interface FileMode {
  /**
   * IsDir reports whether m describes a directory.
   * That is, it tests for the ModeDir bit being set in m.
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
   * Perm returns the Unix permission bits in m (m & ModePerm).
   */
  perm(): FileMode
 }
 interface FileMode {
  /**
   * Type returns type bits in m (m & ModeType).
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
  * WalkDirFunc is the type of the function called by WalkDir to visit
  * each file or directory.
  * 
  * The path argument contains the argument to WalkDir as a prefix.
  * That is, if WalkDir is called with root argument "dir" and finds a file
  * named "a" in that directory, the walk function will be called with
  * argument "dir/a".
  * 
  * The d argument is the fs.DirEntry for the named path.
  * 
  * The error result returned by the function controls how WalkDir
  * continues. If the function returns the special value SkipDir, WalkDir
  * skips the current directory (path if d.IsDir() is true, otherwise
  * path's parent directory). Otherwise, if the function returns a non-nil
  * error, WalkDir stops entirely and returns that error.
  * 
  * The err argument reports an error related to path, signaling that
  * WalkDir will not walk into that directory. The function can decide how
  * to handle that error; as described earlier, returning the error will
  * cause WalkDir to stop walking the entire tree.
  * 
  * WalkDir calls the function with a non-nil err argument in two cases.
  * 
  * First, if the initial fs.Stat on the root directory fails, WalkDir
  * calls the function with path set to root, d set to nil, and err set to
  * the error from fs.Stat.
  * 
  * Second, if a directory's ReadDir method fails, WalkDir calls the
  * function with path set to the directory's path, d set to an
  * fs.DirEntry describing the directory, and err set to the error from
  * ReadDir. In this second case, the function is called twice with the
  * path of the directory: the first call is before the directory read is
  * attempted and has err set to nil, giving the function a chance to
  * return SkipDir and avoid the ReadDir entirely. The second call is
  * after a failed ReadDir and reports the error from ReadDir.
  * (If ReadDir succeeds, there is no second call.)
  * 
  * The differences between WalkDirFunc compared to filepath.WalkFunc are:
  * 
  * ```
  *   - The second argument has type fs.DirEntry instead of fs.FileInfo.
  *   - The function is called before reading a directory, to allow SkipDir
  *     to bypass the directory read entirely.
  *   - If a directory read fails, the function is called a second time
  *     for that directory to report the error.
  * ```
  */
 interface WalkDirFunc {(path: string, d: DirEntry, err: Error): void }
}

/**
 * Package bytes implements functions for the manipulation of byte slices.
 * It is analogous to the facilities of the strings package.
 */
namespace bytes {
 /**
  * A Reader implements the io.Reader, io.ReaderAt, io.WriterTo, io.Seeker,
  * io.ByteScanner, and io.RuneScanner interfaces by reading from
  * a byte slice.
  * Unlike a Buffer, a Reader is read-only and supports seeking.
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
   * Size is the number of bytes available for reading via ReadAt.
   * The returned value is always the same and is not affected by calls
   * to any other method.
   */
  size(): number
 }
 interface Reader {
  /**
   * Read implements the io.Reader interface.
   */
  read(b: string): number
 }
 interface Reader {
  /**
   * ReadAt implements the io.ReaderAt interface.
   */
  readAt(b: string, off: number): number
 }
 interface Reader {
  /**
   * ReadByte implements the io.ByteReader interface.
   */
  readByte(): string
 }
 interface Reader {
  /**
   * UnreadByte complements ReadByte in implementing the io.ByteScanner interface.
   */
  unreadByte(): void
 }
 interface Reader {
  /**
   * ReadRune implements the io.RuneReader interface.
   */
  readRune(): [string, number]
 }
 interface Reader {
  /**
   * UnreadRune complements ReadRune in implementing the io.RuneScanner interface.
   */
  unreadRune(): void
 }
 interface Reader {
  /**
   * Seek implements the io.Seeker interface.
   */
  seek(offset: number, whence: number): number
 }
 interface Reader {
  /**
   * WriteTo implements the io.WriterTo interface.
   */
  writeTo(w: io.Writer): number
 }
 interface Reader {
  /**
   * Reset resets the Reader to be reading from b.
   */
  reset(b: string): void
 }
}

/**
 * Package context defines the Context type, which carries deadlines,
 * cancellation signals, and other request-scoped values across API boundaries
 * and between processes.
 * 
 * Incoming requests to a server should create a Context, and outgoing
 * calls to servers should accept a Context. The chain of function
 * calls between them must propagate the Context, optionally replacing
 * it with a derived Context created using WithCancel, WithDeadline,
 * WithTimeout, or WithValue. When a Context is canceled, all
 * Contexts derived from it are also canceled.
 * 
 * The WithCancel, WithDeadline, and WithTimeout functions take a
 * Context (the parent) and return a derived Context (the child) and a
 * CancelFunc. Calling the CancelFunc cancels the child and its
 * children, removes the parent's reference to the child, and stops
 * any associated timers. Failing to call the CancelFunc leaks the
 * child and its children until the parent is canceled or the timer
 * fires. The go vet tool checks that CancelFuncs are used on all
 * control-flow paths.
 * 
 * Programs that use Contexts should follow these rules to keep interfaces
 * consistent across packages and enable static analysis tools to check context
 * propagation:
 * 
 * Do not store Contexts inside a struct type; instead, pass a Context
 * explicitly to each function that needs it. The Context should be the first
 * parameter, typically named ctx:
 * 
 * ```
 * 	func DoSomething(ctx context.Context, arg Arg) error {
 * 		// ... use ctx ...
 * 	}
 * ```
 * 
 * Do not pass a nil Context, even if a function permits it. Pass context.TODO
 * if you are unsure about which Context to use.
 * 
 * Use context Values only for request-scoped data that transits processes and
 * APIs, not for passing optional parameters to functions.
 * 
 * The same Context may be passed to functions running in different goroutines;
 * Contexts are safe for simultaneous use by multiple goroutines.
 * 
 * See https://blog.golang.org/context for example code for a server that uses
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
   * Canceled if the context was canceled
   * or DeadlineExceeded if the context's deadline passed.
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
 * dangerous input, or use the path/filepath package's Glob function.
 * To expand environment variables, use package os's ExpandEnv.
 * 
 * Note that the examples in this package assume a Unix system.
 * They may not run on Windows, and they do not run in the Go Playground
 * used by golang.org and godoc.org.
 */
namespace exec {
 /**
  * Cmd represents an external command being prepared or run.
  * 
  * A Cmd cannot be reused after calling its Run, Output or CombinedOutput
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
   */
  env: Array<string>
  /**
   * Dir specifies the working directory of the command.
   * If Dir is the empty string, Run runs the command in the
   * calling process's current directory.
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
   * (EOF or a read error) or because writing to the pipe returned an error.
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
   * goroutine reaches EOF or encounters an error.
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
   * ProcessState contains information about an exited process,
   * available after a call to Wait or Run.
   */
  processState?: os.ProcessState
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
   * type *ExitError. Other error types may be returned for other situations.
   * 
   * If the calling goroutine has locked the operating system thread
   * with runtime.LockOSThread and modified any inheritable OS-level
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
   * The Wait method will return the exit code and release associated resources
   * once the command exits.
   */
  start(): void
 }
 interface Cmd {
  /**
   * Wait waits for the command to exit and waits for any copying to
   * stdin or copying from stdout or stderr to complete.
   * 
   * The command must have been started by Start.
   * 
   * The returned error is nil if the command runs, has no problems
   * copying stdin, stdout, and stderr, and exits with a zero exit
   * status.
   * 
   * If the command fails to run or doesn't complete successfully, the
   * error is of type *ExitError. Other error types may be
   * returned for I/O problems.
   * 
   * If any of c.Stdin, c.Stdout or c.Stderr are not an *os.File, Wait also waits
   * for the respective I/O loop copying to or from the process to complete.
   * 
   * Wait releases any resources associated with the Cmd.
   */
  wait(): void
 }
 interface Cmd {
  /**
   * Output runs the command and returns its standard output.
   * Any returned error will usually be of type *ExitError.
   * If c.Stderr was nil, Output populates ExitError.Stderr.
   */
  output(): string
 }
 interface Cmd {
  /**
   * CombinedOutput runs the command and returns its combined standard
   * output and standard error.
   */
  combinedOutput(): string
 }
 interface Cmd {
  /**
   * StdinPipe returns a pipe that will be connected to the command's
   * standard input when the command starts.
   * The pipe will be closed automatically after Wait sees the command exit.
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
   * Wait will close the pipe after seeing the command exit, so most callers
   * need not close the pipe themselves. It is thus incorrect to call Wait
   * before all reads from the pipe have completed.
   * For the same reason, it is incorrect to call Run when using StdoutPipe.
   * See the example for idiomatic usage.
   */
  stdoutPipe(): io.ReadCloser
 }
 interface Cmd {
  /**
   * StderrPipe returns a pipe that will be connected to the command's
   * standard error when the command starts.
   * 
   * Wait will close the pipe after seeing the command exit, so most callers
   * need not close the pipe themselves. It is thus incorrect to call Wait
   * before all reads from the pipe have completed.
   * For the same reason, it is incorrect to use Run when using StderrPipe.
   * See the StdoutPipe example for idiomatic usage.
   */
  stderrPipe(): io.ReadCloser
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
  * TxOptions holds the transaction options to be used in DB.BeginTx.
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
  * DB is a database handle representing a pool of zero or more
  * underlying connections. It's safe for concurrent use by multiple
  * goroutines.
  * 
  * The sql package creates and frees connections automatically; it
  * also maintains a free pool of idle connections. If the database has
  * a concept of per-connection state, such state can be reliably observed
  * within a transaction (Tx) or connection (Conn). Once DB.Begin is called, the
  * returned Tx is bound to a single connection. Once Commit or
  * Rollback is called on the transaction, that transaction's
  * connection is returned to DB's idle connection pool. The pool size
  * can be controlled with SetMaxIdleConns.
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
   * Ping uses context.Background internally; to specify the context, use
   * PingContext.
   */
  ping(): void
 }
 interface DB {
  /**
   * Close closes the database and prevents new queries from starting.
   * Close then waits for all queries that have started processing on the server
   * to finish.
   * 
   * It is rare to Close a DB, as the DB handle is meant to be
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
   * The caller must call the statement's Close method
   * when the statement is no longer needed.
   * 
   * The provided context is used for the preparation of the statement, not for the
   * execution of the statement.
   */
  prepareContext(ctx: context.Context, query: string): (Stmt | undefined)
 }
 interface DB {
  /**
   * Prepare creates a prepared statement for later queries or executions.
   * Multiple queries or executions may be run concurrently from the
   * returned statement.
   * The caller must call the statement's Close method
   * when the statement is no longer needed.
   * 
   * Prepare uses context.Background internally; to specify the context, use
   * PrepareContext.
   */
  prepare(query: string): (Stmt | undefined)
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
   * Exec uses context.Background internally; to specify the context, use
   * ExecContext.
   */
  exec(query: string, ...args: any[]): Result
 }
 interface DB {
  /**
   * QueryContext executes a query that returns rows, typically a SELECT.
   * The args are for any placeholder parameters in the query.
   */
  queryContext(ctx: context.Context, query: string, ...args: any[]): (Rows | undefined)
 }
 interface DB {
  /**
   * Query executes a query that returns rows, typically a SELECT.
   * The args are for any placeholder parameters in the query.
   * 
   * Query uses context.Background internally; to specify the context, use
   * QueryContext.
   */
  query(query: string, ...args: any[]): (Rows | undefined)
 }
 interface DB {
  /**
   * QueryRowContext executes a query that is expected to return at most one row.
   * QueryRowContext always returns a non-nil value. Errors are deferred until
   * Row's Scan method is called.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, query: string, ...args: any[]): (Row | undefined)
 }
 interface DB {
  /**
   * QueryRow executes a query that is expected to return at most one row.
   * QueryRow always returns a non-nil value. Errors are deferred until
   * Row's Scan method is called.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   * 
   * QueryRow uses context.Background internally; to specify the context, use
   * QueryRowContext.
   */
  queryRow(query: string, ...args: any[]): (Row | undefined)
 }
 interface DB {
  /**
   * BeginTx starts a transaction.
   * 
   * The provided context is used until the transaction is committed or rolled back.
   * If the context is canceled, the sql package will roll back
   * the transaction. Tx.Commit will return an error if the context provided to
   * BeginTx is canceled.
   * 
   * The provided TxOptions is optional and may be nil if defaults should be used.
   * If a non-default isolation level is used that the driver doesn't support,
   * an error will be returned.
   */
  beginTx(ctx: context.Context, opts: TxOptions): (Tx | undefined)
 }
 interface DB {
  /**
   * Begin starts a transaction. The default isolation level is dependent on
   * the driver.
   * 
   * Begin uses context.Background internally; to specify the context, use
   * BeginTx.
   */
  begin(): (Tx | undefined)
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
   * calling Conn.Close.
   */
  conn(ctx: context.Context): (Conn | undefined)
 }
 /**
  * Tx is an in-progress database transaction.
  * 
  * A transaction must end with a call to Commit or Rollback.
  * 
  * After a call to Commit or Rollback, all operations on the
  * transaction fail with ErrTxDone.
  * 
  * The statements prepared for a transaction by calling
  * the transaction's Prepare or Stmt methods are closed
  * by the call to Commit or Rollback.
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
   * To use an existing prepared statement on this transaction, see Tx.Stmt.
   * 
   * The provided context will be used for the preparation of the context, not
   * for the execution of the returned statement. The returned statement
   * will run in the transaction context.
   */
  prepareContext(ctx: context.Context, query: string): (Stmt | undefined)
 }
 interface Tx {
  /**
   * Prepare creates a prepared statement for use within a transaction.
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   * 
   * To use an existing prepared statement on this transaction, see Tx.Stmt.
   * 
   * Prepare uses context.Background internally; to specify the context, use
   * PrepareContext.
   */
  prepare(query: string): (Stmt | undefined)
 }
 interface Tx {
  /**
   * StmtContext returns a transaction-specific prepared statement from
   * an existing statement.
   * 
   * Example:
   *  updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
   *  ...
   *  tx, err := db.Begin()
   *  ...
   *  res, err := tx.StmtContext(ctx, updateMoney).Exec(123.45, 98293203)
   * 
   * The provided context is used for the preparation of the statement, not for the
   * execution of the statement.
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   */
  stmtContext(ctx: context.Context, stmt: Stmt): (Stmt | undefined)
 }
 interface Tx {
  /**
   * Stmt returns a transaction-specific prepared statement from
   * an existing statement.
   * 
   * Example:
   *  updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
   *  ...
   *  tx, err := db.Begin()
   *  ...
   *  res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
   * 
   * The returned statement operates within the transaction and will be closed
   * when the transaction has been committed or rolled back.
   * 
   * Stmt uses context.Background internally; to specify the context, use
   * StmtContext.
   */
  stmt(stmt: Stmt): (Stmt | undefined)
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
   * Exec uses context.Background internally; to specify the context, use
   * ExecContext.
   */
  exec(query: string, ...args: any[]): Result
 }
 interface Tx {
  /**
   * QueryContext executes a query that returns rows, typically a SELECT.
   */
  queryContext(ctx: context.Context, query: string, ...args: any[]): (Rows | undefined)
 }
 interface Tx {
  /**
   * Query executes a query that returns rows, typically a SELECT.
   * 
   * Query uses context.Background internally; to specify the context, use
   * QueryContext.
   */
  query(query: string, ...args: any[]): (Rows | undefined)
 }
 interface Tx {
  /**
   * QueryRowContext executes a query that is expected to return at most one row.
   * QueryRowContext always returns a non-nil value. Errors are deferred until
   * Row's Scan method is called.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, query: string, ...args: any[]): (Row | undefined)
 }
 interface Tx {
  /**
   * QueryRow executes a query that is expected to return at most one row.
   * QueryRow always returns a non-nil value. Errors are deferred until
   * Row's Scan method is called.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   * 
   * QueryRow uses context.Background internally; to specify the context, use
   * QueryRowContext.
   */
  queryRow(query: string, ...args: any[]): (Row | undefined)
 }
 /**
  * Stmt is a prepared statement.
  * A Stmt is safe for concurrent use by multiple goroutines.
  * 
  * If a Stmt is prepared on a Tx or Conn, it will be bound to a single
  * underlying connection forever. If the Tx or Conn closes, the Stmt will
  * become unusable and all operations will return an error.
  * If a Stmt is prepared on a DB, it will remain usable for the lifetime of the
  * DB. When the Stmt needs to execute on a new underlying connection, it will
  * prepare itself on the new connection automatically.
  */
 interface Stmt {
 }
 interface Stmt {
  /**
   * ExecContext executes a prepared statement with the given arguments and
   * returns a Result summarizing the effect of the statement.
   */
  execContext(ctx: context.Context, ...args: any[]): Result
 }
 interface Stmt {
  /**
   * Exec executes a prepared statement with the given arguments and
   * returns a Result summarizing the effect of the statement.
   * 
   * Exec uses context.Background internally; to specify the context, use
   * ExecContext.
   */
  exec(...args: any[]): Result
 }
 interface Stmt {
  /**
   * QueryContext executes a prepared query statement with the given arguments
   * and returns the query results as a *Rows.
   */
  queryContext(ctx: context.Context, ...args: any[]): (Rows | undefined)
 }
 interface Stmt {
  /**
   * Query executes a prepared query statement with the given arguments
   * and returns the query results as a *Rows.
   * 
   * Query uses context.Background internally; to specify the context, use
   * QueryContext.
   */
  query(...args: any[]): (Rows | undefined)
 }
 interface Stmt {
  /**
   * QueryRowContext executes a prepared query statement with the given arguments.
   * If an error occurs during the execution of the statement, that error will
   * be returned by a call to Scan on the returned *Row, which is always non-nil.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, ...args: any[]): (Row | undefined)
 }
 interface Stmt {
  /**
   * QueryRow executes a prepared query statement with the given arguments.
   * If an error occurs during the execution of the statement, that error will
   * be returned by a call to Scan on the returned *Row, which is always non-nil.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   * 
   * Example usage:
   * 
   *  var name string
   *  err := nameByUseridStmt.QueryRow(id).Scan(&name)
   * 
   * QueryRow uses context.Background internally; to specify the context, use
   * QueryRowContext.
   */
  queryRow(...args: any[]): (Row | undefined)
 }
 interface Stmt {
  /**
   * Close closes the statement.
   */
  close(): void
 }
 /**
  * Rows is the result of a query. Its cursor starts before the first row
  * of the result set. Use Next to advance from row to row.
  */
 interface Rows {
 }
 interface Rows {
  /**
   * Next prepares the next result row for reading with the Scan method. It
   * returns true on success, or false if there is no next result row or an error
   * happened while preparing it. Err should be consulted to distinguish between
   * the two cases.
   * 
   * Every call to Scan, even the first one, must be preceded by a call to Next.
   */
  next(): boolean
 }
 interface Rows {
  /**
   * NextResultSet prepares the next result set for reading. It reports whether
   * there is further result sets, or false if there is no further result set
   * or if there is an error advancing to it. The Err method should be consulted
   * to distinguish between the two cases.
   * 
   * After calling NextResultSet, the Next method should always be called before
   * scanning. If there are further result sets they may not have rows in the result
   * set.
   */
  nextResultSet(): boolean
 }
 interface Rows {
  /**
   * Err returns the error, if any, that was encountered during iteration.
   * Err may be called after an explicit or implicit Close.
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
   * number of columns in Rows.
   * 
   * Scan converts columns read from the database into the following
   * common Go types and special types provided by the sql package:
   * 
   * ```
   *    *string
   *    *[]byte
   *    *int, *int8, *int16, *int32, *int64
   *    *uint, *uint8, *uint16, *uint32, *uint64
   *    *bool
   *    *float32, *float64
   *    *interface{}
   *    *RawBytes
   *    *Rows (cursor value)
   *    any type implementing Scanner (see Scanner docs)
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
   * using an argument of type *RawBytes instead; see the documentation
   * for RawBytes for restrictions on its use.
   * 
   * If an argument has type *interface{}, Scan copies the value
   * provided by the underlying driver without conversion. When scanning
   * from a source value of type []byte to *interface{}, a copy of the
   * slice is made and the caller owns the result.
   * 
   * Source values of type time.Time may be scanned into values of type
   * *time.Time, *interface{}, *string, or *[]byte. When converting to
   * the latter two, time.RFC3339Nano is used.
   * 
   * Source values of type bool may be scanned into types *bool,
   * *interface{}, *string, *[]byte, or *RawBytes.
   * 
   * For scanning into *bool, the source may be true, false, 1, 0, or
   * string inputs parseable by strconv.ParseBool.
   * 
   * Scan can also convert a cursor returned from a query, such as
   * "select cursor(select * from my_table) from dual", into a
   * *Rows value that can itself be scanned from. The parent
   * select query will close any cursor *Rows if the parent *Rows is closed.
   * 
   * If any of the first arguments implementing Scanner returns an error,
   * that error will be wrapped in the returned error
   */
  scan(...dest: any[]): void
 }
 interface Rows {
  /**
   * Close closes the Rows, preventing further enumeration. If Next is called
   * and returns false and there are no further result sets,
   * the Rows are closed automatically and it will suffice to check the
   * result of Err. Close is idempotent and does not affect the result of Err.
   */
  close(): void
 }
 /**
  * A Result summarizes an executed SQL command.
  */
 interface Result {
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

namespace migrate {
 /**
  * MigrationsList defines a list with migration definitions
  */
 interface MigrationsList {
 }
 interface MigrationsList {
  /**
   * Item returns a single migration from the list by its index.
   */
  item(index: number): (Migration | undefined)
 }
 interface MigrationsList {
  /**
   * Items returns the internal migrations list slice.
   */
  items(): Array<(Migration | undefined)>
 }
 interface MigrationsList {
  /**
   * Register adds new migration definition to the list.
   * 
   * If `optFilename` is not provided, it will try to get the name from its .go file.
   * 
   * The list will be sorted automatically based on the migrations file name.
   */
  register(up: (db: dbx.Builder) => void, down: (db: dbx.Builder) => void, ...optFilename: string[]): void
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
  validArgs: Array<string>
  /**
   * ValidArgsFunction is an optional function that provides valid non-flag arguments for shell completion.
   * It is a dynamic version of using ValidArgs.
   * Only one of ValidArgs and ValidArgsFunction can be used for a command.
   */
  validArgsFunction: (cmd: Command, args: Array<string>, toComplete: string) => [Array<string>, ShellCompDirective]
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
   * group commands.
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
   */
  usageTemplate(): string
 }
 interface Command {
  /**
   * HelpTemplate return help template for the command.
   */
  helpTemplate(): string
 }
 interface Command {
  /**
   * VersionTemplate return version template for the command.
   */
  versionTemplate(): string
 }
 interface Command {
  /**
   * Find the target command given the args and command tree
   * Meant to be run on the highest node. Only searches down.
   */
  find(args: Array<string>): [(Command | undefined), Array<string>]
 }
 interface Command {
  /**
   * Traverse the command tree to find the command, and parse args for
   * each parent.
   */
  traverse(args: Array<string>): [(Command | undefined), Array<string>]
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
  root(): (Command | undefined)
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
  executeContextC(ctx: context.Context): (Command | undefined)
 }
 interface Command {
  /**
   * ExecuteC executes the command.
   */
  executeC(): (Command | undefined)
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
  flags(): (any | undefined)
 }
 interface Command {
  /**
   * LocalNonPersistentFlags are flags specific to this command which will NOT persist to subcommands.
   */
  localNonPersistentFlags(): (any | undefined)
 }
 interface Command {
  /**
   * LocalFlags returns the local FlagSet specifically set in the current command.
   */
  localFlags(): (any | undefined)
 }
 interface Command {
  /**
   * InheritedFlags returns all flags which were inherited from parent commands.
   */
  inheritedFlags(): (any | undefined)
 }
 interface Command {
  /**
   * NonInheritedFlags returns all flags which were not inherited from parent commands.
   */
  nonInheritedFlags(): (any | undefined)
 }
 interface Command {
  /**
   * PersistentFlags returns the persistent FlagSet specifically set in the current command.
   */
  persistentFlags(): (any | undefined)
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
  flag(name: string): (any | undefined)
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
  parent(): (Command | undefined)
 }
 interface Command {
  /**
   * RegisterFlagCompletionFunc should be called to register a function to provide completion for a flag.
   */
  registerFlagCompletionFunc(flagName: string, f: (cmd: Command, args: Array<string>, toComplete: string) => [Array<string>, ShellCompDirective]): void
 }
 interface Command {
  /**
   * InitDefaultCompletionCmd adds a default 'completion' command to c.
   * This function will do nothing if any of the following is true:
   * 1- the feature has been explicitly disabled by the program,
   * 2- c has no subcommands (to avoid creating one),
   * 3- c already has a 'completion' command provided by the program.
   */
  initDefaultCompletionCmd(): void
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
   * MarkFlagsMutuallyExclusive marks the given flags with annotations so that Cobra errors
   * if the command is invoked with more than one flag from the given set of flags.
   */
  markFlagsMutuallyExclusive(...flagNames: string[]): void
 }
 interface Command {
  /**
   * ValidateFlagGroups validates the mutuallyExclusive/requiredAsGroup logic and returns the
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

/**
 * Package jwt is a Go implementation of JSON Web Tokens: http://self-issued.info/docs/draft-jones-json-web-token.html
 * 
 * See README.md for more info.
 */
namespace jwt {
 /**
  * MapClaims is a claims type that uses the map[string]interface{} for JSON decoding.
  * This is the default claims type if you don't supply one
  */
 interface MapClaims extends _TygojaDict{}
 interface MapClaims {
  /**
   * VerifyAudience Compares the aud claim against cmp.
   * If required is false, this method will return true if the value matches or is unset
   */
  verifyAudience(cmp: string, req: boolean): boolean
 }
 interface MapClaims {
  /**
   * VerifyExpiresAt compares the exp claim against cmp (cmp <= exp).
   * If req is false, it will return true, if exp is unset.
   */
  verifyExpiresAt(cmp: number, req: boolean): boolean
 }
 interface MapClaims {
  /**
   * VerifyIssuedAt compares the exp claim against cmp (cmp >= iat).
   * If req is false, it will return true, if iat is unset.
   */
  verifyIssuedAt(cmp: number, req: boolean): boolean
 }
 interface MapClaims {
  /**
   * VerifyNotBefore compares the nbf claim against cmp (cmp >= nbf).
   * If req is false, it will return true, if nbf is unset.
   */
  verifyNotBefore(cmp: number, req: boolean): boolean
 }
 interface MapClaims {
  /**
   * VerifyIssuer compares the iss claim against cmp.
   * If required is false, this method will return true if the value matches or is unset
   */
  verifyIssuer(cmp: string, req: boolean): boolean
 }
 interface MapClaims {
  /**
   * Valid validates time based claims "exp, iat, nbf".
   * There is no accounting for clock skew.
   * As well, if any of the above claims are not in the token, it will still
   * be considered a valid claim.
   */
  valid(): void
 }
}

/**
 * Package multipart implements MIME multipart parsing, as defined in RFC
 * 2046.
 * 
 * The implementation is sufficient for HTTP (RFC 2388) and the multipart
 * bodies generated by popular browsers.
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
   * Open opens and returns the FileHeader's associated File.
   */
  open(): File
 }
}

/**
 * Package http provides HTTP client and server implementations.
 * 
 * Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:
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
 * The client must close the response body when finished with it:
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
 * For control over HTTP client headers, redirect policy, and other
 * settings, create a Client:
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
 * compression, and other settings, create a Transport:
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
 * ListenAndServe starts an HTTP server with a given address and handler.
 * The handler is usually nil, which means to use DefaultServeMux.
 * Handle and HandleFunc add handlers to DefaultServeMux:
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
 * Starting with Go 1.6, the http package has transparent support for the
 * HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2
 * can do so by setting Transport.TLSNextProto (for clients) or
 * Server.TLSNextProto (for servers) to a non-nil, empty
 * map. Alternatively, the following GODEBUG environment variables are
 * currently supported:
 * 
 * ```
 * 	GODEBUG=http2client=0  # disable HTTP/2 client support
 * 	GODEBUG=http2server=0  # disable HTTP/2 server support
 * 	GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
 * 	GODEBUG=http2debug=2   # ... even more verbose, with frame dumps
 * ```
 * 
 * The GODEBUG variables are not covered by Go's API compatibility
 * promise. Please report any issues before disabling HTTP/2
 * support: https://golang.org/s/http2bug
 * 
 * The http package's Transport and Server both automatically enable
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
 // @ts-ignore
 import urlpkg = url
 /**
  * A Request represents an HTTP request received by a server
  * or to be sent by a client.
  * 
  * The field semantics differ slightly between client and server
  * usage. In addition to the notes on the fields below, see the
  * documentation for Request.Write and RoundTripper.
  */
 interface Request {
  /**
   * Method specifies the HTTP method (GET, POST, PUT, etc.).
   * For client requests, an empty string means GET.
   * 
   * Go's HTTP client does not support sending a request with
   * the CONNECT method. See the documentation on Transport for
   * details.
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
 }
 interface Request {
  /**
   * Context returns the request's context. To change the context, use
   * WithContext.
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
   * To create a new request with a context, use NewRequestWithContext.
   * To change the context of a request, such as an incoming request you
   * want to modify before sending back out, use Request.Clone. Between
   * those two uses, it's rare to need WithContext.
   */
  withContext(ctx: context.Context): (Request | undefined)
 }
 interface Request {
  /**
   * Clone returns a deep copy of r with its context changed to ctx.
   * The provided ctx must be non-nil.
   * 
   * For an outgoing client request, the context controls the entire
   * lifetime of a request and its response: obtaining a connection,
   * sending the request, and reading the response headers and body.
   */
  clone(ctx: context.Context): (Request | undefined)
 }
 interface Request {
  /**
   * ProtoAtLeast reports whether the HTTP protocol used
   * in the request is at least major.minor.
   */
  protoAtLeast(major: number): boolean
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
   * Cookie returns the named cookie provided in the request or
   * ErrNoCookie if not found.
   * If multiple cookies match the given name, only one cookie will
   * be returned.
   */
  cookie(name: string): (Cookie | undefined)
 }
 interface Request {
  /**
   * AddCookie adds a cookie to the request. Per RFC 6265 section 5.4,
   * AddCookie does not attach more than one Cookie header field. That
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
   * Header map as Header["Referer"]; the benefit of making it available
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
   * Use this function instead of ParseMultipartForm to
   * process the request body as a stream.
   */
  multipartReader(): (multipart.Reader | undefined)
 }
 interface Request {
  /**
   * Write writes an HTTP/1.1 request, which is the header and body, in wire format.
   * This method consults the following fields of the request:
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
   * If Body is present, Content-Length is <= 0 and TransferEncoding
   * hasn't been set to "identity", Write adds "Transfer-Encoding:
   * chunked" to the header. Body is closed after it is sent.
   */
  write(w: io.Writer): void
 }
 interface Request {
  /**
   * WriteProxy is like Write but writes the request in the form
   * expected by an HTTP proxy. In particular, WriteProxy writes the
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
  basicAuth(): [string, boolean]
 }
 interface Request {
  /**
   * SetBasicAuth sets the request's Authorization header to use HTTP
   * Basic Authentication with the provided username and password.
   * 
   * With HTTP Basic Authentication the provided username and password
   * are not encrypted.
   * 
   * Some protocols may impose additional requirements on pre-escaping the
   * username and password. For instance, when used with OAuth2, both arguments
   * must be URL encoded first with url.QueryEscape.
   */
  setBasicAuth(username: string): void
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
   * If the request Body's size has not already been limited by MaxBytesReader,
   * the size is capped at 10MB.
   * 
   * For other HTTP methods, or when the Content-Type is not
   * application/x-www-form-urlencoded, the request Body is not read, and
   * r.PostForm is initialized to a non-nil, empty value.
   * 
   * ParseMultipartForm calls ParseForm automatically.
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
   * ParseMultipartForm calls ParseForm if necessary.
   * If ParseForm returns an error, ParseMultipartForm returns it but also
   * continues parsing the request body.
   * After one call to ParseMultipartForm, subsequent calls have no effect.
   */
  parseMultipartForm(maxMemory: number): void
 }
 interface Request {
  /**
   * FormValue returns the first value for the named component of the query.
   * POST and PUT body parameters take precedence over URL query string values.
   * FormValue calls ParseMultipartForm and ParseForm if necessary and ignores
   * any errors returned by these functions.
   * If key is not present, FormValue returns the empty string.
   * To access multiple values of the same key, call ParseForm and
   * then inspect Request.Form directly.
   */
  formValue(key: string): string
 }
 interface Request {
  /**
   * PostFormValue returns the first value for the named component of the POST,
   * PATCH, or PUT request body. URL query parameters are ignored.
   * PostFormValue calls ParseMultipartForm and ParseForm if necessary and ignores
   * any errors returned by these functions.
   * If key is not present, PostFormValue returns the empty string.
   */
  postFormValue(key: string): string
 }
 interface Request {
  /**
   * FormFile returns the first file for the provided form key.
   * FormFile calls ParseMultipartForm and ParseForm if necessary.
   */
  formFile(key: string): [multipart.File, (multipart.FileHeader | undefined)]
 }
 /**
  * A ResponseWriter interface is used by an HTTP handler to
  * construct an HTTP response.
  * 
  * A ResponseWriter may not be used after the Handler.ServeHTTP method
  * has returned.
  */
 interface ResponseWriter {
  /**
   * Header returns the header map that will be sent by
   * WriteHeader. The Header map also is the mechanism with which
   * Handlers can set HTTP trailers.
   * 
   * Changing the header map after a call to WriteHeader (or
   * Write) has no effect unless the modified headers are
   * trailers.
   * 
   * There are two ways to set Trailers. The preferred way is to
   * predeclare in the headers which trailers you will later
   * send by setting the "Trailer" header to the names of the
   * trailer keys which will come later. In this case, those
   * keys of the Header map are treated as if they were
   * trailers. See the example. The second way, for trailer
   * keys not known to the Handler until after the first Write,
   * is to prefix the Header map keys with the TrailerPrefix
   * constant value. See TrailerPrefix.
   * 
   * To suppress automatic response headers (such as "Date"), set
   * their value to nil.
   */
  header(): Header
  /**
   * Write writes the data to the connection as part of an HTTP reply.
   * 
   * If WriteHeader has not yet been called, Write calls
   * WriteHeader(http.StatusOK) before writing the data. If the Header
   * does not contain a Content-Type line, Write adds a Content-Type set
   * to the result of passing the initial 512 bytes of written data to
   * DetectContentType. Additionally, if the total size of all written
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
  write(_arg0: string): number
  /**
   * WriteHeader sends an HTTP response header with the provided
   * status code.
   * 
   * If WriteHeader is not called explicitly, the first call to Write
   * will trigger an implicit WriteHeader(http.StatusOK).
   * Thus explicit calls to WriteHeader are mainly used to
   * send error codes.
   * 
   * The provided code must be a valid HTTP 1xx-5xx status code.
   * Only one header may be written. Go does not currently
   * support sending user-defined 1xx informational headers,
   * with the exception of 100-continue response header that the
   * Server sends automatically when the Request.Body is read.
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
   * is considered too slow for the body. If ReadHeaderTimeout
   * is zero, the value of ReadTimeout is used. If both are
   * zero, there is no timeout.
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
   * next request when keep-alives are enabled. If IdleTimeout
   * is zero, the value of ReadTimeout is used. If both are
   * zero, there is no timeout.
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
 }
 interface Server {
  /**
   * Close immediately closes all active net.Listeners and any
   * connections in state StateNew, StateActive, or StateIdle. For a
   * graceful shutdown, use Shutdown.
   * 
   * Close does not attempt to close (and does not even know about)
   * any hijacked connections, such as WebSockets.
   * 
   * Close returns any error returned from closing the Server's
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
   * error returned from closing the Server's underlying Listener(s).
   * 
   * When Shutdown is called, Serve, ListenAndServe, and
   * ListenAndServeTLS immediately return ErrServerClosed. Make sure the
   * program doesn't exit and waits instead for Shutdown to return.
   * 
   * Shutdown does not attempt to close nor wait for hijacked
   * connections such as WebSockets. The caller of Shutdown should
   * separately notify such long-lived connections of shutdown and wait
   * for them to close, if desired. See RegisterOnShutdown for a way to
   * register shutdown notification functions.
   * 
   * Once Shutdown has been called on a server, it may not be reused;
   * future calls to methods such as Serve will return ErrServerClosed.
   */
  shutdown(ctx: context.Context): void
 }
 interface Server {
  /**
   * RegisterOnShutdown registers a function to call on Shutdown.
   * This can be used to gracefully shutdown connections that have
   * undergone ALPN protocol upgrade or that have been hijacked.
   * This function should start protocol-specific graceful shutdown,
   * but should not wait for shutdown to complete.
   */
  registerOnShutdown(f: () => void): void
 }
 interface Server {
  /**
   * ListenAndServe listens on the TCP network address srv.Addr and then
   * calls Serve to handle requests on incoming connections.
   * Accepted connections are configured to enable TCP keep-alives.
   * 
   * If srv.Addr is blank, ":http" is used.
   * 
   * ListenAndServe always returns a non-nil error. After Shutdown or Close,
   * the returned error is ErrServerClosed.
   */
  listenAndServe(): void
 }
 interface Server {
  /**
   * Serve accepts incoming connections on the Listener l, creating a
   * new service goroutine for each. The service goroutines read requests and
   * then call srv.Handler to reply to them.
   * 
   * HTTP/2 support is only enabled if the Listener returns *tls.Conn
   * connections and they were configured with "h2" in the TLS
   * Config.NextProtos.
   * 
   * Serve always returns a non-nil error and closes l.
   * After Shutdown or Close, the returned error is ErrServerClosed.
   */
  serve(l: net.Listener): void
 }
 interface Server {
  /**
   * ServeTLS accepts incoming connections on the Listener l, creating a
   * new service goroutine for each. The service goroutines perform TLS
   * setup and then read requests, calling srv.Handler to reply to them.
   * 
   * Files containing a certificate and matching private key for the
   * server must be provided if neither the Server's
   * TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
   * If the certificate is signed by a certificate authority, the
   * certFile should be the concatenation of the server's certificate,
   * any intermediates, and the CA's certificate.
   * 
   * ServeTLS always returns a non-nil error. After Shutdown or Close, the
   * returned error is ErrServerClosed.
   */
  serveTLS(l: net.Listener, certFile: string): void
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
   * ListenAndServeTLS listens on the TCP network address srv.Addr and
   * then calls ServeTLS to handle requests on incoming TLS connections.
   * Accepted connections are configured to enable TCP keep-alives.
   * 
   * Filenames containing a certificate and matching private key for the
   * server must be provided if neither the Server's TLSConfig.Certificates
   * nor TLSConfig.GetCertificate are populated. If the certificate is
   * signed by a certificate authority, the certFile should be the
   * concatenation of the server's certificate, any intermediates, and
   * the CA's certificate.
   * 
   * If srv.Addr is blank, ":https" is used.
   * 
   * ListenAndServeTLS always returns a non-nil error. After Shutdown or
   * Close, the returned error is ErrServerClosed.
   */
  listenAndServeTLS(certFile: string): void
 }
}

namespace auth {
 /**
  * AuthUser defines a standardized oauth2 user data structure.
  */
 interface AuthUser {
  id: string
  name: string
  username: string
  email: string
  avatarUrl: string
  rawUser: _TygojaDict
  accessToken: string
  refreshToken: string
 }
 /**
  * Provider defines a common interface for an OAuth2 client.
  */
 interface Provider {
  /**
   * Scopes returns the context associated with the provider (if any).
   */
  context(): context.Context
  /**
   * SetContext assigns the specified context to the current provider.
   */
  setContext(ctx: context.Context): void
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
   * RedirectUrl returns the end address to redirect the user
   * going through the OAuth flow.
   */
  redirectUrl(): string
  /**
   * SetRedirectUrl sets the provider's RedirectUrl.
   */
  setRedirectUrl(url: string): void
  /**
   * AuthUrl returns the provider's authorization service url.
   */
  authUrl(): string
  /**
   * SetAuthUrl sets the provider's AuthUrl.
   */
  setAuthUrl(url: string): void
  /**
   * TokenUrl returns the provider's token exchange service url.
   */
  tokenUrl(): string
  /**
   * SetTokenUrl sets the provider's TokenUrl.
   */
  setTokenUrl(url: string): void
  /**
   * UserApiUrl returns the provider's user info api url.
   */
  userApiUrl(): string
  /**
   * SetUserApiUrl sets the provider's UserApiUrl.
   */
  setUserApiUrl(url: string): void
  /**
   * Client returns an http client using the provided token.
   */
  client(token: oauth2.Token): (any | undefined)
  /**
   * BuildAuthUrl returns a URL to the provider's consent page
   * that asks for permissions for the required scopes explicitly.
   */
  buildAuthUrl(state: string, ...opts: oauth2.AuthCodeOption[]): string
  /**
   * FetchToken converts an authorization code to token.
   */
  fetchToken(code: string, ...opts: oauth2.AuthCodeOption[]): (oauth2.Token | undefined)
  /**
   * FetchRawUserData requests and marshalizes into `result` the
   * the OAuth user api response.
   */
  fetchRawUserData(token: oauth2.Token): string
  /**
   * FetchAuthUser is similar to FetchRawUserData, but normalizes and
   * marshalizes the user api response into a standardized AuthUser struct.
   */
  fetchAuthUser(token: oauth2.Token): (AuthUser | undefined)
 }
}

/**
 * Package echo implements high performance, minimalist Go web framework.
 * 
 * Example:
 * 
 * ```
 * 	  package main
 * 
 * 		import (
 * 			"github.com/labstack/echo/v5"
 * 			"github.com/labstack/echo/v5/middleware"
 * 			"log"
 * 			"net/http"
 * 		)
 * 
 * 	  // Handler
 * 	  func hello(c echo.Context) error {
 * 	    return c.String(http.StatusOK, "Hello, World!")
 * 	  }
 * 
 * 	  func main() {
 * 	    // Echo instance
 * 	    e := echo.New()
 * 
 * 	    // Middleware
 * 	    e.Use(middleware.Logger())
 * 	    e.Use(middleware.Recover())
 * 
 * 	    // Routes
 * 	    e.GET("/", hello)
 * 
 * 	    // Start server
 * 	    if err := e.Start(":8080"); err != http.ErrServerClosed {
 * 			  log.Fatal(err)
 * 		  }
 * 	  }
 * ```
 * 
 * Learn more at https://echo.labstack.com
 */
namespace echo {
 /**
  * Context represents the context of the current HTTP request. It holds request and
  * response objects, path, path parameters, data and registered handler.
  */
 interface Context {
  /**
   * Request returns `*http.Request`.
   */
  request(): (http.Request | undefined)
  /**
   * SetRequest sets `*http.Request`.
   */
  setRequest(r: http.Request): void
  /**
   * SetResponse sets `*Response`.
   */
  setResponse(r: Response): void
  /**
   * Response returns `*Response`.
   */
  response(): (Response | undefined)
  /**
   * IsTLS returns true if HTTP connection is TLS otherwise false.
   */
  isTLS(): boolean
  /**
   * IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
   */
  isWebSocket(): boolean
  /**
   * Scheme returns the HTTP protocol scheme, `http` or `https`.
   */
  scheme(): string
  /**
   * RealIP returns the client's network address based on `X-Forwarded-For`
   * or `X-Real-IP` request header.
   * The behavior can be configured using `Echo#IPExtractor`.
   */
  realIP(): string
  /**
   * RouteInfo returns current request route information. Method, Path, Name and params if they exist for matched route.
   * In case of 404 (route not found) and 405 (method not allowed) RouteInfo returns generic struct for these cases.
   */
  routeInfo(): RouteInfo
  /**
   * Path returns the registered path for the handler.
   */
  path(): string
  /**
   * PathParam returns path parameter by name.
   */
  pathParam(name: string): string
  /**
   * PathParamDefault returns the path parameter or default value for the provided name.
   * 
   * Notes for DefaultRouter implementation:
   * Path parameter could be empty for cases like that:
   * * route `/release-:version/bin` and request URL is `/release-/bin`
   * * route `/api/:version/image.jpg` and request URL is `/api//image.jpg`
   * but not when path parameter is last part of route path
   * * route `/download/file.:ext` will not match request `/download/file.`
   */
  pathParamDefault(name: string, defaultValue: string): string
  /**
   * PathParams returns path parameter values.
   */
  pathParams(): PathParams
  /**
   * SetPathParams sets path parameters for current request.
   */
  setPathParams(params: PathParams): void
  /**
   * QueryParam returns the query param for the provided name.
   */
  queryParam(name: string): string
  /**
   * QueryParamDefault returns the query param or default value for the provided name.
   */
  queryParamDefault(name: string): string
  /**
   * QueryParams returns the query parameters as `url.Values`.
   */
  queryParams(): url.Values
  /**
   * QueryString returns the URL query string.
   */
  queryString(): string
  /**
   * FormValue returns the form field value for the provided name.
   */
  formValue(name: string): string
  /**
   * FormValueDefault returns the form field value or default value for the provided name.
   */
  formValueDefault(name: string): string
  /**
   * FormValues returns the form field values as `url.Values`.
   */
  formValues(): url.Values
  /**
   * FormFile returns the multipart form file for the provided name.
   */
  formFile(name: string): (multipart.FileHeader | undefined)
  /**
   * MultipartForm returns the multipart form.
   */
  multipartForm(): (multipart.Form | undefined)
  /**
   * Cookie returns the named cookie provided in the request.
   */
  cookie(name: string): (http.Cookie | undefined)
  /**
   * SetCookie adds a `Set-Cookie` header in HTTP response.
   */
  setCookie(cookie: http.Cookie): void
  /**
   * Cookies returns the HTTP cookies sent with the request.
   */
  cookies(): Array<(http.Cookie | undefined)>
  /**
   * Get retrieves data from the context.
   */
  get(key: string): {
 }
  /**
   * Set saves data in the context.
   */
  set(key: string, val: {
  }): void
  /**
   * Bind binds path params, query params and the request body into provided type `i`. The default binder
   * binds body based on Content-Type header.
   */
  bind(i: {
  }): void
  /**
   * Validate validates provided `i`. It is usually called after `Context#Bind()`.
   * Validator must be registered using `Echo#Validator`.
   */
  validate(i: {
  }): void
  /**
   * Render renders a template with data and sends a text/html response with status
   * code. Renderer must be registered using `Echo.Renderer`.
   */
  render(code: number, name: string, data: {
  }): void
  /**
   * HTML sends an HTTP response with status code.
   */
  html(code: number, html: string): void
  /**
   * HTMLBlob sends an HTTP blob response with status code.
   */
  htmlBlob(code: number, b: string): void
  /**
   * String sends a string response with status code.
   */
  string(code: number, s: string): void
  /**
   * JSON sends a JSON response with status code.
   */
  json(code: number, i: {
  }): void
  /**
   * JSONPretty sends a pretty-print JSON with status code.
   */
  jsonPretty(code: number, i: {
  }, indent: string): void
  /**
   * JSONBlob sends a JSON blob response with status code.
   */
  jsonBlob(code: number, b: string): void
  /**
   * JSONP sends a JSONP response with status code. It uses `callback` to construct
   * the JSONP payload.
   */
  jsonp(code: number, callback: string, i: {
  }): void
  /**
   * JSONPBlob sends a JSONP blob response with status code. It uses `callback`
   * to construct the JSONP payload.
   */
  jsonpBlob(code: number, callback: string, b: string): void
  /**
   * XML sends an XML response with status code.
   */
  xml(code: number, i: {
  }): void
  /**
   * XMLPretty sends a pretty-print XML with status code.
   */
  xmlPretty(code: number, i: {
  }, indent: string): void
  /**
   * XMLBlob sends an XML blob response with status code.
   */
  xmlBlob(code: number, b: string): void
  /**
   * Blob sends a blob response with status code and content type.
   */
  blob(code: number, contentType: string, b: string): void
  /**
   * Stream sends a streaming response with status code and content type.
   */
  stream(code: number, contentType: string, r: io.Reader): void
  /**
   * File sends a response with the content of the file.
   */
  file(file: string): void
  /**
   * FileFS sends a response with the content of the file from given filesystem.
   */
  fileFS(file: string, filesystem: fs.FS): void
  /**
   * Attachment sends a response as attachment, prompting client to save the
   * file.
   */
  attachment(file: string, name: string): void
  /**
   * Inline sends a response as inline, opening the file in the browser.
   */
  inline(file: string, name: string): void
  /**
   * NoContent sends a response with no body and a status code.
   */
  noContent(code: number): void
  /**
   * Redirect redirects the request to a provided URL with status code.
   */
  redirect(code: number, url: string): void
  /**
   * Error invokes the registered global HTTP error handler. Generally used by middleware.
   * A side-effect of calling global error handler is that now Response has been committed (sent to the client) and
   * middlewares up in chain can not change Response status code or Response body anymore.
   * 
   * Avoid using this method in handlers as no middleware will be able to effectively handle errors after that.
   * Instead of calling this method in handler return your error and let it be handled by middlewares or global error handler.
   */
  error(err: Error): void
  /**
   * Echo returns the `Echo` instance.
   * 
   * WARNING: Remember that Echo public fields and methods are coroutine safe ONLY when you are NOT mutating them
   * anywhere in your code after Echo server has started.
   */
  echo(): (Echo | undefined)
 }
 // @ts-ignore
 import stdContext = context
 /**
  * Echo is the top-level framework instance.
  * 
  * Goroutine safety: Do not mutate Echo instance fields after server has started. Accessing these
  * fields from handlers/middlewares and changing field values at the same time leads to data-races.
  * Same rule applies to adding new routes after server has been started - Adding a route is not Goroutine safe action.
  */
 interface Echo {
  /**
   * NewContextFunc allows using custom context implementations, instead of default *echo.context
   */
  newContextFunc: (e: Echo, pathParamAllocSize: number) => ServableContext
  debug: boolean
  httpErrorHandler: HTTPErrorHandler
  binder: Binder
  jsonSerializer: JSONSerializer
  validator: Validator
  renderer: Renderer
  logger: Logger
  ipExtractor: IPExtractor
  /**
   * Filesystem is file system used by Static and File handlers to access files.
   * Defaults to os.DirFS(".")
   * 
   * When dealing with `embed.FS` use `fs := echo.MustSubFS(fs, "rootDirectory") to create sub fs which uses necessary
   * prefix for directory path. This is necessary as `//go:embed assets/images` embeds files with paths
   * including `assets/images` as their prefix.
   */
  filesystem: fs.FS
  /**
   * OnAddRoute is called when Echo adds new route to specific host router. Handler is called for every router
   * and before route is added to the host router.
   */
  onAddRoute: (host: string, route: Routable) => void
 }
 /**
  * HandlerFunc defines a function to serve HTTP requests.
  */
 interface HandlerFunc {(c: Context): void }
 /**
  * MiddlewareFunc defines a function to process middleware.
  */
 interface MiddlewareFunc {(next: HandlerFunc): HandlerFunc }
 interface Echo {
  /**
   * NewContext returns a new Context instance.
   * 
   * Note: both request and response can be left to nil as Echo.ServeHTTP will call c.Reset(req,resp) anyway
   * these arguments are useful when creating context for tests and cases like that.
   */
  newContext(r: http.Request, w: http.ResponseWriter): Context
 }
 interface Echo {
  /**
   * Router returns the default router.
   */
  router(): Router
 }
 interface Echo {
  /**
   * Routers returns the new map of host => router.
   */
  routers(): _TygojaDict
 }
 interface Echo {
  /**
   * RouterFor returns Router for given host. When host is left empty the default router is returned.
   */
  routerFor(host: string): [Router, boolean]
 }
 interface Echo {
  /**
   * ResetRouterCreator resets callback for creating new router instances.
   * Note: current (default) router is immediately replaced with router created with creator func and vhost routers are cleared.
   */
  resetRouterCreator(creator: (e: Echo) => Router): void
 }
 interface Echo {
  /**
   * Pre adds middleware to the chain which is run before router tries to find matching route.
   * Meaning middleware is executed even for 404 (not found) cases.
   */
  pre(...middleware: MiddlewareFunc[]): void
 }
 interface Echo {
  /**
   * Use adds middleware to the chain which is run after router has found matching route and before route/request handler method is executed.
   */
  use(...middleware: MiddlewareFunc[]): void
 }
 interface Echo {
  /**
   * CONNECT registers a new CONNECT route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  connect(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * DELETE registers a new DELETE route for a path with matching handler in the router
   * with optional route-level middleware. Panics on error.
   */
  delete(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * GET registers a new GET route for a path with matching handler in the router
   * with optional route-level middleware. Panics on error.
   */
  get(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * HEAD registers a new HEAD route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  head(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * OPTIONS registers a new OPTIONS route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  options(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * PATCH registers a new PATCH route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  patch(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * POST registers a new POST route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  post(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * PUT registers a new PUT route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  put(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * TRACE registers a new TRACE route for a path with matching handler in the
   * router with optional route-level middleware. Panics on error.
   */
  trace(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * RouteNotFound registers a special-case route which is executed when no other route is found (i.e. HTTP 404 cases)
   * for current request URL.
   * Path supports static and named/any parameters just like other http method is defined. Generally path is ended with
   * wildcard/match-any character (`/*`, `/download/*` etc).
   * 
   * Example: `e.RouteNotFound("/*", func(c echo.Context) error { return c.NoContent(http.StatusNotFound) })`
   */
  routeNotFound(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * Any registers a new route for all HTTP methods (supported by Echo) and path with matching handler
   * in the router with optional route-level middleware.
   * 
   * Note: this method only adds specific set of supported HTTP methods as handler and is not true
   * "catch-any-arbitrary-method" way of matching requests.
   */
  any(path: string, handler: HandlerFunc, ...middleware: MiddlewareFunc[]): Routes
 }
 interface Echo {
  /**
   * Match registers a new route for multiple HTTP methods and path with matching
   * handler in the router with optional route-level middleware. Panics on error.
   */
  match(methods: Array<string>, path: string, handler: HandlerFunc, ...middleware: MiddlewareFunc[]): Routes
 }
 interface Echo {
  /**
   * Static registers a new route with path prefix to serve static files from the provided root directory.
   */
  static(pathPrefix: string): RouteInfo
 }
 interface Echo {
  /**
   * StaticFS registers a new route with path prefix to serve static files from the provided file system.
   * 
   * When dealing with `embed.FS` use `fs := echo.MustSubFS(fs, "rootDirectory") to create sub fs which uses necessary
   * prefix for directory path. This is necessary as `//go:embed assets/images` embeds files with paths
   * including `assets/images` as their prefix.
   */
  staticFS(pathPrefix: string, filesystem: fs.FS): RouteInfo
 }
 interface Echo {
  /**
   * FileFS registers a new route with path to serve file from the provided file system.
   */
  fileFS(path: string, filesystem: fs.FS, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * File registers a new route with path to serve a static file with optional route-level middleware. Panics on error.
   */
  file(path: string, ...middleware: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * AddRoute registers a new Route with default host Router
   */
  addRoute(route: Routable): RouteInfo
 }
 interface Echo {
  /**
   * Add registers a new route for an HTTP method and path with matching handler
   * in the router with optional route-level middleware.
   */
  add(method: string, handler: HandlerFunc, ...middleware: MiddlewareFunc[]): RouteInfo
 }
 interface Echo {
  /**
   * Host creates a new router group for the provided host and optional host-level middleware.
   */
  host(name: string, ...m: MiddlewareFunc[]): (Group | undefined)
 }
 interface Echo {
  /**
   * Group creates a new router group with prefix and optional group-level middleware.
   */
  group(prefix: string, ...m: MiddlewareFunc[]): (Group | undefined)
 }
 interface Echo {
  /**
   * AcquireContext returns an empty `Context` instance from the pool.
   * You must return the context by calling `ReleaseContext()`.
   */
  acquireContext(): Context
 }
 interface Echo {
  /**
   * ReleaseContext returns the `Context` instance back to the pool.
   * You must call it after `AcquireContext()`.
   */
  releaseContext(c: Context): void
 }
 interface Echo {
  /**
   * ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
   */
  serveHTTP(w: http.ResponseWriter, r: http.Request): void
 }
 interface Echo {
  /**
   * Start stars HTTP server on given address with Echo as a handler serving requests. The server can be shutdown by
   * sending os.Interrupt signal with `ctrl+c`.
   * 
   * Note: this method is created for use in examples/demos and is deliberately simple without providing configuration
   * options.
   * 
   * In need of customization use:
   * 
   * ```
   * 	sc := echo.StartConfig{Address: ":8080"}
   * 	if err := sc.Start(e); err != http.ErrServerClosed {
   * 		log.Fatal(err)
   * 	}
   * ```
   * 
   * // or standard library `http.Server`
   * 
   * ```
   * 	s := http.Server{Addr: ":8080", Handler: e}
   * 	if err := s.ListenAndServe(); err != http.ErrServerClosed {
   * 		log.Fatal(err)
   * 	}
   * ```
   */
  start(address: string): void
 }
}

/**
 * Package blob provides an easy and portable way to interact with blobs
 * within a storage location. Subpackages contain driver implementations of
 * blob for supported services.
 * 
 * See https://gocloud.dev/howto/blob/ for a detailed how-to guide.
 * 
 * *blob.Bucket implements io/fs.FS and io/fs.SubFS, so it can be used with
 * functions in that package.
 * 
 * # Errors
 * 
 * The errors returned from this package can be inspected in several ways:
 * 
 * The Code function from gocloud.dev/gcerrors will return an error code, also
 * defined in that package, when invoked on an error.
 * 
 * The Bucket.ErrorAs method can retrieve the driver error underlying the returned
 * error.
 * 
 * # OpenCensus Integration
 * 
 * OpenCensus supports tracing and metric collection for multiple languages and
 * backend providers. See https://opencensus.io.
 * 
 * This API collects OpenCensus traces and metrics for the following methods:
 * ```
 *   - Attributes
 *   - Copy
 *   - Delete
 *   - ListPage
 *   - NewRangeReader, from creation until the call to Close. (NewReader and ReadAll
 *     are included because they call NewRangeReader.)
 *   - NewWriter, from creation until the call to Close.
 * ```
 * 
 * All trace and metric names begin with the package import path.
 * The traces add the method name.
 * For example, "gocloud.dev/blob/Attributes".
 * The metrics are "completed_calls", a count of completed method calls by driver,
 * method and status (error code); and "latency", a distribution of method latency
 * by driver and method.
 * For example, "gocloud.dev/blob/latency".
 * 
 * It also collects the following metrics:
 * ```
 *   - gocloud.dev/blob/bytes_read: the total number of bytes read, by driver.
 *   - gocloud.dev/blob/bytes_written: the total number of bytes written, by driver.
 * ```
 * 
 * To enable trace collection in your application, see "Configure Exporter" at
 * https://opencensus.io/quickstart/go/tracing.
 * To enable metric collection in your application, see "Exporting stats" at
 * https://opencensus.io/quickstart/go/metrics.
 */
namespace blob {
 /**
  * Reader reads bytes from a blob.
  * It implements io.ReadSeekCloser, and must be closed after
  * reads are finished.
  */
 interface Reader {
 }
 interface Reader {
  /**
   * Read implements io.Reader (https://golang.org/pkg/io/#Reader).
   */
  read(p: string): number
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
   * As converts i to driver-specific types.
   * See https://gocloud.dev/concepts/as/ for background information, the "As"
   * examples in this package for examples, and the driver package
   * documentation for the specific types supported for that driver.
   */
  as(i: {
   }): boolean
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
  md5: string
  /**
   * ETag for the blob; see https://en.wikipedia.org/wiki/HTTP_ETag.
   */
  eTag: string
 }
 interface Attributes {
  /**
   * As converts i to driver-specific types.
   * See https://gocloud.dev/concepts/as/ for background information, the "As"
   * examples in this package for examples, and the driver package
   * documentation for the specific types supported for that driver.
   */
  as(i: {
   }): boolean
 }
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
  md5: string
  /**
   * IsDir indicates that this result represents a "directory" in the
   * hierarchical namespace, ending in ListOptions.Delimiter. Key can be
   * passed as ListOptions.Prefix to list items in the "directory".
   * Fields other than Key and IsDir will not be set if IsDir is true.
   */
  isDir: boolean
 }
 interface ListObject {
  /**
   * As converts i to driver-specific types.
   * See https://gocloud.dev/concepts/as/ for background information, the "As"
   * examples in this package for examples, and the driver package
   * documentation for the specific types supported for that driver.
   */
  as(i: {
   }): boolean
 }
}

/**
 * Package types implements some commonly used db serializable types
 * like datetime, json, etc.
 */
namespace types {
 /**
  * JsonArray defines a slice that is safe for json and db read/write.
  */
 interface JsonArray<T> extends Array<T>{}
 interface JsonArray<T> {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string
 }
 interface JsonArray<T> {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface JsonArray<T> {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current JsonArray[T] instance.
   */
  scan(value: any): void
 }
 /**
  * JsonMap defines a map that is safe for json and db read/write.
  */
 interface JsonMap extends _TygojaDict{}
 interface JsonMap {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string
 }
 interface JsonMap {
  /**
   * Get retrieves a single value from the current JsonMap.
   * 
   * This helper was added primarily to assist the goja integration since custom map types
   * don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
   */
  get(key: string): any
 }
 interface JsonMap {
  /**
   * Set sets a single value in the current JsonMap.
   * 
   * This helper was added primarily to assist the goja integration since custom map types
   * don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
   */
  set(key: string, value: any): void
 }
 interface JsonMap {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface JsonMap {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current `JsonMap` instance.
   */
  scan(value: any): void
 }
}

namespace settings {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * Settings defines common app configuration options.
  */
 interface Settings {
  meta: MetaConfig
  logs: LogsConfig
  smtp: SmtpConfig
  s3: S3Config
  backups: BackupsConfig
  adminAuthToken: TokenConfig
  adminPasswordResetToken: TokenConfig
  adminFileToken: TokenConfig
  recordAuthToken: TokenConfig
  recordPasswordResetToken: TokenConfig
  recordEmailChangeToken: TokenConfig
  recordVerificationToken: TokenConfig
  recordFileToken: TokenConfig
  /**
   * Deprecated: Will be removed in v0.9+
   */
  emailAuth: EmailAuthConfig
  googleAuth: AuthProviderConfig
  facebookAuth: AuthProviderConfig
  githubAuth: AuthProviderConfig
  gitlabAuth: AuthProviderConfig
  discordAuth: AuthProviderConfig
  twitterAuth: AuthProviderConfig
  microsoftAuth: AuthProviderConfig
  spotifyAuth: AuthProviderConfig
  kakaoAuth: AuthProviderConfig
  twitchAuth: AuthProviderConfig
  stravaAuth: AuthProviderConfig
  giteeAuth: AuthProviderConfig
  livechatAuth: AuthProviderConfig
  giteaAuth: AuthProviderConfig
  oidcAuth: AuthProviderConfig
  oidc2Auth: AuthProviderConfig
  oidc3Auth: AuthProviderConfig
  appleAuth: AuthProviderConfig
  instagramAuth: AuthProviderConfig
  vkAuth: AuthProviderConfig
  yandexAuth: AuthProviderConfig
 }
 interface Settings {
  /**
   * Validate makes Settings validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface Settings {
  /**
   * Merge merges `other` settings into the current one.
   */
  merge(other: Settings): void
 }
 interface Settings {
  /**
   * Clone creates a new deep copy of the current settings.
   */
  clone(): (Settings | undefined)
 }
 interface Settings {
  /**
   * RedactClone creates a new deep copy of the current settings,
   * while replacing the secret values with `******`.
   */
  redactClone(): (Settings | undefined)
 }
 interface Settings {
  /**
   * NamedAuthProviderConfigs returns a map with all registered OAuth2
   * provider configurations (indexed by their name identifier).
   */
  namedAuthProviderConfigs(): _TygojaDict
 }
}

/**
 * Package schema implements custom Schema and SchemaField datatypes
 * for handling the Collection schema definitions.
 */
namespace schema {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * Schema defines a dynamic db schema as a slice of `SchemaField`s.
  */
 interface Schema {
 }
 interface Schema {
  /**
   * Fields returns the registered schema fields.
   */
  fields(): Array<(SchemaField | undefined)>
 }
 interface Schema {
  /**
   * InitFieldsOptions calls `InitOptions()` for all schema fields.
   */
  initFieldsOptions(): void
 }
 interface Schema {
  /**
   * Clone creates a deep clone of the current schema.
   */
  clone(): (Schema | undefined)
 }
 interface Schema {
  /**
   * AsMap returns a map with all registered schema field.
   * The returned map is indexed with each field name.
   */
  asMap(): _TygojaDict
 }
 interface Schema {
  /**
   * GetFieldById returns a single field by its id.
   */
  getFieldById(id: string): (SchemaField | undefined)
 }
 interface Schema {
  /**
   * GetFieldByName returns a single field by its name.
   */
  getFieldByName(name: string): (SchemaField | undefined)
 }
 interface Schema {
  /**
   * RemoveField removes a single schema field by its id.
   * 
   * This method does nothing if field with `id` doesn't exist.
   */
  removeField(id: string): void
 }
 interface Schema {
  /**
   * AddField registers the provided newField to the current schema.
   * 
   * If field with `newField.Id` already exist, the existing field is
   * replaced with the new one.
   * 
   * Otherwise the new field is appended to the other schema fields.
   */
  addField(newField: SchemaField): void
 }
 interface Schema {
  /**
   * Validate makes Schema validatable by implementing [validation.Validatable] interface.
   * 
   * Internally calls each individual field's validator and additionally
   * checks for invalid renamed fields and field name duplications.
   */
  validate(): void
 }
 interface Schema {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string
 }
 interface Schema {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   * 
   * On success, all schema field options are auto initialized.
   */
  unmarshalJSON(data: string): void
 }
 interface Schema {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface Schema {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current Schema instance.
   */
  scan(value: any): void
 }
}

/**
 * Package models implements all PocketBase DB models and DTOs.
 */
namespace models {
 type _subAzCSd = BaseModel
 interface Admin extends _subAzCSd {
  avatar: number
  email: string
  tokenKey: string
  passwordHash: string
  lastResetSentAt: types.DateTime
 }
 interface Admin {
  /**
   * TableName returns the Admin model SQL table name.
   */
  tableName(): string
 }
 interface Admin {
  /**
   * ValidatePassword validates a plain password against the model's password.
   */
  validatePassword(password: string): boolean
 }
 interface Admin {
  /**
   * SetPassword sets cryptographically secure string to `model.Password`.
   * 
   * Additionally this method also resets the LastResetSentAt and the TokenKey fields.
   */
  setPassword(password: string): void
 }
 interface Admin {
  /**
   * RefreshTokenKey generates and sets new random token key.
   */
  refreshTokenKey(): void
 }
 // @ts-ignore
 import validation = ozzo_validation
 type _subdySrd = BaseModel
 interface Collection extends _subdySrd {
  name: string
  type: string
  system: boolean
  schema: schema.Schema
  indexes: types.JsonArray<string>
  /**
   * rules
   */
  listRule?: string
  viewRule?: string
  createRule?: string
  updateRule?: string
  deleteRule?: string
  options: types.JsonMap
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
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string
 }
 interface Collection {
  /**
   * BaseOptions decodes the current collection options and returns them
   * as new [CollectionBaseOptions] instance.
   */
  baseOptions(): CollectionBaseOptions
 }
 interface Collection {
  /**
   * AuthOptions decodes the current collection options and returns them
   * as new [CollectionAuthOptions] instance.
   */
  authOptions(): CollectionAuthOptions
 }
 interface Collection {
  /**
   * ViewOptions decodes the current collection options and returns them
   * as new [CollectionViewOptions] instance.
   */
  viewOptions(): CollectionViewOptions
 }
 interface Collection {
  /**
   * NormalizeOptions updates the current collection options with a
   * new normalized state based on the collection type.
   */
  normalizeOptions(): void
 }
 interface Collection {
  /**
   * DecodeOptions decodes the current collection options into the
   * provided "result" (must be a pointer).
   */
  decodeOptions(result: any): void
 }
 interface Collection {
  /**
   * SetOptions normalizes and unmarshals the specified options into m.Options.
   */
  setOptions(typedOptions: any): void
 }
 type _subsGIYZ = BaseModel
 interface ExternalAuth extends _subsGIYZ {
  collectionId: string
  recordId: string
  provider: string
  providerId: string
 }
 interface ExternalAuth {
  tableName(): string
 }
 type _subntRqh = BaseModel
 interface Record extends _subntRqh {
 }
 interface Record {
  /**
   * TableName returns the table name associated to the current Record model.
   */
  tableName(): string
 }
 interface Record {
  /**
   * Collection returns the Collection model associated to the current Record model.
   */
  collection(): (Collection | undefined)
 }
 interface Record {
  /**
   * OriginalCopy returns a copy of the current record model populated
   * with its ORIGINAL data state (aka. the initially loaded) and
   * everything else reset to the defaults.
   */
  originalCopy(): (Record | undefined)
 }
 interface Record {
  /**
   * CleanCopy returns a copy of the current record model populated only
   * with its LATEST data state and everything else reset to the defaults.
   */
  cleanCopy(): (Record | undefined)
 }
 interface Record {
  /**
   * Expand returns a shallow copy of the current Record model expand data.
   */
  expand(): _TygojaDict
 }
 interface Record {
  /**
   * SetExpand shallow copies the provided data to the current Record model's expand.
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
   * SchemaData returns a shallow copy ONLY of the defined record schema fields data.
   */
  schemaData(): _TygojaDict
 }
 interface Record {
  /**
   * UnknownData returns a shallow copy ONLY of the unknown record fields data,
   * aka. fields that are neither one of the base and special system ones,
   * nor defined by the collection schema.
   */
  unknownData(): _TygojaDict
 }
 interface Record {
  /**
   * IgnoreEmailVisibility toggles the flag to ignore the auth record email visibility check.
   */
  ignoreEmailVisibility(state: boolean): void
 }
 interface Record {
  /**
   * WithUnknownData toggles the export/serialization of unknown data fields
   * (false by default).
   */
  withUnknownData(state: boolean): void
 }
 interface Record {
  /**
   * Set sets the provided key-value data pair for the current Record model.
   * 
   * If the record collection has field with name matching the provided "key",
   * the value will be further normalized according to the field rules.
   */
  set(key: string, value: any): void
 }
 interface Record {
  /**
   * Get returns a normalized single record model data value for "key".
   */
  get(key: string): any
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
   * GetTime returns the data value for "key" as a [time.Time] instance.
   */
  getTime(key: string): time.Time
 }
 interface Record {
  /**
   * GetDateTime returns the data value for "key" as a DateTime instance.
   */
  getDateTime(key: string): types.DateTime
 }
 interface Record {
  /**
   * GetStringSlice returns the data value for "key" as a slice of unique strings.
   */
  getStringSlice(key: string): Array<string>
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
  expandedOne(relField: string): (Record | undefined)
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
   * BaseFilesPath returns the storage dir path used by the record.
   */
  baseFilesPath(): string
 }
 interface Record {
  /**
   * FindFileFieldByFile returns the first file type field for which
   * any of the record's data contains the provided filename.
   */
  findFileFieldByFile(filename: string): (schema.SchemaField | undefined)
 }
 interface Record {
  /**
   * Load bulk loads the provided data into the current Record model.
   */
  load(data: _TygojaDict): void
 }
 interface Record {
  /**
   * ColumnValueMap implements [ColumnValueMapper] interface.
   */
  columnValueMap(): _TygojaDict
 }
 interface Record {
  /**
   * PublicExport exports only the record fields that are safe to be public.
   * 
   * For auth records, to force the export of the email field you need to set
   * `m.IgnoreEmailVisibility(true)`.
   */
  publicExport(): _TygojaDict
 }
 interface Record {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   * 
   * Only the data exported by `PublicExport()` will be serialized.
   */
  marshalJSON(): string
 }
 interface Record {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   */
  unmarshalJSON(data: string): void
 }
 interface Record {
  /**
   * ReplaceModifers returns a new map with applied modifier
   * values based on the current record and the specified data.
   * 
   * The resolved modifier keys will be removed.
   * 
   * Multiple modifiers will be applied one after another,
   * while reusing the previous base key value result (eg. 1; -5; +2 => -2).
   * 
   * Example usage:
   * 
   * ```
   * 	 newData := record.ReplaceModifers(data)
   * 		// record:  {"field": 10}
   * 		// data:    {"field+": 5}
   * 		// newData: {"field": 15}
   * ```
   */
  replaceModifers(data: _TygojaDict): _TygojaDict
 }
 interface Record {
  /**
   * Username returns the "username" auth record data value.
   */
  username(): string
 }
 interface Record {
  /**
   * SetUsername sets the "username" auth record data value.
   * 
   * This method doesn't check whether the provided value is a valid username.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setUsername(username: string): void
 }
 interface Record {
  /**
   * Email returns the "email" auth record data value.
   */
  email(): string
 }
 interface Record {
  /**
   * SetEmail sets the "email" auth record data value.
   * 
   * This method doesn't check whether the provided value is a valid email.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setEmail(email: string): void
 }
 interface Record {
  /**
   * Verified returns the "emailVisibility" auth record data value.
   */
  emailVisibility(): boolean
 }
 interface Record {
  /**
   * SetEmailVisibility sets the "emailVisibility" auth record data value.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setEmailVisibility(visible: boolean): void
 }
 interface Record {
  /**
   * Verified returns the "verified" auth record data value.
   */
  verified(): boolean
 }
 interface Record {
  /**
   * SetVerified sets the "verified" auth record data value.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setVerified(verified: boolean): void
 }
 interface Record {
  /**
   * TokenKey returns the "tokenKey" auth record data value.
   */
  tokenKey(): string
 }
 interface Record {
  /**
   * SetTokenKey sets the "tokenKey" auth record data value.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setTokenKey(key: string): void
 }
 interface Record {
  /**
   * RefreshTokenKey generates and sets new random auth record "tokenKey".
   * 
   * Returns an error if the record is not from an auth collection.
   */
  refreshTokenKey(): void
 }
 interface Record {
  /**
   * LastResetSentAt returns the "lastResentSentAt" auth record data value.
   */
  lastResetSentAt(): types.DateTime
 }
 interface Record {
  /**
   * SetLastResetSentAt sets the "lastResentSentAt" auth record data value.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setLastResetSentAt(dateTime: types.DateTime): void
 }
 interface Record {
  /**
   * LastVerificationSentAt returns the "lastVerificationSentAt" auth record data value.
   */
  lastVerificationSentAt(): types.DateTime
 }
 interface Record {
  /**
   * SetLastVerificationSentAt sets an "lastVerificationSentAt" auth record data value.
   * 
   * Returns an error if the record is not from an auth collection.
   */
  setLastVerificationSentAt(dateTime: types.DateTime): void
 }
 interface Record {
  /**
   * PasswordHash returns the "passwordHash" auth record data value.
   */
  passwordHash(): string
 }
 interface Record {
  /**
   * ValidatePassword validates a plain password against the auth record password.
   * 
   * Returns false if the password is incorrect or record is not from an auth collection.
   */
  validatePassword(password: string): boolean
 }
 interface Record {
  /**
   * SetPassword sets cryptographically secure string to the auth record "password" field.
   * This method also resets the "lastResetSentAt" and the "tokenKey" fields.
   * 
   * Returns an error if the record is not from an auth collection or
   * an empty password is provided.
   */
  setPassword(password: string): void
 }
 /**
  * RequestInfo defines a HTTP request data struct, usually used
  * as part of the `@request.*` filter resolver.
  */
 interface RequestInfo {
  method: string
  query: _TygojaDict
  data: _TygojaDict
  headers: _TygojaDict
  authRecord?: Record
  admin?: Admin
 }
 interface RequestInfo {
  /**
   * HasModifierDataKeys loosely checks if the current struct has any modifier Data keys.
   */
  hasModifierDataKeys(): boolean
 }
}

/**
 * Package daos handles common PocketBase DB model manipulations.
 * 
 * Think of daos as DB repository and service layer in one.
 */
namespace daos {
 interface Dao {
  /**
   * AdminQuery returns a new Admin select query.
   */
  adminQuery(): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindAdminById finds the admin with the provided id.
   */
  findAdminById(id: string): (models.Admin | undefined)
 }
 interface Dao {
  /**
   * FindAdminByEmail finds the admin with the provided email address.
   */
  findAdminByEmail(email: string): (models.Admin | undefined)
 }
 interface Dao {
  /**
   * FindAdminByToken finds the admin associated with the provided JWT token.
   * 
   * Returns an error if the JWT token is invalid or expired.
   */
  findAdminByToken(token: string, baseTokenKey: string): (models.Admin | undefined)
 }
 interface Dao {
  /**
   * TotalAdmins returns the number of existing admin records.
   */
  totalAdmins(): number
 }
 interface Dao {
  /**
   * IsAdminEmailUnique checks if the provided email address is not
   * already in use by other admins.
   */
  isAdminEmailUnique(email: string, ...excludeIds: string[]): boolean
 }
 interface Dao {
  /**
   * DeleteAdmin deletes the provided Admin model.
   * 
   * Returns an error if there is only 1 admin.
   */
  deleteAdmin(admin: models.Admin): void
 }
 interface Dao {
  /**
   * SaveAdmin upserts the provided Admin model.
   */
  saveAdmin(admin: models.Admin): void
 }
 /**
  * Dao handles various db operations.
  * 
  * You can think of Dao as a repository and service layer in one.
  */
 interface Dao {
  /**
   * MaxLockRetries specifies the default max "database is locked" auto retry attempts.
   */
  maxLockRetries: number
  /**
   * ModelQueryTimeout is the default max duration of a running ModelQuery().
   * 
   * This field has no effect if an explicit query context is already specified.
   */
  modelQueryTimeout: time.Duration
  /**
   * write hooks
   */
  beforeCreateFunc: (eventDao: Dao, m: models.Model, action: () => void) => void
  afterCreateFunc: (eventDao: Dao, m: models.Model) => void
  beforeUpdateFunc: (eventDao: Dao, m: models.Model, action: () => void) => void
  afterUpdateFunc: (eventDao: Dao, m: models.Model) => void
  beforeDeleteFunc: (eventDao: Dao, m: models.Model, action: () => void) => void
  afterDeleteFunc: (eventDao: Dao, m: models.Model) => void
 }
 interface Dao {
  /**
   * DB returns the default dao db builder (*dbx.DB or *dbx.TX).
   * 
   * Currently the default db builder is dao.concurrentDB but that may change in the future.
   */
  db(): dbx.Builder
 }
 interface Dao {
  /**
   * ConcurrentDB returns the dao concurrent (aka. multiple open connections)
   * db builder (*dbx.DB or *dbx.TX).
   * 
   * In a transaction the concurrentDB and nonconcurrentDB refer to the same *dbx.TX instance.
   */
  concurrentDB(): dbx.Builder
 }
 interface Dao {
  /**
   * NonconcurrentDB returns the dao nonconcurrent (aka. single open connection)
   * db builder (*dbx.DB or *dbx.TX).
   * 
   * In a transaction the concurrentDB and nonconcurrentDB refer to the same *dbx.TX instance.
   */
  nonconcurrentDB(): dbx.Builder
 }
 interface Dao {
  /**
   * Clone returns a new Dao with the same configuration options as the current one.
   */
  clone(): (Dao | undefined)
 }
 interface Dao {
  /**
   * WithoutHooks returns a new Dao with the same configuration options
   * as the current one, but without create/update/delete hooks.
   */
  withoutHooks(): (Dao | undefined)
 }
 interface Dao {
  /**
   * ModelQuery creates a new preconfigured select query with preset
   * SELECT, FROM and other common fields based on the provided model.
   */
  modelQuery(m: models.Model): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindById finds a single db record with the specified id and
   * scans the result into m.
   */
  findById(m: models.Model, id: string): void
 }
 interface Dao {
  /**
   * RunInTransaction wraps fn into a transaction.
   * 
   * It is safe to nest RunInTransaction calls as long as you use the txDao.
   */
  runInTransaction(fn: (txDao: Dao) => void): void
 }
 interface Dao {
  /**
   * Delete deletes the provided model.
   */
  delete(m: models.Model): void
 }
 interface Dao {
  /**
   * Save persists the provided model in the database.
   * 
   * If m.IsNew() is true, the method will perform a create, otherwise an update.
   * To explicitly mark a model for update you can use m.MarkAsNotNew().
   */
  save(m: models.Model): void
 }
 interface Dao {
  /**
   * CollectionQuery returns a new Collection select query.
   */
  collectionQuery(): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindCollectionsByType finds all collections by the given type.
   */
  findCollectionsByType(collectionType: string): Array<(models.Collection | undefined)>
 }
 interface Dao {
  /**
   * FindCollectionByNameOrId finds a single collection by its name (case insensitive) or id.
   */
  findCollectionByNameOrId(nameOrId: string): (models.Collection | undefined)
 }
 interface Dao {
  /**
   * IsCollectionNameUnique checks that there is no existing collection
   * with the provided name (case insensitive!).
   * 
   * Note: case insensitive check because the name is used also as a table name for the records.
   */
  isCollectionNameUnique(name: string, ...excludeIds: string[]): boolean
 }
 interface Dao {
  /**
   * FindCollectionReferences returns information for all
   * relation schema fields referencing the provided collection.
   * 
   * If the provided collection has reference to itself then it will be
   * also included in the result. To exclude it, pass the collection id
   * as the excludeId argument.
   */
  findCollectionReferences(collection: models.Collection, ...excludeIds: string[]): _TygojaDict
 }
 interface Dao {
  /**
   * DeleteCollection deletes the provided Collection model.
   * This method automatically deletes the related collection records table.
   * 
   * NB! The collection cannot be deleted, if:
   * - is system collection (aka. collection.System is true)
   * - is referenced as part of a relation field in another collection
   */
  deleteCollection(collection: models.Collection): void
 }
 interface Dao {
  /**
   * SaveCollection persists the provided Collection model and updates
   * its related records table schema.
   * 
   * If collecction.IsNew() is true, the method will perform a create, otherwise an update.
   * To explicitly mark a collection for update you can use collecction.MarkAsNotNew().
   */
  saveCollection(collection: models.Collection): void
 }
 interface Dao {
  /**
   * ImportCollections imports the provided collections list within a single transaction.
   * 
   * NB1! If deleteMissing is set, all local collections and schema fields, that are not present
   * in the imported configuration, WILL BE DELETED (including their related records data).
   * 
   * NB2! This method doesn't perform validations on the imported collections data!
   * If you need validations, use [forms.CollectionsImport].
   */
  importCollections(importedCollections: Array<(models.Collection | undefined)>, deleteMissing: boolean, afterSync: (txDao: Dao, mappedImported: _TygojaDict) => void): void
 }
 interface Dao {
  /**
   * ExternalAuthQuery returns a new ExternalAuth select query.
   */
  externalAuthQuery(): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindAllExternalAuthsByRecord returns all ExternalAuth models
   * linked to the provided auth record.
   */
  findAllExternalAuthsByRecord(authRecord: models.Record): Array<(models.ExternalAuth | undefined)>
 }
 interface Dao {
  /**
   * FindExternalAuthByProvider returns the first available
   * ExternalAuth model for the specified provider and providerId.
   */
  findExternalAuthByProvider(provider: string): (models.ExternalAuth | undefined)
 }
 interface Dao {
  /**
   * FindExternalAuthByRecordAndProvider returns the first available
   * ExternalAuth model for the specified record data and provider.
   */
  findExternalAuthByRecordAndProvider(authRecord: models.Record, provider: string): (models.ExternalAuth | undefined)
 }
 interface Dao {
  /**
   * SaveExternalAuth upserts the provided ExternalAuth model.
   */
  saveExternalAuth(model: models.ExternalAuth): void
 }
 interface Dao {
  /**
   * DeleteExternalAuth deletes the provided ExternalAuth model.
   */
  deleteExternalAuth(model: models.ExternalAuth): void
 }
 interface Dao {
  /**
   * ParamQuery returns a new Param select query.
   */
  paramQuery(): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindParamByKey finds the first Param model with the provided key.
   */
  findParamByKey(key: string): (models.Param | undefined)
 }
 interface Dao {
  /**
   * SaveParam creates or updates a Param model by the provided key-value pair.
   * The value argument will be encoded as json string.
   * 
   * If `optEncryptionKey` is provided it will encrypt the value before storing it.
   */
  saveParam(key: string, value: any, ...optEncryptionKey: string[]): void
 }
 interface Dao {
  /**
   * DeleteParam deletes the provided Param model.
   */
  deleteParam(param: models.Param): void
 }
 interface Dao {
  /**
   * RecordQuery returns a new Record select query from a collection model, id or name.
   * 
   * In case a collection id or name is provided and that collection doesn't
   * actually exists, the generated query will be created with a cancelled context
   * and will fail once an executor (Row(), One(), All(), etc.) is called.
   */
  recordQuery(collectionModelOrIdentifier: any): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindRecordById finds the Record model by its id.
   */
  findRecordById(collectionNameOrId: string, recordId: string, ...optFilters: ((q: dbx.SelectQuery) => void)[]): (models.Record | undefined)
 }
 interface Dao {
  /**
   * FindRecordsByIds finds all Record models by the provided ids.
   * If no records are found, returns an empty slice.
   */
  findRecordsByIds(collectionNameOrId: string, recordIds: Array<string>, ...optFilters: ((q: dbx.SelectQuery) => void)[]): Array<(models.Record | undefined)>
 }
 interface Dao {
  /**
   * @todo consider to depricate as it may be easier to just use dao.RecordQuery()
   * 
   * FindRecordsByExpr finds all records by the specified db expression.
   * 
   * Returns all collection records if no expressions are provided.
   * 
   * Returns an empty slice if no records are found.
   * 
   * Example:
   * 
   * ```
   * 	expr1 := dbx.HashExp{"email": "test@example.com"}
   * 	expr2 := dbx.NewExp("LOWER(username) = {:username}", dbx.Params{"username": "test"})
   * 	dao.FindRecordsByExpr("example", expr1, expr2)
   * ```
   */
  findRecordsByExpr(collectionNameOrId: string, ...exprs: dbx.Expression[]): Array<(models.Record | undefined)>
 }
 interface Dao {
  /**
   * FindFirstRecordByData returns the first found record matching
   * the provided key-value pair.
   */
  findFirstRecordByData(collectionNameOrId: string, key: string, value: any): (models.Record | undefined)
 }
 interface Dao {
  /**
   * FindRecordsByFilter returns limit number of records matching the
   * provided string filter.
   * 
   * NB! Use the last "params" argument to bind untrusted user variables!
   * 
   * The sort argument is optional and can be empty string OR the same format
   * used in the web APIs, eg. "-created,title".
   * 
   * If the limit argument is <= 0, no limit is applied to the query and
   * all matching records are returned.
   * 
   * Example:
   * 
   * ```
   * 	dao.FindRecordsByFilter(
   * 		"posts",
   * 		"title ~ {:title} && visible = {:visible}",
   * 		"-created",
   * 		10,
   * 		0,
   * 		dbx.Params{"title": "lorem ipsum", "visible": true}
   * 	)
   * ```
   */
  findRecordsByFilter(collectionNameOrId: string, filter: string, sort: string, limit: number, offset: number, ...params: dbx.Params[]): Array<(models.Record | undefined)>
 }
 interface Dao {
  /**
   * FindFirstRecordByFilter returns the first available record matching the provided filter.
   * 
   * NB! Use the last params argument to bind untrusted user variables!
   * 
   * Example:
   * 
   * ```
   * 	dao.FindFirstRecordByFilter("posts", "slug={:slug} && status='public'", dbx.Params{"slug": "test"})
   * ```
   */
  findFirstRecordByFilter(collectionNameOrId: string, filter: string, ...params: dbx.Params[]): (models.Record | undefined)
 }
 interface Dao {
  /**
   * IsRecordValueUnique checks if the provided key-value pair is a unique Record value.
   * 
   * For correctness, if the collection is "auth" and the key is "username",
   * the unique check will be case insensitive.
   * 
   * NB! Array values (eg. from multiple select fields) are matched
   * as a serialized json strings (eg. `["a","b"]`), so the value uniqueness
   * depends on the elements order. Or in other words the following values
   * are considered different: `[]string{"a","b"}` and `[]string{"b","a"}`
   */
  isRecordValueUnique(collectionNameOrId: string, key: string, value: any, ...excludeIds: string[]): boolean
 }
 interface Dao {
  /**
   * FindAuthRecordByToken finds the auth record associated with the provided JWT token.
   * 
   * Returns an error if the JWT token is invalid, expired or not associated to an auth collection record.
   */
  findAuthRecordByToken(token: string, baseTokenKey: string): (models.Record | undefined)
 }
 interface Dao {
  /**
   * FindAuthRecordByEmail finds the auth record associated with the provided email.
   * 
   * Returns an error if it is not an auth collection or the record is not found.
   */
  findAuthRecordByEmail(collectionNameOrId: string, email: string): (models.Record | undefined)
 }
 interface Dao {
  /**
   * FindAuthRecordByUsername finds the auth record associated with the provided username (case insensitive).
   * 
   * Returns an error if it is not an auth collection or the record is not found.
   */
  findAuthRecordByUsername(collectionNameOrId: string, username: string): (models.Record | undefined)
 }
 interface Dao {
  /**
   * SuggestUniqueAuthRecordUsername checks if the provided username is unique
   * and return a new "unique" username with appended random numeric part
   * (eg. "existingName" -> "existingName583").
   * 
   * The same username will be returned if the provided string is already unique.
   */
  suggestUniqueAuthRecordUsername(collectionNameOrId: string, baseUsername: string, ...excludeIds: string[]): string
 }
 interface Dao {
  /**
   * CanAccessRecord checks if a record is allowed to be accessed by the
   * specified requestInfo and accessRule.
   * 
   * Rule and db checks are ignored in case requestInfo.Admin is set.
   * 
   * The returned error indicate that something unexpected happened during
   * the check (eg. invalid rule or db error).
   * 
   * The method always return false on invalid access rule or db error.
   * 
   * Example:
   * 
   * ```
   * 	requestInfo := apis.RequestInfo(c /* echo.Context *\/)
   * 	record, _ := dao.FindRecordById("example", "RECORD_ID")
   * 	rule := types.Pointer("@request.auth.id != '' || status = 'public'")
   * 	// ... or use one of the record collection's rule, eg. record.Collection().ViewRule
   * 
   * 	if ok, _ := dao.CanAccessRecord(record, requestInfo, rule); ok { ... }
   * ```
   */
  canAccessRecord(record: models.Record, requestInfo: models.RequestInfo, accessRule: string): boolean
 }
 interface Dao {
  /**
   * SaveRecord persists the provided Record model in the database.
   * 
   * If record.IsNew() is true, the method will perform a create, otherwise an update.
   * To explicitly mark a record for update you can use record.MarkAsNotNew().
   */
  saveRecord(record: models.Record): void
 }
 interface Dao {
  /**
   * DeleteRecord deletes the provided Record model.
   * 
   * This method will also cascade the delete operation to all linked
   * relational records (delete or unset, depending on the rel settings).
   * 
   * The delete operation may fail if the record is part of a required
   * reference in another record (aka. cannot be deleted or unset).
   */
  deleteRecord(record: models.Record): void
 }
 interface Dao {
  /**
   * ExpandRecord expands the relations of a single Record model.
   * 
   * If optFetchFunc is not set, then a default function will be used
   * that returns all relation records.
   * 
   * Returns a map with the failed expand parameters and their errors.
   */
  expandRecord(record: models.Record, expands: Array<string>, optFetchFunc: ExpandFetchFunc): _TygojaDict
 }
 interface Dao {
  /**
   * ExpandRecords expands the relations of the provided Record models list.
   * 
   * If optFetchFunc is not set, then a default function will be used
   * that returns all relation records.
   * 
   * Returns a map with the failed expand parameters and their errors.
   */
  expandRecords(records: Array<(models.Record | undefined)>, expands: Array<string>, optFetchFunc: ExpandFetchFunc): _TygojaDict
 }
 // @ts-ignore
 import validation = ozzo_validation
 interface Dao {
  /**
   * SyncRecordTableSchema compares the two provided collections
   * and applies the necessary related record table changes.
   * 
   * If `oldCollection` is null, then only `newCollection` is used to create the record table.
   */
  syncRecordTableSchema(newCollection: models.Collection, oldCollection: models.Collection): void
 }
 interface Dao {
  /**
   * RequestQuery returns a new Request logs select query.
   */
  requestQuery(): (dbx.SelectQuery | undefined)
 }
 interface Dao {
  /**
   * FindRequestById finds a single Request log by its id.
   */
  findRequestById(id: string): (models.Request | undefined)
 }
 interface Dao {
  /**
   * RequestsStats returns hourly grouped requests logs statistics.
   */
  requestsStats(expr: dbx.Expression): Array<(RequestsStatsItem | undefined)>
 }
 interface Dao {
  /**
   * DeleteOldRequests delete all requests that are created before createdBefore.
   */
  deleteOldRequests(createdBefore: time.Time): void
 }
 interface Dao {
  /**
   * SaveRequest upserts the provided Request model.
   */
  saveRequest(request: models.Request): void
 }
 interface Dao {
  /**
   * FindSettings returns and decode the serialized app settings param value.
   * 
   * The method will first try to decode the param value without decryption.
   * If it fails and optEncryptionKey is set, it will try again by first
   * decrypting the value and then decode it again.
   * 
   * Returns an error if it fails to decode the stored serialized param value.
   */
  findSettings(...optEncryptionKey: string[]): (settings.Settings | undefined)
 }
 interface Dao {
  /**
   * SaveSettings persists the specified settings configuration.
   * 
   * If optEncryptionKey is set, then the stored serialized value will be encrypted with it.
   */
  saveSettings(newSettings: settings.Settings, ...optEncryptionKey: string[]): void
 }
 interface Dao {
  /**
   * HasTable checks if a table (or view) with the provided name exists (case insensitive).
   */
  hasTable(tableName: string): boolean
 }
 interface Dao {
  /**
   * TableColumns returns all column names of a single table by its name.
   */
  tableColumns(tableName: string): Array<string>
 }
 interface Dao {
  /**
   * TableInfo returns the `table_info` pragma result for the specified table.
   */
  tableInfo(tableName: string): Array<(models.TableInfoRow | undefined)>
 }
 interface Dao {
  /**
   * TableIndexes returns a name grouped map with all non empty index of the specified table.
   * 
   * Note: This method doesn't return an error on nonexisting table.
   */
  tableIndexes(tableName: string): _TygojaDict
 }
 interface Dao {
  /**
   * DeleteTable drops the specified table.
   * 
   * This method is a no-op if a table with the provided name doesn't exist.
   * 
   * Be aware that this method is vulnerable to SQL injection and the
   * "tableName" argument must come only from trusted input!
   */
  deleteTable(tableName: string): void
 }
 interface Dao {
  /**
   * Vacuum executes VACUUM on the current dao.DB() instance in order to
   * reclaim unused db disk space.
   */
  vacuum(): void
 }
 interface Dao {
  /**
   * DeleteView drops the specified view name.
   * 
   * This method is a no-op if a view with the provided name doesn't exist.
   * 
   * Be aware that this method is vulnerable to SQL injection and the
   * "name" argument must come only from trusted input!
   */
  deleteView(name: string): void
 }
 interface Dao {
  /**
   * SaveView creates (or updates already existing) persistent SQL view.
   * 
   * Be aware that this method is vulnerable to SQL injection and the
   * "selectQuery" argument must come only from trusted input!
   */
  saveView(name: string, selectQuery: string): void
 }
 interface Dao {
  /**
   * CreateViewSchema creates a new view schema from the provided select query.
   * 
   * There are some caveats:
   * - The select query must have an "id" column.
   * - Wildcard ("*") columns are not supported to avoid accidentally leaking sensitive data.
   */
  createViewSchema(selectQuery: string): schema.Schema
 }
 interface Dao {
  /**
   * FindRecordByViewFile returns the original models.Record of the
   * provided view collection file.
   */
  findRecordByViewFile(viewCollectionNameOrId: string, fileFieldName: string, filename: string): (models.Record | undefined)
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
  */
 interface App {
  /**
   * Deprecated:
   * This method may get removed in the near future.
   * It is recommended to access the app db instance from app.Dao().DB() or
   * if you want more flexibility - app.Dao().ConcurrentDB() and app.Dao().NonconcurrentDB().
   * 
   * DB returns the default app database instance.
   */
  db(): (dbx.DB | undefined)
  /**
   * Dao returns the default app Dao instance.
   * 
   * This Dao could operate only on the tables and models
   * associated with the default app database. For example,
   * trying to access the request logs table will result in error.
   */
  dao(): (daos.Dao | undefined)
  /**
   * Deprecated:
   * This method may get removed in the near future.
   * It is recommended to access the logs db instance from app.LogsDao().DB() or
   * if you want more flexibility - app.LogsDao().ConcurrentDB() and app.LogsDao().NonconcurrentDB().
   * 
   * LogsDB returns the app logs database instance.
   */
  logsDB(): (dbx.DB | undefined)
  /**
   * LogsDao returns the app logs Dao instance.
   * 
   * This Dao could operate only on the tables and models
   * associated with the logs database. For example, trying to access
   * the users table from LogsDao will result in error.
   */
  logsDao(): (daos.Dao | undefined)
  /**
   * DataDir returns the app data directory path.
   */
  dataDir(): string
  /**
   * EncryptionEnv returns the name of the app secret env key
   * (used for settings encryption).
   */
  encryptionEnv(): string
  /**
   * IsDebug returns whether the app is in debug mode
   * (showing more detailed error logs, executed sql statements, etc.).
   */
  isDebug(): boolean
  /**
   * Settings returns the loaded app settings.
   */
  settings(): (settings.Settings | undefined)
  /**
   * Cache returns the app internal cache store.
   */
  cache(): (store.Store<any> | undefined)
  /**
   * SubscriptionsBroker returns the app realtime subscriptions broker instance.
   */
  subscriptionsBroker(): (subscriptions.Broker | undefined)
  /**
   * NewMailClient creates and returns a configured app mail client.
   */
  newMailClient(): mailer.Mailer
  /**
   * NewFilesystem creates and returns a configured filesystem.System instance
   * for managing regular app files (eg. collection uploads).
   * 
   * NB! Make sure to call Close() on the returned result
   * after you are done working with it.
   */
  newFilesystem(): (filesystem.System | undefined)
  /**
   * NewBackupsFilesystem creates and returns a configured filesystem.System instance
   * for managing app backups.
   * 
   * NB! Make sure to call Close() on the returned result
   * after you are done working with it.
   */
  newBackupsFilesystem(): (filesystem.System | undefined)
  /**
   * RefreshSettings reinitializes and reloads the stored application settings.
   */
  refreshSettings(): void
  /**
   * IsBootstrapped checks if the application was initialized
   * (aka. whether Bootstrap() was called).
   */
  isBootstrapped(): boolean
  /**
   * Bootstrap takes care for initializing the application
   * (open db connections, load settings, etc.).
   * 
   * It will call ResetBootstrapState() if the application was already bootstrapped.
   */
  bootstrap(): void
  /**
   * ResetBootstrapState takes care for releasing initialized app resources
   * (eg. closing db connections).
   */
  resetBootstrapState(): void
  /**
   * CreateBackup creates a new backup of the current app pb_data directory.
   * 
   * Backups can be stored on S3 if it is configured in app.Settings().Backups.
   * 
   * Please refer to the godoc of the specific core.App implementation
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
   * Please refer to the godoc of the specific core.App implementation
   * for details on the restore procedures.
   * 
   * NB! This feature is experimental and currently is expected to work only on UNIX based systems.
   */
  restoreBackup(ctx: context.Context, name: string): void
  /**
   * Restart restarts the current running application process.
   * 
   * Currently it is relying on execve so it is supported only on UNIX based systems.
   */
  restart(): void
  /**
   * OnBeforeBootstrap hook is triggered before initializing the main
   * application resources (eg. before db open and initial settings load).
   */
  onBeforeBootstrap(): (hook.Hook<BootstrapEvent | undefined> | undefined)
  /**
   * OnAfterBootstrap hook is triggered after initializing the main
   * application resources (eg. after db open and initial settings load).
   */
  onAfterBootstrap(): (hook.Hook<BootstrapEvent | undefined> | undefined)
  /**
   * OnBeforeServe hook is triggered before serving the internal router (echo),
   * allowing you to adjust its options and attach new routes or middlewares.
   */
  onBeforeServe(): (hook.Hook<ServeEvent | undefined> | undefined)
  /**
   * OnBeforeApiError hook is triggered right before sending an error API
   * response to the client, allowing you to further modify the error data
   * or to return a completely different API response.
   */
  onBeforeApiError(): (hook.Hook<ApiErrorEvent | undefined> | undefined)
  /**
   * OnAfterApiError hook is triggered right after sending an error API
   * response to the client.
   * It could be used to log the final API error in external services.
   */
  onAfterApiError(): (hook.Hook<ApiErrorEvent | undefined> | undefined)
  /**
   * OnTerminate hook is triggered when the app is in the process
   * of being terminated (eg. on SIGTERM signal).
   */
  onTerminate(): (hook.Hook<TerminateEvent | undefined> | undefined)
  /**
   * OnModelBeforeCreate hook is triggered before inserting a new
   * model in the DB, allowing you to modify or validate the stored data.
   * 
   * If the optional "tags" list (table names and/or the Collection id for Record models)
   * is specified, then all event handlers registered via the created hook
   * will be triggered and called only if their event data origin matches the tags.
   */
  onModelBeforeCreate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined> | undefined)
  /**
   * OnModelAfterCreate hook is triggered after successfully
   * inserting a new model in the DB.
   * 
   * If the optional "tags" list (table names and/or the Collection id for Record models)
   * is specified, then all event handlers registered via the created hook
   * will be triggered and called only if their event data origin matches the tags.
   */
  onModelAfterCreate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined> | undefined)
  /**
   * OnModelBeforeUpdate hook is triggered before updating existing
   * model in the DB, allowing you to modify or validate the stored data.
   * 
   * If the optional "tags" list (table names and/or the Collection id for Record models)
   * is specified, then all event handlers registered via the created hook
   * will be triggered and called only if their event data origin matches the tags.
   */
  onModelBeforeUpdate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined> | undefined)
  /**
   * OnModelAfterUpdate hook is triggered after successfully updating
   * existing model in the DB.
   * 
   * If the optional "tags" list (table names and/or the Collection id for Record models)
   * is specified, then all event handlers registered via the created hook
   * will be triggered and called only if their event data origin matches the tags.
   */
  onModelAfterUpdate(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined> | undefined)
  /**
   * OnModelBeforeDelete hook is triggered before deleting an
   * existing model from the DB.
   * 
   * If the optional "tags" list (table names and/or the Collection id for Record models)
   * is specified, then all event handlers registered via the created hook
   * will be triggered and called only if their event data origin matches the tags.
   */
  onModelBeforeDelete(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined> | undefined)
  /**
   * OnModelAfterDelete hook is triggered after successfully deleting an
   * existing model from the DB.
   * 
   * If the optional "tags" list (table names and/or the Collection id for Record models)
   * is specified, then all event handlers registered via the created hook
   * will be triggered and called only if their event data origin matches the tags.
   */
  onModelAfterDelete(...tags: string[]): (hook.TaggedHook<ModelEvent | undefined> | undefined)
  /**
   * OnMailerBeforeAdminResetPasswordSend hook is triggered right
   * before sending a password reset email to an admin, allowing you
   * to inspect and customize the email message that is being sent.
   */
  onMailerBeforeAdminResetPasswordSend(): (hook.Hook<MailerAdminEvent | undefined> | undefined)
  /**
   * OnMailerAfterAdminResetPasswordSend hook is triggered after
   * admin password reset email was successfully sent.
   */
  onMailerAfterAdminResetPasswordSend(): (hook.Hook<MailerAdminEvent | undefined> | undefined)
  /**
   * OnMailerBeforeRecordResetPasswordSend hook is triggered right
   * before sending a password reset email to an auth record, allowing
   * you to inspect and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerBeforeRecordResetPasswordSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined> | undefined)
  /**
   * OnMailerAfterRecordResetPasswordSend hook is triggered after
   * an auth record password reset email was successfully sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerAfterRecordResetPasswordSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined> | undefined)
  /**
   * OnMailerBeforeRecordVerificationSend hook is triggered right
   * before sending a verification email to an auth record, allowing
   * you to inspect and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerBeforeRecordVerificationSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined> | undefined)
  /**
   * OnMailerAfterRecordVerificationSend hook is triggered after a
   * verification email was successfully sent to an auth record.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerAfterRecordVerificationSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined> | undefined)
  /**
   * OnMailerBeforeRecordChangeEmailSend hook is triggered right before
   * sending a confirmation new address email to an auth record, allowing
   * you to inspect and customize the email message that is being sent.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerBeforeRecordChangeEmailSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined> | undefined)
  /**
   * OnMailerAfterRecordChangeEmailSend hook is triggered after a
   * verification email was successfully sent to an auth record.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onMailerAfterRecordChangeEmailSend(...tags: string[]): (hook.TaggedHook<MailerRecordEvent | undefined> | undefined)
  /**
   * OnRealtimeConnectRequest hook is triggered right before establishing
   * the SSE client connection.
   */
  onRealtimeConnectRequest(): (hook.Hook<RealtimeConnectEvent | undefined> | undefined)
  /**
   * OnRealtimeDisconnectRequest hook is triggered on disconnected/interrupted
   * SSE client connection.
   */
  onRealtimeDisconnectRequest(): (hook.Hook<RealtimeDisconnectEvent | undefined> | undefined)
  /**
   * OnRealtimeBeforeMessage hook is triggered right before sending
   * an SSE message to a client.
   * 
   * Returning [hook.StopPropagation] will prevent sending the message.
   * Returning any other non-nil error will close the realtime connection.
   */
  onRealtimeBeforeMessageSend(): (hook.Hook<RealtimeMessageEvent | undefined> | undefined)
  /**
   * OnRealtimeBeforeMessage hook is triggered right after sending
   * an SSE message to a client.
   */
  onRealtimeAfterMessageSend(): (hook.Hook<RealtimeMessageEvent | undefined> | undefined)
  /**
   * OnRealtimeBeforeSubscribeRequest hook is triggered before changing
   * the client subscriptions, allowing you to further validate and
   * modify the submitted change.
   */
  onRealtimeBeforeSubscribeRequest(): (hook.Hook<RealtimeSubscribeEvent | undefined> | undefined)
  /**
   * OnRealtimeAfterSubscribeRequest hook is triggered after the client
   * subscriptions were successfully changed.
   */
  onRealtimeAfterSubscribeRequest(): (hook.Hook<RealtimeSubscribeEvent | undefined> | undefined)
  /**
   * OnSettingsListRequest hook is triggered on each successful
   * API Settings list request.
   * 
   * Could be used to validate or modify the response before
   * returning it to the client.
   */
  onSettingsListRequest(): (hook.Hook<SettingsListEvent | undefined> | undefined)
  /**
   * OnSettingsBeforeUpdateRequest hook is triggered before each API
   * Settings update request (after request data load and before settings persistence).
   * 
   * Could be used to additionally validate the request data or
   * implement completely different persistence behavior.
   */
  onSettingsBeforeUpdateRequest(): (hook.Hook<SettingsUpdateEvent | undefined> | undefined)
  /**
   * OnSettingsAfterUpdateRequest hook is triggered after each
   * successful API Settings update request.
   */
  onSettingsAfterUpdateRequest(): (hook.Hook<SettingsUpdateEvent | undefined> | undefined)
  /**
   * OnFileDownloadRequest hook is triggered before each API File download request.
   * 
   * Could be used to validate or modify the file response before
   * returning it to the client.
   */
  onFileDownloadRequest(...tags: string[]): (hook.TaggedHook<FileDownloadEvent | undefined> | undefined)
  /**
   * OnFileBeforeTokenRequest hook is triggered before each file
   * token API request.
   * 
   * If no token or model was submitted, e.Model and e.Token will be empty,
   * allowing you to implement your own custom model file auth implementation.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onFileBeforeTokenRequest(...tags: string[]): (hook.TaggedHook<FileTokenEvent | undefined> | undefined)
  /**
   * OnFileAfterTokenRequest hook is triggered after each
   * successful file token API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onFileAfterTokenRequest(...tags: string[]): (hook.TaggedHook<FileTokenEvent | undefined> | undefined)
  /**
   * OnAdminsListRequest hook is triggered on each API Admins list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onAdminsListRequest(): (hook.Hook<AdminsListEvent | undefined> | undefined)
  /**
   * OnAdminViewRequest hook is triggered on each API Admin view request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onAdminViewRequest(): (hook.Hook<AdminViewEvent | undefined> | undefined)
  /**
   * OnAdminBeforeCreateRequest hook is triggered before each API
   * Admin create request (after request data load and before model persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onAdminBeforeCreateRequest(): (hook.Hook<AdminCreateEvent | undefined> | undefined)
  /**
   * OnAdminAfterCreateRequest hook is triggered after each
   * successful API Admin create request.
   */
  onAdminAfterCreateRequest(): (hook.Hook<AdminCreateEvent | undefined> | undefined)
  /**
   * OnAdminBeforeUpdateRequest hook is triggered before each API
   * Admin update request (after request data load and before model persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onAdminBeforeUpdateRequest(): (hook.Hook<AdminUpdateEvent | undefined> | undefined)
  /**
   * OnAdminAfterUpdateRequest hook is triggered after each
   * successful API Admin update request.
   */
  onAdminAfterUpdateRequest(): (hook.Hook<AdminUpdateEvent | undefined> | undefined)
  /**
   * OnAdminBeforeDeleteRequest hook is triggered before each API
   * Admin delete request (after model load and before actual deletion).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different delete behavior.
   */
  onAdminBeforeDeleteRequest(): (hook.Hook<AdminDeleteEvent | undefined> | undefined)
  /**
   * OnAdminAfterDeleteRequest hook is triggered after each
   * successful API Admin delete request.
   */
  onAdminAfterDeleteRequest(): (hook.Hook<AdminDeleteEvent | undefined> | undefined)
  /**
   * OnAdminAuthRequest hook is triggered on each successful API Admin
   * authentication request (sign-in, token refresh, etc.).
   * 
   * Could be used to additionally validate or modify the
   * authenticated admin data and token.
   */
  onAdminAuthRequest(): (hook.Hook<AdminAuthEvent | undefined> | undefined)
  /**
   * OnAdminBeforeAuthWithPasswordRequest hook is triggered before each Admin
   * auth with password API request (after request data load and before password validation).
   * 
   * Could be used to implement for example a custom password validation
   * or to locate a different Admin identity (by assigning [AdminAuthWithPasswordEvent.Admin]).
   */
  onAdminBeforeAuthWithPasswordRequest(): (hook.Hook<AdminAuthWithPasswordEvent | undefined> | undefined)
  /**
   * OnAdminAfterAuthWithPasswordRequest hook is triggered after each
   * successful Admin auth with password API request.
   */
  onAdminAfterAuthWithPasswordRequest(): (hook.Hook<AdminAuthWithPasswordEvent | undefined> | undefined)
  /**
   * OnAdminBeforeAuthRefreshRequest hook is triggered before each Admin
   * auth refresh API request (right before generating a new auth token).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different auth refresh behavior.
   */
  onAdminBeforeAuthRefreshRequest(): (hook.Hook<AdminAuthRefreshEvent | undefined> | undefined)
  /**
   * OnAdminAfterAuthRefreshRequest hook is triggered after each
   * successful auth refresh API request (right after generating a new auth token).
   */
  onAdminAfterAuthRefreshRequest(): (hook.Hook<AdminAuthRefreshEvent | undefined> | undefined)
  /**
   * OnAdminBeforeRequestPasswordResetRequest hook is triggered before each Admin
   * request password reset API request (after request data load and before sending the reset email).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different password reset behavior.
   */
  onAdminBeforeRequestPasswordResetRequest(): (hook.Hook<AdminRequestPasswordResetEvent | undefined> | undefined)
  /**
   * OnAdminAfterRequestPasswordResetRequest hook is triggered after each
   * successful request password reset API request.
   */
  onAdminAfterRequestPasswordResetRequest(): (hook.Hook<AdminRequestPasswordResetEvent | undefined> | undefined)
  /**
   * OnAdminBeforeConfirmPasswordResetRequest hook is triggered before each Admin
   * confirm password reset API request (after request data load and before persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onAdminBeforeConfirmPasswordResetRequest(): (hook.Hook<AdminConfirmPasswordResetEvent | undefined> | undefined)
  /**
   * OnAdminAfterConfirmPasswordResetRequest hook is triggered after each
   * successful confirm password reset API request.
   */
  onAdminAfterConfirmPasswordResetRequest(): (hook.Hook<AdminConfirmPasswordResetEvent | undefined> | undefined)
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
  onRecordAuthRequest(...tags: string[]): (hook.TaggedHook<RecordAuthEvent | undefined> | undefined)
  /**
   * OnRecordBeforeAuthWithPasswordRequest hook is triggered before each Record
   * auth with password API request (after request data load and before password validation).
   * 
   * Could be used to implement for example a custom password validation
   * or to locate a different Record model (by reassigning [RecordAuthWithPasswordEvent.Record]).
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeAuthWithPasswordRequest(...tags: string[]): (hook.TaggedHook<RecordAuthWithPasswordEvent | undefined> | undefined)
  /**
   * OnRecordAfterAuthWithPasswordRequest hook is triggered after each
   * successful Record auth with password API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterAuthWithPasswordRequest(...tags: string[]): (hook.TaggedHook<RecordAuthWithPasswordEvent | undefined> | undefined)
  /**
   * OnRecordBeforeAuthWithOAuth2Request hook is triggered before each Record
   * OAuth2 sign-in/sign-up API request (after token exchange and before external provider linking).
   * 
   * If the [RecordAuthWithOAuth2Event.Record] is not set, then the OAuth2
   * request will try to create a new auth Record.
   * 
   * To assign or link a different existing record model you can
   * change the [RecordAuthWithOAuth2Event.Record] field.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeAuthWithOAuth2Request(...tags: string[]): (hook.TaggedHook<RecordAuthWithOAuth2Event | undefined> | undefined)
  /**
   * OnRecordAfterAuthWithOAuth2Request hook is triggered after each
   * successful Record OAuth2 API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterAuthWithOAuth2Request(...tags: string[]): (hook.TaggedHook<RecordAuthWithOAuth2Event | undefined> | undefined)
  /**
   * OnRecordBeforeAuthRefreshRequest hook is triggered before each Record
   * auth refresh API request (right before generating a new auth token).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different auth refresh behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeAuthRefreshRequest(...tags: string[]): (hook.TaggedHook<RecordAuthRefreshEvent | undefined> | undefined)
  /**
   * OnRecordAfterAuthRefreshRequest hook is triggered after each
   * successful auth refresh API request (right after generating a new auth token).
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterAuthRefreshRequest(...tags: string[]): (hook.TaggedHook<RecordAuthRefreshEvent | undefined> | undefined)
  /**
   * OnRecordListExternalAuthsRequest hook is triggered on each API record external auths list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordListExternalAuthsRequest(...tags: string[]): (hook.TaggedHook<RecordListExternalAuthsEvent | undefined> | undefined)
  /**
   * OnRecordBeforeUnlinkExternalAuthRequest hook is triggered before each API record
   * external auth unlink request (after models load and before the actual relation deletion).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different delete behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeUnlinkExternalAuthRequest(...tags: string[]): (hook.TaggedHook<RecordUnlinkExternalAuthEvent | undefined> | undefined)
  /**
   * OnRecordAfterUnlinkExternalAuthRequest hook is triggered after each
   * successful API record external auth unlink request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterUnlinkExternalAuthRequest(...tags: string[]): (hook.TaggedHook<RecordUnlinkExternalAuthEvent | undefined> | undefined)
  /**
   * OnRecordBeforeRequestPasswordResetRequest hook is triggered before each Record
   * request password reset API request (after request data load and before sending the reset email).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different password reset behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeRequestPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordRequestPasswordResetEvent | undefined> | undefined)
  /**
   * OnRecordAfterRequestPasswordResetRequest hook is triggered after each
   * successful request password reset API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterRequestPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordRequestPasswordResetEvent | undefined> | undefined)
  /**
   * OnRecordBeforeConfirmPasswordResetRequest hook is triggered before each Record
   * confirm password reset API request (after request data load and before persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeConfirmPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmPasswordResetEvent | undefined> | undefined)
  /**
   * OnRecordAfterConfirmPasswordResetRequest hook is triggered after each
   * successful confirm password reset API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterConfirmPasswordResetRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmPasswordResetEvent | undefined> | undefined)
  /**
   * OnRecordBeforeRequestVerificationRequest hook is triggered before each Record
   * request verification API request (after request data load and before sending the verification email).
   * 
   * Could be used to additionally validate the loaded request data or implement
   * completely different verification behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeRequestVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordRequestVerificationEvent | undefined> | undefined)
  /**
   * OnRecordAfterRequestVerificationRequest hook is triggered after each
   * successful request verification API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterRequestVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordRequestVerificationEvent | undefined> | undefined)
  /**
   * OnRecordBeforeConfirmVerificationRequest hook is triggered before each Record
   * confirm verification API request (after request data load and before persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeConfirmVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmVerificationEvent | undefined> | undefined)
  /**
   * OnRecordAfterConfirmVerificationRequest hook is triggered after each
   * successful confirm verification API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterConfirmVerificationRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmVerificationEvent | undefined> | undefined)
  /**
   * OnRecordBeforeRequestEmailChangeRequest hook is triggered before each Record request email change API request
   * (after request data load and before sending the email link to confirm the change).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different request email change behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeRequestEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEmailChangeEvent | undefined> | undefined)
  /**
   * OnRecordAfterRequestEmailChangeRequest hook is triggered after each
   * successful request email change API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterRequestEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordRequestEmailChangeEvent | undefined> | undefined)
  /**
   * OnRecordBeforeConfirmEmailChangeRequest hook is triggered before each Record
   * confirm email change API request (after request data load and before persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeConfirmEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmEmailChangeEvent | undefined> | undefined)
  /**
   * OnRecordAfterConfirmEmailChangeRequest hook is triggered after each
   * successful confirm email change API request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterConfirmEmailChangeRequest(...tags: string[]): (hook.TaggedHook<RecordConfirmEmailChangeEvent | undefined> | undefined)
  /**
   * OnRecordsListRequest hook is triggered on each API Records list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordsListRequest(...tags: string[]): (hook.TaggedHook<RecordsListEvent | undefined> | undefined)
  /**
   * OnRecordViewRequest hook is triggered on each API Record view request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordViewRequest(...tags: string[]): (hook.TaggedHook<RecordViewEvent | undefined> | undefined)
  /**
   * OnRecordBeforeCreateRequest hook is triggered before each API Record
   * create request (after request data load and before model persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeCreateRequest(...tags: string[]): (hook.TaggedHook<RecordCreateEvent | undefined> | undefined)
  /**
   * OnRecordAfterCreateRequest hook is triggered after each
   * successful API Record create request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterCreateRequest(...tags: string[]): (hook.TaggedHook<RecordCreateEvent | undefined> | undefined)
  /**
   * OnRecordBeforeUpdateRequest hook is triggered before each API Record
   * update request (after request data load and before model persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeUpdateRequest(...tags: string[]): (hook.TaggedHook<RecordUpdateEvent | undefined> | undefined)
  /**
   * OnRecordAfterUpdateRequest hook is triggered after each
   * successful API Record update request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterUpdateRequest(...tags: string[]): (hook.TaggedHook<RecordUpdateEvent | undefined> | undefined)
  /**
   * OnRecordBeforeDeleteRequest hook is triggered before each API Record
   * delete request (after model load and before actual deletion).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different delete behavior.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordBeforeDeleteRequest(...tags: string[]): (hook.TaggedHook<RecordDeleteEvent | undefined> | undefined)
  /**
   * OnRecordAfterDeleteRequest hook is triggered after each
   * successful API Record delete request.
   * 
   * If the optional "tags" list (Collection ids or names) is specified,
   * then all event handlers registered via the created hook will be
   * triggered and called only if their event data origin matches the tags.
   */
  onRecordAfterDeleteRequest(...tags: string[]): (hook.TaggedHook<RecordDeleteEvent | undefined> | undefined)
  /**
   * OnCollectionsListRequest hook is triggered on each API Collections list request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onCollectionsListRequest(): (hook.Hook<CollectionsListEvent | undefined> | undefined)
  /**
   * OnCollectionViewRequest hook is triggered on each API Collection view request.
   * 
   * Could be used to validate or modify the response before returning it to the client.
   */
  onCollectionViewRequest(): (hook.Hook<CollectionViewEvent | undefined> | undefined)
  /**
   * OnCollectionBeforeCreateRequest hook is triggered before each API Collection
   * create request (after request data load and before model persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onCollectionBeforeCreateRequest(): (hook.Hook<CollectionCreateEvent | undefined> | undefined)
  /**
   * OnCollectionAfterCreateRequest hook is triggered after each
   * successful API Collection create request.
   */
  onCollectionAfterCreateRequest(): (hook.Hook<CollectionCreateEvent | undefined> | undefined)
  /**
   * OnCollectionBeforeUpdateRequest hook is triggered before each API Collection
   * update request (after request data load and before model persistence).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different persistence behavior.
   */
  onCollectionBeforeUpdateRequest(): (hook.Hook<CollectionUpdateEvent | undefined> | undefined)
  /**
   * OnCollectionAfterUpdateRequest hook is triggered after each
   * successful API Collection update request.
   */
  onCollectionAfterUpdateRequest(): (hook.Hook<CollectionUpdateEvent | undefined> | undefined)
  /**
   * OnCollectionBeforeDeleteRequest hook is triggered before each API
   * Collection delete request (after model load and before actual deletion).
   * 
   * Could be used to additionally validate the request data or implement
   * completely different delete behavior.
   */
  onCollectionBeforeDeleteRequest(): (hook.Hook<CollectionDeleteEvent | undefined> | undefined)
  /**
   * OnCollectionAfterDeleteRequest hook is triggered after each
   * successful API Collection delete request.
   */
  onCollectionAfterDeleteRequest(): (hook.Hook<CollectionDeleteEvent | undefined> | undefined)
  /**
   * OnCollectionsBeforeImportRequest hook is triggered before each API
   * collections import request (after request data load and before the actual import).
   * 
   * Could be used to additionally validate the imported collections or
   * to implement completely different import behavior.
   */
  onCollectionsBeforeImportRequest(): (hook.Hook<CollectionsImportEvent | undefined> | undefined)
  /**
   * OnCollectionsAfterImportRequest hook is triggered after each
   * successful API collections import request.
   */
  onCollectionsAfterImportRequest(): (hook.Hook<CollectionsImportEvent | undefined> | undefined)
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
  * ReadCloser is the interface that groups the basic Read and Close methods.
  */
 interface ReadCloser {
 }
 /**
  * WriteCloser is the interface that groups the basic Write and Close methods.
  */
 interface WriteCloser {
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
 * On most systems, that error has type syscall.Errno.
 * 
 * Deprecated: this package is locked down. Callers should use the
 * corresponding package in the golang.org/x/sys repository instead.
 * That is also where updates required by new systems or versions
 * should be applied. See https://golang.org/s/go1.4-syscall for more
 * information.
 */
namespace syscall {
 /**
  * SysProcIDMap holds Container ID to Host ID mappings used for User Namespaces in Linux.
  * See user_namespaces(7).
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
  * by a child process started by StartProcess.
  */
 interface Credential {
  uid: number // User ID.
  gid: number // Group ID.
  groups: Array<number> // Supplementary group IDs.
  noSetGroups: boolean // If true, don't set supplementary groups
 }
 /**
  * A Signal is a number describing a process signal.
  * It implements the os.Signal interface.
  */
 interface Signal extends Number{}
 interface Signal {
  signal(): void
 }
 interface Signal {
  string(): string
 }
}

/**
 * Package time provides functionality for measuring and displaying time.
 * 
 * The calendrical calculations always assume a Gregorian calendar, with
 * no leap seconds.
 * 
 * Monotonic Clocks
 * 
 * Operating systems provide both a “wall clock,” which is subject to
 * changes for clock synchronization, and a “monotonic clock,” which is
 * not. The general rule is that the wall clock is for telling time and
 * the monotonic clock is for measuring time. Rather than split the API,
 * in this package the Time returned by time.Now contains both a wall
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
 * Other idioms, such as time.Since(start), time.Until(deadline), and
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
 * t.After(u), t.Before(u), t.Equal(u), and t.Sub(u) are carried out
 * using the monotonic clock readings alone, ignoring the wall clock
 * readings. If either t or u contains no monotonic clock reading, these
 * operations fall back to using the wall clock readings.
 * 
 * On some systems the monotonic clock will stop if the computer goes to sleep.
 * On such a system, t.Sub(u) may not accurately reflect the actual
 * time that passed between t and u.
 * 
 * Because the monotonic clock reading has no meaning outside
 * the current process, the serialized forms generated by t.GobEncode,
 * t.MarshalBinary, t.MarshalJSON, and t.MarshalText omit the monotonic
 * clock reading, and t.Format provides no format for it. Similarly, the
 * constructors time.Date, time.Parse, time.ParseInLocation, and time.Unix,
 * as well as the unmarshalers t.GobDecode, t.UnmarshalBinary.
 * t.UnmarshalJSON, and t.UnmarshalText always create times with
 * no monotonic clock reading.
 * 
 * Note that the Go == operator compares not just the time instant but
 * also the Location and the monotonic clock reading. See the
 * documentation for the Time type for a discussion of equality
 * testing for Time values.
 * 
 * For debugging, the result of t.String does include the monotonic
 * clock reading if present. If t != u because of different monotonic clock readings,
 * that difference will be visible when printing t.String() and u.String().
 */
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
  */
 interface Location {
 }
 interface Location {
  /**
   * String returns a descriptive name for the time zone information,
   * corresponding to the name argument to LoadLocation or FixedZone.
   */
  string(): string
 }
}

/**
 * Package fs defines basic interfaces to a file system.
 * A file system can be provided by the host operating system
 * but also by other packages.
 */
namespace fs {
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
  * Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/.
  * A consequence is that it is impossible to tell which slashes in the Path were
  * slashes in the raw URL and which were %2f. This distinction is rarely important,
  * but when it is, the code should use RawPath, an optional field which only gets
  * set if the default encoding is different from Path.
  * 
  * URL's String method uses the EscapedPath method to obtain the path. See the
  * EscapedPath method for more details.
  */
 interface URL {
  scheme: string
  opaque: string // encoded opaque data
  user?: Userinfo // username and password information
  host: string // host or host:port
  path: string // path (relative paths may omit leading slash)
  rawPath: string // encoded path hint (see EscapedPath method)
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
   * The String and RequestURI methods use EscapedPath to construct
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
   * The String method uses EscapedFragment to construct its result.
   * In general, code should call EscapedFragment instead of
   * reading u.RawFragment directly.
   */
  escapedFragment(): string
 }
 interface URL {
  /**
   * String reassembles the URL into a valid URL string.
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
   * 	- if u.Scheme is empty, scheme: is omitted.
   * 	- if u.User is nil, userinfo@ is omitted.
   * 	- if u.Host is empty, host/ is omitted.
   * 	- if u.Scheme and u.Host are empty and u.User is nil,
   * 	   the entire scheme://userinfo@host/ is omitted.
   * 	- if u.Host is non-empty and u.Path begins with a /,
   * 	   the form host/path does not add its own /.
   * 	- if u.RawQuery is empty, ?query is omitted.
   * 	- if u.Fragment is empty, #fragment is omitted.
   * ```
   */
  string(): string
 }
 interface URL {
  /**
   * Redacted is like String but replaces any password with "xxxxx".
   * Only the password in u.URL is redacted.
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
  set(key: string): void
 }
 interface Values {
  /**
   * Add adds the value to key. It appends to any existing
   * values associated with key.
   */
  add(key: string): void
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
   * Encode encodes the values into ``URL encoded'' form
   * ("bar=baz&foo=quux") sorted by key.
   */
  encode(): string
 }
 interface URL {
  /**
   * IsAbs reports whether the URL is absolute.
   * Absolute means that it has a non-empty scheme.
   */
  isAbs(): boolean
 }
 interface URL {
  /**
   * Parse parses a URL in the context of the receiver. The provided URL
   * may be relative or absolute. Parse returns nil, err on parse
   * failure, otherwise its return value is the same as ResolveReference.
   */
  parse(ref: string): (URL | undefined)
 }
 interface URL {
  /**
   * ResolveReference resolves a URI reference to an absolute URI from
   * an absolute base URI u, per RFC 3986 Section 5.2. The URI reference
   * may be relative or absolute. ResolveReference always returns a new
   * URL instance, even if the returned URL is identical to either the
   * base or reference. If ref is an absolute URL, then ResolveReference
   * ignores base and returns a copy of ref.
   */
  resolveReference(ref: URL): (URL | undefined)
 }
 interface URL {
  /**
   * Query parses RawQuery and returns the corresponding values.
   * It silently discards malformed value pairs.
   * To check errors use ParseQuery.
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
  marshalBinary(): string
 }
 interface URL {
  unmarshalBinary(text: string): void
 }
}

/**
 * Package context defines the Context type, which carries deadlines,
 * cancellation signals, and other request-scoped values across API boundaries
 * and between processes.
 * 
 * Incoming requests to a server should create a Context, and outgoing
 * calls to servers should accept a Context. The chain of function
 * calls between them must propagate the Context, optionally replacing
 * it with a derived Context created using WithCancel, WithDeadline,
 * WithTimeout, or WithValue. When a Context is canceled, all
 * Contexts derived from it are also canceled.
 * 
 * The WithCancel, WithDeadline, and WithTimeout functions take a
 * Context (the parent) and return a derived Context (the child) and a
 * CancelFunc. Calling the CancelFunc cancels the child and its
 * children, removes the parent's reference to the child, and stops
 * any associated timers. Failing to call the CancelFunc leaks the
 * child and its children until the parent is canceled or the timer
 * fires. The go vet tool checks that CancelFuncs are used on all
 * control-flow paths.
 * 
 * Programs that use Contexts should follow these rules to keep interfaces
 * consistent across packages and enable static analysis tools to check context
 * propagation:
 * 
 * Do not store Contexts inside a struct type; instead, pass a Context
 * explicitly to each function that needs it. The Context should be the first
 * parameter, typically named ctx:
 * 
 * ```
 * 	func DoSomething(ctx context.Context, arg Arg) error {
 * 		// ... use ctx ...
 * 	}
 * ```
 * 
 * Do not pass a nil Context, even if a function permits it. Pass context.TODO
 * if you are unsure about which Context to use.
 * 
 * Use context Values only for request-scoped data that transits processes and
 * APIs, not for passing optional parameters to functions.
 * 
 * The same Context may be passed to functions running in different goroutines;
 * Contexts are safe for simultaneous use by multiple goroutines.
 * 
 * See https://blog.golang.org/context for example code for a server that uses
 * Contexts.
 */
namespace context {
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
  * IsolationLevel is the transaction isolation level used in TxOptions.
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
  * connections. Prefer running queries from DB unless there is a specific
  * need for a continuous single database connection.
  * 
  * A Conn must call Close to return the connection to the database pool
  * and may do so concurrently with a running query.
  * 
  * After a call to Close, all operations on the
  * connection fail with ErrConnDone.
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
  queryContext(ctx: context.Context, query: string, ...args: any[]): (Rows | undefined)
 }
 interface Conn {
  /**
   * QueryRowContext executes a query that is expected to return at most one row.
   * QueryRowContext always returns a non-nil value. Errors are deferred until
   * Row's Scan method is called.
   * If the query selects no rows, the *Row's Scan will return ErrNoRows.
   * Otherwise, the *Row's Scan scans the first selected row and discards
   * the rest.
   */
  queryRowContext(ctx: context.Context, query: string, ...args: any[]): (Row | undefined)
 }
 interface Conn {
  /**
   * PrepareContext creates a prepared statement for later queries or executions.
   * Multiple queries or executions may be run concurrently from the
   * returned statement.
   * The caller must call the statement's Close method
   * when the statement is no longer needed.
   * 
   * The provided context is used for the preparation of the statement, not for the
   * execution of the statement.
   */
  prepareContext(ctx: context.Context, query: string): (Stmt | undefined)
 }
 interface Conn {
  /**
   * Raw executes f exposing the underlying driver connection for the
   * duration of f. The driverConn must not be used outside of f.
   * 
   * Once f returns and err is not driver.ErrBadConn, the Conn will continue to be usable
   * until Conn.Close is called.
   */
  raw(f: (driverConn: any) => void): void
 }
 interface Conn {
  /**
   * BeginTx starts a transaction.
   * 
   * The provided context is used until the transaction is committed or rolled back.
   * If the context is canceled, the sql package will roll back
   * the transaction. Tx.Commit will return an error if the context provided to
   * BeginTx is canceled.
   * 
   * The provided TxOptions is optional and may be nil if defaults should be used.
   * If a non-default isolation level is used that the driver doesn't support,
   * an error will be returned.
   */
  beginTx(ctx: context.Context, opts: TxOptions): (Tx | undefined)
 }
 interface Conn {
  /**
   * Close returns the connection to the connection pool.
   * All operations after a Close will return with ErrConnDone.
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
   * be math.MaxInt64 (any database limits will still apply).
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
  decimalSize(): [number, boolean]
 }
 interface ColumnType {
  /**
   * ScanType returns a Go type suitable for scanning into using Rows.Scan.
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
  nullable(): boolean
 }
 interface ColumnType {
  /**
   * DatabaseTypeName returns the database system name of the column type. If an empty
   * string is returned, then the driver type name is not supported.
   * Consult your driver documentation for a list of driver data types. Length specifiers
   * are not included.
   * Common type names include "VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL",
   * "INT", and "BIGINT".
   */
  databaseTypeName(): string
 }
 /**
  * Row is the result of calling QueryRow to select a single row.
  */
 interface Row {
 }
 interface Row {
  /**
   * Scan copies the columns from the matched row into the values
   * pointed at by dest. See the documentation on Rows.Scan for details.
   * If more than one row matches the query,
   * Scan uses the first row and discards the rest. If no row matches
   * the query, Scan returns ErrNoRows.
   */
  scan(...dest: any[]): void
 }
 interface Row {
  /**
   * Err provides a way for wrapping packages to check for
   * query errors without calling Scan.
   * Err returns the error, if any, that was encountered while running the query.
   * If this error is not nil, this error will also be returned from Scan.
   */
  err(): void
 }
}

namespace migrate {
 interface Migration {
  file: string
  up: (db: dbx.Builder) => void
  down: (db: dbx.Builder) => void
 }
}

/**
 * Package net provides a portable interface for network I/O, including
 * TCP/IP, UDP, domain name resolution, and Unix domain sockets.
 * 
 * Although the package provides access to low-level networking
 * primitives, most clients will need only the basic interface provided
 * by the Dial, Listen, and Accept functions and the associated
 * Conn and Listener interfaces. The crypto/tls package uses
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
 * Name Resolution
 * 
 * The method for resolving domain names, whether indirectly with functions like Dial
 * or directly with functions like LookupHost and LookupAddr, varies by operating system.
 * 
 * On Unix systems, the resolver has two options for resolving names.
 * It can use a pure Go resolver that sends DNS requests directly to the servers
 * listed in /etc/resolv.conf, or it can use a cgo-based resolver that calls C
 * library routines such as getaddrinfo and getnameinfo.
 * 
 * By default the pure Go resolver is used, because a blocked DNS request consumes
 * only a goroutine, while a blocked C call consumes an operating system thread.
 * When cgo is available, the cgo-based resolver is used instead under a variety of
 * conditions: on systems that do not let programs make direct DNS requests (OS X),
 * when the LOCALDOMAIN environment variable is present (even if empty),
 * when the RES_OPTIONS or HOSTALIASES environment variable is non-empty,
 * when the ASR_CONFIG environment variable is non-empty (OpenBSD only),
 * when /etc/resolv.conf or /etc/nsswitch.conf specify the use of features that the
 * Go resolver does not implement, and when the name being looked up ends in .local
 * or is an mDNS name.
 * 
 * The resolver decision can be overridden by setting the netdns value of the
 * GODEBUG environment variable (see package runtime) to go or cgo, as in:
 * 
 * ```
 * 	export GODEBUG=netdns=go    # force pure Go resolver
 * 	export GODEBUG=netdns=cgo   # force cgo resolver
 * ```
 * 
 * The decision can also be forced while building the Go source tree
 * by setting the netgo or netcgo build tag.
 * 
 * A numeric netdns setting, as in GODEBUG=netdns=1, causes the resolver
 * to print debugging information about its decisions.
 * To force a particular resolver while also printing debugging information,
 * join the two settings by a plus sign, as in GODEBUG=netdns=go+1.
 * 
 * On Plan 9, the resolver always accesses /net/cs and /net/dns.
 * 
 * On Windows, the resolver always uses C library functions, such as GetAddrInfo and DnsQuery.
 */
namespace net {
 /**
  * Conn is a generic stream-oriented network connection.
  * 
  * Multiple goroutines may invoke methods on a Conn simultaneously.
  */
 interface Conn {
  /**
   * Read reads data from the connection.
   * Read can be made to time out and return an error after a fixed
   * time limit; see SetDeadline and SetReadDeadline.
   */
  read(b: string): number
  /**
   * Write writes data to the connection.
   * Write can be made to time out and return an error after a fixed
   * time limit; see SetDeadline and SetWriteDeadline.
   */
  write(b: string): number
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
 * Package cobra is a commander providing a simple interface to create powerful modern CLI interfaces.
 * In addition to providing an interface, Cobra simultaneously provides a controller to organize your application code.
 */
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
  * ShellCompDirective is a bit map representing the different behaviors the shell
  * can be instructed to have once completions have been provided.
  */
 interface ShellCompDirective extends Number{}
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
 }
}

/**
 * Package textproto implements generic support for text-based request/response
 * protocols in the style of HTTP, NNTP, and SMTP.
 * 
 * The package provides:
 * 
 * Error, which represents a numeric error response from
 * a server.
 * 
 * Pipeline, to manage pipelined requests and responses
 * in a client.
 * 
 * Reader, to read numeric response code lines,
 * key: value headers, lines wrapped with leading spaces
 * on continuation lines, and whole text blocks ending
 * with a dot on a line by itself.
 * 
 * Writer, to write dot-encoded text blocks.
 * 
 * Conn, a convenient packaging of Reader, Writer, and Pipeline for use
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
  add(key: string): void
 }
 interface MIMEHeader {
  /**
   * Set sets the header entries associated with key to
   * the single element value. It replaces any existing
   * values associated with key.
   */
  set(key: string): void
 }
 interface MIMEHeader {
  /**
   * Get gets the first value associated with the given key.
   * It is case insensitive; CanonicalMIMEHeaderKey is used
   * to canonicalize the provided key.
   * If there are no values associated with the key, Get returns "".
   * To use non-canonical keys, access the map directly.
   */
  get(key: string): string
 }
 interface MIMEHeader {
  /**
   * Values returns all values associated with the given key.
   * It is case insensitive; CanonicalMIMEHeaderKey is
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

/**
 * Package multipart implements MIME multipart parsing, as defined in RFC
 * 2046.
 * 
 * The implementation is sufficient for HTTP (RFC 2388) and the multipart
 * bodies generated by popular browsers.
 */
namespace multipart {
 interface Reader {
  /**
   * ReadForm parses an entire multipart message whose parts have
   * a Content-Disposition of "form-data".
   * It stores up to maxMemory bytes + 10MB (reserved for non-file parts)
   * in memory. File parts which can't be stored in memory will be stored on
   * disk in temporary files.
   * It returns ErrMessageTooLarge if all non-file parts can't be stored in
   * memory.
   */
  readForm(maxMemory: number): (Form | undefined)
 }
 /**
  * Form is a parsed multipart form.
  * Its File parts are stored either in memory or on disk,
  * and are accessible via the *FileHeader's Open method.
  * Its Value parts are stored as strings.
  * Both are keyed by field name.
  */
 interface Form {
  value: _TygojaDict
  file: _TygojaDict
 }
 interface Form {
  /**
   * RemoveAll removes any temporary files associated with a Form.
   */
  removeAll(): void
 }
 /**
  * File is an interface to access the file part of a multipart message.
  * Its contents may be either stored in memory or on disk.
  * If stored on disk, the File's underlying concrete type will be an *os.File.
  */
 interface File {
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
   * When there are no more parts, the error io.EOF is returned.
   * 
   * As a special case, if the "Content-Transfer-Encoding" header
   * has a value of "quoted-printable", that header is instead
   * hidden and the body is transparently decoded during Read calls.
   */
  nextPart(): (Part | undefined)
 }
 interface Reader {
  /**
   * NextRawPart returns the next part in the multipart or an error.
   * When there are no more parts, the error io.EOF is returned.
   * 
   * Unlike NextPart, it does not have special handling for
   * "Content-Transfer-Encoding: quoted-printable".
   */
  nextRawPart(): (Part | undefined)
 }
}

/**
 * Package http provides HTTP client and server implementations.
 * 
 * Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:
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
 * The client must close the response body when finished with it:
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
 * For control over HTTP client headers, redirect policy, and other
 * settings, create a Client:
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
 * compression, and other settings, create a Transport:
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
 * ListenAndServe starts an HTTP server with a given address and handler.
 * The handler is usually nil, which means to use DefaultServeMux.
 * Handle and HandleFunc add handlers to DefaultServeMux:
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
 * Starting with Go 1.6, the http package has transparent support for the
 * HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2
 * can do so by setting Transport.TLSNextProto (for clients) or
 * Server.TLSNextProto (for servers) to a non-nil, empty
 * map. Alternatively, the following GODEBUG environment variables are
 * currently supported:
 * 
 * ```
 * 	GODEBUG=http2client=0  # disable HTTP/2 client support
 * 	GODEBUG=http2server=0  # disable HTTP/2 server support
 * 	GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
 * 	GODEBUG=http2debug=2   # ... even more verbose, with frame dumps
 * ```
 * 
 * The GODEBUG variables are not covered by Go's API compatibility
 * promise. Please report any issues before disabling HTTP/2
 * support: https://golang.org/s/http2bug
 * 
 * The http package's Transport and Server both automatically enable
 * HTTP/2 support for simple configurations. To enable HTTP/2 for more
 * complex configurations, to use lower-level HTTP/2 features, or to use
 * a newer version of Go's http2 package, import "golang.org/x/net/http2"
 * directly and use its ConfigureTransport and/or ConfigureServer
 * functions. Manually configuring HTTP/2 via the golang.org/x/net/http2
 * package takes precedence over the net/http package's built-in HTTP/2
 * support.
 */
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
  raw: string
  unparsed: Array<string> // Raw text of unparsed attribute-value pairs
 }
 interface Cookie {
  /**
   * String returns the serialization of the cookie for use in a Cookie
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
  * CanonicalHeaderKey.
  */
 interface Header extends _TygojaDict{}
 interface Header {
  /**
   * Add adds the key, value pair to the header.
   * It appends to any existing values associated with key.
   * The key is case insensitive; it is canonicalized by
   * CanonicalHeaderKey.
   */
  add(key: string): void
 }
 interface Header {
  /**
   * Set sets the header entries associated with key to the
   * single element value. It replaces any existing values
   * associated with key. The key is case insensitive; it is
   * canonicalized by textproto.CanonicalMIMEHeaderKey.
   * To use non-canonical keys, assign to the map directly.
   */
  set(key: string): void
 }
 interface Header {
  /**
   * Get gets the first value associated with the given key. If
   * there are no values associated with the key, Get returns "".
   * It is case insensitive; textproto.CanonicalMIMEHeaderKey is
   * used to canonicalize the provided key. To use non-canonical keys,
   * access the map directly.
   */
  get(key: string): string
 }
 interface Header {
  /**
   * Values returns all values associated with the given key.
   * It is case insensitive; textproto.CanonicalMIMEHeaderKey is
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
   * CanonicalHeaderKey.
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
 // @ts-ignore
 import urlpkg = url
 /**
  * Response represents the response from an HTTP request.
  * 
  * The Client and Transport return Responses from servers once
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
   * the Response's Request. ErrNoLocation is returned if no
   * Location header is present.
   */
  location(): (url.URL | undefined)
 }
 interface Response {
  /**
   * ProtoAtLeast reports whether the HTTP protocol used
   * in the response is at least major.minor.
   */
  protoAtLeast(major: number): boolean
 }
 interface Response {
  /**
   * Write writes r to w in the HTTP/1.x server response format,
   * including the status line, headers, body, and optional trailer.
   * 
   * This method consults the following fields of the response r:
   * 
   *  StatusCode
   *  ProtoMajor
   *  ProtoMinor
   *  Request.Method
   *  TransferEncoding
   *  Trailer
   *  Body
   *  ContentLength
   *  Header, values for non-canonical keys will have unpredictable behavior
   * 
   * The Response Body is closed after it is sent.
   */
  write(w: io.Writer): void
 }
 /**
  * A Handler responds to an HTTP request.
  * 
  * ServeHTTP should write reply headers and data to the ResponseWriter
  * and then return. Returning signals that the request is finished; it
  * is not valid to use the ResponseWriter or read from the
  * Request.Body after or concurrently with the completion of the
  * ServeHTTP call.
  * 
  * Depending on the HTTP client software, HTTP protocol version, and
  * any intermediaries between the client and the Go server, it may not
  * be possible to read from the Request.Body after writing to the
  * ResponseWriter. Cautious handlers should read the Request.Body
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
  * an error, panic with the value ErrAbortHandler.
  */
 interface Handler {
  serveHTTP(_arg0: ResponseWriter, _arg1: Request): void
 }
 /**
  * A ConnState represents the state of a client connection to a server.
  * It's used by the optional Server.ConnState hook.
  */
 interface ConnState extends Number{}
 interface ConnState {
  string(): string
 }
}

namespace store {
 /**
  * Store defines a concurrent safe in memory key-value data store.
  */
 interface Store<T> {
 }
 interface Store<T> {
  /**
   * Reset clears the store and replaces the store data with a
   * shallow copy of the provided newData.
   */
  reset(newData: _TygojaDict): void
 }
 interface Store<T> {
  /**
   * Length returns the current number of elements in the store.
   */
  length(): number
 }
 interface Store<T> {
  /**
   * RemoveAll removes all the existing store entries.
   */
  removeAll(): void
 }
 interface Store<T> {
  /**
   * Remove removes a single entry from the store.
   * 
   * Remove does nothing if key doesn't exist in the store.
   */
  remove(key: string): void
 }
 interface Store<T> {
  /**
   * Has checks if element with the specified key exist or not.
   */
  has(key: string): boolean
 }
 interface Store<T> {
  /**
   * Get returns a single element value from the store.
   * 
   * If key is not set, the zero T value is returned.
   */
  get(key: string): T
 }
 interface Store<T> {
  /**
   * GetAll returns a shallow copy of the current store data.
   */
  getAll(): _TygojaDict
 }
 interface Store<T> {
  /**
   * Set sets (or overwrite if already exist) a new value for key.
   */
  set(key: string, value: T): void
 }
 interface Store<T> {
  /**
   * SetIfLessThanLimit sets (or overwrite if already exist) a new value for key.
   * 
   * This method is similar to Set() but **it will skip adding new elements**
   * to the store if the store length has reached the specified limit.
   * false is returned if maxAllowedElements limit is reached.
   */
  setIfLessThanLimit(key: string, value: T, maxAllowedElements: number): boolean
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
  marshalJSON(): string
 }
 interface DateTime {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   */
  unmarshalJSON(b: string): void
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
}

/**
 * Package schema implements custom Schema and SchemaField datatypes
 * for handling the Collection schema definitions.
 */
namespace schema {
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * SchemaField defines a single schema field structure.
  */
 interface SchemaField {
  system: boolean
  id: string
  name: string
  type: string
  required: boolean
  /**
   * Presentable indicates whether the field is suitable for
   * visualization purposes (eg. in the Admin UI relation views).
   */
  presentable: boolean
  /**
   * Deprecated: This field is no-op and will be removed in future versions.
   * Please use the collection.Indexes field to define a unique constraint.
   */
  unique: boolean
  options: any
 }
 interface SchemaField {
  /**
   * ColDefinition returns the field db column type definition as string.
   */
  colDefinition(): string
 }
 interface SchemaField {
  /**
   * String serializes and returns the current field as string.
   */
  string(): string
 }
 interface SchemaField {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string
 }
 interface SchemaField {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   * 
   * The schema field options are auto initialized on success.
   */
  unmarshalJSON(data: string): void
 }
 interface SchemaField {
  /**
   * Validate makes `SchemaField` validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface SchemaField {
  /**
   * InitOptions initializes the current field options based on its type.
   * 
   * Returns error on unknown field type.
   */
  initOptions(): void
 }
 interface SchemaField {
  /**
   * PrepareValue returns normalized and properly formatted field value.
   */
  prepareValue(value: any): any
 }
 interface SchemaField {
  /**
   * PrepareValueWithModifier returns normalized and properly formatted field value
   * by "merging" baseValue with the modifierValue based on the specified modifier (+ or -).
   */
  prepareValueWithModifier(baseValue: any, modifier: string, modifierValue: any): any
 }
}

/**
 * Package models implements all PocketBase DB models and DTOs.
 */
namespace models {
 /**
  * Model defines an interface with common methods that all db models should have.
  */
 interface Model {
  tableName(): string
  isNew(): boolean
  markAsNew(): void
  markAsNotNew(): void
  hasId(): boolean
  getId(): string
  setId(id: string): void
  getCreated(): types.DateTime
  getUpdated(): types.DateTime
  refreshId(): void
  refreshCreated(): void
  refreshUpdated(): void
 }
 /**
  * BaseModel defines common fields and methods used by all other models.
  */
 interface BaseModel {
  id: string
  created: types.DateTime
  updated: types.DateTime
 }
 interface BaseModel {
  /**
   * HasId returns whether the model has a nonzero id.
   */
  hasId(): boolean
 }
 interface BaseModel {
  /**
   * GetId returns the model id.
   */
  getId(): string
 }
 interface BaseModel {
  /**
   * SetId sets the model id to the provided string value.
   */
  setId(id: string): void
 }
 interface BaseModel {
  /**
   * MarkAsNew marks the model as "new" (aka. enforces m.IsNew() to be true).
   */
  markAsNew(): void
 }
 interface BaseModel {
  /**
   * MarkAsNotNew marks the model as "not new" (aka. enforces m.IsNew() to be false)
   */
  markAsNotNew(): void
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
   * GetCreated returns the model Created datetime.
   */
  getCreated(): types.DateTime
 }
 interface BaseModel {
  /**
   * GetUpdated returns the model Updated datetime.
   */
  getUpdated(): types.DateTime
 }
 interface BaseModel {
  /**
   * RefreshId generates and sets a new model id.
   * 
   * The generated id is a cryptographically random 15 characters length string.
   */
  refreshId(): void
 }
 interface BaseModel {
  /**
   * RefreshCreated updates the model Created field with the current datetime.
   */
  refreshCreated(): void
 }
 interface BaseModel {
  /**
   * RefreshUpdated updates the model Updated field with the current datetime.
   */
  refreshUpdated(): void
 }
 interface BaseModel {
  /**
   * PostScan implements the [dbx.PostScanner] interface.
   * 
   * It is executed right after the model was populated with the db row values.
   */
  postScan(): void
 }
 // @ts-ignore
 import validation = ozzo_validation
 /**
  * CollectionBaseOptions defines the "base" Collection.Options fields.
  */
 interface CollectionBaseOptions {
 }
 interface CollectionBaseOptions {
  /**
   * Validate implements [validation.Validatable] interface.
   */
  validate(): void
 }
 /**
  * CollectionAuthOptions defines the "auth" Collection.Options fields.
  */
 interface CollectionAuthOptions {
  manageRule?: string
  allowOAuth2Auth: boolean
  allowUsernameAuth: boolean
  allowEmailAuth: boolean
  requireEmail: boolean
  exceptEmailDomains: Array<string>
  onlyEmailDomains: Array<string>
  minPasswordLength: number
 }
 interface CollectionAuthOptions {
  /**
   * Validate implements [validation.Validatable] interface.
   */
  validate(): void
 }
 /**
  * CollectionViewOptions defines the "view" Collection.Options fields.
  */
 interface CollectionViewOptions {
  query: string
 }
 interface CollectionViewOptions {
  /**
   * Validate implements [validation.Validatable] interface.
   */
  validate(): void
 }
 type _subztSzL = BaseModel
 interface Param extends _subztSzL {
  key: string
  value: types.JsonRaw
 }
 interface Param {
  tableName(): string
 }
 type _subOyDRN = BaseModel
 interface Request extends _subOyDRN {
  url: string
  method: string
  status: number
  auth: string
  userIp: string
  remoteIp: string
  referer: string
  userAgent: string
  meta: types.JsonMap
 }
 interface Request {
  tableName(): string
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
  defaultValue: types.JsonRaw
 }
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
   * If zero, TokenSource implementations will reuse the same
   * token forever and RefreshToken or equivalent
   * mechanisms for that TokenSource will not be used.
   */
  expiry: time.Time
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
   * This method is unnecessary when using Transport or an HTTP Client
   * returned by this package.
   */
  setAuthHeader(r: http.Request): void
 }
 interface Token {
  /**
   * WithExtra returns a new Token that's a clone of t, but using the
   * provided raw extra map. This is only intended for use by packages
   * implementing derivative OAuth2 flows.
   */
  withExtra(extra: {
   }): (Token | undefined)
 }
 interface Token {
  /**
   * Extra returns an extra field.
   * Extra fields are key-value pairs returned by the server as a
   * part of the token retrieval response.
   */
  extra(key: string): {
 }
 }
 interface Token {
  /**
   * Valid reports whether t is non-nil, has an AccessToken, and is not expired.
   */
  valid(): boolean
 }
}

namespace mailer {
 /**
  * Mailer defines a base mail client interface.
  */
 interface Mailer {
  /**
   * Send sends an email with the provided Message.
   */
  send(message: Message): void
 }
}

/**
 * Package echo implements high performance, minimalist Go web framework.
 * 
 * Example:
 * 
 * ```
 * 	  package main
 * 
 * 		import (
 * 			"github.com/labstack/echo/v5"
 * 			"github.com/labstack/echo/v5/middleware"
 * 			"log"
 * 			"net/http"
 * 		)
 * 
 * 	  // Handler
 * 	  func hello(c echo.Context) error {
 * 	    return c.String(http.StatusOK, "Hello, World!")
 * 	  }
 * 
 * 	  func main() {
 * 	    // Echo instance
 * 	    e := echo.New()
 * 
 * 	    // Middleware
 * 	    e.Use(middleware.Logger())
 * 	    e.Use(middleware.Recover())
 * 
 * 	    // Routes
 * 	    e.GET("/", hello)
 * 
 * 	    // Start server
 * 	    if err := e.Start(":8080"); err != http.ErrServerClosed {
 * 			  log.Fatal(err)
 * 		  }
 * 	  }
 * ```
 * 
 * Learn more at https://echo.labstack.com
 */
namespace echo {
 /**
  * Binder is the interface that wraps the Bind method.
  */
 interface Binder {
  bind(c: Context, i: {
  }): void
 }
 /**
  * ServableContext is interface that Echo context implementation must implement to be usable in middleware/handlers and
  * be able to be routed by Router.
  */
 interface ServableContext {
  /**
   * Reset resets the context after request completes. It must be called along
   * with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
   * See `Echo#ServeHTTP()`
   */
  reset(r: http.Request, w: http.ResponseWriter): void
 }
 // @ts-ignore
 import stdContext = context
 /**
  * JSONSerializer is the interface that encodes and decodes JSON to and from interfaces.
  */
 interface JSONSerializer {
  serialize(c: Context, i: {
  }, indent: string): void
  deserialize(c: Context, i: {
  }): void
 }
 /**
  * HTTPErrorHandler is a centralized HTTP error handler.
  */
 interface HTTPErrorHandler {(c: Context, err: Error): void }
 /**
  * Validator is the interface that wraps the Validate function.
  */
 interface Validator {
  validate(i: {
  }): void
 }
 /**
  * Renderer is the interface that wraps the Render function.
  */
 interface Renderer {
  render(_arg0: io.Writer, _arg1: string, _arg2: {
  }, _arg3: Context): void
 }
 /**
  * Group is a set of sub-routes for a specified route. It can be used for inner
  * routes that share a common middleware or functionality that should be separate
  * from the parent echo instance while still inheriting from it.
  */
 interface Group {
 }
 interface Group {
  /**
   * Use implements `Echo#Use()` for sub-routes within the Group.
   * Group middlewares are not executed on request when there is no matching route found.
   */
  use(...middleware: MiddlewareFunc[]): void
 }
 interface Group {
  /**
   * CONNECT implements `Echo#CONNECT()` for sub-routes within the Group. Panics on error.
   */
  connect(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * DELETE implements `Echo#DELETE()` for sub-routes within the Group. Panics on error.
   */
  delete(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * GET implements `Echo#GET()` for sub-routes within the Group. Panics on error.
   */
  get(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * HEAD implements `Echo#HEAD()` for sub-routes within the Group. Panics on error.
   */
  head(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * OPTIONS implements `Echo#OPTIONS()` for sub-routes within the Group. Panics on error.
   */
  options(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * PATCH implements `Echo#PATCH()` for sub-routes within the Group. Panics on error.
   */
  patch(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * POST implements `Echo#POST()` for sub-routes within the Group. Panics on error.
   */
  post(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * PUT implements `Echo#PUT()` for sub-routes within the Group. Panics on error.
   */
  put(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * TRACE implements `Echo#TRACE()` for sub-routes within the Group. Panics on error.
   */
  trace(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * Any implements `Echo#Any()` for sub-routes within the Group. Panics on error.
   */
  any(path: string, handler: HandlerFunc, ...middleware: MiddlewareFunc[]): Routes
 }
 interface Group {
  /**
   * Match implements `Echo#Match()` for sub-routes within the Group. Panics on error.
   */
  match(methods: Array<string>, path: string, handler: HandlerFunc, ...middleware: MiddlewareFunc[]): Routes
 }
 interface Group {
  /**
   * Group creates a new sub-group with prefix and optional sub-group-level middleware.
   * Important! Group middlewares are only executed in case there was exact route match and not
   * for 404 (not found) or 405 (method not allowed) cases. If this kind of behaviour is needed then add
   * a catch-all route `/*` for the group which handler returns always 404
   */
  group(prefix: string, ...middleware: MiddlewareFunc[]): (Group | undefined)
 }
 interface Group {
  /**
   * Static implements `Echo#Static()` for sub-routes within the Group.
   */
  static(pathPrefix: string): RouteInfo
 }
 interface Group {
  /**
   * StaticFS implements `Echo#StaticFS()` for sub-routes within the Group.
   * 
   * When dealing with `embed.FS` use `fs := echo.MustSubFS(fs, "rootDirectory") to create sub fs which uses necessary
   * prefix for directory path. This is necessary as `//go:embed assets/images` embeds files with paths
   * including `assets/images` as their prefix.
   */
  staticFS(pathPrefix: string, filesystem: fs.FS): RouteInfo
 }
 interface Group {
  /**
   * FileFS implements `Echo#FileFS()` for sub-routes within the Group.
   */
  fileFS(path: string, filesystem: fs.FS, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * File implements `Echo#File()` for sub-routes within the Group. Panics on error.
   */
  file(path: string, ...middleware: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * RouteNotFound implements `Echo#RouteNotFound()` for sub-routes within the Group.
   * 
   * Example: `g.RouteNotFound("/*", func(c echo.Context) error { return c.NoContent(http.StatusNotFound) })`
   */
  routeNotFound(path: string, h: HandlerFunc, ...m: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * Add implements `Echo#Add()` for sub-routes within the Group. Panics on error.
   */
  add(method: string, handler: HandlerFunc, ...middleware: MiddlewareFunc[]): RouteInfo
 }
 interface Group {
  /**
   * AddRoute registers a new Routable with Router
   */
  addRoute(route: Routable): RouteInfo
 }
 /**
  * IPExtractor is a function to extract IP addr from http.Request.
  * Set appropriate one to Echo#IPExtractor.
  * See https://echo.labstack.com/guide/ip-address for more details.
  */
 interface IPExtractor {(_arg0: http.Request): string }
 /**
  * Logger defines the logging interface that Echo uses internally in few places.
  * For logging in handlers use your own logger instance (dependency injected or package/public variable) from logging framework of your choice.
  */
 interface Logger {
  /**
   * Write provides writer interface for http.Server `ErrorLog` and for logging startup messages.
   * `http.Server.ErrorLog` logs errors from accepting connections, unexpected behavior from handlers,
   * and underlying FileSystem errors.
   * `logger` middleware will use this method to write its JSON payload.
   */
  write(p: string): number
  /**
   * Error logs the error
   */
  error(err: Error): void
 }
 /**
  * Response wraps an http.ResponseWriter and implements its interface to be used
  * by an HTTP handler to construct an HTTP response.
  * See: https://golang.org/pkg/net/http/#ResponseWriter
  */
 interface Response {
  writer: http.ResponseWriter
  status: number
  size: number
  committed: boolean
 }
 interface Response {
  /**
   * Header returns the header map for the writer that will be sent by
   * WriteHeader. Changing the header after a call to WriteHeader (or Write) has
   * no effect unless the modified headers were declared as trailers by setting
   * the "Trailer" header before the call to WriteHeader (see example)
   * To suppress implicit response headers, set their value to nil.
   * Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
   */
  header(): http.Header
 }
 interface Response {
  /**
   * Before registers a function which is called just before the response is written.
   */
  before(fn: () => void): void
 }
 interface Response {
  /**
   * After registers a function which is called just after the response is written.
   * If the `Content-Length` is unknown, none of the after function is executed.
   */
  after(fn: () => void): void
 }
 interface Response {
  /**
   * WriteHeader sends an HTTP response header with status code. If WriteHeader is
   * not called explicitly, the first call to Write will trigger an implicit
   * WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
   * used to send error codes.
   */
  writeHeader(code: number): void
 }
 interface Response {
  /**
   * Write writes the data to the connection as part of an HTTP reply.
   */
  write(b: string): number
 }
 interface Response {
  /**
   * Flush implements the http.Flusher interface to allow an HTTP handler to flush
   * buffered data to the client.
   * See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
   */
  flush(): void
 }
 interface Response {
  /**
   * Hijack implements the http.Hijacker interface to allow an HTTP handler to
   * take over the connection.
   * See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
   */
  hijack(): [net.Conn, (bufio.ReadWriter | undefined)]
 }
 interface Response {
  /**
   * Unwrap returns the original http.ResponseWriter.
   * ResponseController can be used to access the original http.ResponseWriter.
   * See [https://go.dev/blog/go1.20]
   */
  unwrap(): http.ResponseWriter
 }
 interface Routes {
  /**
   * Reverse reverses route to URL string by replacing path parameters with given params values.
   */
  reverse(name: string, ...params: {
   }[]): string
 }
 interface Routes {
  /**
   * FindByMethodPath searched for matching route info by method and path
   */
  findByMethodPath(method: string, path: string): RouteInfo
 }
 interface Routes {
  /**
   * FilterByMethod searched for matching route info by method
   */
  filterByMethod(method: string): Routes
 }
 interface Routes {
  /**
   * FilterByPath searched for matching route info by path
   */
  filterByPath(path: string): Routes
 }
 interface Routes {
  /**
   * FilterByName searched for matching route info by name
   */
  filterByName(name: string): Routes
 }
 /**
  * Router is interface for routing request contexts to registered routes.
  * 
  * Contract between Echo/Context instance and the router:
  * ```
  *   - all routes must be added through methods on echo.Echo instance.
  *     Reason: Echo instance uses RouteInfo.Params() length to allocate slice for paths parameters (see `Echo.contextPathParamAllocSize`).
  *   - Router must populate Context during Router.Route call with:
  *   - RoutableContext.SetPath
  *   - RoutableContext.SetRawPathParams (IMPORTANT! with same slice pointer that c.RawPathParams() returns)
  *   - RoutableContext.SetRouteInfo
  *     And optionally can set additional information to Context with RoutableContext.Set
  * ```
  */
 interface Router {
  /**
   * Add registers Routable with the Router and returns registered RouteInfo
   */
  add(routable: Routable): RouteInfo
  /**
   * Remove removes route from the Router
   */
  remove(method: string, path: string): void
  /**
   * Routes returns information about all registered routes
   */
  routes(): Routes
  /**
   * Route searches Router for matching route and applies it to the given context. In case when no matching method
   * was not found (405) or no matching route exists for path (404), router will return its implementation of 405/404
   * handler function.
   */
  route(c: RoutableContext): HandlerFunc
 }
 /**
  * Routable is interface for registering Route with Router. During route registration process the Router will
  * convert Routable to RouteInfo with ToRouteInfo method. By creating custom implementation of Routable additional
  * information about registered route can be stored in Routes (i.e. privileges used with route etc.)
  */
 interface Routable {
  /**
   * ToRouteInfo converts Routable to RouteInfo
   * 
   * This method is meant to be used by Router after it parses url for path parameters, to store information about
   * route just added.
   */
  toRouteInfo(params: Array<string>): RouteInfo
  /**
   * ToRoute converts Routable to Route which Router uses to register the method handler for path.
   * 
   * This method is meant to be used by Router to get fields (including handler and middleware functions) needed to
   * add Route to Router.
   */
  toRoute(): Route
  /**
   * ForGroup recreates routable with added group prefix and group middlewares it is grouped to.
   * 
   * Is necessary for Echo.Group to be able to add/register Routable with Router and having group prefix and group
   * middlewares included in actually registered Route.
   */
  forGroup(pathPrefix: string, middlewares: Array<MiddlewareFunc>): Routable
 }
 /**
  * Routes is collection of RouteInfo instances with various helper methods.
  */
 interface Routes extends Array<RouteInfo>{}
 /**
  * RouteInfo describes registered route base fields.
  * Method+Path pair uniquely identifies the Route. Name can have duplicates.
  */
 interface RouteInfo {
  method(): string
  path(): string
  name(): string
  params(): Array<string>
  /**
   * Reverse reverses route to URL string by replacing path parameters with given params values.
   */
  reverse(...params: {
  }[]): string
 }
 /**
  * PathParams is collections of PathParam instances with various helper methods
  */
 interface PathParams extends Array<PathParam>{}
 interface PathParams {
  /**
   * Get returns path parameter value for given name or default value.
   */
  get(name: string, defaultValue: string): string
 }
}

namespace settings {
 // @ts-ignore
 import validation = ozzo_validation
 interface TokenConfig {
  secret: string
  duration: number
 }
 interface TokenConfig {
  /**
   * Validate makes TokenConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface SmtpConfig {
  enabled: boolean
  host: string
  port: number
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
 interface SmtpConfig {
  /**
   * Validate makes SmtpConfig validatable by implementing [validation.Validatable] interface.
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
  appUrl: string
  hideControls: boolean
  senderName: string
  senderAddress: string
  verificationTemplate: EmailTemplate
  resetPasswordTemplate: EmailTemplate
  confirmEmailChangeTemplate: EmailTemplate
 }
 interface MetaConfig {
  /**
   * Validate makes MetaConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface LogsConfig {
  maxDays: number
 }
 interface LogsConfig {
  /**
   * Validate makes LogsConfig validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface AuthProviderConfig {
  enabled: boolean
  clientId: string
  clientSecret: string
  authUrl: string
  tokenUrl: string
  userApiUrl: string
 }
 interface AuthProviderConfig {
  /**
   * Validate makes `ProviderConfig` validatable by implementing [validation.Validatable] interface.
   */
  validate(): void
 }
 interface AuthProviderConfig {
  /**
   * SetupProvider loads the current AuthProviderConfig into the specified provider.
   */
  setupProvider(provider: auth.Provider): void
 }
 /**
  * Deprecated: Will be removed in v0.9+
  */
 interface EmailAuthConfig {
  enabled: boolean
  exceptDomains: Array<string>
  onlyDomains: Array<string>
  minPasswordLength: number
 }
 interface EmailAuthConfig {
  /**
   * Deprecated: Will be removed in v0.9+
   */
  validate(): void
 }
}

/**
 * Package daos handles common PocketBase DB model manipulations.
 * 
 * Think of daos as DB repository and service layer in one.
 */
namespace daos {
 /**
  * ExpandFetchFunc defines the function that is used to fetch the expanded relation records.
  */
 interface ExpandFetchFunc {(relCollection: models.Collection, relIds: Array<string>): Array<(models.Record | undefined)> }
 // @ts-ignore
 import validation = ozzo_validation
 interface RequestsStatsItem {
  total: number
  date: types.DateTime
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
   * Unregister removes a single client by its id.
   * 
   * If client with clientId doesn't exist, this method does nothing.
   */
  unregister(clientId: string): void
 }
}

namespace hook {
 /**
  * Hook defines a concurrent safe structure for handling event hooks
  * (aka. callbacks propagation).
  */
 interface Hook<T> {
 }
 interface Hook<T> {
  /**
   * PreAdd registers a new handler to the hook by prepending it to the existing queue.
   * 
   * Returns an autogenerated hook id that could be used later to remove the hook with Hook.Remove(id).
   */
  preAdd(fn: Handler<T>): string
 }
 interface Hook<T> {
  /**
   * Add registers a new handler to the hook by appending it to the existing queue.
   * 
   * Returns an autogenerated hook id that could be used later to remove the hook with Hook.Remove(id).
   */
  add(fn: Handler<T>): string
 }
 interface Hook<T> {
  /**
   * Remove removes a single hook handler by its id.
   */
  remove(id: string): void
 }
 interface Hook<T> {
  /**
   * RemoveAll removes all registered handlers.
   */
  removeAll(): void
 }
 interface Hook<T> {
  /**
   * Trigger executes all registered hook handlers one by one
   * with the specified `data` as an argument.
   * 
   * Optionally, this method allows also to register additional one off
   * handlers that will be temporary appended to the handlers queue.
   * 
   * The execution stops when:
   * - hook.StopPropagation is returned in one of the handlers
   * - any non-nil error is returned in one of the handlers
   */
  trigger(data: T, ...oneOffHandlers: Handler<T>[]): void
 }
 /**
  * TaggedHook defines a proxy hook which register handlers that are triggered only
  * if the TaggedHook.tags are empty or includes at least one of the event data tag(s).
  */
 type _subUCeYn<T> = mainHook<T>
 interface TaggedHook<T> extends _subUCeYn<T> {
 }
 interface TaggedHook<T> {
  /**
   * CanTriggerOn checks if the current TaggedHook can be triggered with
   * the provided event data tags.
   */
  canTriggerOn(tags: Array<string>): boolean
 }
 interface TaggedHook<T> {
  /**
   * PreAdd registers a new handler to the hook by prepending it to the existing queue.
   * 
   * The fn handler will be called only if the event data tags satisfy h.CanTriggerOn.
   */
  preAdd(fn: Handler<T>): string
 }
 interface TaggedHook<T> {
  /**
   * Add registers a new handler to the hook by appending it to the existing queue.
   * 
   * The fn handler will be called only if the event data tags satisfy h.CanTriggerOn.
   */
  add(fn: Handler<T>): string
 }
}

/**
 * Package core is the backbone of PocketBase.
 * 
 * It defines the main PocketBase App interface and its base implementation.
 */
namespace core {
 interface BootstrapEvent {
  app: App
 }
 interface TerminateEvent {
  app: App
 }
 interface ServeEvent {
  app: App
  router?: echo.Echo
  server?: http.Server
  certManager?: any
 }
 interface ApiErrorEvent {
  httpContext: echo.Context
  error: Error
 }
 type _subKXAIm = BaseModelEvent
 interface ModelEvent extends _subKXAIm {
  dao?: daos.Dao
 }
 type _subALxYL = BaseCollectionEvent
 interface MailerRecordEvent extends _subALxYL {
  mailClient: mailer.Mailer
  message?: mailer.Message
  record?: models.Record
  meta: _TygojaDict
 }
 interface MailerAdminEvent {
  mailClient: mailer.Mailer
  message?: mailer.Message
  admin?: models.Admin
  meta: _TygojaDict
 }
 interface RealtimeConnectEvent {
  httpContext: echo.Context
  client: subscriptions.Client
  idleTimeout: time.Duration
 }
 interface RealtimeDisconnectEvent {
  httpContext: echo.Context
  client: subscriptions.Client
 }
 interface RealtimeMessageEvent {
  httpContext: echo.Context
  client: subscriptions.Client
  message?: subscriptions.Message
 }
 interface RealtimeSubscribeEvent {
  httpContext: echo.Context
  client: subscriptions.Client
  subscriptions: Array<string>
 }
 interface SettingsListEvent {
  httpContext: echo.Context
  redactedSettings?: settings.Settings
 }
 interface SettingsUpdateEvent {
  httpContext: echo.Context
  oldSettings?: settings.Settings
  newSettings?: settings.Settings
 }
 type _subqgWfZ = BaseCollectionEvent
 interface RecordsListEvent extends _subqgWfZ {
  httpContext: echo.Context
  records: Array<(models.Record | undefined)>
  result?: search.Result
 }
 type _subaGZvN = BaseCollectionEvent
 interface RecordViewEvent extends _subaGZvN {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subnDhIp = BaseCollectionEvent
 interface RecordCreateEvent extends _subnDhIp {
  httpContext: echo.Context
  record?: models.Record
  uploadedFiles: _TygojaDict
 }
 type _subMTcMC = BaseCollectionEvent
 interface RecordUpdateEvent extends _subMTcMC {
  httpContext: echo.Context
  record?: models.Record
  uploadedFiles: _TygojaDict
 }
 type _subCcvUq = BaseCollectionEvent
 interface RecordDeleteEvent extends _subCcvUq {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subXtzhW = BaseCollectionEvent
 interface RecordAuthEvent extends _subXtzhW {
  httpContext: echo.Context
  record?: models.Record
  token: string
  meta: any
 }
 type _subPbBjK = BaseCollectionEvent
 interface RecordAuthWithPasswordEvent extends _subPbBjK {
  httpContext: echo.Context
  record?: models.Record
  identity: string
  password: string
 }
 type _subNhxHA = BaseCollectionEvent
 interface RecordAuthWithOAuth2Event extends _subNhxHA {
  httpContext: echo.Context
  providerName: string
  providerClient: auth.Provider
  record?: models.Record
  oAuth2User?: auth.AuthUser
  isNewRecord: boolean
 }
 type _subbCsYs = BaseCollectionEvent
 interface RecordAuthRefreshEvent extends _subbCsYs {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subszMYh = BaseCollectionEvent
 interface RecordRequestPasswordResetEvent extends _subszMYh {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subnRVNn = BaseCollectionEvent
 interface RecordConfirmPasswordResetEvent extends _subnRVNn {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subJpyCD = BaseCollectionEvent
 interface RecordRequestVerificationEvent extends _subJpyCD {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subHNoIf = BaseCollectionEvent
 interface RecordConfirmVerificationEvent extends _subHNoIf {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subCPOmq = BaseCollectionEvent
 interface RecordRequestEmailChangeEvent extends _subCPOmq {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subNMhOI = BaseCollectionEvent
 interface RecordConfirmEmailChangeEvent extends _subNMhOI {
  httpContext: echo.Context
  record?: models.Record
 }
 type _subJkPWJ = BaseCollectionEvent
 interface RecordListExternalAuthsEvent extends _subJkPWJ {
  httpContext: echo.Context
  record?: models.Record
  externalAuths: Array<(models.ExternalAuth | undefined)>
 }
 type _submsDbp = BaseCollectionEvent
 interface RecordUnlinkExternalAuthEvent extends _submsDbp {
  httpContext: echo.Context
  record?: models.Record
  externalAuth?: models.ExternalAuth
 }
 interface AdminsListEvent {
  httpContext: echo.Context
  admins: Array<(models.Admin | undefined)>
  result?: search.Result
 }
 interface AdminViewEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface AdminCreateEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface AdminUpdateEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface AdminDeleteEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface AdminAuthEvent {
  httpContext: echo.Context
  admin?: models.Admin
  token: string
 }
 interface AdminAuthWithPasswordEvent {
  httpContext: echo.Context
  admin?: models.Admin
  identity: string
  password: string
 }
 interface AdminAuthRefreshEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface AdminRequestPasswordResetEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface AdminConfirmPasswordResetEvent {
  httpContext: echo.Context
  admin?: models.Admin
 }
 interface CollectionsListEvent {
  httpContext: echo.Context
  collections: Array<(models.Collection | undefined)>
  result?: search.Result
 }
 type _subqyaHS = BaseCollectionEvent
 interface CollectionViewEvent extends _subqyaHS {
  httpContext: echo.Context
 }
 type _subhPtnc = BaseCollectionEvent
 interface CollectionCreateEvent extends _subhPtnc {
  httpContext: echo.Context
 }
 type _subYyOGv = BaseCollectionEvent
 interface CollectionUpdateEvent extends _subYyOGv {
  httpContext: echo.Context
 }
 type _subHHXdk = BaseCollectionEvent
 interface CollectionDeleteEvent extends _subHHXdk {
  httpContext: echo.Context
 }
 interface CollectionsImportEvent {
  httpContext: echo.Context
  collections: Array<(models.Collection | undefined)>
 }
 type _subHVSrj = BaseModelEvent
 interface FileTokenEvent extends _subHVSrj {
  httpContext: echo.Context
  token: string
 }
 type _subGAPsI = BaseCollectionEvent
 interface FileDownloadEvent extends _subGAPsI {
  httpContext: echo.Context
  record?: models.Record
  fileField?: schema.SchemaField
  servedPath: string
  servedName: string
 }
}

/**
 * Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer
 * object, creating another object (Reader or Writer) that also implements
 * the interface but provides buffering and some help for textual I/O.
 */
namespace bufio {
 /**
  * ReadWriter stores pointers to a Reader and a Writer.
  * It implements io.ReadWriter.
  */
 type _subLYkix = Reader&Writer
 interface ReadWriter extends _subLYkix {
 }
}

/**
 * Package net provides a portable interface for network I/O, including
 * TCP/IP, UDP, domain name resolution, and Unix domain sockets.
 * 
 * Although the package provides access to low-level networking
 * primitives, most clients will need only the basic interface provided
 * by the Dial, Listen, and Accept functions and the associated
 * Conn and Listener interfaces. The crypto/tls package uses
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
 * Name Resolution
 * 
 * The method for resolving domain names, whether indirectly with functions like Dial
 * or directly with functions like LookupHost and LookupAddr, varies by operating system.
 * 
 * On Unix systems, the resolver has two options for resolving names.
 * It can use a pure Go resolver that sends DNS requests directly to the servers
 * listed in /etc/resolv.conf, or it can use a cgo-based resolver that calls C
 * library routines such as getaddrinfo and getnameinfo.
 * 
 * By default the pure Go resolver is used, because a blocked DNS request consumes
 * only a goroutine, while a blocked C call consumes an operating system thread.
 * When cgo is available, the cgo-based resolver is used instead under a variety of
 * conditions: on systems that do not let programs make direct DNS requests (OS X),
 * when the LOCALDOMAIN environment variable is present (even if empty),
 * when the RES_OPTIONS or HOSTALIASES environment variable is non-empty,
 * when the ASR_CONFIG environment variable is non-empty (OpenBSD only),
 * when /etc/resolv.conf or /etc/nsswitch.conf specify the use of features that the
 * Go resolver does not implement, and when the name being looked up ends in .local
 * or is an mDNS name.
 * 
 * The resolver decision can be overridden by setting the netdns value of the
 * GODEBUG environment variable (see package runtime) to go or cgo, as in:
 * 
 * ```
 * 	export GODEBUG=netdns=go    # force pure Go resolver
 * 	export GODEBUG=netdns=cgo   # force cgo resolver
 * ```
 * 
 * The decision can also be forced while building the Go source tree
 * by setting the netgo or netcgo build tag.
 * 
 * A numeric netdns setting, as in GODEBUG=netdns=1, causes the resolver
 * to print debugging information about its decisions.
 * To force a particular resolver while also printing debugging information,
 * join the two settings by a plus sign, as in GODEBUG=netdns=go+1.
 * 
 * On Plan 9, the resolver always accesses /net/cs and /net/dns.
 * 
 * On Windows, the resolver always uses C library functions, such as GetAddrInfo and DnsQuery.
 */
namespace net {
 /**
  * Addr represents a network end point address.
  * 
  * The two methods Network and String conventionally return strings
  * that can be passed as the arguments to Dial, but the exact form
  * and meaning of the strings is up to the implementation.
  */
 interface Addr {
  network(): string // name of the network (for example, "tcp", "udp")
  string(): string // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
 }
}

/**
 * Package url parses URLs and implements query escaping.
 */
namespace url {
 /**
  * The Userinfo type is an immutable encapsulation of username and
  * password details for a URL. An existing Userinfo value is guaranteed
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

/**
 * Package multipart implements MIME multipart parsing, as defined in RFC
 * 2046.
 * 
 * The implementation is sufficient for HTTP (RFC 2388) and the multipart
 * bodies generated by popular browsers.
 */
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
   * FileName returns the filename parameter of the Part's Content-Disposition
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
  read(d: string): number
 }
 interface Part {
  close(): void
 }
}

/**
 * Package http provides HTTP client and server implementations.
 * 
 * Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:
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
 * The client must close the response body when finished with it:
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
 * For control over HTTP client headers, redirect policy, and other
 * settings, create a Client:
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
 * compression, and other settings, create a Transport:
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
 * ListenAndServe starts an HTTP server with a given address and handler.
 * The handler is usually nil, which means to use DefaultServeMux.
 * Handle and HandleFunc add handlers to DefaultServeMux:
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
 * Starting with Go 1.6, the http package has transparent support for the
 * HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2
 * can do so by setting Transport.TLSNextProto (for clients) or
 * Server.TLSNextProto (for servers) to a non-nil, empty
 * map. Alternatively, the following GODEBUG environment variables are
 * currently supported:
 * 
 * ```
 * 	GODEBUG=http2client=0  # disable HTTP/2 client support
 * 	GODEBUG=http2server=0  # disable HTTP/2 server support
 * 	GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
 * 	GODEBUG=http2debug=2   # ... even more verbose, with frame dumps
 * ```
 * 
 * The GODEBUG variables are not covered by Go's API compatibility
 * promise. Please report any issues before disabling HTTP/2
 * support: https://golang.org/s/http2bug
 * 
 * The http package's Transport and Server both automatically enable
 * HTTP/2 support for simple configurations. To enable HTTP/2 for more
 * complex configurations, to use lower-level HTTP/2 features, or to use
 * a newer version of Go's http2 package, import "golang.org/x/net/http2"
 * directly and use its ConfigureTransport and/or ConfigureServer
 * functions. Manually configuring HTTP/2 via the golang.org/x/net/http2
 * package takes precedence over the net/http package's built-in HTTP/2
 * support.
 */
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

/**
 * Package echo implements high performance, minimalist Go web framework.
 * 
 * Example:
 * 
 * ```
 * 	  package main
 * 
 * 		import (
 * 			"github.com/labstack/echo/v5"
 * 			"github.com/labstack/echo/v5/middleware"
 * 			"log"
 * 			"net/http"
 * 		)
 * 
 * 	  // Handler
 * 	  func hello(c echo.Context) error {
 * 	    return c.String(http.StatusOK, "Hello, World!")
 * 	  }
 * 
 * 	  func main() {
 * 	    // Echo instance
 * 	    e := echo.New()
 * 
 * 	    // Middleware
 * 	    e.Use(middleware.Logger())
 * 	    e.Use(middleware.Recover())
 * 
 * 	    // Routes
 * 	    e.GET("/", hello)
 * 
 * 	    // Start server
 * 	    if err := e.Start(":8080"); err != http.ErrServerClosed {
 * 			  log.Fatal(err)
 * 		  }
 * 	  }
 * ```
 * 
 * Learn more at https://echo.labstack.com
 */
namespace echo {
 // @ts-ignore
 import stdContext = context
 /**
  * Route contains information to adding/registering new route with the router.
  * Method+Path pair uniquely identifies the Route. It is mandatory to provide Method+Path+Handler fields.
  */
 interface Route {
  method: string
  path: string
  handler: HandlerFunc
  middlewares: Array<MiddlewareFunc>
  name: string
 }
 interface Route {
  /**
   * ToRouteInfo converts Route to RouteInfo
   */
  toRouteInfo(params: Array<string>): RouteInfo
 }
 interface Route {
  /**
   * ToRoute returns Route which Router uses to register the method handler for path.
   */
  toRoute(): Route
 }
 interface Route {
  /**
   * ForGroup recreates Route with added group prefix and group middlewares it is grouped to.
   */
  forGroup(pathPrefix: string, middlewares: Array<MiddlewareFunc>): Routable
 }
 /**
  * RoutableContext is additional interface that structures implementing Context must implement. Methods inside this
  * interface are meant for request routing purposes and should not be used in middlewares.
  */
 interface RoutableContext {
  /**
   * Request returns `*http.Request`.
   */
  request(): (http.Request | undefined)
  /**
   * RawPathParams returns raw path pathParams value. Allocation of PathParams is handled by Context.
   */
  rawPathParams(): (PathParams | undefined)
  /**
   * SetRawPathParams replaces any existing param values with new values for this context lifetime (request).
   * Do not set any other value than what you got from RawPathParams as allocation of PathParams is handled by Context.
   */
  setRawPathParams(params: PathParams): void
  /**
   * SetPath sets the registered path for the handler.
   */
  setPath(p: string): void
  /**
   * SetRouteInfo sets the route info of this request to the context.
   */
  setRouteInfo(ri: RouteInfo): void
  /**
   * Set saves data in the context. Allows router to store arbitrary (that only router has access to) data in context
   * for later use in middlewares/handler.
   */
  set(key: string, val: {
  }): void
 }
 /**
  * PathParam is tuple pf path parameter name and its value in request path
  */
 interface PathParam {
  name: string
  value: string
 }
}

namespace store {
}

/**
 * Package types implements some commonly used db serializable types
 * like datetime, json, etc.
 */
namespace types {
 /**
  * JsonRaw defines a json value type that is safe for db read/write.
  */
 interface JsonRaw extends String{}
 interface JsonRaw {
  /**
   * String returns the current JsonRaw instance as a json encoded string.
   */
  string(): string
 }
 interface JsonRaw {
  /**
   * MarshalJSON implements the [json.Marshaler] interface.
   */
  marshalJSON(): string
 }
 interface JsonRaw {
  /**
   * UnmarshalJSON implements the [json.Unmarshaler] interface.
   */
  unmarshalJSON(b: string): void
 }
 interface JsonRaw {
  /**
   * Value implements the [driver.Valuer] interface.
   */
  value(): any
 }
 interface JsonRaw {
  /**
   * Scan implements [sql.Scanner] interface to scan the provided value
   * into the current JsonRaw instance.
   */
  scan(value: {
   }): void
 }
}

namespace mailer {
 /**
  * Message defines a generic email message struct.
  */
 interface Message {
  from: mail.Address
  to: Array<mail.Address>
  bcc: Array<mail.Address>
  cc: Array<mail.Address>
  subject: string
  html: string
  text: string
  headers: _TygojaDict
  attachments: _TygojaDict
 }
}

namespace search {
 /**
  * Result defines the returned search result structure.
  */
 interface Result {
  page: number
  perPage: number
  totalItems: number
  totalPages: number
  items: any
 }
}

namespace settings {
 // @ts-ignore
 import validation = ozzo_validation
 interface EmailTemplate {
  body: string
  subject: string
  actionUrl: string
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
  resolve(appName: string, appUrl: string): string
 }
}

namespace hook {
 /**
  * Handler defines a hook handler function.
  */
 interface Handler<T> {(e: T): void }
 /**
  * wrapped local Hook embedded struct to limit the public API surface.
  */
 type _subvTioz<T> = Hook<T>
 interface mainHook<T> extends _subvTioz<T> {
 }
}

namespace subscriptions {
 /**
  * Message defines a client's channel data.
  */
 interface Message {
  name: string
  data: string
 }
 /**
  * Client is an interface for a generic subscription client.
  */
 interface Client {
  /**
   * Id Returns the unique id of the client.
   */
  id(): string
  /**
   * Channel returns the client's communication channel.
   */
  channel(): undefined
  /**
   * Subscriptions returns all subscriptions to which the client has subscribed to.
   */
  subscriptions(): _TygojaDict
  /**
   * Subscribe subscribes the client to the provided subscriptions list.
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
   * Discard marks the client as "discarded", meaning that it
   * shouldn't be used anymore for sending new messages.
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
}

/**
 * Package core is the backbone of PocketBase.
 * 
 * It defines the main PocketBase App interface and its base implementation.
 */
namespace core {
 interface BaseModelEvent {
  model: models.Model
 }
 interface BaseModelEvent {
  tags(): Array<string>
 }
 interface BaseCollectionEvent {
  collection?: models.Collection
 }
 interface BaseCollectionEvent {
  tags(): Array<string>
 }
}

/**
 * Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer
 * object, creating another object (Reader or Writer) that also implements
 * the interface but provides buffering and some help for textual I/O.
 */
namespace bufio {
 /**
  * Reader implements buffering for an io.Reader object.
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
   * Calling Reset on the zero value of Reader initializes the internal buffer
   * to the default size.
   */
  reset(r: io.Reader): void
 }
 interface Reader {
  /**
   * Peek returns the next n bytes without advancing the reader. The bytes stop
   * being valid at the next read call. If Peek returns fewer than n bytes, it
   * also returns an error explaining why the read is short. The error is
   * ErrBufferFull if n is larger than b's buffer size.
   * 
   * Calling Peek prevents a UnreadByte or UnreadRune call from succeeding
   * until the next read operation.
   */
  peek(n: number): string
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
   * The bytes are taken from at most one Read on the underlying Reader,
   * hence n may be less than len(p).
   * To read exactly len(p) bytes, use io.ReadFull(b, p).
   * At EOF, the count will be zero and err will be io.EOF.
   */
  read(p: string): number
 }
 interface Reader {
  /**
   * ReadByte reads and returns a single byte.
   * If no byte is available, returns an error.
   */
  readByte(): string
 }
 interface Reader {
  /**
   * UnreadByte unreads the last byte. Only the most recently read byte can be unread.
   * 
   * UnreadByte returns an error if the most recent method called on the
   * Reader was not a read operation. Notably, Peek, Discard, and WriteTo are not
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
  readRune(): [string, number]
 }
 interface Reader {
  /**
   * UnreadRune unreads the last rune. If the most recent method called on
   * the Reader was not a ReadRune, UnreadRune returns an error. (In this
   * regard it is stricter than UnreadByte, which will unread the last byte
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
   * ReadSlice fails with error ErrBufferFull if the buffer fills without a delim.
   * Because the data returned from ReadSlice will be overwritten
   * by the next I/O operation, most clients should use
   * ReadBytes or ReadString instead.
   * ReadSlice returns err != nil if and only if line does not end in delim.
   */
  readSlice(delim: string): string
 }
 interface Reader {
  /**
   * ReadLine is a low-level line-reading primitive. Most callers should use
   * ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
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
   * Calling UnreadByte after ReadLine will always unread the last byte read
   * (possibly a character belonging to the line end) even if that byte is not
   * part of the line returned by ReadLine.
   */
  readLine(): [string, boolean]
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
  readBytes(delim: string): string
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
  readString(delim: string): string
 }
 interface Reader {
  /**
   * WriteTo implements io.WriterTo.
   * This may make multiple calls to the Read method of the underlying Reader.
   * If the underlying reader supports the WriteTo method,
   * this calls the underlying WriteTo without buffering.
   */
  writeTo(w: io.Writer): number
 }
 /**
  * Writer implements buffering for an io.Writer object.
  * If an error occurs writing to a Writer, no more data will be
  * accepted and all subsequent writes, and Flush, will return the error.
  * After all data has been written, the client should call the
  * Flush method to guarantee all data has been forwarded to
  * the underlying io.Writer.
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
   * Calling Reset on the zero value of Writer initializes the internal buffer
   * to the default size.
   */
  reset(w: io.Writer): void
 }
 interface Writer {
  /**
   * Flush writes any buffered data to the underlying io.Writer.
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
   * passed to an immediately succeeding Write call.
   * The buffer is only valid until the next write operation on b.
   */
  availableBuffer(): string
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
  write(p: string): number
 }
 interface Writer {
  /**
   * WriteByte writes a single byte.
   */
  writeByte(c: string): void
 }
 interface Writer {
  /**
   * WriteRune writes a single Unicode code point, returning
   * the number of bytes written and any error.
   */
  writeRune(r: string): number
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
   * ReadFrom implements io.ReaderFrom. If the underlying writer
   * supports the ReadFrom method, this calls the underlying ReadFrom.
   * If there is buffered data and an underlying ReadFrom, this fills
   * the buffer and writes it before calling ReadFrom.
   */
  readFrom(r: io.Reader): number
 }
}

namespace subscriptions {
}

/**
 * Package mail implements parsing of mail messages.
 * 
 * For the most part, this package follows the syntax as specified by RFC 5322 and
 * extended by RFC 6532.
 * Notable divergences:
 * ```
 * 	* Obsolete address formats are not parsed, including addresses with
 * 	  embedded route information.
 * 	* The full range of spacing (the CFWS syntax element) is not supported,
 * 	  such as breaking addresses across lines.
 * 	* No unicode normalization is performed.
 * 	* The special characters ()[]:;@\, are allowed to appear unquoted in names.
 * ```
 */
namespace mail {
 /**
  * Address represents a single mail address.
  * An address such as "Barry Gibbs <bg@example.com>" is represented
  * as Address{Name: "Barry Gibbs", Address: "bg@example.com"}.
  */
 interface Address {
  name: string // Proper name; may be empty.
  address: string // user@domain
 }
 interface Address {
  /**
   * String formats the address as a valid RFC 5322 address.
   * If the address's name contains non-ASCII characters
   * the name will be rendered according to RFC 2047.
   */
  string(): string
 }
}

namespace search {
}
