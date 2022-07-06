<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";

    export let date = "";
</script>

{#if date}
    <span class="txt" use:tooltip={CommonHelper.formatToLocalDate(date) + " Local"}>
        {CommonHelper.formatToUTCDate(date)} UTC
    </span>
{:else}
    <span class="txt txt-hint">N/A</span>
{/if}
