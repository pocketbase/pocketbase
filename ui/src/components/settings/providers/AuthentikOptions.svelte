<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">OpenID Configuration Issuer</div>
<div class="grid spacing">
    <div class="col-lg-8">
        <Field class="form-field" name="{key}.endpoint" let:uniqueId>
            <label for={uniqueId}>OpenID Configuration Issuer URL</label>
            <input type="url" id={uniqueId} bind:value={config.endpoint} />
        </Field>
    </div>
</div>
<div class="section-title">Configuration URL (optional)</div>
<div class="grid spacing">
    <div class="col-lg-8">
        <Field class="form-field" name="{key}.configUrl" let:uniqueId>
            <label for={uniqueId}>OpenID configuration URL</label>
            <input type="url" id={uniqueId} bind:value={config.configUrl} />
        </Field>
    </div>
</div>
<div class="section-title">Endpoint configuration (optional)</div>
<div class="grid spacing">
    <div class="col-lg-4">
        <Field class="form-field" name="{key}.authUrl" let:uniqueId>
            <label for={uniqueId}>Authorize URL</label>
            <input type="url" id={uniqueId} bind:value={config.authUrl} />
        </Field>
    </div>
    <div class="col-lg-4">
        <Field class="form-field" name="{key}.tokenUrl" let:uniqueId>
            <label for={uniqueId}>Token URL</label>
            <input type="url" id={uniqueId} bind:value={config.tokenUrl} />
        </Field>
    </div>
    <div class="col-lg-4">
        <Field class="form-field" name="{key}.userApiUrl" let:uniqueId>
            <label for={uniqueId}>Userinfo URL</label>
            <input type="url" id={uniqueId} bind:value={config.userApiUrl} />
        </Field>
    </div>
</div>
