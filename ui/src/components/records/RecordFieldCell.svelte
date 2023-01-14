<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import IdLabel from "@/components/base/IdLabel.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import RecordFilePreview from "@/components/records/RecordFilePreview.svelte";

    export let record;
    export let field;

    // rough text cut to avoid rendering large chunk of texts
    function cutText(text) {
        text = text || "";
        return text.length > 200 ? text.substring(0, 200) : text;
    }
</script>

<td class="col-type-{field.type} col-field-{field.name}">
    {#if CommonHelper.isEmpty(record[field.name])}
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
            {record[field.name]}
        </a>
    {:else if field.type === "date"}
        <FormattedDate date={record[field.name]} />
    {:else if field.type === "json"}
        <span class="txt txt-ellipsis">
            {cutText(JSON.stringify(record[field.name]))}
        </span>
    {:else if field.type === "select"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as item, i (i + item)}
                <span class="label">{item}</span>
            {/each}
        </div>
    {:else if field.type === "relation" || field.type === "user"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]).slice(0, 20) as item, i (i + item)}
                <IdLabel id={item} />
            {/each}
            {#if CommonHelper.toArray(record[field.name]).length > 20}
                ...
            {/if}
        </div>
    {:else if field.type === "file"}
        <div class="inline-flex">
            {#each CommonHelper.toArray(record[field.name]) as filename, i (i + filename)}
                <figure class="thumb thumb-sm">
                    <RecordFilePreview {record} {filename} />
                </figure>
            {/each}
        </div>
    {:else}
        <span class="txt txt-ellipsis" title={cutText(record[field.name])}>
            {cutText(record[field.name])}
        </span>
    {/if}
</td>
