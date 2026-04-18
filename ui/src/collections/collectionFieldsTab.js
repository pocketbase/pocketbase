export function collectionFieldsTab(upsertData) {
    return t.div(
        { className: "collection-tab-content collection-fields-tab-content" },
        t.div(
            { className: "collection-fields-list" },
            app.components.sortable({
                handle: ".sort-handle",
                data: () => (upsertData.collection?.fields || [])?.filter((f) => !!app.fieldTypes[f.type]?.settings),
                dataItem: (field, _) => {
                    return app.fieldTypes[field.type].settings({
                        field: field,
                        get originalCollection() {
                            return upsertData.originalCollection;
                        },
                        get collection() {
                            return upsertData.collection;
                        },
                        get originalField() {
                            return upsertData.originalCollection?.fields?.find((f) => field.id && f.id == field.id);
                        },
                        get fieldIndex() {
                            return upsertData.collection.fields?.findIndex((f) =>
                                field.id ? f.id == field.id : f == field
                            );
                        },
                    });
                },
                onchange: (sortedList, fromIndex, toIndex) => {
                    upsertData.collection.fields = sortedList;
                },
            }),
        ),
        () => app.components.addCollectionFieldButton(upsertData.collection),
        // indexes
        t.hr(),
        t.p(
            { className: "txt-bold" },
            "Unique constraints and indexes (",
            () => upsertData.collection.indexes?.length,
            ")",
        ),
        app.components.sortable({
            className: "indexes-list",
            data: () => upsertData.collection.indexes || [],
            onchange: function(sortedList) {
                upsertData.collection.indexes = sortedList;
            },
            dataItem: (index, i) => {
                const parsed = app.utils.parseIndex(index);

                return t.button(
                    {
                        type: "button",
                        className: () => {
                            const errMsg = app.store.errors?.indexes?.[i]?.message || "";
                            return `label handle ${errMsg ? "danger error" : "success"}`;
                        },
                        ariaDescription: app.attrs.tooltip(() => app.store.errors?.indexes?.[i]?.message || ""),
                        onclick: () => app.modals.openIndexUpsert(upsertData.collection, index),
                    },
                    () => {
                        if (parsed.unique) {
                            return t.strong(null, "Unique:");
                        }
                    },
                    t.span({ className: "txt" }, () => parsed.columns?.map((c) => c.name).join(", ")),
                );
            },
            after: () => {
                return t.button(
                    {
                        type: "button",
                        className: "label handle",
                        onclick: () => app.modals.openIndexUpsert(upsertData.collection),
                    },
                    t.i({ className: "ri-add-line" }),
                    t.span({ className: "txt" }, "New index"),
                );
            },
        }),
    );
}
