<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
    export let required = false;
    export let title = "Provider endpoints";

    $: isRequired = required && config?.enabled;
</script>

<div class="section-title">{title}</div>
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
