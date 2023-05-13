<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import TokenField from "@/components/settings/TokenField.svelte";

    const recordTokensList = [
        { key: "recordAuthToken", label: "Auth record authentication token" },
        { key: "recordVerificationToken", label: "Auth record email verification token" },
        { key: "recordPasswordResetToken", label: "Auth record password reset token" },
        { key: "recordEmailChangeToken", label: "Auth record email change token" },
        { key: "recordFileToken", label: "Records protected file access token" },
    ];

    const adminTokensList = [
        { key: "adminAuthToken", label: "Admins auth token" },
        { key: "adminPasswordResetToken", label: "Admins password reset token" },
        { key: "adminFileToken", label: "Admins protected file access token" },
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
            const result = await ApiClient.settings.update(CommonHelper.filterRedactedProps(formSettings));
            initSettings(result);
            addSuccessToast("Successfully saved tokens options.");
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function initSettings(data) {
        data = data || {};
        formSettings = {};

        const tokensList = recordTokensList.concat(adminTokensList);

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
                <h3 class="section-title">Record tokens</h3>
                {#each recordTokensList as token (token.key)}
                    <TokenField
                        key={token.key}
                        label={token.label}
                        bind:duration={formSettings[token.key].duration}
                        bind:secret={formSettings[token.key].secret}
                    />
                {/each}

                <hr />

                <h3 class="section-title">Admin tokens</h3>
                {#each adminTokensList as token (token.key)}
                    <TokenField
                        key={token.key}
                        label={token.label}
                        bind:duration={formSettings[token.key].duration}
                        bind:secret={formSettings[token.key].secret}
                    />
                {/each}

                <div class="flex">
                    <div class="flex-fill" />
                    {#if hasChanges}
                        <button
                            type="button"
                            class="btn btn-transparent btn-hint"
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
