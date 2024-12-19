<script>
    import { onMount } from "svelte";
    import { slide } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import BackupCreatePanel from "@/components/settings/BackupCreatePanel.svelte";
    import BackupRestorePanel from "@/components/settings/BackupRestorePanel.svelte";

    let createPanel;
    let restorePanel;
    let backups = [];
    let isLoading = false;
    let isDownloading = {};
    let isDeleting = {};
    let canBackup = true;

    loadBackups();
    loadCanBackup();

    export async function loadBackups() {
        isLoading = true;

        try {
            backups = await ApiClient.backups.getFullList();

            // sort backups DESC by their modified date
            backups.sort((a, b) => {
                if (a.modified < b.modified) {
                    return 1;
                }

                if (a.modified > b.modified) {
                    return -1;
                }

                return 0;
            });

            isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
                isLoading = false;
            }
        }
    }

    async function download(name) {
        if (isDownloading[name]) {
            return;
        }

        isDownloading[name] = true;

        try {
            const token = await ApiClient.getSuperuserFileToken();
            CommonHelper.download(ApiClient.backups.getDownloadURL(token, name));
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
            }
        }

        delete isDownloading[name];
        isDownloading = isDownloading;
    }

    function deleteConfirm(name) {
        confirm(`Do you really want to delete ${name}?`, () => deleteBackup(name));
    }

    async function deleteBackup(name) {
        if (isDeleting[name]) {
            return;
        }

        isDeleting[name] = true;

        try {
            await ApiClient.backups.delete(name);
            CommonHelper.removeByKey(backups, "name", name);
            loadBackups();
            addSuccessToast(`Successfully deleted ${name}.`);
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
            }
        }

        delete isDeleting[name];
        isDeleting = isDeleting;
    }

    async function loadCanBackup() {
        try {
            const health = await ApiClient.health.check({ $autoCancel: false });
            const oldCanBackup = canBackup;
            canBackup = health?.data?.canBackup || false;

            // reload backups list
            if (oldCanBackup != canBackup && canBackup) {
                loadBackups();
            }
        } catch (_) {}
    }

    onMount(() => {
        let canBackupIntervalId = setInterval(() => {
            loadCanBackup();
        }, 3000);

        return () => {
            clearInterval(canBackupIntervalId);
        };
    });
</script>

<div class="list list-compact">
    <div class="list-content">
        {#if isLoading}
            {#each Array(backups.length || 1) as i}
                <div class="list-item list-item-loader">
                    <span class="skeleton-loader" />
                </div>
            {/each}
        {:else}
            {#each backups as backup (backup.key)}
                <div class="list-item" transition:slide={{ duration: 150 }}>
                    <i class="ri-folder-zip-line" />
                    <div class="content">
                        <span class="name backup-name" title={backup.key}>{backup.key}</span>
                        <span class="size txt-hint txt-nowrap">
                            ({CommonHelper.formattedFileSize(backup.size)})
                        </span>
                    </div>
                    <div class="actions nonintrusive">
                        <button
                            type="button"
                            class="btn btn-sm btn-circle btn-hint btn-transparent"
                            class:btn-loading={isDownloading[backup.key]}
                            disabled={isDeleting[backup.key] || isDownloading[backup.key]}
                            aria-label="Download"
                            use:tooltip={"Download"}
                            on:click|preventDefault={() => download(backup.key)}
                        >
                            <i class="ri-download-line" />
                        </button>
                        <button
                            type="button"
                            class="btn btn-sm btn-circle btn-hint btn-transparent"
                            disabled={isDeleting[backup.key]}
                            aria-label="Restore"
                            use:tooltip={"Restore"}
                            on:click|preventDefault={() => restorePanel.show(backup.key)}
                        >
                            <i class="ri-restart-line" />
                        </button>
                        <button
                            type="button"
                            class="btn btn-sm btn-circle btn-hint btn-transparent"
                            class:btn-loading={isDeleting[backup.key]}
                            disabled={isDeleting[backup.key]}
                            aria-label="Delete"
                            use:tooltip={"Delete"}
                            on:click|preventDefault={() => deleteConfirm(backup.key)}
                        >
                            <i class="ri-delete-bin-7-line" />
                        </button>
                    </div>
                </div>
            {:else}
                <div class="list-item list-item-placeholder">
                    <span class="txt">No backups yet.</span>
                </div>
            {/each}
        {/if}
    </div>

    <div class="list-item list-item-btn">
        <button
            type="button"
            class="btn btn-block btn-transparent"
            disabled={isLoading || !canBackup}
            on:click={() => createPanel?.show()}
        >
            {#if canBackup}
                <i class="ri-play-circle-line" />
                <span class="txt">Initialize new backup</span>
            {:else}
                <span class="loader loader-sm" />
                <span class="txt">Backup/restore operation is in process</span>
            {/if}
        </button>
    </div>
</div>

<BackupCreatePanel
    bind:this={createPanel}
    on:submit={() => {
        loadBackups();
    }}
/>

<BackupRestorePanel bind:this={restorePanel} />

<style lang="scss">
    .list-content {
        overflow: auto;
        max-height: 342px;
        .list-item {
            min-height: 49px;
        }
    }
    .backup-name {
        max-width: 300px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
    }
</style>
