<script>
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;

    let serialized = JSON.stringify(typeof value === "undefined" ? null : value, null, 2);

    $: if (value !== serialized?.trim()) {
        serialized = JSON.stringify(typeof value === "undefined" ? null : value, null, 2);
        value = serialized;
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <textarea
        id={uniqueId}
        class="txt-mono"
        required={field.required}
        value={serialized}
        on:input={(e) => {
            serialized = e.target.value;
            value = e.target.value.trim(); // trim the submitted value
        }}
    />
</Field>
