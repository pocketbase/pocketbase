export function backupUploadButton(onSuccess = null) {
    const uniqueId = "backup_upload_" + app.utils.randomString();

    const data = store({
        isUploading: false,
    });

    function uploadConfirm(file) {
        if (!file) {
            return;
        }

        app.modals.confirm(
            `Note that we don't perform validations for the uploaded backup files. Proceed with extreme caution and only if you trust the source.\n\n`
                + `Do you really want to upload "${file.name}"?`,
            () => {
                uploadBackup(file);
            },
            () => {
                resetSelectedFile();
            },
        );
    }

    async function uploadBackup(file) {
        if (!file || data.isUploading) {
            return;
        }

        data.isUploading = true;

        try {
            const formData = new FormData();
            formData.set("file", file);

            await app.pb.backups.upload(formData, { requestKey: uniqueId });

            data.isUploading = false;

            onSuccess(file);

            app.toasts.success("Successfully uploaded a new backup.");
        } catch (err) {
            if (!err.isAbort) {
                data.isUploading = false;
                if (err.response?.formData?.file?.message) {
                    app.toasts.error(err.response.formData.file.message);
                } else {
                    app.checkApiError(err);
                }
            }
        }

        resetSelectedFile();
    }

    function resetSelectedFile() {
        if (fileInput) {
            fileInput.value = "";
        }
    }

    const fileInput = t.input({
        type: "file",
        accept: "application/zip",
        className: "hidden",
        onchange: (e) => {
            uploadConfirm(e.target?.files?.[0]);
        },
    });

    return t.div(
        null,
        t.button(
            {
                type: "button",
                ariaLabel: app.attrs.tooltip("Upload backup"),
                className: () => `btn sm transparent secondary circle ${data.isUploading ? "loading" : ""}`,
                disabled: () => data.isUploading,
                onclick: () => fileInput?.click(),
                onunmount: () => {
                    app.pb.cancelRequest(uniqueId);
                },
            },
            t.i({ className: "ri-upload-cloud-line", ariaHidden: true }),
        ),
        fileInput,
    );
}
