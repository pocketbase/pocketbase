window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openProviderPicker = function(settings = {
    exclude: [],
    // ---
    onbeforeopen: null,
    onafteropen: null,
    onbeforeclose: null,
    onafterclose: null,
    onselect: null, // (providerInfo) => {},
}) {
    const modal = providerPickerModal(settings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);
    app.modals.open(modal);
};

function providerPickerModal(settings = {}) {
    let modal;

    const data = store({
        searchTerm: "",
        get filteredProviders() {
            const search = data.searchTerm.trim().toLowerCase().replaceAll(" ", "");

            return app.store.oauth2Providers.filter((p) => {
                return (
                    !settings.exclude?.includes(p.name)
                    && (p.name + p.displayName).toLowerCase().replaceAll(" ", "").includes(search)
                );
            });
        },
    });

    function clearSearch() {
        data.searchTerm = "";
    }

    modal = t.div(
        {
            pbEvent: "providerPickerModal",
            className: "modal provider-picker-modal",
            onbeforeopen: (el) => {
                return settings.onbeforeopen?.(el);
            },
            onafteropen: (el) => {
                settings.onafteropen?.(el);
            },
            onbeforeclose: (el) => {
                return settings.onbeforeclose?.(el);
            },
            onafterclose: (el) => {
                settings.onafterclose?.(el);
                el?.remove();
            },
        },
        t.header(
            { className: "modal-header" },
            t.h6({ className: "modal-title" }, t.span({ className: "txt" }, "Select OAuth2 provider")),
        ),
        t.div(
            { className: "modal-content" },
            t.div(
                { className: "grid sm" },
                // search
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "fields searchbar" },
                        t.div(
                            { className: "field" },
                            t.input({
                                placeholder: "Search...",
                                className: "p-l-20",
                                value: () => data.searchTerm,
                                oninput: (e) => data.searchTerm = e.target.value,
                            }),
                        ),
                        () => {
                            if (!data.searchTerm) {
                                return;
                            }

                            return t.div(
                                { rid: "search-ctrls", className: "field addon p-r-5" },
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn sm pill secondary transparent",
                                        onclick: () => clearSearch(),
                                    },
                                    "Clear",
                                ),
                            );
                        },
                    ),
                ),
                // no providers
                () => {
                    if (app.store.isLoadingOAuth2Providers || data.filteredProviders.length) {
                        return;
                    }

                    return t.div(
                        { rid: "notfound", className: "block txt-center txt-hint" },
                        t.p(null, "No providers found."),
                        t.button({
                            type: "button",
                            className: "btn sm secondary",
                            textContent: "Clear search",
                            onclick: () => clearSearch(),
                        }),
                    );
                },
                // list
                () => {
                    if (app.store.isLoadingOAuth2Providers) {
                        return t.div({ className: "col-12 txt-center" }, t.span({ className: "loader active" }));
                    }

                    return data.filteredProviders.map((provider) => {
                        return t.div(
                            { className: "col-sm-6" },
                            t.button(
                                {
                                    type: "button",
                                    className: "provider-card handle",
                                    onclick: () => {
                                        app.modals.close(modal);
                                        settings.onselect?.(provider);
                                    },
                                },
                                t.figure(
                                    { className: "provider-logo" },
                                    () => {
                                        if (provider.logo) {
                                            return t.img({
                                                src: "data:image/svg+xml;base64," + btoa(provider.logo),
                                                alt: provider.name + " logo",
                                            });
                                        }

                                        return t.i({ className: app.utils.fallbackProviderIcon, ariaHidden: true });
                                    },
                                ),
                                t.div(
                                    { className: "content" },
                                    t.span({ className: "primadry-txt" }, provider.displayName || provider.name),
                                    t.span({ className: "secondary-txt" }, provider.name),
                                ),
                            ),
                        );
                    });
                },
            ),
        ),
        t.footer(
            { className: "modal-footer gap-base" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => app.modals.close(modal),
                },
                t.span({ className: "txt" }, "Close"),
            ),
        ),
    );

    return modal;
}
