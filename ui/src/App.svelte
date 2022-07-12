<script>
    import "./scss/main.scss";

    import Router, { replace, link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
    import routes from "./routes";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Toasts from "@/components/base/Toasts.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import Confirmation from "@/components/base/Confirmation.svelte";
    import { admin } from "@/stores/admin";
    import { setErrors } from "@/stores/errors";
    import { resetConfirmation } from "@/stores/confirmation";
    import { _, setupI18n, locale, dir } from '@/services/i18n';
    import LocaleSwitcher from '@/components/base/LocaleSwitcher.svelte';
    $: { document.dir = $dir; }

    let oldLocation = undefined;

    let showAppSidebar = false;

    function handleRouteLoading(e) {
        if (e?.detail?.location === oldLocation) {
            return; // not an actual change
        }

        showAppSidebar = !!e?.detail?.userData?.showAppSidebar;

        oldLocation = e?.detail?.location;

        // resets
        CommonHelper.setDocumentTitle("");
        setErrors({});
        resetConfirmation();
    }

    function handleRouteFailure() {
        replace("/");
    }

    function logout() {
        ApiClient.logout();
    }
</script>

<div class="app-layout">
    {#if $admin?.id && showAppSidebar}
        <aside class="app-sidebar">
            <a href="/" class="logo logo-sm" use:link>
                <img
                    src="{import.meta.env.BASE_URL}images/logo.svg"
                    alt="{$_('app.menu.logo')}"
                    width="40"
                    height="40"
                />
            </a>
            <LocaleSwitcher value={$locale} on:locale-changed={e => setupI18n({ withLocale: e.detail })} />

            <nav class="main-menu">
                <a
                    href="/collections"
                    class="menu-item"
                    aria-label="{$_('app.menu.collections')}"
                    use:link
                    use:active={{ path: "/collections/?.*", className: "current-route" }}
                    use:tooltip={{ text: $_('app.menu.collections'), position: "right" }}
                >
                    <i class="ri-database-2-line" />
                </a>
                <a
                    href="/users"
                    class="menu-item"
                    aria-label="{$_('app.menu.users')}"
                    use:link
                    use:active={{ path: "/users/?.*", className: "current-route" }}
                    use:tooltip={{ text: $_('app.menu.users'), position: "right" }}
                >
                    <i class="ri-group-line" />
                </a>
                <a
                    href="/logs"
                    class="menu-item"
                    aria-label="{$_('app.menu.logs')}"
                    use:link
                    use:active={{ path: "/logs/?.*", className: "current-route" }}
                    use:tooltip={{ text: $_('app.menu.logs'), position: "right" }}
                >
                    <i class="ri-line-chart-line" />
                </a>
                <a
                    href="/settings"
                    class="menu-item"
                    aria-label="{$_('app.menu.settings')}"
                    use:link
                    use:active={{ path: "/settings/?.*", className: "current-route" }}
                    use:tooltip={{ text: $_('app.menu.settings'), position: "right" }}
                >
                    <i class="ri-tools-line" />
                </a>
            </nav>

            <figure class="thumb thumb-circle link-hint closable">
                <img
                    src="{import.meta.env.BASE_URL}images/avatars/avatar{$admin?.avatar || 0}.svg"
                    alt="Avatar"
                />
                <Toggler class="dropdown dropdown-nowrap dropdown-upside dropdown-left">
                    <a href="/settings/admins" class="dropdown-item closable" use:link>
                        <i class="ri-shield-user-line" />
                        <span class="txt">{$_('app.menu.admins')}</span>
                    </a>
                    <hr />
                    <div tabindex="0" class="dropdown-item closable" on:click={logout}>
                        <i class="ri-logout-circle-line" />
                        <span class="txt">{$_('app.menu.logout')}</span>
                    </div>
                </Toggler>
            </figure>
        </aside>
    {/if}

    <div class="app-body">
        <Router {routes} on:routeLoading={handleRouteLoading} on:conditionsFailed={handleRouteFailure} />
    </div>
</div>

<Toasts />

<Confirmation />
