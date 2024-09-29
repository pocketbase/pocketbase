<script>
    import tooltip from "@/actions/tooltip";
    import { collections } from "@/stores/collections";
    import CommonHelper from "@/utils/CommonHelper";

    const detailedDateFormat = "yyyy-MM-dd HH:mm:ss.SSS";

    export let record;

    let tooltipDates = [];

    $: collection = record && $collections.find((c) => c.id == record.collectionId);

    $: if (record) {
        refreshTooltipDates();
    }

    function refreshTooltipDates() {
        tooltipDates = [];

        const fields = collection.fields || [];

        for (let field of fields) {
            if (field.type != "autodate") {
                continue;
            }

            tooltipDates.push(
                field.name +
                    ": " +
                    CommonHelper.formatToLocalDate(record[field.name], detailedDateFormat) +
                    " Local",
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
