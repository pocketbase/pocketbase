const TEST_REQUEST_KEY = "test_view_query";

export function collectionViewQueryTab(upsertData) {
    const uniqueId = "query_" + app.utils.randomString();

    // dprint-ignore
    const autocomplete = [
        "SELECT", "FROM", "WHERE", "LEFT JOIN", "INNER JOIN", "ON",
        "GROUP BY", "HAVING", "ORDER BY", "LIMIT", "OFFSET", "AS",
        "WITH", "NOT", "IN", "EXISTS", "LIKE", "CAST",
    ];

    const local = store({
        testRecords: [],
        testError: "",
        isTesting: false,
    });

    async function dryRunViewQuery(query) {
        local.isTesting = true;

        local.testRecords = [];

        // reset form errors related to the query
        if (app.store.errors?.viewQuery || app.store.errors?.fields) {
            delete app.store.errors.viewQuery;
            delete app.store.errors.fields;
        }

        if (!query) {
            local.testError = "";
            local.isTesting = false;
            return;
        }

        try {
            // @todo replace with SDK method
            const result = await app.pb.send("/api/collections/meta/dry-run-view", {
                method: "POST",
                body: { "query": query },
                requestKey: TEST_REQUEST_KEY,
            });

            if (upsertData.collection?.id) {
                // replace the collection meta fields
                local.testRecords = result.sample.map((r) => {
                    r.collectionId = upsertData.collection?.id;
                    r.collectionName = upsertData.collection?.name;
                    return r;
                });
            } else {
                local.testRecords = result.sample;
            }

            local.testError = "";
            local.isTesting = false;
        } catch (err) {
            if (!err.isAbort) {
                local.testError = err.message || "Invalid query.";
                local.isTesting = false;
            }
        }
    }

    let testDebounceId;

    const watchers = [
        watch(() => upsertData.collection?.viewQuery, (newQuery) => {
            clearTimeout(testDebounceId);
            testDebounceId = setTimeout(() => dryRunViewQuery(newQuery), 200);
        }),
    ];

    return t.div(
        {
            pbEvent: "collectionViewQueryTabContent",
            className: "collection-tab-content collection-view-query-tab-content",
            onunmount: () => {
                clearTimeout(testDebounceId);
                app.pb.cancelRequest(TEST_REQUEST_KEY);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: "grid" },
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "txt-right txt-sm m-b-10" },
                    t.button(
                        {
                            type: "button",
                            className: "txt-bold link-hint",
                            "html-popovertarget": uniqueId + "caveats_dropdown",
                        },
                        () => "Query caveats",
                    ),
                ),
                t.div(
                    {
                        id: uniqueId + "caveats_dropdown",
                        className: "dropdown sm query-caveats-dropdown",
                        popover: "auto",
                    },
                    t.ul(
                        null,
                        t.li(null, "Wildcard columns (*) are not supported."),
                        t.li(
                            null,
                            "The query must have a unique ",
                            t.code(null, "id"),
                            " column.",
                            t.br(),
                            "If your query doesn't have a suitable one, you can use the universal ",
                            t.code(null, "(ROW_NUMBER() OVER()) as id"),
                            ".",
                        ),
                        t.li(
                            null,
                            "Expressions must be aliased with a valid formatted field name, e.g. ",
                            t.code(null, "MAX(balance) as maxBalance"),
                            ".",
                        ),
                        t.li(
                            null,
                            "Combined/multi-spaced expressions must be wrapped in parenthesis, e.g. ",
                            t.code(null, "(MAX(balance) + 1) as maxBalance"),
                            ".",
                        ),
                        t.li(
                            null,
                            "UNION expressions are supported but the entire query must be wrapped in parenthesis.",
                        ),
                    ),
                ),
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: uniqueId + ".viewQuery" },
                        t.span({ className: "txt" }, "Select query"),
                        t.span(
                            {
                                hidden: () => !local.testError,
                                className: "query-state",
                                ariaDescription: app.attrs.tooltip("Invalid query", "left"),
                            },
                            t.i({ className: "ri-error-warning-fill txt-danger", ariaHidden: true }),
                        ),
                        t.span(
                            {
                                hidden: () => !!local.testError,
                                className: "query-state",
                                ariaDescription: app.attrs.tooltip("Valid query", "left"),
                            },
                            t.i({ className: "ri-checkbox-circle-fill txt-success", ariaHidden: true }),
                        ),
                    ),
                    app.components.codeEditor({
                        id: uniqueId + ".viewQuery",
                        name: "viewQuery",
                        language: "sql",
                        required: true,
                        autocomplete: autocomplete,
                        className: "inline-error",
                        value: () => upsertData.collection.viewQuery || "",
                        oninput: (newVal) => {
                            upsertData.collection.viewQuery = newVal;
                        },
                    }),
                ),
            ),
            t.div(
                { className: "col-12" },
                t.p(
                    { className: "txt-sm txt-bold" },
                    "Sample output:",
                ),
                t.div(
                    { className: "view-query-sample-wrapper" },
                    t.span({ hidden: () => !local.isTesting, className: "loader sm" }),
                    app.components.codeBlock({
                        language: () => local.testError ? "plain" : "js",
                        className: () => `view-query-sample ${local.testError ? "txt-danger" : ""}`,
                        value: () => {
                            if (local.testRecords?.length) {
                                return JSON.stringify(local.testRecords, null, 2);
                            }

                            return local.testError || "N/A";
                        },
                    }),
                ),
            ),
        ),
    );
}
