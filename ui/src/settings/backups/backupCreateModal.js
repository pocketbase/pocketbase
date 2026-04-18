export function openBackupCreateModal(settings = {
    oncreated: null,
}) {
    const modal = backupCreateModal(settings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);

    app.modals.open(modal);
}

function backupCreateModal(settings) {
    let modal;

    const uniqueId = "backup_create_" + app.utils.randomString();

    const data = store({
        name: "",
        isSubmitting: false,
    });

    let submitTimeoutId;

    async function submit() {
        if (data.isSubmitting) {
            return;
        }

        data.isSubmitting = true;

        clearTimeout(submitTimeoutId);
        submitTimeoutId = setTimeout(() => {
            app.modals.close(modal);
        }, 1500);

        try {
            await app.pb.backups.create(data.name, { requestKey: uniqueId });

            data.isSubmitting = false;

            if (settings.oncreated) {
                settings.oncreated(data.name);
            }

            app.toasts.success("Successfully generated new backup.");

            app.modals.close(modal);
        } catch (err) {
            if (!err.isAbort) {
                clearTimeout(submitTimeoutId);
                data.isSubmitting = false;
                app.checkApiError(err);
            }
        }
    }

    modal = t.div(
        {
            pbEvent: "backupCreateModal",
            className: "modal popup backup-create-modal",
            onbeforeclose: () => {
                if (data.isSubmitting) {
                    app.toasts.info(
                        "The backup was started but may take a while to complete. You can come back later.",
                    );
                }
            },
            onafterclose: (el) => {
                clearTimeout(submitTimeoutId);
                el?.remove();
            },
        },
        t.header(
            { className: "modal-header" },
            t.h5({ className: "m-auto txt-center" }, "Initialize new backup"),
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
                        { className: "alert warning" },
                        t.div(
                            { className: "content" },
                            t.p(
                                null,
                                `Please note that during the backup other concurrent write requests may fail since the database will be temporary "locked" (this usually happens only during the ZIP generation).`,
                            ),
                            t.p(
                                { className: "txt-bold" },
                                `If you are using S3 storage for the collections file upload, you'll have to backup them separately since they are not locally stored and they will not be included in the generated backup!`,
                            ),
                        ),
                    ),
                ),
                t.div(
                    { className: "col-lg-12" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + "_name" }, "Backup name"),
                        t.input({
                            id: uniqueId + "_name",
                            name: "name",
                            type: "text",
                            pattern: "^[a-z0-9_-]+\.zip$",
                            placeholder: "Leave empty to autogenerate",
                            value: () => data.name,
                            oninput: (e) => (data.name = e.target.value),
                        }),
                    ),
                    t.div({ className: "field-help" }, "Must be in the format [a-z0-9_-].zip"),
                ),
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    disabled: () => data.isSubmitting,
                    onclick: () => app.modals.close(modal),
                },
                t.span({ className: "txt" }, "Cancel"),
            ),
            t.button(
                {
                    "html-form": uniqueId,
                    type: "submit",
                    className: () => `btn ${data.isSubmitting ? "loading" : ""}`,
                    disabled: () => data.isSubmitting,
                },
                t.span({ className: "txt" }, "Start backup"),
            ),
        ),
    );

    return modal;
}
