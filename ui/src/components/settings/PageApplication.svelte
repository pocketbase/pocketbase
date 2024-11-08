<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { addSuccessToast } from "@/stores/toasts";
    import { appName, hideControls, pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import Field from "@/components/base/Field.svelte";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import BatchAccordion from "@/components/settings/BatchAccordion.svelte";
    import TrustedProxyAccordion from "@/components/settings/TrustedProxyAccordion.svelte";
    import RateLimitAccordion from "@/components/settings/RateLimitAccordion.svelte";

    $pageTitle = "Application settings";

    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";
    let healthData = {};

    $: initialHash = JSON.stringify(originalFormSettings);

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    loadSettings();

    async function loadHealthData() {
        try {
            healthData = ((await ApiClient.health.check()) || {})?.data || {};
        } catch (err) {
            console.warn("Health check failed:", err);
        }
    }

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.settings.getAll()) || {};
            init(settings);

            await loadHealthData();
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

            await loadHealthData();

            setErrors({});

            addSuccessToast("Successfully saved application settings.");
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        $appName = settings?.meta?.appName;
        $hideControls = !!settings?.meta?.hideControls;

        formSettings = {
            meta: settings?.meta || {},
            batch: settings.batch || {},
            trustedProxy: settings.trustedProxy || { headers: [] },
            rateLimits: settings.rateLimits || { tags: [] },
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
            <div class="breadcrumb-item">Application</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="grid">
                    <div class="col-lg-6">
                        <Field class="form-field required" name="meta.appName" let:uniqueId>
                            <label for={uniqueId}>Application name</label>
                            <input
                                type="text"
                                id={uniqueId}
                                required
                                bind:value={formSettings.meta.appName}
                            />
                        </Field>
                    </div>

                    <div class="col-lg-6">
                        <Field class="form-field required" name="meta.appURL" let:uniqueId>
                            <label for={uniqueId}>Application URL</label>
                            <input type="text" id={uniqueId} required bind:value={formSettings.meta.appURL} />
                        </Field>
                    </div>
                    <div class="col-lg-12">
                        <div class="accordions">
                            <TrustedProxyAccordion bind:formSettings {healthData} />
                            <RateLimitAccordion bind:formSettings />
                            <BatchAccordion bind:formSettings />
                        </div>
                    </div>
                    <div class="col-lg-12">
                        <Field class="form-field form-field-toggle m-0" name="meta.hideControls" let:uniqueId>
                            <input
                                type="checkbox"
                                id={uniqueId}
                                bind:checked={formSettings.meta.hideControls}
                            />
                            <label for={uniqueId}>
                                <span class="txt">Hide collection create and edit controls</span>
                                <i
                                    class="ri-information-line link-hint"
                                    use:tooltip={{
                                        text: `This could prevent making accidental schema changes when in production environment.`,
                                        position: "right",
                                    }}
                                />
                            </label>
                        </Field>
                    </div>
                </div>

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
                        on:click={() => save()}
                    >
                        <span class="txt">Save changes</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</PageWrapper>
