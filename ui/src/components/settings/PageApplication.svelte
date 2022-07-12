<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { _ } from '@/services/i18n';

    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    CommonHelper.setDocumentTitle($_("settings.app.pagetitle"));

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
            addSuccessToast($_("settings.app.tips.saved"));
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        formSettings = {
            meta: settings?.meta || {},
            logs: settings?.logs || {},
        };
        initialHash = JSON.stringify(formSettings);
    }
</script>

<SettingsSidebar />

<main class="page-wrapper">
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">{$_("app.breadcrumb.settings")}</div>
            <div class="breadcrumb-item">{$_("app.breadcrumb.application")}</div>
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
                            <label for={uniqueId}>{$_("settings.app.form.appname")}</label>
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
                            <label for={uniqueId}>{$_("settings.app.form.appurl")}</label>
                            <input type="text" id={uniqueId} required bind:value={formSettings.meta.appUrl} />
                        </Field>
                    </div>

                    <Field class="form-field required" name="logs.maxDays" let:uniqueId>
                        <label for={uniqueId}>{$_("settings.app.form.maxdays")}</label>
                        <input type="number" id={uniqueId} required bind:value={formSettings.logs.maxDays} />
                    </Field>

                    <div class="col-lg-12 flex">
                        <div class="flex-fill" />
                        <button
                            type="submit"
                            class="btn btn-expanded"
                            class:btn-loading={isSaving}
                            disabled={!hasChanges || isSaving}
                            on:click={() => save()}
                        >
                            <span class="txt">{$_("settings.app.form.save")}</span>
                        </button>
                    </div>
                </div>
            {/if}
        </form>
    </div>
</main>
