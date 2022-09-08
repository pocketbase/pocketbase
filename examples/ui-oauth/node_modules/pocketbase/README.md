PocketBase JavaScript SDK
======================================================================

Official JavaScript SDK (browser and node) for interacting with the [PocketBase API](https://pocketbase.io/docs).

- [Installation](#installation)
- [Examples](#examples)
- [Caveats](#caveats)
- [Definitions](#definitions)
- [Development](#development)


## Installation

#### Browser (manually via script tag)

```html
<script src="/path/to/dist/pocketbase.umd.js"></script>
```

_OR if you are using ES modules:_
```html
<script type="module">
    import PocketBase from '/path/to/dist/pocketbase.es.mjs'
    ...
</script>
```

#### Node.js (via npm)

```sh
npm install pocketbase --save
```

```js
// Using ES modules (default)
import PocketBase from 'pocketbase'

// OR if you are using CommonJS modules
const PocketBase = require('pocketbase/cjs')
```

> ⚠️ For Node < 17 you may need to add a `fetch()` polyfill.
    I recommend [cross-fetch](https://github.com/lquixada/cross-fetch):
    ```js
    import 'cross-fetch/polyfill';
    ```


## Examples

```js
import PocketBase from 'pocketbase';

const client = new PocketBase('http://localhost:8090');

...

// list and filter "demo" collection records
const result = await client.Records.getList("demo", 1, 20, {
    filter: "status = true && totalComments > 0"
});

// authenticate as regular user
const userData = await client.Users.authViaEmail("test@example.com", "123456");


// or as admin
const adminData = await client.Admins.authViaEmail("test@example.com", "123456");

// and much more...
```
> More detailed API docs and copy-paste examples could be found in the [API documentation for each service](https://pocketbase.io/docs/api-authentication).


## Caveats

#### Errors handling

Each request error response is wrapped and in a normalized `ClientResponseError` object with the following fields:
```js
ClientResponseError {
    url:           string,     // requested url
    status:        number,     // response status code
    data:          { ... },    // the api json error response
    isAbort:       boolean,    // is abort/cancellation error
    originalError: Error|null, // the original non-normalized error
}
```

#### Auto cancellation

The SDK client auto cancel duplicated pending requests.
For example, if you have the following 3 duplicated calls, only the last one will be executed, while the first 2 will be cancelled with error `null`:

```js
client.Records.getList("demo", 1, 20) // cancelled
client.Records.getList("demo", 1, 20) // cancelled
client.Records.getList("demo", 1, 20) // executed
```

To change this behavior, you could make use of 2 special query parameters:

- `$autoCancel ` - set it to `false` to disable auto cancellation for this request
- `$cancelKey` - set it to a string that will be used as request identifier and based on which pending duplicated requests will be matched (default to `HTTP_METHOD + url`, eg. "get /api/users?page=1")

Example:

```js
client.Records.getList("demo", 1, 20);                           // cancelled
client.Records.getList("demo", 1, 20);                           // executed
client.Records.getList("demo", 1, 20, { "$autoCancel": false }); // executed
client.Records.getList("demo", 1, 20, { "$autoCancel": false })  // executed
client.Records.getList("demo", 1, 20, { "$cancelKey": "test" })  // cancelled
client.Records.getList("demo", 1, 20, { "$cancelKey": "test" })  // executed
```

To manually cancel pending requests, you could use `client.cancelAllRequests()` or `client.cancelRequest(cancelKey)`.

#### AuthStore

The SDK keeps track of the authenticated token and auth model for you via the `client.AuthStore` instance.
It has the following public fields that you can use:

```js
AuthStore {
    token:   string,     // the authenticated token
    model:   User|Admin, // the authenticated User or Admin model
    isValid: boolean,    // checks if the store has existing and unexpired token
    clear(),             // "logout" the authenticated User or Admin
    save(token, model),  // update the store with the new auth data
}
```

By default the SDK initialize a `LocalAuthStore`, which uses the browser's `LocalStorage` if available, otherwise - will fallback to runtime/memory store (aka. on page refresh or service restart you'll have to authenticate again).

To "logout" an authenticated user or admin, you can just call `client.AuthStore.clear()`.

## Definitions

#### Creating new client instance

```js
const client = new PocketBase(
    baseUrl = '/',
    lang = 'en-US',
    authStore = LocalAuthStore
);
```

#### Instance methods

> Each instance method returns the `PocketBase` instance allowing chaining.

| Method                                  | Description                                                         |
|:----------------------------------------|:--------------------------------------------------------------------|
| `client.send(path, reqConfig = {})`     | Sends an api http request.                                          |
| `client.cancelAllRequests()`            | Cancels all pending requests.                                       |
| `client.cancelRequest(cancelKey)`       | Cancels single request by its cancellation token key.               |
| `client.buildUrl(path, reqConfig = {})` | Builds a full client url by safely concatenating the provided path. |


#### API services

> Each service call returns a `Promise` object with the API response.


| Resource                                                                                                                                                                                          | Description                                                                                           |
|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| **[Admins](https://pocketbase.io/docs/api-admins)**                                                                                                                                               |                                                                                                       |
| 🔓`client.Admins.authViaEmail(email, password, bodyParams = {}, queryParams = {})`                                                                                                                 | Authenticate an admin account by its email and password and returns a new admin token and admin data. |
| 🔐`client.Admins.refresh(bodyParams = {}, queryParams = {})`                                                                                                                                       | Refreshes the current admin authenticated instance and returns a new admin token and admin data.      |
| 🔓`client.Admins.requestPasswordReset(email, bodyParams = {}, queryParams = {})`                                                                                                                   | Sends admin password reset request.                                                                   |
| 🔓`client.Admins.confirmPasswordReset(passwordResetToken, password, passwordConfirm, bodyParams = {}, queryParams = {})`                                                                           | Confirms admin password reset request.                                                                |
| 🔐`client.Admins.getList(page = 1, perPage = 30, queryParams = {})`                                                                                                                                | Returns paginated admins list.                                                                        |
| 🔐`client.Admins.getFullList(batchSize = 100, queryParams = {})`                                                                                                                                   | Returns a list with all admins batch fetched at once.                                                 |
| 🔐`client.Admins.getOne(id, queryParams = {})`                                                                                                                                                     | Returns single admin by its id.                                                                       |
| 🔐`client.Admins.create(bodyParams = {}, queryParams = {})`                                                                                                                                        | Creates a new admin.                                                                                  |
| 🔐`client.Admins.update(id, bodyParams = {}, queryParams = {})`                                                                                                                                    | Updates an existing admin by its id.                                                                  |
| 🔐`client.Admins.delete(id, bodyParams = {}, queryParams = {})`                                                                                                                                    | Deletes an existing admin by its id.                                                                  |
| **[Users](https://pocketbase.io/docs/api-users)**                                                                                                                                                 |                                                                                                       |
| 🔓`client.Users.listAuthMethods(queryParams = {})`                                                                                                                                                 | Returns all available application auth methods.                                                       |
| 🔓`client.Users.authViaEmail(email, password, bodyParams = {}, queryParams = {})`                                                                                                                  | Authenticate a user account by its email and password and returns a new user token and user data.     |
| 🔓`client.Users.authViaOAuth2(clientName, code, codeVerifier, redirectUrl, bodyParams = {}, queryParams = {})`                                                                                     | Authenticate a user via OAuth2 client provider.                                                       |
| 🔐`client.Users.refresh(bodyParams = {}, queryParams = {})`                                                                                                                                        | Refreshes the current user authenticated instance and returns a new user token and user data.         |
| 🔓`client.Users.requestPasswordReset(email, bodyParams = {}, queryParams = {})`                                                                                                                    | Sends user password reset request.                                                                    |
| 🔓`client.Users.confirmPasswordReset(passwordResetToken, password, passwordConfirm, bodyParams = {}, queryParams = {})`                                                                            | Confirms user password reset request.                                                                 |
| 🔓`client.Users.requestVerification(email, bodyParams = {}, queryParams = {})`                                                                                                                     | Sends user verification email request.                                                                |
| 🔓`client.Users.confirmVerification(verificationToken, bodyParams = {}, queryParams = {})`                                                                                                         | Confirms user email verification request.                                                             |
| 🔐`client.Users.requestEmailChange(newEmail, bodyParams = {}, queryParams = {})`                                                                                                                   | Sends an email change request to the authenticated user.                                              |
| 🔓`client.Users.confirmEmailChange(emailChangeToken, password, bodyParams = {}, queryParams = {})`                                                                                                 | Confirms user new email address.                                                                      |
| 🔐`client.Users.getList(page = 1, perPage = 30, queryParams = {})`                                                                                                                                 | Returns paginated users list.                                                                         |
| 🔐`client.Users.getFullList(batchSize = 100, queryParams = {})`                                                                                                                                    | Returns a list with all users batch fetched at once.                                                  |
| 🔐`client.Users.getOne(id, queryParams = {})`                                                                                                                                                      | Returns single user by its id.                                                                        |
| 🔐`client.Users.create(bodyParams = {}, queryParams = {})`                                                                                                                                         | Creates a new user.                                                                                   |
| 🔐`client.Users.update(id, bodyParams = {}, queryParams = {})`                                                                                                                                     | Updates an existing user by its id.                                                                   |
| 🔐`client.Users.delete(id, bodyParams = {}, queryParams = {})`                                                                                                                                     | Deletes an existing user by its id.                                                                   |
| **[Realtime](https://pocketbase.io/docs/api-realtime)** <br/> _(for node environments you'll have to install an EventSource polyfill beforehand, eg. https://github.com/EventSource/eventsource)_ |                                                                                                       |
| 🔓`client.Realtime.subscribe(subscription, callback)`                                                                                                                                              | Inits the sse connection (if not already) and register the subscription.                              |
| 🔓`client.Realtime.unsubscribe(subscription = "")`                                                                                                                                                 | Unsubscribe from a subscription (if empty - unsubscibe from all registered subscriptions).            |
| **[Collections](https://pocketbase.io/docs/api-collections)**                                                                                                                                     |                                                                                                       |
| 🔐`client.Collections.getList(page = 1, perPage = 30, queryParams = {})`                                                                                                                           | Returns paginated collections list.                                                                   |
| 🔐`client.Collections.getFullList(batchSize = 100, queryParams = {})`                                                                                                                              | Returns a list with all collections batch fetched at once.                                            |
| 🔐`client.Collections.getOne(id, queryParams = {})`                                                                                                                                                | Returns single collection by its id.                                                                  |
| 🔐`client.Collections.create(bodyParams = {}, queryParams = {})`                                                                                                                                   | Creates a new collection.                                                                             |
| 🔐`client.Collections.update(id, bodyParams = {}, queryParams = {})`                                                                                                                               | Updates an existing collection by its id.                                                             |
| 🔐`client.Collections.delete(id, bodyParams = {}, queryParams = {})`                                                                                                                               | Deletes an existing collection by its id.                                                             |
| **[Records](https://pocketbase.io/docs/api-records)**                                                                                                                                             |                                                                                                       |
| 🔓`client.Records.getList(collectionIdOrName, page = 1, perPage = 30, queryParams = {})`                                                                                                           | Returns paginated records list.                                                                       |
| 🔓`client.Records.getFullList(collectionIdOrName, batchSize = 100, queryParams = {})`                                                                                                              | Returns a list with all records batch fetched at once.                                                |
| 🔓`client.Records.getOne(collectionIdOrName, id, queryParams = {})`                                                                                                                                | Returns single record by its id.                                                                      |
| 🔓`client.Records.create(collectionIdOrName, bodyParams = {}, queryParams = {})`                                                                                                                   | Creates a new record.                                                                                 |
| 🔓`client.Records.update(collectionIdOrName, id, bodyParams = {}, queryParams = {})`                                                                                                               | Updates an existing record by its id.                                                                 |
| 🔓`client.Records.delete(collectionIdOrName, id, bodyParams = {}, queryParams = {})`                                                                                                               | Deletes an existing record by its id.                                                                 |
| **[Logs](https://pocketbase.io/docs/api-logs)**                                                                                                                                                   |                                                                                                       |
| 🔐`client.Logs.getRequestsList(page = 1, perPage = 30, queryParams = {})`                                                                                                                          | Returns paginated request logs list.                                                                  |
| 🔐`client.Logs.getRequest(id, queryParams = {})`                                                                                                                                                   | Returns a single request log by its id.                                                               |
| 🔐`client.Logs.getRequestsStats(queryParams = {})`                                                                                                                                                 | Returns request logs statistics.                                                                      |
| **[Settings](https://pocketbase.io/docs/api-records)**                                                                                                                                            |                                                                                                       |
| 🔐`client.Settings.getAll(queryParams = {})`                                                                                                                                                       | Fetch all available app settings.                                                                     |
| 🔐`client.Settings.update(bodyParams = {}, queryParams = {})`                                                                                                                                      | Bulk updates app settings.                                                                            |


## Development
```sh
# run unit tests
npm test

# build and minify for production
npm run build
```
