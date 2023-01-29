<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";

    export let record;
    export let field;
</script>

<td class="col-type-{field.type} col-field-{field.name}">
    {#if field.type === "json"}
        <span class="txt txt-ellipsis">
            {CommonHelper.truncate(JSON.stringify(record[field.name]))}
        </span>
    {:else if CommonHelper.isEmpty(record[field.name])}
        <span class="txt-hint">N/A</span>
    {:else if field.type === "bool"}
        <span class="txt">{record[field.name] ? "True" : "False"}</span>
    {:else if field.type === "number"}
        <span class="txt">{record[field.name]}</span>
    {:else if field.type === "url"}
        <a
            class="txt-ellipsis"
            href={record[field.name]}
            target="_blank"
            rel="noopener noreferrer"
            use:tooltip={"Open in new tab"}
            on:click|stopPropagation
        >
            {CommonHelper.truncate(record[field.name])}
        </a>
    {:else if field.type === "editor"}
        <span class="txt">
            {CommonHelper.truncate(CommonHelper.plainText(record[field.name]), 300, true)}
        </span>
    {:else if field.type === "date"}
        <FormattedDate date={record[field.name]} />
    {:else if field.type === "select"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as item, i (i + item)}
                <span class="label">{item}</span>
            {/each}
        </div>
    {:else if field.type === "relation" || field.type === "user"}
        {@const relations = CommonHelper.toArray(record[field.name])}
        {@const expanded = CommonHelper.toArray(record.expand[field.name])}
        <div class="inline-flex">
            {#if expanded.length}
                {#each expanded.slice(0, 20) as item, i (i + item)}
                    <span class="label">
                        <RecordInfo record={item} displayFields={field.options?.displayFields} />
                    </span>
                {/each}
            {:else}
                {#each relations.slice(0, 20) as id}
                    <span class="label">{id}</span>
                {/each}
            {/if}
            {#if relations.length > 20}
                ...
            {/if}
        </div>
    {:else if field.type === "file"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as filename, i (i + filename)}
                <RecordFileThumb {record} {filename} size="sm" />
            {/each}
        </div>
    {:else}
        <span class="txt txt-ellipsis" title={CommonHelper.truncate(record[field.name])}>
            {CommonHelper.truncate(record[field.name])}
        </span>
    {/if}
</td>

<style>
    .filename {
        max-width: 200px;
    }
</style>
