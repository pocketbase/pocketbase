<script>
    import { tick } from "svelte";
    import { querystring } from "svelte-spa-router";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import CollectionDocsPanel from "@/components/collections/CollectionDocsPanel.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import CollectionsSidebar from "@/components/collections/CollectionsSidebar.svelte";
    import RecordPreviewPanel from "@/components/records/RecordPreviewPanel.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";
    import RecordsCount from "@/components/records/RecordsCount.svelte";
    import RecordsList from "@/components/records/RecordsList.svelte";
    import { hideControls, pageTitle } from "@/stores/app";
    import {
        activeCollection,
        changeActiveCollectionByIdOrName,
        collections,
        isCollectionsLoading,
        loadCollections,
    } from "@/stores/collections";

    const initialQueryParams = new URLSearchParams($querystring);

    let collectionUpsertPanel;
    let collectionDocsPanel;
    let recordUpsertPanel;
    let recordPreviewPanel;
    let recordsList;
    let recordsCount;
    let filter = initialQueryParams.get("filter") || "";
    let sort = initialQueryParams.get("sort") || "-@rowid";
    let selectedCollectionIdOrName = initialQueryParams.get("collection") || $activeCollection?.id;
    let totalCount = 0; // used to manully change the count without the need of reloading the recordsCount component

    loadCollections(selectedCollectionIdOrName);

    $: reactiveParams = new URLSearchParams($querystring);

    $: collectionQueryParam = reactiveParams.get("collection");

    $: if (
        !$isCollectionsLoading &&
        collectionQueryParam &&
        collectionQueryParam != selectedCollectionIdOrName &&
        collectionQueryParam != $activeCollection?.id &&
        collectionQueryParam != $activeCollection?.name
    ) {
        changeActiveCollectionByIdOrName(collectionQueryParam);
    }

    // reset filter and sort on collection change
    $: if (
        $activeCollection?.id &&
        selectedCollectionIdOrName != $activeCollection.id &&
        selectedCollectionIdOrName != $activeCollection.name
    ) {
        reset();
    }

    $: if ($activeCollection?.id) {
        normalizeSort();
    }

    $: if (!$isCollectionsLoading && initialQueryParams.get("recordId")) {
        showRecordById(initialQueryParams.get("recordId"));
    }

    // keep the url params in sync
    $: if (!$isCollectionsLoading && (sort || filter || $activeCollection?.id)) {
        updateQueryParams();
    }

    $: $pageTitle = $activeCollection?.name || "Collections";

    async function showRecordById(recordId) {
        await tick(); // ensure that the reactive component params are resolved

        $activeCollection?.type === "view"
            ? recordPreviewPanel.show(recordId)
            : recordUpsertPanel?.show(recordId);
    }

    function reset() {
        selectedCollectionIdOrName = $activeCollection?.id;
        filter = "";
        sort = "-@rowid";

        normalizeSort();

        updateQueryParams({ recordId: null });

        // close any open collection panels
        collectionUpsertPanel?.forceHide();
        collectionDocsPanel?.hide();
    }

    // ensures that the sort fields exist in the collection
    async function normalizeSort() {
        if (!sort) {
            return; // nothing to normalize
        }

        const collectionFields = CommonHelper.getAllCollectionIdentifiers($activeCollection);

        const sortFields = sort.split(",").map((f) => {
            if (f.startsWith("+") || f.startsWith("-")) {
                return f.substring(1);
            }
            return f;
        });

        // invalid sort expression or missing sort field
        if (sortFields.filter((f) => collectionFields.includes(f)).length != sortFields.length) {
            if ($activeCollection?.type != "view") {
                sort = "-@rowid"; // all collections with exception to the view has this field
            } else if (collectionFields.includes("created")) {
                // common autodate field
                sort = "-created";
            } else {
                sort = "";
            }
        }
    }

    function updateQueryParams(extra = {}) {
        const queryParams = Object.assign(
            {
                collection: $activeCollection?.id || "",
                filter: filter,
                sort: sort,
            },
            extra,
        );

        CommonHelper.replaceHashQueryParams(queryParams);
    }

    // 导出记录为CSV文件
    async function exportRecords() {
        try {
            // 发送导出请求，包含查询参数
            const url = `/api/collections/${$activeCollection.id}/records/export`;
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Accept': 'text/csv',
                },
            });

            if (!response.ok) {
                // 尝试获取服务器返回的具体错误信息
                const errorText = await response.text().catch(() => 'Export failed');
                throw new Error(errorText || 'Export failed');
            }

            // 处理下载文件
            const blob = await response.blob();
            const downloadUrl = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.style.display = 'none';
            link.href = downloadUrl;
            link.download = `${$activeCollection.name}-records-${new Date().toISOString().slice(0, 10)}.csv`;
            document.body.appendChild(link);
            link.click();
            
            // 清理资源
            setTimeout(() => {
                window.URL.revokeObjectURL(downloadUrl);
                document.body.removeChild(link);
            }, 100);
            
        } catch (error) {
            console.error('Export records error:', error);
            // 显示更详细的错误信息
            alert(`Failed to export records: ${error.message || 'Please try again.'}`);
        }
    }

    // 导入记录从CSV文件
    function importRecords() {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.csv';
        input.onchange = async (event) => {
            const file = event.target.files?.[0];
            if (!file) return;

            try {
                const formData = new FormData();
                formData.append('file', file);

                const response = await fetch(`/api/collections/${$activeCollection.id}/records/import`, {
                    method: 'POST',
                    body: formData,
                });

                if (!response.ok) {
                    throw new Error('Import failed');
                }

                // 导入成功后刷新记录列表
                recordsList?.load();
                recordsCount?.reload();
                alert('Records imported successfully!');
            } catch (error) {
                console.error('Import records error:', error);
                alert('Failed to import records. Please check the CSV format and try again.');
            }
        };
        input.click();
    }
