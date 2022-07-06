import { wrap }          from "svelte-spa-router/wrap";
import ApiClient         from "@/utils/ApiClient";
import PageLogs          from "@/components/logs/PageLogs.svelte";
import PageRecords       from "@/components/records/PageRecords.svelte";
import PageUsers         from "@/components/users/PageUsers.svelte";
import PageAdmins        from "@/components/admins/PageAdmins.svelte";
import PageAdminLogin    from "@/components/admins/PageAdminLogin.svelte";
import PageApplication   from "@/components/settings/PageApplication.svelte";
import PageMail          from "@/components/settings/PageMail.svelte";
import PageStorage       from "@/components/settings/PageStorage.svelte";
import PageAuthProviders from "@/components/settings/PageAuthProviders.svelte";
import PageTokenOptions  from "@/components/settings/PageTokenOptions.svelte";

const routes = {
    "/_elements": wrap({
        asyncComponent: () => import("@/components/Elements.svelte"),
    }),

    "/login": wrap({
        component:  PageAdminLogin,
        conditions: [(_) => !ApiClient.AuthStore.isValid],
    }),

    "/request-password-reset": wrap({
        asyncComponent:  () => import("@/components/admins/PageAdminRequestPasswordReset.svelte"),
        conditions: [(_) => !ApiClient.AuthStore.isValid],
    }),

    "/confirm-password-reset/:token": wrap({
        asyncComponent:  () => import("@/components/admins/PageAdminConfirmPasswordReset.svelte"),
        conditions: [(_) => !ApiClient.AuthStore.isValid],
    }),

    "/collections": wrap({
        component:  PageRecords,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/logs": wrap({
        component: PageLogs,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/users": wrap({
        component:  PageUsers,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/users/confirm-password-reset/:token": wrap({
        asyncComponent:  () => import("@/components/users/PageUserConfirmPasswordReset.svelte"),
        conditions: [
            () => {
                // ensure that there is no authenticated user/admin model
                ApiClient.logout(false);
                return true;
            },
        ],
    }),

    "/users/confirm-verification/:token": wrap({
        asyncComponent:  () => import("@/components/users/PageUserConfirmVerification.svelte"),
        conditions: [
            () => {
                // ensure that there is no authenticated user/admin model
                ApiClient.logout(false);
                return true;
            },
        ],
    }),

    "/users/confirm-email-change/:token": wrap({
        asyncComponent:  () => import("@/components/users/PageUserConfirmEmailChange.svelte"),
        conditions: [
            () => {
                // ensure that there is no authenticated user/admin model
                ApiClient.logout(false);
                return true;
            },
        ],
    }),

    "/settings": wrap({
        component:  PageApplication,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/settings/admins": wrap({
        component:  PageAdmins,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/settings/mail": wrap({
        component:  PageMail,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/settings/storage": wrap({
        component:  PageStorage,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/settings/auth-providers": wrap({
        component:  PageAuthProviders,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    "/settings/tokens": wrap({
        component:  PageTokenOptions,
        conditions: [(_) => ApiClient.AuthStore.isValid],
    }),

    // fallback
    "*": wrap({
        asyncComponent: () => import("@/components/NotFoundPage.svelte"),
    }),
};

export default routes;
