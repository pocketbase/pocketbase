<script>
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import UserSelect from "@/components/users/UserSelect.svelte";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;

    // to prevent accidental changes, disable editing system user field values from the UI
    $: isDisabled = !CommonHelper.isEmpty(value) && field.system;

    $: isMultiple = field.options?.maxSelect > 1;

    $: if (isMultiple && Array.isArray(value) && value.length > field.options.maxSelect) {
        value = value.slice(field.options.maxSelect - 1);
    }
</script>

<Field
    class="form-field {field.required ? 'required' : ''} {isDisabled ? 'disabled' : ''}"
    name={field.name}
    let:uniqueId
>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <UserSelect toggle id={uniqueId} multiple={isMultiple} disabled={isDisabled} bind:keyOfSelected={value} />
    {#if field.options?.maxSelect > 1}
        <div class="help-block">Select up to {field.options.maxSelect} users.</div>
    {/if}
</Field>
