<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
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
    let isSingle = field.options?.maxSelect == 1;
    let oldIsSingle = isSingle;

    $: selectCollections = $collections.filter((c) => c.type != "view");

    // load defaults
    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.options.minSelect = null;
            field.options.maxSelect = 1;
        } else {
            field.options.maxSelect = null;
        }
    }

    $: selectedColection = $collections.find((c) => c.id == field.options.collectionId) || null;

    function loadDefaults() {
        field.options = {
            maxSelect: 1,
            collectionId: null,
            cascadeDelete: false,
        };
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
            name="schema.{key}.options.collectionId"
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
                bind:keyOfSelected={field.options.collectionId}
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
                    <Field class="form-field" name="schema.{key}.options.minSelect" let:uniqueId>
                        <label for={uniqueId}>Min select</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            min="1"
                            placeholder="No min limit"
                            bind:value={field.options.minSelect}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field class="form-field" name="schema.{key}.options.maxSelect" let:uniqueId>
                        <label for={uniqueId}>Max select</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            placeholder="No max limit"
                            min={field.options.minSelect || 2}
                            bind:value={field.options.maxSelect}
                        />
                    </Field>
                </div>
            {/if}

            <div class="col-sm-12">
                <Field class="form-field" name="schema.{key}.options.cascadeDelete" let:uniqueId>
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
                        bind:keyOfSelected={field.options.cascadeDelete}
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
            field.options.collectionId = e.detail.collection.id;
        }
    }}
/>
