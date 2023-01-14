<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Authentik endpoints</div>
<div class="grid">
    <div class="col-lg-4">
        <Field class="form-field" name="{key}.authUrl" let:uniqueId>
            <label for={uniqueId}>Authorize URL</label>
            <input type="url" id={uniqueId} bind:value={config.authUrl} />
        </Field>
    </div>
    <div class="col-lg-4">
        <Field class="form-field" name="{key}.tokenUrl" let:uniqueId>
            <label for={uniqueId}>Token URL</label>
            <input type="text" id={uniqueId} bind:value={config.tokenUrl} />
        </Field>
    </div>
    <div class="col-lg-4">
        <Field class="form-field" name="{key}.userApiUrl" let:uniqueId>
            <label for={uniqueId}>Userinfo URL</label>
            <input type="text" id={uniqueId} bind:value={config.userApiUrl} />
        </Field>
    </div>
</div>
