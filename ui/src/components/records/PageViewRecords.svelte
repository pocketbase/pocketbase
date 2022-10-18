<script>
    import tooltip from "@/actions/tooltip";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RecordsList from "@/components/views/RecordsList.svelte";
    import ViewsSidebar from "@/components/views/ViewsSidebar.svelte";
    import ViewUpsertPanel from "@/components/views/ViewUpsertPanel.svelte";
    import { hideControls, pageTitle } from "@/stores/app";
    import { activeView, isViewLoading, loadViews, views } from "@/stores/views";
    import { querystring, replace } from "svelte-spa-router";
    import ViewRecordInfo from "../views/ViewRecordInfo.svelte";

    $pageTitle = "Views";

    const queryParams = new URLSearchParams($querystring);

    let viewUpsertPanel;
    let recordPanel;
    let recordsList;
    let filter = queryParams.get("filter") || "";
    let sort = queryParams.get("sort") || "";
    let selectedCollectionId = queryParams.get("viewId") || "";

    // reset filter and sort on collection change
    $: if ($activeView?.id && selectedCollectionId != $activeView.id) {
        reset();
    }

    // keep the url params in sync
    $: if (sort || filter || $activeView?.id) {
        const query = new URLSearchParams({
            collectionId: $activeView?.id || "",
            filter: filter,
            sort: sort,
        }).toString();
        replace("/views?" + query);
    }

    function reset() {
        selectedCollectionId = $activeView.id;
        sort = "";
        filter = "";
    }

    loadViews(selectedCollectionId);
</script>

{#if $isViewLoading}
    <PageWrapper center>
        <div class="placeholder-section m-b-base">
            <span class="loader loader-lg" />
            <h1>Loading views...</h1>
        </div>
    </PageWrapper>
{:else if !$views.length}
    <PageWrapper center>
        <div class="placeholder-section m-b-base">
            <div class="icon">
                <i class="ri-database-2-line" />
            </div>
            {#if $hideControls}
                <h1 class="m-b-10">You don't have any views yet.</h1>
            {:else}
                <h1 class="m-b-10">Create your first view!</h1>
                <button
                    type="button"
                    class="btn btn-expanded-lg btn-lg"
                    on:click={() => viewUpsertPanel?.show()}
                >
                    <i class="ri-add-line" />
                    <span class="txt">Create new View</span>
                </button>
            {/if}
        </div>
    </PageWrapper>
{:else}
    <ViewsSidebar />

    <PageWrapper>
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">Views</div>
                <div class="breadcrumb-item">{$activeView.name}</div>
            </nav>

            <div class="inline-flex gap-5">
                {#if !$hideControls}
                    <button
                        type="button"
                        class="btn btn-secondary btn-circle"
                        use:tooltip={{ text: "Edit View", position: "right" }}
                        on:click={() => viewUpsertPanel?.show($activeView)}
                    >
                        <i class="ri-settings-4-line" />
                    </button>
                {/if}

                <RefreshButton on:refresh={() => recordsList?.load()} />
            </div>

            <div class="btns-group">
                <!-- <button
                    type="button"
                    class="btn btn-outline"
                    on:click={() => collectionDocsPanel?.show($activeView)}
                >
                    <i class="ri-code-s-slash-line" />
                    <span class="txt">API Preview</span>
                </button> -->
            </div>
        </header>

        <Searchbar
            value={filter}
            autocompleteCollection={$activeView}
            on:submit={(e) => (filter = e.detail)}
        />

        <RecordsList
            bind:this={recordsList}
            view={$activeView}
            bind:filter
            bind:sort
            on:select={(e) => recordPanel?.show(e?.detail)}
        />
    </PageWrapper>
{/if}

<ViewUpsertPanel bind:this={viewUpsertPanel} />

<ViewRecordInfo bind:this={recordPanel} view={$activeView} />
