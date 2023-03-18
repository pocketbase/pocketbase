<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { collections } from "@/stores/collections";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";

    export let record;
    export let displayFields = [];

    $: collection = $collections?.find((item) => item.id == record?.collectionId);

    $: fileDisplayFields =
        displayFields?.filter((name) => {
            return !!collection?.schema?.find((field) => field.name == name && field.type == "file");
        }) || [];

    $: textDisplayFields =
        (!fileDisplayFields.length
            ? displayFields
            : displayFields?.filter((name) => !fileDisplayFields.includes(name))) || [];
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

    {#each fileDisplayFields as name}
        {@const filenames = CommonHelper.toArray(record[name]).slice(0, 5)}
        {#each filenames as filename}
            {#if !CommonHelper.isEmpty(filename)}
                <RecordFileThumb {record} {filename} size="xs" />
            {/if}
        {/each}
    {/each}

    <span class="txt txt-ellipsis">
        {CommonHelper.truncate(CommonHelper.displayValue(record, textDisplayFields), 70)}
    </span>
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
        :global(.thumb) {
            box-shadow: none;
        }
    }
</style>
