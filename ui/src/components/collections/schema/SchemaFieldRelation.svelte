<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import { collections } from "@/stores/collections";

    export let field;
    export let key = "";

    const isSingleOptions = [
        { label: "Single", value: true },
        { label: "Multiple", value: false },
    ];

    const defaultOptions = [
        { label: "False", value: false },
        { label: "True", value: true },
    ];

    let upsertPanel = null;
    let isSingle = field.maxSelect <= 1;
    let oldIsSingle = isSingle;

    $: selectCollections = $collections.filter((c) => !c.system && c.type != "view");

    // load defaults
    $: if (typeof field.maxSelect == "undefined") {
        loadDefaults();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.minSelect = 0;
            field.maxSelect = 1;
        } else {
            field.maxSelect = 999;
        }
    }

    $: selectedColection = $collections.find((c) => c.id == field.collectionId) || null;

    function loadDefaults() {
        field.maxSelect = 1;
        field.collectionId = null;
        field.cascadeDelete = false;
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
            name="fields.{key}.collectionId"
            let:uniqueId
        >
            <ObjectSelect
                id={uniqueId}
                searchable={selectCollections.length > 5}
                selectPlaceholder={"Select collection *"}
                noOptionsText="No collections found"
                selectionKey="id"
                items={selectCollections}
                readonly={!interactive || field.id}
                bind:keyOfSelected={field.collectionId}
            >
                <svelte:fragment slot="afterOptions">
                    <hr />
                    <button
                        type="button"
                        class="btn btn-transparent btn-block btn-sm"
                        on:click={() => upsertPanel?.show()}
                    >
                        <i class="ri-add-line" />
                        <span class="txt">New collection</span>
                    </button>
                </svelte:fragment>
            </ObjectSelect>
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
        <div class="grid grid-sm">
            {#if !isSingle}
                <div class="col-sm-6">
                    <Field class="form-field" name="fields.{key}.minSelect" let:uniqueId>
                        <label for={uniqueId}>Min select</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            min="0"
                            placeholder="No min limit"
                            value={field.minSelect || ""}
                            on:input={(e) => (field.minSelect = e.target.value << 0)}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field class="form-field" name="fields.{key}.maxSelect" let:uniqueId>
                        <label for={uniqueId}>Max select</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            placeholder="Default to single"
                            min={field.minSelect || 1}
                            bind:value={field.maxSelect}
                        />
                    </Field>
                </div>
            {/if}

            <div class="col-sm-12">
                <Field class="form-field" name="fields.{key}.cascadeDelete" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Cascade delete</span>
                        <!-- prettier-ignore -->
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: [
                                    `Whether on ${selectedColection?.name || "relation"} record deletion to delete also the current corresponding collection record(s).`,
                                    !isSingle ? `For "Multiple" relation fields the cascade delete is triggered only when all ${selectedColection?.name || "relation"} ids are removed from the corresponding record.` : null
                                ].filter(Boolean).join("\n\n"),
                                position: "top",
                            }}
                        />
                    </label>
                    <ObjectSelect
                        id={uniqueId}
                        items={defaultOptions}
                        bind:keyOfSelected={field.cascadeDelete}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>

<CollectionUpsertPanel
    bind:this={upsertPanel}
    on:save={(e) => {
        if (e?.detail?.collection?.id && e.detail.collection.type != "view") {
            field.collectionId = e.detail.collection.id;
        }
    }}
/>
