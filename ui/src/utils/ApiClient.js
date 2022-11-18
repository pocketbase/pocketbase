import PocketBase, { LocalAuthStore, Admin } from "pocketbase";
// ---
import CommonHelper      from "@/utils/CommonHelper";
import { replace }       from "svelte-spa-router";
import { addErrorToast } from "@/stores/toasts";
import { setErrors }     from "@/stores/errors";
import { setAdmin }      from "@/stores/admin";

/**
 * Clears the authorized state and redirects to the login page.
 *
 * @param {Boolean} [redirect] Whether to redirect to the login page.
 */
PocketBase.prototype.logout = function(redirect = true) {
    this.authStore.clear();

    if (redirect) {
        replace('/login');
    }
};

/**
 * Generic API error response handler.
 *
 * @param  {Error}   err        The API error itself.
 * @param  {Boolean} notify     Whether to add a toast notification.
 * @param  {String}  defaultMsg Default toast notification message if the error doesn't have one.
 */
PocketBase.prototype.errorResponseHandler = function(err, notify = true, defaultMsg = '') {
    if (!err || !(err instanceof Error) || err.isAbort) {
        return;
    }

    const statusCode = (err?.status << 0) || 400;
    const responseData = err?.data || {};

    // add toast error notification
    if (
        notify &&          // notifications are enabled
        statusCode !== 404 // is not 404
    ) {
        let msg = responseData.message || err.message || defaultMsg;
        if (msg) {
            addErrorToast(msg);
        }
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
        return replace('/');
    }
};

// Custom auth store to sync the svelte admin store state with the authorized admin instance.
class AppAuthStore extends LocalAuthStore {
    /**
     * @inheritdoc
     */
    save(token, model) {
        super.save(token, model);

        if (model instanceof Admin) {
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

const client = new PocketBase(
    import.meta.env.PB_BACKEND_URL,
    new AppAuthStore("pb_admin_auth")
);

if (client.authStore.model instanceof Admin) {
    setAdmin(client.authStore.model);
}

export default client;
