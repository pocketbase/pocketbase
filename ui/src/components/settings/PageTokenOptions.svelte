<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Field from "@/components/base/Field.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    const tokensList = [
        { key: "recordAuthToken", label: "Auth record authentication token" },
        { key: "recordVerificationToken", label: "Auth record email verification token" },
        { key: "recordPasswordResetToken", label: "Auth record password reset token" },
        { key: "recordEmailChangeToken", label: "Auth record email change token" },
        { key: "adminAuthToken", label: "Admins auth token" },
        { key: "adminPasswordResetToken", label: "Admins password reset token" },
    ];

    $pageTitle = "Token options";

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
            const result = (await ApiClient.settings.getAll()) || {};
            initSettings(result);
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
            const result = await ApiClient.settings.update(CommonHelper.filterRedactedProps(formSettings));
            initSettings(result);
            addSuccessToast("Successfully saved tokens options.");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function initSettings(data) {
        data = data || {};
        formSettings = {};

        for (const listItem of tokensList) {
            formSettings[listItem.key] = {
                duration: data[listItem.key]?.duration || 0,
            };
        }

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
        <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
            <div class="content m-b-sm txt-xl">
                <p>Adjust common token options.</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                {#each tokensList as token (token.key)}
                    <Field class="form-field required" name="{token.key}.duration" let:uniqueId>
                        <label for={uniqueId}>{token.label} duration (in seconds)</label>
                        <input
                            type="number"
                            id={uniqueId}
                            required
                            bind:value={formSettings[token.key].duration}
                        />
                        <div class="help-block">
                            <span
                                class="link-primary"
                                class:txt-success={formSettings[token.key].secret}
                                on:click={() => {
                                    // toggle
                                    if (formSettings[token.key].secret) {
                                        delete formSettings[token.key].secret;
                                        formSettings[token.key] = formSettings[token.key];
                                    } else {
                                        formSettings[token.key].secret = CommonHelper.randomString(50);
                                    }
                                }}
                            >
                                Invalidate all previously issued tokens
                            </span>
                        </div>
                    </Field>
                {/each}

                <div class="flex">
                    <div class="flex-fill" />
                    {#if hasChanges}
                        <button
                            type="button"
                            class="btn btn-secondary btn-hint"
                            disabled={isSaving}
                            on:click={() => reset()}
                        >
                            <span class="txt">Cancel</span>
                        </button>
                    {/if}
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
