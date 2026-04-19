import { fieldsWithExcerpt } from "@/fields/relation/view";

window.app = window.app || {};
window.app.components = window.app.components || {};

const perPage = 40;

const sortRegex = /^([\+\-])?(\w+)$/;

window.app.consts = window.app.consts || {};
window.app.consts.COLUMNS_STORAGE_PREFIX = "pbColumns_";

/**
 * Creates new page records listing element.
 *
 * @example
 * ```js
 * app.components.recordsList({
 *     collection: () => data.activeCollection,
 * })
 * ```
 *
 * @param  {Object} propsArg
 * @return {Element}
 */
window.app.components.recordsList = function(propsArg = {}) {
    const uniqueId = "records_list_" + app.utils.randomString();

    const props = store({
        collection: {},
        filter: "",
        sort: "",
        reset: undefined,
        // ---
        rid: undefined,
        id: undefined,
        hidden: undefined,
        className: "",
        onchange: (newFilter, newSort) => {},
        onselect: (record) => {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        isLoading: false,
        records: [],
        lastPage: 0,
        lastTotalItems: 0,
        bulkSelected: {},
        columnsPreferences: {},
        get canLoadMore() {
            return data.lastTotalItems >= perPage;
        },
        get totalSelected() {
            return Object.keys(data.bulkSelected).length;
        },
        get areAllSelected() {
            return data.records.length && data.records.length == data.totalSelected;
        },
        get firstAutoUpdatedField() {
            return props.collection?.fields?.find((f) => f.type == "autodate" && f.onUpdate);
        },
        get isSuperusersCollection() {
            return props.collection?.type == "auth" && props.collection?.name == "_superusers";
        },
    });

    async function clearList() {
        data.records = [];
        data.lastPage = 0;
        data.lastTotalItems = 0;
        data.bulkSelected = {};
    }

    function triggerOnchange() {
        props.onchange?.(props.filter, props.sort);
    }

    async function loadRecords(reset = false) {
        if (!props.collection?.id) {
            return;
        }

        data.isLoading = true;

        try {
            // (note if changed update the related counter query too!)
            const normalizedFilter = app.utils.normalizeSearchFilter(
                props.filter,
                props.collection.fields.filter((f) => !f.hidden).map((f) => f.name),
            );

            // eagerly expand first level relations
            // (to prevent too many relation queries)
            const relExpands = [];
            const relationFields = props.collection.fields.filter(
                (f) => !f.hidden && f.type == "relation",
            );
            for (const field of relationFields) {
                relExpands.push(field.name);
            }

            let requestFields = fieldsWithExcerpt(props.collection.id, relationFields);

            // allow sorting by the top level relation presentable fields
            let normalizedSort = props.sort || undefined;
            const sortMatch = normalizedSort?.match(sortRegex);
            const sortField = sortMatch
                ? props.collection.fields.find((f) => !f.hidden && f.name === sortMatch[2])
                : null;
            if (!sortField) {
                // default fallback to -@rowid when available
                normalizedSort = props.collection.type != "view" ? "-@rowid" : undefined;
            } else if (sortField?.type == "relation") {
                normalizedSort = app.store.collections
                    ?.find((c) => c.id == sortField.collectionId)
                    ?.fields?.filter((f) => f.presentable)
                    ?.map((f) => (sortMatch[1] || "") + sortMatch[2] + "." + f.name)
                    ?.join(",");
            }

            const page = reset ? 1 : data.lastPage + 1;

            const result = await app.pb.collection(props.collection.name).getList(page, perPage, {
                requestKey: uniqueId,
                skipTotal: 1,
                filter: normalizedFilter,
                sort: normalizedSort,
                expand: relExpands.join(",") || undefined,
                fields: requestFields,
            });

            if (result.page == 1) {
                clearList();
            }

            data.lastPage = result.page;
            data.lastTotalItems = result.items.length;

            for (let i = 0; i < result.items.length; i++) {
                app.utils.pushOrReplaceObject(data.records, result.items[i]);

                // yield to main (with room to "breathe")
                if (i > 1 && i % 15 == 0) {
                    await new Promise((r) => setTimeout(r, 20));
                }
            }

            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                data.isLoading = false;
                clearList();
                app.checkApiError(err);
            }
        }
    }

    function selectAll(state = true) {
        // note: always assign a new object to trigger the getter's Object.keys
        const selected = {};
        if (state) {
            for (let record of data.records) {
                selected[record.id] = record;
            }
        }
        data.bulkSelected = selected;
    }

    function downloadSelected() {
        const selected = JSON.parse(JSON.stringify(Object.values(data.bulkSelected)));
        if (!selected.length) {
            return; // nothing to download
        }

        // unset expand
        for (const record of selected) {
            if (record.expand) {
                delete record.expand;
            }
        }

        if (selected.length == 1) {
            return app.utils.downloadJSON(selected[0], props.collection.name + "_" + selected[0].id + ".json");
        }

        return app.utils.downloadJSON(selected, `${selected.length}_${props.collection.name}_records.json`);
    }

    async function deleteSelected() {
        const idsToDelete = Object.keys(data.bulkSelected);
        if (!idsToDelete.length) {
            return; // nothing to delete
        }

        const remainingIdsToDelete = idsToDelete.slice();

        // delete requests in batches to avoid sending too many requests
        while (remainingIdsToDelete.length) {
            const ids = remainingIdsToDelete.splice(0, 100);
            const promises = [];
            for (const id of ids) {
                promises.push(app.pb.collection(props.collection.name).delete(id));
            }
            try {
                await Promise.all(promises);
            } catch (err) {
                app.checkApiError(err);
                selectAll(false);
                loadRecords(true);
                return;
            }
        }

        selectAll(false);

        app.toasts.success(
            `Successfully deleted ${idsToDelete.length} ${idsToDelete.length == 1 ? "record" : "records"}.`,
        );
    }

    function recordRid(record) {
        if (data.firstAutoUpdatedField) {
            // - the collection update is added in case the collection fields have changed
            // - the record keys are added in case of a record field rename
            //   (the collection update and the refreshed records load doesn't happen at the same time)
            return record.id + record[data.firstAutoUpdatedField.name] + props.collection?.updated
                + Object.keys(record);
        }

        return JSON.stringify(record) + props.collection?.updated;
    }

    function isFieldColumnHidden(field) {
        if (typeof data.columnsPreferences[field.id] != "undefined") {
            return !data.columnsPreferences[field.id];
        }

        return field.hidden;
    }

    let deleteRefreshTimeoutId;

    const documentEvents = {
        "record:save": (e) => {
            if (e.detail.collectionId != props.collection?.id) {
                return;
            }

            // optimistically merge with existing to minimize flickering
            const found = data.records.find((r) => r.id == e.detail.id);
            if (found) {
                Object.assign(found, JSON.parse(JSON.stringify(e.detail)));
            }

            loadRecords(true);
        },
        "record:delete": (e) => {
            if (
                // check both because for delete we don't know which one was assigned to
                e.detail.collectionId != props.collection?.id
                && e.detail.collectionName != props.collection?.name
            ) {
                return;
            }

            delete data.bulkSelected[e.detail.id];
            app.utils.removeByKey(data.records, "id", e.detail.id);

            clearTimeout(deleteRefreshTimeoutId);
            deleteRefreshTimeoutId = setTimeout(() => {
                if (!data.records?.length) {
                    loadRecords(true);
                }
            }, 100);
        },
    };

    return t.div(
        {
            pbEvent: "recordsList",
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            className: () => `page-table-wrapper ${props.className}`,
            onmount: (el) => {
                for (let event in documentEvents) {
                    document.addEventListener(event, documentEvents[event]);
                }

                watchers.push(
                    watch(
                        () => props.collection?.id,
                        (newId, oldId) => {
                            data.columnsPreferences = app.utils.getLocalHistory(
                                app.consts.COLUMNS_STORAGE_PREFIX + newId,
                                {},
                            );

                            if (oldId && oldId != newId) {
                                clearList();
                            }
                        },
                    ),
                );

                // trigger load on props change
                watchers.push(
                    watch(
                        () =>
                            (props.collection?.id || "") + (props.filter || "") + (props.sort || "")
                            + (props.reset || ""),
                        (newVal, oldVal) => {
                            if (newVal != oldVal) {
                                loadRecords(true);
                            }
                        },
                    ),
                );

                // always scroll to top on first page load
                watchers.push(
                    watch(
                        () => data.lastPage,
                        (page) => {
                            if (page == 1 && el) {
                                el.scrollTop = 0;
                            }
                        },
                    ),
                );

                watchers.push(
                    watch(
                        () => JSON.stringify(data.columnsPreferences),
                        (newVal, oldVal) => {
                            if (props.collection?.id && oldVal) {
                                app.utils.saveLocalHistory(
                                    app.consts.COLUMNS_STORAGE_PREFIX + props.collection.id,
                                    data.columnsPreferences,
                                );
                            }
                        },
                    ),
                );
            },
            onunmount: () => {
                app.pb.cancelRequest(uniqueId);

                clearTimeout(deleteRefreshTimeoutId);

                watchers.forEach((w) => w?.unwatch());

                for (let event in documentEvents) {
                    document.removeEventListener(event, documentEvents[event]);
                }
            },
        },
        t.table(
            {
                pbEvent: "recordsListTable",
                className: () => `records-table responsive-table ${data.records.length > perPage ? "optimize" : ""}`,
            },
            t.thead(
                { className: "sticky" },
                t.tr(
                    null,
                    t.th(
                        { className: "col-bulk-select" },
                        t.div(
                            {
                                className: "field",
                                hidden: () => data.isLoading,
                            },
                            t.input({
                                id: "all_" + uniqueId,
                                type: "checkbox",
                                disabled: () => !data.records.length,
                                checked: () => data.areAllSelected,
                                onchange: (e) => selectAll(e.target.checked),
                            }),
                            t.label({ htmlFor: "all_" + uniqueId }),
                        ),
                        t.span({
                            className: "loader",
                            hidden: () => !data.isLoading,
                        }),
                    ),
                    () => {
                        const fields = props.collection?.fields || [];

                        const columns = [];

                        for (const field of fields) {
                            if (
                                !app.fieldTypes[field.type]?.view
                                // superusers are always verified
                                || (data.isSuperusersCollection && field.name == "verified")
                            ) {
                                continue;
                            }

                            columns.push(
                                t.th(
                                    {
                                        hidden: () => isFieldColumnHidden(field),
                                        className: () => {
                                            let sortDir = "";
                                            if (props.sort == field.name || props.sort == "+" + field.name) {
                                                sortDir = "asc";
                                            } else if (props.sort == "-" + field.name) {
                                                sortDir = "desc";
                                            }

                                            return `sort-handle ${sortDir} col-field-type-${field.type} col-field-name-${field.name}`;
                                        },
                                        onclick: (e) => {
                                            let newSort = "-" + field.name;
                                            if (props.sort == newSort) {
                                                newSort = field.name;
                                            }
                                            props.sort = newSort;
                                            triggerOnchange();
                                        },
                                    },
                                    t.div(
                                        { className: "inline-flex gap-5" },
                                        t.i({
                                            ariaHidden: true,
                                            className: () => {
                                                if (field.primaryKey) {
                                                    return "ri-key-line";
                                                }

                                                return app.fieldTypes[field.type]?.icon || app.utils.fallbackFieldIcon;
                                            },
                                        }),
                                        t.span({ className: "txt", textContent: field.name }),
                                    ),
                                ),
                            );
                        }

                        return columns;
                    },
                    t.th({ className: "col-meta" }, () => columnsDropdown(props, data)),
                ),
            ),
            t.tbody(
                null,
                () => {
                    if (!data.records.length) {
                        return t.tr(
                            null,
                            t.td({ colSpan: 99, style: "height:59px" }, () => {
                                if (data.isLoading) {
                                    return t.span({ className: "skeleton-loader" });
                                }

                                return t.div(
                                    { className: "sticky-content txt-center txt-hint" },
                                    t.p({ className: "txt-bold" }, "No records found."),
                                    t.button(
                                        {
                                            hidden: () => props.filter?.length || props.collection?.type == "view",
                                            type: "button",
                                            className: "btn secondary expanded-lg",
                                            onclick() {
                                                app.modals.openRecordUpsert(props.collection);
                                            },
                                        },
                                        t.i({ className: "ri-add-line" }),
                                        t.span({ className: "txt" }, "New record"),
                                    ),
                                    t.button(
                                        {
                                            hidden: () => !props.filter?.length,
                                            type: "button",
                                            className: "btn secondary expanded-lg",
                                            onclick() {
                                                props.filter = "";
                                                triggerOnchange();
                                            },
                                        },
                                        t.span({ className: "txt" }, "Clear search"),
                                    ),
                                );
                            }),
                        );
                    }

                    return data.records.map((record, i) => {
                        return t.tr(
                            {
                                rid: recordRid(record),
                                tabIndex: 0,
                                className: "handle",
                                // disable out-of-view detection for now as it can cause more issues than help
                                // onmount: (el) => {
                                //     el._intersectionObserver?.disconnect();
                                //     el._intersectionObserver = new IntersectionObserver((entries) => {
                                //         if (entries[0].intersectionRatio <= 0) {
                                //             el?.classList?.add("out-of-view")
                                //             return
                                //         }
                                //         el?.classList?.remove("out-of-view")
                                //     });
                                //     el._intersectionObserver.observe(el);
                                // },
                                // onunmount: (el) => {
                                //     el._intersectionObserver.disconnect();
                                //     el._intersectionObserver = null;
                                // },
                                onclick: (e) => {
                                    e.preventDefault();
                                    props.onselect(record);
                                },
                                onkeypress: (e) => {
                                    if (e.key == "Enter" || e.key == " ") {
                                        e.preventDefault();
                                        props.onselect(record);
                                    }
                                },
                            },
                            t.td(
                                {
                                    className: "col-bulk-select",
                                    onclick: (e) => e.stopPropagation(),
                                    onkeypress: (e) => e.stopPropagation(),
                                },
                                t.div(
                                    { className: "field" },
                                    t.input({
                                        type: "checkbox",
                                        id: () => uniqueId + record.id,
                                        checked: () => !!data.bulkSelected[record.id],
                                        onchange: (e) => {
                                            const bulkSelected = JSON.parse(JSON.stringify(data.bulkSelected));
                                            if (e.target.checked) {
                                                bulkSelected[record.id] = record;
                                            } else {
                                                delete bulkSelected[record.id];
                                            }

                                            // reassign to trigger the getter's Object.keys
                                            data.bulkSelected = bulkSelected;
                                        },
                                    }),
                                    t.label({ htmlFor: uniqueId + record.id }),
                                ),
                            ),
                            () => {
                                const columns = [];

                                // prepopulate outside of the tr to ensure that in case of a collection updated
                                // it will still reflect the column value change even if the record itself didn't get an update event
                                // (e.g. on field rename or single/multiple normalization)
                                const fields = props.collection?.fields || [];
                                for (const field of fields) {
                                    const viewFunc = app.fieldTypes[field.type]?.view;
                                    if (!viewFunc) {
                                        continue;
                                    }

                                    // superusers are always verified
                                    if (data.isSuperusersCollection && field.name == "verified") {
                                        continue;
                                    }

                                    columns.push(
                                        t.td(
                                            {
                                                "html-data-name": field.name,
                                                hidden: () => isFieldColumnHidden(field),
                                                className: `col-field-type-${field.type} col-field-name-${field.name}`,
                                            },
                                            () => {
                                                return viewFunc({
                                                    short: true,
                                                    get record() {
                                                        return record;
                                                    },
                                                    get field() {
                                                        return field;
                                                    },
                                                });
                                            },
                                        ),
                                    );
                                }

                                return columns;
                            },
                            // columns,
                            t.td(
                                { className: "col-meta" },
                                t.i({ className: "ri-arrow-right-line m-r-10", ariaHidden: true }),
                            ),
                        );
                    });
                },
                // load more btn
                t.tr(
                    { hidden: () => !data.canLoadMore },
                    t.td(
                        { colSpan: 99 },
                        t.button(
                            {
                                className: () =>
                                    `btn lg secondary load-more-btn ${data.isLoading ? "transparent loading" : ""}`,
                                disabled: () => data.isLoading,
                                onclick: () => loadRecords(),
                            },
                            t.span({ className: "txt" }, "Load more"),
                        ),
                    ),
                ),
            ),
        ),
        t.div(
            { className: "bulkbar-wrapper" },
            t.div(
                {
                    hidden: () => !data.totalSelected,
                    className: "bulkbar records-bulkbar",
                },
                t.span(
                    { className: "txt" },
                    "Selected ",
                    t.strong(null, () => data.totalSelected),
                    () => ` ${data.totalSelected == 1 ? "record" : "records"}`,
                ),
                t.button(
                    {
                        type: "button",
                        className: "btn sm secondary pill m-r-auto",
                        onclick: () => selectAll(false),
                    },
                    t.span({ className: "txt" }, "Reset"),
                ),
                () => {
                    if (props.collection?.type == "view") {
                        return;
                    }
                    return t.button(
                        {
                            type: "button",
                            className: "btn sm pill outline danger",
                            onclick: () => {
                                app.modals.confirm(
                                    "Do you really want to delete the selected records?",
                                    deleteSelected,
                                );
                            },
                        },
                        t.i({ className: "ri-delete-bin-7-line", ariaHidden: true }),
                        t.span({ className: "txt" }, "Delete"),
                    );
                },
                t.button(
                    {
                        type: "button",
                        className: "btn sm pill",
                        onclick: () => downloadSelected(),
                    },
                    t.i({ className: "ri-download-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "JSON"),
                ),
            ),
        ),
    );
};

function columnsDropdown(props, data) {
    const uniqueId = "cols_" + app.utils.randomString();

    const dropdown = t.div(
        { className: "dropdown sm nowrap records-list-columns-dropdown gap-0", popover: "auto" },
        () => {
            if (!props.collection?.fields) {
                return;
            }

            const result = [];

            const isAuth = props.collection.type == "auth";

            for (const field of props.collection.fields) {
                if (
                    field.primaryKey
                    || !app.fieldTypes[field.type].view
                    || (isAuth && field.name == "tokenKey")
                ) {
                    continue;
                }

                result.push(
                    t.div(
                        {
                            className: "dropdown-item",
                            onclick: (e) => {
                                // workaround for clicking on the padded area
                                e.target.querySelector("label")?.click();
                            },
                        },
                        t.div(
                            { className: "field" },
                            t.input({
                                type: "checkbox",
                                className: "switch sm",
                                id: () => uniqueId + field.name,
                                checked: () => {
                                    if (typeof data.columnsPreferences[field.id] != "undefined") {
                                        return !!data.columnsPreferences[field.id];
                                    }

                                    // no explicit preference
                                    return !field.hidden;
                                },
                                onchange: (e) => {
                                    data.columnsPreferences[field.id] = e.target.checked;
                                },
                            }),
                            t.label({ htmlFor: () => uniqueId + field.name }, field.name),
                        ),
                    ),
                );
            }

            return result;
        },
    );

    return t.button(
        {
            hidden: () => props.collection?.fields.length <= 1,
            type: "button",
            title: "Toggle columns",
            className: "btn sm secondary transparent circle",
            popoverTargetElement: dropdown,
        },
        t.i({ className: "ri-more-2-line", ariaHidden: true }),
        dropdown,
    );
}
