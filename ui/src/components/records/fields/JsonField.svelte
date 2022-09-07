<script>
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;

    $: if (typeof value !== "undefined" && typeof value !== "string" && value !== null) {
        // the JSON field support both js primitives and encoded JSON string
        // so we are normalizing the value to only a string
        value = JSON.stringify(value, null, 2);
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <textarea id={uniqueId} required={field.required} class="txt-mono" bind:value />
</Field>
