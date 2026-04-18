window.app = window.app || {};
window.app.modals = window.app.modals || {};

const recordsPerPage = 50;
const selectedBatchSize = 100;

const RECORDS_REQUEST_KEY = "listRelationPickerRecords";

const defaultSettings = {
    collection: "", // model, id or name
    selectedIds: [],
    maxSelect: 1,
    btnText: "Set selection",
    onselect: function(records) {},
};

/**
 * Opens a new records relation picker.
 *
 * @example
 *
 * ```js
 * app.modals.openRecordsPicker({
 *     collection:  "yourCollection",
 *     selectedIds: ["id1", "id2"],
 *     maxSelect:   1,
 *     onselect:    (records) => { ... }
 * })
 * ```
 *
 * @param {Object} settings
 */
window.app.modals.openRecordsPicker = function(settings = {}) {
    settings = Object.assign({}, defaultSettings, settings);

    const modal = recordsPickerModal(settings);

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function recordsPickerModal(settings = defaultSettings) {
    let modal;

    const data = store({
        searchTerm: "",
        selected: [],

        preselected: [],
        isLoadingPreselected: false,

        records: [],
        isLoadingRecords: false,
        lastRecordsPage: 1,
        lastRecordsTotal: 0,

        get collection() {
            let idOrName = settings.collection;
            if (typeof settings.collection == "object" && settings.collection?.id) {
                idOrName = settings.collection?.id;
            }

            return app.store.collections.find((c) => c.id == idOrName || c.name == idOrName);
        },
        get isLoading() {
            return data.isLoadingPreselected || data.isLoadingRecords;
        },
        get canLoadMore() {
            return !data.isLoadingRecords && data.lastRecordsTotal == recordsPerPage;
        },
    });

    const watchers = [
        watch(
            () => [settings.collection, settings.selectedIds],
            () => {
                loadSelected();
            },
        ),

        // load initial records and reload on search
        watch(
            () => [data.collection, data.searchTerm],
            () => {
                loadRecords(true);
            },
        ),
    ];

    function close() {
        setTimeout(() => app.modals.close(modal), 0); // the popup may not be yet opened
    }

    async function loadSelected() {
        const selectedIds = app.utils.toArray(settings.selectedIds);

        const collectionId = settings.collection?.id || settings.collection;

        if (!collectionId || !selectedIds.length) {
            return;
        }

        data.isLoadingSelected = true;

        let loaded = [];

        // batch load all selected records to avoid filter length errors
        const loadIds = selectedIds.slice();
        const loadPromises = [];
        while (loadIds.length > 0) {
            const filters = [];
            const ids = loadIds.splice(0, selectedBatchSize);
            for (const id of ids) {
                filters.push(`id="${id}"`);
            }

            loadPromises.push(
                app.pb.collection(collectionId).getFullList({
                    requestKey: null,
                    filter: filters.join("||"),
                }),
            );
        }

        try {
            await Promise.all(loadPromises).then((values) => {
                loaded = loaded.concat(...values);
            });

            // preserve selected order
            const orderedSelected = [];
            for (const id of selectedIds) {
                const record = loaded.find((r) => r.id == id);
                if (record) {
                    orderedSelected.push(record);
                }
            }

            // add the ordered selected models to the list (if not already)
            if (!data.searchTerm.trim()) {
                data.records = app.utils.filterDuplicatesByKey(orderedSelected.concat(data.records));
            }

            data.selected = orderedSelected;
            data.isLoadingSelected = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                data.isLoadingSelected = false;
            }
        }
    }

    async function loadRecords(reset = false) {
        if (!data.collection?.id) {
            return;
        }

        if (reset) {
            resetList();

            if (!data.searchTerm.trim()) {
                // prepend the loaded selected items
                data.records = data.selected.slice();
            }
        }

        data.isLoadingRecords = true;

        try {
            const page = reset ? 1 : data.lastRecordsPage + 1;

            const fallbackSearchFields = app.utils.getAllCollectionIdentifiers(data.collection);

            let normalizedFilter = app.utils.normalizeSearchFilter(data.searchTerm, fallbackSearchFields) || "";

            const result = await app.pb.collection(data.collection.id).getList(page, recordsPerPage, {
                requestKey: RECORDS_REQUEST_KEY,
                filter: normalizedFilter,
                skipTotal: 1,
                sort: data.collection.type != "view" ? "-@rowid" : "",
            });

            data.lastRecordsPage = result.page;
            data.lastRecordsTotal = result.items.length;
            data.records = app.utils.filterDuplicatesByKey(data.records.concat(result.items));
            data.isLoadingRecords = false;
        } catch (err) {
            if (!err.isAbort) {
                data.isLoadingRecords = false;
                app.checkApiError(err);
            }
        }
    }

    function resetList() {
        app.pb.cancelRequest(RECORDS_REQUEST_KEY);
        data.isLoadingRecords = false;
        data.records = [];
        data.lastTotalRecords = 0;
        data.lastRecordsPage = 1;
    }

    function scrollHandler(e) {
        const offset = e.target.scrollHeight - e.target.clientHeight - e.target.scrollTop;

        if (offset <= 100 && data.canLoadMore) {
            loadRecords();
        }
    }

    function toggleSelected(record) {
        const index = data.selected.findIndex((r) => r.id == record.id);

        if (index >= 0) {
            data.selected.splice(index, 1);
        } else {
            const maxSelect = settings.maxSelect || 1;

            // clear last redundant elements (leaving place for the new selected)
            let toRemove = data.selected.length - maxSelect;
            while (toRemove >= 0) {
                data.selected.pop();
                toRemove--;
            }

            data.selected.push(record);
        }
    }

    function isSelected(record) {
        return data.selected.findIndex((r) => r.id == record.id) >= 0;
    }

    const documentEvents = {
        "record:save": (e) => {
            if (e.detail.collectionId != data.collection?.id) {
                return;
            }

            const selectedIndex = data.selected?.findIndex((r) => r.id == e.detail.id);
            if (selectedIndex >= 0) {
                data.selected[selectedIndex] = e.detail;
            }

            app.utils.pushOrReplaceObject(data.records, e.detail);
            loadRecords(true);
        },
        "record:delete": (e) => {
            if (
                // check both because for delete we don't know which one was assigned to
                e.detail.collectionId != data.collection?.id
                && e.detail.collectionName != data.collection?.name
            ) {
                return;
            }

            if (isSelected(e.detail)) {
                toggleSelected(e.detail);
            }

            app.utils.removeByKey(data.records, "id", e.detail.id);
            loadRecords(true);
        },
    };

    modal = t.div(
        {
            className: "modal popup lg records-picker-modal",
            onafterclose: (el) => {
                el.remove();
            },
            onmount: (el) => {
                for (let event in documentEvents) {
                    document.addEventListener(event, documentEvents[event]);
                }
            },
            onunmount: (el) => {
                watchers.forEach((w) => w?.unwatch());

                for (let event in documentEvents) {
                    document.removeEventListener(event, documentEvents[event]);
                }
            },
        },
        t.header(
            { className: "modal-header" },
            t.h6({ className: "collection-name" }, () => data.collection.name),
            app.components.recordsSearchbar({
                disabled: () => !data.collection?.id,
                collection: () => data.collection,
                value: () => data.searchTerm,
                onsubmit: (newFilter) => (data.searchTerm = newFilter),
            }),
            t.button(
                {
                    type: "button",
                    className: "btn circle transparent",
                    ariaDescription: app.attrs.tooltip("Add new record"),
                    onclick: () => {
                        app.modals.openRecordUpsert(data.collection);
                    },
                },
                t.i({ className: "ri-add-line txt-hint" }),
            ),
        ),
        t.div(
            { className: "modal-content", hidden: () => data.isLoadingCollection },
            t.div(
                {
                    className: "list records-picker-list",
                    onscroll: scrollHandler,
                    onresize: scrollHandler,
                },
                () => {
                    return data.records.map((record) => {
                        return t.div(
                            {
                                tabIndex: 0,
                                className: "list-item handle",
                                onclick: () => {
                                    toggleSelected(record);
                                    document.activeElement?.blur();
                                },
                            },
                            t.div(
                                { className: "content" },
                                t.span(
                                    { className: "state-icon" },
                                    t.i({
                                        className: () =>
                                            isSelected(record)
                                                ? "ri-checkbox-circle-fill txt-success"
                                                : "ri-checkbox-blank-circle-line txt-disabled",
                                    }),
                                ),
                                () => app.components.recordSummary(record),
                            ),
                            t.div(
                                { className: "actions autohide" },
                                t.button(
                                    {
                                        className: "btn sm secondary transparent circle",
                                        ariaDescription: app.attrs.tooltip("Edit"),
                                        onclick: (e) => {
                                            e.stopPropagation();
                                            app.modals.openRecordUpsert(data.collection, record);
                                        },
                                    },
                                    t.i({ className: "ri-pencil-line" }),
                                ),
                            ),
                        );
                    });
                },
                // loader
                t.div(
                    {
                        className: "list-item",
                        hidden: () => !data.isLoading,
                    },
                    t.div({ className: "skeleton-loader" }),
                ),
                // no records
                t.div(
                    {
                        className: "list-item",
                        hidden: () => data.records.length || data.isLoading,
                    },
                    t.div(
                        { className: "content txt-hint" },
                        t.span({ className: "txt" }, "No records found."),
                        t.button({
                            type: "button",
                            className: "btn sm secondary",
                            textContent: "Clear search",
                            hidden: () => !data.searchTerm.trim().length,
                            onclick: () => {
                                data.searchTerm = "";
                            },
                        }),
                    ),
                ),
            ),
            t.div(
                { className: "block m-t-base" },
                t.p(
                    { className: "txt-bold" },
                    () => `Selected (${data.selected.length} of max ${settings.maxSelect || 1})`,
                ),
                t.span({ className: "txt-hint", hidden: () => data.selected }, "No selected records."),
                app.components.sortable({
                    className: "records-picker-selected-list",
                    data: () => data.selected,
                    dataItem: (record, i) => {
                        return t.div(
                            { rid: record, className: "label handle" },
                            () => app.components.recordSummary(record, [], true),
                            t.span(
                                {
                                    className: "link-hint",
                                    title: "Remove",
                                    role: "button",
                                    onclick: () => toggleSelected(record),
                                },
                                t.i({ className: "ri-close-line", ariaHidden: true }),
                            ),
                        );
                    },
                    onchange: (sortedList, fromIndex, toIndex) => {
                        data.selected = sortedList;
                    },
                }),
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => close(),
                },
                t.span({ className: "txt" }, "Close"),
            ),
            // image thumb selector
            () => {
                if (!data.selectedFile?.name || !app.utils.hasImageExtension(data.selectedFile.name)) {
                    return;
                }

                const options = [
                    { value: "", label: "Original size" },
                    { value: "100x100", label: "100x100 thumb" },
                ];

                // find the related field and its thumbs
                const fileField = data.activeCollectionFileFields.find((f) => {
                    return data.selectedFile.record[f.name].includes(data.selectedFile.name);
                });
                const thumbs = app.utils.toArray(fileField.thumbs);
                for (let thumb of thumbs) {
                    options.push({
                        value: thumb,
                        label: `${thumb} thumb`,
                    });
                }

                return t.div(
                    { className: "record-file-picker-thumb-select" },
                    app.components.select({
                        required: true,
                        value: data.selectedFile.thumb || "",
                        options: options,
                        onchange: (opts) => {
                            data.selectedFile.thumb = opts?.[0].value;
                        },
                    }),
                );
            },
            // submit selected
            t.button(
                {
                    type: "button",
                    className: "btn expanded",
                    disabled: () => data.isLoadingCollection,
                    onclick: () => {
                        const selected = JSON.parse(JSON.stringify(data.selected));

                        if (settings.onselect && settings.onselect(selected) === false) {
                            return false;
                        }

                        close();
                    },
                },
                t.span({ className: "txt" }, settings.btnText || defaultSettings.btnText),
            ),
        ),
    );

    return modal;
}
