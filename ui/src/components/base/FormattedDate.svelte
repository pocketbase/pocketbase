<script>
    export let date = "";

    // note: manual trim the ms without converting to DateTime
    // to help improving the rendering performance in large data sets
    let shortDate = date.length > 19 ? date.substring(0, 19) : date;
</script>

{#if date}
    <span class="txt">{shortDate} UTC</span>
{:else}
    <span class="txt txt-hint">N/A</span>
{/if}
