<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";
    import TinyMCE from "@tinymce/tinymce-svelte";

    export let record;
    export let field;
    export let short = false;

    $: rawValue = record[field.name];
</script>

{#if field.type === "json"}
    <span class="txt txt-ellipsis">
        {short
            ? CommonHelper.truncate(JSON.stringify(rawValue))
            : CommonHelper.truncate(JSON.stringify(rawValue, null, 2), 2000, true)}
    </span>
{:else if CommonHelper.isEmpty(rawValue)}
    <span class="txt-hint">N/A</span>
{:else if field.type === "bool"}
    <span class="txt">{rawValue ? "True" : "False"}</span>
{:else if field.type === "number"}
    <span class="txt">{rawValue}</span>
{:else if field.type === "url"}
    <a
        class="txt-ellipsis"
        href={rawValue}
        target="_blank"
        rel="noopener noreferrer"
        use:tooltip={"Open in new tab"}
        on:click|stopPropagation
    >
        {CommonHelper.truncate(rawValue)}
    </a>
{:else if field.type === "editor"}
    {#if short}
        <span class="txt">
            {CommonHelper.truncate(CommonHelper.plainText(rawValue), 250)}
        </span>
    {:else}
        <TinyMCE
            scriptSrc="{import.meta.env.BASE_URL}libs/tinymce/tinymce.min.js"
            cssClass="tinymce-preview"
            conf={{
                branding: false,
                promotion: false,
                menubar: false,
                min_height: 30,
                statusbar: false,
                height: 59,
                max_height: 500,
                autoresize_bottom_margin: 5,
                resize: false,
                skin: "pocketbase",
                content_style: "body { font-size: 14px }",
                toolbar: "",
                plugins: ["autoresize"],
            }}
            value={rawValue}
            disabled
        />
    {/if}
{:else if field.type === "date"}
    <FormattedDate date={rawValue} />
{:else if field.type === "select"}
    <div class="inline-flex">
        {#each CommonHelper.toArray(rawValue) as item, i (i + item)}
            <span class="label">{item}</span>
        {/each}
    </div>
{:else if field.type === "relation"}
    {@const relations = CommonHelper.toArray(rawValue)}
    {@const expanded = CommonHelper.toArray(record.expand[field.name])}
    {@const relLimit = short ? 20 : 200}
    <div class="inline-flex">
        {#if expanded.length}
            {#each expanded.slice(0, relLimit) as item, i (i + item)}
                <span class="label">
                    <RecordInfo record={item} displayFields={field.options?.displayFields} />
                </span>
            {/each}
        {:else}
            {#each relations.slice(0, relLimit) as id}
                <span class="label">{id}</span>
            {/each}
        {/if}
        {#if relations.length > relLimit}
            ...
        {/if}
    </div>
{:else if field.type === "file"}
    <div class="inline-flex">
        {#each CommonHelper.toArray(rawValue) as filename, i (i + filename)}
            <RecordFileThumb {record} {filename} size="sm" />
        {/each}
    </div>
{:else if short}
    <span class="txt txt-ellipsis" title={CommonHelper.truncate(rawValue)}>
        {CommonHelper.truncate(rawValue)}
    </span>
{:else}
    <span class="block txt-break">{CommonHelper.truncate(rawValue, 2000)}</span>
{/if}
