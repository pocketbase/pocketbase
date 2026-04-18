const excludedFieldNames = ["id", "email", "emailVisibility", "verified", "tokenKey", "password"];
const allowedRegularTypes = ["text", "editor", "url", "email", "json"];

export function oauth2Accordion(collection) {
    const uniqueId = "oauth2_" + app.utils.randomString();

    const data = store({
        get config() {
            if (!collection.oauth2) {
                collection.oauth2 = {
                    enabled: false,
                    mappedFields: {},
                    providers: [],
                };
            }

            return collection.oauth2;
        },
        get regularFieldOptions() {
            return collection.fields
                ?.filter((f) => {
                    return allowedRegularTypes.includes(f.type) && !excludedFieldNames.includes(f.name);
                })
                .map((f) => {
                    return { value: f.name };
                });
        },
        get regularAndFileFieldOptions() {
            return collection.fields
                ?.filter((f) => {
                    return (
                        (f.type == "file" || allowedRegularTypes.includes(f.type))
                        && !excludedFieldNames.includes(f.name)
                    );
                })
                .map((f) => {
                    return { value: f.name };
                });
        },
        showMapping: false,
    });

    function clearProviderErrors(index) {
        app.utils.deleteByPath(app.store.errors, "oauth2.providers." + index);
    }

    return t.details(
        {
            pbEvent: "oauth2Accordion",
            name: "auth-methods",
            className: "accordion oauth2-accordion",
        },
        t.summary(
            null,
            t.i({ className: "ri-profile-line" }),
            t.span({ className: "txt", textContent: "OAuth2" }),
            t.span({
                className: () => `label m-l-auto ${data.config.enabled ? "success" : ""}`,
                textContent: () => (data.config.enabled ? "Enabled" : "Disabled"),
            }),
            () => {
                if (app.utils.isEmpty(app.store.errors?.oauth2)) {
                    return;
                }

                return t.i({
                    className: "ri-error-warning-fill txt-danger",
                    ariaDescription: app.attrs.tooltip("Has errors", "left"),
                });
            },
        ),
        t.div(
            { className: "grid sm" },
            t.div(
                { className: "col-sm-12" },
                t.div(
                    { className: "field" },
                    t.input({
                        type: "checkbox",
                        id: uniqueId + ".enabled",
                        name: "oauth2.enabled",
                        className: "switch",
                        checked: () => data.config.enabled,
                        onchange: (e) => (data.config.enabled = e.target.checked),
                    }),
                    t.label({
                        htmlFor: uniqueId + ".enabled",
                        textContent: "Enable",
                    }),
                ),
            ),
            () => {
                return data.config.providers.map((providerConfig, configIndex) => {
                    const providerId = uniqueId + providerConfig.name;
                    const providerInfo = app.store.oauth2Providers?.find((p) => p.name == providerConfig.name) || {};

                    return t.div(
                        { className: "col-sm-6" },
                        t.div(
                            {
                                className: () => {
                                    let result = "provider-card";

                                    if (!app.utils.isEmpty(app.store.errors?.oauth2?.providers?.[configIndex])) {
                                        result += " error";
                                    }

                                    return result;
                                },
                            },
                            t.figure(
                                { className: "provider-logo" },
                                () => {
                                    if (providerInfo.logo) {
                                        return t.img({
                                            src: "data:image/svg+xml;base64," + btoa(providerInfo.logo),
                                            alt: providerConfig.name + " logo",
                                        });
                                    }

                                    return t.i({ className: app.utils.fallbackProviderIcon });
                                },
                            ),
                            t.div(
                                { className: "content" },
                                t.span(
                                    { className: "primary-txt" },
                                    () => providerConfig.displayName || providerInfo.displayName || providerInfo.name,
                                ),
                                t.span({ className: "secondary-txt" }, () => providerConfig.name || providerInfo.name),
                            ),
                            t.div(
                                { className: "actions" },
                                t.button(
                                    {
                                        "type": "button",
                                        "className": "btn secondary transparent sm circle",
                                        "html-popovertarget": providerId + "dropdown",
                                    },
                                    t.i({ className: "ri-more-2-line" }),
                                ),
                                t.div(
                                    {
                                        id: providerId + "dropdown",
                                        className: "dropdown sm",
                                        popover: "auto",
                                    },
                                    t.button(
                                        {
                                            type: "button",
                                            className: "dropdown-item",
                                            onclick: (e) => {
                                                e.target.closest(".dropdown").hidePopover();
                                                app.modals.openProviderSettings(providerConfig, {
                                                    namePrefix: "oauth2.providers." + configIndex,
                                                    onsubmit: (providerInfo, providerConfig) => {
                                                        data.config.providers[configIndex] = providerConfig;
                                                        clearProviderErrors(configIndex);
                                                    },
                                                });
                                            },
                                        },
                                        t.span({ className: "txt" }, "Settings"),
                                    ),
                                    t.hr(),
                                    t.button(
                                        {
                                            type: "button",
                                            className: "dropdown-item",
                                            onclick: (e) => {
                                                e.target.closest(".dropdown").hidePopover();
                                                app.modals.confirm(
                                                    `Do you really want to remove provider "${
                                                        providerConfig.displayName || providerInfo.displayName
                                                        || providerInfo.name
                                                    }"?`,
                                                    () => {
                                                        clearProviderErrors(configIndex);

                                                        data.config.providers.splice(configIndex, 1);

                                                        if (data.config.providers.length == 0) {
                                                            data.config.enabled = false;
                                                        }
                                                    },
                                                );
                                            },
                                        },
                                        t.span({ className: "txt" }, "Remove"),
                                    ),
                                ),
                            ),
                        ),
                    );
                });
            },
            t.div(
                { className: "col-sm-6" },
                t.button(
                    {
                        type: "button",
                        className: "btn lg block secondary add-provider-btn",
                        onclick: () => {
                            app.modals.openProviderPicker({
                                exclude: data.config.providers.map((p) => p.name),
                                onselect: (providerInfo) => {
                                    app.modals.openProviderSettings({ name: providerInfo.name }, {
                                        onsubmit: (providerInfo, providerConfig) => {
                                            if (data.config.providers.length == 0) {
                                                data.config.enabled = true;
                                            }

                                            data.config.providers.push(providerConfig);
                                        },
                                    });
                                },
                            });
                        },
                    },
                    t.i({ className: "ri-add-line" }),
                    t.span({ className: "txt " }, "Add provider"),
                ),
            ),
            t.div(
                { className: "col-sm-12" },
                t.button(
                    {
                        type: "button",
                        className: () => `btn secondary sm ${data.showMapping ? "" : "transparent"}`,
                        onclick: () => (data.showMapping = !data.showMapping),
                    },
                    t.span({ className: "txt" }, "Optional users create fields mapping"),
                    t.i({
                        className: () => (data.showMapping ? "ri-arrow-drop-up-line" : "ri-arrow-drop-down-line"),
                    }),
                ),
                app.components.slide(
                    () => data.showMapping,
                    t.div(
                        { className: "grid sm m-t-sm" },
                        t.div(
                            { className: "col-sm-6" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: uniqueId + ".mappedFields.name" }, "OAuth2 full name"),
                                app.components.select({
                                    id: uniqueId + ".mappedFields.name",
                                    name: "oauth2.mappedFields.name",
                                    placeholder: "Select field",
                                    options: () => data.regularFieldOptions,
                                    value: () => collection.oauth2.mappedFields.name,
                                    onchange: (selectedOpts) => {
                                        collection.oauth2.mappedFields.name = selectedOpts?.[0]?.value || "";
                                    },
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-sm-6" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: uniqueId + ".mappedFields.avatarURL" }, "OAuth2 avatar"),
                                app.components.select({
                                    id: uniqueId + ".mappedFields.avatarURL",
                                    name: "oauth2.mappedFields.avatarURL",
                                    placeholder: "Select field",
                                    options: () => data.regularAndFileFieldOptions,
                                    value: () => collection.oauth2.mappedFields.avatarURL,
                                    onchange: (selectedOpts) => {
                                        collection.oauth2.mappedFields.avatarURL = selectedOpts?.[0]?.value || "";
                                    },
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-sm-6" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: uniqueId + ".mappedFields.id" }, "OAuth2 id"),
                                app.components.select({
                                    id: uniqueId + ".mappedFields.id",
                                    name: "oauth2.mappedFields.id",
                                    placeholder: "Select field",
                                    options: () => data.regularFieldOptions,
                                    value: () => collection.oauth2.mappedFields.id,
                                    onchange: (selectedOpts) => {
                                        collection.oauth2.mappedFields.id = selectedOpts?.[0]?.value || "";
                                    },
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-sm-6" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: uniqueId + ".mappedFields.username" }, "OAuth2 username"),
                                app.components.select({
                                    id: uniqueId + ".mappedFields.username",
                                    name: "oauth2.mappedFields.username",
                                    placeholder: "Select field",
                                    options: () => data.regularFieldOptions,
                                    value: () => collection.oauth2.mappedFields.username,
                                    onchange: (selectedOpts) => {
                                        collection.oauth2.mappedFields.username = selectedOpts?.[0]?.value || "";
                                    },
                                }),
                            ),
                        ),
                    ),
                ),
            ),
        ),
    );
}
