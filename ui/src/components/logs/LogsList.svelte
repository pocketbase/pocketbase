<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import SortHeader from "@/components/base/SortHeader.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import HorizontalScroller from "@/components/base/HorizontalScroller.svelte";

    const dispatch = createEventDispatcher();
    const labelMethodClass = {
        get: "label-info",
        post: "label-success",
        patch: "label-warning",
        delete: "label-danger",
    };

    export let filter = "";
    export let presets = "";
    export let sort = "-rowid";

    let items = [];
    let currentPage = 1;
    let totalItems = 0;
    let isLoading = false;
    let yieldedItemsId = 0;

    $: if (typeof sort !== "undefined" || typeof filter !== "undefined" || typeof presets !== "undefined") {
        clearList();
        load(1);
    }

    $: canLoadMore = totalItems > items.length;

    export async function load(page = 1, breakTasks = true) {
        isLoading = true;

        return ApiClient.logs
            .getRequestsList(page, 30, {
                sort: sort,
                filter: [presets, filter].filter(Boolean).join("&&"),
            })
            .then(async (result) => {
                if (page <= 1) {
                    clearList();
                }

                isLoading = false;
                currentPage = result.page;
                totalItems = result.totalItems;
                dispatch("load", items.concat(result.items));

                // optimize the items listing by rendering the rows in task batches
                if (breakTasks) {
                    const currentYieldId = ++yieldedItemsId;
                    while (result.items.length) {
                        if (yieldedItemsId != currentYieldId) {
                            break; // new yeild has been started
                        }

                        items = items.concat(result.items.splice(0, 10));

                        await CommonHelper.yieldToMain();
                    }
                } else {
                    items = items.concat(result.items);
                }
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoading = false;
                    console.warn(err);
                    clearList();
                    ApiClient.error(err, false);
                }
            });
    }

    function clearList() {
        items = [];
        currentPage = 1;
        totalItems = 0;
    }
</script>

<HorizontalScroller class="table-wrapper">
    <table class="table" class:table-loading={isLoading}>
        <thead>
            <tr>
                <SortHeader disable class="col-field-method" name="method" bind:sort>
                    <div class="col-header-content">
                        <i class="ri-global-line" />
                        <span class="txt">Method</span>
                    </div>
                </SortHeader>

                <SortHeader disable class="col-type-text col-field-url" name="url" bind:sort>
                    <div class="col-header-content">
                        <i class={CommonHelper.getFieldTypeIcon("url")} />
                        <span class="txt">URL</span>
                    </div>
                </SortHeader>

                <SortHeader disable class="col-type-text col-field-referer" name="referer" bind:sort>
                    <div class="col-header-content">
                        <i class={CommonHelper.getFieldTypeIcon("url")} />
                        <span class="txt">Referer</span>
                    </div>
                </SortHeader>

                <SortHeader disable class="col-type-number col-field-userIp" name="userIp" bind:sort>
                    <div class="col-header-content">
                        <i class={CommonHelper.getFieldTypeIcon("number")} />
                        <span class="txt">User IP</span>
                    </div>
                </SortHeader>

                <SortHeader disable class="col-type-number col-field-status" name="status" bind:sort>
                    <div class="col-header-content">
                        <i class={CommonHelper.getFieldTypeIcon("number")} />
                        <span class="txt">Status</span>
                    </div>
                </SortHeader>

                <SortHeader disable class="col-type-date col-field-created" name="created" bind:sort>
                    <div class="col-header-content">
                        <i class={CommonHelper.getFieldTypeIcon("date")} />
                        <span class="txt">Created</span>
                    </div>
                </SortHeader>

                <th class="col-type-action min-width" />
            </tr>
        </thead>
        <tbody>
            {#each items as item (item.id)}
                <tr
                    tabindex="0"
                    class="row-handle"
                    on:click={() => dispatch("select", item)}
                    on:keydown={(e) => {
                        if (e.code === "Enter") {
                            e.preventDefault();
                            dispatch("select", item);
                        }
                    }}
                >
                    <td class="col-type-text col-field-method min-width">
                        <span class="label txt-uppercase {labelMethodClass[item.method.toLowerCase()]}">
                            {item.method?.toUpperCase()}
                        </span>
                    </td>

                    <td class="col-type-text col-field-url">
                        <span class="txt txt-ellipsis" title={item.url}>
                            {item.url}
                        </span>
                        {#if item.meta?.errorMessage || item.meta?.errorData}
                            <i class="ri-error-warning-line txt-danger m-l-5 m-r-5" title="Error" />
                        {/if}
                    </td>

                    <td class="col-type-text col-field-referer">
                        <span class="txt txt-ellipsis" class:txt-hint={!item.referer} title={item.referer}>
                            {item.referer || "N/A"}
                        </span>
                    </td>

                    <td class="col-type-number col-field-userIp">
                        <span class="txt txt-ellipsis" class:txt-hint={!item.userIp} title={item.userIp}>
                            {item.userIp || "N/A"}
                        </span>
                    </td>

                    <td class="col-type-number col-field-status">
                        <span class="label" class:label-danger={item.status >= 400}>
                            {item.status}
                        </span>
                    </td>

                    <td class="col-type-date col-field-created">
                        <FormattedDate date={item.created} />
                    </td>

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
                            <h6>No logs found.</h6>
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
</HorizontalScroller>

{#if items.length}
    <small class="block txt-hint txt-right m-t-sm">Showing {items.length} of {totalItems}</small>
{/if}

{#if items.length && canLoadMore}
    <div class="block txt-center m-t-xs">
        <button
            type="button"
            class="btn btn-lg btn-secondary btn-expanded"
            class:btn-loading={isLoading}
            class:btn-disabled={isLoading}
            on:click={() => load(currentPage + 1)}
        >
            <span class="txt">Load more ({totalItems - items.length})</span>
        </button>
    </div>
{/if}
