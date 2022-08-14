<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import EmailAuthAccordion from "@/components/settings/EmailAuthAccordion.svelte";
    import AuthProviderAccordion from "@/components/settings/AuthProviderAccordion.svelte";

    $pageTitle = "Auth providers";

    let emailAuthAccordion;
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

        formSettings = {
            emailAuth: Object.assign({ enabled: true }, data.emailAuth),
        };

        const providers = ["googleAuth", "facebookAuth", "githubAuth", "gitlabAuth"];
        for (const provider of providers) {
            formSettings[provider] = Object.assign(
                { enabled: false, allowRegistrations: true },
                data[provider]
            );
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
            <h6 class="m-b-base">Manage the allowed users sign-in/sign-up methods.</h6>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="accordions">
                    <EmailAuthAccordion
                        bind:this={emailAuthAccordion}
                        single
                        bind:config={formSettings.emailAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="googleAuth"
                        title="Google"
                        icon="ri-google-line"
                        bind:config={formSettings.googleAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="facebookAuth"
                        title="Facebook"
                        icon="ri-facebook-line"
                        bind:config={formSettings.facebookAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="githubAuth"
                        title="GitHub"
                        icon="ri-github-line"
                        bind:config={formSettings.githubAuth}
                    />
                    <AuthProviderAccordion
                        single
                        key="gitlabAuth"
                        title="GitLab"
                        icon="ri-gitlab-line"
                        showSelfHostedFields
                        bind:config={formSettings.gitlabAuth}
                    />
                </div>

                <div class="flex m-t-base">
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
