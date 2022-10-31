<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Azure AD endpoints</div>
<div class="grid">
    <div class="col-lg-12">
        <Field class="form-field required" name="{key}.authUrl" let:uniqueId>
            <label for={uniqueId}>Auth URL</label>
            <input
                type="url"
                id={uniqueId}
                required
                placeholder="https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/authorize"
                bind:value={config.authUrl}
            />
            <div class="help-block">
                Eg. {`https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/authorize`}
            </div>
        </Field>
    </div>
    <div class="col-lg-12">
        <Field class="form-field required" name="{key}.tokenUrl" let:uniqueId>
            <label for={uniqueId}>Token URL</label>
            <input
                type="text"
                id={uniqueId}
                required
                placeholder="https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/token"
                bind:value={config.tokenUrl}
            />
            <div class="help-block">
                Eg. {`https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/token`}
            </div>
        </Field>
    </div>
</div>
