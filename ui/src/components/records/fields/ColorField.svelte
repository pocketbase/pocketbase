<script>
    import Field from "@/components/base/Field.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let field;
    export let value = undefined;

    $: isRequired = field.required;
</script>

<Field class="form-field {isRequired ? 'required' : ''}" name={field.name} let:uniqueId>

    <div class="color-input-wrapper">
        <FieldLabel {uniqueId} {field} />
        <label for="{uniqueId}-color-picker" class="color-swatch-wrapper" style="border-radius: 0;">
            <span class="color-swatch" style="background-color: {value ? value : '#000000'}" title={value}></span>
            <input
                type="color"
                id="{uniqueId}-color-picker"
                bind:value
                class="color-picker-hidden"
            />
        </label>
        <input
            type="text"
            id={uniqueId}
            required={isRequired}
            bind:value
            placeholder="#000000"
        />
    </div>

</Field>

<style>
    .color-swatch {
        display: inline-block;
        border-radius: var(--baseRadius);
        width: 100% ;
        height: 20px;
        vertical-align: middle;
    }
    .color-picker-hidden {
        display: none;
    }
</style>
