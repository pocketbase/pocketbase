<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import AutoExpandTextarea from "@/components/base/AutoExpandTextarea.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let original;
    export let field;
    export let value = undefined;

    $: hasAutogenerate = !CommonHelper.isEmpty(field.autogeneratePattern) && !original?.id;

    $: isRequired = field.required && !hasAutogenerate;
</script>

<Field class="form-field {isRequired ? 'required' : ''}" name={field.name} let:uniqueId>
    <FieldLabel {uniqueId} {field} />

    <AutoExpandTextarea
        id={uniqueId}
        required={isRequired}
        placeholder={hasAutogenerate ? "Leave empty to autogenerate..." : ""}
        bind:value
    />
</Field>
