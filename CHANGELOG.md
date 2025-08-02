## v0.29.1

- Updated the X/Twitter provider to return the `confirmed_email` field and to use the `x.com` domain ([#7035](https://github.com/pocketbase/pocketbase/issues/7035)).

- Added Box.com OAuth2 provider ([#7056](https://github.com/pocketbase/pocketbase/pull/7056); thanks @blakepatteson).

- Updated `modernc.org/sqlite` to 1.38.2 (SQLite 3.50.3).

- Fixed example List API response ([#7049](https://github.com/pocketbase/pocketbase/pull/7049); thanks @williamtguerra).


## v0.29.0

- Enabled calling the `/auth-refresh` endpoint with nonrenewable tokens.
    _When used with nonrenewable tokens (e.g. impersonate) the endpoint will simply return the same token with the up-to-date user data associated with it._

- Added the triggered rate rimit rule in the error log `details`.

- Added optional `ServeEvent.Listener` field to initialize a custom network listener (e.g. `unix`) instead of the default `tcp` ([#3233](https://github.com/pocketbase/pocketbase/discussions/3233)).

- Fixed request data unmarshalization for the `DynamicModel` array/object fields ([#7022](https://github.com/pocketbase/pocketbase/discussions/7022)).

- Fixed Dashboard page title `-` escaping ([#6982](https://github.com/pocketbase/pocketbase/issues/6982)).

- Other minor improvements (updated first superuser console text when running with `go run`, clarified trusted IP proxy header label, wrapped the backup restore in a transaction as an extra precaution, updated deps, etc.).


## v0.28.4

- Added global JSVM `toBytes()` helper to return the bytes slice representation of a value such as io.Reader or string, _other types are first serialized to Go string_ ([#6935](https://github.com/pocketbase/pocketbase/issues/6935)).

- Fixed `security.RandomStringByRegex` random distribution ([#6947](https://github.com/pocketbase/pocketbase/pull/6947); thanks @yerTools).

- Minor docs and typos fixes.


## v0.28.3

- Skip sending empty `Range` header when fetching blobs from S3 ([#6914](https://github.com/pocketbase/pocketbase/pull/6914)).

- Updated Go deps and particularly `modernc.org/sqlite` to 1.38.0 (SQLite 3.50.1).

- Bumped GitHub action min Go version to 1.23.10 as it comes with some [minor security `net/http` fixes](https://github.com/golang/go/issues?q=milestone%3AGo1.23.10+label%3ACherryPickApproved).


## v0.28.2

- Loaded latin-ext charset for the default text fonts ([#6869](https://github.com/pocketbase/pocketbase/issues/6869)).

- Updated view query CAST regex to properly recognize multiline expressions ([#6860](https://github.com/pocketbase/pocketbase/pull/6860); thanks @azat-ismagilov).

- Updated Go and npm dependencies.


## v0.28.1

- Fixed `json_each`/`json_array_length` normalizations to properly check for array values ([#6835](https://github.com/pocketbase/pocketbase/issues/6835)).


## v0.28.0

- Write the default response body of `*Request` hooks that are wrapped in a transaction after the related transaction completes to allow propagating the transaction error ([#6462](https://github.com/pocketbase/pocketbase/discussions/6462#discussioncomment-12207818)).

- Updated `app.DB()` to automatically routes raw write SQL statements to the nonconcurrent db pool ([#6689](https://github.com/pocketbase/pocketbase/discussions/6689)).
    _For the rare cases when it is needed users still have the option to explicitly target the specific pool they want using `app.ConcurrentDB()`/`app.NonconcurrentDB()`._

- ⚠️ Changed the default `json` field max size to 1MB.
    _Users still have the option to adjust the default limit from the collection field options but keep in mind that storing large strings/blobs in the database is known to cause performance issues and should be avoided when possible._

- ⚠️ Soft-deprecated and replaced `filesystem.System.GetFile(fileKey)` with `filesystem.System.GetReader(fileKey)` to avoid the confusion with `filesystem.File`.
    _The old method will still continue to work for at least until v0.29.0 but you'll get a console warning to replace it with `GetReader`._

- Added new `filesystem.System.GetReuploadableFile(fileKey, preserveName)` method to return an existing blob as a `*filesystem.File` value ([#6792](https://github.com/pocketbase/pocketbase/discussions/6792)).
    _This method could be useful in case you want to clone an existing Record file and assign it to a new Record (e.g. in a Record duplicate action)._

- Other minor improvements (updated the GitHub release min Go version to 1.23.9, updated npm and Go deps, etc.)


## v0.27.2

- Added workers pool when cascade deleting record files to minimize _"thread exhaustion"_ errors ([#6780](https://github.com/pocketbase/pocketbase/discussions/6780)).

- Updated the `:excerpt` fields modifier to properly account for multibyte characters ([#6778](https://github.com/pocketbase/pocketbase/issues/6778)).

- Use `rowid` as count column for non-view collections to minimize the need of having the id field in a covering index ([#6739](https://github.com/pocketbase/pocketbase/discussions/6739))


## v0.27.1

- Updated example `geoPoint` API preview body data.

- Added JSVM `new GeoPointField({ ... })` constructor.

- Added _partial_ WebP thumbs generation (_the thumbs will be stored as PNG_; [#6744](https://github.com/pocketbase/pocketbase/pull/6744)).

- Updated npm dev dependencies.


## v0.27.0

- ⚠️ Moved the Create and Manage API rule checks out of the `OnRecordCreateRequest` hook finalizer, **aka. now all CRUD API rules are checked BEFORE triggering their corresponding `*Request` hook**.
    This was done to minimize the confusion regarding the firing order of the request operations, making it more predictable and consistent with the other record List/View/Update/Delete request actions.
    It could be a minor breaking change if you are relying on the old behavior and have a Go `tests.ApiScenario` that is testing a Create API rule failure and expect `OnRecordCreateRequest` to be fired. In that case for example you may have to update your test scenario like:
    ```go
    tests.ApiScenario{
        Name:   "Example test that checks a Create API rule failure"
        Method: http.MethodPost,
        URL:    "/api/collections/example/records",
        ...
        // old:
        ExpectedEvents:  map[string]int{
            "*":                     0,
            "OnRecordCreateRequest": 1,
        },
        // new:
        ExpectedEvents:  map[string]int{"*": 0},
    }
    ```
    If you are having difficulties adjusting your code, feel free to open a [Q&A discussion](https://github.com/pocketbase/pocketbase/discussions) with the failing/problematic code sample.

- Added [new `geoPoint` field](https://pocketbase.io/docs/collections/#geopoint) for storing `{"lon":x,"lat":y}` geographic coordinates.
    In addition, a new [`geoDistance(lonA, lotA, lonB, lotB)` function](htts://pocketbase.io/docs/api-rules-and-filters/#geodistancelona-lata-lonb-latb) was also implemented that could be used to apply an API rule or filter constraint based on the distance (in km) between 2 geo points.

- Updated the `select` field UI to accommodate better larger lists and RTL languages ([#4674](https://github.com/pocketbase/pocketbase/issues/4674)).

- Updated the mail attachments auto MIME type detection to use `gabriel-vasile/mimetype` for consistency and broader sniffing signatures support.

- Forced `text/javascript` Content-Type when serving `.js`/`.mjs` collection uploaded files with the `/api/files/...` endpoint ([#6597](https://github.com/pocketbase/pocketbase/issues/6597)).

- Added second optional JSVM `DateTime` constructor argument for specifying a default timezone as TZ identifier when parsing the date string as alternative to a fixed offset in order to better handle daylight saving time nuances ([#6688](https://github.com/pocketbase/pocketbase/discussions/6688)):
    ```js
    // the same as with CET offset: new DateTime("2025-10-26 03:00:00 +01:00")
    new DateTime("2025-10-26 03:00:00", "Europe/Amsterdam") // 2025-10-26 02:00:00.000Z

    // the same as with CEST offset: new DateTime("2025-10-26 01:00:00 +02:00")
    new DateTime("2025-10-26 01:00:00", "Europe/Amsterdam") // 2025-10-25 23:00:00.000Z
    ```

- Soft-deprecated the `$http.send`'s `result.raw` field in favor of `result.body` that contains the response body as plain bytes slice to avoid the discrepancies between Go and the JSVM when casting binary data to string.

- Updated `modernc.org/sqlite` to 1.37.0.

- Other minor improvements (_removed the superuser fields from the auth record create/update body examples, allowed programmatically updating the auth record password from the create/update hooks, fixed collections import error response, etc._).


## v0.26.6

- Allow OIDC `email_verified` to be int or boolean string since some OIDC providers like AWS Cognito has non-standard userinfo response ([#6657](https://github.com/pocketbase/pocketbase/pull/6657)).

- Updated `modernc.org/sqlite` to 1.36.3.


## v0.26.5

- Fixed canonical URI parts escaping when generating the S3 request signature ([#6654](https://github.com/pocketbase/pocketbase/issues/6654)).


## v0.26.4

- Fixed `RecordErrorEvent.Error` and `CollectionErrorEvent.Error` sync with `ModelErrorEvent.Error` ([#6639](https://github.com/pocketbase/pocketbase/issues/6639)).

- Fixed logs details copy to clipboard action.

- Updated `modernc.org/sqlite` to 1.36.2.


## v0.26.3

- Fixed and normalized logs error serialization across common types for more consistent logs error output ([#6631](https://github.com/pocketbase/pocketbase/issues/6631)).


## v0.26.2

- Updated `golang-jwt/jwt` dependency because it comes with a [minor security fix](https://github.com/golang-jwt/jwt/security/advisories/GHSA-mh63-6h87-95cp).


## v0.26.1

- Removed the wrapping of `io.EOF` error when reading files since currently `io.ReadAll` doesn't check for wrapped errors ([#6600](https://github.com/pocketbase/pocketbase/issues/6600)).


## v0.26.0

- ⚠️ Replaced `aws-sdk-go-v2` and `gocloud.dev/blob` with custom lighter implementation ([#6562](https://github.com/pocketbase/pocketbase/discussions/6562)).
    As a side-effect of the dependency removal, the binary size has been reduced with ~10MB and builds ~30% faster.
    _Although the change is expected to be backward-compatible, I'd recommend to test first locally the new version with your S3 provider (if you use S3 for files storage and backups)._

- ⚠️ Prioritized the user submitted non-empty `createData.email` (_it will be unverified_) when creating the PocketBase user during the first OAuth2 auth.

- Load the request info context during password/OAuth2/OTP authentication ([#6402](https://github.com/pocketbase/pocketbase/issues/6402)).
    This could be useful in case you want to target the auth method as part of the MFA and Auth API rules.
    For example, to disable MFA for the OAuth2 auth could be expressed as `@request.context != "oauth2"` MFA rule.

- Added `store.Store.SetFunc(key, func(old T) new T)` to set/update a store value with the return result of the callback in a concurrent safe manner.

- Added `subscription.Message.WriteSSE(w, id)` for writing an SSE formatted message into the provided writer interface (_used mostly to assist with the unit testing_).

- Added `$os.stat(file)` JSVM helper ([#6407](https://github.com/pocketbase/pocketbase/discussions/6407)).

- Added log warning for `async` marked JSVM handlers and resolve when possible the returned `Promise` as fallback ([#6476](https://github.com/pocketbase/pocketbase/issues/6476)).

- Allowed calling `cronAdd`, `cronRemove` from inside other JSVM handlers ([#6481](https://github.com/pocketbase/pocketbase/discussions/6481)).

- Bumped the default request read and write timeouts to 5mins (_old 3mins_) to accommodate slower internet connections and larger file uploads/downloads.
    _If you want to change them you can modify the `OnServe` hook's `ServeEvent.ReadTimeout/WriteTimeout` fields as shown in [#6550](https://github.com/pocketbase/pocketbase/discussions/6550#discussioncomment-12364515)._

- Normalized the `@request.auth.*` and `@request.body.*` back relations resolver to always return `null` when the relation field is pointing to a different collection ([#6590](https://github.com/pocketbase/pocketbase/discussions/6590#discussioncomment-12496581)).

- Other minor improvements (_fixed query dev log nested parameters output, reintroduced `DynamicModel` object/array props reflect types caching, updated Go and npm deps, etc._)


## v0.25.9

- Fixed `DynamicModel` object/array props reflect type caching ([#6563](https://github.com/pocketbase/pocketbase/discussions/6563)).


## v0.25.8

- Added a default leeway of 5 minutes for the Apple/OIDC `id_token` timestamp claims check to account for clock-skew ([#6529](https://github.com/pocketbase/pocketbase/issues/6529)).
    It can be further customized if needed with the `PB_ID_TOKEN_LEEWAY` env variable (_the value must be in seconds, e.g. "PB_ID_TOKEN_LEEWAY=60" for 1 minute_).


## v0.25.7

- Fixed `@request.body.jsonObjOrArr.*` values extraction ([#6493](https://github.com/pocketbase/pocketbase/discussions/6493)).


## v0.25.6

- Restore the missing `meta.isNew` field of the OAuth2 success response ([#6490](https://github.com/pocketbase/pocketbase/issues/6490)).

- Updated npm dependencies.


## v0.25.5

- Set the current working directory as a default goja script path when executing inline JS strings to allow `require(m)` traversing parent `node_modules` directories.

- Updated `modernc.org/sqlite` and `modernc.org/libc` dependencies.


## v0.25.4

- Downgraded `aws-sdk-go-v2` to the version before the default data integrity checks because there have been reports for non-AWS S3 providers in addition to Backblaze (IDrive, R2) that no longer or partially work with the latest AWS SDK changes.

    While we try to enforce `when_required` by default, it is not enough to disable the new AWS SDK integrity checks entirely and some providers will require additional manual adjustments to make them compatible with the latest AWS SDK (e.g. removing the `x-aws-checksum-*` headers, unsetting the checksums calculation or reinstantiating the old MD5 checksums for some of the required operations, etc.) which as a result leads to a configuration mess that I'm not sure it would be a good idea to introduce.

    This unfornuatelly is not a PocketBase or Go specific issue and the official AWS SDKs for other languages are in the same situation (even the latest aws-cli).

    For those of you that extend PocketBase with Go: if your S3 vendor doesn't support the [AWS Data integrity checks](https://docs.aws.amazon.com/sdkref/latest/guide/feature-dataintegrity.html) and you are updating with `go get -u`, then make sure that the `aws-sdk-go-v2` dependencies in your `go.mod` are the same as in the repo:
    ```
    // go.mod
    github.com/aws/aws-sdk-go-v2 v1.36.1
    github.com/aws/aws-sdk-go-v2/config v1.28.10
    github.com/aws/aws-sdk-go-v2/credentials v1.17.51
    github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.48
    github.com/aws/aws-sdk-go-v2/service/s3 v1.72.2

    // after that run
    go clean -modcache && go mod tidy
    ```
    _The versions pinning is temporary until the non-AWS S3 vendors patch their implementation or until I manage to find time to remove/replace the `aws-sdk-go-v2` dependency (I'll consider prioritizing it for the v0.26 or v0.27 release)._


## v0.25.3

- Added a temporary exception for Backblaze S3 endpoints to exclude the new `aws-sdk-go-v2` checksum headers ([#6440](https://github.com/pocketbase/pocketbase/discussions/6440)).


## v0.25.2

- Fixed realtime delete event not being fired for `RecordProxy`-ies and added basic realtime record resolve automated tests ([#6433](https://github.com/pocketbase/pocketbase/issues/6433)).


## v0.25.1

- Fixed the batch API Preview success sample response.

- Bumped GitHub action min Go version to 1.23.6 as it comes with a [minor security fix](https://github.com/golang/go/issues?q=milestone%3AGo1.23.6+label%3ACherryPickApproved) for the ppc64le build.


## v0.25.0

- ⚠️ Upgraded Google OAuth2 auth, token and userinfo endpoints to their latest versions.
    _For users that don't do anything custom with the Google OAuth2 data or the OAuth2 auth URL, this should be a non-breaking change. The exceptions that I could find are:_
    - `/v3/userinfo` auth response changes:
        ```
        meta.rawUser.id             => meta.rawUser.sub
        meta.rawUser.verified_email => meta.rawUser.email_verified
        ```
    - `/v2/auth` query parameters changes:
        If you are specifying custom `approval_prompt=force` query parameter for the OAuth2 auth URL, you'll have to replace it with **`prompt=consent`**.

- Added Trakt OAuth2 provider ([#6338](https://github.com/pocketbase/pocketbase/pull/6338); thanks @aidan-)

- Added support for case-insensitive password auth based on the related UNIQUE index field collation ([#6337](https://github.com/pocketbase/pocketbase/discussions/6337)).

- Enforced `when_required` for the new AWS SDK request and response checksum validations to allow other non-AWS vendors to catch up with new AWS SDK changes (see [#6313](https://github.com/pocketbase/pocketbase/discussions/6313) and [aws/aws-sdk-go-v2#2960](https://github.com/aws/aws-sdk-go-v2/discussions/2960)).
    _You can set the environment variables `AWS_REQUEST_CHECKSUM_CALCULATION` and `AWS_RESPONSE_CHECKSUM_VALIDATION` to `when_supported` if your S3 vendor supports the [new default integrity protections](https://docs.aws.amazon.com/sdkref/latest/guide/feature-dataintegrity.html)._

- Soft-deprecated `Record.GetUploadedFiles` in favor of `Record.GetUnsavedFiles` to minimize the ambiguities what the method do ([#6269](https://github.com/pocketbase/pocketbase/discussions/6269)).

- Replaced archived `github.com/AlecAivazis/survey` dependency with a simpler  `osutils.YesNoPrompt(message, fallback)` helper.

- Upgraded to `golang-jwt/jwt/v5`.

- Added JSVM `new Timezone(name)` binding for constructing `time.Location` value ([#6219](https://github.com/pocketbase/pocketbase/discussions/6219)).

- Added `inflector.Camelize(str)` and `inflector.Singularize(str)` helper methods.

- Use the non-transactional app instance during the realtime records delete access checks to ensure that cascade deleted records with API rules relying on the parent will be resolved.

- Other minor improvements (_replaced all `bool` exists db scans with `int` for broader drivers compatibility, updated API Preview sample error responses, updated UI dependencies, etc._)


## v0.24.4

- Fixed fields extraction for view query with nested comments ([#6309](https://github.com/pocketbase/pocketbase/discussions/6309)).

- Bumped GitHub action min Go version to 1.23.5 as it comes with some [minor security fixes](https://github.com/golang/go/issues?q=milestone%3AGo1.23.5).


## v0.24.3

- Fixed incorrectly reported unique validator error for fields starting with name of another field ([#6281](https://github.com/pocketbase/pocketbase/pull/6281); thanks @svobol13).

- Reload the created/edited records data in the RecordsPicker UI.

- Updated Go dependencies.


## v0.24.2

- Fixed display fields extraction when there are multiple "Presentable" `relation` fields in a single related collection ([#6229](https://github.com/pocketbase/pocketbase/issues/6229)).


## v0.24.1

- Added missing time macros in the UI autocomplete.

- Fixed JSVM types for structs and functions with multiple generic parameters.


## v0.24.0

- ⚠️ Removed the "dry submit" when executing the collections Create API rule
    (you can find more details why this change was introduced and how it could affect your app in https://github.com/pocketbase/pocketbase/discussions/6073).
    For most users it should be non-breaking change, BUT if you have Create API rules that uses self-references or view counters you may have to adjust them manually.
    With this change the "multi-match" operators are also normalized in case the targeted collection doesn't have any records
    (_or in other words, `@collection.example.someField != "test"` will result to `true` if `example` collection has no records because it satisfies the condition that all available "example" records mustn't have `someField` equal to "test"_).
    As a side-effect of all of the above minor changes, the record create API performance has been also improved ~4x times in high concurrent scenarios (500 concurrent clients inserting total of 50k records - [old (58.409064001s)](https://github.com/pocketbase/benchmarks/blob/54140be5fb0102f90034e1370c7f168fbcf0ddf0/results/hetzner_cax41_cgo.md#creating-50000-posts100k-reqs50000-conc500-rulerequestauthid----requestdatapublicisset--true) vs [new (13.580098262s)](https://github.com/pocketbase/benchmarks/blob/7df0466ac9bd62fe0a1056270d20ef82012f0234/results/hetzner_cax41_cgo.md#creating-50000-posts100k-reqs50000-conc500-rulerequestauthid----requestbodypublicisset--true)).

- ⚠️ Changed the type definition of `store.Store[T any]` to `store.Store[K comparable, T any]` to allow support for custom store key types.
    For most users it should be non-breaking change, BUT if you are calling `store.New[any](nil)` instances you'll have to specify the store key type, aka. `store.New[string, any](nil)`.

- Added `@yesterday` and `@tomorrow` datetime filter macros.

- Added `:lower` filter modifier (e.g. `title:lower = "lorem"`).

- Added `mailer.Message.InlineAttachments` field for attaching inline files to an email (_aka. `cid` links_).

- Added cache for the JSVM `arrayOf(m)`, `DynamicModel`, etc. dynamic `reflect` created types.

- Added auth collection select for the settings "Send test email" popup ([#6166](https://github.com/pocketbase/pocketbase/issues/6166)).

- Added `record.SetRandomPassword()` to simplify random password generation usually used in the OAuth2 or OTP record creation flows.
    _The generated ~30 chars random password is assigned directly as bcrypt hash and ignores the `password` field plain value validators like min/max length or regex pattern._

- Added option to list and trigger the registered app level cron jobs via the Web API and UI.

- Added extra validators for the collection field `int64` options (e.g. `FileField.MaxSize`) restricting them to the max safe JSON number (2^53-1).

- Added option to unset/overwrite the default PocketBase superuser installer using `ServeEvent.InstallerFunc`.

- Added `app.FindCachedCollectionReferences(collection, excludeIds)` to speedup records cascade delete almost twice for projects with many collections.

- Added `tests.NewTestAppWithConfig(config)` helper if you need more control over the test configurations like `IsDev`, the number of allowed connections, etc.

- Invalidate all record tokens when the auth record email is changed programmatically or by a superuser ([#5964](https://github.com/pocketbase/pocketbase/issues/5964)).

- Eagerly interrupt waiting for the email alert send in case it takes longer than 15s.

- Normalized the hidden fields filter checks and allow targetting hidden fields in the List API rule.

- Fixed "Unique identify fields" input not refreshing on unique indexes change ([#6184](https://github.com/pocketbase/pocketbase/issues/6184)).


## v0.23.12

- Added warning logs in case of mismatched `modernc.org/sqlite` and `modernc.org/libc` versions ([#6136](https://github.com/pocketbase/pocketbase/issues/6136#issuecomment-2556336962)).

- Skipped the default body size limit middleware for the backup upload endpoint ([#6152](https://github.com/pocketbase/pocketbase/issues/6152)).


## v0.23.11

- Upgraded `golang.org/x/net` to 0.33.0 to fix [CVE-2024-45338](https://www.cve.org/CVERecord?id=CVE-2024-45338).
  _PocketBase uses the vulnerable functions primarily for the auto html->text mail generation, but most applications shouldn't be affected unless you are manually embedding unrestricted user provided value in your mail templates._


## v0.23.10

- Renew the superuser file token cache when clicking on the thumb preview or download link ([#6137](https://github.com/pocketbase/pocketbase/discussions/6137)).

- Upgraded `modernc.org/sqlite` to 1.34.3 to fix "disk io" error on arm64 systems.
    _If you are extending PocketBase with Go and upgrading with `go get -u` make sure to manually set in your go.mod the `modernc.org/libc` indirect dependency to v1.55.3, aka. the exact same version the driver is using._


## v0.23.9

- Replaced `strconv.Itoa` with `strconv.FormatInt` to avoid the int64->int conversion overflow on 32-bit platforms ([#6132](https://github.com/pocketbase/pocketbase/discussions/6132)).


## v0.23.8

- Fixed Model->Record and Model->Collection hook events sync for nested and/or inner-hook transactions ([#6122](https://github.com/pocketbase/pocketbase/discussions/6122)).

- Other minor improvements (updated Go and npm deps, added extra escaping for the default mail record params in case the emails are stored as html files, fixed code comment typos, etc.).


## v0.23.7

- Fixed JSVM exception -> Go error unwrapping when throwing errors from non-request hooks ([#6102](https://github.com/pocketbase/pocketbase/discussions/6102)).


## v0.23.6

- Fixed `$filesystem.fileFromURL` documentation and generated type ([#6058](https://github.com/pocketbase/pocketbase/issues/6058)).

- Fixed `X-Forwarded-For` header typo in the suggested UI "Common trusted proxy" headers ([#6063](https://github.com/pocketbase/pocketbase/pull/6063)).

- Updated the `text` field max length validator error message to make it more clear ([#6066](https://github.com/pocketbase/pocketbase/issues/6066)).

- Other minor fixes (updated Go deps, skipped unnecessary validator check when the default primary key pattern is used, updated JSVM types, etc.).


## v0.23.5

- Fixed UI logs search not properly accounting for the "Include requests by superusers" toggle when multiple search expressions are used.

- Fixed `text` field max validation error message ([#6053](https://github.com/pocketbase/pocketbase/issues/6053)).

- Other minor fixes (comment typos, JSVM types update).

- Updated Go deps and the min Go release GitHub action version to 1.23.4.


## v0.23.4

- Fixed `autodate` fields not refreshing when calling `Save` multiple times on the same `Record` instance ([#6000](https://github.com/pocketbase/pocketbase/issues/6000)).

- Added more descriptive test OTP id and failure log message ([#5982](https://github.com/pocketbase/pocketbase/discussions/5982)).

- Moved the default UI CSP from meta tag to response header ([#5995](https://github.com/pocketbase/pocketbase/discussions/5995)).

- Updated Go and npm dependencies.


## v0.23.3

- Fixed Gzip middleware not applying when serving static files.

- Fixed `Record.Fresh()`/`Record.Clone()` methods not properly cloning `autodate` fields ([#5973](https://github.com/pocketbase/pocketbase/discussions/5973)).


## v0.23.2

- Fixed `RecordQuery()` custom struct scanning ([#5958](https://github.com/pocketbase/pocketbase/discussions/5958)).

- Fixed `--dev` log query print formatting.

- Added support for passing more than one id in the `Hook.Unbind` method for consistency with the router.

- Added collection rules change list in the confirmation popup
  (_to avoid getting anoying during development, the rules confirmation currently is enabled only when using https_).


## v0.23.1

- Added `RequestEvent.Blob(status, contentType, bytes)` response write helper ([#5940](https://github.com/pocketbase/pocketbase/discussions/5940)).

- Added more descriptive error messages.


## v0.23.0

> [!NOTE]
> You don't have to upgrade to PocketBase v0.23.0 if you are not planning further developing
> your existing app and/or are satisfied with the v0.22.x features set. There are no identified critical issues
> with PocketBase v0.22.x yet and in the case of critical bugs and security vulnerabilities, the fixes
> will be backported for at least until Q1 of 2025 (_if not longer_).
>
> **If you don't plan upgrading make sure to pin the SDKs version to their latest PocketBase v0.22.x compatible:**
> - JS SDK: `<0.22.0`
> - Dart SDK: `<0.19.0`

> [!CAUTION]
> This release introduces many Go/JSVM and Web APIs breaking changes!
>
> Existing `pb_data` will be automatically upgraded with the start of the new executable,
> but custom Go or JSVM (`pb_hooks`, `pb_migrations`) and JS/Dart SDK code will have to be migrated manually.
> Please refer to the below upgrade guides:
> - Go:   https://pocketbase.io/v023upgrade/go/.
> - JSVM: https://pocketbase.io/v023upgrade/jsvm/.
>
> If you had already switched to some of the earlier `<v0.23.0-rc14` versions and have generated a full collections snapshot migration (aka. `./pocketbase migrate collections`), then you may have to regenerate the migration file to ensure that it includes the latest changes.

PocketBase v0.23.0 is a major refactor of the internals with the overall goal of making PocketBase an easier to use Go framework.
There are a lot of changes but to highlight some of the most notable ones:

- New and more [detailed documentation](https://pocketbase.io/docs/).
  _The old documentation could be accessed at [pocketbase.io/old](https://pocketbase.io/old/)._
- Replaced `echo` with a new router built on top of the Go 1.22 `net/http` mux enhancements.
- Merged `daos` packages in `core.App` to simplify the DB operations (_the `models` package structs are also migrated in `core`_).
- Option to specify custom `DBConnect` function as part of the app configuration to allow different `database/sql` SQLite drivers (_turso/libsql, sqlcipher, etc._) and custom builds.
  _Note that we no longer loads the `mattn/go-sqlite3` driver by default when building with `CGO_ENABLED=1` to avoid `multiple definition` linker errors in case different CGO SQLite drivers or builds are used. You can find an example how to enable it back if you want to in the [new documentation](https://pocketbase.io/docs/go-overview/#github-commattngo-sqlite3)._
- New hooks allowing better control over the execution chain and error handling (_including wrapping an entire hook chain in a single DB transaction_).
- Various `Record` model improvements (_support for get/set modifiers, simplfied file upload by treating the file(s) as regular field value like `record.Set("document", file)`, etc._).
- Dedicated fields structs with safer defaults to make it easier creating/updating collections programmatically.
- Option to mark field as "Hidden", disallowing regular users to read or modify it (_there is also a dedicated Record hook to hide/unhide Record fields programmatically from a single place_).
- Option to customize the default system collection fields (`id`, `email`, `password`, etc.).
- Admins are now system `_superusers` auth records.
- Builtin rate limiter (_supports tags, wildcards and exact routes matching_).
- Batch/transactional Web API endpoint.
- Impersonate Web API endpoint (_it could be also used for generating fixed/nonrenewable superuser tokens, aka. "API keys"_).
- Support for custom user request activity log attributes.
- One-Time Password (OTP) auth method (_via email code_).
- Multi-Factor Authentication (MFA) support (_currently requires any 2 different auth methods to be used_).
- Support for Record "proxy/projection" in preparation for the planned autogeneration of typed Go record models.
- Linear OAuth2 provider ([#5909](https://github.com/pocketbase/pocketbase/pull/5909); thanks @chnfyi).
- WakaTime OAuth2 provider ([#5829](https://github.com/pocketbase/pocketbase/pull/5829); thanks @tigawanna).
- Notion OAuth2 provider ([#4999](https://github.com/pocketbase/pocketbase/pull/4999); thanks @s-li1).
- monday.com OAuth2 provider ([#5346](https://github.com/pocketbase/pocketbase/pull/5346); thanks @Jaytpa01).
- New Instagram provider compatible with the new Instagram Login APIs ([#5588](https://github.com/pocketbase/pocketbase/pull/5588); thanks @pnmcosta).
    _The provider key is `instagram2` to prevent conflicts with existing linked users._
- Option to retrieve the OIDC OAuth2 user info from the `id_token` payload for the cases when the provider doesn't have a dedicated user info endpoint.
- Various minor UI improvements (_recursive `Presentable` view, slightly different collection options organization, zoom/pan for the logs chart, etc._)
- and many more...

#### Go/JSVM APIs changes

> - Go:   https://pocketbase.io/v023upgrade/go/.
> - JSVM: https://pocketbase.io/v023upgrade/jsvm/.

#### SDKs changes

- [JS SDK v0.22.0](https://github.com/pocketbase/js-sdk/blob/master/CHANGELOG.md)
- [Dart SDK v0.19.0](https://github.com/pocketbase/dart-sdk/blob/master/CHANGELOG.md)

#### Web APIs changes

- New `POST /api/batch` endpoint.

- New `GET /api/collections/meta/scaffolds` endpoint.

- New `DELETE /api/collections/{collection}/truncate` endpoint.

- New `POST /api/collections/{collection}/request-otp` endpoint.

- New `POST /api/collections/{collection}/auth-with-otp` endpoint.

- New `POST /api/collections/{collection}/impersonate/{id}` endpoint.

- ⚠️ If you are constructing requests to `/api/*` routes manually remove the trailing slash (_there is no longer trailing slash removal middleware registered by default_).

- ⚠️ Removed `/api/admins/*` endpoints because admins are converted to `_superusers` auth collection records.

- ⚠️ Previously when uploading new files to a multiple `file` field, new files were automatically appended to the existing field values.
     This behaviour has changed with v0.23+ and for consistency with the other multi-valued fields when uploading new files they will replace the old ones. If you want to prepend or append new files to an existing multiple `file` field value you can use the `+` prefix or suffix:
     ```js
     "documents": [file1, file2]  // => [file1_name, file2_name]
     "+documents": [file1, file2] // => [file1_name, file2_name, old1_name, old2_name]
     "documents+": [file1, file2] // => [old1_name, old2_name, file1_name, file2_name]
     ```

- ⚠️ Removed `GET /records/{id}/external-auths` and `DELETE /records/{id}/external-auths/{provider}` endpoints because this is now handled by sending list and delete requests to the `_externalAuths` collection.

- ⚠️ Changes to the app settings model fields and response (+new options such as `trustedProxy`, `rateLimits`, `batch`, etc.). The app settings Web APIs are mostly used by the Dashboard UI and rarely by the end users, but if you want to check all settings changes please refer to the [Settings Go struct](https://github.com/pocketbase/pocketbase/blob/develop/core/settings_model.go#L121).

- ⚠️ New flatten Collection model and fields structure. The Collection model Web APIs are mostly used by the Dashboard UI and rarely by the end users, but if you want to check all changes please refer to the [Collection Go struct](https://github.com/pocketbase/pocketbase/blob/develop/core/collection_model.go#L308).

- ⚠️ The top level error response `code` key was renamed to `status` for consistency with the Go APIs.
    The error field key remains `code`:
    ```js
    {
        "status": 400, // <-- old: "code"
        "message": "Failed to create record.",
        "data": {
            "title": {
                "code": "validation_required",
                "message": "Missing required value."
            }
        }
    }
    ```

- ⚠️ New fields in the `GET /api/collections/{collection}/auth-methods` response.
    _The old `authProviders`, `usernamePassword`, `emailPassword` fields are still returned in the response but are considered deprecated and will be removed in the future._
    ```js
    {
        "mfa": {
            "duration": 100,
            "enabled": true
        },
        "otp": {
            "duration": 0,
            "enabled": false
        },
        "password": {
            "enabled": true,
            "identityFields": ["email", "username"]
        },
        "oauth2": {
            "enabled": true,
            "providers": [{"name": "gitlab", ...}, {"name": "google", ...}]
        },
        // old fields...
    }
    ```

- ⚠️ Soft-deprecated the OAuth2 auth success `meta.avatarUrl` field in favour of `meta.avatarURL`.
