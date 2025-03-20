<script>
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import DynamicOptionsSelect from "@/components/base/DynamicOptionsSelect.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    const isSingleOptions = [
        { label: "Single", value: true },
        { label: "Multiple", value: false },
    ];

    let isSingle = field.maxSelect <= 1;
    let oldIsSingle = isSingle;

    $: if (typeof field.maxSelect == "undefined") {
        loadDefaults();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.maxSelect = 1;
        } else {
            field.maxSelect = field.values?.length || 2;
        }
    }

    function loadDefaults() {
        field.maxSelect = 1;
        field.values = [];
        isSingle = true;
        oldIsSingle = isSingle;
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment let:interactive>
        <div class="separator" />

        <Field
            class="form-field required {!interactive ? 'readonly' : ''}"
            inlineError
            name="fields.{key}.values"
            let:uniqueId
        >
            <DynamicOptionsSelect
                id={uniqueId}
                emptyPlaceholder={"Add choices *"}
                bind:items={field.values}
            />
        </Field>

        <div class="separator" />

        <Field
            class="form-field form-field-single-multiple-select {!interactive ? 'readonly' : ''}"
            inlineError
            let:uniqueId
        >
            <ObjectSelect
                id={uniqueId}
                items={isSingleOptions}
                readonly={!interactive}
                bind:keyOfSelected={isSingle}
            />
        </Field>

        <div class="separator" />
    </svelte:fragment>

    <svelte:fragment slot="options">
        {#if !isSingle}
            <Field class="form-field" name="fields.{key}.maxSelect" let:uniqueId>
                <label for={uniqueId}>Max select</label>
                <input
                    id={uniqueId}
                    type="number"
                    step="1"
                    min="2"
                    max={field.values.length}
                    placeholder="Default to single"
                    bind:value={field.maxSelect}
                />
            </Field>
        {/if}
    </svelte:fragment>
</SchemaField>
