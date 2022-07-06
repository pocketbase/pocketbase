<script>
    export let id = "";

    let shortId = id;

    $: if (typeof id === "string" && id.length > 27) {
        shortId = id.substring(0, 5) + "..." + id.substring(id.length - 10);
    }
</script>

{#if id}
    <span class="label txt-base txt-mono" title={id}>{shortId}</span>
{:else}
    <span class="txt txt-hint">N/A</span>
{/if}
