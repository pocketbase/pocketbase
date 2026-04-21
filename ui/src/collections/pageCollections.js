import { collectionsSidebar } from "./collectionsSidebar";

const SORT_QUERY_KEY = "sort";
const FILTER_QUERY_KEY = "filter";
const COLLECTION_QUERY_KEY = "collection";
const RECORD_QUERY_KEY = "record";
const LAST_ACTIVE_STORAGE_KEY = "pbLastActiveCollection";

export function pageCollections(route) {
    const uniqueId = "page_collections_" + app.utils.randomString();

    app.store.activeCollection = route.query[COLLECTION_QUERY_KEY]?.[0]
        || window.localStorage.getItem(LAST_ACTIVE_STORAGE_KEY);

    const pageData = store({
        reset: null,
        activeRecordIdOrModel: route.query[RECORD_QUERY_KEY]?.[0] || "",
        sort: route.query[SORT_QUERY_KEY]?.[0] || "",
        filter: route.query[FILTER_QUERY_KEY]?.[0] || "",
        totalCount: 0,
        isTotalCountLoading: false,
    });

    async function loadTotalCount() {
        if (!app.store.activeCollection?.id) {
            return;
        }

        pageData.isTotalCountLoading = true;

        try {
            const normalizedFilter = app.utils.normalizeSearchFilter(
                pageData.filter,
                app.store.activeCollection.fields.filter((f) => !f.hidden).map((f) => f.name),
            );

            const result = await app.pb.collection(app.store.activeCollection.name).getList(1, 1, {
                // use a per page unique id and not a global constant to prevent race issues with the async unmount
                requestKey: uniqueId,
                filter: normalizedFilter,
                fields: "id",
            });

            pageData.totalCount = result.totalItems;
        } catch (err) {
            if (!err.isAbort) {
                pageData.totalCount = 0;
                console.warn("failed to load total count:", err);
            }
        }

        pageData.isTotalCountLoading = false;
    }

    function refreshRecordsList() {
        pageData.reset = Date.now();
    }

    const watchers = [
        watch(
            () => (app.store.activeCollection?.name || "") + (app.store.activeCollection?.updated || ""),
            (newVal, oldVal) => {
                app.store.title = app.store.activeCollection?.name || "Collections";

                // skip unnecessery initial params replacement
                if (!oldVal) {
                    return;
                }

                // reset filter and sort params on collection change
                if (oldVal != newVal) {
                    pageData.filter = "";
                    pageData.sort = "";
                }

                app.utils.replaceHashQueryParams({
                    [COLLECTION_QUERY_KEY]: app.store.activeCollection?.name,
                    [FILTER_QUERY_KEY]: pageData.filter || null,
                    [SORT_QUERY_KEY]: pageData.sort || null,
                }, newVal != oldVal ? true : null);

                if (app.store.activeCollection?.id) {
                    window.localStorage.setItem(LAST_ACTIVE_STORAGE_KEY, app.store.activeCollection.id);
                } else {
                    window.localStorage.removeItem(LAST_ACTIVE_STORAGE_KEY);
                }
            },
        ),
        watch(
            () => [pageData.filter, pageData.sort],
            (newVal, oldVal) => {
                if (!oldVal) {
                    return;
                }

                app.utils.replaceHashQueryParams({
                    [FILTER_QUERY_KEY]: pageData.filter || null,
                    [SORT_QUERY_KEY]: pageData.sort || null,
                });
            },
        ),
        watch(
            () => (pageData.activeRecordIdOrModel || "") + (app.store.activeCollection?.id || ""),
            (newVal, oldVal) => {
                if (!pageData.activeRecordIdOrModel) {
                    app.utils.replaceHashQueryParams({
                        [RECORD_QUERY_KEY]: null,
                    });
                    return;
                }

                // no change or the collection model is still loading
                if (newVal == oldVal || !app.store.activeCollection?.id) {
                    return;
                }

                const recordData = typeof pageData.activeRecordIdOrModel == "string"
                    ? {
                        id: pageData.activeRecordIdOrModel,
                        collectionId: app.store.activeCollection?.id,
                        collectionName: app.store.activeCollection?.name,
                    }
                    : pageData.activeRecordIdOrModel;

                app.utils.replaceHashQueryParams({
                    [RECORD_QUERY_KEY]: recordData.id || null,
                });

                // force close any previous modal
                app.modals.close(null, true);

                if (app.store.activeCollection?.type == "view") {
                    app.modals.openRecordPreview(recordData, {
                        onafterclose: () => {
                            pageData.activeRecordIdOrModel = "";
                        },
                    });
                } else {
                    app.modals.openRecordUpsert(app.store.activeCollection, recordData, {
                        onafterclose: () => {
                            pageData.activeRecordIdOrModel = "";
                        },
                    });
                }
            },
        ),
        watch(
            () => [app.store.activeCollection?.id, pageData.filter, pageData.reset],
            () => loadTotalCount(),
        ),
    ];

    const documentEvents = {
        "record:save": (e) => {
            if (e.detail.collectionId != app.store.activeCollection?.id) {
                return;
            }

            pageData.totalCount++;
        },
        "record:delete": (e) => {
            if (
                // check both because for delete we don't know which one was assigned to
                e.detail.collectionId != app.store.activeCollection?.id
                && e.detail.collectionName != app.store.activeCollection?.name
            ) {
                return;
            }

            pageData.totalCount--;
        },
    };

    return t.div(
        {
            pbEvent: "pageCollections",
            className: "page",
            onmount: () => {
                // refresh if necesser the cached collections in the background
                if (!app.store.isLoadingCollections) {
                    app.store.silentlyReloadCollections();
                }

                for (let event in documentEvents) {
                    document.addEventListener(event, documentEvents[event]);
                }
            },
            onunmount: () => {
                app.pb.cancelRequest(uniqueId);

                watchers.forEach((w) => w?.unwatch());

                for (let event in documentEvents) {
                    document.removeEventListener(event, documentEvents[event]);
                }
            },
        },
        () => collectionsSidebar(),
        t.div(
            { className: "page-content full-height" },
            t.header(
                { className: "page-header compact flex-nowrap" },
                t.nav(
                    { className: "breadcrumbs" },
                    t.div(null, "Collections"),
                    () => {
                        if (app.store.activeCollection?.name) {
                            return t.div({
                                title: app.store.activeCollection.name,
                                textContent: app.store.activeCollection.name,
                            });
                        }
                    },
                ),
                t.div(
                    {
                        hidden: () => !app.store.activeCollection?.id,
                        pbEvent: "pageHeaderSecondaryBtns",
                        className: "page-header-secondary-btns",
                    },
                    t.button(
                        {
                            type: "button",
                            className: "btn circle transparent secondary tooltip-bottom btn-collection-settings",
                            ariaLabel: app.attrs.tooltip("Collection settings"),
                            onclick: () => {
                                app.modals.openCollectionUpsert(app.store.activeCollection, {
                                    ontruncate: () => refreshRecordsList(),
                                    onsave: (collection, isNew) => {
                                        if (isNew) {
                                            // e.g. in case of a duplicate or modal state reset
                                            app.store.activeCollection = collection.id;
                                        } else {
                                            refreshRecordsList();
                                        }
                                    },
                                });
                            },
                        },
                        t.i({ className: "ri-settings-3-line", ariaHidden: true }),
                    ),
                    app.components.refreshButton({
                        onclick: () => refreshRecordsList(),
                    }),
                ),
                t.div(
                    {
                        hidden: () => !app.store.activeCollection?.id,
                        pbEvent: "pageHeaderPrimaryBtns",
                        className: "page-header-primary-btns",
                    },
                    t.button(
                        {
                            type: "button",
                            className: "btn outline api-preview-btn",
                            onclick: () => app.modals.openApiPreview(app.store.activeCollection),
                        },
                        t.i({ className: "ri-code-s-slash-line", ariaHidden: true }),
                        t.span({ className: "txt", textContent: "API preview" }),
                    ),
                    () => {
                        if (app.store.activeCollection?.type == "view") {
                            return;
                        }

                        return t.button(
                            {
                                type: "button",
                                className: "btn new-record-btn",
                                onclick: () => app.modals.openRecordUpsert(app.store.activeCollection),
                            },
                            t.i({ className: "ri-add-line", ariaHidden: true }),
                            t.span({ className: "txt", textContent: "New Record" }),
                        );
                    },
                ),
            ),
            // page loader
            t.div(
                {
                    hidden: () => !app.store.isLoadingCollections || app.store.activeCollection?.id,
                    className: "block txt-center p-base",
                },
                t.span({ className: "loader lg" }),
            ),
            // no selected collection
            t.div(
                {
                    hidden: () => app.store.isLoadingCollections || app.store.activeCollection?.id,
                    className: "block txt-center p-base",
                },
                t.h6(
                    { className: "txt" },
                    () => {
                        if (app.store.collections?.length) {
                            return "Select collection from the sidebar.";
                        }
                        return "No collections found.";
                    },
                ),
            ),
            // records list
            app.components.recordsSearchbar({
                hidden: () => !app.store.activeCollection?.id,
                collection: () => app.store.activeCollection,
                value: () => pageData.filter,
                onsubmit: (newFilter) => (pageData.filter = newFilter),
            }),
            app.components.recordsList({
                className: "m-t-sm",
                reset: () => pageData.reset,
                hidden: () => !app.store.activeCollection?.id,
                collection: () => app.store.activeCollection,
                filter: () => pageData.filter,
                sort: () => pageData.sort,
                onselect: (record) => {
                    pageData.activeRecordIdOrModel = record;
                },
                onchange: (newFilter, newSort) => {
                    pageData.filter = newFilter;
                    pageData.sort = newSort;
                },
            }),
            t.footer(
                { className: "page-footer" },
                t.span(
                    {
                        className: () => `total-count ${pageData.isTotalCountLoading ? "faded" : ""}`,
                    },
                    "Total: ",
                    () => pageData.totalCount,
                ),
                app.components.credits(),
            ),
        ),
    );
}
