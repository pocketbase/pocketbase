<script>
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import RecordsPicker from "@/components/records/RecordsPicker.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";

    const batchSize = 100;

    export let value;
    export let picker;
    export let field = new SchemaField();

    let fieldRef;
    let list = [];
    let isLoading = false;

    $: isMultiple = field.options?.maxSelect != 1;

    $: if (typeof value != "undefined") {
        fieldRef?.changed();
    }

    load();

    async function load() {
        const ids = CommonHelper.toArray(value);

        if (!field?.options?.collectionId || !ids.length) {
            list = [];
            isLoading = false;
            return;
        }

        isLoading = true;

        // batch load all selected records to avoid parser stack overflow errors
        const filterIds = ids.slice();
        const loadPromises = [];
        while (filterIds.length > 0) {
            const filters = [];
            for (const id of filterIds.splice(0, batchSize)) {
                filters.push(`id="${id}"`);
            }

            loadPromises.push(
                ApiClient.collection(field?.options?.collectionId).getFullList(batchSize, {
                    filter: filters.join("||"),
                    $autoCancel: false,
                })
            );
        }

        try {
            let loadedItems = [];
            await Promise.all(loadPromises).then((values) => {
                loadedItems = loadedItems.concat(...values);
            });

            // preserve selected order
            for (const id of ids) {
                const rel = CommonHelper.findByKey(loadedItems, "id", id);
                if (rel) {
                    list.push(rel);
                }
            }

            list = list;
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }

    function remove(rel) {
        CommonHelper.removeByKey(list, "id", rel.id);
        list = list;

        if (isMultiple) {
            value = list.map((r) => r.id);
        } else {
            value = list[0]?.id || "";
        }
    }
</script>

<Field
    bind:this={fieldRef}
    class="form-field form-field-list {field.required ? 'required' : ''}"
    name={field.name}
    let:uniqueId
>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>

    <div class="list">
        <div class="relations-list">
            {#each list as record}
                <div class="list-item">
                    <div class="content">
                        <RecordInfo {record} displayFields={field.options?.displayFields} />
                    </div>
                    <div class="actions">
                        <button
                            type="button"
                            class="btn btn-transparent btn-hint btn-sm btn-circle btn-remove"
                            use:tooltip={"Remove"}
                            on:click={() => remove(record)}
                        >
                            <i class="ri-close-line" />
                        </button>
                    </div>
                </div>
            {:else}
                {#if isLoading}
                    {#each CommonHelper.toArray(value).slice(0, 10) as _}
                        <div class="list-item">
                            <div class="skeleton-loader" />
                        </div>
                    {/each}
                {/if}
            {/each}
        </div>

        <div class="list-item list-item-btn">
            <button
                type="button"
                class="btn btn-transparent btn-sm btn-block"
                on:click={() => picker?.show()}
            >
                <i class="ri-magic-line" />
                <!-- <i class="ri-layout-line" /> -->
                <span class="txt">Open picker</span>
            </button>
        </div>
    </div>
</Field>

<RecordsPicker
    bind:this={picker}
    {value}
    {field}
    on:save={(e) => {
        list = e.detail || [];
        value = isMultiple ? list.map((r) => r.id) : list[0]?.id || "";
    }}
/>

<style lang="scss">
    .relations-list {
        max-height: 300px;
        overflow: auto; /* fallback */
        overflow: overlay;
    }
</style>
