<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.exceptDomains" let:uniqueId>
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
                        disabled={!CommonHelper.isEmpty(field.onlyDomains)}
                        bind:value={field.exceptDomains}
                    />
                    <div class="help-block">Use comma as separator.</div>
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.onlyDomains" let:uniqueId>
                    <label for="{uniqueId}.onlyDomains">
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
                        id="{uniqueId}.onlyDomains"
                        disabled={!CommonHelper.isEmpty(field.exceptDomains)}
                        bind:value={field.onlyDomains}
                    />
                    <div class="help-block">Use comma as separator.</div>
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
