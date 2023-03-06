<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import AuthProviderAccordion from "@/components/settings/AuthProviderAccordion.svelte";
    import providersList from "@/providers.js";

    $pageTitle = "Auth providers";

    let accordions = {};
    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let showHidden = false;

    $: initialHash = JSON.stringify(originalFormSettings);

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    $: totalHidden = Object.values(providersList).filter((provider) => provider.hidden).length;

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
            setErrors({});

            accordions[Object.keys(accordions)[0]]?.collapseSiblings();
            addSuccessToast("Successfully updated auth providers.");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function initSettings(data) {
        data = data || {};

        formSettings = {};

        for (const providerKey in providersList) {
            formSettings[providerKey] = Object.assign({ enabled: false }, data[providerKey]);
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
            <h6 class="m-b-base">Manage the allowed users OAuth2 sign-in/sign-up methods.</h6>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="accordions">
                    {#each Object.entries(providersList) as [key, provider]}
                        {#if showHidden || !provider.hidden || formSettings[key]?.enabled}
                            <AuthProviderAccordion
                                bind:this={accordions[key]}
                                single
                                {key}
                                title={provider.title}
                                icon={provider.icon || "ri-fingerprint-line"}
                                optionsComponent={provider.optionsComponent}
                                bind:config={formSettings[key]}
                            />
                        {/if}
                    {/each}
                </div>

                {#if !showHidden}
                    <button
                        type="button"
                        class="btn btn-sm btn-transparent btn-hint m-t-10"
                        on:click={() => (showHidden = true)}
                    >
                        <i class="ri-arrow-down-s-line" />
                        <span class="txt">Show all ({totalHidden})</span>
                    </button>
                {/if}

                <div class="flex m-t-base">
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
                    >
                        <span class="txt">Save changes</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</PageWrapper>
