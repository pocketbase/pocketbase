## v0.8.0

**⚠️ This release contains breaking changes and requires some manual migration steps!**

The biggest change is the merge of the `User` models and the `profiles` collection per [#376](https://github.com/pocketbase/pocketbase/issues/376).
There is no longer `user` type field and the users are just an "auth" collection (we now support **collection types**, currently only "base" and "auth").
This should simplify the users management and at the same time allow us to have unlimited multiple "auth" collections each with their own custom fields and authentication options (eg. staff, client, etc.).

In addition to the `Users` and `profiles` merge, this release comes with several other improvements:

- Added indirect expand support [#312](https://github.com/pocketbase/pocketbase/issues/312#issuecomment-1242893496).

- The `json` field type now supports filtering and sorting [#423](https://github.com/pocketbase/pocketbase/issues/423#issuecomment-1258302125).

- The `relation` field now allows unlimitted `maxSelect` (aka. without upper limit).

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

- Marked as "Deprecated" and will be removed in v0.9:
    ```
    core.Settings.EmailAuth{}
    core.EmailAuthConfig{}
    schema.FieldTypeUser
    schema.UserOptions{}
    ```

- The second argument of `apis.StaticDirectoryHandler(fileSystem, enableIndexFallback)` now is used to enable/disable index.html forwarding on missing file (eg. in case of SPA).
