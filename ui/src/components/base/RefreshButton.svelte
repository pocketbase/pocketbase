<script>
    import { createEventDispatcher, onMount } from "svelte";
    import tooltip from "@/actions/tooltip";

    const dispatch = createEventDispatcher();

    let tooltipData = { text: "Refresh", position: "right" };
    export { tooltipData as tooltip };

    let classes = "";
    export { classes as class };

    let refreshTimeoutId = null;

    function refresh() {
        dispatch("refresh");

        // clear tooltip
        const oldTooltipData = tooltipData;
        tooltipData = null;

        clearTimeout(refreshTimeoutId);
        refreshTimeoutId = setTimeout(() => {
            refreshTimeoutId = null;
            tooltipData = oldTooltipData;
        }, 150);
    }

    onMount(() => {
        return () => clearTimeout(refreshTimeoutId);
    });
</script>

<button
    type="button"
    aria-label="Refresh"
    class="btn btn-transparent btn-circle {classes}"
    class:refreshing={refreshTimeoutId}
    use:tooltip={tooltipData}
    on:click={refresh}
>
    <i class="ri-refresh-line" />
</button>

<style>
    @keyframes refresh {
        100% {
            transform: rotate(180deg);
        }
    }
    .btn.refreshing i {
        animation: refresh 150ms ease-out;
    }
</style>
