<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { collections } from "@/stores/collections";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";

    export let record;

    let fileDisplayFields = [];
    let nonFileDisplayFields = [];

    $: collection = $collections?.find((item) => item.id == record?.collectionId);

    $: if (collection) {
        loadDisplayFields();
    }

    function loadDisplayFields() {
        const fields = collection?.schema || [];

        // reset
        fileDisplayFields = fields.filter((f) => f.presentable && f.type == "file").map((f) => f.name);
        nonFileDisplayFields = fields.filter((f) => f.presentable && f.type != "file").map((f) => f.name);

        // fallback to the first single file field that accept images
        // if no presentable field is available
        if (!fileDisplayFields.length && !nonFileDisplayFields.length) {
            const fallbackFileField = fields.find((f) => {
                return (
                    f.type == "file" &&
                    f.options?.maxSelect == 1 &&
                    f.options?.mimeTypes?.find((t) => t.startsWith("image/"))
                );
            });
            if (fallbackFileField) {
                fileDisplayFields.push(fallbackFileField.name);
            }
        }
    }
</script>

<div class="record-info">
    <i
        class="link-hint txt-sm ri-information-line"
        use:tooltip={{
            text: CommonHelper.truncate(
                JSON.stringify(CommonHelper.truncateObject(record), null, 2),
                800,
                true,
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
        {CommonHelper.truncate(CommonHelper.displayValue(record, nonFileDisplayFields), 70)}
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
        :global(.thumb) {
            box-shadow: none;
        }
    }
</style>
