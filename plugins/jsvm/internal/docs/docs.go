package main

import (
	"log"
	"os"
	"reflect"

	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/tygoja"
)

const heading = `
// -------------------------------------------------------------------
// baseBinds
// -------------------------------------------------------------------

declare var $app: pocketbase.PocketBase

/**
 * $arrayOf creates a placeholder array of the specified models.
 * Usually used to populate DB result into an array of models.
 *
 * Example:
 *
 * ` + "```" + `js
 * const records = $arrayOf(new Record)
 *
 * $app.dao().recordQuery(collection).limit(10).all(records)
 * ` + "```" + `
 */
declare function $arrayOf<T>(model: T): Array<T>;

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
 */
declare class DynamicModel {
  constructor(shape?: { [key:string]: any })
}

interface Record extends models.Record{} // merge
declare class Record implements models.Record {
  constructor(collection?: models.Collection, data?: { [key:string]: any })
}

interface Collection extends models.Collection{} // merge
declare class Collection implements models.Collection {
  constructor(data?: Partial<models.Collection>)
}

interface Admin extends models.Admin{} // merge
declare class Admin implements models.Admin {
  constructor(data?: Partial<models.Admin>)
}

interface Schema extends schema.Schema{} // merge
declare class Schema implements schema.Schema {
  constructor(data?: Partial<schema.Schema>)
}

interface SchemaField extends schema.SchemaField{} // merge
declare class SchemaField implements schema.SchemaField {
  constructor(data?: Partial<schema.SchemaField>)
}

interface MailerMessage extends mailer.Message{} // merge
declare class MailerMessage implements mailer.Message {
  constructor(message?: Partial<mailer.Message>)
}

interface Command extends cobra.Command{} // merge
declare class Command implements cobra.Command {
  constructor(cmd?: Partial<cobra.Command>)
}

interface ValidationError extends ozzo_validation.Error{} // merge
declare class ValidationError implements ozzo_validation.Error {
  constructor(code?: number, message?: string)
}

interface Dao extends daos.Dao{} // merge
declare class Dao implements daos.Dao {
  constructor(concurrentDB?: dbx.Builder, nonconcurrentDB?: dbx.Builder)
}

// -------------------------------------------------------------------
// dbxBinds
// -------------------------------------------------------------------

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

declare namespace $security {
  let randomString:                   security.randomString
  let randomStringWithAlphabet:       security.randomStringWithAlphabet
  let pseudorandomString:             security.pseudorandomString
  let pseudorandomStringWithAlphabet: security.pseudorandomStringWithAlphabet
  let parseUnverifiedToken:           security.parseUnverifiedJWT
  let parseToken:                     security.parseJWT
  let createToken:                    security.newToken
}

// -------------------------------------------------------------------
// filesystemBinds
// -------------------------------------------------------------------

declare namespace $filesystem {
  let fileFromPath:      filesystem.newFileFromPath
  let fileFromBytes:     filesystem.newFileFromBytes
  let fileFromMultipart: filesystem.newFileFromMultipart
}

// -------------------------------------------------------------------
// formsBinds
// -------------------------------------------------------------------

interface AdminLoginForm extends forms.AdminLogin{} // merge
declare class AdminLoginForm implements forms.AdminLogin {
  constructor(app: core.App)
}

interface AdminPasswordResetConfirmForm extends forms.AdminPasswordResetConfirm{} // merge
declare class AdminPasswordResetConfirmForm implements forms.AdminPasswordResetConfirm {
  constructor(app: core.App)
}

interface AdminPasswordResetRequestForm extends forms.AdminPasswordResetRequest{} // merge
declare class AdminPasswordResetRequestForm implements forms.AdminPasswordResetRequest {
  constructor(app: core.App)
}

interface AdminUpsertForm extends forms.AdminUpsert{} // merge
declare class AdminUpsertForm implements forms.AdminUpsert {
  constructor(app: core.App, admin: models.Admin)
}

interface AppleClientSecretCreateForm extends forms.AppleClientSecretCreate{} // merge
declare class AppleClientSecretCreateForm implements forms.AppleClientSecretCreate {
  constructor(app: core.App)
}

interface CollectionUpsertForm extends forms.CollectionUpsert{} // merge
declare class CollectionUpsertForm implements forms.CollectionUpsert {
  constructor(app: core.App, collection: models.Collection)
}

interface CollectionsImportForm extends forms.CollectionsImport{} // merge
declare class CollectionsImportForm implements forms.CollectionsImport {
  constructor(app: core.App)
}

interface RealtimeSubscribeForm extends forms.RealtimeSubscribe{} // merge
declare class RealtimeSubscribeForm implements forms.RealtimeSubscribe {}

interface RecordEmailChangeConfirmForm extends forms.RecordEmailChangeConfirm{} // merge
declare class RecordEmailChangeConfirmForm implements forms.RecordEmailChangeConfirm {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordEmailChangeRequestForm extends forms.RecordEmailChangeRequest{} // merge
declare class RecordEmailChangeRequestForm implements forms.RecordEmailChangeRequest {
  constructor(app: core.App, record: models.Record)
}

interface RecordOAuth2LoginForm extends forms.RecordOAuth2Login{} // merge
declare class RecordOAuth2LoginForm implements forms.RecordOAuth2Login {
  constructor(app: core.App, collection: models.Collection, optAuthRecord?: models.Record)
}

interface RecordPasswordLoginForm extends forms.RecordPasswordLogin{} // merge
declare class RecordPasswordLoginForm implements forms.RecordPasswordLogin {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordPasswordResetConfirmForm extends forms.RecordPasswordResetConfirm{} // merge
declare class RecordPasswordResetConfirmForm implements forms.RecordPasswordResetConfirm {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordPasswordResetRequestForm extends forms.RecordPasswordResetRequest{} // merge
declare class RecordPasswordResetRequestForm implements forms.RecordPasswordResetRequest {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordUpsertForm extends forms.RecordUpsert{} // merge
declare class RecordUpsertForm implements forms.RecordUpsert {
  constructor(app: core.App, record: models.Record)
}

interface RecordVerificationConfirmForm extends forms.RecordVerificationConfirm{} // merge
declare class RecordVerificationConfirmForm implements forms.RecordVerificationConfirm {
  constructor(app: core.App, collection: models.Collection)
}

interface RecordVerificationRequestForm extends forms.RecordVerificationRequest{} // merge
declare class RecordVerificationRequestForm implements forms.RecordVerificationRequest {
  constructor(app: core.App, collection: models.Collection)
}

interface SettingsUpsertForm extends forms.SettingsUpsert{} // merge
declare class SettingsUpsertForm implements forms.SettingsUpsert {
  constructor(app: core.App)
}

interface TestEmailSendForm extends forms.TestEmailSend{} // merge
declare class TestEmailSendForm implements forms.TestEmailSend {
  constructor(app: core.App)
}

interface TestS3FilesystemForm extends forms.TestS3Filesystem{} // merge
declare class TestS3FilesystemForm implements forms.TestS3Filesystem {
  constructor(app: core.App)
}

// -------------------------------------------------------------------
// apisBinds
// -------------------------------------------------------------------

interface Route extends echo.Route{} // merge
declare class Route implements echo.Route {
  constructor(data?: Partial<echo.Route>)
}

interface ApiError extends apis.ApiError{} // merge
declare class ApiError implements apis.ApiError {
  constructor(status?: number, message?: string, data?: any)
}

interface NotFoundError extends apis.ApiError{} // merge
declare class NotFoundError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

interface BadRequestError extends apis.ApiError{} // merge
declare class BadRequestError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

interface ForbiddenError extends apis.ApiError{} // merge
declare class ForbiddenError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

interface UnauthorizedError extends apis.ApiError{} // merge
declare class UnauthorizedError implements apis.ApiError {
  constructor(message?: string, data?: any)
}

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
// migrate only
// -------------------------------------------------------------------

/**
 * Migrate defines a single migration upgrade/downgrade action.
 *
 * Note that this method is available only in pb_migrations context.
 */
declare function migrate(
  up: (db: dbx.Builder) => void,
  down?: (db: dbx.Builder) => void
): void;
`

func main() {
	mapper := &jsvm.FieldMapper{}

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
		Heading:              heading,
	})

	result, err := gen.Generate()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./generated/types.d.ts", []byte(result), 0644); err != nil {
		log.Fatal(err)
	}
}
