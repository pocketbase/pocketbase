import PocketBase, { isTokenExpired, LocalAuthStore } from "pocketbase";

const LOGIN_PATH = "#/login";

const currentPath = window.location.pathname.endsWith("/")
    ? window.location.pathname.substring(0, window.location.pathname.length - 1)
    : window.location.pathname;

window.app = window.app || {};
window.app.pb = new PocketBase(
    import.meta.env.PB_BACKEND_URL,
    // concatenate the path in case hosted under subpath alongside other apps
    new LocalAuthStore("__pb_superusers__" + currentPath),
);

// add UI specific header to all requests
app.pb.beforeSend = function(url, options) {
    options.headers["x-request-source"] = "pbui";
    return { url, options };
};

app.store.superuser = app.pb.authStore.record;
app.pb.authStore.onChange((_, record) => {
    if (!record && window.location.hash != LOGIN_PATH) {
        app.modals.close();
        window.location.hash = LOGIN_PATH;
    }

    app.store.superuser = record;
});

// refresh the token in the background
if (app.pb.authStore.isValid) {
    app.pb.collection(app.pb.authStore.record?.collectionName || "_superusers")
        .authRefresh()
        .catch((err) => {
            console.warn("Failed to refresh the existing auth token:", err);

            // clear the store only on invalidated/expired token
            const status = err?.status << 0;
            if (status == 401 || status == 403) {
                app.utils.rememberPath();
                app.pb.cancelAllRequests();
                app.pb.authStore.clear();
            }
        });
}

// load initial store data
app.pb.authStore.onChange((_, record) => {
    if (record?.id) {
        app.store.loadCollections();
        app.store.loadSettings();
        app.store.loadOAuth2Providers();
    }
});

// Modify the default RecordService to fire global events on record
// create, update, delete without relying on the realtime service
// -------------------------------------------------------------------

const originalRecordService = app.pb.collection;
app.pb.collection = function(idOrName) {
    const service = originalRecordService.call(this, idOrName);

    bindRecordServiceEvents(service);

    return service;
};

function bindRecordServiceEvents(service) {
    if (service.__customUIEvents) {
        return;
    }

    service.__customUIEvents = true;

    const originalCreate = service.create;
    service.create = function() {
        return originalCreate.apply(service, arguments).then((r) => {
            setTimeout(() => {
                document.dispatchEvent(new CustomEvent("record:create", { detail: r }));
                document.dispatchEvent(new CustomEvent("record:save", { detail: r }));
            }, 0);
            return r;
        });
    };

    const originalUpdate = service.update;
    service.update = function() {
        return originalUpdate.apply(service, arguments).then((r) => {
            setTimeout(() => {
                document.dispatchEvent(new CustomEvent("record:update", { detail: r }));
                document.dispatchEvent(new CustomEvent("record:save", { detail: r }));
            }, 0);
            return r;
        });
    };

    const originalDelete = service.delete;
    service.delete = function() {
        return originalDelete.apply(service, arguments).then((r) => {
            const minimalRecord = {
                id: arguments[0],
                collectionId: service.collectionIdOrName,
                collectionName: service.collectionIdOrName,
            };
            setTimeout(() => {
                document.dispatchEvent(new CustomEvent("record:delete", { detail: minimalRecord }));
            }, 0);

            return r;
        });
    };
}

// File token helpers
// -------------------------------------------------------------------

const LAST_FILE_TOKEN_KEY = "pbLastFileToken";
let isFileTokenLoading = false;
let fileTokenPromises = [];

// clear stored token on logout
app.pb.authStore.onChange((_, record) => {
    if (!record?.id) {
        window.localStorage.removeItem(LAST_FILE_TOKEN_KEY);
    }
});

/**
 * Return a superuser file token.
 * Optionally you can provide a collection identifier, to avoid unnecessery
 * calls in case the collection doesn't have protected files.
 *
 * @param  {String} optCollectionIdORName
 * @return {Promise<string>}
 */
window.app.getFileToken = async function(optCollectionIdORName = "") {
    // check if the collection needs a file token
    const collection = optCollectionIdORName
        && app.store.collections?.find((c) => c.id == optCollectionIdORName || c.name == optCollectionIdORName);
    if (collection) {
        const hasProtectedFile = collection.fields?.find((f) => f.type == "file" && f.protected);
        if (!hasProtectedFile) {
            return;
        }
    }

    let token = window.localStorage.getItem(LAST_FILE_TOKEN_KEY);

    if (!token || isTokenExpired(token, 60)) {
        token = await fetchFileToken();
    }

    return token;
};

async function fetchFileToken() {
    return new Promise(async (resolve, reject) => {
        fileTokenPromises.push({ resolve, reject });

        if (isFileTokenLoading) {
            return;
        }

        isFileTokenLoading = true;

        try {
            const token = await app.pb.files.getToken();

            window.localStorage.setItem(LAST_FILE_TOKEN_KEY, token);

            fileTokenPromises.forEach((p) => p.resolve(token));
        } catch (err) {
            fileTokenPromises.forEach((p) => p.reject(err));
        }

        isFileTokenLoading = false;
        fileTokenPromises = [];
    });
}

// Generic API error handler
// -------------------------------------------------------------------

/**
 * Helper to parse a response error and to show an optional toast message.
 * In case of 401 it clears the auth store and redirects to the home page.
 * In case of 403 it redirects to the home or login page.
 *
 * Example:
 *
 * ```js
 * try {
 *     await app.pb.collection("example").getFullList()
 * } catch (err) {
 *     if (!err?.isAbort) {
 *         app.checkApiError(err)
 *     }
 * }
 * ```
 *
 * @param {Error}  err
 * @param {boolean} showToast
 */
window.app.checkApiError = function(err, showToast = true) {
    if (!err || !(err instanceof Error) || err.isAbort) {
        console.warn("checkApiError - unexpected error type:", err);
        return;
    }

    const statusCode = err?.status << 0;
    const response = err?.response || {};

    // add toast error notification
    let msg = showToast && (response.message || err.message || "Something went wrong!");
    if (msg) {
        app.toasts.error(msg);
    }

    // unknown client-side error
    if (statusCode == 0) {
        console.log(err);
    }

    // populate form field errors
    if (!app.utils.isEmpty(response.data)) {
        app.store.errors = response.data;
    }

    // unauthorized
    if (statusCode === 401 && window.location.hash != LOGIN_PATH) {
        app.utils.rememberPath();
        app.pb.cancelAllRequests();
        return app.pb.authStore.clear();
    }

    // forbidden
    if (statusCode === 403) {
        app.pb.cancelAllRequests();
        if (window.location.hash != LOGIN_PATH) {
            window.location.hash = LOGIN_PATH;
        }
    }
};
