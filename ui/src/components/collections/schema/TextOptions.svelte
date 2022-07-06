<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let options = {};
</script>

<div class="grid">
    <div class="col-sm-6">
        <Field class="form-field" name="schema.{key}.options.min" let:uniqueId>
            <label for={uniqueId}>Min length</label>
            <input type="number" id={uniqueId} step="1" min="0" bind:value={options.min} />
        </Field>
    </div>

    <div class="col-sm-6">
        <Field class="form-field" name="schema.{key}.options.max" let:uniqueId>
            <label for={uniqueId}>Max length</label>
            <input type="number" id={uniqueId} step="1" min={options.min || 0} bind:value={options.max} />
        </Field>
    </div>

    <div class="col-sm-12">
        <Field class="form-field" name="schema.{key}.options.pattern" let:uniqueId>
            <label for={uniqueId}>Regex pattern</label>
            <input type="text" id={uniqueId} bind:value={options.pattern} />
            <div class="help-block">Valid Go regular expression, eg. <code>^\w+$</code>.</div>
        </Field>
    </div>
</div>
