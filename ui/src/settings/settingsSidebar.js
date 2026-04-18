export function settingsSidebar() {
    return app.components.pageSidebar(
        {
            pbEvent: "settingsSidebar",
            className: "settings-sidebar",
        },
        t.nav(
            { className: "sidebar-content scrollable" },
            () => {
                const result = [];

                for (const groupName in app.store.settingsNavGroups) {
                    const children = app.store.settingsNavGroups[groupName];

                    const groupEl = t.details(
                        { className: "nav-group", "html-data-group": groupName, open: true },
                        t.summary(
                            { tabIndex: -1, onfocusout: () => false, onclick: () => false, onkeyup: () => false },
                            groupName,
                        ),
                        () => {
                            return children.map((link) => {
                                const isLocal = link.href.startsWith("#/");

                                return t.a(
                                    {
                                        href: () => link.href,
                                        target: () => !isLocal ? "_blank" : undefined,
                                        rel: () => !isLocal ? "noopener noreferrer" : undefined,
                                        className: (el) => {
                                            const isActive = link.isActive?.(el)
                                                || app.utils.isActivePath(link.href, false);
                                            return `nav-item ${isActive ? "active" : ""}`;
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
                    );

                    result.push(groupEl);
                }

                return result;
            },
        ),
    );
}
