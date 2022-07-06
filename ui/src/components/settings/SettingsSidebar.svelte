<script>
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
</script>

<aside class="page-sidebar settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">System</div>
        <a href="/settings" class="sidebar-list-item" use:active={{ path: "/settings" }} use:link>
            <i class="ri-home-gear-line" />
            <span class="txt">Application</span>
        </a>
        <a
            href="/settings/mail"
            class="sidebar-list-item"
            use:active={{ path: "/settings/mail/?.*" }}
            use:link
        >
            <i class="ri-send-plane-2-line" />
            <span class="txt">Mail settings</span>
        </a>
        <a
            href="/settings/storage"
            class="sidebar-list-item"
            use:active={{ path: "/settings/storage/?.*" }}
            use:link
        >
            <i class="ri-archive-drawer-line" />
            <span class="txt">Files storage</span>
        </a>

        <div class="sidebar-title">Authentication</div>
        <a
            href="/settings/auth-providers"
            class="sidebar-list-item"
            use:active={{ path: "/settings/auth-providers/?.*" }}
            use:link
        >
            <i class="ri-lock-password-line" />
            <span class="txt">Auth providers</span>
        </a>
        <a
            href="/settings/tokens"
            class="sidebar-list-item"
            use:active={{ path: "/settings/tokens/?.*" }}
            use:link
        >
            <i class="ri-key-line" />
            <span class="txt">Token options</span>
        </a>
        <a
            href="/settings/admins"
            class="sidebar-list-item"
            use:active={{ path: "/settings/admins/?.*" }}
            use:link
        >
            <i class="ri-shield-user-line" />
            <span class="txt">Admins</span>
        </a>
    </div>
</aside>
