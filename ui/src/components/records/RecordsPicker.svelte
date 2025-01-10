<script>
    import scrollend from "@/actions/scrollend";
    import tooltip from "@/actions/tooltip";
    import Draggable from "@/components/base/Draggable.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";
    import { collections } from "@/stores/collections";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();
    const uniqueId = "picker_" + CommonHelper.randomString(5);
    const batchSize = 50;

    export let value;
    export let field;

    let pickerPanel;
    let upsertPanel;
    let filter = "";
    let list = [];
    let selected = [];
    let currentPage = 1;
    let lastItemsCount = 0;
    let isLoadingList = false;
    let isLoadingSelected = false;
    let isReloadingRecord = {};

    $: maxSelect = field?.maxSelect || null;

    $: collectionId = field?.collectionId;

    $: collection = $collections.find((c) => c.id == collectionId) || null;

    $: if (typeof filter !== "undefined" && pickerPanel?.isActive()) {
        loadList(true); // reset list on filter change
    }

    $: isView = collection?.type === "view";

    $: isLoading = isLoadingList || isLoadingSelected;

    $: canLoadMore = lastItemsCount == batchSize;

    $: canSelectMore = maxSelect <= 0 || maxSelect > selected.length;

    export function show() {
        filter = "";
        list = [];
        selected = [];
        loadSelected();
        loadList(true);

        return pickerPanel?.show();
    }

    export function hide() {
        return pickerPanel?.hide();
    }

    function getExpand() {
        let expands = [];

        const presentableRelFields = collection?.fields?.filter(
            (f) => !f.hidden && f.presentable && f.type == "relation",
        );
        for (const field of presentableRelFields) {
            expands = expands.concat(CommonHelper.getExpandPresentableRelFields(field, $collections, 2));
        }

        return expands.join(",");
    }

    async function loadSelected() {
        const selectedIds = CommonHelper.toArray(value);

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
            for (const id of filterIds.splice(0, batchSize)) {
                filters.push(`id="${id}"`);
            }

            loadPromises.push(
                ApiClient.collection(collectionId).getFullList({
                    batch: batchSize,
                    filter: filters.join("||"),
                    fields: "*:excerpt(200)",
                    expand: getExpand(),
                    requestKey: null,
                }),
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

            if (!filter.trim()) {
                // add the selected models to the list (if not already)
                list = CommonHelper.filterDuplicatesByKey(selected.concat(list));
            }

            isLoadingSelected = false;
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
                isLoadingSelected = false;
            }
        }
    }

    async function loadList(reset = false) {
        if (!collectionId) {
            return;
        }

        isLoadingList = true;

        if (reset) {
            if (!filter.trim()) {
                // prepend the loaded selected items
                list = CommonHelper.toArray(selected).slice();
            } else {
                list = [];
            }
        }

        try {
            const page = reset ? 1 : currentPage + 1;

            const fallbackSearchFields = CommonHelper.getAllCollectionIdentifiers(collection);

            let sort = "";
            if (!isView) {
                sort = "-@rowid"; // all collections with exception to the view has this field
            }

            const result = await ApiClient.collection(collectionId).getList(page, batchSize, {
                filter: CommonHelper.normalizeSearchFilter(filter, fallbackSearchFields),
                sort: sort,
                fields: "*:excerpt(200)",
                skipTotal: 1,
                expand: getExpand(),
                requestKey: uniqueId + "loadList",
            });

            list = CommonHelper.filterDuplicatesByKey(list.concat(result.items));
            currentPage = result.page;
            lastItemsCount = result.items.length;

            isLoadingList = false;
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
                isLoadingList = false;
            }
        }
    }

    async function reloadRecord(record) {
        if (!record?.id) {
            return;
        }

        isReloadingRecord[record.id] = true;

        try {
            const reloaded = await ApiClient.collection(collectionId).getOne(record.id, {
                fields: "*:excerpt(200)",
                expand: getExpand(),
                requestKey: uniqueId + "reload" + record.id,
            });

            CommonHelper.pushOrReplaceByKey(selected, reloaded);
            CommonHelper.pushOrReplaceByKey(list, reloaded);
            selected = selected;
            list = list;

            isReloadingRecord[record.id] = false;
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
                isReloadingRecord[record.id] = false;
            }
        }
    }

    $: isSelected = function (record) {
        return CommonHelper.findByKey(selected, "id", record.id);
    };

    function select(record) {
        if (maxSelect == 1) {
            selected = [record];
        } else if (canSelectMore) {
            CommonHelper.pushOrReplaceByKey(selected, record);
            selected = selected;
        }
    }

    function deselect(record) {
        CommonHelper.removeByKey(selected, "id", record.id);
        selected = selected;
    }

    function toggle(record) {
        if (isSelected(record)) {
            deselect(record);
        } else {
            select(record);
        }
    }

    function save() {
        if (maxSelect != 1) {
            value = selected.map((r) => r.id);
        } else {
            value = selected?.[0]?.id || "";
        }

        dispatch("save", selected);
        hide();
    }
