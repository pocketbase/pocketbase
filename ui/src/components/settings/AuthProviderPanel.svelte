<script>
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    const formId = "provider_popup_" + CommonHelper.randomString(5);

    let panel;
    let provider = {};
    let config = {};
    let isSubmitting = false;
    let initialHash = "";

    $: hasChanges = JSON.stringify(config) != initialHash;

    export function show(showProvider, showConfig) {
        setErrors({}); // reset any previous errors

        provider = Object.assign({}, showProvider);
        config = Object.assign({ enabled: true }, showConfig);
        initialHash = JSON.stringify(config);

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    async function submit() {
        isSubmitting = true;

        try {
            const data = {};
            data[provider.key] = CommonHelper.filterRedactedProps(config);

            const result = await ApiClient.settings.update(data);

            setErrors({});

            addSuccessToast("Successfully updated provider settings.");

            dispatch("submit", result);

            hide();
        } catch (err) {
            ApiClient.error(err);
        }

        isSubmitting = false;
    }

    function clear() {
        for (let k in config) {
            config[k] = CommonHelper.zeroValue(config[k]);
        }

        // set to false only for the oidc providers
        // (@todo remove after the refactoring)
        if (provider.key?.startsWith("oidc")) {
            config.pkce = false;
        } else {
            config.pkce = null;
        }
    }
</script>

<OverlayPanel bind:this={panel} overlayClose={!isSubmitting} escClose={!isSubmitting} on:show on:hide>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">{provider.title || provider.key} provider</h4>
    </svelte:fragment>

    <form id={formId} autocomplete="off" on:submit|preventDefault={() => submit()}>
        <div class="flex m-b-base">
            <Field class="form-field form-field-toggle m-b-0" name="{provider.key}.enabled" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={config.enabled} />
                <label for={uniqueId}>Enable</label>
            </Field>

            <button type="button" class="btn btn-sm btn-transparent btn-hint m-l-auto" on:click={clear}>
                <span class="txt">Clear all fields</span>
            </button>
        </div>

        <Field
            class="form-field {config.enabled ? 'required' : ''}"
            name="{provider.key}.clientId"
            let:uniqueId
        >
            <label for={uniqueId}>Client ID</label>
            <input type="text" id={uniqueId} bind:value={config.clientId} required={config.enabled} />
        </Field>

        <Field
            class="form-field {config.enabled ? 'required' : ''}"
            name="{provider.key}.clientSecret"
            let:uniqueId
        >
            <label for={uniqueId}>Client secret</label>
            <RedactedPasswordInput bind:value={config.clientSecret} id={uniqueId} required={config.enabled} />
        </Field>

        {#if provider.optionsComponent}
            <div class="col-lg-12">
                <svelte:component
                    this={provider.optionsComponent}
                    key={provider.key}
                    bind:config
                    {...provider.optionsComponentProps || {}}
                />
            </div>
        {/if}
    </form>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide} disabled={isSubmitting}>
            Close
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSubmitting}
            disabled={!hasChanges || isSubmitting}
        >
            <span class="txt">Save changes</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
