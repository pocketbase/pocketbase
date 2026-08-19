import { defaultLogLevels } from "./defaultLogLevels";

window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openLogsSettings = function(modalSettings = {
    onbeforeopen: null,
    onafteropen: null,
    onbeforeclose: null,
    onafterclose: null,
    onsave: null,
    ondelete: null,
}) {
    const modal = logsSettingsModal(modalSettings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function logsSettingsModal(modalSettings) {
    let modal;

    const data = store({
        isLoading: false,
        isSaving: false,
        isDeleting: false,
        formSettings: {},
        initFormSettingsHash: "",
        get hasChanges() {
            return data.initFormSettingsHash != JSON.stringify(data.formSettings);
        },
    });

    function init(settings = {}) {
        data.formSettings = {
            logs: settings?.logs || {},
        };

        data.initFormSettingsHash = JSON.stringify(data.formSettings);
    }

    async function loadSettings() {
        data.isLoading = true;

        try {
            const settings = await app.pb.settings.getAll({
                requestKey: "logsSettings",
            });

            init(settings);

            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                data.isLoading = false;
                app.checkApiError(err);
            }
        }
    }

    async function save() {
        if (!data.hasChanges) {
            return;
        }

        data.isSaving = true;

        try {
            const settings = await app.pb.settings.update(app.utils.filterRedactedProps(data.formSettings));

            modalSettings.onsave?.(settings);

            init(settings);

            app.toasts.success("Successfully saved logs settings.");

            data.isSaving = false;

            app.modals.close(modal);
        } catch (err) {
            if (!err.isAbort) {
                data.isSaving = false;
                app.checkApiError(err);
            }
        }
    }

    async function deleteLogs() {
        data.isDeleting = true;

        try {
            await app.pb.logs.truncate();

            modalSettings.ondelete?.();

            app.toasts.success("Successfully deleted all logs.");

            data.isDeleting = false;

            if (!data.hasChanges) {
                app.modals.close(modal);
            }
        } catch (err) {
            if (!err.isAbort) {
                data.isDeleting = false;
                app.checkApiError(err);
            }
        }
    }

    function confirmLogsDelete() {
        app.modals.confirm(
            "Do you really want to delete all logs?",
            () => deleteLogs(),
            null,
            { yesButton: "Yes, delete" },
        );
    }

    modal = t.div(
        {
            pbEvent: "logsSettingsModal",
            className: "modal popup sm logs-settings-modal",
            onbeforeopen: (el) => {
                loadSettings();
                return modalSettings.onbeforeopen?.(el);
            },
            onafteropen: (el) => {
                modalSettings.onafteropen?.(el);
            },
            onbeforeclose: (el) => {
                return modalSettings.onbeforeclose?.(el);
            },
            onafterclose: (el) => {
                modalSettings.onafterclose?.(el);
                el?.remove();
            },
        },
        t.header({ className: "modal-header" }, t.h5({ className: "m-auto" }, "Logs settings")),
        () => {
            if (data.isLoading) {
                return t.div(
                    { className: "modal-content flex", style: "min-height: 200px" },
                    t.span({ className: "loader m-auto" }),
                );
            }

            return [
                t.form(
                    {
                        pbEvent: "logsSettingsForm",
                        id: "logsSettingsForm",
                        className: "modal-content logs-settings-form",
                        onsubmit: (e) => {
                            e.preventDefault();
                            save();
                        },
                    },
                    t.div(
                        { className: "grid" },
                        t.div(
                            { className: "col-lg-12" },
                            t.field(
                                { className: "field" },
                                t.label({ htmlFor: "logs.maxDays" }, "Max days retention"),
                                t.input({
                                    type: "number",
                                    id: "logs.maxDays",
                                    name: "logs.maxDays",
                                    min: 0,
                                    required: true,
                                    value: () => data.formSettings.logs.maxDays,
                                    oninput: (e) => (data.formSettings.logs.maxDays = e.target.value << 0),
                                }),
                            ),
                            t.div(
                                { className: "field-help" },
                                "Set to ",
                                t.code(null, 0),
                                " to delete all logs and disable logs persistence.",
                            ),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            t.field(
                                { className: "field" },
                                t.label({ htmlFor: "logs.minLevel" }, "Min log level"),
                                t.input({
                                    type: "number",
                                    id: "logs.minLevel",
                                    name: "logs.minLevel",
                                    min: -100,
                                    max: 100,
                                    required: true,
                                    value: () => data.formSettings.logs.minLevel,
                                    oninput: (e) => (data.formSettings.logs.minLevel = e.target.value << 0),
                                }),
                            ),
                            t.div(
                                { className: "field-help" },
                                t.div(null, "Logs with level below the minimum will be ignored."),
                                defaultLogLevels(),
                            ),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            t.field(
                                { className: "field" },
                                t.label(
                                    { htmlFor: "logs.maxDataSize" },
                                    t.span({ className: "txt" }, "Max data size"),
                                    t.small(null, "(bytes)"),
                                    t.i({
                                        className: "ri-information-line link-hint",
                                        ariaDescription: app.attrs.tooltip(
                                            `The max size in bytes of the serialized log JSON data field (truncated on "best-effort" basis).`,
                                        ),
                                    }),
                                ),
                                t.input({
                                    type: "number",
                                    id: "logs.maxDataSize",
                                    name: "logs.maxDataSize",
                                    min: 0,
                                    max: Number.MAX_SAFE_INTEGER,
                                    placeholder: "Default to ~16KB",
                                    value: () => data.formSettings.logs.maxDataSize || "",
                                    oninput: (e) => {
                                        if (e.target.value <= 0) {
                                            data.formSettings.logs.maxDataSize = 0;
                                        } else {
                                            data.formSettings.logs.maxDataSize = Number(e.target.value);
                                        }
                                    },
                                    onchange: () => {
                                        // force placeholder rendering
                                        if (!data.formSettings.logs.maxDataSize) {
                                            data.formSettings.logs.maxDataSize = null;
                                            data.formSettings.logs.maxDataSize = 0;
                                        }
                                    },
                                }),
                            ),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            t.field(
                                { className: "field" },
                                t.input({
                                    type: "checkbox",
                                    id: "logs.logIP",
                                    name: "logs.logIP",
                                    className: "switch",
                                    checked: () => data.formSettings.logs.logIP,
                                    onchange: (e) => (data.formSettings.logs.logIP = e.target.checked),
                                }),
                                t.label({ htmlFor: "logs.logIP" }, "Enable IP logging"),
                            ),
                        ),
                        t.div(
                            { className: "col-lg-12" },
                            t.field(
                                { className: "field" },
                                t.input({
                                    type: "checkbox",
                                    id: "logs.logAuthId",
                                    name: "logs.logAuthId",
                                    className: "switch",
                                    checked: () => data.formSettings.logs.logAuthId,
                                    onchange: (e) => (data.formSettings.logs.logAuthId = e.target.checked),
                                }),
                                t.label({ htmlFor: "logs.logAuthId" }, "Enable Auth Id logging"),
                            ),
                        ),
                    ),
                ),
                t.footer(
                    { className: "modal-footer" },
                    t.button(
                        {
                            type: "button",
                            className: "btn transparent m-r-auto",
                            onclick: () => app.modals.close(modal),
                            disabled: () => data.isSaving,
                        },
                        t.span({ className: "txt" }, "Close"),
                    ),
                    t.button(
                        {
                            type: "button",
                            ariaLabel: app.attrs.tooltip("Delete all logs", "left"),
                            className: () =>
                                `btn circle sm secondary transparent link-faded ${data.isDeleting ? "loading" : ""}`,
                            disabled: () => data.isDeleting || data.isSaving,
                            onclick: () => {
                                confirmLogsDelete();
                            },
                        },
                        t.i({ className: "ri-delete-bin-7-line", ariaHidden: true }),
                    ),
                    t.button(
                        {
                            type: "submit",
                            "html-form": "logsSettingsForm",
                            className: () => `btn ${data.isSaving ? "loading" : ""}`,
                            disabled: () => !data.hasChanges || data.isSaving,
                        },
                        t.span({ className: "txt" }, "Save changes"),
                    ),
                ),
            ];
        },
    );

    return modal;
}