</script>

<OverlayPanel bind:this={pickerPanel} popup class="overlay-panel-xl" on:hide on:show {...$$restProps}>
    <svelte:fragment slot="header">
        <h4>
            Select <strong>{collection?.name || ""}</strong> records
        </h4>
    </svelte:fragment>

    <div class="flex m-b-base flex-gap-10">
        <Searchbar
            value={filter}
            autocompleteCollection={collection}
            on:submit={(e) => (filter = e.detail)}
        />
        {#if !isView}
            <button
                type="button"
                class="btn btn-pill btn-transparent btn-hint p-l-xs p-r-xs"
                on:click={() => upsertPanel?.show()}
            >
                <div class="txt">New record</div>
            </button>
        {/if}
    </div>

    <div
        class="list picker-list m-b-base"
        use:scrollend={() => {
            if (canLoadMore && !isLoadingList) {
                loadList();
            }
        }}
    >
        {#each list as record (record.id)}
            {@const selected = isSelected(record)}

            <!-- svelte-ignore a11y-no-static-element-interactions -->
            <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
            <div
                tabindex="0"
                class="list-item handle"
                class:selected
                class:disabled={isReloadingRecord[record.id] ||
                    (!selected && maxSelect > 1 && !canSelectMore)}
                on:click={() => toggle(record)}
                on:keydown={(e) => {
                    if (e.code === "Enter" || e.code === "Space") {
                        e.preventDefault();
                        e.stopPropagation();
                        toggle(record);
                    }
                }}
            >
                {#if selected}
                    <i class="ri-checkbox-circle-fill txt-success" />
                {:else}
                    <i class="ri-checkbox-blank-circle-line txt-disabled" />
                {/if}
                <div class="content">
                    {#if isReloadingRecord[record.id]}
                        <span class="loader loader-xs active"></span>
                    {:else}
                        <RecordInfo {record} />
                    {/if}
                </div>
                {#if !isView}
                    <div class="actions nonintrusive">
                        <button
                            type="button"
                            class="btn btn-sm btn-circle btn-transparent btn-hint m-l-auto"
                            use:tooltip={"Edit"}
                            on:keydown|stopPropagation
                            on:click|stopPropagation={() => upsertPanel?.show(record.id)}
                        >
                            <i class="ri-pencil-line" />
                        </button>
                    </div>
                {/if}
            </div>
        {:else}
            {#if !isLoading}
                <div class="list-item">
                    <span class="txt txt-hint">No records found.</span>
                    {#if filter?.length}
                        <button type="button" class="btn btn-hint btn-sm" on:click={() => (filter = "")}>
                            <span class="txt">Clear filters</span>
                        </button>
                    {/if}
                </div>
            {/if}
        {/each}

        {#if isLoading}
            <div class="list-item">
                <div class="block txt-center">
                    <span class="loader loader-sm active" />
                </div>
            </div>
        {/if}
    </div>

    <h5 class="section-title">
        Selected
        {#if maxSelect > 1}
            ({selected.length} of MAX {maxSelect})
        {/if}
    </h5>
    {#if selected.length}
        <div class="selected-list">
            {#each selected as record, i}
                <Draggable bind:list={selected} index={i} let:dragging let:dragover>
                    <span class="label" class:label-danger={dragging} class:label-warning={dragover}>
                        {#if isReloadingRecord[record.id]}
                            <span class="loader loader-xs active"></span>
                        {:else}
                            <RecordInfo {record} />
                        {/if}
                        <button
                            type="button"
                            title="Remove"
                            class="btn btn-circle btn-transparent btn-hint btn-xs"
                            on:click={() => deselect(record)}
                        >
                            <i class="ri-close-line" />
                        </button>
                    </span>
                </Draggable>
            {/each}
        </div>
    {:else}
        <p class="txt-hint">No selected records.</p>
    {/if}

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button type="button" class="btn" on:click={() => save()}>
            <span class="txt">Set selection</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<RecordUpsertPanel
    bind:this={upsertPanel}
    {collection}
    on:save={(e) => {
        CommonHelper.removeByKey(list, "id", e.detail.record.id);
        list.unshift(e.detail.record);
        list = list;

        select(e.detail.record);

        reloadRecord(e.detail.record);
    }}
    on:delete={(e) => {
        CommonHelper.removeByKey(list, "id", e.detail.id);
        list = list;

        deselect(e.detail);
    }}
/>

<style lang="scss">
    .picker-list {
        max-height: 380px;
    }
    .selected-list {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 10px;
        max-height: 220px;
        overflow: auto;
    }
</style>
