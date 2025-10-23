<script>
    import CommonHelper from "@/utils/CommonHelper";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";
    import RecordInfoContent from "@/components/records/RecordInfoContent.svelte";
    import GeoPointValue from "@/components/records/fields/GeoPointValue.svelte";
    import { collections } from "@/stores/collections";

    export let record;

    let fileDisplayFields = [];
    let nonFileDisplayFields = [];

    $: collection = $collections?.find((item) => item.id == record?.collectionId);

    $: if (collection) {
        loadDisplayFields();
    }

    function loadDisplayFields() {
        const fields = collection?.fields || [];

        fileDisplayFields = fields.filter((f) => !f.hidden && f.presentable && f.type == "file");
        nonFileDisplayFields = fields.filter((f) => !f.hidden && f.presentable && f.type != "file");

        // fallback to the first single file field that accept images
        // if no presentable field is available
        if (!fileDisplayFields.length && !nonFileDisplayFields.length) {
            const fallbackFileField = fields.find((f) => {
                return (
                    !f.hidden &&
                    f.type == "file" &&
                    f.maxSelect == 1 &&
                    f.mimeTypes?.find((t) => t.startsWith("image/"))
                );
            });
            if (fallbackFileField) {
                fileDisplayFields.push(fallbackFileField);
            }
        }
    }
</script>

{#each fileDisplayFields as field}
    {@const filenames = CommonHelper.toArray(record[field.name]).slice(0, 5)}
    {#each filenames as filename}
        {#if !CommonHelper.isEmpty(filename)}
            <RecordFileThumb {record} {filename} size="xs" />
        {/if}
    {/each}
{/each}

{#each nonFileDisplayFields as field, i}
    {#if i > 0}
        <span class="delimiter">•</span>
    {/if}
    {#if field.type == "relation" && record.expand?.[field.name]}
        {@const isMultiple = Array.isArray(record.expand?.[field.name])}
        {@const expanded = CommonHelper.toArray(record.expand?.[field.name])}
        {#if isMultiple}<span class="expand-start">{"["}</span>{/if}
        {#each expanded.slice(0, 2) as expandRecord, j}
            {#if j > 0}<span class="delimiter">|</span>{/if}
            <RecordInfoContent record={expandRecord} />
        {/each}
        {#if expanded.length > 2}
            <span class="delimiter">|</span>
            <small class="delimiter txt-hint">({expanded.length - 2} more)</small>
        {/if}
        {#if isMultiple}<span class="expand-end">{"]"}</span>{/if}
    {:else if field.type == "geoPoint"}
        <GeoPointValue value={record[field.name]} />
    {:else}
        <span class="txt">{CommonHelper.truncate(CommonHelper.displayValue(record, [field.name]), 70)}</span>
    {/if}
{:else}
    <span class="txt">{CommonHelper.truncate(CommonHelper.displayValue(record, []), 70)}</span>
{/each}
