import PageIndex from "@/components/PageIndex.svelte";
import PageLogs from "@/components/logs/PageLogs.svelte";
import PageRecords from "@/components/records/PageRecords.svelte";
import PageApplication from "@/components/settings/PageApplication.svelte";
import PageBackups from "@/components/settings/PageBackups.svelte";
import PageCrons from "@/components/settings/PageCrons.svelte";
import PageExportCollections from "@/components/settings/PageExportCollections.svelte";
import PageImportCollections from "@/components/settings/PageImportCollections.svelte";
import PageMail from "@/components/settings/PageMail.svelte";
import PageStorage from "@/components/settings/PageStorage.svelte";
import PageSuperuserLogin from "@/components/superusers/PageSuperuserLogin.svelte";
import ApiClient from "@/utils/ApiClient";
import { isTokenExpired } from "pocketbase";
import { wrap } from "svelte-spa-router/wrap";

const routes = {
    "/pbinstal/:token": wrap({
        asyncComponent: () => import("@/components/base/PageInstaller.svelte"),
        conditions: [(details) => {
            return details.params.token && !isTokenExpired(details.params.token)
        }],
        userData: { showAppSidebar: false },
    }),

    "/login": wrap({
        component: PageSuperuserLogin,
        conditions: [(_) => !ApiClient.authStore.isValid],
        userData: { showAppSidebar: false },
    }),

    "/request-password-reset": wrap({
        asyncComponent: () => import("@/components/superusers/PageSuperuserRequestPasswordReset.svelte"),
        conditions: [(_) => !ApiClient.authStore.isValid],
        userData: { showAppSidebar: false },
    }),

    "/confirm-password-reset/:token": wrap({
        asyncComponent: () => import("@/components/superusers/PageSuperuserConfirmPasswordReset.svelte"),
        conditions: [(_) => !ApiClient.authStore.isValid],
        userData: { showAppSidebar: false },
    }),

    "/collections": wrap({
        component: PageRecords,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/logs": wrap({
        component: PageLogs,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings": wrap({
        component: PageApplication,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings/mail": wrap({
        component: PageMail,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings/storage": wrap({
        component: PageStorage,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings/export-collections": wrap({
        component: PageExportCollections,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings/import-collections": wrap({
        component: PageImportCollections,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings/backups": wrap({
        component: PageBackups,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    "/settings/crons": wrap({
        component: PageCrons,
        conditions: [(_) => ApiClient.authStore.isValid],
        userData: { showAppSidebar: true },
    }),

    // ---------------------------------------------------------------
    // Records email confirmation actions
    // ---------------------------------------------------------------

    // @deprecated
    "/users/confirm-password-reset/:token": wrap({
        asyncComponent: () => import("@/components/records/PageRecordConfirmPasswordReset.svelte"),
        userData: { showAppSidebar: false },
    }),
    "/auth/confirm-password-reset/:token": wrap({
        asyncComponent: () => import("@/components/records/PageRecordConfirmPasswordReset.svelte"),
        userData: { showAppSidebar: false },
    }),

    // @deprecated
    "/users/confirm-verification/:token": wrap({
        asyncComponent: () => import("@/components/records/PageRecordConfirmVerification.svelte"),
        userData: { showAppSidebar: false },
    }),
    "/auth/confirm-verification/:token": wrap({
        asyncComponent: () => import("@/components/records/PageRecordConfirmVerification.svelte"),
        userData: { showAppSidebar: false },
    }),

    // @deprecated
    "/users/confirm-email-change/:token": wrap({
        asyncComponent: () => import("@/components/records/PageRecordConfirmEmailChange.svelte"),
        userData: { showAppSidebar: false },
    }),
    "/auth/confirm-email-change/:token": wrap({
        asyncComponent: () => import("@/components/records/PageRecordConfirmEmailChange.svelte"),
        userData: { showAppSidebar: false },
    }),

    "/auth/oauth2-redirect-success": wrap({
        asyncComponent: () => import("@/components/records/PageOAuth2RedirectSuccess.svelte"),
        userData: { showAppSidebar: false },
    }),

    "/auth/oauth2-redirect-failure": wrap({
        asyncComponent: () => import("@/components/records/PageOAuth2RedirectFailure.svelte"),
        userData: { showAppSidebar: false },
    }),

    // catch-all fallback
    "*": wrap({
        component: PageIndex,
        userData: { showAppSidebar: false },
    }),
};

export default routes;
