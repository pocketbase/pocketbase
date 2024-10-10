<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
    import OAuth2ProviderPanel from "@/components/collections/OAuth2ProviderPanel.svelte";
    import OAuth2ProvidersListPanel from "@/components/collections/OAuth2ProvidersListPanel.svelte";
    import providersUIList from "@/providers.js";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale, slide } from "svelte/transition";

    export let collection;

    const excludedFieldNames = ["id", "email", "emailVisibility", "verified", "tokenKey", "password"];
    const allowedRegularTypes = ["text", "editor", "url", "email", "json"];
    const allowedRegularAndFileTypes = allowedRegularTypes.concat("file");

    let providersListPanel;
    let providerPanel;
    let showMappedFields = false;
    let regularFieldOptions = [];
    let regularAndFileFieldOptions = [];

    $: refreshFieldOptions(collection.fields);

    $: if (CommonHelper.isEmpty(collection.oauth2)) {
        collection.oauth2 = {
            enabled: false,
            mappedFields: {},
            providers: [],
        };
    }

    $: hasErrors = !CommonHelper.isEmpty($errors?.oauth2);

    $: totalProviders = collection.oauth2?.providers?.length || 0;

    function refreshFieldOptions(fields = []) {
        regularFieldOptions =
            fields
                ?.filter((f) => allowedRegularTypes.includes(f.type) && !excludedFieldNames.includes(f.name))
                ?.map((f) => f.name) || [];

        regularAndFileFieldOptions =
            fields
                ?.filter(
                    (f) =>
                        allowedRegularAndFileTypes.includes(f.type) && !excludedFieldNames.includes(f.name),
                )
                ?.map((f) => f.name) || [];
    }

    function getProviderUIOptions(key) {
        for (let item of providersUIList) {
            if (item.key == key) {
                return item;
            }
        }
        return null;
    }
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-pass-expired-line"></i>
            <span class="txt">OAuth2</span>
        </div>

        <div class="flex-fill" />

        {#if collection.oauth2.enabled}
            <span class="label" class:label-warning={!totalProviders} class:label-info={totalProviders > 0}>
                {totalProviders}
                {totalProviders == 1 ? "provider" : "providers"}
            </span>

            <span class="label label-success">Enabled</span>
        {:else}
            <span class="label">Disabled</span>
        {/if}

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "Has errors", position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle" name="oauth2.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={collection.oauth2.enabled} />
        <label for={uniqueId}>Enable</label>
    </Field>

    <div class="grid grid-sm">
        {#each collection.oauth2.providers as providerConfig, i (providerConfig.name)}
            {@const uiOptions = getProviderUIOptions(providerConfig.name)}

            <div class="col-lg-6">
                <div
                    class="provider-card"
                    class:error={!CommonHelper.isEmpty($errors?.oauth2?.providers?.[i])}
                >
                    <figure class="provider-logo">
                        {#if uiOptions?.logo}
                            <img
                                src="{import.meta.env.BASE_URL}images/oauth2/{uiOptions.logo}"
                                alt="{uiOptions.title} logo"
                            />
                        {:else}
                            <i class="ri-puzzle-line txt-sm txt-hint"></i>
                        {/if}
                    </figure>
                    <div class="content">
                        <div class="title">{providerConfig.displayName || uiOptions?.title || "Custom"}</div>
                        <em class="txt-hint txt-sm m-r-auto">{providerConfig.name}</em>
                    </div>
                    {#if uiOptions}
                        <button
                            type="button"
                            class="btn btn-circle btn-hint btn-transparent"
                            aria-label="Provider settings"
                            use:tooltip={{ text: "Edit config", position: "left" }}
                            on:click={() => {
                                providerPanel?.show(uiOptions, providerConfig, i);
                            }}
                        >
                            <i class="ri-settings-4-line" />
                        </button>
                    {/if}
                </div>
            </div>
        {/each}
        <div class="col-lg-6">
            <button
                class="btn btn-block btn-lg btn-secondary txt-base"
                on:click={() => providersListPanel?.show()}
            >
                <i class="ri-add-line"></i>
                <span class="txt">Add provider</span>
            </button>
        </div>
    </div>

    <button
        type="button"
        class="m-t-25 btn btn-sm {showMappedFields ? 'btn-secondary' : 'btn-hint btn-transparent'}"
        on:click={() => (showMappedFields = !showMappedFields)}
    >
        <strong class="txt">Optional {collection.name} create fields map</strong>
        {#if showMappedFields}
            <i class="ri-arrow-up-s-line txt-sm" />
        {:else}
            <i class="ri-arrow-down-s-line txt-sm" />
        {/if}
    </button>
    {#if showMappedFields}
        <div class="block" transition:slide={{ duration: 150 }}>
            <div class="grid grid-sm p-t-xs">
                <div class="col-sm-6">
                    <Field class="form-field form-field-toggle" name="oauth2.mappedFields.name" let:uniqueId>
                        <label for={uniqueId}>OAuth2 full name</label>
                        <Select
                            id={uniqueId}
                            items={regularFieldOptions}
                            toggle={true}
                            zeroFunc={() => ""}
                            selectPlaceholder={"Select field"}
                            bind:selected={collection.oauth2.mappedFields.name}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field
                        class="form-field form-field-toggle"
                        name="oauth2.mappedFields.avatarURL"
                        let:uniqueId
                    >
                        <label for={uniqueId}>OAuth2 avatar</label>
                        <Select
                            id={uniqueId}
                            items={regularAndFileFieldOptions}
                            toggle={true}
                            zeroFunc={() => ""}
                            selectPlaceholder={"Select field"}
                            bind:selected={collection.oauth2.mappedFields.avatarURL}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field class="form-field form-field-toggle" name="oauth2.mappedFields.id" let:uniqueId>
                        <label for={uniqueId}>OAuth2 id</label>
                        <Select
                            id={uniqueId}
                            items={regularFieldOptions}
                            toggle={true}
                            zeroFunc={() => ""}
                            selectPlaceholder={"Select field"}
                            bind:selected={collection.oauth2.mappedFields.id}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field
                        class="form-field form-field-toggle"
                        name="oauth2.mappedFields.username"
                        let:uniqueId
                    >
                        <label for={uniqueId}>OAuth2 username</label>
                        <Select
                            id={uniqueId}
                            items={regularFieldOptions}
                            toggle={true}
                            zeroFunc={() => ""}
                            selectPlaceholder={"Select field"}
                            bind:selected={collection.oauth2.mappedFields.username}
                        />
                    </Field>
                </div>
            </div>
        </div>
    {/if}
</Accordion>

<OAuth2ProvidersListPanel
    bind:this={providersListPanel}
    disabled={collection.oauth2?.providers?.map((p) => p.name) || []}
    on:select={(e) => {
        providerPanel.show(e.detail, {}, collection.oauth2?.providers?.length || 0);
    }}
/>

<OAuth2ProviderPanel
    bind:this={providerPanel}
    on:remove={(e) => {
        const uiOptions = e.detail.uiOptions;
        CommonHelper.removeByKey(collection.oauth2.providers, "name", uiOptions.key);
        collection.oauth2.providers = collection.oauth2.providers;
    }}
    on:submit={(e) => {
        const uiOptions = e.detail.uiOptions;
        const config = e.detail.config;
        collection.oauth2.providers = collection.oauth2.providers || [];
        CommonHelper.pushOrReplaceByKey(
            collection.oauth2.providers,
            Object.assign({ name: uiOptions.key }, config),
            "name",
        );
        collection.oauth2.providers = collection.oauth2.providers;
    }}
/>
