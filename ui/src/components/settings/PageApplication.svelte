<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import setUiTheme from "@/utils/UiTheme";
    import { pageTitle, appName, hideControls, uiTheme } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import tooltip from "@/actions/tooltip";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    $pageTitle = "Application settings";

    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";
    let themeOptions = ["System", "Light", "Dark"];

    $: initialHash = JSON.stringify(originalFormSettings);

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
            addSuccessToast("Successfully saved application settings.");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        $appName = settings?.meta?.appName;
        $hideControls = !!settings?.meta?.hideControls;
        $uiTheme = settings?.meta?.uiTheme;
        setUiTheme($uiTheme);

        formSettings = {
            meta: settings?.meta || {},
            logs: settings?.logs || {},
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
                        <Field class="form-field required" name="meta.appUrl" let:uniqueId>
                            <label for={uniqueId}>Application url</label>
                            <input type="text" id={uniqueId} required bind:value={formSettings.meta.appUrl} />
                        </Field>
                    </div>

                    <Field class="form-field required" name="logs.maxDays" let:uniqueId>
                        <label for={uniqueId}>Logs max days retention</label>
                        <input type="number" id={uniqueId} required bind:value={formSettings.logs.maxDays} />
                    </Field>

                    <Field class="form-field" name="meta.uiTheme" let:uniqueId>
                        <label for={uniqueId}>
                            <span class="txt">Set UI Theme</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: "You can either choose to let the system handle the theme automatically or set it to your preferred one.",
                                    position: "top",
                                }}
                            />
                        </label>
                        <Select
                            id={uniqueId}
                            items={themeOptions}
                            bind:selected={formSettings.meta.uiTheme}
                        />
                    </Field>

                    <Field class="form-field form-field-toggle" name="meta.hideControls" let:uniqueId>
                        <input type="checkbox" id={uniqueId} bind:checked={formSettings.meta.hideControls} />
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

                    <div class="col-lg-12 flex">
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
                </div>
            {/if}
        </form>
    </div>
</PageWrapper>
