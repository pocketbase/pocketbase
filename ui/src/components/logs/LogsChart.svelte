<script>
    import { onMount } from "svelte";
    import { scale } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import {
        Chart,
        LineElement,
        PointElement,
        LineController,
        LinearScale,
        TimeScale,
        Filler,
        Tooltip,
    } from "chart.js";
    import "chartjs-adapter-luxon";
    import zoomPlugin from "chartjs-plugin-zoom";

    export let filter = "";
    export let zoom = {};
    export let presets = "";

    let chartCanvas;
    let chartInst;
    let chartData = [];
    let totalLogs = 0;
    let isLoading = false;
    let isZoomedOrPanned = false;

    $: if (typeof filter !== "undefined" || typeof presets !== "undefined") {
        load();
    }

    $: if (typeof chartData !== "undefined" && chartInst) {
        chartInst.data.datasets[0].data = chartData;
        chartInst.update();
    }

    export async function load() {
        isLoading = true;

        const normalizedFilter = [presets, CommonHelper.normalizeLogsFilter(filter)]
            .filter(Boolean)
            .map((f) => "(" + f + ")")
            .join("&&");

        return ApiClient.logs
            .getStats({
                filter: normalizedFilter,
            })
            .then((result) => {
                resetData();

                result = CommonHelper.toArray(result);

                for (let item of result) {
                    chartData.push({
                        x: new Date(item.date),
                        y: item.total,
                    });
                    totalLogs += item.total;
                }
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    resetData();
                    console.warn(err);
                    ApiClient.error(err, !normalizedFilter || err?.status != 400); // silence filter errors
                }
            })
            .finally(() => {
                isLoading = false;
            });
    }

    function resetData() {
        chartData = [];
        totalLogs = 0;
    }

    function resetZoom() {
        chartInst?.resetZoom();
    }

    onMount(() => {
        Chart.register(LineElement, PointElement, LineController, LinearScale, TimeScale, Filler, Tooltip);
        Chart.register(zoomPlugin);

        chartInst = new Chart(chartCanvas, {
            type: "line",
            data: {
                datasets: [
                    {
                        label: "Total requests",
                        data: chartData,
                        borderColor: "#e34562",
                        pointBackgroundColor: "#e34562",
                        backgroundColor: "rgb(239,69,101,0.05)",
                        borderWidth: 2,
                        pointRadius: 1,
                        pointBorderWidth: 0,
                        fill: true,
                    },
                ],
            },
            options: {
                resizeDelay: 250,
                maintainAspectRatio: false,
                animation: false,
                interaction: {
                    intersect: false,
                    mode: "index",
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        grid: {
                            color: "#edf0f3",
                        },
                        border: {
                            color: "#e4e9ec",
                        },
                        ticks: {
                            precision: 0,
                            maxTicksLimit: 4,
                            autoSkip: true,
                            color: "#666f75",
                        },
                    },
                    x: {
                        type: "time",
                        time: {
                            unit: "hour",
                            tooltipFormat: "DD h a",
                        },
                        grid: {
                            color: (c) => (c.tick?.major ? "#edf0f3" : ""),
                        },
                        color: "#e4e9ec",
                        ticks: {
                            maxTicksLimit: 15,
                            autoSkip: true,
                            maxRotation: 0,
                            major: {
                                enabled: true,
                            },
                            color: (c) => (c.tick?.major ? "#16161a" : "#666f75"),
                        },
                    },
                },
                plugins: {
                    legend: {
                        display: false,
                    },
                    zoom: {
                        enabled: true,
                        zoom: {
                            mode: "x",
                            pinch: {
                                enabled: true,
                            },
                            drag: {
                                enabled: true,
                                backgroundColor: "rgba(255, 99, 132, 0.2)",
                                borderWidth: 0,
                                threshold: 10,
                            },
                            limits: {
                                x: { minRange: 100000000 },
                                y: { minRange: 100000000 },
                            },
                            onZoomComplete: ({ chart }) => {
                                isZoomedOrPanned = chart.isZoomedOrPanned();
                                if (!isZoomedOrPanned) {
                                    if (zoom.min || zoom.max) {
                                        zoom = {}; // reset
                                    }
                                } else {
                                    // trim minutes and seconds since the statistic is hourly based
                                    zoom.min =
                                        CommonHelper.formatToUTCDate(chart.scales.x.min, "yyyy-MM-dd HH") +
                                        ":00:00.000Z";
                                    zoom.max =
                                        CommonHelper.formatToUTCDate(chart.scales.x.max, "yyyy-MM-dd HH") +
                                        ":59:59.999Z";
                                }
                            },
                        },
                    },
                },
            },
        });

        return () => chartInst?.destroy();
    });
</script>

<div class="chart-wrapper" class:loading={isLoading}>
    <div class="total-logs entrance-right" class:hidden={isLoading}>
        Found {totalLogs}
        {totalLogs == 1 ? "log" : "logs"}
    </div>

    {#if isLoading}
        <div class="chart-loader loader" transition:scale={{ duration: 150 }} />
    {/if}

    <canvas bind:this={chartCanvas} class="chart-canvas" on:dblclick={resetZoom} />

    {#if isZoomedOrPanned}
        <button type="button" class="btn btn-secondary btn-sm btn-chart-zoom" on:click={resetZoom}>
            Reset zoom
        </button>
    {/if}
</div>

<style>
    .chart-wrapper {
        position: relative;
        display: block;
        width: 100%;
        height: 170px;
    }
    .chart-wrapper.loading .chart-canvas {
        pointer-events: none;
        opacity: 0.5;
    }
    .chart-loader {
        position: absolute;
        z-index: 999;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
    .total-logs {
        position: absolute;
        right: 0;
        top: -50px;
        font-size: var(--smFontSize);
        color: var(--txtHintColor);
    }
    .btn-chart-zoom {
        position: absolute;
        right: 10px;
        top: 20px;
    }
</style>
