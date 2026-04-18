// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "rel_" + app.utils.randomString();

    // trigger custom change event for clearing field errors
    function triggerChangeEvent() {
        fieldEl?.dispatchEvent(
            new CustomEvent("change", {
                detail: { data: props },
                bubbles: true,
            }),
        );
    }

    const local = store({
        selected: [],
        isLoading: true,
        get maxReached() {
            const maxSelect = props.field.maxSelect || 1;
            return app.utils.toArray(props.record[props.field.name]).length >= maxSelect;
        },
    });

    async function loadSelected() {
        local.isLoading = true;

        const ids = app.utils.toArray(props.record[props.field.name]);
        if (!ids.length) {
            local.selected = [];
            local.isLoading = false;
            return;
        }

        try {
            const fieldCollection = app.store.collections.find((c) => c.id == props.field.collectionId);

            // eagerly expand first level presentable relations (if any and the collections are loaded)
            const relExpands = [];
            const presentableRelationFields = fieldCollection?.fields?.filter(
                (f) => !f.hidden && f.presentable && f.type == "relation",
            ) || [];
            for (const field of presentableRelationFields) {
                relExpands.push(field.name);
            }

            const records = await app.pb.collection(props.field.collectionId).getFullList({
                requestKey: null,
                filter: ids.map((id) => app.pb.filter("id={:id}", { id })).join("||"),
                expand: relExpands.join(",") || undefined,
            });

            // preserve the original order
            const orderedRecords = [];
            for (let id of ids) {
                const record = records.find((r) => r.id == id);
                if (record) {
                    orderedRecords.push(record);
                }
            }

            local.selected = orderedRecords;
            local.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                local.isLoading = false;
            }
        }
    }

    function remove(id) {
        const ids = app.utils.toArray(props.record[props.field.name]);
        const propIndex = ids.indexOf(id);
        if (propIndex >= 0) {
            ids.splice(propIndex, 1);
            updateRecordValue(ids);
        }

        const selectedIndex = local.selected.findIndex((r) => r.id == id);
        local.selected.splice(selectedIndex, 1);
    }

    function updateRecordValue(ids = []) {
        props.record[props.field.name] = props.field.maxSelect > 1 ? ids : ids?.[0] || "";
    }

    const watchers = [
        watch(
            () => props.record[props.field.name],
            () => loadSelected(),
        ),
    ];

    const fieldEl = t.div(
        {
            className: "record-field-input field-type-relation",
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: () => `field ${props.field.required ? "required" : ""}` },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.relation.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            t.output(
                {
                    className: "field-content",
                    name: () => props.field.name,
                },
                // loader
                t.div(
                    {
                        hidden: () => !local.isLoading,
                        className: "list",
                    },
                    () => {
                        const ids = app.utils.toArray(props.record[props.field.name]);
                        return ids.map(() => {
                            return t.div({ className: "list-item" }, t.span({ className: "skeleton-loader" }));
                        });
                    },
                ),
                // list
                app.components.sortable({
                    className: "list",
                    hidden: () => local.isLoading,
                    data: () => local.selected,
                    onchange: (sortedList) => {
                        local.selected = sortedList;
                        updateRecordValue(sortedList.map((r) => r.id));
                        triggerChangeEvent();
                    },
                    dataItem: (record, relIndex) => {
                        return t.div(
                            {
                                rid: record,
                                className: "list-item highlight",
                            },
                            t.div({ className: "content" }, () => app.components.recordSummary(record)),
                            t.div(
                                { className: "actions" },
                                t.button(
                                    {
                                        className: "btn sm secondary transparent circle",
                                        ariaLabel: app.attrs.tooltip("Remove"),
                                        onclick: () => remove(record.id),
                                    },
                                    t.i({ className: "ri-close-line", ariaHidden: true }),
                                ),
                            ),
                        );
                    },
                }),
                // picker btn
                t.hr({
                    hidden: () => !app.utils.isEmpty(props.record[props.field.name]),
                    className: "m-t-5 m-b-0",
                }),
                t.button(
                    {
                        type: "button",
                        className: "btn sm secondary block",
                        disabled: () => local.isLoading,
                        onclick: (e) => {
                            app.modals.openRecordsPicker({
                                collection: props.field.collectionId,
                                selectedIds: app.utils.toArray(props.record[props.field.name]),
                                maxSelect: props.field.maxSelect,
                                onselect: (records) => {
                                    local.selected = records;
                                    updateRecordValue(records.map((r) => r.id));
                                },
                            });
                        },
                    },
                    t.i({ className: "ri-magic-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Open records picker"),
                ),
            ),
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );

    return fieldEl;
}
