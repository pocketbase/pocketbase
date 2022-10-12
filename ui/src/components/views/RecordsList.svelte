<script>
    import SortHeader from "@/components/base/SortHeader.svelte";
    import RecordFieldCell from "@/components/views/RecordFieldCell.svelte";
    import ApiClient from "@/utils/ApiClient";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let view;
    export let sort = "";
    export let filter = "";

    let records = [];
    let recordPanel;
    let currentPage = 1;
    let totalRecords = 0;
    let isLoading = true;

    $: if (view?.id) {
        clearList();
    }

    $: if (view?.id && sort !== -1 && filter !== -1) {
        load(1);
    }

    $: canLoadMore = totalRecords > records.length;

    $: fields = view?.schema || [];

    export async function load(page = 1) {
        if (!view?.id) {
            return;
        }

        isLoading = true;

        return ApiClient.views
            .getRecordsList(view.id, page, 50, {
                // sort: sort,
                filter: filter,
            })
            .then((result) => {
                if (page <= 1) {
                    clearList();
                }

                isLoading = false;
                records = records.concat(result.items);
                currentPage = result.page;
                totalRecords = result.totalItems;

                dispatch("load", records);
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoading = false;
                    console.warn(err);
                    clearList();
                    ApiClient.errorResponseHandler(err, false);
                }
            });
    }

    function clearList() {
        records = [];
        currentPage = 1;
        totalRecords = 0;
    }
</script>

<div class="table-wrapper">
    <table class="table" class:table-loading={isLoading}>
        <thead>
            <tr>
                <th class="min-width">
                    {#if isLoading}
                        <span class="loader loader-sm" />
                    {/if}
                </th>
                {#each fields as field (field.name)}
                    <SortHeader
                        class="col-type-{field.type} col-field-{field.name}"
                        name={field.name}
                        bind:sort
                    >
                        <div class="col-header-content">
                            <span class="txt">{field.name}</span>
                        </div>
                    </SortHeader>
                {/each}
            </tr>
        </thead>
        <tbody>
            {#each records as record (record.id)}
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
                    <td class="min-width" />

                    {#each fields as field (field.id)}
                        <RecordFieldCell {record} {field} />
                    {/each}

                    <td class="col-type-action min-width">
                        <i class="ri-arrow-right-line" />
                    </td>
                </tr>
            {:else}
                {#if isLoading}
                    <tr>
                        <td colspan="99" class="p-xs">
                            <span class="skeleton-loader" />
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
                            {/if}
                        </td>
                    </tr>
                {/if}
            {/each}
        </tbody>
    </table>
</div>

{#if records.length}
    <small class="block txt-hint txt-right m-t-sm">Showing {records.length} of {totalRecords}</small>
{/if}

{#if records.length && canLoadMore}
    <div class="block txt-center m-t-xs">
        <button
            type="button"
            class="btn btn-lg btn-secondary btn-expanded"
            class:btn-loading={isLoading}
            class:btn-disabled={isLoading}
            on:click={() => load(currentPage + 1)}
        >
            <span class="txt">Load more ({totalRecords - records.length})</span>
        </button>
    </div>
{/if}
