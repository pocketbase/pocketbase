const PINNED_STORAGE_KEY = "pbPinnedCollections";

const compactThreshold = 12;

export function collectionsSidebar() {
    const data = store({
        search: "",
        pinned: app.utils.getLocalHistory(PINNED_STORAGE_KEY, []),
        get filteredCollections() {
            if (!data.search.length) {
                return app.store.collections;
            }

            const normalizedSearch = data.search.replaceAll(" ", "").toLowerCase();

            return app.store.collections.filter((c) => {
                return (c.name + c.id + c.type).toLowerCase().includes(normalizedSearch);
            });
        },
        get systemCollections() {
            return data.filteredCollections.filter((c) => c.system && !data.pinned.includes(c.id));
        },
        get regularCollections() {
            return data.filteredCollections.filter((c) => !c.system && !data.pinned.includes(c.id));
        },
        get pinnedCollections() {
            if (!data.pinned.length) {
                return [];
            }

            return data.filteredCollections.filter((c) => data.pinned.includes(c.id));
        },
    });

    function clearSearch() {
        data.search = "";
    }

    const watchers = [];

    return app.components.pageSidebar(
        {
            className: () => `collections-sidebar ${data.responsiveShow ? "active" : ""}`,
            onmount: (el) => {
                // init and persist pinned changes
                watchers.push(watch(() => {
                    app.utils.saveLocalHistory(PINNED_STORAGE_KEY, JSON.stringify(data.pinned));
                }));

                // scroll to the active item
                watchers.push(watch(
                    () => app.store.activeCollection?.id,
                    async () => {
                        await new Promise((r) => setTimeout(r, 0));

                        const activeNavItem = el?.querySelector(".nav-item.active");
                        const details = activeNavItem?.closest("details");
                        if (details) {
                            details.open = true;
                            activeNavItem?.scrollIntoView({ block: "nearest" });
                        }
                    },
                ));
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: "sidebar-search" },
            t.div(
                { className: "fields" },
                t.div(
                    { className: "field" },
                    t.input({
                        className: "p-r-5",
                        type: "text",
                        placeholder: "Search collections...",
                        value: () => data.search,
                        oninput: (e) => data.search = e.target.value,
                    }),
                ),
                t.div(
                    { className: "field addon p-l-0 p-r-5 gap-0" },
                    t.button(
                        {
                            hidden: () => !data.search.length,
                            type: "button",
                            className: "btn sm circle transparent secondary",
                            ariaDescription: app.attrs.tooltip("Clear", "left"),
                            onclick: clearSearch,
                        },
                        t.i({ className: "ri-close-line", ariaHidden: true }),
                    ),
                    t.button(
                        {
                            type: "button",
                            disabled: () => app.store.isLoadingCollections,
                            className: () =>
                                `btn sm circle transparent secondary link-faded ${
                                    app.store.isLoadingCollections ? "loading" : ""
                                }`,
                            ariaDescription: app.attrs.tooltip("Collections overview", "left"),
                            onclick: () => app.modals.openCollectionsOverview(),
                        },
                        t.i({ className: "ri-organization-chart", ariaHidden: true }),
                    ),
                ),
            ),
        ),
        () => {
            if (
                !data.search.length
                || !!data.filteredCollections.length
                || app.store.isLoadingCollections
            ) {
                return;
            }

            return t.div(
                { className: "block p-t-base txt-center txt-hint" },
                t.p(null, "No collections found."),
                t.button({
                    type: "button",
                    className: "btn sm secondary",
                    textContent: "Clear search",
                    onclick: () => clearSearch(),
                }),
            );
        },
        // show the standalone loader only when there are no other collections loaded
        () => {
            if (app.store.isLoadingCollections && !data.filteredCollections.length) {
                return t.div({ className: "sidebar-content txt-center" }, t.span({ className: "loader sm" }));
            }
        },
        () => {
            return [
                t.nav(
                    {
                        className: () =>
                            `sidebar-content collections-list scrollable ${
                                data.regularCollections.length + data.pinnedCollections >= compactThreshold
                                    ? "compact"
                                    : ""
                            }`,
                    },
                    t.details(
                        {
                            hidden: () => !data.pinnedCollections.length,
                            className: () => `nav-group nav-group-pinned-collections`,
                            open: true,
                        },
                        t.summary(
                            { tabIndex: -1, onfocusout: () => false, onclick: () => false, onkeyup: () => false },
                            "Pinned",
                        ),
                        () => data.pinnedCollections.map((c) => collectionItem(c, data)),
                    ),
                    t.details(
                        {
                            hidden: () => !data.regularCollections.length,
                            className: "nav-group nav-group-regular-collections",
                            open: true,
                        },
                        t.summary(
                            { tabIndex: -1, onfocusout: () => false, onclick: () => false, onkeyup: () => false },
                            () => data.pinnedCollections.length ? "Others" : "Collections",
                        ),
                        () => data.regularCollections.map((c) => collectionItem(c, data)),
                    ),
                    t.details(
                        {
                            hidden: () => !data.systemCollections.length,
                            className: "nav-group nav-group-system-collections",
                            open: () => data.search.length,
                        },
                        t.summary(null, "System"),
                        () => data.systemCollections.map((c) => collectionItem(c, data)),
                    ),
                ),
                t.div(
                    {
                        hidden: () => data.search.length && !data.filteredCollections.length,
                        className: "sidebar-content new-collection",
                    },
                    t.button(
                        {
                            type: "button",
                            className: "btn outline block",
                            onclick: () => {
                                app.modals.openCollectionUpsert({}, {
                                    onsave: (newCollection) => {
                                        app.store.activeCollection = newCollection.id;
                                    },
                                });
                            },
                        },
                        t.i({ className: "ri-add-line", ariaHidden: true }),
                        t.span({ textContent: "New collection" }),
                    ),
                ),
            ];
        },
    );
}

