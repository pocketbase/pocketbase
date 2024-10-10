<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { setErrors } from "@/stores/errors";
    import { addErrorToast } from "@/stores/toasts";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Field from "@/components/base/Field.svelte";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
    import { onDestroy } from "svelte";

    const formId = "backup_restore_" + CommonHelper.randomString(5);

    let panel;
    let name = "";
    let nameConfirm = "";
    let isSubmitting = false;
    let reloadTimeoutId = null;

    $: canSubmit = nameConfirm != "" && name == nameConfirm;

    export function show(backupName) {
        setErrors({});
        nameConfirm = "";
        name = backupName;
        isSubmitting = false;
        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    async function submit() {
        if (!canSubmit || isSubmitting) {
            return;
        }

        clearTimeout(reloadTimeoutId);

        isSubmitting = true;

        try {
            await ApiClient.backups.restore(name);

            // optimistic restore page reload
            reloadTimeoutId = setTimeout(() => {
                window.location.reload();
            }, 2000);
        } catch (err) {
            clearTimeout(reloadTimeoutId);

            if (!err?.isAbort) {
                isSubmitting = false;
                addErrorToast(err.response?.message || err.message);
            }
        }
    }

    onDestroy(() => {
        clearTimeout(reloadTimeoutId);
    });
</script>

<OverlayPanel
    bind:this={panel}
    class="backup-restore-panel"
    overlayClose={!isSubmitting}
    escClose={!isSubmitting}
    beforeHide={() => !isSubmitting}
    popup
    on:show
    on:hide
>
    <svelte:fragment slot="header">
        <h4 class="popup-title txt-ellipsis">Restore <strong>{name}</strong></h4>
    </svelte:fragment>

    <div class="alert alert-danger">
        <div class="icon">
            <i class="ri-alert-line" />
        </div>
        <div class="content">
            <p class="txt-bold">Please proceed with caution and use it only with trusted backups!</p>

            <p>Backup restore is experimental and works only on UNIX based systems.</p>
            <p>
                The restore operation will attempt to replace your existing <code>pb_data</code> with the one from
                the backup and will restart the application process.
            </p>
            <p>
                This means that on success all of your data (including app settings, users, superusers, etc.) will
                be replaced with the ones from the backup.
            </p>
            <p>
                Nothing will happen if the backup is invalid or incompatible (ex. missing
                <code>data.db</code> file).
            </p>
        </div>
    </div>

    <div class="content m-b-xs">
        Type the backup name
        <div class="label">
            <span class="txt">{name}</span>
            <CopyIcon value={name} />
        </div>
        to confirm:
    </div>

    <form id={formId} autocomplete="off" on:submit|preventDefault={submit}>
        <Field class="form-field required m-0" name="name" let:uniqueId>
            <label for={uniqueId}>Backup name</label>
            <input type="text" id={uniqueId} required bind:value={nameConfirm} />
        </Field>
    </form>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide} disabled={isSubmitting}>
            Cancel
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSubmitting}
            disabled={!canSubmit || isSubmitting}
        >
            <span class="txt">Restore backup</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style>
    .popup-title {
        max-width: 80%;
    }
</style>
