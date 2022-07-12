<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { _, locale } from '@/services/i18n';

    let tokensList = [];

    let tokenSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";

    $: hasChanges = initialHash != JSON.stringify(tokenSettings);

    CommonHelper.setDocumentTitle($_("settings.token.pagetitle"));

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const result = (await ApiClient.Settings.getAll()) || {};
            updateTokenList()
            initSettings(result);
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }

    locale.subscribe(() => updateTokenList())

    function updateTokenList(){
        tokensList =  [
                { key: "userAuthToken", label: $_("settings.token.label.userAuthToken") },
                { key: "userVerificationToken", label: $_("settings.token.label.userVerificationToken") },
                { key: "userPasswordResetToken", label: $_("settings.token.label.userPasswordResetToken") },
                { key: "userEmailChangeToken", label: $_("settings.token.label.userEmailChangeToken") },
                { key: "adminAuthToken", label: $_("settings.token.label.adminAuthToken") },
                { key: "adminPasswordResetToken", label: $_("settings.token.label.adminPasswordResetToken") },
            ];
    }

    async function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        try {
            const result = await ApiClient.Settings.update(CommonHelper.filterRedactedProps(tokenSettings));
            initSettings(result);
            addSuccessToast($_("settings.token.tips.saved"));
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isSaving = false;
    }

    function initSettings(data) {
        data = data || {};
        tokenSettings = {};

        for (const listItem of tokensList) {
            tokenSettings[listItem.key] = {
                duration: data[listItem.key]?.duration || 0,
            };
        }

        initialHash = JSON.stringify(tokenSettings);
    }
</script>

<SettingsSidebar />

<main class="page-wrapper">
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">{$_("app.breadcrumb.settings")}</div>
            <div class="breadcrumb-item">{$_("app.breadcrumb.token")}</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
            <div class="content m-b-sm txt-xl">
                <p>{$_("settings.token.title")}</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                {#each tokensList as token (token.key)}
                    <Field class="form-field required" name="{token.key}.duration" let:uniqueId>
                        <label for={uniqueId}>{token.label} {$_("settings.token.tips.duration")}</label>
                        <input
                            type="number"
                            id={uniqueId}
                            required
                            bind:value={tokenSettings[token.key].duration}
                        />
                        <div class="help-block">
                            <span
                                class="link-primary"
                                class:txt-success={tokenSettings[token.key].secret}
                                on:click={() => {
                                    // toggle
                                    if (tokenSettings[token.key].secret) {
                                        delete tokenSettings[token.key].secret;
                                        tokenSettings[token.key] = tokenSettings[token.key];
                                    } else {
                                        tokenSettings[token.key].secret = CommonHelper.randomString(50);
                                    }
                                }}
                            >
                            {$_("settings.token.tips.invalidate_all")}
                            </span>
                        </div>
                    </Field>
                {/each}

                <div class="flex">
                    <div class="flex-fill" />
                    <button
                        type="submit"
                        class="btn btn-expanded"
                        class:btn-loading={isSaving}
                        disabled={!hasChanges || isSaving}
                        on:click={() => save()}
                    >
                        <span class="txt">{$_("settings.token.tips.save")}</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</main>
