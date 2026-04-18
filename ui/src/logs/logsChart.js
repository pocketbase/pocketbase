function getChartHeight() {
    return document.body.clientWidth > 600 && document.body.clientHeight >= 720 ? 150 : 100;
}

export function logsChart(logsSettings) {
    const data = store({
        stats: [],
    });

    async function loadAndInit(el) {
        if (!el) {
            return;
        }

        logsSettings.isChartLoading = true;

        try {
            const normalizedFilter = (logsSettings.presets || []).concat(
                app.utils.normalizeSearchFilter(logsSettings.filter, ["level", "message", "data"]),
            );

            const stats = await app.pb.logs.getStats({
                filter: normalizedFilter
                    .filter(Boolean)
                    .map((f) => "(" + f + ")")
                    .join("&&"),
            });

            const totalStats = stats.length;

            const timestamps = [];
            const totals = [];

            logsSettings.totalFound = 0;

            for (let i = 0; i < totalStats; i++) {
                const unix = new Date(stats[i].date.replace(" ", "T")).getTime() / 1000;
                timestamps.push(unix);
                totals.push(stats[i].total);

                logsSettings.totalFound += stats[i].total;

                // if between the current and next point there is more than 1h difference
                // insert a 0 point for a flat 1h line in the stepped chart visualization
                if (stats[i + 1]?.date) {
                    const unixNext = new Date(stats[i + 1].date.replace(" ", "T")).getTime() / 1000;
                    if (unixNext - unix > 3600) {
                        timestamps.push(unix + 3600);
                        totals.push(0);
                    }
                } else if (i + 1 == totalStats) {
                    timestamps.push(unix + 3600);
                    totals.push(0);
                }
            }

            data.stats = stats;
            logsSettings.isChartLoading = false;

            initChart(el, [timestamps, totals], logsSettings);
        } catch (err) {
            if (!err?.isAbort) {
                logsSettings.isChartLoading = false;
                // only log to avoid showing multiple errors with the logs listing
                // app.checkApiError(err)
                console.warn("failed to load logs chart:", err);
            }
        }
    }

    const watchers = [];

    return t.div(
        {
            pbEvent: "logsChart",
            className: () =>
                [
                    "logs-chart",
                    logsSettings.isChartLoading ? "loading" : null,
                    logsSettings.zoom?.min && logsSettings.zoom?.max ? "zoomed" : "",
                    !data.stats.length || !logsSettings.isFirstLoadReady ? "nodata" : null,
                ].filter(Boolean).join(" "),
            onmount: (el) => {
                // init and refresh chart
                watchers.push(
                    watch(
                        () => [logsSettings.reset, logsSettings.filter, logsSettings.presets?.length],
                        () => loadAndInit(el),
                    ),
                );

                el._resizeChartFunc = () => {
                    clearTimeout(el._resizeTimeoutId);
                    el._resizeTimeoutId = setTimeout(() => {
                        if (!el?._uplot) {
                            return;
                        }

                        el._uplot.setSize({
                            width: el.clientWidth,
                            height: getChartHeight(),
                        });
                    }, 100);
                };
                window.addEventListener("resize", el._resizeChartFunc);
            },
            onunmount: (el) => {
                watchers.forEach((w) => w?.unwatch());

                el._uplot?.destroy();

                if (el._resizeChartFunc) {
                    clearTimeout(el._resizeTimeoutId);
                    window.removeEventListener("resize", el._resizeChartFunc);
                    el._resizeChartFunc = null;
                    el._resizeTimeoutId = null;
                }
            },
        },
        t.button(
            {
                type: "button",
                className: () =>
                    `logs-reset-zoom-ctrl ${logsSettings.zoom?.min && logsSettings.zoom?.max ? "" : "hidden"}`,
                onclick() {
                    logsSettings.zoom = {};
                },
            },
            t.div({ className: "content-primary" }, "Reset zoom"),
            t.div({ className: "content-secondary" }, "(drag the timeline to pan)"),
        ),
        t.span({
            hidden: () => !logsSettings.isChartLoading,
            className: () => "loader logs-chart-loader",
        }),
    );
}

function resetChartZoom(chart) {
    const data = chart?.data?.[0];
    if (!data) {
        return;
    }

    chart.setScale("x", {
        min: data[0],
        max: data[data.length - 1],
    });
}

