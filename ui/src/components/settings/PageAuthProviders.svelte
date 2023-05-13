<script>
    import ApiClient from "@/utils/ApiClient";
    import { pageTitle } from "@/stores/app";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import AuthProviderCard from "@/components/settings/AuthProviderCard.svelte";
    import providersList from "@/providers.js";

    $pageTitle = "Auth providers";

    let isLoading = false;
    let formSettings = {};

    $: enabledProviders = providersList.filter((provider) => formSettings[provider.key]?.enabled);

    $: disabledProviders = providersList.filter((provider) => !formSettings[provider.key]?.enabled);

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

    function initSettings(data) {
        data = data || {};
        formSettings = {};

        for (const provider of providersList) {
            formSettings[provider.key] = Object.assign({ enabled: false }, data[provider.key]);
        }
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
        <div class="panel">
            <h6 class="m-b-base">Manage the allowed users OAuth2 sign-in/sign-up methods.</h6>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="grid grid-sm">
                    {#each enabledProviders as provider (provider.key)}
                        <div class="col-lg-6">
                            <AuthProviderCard {provider} bind:config={formSettings[provider.key]} />
                        </div>
                    {/each}
                </div>

                {#if enabledProviders.length > 0 && disabledProviders.length > 0}
                    <hr />
                {/if}

                <div class="grid grid-sm">
                    {#each disabledProviders as provider (provider.key)}
                        <div class="col-lg-6">
                            <AuthProviderCard {provider} bind:config={formSettings[provider.key]} />
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
</PageWrapper>
