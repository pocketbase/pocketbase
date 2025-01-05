<script>
    import { createEventDispatcher } from "svelte";
    import { fly } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { collections, isCollectionsLoading } from "@/stores/collections";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Scroller from "@/components/base/Scroller.svelte";
    import SortHeader from "@/components/base/SortHeader.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import RecordFieldValue from "@/components/records/RecordFieldValue.svelte";

    const dispatch = createEventDispatcher();
    const sortRegex = /^([\+\-])?(\w+)$/;
    const perPage = 40;

    export let collection;
    export let sort = "";
    export let filter = "";

    let scrollWrapper;
    let records = [];
    let currentPage = 1;
    let lastTotal = 0;
    let bulkSelected = {};
    let isLoading = true;
    let isDeleting = false;
    let yieldedRecordsId = 0;
    let columnsTrigger;
    let hiddenColumns = [];
    let collumnsToHide = [];
    let hiddenColumnsKey = "";

    const unusedSuperusersFields = ["verified", "emailVisibility"];

    $: if (collection?.id) {
        hiddenColumnsKey = collection.id + "@hiddenColumns";
        loadStoredHiddenColumns();
        clearList();
    }

    $: isView = collection?.type === "view";

    $: isSuperusers = collection?.type === "auth" && collection.name === "_superusers";

    // skip unused superusers fields
    $: fields = (collection?.fields || []).filter(
        (f) => !f.hidden && (!isSuperusers || !unusedSuperusersFields.includes(f.name)),
    );

    $: editorFields = fields.filter((f) => f.type === "editor");

    $: relFields = fields.filter((f) => f.type === "relation");

    $: visibleFields = fields.filter((f) => !hiddenColumns.includes(f.id));

    $: if (!$isCollectionsLoading && collection?.id && sort !== -1 && filter !== -1) {
        load(1);
    }

    $: canLoadMore = lastTotal >= perPage;

    $: totalBulkSelected = Object.keys(bulkSelected).length;

    $: areAllRecordsSelected = records.length && totalBulkSelected === records.length;

    $: if (hiddenColumns !== -1) {
        updateStoredHiddenColumns();
    }

    $: collumnsToHide = fields
        .filter((f) => !f.primaryKey)
        .map((f) => {
            return { id: f.id, name: f.name };
        });

    function updateStoredHiddenColumns() {
        if (!collection?.id) {
            return;
        }

        if (hiddenColumns.length) {
            localStorage.setItem(hiddenColumnsKey, JSON.stringify(hiddenColumns));
        } else {
            localStorage.removeItem(hiddenColumnsKey);
        }
    }

    function loadStoredHiddenColumns() {
        hiddenColumns = [];

        if (!collection?.id) {
            return;
        }

        try {
            const encoded = localStorage.getItem(hiddenColumnsKey);
            if (encoded) {
                hiddenColumns = JSON.parse(encoded) || [];
            }
        } catch (_) {}
    }

    export function hasRecord(id) {
        return !!records.find((r) => r.id == id);
    }

    export async function reloadLoadedPages() {
        const loadedPages = currentPage;

        for (let i = 1; i <= loadedPages; i++) {
            if (i === 1 || canLoadMore) {
                await load(i, false);
            }
        }
    }

    export async function load(page = 1, breakTasks = true) {
        if (!collection?.id) {
            return;
        }

        isLoading = true;

        // allow sorting by the relation display fields
        let listSort = sort;
        const sortMatch = listSort.match(sortRegex);
        const sortRelField = sortMatch ? relFields.find((f) => f.name === sortMatch[2]) : null;
        if (sortMatch && sortRelField) {
            const relPresentableFields =
                $collections
                    ?.find((c) => c.id == sortRelField.collectionId)
                    ?.fields?.filter((f) => f.presentable)
                    ?.map((f) => f.name) || [];

            const parts = [];
            for (const presentableField of relPresentableFields) {
                parts.push((sortMatch[1] || "") + sortMatch[2] + "." + presentableField);
            }
            if (parts.length > 0) {
                listSort = parts.join(",");
            }
        }

        const fallbackSearchFields = CommonHelper.getAllCollectionIdentifiers(collection);

        const listFields = editorFields
            .map((f) => f.name + ":excerpt(200)")
            .concat(relFields.map((field) => "expand." + field.name + ".*:excerpt(200)"));
        if (listFields.length) {
            listFields.unshift("*");
        }

        let expandFields = [];
        for (const field of relFields) {
            expandFields = expandFields.concat(
                CommonHelper.getExpandPresentableRelFields(field, $collections, 2),
            );
        }

        return ApiClient.collection(collection.id)
            .getList(page, perPage, {
                sort: listSort,
                skipTotal: 1,
                filter: CommonHelper.normalizeSearchFilter(filter, fallbackSearchFields),
                expand: expandFields.join(","),
                fields: listFields.join(","),
                requestKey: "records_list",
            })
            .then(async (result) => {
                if (page <= 1) {
                    clearList();
                }

                isLoading = false;
                currentPage = result.page;
                lastTotal = result.items.length;
                dispatch("load", records.concat(result.items));

                // mark the records as "partial" because of the excerpt
                if (editorFields.length) {
                    for (let record of result.items) {
                        record._partial = true;
                    }
                }

                // optimize the records listing by rendering the rows in task batches
                if (breakTasks) {
                    const currentYieldId = ++yieldedRecordsId;
                    while (result.items?.length) {
                        if (yieldedRecordsId != currentYieldId) {
                            break; // new yield has been started
                        }

                        const subset = result.items.splice(0, 20);
                        for (let item of subset) {
                            CommonHelper.pushOrReplaceByKey(records, item);
                        }
                        records = records;

                        await CommonHelper.yieldToMain();
                    }
                } else {
                    for (let item of result.items) {
                        CommonHelper.pushOrReplaceByKey(records, item);
                    }
                    records = records;
                }
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoading = false;
                    console.warn(err);
                    clearList();
                    ApiClient.error(err, !filter || err?.status != 400); // silence filter errors
                }
            });
    }

    function clearList() {
        scrollWrapper?.resetVerticalScroll();
        records = [];
        currentPage = 1;
        lastTotal = 0;
        bulkSelected = {};
    }

    function toggleSelectAllRecords() {
        if (areAllRecordsSelected) {
            deselectAllRecords();
        } else {
            selectAllRecords();
        }
    }

    function deselectAllRecords() {
        bulkSelected = {};
    }

    function selectAllRecords() {
        for (const record of records) {
            bulkSelected[record.id] = record;
        }
        bulkSelected = bulkSelected;
    }

    function toggleSelectRecord(record) {
        if (!bulkSelected[record.id]) {
            bulkSelected[record.id] = record;
        } else {
            delete bulkSelected[record.id];
        }

        bulkSelected = bulkSelected; // trigger reactivity
    }

    function deleteSelectedConfirm() {
        const msg = `Do you really want to delete the selected ${
            totalBulkSelected === 1 ? "record" : "records"
        }?`;

        confirm(msg, deleteSelected);
    }

    async function deleteSelected() {
        if (isDeleting || !totalBulkSelected || !collection?.id) {
            return;
        }

        let promises = [];
        for (const recordId of Object.keys(bulkSelected)) {
            promises.push(ApiClient.collection(collection.id).delete(recordId));
        }

        isDeleting = true;

        return Promise.all(promises)
            .then(() => {
                addSuccessToast(
                    `Successfully deleted the selected ${totalBulkSelected === 1 ? "record" : "records"}.`,
                );

                dispatch("delete", bulkSelected);

                deselectAllRecords();
            })
            .catch((err) => {
                ApiClient.error(err);
            })
            .finally(() => {
                isDeleting = false;

                // always reload because some of the records may not be deletable
                return reloadLoadedPages();
            });
    }
