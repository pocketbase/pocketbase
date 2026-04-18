window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openIndexUpsert = function(
    collection,
    index = "",
    settings = {
        onsave: () => {},
        ondelete: () => {},
    },
) {
    const modal = indexUpsertModal(collection, index, settings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);
    app.modals.open(modal);
};

function indexUpsertModal(collection, index = "", settings = {}) {
    if (!collection) {
        console.warn("[indexUpsertModal] missing required collection argument");
        return;
    }

    let modal;

    const uniqueId = app.utils.randomString();

    const data = store({
        originalIndex: "",
        index: "",
        get isNew() {
            return data.originalIndex == "";
        },
        get indexParts() {
            return app.utils.parseIndex(data.index);
        },
        get lowerCasedIndexColumnNames() {
            return data.indexParts.columns.map((c) => c.name.toLowerCase());
        },
        get canSave() {
            return data.lowerCasedIndexColumnNames.length > 0;
        },
    });

    const presetColumns = collection?.fields?.filter((f) => !f["@toDelete"] && f.name != "id")?.map((f) => f.name)
        || [];

    function loadIndex(index) {
        data.originalIndex = index || "";

        if (!index) {
            const parsed = app.utils.parseIndex("");
            parsed.tableName = collection?.name || "";
            index = app.utils.buildIndex(parsed);
        }

        data.index = index;
    }

    function saveIndex() {
        if (!collection || !data.canSave) {
            console.warn("[saveIndex] no collection or invalid save state:", collection, data.canSave);
            return;
        }

        collection.indexes = collection.indexes || [];

        // search for existing
        const pos = collection.indexes.findIndex((index) => index == data.originalIndex);
        if (pos >= 0) {
            // replace
            collection.indexes[pos] = data.index;
            app.utils.deleteByPath(app.store.errors, "indexes." + pos);
        } else {
            // push missing
            collection.indexes.push(data.index);
        }

        if (typeof settings?.onsave == "function") {
            settings.onsave({
                collection: collection,
                index: data.index,
                oldIndex: data.originalIndex,
            });
        }

        clearIndexError();

        app.modals.close(modal);
    }

    function deleteIndex() {
        if (!collection || !data.originalIndex) {
            console.warn("[deleteIndex] no collection or index:", collection, data.originalIndex);
            return;
        }

        const pos = collection.indexes?.findIndex((index) => index == data.originalIndex);
        if (pos == -1) {
            console.warn("[deleteIndex] missing index:", data.originalIndex);
            return;
        }

        collection.indexes.splice(pos, 1);
        app.utils.deleteByPath(app.store.errors, "indexes." + pos);

        if (typeof settings?.ondelete == "function") {
            settings.ondelete({
                collection: collection,
                position: pos,
                index: data.originalIndex,
            });
        }

        clearIndexError();

        app.modals.close(modal);
    }

    function toggleColumn(column) {
        const clone = JSON.parse(JSON.stringify(data.indexParts));
        clone.tableName = collection?.name || "";

        const colLowerCased = column.toLowerCase();

        const i = clone.columns.findIndex((c) => c.name.toLowerCase() == colLowerCased);
        if (i >= 0) {
            clone.columns.splice(i, 1);
        } else {
            app.utils.pushUnique(clone.columns, { name: column });
        }

        data.index = app.utils.buildIndex(clone);

        clearIndexError();
    }

    function clearIndexError() {
        if (app.store.errors?.indexes) {
            const pos = collection.indexes.findIndex((idx) => idx == data.originalIndex);
            app.utils.deleteByPath(app.store.errors, "indexes." + pos);
        }
    }

    modal = t.div(
        {
            className: "modal popup index-upsert-modal",
            onbeforeopen: () => {
                loadIndex(index);
            },
            onafteropen: () => {
                // retrigger indexes error (if any)
                if (app.store.errors?.indexes) {
                    app.store.errors.indexes = JSON.parse(JSON.stringify(app.store.errors.indexes));
                }
            },
            onafterclose: (el) => {
                el?.remove();
            },
        },
        t.header(
            { className: "modal-header" },
            t.h6(
                { className: "modal-title" },
                t.span({ className: "txt" }, () => (data.isNew ? "Create index" : "Update index")),
            ),
        ),
        t.div(
            { className: "modal-content" },
            t.form(
                {
                    id: uniqueId + "form",
                    className: "grid sm index-upsert-form",
                    onsubmit: (e) => {
                        e.preventDefault();
                        saveIndex();
                    },
                },
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "field" },
                        t.input({
                            type: "checkbox",
                            className: "switch",
                            id: uniqueId + "checkbox_unique",
                            checked: () => data.indexParts.unique,
                            onchange: (e) => {
                                const newIndexParts = JSON.parse(JSON.stringify(data.indexParts));
                                newIndexParts.unique = e.target.checked;
                                newIndexParts.tableName = newIndexParts.tableName || collection?.name || "";
                                data.index = app.utils.buildIndex(newIndexParts);
                            },
                        }),
                        t.label({ htmlFor: uniqueId + "checkbox_unique" }, "Unique"),
                    ),
                ),
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "field" },
                        app.components.codeEditor({
                            required: true,
                            className: "collection-index-input pre-wrap",
                            name: () => "indexes." + collection.indexes?.findIndex((idx) => idx == data.originalIndex),
                            placeholder: () => `e.g. CREATE INDEX idx_test on ${collection?.name || "X"} (created)`,
                            value: () => data.index,
                            oninput: (val) => (data.index = val),
                        }),
                    ),
                    t.div(
                        { hidden: () => !presetColumns.length, className: "field-help m-t-sm" },
                        t.div(
                            { className: "flex flex-wrap gap-5" },
                            t.span({ className: "txt", textContent: "Presets:" }),
                            () => {
                                return presetColumns?.map((col) => {
                                    const isSelected = data.lowerCasedIndexColumnNames.includes(col.toLowerCase());
                                    return t.button({
                                        type: "button",
                                        textContent: col,
                                        className: () => `label handle ${isSelected ? "success" : ""}`,
                                        onclick: () => toggleColumn(col),
                                    });
                                });
                            },
                        ),
                    ),
                ),
            ),
        ),
        t.footer(
            { className: "modal-footer gap-base" },
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
                    hidden: () => data.isNew,
                    type: "button",
                    className: () => "btn sm circle transparent secondary",
                    ariaDescription: app.attrs.tooltip("Delete index", "left"),
                    onclick: () => {
                        app.modals.confirm(
                            "Do you really want to remove the selected index from the collection?",
                            deleteIndex,
                        );
                    },
                },
                t.i({ className: "ri-delete-bin-7-line" }),
            ),
            t.button(
                {
                    "type": "submit",
                    "html-form": uniqueId + "form",
                    "disabled": () => !data.canSave,
                    "className": () => "btn expanded",
                },
                t.span({ className: "txt" }, "Set index"),
            ),
        ),
    );

    return modal;
}
