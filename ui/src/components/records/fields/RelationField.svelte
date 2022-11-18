<script>
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import RecordSelect from "@/components/records/RecordSelect.svelte";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;

    $: isMultiple = field.options?.maxSelect != 1;

    $: if (
        isMultiple &&
        Array.isArray(value) &&
        field.options?.maxSelect &&
        value.length > field.options.maxSelect
    ) {
        value = value.slice(field.options.maxSelect - 1);
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <RecordSelect
        toggle
        id={uniqueId}
        multiple={isMultiple}
        collectionId={field.options?.collectionId}
        bind:keyOfSelected={value}
    />
    {#if field.options?.maxSelect > 1}
        <div class="help-block">Select up to {field.options.maxSelect} items.</div>
    {/if}
</Field>
