<script>
    import { slide } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { _ } from '@/services/i18n';

    const tlsOptions = [
        { label: $_("settings.mail.form.smtp.StartTLS"), value: false },
        { label: $_("settings.mail.form.smtp.AlwaysTLS"), value: true },
    ];

    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    CommonHelper.setDocumentTitle($_("settings.mail.pagetitle"));

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
            const settings = await ApiClient.Settings.update(CommonHelper.filterRedactedProps(formSettings));
            init(settings);
            addSuccessToast($_("settings.mail.tips.saved"));
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

<main class="page-wrapper">
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">{$_("app.breadcrumb.settings")}</div>
            <div class="breadcrumb-item">{$_("app.breadcrumb.mail_settings")}</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={() => save()}>
            <div class="content txt-xl m-b-base">
                <p>{$_("settings.mail.tips.title")}</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="grid">
                    <div class="col-lg-6">
                        <Field class="form-field required" name="meta.senderName" let:uniqueId>
                            <label for={uniqueId}>{$_("settings.mail.form.senderName")}</label>
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
                            <label for={uniqueId}>{$_("settings.mail.form.senderAddress")}</label>
                            <input
                                type="email"
                                id={uniqueId}
                                required
                                bind:value={formSettings.meta.senderAddress}
                            />
                        </Field>
                    </div>

                    <Field class="form-field required" name="meta.userVerificationUrl" let:uniqueId>
                        <label for={uniqueId}>{$_("settings.mail.form.userVerificationUrl")}</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.meta.userVerificationUrl}
                        />
                        <div class="help-block">
                            {$_("settings.mail.tips.help.userVerificationUrl")}:
                            <code>%APP_URL%</code>, <code>%TOKEN%</code>.
                        </div>
                    </Field>

                    <Field class="form-field required" name="meta.userResetPasswordUrl" let:uniqueId>
                        <label for={uniqueId}>{$_("settings.mail.form.userResetPasswordUrl")}</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.meta.userResetPasswordUrl}
                        />
                        <div class="help-block">
                            {$_("settings.mail.tips.help.userResetPasswordUrl")}:
                            <code>%APP_URL%</code>, <code>%TOKEN%</code>.
                        </div>
                    </Field>

                    <Field class="form-field required" name="meta.userConfirmEmailChangeUrl" let:uniqueId>
                        <label for={uniqueId}>{$_("settings.mail.form.userConfirmEmailChangeUrl")}</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.meta.userConfirmEmailChangeUrl}
                        />
                        <div class="help-block">
                            {$_("settings.mail.tips.help.userConfirmEmailChangeUrl")}:
                            <code>%APP_URL%</code>, <code>%TOKEN%</code>.
                        </div>
                    </Field>
                </div>

                <hr />

                <div class="content m-b-sm">
                    <p>
                        {$_("settings.mail.tips.help.smtp")}
                        <br />
                        <strong class="txt-bold">
                            {$_("settings.mail.tips.help.smtpenable")}
                        </strong>
                    </p>
                </div>

                <Field class="form-field form-field-toggle" let:uniqueId>
                    <input type="checkbox" id={uniqueId} required bind:checked={formSettings.smtp.enabled} />
                    <label for={uniqueId}>{$_("settings.mail.form.usesmtp")}</label>
                </Field>

                {#if formSettings.smtp.enabled}
                    <div class="grid" transition:slide|local={{ duration: 150 }}>
                        <div class="col-lg-6">
                            <Field class="form-field required" name="smtp.host" let:uniqueId>
                                <label for={uniqueId}>{$_("settings.mail.form.smtp.host")}</label>
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
                                <label for={uniqueId}>{$_("settings.mail.form.smtp.port")}</label>
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
                                <label for={uniqueId}>{$_("settings.mail.form.smtp.tls")}</label>
                                <ObjectSelect
                                    id={uniqueId}
                                    items={tlsOptions}
                                    bind:keyOfSelected={formSettings.smtp.tls}
                                />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field" name="smtp.username" let:uniqueId>
                                <label for={uniqueId}>{$_("settings.mail.form.smtp.username")}</label>
                                <input type="text" id={uniqueId} bind:value={formSettings.smtp.username} />
                            </Field>
                        </div>
                        <div class="col-lg-6">
                            <Field class="form-field" name="smtp.password" let:uniqueId>
                                <label for={uniqueId}>{$_("settings.mail.form.smtp.password")}</label>
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
                        <span class="txt">{$_("settings.mail.form.save")}</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</main>
