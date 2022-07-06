<script>
    import Flatpickr from "svelte-flatpickr";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let options = {};
</script>

<div class="grid">
    <div class="col-sm-6">
        <Field class="form-field" name="schema.{key}.options.min" let:uniqueId>
            <label for={uniqueId}>Min date (UTC)</label>
            <Flatpickr
                id={uniqueId}
                options={CommonHelper.defaultFlatpickrOptions()}
                value={options.min}
                bind:formattedValue={options.min}
            />
        </Field>
    </div>

    <div class="col-sm-6">
        <Field class="form-field" name="schema.{key}.options.max" let:uniqueId>
            <label for={uniqueId}>Max date (UTC)</label>
            <Flatpickr
                id={uniqueId}
                options={CommonHelper.defaultFlatpickrOptions()}
                value={options.max}
                bind:formattedValue={options.max}
            />
        </Field>
    </div>
</div>
