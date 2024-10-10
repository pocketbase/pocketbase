<script>
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import Flatpickr from "svelte-flatpickr";

    export let field;
    export let key = "";

    let pickerMinValue = field?.min;
    let pickerMaxValue = field?.max;

    $: if (pickerMinValue != field?.min) {
        pickerMinValue = field?.min;
    }

    $: if (pickerMaxValue != field?.max) {
        pickerMaxValue = field?.max;
    }

    // ensure that value is set even on manual input edit
    function onClose(e, key) {
        if (e.detail && e.detail.length == 3) {
            field[key] = e.detail[1];
        }
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.min" let:uniqueId>
                    <label for={uniqueId}>Min date (UTC)</label>
                    <Flatpickr
                        id={uniqueId}
                        options={CommonHelper.defaultFlatpickrOptions()}
                        bind:value={pickerMinValue}
                        bind:formattedValue={field.min}
                        on:close={(e) => onClose(e, "min")}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.max" let:uniqueId>
                    <label for={uniqueId}>Max date (UTC)</label>
                    <Flatpickr
                        id={uniqueId}
                        options={CommonHelper.defaultFlatpickrOptions()}
                        bind:value={pickerMaxValue}
                        bind:formattedValue={field.max}
                        on:close={(e) => onClose(e, "max")}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
