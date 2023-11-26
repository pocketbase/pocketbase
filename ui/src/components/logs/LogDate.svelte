<script>
    import tooltip from "@/actions/tooltip";
    import CommonHelper from "@/utils/CommonHelper";

    export let date;

    const tooltipData = {
        // generate the tooltip text as getter to speed up the initial load
        // in case the component is used with large number of items
        get text() {
            return CommonHelper.formatToLocalDate(date, "yyyy-MM-dd HH:mm:ss.SSS") + " Local";
        },
    };
</script>

<span class="txt-nowrap" use:tooltip={tooltipData}>
    {date.replace("Z", " UTC")}
</span>