function initChart(el, dataPoints, logsSettings) {
    const computedStyle = window.getComputedStyle(el);
    const gridTxtColor = computedStyle.getPropertyValue("--surfaceTxtHintColor");
    const gridColor = computedStyle.getPropertyValue("--surfaceAlt1Color");
    const strokeColor = computedStyle.getPropertyValue("--surfaceAlt4Color");
    const fillColor = computedStyle.getPropertyValue("--surfaceAlt2Color");

    const opts = {
        width: el.clientWidth,
        height: getChartHeight(),
        legend: {
            show: false,
        },
        cursor: {
            x: false,
            y: false,
        },
        scales: {
            x: {
                range: (self, newMin, newMax) => {
                    // disallow pan beyond the dataPoints edges
                    if (newMin < dataPoints[0][0] || newMax > dataPoints[0][dataPoints[0].length - 1]) {
                        return [self.scales.x.min, self.scales.x.max];
                    }

                    // allow zoom
                    return [newMin, newMax];
                },
            },
            y: {
                range: {
                    min: {
                        pad: 0.2,
                        soft: 0,
                        mode: 1,
                    },
                    max: {
                        pad: 0.2,
                        soft: 0,
                        mode: 2,
                    },
                },
            },
        },
        series: [
            {},
            {
                paths: uPlot.paths.stepped({ align: 1 }),
                points: {
                    show: false,
                    size: 1,
                },
                width: 2,
                fill: fillColor,
                stroke: strokeColor,
            },
        ],
        axes: [
            {
                show: true,
                size: 35,
                lineGap: 1,
                space: 60,
                stroke: gridTxtColor,
                incrs: [
                    // minutes
                    60 * 5,
                    60 * 10,
                    60 * 15,
                    60 * 30,
                    // hours
                    3600,
                    3600 * 2,
                    3600 * 3,
                    3600 * 4,
                    3600 * 6,
                    3600 * 12,
                    // days
                    86400,
                    86400 * 2,
                ],
                // dprint-ignore
                values: [
                    // incr  default         year            month  day             hour  min   sec   mode
                    [3600,   "{h}{aa}",      "\n{MMM} {DD}", null,  "\n{MMM} {DD}", null, null, null, 1],
                    [60,     "{h}:{mm}{aa}", "\n{MMM} {DD}", null,  "\n{MMM} {DD}", null, null, null, 1],
                ],
                grid: {
                    show: true,
                    stroke: gridColor,
                    width: 1,
                },
                ticks: {
                    show: true,
                    stroke: gridColor,
                    width: 1,
                    size: 5,
                },
            },
            {
                show: true,
                stroke: gridTxtColor,
                grid: {
                    show: true,
                    stroke: gridColor,
                    width: 1,
                },
                ticks: {
                    show: true,
                    stroke: gridColor,
                    width: 1,
                    size: 5,
                },
            },
        ],
        plugins: [tooltipsPlugin(), zoomSetPlugin(logsSettings), xPanPlugin(logsSettings)],
    };

    el._uplot?.destroy();
    el._uplot = new uPlot(opts, dataPoints, el);
}

// based on view-source:https://leeoniya.github.io/uPlot/demos/zoom-fetch.html
function zoomSetPlugin(logsSettings) {
    let zoomWatcher;

    return {
        hooks: {
            init: (u) => {
                u.over.ondblclick = (e) => {
                    logsSettings.zoom = {};
                };

                zoomWatcher = watch(() => {
                    if (!logsSettings.zoom?.min || !logsSettings.zoom?.max) {
                        resetChartZoom(u);
                    } else {
                        u.setScale("x", {
                            min: logsSettings.zoom.min,
                            max: logsSettings.zoom.max,
                        });
                    }
                });
            },
            destroy: (u) => {
                zoomWatcher?.unwatch();
            },
            setSelect: (u) => {
                if (u.select.width > 0) {
                    logsSettings.zoom = {
                        min: u.posToVal(u.select.left, "x"),
                        max: u.posToVal(u.select.left + u.select.width, "x"),
                    };
                }
            },
        },
    };
}

