<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let file; // File() instance
    export let size = 50; // preview thumb size (if file is image)

    $: if (typeof file !== "undefined") {
        loadPreviewUrl();
    }

    $: previewUrl = "";

    function loadPreviewUrl() {
        previewUrl = "";

        if (CommonHelper.hasImageExtension(file?.name)) {
            CommonHelper.generateThumb(file, size, size)
                .then((url) => {
                    previewUrl = url;
                })
                .catch((err) => {
                    console.warn("Unable to generate thumb: ", err);
                });
        }
    }
</script>

{#if previewUrl}
    <img src={previewUrl} width={size} height={size} alt={file.name} />
{:else}
    <i class="ri-file-line" alt={file.name} />
{/if}
