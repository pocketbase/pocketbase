export function collectionsDiffTable(propsArg = {}) {
    const props = store({
        rid: undefined,
        collectionA: null,
        collectionB: null,
        deleteMissing: false,
        className: "",
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        hasAnyChange: false,
        get isDeleteDiff() {
            return !props.collectionB?.id && !props.collectionB?.name;
        },
        get isCreateDiff() {
            return !data.isDeleteDiff && !props.collectionA?.id;
        },
        get hasAnyChange() {
            return app.utils.hasCollectionChanges(props.collectionA, props.collectionB, props.deleteMissing);
        },
        get fieldsListA() {
            return Array.isArray(props.collectionA?.fields) ? props.collectionA?.fields : [];
        },
        get fieldsListB() {
            let fieldsB = Array.isArray(props.collectionB?.fields) ? props.collectionB?.fields : [];

            if (!props.deleteMissing) {
                fieldsB = fieldsB.concat(
                    props.collectionA?.fields?.filter((a) => {
                        return !fieldsB.find((b) => a.id == b.id);
                    }) || [],
                );
            }

            return fieldsB;
        },
        get mainModelProps() {
            return app.utils
                .mergeUnique(Object.keys(props.collectionA || {}), Object.keys(props.collectionB || {}))
                .filter((key) => {
                    return !["fields", "created", "updated"].includes(key);
                });
        },
        get removedFields() {
            return data.fieldsListA.filter((a) => {
                return !data.fieldsListB.find((b) => a.id == b.id);
            });
        },
        get sharedFields() {
            return data.fieldsListB.filter((b) => {
                return data.fieldsListA.find((a) => a.id == b.id);
            });
        },
        get addedFields() {
            return data.fieldsListB.filter((b) => {
                return !data.fieldsListA.find((a) => a.id == b.id);
            });
        },
    });

    function stringify(value) {
        if (typeof value == "undefined") {
            return "";
        }

        return app.utils.isObject(value) ? JSON.stringify(value, null, 4) : "" + value;
    }

    function isDifferent(valA, valB) {
        if (valA === valB) {
            return false; // direct match
        }

        return JSON.stringify(valA) != JSON.stringify(valB);
    }

    function getFieldById(fields, id) {
        return (fields || []).find((f) => f.id == id);
    }

    return t.div(
        {
            rid: props.rid,
            pbEvent: "collectionsDiffTableWrapper",
            className: () => `collections-diff-table-wrapper ${props.className}`,
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: "collections-diff-table-title" },
            () => {
                if (!props.collectionA?.id) {
                    return [
                        t.span({
                            className: "label import-change-label success",
                            textContent: "Added",
                        }),
                        t.strong({ textContent: () => props.collectionB?.name }),
                    ];
                }

                if (!props.collectionB?.id) {
                    return [
                        t.span({
                            className: "label import-change-label danger",
                            textContent: "Deleted",
                        }),
                        t.strong({ textContent: () => props.collectionA?.name }),
                    ];
                }

                return [
                    t.span({
                        hidden: () => !data.hasAnyChange,
                        className: "label import-change-label warning",
                        textContent: "Changed",
                    }),
                    t.div(
                        { className: "inline-flex gap-5" },
                        () => {
                            if (props.collectionA?.name == props.collectionB?.name) {
                                return;
                            }

                            return [
                                t.strong({
                                    className: "txt-strikethrough txt-hint",
                                    textContent: props.collectionA?.name,
                                }),
                                t.i({
                                    className: "ri-arrow-right-line txt-sm",
                                    ariaHidden: true,
                                }),
                            ];
                        },
                        t.strong({ textContent: () => props.collectionB?.name }),
                    ),
                ];
            },
        ),
        t.table(
            { className: "collections-diff-table" },
            t.thead(
                null,
                t.tr(
                    null,
                    t.th({ className: "min-width" }, "Props"),
                    t.th({ width: "40%" }, "Old"),
                    t.th({ width: "40%" }, "New"),
                ),
            ),
            t.tbody(
                null,
                () => {
                    return data.mainModelProps.map((p) => {
                        const isDiff = isDifferent(props.collectionA?.[p], props.collectionB?.[p]);
                        return t.tr(
                            { className: isDiff ? "txt-primary" : "" },
                            t.td({ className: "min-width" }, p),
                            t.td(
                                {
                                    className: () => {
                                        if (data.isCreateDiff) {
                                            return "changed-non-col";
                                        }
                                        if (isDiff) {
                                            return "changed-old-col";
                                        }
                                        return "";
                                    },
                                },
                                t.pre({ className: "txt diff-value" }, stringify(props.collectionA?.[p])),
                            ),
                            t.td(
                                {
                                    className: () => {
                                        if (data.isDeleteDiff) {
                                            return "changed-non-col";
                                        }
                                        if (isDiff) {
                                            return "changed-new-col";
                                        }
                                        return "";
                                    },
                                },
                                t.pre({ className: "txt diff-value" }, stringify(props.collectionB?.[p])),
                            ),
                        );
                    });
                },
                () => {
                    if (!props.deleteMissing && !data.isDeleteDiff) {
                        return;
                    }

                    const rows = [];

                    for (let field of data.removedFields) {
                        rows.push(
                            t.tr(
                                null,
                                t.th(
                                    { className: "min-width", colSpan: 3 },
                                    t.span({ className: "txt" }, "field: ", field.name),
                                    t.span(
                                        { className: "label danger m-l-5" },
                                        "Deleted - ",
                                        t.small(null, `All stored data related to '${field.name}' will be deleted!`),
                                    ),
                                ),
                            ),
                        );

                        for (let key in field) {
                            const val = field[key];

                            rows.push(
                                t.tr(
                                    null,
                                    t.td({ className: "min-width field-key-col" }, key),
                                    t.td(
                                        { className: "changed-old-col" },
                                        t.pre({ className: "txt" }, stringify(val)),
                                    ),
                                    t.td({ className: "changed-none-col" }),
                                ),
                            );
                        }
                    }

                    return rows;
                },
                () => {
                    const rows = [];

                    for (let field of data.sharedFields) {
                        const fieldA = getFieldById(data.fieldsListA, field.id);
                        const fieldB = getFieldById(data.fieldsListB, field.id);

                        const hasFieldChanged = isDifferent(fieldA, fieldB);

                        rows.push(
                            t.tr(
                                null,
                                t.th(
                                    { className: "min-width", colSpan: 3 },
                                    t.span({ className: "txt" }, "field: ", field.name),
                                    t.span({
                                        className: `label warning m-l-5 ${!hasFieldChanged ? "hidden" : ""}`,
                                        textContent: "Changed",
                                    }),
                                ),
                            ),
                        );

                        for (let key in field) {
                            const newValue = field[key];

                            const isDiff = isDifferent(fieldA?.[key], newValue);

                            rows.push(
                                t.tr(
                                    { className: isDiff ? "txt-primary" : "" },
                                    t.td({ className: "min-width field-key-col" }, key),
                                    t.td(
                                        { className: isDiff ? "changed-old-col" : "" },
                                        t.pre({ className: "txt" }, stringify(fieldA?.[key])),
                                    ),
                                    t.td(
                                        { className: isDiff ? "changed-new-col" : "" },
                                        t.pre({ className: "txt" }, stringify(newValue)),
                                    ),
                                ),
                            );
                        }
                    }

                    return rows;
                },
                () => {
                    const rows = [];

                    for (let field of data.addedFields) {
                        rows.push(
                            t.tr(
                                null,
                                t.th(
                                    { className: "min-width", colSpan: 3 },
                                    t.span({ className: "txt" }, "field: ", field.name),
                                    t.span({ className: "label success m-l-5" }, "Added"),
                                ),
                            ),
                        );

                        for (let key in field) {
                            const val = field[key];

                            rows.push(
                                t.tr(
                                    { className: "txt-primary" },
                                    t.td({ className: "min-width field-key-col" }, key),
                                    t.td({ className: "changed-none-col" }),
                                    t.td(
                                        { className: "changed-new-col" },
                                        t.pre({ className: "txt" }, stringify(val)),
                                    ),
                                ),
                            );
                        }
                    }

                    return rows;
                },
            ),
        ),
    );
}
