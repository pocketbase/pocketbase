<script>
    import { onMount } from "svelte";
    import { scale } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
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

    export let filter = "";
    export let presets = "";

    let chartCanvas;
    let chartInst;
    let chartData = [];
    let totalRequests = 0;
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
            .getRequestsStats({
                filter: [presets, filter].filter(Boolean).join("&&"),
            })
            .then((result) => {
                resetData();
                for (let item of result) {
                    chartData.push({
                        x: new Date(item.date),
                        y: item.total,
                    });
                    totalRequests += item.total;
                }

                // add current time marker to the chart
                chartData.push({
                    x: new Date(),
                    y: undefined,
                });
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    resetData();
                    console.warn(err);
                    ApiClient.error(err, false);
                }
            })
            .finally(() => {
                isLoading = false;
            });
    }

    function resetData() {
        totalRequests = 0;
        chartData = [];
    }

    onMount(() => {
        Chart.register(LineElement, PointElement, LineController, LinearScale, TimeScale, Filler, Tooltip);

        chartInst = new Chart(chartCanvas, {
            type: "line",
            data: {
                datasets: [
                    {
                        label: "Total requests",
                        data: chartData,
                        borderColor: "#ef4565",
                        pointBackgroundColor: "#ef4565",
                        backgroundColor: "rgb(239,69,101,0.05)",
                        borderWidth: 2,
                        pointRadius: 1,
                        pointBorderWidth: 0,
                        fill: true,
                    },
                ],
            },
            options: {
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
                            borderColor: "#dee3e8",
                        },
                        ticks: {
                            precision: 0,
                            maxTicksLimit: 6,
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
                            borderColor: "#dee3e8",
                            color: (c) => (c.tick.major ? "#edf0f3" : ""),
                        },
                        ticks: {
                            maxTicksLimit: 15,
                            autoSkip: true,
                            maxRotation: 0,
                            major: {
                                enabled: true,
                            },
                            color: (c) => (c.tick.major ? "#16161a" : "#666f75"),
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
    {#if isLoading}
        <div class="chart-loader loader" transition:scale|local={{ duration: 150 }} />
    {/if}
    <canvas bind:this={chartCanvas} class="chart-canvas" style="height: 250px; width: 100%;" />
</div>

<div class="txt-hint m-t-xs txt-right">
    {#if isLoading}
        Loading...
    {:else}
        {totalRequests}
        {totalRequests === 1 ? "log" : "logs"}
    {/if}
</div>

<style>
    .chart-wrapper {
        position: relative;
        display: block;
        width: 100%;
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
</style>
