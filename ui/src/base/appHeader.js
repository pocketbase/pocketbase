export function appHeader() {
    return () => {
        if (!app.store._ready || !app.store.showHeader || !app.store.superuser?.id) {
            return;
        }

        return t.header(
            {
                pbEvent: "appHeader",
                rid: "appHeader",
                className: "app-header accent-surface",
                onmount: async (el) => {
                    await new Promise((r) => setTimeout(r, 0));

                    el._scrollToActiveMenuItem = function() {
                        el?.querySelector(".app-main-nav .header-link.active")?.scrollIntoView();
                    };
                    el._scrollToActiveMenuItem();
                    window.addEventListener("hashchange", el._scrollToActiveMenuItem);
                },
                onunmount: (el) => {
                    window.removeEventListener("hashchange", el?._scrollToActiveMenuItem);
                },
            },
            t.a(
                { href: "#/", className: "logo" },
                t.img({ src: () => app.store.headerLogo, alt: "App logo" }),
            ),
            t.nav(
                {
                    pbEvent: "mainNav",
                    className: "app-main-nav",
                },
                () => {
                    return app.store.headerLinks.map((link) => {
                        const isLocal = link.href.startsWith("#/");

                        return t.a(
                            {
                                href: () => link.href,
                                target: () => !isLocal ? "_blank" : undefined,
                                rel: () => !isLocal ? "noopener noreferrer" : undefined,
                                className: (el) => {
                                    const isActive = link.isActive?.(el) || app.utils.isActivePath(link.href);
                                    return `header-link ${isActive ? "active" : ""}`;
                                },
                            },
                            () => {
                                if (link.icon) {
                                    return t.i({ className: link.icon });
                                }
                            },
                            t.span({ className: "txt" }, () => link.label),
                        );
                    });
                },
            ),
            t.div({ className: "flex-fill app-header-separator" }),
            colorSchemeButton(),
            t.button(
                {
                    className: "header-link logged-user txt-normal",
                    "html-popovertarget": "logged-user-dropdown",
                },
                t.span({ className: "superuser-name txt-ellipsis" }, () => app.store.superuser?.email),
                t.i({ className: "ri-arrow-drop-down-line" }),
            ),
            t.div(
                {
                    pbEvent: "loggedUserDropdown",
                    id: "logged-user-dropdown",
                    className: "dropdown sm nowrap logged-user-dropdown",
                    popover: "auto",
                },
                t.a(
                    {
                        className: "dropdown-item dropdown-item-manage",
                        href: "#/collections?collection=_superusers",
                        onclick: (e) => {
                            e.target.closest(".dropdown").hidePopover();
                        },
                    },
                    t.i({ className: "ri-group-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Manage superusers"),
                ),
                t.hr(),
                t.button(
                    {
                        type: "button",
                        className: "dropdown-item txt-danger dropdown-item-logout",
                        onclick: () => app.pb.authStore.clear(),
                    },
                    t.i({ className: "ri-logout-circle-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Logout"),
                ),
            ),
        );
    };
}

function colorSchemeButton() {
    const options = [
        { value: "light", icon: "ri-sun-line", label: "Light" },
        { value: "dark", icon: "ri-moon-line", label: "Dark" },
        { value: "", icon: "ri-subtract-line", label: "Auto" },
    ];

    return [
        t.button(
            {
                className: "header-link color-scheme-picker",
                "html-popovertarget": "color-scheme-dropdown",
                title: "Color scheme",
            },
            t.i({
                className: () => app.store.activeColorScheme == "dark" ? "ri-moon-line" : "ri-sun-line",
                ariaHidden: true,
            }),
        ),
        t.div(
            {
                pbEvent: "colorSchemeDropdown",
                id: "color-scheme-dropdown",
                className: "dropdown sm nowrap color-scheme-dropdown",
                popover: "auto",
            },
            () => {
                return options.map((opt) => {
                    return t.button(
                        {
                            type: "button",
                            className: () =>
                                `dropdown-item dropdown-item-light ${
                                    app.store.userColorScheme == opt.value ? "active" : ""
                                }`,
                            onclick: (e) => {
                                e.target.closest(".dropdown").hidePopover();
                                app.store.userColorScheme = opt.value;
                            },
                        },
                        t.i({ className: opt.icon, ariaHidden: true }),
                        t.span({ className: "txt" }, opt.label),
                    );
                });
            },
        ),
    ];
}
