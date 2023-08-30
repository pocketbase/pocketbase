<script>
    import { createEventDispatcher, onDestroy } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import tooltip from "@/actions/tooltip";

    const dispatch = createEventDispatcher();
    const backupRequestKey = "upload_backup";

    let classes = "";
    export { classes as class };

    let fileInput;
    let isUploading = false;

    async function upload(e) {
        if (isUploading || !e?.target?.files?.length) {
            return;
        }

        isUploading = true;

        const data = new FormData();
        data.set("file", e.target.files[0]);

        try {
            await ApiClient.backups.upload(data, { requestKey: backupRequestKey });
            isUploading = false;

            dispatch("success");
            addSuccessToast("Successfully uploaded a new backup.");
        } catch (err) {
            if (!err.isAbort) {
                isUploading = false;
                if (err.response?.data?.file?.message) {
                    addErrorToast(err.response.data.file.message);
                } else {
                    ApiClient.error(err);
                }
            }
        }
    }

    onDestroy(() => {
        ApiClient.cancelRequest(backupRequestKey);
    });
</script>

<button
    type="button"
    class="btn btn-circle btn-transparent {classes}"
    class:btn-loading={isUploading}
    class:btn-disabled={isUploading}
    aria-label="Upload backup"
    use:tooltip={"Upload backup"}
    on:click={() => fileInput?.click()}
>
    <i class="ri-upload-cloud-line" />
</button>

<input bind:this={fileInput} type="file" accept="application/zip" class="hidden" on:change={upload} />
