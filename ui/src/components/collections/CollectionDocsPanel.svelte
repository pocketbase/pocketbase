<script>
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const baseTabs = {
        list: {
            label: "List/Search",
            component: import("@/components/collections/docs/ListApiDocs.svelte"),
        },
        view: {
            label: "View",
            component: import("@/components/collections/docs/ViewApiDocs.svelte"),
        },
        create: {
            label: "Create",
            component: import("@/components/collections/docs/CreateApiDocs.svelte"),
        },
        update: {
            label: "Update",
            component: import("@/components/collections/docs/UpdateApiDocs.svelte"),
        },
        delete: {
            label: "Delete",
            component: import("@/components/collections/docs/DeleteApiDocs.svelte"),
        },
        realtime: {
            label: "Realtime",
            component: import("@/components/collections/docs/RealtimeApiDocs.svelte"),
        },
    };

    const authTabs = {
        "auth-with-password": {
            label: "Auth with password",
            component: import("@/components/collections/docs/AuthWithPasswordDocs.svelte"),
        },
        "auth-with-oauth2": {
            label: "Auth with OAuth2",
            component: import("@/components/collections/docs/AuthWithOAuth2Docs.svelte"),
        },
        refresh: {
            label: "Auth refresh",
            component: import("@/components/collections/docs/AuthRefreshDocs.svelte"),
        },
        "request-verification": {
            label: "Request verification",
            component: import("@/components/collections/docs/RequestVerificationDocs.svelte"),
        },
        "confirm-verification": {
            label: "Confirm verification",
            component: import("@/components/collections/docs/ConfirmVerificationDocs.svelte"),
        },
        "request-password-reset": {
            label: "Request password reset",
            component: import("@/components/collections/docs/RequestPasswordResetDocs.svelte"),
        },
        "confirm-password-reset": {
            label: "Confirm password reset",
            component: import("@/components/collections/docs/ConfirmPasswordResetDocs.svelte"),
        },
        "request-email-change": {
            label: "Request email change",
            component: import("@/components/collections/docs/RequestEmailChangeDocs.svelte"),
        },
        "confirm-email-change": {
            label: "Confirm email change",
            component: import("@/components/collections/docs/ConfirmEmailChangeDocs.svelte"),
        },
        "list-auth-methods": {
            label: "List auth methods",
            component: import("@/components/collections/docs/AuthMethodsDocs.svelte"),
        },
        "list-linked-accounts": {
            label: "List OAuth2 accounts",
            component: import("@/components/collections/docs/ListExternalAuthsDocs.svelte"),
        },
        "unlink-account": {
            label: "Unlink OAuth2 account",
            component: import("@/components/collections/docs/UnlinkExternalAuthDocs.svelte"),
        },
    };

    let docsPanel;
    let collection = {};
    let activeTab;
    let tabs = [];

    $: if (collection.type === "auth") {
        tabs = Object.assign({}, baseTabs, authTabs);
        if (!collection.options.allowUsernameAuth && !collection.options.allowEmailAuth) {
            delete tabs["auth-with-password"];
        }
        if (!collection.options.allowOAuth2Auth) {
            delete tabs["auth-with-oauth2"];
        }
    } else if (collection.type === "view") {
        tabs = Object.assign({}, baseTabs);
        delete tabs.create;
        delete tabs.update;
        delete tabs.delete;
        delete tabs.realtime;
    } else {
        tabs = Object.assign({}, baseTabs);
    }

    // reset active tab on tabs list change
    if (tabs.length) {
        activeTab = Object.keys(tabs)[0];
    }

    export function show(model) {
        collection = model;

        changeTab(Object.keys(tabs)[0]);

        return docsPanel?.show();
    }

    export function hide() {
        return docsPanel?.hide();
    }

    export function changeTab(newTab) {
        activeTab = newTab;
    }
</script>

<OverlayPanel bind:this={docsPanel} on:hide on:show class="docs-panel">
    <div class="docs-content-wrapper">
        <aside class="docs-sidebar" class:compact={collection?.type === "auth"}>
            <nav class="sidebar-content">
                {#each Object.entries(tabs) as [key, tab], i (key)}
                    <!-- add a separator before the first auth tab -->
                    {#if i === Object.keys(baseTabs).length}
                        <hr class="m-t-sm m-b-sm" />
                    {/if}

                    <button
                        type="button"
                        class="sidebar-item"
                        class:active={activeTab === key}
                        on:click={() => changeTab(key)}
                    >
                        {tab.label}
                    </button>
                {/each}
            </nav>
        </aside>

        <div class="docs-content">
            {#each Object.entries(tabs) as [key, tab] (key)}
                {#if activeTab === key}
                    {#await tab.component then { default: TabComponent }}
                        <TabComponent {collection} />
                    {/await}
                {/if}
            {/each}
        </div>
    </div>

    <!-- visible only on small screens -->
    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Close</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
