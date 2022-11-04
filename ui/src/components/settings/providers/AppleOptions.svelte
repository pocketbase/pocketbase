<script>
    import Field from "@/components/base/Field.svelte";
    import RedactedTextarea from "@/components/base/RedactedTextarea.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Client Secret Generation</div>
<div class="grid">
    <div class="col-lg-12">
        <Field class="form-field" name="{key}.teamId" let:uniqueId>
            <label for={uniqueId}>Team ID</label>
            <input type="text" id={uniqueId} placeholder="ABCDEF1234" bind:value={config.teamId} />
        </Field>
    </div>
    <div class="col-lg-12">
        <Field class="form-field" name="{key}.keyId" let:uniqueId>
            <label for={uniqueId}>Key ID</label>
            <input type="text" id={uniqueId} placeholder="ABCDEF1234" bind:value={config.keyId} />
        </Field>
    </div>
    <div class="col-lg-12">
        <Field class="form-field" name="{key}.signingKey" let:uniqueId>
            <label for={uniqueId}>Private Key</label>
            <RedactedTextarea id={uniqueId} bind:value={config.signingKey} />
            <div class="help-block">
                Should be wrapped in <code>-----BEGIN PRIVATE KEY-----</code> and
                <code>-----END PRIVATE KEY-----</code>
            </div>
        </Field>
    </div>
</div>
