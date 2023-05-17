<script>
    import { onDestroy } from "svelte";
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Draggable from "@/components/base/Draggable.svelte";
    import RecordsPicker from "@/components/records/RecordsPicker.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";

    const batchSize = 100;

    export let value;
    export let picker;
    export let field = new SchemaField();

    let fieldRef;
    let list = [];
    let isLoading = false;
    let loadTimeoutId;

    $: isMultiple = field.options?.maxSelect != 1;

    $: if (typeof value != "undefined") {
        fieldRef?.changed();
    }

    $: if (needLoad(list, value)) {
        // Move the load function to the end of the execution queue.
        //
        // It helps reducing the layout shifts (the relation field has fixed height skeleton loader)
        // and allows the other form fields to load sooner.
        clearTimeout(loadTimeoutId);
        loadTimeoutId = setTimeout(load, 0);
    }

    function needLoad() {
        if (isLoading) {
            return false;
        }

        const ids = CommonHelper.toArray(value);

        list = list.filter((item) => ids.includes(item.id));

        return ids.length != list.length;
    }

    async function load() {
        const ids = CommonHelper.toArray(value);

        list = []; // reset

        if (!field?.options?.collectionId || !ids.length) {
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
            ApiClient.error(err);
        }

        isLoading = false;
    }

    function remove(rel) {
        CommonHelper.removeByKey(list, "id", rel.id);
        list = list;

        listToValue();
    }

    function listToValue() {
        if (isMultiple) {
            value = list.map((r) => r.id);
        } else {
            value = list[0]?.id || "";
        }
    }

    onDestroy(() => {
        clearTimeout(loadTimeoutId);
    });
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
            {#each list as record, i (record.id)}
                <Draggable
                    bind:list
                    group={field.name + "_relation"}
                    index={i}
                    disabled={!isMultiple}
                    let:dragging
                    let:dragover
                    on:sort={() => {
                        listToValue();
                    }}
                >
                    <div class="list-item" class:dragging class:dragover>
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
                </Draggable>
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
