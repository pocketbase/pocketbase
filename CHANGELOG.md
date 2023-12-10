## v0.20.0

- Added `expand`, `filter`, `fields`, custom query and headers parameters support for the realtime subscriptions.
    _Requires JS SDK v0.20.0+ or Dart SDK v0.17.0+._

    ```js
    // JS SDK v0.20.0
    pb.collection("example").subscribe("*", (e) => {
      ...
    }, {
      expand: "someRelField",
      filter: "status = 'active'",
      fields: "id,expand.someRelField.*:excerpt(100)",
    })
    ```

    ```dart
    // Dart SDK v0.17.0
    pb.collection("example").subscribe("*", (e) {
        ...
      },
      expand: "someRelField",
      filter: "status = 'active'",
      fields: "id,expand.someRelField.*:excerpt(100)",
    )
    ```

- Generalized the logs to allow any kind of application logs, not just requests.

    The new `app.Logger()` implements the standard [`log/slog` interfaces](https://pkg.go.dev/log/slog) available with Go 1.21.
    ```
    // Go: https://pocketbase.io/docs/go-logging/
    app.Logger().Info("Example message", "total", 123, "details", "lorem ipsum...")

    // JS: https://pocketbase.io/docs/js-logging/
    $app.logger().info("Example message", "total", 123, "details", "lorem ipsum...")
    ```

    For better performance and to minimize blocking on hot paths, logs are currently written with
    debounce and on batches:

      - 3 seconds after the last debounced log write
      - when the batch threshold is reached (currently 200)
      - right before app termination to attempt saving everything from the existing logs queue

    Some notable log related changes:

      - ⚠️ Bumped the minimum required Go version to 1.21.

      - ⚠️ Removed `_requests` table in favor of the generalized `_logs`.
        _Note that existing logs will be deleted!_

      - ⚠️ Renamed the following `Dao` log methods:
        ```go
        Dao.RequestQuery(...)      -> Dao.LogQuery(...)
        Dao.FindRequestById(...)   -> Dao.FindLogById(...)
        Dao.RequestsStats(...)     -> Dao.LogsStats(...)
        Dao.DeleteOldRequests(...) -> Dao.DeleteOldLogs(...)
        Dao.SaveRequest(...)       -> Dao.SaveLog(...)
        ```
      - ⚠️ Removed `app.IsDebug()` and the `--debug` flag.
        This was done to avoid the confusion with the new logger and its debug severity level.
        If you want to store debug logs you can set `-4` as min log level from the Admin UI.

      - Refactored Admin UI Logs:
        - Added new logs table listing.
        - Added log settings option to toggle the IP logging for the activity logger.
        - Added log settings option to specify a minimum log level.
        - Added controls to export individual or bulk selected logs as json.
        - Other minor improvements and fixes.

- Added new `filesystem/System.Copy(src, dest)` method to copy existing files from one location to another.
  _This is usually useful when duplicating records with `file` field(s) programmatically._

- Added `filesystem.NewFileFromUrl(ctx, url)` helper method to construct a `*filesystem.BytesReader` file from the specified url.

- OAuth2 related additions:

    - Added new `PKCE()` and `SetPKCE(enable)` OAuth2 methods to indicate whether the PKCE flow is supported or not.
      _The PKCE value is currently configurable from the UI only for the OIDC providers._
      _This was added to accommodate OIDC providers that may throw an error if unsupported PKCE params are submitted with the auth request (eg. LinkedIn; see [#3799](https://github.com/pocketbase/pocketbase/discussions/3799#discussioncomment-7640312))._

    - Added new `displayName` field for each `listAuthMethods()` OAuth2 provider item.
      _The value of the `displayName` property is currently configurable from the UI only for the OIDC providers._

    - Added `expiry` field to the OAuth2 user response containing the _optional_ expiration time of the OAuth2 access token ([#3617](https://github.com/pocketbase/pocketbase/discussions/3617)).

    - Allow a single OAuth2 user to be used for authentication in multiple auth collection.
      _⚠️ Because now you can have more than one external provider with `collectionId-provider-providerId` pair, `Dao.FindExternalAuthByProvider(provider, providerId)` method was removed in favour of the more generic `Dao.FindFirstExternalAuthByExpr(expr)`._

- Added `onlyVerified` auth collection option to globally disallow authentication requests for unverified users.

- Added support for single line comments (ex. `// your comment`) in the API rules and filter expressions.

- Added support for specifying a collection alias in `@collection.someCollection:alias.*`.

- Soft-deprecated and renamed `app.Cache()` with `app.Store()`.

- Minor JSVM updates and fixes:

    - Updated `$security.parseUnverifiedJWT(token)` and `$security.parseJWT(token, key)` to return the token payload result as plain object.

    - Added `$apis.requireGuestOnly()` middleware JSVM binding ([#3896](https://github.com/pocketbase/pocketbase/issues/3896)).

- Use `IS NOT` instead of `!=` as not-equal SQL query operator to handle the cases when comparing with nullable columns or expressions (eg. `json_extract` over `json` field).
  _Based on my local dataset I wasn't able to find a significant difference in the performance between the 2 operators, but if you stumble on a query that you think may be affected negatively by this, please report it and I'll test it further._

- Added `MaxSize` `json` field option to prevent storing large json data in the db ([#3790](https://github.com/pocketbase/pocketbase/issues/3790)).
  _Existing `json` fields are updated with a system migration to have a ~2MB size limit (it can be adjusted from the Admin UI)._

- Fixed negative string number normalization support for the `json` field type.

- Trigger the `app.OnTerminate()` hook on `app.Restart()` call.
  _A new bool `IsRestart` field was also added to the `core.TerminateEvent` event._

- Fixed graceful shutdown handling and speed up a little the app termination time.

- Limit the concurrent thumbs generation to avoid high CPU and memory usage in spiky scenarios ([#3794](https://github.com/pocketbase/pocketbase/pull/3794); thanks @t-muehlberger).
  _Currently the max concurrent thumbs generation processes are limited to "total of logical process CPUs + 1"._
  _This is arbitrary chosen and may change in the future depending on the users feedback and usage patterns._
  _If you are experiencing OOM errors during large image thumb generations, especially in container environment, you can try defining the `GOMEMLIMIT=500MiB` env variable before starting the executable._

- Slightly speed up (~10%) the thumbs generation by changing from cubic (`CatmullRom`) to bilinear (`Linear`) resampling filter (_the quality difference is very little_).

- Added a default red colored Stderr output in case of a console command error.
  _You can now also silence individually custom commands errors using the `cobra.Command.SilenceErrors` field._

- Fixed links formatting in the autogenerated html->text mail body.

- Removed incorrectly imported empty `local('')` font-face declarations.


## v0.19.4

- Fixed TinyMCE source code viewer textarea styles ([#3715](https://github.com/pocketbase/pocketbase/issues/3715)).

- Fixed `text` field min/max validators to properly count multi-byte characters ([#3735](https://github.com/pocketbase/pocketbase/issues/3735)).

- Allowed hyphens in `username` ([#3697](https://github.com/pocketbase/pocketbase/issues/3697)).
  _More control over the system fields settings will be available in the future._

- Updated the JSVM generated types to use directly the value type instead of `* | undefined` union in functions/methods return declarations.


## v0.19.3

- Added the release notes to the console output of `./pocketbase update` ([#3685](https://github.com/pocketbase/pocketbase/discussions/3685)).

- Added missing documention for the JSVM `$mails.*` bindings.

- Relaxed the OAuth2 redirect url validation to allow any string value ([#3689](https://github.com/pocketbase/pocketbase/pull/3689); thanks @sergeypdev).
  _Note that the redirect url format is still bound to the accepted values by the specific OAuth2 provider._


## v0.19.2

- Updated the JSVM generated types ([#3627](https://github.com/pocketbase/pocketbase/issues/3627), [#3662](https://github.com/pocketbase/pocketbase/issues/3662)).


## v0.19.1

- Fixed `tokenizer.Scan()/ScanAll()` to ignore the separators from the default trim cutset.
  An option to return also the empty found tokens was also added via `Tokenizer.KeepEmptyTokens(true)`.
  _This should fix the parsing of whitespace characters around view query column names when no quotes are used ([#3616](https://github.com/pocketbase/pocketbase/discussions/3616#discussioncomment-7398564))._

- Fixed the `:excerpt(max, withEllipsis?)` `fields` query param modifier to properly add space to the generated text fragment after block tags.


## v0.19.0

- Added Patreon OAuth2 provider ([#3323](https://github.com/pocketbase/pocketbase/pull/3323); thanks @ghostdevv).

- Added mailcow OAuth2 provider ([#3364](https://github.com/pocketbase/pocketbase/pull/3364); thanks @thisni1s).

- Added support for `:excerpt(max, withEllipsis?)` `fields` modifier that will return a short plain text version of any string value (html tags are stripped).
    This could be used to minimize the downloaded json data when listing records with large `editor` html values.
    ```js
    await pb.collection("example").getList(1, 20, {
      "fields": "*,description:excerpt(100)"
    })
    ```

- Several Admin UI improvements:
  - Count the total records separately to speed up the query execution for large datasets ([#3344](https://github.com/pocketbase/pocketbase/issues/3344)).
  - Enclosed the listing scrolling area within the table so that the horizontal scrollbar and table header are always reachable ([#2505](https://github.com/pocketbase/pocketbase/issues/2505)).
  - Allowed opening the record preview/update form via direct URL ([#2682](https://github.com/pocketbase/pocketbase/discussions/2682)).
  - Reintroduced the local `date` field tooltip on hover.
  - Speed up the listing loading times for records with large `editor` field values by initially fetching only a partial of the records data (the complete record data is loaded on record preview/update).
  - Added "Media library" (collection images picker) support for the TinyMCE `editor` field.
  - Added support to "pin" collections in the sidebar.
  - Added support to manually resize the collections sidebar.
  - More clear "Nonempty" field label style.
  - Removed the legacy `.woff` and `.ttf` fonts and keep only `.woff2`.

- Removed the explicit `Content-Type` charset from the realtime response due to compatibility issues with IIS ([#3461](https://github.com/pocketbase/pocketbase/issues/3461)).
  _The `Connection:keep-alive` realtime response header was also removed as it is not really used with HTTP2 anyway._

- Added new JSVM bindings:
  - `new Cookie({ ... })` constructor for creating `*http.Cookie` equivalent value.
  - `new SubscriptionMessage({ ... })` constructor for creating a custom realtime subscription payload.
  - Soft-deprecated `$os.exec()` in favour of `$os.cmd()` to make it more clear that the call only prepares the command and doesn't execute it.

- ⚠️ Bumped the min required Go version to 1.19.


## v0.18.10

- Added global `raw` template function to allow outputting raw/verbatim HTML content in the JSVM templates ([#3476](https://github.com/pocketbase/pocketbase/discussions/3476)).
  ```
  {{.description|raw}}
  ```

- Trimmed view query semicolon and allowed single quotes for column aliases ([#3450](https://github.com/pocketbase/pocketbase/issues/3450#issuecomment-1748044641)).
  _Single quotes are usually [not a valid identifier quote characters](https://www.sqlite.org/lang_keywords.html), but for resilience and compatibility reasons SQLite allows them in some contexts where only an identifier is expected._

- Bumped the GitHub action to use [min Go 1.21.2](https://github.com/golang/go/issues?q=milestone%3AGo1.21.2) (_the fixed issues are not critical as they are mostly related to the compiler/build tools_).


## v0.18.9

- Fixed empty thumbs directories not getting deleted on Windows after deleting a record img file ([#3382](https://github.com/pocketbase/pocketbase/issues/3382)).

- Updated the generated JSVM typings to silent the TS warnings when trying to access a field/method in a Go->TS interface.


## v0.18.8

- Minor fix for the View collections API Preview and Admin UI listings incorrectly showing the `created` and `updated` fields as `N/A` when the view query doesn't have them.


## v0.18.7

- Fixed JS error in the Admin UI when listing records with invalid `relation` field value ([#3372](https://github.com/pocketbase/pocketbase/issues/3372)).
  _This could happen usually only during custom SQL import scripts or when directly modifying the record field value without data validations._

- Updated Go deps and the generated JSVM types.


## v0.18.6

- Return the response headers and cookies in the `$http.send()` result ([#3310](https://github.com/pocketbase/pocketbase/discussions/3310)).

- Added more descriptive internal error message for missing user/admin email on password reset requests.

- Updated Go deps.


## v0.18.5

- Fixed minor Admin UI JS error in the auth collection options panel introduced with the change from v0.18.4.


## v0.18.4

- Added escape character (`\`) support in the Admin UI to allow using `select` field values with comma ([#2197](https://github.com/pocketbase/pocketbase/discussions/2197)).


## v0.18.3

- Exposed a global JSVM `readerToString(reader)` helper function to allow reading Go `io.Reader` values ([#3273](https://github.com/pocketbase/pocketbase/discussions/3273)).

- Bumped the GitHub action to use [min Go 1.21.1](https://github.com/golang/go/issues?q=milestone%3AGo1.21.1+label%3ACherryPickApproved) for the prebuilt executable since it contains some minor `html/template` and `net/http` security fixes.


## v0.18.2

- Prevent breaking the record form in the Admin UI in case the browser's localStorage quota has been exceeded when uploading or storing large `editor` values ([#3265](https://github.com/pocketbase/pocketbase/issues/3265)).

- Updated docs and missing JSVM typings.

- Exposed additional crypto primitives under the `$security.*` JSVM namespace ([#3273](https://github.com/pocketbase/pocketbase/discussions/3273)):
  ```js
  // HMAC with SHA256
  $security.hs256("hello", "secret")

  // HMAC with SHA512
  $security.hs512("hello", "secret")

  // compare 2 strings with a constant time
  $security.equal(hash1, hash2)
  ```


## v0.18.1

- Excluded the local temp dir from the backups ([#3261](https://github.com/pocketbase/pocketbase/issues/3261)).


## v0.18.0

- Simplified the `serve` command to accept domain name(s) as argument to reduce any additional manual hosts setup that sometimes previously was needed when deploying on production ([#3190](https://github.com/pocketbase/pocketbase/discussions/3190)).
  ```sh
  ./pocketbase serve yourdomain.com
  ```

- Added `fields` wildcard (`*`) support.

- Added option to upload a backup file from the Admin UI ([#2599](https://github.com/pocketbase/pocketbase/issues/2599)).

- Registered a custom Deflate compressor to speedup (_nearly 2-3x_) the backups generation for the sake of a small zip size increase.
  _Based on several local tests, `pb_data` of ~500MB (from which ~350MB+ are several hundred small files) results in a ~280MB zip generated for ~11s (previously it resulted in ~250MB zip but for ~35s)._

- Added the application name as part of the autogenerated backup name for easier identification ([#3066](https://github.com/pocketbase/pocketbase/issues/3066)).

- Added new `SmtpConfig.LocalName` option to specify a custom domain name (or IP address) for the initial EHLO/HELO exchange ([#3097](https://github.com/pocketbase/pocketbase/discussions/3097)).
  _This is usually required for verification purposes only by some SMTP providers, such as on-premise [Gmail SMTP-relay](https://support.google.com/a/answer/2956491)._

- Added `NoDecimal` `number` field option.

- `editor` field improvements:
    - Added new "Strip urls domain" option to allow controlling the default TinyMCE urls behavior (_default to `false` for new content_).
    - Normalized pasted text while still preserving links, lists, tables, etc. formatting ([#3257](https://github.com/pocketbase/pocketbase/issues/3257)).

- Added option to auto generate admin and auth record passwords from the Admin UI.

- Added JSON validation and syntax highlight for the `json` field in the Admin UI ([#3191](https://github.com/pocketbase/pocketbase/issues/3191)).

- Added datetime filter macros:
  ```
  // all macros are UTC based
  @second     - @now second number (0-59)
  @minute     - @now minute number (0-59)
  @hour       - @now hour number (0-23)
  @weekday    - @now weekday number (0-6)
  @day        - @now day number
  @month      - @now month number
  @year       - @now year number
  @todayStart - beginning of the current day as datetime string
  @todayEnd   - end of the current day as datetime string
  @monthStart - beginning of the current month as datetime string
  @monthEnd   - end of the current month as datetime string
  @yearStart  - beginning of the current year as datetime string
  @yearEnd    - end of the current year as datetime string
  ```

- Added cron expression macros ([#3132](https://github.com/pocketbase/pocketbase/issues/3132)):
  ```
  @yearly   - "0 0 1 1 *"
  @annually - "0 0 1 1 *"
  @monthly  - "0 0 1 * *"
  @weekly   - "0 0 * * 0"
  @daily    - "0 0 * * *"
  @midnight - "0 0 * * *"
  @hourly   - "0 * * * *"
  ```

- ⚠️ Added offset argument `Dao.FindRecordsByFilter(collection, filter, sort, limit, offset, [params...])`.
  _If you don't need an offset, you can set it to `0`._

- To minimize the footguns with `Dao.FindFirstRecordByFilter()` and `Dao.FindRecordsByFilter()`, the functions now supports an optional placeholder params argument that is safe to be populated with untrusted user input.
  The placeholders are in the same format as when binding regular SQL parameters.
  ```go
  // unsanitized and untrusted filter variables
  status := "..."
  author := "..."

  app.Dao().FindFirstRecordByFilter("articles", "status={:status} && author={:author}", dbx.Params{
    "status": status,
    "author": author,
  })

  app.Dao().FindRecordsByFilter("articles", "status={:status} && author={:author}", "-created", 10, 0, dbx.Params{
    "status": status,
    "author": author,
  })
  ```

- Added JSVM `$mails.*` binds for the corresponding Go [mails package](https://pkg.go.dev/github.com/pocketbase/pocketbase/mails) functions.

- Added JSVM helper crypto primitives under the `$security.*` namespace:
  ```js
  $security.md5(text)
  $security.sha256(text)
  $security.sha512(text)
  ```

- ⚠️ Deprecated `RelationOptions.DisplayFields` in favor of the new `SchemaField.Presentable` option to avoid the duplication when a single collection is referenced more than once and/or by multiple other collections.

- ⚠️ Fill the `LastVerificationSentAt` and `LastResetSentAt` fields only after a successfull email send ([#3121](https://github.com/pocketbase/pocketbase/issues/3121)).

- ⚠️ Skip API `fields` json transformations for non 20x responses ([#3176](https://github.com/pocketbase/pocketbase/issues/3176)).

- ⚠️ Changes to `tests.ApiScenario` struct:

    - The `ApiScenario.AfterTestFunc` now receive as 3rd argument `*http.Response` pointer instead of `*echo.Echo` as the latter is not really useful in this context.
      ```go
      // old
      AfterTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo)

      // new
      AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response)
      ```

    - The `ApiScenario.TestAppFactory` now accept the test instance as argument and no longer expect an error as return result ([#3025](https://github.com/pocketbase/pocketbase/discussions/3025#discussioncomment-6592272)).
      ```go
      // old
      TestAppFactory: func() (*tests.TestApp, error)

      // new
      TestAppFactory: func(t *testing.T) *tests.TestApp
      ```
      _Returning a `nil` app instance from the factory results in test failure. You can enforce a custom test failure by calling `t.Fatal(err)` inside the factory._

- Bumped the min required TLS version to 1.2 in order to improve the cert reputation score.

- Reduced the default JSVM prewarmed pool size to 25 to reduce the initial memory consumptions (_you can manually adjust the pool size with `--hooksPool=50` if you need to, but the default should suffice for most cases_).

- Update `gocloud.dev` dependency to v0.34 and explicitly set the new `NoTempDir` fileblob option to prevent the cross-device link error introduced with v0.33.

- Other minor Admin UI and docs improvements.


## v0.17.7

- Fixed the autogenerated `down` migrations to properly revert the old collection rules in case a change was made in `up` ([#3192](https://github.com/pocketbase/pocketbase/pull/3192); thanks @impact-merlinmarek).
  _Existing `down` migrations can't be fixed but that should be ok as usually the `down` migrations are rarely used against prod environments since they can cause data loss and, while not ideal, the previous old behavior of always setting the rules to `null/nil` is safer than not updating the rules at all._

- Updated some Go deps.


## v0.17.6

- Fixed JSVM `require()` file path error when using Windows-style path delimiters ([#3163](https://github.com/pocketbase/pocketbase/issues/3163#issuecomment-1685034438)).


## v0.17.5

- Added quotes around the wrapped view query columns introduced with v0.17.4.


## v0.17.4

- Fixed Views record retrieval when numeric id is used ([#3110](https://github.com/pocketbase/pocketbase/issues/3110)).
  _With this fix we also now properly recognize `CAST(... as TEXT)` and `CAST(... as BOOLEAN)` as `text` and `bool` fields._

- Fixed `relation` "Cascade delete" tooltip message ([#3098](https://github.com/pocketbase/pocketbase/issues/3098)).

- Fixed jsvm error message prefix on failed migrations ([#3103](https://github.com/pocketbase/pocketbase/pull/3103); thanks @nzhenev).

- Disabled the initial Admin UI admins counter cache when there are no initial admins to allow detecting externally created accounts (eg. with the `admin` command) ([#3106](https://github.com/pocketbase/pocketbase/issues/3106)).

- Downgraded `google/go-cloud` dependency to v0.32.0 until v0.34.0 is released to prevent the `os.TempDir` `cross-device link` errors as too many users complained about it.


## v0.17.3

- Fixed Docker `cross-device link` error when creating `pb_data` backups on a local mounted volume ([#3089](https://github.com/pocketbase/pocketbase/issues/3089)).

- Fixed the error messages for relation to views ([#3090](https://github.com/pocketbase/pocketbase/issues/3090)).

- Always reserve space for the scrollbar to reduce the layout shifts in the Admin UI records listing due to the deprecated `overflow: overlay`.

- Enabled lazy loading for the Admin UI thumb images.


## v0.17.2

- Soft-deprecated `$http.send({ data: object, ... })` in favour of `$http.send({ body: rawString, ... })`
  to allow sending non-JSON body with the request ([#3058](https://github.com/pocketbase/pocketbase/discussions/3058)).
  The existing `data` prop will still work, but it is recommended to use `body` instead (_to send JSON you can use `JSON.stringify(...)` as body value_).

- Added `core.RealtimeConnectEvent.IdleTimeout` field to allow specifying a different realtime idle timeout duration per client basis ([#3054](https://github.com/pocketbase/pocketbase/discussions/3054)).

- Fixed `apis.RequestData` deprecation log note ([#3068](https://github.com/pocketbase/pocketbase/pull/3068); thanks @gungjodi).


## v0.17.1

- Use relative path when redirecting to the OAuth2 providers page in the Admin UI to support subpath deployments ([#3026](https://github.com/pocketbase/pocketbase/pull/3026); thanks @sonyarianto).

- Manually trigger the `OnBeforeServe` hook for `tests.ApiScenario` ([#3025](https://github.com/pocketbase/pocketbase/discussions/3025)).

- Trigger the JSVM `cronAdd()` handler only on app `serve` to prevent unexpected (and eventually duplicated) cron handler calls when custom console commands are used ([#3024](https://github.com/pocketbase/pocketbase/discussions/3024#discussioncomment-6592703)).

- The `console.log()` messages are now written to the `stdout` instead of `stderr`.


## v0.17.0

- New more detailed guides for using PocketBase as framework (both Go and JS).
  _If you find any typos or issues with the docs please report them in https://github.com/pocketbase/site._

- Added new experimental JavaScript app hooks binding via [goja](https://github.com/dop251/goja).
  They are available by default with the prebuilt executable if you create `*.pb.js` file(s) in the `pb_hooks` directory.
  Lower your expectations because the integration comes with some limitations. For more details please check the [Extend with JavaScript](https://pocketbase.io/docs/js-overview/) guide.
  Optionally, you can also enable the JS app hooks as part of a custom Go build for dynamic scripting but you need to register the `jsvm` plugin manually:
  ```go
  jsvm.MustRegister(app core.App, config jsvm.Config{})
  ```

- Added Instagram OAuth2 provider ([#2534](https://github.com/pocketbase/pocketbase/pull/2534); thanks @pnmcosta).

- Added VK OAuth2 provider ([#2533](https://github.com/pocketbase/pocketbase/pull/2533); thanks @imperatrona).

- Added Yandex OAuth2 provider ([#2762](https://github.com/pocketbase/pocketbase/pull/2762); thanks @imperatrona).

- Added new fields to `core.ServeEvent`:
  ```go
  type ServeEvent struct {
    App    App
    Router *echo.Echo
    // new fields
    Server      *http.Server      // allows adjusting the HTTP server config (global timeouts, TLS options, etc.)
    CertManager *autocert.Manager // allows adjusting the autocert options (cache dir, host policy, etc.)
  }
  ```

- Added `record.ExpandedOne(rel)` and `record.ExpandedAll(rel)` helpers to retrieve casted single or multiple expand relations from the already loaded "expand" Record data.

- Added rule and filter record `Dao` helpers:
  ```go
  app.Dao().FindRecordsByFilter("posts", "title ~ 'lorem ipsum' && visible = true", "-created", 10)
  app.Dao().FindFirstRecordByFilter("posts", "slug='test' && active=true")
  app.Dao().CanAccessRecord(record, requestInfo, rule)
  ```

- Added `Dao.WithoutHooks()` helper to create a new `Dao` from the current one but without the create/update/delete hooks.

- Use a default fetch function that will return all relations in case the `fetchFunc` argument of `Dao.ExpandRecord(record, expands, fetchFunc)` and `Dao.ExpandRecords(records, expands, fetchFunc)` is `nil`.

- For convenience it is now possible to call `Dao.RecordQuery(collectionModelOrIdentifier)` with just the collection id or name.
  In case an invalid collection id/name string is passed the query will be resolved with cancelled context error.

- Refactored `apis.ApiError` validation errors serialization to allow `map[string]error` and `map[string]any` when generating the public safe formatted `ApiError.Data`.

- Added support for wrapped API errors (_in case Go 1.20+ is used with multiple wrapped errors, the first `apis.ApiError` takes precedence_).

- Added `?download=1` file query parameter to the file serving endpoint to force the browser to always download the file and not show its preview.

- Added new utility `github.com/pocketbase/pocketbase/tools/template` subpackage to assist with rendering HTML templates using the standard Go `html/template` and `text/template` syntax.

- Added `types.JsonMap.Get(k)` and `types.JsonMap.Set(k, v)` helpers for the cases where the type aliased direct map access is not allowed (eg. in [goja](https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods)).

- Soft-deprecated `security.NewToken()` in favor of `security.NewJWT()`.

- `Hook.Add()` and `Hook.PreAdd` now returns a unique string identifier that could be used to remove the registered hook handler via `Hook.Remove(handlerId)`.

- Changed the after* hooks to be called right before writing the user response, allowing users to return response errors from the after hooks.
  There is also no longer need for returning explicitly `hook.StopPropagtion` when writing custom response body in a hook because we will skip the finalizer response body write if a response was already "committed".

- ⚠️ Renamed `*Options{}` to `Config{}` for consistency and replaced the unnecessary pointers with their value equivalent to keep the applied configuration defaults isolated within their function calls:
  ```go
  old: pocketbase.NewWithConfig(config *pocketbase.Config) *pocketbase.PocketBase
  new: pocketbase.NewWithConfig(config pocketbase.Config) *pocketbase.PocketBase

  old: core.NewBaseApp(config *core.BaseAppConfig) *core.BaseApp
  new: core.NewBaseApp(config core.BaseAppConfig) *core.BaseApp

  old: apis.Serve(app core.App, options *apis.ServeOptions) error
  new: apis.Serve(app core.App, config apis.ServeConfig) (*http.Server, error)

  old: jsvm.MustRegisterMigrations(app core.App, options *jsvm.MigrationsOptions)
  new: jsvm.MustRegister(app core.App, config jsvm.Config)

  old: ghupdate.MustRegister(app core.App, rootCmd *cobra.Command, options *ghupdate.Options)
  new: ghupdate.MustRegister(app core.App, rootCmd *cobra.Command, config ghupdate.Config)

  old: migratecmd.MustRegister(app core.App, rootCmd *cobra.Command, options *migratecmd.Options)
  new: migratecmd.MustRegister(app core.App, rootCmd *cobra.Command, config migratecmd.Config)
  ```

- ⚠️ Changed the type of `subscriptions.Message.Data` from `string` to `[]byte` because `Data` usually is a json bytes slice anyway.

- ⚠️ Renamed `models.RequestData` to `models.RequestInfo` and soft-deprecated `apis.RequestData(c)` in favor of `apis.RequestInfo(c)` to avoid the stuttering with the `Data` field.
  _The old `apis.RequestData()` method still works to minimize the breaking changes but it is recommended to replace it with `apis.RequestInfo(c)`._

- ⚠️ Changes to the List/Search APIs
    - Added new query parameter `?skipTotal=1` to skip the `COUNT` query performed with the list/search actions ([#2965](https://github.com/pocketbase/pocketbase/discussions/2965)).
      If `?skipTotal=1` is set, the response fields `totalItems` and `totalPages` will have `-1` value (this is to avoid having different JSON responses and to differentiate from the zero default).
      With the latest JS SDK 0.16+ and Dart SDK v0.11+ versions `skipTotal=1` is set by default for the `getFirstListItem()` and `getFullList()` requests.

    - The count and regular select statements also now executes concurrently, meaning that we no longer perform normalization over the `page` parameter and in case the user
      request a page that doesn't exist (eg. `?page=99999999`) we'll return empty `items` array.

    - Reverted the default `COUNT` column to `id` as there are some common situations where it can negatively impact the query performance.
      Additionally, from this version we also set `PRAGMA temp_store = MEMORY` so that also helps with the temp B-TREE creation when `id` is used.
      _There are still scenarios where `COUNT` queries with `rowid` executes faster, but the majority of the time when nested relations lookups are used it seems to have the opposite effect (at least based on the benchmarks dataset)._

- ⚠️ Disallowed relations to views **from non-view** collections ([#3000](https://github.com/pocketbase/pocketbase/issues/3000)).
  The change was necessary because I wasn't able to find an efficient way to track view changes and the previous behavior could have too many unexpected side-effects (eg. view with computed ids).
  There is a system migration that will convert the existing view `relation` fields to `json` (multiple) and `text` (single) fields.
  This could be a breaking change if you have `relation` to view and use `expand` or some of the `relation` view fields as part of a collection rule.

- ⚠️ Added an extra `action` argument to the `Dao` hooks to allow skipping the default persist behavior.
  In preparation for the logs generalization, the `Dao.After*Func` methods now also allow returning an error.

- Allowed `0` as `RelationOptions.MinSelect` value to avoid the ambiguity between 0 and non-filled input value ([#2817](https://github.com/pocketbase/pocketbase/discussions/2817)).

- Fixed zero-default value not being used if the field is not explicitly set when manually creating records ([#2992](https://github.com/pocketbase/pocketbase/issues/2992)).
  Additionally, `record.Get(field)` will now always return normalized value (the same as in the json serialization) for consistency and to avoid ambiguities with what is stored in the related DB table.
  The schema fields columns `DEFAULT` definition was also updated for new collections to ensure that `NULL` values can't be accidentally inserted.

- Fixed `migrate down` not returning the correct `lastAppliedMigrations()` when the stored migration applied time is in seconds.

- Fixed realtime delete event to be called after the record was deleted from the DB (_including transactions and cascade delete operations_).

- Other minor fixes and improvements (typos and grammar fixes, updated dependencies, removed unnecessary 404 error check in the Admin UI, etc.).


## v0.16.10

- Added multiple valued fields (`relation`, `select`, `file`) normalizations to ensure that the zero-default value of a newly created multiple field is applied for already existing data ([#2930](https://github.com/pocketbase/pocketbase/issues/2930)).


## v0.16.9

- Register the `eagerRequestInfoCache` middleware only for the internal `api` group routes to avoid conflicts with custom route handlers ([#2914](https://github.com/pocketbase/pocketbase/issues/2914)).


## v0.16.8

- Fixed unique validator detailed error message not being returned when camelCase field name is used ([#2868](https://github.com/pocketbase/pocketbase/issues/2868)).

- Updated the index parser to allow no space between the table name and the columns list ([#2864](https://github.com/pocketbase/pocketbase/discussions/2864#discussioncomment-6373736)).

- Updated go deps.


## v0.16.7

- Minor optimization for the list/search queries to use `rowid` with the `COUNT` statement when available.
  _This eliminates the temp B-TREE step when executing the query and for large datasets (eg. 150k) it could have 10x improvement (from ~580ms to ~60ms)._


## v0.16.6

- Fixed collection index column sort normalization in the Admin UI ([#2681](https://github.com/pocketbase/pocketbase/pull/2681); thanks @SimonLoir).

- Removed unnecessary admins count in `apis.RequireAdminAuthOnlyIfAny()` middleware ([#2726](https://github.com/pocketbase/pocketbase/pull/2726); thanks @svekko).

- Fixed `multipart/form-data` request bind not populating map array values ([#2763](https://github.com/pocketbase/pocketbase/discussions/2763#discussioncomment-6278902)).

- Upgraded npm and Go dependencies.


## v0.16.5

- Fixed the Admin UI serialization of implicit relation display fields ([#2675](https://github.com/pocketbase/pocketbase/issues/2675)).

- Reset the Admin UI sort in case the active sort collection field is renamed or deleted.


## v0.16.4

- Fixed the selfupdate command not working on Windows due to missing `.exe` in the extracted binary path ([#2589](https://github.com/pocketbase/pocketbase/discussions/2589)).
  _Note that the command on Windows will work from v0.16.4+ onwards, meaning that you still will have to update manually one more time to v0.16.4._

- Added `int64`, `int32`, `uint`, `uint64` and `uint32` support when scanning `types.DateTime` ([#2602](https://github.com/pocketbase/pocketbase/discussions/2602))

- Updated dependencies.


## v0.16.3

- Fixed schema fields sort not working on Safari/Gnome Web ([#2567](https://github.com/pocketbase/pocketbase/issues/2567)).

- Fixed default `PRAGMA`s not being applied for new connections ([#2570](https://github.com/pocketbase/pocketbase/discussions/2570)).


## v0.16.2

- Fixed backups archive not excluding the local `backups` directory on Windows ([#2548](https://github.com/pocketbase/pocketbase/discussions/2548#discussioncomment-5979712)).

- Changed file field to not use `dataTransfer.effectAllowed` when dropping files since it is not reliable and consistent across different OS and browsers ([#2541](https://github.com/pocketbase/pocketbase/issues/2541)).

- Auto register the initial generated snapshot migration to prevent incorrectly reapplying the snapshot on Docker restart ([#2551](https://github.com/pocketbase/pocketbase/discussions/2551)).

- Fixed missing view id field error message typo.


## v0.16.1

- Fixed backup restore not working in a container environment when `pb_data` is mounted as volume ([#2519](https://github.com/pocketbase/pocketbase/issues/2519)).

- Fixed Dart SDK realtime API preview example ([#2523](https://github.com/pocketbase/pocketbase/pull/2523); thanks @xFrann).

- Fixed typo in the backups create panel ([#2526](https://github.com/pocketbase/pocketbase/pull/2526); thanks @dschissler).

- Removed unnecessary slice length check in `list.ExistInSlice` ([#2527](https://github.com/pocketbase/pocketbase/pull/2527); thanks @KunalSin9h).

- Avoid mutating the cached request data on OAuth2 user create ([#2535](https://github.com/pocketbase/pocketbase/discussions/2535)).

- Fixed Export Collections "Download as JSON" ([#2540](https://github.com/pocketbase/pocketbase/issues/2540)).

- Fixed file field drag and drop not working in Firefox and Safari ([#2541](https://github.com/pocketbase/pocketbase/issues/2541)).


## v0.16.0

- Added automated backups (_+ cron rotation_) APIs and UI for the `pb_data` directory.
  The backups can be also initialized programmatically using `app.CreateBackup("backup.zip")`.
  There is also experimental restore method - `app.RestoreBackup("backup.zip")` (_currently works only on UNIX systems as it relies on execve_).
  The backups can be stored locally or in external S3 storage (_it has its own configuration, separate from the file uploads storage filesystem_).

- Added option to limit the returned API fields using the `?fields` query parameter.
  The "fields picker" is applied for `SearchResult.Items` and every other JSON response. For example:
  ```js
  // original: {"id": "RECORD_ID", "name": "abc", "description": "...something very big...", "items": ["id1", "id2"], "expand": {"items": [{"id": "id1", "name": "test1"}, {"id": "id2", "name": "test2"}]}}
  // output:   {"name": "abc", "expand": {"items": [{"name": "test1"}, {"name": "test2"}]}}
  const result = await pb.collection("example").getOne("RECORD_ID", {
    expand: "items",
    fields: "name,expand.items.name",
  })
  ```

- Added new `./pocketbase update` command to selfupdate the prebuilt executable (with option to generate a backup of your `pb_data`).

- Added new `./pocketbase admin` console command:
  ```sh
  // creates new admin account
  ./pocketbase admin create test@example.com 123456890

  // changes the password of an existing admin account
  ./pocketbase admin update test@example.com 0987654321

  // deletes single admin account (if exists)
  ./pocketbase admin delete test@example.com
  ```

- Added `apis.Serve(app, options)` helper to allow starting the API server programmatically.

- Updated the schema fields Admin UI for "tidier" fields visualization.

- Updated the logs "real" user IP to check for `Fly-Client-IP` header and changed the `X-Forward-For` header to use the first non-empty leftmost-ish IP as it the closest to the "real IP".

- Added new `tools/archive` helper subpackage for managing archives (_currently works only with zip_).

- Added new `tools/cron` helper subpackage for scheduling task using cron-like syntax (_this eventually may get exported in the future in a separate repo_).

- Added new `Filesystem.List(prefix)` helper to retrieve a flat list with all files under the provided prefix.

- Added new `App.NewBackupsFilesystem()` helper to create a dedicated filesystem abstraction for managing app data backups.

- Added new `App.OnTerminate()` hook (_executed right before app termination, eg. on `SIGTERM` signal_).

- Added `accept` file field attribute with the field MIME types ([#2466](https://github.com/pocketbase/pocketbase/pull/2466); thanks @Nikhil1920).

- Added support for multiple files sort in the Admin UI ([#2445](https://github.com/pocketbase/pocketbase/issues/2445)).

- Added support for multiple relations sort in the Admin UI.

- Added `meta.isNew` to the OAuth2 auth JSON response to indicate a newly OAuth2 created PocketBase user.


## v0.15.3

- Updated the Admin UI to use the latest JS SDK to resolve the `isNew` record field conflict ([#2385](https://github.com/pocketbase/pocketbase/discussions/2385)).

- Fixed `editor` field fullscreen `z-index` ([#2410](https://github.com/pocketbase/pocketbase/issues/2410)).

- Inserts the default app settings as part of the system init migration so that they are always available when accessed from within a user defined migration ([#2423](https://github.com/pocketbase/pocketbase/discussions/2423)).


## v0.15.2

- Fixed View query `SELECT DISTINCT` identifiers parsing ([#2349-5706019](https://github.com/pocketbase/pocketbase/discussions/2349#discussioncomment-5706019)).

- Fixed View collection schema incorrectly resolving multiple aliased fields originating from the same field source ([#2349-5707675](https://github.com/pocketbase/pocketbase/discussions/2349#discussioncomment-5707675)).

- Added OAuth2 redirect fallback message to notify the user to go back to the app in case the browser window is not auto closed.


## v0.15.1

- Trigger the related `Record` model realtime subscription events on [custom model struct](https://pocketbase.io/docs/custom-models/) save ([#2325](https://github.com/pocketbase/pocketbase/discussions/2325)).

- Fixed `Ctrl + S` in the `editor` field not propagating the quick save shortcut to the parent form.

- Added `⌘ + S` alias for the record quick save shortcut (_I have no Mac device to test it but it should work based on [`e.metaKey` docs](https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/metaKey)_).

- Enabled RTL for the TinyMCE editor ([#2327](https://github.com/pocketbase/pocketbase/issues/2327)).

- Reduced the record form vertical layout shifts and slightly improved the rendering speed when loading multiple `relation` fields.

- Enabled Admin UI assets cache.


## v0.15.0

- Simplified the OAuth2 authentication flow in a single "all in one" call ([#55](https://github.com/pocketbase/pocketbase/issues/55)).
  Requires JS SDK v0.14.0+ or Dart SDK v0.9.0+.
  The manual code-token exchange flow is still supported but the SDK method is renamed to `authWithOAuth2Code()` (_to minimize the breaking changes the JS SDK has a function overload that will proxy the existing `authWithOauth2` calls to `authWithOAuth2Code`_).
  For more details and example, you could check https://pocketbase.io/docs/authentication/#oauth2-integration.

- Added support for protected files ([#215](https://github.com/pocketbase/pocketbase/issues/215)).
  Requires JS SDK v0.14.0+ or Dart SDK v0.9.0+.
  It works with a short lived (~5min) file token passed as query param with the file url.
  For more details and example, you could check https://pocketbase.io/docs/files-handling/#protected-files.

- ⚠️ Fixed typo in `Record.WithUnkownData()` -> `Record.WithUnknownData()`.

- Added simple loose wildcard search term support in the Admin UI.

- Added auto "draft" to allow restoring previous record state in case of accidental reload or power outage.

- Added `Ctrl + S` shortcut to save the record changes without closing the panel.

- Added "drop files" support for the file upload field.

- Refreshed the OAuth2 Admin UI.


## v0.14.5

- Added checks for `nil` hooks in `forms.RecordUpsert` when used with custom `Dao` ([#2277](https://github.com/pocketbase/pocketbase/issues/2277)).

- Fixed unique detailed field error not returned on record create failure ([#2287](https://github.com/pocketbase/pocketbase/discussions/2287)).


## v0.14.4

- Fixed concurrent map write pannic on `list.ExistInSliceWithRegex()` cache ([#2272](https://github.com/pocketbase/pocketbase/issues/2272)).


## v0.14.3

- Fixed Admin UI Logs `meta` visualization in Firefox ([#2221](https://github.com/pocketbase/pocketbase/issues/2221)).

- Downgraded to v1 of the `aws/aws-sdk-go` package since v2 has compatibility issues with GCS ([#2231](https://github.com/pocketbase/pocketbase/issues/2231)).

- Upgraded the GitHub action to use [min Go 1.20.3](https://github.com/golang/go/issues?q=milestone%3AGo1.20.3+label%3ACherryPickApproved) for the prebuilt executable since it contains some minor `net/http` security fixes.


## v0.14.2

- Reverted part of the old `COALESCE` handling as a fallback to support empty string comparison with missing joined relation fields.


## v0.14.1

- Fixed realtime events firing before the files upload completion.

- Updated the underlying S3 lib to use `aws-sdk-go-v2` ([#1346](https://github.com/pocketbase/pocketbase/pull/1346); thanks @yuxiang-gao).

- Updated TinyMCE to v6.4.1.

- Updated the godoc of `Dao.Save*` methods.


## v0.14.0

- Added _experimental_ Apple OAuth2 integration.

- Added `@request.headers.*` filter rule support.

- Added support for advanced unique constraints and indexes management ([#345](https://github.com/pocketbase/pocketbase/issues/345), [#544](https://github.com/pocketbase/pocketbase/issues/544))

- Simplified the collections fields UI to allow easier and quicker scaffolding of the data schema.

- Deprecated `SchemaField.Unique`. Unique constraints are now managed via indexes.
  The `Unique` field is a no-op and will be removed in future version.

- Removed the `COALESCE` wrapping from some of the generated filter conditions to make better use of the indexes ([#1939](https://github.com/pocketbase/pocketbase/issues/1939)).

- Detect `id` aliased view columns as single `relation` fields ([#2029](https://github.com/pocketbase/pocketbase/discussions/2029)).

- Optimized single relation lookups.

- Normalized record values on `maxSelect` field option change (`select`, `file`, `relation`).
  When changing **from single to multiple** all already inserted single values are converted to an array.
  When changing **from multiple to single** only the last item of the already inserted array items is kept.

- Changed the cost/round factor of bcrypt hash generation from 13 to 12 since several users complained about the slow authWithPassword responses on lower spec hardware.
  _The change will affect only new users. Depending on the demand, we might make it configurable from the auth options._

- Simplified the default mail template styles to allow more control over the template layout ([#1904](https://github.com/pocketbase/pocketbase/issues/1904)).

- Added option to explicitly set the record id from the Admin UI ([#2118](https://github.com/pocketbase/pocketbase/issues/2118)).

- Added `migrate history-sync` command to clean `_migrations` history table from deleted migration files references.

- Added new fields to the `core.RecordAuthWithOAuth2Event` struct:
    ```
    IsNewRecord     bool,          // boolean field indicating whether the OAuth2 action created a new auth record
    ProviderName    string,        // the name of the OAuth2 provider (eg. "google")
    ProviderClient  auth.Provider, // the loaded Provider client instance
    ```

- Added CGO linux target for the prebuilt executable.

- ⚠️ Renamed `daos.GetTableColumns()` to `daos.TableColumns()` for consistency with the other Dao table related helpers.

- ⚠️ Renamed `daos.GetTableInfo()` to `daos.TableInfo()` for consistency with the other Dao table related helpers.

- ⚠️ Changed `types.JsonArray` to support specifying a generic type, aka. `types.JsonArray[T]`.
  If you have previously used `types.JsonArray`, you'll have to update it to `types.JsonArray[any]`.

- ⚠️ Registered the `RemoveTrailingSlash` middleware only for the `/api/*` routes since it is causing issues with subpath file serving endpoints ([#2072](https://github.com/pocketbase/pocketbase/issues/2072)).

- ⚠️ Changed the request logs `method` value to UPPERCASE, eg. "get" => "GET" ([#1956](https://github.com/pocketbase/pocketbase/discussions/1956)).

- Other minor UI improvements.


## v0.13.4

- Removed eager unique collection name check to support lazy validation during bulk import.


## v0.13.3

- Fixed view collections import ([#2044](https://github.com/pocketbase/pocketbase/issues/2044)).

- Updated the records picker Admin UI to show properly view collection relations.


## v0.13.2

- Fixed Admin UI js error when selecting multiple `file` field as `relation` "Display fields" ([#1989](https://github.com/pocketbase/pocketbase/issues/1989)).


## v0.13.1

- Added `HEAD` request method support for the `/api/files/:collection/:recordId/:filename` route ([#1976](https://github.com/pocketbase/pocketbase/discussions/1976)).


## v0.13.0

- Added new "View" collection type allowing you to create a read-only collection from a custom SQL `SELECT` statement. It supports:
  - aggregations (`COUNT()`, `MIN()`, `MAX()`, `GROUP BY`, etc.)
  - column and table aliases
  - CTEs and subquery expressions
  - auto `relation` fields association
  - `file` fields proxying (up to 5 linked relations, eg. view1->view2->...->base)
  - `filter`, `sort` and `expand`
  - List and View API rules

- Added auto fail/retry (default to 8 attempts) for the `SELECT` queries to gracefully handle the `database is locked` errors ([#1795](https://github.com/pocketbase/pocketbase/discussions/1795#discussioncomment-4882169)).
  _The default max attempts can be accessed or changed via `Dao.MaxLockRetries`._

- Added default max query execution timeout (30s).
  _The default timeout can be accessed or changed via `Dao.ModelQueryTimeout`._
  _For the prebuilt executables it can be also changed via the `--queryTimeout=10` flag._

- Added support for `dao.RecordQuery(collection)` to scan directly the `One()` and `All()` results in `*models.Record` or `[]*models.Record` without the need of explicit `NullStringMap`.

- Added support to overwrite the default file serve headers if an explicit response header is set.

- Added file thumbs when visualizing `relation` display file fields.

- Added "Min select" `relation` field option.

- Enabled `process.env` in JS migrations to allow accessing `os.Environ()`.

- Added `UploadedFiles` field to the `RecordCreateEvent` and `RecordUpdateEvent` event structs.

- ⚠️ Moved file upload after the record persistent to allow setting custom record id safely from the `OnModelBeforeCreate` hook.

- ⚠️ Changed `System.GetFile()` to return directly `*blob.Reader` instead of the `io.ReadCloser` interface.

- ⚠️ Changed `To`, `Cc` and `Bcc` of `mailer.Message` to `[]mail.Address` for consistency and to allow multiple recipients and optional name.

    If you are sending custom emails, you'll have to replace:
    ```go
    message := &mailer.Message{
      ...

      // (old) To: mail.Address{Address: "to@example.com"}
      To: []mail.Address{{Address: "to@example.com", Name: "Some optional name"}},

      // (old) Cc: []string{"cc@example.com"}
      Cc: []mail.Address{{Address: "cc@example.com", Name: "Some optional name"}},

      // (old) Bcc: []string{"bcc@example.com"}
      Bcc: []mail.Address{{Address: "bcc@example.com", Name: "Some optional name"}},

      ...
    }
    ```

- ⚠️ Refactored the Authentik integration as a more generic "OpenID Connect" provider (`oidc`) to support any OIDC provider (Okta, Keycloak, etc.).
  _If you've previously used Authentik, make sure to rename the provider key in your code to `oidc`._
  _To enable more than one OIDC provider you can use the additional `oidc2` and `oidc3` provider keys._

- ⚠️ Removed the previously deprecated `Dao.Block()` and `Dao.Continue()` helpers in favor of `Dao.NonconcurrentDB()`.

- Updated the internal redirects to allow easier subpath deployment when behind a reverse proxy.

- Other minor Admin UI improvements.


## v0.12.3

- Fixed "Toggle column" reactivity when navigating between collections ([#1836](https://github.com/pocketbase/pocketbase/pull/1836)).

- Logged the current datetime on server start ([#1822](https://github.com/pocketbase/pocketbase/issues/1822)).


## v0.12.2

- Fixed the "Clear" button of the datepicker component not clearing the value ([#1730](https://github.com/pocketbase/pocketbase/discussions/1730)).

- Increased slightly the fields contrast ([#1742](https://github.com/pocketbase/pocketbase/issues/1742)).

- Auto close the multi-select dropdown if "Max select" is reached.


## v0.12.1

- Fixed js error on empty relation save.

- Fixed `overlay-active` css class not being removed on nested overlay panel close ([#1718](https://github.com/pocketbase/pocketbase/issues/1718)).

- Added the collection name in the page title ([#1711](https://github.com/pocketbase/pocketbase/issues/1711)).


## v0.12.0

- Refactored the relation picker UI to allow server-side search, sort, create, update and delete of relation records ([#976](https://github.com/pocketbase/pocketbase/issues/976)).

- Added new `RelationOptions.DisplayFields` option to specify custom relation field(s) visualization in the Admin UI.

- Added Authentik OAuth2 provider ([#1377](https://github.com/pocketbase/pocketbase/pull/1377); thanks @pr0ton11).

- Added LiveChat OAuth2 provider ([#1573](https://github.com/pocketbase/pocketbase/pull/1573); thanks @mariosant).

- Added Gitea OAuth2 provider ([#1643](https://github.com/pocketbase/pocketbase/pull/1643); thanks @hlanderdev).

- Added PDF file previews ([#1548](https://github.com/pocketbase/pocketbase/pull/1548); thanks @mjadobson).

- Added video and audio file previews.

- Added rich text editor (`editor`) field for HTML content based on TinyMCE ([#370](https://github.com/pocketbase/pocketbase/issues/370)).
  _Currently the new field doesn't have any configuration options or validations but this may change in the future depending on how devs ended up using it._

- Added "Duplicate" Collection and Record options in the Admin UI ([#1656](https://github.com/pocketbase/pocketbase/issues/1656)).

- Added `filesystem.GetFile()` helper to read files through the FileSystem abstraction ([#1578](https://github.com/pocketbase/pocketbase/pull/1578); thanks @avarabyeu).

- Added new auth event hooks for finer control and more advanced auth scenarios handling:

  ```go
  // auth record
  OnRecordBeforeAuthWithPasswordRequest()
  OnRecordAfterAuthWithPasswordRequest()
  OnRecordBeforeAuthWithOAuth2Request()
  OnRecordAfterAuthWithOAuth2Request()
  OnRecordBeforeAuthRefreshRequest()
  OnRecordAfterAuthRefreshRequest()

  // admin
  OnAdminBeforeAuthWithPasswordRequest()
  OnAdminAfterAuthWithPasswordRequest()
  OnAdminBeforeAuthRefreshRequest()
  OnAdminAfterAuthRefreshRequest()
  OnAdminBeforeRequestPasswordResetRequest()
  OnAdminAfterRequestPasswordResetRequest()
  OnAdminBeforeConfirmPasswordResetRequest()
  OnAdminAfterConfirmPasswordResetRequest()
  ```

- Added `models.Record.CleanCopy()` helper that creates a new record copy with only the latest data state of the existing one and all other options reset to their defaults.

- Added new helper `apis.RecordAuthResponse(app, httpContext, record, meta)` to return a standard Record auth API response ([#1623](https://github.com/pocketbase/pocketbase/issues/1623)).

- Refactored `models.Record` expand and data change operations to be concurrent safe.

- Refactored all `forms` Submit interceptors to use a generic data type as their payload.

- Added several `store.Store` helpers:
  ```go
  store.Reset(newData map[string]T)
  store.Length() int
  store.GetAll() map[string]T
  ```

- Added "tags" support for all Record and Model related event hooks.

    The "tags" allow registering event handlers that will be called only on matching table name(s) or colleciton id(s)/name(s).
    For example:
    ```go
    app.OnRecordBeforeCreateRequest("articles").Add(func(e *core.RecordCreateEvent) error {
      // called only on "articles" record creation
      log.Println(e.Record)
      return nil
    })
    ```
    For all those event hooks `*hook.Hook` was replaced with `*hooks.TaggedHook`, but the hook methods signatures are the same so it should behave as it was previously if no tags were specified.

- ⚠️ Fixed the `json` field **string** value normalization ([#1703](https://github.com/pocketbase/pocketbase/issues/1703)).

    In order to support seamlessly both `application/json` and `multipart/form-data`
    requests, the following normalization rules are applied if the `json` field is a
    **plain string value**:

    - "true" is converted to the json `true`
    - "false" is converted to the json `false`
    - "null" is converted to the json `null`
    - "[1,2,3]" is converted to the json `[1,2,3]`
    - "{\"a\":1,\"b\":2}" is converted to the json `{"a":1,"b":2}`
    - numeric strings are converted to json number
    - double quoted strings are left as they are (aka. without normalizations)
    - any other string (empty string too) is double quoted

    Additionally, the "Nonempty" `json` field constraint now checks for `null`, `[]`, `{}` and `""` (empty string).

- Added `aria-label` to some of the buttons in the Admin UI for better accessibility ([#1702](https://github.com/pocketbase/pocketbase/pull/1702); thanks @ndarilek).

- Updated the filename extension checks in the Admin UI to be case-insensitive ([#1707](https://github.com/pocketbase/pocketbase/pull/1707); thanks @hungcrush).

- Other minor improvements (more detailed API file upload errors, UI optimizations, docs improvements, etc.)


## v0.11.4

- Fixed cascade delete for rel records with the same id as the main record ([#1689](https://github.com/pocketbase/pocketbase/issues/1689)).


## v0.11.3

- Fix realtime API panic on concurrent clients iteration ([#1628](https://github.com/pocketbase/pocketbase/issues/1628))

  - `app.SubscriptionsBroker().Clients()` now returns a shallow copy of the underlying map.

  - Added `Discard()` and `IsDiscarded()` helper methods to the `subscriptions.Client` interface.

  - Slow clients should no longer "block" the main action completion.


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

- ⚠️ Added support for multi-match/match-all request data and collection multi-valued fields (`select`, `relation`) conditions.
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

- ⚠️ Renamed the metadata key `original_filename` to `original-filename` due to an S3 file upload error caused by the underscore character ([#1343](https://github.com/pocketbase/pocketbase/pull/1343); thanks @yuxiang-gao).

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

- ⚠️ Changed `pocketbase.NewWithConfig(config Config)` to `pocketbase.NewWithConfig(config *Config)` and added 4 new config settings:
  ```go
  DataMaxOpenConns int // default to core.DefaultDataMaxOpenConns
  DataMaxIdleConns int // default to core.DefaultDataMaxIdleConns
  LogsMaxOpenConns int // default to core.DefaultLogsMaxOpenConns
  LogsMaxIdleConns int // default to core.DefaultLogsMaxIdleConns
  ```

- Added new helper method `core.App.IsBootstrapped()` to check the current app bootstrap state.

- ⚠️ Changed `core.NewBaseApp(dir, encryptionEnv, isDebug)` to `NewBaseApp(config *BaseAppConfig)`.

- ⚠️ Removed `rest.UploadedFile` struct (see below `filesystem.File`).

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
  func (m *Record) WithUnknownData(state bool)
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
