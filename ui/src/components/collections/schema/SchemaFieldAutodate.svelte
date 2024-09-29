<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    const ON_CREATE = 1;
    const ON_UPDATE = 2;
    const ON_CREATE_UPDATE = 3;

    const options = [
        { label: "Create", value: ON_CREATE },
        { label: "Update", value: ON_UPDATE },
        { label: "Create/Update", value: ON_CREATE_UPDATE },
    ];

    export let field;
    export let key = "";

    let selectedOption = optionFromField();

    $: updateField(selectedOption);

    function optionFromField() {
        if (field.onCreate && field.onUpdate) {
            return ON_CREATE_UPDATE;
        }

        if (field.onUpdate) {
            return ON_UPDATE;
        }

        return ON_CREATE;
    }

    function updateField(option) {
        switch (option) {
            case ON_CREATE:
                field.onCreate = true;
                field.onUpdate = false;
                break;
            case ON_UPDATE:
                field.onCreate = false;
                field.onUpdate = true;
                break;
            case ON_CREATE_UPDATE:
                field.onCreate = true;
                field.onUpdate = true;
                break;
        }
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment let:interactive>
        <div class="separator" />

        <Field
            class="form-field form-field-single-multiple-select form-field-autodate-select {!interactive
                ? 'readonly'
                : ''}"
            inlineError
            let:uniqueId
        >
            <div use:tooltip={{ text: "Auto set on:", position: "top" }}>
                <ObjectSelect
                    id={uniqueId}
                    items={options}
                    disabled={field.system}
                    readonly={!interactive}
                    bind:keyOfSelected={selectedOption}
                />
            </div>
        </Field>

        <div class="separator" />
    </svelte:fragment>
</SchemaField>
