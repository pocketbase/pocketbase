import { settingsSidebar } from "../settingsSidebar";

export function pageMailSettings(route) {
    app.store.title = "Mail settings";

    const tlsOptions = [
        { label: "Auto (StartTLS)", value: false },
        { label: "Always", value: true },
    ];

    const authMethods = [
        { label: "PLAIN (default)", value: "PLAIN" },
        { label: "LOGIN", value: "LOGIN" },
    ];

    const data = store({
        isLoading: false,
        isSaving: false,
        formSettings: null,
        initSerialized: "null",
        showMoreOptions: false,
        get hasChanges() {
            return data.initSerialized != JSON.stringify(data.formSettings);
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

        try {
            const redacted = app.utils.filterRedactedProps(data.formSettings);
            const settings = await app.pb.settings.update(redacted);
            init(settings);

            app.toasts.success("Successfully saved mail settings.");
        } catch (err) {
            app.checkApiError(err);
        }

        data.isSaving = false;
    }

    function init(settings = {}) {
        // refresh local app settings
        app.store.settings = JSON.parse(JSON.stringify(settings));

        data.formSettings = {
            meta: settings?.meta || {},
            smtp: settings?.smtp || {},
        };

        if (!data.formSettings.smtp.authMethod) {
            data.formSettings.smtp.authMethod = authMethods[0].value;
        }

        data.initSerialized = JSON.stringify(data.formSettings);
    }

    function reset() {
        data.formSettings = JSON.parse(data.initSerialized);
    }

    return t.div(
        { pbEvent: "pageMailSettings", className: "page page-mail-settings" },
        settingsSidebar(),
        t.div(
            { className: "page-content full-height" },
            t.header(
                { className: "page-header" },
                t.nav(
                    { className: "breadcrumbs" },
                    t.div({ className: "breadcrumb-item" }, "Settings"),
                    t.div({ className: "breadcrumb-item" }, () => app.store.title),
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
                            pbEvent: "mailSettingsForm",
                            className: "grid mail-settings-form",
                            inert: () => data.isSaving,
                            onsubmit: (e) => {
                                e.preventDefault();
                                save();
                            },
                        },
                        t.div(
                            { className: "col-lg-12 txt-lg" },
                            t.p(null, "Configure common settings for sending emails."),
                        ),
                        t.div(
                            { className: "col-lg-6" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: "meta.senderName" }, "Sender name"),
                                t.input({
                                    id: "meta.senderName",
                                    name: "meta.senderName",
                                    type: "text",
                                    required: true,
                                    value: () => data.formSettings.meta.senderName || "",
                                    oninput: (e) => (data.formSettings.meta.senderName = e.target.value),
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-lg-6" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: "meta.senderAddress" }, "Sender address"),
                                t.input({
                                    id: "meta.senderAddress",
                                    name: "meta.senderAddress",
                                    type: "email",
                                    required: true,
                                    value: () => data.formSettings.meta.senderAddress || "",
                                    oninput: (e) => (data.formSettings.meta.senderAddress = e.target.value),
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            t.div(
                                { className: "field" },
                                t.input({
                                    id: "smtp.enabled",
                                    name: "smtp.enabled",
                                    type: "checkbox",
                                    className: "switch",
                                    checked: () => !!data.formSettings.smtp.enabled,
                                    onchange: (e) => (data.formSettings.smtp.enabled = e.target.checked),
                                }),
                                t.label(
                                    { htmlFor: "smtp.enabled" },
                                    t.span(
                                        { className: "txt" },
                                        "Use SMTP mail server ",
                                        t.strong(null, "(recommended)"),
                                    ),
                                    t.i({
                                        className: "ri-information-line link-faded",
                                        ariaDescription: app.attrs.tooltip(
                                            `By default PocketBase uses the unix "sendmail" command for sending emails. For better emails deliverability it is recommended to use a SMTP mail server.`,
                                        ),
                                    }),
                                ),
                            ),
                            // SMTP
                            app.components.slide(
                                () => data.formSettings.smtp.enabled,
                                t.div(
                                    { className: "grid m-t-sm" },
                                    t.div(
                                        { className: "col-lg-4" },
                                        t.div(
                                            { className: "field" },
                                            t.label({ htmlFor: "smtp.host" }, "SMTP server host"),
                                            t.input({
                                                id: "smtp.host",
                                                name: "smtp.host",
                                                type: "text",
                                                required: () => data.formSettings.smtp.enabled,
                                                value: () => data.formSettings.smtp.host || "",
                                                oninput: (e) => data.formSettings.smtp.host = e.target.value,
                                            }),
                                        ),
                                    ),
                                    t.div(
                                        { className: "col-lg-2" },
                                        t.div(
                                            { className: "field" },
                                            t.label({ htmlFor: "smtp.port" }, "Port"),
                                            t.input({
                                                id: "smtp.port",
                                                name: "smtp.port",
                                                type: "number",
                                                min: 0,
                                                step: 1,
                                                required: () => data.formSettings.smtp.enabled,
                                                value: () => data.formSettings.smtp.port || "",
                                                oninput: (e) =>
                                                    data.formSettings.smtp.port = parseInt(e.target.value, 10),
                                            }),
                                        ),
                                    ),
                                    t.div(
                                        { className: "col-lg-3" },
                                        t.div(
                                            { className: "field" },
                                            t.label({ htmlFor: "smtp.username" }, "Username"),
                                            t.input({
                                                id: "smtp.username",
                                                name: "smtp.username",
                                                type: "text",
                                                autocomplete: "off",
                                                value: () => data.formSettings.smtp.username || "",
                                                oninput: (e) => data.formSettings.smtp.username = e.target.value,
                                            }),
                                        ),
                                    ),
                                    t.div(
                                        { className: "col-lg-3" },
                                        t.div(
                                            { className: "field" },
                                            t.label({ htmlFor: "smtp.password" }, "Password"),
                                            t.input({
                                                id: "smtp.password",
                                                name: "smtp.password",
                                                type: "password",
                                                autocomplete: "new-password",
                                                value: () => data.formSettings.smtp.password || "",
                                                oninput: (e) => data.formSettings.smtp.password = e.target.value,
                                                onkeyup: (e) => {
                                                    if (
                                                        e.key == "Backspace"
                                                        && typeof data.formSettings.smtp.password === "undefined"
                                                    ) {
                                                        data.formSettings.smtp.password = "";
                                                    }
                                                },
                                                placeholder: () =>
                                                    typeof data.formSettings.smtp.password !== "undefined"
                                                        ? ""
                                                        : "* * * * * *",
                                            }),
                                        ),
                                    ),
                                ),
                                // additional options
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn secondary sm m-t-sm",
                                        onclick: () => data.showMoreOptions = !data.showMoreOptions,
                                    },
                                    t.span(
                                        { className: "txt" },
                                        () => data.showMoreOptions ? "Hide more options" : "Show more options",
                                    ),
                                    t.i({
                                        className: () =>
                                            data.showMoreOptions ? "ri-arrow-drop-up-line" : "ri-arrow-drop-down-line",
                                    }),
                                ),
                                app.components.slide(
                                    () => data.showMoreOptions,
                                    t.div(
                                        { className: "grid m-t-sm" },
                                        t.div(
                                            { className: "col-lg-3" },
                                            t.div(
                                                { className: "field" },
                                                t.label({ htmlFor: "smtp.tls" }, "TLS encryption"),
                                                app.components.select({
                                                    id: "smtp.tls",
                                                    name: "smtp.tls",
                                                    required: true,
                                                    options: tlsOptions,
                                                    value: () => data.formSettings.smtp.tls || false,
                                                    onchange: (selected) => {
                                                        data.formSettings.smtp.tls = selected?.[0]?.value;
                                                    },
                                                }),
                                            ),
                                        ),
                                        t.div(
                                            { className: "col-lg-3" },
                                            t.div(
                                                { className: "field" },
                                                t.label({ htmlFor: "smtp.authMethod" }, "AUTH method"),
                                                app.components.select({
                                                    id: "smtp.authMethod",
                                                    name: "smtp.authMethod",
                                                    required: true,
                                                    options: authMethods,
                                                    value: () =>
                                                        data.formSettings.smtp.authMethod || authMethods[0].value,
                                                    onchange: (selected) => {
                                                        data.formSettings.smtp.authMethod = selected?.[0]?.value;
                                                    },
                                                }),
                                            ),
                                        ),
                                        t.div(
                                            { className: "col-lg-6" },
                                            t.div(
                                                { className: "field" },
                                                t.label(
                                                    { htmlFor: "smtp.localName" },
                                                    t.span({ className: "txt" }, "EHLO/HELO domain"),
                                                    t.i({
                                                        className: "ri-information-line link-hint tooltip-top",
                                                        ariaDescription: app.attrs.tooltip(
                                                            "Some SMTP servers, such as the Gmail SMTP-relay, requires a proper domain name in the inital EHLO/HELO exchange and will reject attempts to use localhost.",
                                                        ),
                                                    }),
                                                ),
                                                t.input({
                                                    id: "smtp.localName",
                                                    name: "smtp.localName",
                                                    type: "text",
                                                    placeholder: "Default to localhost",
                                                    value: () => data.formSettings.smtp.localName || "",
                                                    oninput: (e) => data.formSettings.smtp.localName = e.target.value,
                                                }),
                                            ),
                                        ),
                                    ),
                                ),
                            ),
                        ),
                        t.div({ className: "col-lg-12" }, t.hr()),
                        t.div(
                            { className: "col-lg-12" },
                            t.div(
                                { className: "flex" },
                                t.div({ className: "m-r-auto" }),
                                () => {
                                    if (data.hasChanges) {
                                        return [
                                            t.button(
                                                {
                                                    type: "button",
                                                    className: "btn transparent secondary",
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
                                        ];
                                    }

                                    return t.button(
                                        {
                                            type: "button",
                                            className: () => `btn expanded outline`,
                                            onclick: () => app.modals.openMailTest(),
                                        },
                                        t.i({ className: "ri-mail-check-line", ariaHidden: true }),
                                        t.span({ className: "txt" }, "Send test email"),
                                    );
                                },
                            ),
                        ),
                    );
                },
            ),
            t.footer({ className: "page-footer" }, app.components.credits()),
        ),
    );
}
