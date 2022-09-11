<script>
    import { hideControls } from "@/stores/app";
    import { collections, activeCollection } from "@/stores/collections";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";

    let collectionPanel;
    let searchTerm = "";

    $: normalizedSearch = searchTerm.replace(/\s+/g, "").toLowerCase();

    $: hasSearch = searchTerm !== "";

    $: filteredCollections = $collections.filter((collection) => {
        return (
            collection.name != import.meta.env.PB_PROFILE_COLLECTION &&
            (collection.id == searchTerm ||
                collection.name.replace(/\s+/g, "").toLowerCase().includes(normalizedSearch))
        );
    });

    function selectCollection(collection) {
        $activeCollection = collection;
    }
</script>

<aside class="page-sidebar collection-sidebar">
    <header class="sidebar-header">
        <div class="form-field search" class:active={hasSearch}>
            <div class="form-field-addon">
                <button
                    type="button"
                    class="btn btn-xs btn-secondary btn-circle btn-clear"
                    class:hidden={!hasSearch}
                    on:click={() => (searchTerm = "")}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
            <input type="text" placeholder="Search collections..." bind:value={searchTerm} />
        </div>
    </header>

    <hr class="m-t-5 m-b-xs" />

    <div class="sidebar-content">
        {#each filteredCollections as collection (collection.id)}
            <div
                tabindex="0"
                class="sidebar-list-item"
                class:active={$activeCollection?.id === collection.id}
                on:click={() => selectCollection(collection)}
            >
                {#if $activeCollection?.id === collection.id}
                    <i class="ri-folder-open-line" />
                {:else}
                    <i class="ri-folder-2-line" />
                {/if}
                <span class="txt">{collection.name}</span>
            </div>
        {:else}
            {#if normalizedSearch.length}
                <p class="txt-hint m-t-10 m-b-10 txt-center">No collections found.</p>
            {/if}
        {/each}
    </div>

    {#if !$hideControls}
        <footer class="sidebar-footer">
            <button type="button" class="btn btn-block btn-outline" on:click={() => collectionPanel?.show()}>
                <i class="ri-add-line" />
                <span class="txt">New collection</span>
            </button>
        </footer>
    {/if}
</aside>

<CollectionUpsertPanel bind:this={collectionPanel} />
