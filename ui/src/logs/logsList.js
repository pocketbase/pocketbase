import { logLevel } from "./logLevel";

const perPage = 50;

export function logsList(logsSettings) {
    const data = store({
        logs: [],
        lastLoadCount: 0,
        lastPage: 1,
        bulkSelected: {},
        get canLoadMore() {
            return data.lastLoadCount >= perPage;
        },
        get totalSelected() {
            return Object.keys(data.bulkSelected).length;
        },
        get areAllSelected() {
            return data.logs.length && data.logs.length == data.totalSelected;
        },
    });

    async function load(reset = false) {
        logsSettings.isListLoading = true;

        try {
            const page = reset ? 1 : data.lastPage + 1;

            const normalizedFilter = (logsSettings.presets || []).concat(
                app.utils.normalizeSearchFilter(logsSettings.filter, ["level", "message", "data"]),
            );

            if (logsSettings.zoom?.min && logsSettings.zoom?.max) {
                const minDate = new Date(logsSettings.zoom.min * 1000);
                minDate.setSeconds(0);
                minDate.setMilliseconds(0);
                const min = app.utils.toRFC3339Datetime(minDate);

                const maxDate = new Date(logsSettings.zoom.max * 1000);
                maxDate.setSeconds(59);
                maxDate.setMilliseconds(999);
                const max = app.utils.toRFC3339Datetime(maxDate);

                normalizedFilter.push(`created >= "${min}" && created <= "${max}"`);
            } else if (page > 1) {
                // minimize duplicates in case there were new logs that push the old ones to later pages
                normalizedFilter.push(`created <= "${data.logs[data.logs.length - 1].created}"`);
            }

            const result = await app.pb.logs.getList(page, perPage, {
                skipTotal: 1,
                sort: "-@rowid",
                requestKey: "logs_list",
                filter: normalizedFilter
                    .filter(Boolean)
                    .map((f) => "(" + f + ")")
                    .join("&&"),
            });

            if (result.page == 1) {
                data.logs = [];
                data.bulkSelected = {};
            }

            data.lastPage = result.page;
            data.lastLoadCount = result.items.length;

            for (let i = 0; i < result.items.length; i++) {
                app.utils.pushOrReplaceObject(data.logs, result.items[i]);

                // yield to main (with room to "breathe")
                if (i > 1 && i % 20 == 0) {
                    await new Promise((r) => setTimeout(r, 20));
                }
            }

            logsSettings.isListLoading = false;

            if (!logsSettings.isFirstLoadReady) {
                logsSettings.isFirstLoadReady = true;
            }
        } catch (err) {
            if (!err.isAbort) {
                logsSettings.isListLoading = false;
                app.checkApiError(err);
            }
        }
    }

    // note: allways assign a new object to trigger the getter's Object.keys
    function selectAll(state = true) {
        const selected = {};
        if (state) {
            for (let log of data.logs) {
                selected[log.id] = log;
            }
        }
        data.bulkSelected = selected;
    }

    function getLogPreviewKeys(log) {
        let keys = [];

        if (!log.data) {
            return keys;
        }

        if (log.data.type == "request") {
            const requestKeys = ["status", "execTime", "auth", "authId", "userIP"];
            for (let key of requestKeys) {
                if (typeof log.data[key] != "undefined") {
                    keys.push({ key });
                }
            }

            // add the referer if it is from a different source
            if (log.data.referer && !log.data.referer.includes(window.location.host)) {
                keys.push({ key: "referer" });
            }
        } else {
            // extract the first 6 keys (excluding the error and details)
            const allKeys = Object.keys(log.data);
            for (const key of allKeys) {
                if (key != "error" && key != "details" && keys.length < 6) {
                    keys.push({ key });
                }
            }
        }

        // ensure that the error and detail keys are last
        if (log.data.error) {
            keys.push({ key: "error", label: "danger" });
        }
        if (log.data.details) {
            keys.push({ key: "details", label: "warning" });
        }

        return keys;
    }

    const dateFilenameRegex = /[-:\. ]/gi;

    function downloadSelected() {
        // extract the bulk selected log objects sorted desc
        const selected = Object.values(data.bulkSelected).sort((a, b) => {
            if (a.created < b.created) {
                return 1;
            }

            if (a.created > b.created) {
                return -1;
            }

            return 0;
        });

        if (!selected.length) {
            return; // nothing to download
        }

        if (selected.length == 1) {
            return app.utils.downloadJSON(
                selected[0],
                "log_" + selected[0].created.replaceAll(dateFilenameRegex, "") + ".json",
            );
        }

        const to = selected[0].created.replaceAll(dateFilenameRegex, "");
        const from = selected[selected.length - 1].created.replaceAll(dateFilenameRegex, "");

        return app.utils.downloadJSON(selected, `${selected.length}_logs_${from}_to_${to}.json`);
    }

    const watchers = [];

    return t.div(
        {
            pbEvent: "logsList",
            className: "page-table-wrapper",
            onmount(el) {
                watchers.push(
                    // always reset zoom on filter or preset change
                    watch(
                        () => [logsSettings.filter, logsSettings.presets?.length],
                        () => {
                            logsSettings.zoom = {};
                        },
                    ),
                    watch(
                        () => [
                            logsSettings.reset,
                            logsSettings.filter,
                            logsSettings.presets?.length,
                            logsSettings.zoom?.min,
                            logsSettings.zoom?.max,
                        ],
                        () => {
                            load(true);

                            if (el) {
                                el.scrollTop = 0;
                            }
                        },
                    ),
                );
            },
            onunmount() {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.table(
            { className: () => `logs-table ${data.logs?.length > perPage ? "optimize" : ""}` },
            t.thead(
                null,
                t.tr(
                    null,
                    t.th(
                        { className: "col-bulk-select" },
                        t.div(
                            {
                                className: "field",
                                hidden: () => logsSettings.isLoading,
                            },
                            t.input({
                                id: "logs_select_all",
                                type: "checkbox",
                                checked: () => data.areAllSelected,
                                onchange: (e) => selectAll(e.target.checked),
                            }),
                            t.label({ htmlFor: "logs_select_all" }),
                        ),
                        t.span({
                            className: "loader",
                            hidden: () => !logsSettings.isLoading,
                        }),
                    ),
                    t.th(
                        { className: "col-field-name-level" },
                        t.div(
                            { className: "inline-flex gap-5" },
                            t.i({ className: "ri-bookmark-line", ariaHidden: true }),
                            t.span({ textContent: "Level" }),
                        ),
                    ),
                    t.th(
                        { className: "col-field-name-message" },
                        t.div(
                            { className: "inline-flex gap-5" },
                            t.i({ className: "ri-file-list-2-line", ariaHidden: true }),
                            t.span({ textContent: "Message" }),
                        ),
                    ),
                    t.th(
                        { className: "col-field-type-date col-field-name-created" },
                        t.div(
                            { className: "inline-flex gap-5" },
                            t.i({ className: "ri-calendar-line", ariaHidden: true }),
                            t.span({ textContent: "Created" }),
                        ),
                    ),
                    t.th({ className: "col-meta" }),
                ),
            ),
            t.tbody(
                null,
                () => {
                    if (!data.logs?.length) {
                        return t.tr(
                            null,
                            t.td(
                                { colSpan: 99 },
                                () => {
                                    if (logsSettings.isListLoading) {
                                        return t.span({ className: "skeleton-loader" });
                                    }

                                    return t.div(
                                        { className: "sticky-content txt-center txt-hint" },
                                        t.p({ className: "txt-bold" }, "No logs found."),
                                        t.button(
                                            {
                                                hidden: () =>
                                                    logsSettings.filter?.length
                                                    || app.utils.isEmpty(logsSettings.zoom),
                                                type: "button",
                                                className: "btn secondary expanded-lg",
                                                onclick() {
                                                    logsSettings.zoom = {};
                                                },
                                            },
                                            t.span({ className: "txt" }, "Reset zoom"),
                                        ),
                                        t.button(
                                            {
                                                hidden: () => !logsSettings.filter?.length,
                                                type: "button",
                                                className: "btn secondary expanded-lg",
                                                onclick() {
                                                    logsSettings.filter = "";
                                                },
                                            },
                                            t.span({ className: "txt" }, "Clear search"),
                                        ),
                                    );
                                },
                            ),
                        );
                    }

                    return data.logs.map((log) => {
                        return t.tr(
                            {
                                rid: log.id,
                                tabIndex: 0,
                                role: "button",
                                className: () => `handle ${log.data.type == "request" ? "log-request" : ""}`,
                                onclick: () => {
                                    logsSettings.activeLogIdOrModel = log;
                                },
                                onkeypress: (e) => {
                                    if (e.key == "Enter" || e.key == " ") {
                                        e.preventDefault();
                                        logsSettings.activeLogIdOrModel = log;
                                    }
                                },
                            },
                            () => {
                                return [
                                    t.td(
                                        {
                                            className: "col-bulk-select",
                                            onclick: (e) => e.stopPropagation(),
                                            onkeypress: (e) => e.stopPropagation(),
                                        },
                                        t.div(
                                            { className: "field" },
                                            t.input({
                                                id: "cb_" + log.id,
                                                type: "checkbox",
                                                checked: () => !!data.bulkSelected[log.id],
                                                onchange: (e) => {
                                                    const bulkSelected = JSON.parse(
                                                        JSON.stringify(data.bulkSelected),
                                                    );
                                                    if (e.target.checked) {
                                                        bulkSelected[log.id] = true;
                                                    } else {
                                                        delete bulkSelected[log.id];
                                                    }

                                                    // reassign to trigger the getter's Object.keys
                                                    data.bulkSelected = bulkSelected;
                                                },
                                            }),
                                            t.label({ htmlFor: "cb_" + log.id }),
                                        ),
                                    ),
                                    t.td({ className: "col-field-name-level" }, logLevel(log)),
                                    t.td(
                                        { className: "col-field-name-message" },
                                        t.div(
                                            { className: "content-primary" },
                                            () => app.utils.truncate(log.message, 1000),
                                        ),
                                        t.div(
                                            {
                                                className: "content-secondary",
                                            },
                                            () => {
                                                const labels = [];

                                                const keyItems = getLogPreviewKeys(log);
                                                for (const keyItem of keyItems) {
                                                    let value;
                                                    if (app.utils.logDataFormatters[keyItem.key]) {
                                                        value = app.utils.logDataFormatters[keyItem.key](log);
                                                    } else {
                                                        value = app.utils.stringifyValue(
                                                            log.data[keyItem.key],
                                                            "N/A",
                                                            80,
                                                        );
                                                    }

                                                    labels.push(
                                                        t.span(
                                                            {
                                                                className: `label sm ${keyItem.label || ""}`,
                                                            },
                                                            `${keyItem.key}: ${value}`,
                                                        ),
                                                    );
                                                }

                                                return labels;
                                            },
                                        ),
                                    ),
                                    t.td(
                                        {
                                            className: "col-field-type-date col-field-name-created",
                                        },
                                        app.components.formattedDate({
                                            value: () => log.created,
                                            short: false,
                                        }),
                                    ),
                                    t.td(
                                        { className: "col-meta" },
                                        t.i({ className: "ri-arrow-right-line", ariaHidden: true }),
                                    ),
                                ];
                            },
                        );
                    });
                },
                t.tr(
                    { hidden: () => !data.canLoadMore },
                    t.td(
                        { colSpan: 99 },
                        t.button(
                            {
                                className: () =>
                                    `btn lg secondary load-more-btn ${
                                        logsSettings.isListLoading ? "transparent loading" : ""
                                    }`,
                                disabled: () => logsSettings.isListLoading,
                                onclick: () => load(),
                            },
                            t.span({ className: "txt" }, "Load older"),
                        ),
                    ),
                ),
            ),
        ),
        t.div(
            { className: "bulkbar-wrapper" },
            t.div(
                {
                    hidden: () => data.totalSelected == 0,
                    className: "bulkbar logs-bulkbar",
                },
                t.span(
                    { className: "txt" },
                    "Selected ",
                    t.strong(null, () => data.totalSelected),
                    () => ` ${data.totalSelected == 1 ? "log" : "logs"}`,
                ),
                t.button(
                    {
                        type: "button",
                        className: "btn sm secondary pill m-r-auto",
                        onclick: () => selectAll(false),
                    },
                    t.span({ className: "txt" }, "Reset"),
                ),
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
}
