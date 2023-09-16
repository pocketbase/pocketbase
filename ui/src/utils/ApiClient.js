import PocketBase, { LocalAuthStore, isTokenExpired } from "pocketbase";
// ---
import CommonHelper                     from "@/utils/CommonHelper";
import { replace }                      from "svelte-spa-router";
import { get }                          from "svelte/store";
import { addErrorToast }                from "@/stores/toasts";
import { setErrors }                    from "@/stores/errors";
import { setAdmin }                     from "@/stores/admin";
import { protectedFilesCollectionsCache } from "@/stores/collections";

const adminFileTokenKey = "pb_admin_file_token";

/**
 * Clears the authorized state and redirects to the login page.
 *
 * @param {Boolean} [redirect] Whether to redirect to the login page.
 */
PocketBase.prototype.logout = function(redirect = true) {
    this.authStore.clear();

    if (redirect) {
        replace("/login");
    }
};

/**
 * Generic API error response handler.
 *
 * @param  {Error}   err        The API error itself.
 * @param  {Boolean} notify     Whether to add a toast notification.
 * @param  {String}  defaultMsg Default toast notification message if the error doesn't have one.
 */
PocketBase.prototype.error = function(err, notify = true, defaultMsg = "") {
    if (!err || !(err instanceof Error) || err.isAbort) {
        return;
    }

    const statusCode = (err?.status << 0) || 400;
    const responseData = err?.data || {};
    const msg = responseData.message || err.message || defaultMsg;

    // add toast error notification
    if (notify && msg) {
        addErrorToast(msg);
    }

    // populate form field errors
    if (!CommonHelper.isEmpty(responseData.data)) {
        setErrors(responseData.data);
    }

    // unauthorized
    if (statusCode === 401) {
        this.cancelAllRequests();
        return this.logout();
    }

    // forbidden
    if (statusCode === 403) {
        this.cancelAllRequests();
        return replace("/");
    }
};

/**
 * @return {Promise<String>}
 */
PocketBase.prototype.getAdminFileToken = async function(collectionId = "") {
    let needToken = true;

    if (collectionId) {
        const protectedCollections = get(protectedFilesCollectionsCache);
        needToken = typeof protectedCollections[collectionId] !== "undefined"
            ? protectedCollections[collectionId]
            : true;
    }

    if (!needToken) {
        return "";
    }

    let token = localStorage.getItem(adminFileTokenKey) || "";

    // request a new token only if the previous one is missing or will expire soon
    if (!token || isTokenExpired(token, 10)) {
        // remove previously stored token (if any)
        token && localStorage.removeItem(adminFileTokenKey);

        if (!this._adminFileTokenRequest) {
            this._adminFileTokenRequest = this.files.getToken();
        }

        token = await this._adminFileTokenRequest;
        localStorage.setItem(adminFileTokenKey, token);
        this._adminFileTokenRequest = null;
    }

    return token;
}

// Custom auth store to sync the svelte admin store state with the authorized admin instance.
class AppAuthStore extends LocalAuthStore {
    /**
     * @inheritdoc
     */
    save(token, model) {
        super.save(token, model);

        if (model && !model.collectionId) { // not an auth record
            setAdmin(model);
        }
    }

    /**
     * @inheritdoc
     */
    clear() {
        super.clear();

        setAdmin(null);
    }
}

const pb = new PocketBase(
    import.meta.env.PB_BACKEND_URL,
    new AppAuthStore("pb_admin_auth")
);

if (pb.authStore.model && !pb.authStore.model.collectionId) { // not an auth record
    setAdmin(pb.authStore.model);
}

export default pb;
