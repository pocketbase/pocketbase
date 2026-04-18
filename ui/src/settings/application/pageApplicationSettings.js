import { settingsSidebar } from "../settingsSidebar";
import { batchAccordion } from "./batchAccordion";
import { rateLimitAccordion, sortRules } from "./rateLimitAccordion";
import { trustedProxyAccordion } from "./trustedProxyAccordion";

export function pageApplicationSettings() {
    app.store.title = "Application settings";

    const data = store({
        isLoading: false,
        isSaving: false,
        formSettings: null,
        originalFormSettings: null,
        get originalFormSettingsHash() {
            return JSON.stringify(data.originalFormSettings);
        },
        get hasChanges() {
            return data.originalFormSettingsHash != JSON.stringify(data.formSettings);
        },
    });

    loadSettings();

    async function loadSettings() {
        data.isLoading = true;

        try {
            const settings = await app.pb.settings.getAll();
            init(settings);

            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                // data.isLoading = false; don't reset in case of a server error
            }
        }
    }

    async function save() {
        if (data.isSaving || !data.hasChanges) {
            return;
        }

        data.isSaving = true;

        data.formSettings.rateLimits.rules = sortRules(data.formSettings.rateLimits.rules);

        try {
            const redacted = app.utils.filterRedactedProps(data.formSettings);
            const settings = await app.pb.settings.update(redacted);
            init(settings);

            app.toasts.success("Successfully saved application settings.");
        } catch (err) {
            app.checkApiError(err);
        }

        data.isSaving = false;
    }

    function init(settings = {}) {
        // refresh local app settings
        app.store.settings = JSON.parse(JSON.stringify(settings));

        // load from the css style as fallback
        if (!settings.meta?.accentColor) {
            const cssColor = window.getComputedStyle(document.documentElement)?.getPropertyValue("--accentColor");
            if (cssColor?.startsWith("#")) {
                settings.meta = settings.meta || {};
                settings.meta.accentColor = cssColor.toLowerCase() || "";
            }
        }

        data.originalFormSettings = {
            meta: settings.meta || {},
            batch: settings.batch || {},
            trustedProxy: settings.trustedProxy || { headers: [] },
            rateLimits: settings.rateLimits || { rules: [] },
        };

        sortRules(data.originalFormSettings.rateLimits.rules);

        data.formSettings = JSON.parse(JSON.stringify(data.originalFormSettings));
    }

    function reset() {
        data.formSettings = JSON.parse(data.originalFormSettingsHash);
    }

    return t.div(
        {
            pbEvent: "pageApplicationSettings",
            className: "page page-application-settings",
        },
        settingsSidebar(),
        t.div(
            { className: "page-content full-height" },
            t.header(
                { className: "page-header" },
                t.nav(
                    { className: "breadcrumbs" },
                    t.div({ className: "breadcrumb-item" }, "Settings"),
                    t.div({ className: "breadcrumb-item" }, "Application"),
                ),
            ),
            t.div(
                { className: "wrapper m-b-base" },
                () => {
                    if (data.isLoading) {
                        return t.div({ className: "block txt-center" }, t.span({ className: "loader lg" }));
                    }

                    return t.form(
                        {
                            pbEvent: "applicationSettingsForm",
                            className: "grid application-settings-form",
                            inert: () => data.isSaving,
                            onsubmit: (e) => {
                                e.preventDefault();
                                save();
                            },
                        },
                        t.div(
                            { className: "col-md-5" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: "meta.appName" }, "Application name"),
                                t.input({
                                    id: "meta.appName",
                                    name: "meta.appName",
                                    type: "text",
                                    required: true,
                                    value: () => data.formSettings.meta.appName || "",
                                    oninput: (e) => (data.formSettings.meta.appName = e.target.value),
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-md-5" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: "meta.appURL" }, "Application URL"),
                                t.input({
                                    id: "meta.appURL",
                                    name: "meta.appURL",
                                    type: "url",
                                    required: true,
                                    value: () => data.formSettings.meta.appURL || "",
                                    oninput: (e) => (data.formSettings.meta.appURL = e.target.value),
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-md-2" },
                            // pass isSaving to ensure that it will be rerendered after save
                            () => accentColorField(data, data.isSaving),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            () => trustedProxyAccordion(data),
                            () => rateLimitAccordion(data),
                            () => batchAccordion(data),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            t.div(
                                { className: "field" },
                                t.input({
                                    id: "meta.hideControls",
                                    name: "meta.hideControls",
                                    type: "checkbox",
                                    className: "switch",
                                    checked: () => data.formSettings.meta.hideControls,
                                    onchange: (e) => (data.formSettings.meta.hideControls = e.target.checked),
                                }),
                                t.label(
                                    { htmlFor: "meta.hideControls" },
                                    t.span({ className: "txt" }, "Hide/Lock collection and record controls"),
                                    t.i({
                                        className: "ri-information-line link-hint",
                                        ariaDescription: app.attrs.tooltip("To prevent accidental changes when in production environment, collections create and update buttons will be hidden.\nRecords update will also require an extra unlock step before save.")
                                    }),
                                ),
                            ),
                        ),
                        t.div({ className: "col-lg-12" }, t.hr()),
                        t.div(
                            { className: "col-lg-12" },
                            t.div(
                                { className: "flex" },
                                t.div({ className: "m-r-auto" }),
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn transparent secondary",
                                        disabled: () => data.isSaving,
                                        hidden: () => !data.hasChanges,
                                        onclick: reset,
                                    },
                                    t.span({ className: "txt" }, "Cancel"),
                                ),
                                t.button(
                                    {
                                        className: () => `btn expanded ${data.isSaving ? "loading" : ""}`,
                                        disabled: () => !data.hasChanges || data.isSaving,
                                    },
                                    t.span({ className: "txt" }, "Save changes"),
                                ),
                            ),
                        ),
                    );
                },
            ),
            t.footer({ className: "page-footer" }, app.components.credits()),
        ),
    );
}

