<script>
    import Field from "@/components/base/Field.svelte";

    export let field;
    export let value = false;
</script>

<Field class="form-field form-field-toggle {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={value} />
    <label for={uniqueId}>{field.name}</label>
</Field>
