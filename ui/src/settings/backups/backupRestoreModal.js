export function openBackupRestoreModal(key) {
    const modal = backupRestoreModal(key);

    document.body.appendChild(modal);

    app.modals.open(modal);
}

function backupRestoreModal(key) {
    const uniqueId = "backup_restore_" + app.utils.randomString();

    const data = store({
        key: key,
        keyConfirm: "",
        isSubmitting: false,
        get canSubmit() {
            return data.key && data.key == data.keyConfirm;
        },
    });

    let reloadTimeoutId;

    async function submit() {
        if (data.isSubmitting || !data.canSubmit) {
            return;
        }

        clearTimeout(reloadTimeoutId);

        data.isSubmitting = true;

        try {
            await app.pb.backups.restore(data.keyConfirm);

            // optimistic restore page reload
            reloadTimeoutId = setTimeout(() => {
                window.location.reload();
                data.isSubmitting = false;
            }, 2000);
        } catch (err) {
            clearTimeout(reloadTimeoutId);

            if (!err?.isAbort) {
                data.isSubmitting = false;
                app.checkApiError(err);
            }
        }
    }

    return t.div(
        {
            pbEvent: "backupRestoreModal",
            className: "modal popup backup-restore-modal",
            onbeforeclose: () => {
                return !data.isSubmitting;
            },
            onafterclose: (el) => {
                el?.remove();
            },
            onunmount: () => {
                clearTimeout(reloadTimeoutId);
            },
        },
        t.header(
            { className: "modal-header" },
            t.h5(
                { className: "m-auto txt-center" },
                "Restore ",
                t.strong(null, () => data.key),
            ),
        ),
        t.form(
            {
                id: uniqueId,
                className: "modal-content backup-restore-form",
                autocomplete: "off",
                onsubmit: (e) => {
                    e.preventDefault();
                    submit();
                },
            },
            t.div(
                { className: "grid" },
                t.div(
                    { className: "col-lg-12" },
                    t.div(
                        { className: "alert danger" },
                        t.div(
                            { className: "content" },
                            t.p(
                                { className: "txt-bold" },
                                "Please proceed with extreme caution and use it only with trusted backups!",
                            ),
                            t.p(null, "Backup restore currently works only on UNIX based systems."),
                            t.p(
                                null,
                                "The restore operation will attempt to replace your existing ",
                                t.code(null, "pb_data"),
                                " with the one from the backup and will restart the application process.",
                            ),
                            t.p(
                                null,
                                "This means that on success all of your data (including app settings, users, superusers, etc.) will be replaced with the ones from the backup.",
                            ),
                            t.p(
                                null,
                                "The operation will be reverted if the backup is invalid (ex. missing ",
                                t.code(null, "data.db"),
                                " file).",
                            ),
                            t.p(null, "Below is an oversimplified version of the restore flow:"),
                            t.ol(
                                null,
                                t.li(
                                    null,
                                    "Replaces the current ",
                                    t.code(null, "pb_data"),
                                    " with the content from the backup.",
                                ),
                                t.li(null, "Triggers app restart."),
                                t.li(
                                    null,
                                    "Applies all migrations that are missing in the restored ",
                                    t.code(null, "pb_data"),
                                    ".",
                                ),
                                t.li(null, "Initializes the app server as usual."),
                            ),
                        ),
                    ),
                ),
                t.div(
                    { className: "col-lg-12" },
                    t.div(
                        { className: "confirm-key-label m-b-sm" },
                        "Type the backup name ",
                        t.div(
                            { className: "label" },
                            () => data.key,
                            app.components.copyButton(() => data.key),
                        ),
                        " to confirm:",
                    ),
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + "_key" }, "Backup name"),
                        t.input({
                            id: uniqueId + "_key",
                            name: "key",
                            type: "text",
                            required: true,
                            value: () => data.keyConfirm,
                            oninput: (e) => (data.keyConfirm = e.target.value),
                        }),
                    ),
                ),
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => app.modals.close(),
                    disabled: () => data.isSubmitting,
                },
                t.span({ className: "txt" }, "Cancel"),
            ),
            t.button(
                {
                    "html-form": uniqueId,
                    type: "submit",
                    className: () => `btn ${data.isSubmitting ? "loading" : ""}`,
                    disabled: () => data.isSubmitting || !data.canSubmit,
                },
                t.span({ className: "txt" }, "Restore backup"),
            ),
        ),
    );
}
