<script>
    import { onMount } from "svelte";
    import { scale } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import {
        Chart,
        BarController,
        BarElement,
        CategoryScale,
        LinearScale,
        TimeScale,
        Filler,
        Tooltip,
    } from "chart.js";
    import "chartjs-adapter-luxon";

    export let filter = "";
    export let presets = "";

    let chartCanvas;
    let chartInst;
    let chartData = [];
    let totalLogs = 0;
    let isLoading = false;

    $: if (typeof filter !== "undefined" || typeof presets !== "undefined") {
        load();
    }

    $: if (typeof chartData !== "undefined" && chartInst) {
        chartInst.data.datasets[0].data = chartData;
        chartInst.update();
    }

    export async function load() {
        isLoading = true;

        return ApiClient.logs
            .getStats({
                filter: [presets, CommonHelper.normalizeLogsFilter(filter)].filter(Boolean).join("&&"),
            })
            .then((result) => {
                resetData();

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
                    ApiClient.error(err, err?.status != 400); // silence filter errors
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

    onMount(() => {
        Chart.register(BarController, BarElement, CategoryScale, LinearScale, TimeScale, Filler, Tooltip);

        chartInst = new Chart(chartCanvas, {
            type: "bar",
            data: {
                datasets: [
                    {
                        label: "Total requests",
                        data: chartData,
                        backgroundColor: "#e34562",
                        maxBarThickness: 40,
                        borderRadius: 2,
                        minBarLength: 7,
                        hoverBackgroundColor: "#e34562",
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
                        // offset: false,
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
    <canvas bind:this={chartCanvas} class="chart-canvas" />
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
</style>
