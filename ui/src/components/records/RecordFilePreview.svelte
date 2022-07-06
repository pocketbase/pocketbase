<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let record;
    export let filename;

    let previewUrl = "";

    $: if (CommonHelper.hasImageExtension(filename)) {
        previewUrl = ApiClient.Records.getFileUrl(record, `${filename}?thumb=100x100`);
    }

    function onError() {
        previewUrl = "";
    }
</script>

{#if previewUrl}
    <img src={previewUrl} alt={filename} on:error={onError} />
{:else}
    <i class="ri-file-line" />
{/if}
