<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    }

    function loadDefaults() {
        field.options = {
            convertUrls: false,
        };
    }
</script>

<SchemaField bind:field {key} on:rename on:remove {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-12">
                <Field class="form-field" name="schema.{key}.title" let:title>
                    <label for={title}>Display Name</label>
                    <input type="text" id={title} bind:value={field.title} />
                </Field>
            </div>
        </div>
    </svelte:fragment>

    <svelte:fragment slot="optionsFooter">
        <Field class="form-field form-field-toggle" name="schema.{key}.options.convertUrls" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={field.options.convertUrls} />
            <label for={uniqueId}>
                <span class="txt">Strip urls domain</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `This could help making the editor content more portable between environments since there will be no local base url to replace.`,
                    }}
                />
            </label>
        </Field>
    </svelte:fragment>
</SchemaField>
