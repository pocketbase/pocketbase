import { settingsSidebar } from "../settingsSidebar";
import { backupsForm } from "./backupsForm";
import { backupsList } from "./backupsList";
import { backupUploadButton } from "./backupUploadButton";

export function pageBackupsSettings(route) {
    app.store.title = "Backups";

    const data = store({
        resetList: null,
    });

    function resetBackupsList() {
        data.resetList = Date.now();
    }

    return t.div(
        { pbEvent: "pageBackupsSettings", className: "page page-backups-settings" },
        settingsSidebar(),
        t.div(
            { className: "page-content full-height" },
            t.header(
                { className: "page-header" },
                t.nav(
                    { className: "breadcrumbs" },
                    t.div({ className: "breadcrumb-item" }, "Settings"),
                    t.div({ className: "breadcrumb-item" }, () => app.store.title),
                ),
            ),
            t.div(
                { className: "wrapper m-b-base" },
                t.div(
                    {
                        className: "grid",
                    },
                    t.div(
                        { className: "col-lg-12" },
                        t.div(
                            { className: "flex gap-10 m-b-sm" },
                            t.div({ className: "txt-lg" }, "Backup and restore your PocketBase data"),
                            app.components.refreshButton({
                                className: "btn sm transparent secondary circle tooltip-bottom",
                                onclick: resetBackupsList,
                            }),
                            backupUploadButton(resetBackupsList),
                        ),
                        backupsList({
                            reset: () => data.resetList,
                        }),
                    ),
                    t.div(
                        { className: "col-lg-12" },
                        backupsForm({
                            onsave: () => resetBackupsList(),
                        }),
                    ),
                ),
            ),
            t.footer({ className: "page-footer" }, app.components.credits()),
        ),
    );
}
