window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openProviderSettings = function(
    config = {},
    settings = {
        namePrefix: "", // e.g. 'oauth2.providers.0'
        // ---
        onbeforeopen: null,
        onafteropen: null,
        onbeforeclose: null,
        onafterclose: null,
        onsubmit: null, // (providerInfo, newConfig) => {},
    },
) {
    const modal = providerSettingsModal(config, settings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);
    app.modals.open(modal);
};

function providerSettingsModal(providerConfig, settings) {
    let modal;

    const uniqueId = "provider_" + app.utils.randomString();

    providerConfig = providerConfig || {};

    const isNew = !providerConfig.clientId;

    const initialHash = JSON.stringify(providerConfig);

    const providerInfo = app.store.oauth2Providers?.find((p) => p.name == providerConfig.name);
    if (!providerInfo) {
        console.warn("missing provider for config", providerConfig);
        return;
    }

    const data = store({
        config: JSON.parse(initialHash),
        get hasChanges() {
            return initialHash != JSON.stringify(data.config);
        },
        onsubmit: (providerInfo, newConfig) => {},
    });

    function submit() {
        if (!data.hasChanges) {
            return;
        }

        settings.onsubmit?.(providerInfo, JSON.parse(JSON.stringify(data.config)));

        app.modals.close(modal);
    }

    modal = t.div(
        {
            pbEvent: "providerSettingsModal",
            className: "modal provider-settings-modal",
            onbeforeopen: (el) => {
                return settings.onbeforeopen?.(el);
            },
            onafteropen: (el) => {
                settings.onafteropen?.(el);

                // retriger errors (if any)
                setTimeout(() => {
                    if (app.store.errors?.oauth2) {
                        app.store.errors.oauth2 = JSON.parse(JSON.stringify(app.store.errors.oauth2));
                    }
                }, 0);
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
            t.figure(
                { className: "provider-logo" },
                () => {
                    if (providerInfo.logo) {
                        return t.img({
                            src: "data:image/svg+xml;base64," + btoa(providerInfo.logo),
                            alt: providerInfo.name + " logo",
                        });
                    }

                    return t.i({ className: app.utils.fallbackProviderIcon });
                },
            ),
            t.h6(
                { className: "modal-title" },
                providerConfig.displayName || providerInfo.displayName || providerInfo.name,
                t.small({ className: "txt-hint" }, " (", providerConfig.name, ")"),
            ),
        ),
        t.form(
            {
                pbEvent: "providerSettingsForm",
                id: uniqueId + "form",
                className: "modal-content",
                onsubmit: (e) => {
                    e.preventDefault();
                    submit();
                },
            },
            t.div(
                { className: "grid" },
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "field" },
                        t.label({
                            htmlFor: uniqueId + ".clientId",
                            textContent: "Client ID",
                        }),
                        t.input({
                            type: "text",
                            required: true,
                            id: uniqueId + ".clientId",
                            autocomplete: "off",
                            name: () => settings.namePrefix + ".clientId",
                            value: () => data.config.clientId || "",
                            oninput: (e) => (data.config.clientId = e.target.value.trim()),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "field" },
                        t.label({
                            htmlFor: uniqueId + ".clientSecret",
                            textContent: "Client secret",
                        }),
                        t.input({
                            type: "password",
                            id: uniqueId + ".clientSecret",
                            autocomplete: "new-password",
                            required: () => isNew || typeof data.config.clientSecret != "undefined",
                            name: () => settings.namePrefix + ".clientSecret",
                            value: () => data.config.clientSecret || "",
                            oninput: (e) => (data.config.clientSecret = e.target.value.trim()),
                            onkeyup: (e) => {
                                if (
                                    e.key == "Backspace"
                                    && typeof data.config.clientSecret === "undefined"
                                ) {
                                    data.config.clientSecret = "";
                                }
                            },
                            placeholder:
                                () => (isNew || typeof data.config.clientSecret !== "undefined" ? "" : "* * * * * *"),
                        }),
                    ),
                ),
                // extra fields
                () => {
                    if (typeof app.oauth2?.[providerInfo.name] == "function") {
                        return t.div(
                            { className: "col-12" },
                            app.oauth2[providerInfo.name](providerInfo, settings.namePrefix, data),
                        );
                    }
                },
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => app.modals.close(modal),
                },
                t.span({ className: "txt" }, "Close"),
            ),
            t.button(
                {
                    "html-form": uniqueId + "form",
                    type: "submit",
                    className: "btn",
                    disabled: () => !data.hasChanges,
                },
                t.span({ className: "txt" }, "Set provider config"),
            ),
        ),
    );

    return modal;
}
