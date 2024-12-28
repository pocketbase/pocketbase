<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.min" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Min length</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={"Clear the field or set it to 0 for no limit."}
                        />
                    </label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        min="0"
                        max={Number.MAX_SAFE_INTEGER}
                        placeholder="No min limit"
                        value={field.min || ""}
                        on:input={(e) => (field.min = parseInt(e.target.value, 10))}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.max" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Max length</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={"Clear the field or set it to 0 to fallback to the default limit."}
                        />
                    </label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        placeholder="Default to max 5000 characters"
                        min={field.min || 0}
                        max={Number.MAX_SAFE_INTEGER}
                        value={field.max || ""}
                        on:input={(e) => (field.max = parseInt(e.target.value, 10))}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.pattern" let:uniqueId>
                    <label for={uniqueId}>Validation pattern</label>
                    <input type="text" id={uniqueId} bind:value={field.pattern} />
                    <div class="help-block">
                        <p>Ex. <code>{"^[a-z0-9]+$"}</code></p>
                    </div>
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.autogeneratePattern" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Autogenerate pattern</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={"Set and autogenerate text matching the pattern on missing record create value."}
                        />
                    </label>
                    <input type="text" id={uniqueId} bind:value={field.autogeneratePattern} />
                    <div class="help-block">
                        <p>Ex. <code>{"[a-z0-9]{30}"}</code></p>
                    </div>
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
