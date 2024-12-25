<script>
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import { hideControls } from "@/stores/app";
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
</script>

<PageSidebar class="settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">System</div>
        <a href="/settings" class="sidebar-list-item" use:active={{ path: "/settings" }} use:link>
            <i class="ri-home-gear-line" aria-hidden="true" />
            <span class="txt">Application</span>
        </a>
        <a
            href="/settings/mail"
            class="sidebar-list-item"
            use:active={{ path: "/settings/mail/?.*" }}
            use:link
        >
            <i class="ri-send-plane-2-line" aria-hidden="true" />
            <span class="txt">Mail settings</span>
        </a>
        <a
            href="/settings/storage"
            class="sidebar-list-item"
            use:active={{ path: "/settings/storage/?.*" }}
            use:link
        >
            <i class="ri-archive-drawer-line" aria-hidden="true" />
            <span class="txt">Files storage</span>
        </a>
        <a
            href="/settings/backups"
            class="sidebar-list-item"
            use:active={{ path: "/settings/backups/?.*" }}
            use:link
        >
            <i class="ri-archive-line" aria-hidden="true" />
            <span class="txt">Backups</span>
        </a>
        <a
            href="/settings/crons"
            class="sidebar-list-item"
            use:active={{ path: "/settings/crons/?.*" }}
            use:link
        >
            <i class="ri-time-line" aria-hidden="true"></i>
            <span class="txt">Crons</span>
        </a>

        {#if !$hideControls}
            <div class="sidebar-title">
                <span class="txt">Sync</span>
            </div>
            <a
                href="/settings/export-collections"
                class="sidebar-list-item"
                use:active={{ path: "/settings/export-collections/?.*" }}
                use:link
            >
                <i class="ri-uninstall-line" aria-hidden="true" />
                <span class="txt">Export collections</span>
            </a>
            <a
                href="/settings/import-collections"
                class="sidebar-list-item"
                use:active={{ path: "/settings/import-collections/?.*" }}
                use:link
            >
                <i class="ri-install-line" aria-hidden="true" />
                <span class="txt">Import collections</span>
            </a>
        {/if}
    </div>
</PageSidebar>
