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

        formSettings.rateLimits.rules = sortRules(formSettings.rateLimits.rules);

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
            rateLimits: settings.rateLimits || { rules: [] },
        };

        sortRules(formSettings.rateLimits.rules);

        originalFormSettings = JSON.parse(JSON.stringify(formSettings));
    }

    function reset() {
        formSettings = JSON.parse(JSON.stringify(originalFormSettings || {}));
    }

    // sort the specified rules list in place
    function sortRules(rules) {
        if (!rules) {
            return;
        }

        let compare = [{}, {}];

        rules.sort((a, b) => {
            compare[0].length = a.label.length;
            compare[0].isTag = a.label.includes(":") || !a.label.includes("/");
            compare[0].isWildcardTag = compare[0].isTag && a.label.startsWith("*");
            compare[0].isExactTag = compare[0].isTag && !compare[0].isWildcardTag;
            compare[0].isPrefix = !compare[0].isTag && a.label.endsWith("/");
            compare[0].hasMethod = !compare[0].isTag && a.label.includes(" /");

            compare[1].length = b.label.length;
            compare[1].isTag = b.label.includes(":") || !b.label.includes("/");
            compare[1].isWildcardTag = compare[1].isTag && b.label.startsWith("*");
            compare[1].isExactTag = compare[1].isTag && !compare[1].isWildcardTag;
            compare[1].isPrefix = !compare[1].isTag && b.label.endsWith("/");
            compare[1].hasMethod = !compare[1].isTag && b.label.includes(" /");

            for (let item of compare) {
                item.priority = 0; // reset

                if (item.isTag) {
                    item.priority += 1000;

                    if (item.isExactTag) {
                        item.priority += 10;
                    } else {
                        item.priority += 5;
                    }
                } else {
                    if (item.hasMethod) {
                        item.priority += 10;
                    }

                    if (!item.isPrefix) {
                        item.priority += 5;
                    }
                }
            }
            // sort additionally prefix paths based on their length
            if (
                compare[0].isPrefix &&
                compare[1].isPrefix &&
                ((compare[0].hasMethod && compare[1].hasMethod) ||
                    (!compare[0].hasMethod && !compare[1].hasMethod))
            ) {
                if (compare[0].length > compare[1].length) {
                    compare[0].priority += 1;
                } else if (compare[0].length < compare[1].length) {
                    compare[1].priority += 1;
                }
            }

            if (compare[0].priority > compare[1].priority) {
                return -1;
            }

            if (compare[0].priority < compare[1].priority) {
                return 1;
            }

            return 0;
        });

        return rules;
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
