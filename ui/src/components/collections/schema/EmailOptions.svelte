<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";

    export let key = "";
    export let options = {};
</script>

<div class="grid">
    <div class="col-sm-6">
        <Field class="form-field" name="schema.{key}.options.exceptDomains" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">Except domains</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: 'List of domains that are NOT allowed. \n This field is disabled if "Only domains" is set.',
                        position: "top",
                    }}
                />
            </label>
            <MultipleValueInput
                id={uniqueId}
                disabled={!CommonHelper.isEmpty(options.onlyDomains)}
                bind:value={options.exceptDomains}
            />
            <div class="help-block">Use comma as separator.</div>
        </Field>
    </div>

    <div class="col-sm-6">
        <Field class="form-field" name="schema.{key}.options.onlyDomains" let:uniqueId>
            <label for="{uniqueId}.options.onlyDomains">
                <span class="txt">Only domains</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: 'List of domains that are ONLY allowed. \n This field is disabled if "Except domains" is set.',
                        position: "top",
                    }}
                />
            </label>
            <MultipleValueInput
                id="{uniqueId}.options.onlyDomains"
                disabled={!CommonHelper.isEmpty(options.exceptDomains)}
                bind:value={options.onlyDomains}
            />
            <div class="help-block">Use comma as separator.</div>
        </Field>
    </div>
</div>
