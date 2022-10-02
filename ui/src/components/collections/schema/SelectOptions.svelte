<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";

    export let key = "";
    export let options = {};

    $: if (CommonHelper.isEmpty(options)) {
        // load defaults
        options = {
            maxSelect: 1,
            values: [],
        };
    }

    // note: leave the validation to the api
    // $: if (!CommonHelper.isEmpty(options.values) && options.maxSelect > options.values.length) {
    //     options.maxSelect = options.values.length;
    // }
</script>

<div class="grid">
    <div class="col-sm-9">
        <Field class="form-field required" name="schema.{key}.options.values" let:uniqueId>
            <label for={uniqueId}>Choices</label>
            <MultipleValueInput
                id={uniqueId}
                placeholder="eg. optionA, optionB"
                required
                bind:value={options.values}
            />
            <div class="help-block">Use comma as separator.</div>
        </Field>
    </div>

    <div class="col-sm-3">
        <Field class="form-field required" name="schema.{key}.options.maxSelect" let:uniqueId>
            <label for={uniqueId}>Max select</label>
            <input type="number" id={uniqueId} step="1" min="1" required bind:value={options.maxSelect} />
        </Field>
    </div>
</div>