function collectionItem(collection, data) {
    return t.button(
        {
            "html-data-collection-id": () => collection.id,
            type: "button",
            className: () =>
                `nav-item responsive-close ${collection.id == app.store.activeCollection?.id ? "active" : ""}`,
            title: () => collection.name,
            onauxclick: (e) => {
                e.preventDefault();
                window.open(`#/collections?collection=${collection.name}`, "_blank", "noreferrer,noopener");
            },
            onclick: (e) => {
                e.preventDefault();
                app.store.activeCollection = collection.name;
            },
        },
        t.i({
            className: () => app.collectionTypes[collection.type]?.icon || app.utils.fallbackCollectionIcon,
            ariaHidden: true,
        }),
        t.span({ className: "txt" }, () => collection.name),
        () => {
            if (
                collection.type != "auth"
                || !collection.oauth2?.enabled
                || collection.oauth2?.providers?.length > 0
            ) {
                return;
            }

            return t.i({
                ariaHidden: true,
                className: "ri-alert-line txt-hint txt-sm",
                ariaDescription: app.attrs.tooltip(
                    "OAuth2 auth is enabled but the collection doesn't have any registered providers",
                ),
            });
        },
        () => {
            const pinnedIndex = data.pinned.indexOf(collection.id);

            return t.span(
                {
                    tabIndex: -1,
                    role: "button",
                    className: "pin",
                    title: () => pinnedIndex >= 0 ? "Unpin" : "Pin",
                    onclick: (e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        if (pinnedIndex >= 0) {
                            data.pinned.splice(pinnedIndex, 1);
                        } else {
                            data.pinned.push(collection.id);
                        }
                    },
                },
                t.i({
                    ariaHidden: false,
                    className: () => pinnedIndex >= 0 ? "ri-unpin-line" : "ri-pushpin-line",
                }),
            );
        },
    );
}
