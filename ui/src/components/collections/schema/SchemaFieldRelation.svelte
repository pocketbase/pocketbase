<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import { collections, activeCollection } from "@/stores/collections";

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

    const baseFields = ["id", "created", "updated"];

    const authFields = ["username", "email", "emailVisibility", "verified"];

    let upsertPanel = null;
    let displayFieldsList = [];
    let oldCollectionId = null;
    let isSingle = field.options?.maxSelect == 1;
    let oldIsSingle = isSingle;

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

    $: if (oldCollectionId != field.options.collectionId) {
        oldCollectionId = field.options.collectionId;
        refreshDisplayFieldsList();
    }

    function loadDefaults() {
        field.options = {
            maxSelect: 1,
            collectionId: null,
            cascadeDelete: false,
            displayFields: [],
        };
        isSingle = true;
        oldIsSingle = isSingle;
    }

    function refreshDisplayFieldsList() {
        displayFieldsList = baseFields.slice(0);
        if (!selectedColection) {
            return;
        }

        if (selectedColection.isAuth) {
            displayFieldsList = displayFieldsList.concat(authFields);
        }

        for (const f of selectedColection.schema) {
            displayFieldsList.push(f.name);
        }

        // deselect any missing display field
        if (field.options?.displayFields?.length > 0) {
            for (let i = field.options.displayFields.length - 1; i >= 0; i--) {
                if (!displayFieldsList.includes(field.options.displayFields[i])) {
                    field.options.displayFields.splice(i, 1);
                }
            }
        }
    }
</script>

<SchemaField
    bind:field
    {key}
    on:rename
    on:remove
    on:drop
    on:dragstart
    on:dragenter
    on:dragleave
    {...$$restProps}
>
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
                searchable={$collections.length > 5}
                selectPlaceholder={"Select collection *"}
                noOptionsText="No collections found"
                selectionKey="id"
                items={$collections}
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

            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.displayFields" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Display fields</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: "Optionally select the field(s) that will be used in the listings UI. Leave empty for auto.",
                                position: "top",
                            }}
                        />
                    </label>
                    <Select
                        multiple
                        searchable
                        id={uniqueId}
                        selectPlaceholder="Auto"
                        items={displayFieldsList}
                        bind:selected={field.options.displayFields}
                    />
                </Field>
            </div>
            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.cascadeDelete" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Cascade delete</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: `Whether on ${
                                    selectedColection?.name || "relation"
                                } record deletion to delete also the ${
                                    $activeCollection?.name || "field"
                                } associated records.`,
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
        if (e?.detail?.collection?.id) {
            field.options.collectionId = e.detail.collection.id;
        }
    }}
/>
