import { logsChart } from "./logsChart";
import { logsList } from "./logsList";

export function pageLogs(route) {
    app.store.title = "Logs";

    const LOG_QUERY_KEY = "logId";
    const FILTER_QUERY_KEY = "filter";
    const SUPERUSER_REQUESTS_QUERY_KEY = "superuserRequests";
    const SUPERUSER_REQUESTS_STORAGE_KEY = "pbLogSuperuserRequests";

    const withoutSuperusersPresets = "data.auth!='_superusers'";

    const queryLogId = route.query[LOG_QUERY_KEY]?.[0] || "";

    const queryFilter = route.query[FILTER_QUERY_KEY]?.[0] || "";

    const querySuperuserRequests = !!(route.query[SUPERUSER_REQUESTS_QUERY_KEY]?.[0] << 0)
        || !!(window.localStorage.getItem(SUPERUSER_REQUESTS_STORAGE_KEY) << 0);

    const logsSettings = store({
        reset: null,
        isChartLoading: false,
        isListLoading: false,
        isFirstLoadReady: false, // used by the chart to show itself after the first list load to minimize flickering
        zoom: {}, // only unidirectional from within the chart
        presets: querySuperuserRequests ? [] : [withoutSuperusersPresets],
        filter: queryFilter,
        totalFound: null,
        activeLogIdOrModel: queryLogId,
        get hasIncludeRequestsBySuperusers() {
            return !logsSettings.presets.includes(withoutSuperusersPresets);
        },
        get isLoading() {
            return logsSettings.isListLoading || logsSettings.isChartLoading;
        },
    });

    function getLogId(logIdOrModel) {
        if (!logIdOrModel) {
            return null;
        }

        return typeof logIdOrModel === "string" ? logIdOrModel : logIdOrModel?.id;
    }

    function refreshLogsList() {
        logsSettings.reset = Date.now();
    }

    const watchers = [];

    return [
        t.div(
            { pbEvent: "logsChartContainer", className: "logs-chart-container accent-surface" },
            logsChart(logsSettings),
        ),
        t.div(
            {
                pbEvent: "pageLogs",
                className: "page page-logs",
                onmount(el) {
                    watchers.push(
                        watch(() => {
                            app.utils.replaceHashQueryParams({
                                [FILTER_QUERY_KEY]: logsSettings.filter,
                            });
                        }),
                        watch(() => {
                            const superuserRequests = logsSettings.hasIncludeRequestsBySuperusers ? 1 : 0;

                            app.utils.replaceHashQueryParams({
                                [SUPERUSER_REQUESTS_QUERY_KEY]: superuserRequests,
                            });

                            window.localStorage.setItem(SUPERUSER_REQUESTS_STORAGE_KEY, superuserRequests);
                        }),
                        watch(() => logsSettings.activeLogIdOrModel, () => {
                            app.utils.replaceHashQueryParams({
                                [LOG_QUERY_KEY]: getLogId(logsSettings.activeLogIdOrModel),
                            });

                            if (logsSettings.activeLogIdOrModel) {
                                // force close any previous modal
                                app.modals.close(null, true);

                                app.modals.openLogPreview(logsSettings.activeLogIdOrModel, {
                                    onafterclose: () => {
                                        logsSettings.activeLogIdOrModel = null;
                                    },
                                });
                            }
                        }),
                    );
                },
                onunmount(el) {
                    clearTimeout(el._chartTiemoutId);
                    watchers.forEach((w) => w?.unwatch());
                },
            },
            t.div(
                { className: "page-content full-height" },
                t.header(
                    { className: "page-header" },
                    t.nav({ className: "breadcrumbs" }, t.div(null, "Logs")),
                    t.div(
                        { className: "inline-flex gap-sm" },
                        t.button(
                            {
                                className: "btn circle transparent secondary tooltip-right",
                                ariaDescription: app.attrs.tooltip("Logs settings"),
                                onclick: () =>
                                    app.modals.openLogsSettings({
                                        onsave: () => refreshLogsList(),
                                    }),
                            },
                            t.i({ className: "ri-settings-3-line" }),
                        ),
                        app.components.refreshButton({
                            onclick: refreshLogsList,
                        }),
                    ),
                    app.components.searchbar({
                        className: "logs-searchbar",
                        historyKey: "pbLogsSearchHistory",
                        placeholder: "Search term or filter like `level > 0`",
                        value: () => logsSettings.filter || "",
                        onsubmit: (val) => logsSettings.filter = val,
                        autocomplete: [
                            "id",
                            "level",
                            "message",
                            "created",
                            { value: "data.", label: "data.*" },
                        ],
                    }),
                    t.div(
                        { className: "meta m-l-auto" },
                        t.div(
                            { className: "field logs-include-superuser-requests" },
                            t.input({
                                type: "checkbox",
                                id: "logs_checkbox",
                                className: "switch sm",
                                checked: () => logsSettings.hasIncludeRequestsBySuperusers,
                                onchange: (e) => {
                                    if (e.target.checked) {
                                        app.utils.removeByValue(logsSettings.presets, withoutSuperusersPresets);
                                    } else {
                                        app.utils.pushUnique(logsSettings.presets, withoutSuperusersPresets);
                                    }
                                },
                            }),
                            t.label(
                                { htmlFor: "logs_checkbox" },
                                t.small({ className: "txt" }, "Include requests by superusers"),
                            ),
                        ),
                    ),
                ),
                logsList(logsSettings),
                t.footer(
                    { className: "page-footer" },
                    t.span(
                        { className: "txt total-logs" },
                        "Total: ",
                        () => {
                            if (logsSettings.totalFound == null) {
                                return "...";
                            }
                            return logsSettings.totalFound;
                        },
                    ),
                    app.components.credits(),
                ),
            ),
        ),
    ];
}
