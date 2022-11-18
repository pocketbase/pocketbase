<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import PreviewPopup from "@/components/base/PreviewPopup.svelte";

    export let record;
    export let filename;

    let previewPopup;
    let thumbUrl = "";
    let originalUrl = "";

    $: hasPreview = CommonHelper.hasImageExtension(filename);

    $: if (hasPreview) {
        originalUrl = ApiClient.getFileUrl(record, `${filename}`);
    }

    $: thumbUrl = originalUrl ? originalUrl + "?thumb=100x100" : "";

    function onError() {
        thumbUrl = "";
    }
</script>

{#if hasPreview}
    <img
        src={thumbUrl}
        alt={filename}
        title="Preview {filename}"
        class:link-fade={hasPreview}
        on:click={(e) => {
            e.stopPropagation();
            previewPopup?.show(originalUrl);
        }}
        on:error={onError}
    />
{:else}
    <i class="ri-file-line" />
{/if}

<PreviewPopup bind:this={previewPopup} />
