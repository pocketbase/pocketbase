<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import AutoExpandTextarea from "@/components/base/AutoExpandTextarea.svelte";

    export let field;
    export let value = undefined;
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <AutoExpandTextarea id={uniqueId} required={field.required} bind:value />
</Field>
