<script>
    import { slide } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    const tlsOptions = [
        { label: "Optional (StartTLS)", value: false },
        { label: "Always", value: true },
    ];

    $pageTitle = "Mail settings";

    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.settings.getAll()) || {};
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
            const settings = await ApiClient.settings.update(CommonHelper.filterRedactedProps(formSettings));
            init(settings);
            addSuccessToast("Successfully saved mail settings.");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        formSettings = {
            meta: settings?.meta || {},
            smtp: settings?.smtp || {},
        };
        initialHash = JSON.stringify(formSettings);
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
                <p>Configure common settings for sending emails.</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="grid">
                    <div class="col-lg-6">
                        <Field class="form-field required" name="meta.senderName" let:uniqueId>
                            <label for={uniqueId}>Sender name</label>
                            <input
                                type="text"
                                id={uniqueId}
                                required
                                bind:value={formSettings.meta.senderName}
                            />
                        </Field>
                    </div>

                    <div class="col-lg-6">
                        <Field class="form-field required" name="meta.senderAddress" let:uniqueId>
                            <label for={uniqueId}>Sender address</label>
                            <input
                                type="email"
                                id={uniqueId}
                                required
                                bind:value={formSettings.meta.senderAddress}
                            />
                        </Field>
                    </div>

                    <Field class="form-field required" name="meta.userVerificationUrl" let:uniqueId>
                        <label for={uniqueId}>User verification page url</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.meta.userVerificationUrl}
                        />
                        <div class="help-block">
                            Used in the user verification email. Available placeholder parameters:
                            <code>%APP_URL%</code>, <code>%TOKEN%</code>.
                        </div>
                    </Field>

                    <Field class="form-field required" name="meta.userResetPasswordUrl" let:uniqueId>
                        <label for={uniqueId}>User reset password page url</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.meta.userResetPasswordUrl}
                        />
                        <div class="help-block">
                            Used in the user password reset email. Available placeholder parameters:
                            <code>%APP_URL%</code>, <code>%TOKEN%</code>.
                        </div>
                    </Field>

                    <Field class="form-field required" name="meta.userConfirmEmailChangeUrl" let:uniqueId>
                        <label for={uniqueId}>User confirm email change page url</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.meta.userConfirmEmailChangeUrl}
                        />
                        <div class="help-block">
                            Used in the user email change confirmation email. Available placeholder
                            parameters:
                            <code>%APP_URL%</code>, <code>%TOKEN%</code>.
                        </div>
                    </Field>
                </div>

                <hr />

                <div class="content m-b-sm">
                    <p>
                        By default PocketBase uses the OS <code>sendmail</code> command for sending emails.
                        <br />
                        <strong class="txt-bold">
                            For better emails deliverability it is recommended to enable the SMTP settings
                            below.
                        </strong>
                    </p>
                </div>

                <Field class="form-field form-field-toggle" let:uniqueId>
                    <input type="checkbox" id={uniqueId} required bind:checked={formSettings.smtp.enabled} />
                    <label for={uniqueId}>Use SMTP mail server</label>
                </Field>

                {#if formSettings.smtp.enabled}
                    <div class="grid" transition:slide|local={{ duration: 150 }}>
                        <div class="col-lg-6">
                            <Field class="form-field required" name="smtp.host" let:uniqueId>
                                <label for={uniqueId}>SMTP server host</label>
                                <input
                                    type="text"
                                    id={uniqueId}
                                    required
                                    bind:value={formSettings.smtp.host}
                                />
                            </Field>
                        </div>
                        <div class="col-lg-3">
                            <Field class="form-field required" name="smtp.port" let:uniqueId>
                                <label for={uniqueId}>Port</label>
                                <input
                                    type="number"
                                    id={uniqueId}
                                    required
                                    bind:value={formSettings.smtp.port}
                                />
                            </Field>
                        </div>
                        <div class="col-lg-3">
                            <Field class="form-field required" name="smtp.tls" let:uniqueId>
                                <label for={uniqueId}>TLS Encryption</label>
                                <ObjectSelect
                                    id={uniqueId}
                                    items={tlsOptions}
                                    bind:keyOfSelected={formSettings.smtp.tls}
                                />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field" name="smtp.username" let:uniqueId>
                                <label for={uniqueId}>Username</label>
                                <input type="text" id={uniqueId} bind:value={formSettings.smtp.username} />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field" name="smtp.password" let:uniqueId>
                                <label for={uniqueId}>Password</label>
                                <RedactedPasswordInput
                                    id={uniqueId}
                                    bind:value={formSettings.smtp.password}
                                />
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
</PageWrapper>
