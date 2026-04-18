window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Returns a new element with common helper links, usually shown in the footer.
 *
 * @return {Element}
 */
window.app.components.credits = function() {
    return t.div(
        { pbEvent: "credits", className: "credits" },
        () => {
            return app.store.creditLinks.map((link) => {
                const isLocal = link.href.startsWith("#/");

                return t.a(
                    {
                        href: () => link.href,
                        target: () => !isLocal ? "_blank" : undefined,
                        rel: () => !isLocal ? "noopener noreferrer" : undefined,
                        className: (el) => {
                            const isActive = link.isActive?.(el) || app.utils.isActivePath(link.href, false);
                            return `credit-item ${isActive ? "active" : ""}`;
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
};
