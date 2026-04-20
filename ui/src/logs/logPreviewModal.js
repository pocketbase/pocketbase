import { logLevel } from "./logLevel";

window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openLogPreview = function(logIdOrModel, settings = {
    onbeforeopen: null,
    onafteropen: null,
    onbeforeclose: null,
    onafterclose: null,
}) {
    const modal = logPreviewModal(logIdOrModel, settings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);

    app.modals.open(modal);
};

const priotizedKeys = [
    "execTime",
    "type",
    "auth",
    "authId",
    "status",
    "method",
    "url",
    "referer",
    "remoteIP",
    "userIP",
    "userAgent",
    "error",
    "details",
];

function downloadJSON(log) {
    app.utils.downloadJSON(log, "log_" + log.created.replaceAll(/[-:\. ]/gi, "") + ".json");
}

function copyJSON(log) {
    app.utils.copyToClipboard(JSON.stringify(log, null, 2));
    app.toasts.success("Log copied to clipboard!");
}

function logPreviewModal(logIdOrModel, settings) {
    let modal;

    const data = store({
        isLoading: false,
        log: null,
        get isRequest() {
            return data.log?.data?.type == "request";
        },
        get orderedDataKeys() {
            const result = new Set();

            if (!data.log?.data) {
                return result;
            }

            for (let key of priotizedKeys) {
                if (typeof data.log.data[key] != "undefined") {
                    result.add(key);
                }
            }

            for (let key in data.log.data) {
                result.add(key);
            }

            return result;
        },
    });

    async function load() {
        data.isLoading = true;

        try {
            if (app.utils.isObject(logIdOrModel)) {
                data.log = JSON.parse(JSON.stringify(logIdOrModel));
            } else {
                data.log = await app.pb.logs.getOne(logIdOrModel, {
                    requestKey: "log_preview",
                });
            }

            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                data.isLoading = false;
                app.checkApiError(err);
            }
        }
    }

    modal = t.div(
        {
            pbEvent: "logPreviewModal",
            className: "modal log-preview-modal",
            onbeforeopen: (el) => {
                load();
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
            t.h5(null, "Log details"),
            t.button(
                {
                    className: "btn sm circle transparent m-l-auto",
                    title: "More options",
                    "html-popovertarget": "log-meta-dropdown",
                },
                t.i({ className: "ri-more-line", ariaHidden: true }),
            ),
            t.div({ id: "log-meta-dropdown", className: "dropdown", popover: "auto" }, (el) => {
                return t.button(
                    {
                        className: "dropdown-item",
                        onclick: () => {
                            copyJSON(data.log);
                            el.hidePopover();
                        },
                    },
                    t.i({ className: "ri-braces-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Copy JSON"),
                );
            }),
        ),
        t.div({ className: "modal-content" }, () => {
            if (!data.log || data.isLoading) {
                return t.div({ className: "block txt-center" }, t.span({ className: "loader" }));
            }

            return t.table(
                {
                    pbEvent: "logPreviewTable",
                    className: "log-view-table responsive-table",
                },
                t.tbody(
                    null,
                    t.tr(
                        null,
                        t.th({ className: "col-field-name-id p-r-0" }, "id"),
                        t.td(null, () => data.log.id),
                        t.td({ className: "col-copy min-width" }, app.components.copyButton(data.log.id)),
                    ),
                    t.tr(
                        null,
                        t.th({ className: "col-field-name-level p-r-0" }, "level"),
                        t.td(null, () => logLevel(data.log)),
                        t.td({ className: "col-copy min-width" }, app.components.copyButton(data.log.level)),
                    ),
                    t.tr(
                        null,
                        t.th({ className: "col-field-name-created p-r-0" }, "created"),
                        t.td(
                            null,
                            app.components.formattedDate({
                                value: () => data.log.created,
                                short: false,
                            }),
                        ),
                        t.td({ className: "col-copy min-width" }, app.components.copyButton(data.log.created)),
                    ),
                    () => {
                        if (!data.isRequest) {
                            return t.tr(
                                null,
                                t.th({ className: "col-field-name-message p-r-0" }, "message"),
                                t.td(null, () => app.utils.truncate(data.log.message, 1000)),
                                t.td(
                                    { className: "col-copy min-width" },
                                    app.components.copyButton(data.log.message),
                                ),
                            );
                        }
                    },
                    () => {
                        const rows = [];
                        for (let key of data.orderedDataKeys) {
                            let value = data.log.data?.[key];

                            if (app.utils.logDataFormatters[key]) {
                                value = app.utils.logDataFormatters[key](data.log);
                            }

                            const isEmpty = app.utils.isEmpty(value);
                            const isJSON = !isEmpty && app.utils.isObject(value);
                            if (isJSON) {
                                value = JSON.stringify(value, null, 2);
                            }

                            rows.push(
                                t.tr(
                                    {
                                        rid: "log_data_" + data.log.id + "_" + key,
                                    },
                                    t.th({ className: "min-width p-r-0" }, "data." + key),
                                    t.td(null, () => {
                                        if (isEmpty) {
                                            return t.span({
                                                className: "txt txt-hint",
                                                textContent: "N/A",
                                            });
                                        }

                                        if (key === "error") {
                                            return t.span({
                                                className: `label danger log-error-label ${isJSON ? "txt-code" : ""}`,
                                                textContent: value,
                                            });
                                        }

                                        if (key == "details") {
                                            return t.span({
                                                className: `label warning log-details-label ${
                                                    isJSON ? "txt-code" : ""
                                                }`,
                                                textContent: value,
                                            });
                                        }

                                        if (isJSON) {
                                            return app.components.codeBlock({ value });
                                        }

                                        return t.span({
                                            className: "txt",
                                            textContent: app.utils.displayValue(value, 1000),
                                        });
                                    }),
                                    t.td({ className: "col-copy min-width" }, app.components.copyButton(value)),
                                ),
                            );
                        }
                        return rows;
                    },
                ),
            );
        }),
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
                    type: "button",
                    className: "btn",
                    onclick: () => downloadJSON(data.log),
                },
                t.i({ className: "ri-download-line", ariaHidden: true }),
                t.span({ className: "txt" }, "Download JSON"),
            ),
        ),
    );

    return modal;
}
