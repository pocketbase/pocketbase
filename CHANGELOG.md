## v0.23.0-rc10

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Restore the CRC32 checksum autogeneration for the collection/field ids in order to maintain deterministic default identifier value and minimize conflicts between custom migrations and full collections snapshots.
  _There is a system migration that will attempt to normalize existing system collections ids, but if you already migrated to v0.23.0-rc and have generated a full collections snapshot migration, you have to delete it and regenerate a new one._

- Change the behavior of the default generated collections snapshot migration to act as "extend" instead of "replace" to prevent accidental data deletion.
  _I think this would be rare but if you want the old behaviour you can edit the generated snapshot file and replace the second argument (`deleteMissing`) of `App.ImportCollection/App.ImportCollectionsByMarshaledJSON` from `false` to `true`._


## v0.23.0-rc9

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Fixed auto `www.` redirect due to missing URI schema.

- Fixed collection and field renaming when reusing an old collection/field name ([#5741](https://github.com/pocketbase/pocketbase/issues/5741)).

- Update the "API preview" section to include information about the batch api.

- Exported `core.DefaultDBConnect` function that could be used as a fallback when initializing custom SQLite drivers and builds.

- ⚠️ No longer loads the `mattn/go-sqlite3` driver by default when building with `CGO_ENABLED=1` to avoid `multiple definition ...` linker errors in case different CGO SQLite drivers or builds are used.
    This means that no matter of the `CGO_ENABLED` value, now out of the box PocketBase will always use only the pure Go driver ([`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite)).
    This will be documented properly in the new website but if you want to continue using `mattn/go-sqlite3` (e.g. because of the icu or other builtin extension) you could register it as follow:
    ```go
    package main

    import (
        "database/sql"
        "log"

        "github.com/mattn/go-sqlite3"
        "github.com/pocketbase/dbx"
        "github.com/pocketbase/pocketbase"
    )

    func init() {
        // initialize default PRAGMAs for each new connection
        sql.Register("pb_sqlite3",
            &sqlite3.SQLiteDriver{
                ConnectHook: func(conn *sqlite3.SQLiteConn) error {
                    _, err := conn.Exec(`
                        PRAGMA busy_timeout       = 10000;
                        PRAGMA journal_mode       = WAL;
                        PRAGMA journal_size_limit = 200000000;
                        PRAGMA synchronous        = NORMAL;
                        PRAGMA foreign_keys       = ON;
                        PRAGMA temp_store         = MEMORY;
                        PRAGMA cache_size         = -16000;
                    `, nil)

                    return err
                },
            },
        )

        dbx.BuilderFuncMap["pb_sqlite3"] = dbx.BuilderFuncMap["sqlite3"]
    }

    func main() {
        app := pocketbase.NewWithConfig(pocketbase.Config{
            DBConnect: func(dbPath string) (*dbx.DB, error) {
                return dbx.Open("pb_sqlite3", dbPath)
            },
        })

        // custom hooks and plugins...

        if err := app.Start(); err != nil {
            log.Fatal(err)
        }
    }
    ```
    Also note that if you are not planning to use the `core.DefaultDBConnect` fallback as part of your custom driver registration you can exclude the default pure Go driver from the build with the build tag `-tags no_default_driver` to reduce the binary size a little.

- ⚠️ Removed JSVM `BaseCollection()`, `AuthCollection()`, `ViewCollection()` class aliases for simplicity and to avoid confusion with the accepted constructor arguments (_you can simply use as before `new Collection({ type: "base", ... })`; this will also initialize the default type specific options_).

- Other minor improvements (added validator for duplicated index definitions, updated the impersonate popup styles, added query param support for loading a collection based on its name, etc.).


## v0.23.0-rc8

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Lock the `_otps` and `_mfas` system collections Delete API rule for superusers only.

- Reassign in the JSVM executors the global `$app` variable with the hook scoped `e.app` value to minimize the risk of a deadlock when a hook or middleware is wrapped in a transaction.

- Reuse the OAuth2 created user record pointer to ensure that all its following hooks operate on the same record instance.

- Added tags support for the `OnFileTokenRequest` hook.

- Other minor changes (added index for the `_collections` type column, added more detailed godoc for the collection fields and `core.App` methods, fixed flaky record enrich tests, etc.).


## v0.23.0-rc7

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Register the default panic-recover middleware after the activity logger so that we can log the error.

- Updated the `RequestEvent.BindBody` FormData type inferring rules to convert numeric strings into float64 only if the resulting minimal number string representation matches the initial FormData string value ([#5687](https://github.com/pocketbase/pocketbase/issues/5687)).

- Fixed the JSVM types to include properly generated function declarations when the related Go functions have shortened/combined return values.

- Reorganized the record table fields<->columns syncing to remove the `PRAGMA writable_schema` usage.


## v0.23.0-rc6

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Fixed realtime 403 API error on resubscribe ([#5674](https://github.com/pocketbase/pocketbase/issues/5674)).

- Fixed the auto OAuth2 avatar mapped field assignment when the OAuth2 provider doesn't return an avatar URL ([#5673](https://github.com/pocketbase/pocketbase/pull/5673)).
  _In case the avatar retrieval fails and the mapped record field "Required" option is not set, the error is silenced and only logged with WARN level._

- Added `Router.SEARCH(path, action)` helper method for registering `SEARCH` endpoints.

- Changed all builtin middlewares to return `*hook.Handler[*core.RequestEvent]` with a default middleware id for consistency and to allow removal.
  Or in other words, replace `.BindFunc(apis.Gzip())` with `.Bind(apis.Gzip())`.

- Updated the JSVM types to reflect the recent changes.


## v0.23.0-rc5

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Added Notion OAuth2 provider ([#4999](https://github.com/pocketbase/pocketbase/pull/4999); thanks @s-li1).

- Added monday.com OAuth2 provider ([#5346](https://github.com/pocketbase/pocketbase/pull/5346); thanks @Jaytpa01).

- Added option to retrieve the OIDC OAuth2 user info from the `id_token` payload for the cases when the provider doesn't have a dedicated user info endpoint.

- Fixed the relation record picker to sort by default by `@rowid` instead of the `created` field as the latter is optional ([#5641](https://github.com/pocketbase/pocketbase/discussions/5641)).

- Fixed the UI "Set Superusers only" button click not properly resetting the input state.

- Fixed the OAuth2 providers logo path shown in the "Authorized providers" UI.

- Fixed the single value UI for the `select`, `file` and `relation` fields ([#5646](https://github.com/pocketbase/pocketbase/discussions/5646))


## v0.23.0-rc4

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Fixed the UI settings update form to prevent sending empty string for the mail password or the S3 secret options on resave of the form.

- ⚠️ Added an exception for the `OAuth2` field in the GO->JSVM name mapping rules:
    ```
    // old              -> new
    collection.oAuth2.* -> collection.oauth2.*
    ```

- Added more user friendly view collection truncate error message.

- Added an extra suffix character to the name of autogenerated template migration file for `*test` suffixed collections to prevent acidentally resulting in `_test.go` migration files.

- Added `FieldsList.AddMarshaledJSON([]byte)` helper method to load a serialized json array of objects or a single json object into an existing collection fields list.

- Fixed the autogenerated Go migration template when updating a collection ([#5631](https://github.com/pocketbase/pocketbase/discussions/5631)).

    ⚠️ If you have already used a previous prerelease and have autogenerated Go migration files, please check the migration files named **`{timestamp}_updated_{collection}.go`** and manually change:

    <table width="100%">
        <tr>
            <th width="50%">Old (broken)</th>
            <th width="50%">New</th>
        </tr>
        <tr>
    <td width="50%">

    ```go
    // add field / update field
    if err := json.Unmarshal([]byte(`[{
        ...
    }]`), &collection.Fields); err != nil {
        return err
    }
    ```

    </td>
    <td width="50%">

    ```go
    // add field / update field
    if err := collection.Fields.AddMarshaledJSON([]byte(`{
        ...
    }`)); err != nil {
        return err
    }
    ```

    </td>
        </tr>
    </table>
    To test that your Go migration files work correctly you can try to start PocketBase with a new temp pb_data, e.g.:

    ```go
    go run . serve --dir="pb_data_temp"
    ```


## v0.23.0-rc3

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Make `PRAGMA optimize` statement optional in case it is not supported by the driver ([#5611](https://github.com/pocketbase/pocketbase/discussions/5611)).

- Reapply the minimum required `pb_data/auxiliary.db` migrations if the db file was manually deleted ([#5618](https://github.com/pocketbase/pocketbase/discussions/5618)).

- To avoid confusion and unnecessary casting, the `hook.HandlerFunc[T]` type has been removed and instead everywhere we now use directly the underlying function definition, aka.:
  ```go
  func(T) error
  ```

- Fixed the UI input field type of the OTP.length field ([#5617](https://github.com/pocketbase/pocketbase/issues/5617)).

- Other minor fixes (fixed API preview and examples error message typos, better hint for combined/multi-spaced view query columns, fixed the path for the HTTPS green favicon path, etc.).


## v0.23.0-rc2

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**

- Small update to the earlier v0.23.0-rc that uses `pb_data/auxiliary.db` instead of `pb_data/aux.db` because it seems that on Windows `aux` is disallowed as file name ([#5607](https://github.com/pocketbase/pocketbase/issues/5607)).
   _If you have already upgraded to v0.23.0-rc please rename manually your `pb_data/aux.db` file to `pb_data/auxiliary.db`._


## v0.23.0-rc

> [!CAUTION]
> **This is a prerelease intended for test and experimental purposes only!**
>
> It introduces many Go/JSVM breaking changes and requires manual migration steps.
>
> All new features will be reflected in the new website documentation with the final v0.23.0 release.

> [!NOTE]
>  Please note that you don't have to upgrade to PocketBase v0.23.0 if you are not planning further developing
> your existing app and/or are satisfied with the v0.22.x features set. There are no identified critical issues
> with PocketBase v0.22.x yet and in the case of critical bugs and security vulnerabilities, the fixes
> will be backported for at least until Q1 of 2025 (_if not longer_).
>
> If you don't plan upgrading just make sure to pin the SDKs version to their latest PocketBase v0.22.x compatible:
> - JS SDK: `<0.22.0`
> - Dart SDK: `<0.19.0`

PocketBase v0.23.0-rc is a major refactor of the internals with the overall goal of making PocketBase an easier to use Go framework.

There are many changes but to highlight some of the most notable ones:

- Replaced `echo` with a new router built on top of the Go 1.22 `net/http` mux enhancements.
- Merged `daos` packages in `core.App` to simplify the DB operations (_the `models` package structs are also migrated in `core`_).
- Option to specify custom `DBConnect` function as part of the app configuration to allow different `database/sql` SQLite drivers (_turso/libsql, sqlcipher, etc._) and custom builds.
- New hooks allowing better control over the execution chain and error handling (_including wrapping an entire hook chain in a single DB transaction_).
- Various `Record` model improvements (_support for get/set modifiers, simplfied file upload by treating the file(s) as regular field value like `record.Set("document", file)`, etc._).
- Dedicated fields structs with safer defaults to make it easier creating/updating collections programmatically.
- Option to mark field as Private/Hidden, disallowing regular users to read or modify it (_there is also a dedicated Record hook to hide/unhide Record fields programmatically from a single place_).
- Option to customize the default system collection fields (`id`, `email`, `password`, etc.).
- Admins are now system `_superusers` auth records.
- Builtin rate limiter (_supports tags, wildcards and exact routes matching_).
- Batch/transactional Web API endpoint.
- Impersonate Web API endpoint (_it could be also used for generating fixed/non-refreshable superuser tokens, aka. "API keys"_).
- Support for custom user request activity log attributes.
- One-Time Password (OTP) auth method (_via email code_).
- Multi-Factor Authentication (MFA) support (_currently requires any 2 different auth methods to be used_).
- Support for Record "proxy/projection" in preparation for the planned autogeneration of typed Go record models.
- Various minor UI improvements (_recursive `Presentable` view, slightly different collection options organization, zoom/pan for the logs chart, etc._)
- and many more...

In terms of performance, the Go standard router mux is known to be slightly slower compared to Gin, Echo, etc. implementations, but based on my local tests the difference is negliable.
The [benchmarks repo](https://github.com/pocketbase/benchmarks) will be updated with the final v0.23.0 release (_currently there seems to be ~10% memory consumption increase which I'll have to investigate to see whether it is from the router change or from the new hooks_).

#### Go/JSVM APIs changes

For upgrading to PocketBase v0.23.0, please refer to:

- Go:   https://pocketbase.io/v023upgrade/go/.
- JSVM: https://pocketbase.io/v023upgrade/jsvm/.

#### SDKs changes

- [JS SDK v0.22.0-rc](https://github.com/pocketbase/js-sdk/blob/develop/CHANGELOG.md)
- [Dart SDK v0.19.0-rc](https://github.com/pocketbase/dart-sdk/blob/develop/CHANGELOG.md)

#### Web APIs changes

- New `POST /api/batch` endpoint.

- New `GET /api/collections/meta/scaffolds` endpoint.

- New `DELETE /api/collections/{collection}/truncate` endpoint.

- New `POST /api/collections/{collection}/request-otp` endpoint.

- New `POST /api/collections/{collection}/auth-with-otp` endpoint.

- New `POST /api/collections/{collection}/impersonate/{id}` endpoint.

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
