import { pageSuperuserLogin } from "@/auth/pageSuperuserLogin";
import { pageCollections } from "@/collections/pageCollections";
import { pageLogs } from "@/logs/pageLogs";
import { pageApplicationSettings } from "@/settings/application/pageApplicationSettings";
import { pageBackupsSettings } from "@/settings/backups/pageBackupsSettings";
import { pageCronsSettings } from "@/settings/crons/pageCronsSettings";
import { pageMailSettings } from "@/settings/mail/pageMailSettings";
import { pageStorageSettings } from "@/settings/storage/pageStorageSettings";
import { pageExportCollections } from "@/settings/sync/pageExportCollections";
import { pageImportCollections } from "@/settings/sync/pageImportCollections";

window.app = window.app || {};
window.app.routes = window.app.routes || {};

window.app.routes.fallbackPath = "#/collections";

/**
 * @callback RouteHandler
 * @param {Object} route
 * @param {string} route.path
 * @param {Object} route.query
 * @param {Object} route.params
 * @param {RegExp} route.regex
 * @param {string} route.pattern
 * @return {HtmlElement|Promise<HtmlElement>} The page to render.
 */

/**
 * Registers a new guest-only route.
 * If the user is authenticated, they are redirected to the home page.
 *
 * Example:
 *
 * ```js
 * app.routes.guestOnly("#/guest", (route) => {
 *     return t.div({ className: "page" }, "Guest page")
 * })
 * ```
 *
 * @param {string} path
 * @param {RouteHandler} handler
 */
app.routes.guestOnly = function(path, handler) {
    if (app.store._ready) {
        throw new Error("the router is already initialized");
    }

    routeDefs[path] = async (route) => {
        if (app.pb.authStore.isValid && app.pb.authStore.record?.id) {
            window.location.hash = "#/";
            return;
        }

        app.store.showHeader = false;
        app.store.page = await handler(route);
    };
};

/**
 * Registers a new superuser-only route (the layout has a top-nav header).
 * If the user is not authenticated, they are redirected to the login page.
 *
 * Example:
 *
 * ```js
 * app.routes.superuserOnly("#/example", (route) => {
 *     return t.div({ className: "page" }, "Superuser accessible page")
 * })
 * ```
 *
 * @param {string} path
 * @param {RouteHandler} handler
 */
app.routes.superuserOnly = function(path, handler) {
    if (app.store._ready) {
        throw new Error("the router is already initialized");
    }

    routeDefs[path] = async (route) => {
        if (!app.pb.authStore.isValid || !app.pb.authStore.record?.id) {
            window.location.hash = "#/login";
            return;
        }

        app.store.showHeader = true;
        app.store.page = await handler(route);
    };
};

/**
 * Registers a new route with a blank layout.
 * This route layout doesn't perform any access check to allow more advanced customizations.
 *
 * Example:
 *
 * ```js
 * app.routes.blank("#/blank", (route) => {
 *     return t.div({ className: "page" }, "Blank layout page")
 * })
 * ```
 *
 * @param {string} path
 * @param {RouteHandler} handler
 */
app.routes.blank = function(path, handler) {
    if (app.store._ready) {
        throw new Error("the router is already initialized");
    }

    routeDefs[path] = async (route) => {
        app.store.showHeader = false;
        app.store.page = await handler(route);
    };
};

const routeDefs = {};
let destroyRouter;
export function initRouter() {
    if (destroyRouter) {
        destroyRouter();
    }

    destroyRouter = router(routeDefs, { fallbackPath: app.routes.fallbackPath });
}

// -------------------------------------------------------------------

app.routes.guestOnly("#/pbinstall/{token}", async (route) => {
    const { pageInstaller } = await import("@/auth/pageInstaller");
    return pageInstaller(route);
});
app.routes.guestOnly("#/login", (route) => {
    return pageSuperuserLogin(route);
});
app.routes.guestOnly("#/request-password-reset", async (route) => {
    const { pageRequestSuperuserPasswordReset } = await import("@/auth/pageRequestSuperuserPasswordReset");
    return pageRequestSuperuserPasswordReset(route);
});

// email confirmation actions
app.routes.blank("#/auth/confirm-password-reset/{token}", async (route) => {
    const { pageConfirmPasswordReset } = await import("@/auth/pageConfirmPasswordReset");
    return pageConfirmPasswordReset(route);
});
app.routes.blank("#/auth/confirm-verification/{token}", async (route) => {
    const { pageConfirmVerification } = await import("@/auth/pageConfirmVerification");
    return pageConfirmVerification(route);
});
app.routes.blank("#/auth/confirm-email-change/{token}", async (route) => {
    const { pageConfirmEmailChange } = await import("@/auth/pageConfirmEmailChange");
    return pageConfirmEmailChange(route);
});

// oauth2 redirect pages
app.routes.blank("#/auth/oauth2-redirect-success", async (route) => {
    const { pageOAuth2RedirectSuccess } = await import("@/auth/pageOAuth2RedirectSuccess");
    return pageOAuth2RedirectSuccess(route);
});
app.routes.blank("#/auth/oauth2-redirect-failure", async (route) => {
    const { pageOAuth2RedirectFailure } = await import("@/auth/pageOAuth2RedirectFailure");
    return pageOAuth2RedirectFailure(route);
});

app.routes.superuserOnly("#/collections", pageCollections);
app.routes.superuserOnly("#/logs", pageLogs);
app.routes.superuserOnly("#/settings", pageApplicationSettings);
app.routes.superuserOnly("#/settings/mail", pageMailSettings);
app.routes.superuserOnly("#/settings/storage", pageStorageSettings);
app.routes.superuserOnly("#/settings/backups", pageBackupsSettings);
app.routes.superuserOnly("#/settings/crons", pageCronsSettings);
app.routes.superuserOnly("#/settings/export-collections", pageExportCollections);
app.routes.superuserOnly("#/settings/import-collections", pageImportCollections);
