## v0.40.0 (WIP)

- Propagate console command errors and recovered panics to `app.Start()` so that the program can exit with non-zero code while still ensuring that `app.OnTerminate` hook (responsible for the app graceful shutdown handling) was triggered.
    _⚠️ Note that this could be a slight breaking change in case you are chaining PocketBase commands and relied on the previous `0` exit status for `Command.RunE` returned errors._
    _Or in other words, if you have `./pocketbase invalid && someothercommand` and previously relied that `someothercommand` will be always executed then this is no longer the case and you'll have to adjust it or replace `&&` with `;`._

- Added `filesystem.NewWriter(key, opts)` low-level helper to allow direct file create from an `io.Reader` value.

- Added new `DELETE /api/logs` endpoint and UI control to delete all logs without changing the `maxDays` retention setting.

- Added quotes around the default `Content-Disposition` serving filename in case custom name with special characters is provided.

- Added `Cross-Origin-Opener-Policy:same-origin` to the default security response headers.
    _This is an extra precaution to prevent tab-nabbing in case custom UI plugins use `target="_blank"` without `rel="noopener"`._

- Added new log settings option to limit the max `Log.Data` size that will be saved in the database (default to ~16KB).
    _This is an extra precaution for the cases when logging user supplied data without validating it beforehand._
    _If the resulting `Log.Data` json is above the limit, it is truncated to the last valid decoded character and an extra `"__pb_truncated__":true` log data entry will be added.`_
    _Additionally, for just in case the log message is also truncated at max 8k characters._

- Updated `modernc.org/sqlite` to 1.57 and registered by default the new `_defensive=1` DSN query parameter to enable [SQLite's defensive mode](https://sqlite.org/c3ref/c_dbconfig_defensive.html#sqlitedbconfigdefensive).

- (@todo) Bumped the min Go version to 1.27.0 and migrated to the new `encoding/json/v2` package.
