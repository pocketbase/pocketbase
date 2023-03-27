<script>
    import { replace, querystring } from "svelte-spa-router";
    import CommonHelper from "@/utils/CommonHelper";
    import {
        collections,
        activeCollection,
        isCollectionsLoading,
        loadCollections,
        changeActiveCollectionById,
    } from "@/stores/collections";
    import tooltip from "@/actions/tooltip";
    import { pageTitle, hideControls } from "@/stores/app";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import CollectionsSidebar from "@/components/collections/CollectionsSidebar.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import CollectionDocsPanel from "@/components/collections/CollectionDocsPanel.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";
    import RecordPreviewPanel from "@/components/records/RecordPreviewPanel.svelte";
    import RecordsList from "@/components/records/RecordsList.svelte";

    const queryParams = new URLSearchParams($querystring);

    let collectionUpsertPanel;
    let collectionDocsPanel;
    let recordUpsertPanel;
    let recordPreviewPanel;
    let recordsList;
    let filter = queryParams.get("filter") || "";
    let sort = queryParams.get("sort") || "";
    let selectedCollectionId = queryParams.get("collectionId") || $activeCollection?.id;

    $: reactiveParams = new URLSearchParams($querystring);

    $: if (
        !$isCollectionsLoading &&
        reactiveParams.get("collectionId") &&
        reactiveParams.get("collectionId") != selectedCollectionId
    ) {
        changeActiveCollectionById(reactiveParams.get("collectionId"));
    }

    // reset filter and sort on collection change
    $: if ($activeCollection?.id && selectedCollectionId != $activeCollection.id) {
        reset();
    }

    // keep the url params in sync
    $: if (sort || filter || $activeCollection?.id) {
        const query = new URLSearchParams({
            collectionId: $activeCollection?.id || "",
            filter: filter,
            sort: sort,
        }).toString();
        replace("/collections?" + query);
    }

    $: $pageTitle = $activeCollection?.name || "Collections";

    function reset() {
        selectedCollectionId = $activeCollection?.id;
        filter = "";
        sort = "-created";

        // clear default sort if created field is not available
        if (
            $activeCollection?.$isView &&
            !CommonHelper.extractColumnsFromQuery($activeCollection.options.query).includes("created")
        ) {
            sort = "";
        }
    }

    loadCollections(selectedCollectionId);
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

    <PageWrapper>
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

                <RefreshButton on:refresh={() => recordsList?.load()} />
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

                {#if !$activeCollection.$isView}
                    <button type="button" class="btn btn-expanded" on:click={() => recordUpsertPanel?.show()}>
                        <i class="ri-add-line" />
                        <span class="txt">New record</span>
                    </button>
                {/if}
            </div>
        </header>

        <Searchbar
            value={filter}
            autocompleteCollection={$activeCollection}
            on:submit={(e) => (filter = e.detail)}
        />
        <div class="clearfix m-b-base" />

        <RecordsList
            bind:this={recordsList}
            collection={$activeCollection}
            bind:filter
            bind:sort
            on:select={(e) => {
                $activeCollection.$isView
                    ? recordPreviewPanel.show(e?.detail)
                    : recordUpsertPanel?.show(e?.detail);
            }}
            on:new={() => recordUpsertPanel?.show()}
        />
    </PageWrapper>
{/if}

<CollectionUpsertPanel bind:this={collectionUpsertPanel} />

<CollectionDocsPanel bind:this={collectionDocsPanel} />

<RecordUpsertPanel
    bind:this={recordUpsertPanel}
    collection={$activeCollection}
    on:save={() => recordsList?.reloadLoadedPages()}
    on:delete={() => recordsList?.reloadLoadedPages()}
/>

<RecordPreviewPanel bind:this={recordPreviewPanel} collection={$activeCollection} />
