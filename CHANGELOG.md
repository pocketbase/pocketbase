## v0.40.2

- Return an error when filter params fallback fails to json serialize and optimized params replacement to execute in a single pass.

- Fixed collection index parsing error for indexes with missing name.

- Minor UI autocomplete optimizations _(prefix match, autocomplete debounce, etc.)_.

- Fixed linter warnings and comment typos.

- Bumped goja and its related dependencies _(regex unescaped dash error fix and base64 optimizations)_.

- Bumped the min Go GitHub action version to 1.27.1 as it includes some [minor `database/sql` and `enconding/json/v2` bug fixes](https://github.com/golang/go/issues?q=milestone%3AGo1.27.1).


## v0.40.1

- Fixes for some reported regressions related to the `encoding/json/v2` update:
    - allow mangling invalid UTF8 characters when serializing json data ([#7814](https://github.com/pocketbase/pocketbase/issues/7814))
    - fixed OAuth2 providers config merge incorrectly replacing the entire slice ([#7815](https://github.com/pocketbase/pocketbase/issues/7815))


## v0.40.0

- Propagate console command errors and recovered panics to `app.Start()` so that the program can exit with non-zero code while still ensuring that `app.OnTerminate` hook was triggered _(responsible for the app graceful shutdown handling)_.
    _⚠️ Note that this could be a slight breaking change in case you are chaining PocketBase commands and relied on the previous `0` exit status for `Command.RunE` returned errors._
    _Or in other words, if you have `./pocketbase invalid && someothercommand` and previously relied that `someothercommand` will be always executed then this is no longer the case and you'll have to adjust it or replace `&&` with `;`._

- Added quotes around the default `Content-Disposition` serving filename in case custom name with special characters is provided.

- Added `Cross-Origin-Opener-Policy:same-origin` to the default security response headers.
    _This is an extra precaution to prevent tab-nabbing in case custom UI plugins use `target="_blank"` without `rel="noopener"`._

- Added `Record.GetInt64(field)` helper (note that the serializable max safe integer of the `number` field is ~2^53-1).

- Added `Store.Keys()` method that returns a slice with all of the store keys.

- Added new `DELETE /api/logs` endpoint and UI control to delete all logs without changing the `maxDays` retention setting.

- Added new log settings option to limit the max `Log.Data` size that will be saved in the database (default to ~16KB).
    _This is an extra precaution for the cases when logging user supplied data without validating it beforehand._
    _If the resulting `Log.Data` json is above the limit, it is truncated to the last valid decoded character and an extra `"__pb_truncated__":true` log data entry will be added.`_
    _Additionally, for just in case the log message is also truncated at max 8k characters._

- Added new `filesystem` low-level helper methods:
    - `filesystem.NewWriter(key, opts)` to allow direct file create from an `io.Reader` value.
    - `filesystem.OnNewWriter()` hook to allow listening for new/to-be-created files _(it is not exposed in `core.App` instance for now to avoid introducing breaking changes)_.
    - `filesystem.OnDelete()` hook to allow listening for deleted files _(it is not exposed in `core.App` instance for now to avoid introducing breaking changes)_.

- Optimized backups to no longer transaction lock the database during backup generation ([#7799](https://github.com/pocketbase/pocketbase/discussions/7799#discussioncomment-18108244)).

- Updated `modernc.org/sqlite` to 1.57.0 and registered by default the new `_defensive=1` DSN query parameter to enable [SQLite's defensive mode](https://sqlite.org/c3ref/c_dbconfig_defensive.html#sqlitedbconfigdefensive).

- Bumped the min Go version to 1.27.0 and migrated to the new `encoding/json/v2` package.
    _⚠️ Please note that Go 1.27.0 retrofitted `encoding/json` to use the v2 package under the hood but unfortunately is not fully backward compatible._
    _I recommend to not push blindly an update on production and to test your PocketBase application first locally to see if everything works correctly._
