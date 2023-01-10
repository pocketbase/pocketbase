<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import PreviewPopup from "@/components/base/PreviewPopup.svelte";

    export let record = null;
    export let filename = "";
    export let size = ""; // sm/lg/xl

    let previewPopup;
    let thumbUrl = "";
    let originalUrl = ApiClient.getFileUrl(record, filename);

    $: type = CommonHelper.getFileType(filename);

    $: hasPreview = ["image", "audio", "video"].includes(type) || filename.endsWith(".pdf");

    $: thumbUrl = originalUrl ? originalUrl + "?thumb=100x100" : "";

    function onError() {
        thumbUrl = "";
    }
</script>

<a
    class="thumb {size ? `thumb-${size}` : ''}"
    href={originalUrl}
    target="_blank"
    rel="noreferrer"
    title={(hasPreview ? "Preview" : "Download") + " " + filename}
    on:click|stopPropagation={(e) => {
        if (hasPreview) {
            e.preventDefault();
            previewPopup?.show(originalUrl);
        }
    }}
>
    {#if type === "image"}
        <img src={thumbUrl} alt={filename} title="Preview {filename}" on:error={onError} />
    {:else if type === "video" || type === "audio"}
        <i class="ri-video-line" />
    {:else}
        <i class="ri-file-3-line" />
    {/if}
</a>

<PreviewPopup bind:this={previewPopup} />
