<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import IdLabel from "@/components/base/IdLabel.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import RecordFilePreview from "@/components/records/RecordFilePreview.svelte";

    export let record;
    export let field;
</script>

<td class="col-type-{field.type} col-field-{field.name}">
    {#if CommonHelper.isEmpty(record[field.name])}
        <span class="txt-hint">N/A</span>
    {:else if field.type === "bool"}
        <span class="txt">{record[field.name] ? "True" : "False"}</span>
    {:else if field.type === "url"}
        <a
            class="txt-ellipsis"
            href={record[field.name]}
            target="_blank"
            rel="noopener"
            use:tooltip={"Open in new tab"}
            on:click|stopPropagation
        >
            {record[field.name]}
        </a>
    {:else if field.type === "date"}
        <FormattedDate date={record[field.name]} />
    {:else if field.type === "json"}
        <span class="txt txt-ellipsis">{JSON.stringify(record[field.name])}</span>
    {:else if field.type === "select"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as item}
                <span class="label">{item}</span>
            {/each}
        </div>
    {:else if field.type === "relation" || field.type === "user"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as item}
                <IdLabel id={item} />
            {/each}
        </div>
    {:else if field.type === "file"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as filename}
                <figure class="thumb thumb-sm">
                    <RecordFilePreview {record} {filename} />
                </figure>
            {/each}
        </div>
    {:else}
        <span class="txt txt-ellipsis" title={record[field.name]}>
            {record[field.name]}
        </span>
    {/if}
</td>
