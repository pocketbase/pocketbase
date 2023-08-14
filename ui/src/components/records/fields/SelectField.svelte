<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Select from "@/components/base/Select.svelte";
    import Field from "@/components/base/Field.svelte";

    export let field;
    export let value = undefined;

    $: isMultiple = field.options?.maxSelect > 1;

    $: if (typeof value === "undefined") {
        value = isMultiple ? [] : "";
    }

    $: if (isMultiple && Array.isArray(value) && value.length > field.options.maxSelect) {
        value = value.slice(value.length - field.options.maxSelect);
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <Select
        id={uniqueId}
        toggle={!field.required || isMultiple}
        multiple={isMultiple}
        closable={!isMultiple || value?.length >= field.options?.maxSelect}
        items={field.options?.values}
        searchable={field.options?.values?.length > 5}
        bind:selected={value}
    />
    {#if field.options?.maxSelect > 1}
        <div class="help-block">Select up to {field.options.maxSelect} items.</div>
    {/if}
</Field>
