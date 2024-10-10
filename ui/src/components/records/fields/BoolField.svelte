<script>
    import Field from "@/components/base/Field.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let field;
    export let value = false;
</script>

<Field class="form-field form-field-toggle {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={value} />
    <FieldLabel {uniqueId} {field} icon={false} />
</Field>
