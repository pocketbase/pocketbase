<script>
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;

    let colorInput = undefined;

    function focusColorInput() {
        colorInput.click();
    }
</script>

<Field class="form-field form-field-color {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <div class="color-picker">
        <input
            type="text"
            id={uniqueId + "-input"}
            bind:value
            required={field.required}
        />

        <div class="color-picker-preview" on:click={focusColorInput} style="background-color: {value}" />

        <input
            type="color"
            id={uniqueId}
            required={field.required}
            bind:this={colorInput}
            bind:value
        />
    </div>
</Field>
