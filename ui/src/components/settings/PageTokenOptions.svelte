<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    const tokensList = [
        { key: "userAuthToken", label: "Users auth token" },
        { key: "userVerificationToken", label: "Users email verification token" },
        { key: "userPasswordResetToken", label: "Users password reset token" },
        { key: "userEmailChangeToken", label: "Users email change token" },
        { key: "adminAuthToken", label: "Admins auth token" },
        { key: "adminPasswordResetToken", label: "Admins password reset token" },
    ];

    $pageTitle = "Token options";

    let tokenSettings = {};
    let isLoading = false;
    let isSaving = false;
    let initialHash = "";

    $: hasChanges = initialHash != JSON.stringify(tokenSettings);

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
            const result = await ApiClient.settings.update(CommonHelper.filterRedactedProps(tokenSettings));
            initSettings(result);
            addSuccessToast("Successfully saved tokens options.");
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
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>

    <div class="wrapper">
        <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
            <div class="content m-b-sm txt-xl">
                <p>Adjust common token options.</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                {#each tokensList as token (token.key)}
                    <Field class="form-field required" name="{token.key}.duration" let:uniqueId>
                        <label for={uniqueId}>{token.label} duration (in seconds)</label>
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
                                Invalidate all previously issued tokens
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
                        <span class="txt">Save changes</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</main>
