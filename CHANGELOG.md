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

- Updated Go deps and the min Go releleaser GitHub action version to 1.23.4.


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
- Impersonate Web API endpoint (_it could be also used for generating fixed/non-refreshable superuser tokens, aka. "API keys"_).
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
