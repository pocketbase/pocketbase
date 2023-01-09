<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import PreviewPopup from "@/components/base/PreviewPopup.svelte";

    export let record;
    export let filename;
    export let size;

    let previewPopup;
    let thumbUrl = "";
    let originalUrl = ApiClient.getFileUrl(record, `${filename}`);

    $: type = CommonHelper.getFileType(filename);
    $: hasPreview = ["image", "pdf"].includes(type);
    $: thumbUrl = originalUrl ? originalUrl + "?thumb=100x100" : "";

    function onError() {
        thumbUrl = "";
    }
</script>

<a
    class="thumb thumb-{size} link-fade"
    href={originalUrl}
    target="_blank"
    rel="noreferrer"
    on:click|stopPropagation={(e) => {
        if (!hasPreview) return;
        e.preventDefault();
        previewPopup?.show(originalUrl);
    }}
>
    {#if type === "image"}
        <img src={thumbUrl} alt={filename} title="Preview {filename}" on:error={onError} />
    {:else if type === "pdf"}
        <i class="ri-file-pdf-line" />
    {:else}
        <i class="ri-file-3-line" />
    {/if}
</a>

<PreviewPopup bind:this={previewPopup} {type} />
