## v0.11.2

- Fixed `fs.DeleteByPrefix()` hang on invalid S3 settings ([#1575](https://github.com/pocketbase/pocketbase/discussions/1575#discussioncomment-4661089)).

- Updated file(s) delete to run in the background on record/collection delete to avoid blocking the delete model transaction.
  _Currently the cascade files delete operation is treated as "non-critical" and in case of an error it is just logged during debug._
  _This will be improved in the near future with the planned async job queue implementation._


## v0.11.1

- Unescaped path parameter values ([#1552](https://github.com/pocketbase/pocketbase/issues/1552)).


## v0.11.0

- Added `+` and `-` body field modifiers for `number`, `files`, `select` and `relation` fields.
  ```js
  {
    // oldValue + 2
    "someNumber+": 2,

    // oldValue + ["id1", "id2"] - ["id3"]
    "someRelation+": ["id1", "id2"],
    "someRelation-": ["id3"],

    // delete single file by its name (file fields supports only the "-" modifier!)
    "someFile-": "filename.png",
  }
  ```
  _Note1: `@request.data.someField` will contain the final resolved value._

  _Note2: The old index (`"field.0":null`) and filename (`"field.filename.png":null`) based suffixed syntax for deleting files is still supported._

- ! Added support for multi-match/match-all request data and collection multi-valued fields (`select`, `relation`) conditions.
  If you want a "at least one of" type of condition, you can prefix the operator with `?`.
  ```js
  // for each someRelA.someRelB record require the "status" field to be "active"
  someRelA.someRelB.status = "active"

  // OR for "at least one of" condition
  someRelA.someRelB.status ?= "active"
  ```
  _**Note: Previously the behavior for multi-valued fields was as the "at least one of" type.
  The release comes with system db migration that will update your existing API rules (if needed) to preserve the compatibility.
  If you have multi-select or multi-relation filter checks in your client-side code and want to preserve the old behavior, you'll have to prefix with `?` your operators.**_

- Added support for querying `@request.data.someRelField.*` relation fields.
  ```js
  // example submitted data: {"someRel": "REL_RECORD_ID"}
  @request.data.someRel.status = "active"
  ```

- Added `:isset` modifier for the static request data fields.
  ```js
  // prevent changing the "role" field
  @request.data.role:isset = false
  ```

- Added `:length` modifier for the arrayable request data and collection fields (`select`, `file`, `relation`).
  ```js
  // example submitted data: {"someSelectField": ["val1", "val2"]}
  @request.data.someSelectField:length = 2

  // check existing record field length
  someSelectField:length = 2
  ```

- Added `:each` modifier support for the multi-`select` request data and collection field.
  ```js
  // check if all selected rows has "pb_" prefix
  roles:each ~ 'pb_%'
  ```

- Improved the Admin UI filters autocomplete.

- Added `@random` sort key for `RANDOM()` sorted list results.

- Added Strava OAuth2 provider ([#1443](https://github.com/pocketbase/pocketbase/pull/1443); thanks @szsascha).

- Added Gitee OAuth2 provider ([#1448](https://github.com/pocketbase/pocketbase/pull/1448); thanks @yuxiang-gao).

- Added IME status check to the textarea keydown handler ([#1370](https://github.com/pocketbase/pocketbase/pull/1370); thanks @tenthree).

- Added `filesystem.NewFileFromBytes()` helper ([#1420](https://github.com/pocketbase/pocketbase/pull/1420); thanks @dschissler).

- Added support for reordering uploaded multiple files.

- Added `webp` to the default images mime type presets list ([#1469](https://github.com/pocketbase/pocketbase/pull/1469); thanks @khairulhaaziq).

- Added the OAuth2 refresh token to the auth meta response ([#1487](https://github.com/pocketbase/pocketbase/issues/1487)).

- Fixed the text wrapping in the Admin UI listing searchbar ([#1416](https://github.com/pocketbase/pocketbase/issues/1416)).

- Fixed number field value output in the records listing ([#1447](https://github.com/pocketbase/pocketbase/issues/1447)).

- Fixed duplicated settings view pages caused by uncompleted transitions ([#1498](https://github.com/pocketbase/pocketbase/issues/1498)).

- Allowed sending `Authorization` header with the `/auth-with-password` record and admin login requests ([#1494](https://github.com/pocketbase/pocketbase/discussions/1494)).

- `migrate down` now reverts migrations in the applied order.

- Added additional list-bucket check in the S3 config test API.

- Other minor improvements.


## v0.10.4

- Fixed `Record.MergeExpand` panic when the main model expand map is not initialized ([#1365](https://github.com/pocketbase/pocketbase/issues/1365)).


## v0.10.3

- ! Renamed the metadata key `original_filename` to `original-filename` due to an S3 file upload error caused by the underscore character ([#1343](https://github.com/pocketbase/pocketbase/pull/1343); thanks @yuxiang-gao).

- Fixed request verification docs api url ([#1332](https://github.com/pocketbase/pocketbase/pull/1332); thanks @JoyMajumdar2001)

- Excluded `collectionId` and `collectionName` from the displayable relation props list ([1322](https://github.com/pocketbase/pocketbase/issues/1322); thanks @dhall2).


## v0.10.2

- Fixed nested multiple expands with shared path ([#586](https://github.com/pocketbase/pocketbase/issues/586#issuecomment-1357784227)).
  A new helper method `models.Record.MergeExpand(map[string]any)` was also added to simplify the expand handling and unit testing.


## v0.10.1

- Fixed nested transactions deadlock when authenticating with OAuth2 ([#1291](https://github.com/pocketbase/pocketbase/issues/1291)).


## v0.10.0

- Added `/api/health` endpoint (thanks @MarvinJWendt).

- Added support for SMTP `LOGIN` auth for Microsoft/Outlook and other providers that don't support the `PLAIN` auth method ([#1217](https://github.com/pocketbase/pocketbase/discussions/1217#discussioncomment-4387970)).

- Reduced memory consumption (you can expect ~20% less allocated memory).

- Added support for split (concurrent and nonconcurrent) DB connections pool increasing even further the concurrent throughput without blocking reads on heavy write load.

- Improved record references delete performance.

- Removed the unnecessary parenthesis in the generated filter SQL query, reducing the "_parse stack overflow_" errors.

- Fixed `~` expressions backslash literal escaping ([#1231](https://github.com/pocketbase/pocketbase/discussions/1231)).

- Refactored the `core.app.Bootstrap()` to be called before starting the cobra commands ([#1267](https://github.com/pocketbase/pocketbase/discussions/1267)).

- ! Changed `pocketbase.NewWithConfig(config Config)` to `pocketbase.NewWithConfig(config *Config)` and added 4 new config settings:
  ```go
  DataMaxOpenConns int // default to core.DefaultDataMaxOpenConns
  DataMaxIdleConns int // default to core.DefaultDataMaxIdleConns
  LogsMaxOpenConns int // default to core.DefaultLogsMaxOpenConns
  LogsMaxIdleConns int // default to core.DefaultLogsMaxIdleConns
  ```

- Added new helper method `core.App.IsBootstrapped()` to check the current app bootstrap state.

- ! Changed `core.NewBaseApp(dir, encryptionEnv, isDebug)` to `NewBaseApp(config *BaseAppConfig)`.

- ! Removed `rest.UploadedFile` struct (see below `filesystem.File`).

- Added generic file resource struct that allows loading and uploading file content from
  different sources (at the moment multipart/form-data requests and from the local filesystem).
  ```
  filesystem.File{}
  filesystem.NewFileFromPath(path)
  filesystem.NewFileFromMultipart(multipartHeader)
  filesystem/System.UploadFile(file)
  ```

- Refactored `forms.RecordUpsert` to allow more easily loading and removing files programmatically.
  ```
  forms.RecordUpsert.AddFiles(key, filesystem.File...) // add new filesystem.File to the form for upload
  forms.RecordUpsert.RemoveFiles(key, filenames...)     // marks the filenames for deletion
  ```

- Trigger the `password` validators if any of the others password change fields is set.


## v0.9.2

- Fixed field column name conflict on record deletion ([#1220](https://github.com/pocketbase/pocketbase/discussions/1220)).


## v0.9.1

- Moved the record file upload and delete out of the db transaction to minimize the locking times.

- Added `Dao` query semaphore and base fail/retry handling to improve the concurrent writes throughput ([#1187](https://github.com/pocketbase/pocketbase/issues/1187)).

- Fixed records cascade deletion when there are "A<->B" relation references.

- Replaced `c.QueryString()` with `c.QueryParams().Encode()` to allow loading middleware modified query parameters in the default crud actions ([#1210](https://github.com/pocketbase/pocketbase/discussions/1210)).

- Fixed the datetime field not triggering the `onChange` event on manual field edit and added a "Clear" button ([#1219](https://github.com/pocketbase/pocketbase/issues/1219)).

- Updated the GitHub goreleaser action to use go 1.19.4 since it comes with [some security fixes](https://github.com/golang/go/issues?q=milestone%3AGo1.19.4+label%3ACherryPickApproved).


## v0.9.0

- Fixed concurrent multi-relation cascade update/delete ([#1138](https://github.com/pocketbase/pocketbase/issues/1138)).

- Added the raw OAuth2 user data (`meta.rawUser`) and OAuth2 access token (`meta.accessToken`) to the auth response ([#654](https://github.com/pocketbase/pocketbase/discussions/654)).

- `BaseModel.UnmarkAsNew()` method was renamed to `BaseModel.MarkAsNotNew()`.
  Additionally, to simplify the insert model queries with custom IDs, it is no longer required to call `MarkAsNew()` for manually initialized models with set ID since now this is the default state.
  When the model is populated with values from the database (eg. after row `Scan`) it will be marked automatically as "not new".

- Added `Record.OriginalCopy()` method that returns a new `Record` copy populated with the initially loaded record data (useful if you want to compare old and new field values).

- Added new event hooks:
  ```go
  app.OnBeforeBootstrap()
  app.OnAfterBootstrap()
  app.OnBeforeApiError()
  app.OnAfterApiError()
  app.OnRealtimeDisconnectRequest()
  app.OnRealtimeBeforeMessageSend()
  app.OnRealtimeAfterMessageSend()
  app.OnRecordBeforeRequestPasswordResetRequest()
  app.OnRecordAfterRequestPasswordResetRequest()
  app.OnRecordBeforeConfirmPasswordResetRequest()
  app.OnRecordAfterConfirmPasswordResetRequest()
  app.OnRecordBeforeRequestVerificationRequest()
  app.OnRecordAfterRequestVerificationRequest()
  app.OnRecordBeforeConfirmVerificationRequest()
  app.OnRecordAfterConfirmVerificationRequest()
  app.OnRecordBeforeRequestEmailChangeRequest()
  app.OnRecordAfterRequestEmailChangeRequest()
  app.OnRecordBeforeConfirmEmailChangeRequest()
  app.OnRecordAfterConfirmEmailChangeRequest()
  ```

- The original uploaded file name is now stored as metadata under the `original_filename` key. It could be accessed via:
  ```go
  fs, _ := app.NewFilesystem()
  defer fs.Close()

  attrs, _ := fs.Attributes(fikeKey)
  attrs.Metadata["original_name"]
  ```

- Added support for `Partial/Range` file requests ([#1125](https://github.com/pocketbase/pocketbase/issues/1125)).
  This is a minor breaking change if you are using `filesystem.Serve` (eg. as part of a custom `OnFileDownloadRequest` hook):
  ```go
  // old
  filesystem.Serve(res, e.ServedPath, e.ServedName)

  // new
  filesystem.Serve(res, req, e.ServedPath, e.ServedName)
  ```

- Refactored the `migrate` command to support **external JavaScript migration files** using an embedded JS interpreter ([goja](https://github.com/dop251/goja)).
  This allow writing custom migration scripts such as programmatically creating collections,
  initializing default settings, running data imports, etc., with a JavaScript API very similar to the Go one (_more documentation will be available soon_).

  The `migrate` command is available by default for the prebuilt executable,
  but if you use PocketBase as framework you need register it manually:
  ```go
  migrationsDir := "" // default to "pb_migrations" (for js) and "migrations" (for go)

  // load js files if you want to allow loading external JavaScript migrations
  jsvm.MustRegisterMigrations(app, &jsvm.MigrationsOptions{
    Dir: migrationsDir,
  })

  // register the `migrate` command
  migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
    TemplateLang: migratecmd.TemplateLangJS, // or migratecmd.TemplateLangGo (default)
    Dir:          migrationsDir,
    Automigrate:  true,
  })
  ```

  **The refactoring also comes with automigrations support.**

  If `Automigrate` is enabled (`true` by default for the prebuilt executable; can be disabled with `--automigrate=0`),
  PocketBase will generate seamlessly in the background JS (or Go) migration file with your collection changes.
  **The directory with the JS migrations can be committed to your git repo.**
  All migrations (Go and JS) are automatically executed on server start.
  Also note that the auto generated migrations are granural (in contrast to the `migrate collections` snapshot command)
  and allow multiple developers to do changes on the collections independently (even editing the same collection) miniziming the eventual merge conflicts.
  Here is a sample JS migration file that will be generated if you for example edit a single collection name:
  ```js
  // pb_migrations/1669663597_updated_posts_old.js
  migrate((db) => {
    // up
    const dao = new Dao(db)
    const collection = dao.findCollectionByNameOrId("lngf8rb3dqu86r3")
    collection.name = "posts_new"
    return dao.saveCollection(collection)
  }, (db) => {
    // down
    const dao = new Dao(db)
    const collection = dao.findCollectionByNameOrId("lngf8rb3dqu86r3")
    collection.name = "posts_old"
    return dao.saveCollection(collection)
  })
  ```

- Added new `Dao` helpers to make it easier fetching and updating the app settings from a migration:
  ```go
  dao.FindSettings([optEncryptionKey])
  dao.SaveSettings(newSettings, [optEncryptionKey])
  ```

- Moved `core.Settings` to `models/settings.Settings`:
  ```
  core.Settings{}           -> settings.Settings{}
  core.NewSettings()        -> settings.New()
  core.MetaConfig{}         -> settings.MetaConfig{}
  core.LogsConfig{}         -> settings.LogsConfig{}
  core.SmtpConfig{}         -> settings.SmtpConfig{}
  core.S3Config{}           -> settings.S3Config{}
  core.TokenConfig{}        -> settings.TokenConfig{}
  core.AuthProviderConfig{} -> settings.AuthProviderConfig{}
  ```

- Changed the `mailer.Mailer` interface (**minor breaking if you are sending custom emails**):
  ```go
  // Old:
  app.NewMailClient().Send(from, to, subject, html, attachments?)

  // New:
  app.NewMailClient().Send(&mailer.Message{
    From: from,
    To: to,
    Subject: subject,
    HTML: html,
    Attachments: attachments,
    // new configurable fields
    Bcc: []string{"bcc1@example.com", "bcc2@example.com"},
    Cc: []string{"cc1@example.com", "cc2@example.com"},
    Headers: map[string]string{"Custom-Header": "test"},
    Text: "custom plain text version",
  })
  ```
  The new `*mailer.Message` struct is also now a member of the `MailerRecordEvent` and `MailerAdminEvent` events.

- Other minor UI fixes and improvements


## v0.8.0

**⚠️ This release contains breaking changes and requires some manual migration steps!**

The biggest change is the merge of the `User` models and the `profiles` collection per [#376](https://github.com/pocketbase/pocketbase/issues/376).
There is no longer `user` type field and the users are just an "auth" collection (we now support **collection types**, currently only "base" and "auth").
This should simplify the users management and at the same time allow us to have unlimited multiple "auth" collections each with their own custom fields and authentication options (eg. staff, client, etc.).

In addition to the `Users` and `profiles` merge, this release comes with several other improvements:

- Added indirect expand support [#312](https://github.com/pocketbase/pocketbase/issues/312#issuecomment-1242893496).

- The `json` field type now supports filtering and sorting [#423](https://github.com/pocketbase/pocketbase/issues/423#issuecomment-1258302125).

- The `relation` field now allows unlimited `maxSelect` (aka. without upper limit).

- Added support for combined email/username + password authentication (see below `authWithPassword()`).

- Added support for full _"manager-subordinate"_ users management, including a special API rule to allow directly changing system fields like email, password, etc. without requiring `oldPassword` or other user verification.

- Enabled OAuth2 account linking on authorized request from the same auth collection (_this is useful for example if the OAuth2 provider doesn't return an email and you want to associate it with the current logged in user_).

- Added option to toggle the record columns visibility from the table listing.

- Added support for collection schema fields reordering.

- Added several new OAuth2 providers (Microsoft Azure AD, Spotify, Twitch, Kakao).

- Improved memory usage on large file uploads [#835](https://github.com/pocketbase/pocketbase/discussions/835).

- More detailed API preview docs and site documentation (the repo is located at https://github.com/pocketbase/site).

- Other minor performance improvements (mostly related to the search apis).

### Migrate from v0.7.x

- **[Data](#data)**
- **[SDKs](#sdks)**
- **[API](#api)**
- **[Internals](#internals)**

#### Data

The merge of users and profiles comes with several required db changes.
The easiest way to apply them is to use the new temporary `upgrade` command:

```sh
# make sure to have a copy of your pb_data in case something fails
cp -r ./pb_data ./pb_data_backup

# run the upgrade command
./pocketbase08 upgrade

# start the application as usual
./pocketbase08 serve
```

The upgrade command:

- Creates a new `users` collection with merged fields from the `_users` table and the `profiles` collection.
  The new user records will have the ids from the `profiles` collection.
- Changes all `user` type fields to `relation` and update the references to point to the new user ids.
- Renames all `@collection.profiles.*`, `@request.user.*` and `@request.user.profile.*` filters to `@collection.users.*` and `@request.auth.*`.
- Appends `2` to all **schema field names** and **api filter rules** that conflicts with the new system reserved ones:
  ```
  collectionId   => collectionId2
  collectionName => collectionName2
  expand         => expand2

  // only for the "profiles" collection fields:
  username               => username2
  email                  => email2
  emailVisibility        => emailVisibility2
  verified               => verified2
  tokenKey               => tokenKey2
  passwordHash           => passwordHash2
  lastResetSentAt        => lastResetSentAt2
  lastVerificationSentAt => lastVerificationSentAt2
  ```

#### SDKs

Please check the individual SDK package changelog and apply the necessary changes in your code:

- [**JavaScript SDK changelog**](https://github.com/pocketbase/js-sdk/blob/master/CHANGELOG.md)
  ```sh
  npm install pocketbase@latest --save
  ```

- [**Dart SDK changelog**](https://github.com/pocketbase/dart-sdk/blob/master/CHANGELOG.md)

  ```sh
  dart pub add pocketbase:^0.5.0
  # or with Flutter:
  flutter pub add pocketbase:^0.5.0
  ```

#### API

> _**You don't have to read this if you are using an official SDK.**_

- The authorization schema is no longer necessary. Now it is auto detected from the JWT token payload:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>Authorization: Admin TOKEN</td>
      <td>Authorization: TOKEN</td>
    </tr>
    <tr valign="top">
      <td>Authorization: User TOKEN</td>
      <td>Authorization: TOKEN</td>
    </tr>
  </table>

- All datetime stings are now returned in ISO8601 format - with _Z_ suffix and space as separator between the date and time part:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>2022-01-02 03:04:05.678</td>
      <td>2022-01-02 03:04:05.678<strong>Z</strong></td>
    </tr>
  </table>

- Removed the `@` prefix from the system record fields for easier json parsing:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td><strong>@</strong>collectionId</td>
      <td>collectionId</td>
    </tr>
    <tr valign="top">
      <td><strong>@</strong>collectionName</td>
      <td>collectionName</td>
    </tr>
    <tr valign="top">
      <td><strong>@</strong>expand</td>
      <td>expand</td>
    </tr>
  </table>

- All users api handlers are moved under `/api/collections/:collection/`:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>
        <em>GET /api/<strong>users</strong>/auth-methods</em>
      </td>
      <td>
        <em>GET /api/<strong>collections/:collection</strong>/auth-methods</em>
      </td>
    </tr>
    <tr valign="top">
      <td>
        <em>POST /api/<strong>users/refresh</strong></em>
      </td>
      <td>
        <em>POST /api/<strong>collections/:collection/auth-refresh</strong></em>
      </td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users/auth-via-oauth2</strong></em></td>
      <td>
        <em>POST /api/<strong>collections/:collection/auth-with-oauth2</strong></em>
        <br/>
        <em>You can now also pass optional <code>createData</code> object on OAuth2 sign-up.</em>
        <br/>
        <em>Also please note that now required user/profile fields are properly validated when creating new auth model on OAuth2 sign-up.</em>
      </td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users/auth-via-email</strong></em></td>
      <td>
        <em>POST /api/<strong>collections/:collection/auth-with-password</strong></em>
        <br/>
        <em>Handles username/email + password authentication.</em>
        <br/>
        <code>{"identity": "usernameOrEmail", "password": "123456"}</code>
      </td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong>/request-password-reset</em></td>
      <td><em>POST /api/<strong>collections/:collection</strong>/request-password-reset</em></td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong>/confirm-password-reset</em></td>
      <td><em>POST /api/<strong>collections/:collection</strong>/confirm-password-reset</em></td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong>/request-verification</em></td>
      <td><em>POST /api/<strong>collections/:collection</strong>/request-verification</em></td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong>/confirm-verification</em></td>
      <td><em>POST /api/<strong>collections/:collection</strong>/confirm-verification</em></td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong>/request-email-change</em></td>
      <td><em>POST /api/<strong>collections/:collection</strong>/request-email-change</em></td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong>/confirm-email-change</em></td>
      <td><em>POST /api/<strong>collections/:collection</strong>/confirm-email-change</em></td>
    </tr>
    <tr valign="top">
      <td><em>GET /api/<strong>users</strong></em></td>
      <td><em>GET /api/<strong>collections/:collection/records</strong></em></td>
    </tr>
    <tr valign="top">
      <td><em>GET /api/<strong>users</strong>/:id</em></td>
      <td><em>GET /api/<strong>collections/:collection/records</strong>/:id</em></td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/<strong>users</strong></em></td>
      <td><em>POST /api/<strong>collections/:collection/records</strong></em></td>
    </tr>
    <tr valign="top">
      <td><em>PATCH /api/<strong>users</strong>/:id</em></td>
      <td><em>PATCH /api/<strong>collections/:collection/records</strong>/:id</em></td>
    </tr>
    <tr valign="top">
      <td><em>DELETE /api/<strong>users</strong>/:id</em></td>
      <td><em>DELETE /api/<strong>collections/:collection/records</strong>/:id</em></td>
    </tr>
    <tr valign="top">
      <td><em>GET /api/<strong>users</strong>/:id/external-auths</em></td>
      <td><em>GET /api/<strong>collections/:collection/records</strong>/:id/external-auths</em></td>
    </tr>
    <tr valign="top">
      <td><em>DELETE /api/<strong>users</strong>/:id/external-auths/:provider</em></td>
      <td><em>DELETE /api/<strong>collections/:collection/records</strong>/:id/external-auths/:provider</em></td>
    </tr>
  </table>

  _In relation to the above changes, the `user` property in the auth response is renamed to `record`._

- The admins api was also updated for consistency with the users api changes:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>
        <em>POST /api/admins/<strong>refresh</strong></em>
      </td>
      <td>
        <em>POST /api/admins/<strong>auth-refresh</strong></em>
      </td>
    </tr>
    <tr valign="top">
      <td><em>POST /api/admins/<strong>auth-via-email</strong></em></td>
      <td>
        <em>POST /api/admins/<strong>auth-with-password</strong></em>
        <br />
        <code>{"identity": "test@example.com", "password": "123456"}</code>
        <br />
        (notice that the <code>email</code> body field was renamed to <code>identity</code>)
      </td>
    </tr>
  </table>

- To prevent confusion with the auth method responses, the following endpoints now returns 204 with empty body (previously 200 with token and auth model):
  ```
  POST /api/admins/confirm-password-reset
  POST /api/collections/:collection/confirm-password-reset
  POST /api/collections/:collection/confirm-verification
  POST /api/collections/:collection/confirm-email-change
  ```

- Renamed the "user" related settings fields returned by `GET /api/settings`:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td><strong>user</strong>AuthToken</td>
      <td><strong>record</strong>AuthToken</td>
    </tr>
    <tr valign="top">
      <td><strong>user</strong>PasswordResetToken</td>
      <td><strong>record</strong>PasswordResetToken</td>
    </tr>
    <tr valign="top">
      <td><strong>user</strong>EmailChangeToken</td>
      <td><strong>record</strong>EmailChangeToken</td>
    </tr>
    <tr valign="top">
      <td><strong>user</strong>VerificationToken</td>
      <td><strong>record</strong>VerificationToken</td>
    </tr>
  </table>

#### Internals

> _**You don't have to read this if you are not using PocketBase as framework.**_

- Removed `forms.New*WithConfig()` factories to minimize ambiguities.
  If you need to pass a transaction Dao you can use the new `SetDao(dao)` method available to the form instances.

- `forms.RecordUpsert.LoadData(data map[string]any)` now can bulk load external data from a map.
  To load data from a request instance, you could use `forms.RecordUpsert.LoadRequest(r, optKeysPrefix = "")`.

- `schema.RelationOptions.MaxSelect` has new type `*int` (_you can use the new `types.Pointer(123)` helper to assign pointer values_).

- Renamed the constant `apis.ContextUserKey` (_"user"_) to `apis.ContextAuthRecordKey` (_"authRecord"_).

- Replaced user related middlewares with their auth record alternative:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>apis.Require<strong>User</strong>Auth()</td>
      <td>apis.Require<strong>Record</strong>Auth(<strong>optCollectionNames ...string</strong>)</td>
    </tr>
    <tr valign="top">
      <td>apis.RequireAdminOr<strong>User</strong>Auth()</td>
      <td>apis.RequireAdminOr<strong>Record</strong>Auth(<strong>optCollectionNames ...string</strong>)</td>
    </tr>
    <tr valign="top">
      <td>N/A</td>
      <td>
        <strong>RequireSameContextRecordAuth()</strong>
        <br/>
        <em>(requires the auth record to be from the same context collection)</em>
      </td>
    </tr>
  </table>

- The following record Dao helpers now uses the collection id or name instead of `*models.Collection` instance to reduce the verbosity when fetching records:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>dao.FindRecordById(<strong>collection</strong>, ...)</td>
      <td>dao.FindRecordById(<strong>collectionNameOrId</strong>, ...)</td>
    </tr>
    <tr valign="top">
      <td>dao.FindRecordsByIds(<strong>collection</strong>, ...)</td>
      <td>dao.FindRecordsByIds(<strong>collectionNameOrId</strong>, ...)</td>
    </tr>
    <tr valign="top">
      <td>dao.FindRecordsByExpr(<strong>collection</strong>, ...)</td>
      <td>dao.FindRecordsByExpr(<strong>collectionNameOrId</strong>, ...)</td>
    </tr>
    <tr valign="top">
      <td>dao.FindFirstRecordByData(<strong>collection</strong>, ...)</td>
      <td>dao.FindFirstRecordByData(<strong>collectionNameOrId</strong>, ...)</td>
    </tr>
    <tr valign="top">
      <td>dao.IsRecordValueUnique(<strong>collection</strong>, ...)</td>
      <td>dao.IsRecordValueUnique(<strong>collectionNameOrId</strong>, ...)</td>
    </tr>
  </table>

- Replaced all User related Dao helpers with Record equivalents:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>dao.UserQuery()</td>
      <td>dao.RecordQuery(collection)</td>
    </tr>
    <tr valign="top">
      <td>dao.FindUserById(id)</td>
      <td>dao.FindRecordById(collectionNameOrId, id)</td>
    </tr>
    <tr valign="top">
      <td>dao.FindUserByToken(token, baseKey)</td>
      <td>dao.FindAuthRecordByToken(token, baseKey)</td>
    </tr>
    <tr valign="top">
      <td>dao.FindUserByEmail(email)</td>
      <td>dao.FindAuthRecordByEmail(collectionNameOrId, email)</td>
    </tr>
    <tr valign="top">
      <td>N/A</td>
      <td>dao.FindAuthRecordByUsername(collectionNameOrId, username)</td>
    </tr>
  </table>

- Moved the formatted `ApiError` struct and factories to the `github.com/pocketbase/pocketbase/apis` subpackage:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td colspan="2"><em>Import path</em></td>
    </tr>
    <tr valign="top">
      <td>github.com/pocketbase/pocketbase/<strong>tools/rest</strong></td>
      <td>github.com/pocketbase/pocketbase/<strong>apis</strong></td>
    </tr>
    <tr valign="top">
      <td colspan="2"><em>Fields</em></td>
    </tr>
    <tr valign="top">
      <td><strong>rest</strong>.ApiError{}</td>
      <td><strong>apis</strong>.ApiError{}</td>
    </tr>
    <tr valign="top">
      <td><strong>rest</strong>.NewNotFoundError()</td>
      <td><strong>apis</strong>.NewNotFoundError()</td>
    </tr>
    <tr valign="top">
      <td><strong>rest</strong>.NewBadRequestError()</td>
      <td><strong>apis</strong>.NewBadRequestError()</td>
    </tr>
    <tr valign="top">
      <td><strong>rest</strong>.NewForbiddenError()</td>
      <td><strong>apis</strong>.NewForbiddenError()</td>
    </tr>
    <tr valign="top">
      <td><strong>rest</strong>.NewUnauthorizedError()</td>
      <td><strong>apis</strong>.NewUnauthorizedError()</td>
    </tr>
    <tr valign="top">
      <td><strong>rest</strong>.NewApiError()</td>
      <td><strong>apis</strong>.NewApiError()</td>
    </tr>
  </table>

- Renamed `models.Record` helper getters:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>Set<strong>DataValue</strong></td>
      <td>Set</td>
    </tr>
    <tr valign="top">
      <td>Get<strong>DataValue</strong></td>
      <td>Get</td>
    </tr>
    <tr valign="top">
      <td>GetBool<strong>DataValue</strong></td>
      <td>GetBool</td>
    </tr>
    <tr valign="top">
      <td>GetString<strong>DataValue</strong></td>
      <td>GetString</td>
    </tr>
    <tr valign="top">
      <td>GetInt<strong>DataValue</strong></td>
      <td>GetInt</td>
    </tr>
    <tr valign="top">
      <td>GetFloat<strong>DataValue</strong></td>
      <td>GetFloat</td>
    </tr>
    <tr valign="top">
      <td>GetTime<strong>DataValue</strong></td>
      <td>GetTime</td>
    </tr>
    <tr valign="top">
      <td>GetDateTime<strong>DataValue</strong></td>
      <td>GetDateTime</td>
    </tr>
    <tr valign="top">
      <td>GetStringSlice<strong>DataValue</strong></td>
      <td>GetStringSlice</td>
    </tr>
  </table>

- Added new auth collection `models.Record` helpers:
  ```go
  func (m *Record) Username() string
  func (m *Record) SetUsername(username string) error
  func (m *Record) Email() string
  func (m *Record) SetEmail(email string) error
  func (m *Record) EmailVisibility() bool
  func (m *Record) SetEmailVisibility(visible bool) error
  func (m *Record) IgnoreEmailVisibility(state bool)
  func (m *Record) Verified() bool
  func (m *Record) SetVerified(verified bool) error
  func (m *Record) TokenKey() string
  func (m *Record) SetTokenKey(key string) error
  func (m *Record) RefreshTokenKey() error
  func (m *Record) LastResetSentAt() types.DateTime
  func (m *Record) SetLastResetSentAt(dateTime types.DateTime) error
  func (m *Record) LastVerificationSentAt() types.DateTime
  func (m *Record) SetLastVerificationSentAt(dateTime types.DateTime) error
  func (m *Record) ValidatePassword(password string) bool
  func (m *Record) SetPassword(password string) error
  ```

- Added option to return serialized custom `models.Record` fields data:
  ```go
  func (m *Record) UnknownData() map[string]any
  func (m *Record) WithUnkownData(state bool)
  ```

- Deleted `model.User`. Now the user data is stored as an auth `models.Record`.
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>User.Email</td>
      <td>Record.Email()</td>
    </tr>
    <tr valign="top">
      <td>User.TokenKey</td>
      <td>Record.TokenKey()</td>
    </tr>
    <tr valign="top">
      <td>User.Verified</td>
      <td>Record.Verified()</td>
    </tr>
    <tr valign="top">
      <td>User.SetPassword()</td>
      <td>Record.SetPassword()</td>
    </tr>
    <tr valign="top">
      <td>User.RefreshTokenKey()</td>
      <td>Record.RefreshTokenKey()</td>
    </tr>
    <tr valign="top">
      <td colspan="2"><em>etc.</em></td>
    </tr>
  </table>

- Replaced `User` related event hooks with their `Record` alternative:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>OnMailerBefore<strong>User</strong>ResetPasswordSend() *hook.Hook[*Mailer<strong>User</strong>Event]</td>
      <td>OnMailerBefore<strong>Record</strong>ResetPasswordSend() *hook.Hook[*Mailer<strong>Record</strong>Event]</td>
    </tr>
    <tr valign="top">
      <td>OnMailerAfter<strong>User</strong>ResetPasswordSend() *hook.Hook[*Mailer<strong>User</strong>Event]</td>
      <td>OnMailerAfter<strong>Record</strong>ResetPasswordSend() *hook.Hook[*Mailer<strong>Record</strong>Event]</td>
    </tr>
    <tr valign="top">
      <td>OnMailerBefore<strong>User</strong>VerificationSend() *hook.Hook[*Mailer<strong>User</strong>Event]</td>
      <td>OnMailerBefore<strong>Record</strong>VerificationSend() *hook.Hook[*Mailer<strong>Record</strong>Event]</td>
    </tr>
    <tr valign="top">
      <td>OnMailerAfter<strong>User</strong>VerificationSend() *hook.Hook[*Mailer<strong>User</strong>Event]</td>
      <td>OnMailerAfter<strong>Record</strong>VerificationSend() *hook.Hook[*Mailer<strong>Record</strong>Event]</td>
    </tr>
    <tr valign="top">
      <td>OnMailerBefore<strong>User</strong>ChangeEmailSend() *hook.Hook[*Mailer<strong>User</strong>Event]</td>
      <td>OnMailerBefore<strong>Record</strong>ChangeEmailSend() *hook.Hook[*Mailer<strong>Record</strong>Event]</td>
    </tr>
    <tr valign="top">
      <td>OnMailerAfter<strong>User</strong>ChangeEmailSend() *hook.Hook[*Mailer<strong>User</strong>Event]</td>
      <td>OnMailerAfter<strong>Record</strong>ChangeEmailSend() *hook.Hook[*Mailer<strong>Record</strong>Event]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>Users</strong>ListRequest() *hook.Hook[*<strong>User</strong>ListEvent]</td>
      <td>On<strong>Records</strong>ListRequest() *hook.Hook[*<strong>Records</strong>ListEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>ViewRequest() *hook.Hook[*<strong>User</strong>ViewEvent]</td>
      <td>On<strong>Record</strong>ViewRequest() *hook.Hook[*<strong>Record</strong>ViewEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>BeforeCreateRequest() *hook.Hook[*<strong>User</strong>CreateEvent]</td>
      <td>On<strong>Record</strong>BeforeCreateRequest() *hook.Hook[*<strong>Record</strong>CreateEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>AfterCreateRequest() *hook.Hook[*<strong>User</strong>CreateEvent]</td>
      <td>On<strong>Record</strong>AfterCreateRequest() *hook.Hook[*<strong>Record</strong>CreateEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>BeforeUpdateRequest() *hook.Hook[*<strong>User</strong>UpdateEvent]</td>
      <td>On<strong>Record</strong>BeforeUpdateRequest() *hook.Hook[*<strong>Record</strong>UpdateEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>AfterUpdateRequest() *hook.Hook[*<strong>User</strong>UpdateEvent]</td>
      <td>On<strong>Record</strong>AfterUpdateRequest() *hook.Hook[*<strong>Record</strong>UpdateEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>BeforeDeleteRequest() *hook.Hook[*<strong>User</strong>DeleteEvent]</td>
      <td>On<strong>Record</strong>BeforeDeleteRequest() *hook.Hook[*<strong>Record</strong>DeleteEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>AfterDeleteRequest() *hook.Hook[*<strong>User</strong>DeleteEvent]</td>
      <td>On<strong>Record</strong>AfterDeleteRequest() *hook.Hook[*<strong>Record</strong>DeleteEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>AuthRequest() *hook.Hook[*<strong>User</strong>AuthEvent]</td>
      <td>On<strong>Record</strong>AuthRequest() *hook.Hook[*<strong>Record</strong>AuthEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>ListExternalAuths() *hook.Hook[*<strong>User</strong>ListExternalAuthsEvent]</td>
      <td>On<strong>Record</strong>ListExternalAuths() *hook.Hook[*<strong>Record</strong>ListExternalAuthsEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>BeforeUnlinkExternalAuthRequest() *hook.Hook[*<strong>User</strong>UnlinkExternalAuthEvent]</td>
      <td>On<strong>Record</strong>BeforeUnlinkExternalAuthRequest() *hook.Hook[*<strong>Record</strong>UnlinkExternalAuthEvent]</td>
    </tr>
    <tr valign="top">
      <td>On<strong>User</strong>AfterUnlinkExternalAuthRequest() *hook.Hook[*<strong>User</strong>UnlinkExternalAuthEvent]</td>
      <td>On<strong>Record</strong>AfterUnlinkExternalAuthRequest() *hook.Hook[*<strong>Record</strong>UnlinkExternalAuthEvent]</td>
    </tr>
  </table>

- Replaced `forms.UserEmailLogin{}` with `forms.RecordPasswordLogin{}` (for both username and email depending on which is enabled for the collection).

- Renamed user related `core.Settings` fields:
  <table class="d-table" width="100%">
    <tr>
      <th>Old</th>
      <th>New</th>
    </tr>
    <tr valign="top">
      <td>core.Settings.<strong>User</strong>AuthToken{}</td>
      <td>core.Settings.<strong>Record</strong>AuthToken{}</td>
    </tr>
    <tr valign="top">
      <td>core.Settings.<strong>User</strong>PasswordResetToken{}</td>
      <td>core.Settings.<strong>Record</strong>PasswordResetToken{}</td>
    </tr>
    <tr valign="top">
      <td>core.Settings.<strong>User</strong>EmailChangeToken{}</td>
      <td>core.Settings.<strong>Record</strong>EmailChangeToken{}</td>
    </tr>
    <tr valign="top">
      <td>core.Settings.<strong>User</strong>VerificationToken{}</td>
      <td>core.Settings.<strong>Record</strong>VerificationToken{}</td>
    </tr>
  </table>

- Marked as "Deprecated" and will be removed in v0.9+:
    ```
    core.Settings.EmailAuth{}
    core.EmailAuthConfig{}
    schema.FieldTypeUser
    schema.UserOptions{}
    ```

- The second argument of `apis.StaticDirectoryHandler(fileSystem, enableIndexFallback)` now is used to enable/disable index.html forwarding on missing file (eg. in case of SPA).
