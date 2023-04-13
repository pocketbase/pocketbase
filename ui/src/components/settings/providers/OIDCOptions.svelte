<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Endpoints</div>
<Field class="form-field {config.enabled ? 'required' : ''}" name="{key}.authUrl" let:uniqueId>
    <label for={uniqueId}>Auth URL</label>
    <input type="url" id={uniqueId} bind:value={config.authUrl} required={config.enabled} />
    <div class="help-block">Eg. https://example.com/authorize/</div>
</Field>
<Field class="form-field {config.enabled ? 'required' : ''}" name="{key}.tokenUrl" let:uniqueId>
    <label for={uniqueId}>Token URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenUrl} required={config.enabled} />
    <div class="help-block">Eg. https://example.com/token/</div>
</Field>
<Field class="form-field {config.enabled ? 'required' : ''}" name="{key}.userApiUrl" let:uniqueId>
    <label for={uniqueId}>User API URL</label>
    <input type="url" id={uniqueId} bind:value={config.userApiUrl} required={config.enabled} />
    <div class="help-block">Eg. https://example.com/userinfo/</div>
</Field>
