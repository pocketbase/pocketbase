<script>
    import { SchemaField } from "pocketbase";
    import CodeMirror from "svelte-codemirror-editor"
    import { javascript } from "@codemirror/lang-javascript";
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
    <CodeMirror value={serialized} lang={javascript()} on:change={(e) => {
        serialized = e.detail
        value = e.detail.trim()
    }}/>
</Field>
