<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Azure AD endpoints</div>
<Field class="form-field {config.enabled ? 'required' : ''}" name="{key}.authUrl" let:uniqueId>
    <label for={uniqueId}>Auth URL</label>
    <input type="url" id={uniqueId} bind:value={config.authUrl} required={config.enabled} />
    <div class="help-block">
        Eg. {`https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/authorize`}
    </div>
</Field>
<Field class="form-field {config.enabled ? 'required' : ''}" name="{key}.tokenUrl" let:uniqueId>
    <label for={uniqueId}>Token URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenUrl} required={config.enabled} />
    <div class="help-block">
        Eg. {`https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/token`}
    </div>
</Field>
