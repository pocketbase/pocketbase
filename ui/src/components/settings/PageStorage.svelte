<script>
    import { slide } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    let s3 = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";
    let initialEnabled = false;

    $: hasChanges = initialHash != JSON.stringify(s3);

    CommonHelper.setDocumentTitle("Files storage");

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.Settings.getAll()) || {};
            init(settings);
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }

    async function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        try {
            const settings = await ApiClient.Settings.update(CommonHelper.filterRedactedProps({ s3 }));
            init(settings);
            setErrors({});
            addSuccessToast("Successfully saved Files storage settings.");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        s3 = settings?.s3 || {};
        initialEnabled = s3.enabled;
        initialHash = JSON.stringify(s3);
    }
</script>

<SettingsSidebar />

<main class="page-wrapper">
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">Files storage</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={() => save()}>
            <div class="content txt-xl m-b-base">
                <p>By default PocketBase uses the local file system to store uploaded files.</p>
                <p>
                    If you have limited disk space, you could optionally connect to a S3 compatible storage.
                </p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <Field class="form-field form-field-toggle" let:uniqueId>
                    <input type="checkbox" id={uniqueId} required bind:checked={s3.enabled} />
                    <label for={uniqueId}>Use S3 storage</label>
                </Field>

                {#if initialEnabled != s3.enabled}
                    <div transition:slide={{ duration: 150 }}>
                        <div class="alert alert-warning m-0">
                            <div class="icon">
                                <i class="ri-error-warning-line" />
                            </div>
                            <div class="content">
                                If you have existing uploaded files, you'll have to migrate them manually from
                                the
                                <strong>{initialEnabled ? "S3 storage" : "local file system"}</strong>
                                to the
                                <strong>{s3.enabled ? "S3 storage" : "local file system"}</strong>.
                                <br />
                                There are numerous command line tools that can help you, such as:
                                <a
                                    href="https://github.com/rclone/rclone"
                                    target="_blank"
                                    rel="noopener"
                                    class="txt-bold"
                                >
                                    rclone
                                </a>,
                                <a
                                    href="https://github.com/peak/s5cmd"
                                    target="_blank"
                                    rel="noopener"
                                    class="txt-bold"
                                >
                                    s5cmd
                                </a>, etc.
                            </div>
                        </div>
                        <div class="clearfix m-t-base" />
                    </div>
                {/if}

                {#if s3.enabled}
                    <div class="grid" transition:slide|local={{ duration: 150 }}>
                        <div class="col-lg-12">
                            <Field class="form-field required" name="s3.endpoint" let:uniqueId>
                                <label for={uniqueId}>Endpoint</label>
                                <input type="text" id={uniqueId} required bind:value={s3.endpoint} />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field required" name="s3.bucket" let:uniqueId>
                                <label for={uniqueId}>Bucket</label>
                                <input type="text" id={uniqueId} required bind:value={s3.bucket} />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field required" name="s3.region" let:uniqueId>
                                <label for={uniqueId}>Region</label>
                                <input type="text" id={uniqueId} required bind:value={s3.region} />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field required" name="s3.accessKey" let:uniqueId>
                                <label for={uniqueId}>Access key</label>
                                <input type="text" id={uniqueId} required bind:value={s3.accessKey} />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field required" name="s3.secret" let:uniqueId>
                                <label for={uniqueId}>Secret</label>
                                <RedactedPasswordInput id={uniqueId} required bind:value={s3.secret} />
                            </Field>
                        </div>
                        <!-- margin helper -->
                        <div class="col-lg-12" />
                    </div>
                {/if}

                <div class="flex">
                    <div class="flex-fill" />
                    <button
                        type="submit"
                        class="btn btn-expanded"
                        class:btn-loading={isSaving}
                        disabled={!hasChanges || isSaving}
                        on:click={() => save()}
                    >
                        <span class="txt">Save changes</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</main>
