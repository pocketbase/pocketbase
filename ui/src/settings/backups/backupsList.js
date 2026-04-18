import { openBackupCreateModal } from "./backupCreateModal";
import { openBackupRestoreModal } from "./backupRestoreModal";

export function backupsList(propsArg = {}) {
    const props = store({
        reset: null,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        canBackup: true,
        isLoading: false,
        isDownloading: {},
        isDeleting: {},
        backups: [],
    });

    async function loadBackups() {
        data.isLoading = true;

        try {
            data.backups = await app.pb.backups.getFullList();

            // sort backups DESC by their modified date
            data.backups.sort((a, b) => {
                if (a.modified < b.modified) {
                    return 1;
                }

                if (a.modified > b.modified) {
                    return -1;
                }

                return 0;
            });

            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                data.isLoading = false;
            }
        }
    }

    async function confirmBackupDelete(key) {
        app.modals.confirm(`Do you really want to delete ${key}?`, () => deleteBackup(key));
    }

    async function deleteBackup(key) {
        if (data.isDeleting[key]) {
            return;
        }

        data.isDeleting[key] = true;

        try {
            await app.pb.backups.delete(key);
            loadBackups();
            app.toasts.success(`Successfully deleted ${key}.`);
        } catch (err) {
            app.checkApiError(err);
        }

        delete data.isDeleting[key];
    }

    async function loadCanBackup() {
        try {
            const health = await app.pb.health.check({ requestKey: null });
            const oldCanBackup = data.canBackup;
            data.canBackup = health?.data?.canBackup || false;

            // reload backups list
            if (data.canBackup && oldCanBackup != data.canBackup) {
                loadBackups();
            }
        } catch (err) {
            console.warn("failed to load canBackup checks", err);
        }
    }

    async function downloadBackup(key) {
        if (data.isDownloading[key]) {
            return;
        }

        data.isDownloading[key] = true;

        try {
            const token = await app.pb.files.getToken({ requestKey: null });
            app.utils.download(app.pb.backups.getDownloadURL(token, key));
        } catch (err) {
            app.checkApiError(err);
        }

        delete data.isDownloading[key];
    }

    return t.div(
        {
            className: "list",
            onmount: (el) => {
                watchers.push(watch(() => props.reset, () => {
                    loadBackups();
                }));

                el._canBackupIntervalId = setInterval(() => {
                    loadCanBackup();
                }, 3500);
            },
            onunmount: (el) => {
                clearInterval(el._canBackupIntervalId);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            {
                hidden: () => !data.isLoading || data.backups.length,
                className: "list-item",
            },
            t.div({ className: "skeleton-loader" }),
        ),
        t.div(
            {
                hidden: () => data.isLoading || data.backups.length,
                className: () => "list-item",
            },
            t.div({ className: "content block txt-hint" }, "No backups found."),
        ),
        () => {
            return data.backups.map((backup) => {
                return t.div(
                    { className: () => `list-item ${data.isLoading ? "faded" : ""}` },
                    t.i({ className: "ri-folder-zip-line", ariaHidden: true }),
                    t.div(
                        { className: "content" },
                        t.span({
                            className: "backup-name txt-ellipsis",
                            title: () => backup.key,
                            textContent: () => backup.key,
                        }),
                        t.small(
                            { className: "backup-size txt-hint txt-nowrap" },
                            "(",
                            () => app.utils.formattedFileSize(backup.size),
                            ")",
                        ),
                    ),
                    t.nav(
                        {
                            hidden: () => data.isLoading,
                            className: "actions autohide",
                        },
                        t.button(
                            {
                                type: "button",
                                ariaLabel: app.attrs.tooltip("Download"),
                                className: () =>
                                    `btn sm circle secondary transparent ${
                                        data.isDownloading[backup.key] ? "loading" : ""
                                    }`,
                                disabled: () => data.isDeleting[backup.key] || data.isDownloading[backup.key],
                                onclick: () => downloadBackup(backup.key),
                            },
                            t.i({ className: "ri-download-line", ariaHidden: true }),
                        ),
                        t.button(
                            {
                                type: "button",
                                ariaLabel: app.attrs.tooltip("Restore"),
                                className: () => `btn sm circle secondary transparent`,
                                disabled: () => data.isDeleting[backup.key] || data.isDownloading[backup.key],
                                onclick: () => openBackupRestoreModal(backup.key),
                            },
                            t.i({ className: "ri-restart-line", ariaHidden: true }),
                        ),
                        t.button(
                            {
                                type: "button",
                                ariaLabel: app.attrs.tooltip("Delete"),
                                className: () =>
                                    `btn sm circle secondary transparent ${
                                        data.isDeleting[backup.key] ? "loading" : ""
                                    }`,
                                disabled: () => data.isDeleting[backup.key] || data.isDownloading[backup.key],
                                onclick: () => confirmBackupDelete(backup.key),
                            },
                            t.i({ className: "ri-delete-bin-7-line", ariaHidden: true }),
                        ),
                    ),
                );
            });
        },
        t.div(
            { className: "list-item" },
            t.button(
                {
                    type: "button",
                    className: () => `btn secondary block ${data.isLoading ? "loading" : ""}`,
                    disabled: () => !data.canBackup || data.isLoading,
                    onclick: () => {
                        openBackupCreateModal({
                            oncreated: () => loadBackups(),
                        });
                    },
                },
                () => {
                    if (data.canBackup) {
                        return [
                            t.i({ className: "ri-play-circle-line", ariaHidden: true }),
                            t.span({ className: "txt" }, "Initialize new backup"),
                        ];
                    }

                    return [
                        t.span({ className: "loader sm" }),
                        t.span({ className: "txt" }, "Backup/restore operation is in process"),
                    ];
                },
            ),
        ),
    );
}
