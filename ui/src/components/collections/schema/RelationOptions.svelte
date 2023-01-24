<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import { collections } from "@/stores/collections";

    export let key = "";
    export let options = {};

    const defaultOptions = [
        { label: "False", value: false },
        { label: "True", value: true },
    ];

    const baseFields = ["id", "created", "updated"];

    const authFields = ["username", "email", "emailVisibility", "verified"];

    let upsertPanel = null;
    let displayFieldsList = [];
    let oldCollectionId = null;

    // load defaults
    $: if (CommonHelper.isEmpty(options)) {
        options = {
            maxSelect: 1,
            collectionId: null,
            cascadeDelete: false,
            displayFields: [],
        };
    }

    $: selectedColection = $collections.find((c) => c.id == options.collectionId) || null;

    $: if (oldCollectionId != options.collectionId) {
        oldCollectionId = options.collectionId;
        refreshDisplayFieldsList();
    }

    function refreshDisplayFieldsList() {
        displayFieldsList = baseFields.slice(0);
        if (!selectedColection) {
            return;
        }

        if (selectedColection.isAuth) {
            displayFieldsList = displayFieldsList.concat(authFields);
        }

        for (const field of selectedColection.schema) {
            displayFieldsList.push(field.name);
        }

        // deselect any missing display field
        if (options?.displayFields?.length > 0) {
            for (let i = options.displayFields.length - 1; i >= 0; i--) {
                if (!displayFieldsList.includes(options.displayFields[i])) {
                    options.displayFields.splice(i, 1);
                }
            }
        }
    }
</script>

<div class="grid">
    <div class="col-sm-9">
        <Field class="form-field required" name="schema.{key}.options.collectionId" let:uniqueId>
            <label for={uniqueId}>Collection</label>
            <ObjectSelect
                searchable={$collections.length > 5}
                selectPlaceholder={"Select collection"}
                noOptionsText="No collections found"
                selectionKey="id"
                items={$collections}
                bind:keyOfSelected={options.collectionId}
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
    </div>
    <div class="col-sm-3">
        <Field class="form-field" name="schema.{key}.options.maxSelect" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">Max select</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: "Leave empty for no limit.",
                        position: "top",
                    }}
                />
            </label>
            <input type="number" id={uniqueId} step="1" min="1" bind:value={options.maxSelect} />
        </Field>
    </div>
    <div class="col-sm-9">
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
                bind:selected={options.displayFields}
            />
        </Field>
    </div>
    <div class="col-sm-3">
        <Field class="form-field" name="schema.{key}.options.cascadeDelete" let:uniqueId>
            <label for={uniqueId}>Cascade delete</label>
            <ObjectSelect id={uniqueId} items={defaultOptions} bind:keyOfSelected={options.cascadeDelete} />
        </Field>
    </div>
</div>

<CollectionUpsertPanel
    bind:this={upsertPanel}
    on:save={(e) => {
        if (e?.detail?.collection?.id) {
            options.collectionId = e.detail.collection.id;
        }
    }}
/>
