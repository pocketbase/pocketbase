<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let key = "";
    export let config = {};

    $: isRequired = !!config.enabled;

    if (CommonHelper.isEmpty(config.pkce)) {
        config.pkce = true;
    }

    if (!config.displayName) {
        config.displayName = "OIDC";
    }
</script>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.displayName" let:uniqueId>
    <label for={uniqueId}>Display name</label>
    <input type="text" id={uniqueId} bind:value={config.displayName} required={isRequired} />
</Field>

<div class="section-title">Endpoints</div>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.authUrl" let:uniqueId>
    <label for={uniqueId}>Auth URL</label>
    <input type="url" id={uniqueId} bind:value={config.authUrl} required={isRequired} />
</Field>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.tokenUrl" let:uniqueId>
    <label for={uniqueId}>Token URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenUrl} required={isRequired} />
</Field>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.userApiUrl" let:uniqueId>
    <label for={uniqueId}>User API URL</label>
    <input type="url" id={uniqueId} bind:value={config.userApiUrl} required={isRequired} />
</Field>

<Field class="form-field" name="{key}.pkce" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={config.pkce} />
    <label for={uniqueId}>
        <span class="txt">Support PKCE</span>
        <i
            class="ri-information-line link-hint"
            use:tooltip={{
                text: "Usually it should be safe to be always enabled as most providers will just ignore the extra query parameters if they don't support PKCE.",
                position: "right",
            }}
        />
    </label>
</Field>
