<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    $pageTitle = "Database";

    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;
    let showToken = false; // State to toggle token visibility

    $: initialHash = JSON.stringify(originalFormSettings);
    $: hasChanges = initialHash != JSON.stringify(formSettings);

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.settings.getAll()) || {};
            init(settings);
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
            setErrors({});
            addSuccessToast("Successfully saved database settings.");
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function init(settings = {}) {
        formSettings = {
            libsql: settings?.libsql || {
                url: "libsql://pocketbase-default-db.turso.io",
                token: "testtokenkeylookalike",
            }, // Initialize LibSQL settings with default values
        };

        originalFormSettings = JSON.parse(JSON.stringify(formSettings));
    }

    function reset() {
        formSettings = JSON.parse(JSON.stringify(originalFormSettings || {}));
    }

    function toggleTokenVisibility() {
        showToken = !showToken; // Toggle the visibility state
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
        <form class="panel" autocomplete="off" on:submit|preventDefault={() => save()}>
            <div class="content txt-xl m-b-base">
                <p>Configure your LibSQL database connection settings.</p>
                <p>
                    By default we use Turso as a database origin but you can use any that can support libsql
                    database.
                </p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <div class="grid m-b-base">
                    <div class="col-lg-6">
                        <Field class="form-field required" name="libsql.url" let:uniqueId>
                            <label for={uniqueId}>Database Origin</label>
                            <input type="text" id={uniqueId} required bind:value={formSettings.libsql.url} />
                        </Field>
                    </div>

                    <div class="col-lg-6">
                        <Field class="form-field required" name="libsql.token" let:uniqueId>
                            <label for={uniqueId}>Token</label>
                            <div class="input-group">
                                {#if showToken}
                                    <input
                                        type="text"
                                        id={uniqueId}
                                        required
                                        bind:value={formSettings.libsql.token}
                                    />
                                {:else}
                                    <input
                                        type="password"
                                        id={uniqueId}
                                        required
                                        bind:value={formSettings.libsql.token}
                                    />
                                {/if}

                                <span class="input-group-addon flex gap_custome_1">
                                    {#if showToken}
                                        <i class="ri-checkbox-line" on:click={toggleTokenVisibility}></i>
                                    {:else}
                                        <i class="ri-checkbox-blank-line" on:click={toggleTokenVisibility}
                                        ></i>
                                    {/if}
                                    <p>show token</p>
                                </span>
                            </div>
                        </Field>
                    </div>
                </div>

                <div class="flex">
                    <div class="flex-fill" />

                    <button
                        type="button"
                        class="btn btn-transparent btn-hint"
                        disabled={isSaving}
                        on:click={() => reset()}
                    >
                        <span class="txt">Cancel</span>
                    </button>
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

<style>
    .gap_custome_1 {
        gap: 3px !important;
        /*     margin-left: 15px;*/
    }
</style>
