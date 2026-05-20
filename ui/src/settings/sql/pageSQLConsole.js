import { settingsSidebar } from "../settingsSidebar";

const SQL_HISTORY_STORAGE_KEY = "pbSQLConsoleHistory";

export function pageSQLConsole(route) {
    app.store.title = "SQL console";

    const uniqueId = "sql_console_" + app.utils.randomString();
    const editorId = uniqueId + "editor";
    const requestKey = uniqueId + "executeSQL";
    const defaultMaxRows = 250;

    const pageData = store({
        askedForConfirmationAtLeastOnce: false,
        isExecuting: false,
        maxRows: defaultMaxRows,
        query: "",
        result: {},
        sort: {},
        errorMsg: "",
        executedHistory: app.utils.getLocalHistory(SQL_HISTORY_STORAGE_KEY, []),

        get sortedResultRows() {
            // no client-side sort
            if (typeof pageData.sort?.index == "undefined") {
                const rows = pageData.result?.rows || [];
                return pageData.maxRows >= rows.length ? rows : rows.slice(0, pageData.maxRows);
            }

            const isAsc = !!pageData.sort?.asc;

            const sorted = pageData.result?.rows?.toSorted((rowA, rowB) => {
                const valA = rowA[pageData.sort.index];
                const valB = rowB[pageData.sort.index];

                if (isAsc) {
                    if (valA == valB) {
                        return 0;
                    }
                    if (valA == null) {
                        return -1;
                    }
                    if (valB == null) {
                        return 1;
                    }
                    return valA.localeCompare(valB);
                }

                if (valA == valB) {
                    return 0;
                }
                if (valA == null) {
                    return 1;
                }
                if (valB == null) {
                    return -1;
                }
                return valB.localeCompare(valA);
            }) || [];

            return pageData.maxRows >= sorted.length ? sorted : sorted.slice(0, pageData.maxRows);
        },
        get totalRemainingRows() {
            return (pageData.result?.rows?.length << 0) - pageData.sortedResultRows.length;
        },
    });

    function needConfirmation(query) {
        if (!query?.trim()) {
            return false;
        }

        // if the hideControls option is enabled ask for confirmation
        // at least once no matter of the query
        if (
            !pageData.askedForConfirmationAtLeastOnce
            && !!app.store.settings?.meta?.hideControls
        ) {
            return true;
        }

        query = query?.replace(/[\s\;]/gm, " ").toUpperCase() + " ";

        return !![
            "INSERT ",
            "CREATE ",
            "UPDATE ",
            "DELETE ",
            "DROP ",
            "DETACH ",
            "PRAGMA ",
        ].find((p) => query.includes(p));
    }

    async function executeSQLWithConfirm() {
        if (!needConfirmation(pageData.query)) {
            return executeSQL();
        }

        pageData.askedForConfirmationAtLeastOnce = true;

        return app.modals.confirm(
            t.div(
                { className: "txt-center" },
                t.h6(
                    null,
                    "Be careful and continue only if you really know what you are doing because, depending on the query, the operation could break your application and may not be reversible.",
                ),
            ),
            () => executeSQL(),
            null,
            { yesButton: "Execute", noButton: "Cancel" },
        );
    }

    async function executeSQL() {
        pageData.isExecuting = true;
        pageData.maxRows = defaultMaxRows;
        pageData.result = {};
        pageData.sort = {};
        pageData.errorMsg = "";

        const query = pageData.query.trim();
        if (!query) {
            app.pb.cancelRequest(requestKey);
            pageData.isExecuting = false;
            return;
        }

        try {
            // @todo add method to JS SDK
            pageData.result = await app.pb.send("/api/sql", {
                method: "POST",
                body: { query },
                requestKey: requestKey,
            });

            addToHistory(query);

            pageData.isExecuting = false;
        } catch (err) {
            if (!err?.isAbort) {
                pageData.isExecuting = false;
                pageData.errorMsg = err?.response?.message || err?.message || "Failed to execute query.";
            }
        }
    }

    function removeFromHistory(query) {
        function looseNormalize(str) {
            return str.replace(/[\s\;]/gm, "").toUpperCase();
        }

        const normalized = looseNormalize(query);

        for (let i = pageData.executedHistory.length - 1; i >= 0; i--) {
            if (
                pageData.executedHistory[i] == query
                || looseNormalize(pageData.executedHistory[i]) == normalized
            ) {
                pageData.executedHistory.splice(i, 1);
            }
        }
    }

    function addToHistory(query) {
        removeFromHistory(query);

        pageData.executedHistory.unshift(pageData.query);
        if (pageData.executedHistory.length > 10) {
            pageData.executedHistory.splice(10);
        }
    }

    function downloadCSV() {
        if (!pageData.sortedResultRows.length) {
            return;
        }

        const data = [pageData.result.columns.map((c) => c.name)].concat(pageData.sortedResultRows);

        const name = "export_" + app.utils.toLocalDatetime(new Date()).replace(/[\-\:\. ]/g, "_") + ".csv";

        app.utils.downloadCSV(data, name);
    }

    function toggleSort(index) {
        if (pageData.sort?.index == index) {
            pageData.sort = {
                index: index,
                asc: !pageData.sort.asc,
            };
        } else {
            pageData.sort = {
                index: index,
                asc: true,
            };
        }
    }

    const watchers = [
        watch(() => JSON.stringify(pageData.executedHistory), (newVal, oldVal) => {
            if (typeof oldVal == "undefined") {
                return;
            }

            window.localStorage.setItem(SQL_HISTORY_STORAGE_KEY, newVal);
        }),
    ];

    return t.div(
        {
            pbEvent: "pageSQLConsole",
            className: "page",
            onunmount: () => {
                app.pb.cancelRequest(requestKey);
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
                t.div(
                    { className: "page-header-secondary-btns" },
                    t.button(
                        {
                            type: "button",
                            className: "btn circle transparent secondary",
                            ariaDescription: app.attrs.tooltip("Recently executed queries", "right"),
                            "html-popovertarget": "sql-console-history-dropdown",
                        },
                        t.i({ className: "ri-history-line", ariaHidden: true }),
                    ),
                    t.div(
                        {
                            id: "sql-console-history-dropdown",
                            className: () =>
                                `dropdown left sql-console-history-dropdown ${
                                    !pageData.executedHistory.length ? "no-items" : ""
                                }`,
                            popover: "auto",
                        },
                        () => {
                            if (!pageData.executedHistory.length) {
                                return t.span({ className: "txt txt-hint p-5" }, "No recently executed queries.");
                            }

                            return pageData.executedHistory.map((item) => {
                                return t.button(
                                    {
                                        role: "button",
                                        className: "dropdown-item",
                                        onclick: (e) => {
                                            e.target.closest(".dropdown").hidePopover();
                                            pageData.query = item;

                                            document.getElementById(editorId)?.click();
                                        },
                                    },
                                    t.span(
                                        { className: "query" },
                                        () => app.utils.truncate(item, 500),
                                    ),
                                    t.small(
                                        {
                                            role: "button",
                                            className: "remove-btn link-hint m-l-auto p-l-5 p-r-5",
                                            title: "Clear",
                                            onauxclick: (e) => {
                                                e.stopPropagation();
                                                return false;
                                            },
                                            onclick: (e) => {
                                                e.stopPropagation();
                                                removeFromHistory(item);
                                                return false;
                                            },
                                        },
                                        t.i({ className: "ri-close-line", ariaHidden: true }),
                                    ),
                                );
                            });
                        },
                    ),
                ),
                t.div(
                    { className: "page-header-primary-btns" },
                    t.button(
                        {
                            type: "button",
                            className: () => `btn expanded-lg ${pageData.isExecuting ? "loading" : ""}`,
                            disabled: () => pageData.isExecuting,
                            onclick: () => executeSQLWithConfirm(),
                        },
                        t.i({ className: "ri-play-large-line", ariaHidden: true }),
                        t.span({ className: "txt" }, "Execute"),
                    ),
                ),
            ),
            t.div(
                { className: "field sql-console-field" },
                app.components.codeEditor({
                    id: editorId,
                    language: "sql",
                    required: true,
                    name: "query",
                    placeholder: "e.g. EXPLAIN QUERY PLAN SELECT * from users WHERE verified=true",
                    value: () => pageData.query,
                    oninput: (val) => pageData.query = val,
                    onblur: (val) => pageData.query = val.trim(),
                }),
            ),
            t.div(
                { className: "flex field-help m-b-sm" },
                t.button(
                    {
                        type: "button",
                        className: "link-hint m-l-auto",
                        "html-popovertarget": uniqueId + "caveats_dropdown",
                    },
                    () => "SQL console caveats",
                ),
                t.div(
                    {
                        id: uniqueId + "caveats_dropdown",
                        className: "dropdown sm query-caveats-dropdown",
                        popover: "auto",
                    },
                    t.ul(
                        null,
                        t.li(null, "The returned rows are limited up to 1000."),
                        t.li(null, "The executed queries have a max timeout of 3 minutes."),
                        t.li(null, "The data is returned as byte strings without any additional formatting."),
                        t.li(null, "Multiple queries are supported but only the result of the last one is returned."),
                    ),
                ),
            ),
            // error alert
            t.div(
                {
                    hidden: () => pageData.isExecuting || !pageData.errorMsg,
                    className: "alert danger m-b-sm",
                },
                t.pre(null, () => pageData.errorMsg),
            ),
            // success alert
            t.div(
                {
                    hidden: () =>
                        pageData.isExecuting || pageData.errorMsg || pageData.result?.columns?.length
                        || app.utils.isEmpty(pageData.result),
                    className: "alert success m-b-sm",
                },
                t.p({ className: "txt-bold" }, "Query executed successfully!"),
                t.p(null, "Affected rows: ", () => pageData.result?.affectedRows || 0),
            ),
            // rows
            t.div(
                {
                    hidden: () => pageData.isExecuting || !pageData.result?.columns?.length,
                    className: "page-table-wrapper",
                },
                t.table(
                    { className: "sql-console-table responsive-table optimize" },
                    t.thead(
                        { className: "sticky" },
                        t.tr(null, () => {
                            return pageData.result?.columns?.map((col, i) => {
                                return t.th({
                                    textContent: col.name,
                                    className: () => {
                                        let classes = "sort-handle";

                                        if (pageData.sort?.index == i) {
                                            classes += pageData.sort.asc ? " asc" : " desc";
                                        }

                                        return classes;
                                    },
                                    onclick: () => toggleSort(i),
                                });
                            });
                        }),
                    ),
                    t.tbody(
                        null,
                        () => {
                            if (!pageData.sortedResultRows.length) {
                                return t.tr(
                                    null,
                                    t.td(
                                        { colSpan: pageData.result?.columns?.length || 1, className: "txt-center" },
                                        t.span({ className: "txt-hint" }, "No rows found."),
                                    ),
                                );
                            }

                            return pageData.sortedResultRows.map((rowData) => {
                                return t.tr(
                                    null,
                                    () => {
                                        return pageData.result?.columns?.map((col, j) => {
                                            const val = rowData[j];
                                            return t.td(
                                                { "html-data-name": col.name },
                                                val == null ? "NULL" : app.utils.truncate(val, 2000),
                                            );
                                        });
                                    },
                                );
                            });
                        },
                        // load more btn
                        t.tr(
                            {
                                hidden: () => (
                                    pageData.isExecuting
                                    || !pageData.result?.rows?.length
                                    || pageData.result.rows.length <= pageData.sortedResultRows.length
                                ),
                            },
                            t.td(
                                { colSpan: 99 },
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn lg secondary load-more-btn",
                                        onclick: () => {
                                            pageData.maxRows = pageData.result?.rows?.length || defaultMaxRows;
                                        },
                                    },
                                    t.span({
                                        className: "txt",
                                        textContent: () => `Load remaining (${pageData.totalRemainingRows})`,
                                    }),
                                ),
                            ),
                        ),
                    ),
                ),
            ),
            t.footer(
                { className: "page-footer" },
                t.span(
                    {
                        className: () => `exec-time ${pageData.isExecuting ? "faded" : ""}`,
                    },
                    "Time: ",
                    () => (pageData.result?.execTime || 0) + "ms",
                ),
                t.span(
                    {
                        hidden: () => !pageData.result?.columns?.length,
                        className: () => `total-count ${pageData.isExecuting ? "faded" : ""}`,
                    },
                    "Rows: ",
                    () => pageData.result?.rows?.length || 0,
                    () => {
                        if (!pageData.result?.rows?.length) {
                            return;
                        }

                        return [
                            " (",
                            t.span({
                                role: "button",
                                className: "link-hint",
                                textContent: "Export as CSV",
                                onclick: downloadCSV,
                            }),
                            ")",
                        ];
                    },
                ),
                app.components.credits(),
            ),
        ),
    );
}
