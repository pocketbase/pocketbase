<script>
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
    import { _ } from '@/services/i18n';
</script>

<aside class="page-sidebar settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">{$_('settings.sidebar.system')}</div>
        <a href="/settings" class="sidebar-list-item" use:active={{ path: "/settings" }} use:link>
            <i class="ri-home-gear-line" />
            <span class="txt">{$_('settings.sidebar.application')}</span>
        </a>
        <a
            href="/settings/mail"
            class="sidebar-list-item"
            use:active={{ path: "/settings/mail/?.*" }}
            use:link
        >
            <i class="ri-send-plane-2-line" />
            <span class="txt">{$_('settings.sidebar.mail_setting')}</span>
        </a>
        <a
            href="/settings/storage"
            class="sidebar-list-item"
            use:active={{ path: "/settings/storage/?.*" }}
            use:link
        >
            <i class="ri-archive-drawer-line" />
            <span class="txt">{$_('settings.sidebar.file_storage')}</span>
        </a>

        <div class="sidebar-title">{$_('settings.sidebar.authentication')}</div>
        <a
            href="/settings/auth-providers"
            class="sidebar-list-item"
            use:active={{ path: "/settings/auth-providers/?.*" }}
            use:link
        >
            <i class="ri-lock-password-line" />
            <span class="txt">{$_('settings.sidebar.auth_provider')}</span>
        </a>
        <a
            href="/settings/tokens"
            class="sidebar-list-item"
            use:active={{ path: "/settings/tokens/?.*" }}
            use:link
        >
            <i class="ri-key-line" />
            <span class="txt">{$_('settings.sidebar.token_options')}</span>
        </a>
        <a
            href="/settings/admins"
            class="sidebar-list-item"
            use:active={{ path: "/settings/admins/?.*" }}
            use:link
        >
            <i class="ri-shield-user-line" />
            <span class="txt">{$_('settings.sidebar.admins')}</span>
        </a>
    </div>
</aside>
