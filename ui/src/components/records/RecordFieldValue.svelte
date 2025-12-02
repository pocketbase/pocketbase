<script>
    import tooltip from "@/actions/tooltip";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import TinyMCE from "@/components/base/TinyMCE.svelte";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";
    import GeoPointValue from "@/components/records/fields/GeoPointValue.svelte";
    import { superuser } from "@/stores/superuser";
    import { isDarkMode } from "@/stores/app";
    import CommonHelper from "@/utils/CommonHelper";


    export let record;
    export let field;
    export let short = false;

    $: rawValue = record?.[field.name];
</script>

{#if field.primaryKey}
    <div class="label">
        <CopyIcon value={rawValue} />
        <div class="txt txt-ellipsis">{rawValue}</div>
    </div>
    {#if record.collectionName == "_superusers" && record.id == $superuser.id}
        <span class="label label-warning">You</span>
    {/if}
{:else if field.type === "json"}
    {@const stringifiedJson = CommonHelper.trimQuotedValue(JSON.stringify(rawValue)) || '""'}
    {#if short}
        <span class="txt txt-ellipsis">
            {CommonHelper.truncate(stringifiedJson)}
        </span>
    {:else}
        <span class="txt">
            {CommonHelper.truncate(stringifiedJson, 500, true)}
        </span>
        {#if stringifiedJson.length > 500}
            <CopyIcon value={JSON.stringify(rawValue, null, 2)} />
        {/if}
    {/if}
{:else if CommonHelper.isEmpty(rawValue)}
    <span class="txt-hint">N/A</span>
{:else if field.type === "bool"}
    <span class="label" class:label-success={!!rawValue}>{rawValue ? "True" : "False"}</span>
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
            {CommonHelper.truncate(CommonHelper.plainText(rawValue), 195)}
        </span>
    {:else}
        {#key $isDarkMode}
            <TinyMCE
                cssClass="tinymce-preview"
                conf={{
                    branding: false,
                    promotion: false,
                    menubar: false,
                    statusbar: false,
                    min_height: 30,
                    height: 59,
                    max_height: 500,
                    autoresize_bottom_margin: 5,
                    resize: false,
                    content_style: `
                        body {
                            font-size: 14px;
                            font-family: 'Source Sans 3', sans-serif, emoji;
                            color: ${$isDarkMode ? "#dcdcdc" : "#1a1a24"};
                            background: ${$isDarkMode ? "#1c1c21" : "#ffffff"};
                        }
                        a { color: #5499e8; }
                        code {
                            background-color: ${$isDarkMode ? "#303038" : "#e8e8e8"};
                            color: inherit;
                            border-radius: 3px;
                            padding: 0.1rem 0.2rem;
                        }
                        blockquote {
                            border-left: 2px solid ${$isDarkMode ? "#42424e" : "#ccc"};
                            margin-left: 1.5rem;
                            padding-left: 1rem;
                        }
                        table th, table td {
                            border-color: ${$isDarkMode ? "#42424e" : "#ccc"};
                        }
                        hr {
                            border-color: ${$isDarkMode ? "#42424e" : "#ccc"};
                        }
                    `,
                    toolbar: "",
                    plugins: ["autoresize"],
                    skin: "pocketbase",
                }}
                value={rawValue}
                disabled
            />
        {/key}
    {/if}
{:else if field.type === "date" || field.type === "autodate"}
    <FormattedDate date={rawValue} />
{:else if field.type === "select"}
    <div class="inline-flex">
        {#each CommonHelper.toArray(rawValue) as item, i (i + item)}
            <span class="label">{item}</span>
        {/each}
    </div>
{:else if field.type === "relation"}
    {@const relations = CommonHelper.toArray(rawValue)}
    {@const expanded = CommonHelper.toArray(record?.expand?.[field.name])}
    {@const relLimit = short ? 20 : 500}
    <div class="inline-flex">
        {#if expanded.length}
            {#each expanded.slice(0, relLimit) as item, i (i + item)}
                <span class="label">
                    <RecordInfo record={item} />
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
    {@const files = CommonHelper.toArray(rawValue)}
    {@const filesLimit = short ? 10 : 500}
    <div class="inline-flex" class:multiple={field.maxSelect != 1}>
        {#each files.slice(0, filesLimit) as filename, i (i + filename)}
            <RecordFileThumb {record} {filename} size="sm" />
        {/each}
        {#if files.length > filesLimit}
            ...
        {/if}
    </div>
{:else if field.type === "geoPoint"}
    <div class="label"><GeoPointValue value={rawValue} /></div>
{:else if short}
    <span class="txt txt-ellipsis" title={CommonHelper.truncate(rawValue)}>
        {CommonHelper.truncate(rawValue)}
    </span>
{:else}
    <div class="block txt-break fallback-block">{rawValue}</div>
{/if}

<style>
    .fallback-block {
        max-height: 100px;
        overflow: auto;
    }
    .inline-flex {
        max-width: 100%;
    }
</style>
