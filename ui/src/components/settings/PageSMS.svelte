<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { slide } from "svelte/transition";

    $pageTitle = "SMS settings";

    let testPopup;
    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;

    $: initialHash = JSON.stringify(originalFormSettings);

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.settings.getAll()) || {};
            init(settings);
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    async function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        try {
            const settings = await ApiClient.settings.update(CommonHelper.filterRedactedProps(formSettings));
            init(settings);
            setErrors({});
            addSuccessToast("Successfully saved SMS settings.");
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        formSettings = {
            meta: settings?.meta || {},
            sms: settings?.sms || {},
        };

        originalFormSettings = JSON.parse(JSON.stringify(formSettings));
    }

    function reset() {
        formSettings = JSON.parse(JSON.stringify(originalFormSettings || {}));
    }
</script>

<SettingsSidebar />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={() => save()}>
            <div class="content txt-xl m-b-base">
                <p>Configure common settings for sending SMS.</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <Field class="form-field form-field-toggle m-b-sm" let:uniqueId>
                    <input type="checkbox" id={uniqueId} required bind:checked={formSettings.sms.enabled} />
                    <label for={uniqueId}>
                        <span class="txt">Use Twillio SMS</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: 'By default PocketBase uses the unix "sendmail" command for sending emails. For better emails deliverability it is recommended to use a SMTP mail server.',
                                position: "top",
                            }}
                        />
                    </label>
                </Field>

                {#if formSettings.sms.enabled}
                    <div transition:slide={{ duration: 150 }}>
                        <div class="grid">
                            <div class="col-lg-4">
                                <Field class="form-field required" name="sms.accountSID" let:uniqueId>
                                    <label for={uniqueId}>Twillio Account SID</label>
                                    <input
                                        type="text"
                                        id={uniqueId}
                                        required
                                        bind:value={formSettings.sms.accountSID}
                                    />
                                </Field>
                            </div>
                            <div class="col-lg-4">
                                <Field
                                    class="form-field
                                    required"
                                    name="sms.authToken"
                                    let:uniqueId
                                >
                                    <label for={uniqueId}>Twillio Auth Token</label>
                                    <RedactedPasswordInput
                                        id={uniqueId}
                                        required
                                        bind:value={formSettings.sms.authToken}
                                    />
                                </Field>
                            </div>
                            <div class="col-lg-4">
                                <Field
                                    class="form-field
                                    required"
                                    name="sms.fromNumber"
                                    let:uniqueId
                                >
                                    <label for={uniqueId}>From Number</label>
                                    <input
                                        type="text"
                                        id={uniqueId}
                                        required
                                        bind:value={formSettings.sms.fromNumber}
                                    />
                                </Field>
                            </div>
                        </div>
                    </div>
                {/if}

                <div class="flex">
                    <div class="flex-fill" />
                    <button
                        type="button"
                        class="btn btn-transparent btn-hint"
                        disabled={isSaving}
                        on:click={() => reset()}
                    >
                        <span class="txt">Cancel</span>
                    </button>
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
</PageWrapper>