</script>

<Scroller bind:this={scrollWrapper} class="table-wrapper">
    <svelte:fragment slot="before">
        {#if columnsTrigger}
            <Toggler
                class="dropdown dropdown-right dropdown-nowrap columns-dropdown"
                trigger={columnsTrigger}
            >
                <div class="txt-hint txt-sm p-5 m-b-5">Toggle columns</div>
                {#each collumnsToHide as column (column.id + column.name)}
                    <Field class="form-field form-field-sm form-field-toggle m-0 p-5" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            checked={!hiddenColumns.includes(column.id)}
                            on:change={(e) => {
                                if (e.target.checked) {
                                    CommonHelper.removeByValue(hiddenColumns, column.id);
                                } else {
                                    CommonHelper.pushUnique(hiddenColumns, column.id);
                                }
                                hiddenColumns = hiddenColumns;
                            }}
                        />
                        <label for={uniqueId}>{column.name}</label>
                    </Field>
                {/each}
            </Toggler>
        {/if}
    </svelte:fragment>

    <table class="table" class:table-loading={isLoading}>
        <thead>
            <tr>
                {#if !isView}
                    <th class="bulk-select-col min-width">
                        {#if isLoading}
                            <span class="loader loader-sm" />
                        {:else}
                            <div class="form-field">
                                <input
                                    type="checkbox"
                                    id="checkbox_0"
                                    disabled={!records.length}
                                    checked={areAllRecordsSelected}
                                    on:change={() => toggleSelectAllRecords()}
                                />
                                <label for="checkbox_0" />
                            </div>
                        {/if}
                    </th>
                {/if}

                {#each visibleFields as field (field.id)}
                    <SortHeader
                        class="col-type-{field.type} col-field-{field.name}"
                        name={field.name}
                        bind:sort
                    >
                        <div class="col-header-content">
                            {#if field.primaryKey}
                                <i class={CommonHelper.getFieldTypeIcon("primary")} />
                            {:else}
                                <i class={CommonHelper.getFieldTypeIcon(field.type)} />
                            {/if}
                            <span class="txt">{field.name}</span>
                        </div>
                    </SortHeader>
                {/each}

                <th class="col-type-action min-width">
                    {#if collumnsToHide.length}
                        <button
                            bind:this={columnsTrigger}
                            type="button"
                            aria-label="Toggle columns"
                            class="btn btn-sm btn-transparent p-0"
                        >
                            <i class="ri-more-line" />
                        </button>
                    {/if}
                </th>
            </tr>
        </thead>
        <tbody>
            {#each records as record (!isView ? record.id : record)}
                <tr
                    tabindex="0"
                    class="row-handle"
                    on:click={() => dispatch("select", record)}
                    on:keydown={(e) => {
                        if (e.code === "Enter") {
                            e.preventDefault();
                            dispatch("select", record);
                        }
                    }}
                >
                    {#if !isView}
                        <td class="bulk-select-col min-width">
                            <!-- svelte-ignore a11y-click-events-have-key-events -->
                            <!-- svelte-ignore a11y-no-static-element-interactions -->
                            <div class="form-field" on:click|stopPropagation>
                                <input
                                    type="checkbox"
                                    id="checkbox_{record.id}"
                                    checked={bulkSelected[record.id]}
                                    on:change={() => toggleSelectRecord(record)}
                                />
                                <label for="checkbox_{record.id}" />
                            </div>
                        </td>
                    {/if}

                    {#each visibleFields as field (field.id)}
                        <td class="col-type-{field.type} col-field-{field.name}">
                            <RecordFieldValue short {record} {field} />
                        </td>
                    {/each}

                    <td class="col-type-action min-width">
                        <i class="ri-arrow-right-line" />
                    </td>
                </tr>
            {:else}
                {#if isLoading}
                    <tr>
                        <td colspan="99" class="p-xs">
                            <span class="skeleton-loader m-0" />
                        </td>
                    </tr>
                {:else}
                    <tr>
                        <td colspan="99" class="txt-center txt-hint p-xs">
                            <h6>No records found.</h6>
                            {#if filter?.length}
                                <button
                                    type="button"
                                    class="btn btn-hint btn-expanded m-t-sm"
                                    on:click={() => (filter = "")}
                                >
                                    <span class="txt">Clear filters</span>
                                </button>
                            {:else if !isView}
                                <button
                                    type="button"
                                    class="btn btn-secondary btn-expanded m-t-sm"
                                    on:click={() => dispatch("new")}
                                >
                                    <i class="ri-add-line" />
                                    <span class="txt">New record</span>
                                </button>
                            {/if}
                        </td>
                    </tr>
                {/if}
            {/each}

            {#if records.length && canLoadMore}
                <tr>
                    <td colspan="99" class="txt-center">
                        <button
                            class="btn btn-expanded-lg btn-secondary btn-horizontal-sticky"
                            disabled={isLoading}
                            class:btn-loading={isLoading}
                            on:click|preventDefault={() => load(currentPage + 1)}
                        >
                            <span class="txt">Load more</span>
                        </button>
                    </td>
                </tr>
            {/if}
        </tbody>
    </table>
</Scroller>

{#if totalBulkSelected}
    <div class="bulkbar" transition:fly={{ duration: 150, y: 5 }}>
        <div class="txt">
            Selected <strong>{totalBulkSelected}</strong>
            {totalBulkSelected === 1 ? "record" : "records"}
        </div>
        <button
            type="button"
            class="btn btn-xs btn-transparent btn-outline p-l-5 p-r-5"
            class:btn-disabled={isDeleting}
            on:click={() => deselectAllRecords()}
        >
            <span class="txt">Reset</span>
        </button>
        <div class="flex-fill" />
        <button
            type="button"
            class="btn btn-sm btn-transparent btn-danger"
            class:btn-loading={isDeleting}
            class:btn-disabled={isDeleting}
            on:click={() => deleteSelectedConfirm()}
        >
            <span class="txt">Delete selected</span>
        </button>
    </div>
{/if}
