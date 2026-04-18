import { settingsSidebar } from "../settingsSidebar";

export function pageExportCollections(route) {
    app.store.title = "Export collections";

    const uniqueId = "export_" + app.utils.randomString();

    const data = store({
        isLoading: false,
        collections: [],
        bulkSelected: {},
        get bulkSelectStr() {
            return JSON.stringify(app.utils.sortedCollectionsByType(Object.values(data.bulkSelected)), null, 2);
        },
        get totalSelected() {
            return Object.keys(data.bulkSelected).length;
        },
        get areAllSelected() {
            return data.collections.length && data.collections.length == data.totalSelected;
        },
    });

    loadCollections();

    async function loadCollections() {
        data.isLoading = true;

        try {
            let collections = await app.pb.collections.getFullList({
                requestKey: uniqueId,
            });

            for (let collection of collections) {
                // delete timestamps
                delete collection.created;
                delete collection.updated;

                // unset oauth2 providers
                delete collection.oauth2?.providers;
            }

            data.collections = app.utils.sortedCollectionsByType(collections);

            selectAll();

            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                data.isLoading = false;
            }
        }
    }

    function download() {
        const collectionsArr = app.utils.sortedCollectionsByType(Object.values(data.bulkSelected));
        app.utils.downloadJSON(collectionsArr, "pb_schema");
    }

    function toggleSelectAll() {
        if (data.areAllSelected) {
            deselectAll();
        } else {
            selectAll();
        }
    }

    function deselectAll() {
        data.bulkSelected = {};
    }

    // note: allways assign a new object to trigger the getter's Object.keys
    function selectAll() {
        data.bulkSelected = {};

        for (const collection of data.collections) {
            data.bulkSelected[collection.id] = collection;
        }
    }

    function toggleSelectCollection(collection) {
        const bulkSelected = JSON.parse(JSON.stringify(data.bulkSelected));
        if (!data.bulkSelected[collection.id]) {
            bulkSelected[collection.id] = collection;
        } else {
            delete bulkSelected[collection.id];
        }

        // reassign to trigger the getter's Object.keys
        data.bulkSelected = bulkSelected;
    }

    return t.div(
        {
            pbEvent: "pageExportCollections",
            className: "page page-export-collections",
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
                if (data.isLoading) {
                    return t.div({ className: "txt-center" }, t.span({ className: "loader lg" }));
                }

                return t.div(
                    { className: "grid" },
                    t.div(
                        { className: "col-lg-12" },
                        t.div(
                            { className: "txt-lg" },
                            "Below you'll find your current collections configuration that you could import in another PocketBase environment.",
                        ),
                    ),
                    t.div(
                        { className: "col-lg-12" },
                        t.div(
                            { className: "export-panel" },
                            t.aside(
                                { className: "export-list" },
                                t.div(
                                    { className: "list-item" },
                                    t.div(
                                        { className: "field" },
                                        t.input({
                                            id: uniqueId + ".select_all",
                                            type: "checkbox",
                                            checked: () => data.areAllSelected,
                                            onchange: () => toggleSelectAll(),
                                        }),
                                        t.label({ htmlFor: uniqueId + ".select_all" }, "Select all"),
                                    ),
                                ),
                                () => {
                                    return data.collections.map((collection) => {
                                        const checkboxId = uniqueId + "_c_" + collection.id;
                                        return t.div(
                                            { className: "list-item" },
                                            t.div(
                                                { className: "field" },
                                                t.input({
                                                    id: checkboxId,
                                                    type: "checkbox",
                                                    checked: () => !!data.bulkSelected[collection.id],
                                                    onchange: () => {
                                                        toggleSelectCollection(collection);
                                                    },
                                                }),
                                                t.label({ htmlFor: checkboxId }, collection.name),
                                            ),
                                        );
                                    });
                                },
                            ),
                            t.output(
                                { className: "export-preview" },
                                app.components.codeBlock({
                                    value: () => data.bulkSelectStr,
                                    // disable highlight bacause it can cause
                                    // performance issue with too many collections
                                    language: "plain",
                                }),
                                t.nav(
                                    { className: "ctrls" },
                                    app.components.copyButton(() => data.bulkSelectStr),
                                ),
                            ),
                        ),
                    ),
                    t.div(
                        { className: "col-lg-12 txt-right" },
                        t.button(
                            { className: "btn", onclick: download },
                            t.i({ className: "ri-download-line", ariaHidden: true }),
                            t.span({ className: "txt" }, "Download as JSON"),
                        ),
                    ),
                );
            }),
            t.footer({ className: "page-footer" }, app.components.credits()),
        ),
    );
}
