window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Creates an element with a short representation of the specified record
 * (usually based on the presentable fields of the record).
 *
 * @example
 * ```js
 * app.components.recordSummary(record)
 * ```
 *
 * @param  {Object} record
 * @param  {Object} meta
 * @return {Element}
 */
window.app.components.recordSummary = function(record, meta = null) {
    const local = store({
        get collection() {
            return app.store.collections.find(
                (c) => c.id == record.collectionId || c.name == record.collectionName,
            );
        },
        get presentableFields() {
            if (!local.collection?.id) {
                return [];
            }

            const result = local.collection.fields
                .filter((f) => f.presentable)
                .sort((f1, f2) => {
                    const f1Priority = app.fieldTypes[f1.type].summaryPriority || 0;
                    const f2Priority = app.fieldTypes[f2.type].summaryPriority || 0;
                    if (f1Priority > f2Priority) {
                        return 1;
                    }
                    if (f1Priority < f2Priority) {
                        return -1;
                    }
                    return 0;
                });

            // autoset the first found fallback field as presentable
            if (!result.length) {
                for (let name of app.utils.fallbackPresentableProps) {
                    const field = local.collection?.fields?.find((f) => f.name == name);
                    if (field) {
                        result.push(field);
                        break;
                    }
                }
            }

            return result;
        },
    });

    return t.div(
        { className: "label record-summary" },
        t.i({
            ariaHidden: true,
            className: "ri-eye-line link-hint record-preview-icon",
            onclick: (e) => {
                e.stopImmediatePropagation();
                e.preventDefault();
            },
            onmouseenter: (e) => {
                showRecordSummaryDropdown(e.target, record, 100);
            },
            onmouseleave: (e) => {
                hideRecordSummaryDropdown(e.target, 100);
            },
            onunmount: (el) => {
                hideRecordSummaryDropdown(el, 0);
            },
        }),
        () => {
            const result = [];

            function add(val) {
                if (val == null || val == "") {
                    val = t.span({ className: "missing-value" });
                }

                result.push(val);
            }

            for (const field of local.presentableFields) {
                const viewFunc = app.fieldTypes[field.type]?.view;
                if (viewFunc) {
                    const val = viewFunc({
                        short: true,
                        get record() {
                            return record;
                        },
                        get field() {
                            return field;
                        },
                        get meta() {
                            return meta;
                        },
                    });
                    add(val);
                } else {
                    const values = app.utils.toArray(record[field.name]).splice(0, 3);
                    for (const val of values) {
                        add(val);
                    }
                }
            }

            return result;
        },
    );
};

function hideRecordSummaryDropdown(target, delay = 150) {
    if (!target) {
        return;
    }

    clearTimeout(target._summaryDropdownTimeoutId);

    if (delay <= 0) {
        target?._summaryDropdown?.hidePopover?.();
        return;
    }

    target._summaryDropdownTimeoutId = setTimeout(() => {
        target?._summaryDropdown?.hidePopover?.();
    }, delay);
}

function showRecordSummaryDropdown(target, record, delay = 150) {
    if (!target) {
        return;
    }

    clearTimeout(target._summaryDropdownTimeoutId);

    if (delay <= 0) {
        showRecordSummaryDropdownNoDelay(target, record);
        return;
    }

    target._summaryDropdownTimeoutId = setTimeout(() => {
        showRecordSummaryDropdownNoDelay(target, record);
    }, delay);
}

const showRecordSummaryDropdownNoDelay = function(target, record) {
    if (!target) {
        return;
    }

    if (!target._summaryDropdown) {
        target._summaryDropdown = t.div(
            {
                className: "dropdown record-summary-dropdown",
                popover: "manual",
                onclick: (e) => {
                    e.stopImmediatePropagation();
                    e.preventDefault();
                },
            },
            t.div(
                { className: "record-header" },
                t.a(
                    {
                        className: "link-hint txt-bold m-r-auto",
                        target: "_blank",
                        href: `#/collections?collection=${record.collectionName}&record=${record.id}`,
                        onclick: (e) => {
                            e.stopImmediatePropagation();
                        },
                    },
                    t.span({ className: "txt" }, "Edit relation record"),
                    t.i({ className: "ri-external-link-line" }),
                ),
                t.button(
                    {
                        type: "button",
                        className: "link-hint",
                        title: "Close",
                        onclick: () => hideRecordSummaryDropdown(target, 0),
                    },
                    t.i({ className: "ri-close-line", ariaHidden: true }),
                ),
            ),
            t.hr(),
            t.pre(
                { className: "record-json" },
                () => {
                    const fields = app.store.collections.find((c) =>
                        c.id == record.collectionId || c.name == record.collectionName
                    )?.fields || [];
                    if (!fields.length) {
                        return;
                    }

                    const orderedProps = {
                        collectionId: record.collectionId,
                        collectionName: record.collectionName,
                    };
                    for (const field of fields) {
                        orderedProps[field.name] = record[field.name];
                    }

                    return JSON.stringify(app.utils.truncateObject(orderedProps, 27), null, 2);
                },
            ),
        );
        target.appendChild(target._summaryDropdown);
    }

    target._summaryDropdown?.showPopover({
        source: target,
    });
};
