<script>
	import { onMount } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import ForceGraph from 'force-graph';

    export let filter = "";
    export let presets = "";

    let chartCanvas;
    let graphInst;
    let chartData = {
        nodes: [],
        links: [],
    };
    let isLoading = false;

    $: if (typeof filter !== "undefined" || typeof presets !== "undefined") {
        load();
    }

    $: if (typeof chartData !== "undefined" && graphInst) {
        graphInst.graphData(chartData);
    }

	export async function load() {
        isLoading = true;

        const normalizedFilter = [presets, CommonHelper.normalizeLogsFilter(filter)]
            .filter(Boolean)
            .join("&&");
	

        return ApiClient.collections
            .getFullList()
            .then((result) => {
                resetData();

                result = CommonHelper.toArray(result);
                result.forEach((item) => {
                    chartData.nodes.push({
                        id: item.id,
                        name: item.name,
                        type: item.type,
                    });
                    const link = [];
                    const relationFields = item.schema.filter((field) => field.type === "relation");
                    relationFields.map((relationField) => {
                        link.push({
                            name: relationField.name,
                            source: item.id,
                            target: relationField.options.collectionId,
                        });
                    });
                    chartData.links.push(...link);
                });
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    resetData();
                    console.warn(err);
                    ApiClient.error(err, !normalizedFilter || err?.status != 400);
                }
            })
            .finally(() => {
                isLoading = false;
            });
    }

    function resetData() {
        chartData = {
            nodes: [],
            links: [],
        };
    }

	onMount(() => {
        graphInst = ForceGraph();
        graphInst(chartCanvas)
            .graphData(chartData)
            .nodeAutoColorBy('type')
            .nodeLabel('name')
            .onNodeDragEnd(node => {
                node.fx = node.x;
                node.fy = node.y;
            })
            .linkDirectionalParticles(1)
            .linkWidth(2)
            .linkCanvasObjectMode(() => 'after')
            .linkCanvasObject((link, ctx, globalScale) => {
                const fontSize = 12 / globalScale;
                ctx.font = `${fontSize}px Sans-Serif`;

                const start = link.source;
                const end = link.target;
                if (typeof start !== 'object' || typeof end !== 'object') {
                    return;
                }
                const textPos = Object.assign(...['x', 'y'].map(c => ({
                    [c]: start[c] + (end[c] - start[c]) / 3
                })));

                const relLink = { x: end.x - start.x, y: end.y - start.y };
                let textAngle = Math.atan2(relLink.y, relLink.x);
                if (textAngle > Math.PI / 2) {
                    textAngle = -(Math.PI - textAngle);
                }
                if (textAngle < -Math.PI / 2) {
                    textAngle = -(-Math.PI - textAngle);
                }

                const label = link.name;
                const textWidth = ctx.measureText(label).width;
                const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2);
                ctx.save();
                ctx.translate(textPos.x, textPos.y);
                ctx.rotate(textAngle);

                ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
                ctx.fillRect(- bckgDimensions[0] / 2, - bckgDimensions[1] / 2, ...bckgDimensions);

                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillStyle = 'darkgrey';
                ctx.fillText(label, 0, 0);
                ctx.restore();
            })
            .zoom(5);

    });
</script>

<div bind:this={chartCanvas} />
