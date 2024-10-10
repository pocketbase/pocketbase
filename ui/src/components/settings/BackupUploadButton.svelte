<script>
    import { createEventDispatcher, onDestroy } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";

    const dispatch = createEventDispatcher();
    const backupRequestKey = "upload_backup";

    let classes = "";
    export { classes as class };

    let fileInput;
    let isUploading = false;

    function resetSelectedFile() {
        if (fileInput) {
            fileInput.value = "";
        }
    }

    function uploadConfirm(file) {
        if (!file) {
            return;
        }

        confirm(
            `Note that we don't perform validations for the uploaded backup files. Proceed with caution and only if you trust the source.\n\n` +
                `Do you really want to upload "${file.name}"?`,
            () => {
                upload(file);
            },
            () => {
                resetSelectedFile();
            },
        );
    }

    async function upload(file) {
        if (isUploading || !file) {
            return;
        }

        isUploading = true;

        const data = new FormData();
        data.set("file", file);

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

        resetSelectedFile();
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

<input
    bind:this={fileInput}
    type="file"
    accept="application/zip"
    class="hidden"
    on:change={(e) => {
        uploadConfirm(e?.target?.files?.[0]);
    }}
/>
