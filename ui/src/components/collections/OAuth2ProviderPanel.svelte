<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";
    import { confirm } from "@/stores/confirmation";
    import { removeError } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    const formId = "provider_popup_" + CommonHelper.randomString(5);

    let panel;
    let uiOptions = {};
    let config = {};
    let isNew = false;
    let initialHash = "";
    let maskSecret = false;
    let providerIndex = 0;

    $: hasChanges = JSON.stringify(config) != initialHash;

    $: errPrefix = "oauth2.providers." + providerIndex;

    export function show(showOptions, showConfig, showIndex) {
        providerIndex = showIndex || 0;
        isNew = CommonHelper.isEmpty(showConfig);
        uiOptions = Object.assign({}, showOptions);
        config = Object.assign({}, showConfig);
        maskSecret = !!config.clientId;
        initialHash = JSON.stringify(config);

        panel?.show();
    }

    export function hide() {
        removeError(errPrefix);
        panel?.hide();
    }

    async function submit() {
        dispatch("submit", { uiOptions, config });
        hide();
    }

    async function remove() {
        confirm(
            `Do you really want to remove the "${uiOptions.title}" OAuth2 provider from the collection?`,
            () => {
                dispatch("remove", { uiOptions });
                hide();
            },
        );
    }
</script>

<OverlayPanel bind:this={panel} btnClose={false} on:show on:hide>
    <svelte:fragment slot="header">
        <figure class="provider-logo">
            {#if uiOptions.logoBase64}
                <img
                    src="data:image/svg+xml;base64,{uiOptions.logoBase64}"
                    alt="{uiOptions.title} logo"
                />
            {:else}
                <i class="ri-puzzle-line txt-sm txt-hint"></i>
            {/if}
        </figure>
        <h4 class="center txt-break">{uiOptions.title} <small class="txt-hint">({uiOptions.key})</small></h4>
    </svelte:fragment>

    <form id={formId} autocomplete="off" on:submit|preventDefault={() => submit()}>
        <Field class="form-field required" name="{errPrefix}.clientId" let:uniqueId>
            <label for={uniqueId}>Client ID</label>
            <input type="text" id={uniqueId} bind:value={config.clientId} />
        </Field>

        <Field class="form-field required" name="{errPrefix}.clientSecret" let:uniqueId>
            <label for={uniqueId}>Client secret</label>
            <RedactedPasswordInput id={uniqueId} bind:mask={maskSecret} bind:value={config.clientSecret} />
        </Field>

        {#if uiOptions.optionsComponent}
            <div class="col-lg-12">
                <svelte:component
                    this={uiOptions.optionsComponent}
                    key={errPrefix}
                    bind:config
                    {...uiOptions.optionsComponentProps || {}}
                />
            </div>
        {/if}
    </form>

    <svelte:fragment slot="footer">
        {#if !isNew}
            <button
                type="button"
                class="btn btn-transparent btn-circle btn-hint btn-sm"
                aria-label="Remove provider"
                use:tooltip={{ text: "Remove provider", position: "right" }}
                on:click={remove}
            >
                <i class="ri-delete-bin-7-line" aria-hidden="true" />
            </button>
            <div class="flex-fill"></div>
        {/if}
        <button type="button" class="btn btn-transparent" on:click={hide}>Cancel</button>
        <button type="submit" form={formId} class="btn btn-expanded" disabled={!hasChanges}>
            <span class="txt">Set provider config</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
