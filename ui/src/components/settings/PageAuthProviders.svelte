<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import EmailAuthAccordion from "@/components/settings/EmailAuthAccordion.svelte";
    import AuthProviderAccordion from "@/components/settings/AuthProviderAccordion.svelte";

    let emailAuthAccordion;
    let authSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";

    $: hasChanges = initialHash != JSON.stringify(authSettings);

    CommonHelper.setDocumentTitle("Auth providers");

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const result = (await ApiClient.Settings.getAll()) || {};
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
            const result = await ApiClient.Settings.update(CommonHelper.filterRedactedProps(authSettings));
            initSettings(result);
            setErrors({});
            emailAuthAccordion?.collapseSiblings();
            addSuccessToast("Successfully updated auth providers.");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function initSettings(data) {
        data = data || {};

        authSettings = {};
        authSettings.emailAuth = Object.assign({ enabled: true }, data.emailAuth);

        const providers = ["googleAuth", "facebookAuth", "githubAuth", "gitlabAuth", "stravaAuth"];
        for (const provider of providers) {
            authSettings[provider] = Object.assign(
                { enabled: false, allowRegistrations: true },
                data[provider]
            );
        }

        initialHash = JSON.stringify(authSettings);
    }
</script>

<SettingsSidebar />

<main class="page-wrapper">
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">Auth providers</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
            <h6 class="m-b-base">Manage the allowed users sign-in/sign-up methods.</h6>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="accordions">
                    <EmailAuthAccordion
                        bind:this={emailAuthAccordion}
                        single
                        bind:config={authSettings.emailAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="googleAuth"
                        title="Google"
                        icon="ri-google-line"
                        bind:config={authSettings.googleAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="facebookAuth"
                        title="Facebook"
                        icon="ri-facebook-line"
                        bind:config={authSettings.facebookAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="githubAuth"
                        title="GitHub"
                        icon="ri-github-line"
                        bind:config={authSettings.githubAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="gitlabAuth"
                        title="GitLab"
                        icon="ri-gitlab-line"
                        showSelfHostedFields
                        bind:config={authSettings.gitlabAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="stravaAuth"
                        title="Strava"
                        icon="ri-user-6-line"
                        showSelfHostedFields
                        bind:config={authSettings.stravaAuth}
                    />
                </div>

                <div class="flex m-t-base">
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
