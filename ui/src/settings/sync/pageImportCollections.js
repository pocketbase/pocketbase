import { settingsSidebar } from "../settingsSidebar";

export function pageImportCollections(route) {
    app.store.title = "Import collections";

    const uniqueId = "import_" + app.utils.randomString();

    const watchers = [];

    const data = store({
        rawNewCollections: "",
        oldCollections: [],
        newCollections: [],
        collectionsToUpdate: [],
        deleteMissing: true,
        isLoadingFile: false,
        isLoadingOldCollections: false,
        mergeWithOldCollections: false, // an alternative to the default deleteMissing option
        get isRawValid() {
            return (
                !!data.rawNewCollections
                && data.newCollections?.length > 0
                && data.newCollections.length == data.newCollections.filter((c) => !!c.id && !!c.name).length
            );
        },
        get collectionsToDelete() {
            return data.oldCollections.filter((oldC) => {
                return (
                    data.isRawValid
                    && !data.mergeWithOldCollections
                    && data.deleteMissing
                    && !data.newCollections.find((c) => c.id == oldC.id)
                );
            });
        },
        get collectionsToAdd() {
            return data.newCollections.filter((newC) => {
                return data.isRawValid && !data.oldCollections.find((c) => c.id == newC.id);
            });
        },
        // @see replaceIds()
        get idReplacableCollections() {
            return data.newCollections.filter((collection) => {
                let old = data.oldCollections.find((c) => c.name == collection.name || c.id == collection.id);
                if (!old) {
                    return false; // new
                }

                if (old.id != collection.id) {
                    return true;
                }

                // check for matching schema fields
                const oldFields = Array.isArray(old.fields) ? old.fields : [];
                const newFields = Array.isArray(collection.fields) ? collection.fields : [];
                for (const field of newFields) {
                    const oldFieldById = oldFields.find((f) => f.id == field.id);
                    if (oldFieldById) {
                        continue; // no need to do any replacements
                    }

                    const oldFieldByName = oldFields.find((f) => f.name == field.name);
                    if (oldFieldByName && field.id != oldFieldByName.id) {
                        return true;
                    }
                }

                return false;
            });
        },
        get hasChanges() {
            return (
                !!data.rawNewCollections
                && !!(data.collectionsToDelete.length || data.collectionsToAdd.length
                    || data.collectionsToUpdate.length)
            );
        },
        get canReview() {
            return !data.isLoadingOldCollections && data.isRawValid && data.hasChanges;
        },
    });

    const fileInput = t.input({
        id: uniqueId + "_load_json",
        type: "file",
        className: "hidden",
        accept: ".json",
        onchange: () => {
            loadFile(fileInput.files?.[0]);
        },
    });

    loadOldCollections();

    async function loadOldCollections() {
        data.isLoadingOldCollections = true;

        try {
            const collections = await app.pb.collections.getFullList();
            for (let collection of collections) {
                // delete timestamps
                delete collection.created;
                delete collection.updated;

                // unset oauth2 providers
                delete collection.oauth2?.providers;
            }

            data.oldCollections = collections;
            data.isLoadingOldCollections = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                data.isLoadingOldCollections = false;
            }
        }
    }

    watchers.push(
        watch(
            () => data.rawNewCollections,
            () => {
                loadNewCollections();
            },
        ),
    );

    function loadNewCollections() {
        let collections = [];

        try {
            collections = JSON.parse(data.rawNewCollections);

            if (!Array.isArray(collections)) {
                collections = [];
            } else {
                collections = app.utils.filterDuplicatesByKey(collections);
            }

            // normalizations
            for (let collection of collections) {
                // delete timestamps
                delete collection.created;
                delete collection.updated;

                // merge fields with duplicated ids
                if (collection.fields) {
                    collection.fields = app.utils.filterDuplicatesByKey(collection.fields);
                }
            }
        } catch (_) {}

        data.newCollections = collections;
    }

    watchers.push(
        watch(
            () => [data.newCollections, data.deleteMissing],
            () => {
                loadCollectionsToUpdate();
            },
        ),
    );

    function loadCollectionsToUpdate() {
        data.collectionsToUpdate = [];

        if (!data.isRawValid) {
            return;
        }

        for (let newCollection of data.newCollections) {
            const oldCollection = data.oldCollections.find((c) => c.id == newCollection.id);
            if (
                // no old collection
                !oldCollection?.id
                // no changes
                || !app.utils.hasCollectionChanges(oldCollection, newCollection, data.deleteMissing)
            ) {
                continue;
            }

            data.collectionsToUpdate.push({
                new: newCollection,
                old: oldCollection,
            });
        }
    }

    function replaceIds() {
        for (let collection of data.newCollections) {
            const old = data.oldCollections.find((c) => c.name == collection.name || c.id == collection.id);
            if (!old) {
                continue;
            }

            const originalId = collection.id;
            const replacedId = old.id;
            collection.id = replacedId;

            // replace field ids
            const oldFields = Array.isArray(old.fields) ? old.fields : [];
            const newFields = Array.isArray(collection.fields) ? collection.fields : [];
            for (const field of newFields) {
                const oldField = oldFields.find((f) => f.name == field.name);
                if (oldField && oldField.id) {
                    field.id = oldField.id;
                }
            }

            // update references
            for (let ref of data.newCollections) {
                if (!Array.isArray(ref.fields)) {
                    continue;
                }
                for (let field of ref.fields) {
                    if (field.collectionId && field.collectionId === originalId) {
                        field.collectionId = replacedId;
                    }
                }
            }

            // update index names that contains the collection id
            for (let i = 0; i < collection.indexes?.length; i++) {
                collection.indexes[i] = collection.indexes[i].replace(
                    /create\s+(?:unique\s+)?\s*index\s*(?:if\s+not\s+exists\s+)?(\S*)\s+on/gim,
                    (v) => v.replace(originalId, replacedId),
                );
            }
        }

        data.rawNewCollections = JSON.stringify(data.newCollections, null, 2);
    }

    function clear() {
        data.rawNewCollections = "";
        fileInput.value = "";
        app.store.errors = null;
    }

    function loadFile(file) {
        data.isLoadingFile = true;

        const reader = new FileReader();

        reader.onload = async (event) => {
            data.isLoadingFile = false;
            fileInput.value = ""; // reset

            data.rawNewCollections = event.target.result;

            await new Promise((r) => setTimeout(r, 0));

            if (!data.newCollections.length) {
                app.toasts.error("Invalid collections configuration.");
                clear();
            }
        };

        reader.onerror = (err) => {
            app.toasts.error("Failed to load the imported JSON.");
            console.warn(err);

            data.isLoadingFile = false;
            fileInput.value = ""; // reset
        };

        reader.readAsText(file);
    }

    function review() {
        const collectionsToImport = !data.mergeWithOldCollections
            ? data.newCollections
            : app.utils.filterDuplicatesByKey(data.oldCollections.concat(data.newCollections));

        app.modals.openImportCollectionsReview(data.oldCollections, collectionsToImport, {
            deleteMissing: data.deleteMissing,
            onsubmit: () => {
                clear();
                loadOldCollections();
            },
        });
    }

    return t.div(
        {
            pbEvent: "pageImportCollections",
            className: "page page-import-collections",
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        settingsSidebar(),
        t.div(
            { className: "page-content full-height" },
            t.header(
                { className: "page-header" },
                t.nav(
                    { className: "breadcrumbs" },
                    t.div({ className: "breadcrumb-item" }, "Settings"),
                    t.div({ className: "breadcrumb-item" }, () => app.store.title),
                ),
            ),
            t.div({ className: "wrapper m-b-base" }, () => {
                if (data.isLoadingOldCollections) {
                    return t.div({ className: "block txt-center" }, t.span({ className: "loader lg" }));
                }

                return t.div(
                    { className: "grid" },
                    t.div(
                        { className: "col-lg-12" },
                        t.span(
                            { className: "txt-lg m-r-5" },
                            "Paste below the collections configuration you want to import or",
                        ),
                        t.label(
                            {
                                htmlFor: fileInput.id,
                                className: () => `btn sm outline ${data.isLoadingFile ? "loading" : ""}`,
                            },
                            t.span({ className: "txt" }, "Load from JSON file"),
                        ),
                        fileInput,
                        t.p(
                            { className: "txt-hint" },
                            t.em(
                                null,
                                "You can use the ",
                                t.a({
                                    href: `${import.meta.env.PB_DOCS_URL}/go-migrations/`,
                                    target: "_blank",
                                    rel: "noopener noreferrer",
                                    textContent: "Go",
                                }),
                                " or ",
                                t.a({
                                    href: `${import.meta.env.PB_DOCS_URL}/js-migrations/`,
                                    target: "_blank",
                                    rel: "noopener noreferrer",
                                    textContent: "JS",
                                }),
                                " migrations to manage your collections programmatically in more granular and version controlled manner.",
                            ),
                        ),
                    ),
                    t.div(
                        { className: "col-lg-12" },
                        t.div(
                            { className: "field" },
                            t.label({ htmlFor: uniqueId + "_collections_field" }, "Collections"),
                            t.textarea({
                                id: uniqueId + "_collections_field",
                                name: "collections",
                                rows: 12,
                                className: "txt-code",
                                spellcheck: false,
                                autocorrect: false,
                                autocomplete: "off",
                                autocapitalize: "off",
                                value: () => data.rawNewCollections,
                                oninput: (e) => (data.rawNewCollections = e.target.value),
                            }),
                        ),
                        t.div(
                            {
                                className: () =>
                                    `field-help error ${!!data.rawNewCollections && !data.isRawValid ? "" : "hidden"}`,
                            },
                            "Invalid collections configuration.",
                        ),
                    ),
                    t.div(
                        { className: () => `col-lg-12 ${!data.isRawValid ? "hidden" : ""}` },
                        t.div(
                            { className: "field" },
                            t.input({
                                id: uniqueId + "_merge_checkbox",
                                type: "checkbox",
                                className: "switch",
                                checked: () => data.mergeWithOldCollections,
                                onchange: (e) => (data.mergeWithOldCollections = e.target.checked),
                            }),
                            t.label({ htmlFor: uniqueId + "_merge_checkbox" }, "Merge with the existing collections"),
                        ),
                    ),
                    t.div(
                        {
                            className: () => `col-lg-12 ${data.isRawValid && !data.hasChanges ? "" : "hidden"}`,
                        },
                        t.div(
                            { className: "alert info" },
                            t.div(
                                { className: "content" },
                                t.p(null, "Your collections configuration is already up-to-date!"),
                            ),
                        ),
                    ),
                    t.div(
                        {
                            className: () => `col-lg-12 ${data.isRawValid && data.hasChanges ? "" : "hidden"}`,
                        },
                        t.p({ className: "txt-hint txt-bold" }, "Detected changes"),
                        t.div(
                            { className: "list" },
                            // to delete
                            () => {
                                return data.collectionsToDelete.map((collection) => {
                                    return t.div(
                                        { className: "list-item" },
                                        t.span({
                                            className: "label import-change-label danger",
                                            textContent: "Deleted",
                                        }),
                                        t.div(
                                            { className: "inline-flex gap-5" },
                                            t.strong({ textContent: () => collection.name }),
                                            t.small({
                                                className: () => `txt-hint ${!collection.id ? "hidden" : ""}`,
                                                textContent: () => collection.id,
                                            }),
                                        ),
                                    );
                                });
                            },
                            // to update
                            () => {
                                return data.collectionsToUpdate.map((pair) => {
                                    return t.div(
                                        { className: "list-item" },
                                        t.span({
                                            className: "label import-change-label warning",
                                            textContent: "Changed",
                                        }),
                                        t.div(
                                            { className: "inline-flex gap-5" },
                                            () => {
                                                if (pair.old.name == pair.new.name) {
                                                    return;
                                                }

                                                return [
                                                    t.span({
                                                        className: "txt-strikethrough txt-hint",
                                                        textContent: pair.old.name,
                                                    }),
                                                    t.i({
                                                        className: "ri-arrow-right-line txt-sm",
                                                        ariaHidden: true,
                                                    }),
                                                ];
                                            },
                                            t.strong({ textContent: () => pair.new.name }),
                                            t.small({
                                                className: () => `txt-hint ${!pair.new.id ? "hidden" : ""}`,
                                                textContent: () => pair.new.id,
                                            }),
                                        ),
                                    );
                                });
                            },
                            // to add
                            () => {
                                return data.collectionsToAdd.map((collection) => {
                                    return t.div(
                                        { className: "list-item" },
                                        t.span({
                                            className: "label import-change-label success",
                                            textContent: "Added",
                                        }),
                                        t.div(
                                            { className: "inline-flex gap-5" },
                                            t.strong({ textContent: () => collection.name }),
                                            t.small({
                                                className: () => `txt-hint ${!collection.id ? "hidden" : ""}`,
                                                textContent: () => collection.id,
                                            }),
                                        ),
                                    );
                                });
                            },
                        ),
                    ),
                    t.div(
                        {
                            className: () => `col-lg-12 ${!data.idReplacableCollections?.length ? "hidden" : ""}`,
                        },
                        t.div(
                            { className: "alert warning" },
                            t.div(
                                { className: "content" },
                                t.p(
                                    null,
                                    "Some of the imported collections share the same name and/or fields but are imported with different IDs.",
                                ),
                                t.p(
                                    null,
                                    "You can replace them in the import if you want to:",
                                    t.button({
                                        type: "button",
                                        className: "btn warning sm m-l-10",
                                        textContent: "Replace with original IDs",
                                        onclick: replaceIds,
                                    }),
                                ),
                            ),
                        ),
                    ),
                    t.div(
                        { className: "col-lg-12" },
                        t.div(
                            { className: "flex" },
                            t.button(
                                {
                                    type: "button",
                                    className: () => `btn secondary ${!data.rawNewCollections ? "hidden" : ""}`,
                                    onclick: clear,
                                },
                                t.span({ className: "txt" }, "Clear"),
                            ),
                            t.button(
                                {
                                    type: "button",
                                    className: "btn expanded-lg m-l-auto",
                                    disabled: () => !data.canReview,
                                    onclick: review,
                                },
                                t.span({ className: "txt" }, "Review"),
                            ),
                        ),
                    ),
                );
            }),
            t.footer({ className: "page-footer" }, app.components.credits()),
        ),
    );
}
