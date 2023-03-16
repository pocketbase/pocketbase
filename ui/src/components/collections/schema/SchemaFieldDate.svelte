<script>
    import Flatpickr from "svelte-flatpickr";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;

    export let key = "";
</script>

<SchemaField
    bind:field
    {key}
    on:rename
    on:remove
    on:drop
    on:dragstart
    on:dragenter
    on:dragleave
    {...$$restProps}
>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.min" let:uniqueId>
                    <label for={uniqueId}>Min date (UTC)</label>
                    <Flatpickr
                        id={uniqueId}
                        options={CommonHelper.defaultFlatpickrOptions()}
                        value={field.options.min}
                        bind:formattedValue={field.options.min}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.max" let:uniqueId>
                    <label for={uniqueId}>Max date (UTC)</label>
                    <Flatpickr
                        id={uniqueId}
                        options={CommonHelper.defaultFlatpickrOptions()}
                        value={field.options.max}
                        bind:formattedValue={field.options.max}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