function accentColorField(pageData) {
    const uniqueId = "accent_" + app.utils.randomString();

    const local = store({
        isTooLight: false,
    });

    let colorChangeTimeoutId;
    let tempNoAnimationTimeoutId;

    function changeAccentColor(color) {
        // temporary disable animations to minimize flickering
        clearTimeout(tempNoAnimationTimeoutId);
        document.documentElement.style.setProperty("--animationSpeed", "0");

        if (color) {
            document.documentElement.style.setProperty("--accentColor", color.toLowerCase());
        } else {
            document.documentElement.style.removeProperty("--accentColor");
        }

        // restore animation
        tempNoAnimationTimeoutId = setTimeout(() => {
            document.documentElement.style.removeProperty("--animationSpeed");
        }, 100);
    }

    const watchers = [
        watch(() => pageData.formSettings?.meta?.accentColor, (newColor) => {
            clearTimeout(colorChangeTimeoutId);
            colorChangeTimeoutId = setTimeout(() => {
                changeAccentColor(newColor);
            }, 100);
        }),
    ];

    return t.div(
        {
            className: "field",
            ariaDescription: app.attrs.tooltip(() => local.isTooLight ? "Invalid - color is too light" : ""),
            onunmount: () => {
                clearTimeout(colorChangeTimeoutId);
                changeAccentColor(pageData.formSettings.meta.accentColor);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.label(
            { htmlFor: uniqueId },
            t.span({ className: "txt" }, "Accent"),
            t.i({
                hidden: () => !local.isTooLight,
                className: "txt-warning ri-alert-line",
            }),
        ),
        app.components.colorPicker({
            id: uniqueId,
            name: "meta.accentColor",
            predefinedColors: () => app.store.predefinedAccentColors,
            value: () => pageData.formSettings.meta.accentColor,
            onchange: (color) => {
                // @todo consider removing the constraint once contrast-color is implemented
                local.isTooLight = false;
                if (!app.utils.isDarkEnoughForWhiteText(color)) {
                    local.isTooLight = true;
                    return;
                }

                pageData.formSettings.meta.accentColor = color;
            },
        }),
    );
}
