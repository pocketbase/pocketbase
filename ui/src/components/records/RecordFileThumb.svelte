<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import PreviewPopup from "@/components/base/PreviewPopup.svelte";

    export let record = null;
    export let filename = "";
    export let size = ""; // sm/lg/xl

    let previewPopup;
    let thumbUrl = "";
    let originalUrl = "";
    let token = "";
    let isLoadingToken = true;

    loadFileToken();

    $: type = CommonHelper.getFileType(filename);

    $: hasPreview = ["image", "audio", "video"].includes(type) || filename.endsWith(".pdf");

    $: originalUrl = !isLoadingToken ? ApiClient.files.getUrl(record, filename, { token }) : "";

    $: thumbUrl = !isLoadingToken
        ? ApiClient.files.getUrl(record, filename, { thumb: "100x100", token: token })
        : "";

    async function loadFileToken() {
        isLoadingToken = true;

        try {
            token = await ApiClient.getAdminFileToken(record.collectionId);
        } catch (err) {
            console.warn("File token failure:", err);
        }

        isLoadingToken = false;
    }

    function onError() {
        thumbUrl = "";
    }
</script>

{#if isLoadingToken}
    <div class="thumb {size ? `thumb-${size}` : ''}" />
{:else}
    <a
        draggable={false}
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
            <img
                draggable={false}
                loading="lazy"
                src={thumbUrl}
                alt={filename}
                title="Preview {filename}"
                on:error={onError}
            />
        {:else if type === "video" || type === "audio"}
            <i class="ri-video-line" />
        {:else}
            <i class="ri-file-3-line" />
        {/if}
    </a>
{/if}

{#if hasPreview}
    <PreviewPopup bind:this={previewPopup} />
{/if}
