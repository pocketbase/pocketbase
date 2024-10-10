<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import { addInfoToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";
    import Field from "@/components/base/Field.svelte";

    const dispatch = createEventDispatcher();

    let email = "";
    let password = "";
    let passwordConfirm = "";
    let isLoading = false;
    let isUploading = false;

    let backupFileInput;

    $: isBusy = isLoading || isUploading;

    async function submit() {
        if (isBusy) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.collection("_superusers").create({
                email,
                password,
                passwordConfirm,
            });

            await ApiClient.collection("_superusers").authWithPassword(email, password);

            dispatch("submit");
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    function resetSelectedBackupFile() {
        if (backupFileInput) {
            backupFileInput.value = "";
        }
    }

    function uploadBackupConfirm(file) {
        if (!file) {
            return;
        }

        confirm(
            `Note that we don't perform validations for the uploaded backup files. Proceed with caution and only if you trust the file source.\n\n` +
                `Do you really want to upload and initialize "${file.name}"?`,
            () => {
                uploadBackup(file);
            },
            () => {
                resetSelectedBackupFile();
            },
        );
    }

    async function uploadBackup(file) {
        if (!file || isBusy) {
            return;
        }

        isUploading = true;

        try {
            await ApiClient.backups.upload({ file: file });

            await ApiClient.backups.restore(file.name);

            addInfoToast("Please wait while extracting the uploaded archive!");

            // optimistic restore completion
            await new Promise((r) => setTimeout(r, 2000));

            dispatch("submit");
        } catch (err) {
            ApiClient.error(err);
        }

        resetSelectedBackupFile();

        isUploading = false;
    }
</script>

<form class="block" autocomplete="off" on:submit|preventDefault={submit}>
    <div class="content txt-center m-b-base">
        <h4>Create your first superuser account in order to continue</h4>
    </div>

    <Field class="form-field required" name="email" let:uniqueId>
        <label for={uniqueId}>Email</label>
        <!-- svelte-ignore a11y-autofocus -->
        <input type="email" autocomplete="off" id={uniqueId} bind:value={email} required autofocus />
    </Field>

    <Field class="form-field required" name="password" let:uniqueId>
        <label for={uniqueId}>Password</label>
        <input
            type="password"
            autocomplete="new-password"
            minlength="10"
            id={uniqueId}
            bind:value={password}
            required
        />
        <div class="help-block">Recommended at least 10 characters.</div>
    </Field>

    <Field class="form-field required" name="passwordConfirm" let:uniqueId>
        <label for={uniqueId}>Password confirm</label>
        <input type="password" minlength="10" id={uniqueId} bind:value={passwordConfirm} required />
    </Field>

    <button
        type="submit"
        class="btn btn-lg btn-block btn-next"
        class:btn-disabled={isBusy}
        class:btn-loading={isLoading}
    >
        <span class="txt">Create superuser and login</span>
        <i class="ri-arrow-right-line" />
    </button>
</form>

<hr />

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<label
    for="backupFileInput"
    class="btn btn-lg btn-hint btn-transparent btn-block"
    class:btn-disabled={isBusy}
    class:btn-loading={isUploading}
>
    <i class="ri-upload-cloud-line" />
    <span class="txt">Or initialize from backup</span>
</label>
<input
    bind:this={backupFileInput}
    id="backupFileInput"
    type="file"
    class="hidden"
    accept=".zip"
    on:change={(e) => {
        uploadBackupConfirm(e.target?.files?.[0]);
    }}
/>
