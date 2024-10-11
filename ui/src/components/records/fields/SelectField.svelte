<script>
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let field;
    export let value = undefined;

    $: isMultiple = field.maxSelect > 1;

    $: if (typeof value === "undefined") {
        value = isMultiple ? [] : "";
    }

    $: maxSelect = field.maxSelect || field.values.length;

    $: if (isMultiple && Array.isArray(value)) {
        value = value.filter((v) => field.values.includes(v));
        if (value.length > maxSelect) {
            value = value.slice(value.length - maxSelect);
        }
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <FieldLabel {uniqueId} {field} />

    <Select
        id={uniqueId}
        toggle={!field.required || isMultiple}
        multiple={isMultiple}
        closable={!isMultiple || value?.length >= field.maxSelect}
        items={field.values}
        searchable={field.values?.length > 5}
        bind:selected={value}
    />
    {#if isMultiple}
        <div class="help-block">Select up to {maxSelect} items.</div>
    {/if}
</Field>