</script>

{#if $isCollectionsLoading && !$collections.length}
    <PageWrapper center>
        <div class="placeholder-section m-b-base">
            <span class="loader loader-lg" />
            <h1>Loading collections...</h1>
        </div>
    </PageWrapper>
{:else if !$collections.length}
    <PageWrapper center>
        <div class="placeholder-section m-b-base">
            <div class="icon">
                <i class="ri-database-2-line" />
            </div>
            {#if $hideControls}
                <h1 class="m-b-10">You don't have any collections yet.</h1>
            {:else}
                <h1 class="m-b-10">Create your first collection to add records!</h1>
                <button
                    type="button"
                    class="btn btn-expanded-lg btn-lg"
                    on:click={() => collectionUpsertPanel?.show()}
                >
                    <i class="ri-add-line" />
                    <span class="txt">Create new collection</span>
                </button>
            {/if}
        </div>
    </PageWrapper>
{:else}
    <CollectionsSidebar />

    <PageWrapper class="flex-content">
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">Collections</div>
                <div class="breadcrumb-item">{$activeCollection.name}</div>
            </nav>

            <div class="inline-flex gap-5">
                {#if !$hideControls}
                    <button
                        type="button"
                        aria-label="Edit collection"
                        class="btn btn-transparent btn-circle"
                        use:tooltip={{ text: "Edit collection", position: "right" }}
                        on:click={() => collectionUpsertPanel?.show($activeCollection)}
                    >
                        <i class="ri-settings-4-line" />
                    </button>
                {/if}

                <RefreshButton
                    on:refresh={() => {
                        recordsList?.load();
                        recordsCount?.reload();
                    }}
                />
            </div>

            <div class="btns-group">
                <button
                    type="button"
                    class="btn btn-outline"
                    on:click={() => collectionDocsPanel?.show($activeCollection)}
                >
                    <i class="ri-code-s-slash-line" />
                    <span class="txt">API Preview</span>
                </button>
                <button type="button" class="btn btn-outline" on:click={exportRecords}>
                    <i class="ri-download-2-line" />
                    <span class="txt">Export records</span>
                </button>
                {#if $activeCollection.type !== "view"}
                    <div class="inline-flex gap-2">                        
                        <button type="button" class="btn btn-outline" on:click={importRecords}>
                            <i class="ri-upload-2-line" />
                            <span class="txt">Import records</span>
                        </button>
                        <button type="button" class="btn btn-expanded" on:click={() => recordUpsertPanel?.show()}>
                            <i class="ri-add-line" />
                            <span class="txt">New record</span>
                        </button>
                    </div>
                {/if}
            </div>
        </header>

        <Searchbar
            value={filter}
            autocompleteCollection={$activeCollection}
            on:submit={(e) => (filter = e.detail)}
        />

        <div class="clearfix m-b-sm" />

        <RecordsList
            bind:this={recordsList}
            collection={$activeCollection}
            bind:filter
            bind:sort
            on:select={(e) => {
                updateQueryParams({
                    recordId: e.detail.id,
                });

                let showModel = e.detail._partial ? e.detail.id : e.detail;

                $activeCollection.type === "view"
                    ? recordPreviewPanel?.show(showModel)
                    : recordUpsertPanel?.show(showModel);
            }}
            on:delete={() => {
                recordsCount?.reload();
            }}
            on:new={() => recordUpsertPanel?.show()}
        />

        <svelte:fragment slot="footer">
            <RecordsCount
                bind:this={recordsCount}
                class="m-r-auto txt-sm txt-hint"
                collection={$activeCollection}
                {filter}
                bind:totalCount
            />
        </svelte:fragment>
    </PageWrapper>
{/if}

<CollectionUpsertPanel
    bind:this={collectionUpsertPanel}
    on:truncate={() => {
        recordsList?.load();
        recordsCount?.reload();
    }}
/>

<CollectionDocsPanel bind:this={collectionDocsPanel} />

<RecordUpsertPanel
    bind:this={recordUpsertPanel}
    collection={$activeCollection}
    on:hide={() => {
        updateQueryParams({ recordId: null });
    }}
    on:save={(e) => {
        if (filter) {
            // if there is applied filter, reload the count since we
            // don't know after the save whether the record satisfies it
            recordsCount?.reload();
        } else if (e.detail.isNew) {
            totalCount++;
        }

        recordsList?.reloadLoadedPages();
    }}
    on:delete={(e) => {
        if (!filter || recordsList?.hasRecord(e.detail.id)) {
            totalCount--;
        }

        recordsList?.reloadLoadedPages();
    }}
/>

<RecordPreviewPanel
    bind:this={recordPreviewPanel}
    collection={$activeCollection}
    on:hide={() => {
        updateQueryParams({ recordId: null });
    }}
/>
