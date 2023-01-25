<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";

    export let record;
    export let displayFields = [];

    $: displayValue = CommonHelper.displayValue(record, displayFields);
</script>

<div class="record-info">
    <i
        class="link-hint txt-sm ri-information-line"
        use:tooltip={{
            text: CommonHelper.truncate(
                JSON.stringify(CommonHelper.truncateObject(record), null, 2),
                800,
                true
            ),
            class: "code",
            position: "left",
        }}
    />
    <span class="txt txt-ellipsis">{CommonHelper.truncate(displayValue, 150)}</span>
</div>

<style lang="scss">
    .record-info {
        display: inline-flex;
        vertical-align: top;
        align-items: center;
        max-width: 100%;
        min-width: 0;
        gap: 5px;
        line-height: normal;
        > * {
            line-height: inherit;
        }
    }
</style>
