<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let record;
    export let filename;

    let previewUrl = "";
    let fileType = false;

    if (CommonHelper.canBePreviewed(filename)) {
        previewUrl = ApiClient.Records.getFileUrl(record, filename);
        fileType = CommonHelper.getFileType(filename);
        // Apply thumb size for image preview
        if(fileType === "image") {
            previewUrl += '?thumb=100x100';
        }
    }

    function onError() {
        previewUrl = "";
    }
</script>

{#if fileType === "image"}
    <img src={previewUrl} alt={filename} on:error={onError} />
{:else if fileType === "video"}
    <video src={previewUrl} width="100%" alt={filename} on:error={onError}><track kind="captions"></video>
{:else}
    <i class="ri-file-line" />
{/if}
