<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Authentik endpoints</div>
<div class="grid">
    <div class="col-lg-12">
        <Field class="form-field required" name="{key}.authUrl" let:uniqueId>
            <label for={uniqueId}>Auth URL</label>
            <input type="url" id={uniqueId} bind:value={config.authUrl} required />
            <div class="help-block">Eg. https://YOUR_AUTHENTIK_URL/application/o/authorize/</div>
        </Field>
    </div>
    <div class="col-lg-12">
        <Field class="form-field required" name="{key}.tokenUrl" let:uniqueId>
            <label for={uniqueId}>Token URL</label>
            <input type="text" id={uniqueId} bind:value={config.tokenUrl} required />
            <div class="help-block">Eg. https://YOUR_AUTHENTIK_URL/application/o/token/</div>
        </Field>
    </div>
    <div class="col-lg-12">
        <Field class="form-field required" name="{key}.userApiUrl" let:uniqueId>
            <label for={uniqueId}>User API URL</label>
            <input type="text" id={uniqueId} bind:value={config.userApiUrl} required />
            <div class="help-block">Eg. https://YOUR_AUTHENTIK_URL/application/o/userinfo/</div>
        </Field>
    </div>
</div>
