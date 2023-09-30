<script>
    import tooltip from "@/actions/tooltip";
    import CommonHelper from "@/utils/CommonHelper";

    const detailedDateFormat = "yyyy-MM-dd HH:mm:ss.SSS";

    export let model;

    let tooltipDates = [];

    $: if (model) {
        refreshTooltipDates();
    }

    function refreshTooltipDates() {
        tooltipDates = [];

        if (model.created) {
            tooltipDates.push(
                "Created: " + CommonHelper.formatToLocalDate(model.created, detailedDateFormat) + " Local"
            );
        }

        if (model.updated) {
            tooltipDates.push(
                "Updated: " + CommonHelper.formatToLocalDate(model.updated, detailedDateFormat) + " Local"
            );
        }
    }
</script>

<i
    class="ri-calendar-event-line txt-disabled"
    use:tooltip={{
        text: tooltipDates.join("\n"),
        position: "left",
    }}
/>
