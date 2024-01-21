<script>
    import Flatpickr from "svelte-flatpickr";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    let pickerMinValue = field?.options?.min;
    let pickerMaxValue = field?.options?.max;

    $: if (pickerMinValue != field?.options?.min) {
        pickerMinValue = field?.options?.min;
    }

    $: if (pickerMaxValue != field?.options?.max) {
        pickerMaxValue = field?.options?.max;
    }

    // ensure that value is set even on manual input edit
    function onClose(e, key) {
        if (e.detail && e.detail.length == 3) {
            field.options[key] = e.detail[1];
        }
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.min" let:uniqueId>
                    <label for={uniqueId}>Min date (UTC)</label>
                    <Flatpickr
                        id={uniqueId}
                        options={CommonHelper.defaultFlatpickrOptions()}
                        bind:value={pickerMinValue}
                        bind:formattedValue={field.options.min}
                        on:close={(e) => onClose(e, "min")}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.max" let:uniqueId>
                    <label for={uniqueId}>Max date (UTC)</label>
                    <Flatpickr
                        id={uniqueId}
                        options={CommonHelper.defaultFlatpickrOptions()}
                        bind:value={pickerMaxValue}
                        bind:formattedValue={field.options.max}
                        on:close={(e) => onClose(e, "max")}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
