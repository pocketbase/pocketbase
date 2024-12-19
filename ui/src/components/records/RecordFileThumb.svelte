<script>
    import PreviewPopup from "@/components/base/PreviewPopup.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let record = null;
    export let filename = "";
    export let size = ""; // sm/lg/xl

    let previewPopup;
    let thumbUrl = "";
    let token = "";
    let isLoadingToken = true;

    loadThumbUrlToken();

    $: type = CommonHelper.getFileType(filename);

    $: hasPreview = ["image", "audio", "video"].includes(type) || filename.endsWith(".pdf");

    $: thumbUrl = !isLoadingToken
        ? ApiClient.files.getURL(record, filename, { thumb: "100x100", token: token })
        : "";

    async function loadThumbUrlToken() {
        isLoadingToken = true;

        try {
            token = await ApiClient.getSuperuserFileToken(record.collectionId);
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
    <button
        type="button"
        draggable={false}
        class="handle thumb {size ? `thumb-${size}` : ''}"
        title={(hasPreview ? "Preview" : "Download") + " " + filename}
        on:click|stopPropagation={async () => {
            if (!hasPreview) {
                return;
            }

            try {
                // refetch the token because it could have expired
                previewPopup?.show(async () => {
                    token = await ApiClient.getSuperuserFileToken(record.collectionId);
                    return ApiClient.files.getURL(record, filename, { token });
                });
            } catch (err) {
                if (!err.isAbort) {
                    console.warn("Preview file token failure:", err);
                }
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
    </button>
{/if}

{#if hasPreview}
    <PreviewPopup bind:this={previewPopup} />
{/if}
