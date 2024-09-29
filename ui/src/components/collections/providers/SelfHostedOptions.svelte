<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
    export let required = false;
    export let title = "Provider endpoints";

    $: isRequired = required && config?.enabled;
</script>

<div class="section-title">{title}</div>
<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.authURL" let:uniqueId>
    <label for={uniqueId}>Auth URL</label>
    <input type="url" id={uniqueId} bind:value={config.authURL} required={isRequired} />
</Field>
<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.tokenURL" let:uniqueId>
    <label for={uniqueId}>Token URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenURL} required={isRequired} />
</Field>
<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.userInfoURL" let:uniqueId>
    <label for={uniqueId}>User info URL</label>
    <input type="url" id={uniqueId} bind:value={config.userInfoURL} required={isRequired} />
</Field>
