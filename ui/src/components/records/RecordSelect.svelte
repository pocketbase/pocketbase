<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import RecordSelectOption from "./RecordSelectOption.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";

    const uniqueId = "select_" + CommonHelper.randomString(5);

    // original select props
    export let multiple = false;
    export let selected = [];
    export let keyOfSelected = multiple ? [] : undefined;
    export let selectPlaceholder = "- Select -";
    export let optionComponent = RecordSelectOption; // custom component to use for each dropdown option item

    // custom props
    export let collectionId;

    let list = [];
    let currentPage = 1;
    let totalItems = 0;
    let isLoadingList = false;
    let isLoadingSelected = false;
    let isLoadingCollection = false;
    let collection = null;
    let upsertPanel;

    $: if (collectionId) {
        loadCollection();
        loadSelected().then(() => {
            loadList(true);
        });
    }

    $: isLoading = isLoadingList || isLoadingSelected;

    $: canLoadMore = totalItems > list.length;

    async function loadCollection() {
        if (!collectionId) {
            collection = null;
            isLoadingCollection = false;
            return;
        }

        isLoadingCollection = true;

        try {
            collection = await ApiClient.collections.getOne(collectionId, {
                $cancelKey: "collection_" + uniqueId,
            });
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingCollection = false;
    }

    async function loadSelected() {
        const selectedIds = CommonHelper.toArray(keyOfSelected);

        if (!collectionId || !selectedIds.length) {
            return;
        }

        isLoadingSelected = true;

        let loadedItems = [];

        // batch load all selected records to avoid parser stack overflow errors
        const filterIds = selectedIds.slice();
        const loadPromises = [];
        while (filterIds.length > 0) {
            const filters = [];
            for (const id of filterIds.splice(0, 50)) {
                filters.push(`id="${id}"`);
            }

            loadPromises.push(
                ApiClient.collection(collectionId).getFullList(200, {
                    filter: filters.join("||"),
                    $autoCancel: false,
                })
            );
        }

        try {
            await Promise.all(loadPromises).then((values) => {
                loadedItems = loadedItems.concat(...values);
            });

            // preserve selected order
            selected = [];
            for (const id of selectedIds) {
                const item = CommonHelper.findByKey(loadedItems, "id", id);
                if (item) {
                    selected.push(item);
                }
            }

            // add the selected models to the list (if not already)
            list = CommonHelper.filterDuplicatesByKey(selected.concat(list));
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingSelected = false;
    }

    async function loadList(reset = false) {
        if (!collectionId) {
            return;
        }

        isLoadingList = true;

        try {
            const page = reset ? 1 : currentPage + 1;

            const result = await ApiClient.collection(collectionId).getList(page, 200, {
                sort: "-created",
                $cancelKey: uniqueId + "loadList",
            });

            if (reset) {
                list = CommonHelper.toArray(selected).slice();
            }

            list = CommonHelper.filterDuplicatesByKey(
                list.concat(result.items, CommonHelper.toArray(selected))
            );
            currentPage = result.page;
            totalItems = result.totalItems;
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingList = false;
    }
</script>

<ObjectSelect
    selectPlaceholder={isLoading ? "Loading..." : selectPlaceholder}
    items={list}
    searchable={list.length > 5}
    selectionKey="id"
    labelComponent={optionComponent}
    disabled={isLoading}
    {optionComponent}
    {multiple}
    bind:keyOfSelected
    bind:selected
    on:show
    on:hide
    class="records-select block-options"
    {...$$restProps}
>
    <svelte:fragment slot="afterOptions">
        {#if !isLoadingCollection && collection}
            <button
                type="button"
                class="btn btn-warning btn-block btn-sm m-t-5"
                on:click={() => upsertPanel?.show()}
            >
                <span class="txt">New record</span>
            </button>
        {/if}
        {#if canLoadMore}
            <button
                type="button"
                class="btn btn-block btn-sm m-t-5"
                class:btn-loading={isLoadingList}
                class:btn-disabled={isLoadingList}
                on:click|stopPropagation={() => loadList()}
            >
                <span class="txt">Load more</span>
            </button>
        {/if}
    </svelte:fragment>
</ObjectSelect>

<RecordUpsertPanel
    bind:this={upsertPanel}
    {collection}
    on:save={(e) => {
        if (e?.detail?.id) {
            keyOfSelected = CommonHelper.toArray(keyOfSelected).concat(e.detail.id);
        }
        loadList(true);
    }}
/>
