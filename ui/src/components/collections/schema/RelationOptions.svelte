<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";

    export let key = "";
    export let options = {};

    const defaultOptions = [
        { label: "False", value: false },
        { label: "True", value: true },
    ];

    let isLoading = false;
    let collections = [];
    let upsertPanel = null;

    // load defaults
    $: if (CommonHelper.isEmpty(options)) {
        options = {
            maxSelect: 1,
            collectionId: null,
            cascadeDelete: false,
        };
    }

    $: selectedColection = collections.find((c) => c.id == options.collectionId) || null;

    loadCollections();

    async function loadCollections() {
        isLoading = true;

        try {
            const result = await ApiClient.collections.getFullList(200, {
                sort: "created",
            });

            collections = CommonHelper.sortCollections(result);
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }
</script>

<div class="grid">
    <div class="col-sm-9">
        <Field class="form-field required" name="schema.{key}.options.collectionId" let:uniqueId>
            <label for={uniqueId}>Collection</label>
            <ObjectSelect
                searchable={collections.length > 5}
                selectPlaceholder={isLoading ? "Loading..." : "Select collection"}
                noOptionsText="No collections found"
                selectionKey="id"
                items={collections}
                bind:keyOfSelected={options.collectionId}
            >
                <svelte:fragment slot="afterOptions">
                    <button
                        type="button"
                        class="btn btn-warning btn-block btn-sm m-t-5"
                        on:click={() => upsertPanel?.show()}
                    >
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
    <div class="col-sm-12">
        <Field class="form-field" name="schema.{key}.options.cascadeDelete" let:uniqueId>
            <label for={uniqueId}>
                Delete record on {selectedColection ? selectedColection.name : "relation"} delete
            </label>
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
        loadCollections();
    }}
/>
