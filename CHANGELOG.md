## v0.20.7

- Fixed the Admin UI auto indexes update when renaming fields with a common prefix ([#4160](https://github.com/pocketbase/pocketbase/issues/4160)).


## v0.20.6

- Fixed JSVM types generation for functions with omitted arg types ([#4145](https://github.com/pocketbase/pocketbase/issues/4145)).

- Updated Go deps.


## v0.20.5

- Minor CSS fix for the Admin UI to prevent the searchbar within a popup from expanding too much and pushing the controls out of the visible area ([#4079](https://github.com/pocketbase/pocketbase/issues/4079#issuecomment-1876994116)).


## v0.20.4

- Small fix for a regression introduced with the recent `json` field changes that was causing View collection column expressions recognized as `json` to fail to resolve ([#4072](https://github.com/pocketbase/pocketbase/issues/4072)).


## v0.20.3

- Fixed the `json` field query comparisons to work correctly with plain JSON values like `null`, `bool` `number`, etc. ([#4068](https://github.com/pocketbase/pocketbase/issues/4068)).
  Since there are plans in the future to allow custom SQLite builds and also in some situations it may be useful to be able to distinguish `NULL` from `''`,
  for the `json` fields (and for any other future non-standard field) we no longer apply `COALESCE` by default, aka.:
  ```
  Dataset:
  1) data: json(null)
  2) data: json('')

  For the filter "data = null" only 1) will resolve to TRUE.
  For the filter "data = ''"   only 2) will resolve to TRUE.
  ```

- Minor Go tests improvements
  - Sorted the record cascade delete references to ensure that the delete operation will preserve the order of the fired events when running the tests.
  - Marked some of the tests as safe for parallel execution to speed up a little the GitHub action build times.


## v0.20.2

- Added `sleep(milliseconds)` JSVM binding.
  _It works the same way as Go `time.Sleep()`, aka. it pauses the goroutine where the JSVM code is running._

- Fixed multi-line text paste in the Admin UI search bar ([#4022](https://github.com/pocketbase/pocketbase/discussions/4022)).

- Fixed the monospace font loading in the Admin UI.

- Fixed various reported docs and code comment typos.


## v0.20.1

- Added `--dev` flag and its accompanying `app.IsDev()` method (_in place of the previously removed `--debug`_) to assist during development ([#3918](https://github.com/pocketbase/pocketbase/discussions/3918)).
  The `--dev` flag prints in the console "everything" and more specifically:
  - the data DB SQL statements
  - all `app.Logger().*` logs (debug, info, warning, error, etc.), no matter of the logs persistence settings in the Admin UI

- Minor Admin UI fixes:
  - Fixed the log `error` label text wrapping.
  - Added the log `referer` (_when it is from a different source_) and `details` labels in the logs listing.
  - Removed the blank current time entry from the logs chart because it was causing confusion when used with custom time ranges.
  - Updated the SQL syntax highlighter and keywords autocompletion in the Admin UI to recognize `CAST(x as bool)` expressions.

- Replaced the default API tests timeout with a new `ApiScenario.Timeout` option ([#3930](https://github.com/pocketbase/pocketbase/issues/3930)).
  A negative or zero value means no tests timeout.
  If a single API test takes more than 3s to complete it will have a log message visible when the test fails or when `go test -v` flag is used.

- Added timestamp at the beginning of the generated JSVM types file to avoid creating it everytime with the app startup.


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

- Added missing documentation for the JSVM `$mails.*` bindings.

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
