window.app = window.app || {};
window.app.modals = window.app.modals || {};

const recordsPerPage = 100;

const defaultSettings = {
    btnText: "Insert",
    fileTypes: [], // "image", "document", "video", "audio", "file"
    onselect: function(selectedFile) {},
};

const LAST_SELECTED_STORAGE_KEY = "pbLastRecordFilePickerCollection";
const RECORDS_REQUEST_KEY = "listFilePickerRecords";

/**
 * Opens a new record file picker.
 *
 * @example
 * ```js
 * app.modals.openRecordFilePicker({
 *     onselect: (selectedFile) => { ... }
 * })
 * ```
 *
 * @param {Object} settings
 */
window.app.modals.openRecordFilePicker = function(settings = {}) {
    settings = Object.assign({}, defaultSettings, settings);

    const modal = recordFilePickerModal(settings);

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function recordFilePickerModal(settings = defaultSettings) {
    let modal;

    const uniqueId = "file_picker_" + app.utils.randomString();

    const data = store({
        selectedFile: {},
        records: [],
        activeCollectionId: "",
        searchTerm: "",
        lastRecordsPage: 1,
        lastTotalRecords: 0,
        isLoadingRecords: false,
        get collections() {
            return app.utils.sortedCollections(
                app.store.collections.filter((c) => {
                    if (c.type == "view") {
                        return false;
                    }

                    // has at least one public file key
                    return !!c.fields?.find((f) => {
                        return f.type === "file" && !f.protected;
                    });
                }),
            );
        },
        get activeCollection() {
            const collection = data.collections.find((c) => c.id == data.activeCollectionId);
            if (collection) {
                return collection;
            }

            // always fallback to the first one (if available)
            return data.collections[0];
        },
        get activeCollectionFileFields() {
            return data.activeCollection?.fields?.filter((f) => f.type === "file" && !f.protected) || [];
        },
        get isLoading() {
            return app.store.isLoadingCollections || data.isLoadingRecords;
        },
        get canLoadMore() {
            return !data.isLoadingRecords && data.lastTotalRecords == recordsPerPage;
        },
        get hasAtleastOneFile() {
            return !!data.records.find((r) => extractFiles(r).length > 0);
        },
    });

    const watchers = [];

    // load and sync activeCollectionId from/to localStorage
    watchers.push(
        watch(() => {
            if (!data.activeCollectionId) {
                data.activeCollectionId = window.localStorage.getItem(LAST_SELECTED_STORAGE_KEY);
            } else {
                window.localStorage.setItem(LAST_SELECTED_STORAGE_KEY, data.activeCollectionId);
                data.searchTerm = ""; // reset
            }
        }),
    );

    // reload records on search or active collection change
    watchers.push(
        watch(
            () => [data.activeCollection, data.searchTerm],
            () => loadRecords(true),
        ),
    );

    function resetList() {
        app.pb.cancelRequest(RECORDS_REQUEST_KEY);
        data.isLoadingRecords = false;
        data.records = [];
        data.lastTotalRecords = 0;
        data.lastRecordsPage = 1;
        data.selectedFile = {};
    }

    async function loadRecords(reset = false) {
        if (!data.activeCollection) {
            resetList();
            return;
        }

        if (reset) {
            resetList();
        }

        data.isLoadingRecords = true;

        try {
            const page = reset ? 1 : data.lastRecordsPage + 1;

            const fallbackSearchFields = app.utils.getAllCollectionIdentifiers(data.activeCollection);

            let normalizedFilter = app.utils.normalizeSearchFilter(data.searchTerm, fallbackSearchFields) || "";
            if (normalizedFilter) {
                normalizedFilter += " && ";
            }
            normalizedFilter += "(" + data.activeCollectionFileFields.map((f) => `${f.name}:length>0`).join("||")
                + ")";

            const result = await app.pb.collection(data.activeCollection.id).getList(page, recordsPerPage, {
                requestKey: RECORDS_REQUEST_KEY,
                filter: normalizedFilter,
                skipTotal: 1,
                sort: data.activeCollection.type != "view" ? "-@rowid" : "",
            });

            data.lastRecordsPage = result.page;
            data.lastTotalRecords = result.items.length;
            data.records = app.utils.filterDuplicatesByKey(data.records.concat(result.items));

            data.isLoadingRecords = false;
        } catch (err) {
            if (!err.isAbort) {
                data.isLoadingRecords = false;
                app.checkApiError(err);
            }
        }
    }

    function extractFiles(record) {
        let result = [];

        for (const field of data.activeCollectionFileFields) {
            const names = app.utils.toArray(record[field.name]);
            for (const name of names) {
                if (
                    app.utils.isEmpty(settings.fileTypes)
                    || settings.fileTypes?.includes(app.utils.getFileType(name))
                ) {
                    result.push(name);
                }
            }
        }

        return result;
    }

    function selectFile(record, name) {
        data.selectedFile = { record, name, thumb: "" };
    }

    function isSelected(record, name) {
        return data.selectedFile?.name == name && data.selectedFile?.record?.id == record?.id;
    }

    const documentEvents = {
        "record:create": (e) => {
            if (e.detail.collectionId != data.activeCollection?.id) {
                return;
            }

            if (data.selectedFile?.record?.id == e.detail.id) {
                data.selectedFile.record = e.detail;
            }

            loadRecords(true);
        },
        "record:delete": (e) => {
            if (
                // check both because for delete we don't know which one was assigned to
                e.detail.collectionId != data.activeCollection?.id
                && e.detail.collectionName != data.activeCollection?.name
            ) {
                return;
            }

            if (data.selectedFile?.record?.id == e.detail.id) {
                data.selectedFile = {};
            }

            loadRecords(true);
        },
    };

    modal = t.div(
        {
            className: "modal popup record-file-picker-modal",
            onafterclose: (el) => {
                el?.remove();
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
            // collections select
            t.button(
                {
                    className: () =>
                        `btn primary outline record-file-picker-collection-select-btn ${
                            app.store.isLoadingCollections ? "loading" : ""
                        }`,
                    disabled: () => app.store.isLoadingCollections,
                    "html-popovertarget": "collections_dropdown" + uniqueId,
                },
                t.span(
                    { className: "txt-lg collection-name m-r-auto" },
                    () => data.activeCollection?.name || "Select collection",
                ),
                t.i({ className: "ri-arrow-drop-down-line", ariaHidden: true }),
            ),
            t.div(
                { id: "collections_dropdown" + uniqueId, className: "dropdown", popover: "hint" },
                () => {
                    return data.collections.map((c) => {
                        return t.button(
                            {
                                type: "button",
                                className: () => `dropdown-item ${data.activeCollectionId == c.id ? "active" : ""}`,
                                onclick: (e) => {
                                    data.activeCollectionId = c.id;
                                    e.target?.closest(".dropdown")?.hidePopover();
                                },
                            },
                            c.name,
                        );
                    });
                },
            ),
            // search
            app.components.recordsSearchbar({
                disabled: () => !data.activeCollection?.id,
                collection: () => data.activeCollection,
                value: () => data.searchTerm,
                onsubmit: (newFilter) => (data.searchTerm = newFilter),
            }),
            // new record
            t.button(
                {
                    type: "button",
                    className: "btn circle transparent",
                    ariaLabel: app.attrs.tooltip("Add new record"),
                    onclick: () => app.modals.openRecordUpsert(data.activeCollection),
                },
                t.i({ className: "ri-add-line txt-hint", ariaHidden: true }),
            ),
        ),
        t.div(
            { className: "modal-content" },
            // initial loader
            t.div(
                {
                    className: "block txt-center",
                    hidden: () => data.hasAtleastOneFile || !data.isLoading,
                },
                t.span({ className: "loader" }),
            ),
            // files list
            t.div({ className: "record-file-picker-list" }, () => {
                const result = [];

                for (const record of data.records) {
                    const files = extractFiles(record);
                    for (const name of files) {
                        result.push(
                            t.button(
                                {
                                    rid: record.id + ":" + name,
                                    className: () => `list-item thumb ${isSelected(record, name) ? "success" : ""}`,
                                    ariaDescription: app.attrs.tooltip(name, "bottom"),
                                    onclick: () => selectFile(record, name),
                                },
                                () => {
                                    if (app.utils.hasImageExtension(name)) {
                                        return t.img({
                                            loading: "lazy",
                                            src: app.pb.files.getURL(record, name, { thumb: "100x100" }),
                                            alt: name,
                                        });
                                    }

                                    const ftype = app.utils.getFileType(name);

                                    return t.i({
                                        className: app.utils.fileTypeIcons[ftype] || "ri-file-line",
                                        ariaHidden: true,
                                    });
                                },
                            ),
                        );
                    }
                }

                return result;
            }),
            // load more
            t.div(
                {
                    hidden: () => !data.canLoadMore || !data.hasAtleastOneFile,
                    className: "block txt-center",
                },
                t.button(
                    {
                        className: () => `btn secondary expanded-lg m-t-base ${data.isLoadingRecords ? "loading" : ""}`,
                        disabled: () => data.isLoadingRecords,
                        onclick: () => loadRecords(),
                    },
                    t.span({ className: "txt" }, "Load more"),
                ),
            ),
            // no files
            t.div(
                {
                    className: "block txt-center txt-hint p-t-10 p-b-10",
                    hidden: () => data.hasAtleastOneFile || data.isLoading,
                },
                () => {
                    if (app.utils.isEmpty(settings.fileTypes)) {
                        return t.p(null, "No records with selectable files found.");
                    }
                    return t.p(null, `No "${settings.fileTypes.join("\", \"")}" files found.`);
                },
                t.button({
                    type: "button",
                    className: "btn sm secondary",
                    textContent: "Clear search",
                    hidden: () => !data.searchTerm?.length,
                    onclick: () => {
                        data.searchTerm = "";
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
                    onclick: () => app.modals.close(modal),
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
                    disabled: () => data.isLoading || !data.selectedFile?.name,
                    onclick: () => {
                        const selected = JSON.parse(JSON.stringify(data.selectedFile));

                        if (settings.onselect && settings.onselect(selected) === false) {
                            return false;
                        }

                        app.modals.close(modal);
                    },
                },
                t.span({ className: "txt" }, () => settings.btnText || defaultSettings.btnText),
            ),
        ),
    );

    return modal;
}
