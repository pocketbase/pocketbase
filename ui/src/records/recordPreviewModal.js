window.app = window.app || {};
window.app.modals = window.app.modals || {};

/**
 * Opens a record preview modal.
 *
 * @example
 *
 * ```js
 * // full record model
 * app.modals.openRecordPreview(record)
 *
 * // or partial record
 * app.modals.openRecordPreview({ id: "rId", collectionId: "cId"})
 * ```
 *
 * @param {Object} record
 */
window.app.modals.openRecordPreview = function(record, modalSettings = {
    onbeforeopen: null,
    onafteropen: null,
    onbeforeclose: null,
    onafterclose: null,
}) {
    const modal = recordPreviewModal(record, modalSettings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function downloadJSON(record) {
    // clear expand if any
    if (record.expand) {
        record = Object.assign({}, record);
        delete record.expand;
    }

    app.utils.downloadJSON(record, record.collectionName + "_" + record.id + ".json");
}

function copyJSON(record) {
    // clear expand if any
    if (record.expand) {
        record = Object.assign({}, record);
        delete record.expand;
    }

    app.utils.copyToClipboard(JSON.stringify(record, null, 2));
    app.toasts.success("Record copied to clipboard!");
}

function recordPreviewModal(rawRecord, modalSettings) {
    let modal;

    const uniqueId = app.utils.randomString();

    const data = store({
        isLoading: false,
        record: null,
        get collection() {
            return app.store.collections.find((c) => {
                return c.id == rawRecord.collectionId || c.name == rawRecord.collectionName;
            });
        },
    });

    async function loadRecord() {
        if (!rawRecord?.id) {
            app.toasts.error("Failed to load record.");
            setTimeout(() => app.modals.close(modal), 0);
            console.warn("[recordPreviewModal] missing required record id field:", rawRecord);
            return;
        }

        if (!rawRecord.collectionId && !rawRecord.collectionName) {
            app.toasts.error("Failed to load record.");
            setTimeout(() => app.modals.close(modal), 0);
            console.warn("[recordPreviewModal] missing required collectionId or collectionName field:", rawRecord);
            return;
        }

        data.isLoading = true;

        try {
            // eagerly expand first level presentable relations (if any and the collections are loaded)
            let relExpands = [];
            const presentableRelationFields = data.collection?.fields?.filter(
                (f) => !f.hidden && f.presentable && f.type == "relation",
            ) || [];
            for (let field of presentableRelationFields) {
                relExpands.push(field.name);
            }

            data.record = await app.pb
                .collection(rawRecord.collectionId || rawRecord.collectionName)
                .getOne(rawRecord.id, {
                    requestKey: "record_preview_" + rawRecord.id,
                    expand: relExpands.join(",") || undefined,
                });

            data.isLoading = false;
        } catch (err) {
            if (!err?.isAbort) {
                data.isLoading = false;
                app.checkApiError(err);
                setTimeout(() => app.modals.close(modal), 0);
            }
        }
    }

    modal = t.div(
        {
            pbEvent: "recordPreviewModal",
            className: "modal record-preview-modal",
            onbeforeopen: (el) => {
                loadRecord();
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
            onmount: (el) => {
            },
            onunmount: (el) => {
            },
        },
        t.header(
            { className: "modal-header" },
            t.h6(
                null,
                t.strong(null, () => rawRecord?.collectionName || data.collection?.name),
                " record preview",
            ),
            t.button(
                {
                    title: "More options",
                    className: "btn sm circle transparent m-l-auto",
                    "html-popovertarget": uniqueId + "preview-dropdown",
                },
                t.i({ className: "ri-more-line", ariaHidden: true }),
            ),
            t.div({ id: uniqueId + "preview-dropdown", className: "dropdown", popover: "auto" }, (el) => {
                return t.button(
                    {
                        className: "dropdown-item",
                        onclick: () => {
                            copyJSON(data.record);
                            el.hidePopover();
                        },
                    },
                    t.i({ className: "ri-braces-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Copy JSON"),
                );
            }),
        ),
        t.div({ className: "modal-content" }, () => {
            // loader
            if (data.isLoading || !data.record?.id || !data.collection?.id) {
                return t.table(
                    null,
                    t.tbody(null, () => {
                        const totalRows = data.collection?.fields?.filter((f) => f.type != "password").length || 1;
                        const rows = [];

                        for (let i = 0; i < totalRows; i++) {
                            rows.push(t.tr(null, t.td(null, t.span({ className: "skeleton-loader" }))));
                        }

                        return rows;
                    }),
                );
            }

            // attrs
            return t.table(
                {
                    pbEvent: "recordPreviewTable",
                    className: "record-preview-table responsive-table",
                },
                t.tbody(null, () => {
                    const fields = data.collection?.fields?.filter((f) => f.type != "password") || [];

                    return fields.map((f) => {
                        return t.tr(
                            null,
                            t.th(
                                { className: () => `min-width p-r-0 col-field-name-${f.name}` },
                                f.name,
                            ),
                            t.td(
                                { className: () => `col-field-name-${f.name}` },
                                () => {
                                    if (app.fieldTypes[f.type]?.view) {
                                        return app.fieldTypes[f.type].view({
                                            short: false,
                                            get record() {
                                                return data.record;
                                            },
                                            get field() {
                                                return f;
                                            },
                                        });
                                    }

                                    return app.utils.stringifyValue(data.record[f.name]);
                                },
                            ),
                        );
                    });
                }),
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
                    onclick: () => downloadJSON(data.record),
                },
                t.i({ className: "ri-download-line", ariaHidden: true }),
                t.span({ className: "txt" }, "Download JSON"),
            ),
        ),
    );

    return modal;
}
