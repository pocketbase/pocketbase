export function backupsForm(propsArg = {}) {
    const props = store({
        onsave: null,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const presets = [
        { cron: "0 0 * * *", label: "Every day at 00:00h" },
        { cron: "0 0 * * 0", label: "Every sunday at 00:00h" },
        { cron: "0 0 * * 1,3", label: "Every Mon and Wed at 00:00h" },
        { cron: "0 0 1 * *", label: "Every first day of the month at 00:00h" },
    ];

    const data = store({
        showForm: false,
        isLoading: false,
        isSaving: false,
        formSettings: null,
        initSerialized: "null",
        enableAutoBackups: false,
        get hasChanges() {
            return data.initSerialized != JSON.stringify(data.formSettings);
        },
    });

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

            props.onsave?.(settings);

            init(settings);

            app.toasts.success("Successfully saved backups settings.");
        } catch (err) {
            app.checkApiError(err);
        }

        data.isSaving = false;
    }

    function init(settings = {}) {
        // refresh local app settings
        app.store.settings = JSON.parse(JSON.stringify(settings));

        data.formSettings = {
            backups: settings?.backups || {},
        };

        data.enableAutoBackups = !!data.formSettings.backups.cron;
        data.initSerialized = JSON.stringify(data.formSettings);
    }

    function reset() {
        data.formSettings = JSON.parse(data.initSerialized);
        data.enableAutoBackups = !!data.formSettings.backups.cron;
    }

    watchers.push(
        watch(() => {
            if (!data.enableAutoBackups && data.formSettings?.backups?.cron) {
                data.formSettings.backups.cron = "";
            }
        }),
    );

    return t.div(
        {
            className: "block backups-settings-form-wrapper",
            onmount: () => {
                loadSettings();
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.button(
            {
                type: "button",
                className: () => `btn secondary ${data.isLoading ? "loading" : ""}`,
                disabled: () => data.isLoading || data.hasChanges,
                onclick: () => (data.showForm = !data.showForm),
            },
            t.span({ className: "txt" }, "Backup options"),
            t.i({
                className: () => (data.showForm ? "ri-arrow-up-s-line" : "ri-arrow-down-s-line"),
                ariaHidden: true,
            }),
        ),
        app.components.slide(
            () => data.showForm,
            t.form(
                {
                    pbEvent: "backupsSettingsForm",
                    className: "grid backups-settings-form m-t-base",
                    inert: () => data.isSaving,
                    onsubmit: (e) => {
                        e.preventDefault();
                        save();
                    },
                },
                () => {
                    if (data.isLoading) {
                        return t.div({ className: "col-lg-12 txt-center" }, t.span({ className: "loader lg" }));
                    }

                    return [
                        t.div(
                            { className: "col-lg-12" },
                            t.div(
                                { className: "field" },
                                t.input({
                                    id: "enableAutoBackupsToggle",
                                    type: "checkbox",
                                    className: "switch",
                                    checked: () => data.enableAutoBackups,
                                    onchange: (e) => {
                                        data.enableAutoBackups = e.target.checked;
                                        if (!data.formSettings.backups.cron) {
                                            data.formSettings.backups.cron = presets[0].cron;
                                        }
                                    },
                                }),
                                t.label({ htmlFor: "enableAutoBackupsToggle" }, "Enable auto backups"),
                            ),
                            app.components.slide(
                                () => data.enableAutoBackups,
                                t.div(
                                    { className: "grid m-t-base m-b-base" },
                                    t.div(
                                        { className: "col-lg-6" },
                                        t.div(
                                            { className: "fields" },
                                            t.div(
                                                { className: "field" },
                                                t.label({ htmlFor: "backups.cron" }, "Cron expression"),
                                                t.input({
                                                    id: "backups.cron",
                                                    name: "backups.cron",
                                                    className: "txt-code",
                                                    type: "text",
                                                    placeholder: "e.g. 0 0 * * *",
                                                    required: () => data.enableAutoBackups,
                                                    value: () => data.formSettings.backups.cron,
                                                    oninput: (e) => (data.formSettings.backups.cron = e.target.value),
                                                }),
                                            ),
                                            t.div(
                                                { className: "field addon" },
                                                t.button(
                                                    {
                                                        type: "button",
                                                        className: "btn outline sm",
                                                        "html-popovertarget": "cron-presets-dropdown",
                                                    },
                                                    t.span({ className: "txt" }, "Presets"),
                                                    t.i({ className: "ri-arrow-drop-down-line", ariaHidden: true }),
                                                ),
                                                t.div(
                                                    {
                                                        id: "cron-presets-dropdown",
                                                        className: "dropdown sm txt-nowrap",
                                                        popover: "auto",
                                                    },
                                                    () => {
                                                        return presets.map((preset) => {
                                                            return t.button({
                                                                type: "button",
                                                                className: () =>
                                                                    `dropdown-item ${
                                                                        data.formSettings.backups.cron == preset.cron
                                                                            ? "active"
                                                                            : ""
                                                                    }`,
                                                                textContent: preset.label,
                                                                onclick: (e) => {
                                                                    data.formSettings.backups.cron = preset.cron;
                                                                    e.target.closest(".dropdown").hidePopover();
                                                                },
                                                            });
                                                        });
                                                    },
                                                ),
                                            ),
                                        ),
                                        t.div(
                                            { className: "field-help" },
                                            "Supports numeric list, steps, ranges or ",
                                            t.strong(
                                                {
                                                    className: "link-hint tooltip-bottom",
                                                    ariaDescription: app.attrs.tooltip(
                                                        "@yearly\n@annually\n@monthly\n@weekly\n@daily\n@midnight\n@hourly",
                                                    ),
                                                },
                                                "macros",
                                            ),
                                            ".",
                                            t.br(),
                                            "By default the timezone is in UTC.",
                                        ),
                                    ),
                                    t.div(
                                        { className: "col-lg-6" },
                                        t.div(
                                            { className: "field" },
                                            t.label({ htmlFor: "backups.cronMaxKeep" }, "Max @auto backups to keep"),
                                            t.input({
                                                id: "backups.cronMaxKeep",
                                                name: "backups.cronMaxKeep",
                                                type: "number",
                                                required: () => data.enableAutoBackups,
                                                min: 1,
                                                value: () => data.formSettings.backups.cronMaxKeep,
                                                oninput: (e) => {
                                                    data.formSettings.backups.cronMaxKeep = parseInt(
                                                        e.target.value,
                                                        10,
                                                    );
                                                },
                                            }),
                                        ),
                                    ),
                                ),
                            ),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            app.components.s3ConfigFields({
                                toggleLabel: "Store backups in S3 storage",
                                testFilesystem: "backups",
                                config: () => data.formSettings.backups.s3,
                            }),
                        ),
                        t.div({ className: "col-lg-12" }, t.hr()),
                        t.div(
                            { className: "col-lg-12" },
                            t.div(
                                { className: "flex" },
                                t.div({ className: "m-r-auto" }),
                                t.button(
                                    {
                                        hidden: () => !data.hasChanges,
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
                            ),
                        ),
                    ];
                },
            ),
        ),
    );
}