// based on https://leeoniya.github.io/uPlot/demos/y-scale-drag.html
function xPanPlugin(logsSettings) {
    let debounceTimeout;

    return {
        hooks: {
            init: (u) => {
                let axisElems = u.root.querySelectorAll(".u-axis");
                if (!axisElems.length) {
                    console.warn("xPanPlugin requires x axis to be defined");
                    return;
                }

                axisElems[0].addEventListener("mousedown", (e) => {
                    if (!logsSettings.zoom?.min) {
                        return; // no zoom
                    }

                    let x0 = e.clientX;
                    let scale = u.scales.x;
                    let { min, max } = scale;
                    let dim = u.bbox.width;
                    let unitsPerPx = (max - min) / (dim / uPlot.pxRatio);

                    const mousemoveFunc = (e) => {
                        let diff = x0 - e.clientX;
                        let shiftBy = diff * unitsPerPx;

                        u.setScale("x", {
                            min: min + shiftBy,
                            max: max + shiftBy,
                        });

                        clearTimeout(debounceTimeout);
                        debounceTimeout = setTimeout(() => {
                            if (u?.scales?.x) {
                                logsSettings.zoom = {
                                    min: u.scales.x.min,
                                    max: u.scales.x.max,
                                };
                            }
                        }, 100);
                    };

                    const mouseupFunc = (e) => {
                        document.removeEventListener("mousemove", mousemoveFunc);
                        document.removeEventListener("mouseup", mouseupFunc);
                    };

                    document.addEventListener("mousemove", mousemoveFunc);
                    document.addEventListener("mouseup", mouseupFunc);
                });
            },
            destroy: (u) => {
                if (debounceTimeout) {
                    clearTimeout(debounceTimeout);
                }
            },
        },
    };
}

// based on https://leeoniya.github.io/uPlot/demos/tooltips.html
function tooltipsPlugin(logsSettings) {
    let tooltip;

    return {
        hooks: {
            init: (u) => {
                const over = u.over;

                tooltip = store({
                    date: "",
                    total: 0,
                    left: 0,
                    top: 0,
                    show: false,
                });

                const tooltipEl = t.div(
                    {
                        className: () => `chart-tooltip ${tooltip.show ? "" : "hidden"}`,
                        onmount(el) {
                            el._positionWatcher?.unwatch();
                            el._positionWatcher = watch(
                                () => [tooltip.left, tooltip.top],
                                () => {
                                    if (!el) {
                                        return;
                                    }

                                    const rect = el.getBoundingClientRect();

                                    let left = tooltip.left;
                                    if (left < 0) {
                                        left = 0;
                                    } else if (left + rect.width > over.clientWidth) {
                                        left = over.clientWidth - rect.width;
                                    }
                                    el.style.left = left + "px";

                                    const vOffset = 5;
                                    let top = tooltip.top - rect.height - vOffset;
                                    if (top < 0) {
                                        top = tooltip.top + vOffset;
                                        if (top + rect.high > over.clientHeight) {
                                            top = over.clientHeight - rect.height;
                                        }
                                    }
                                    el.style.top = top + "px";
                                },
                            );
                        },
                        onunmount(el) {
                            el._positionWatcher?.unwatch();
                        },
                    },
                    t.div(
                        { className: "content-primary" },
                        () => `${tooltip.total} ${tooltip.total == 1 ? "request" : "requests"}`,
                    ),
                    t.div({ className: "content-secondary" }, () => tooltip.date),
                );
                over.appendChild(tooltipEl);

                over.addEventListener("mouseleave", () => {
                    if (tooltip) {
                        tooltip.show = false;
                    }
                });
            },
            destroy: () => {
                tooltip.show = false;
            },
            setCursor: (u) => {
                if (!tooltip) {
                    return;
                }

                const xVal = u.data[0][u.cursor.idx] || 0;
                const yVal = u.data[1][u.cursor.idx] || 0;

                // skip zero points
                if (xVal == 0 || yVal == 0) {
                    tooltip.show = false;
                    return;
                }

                tooltip.show = true;
                tooltip.total = yVal;

                const xDateStart = new Date(xVal * 1000);
                const xDateEnd = new Date(xVal * 1000 + 3600000); // all stats are hourly based
                const monthName = xDateStart.toLocaleString("default", { month: "short" });
                const day = xDateStart.getDate().toString().padStart(2, "0");
                tooltip.date = `${monthName} ${day} ${dateHour(xDateStart)}-${dateHour(xDateEnd)}`;

                tooltip.left = Math.round(u.valToPos(xVal, "x"));
                tooltip.top = Math.round(u.valToPos(yVal, "y"));
            },
        },
    };
}

function dateHour(date) {
    let hours = date.getHours();
    let ampm = hours >= 12 ? "pm" : "am";

    // normalize
    hours = hours % 12 || 12;

    return hours + ampm;
}
