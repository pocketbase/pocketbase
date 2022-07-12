<script>
    import { createEventDispatcher, onMount } from "svelte";
    import tooltip from "@/actions/tooltip";

    const dispatch = createEventDispatcher();

    let tooltipData = { text: "Refresh", position: "right" };
    export { tooltipData as tooltip };

    let refreshTimeoutId = null;

    function refresh() {
        dispatch("refresh");

        // clear tooltip
        const oldTooltipData = tooltipData;
        tooltipData = null;

        clearTimeout(refreshTimeoutId);
        refreshTimeoutId = setTimeout(() => {
            clearTimeout(refreshTimeoutId);
            refreshTimeoutId = null;
            tooltipData = oldTooltipData;
        }, 230);
    }

    onMount(() => {
        return () => clearTimeout(refreshTimeoutId);
    });
</script>

<button
    type="button"
    class="btn btn-secondary btn-circle"
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
        animation: refresh 200ms linear infinite;
    }
</style>
