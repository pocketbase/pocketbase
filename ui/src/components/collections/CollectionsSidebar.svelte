<script>
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import CollectionSidebarItem from "@/components/collections/CollectionSidebarItem.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import { hideControls } from "@/stores/app";
    import { activeCollection, collections, isCollectionsLoading } from "@/stores/collections";

    const pinnedStorageKey = "@pinnedCollections";

    let collectionPanel;
    let searchTerm = "";
    let pinnedIds = [];
    let showSystemSection = false;
    let oldCollectionId;

    loadPinned();

    $: if ($collections) {
        syncPinned();
        scrollIntoView();
    }

    $: normalizedSearch = searchTerm.replace(/\s+/g, "").toLowerCase();

    $: hasSearch = searchTerm !== "";

    $: if (pinnedIds) {
        localStorage.setItem(pinnedStorageKey, JSON.stringify(pinnedIds));
    }

    $: filtered = $collections.filter((c) => {
        return c.id == searchTerm || c.name?.replace(/\s+/g, "")?.toLowerCase()?.includes(normalizedSearch);
    });

    $: pinnedCollections = filtered.filter((c) => pinnedIds.includes(c.id));

    $: unpinnedRegularCollections = filtered.filter((c) => !c.system && !pinnedIds.includes(c.id));

    $: unpinnedSystemCollections = filtered.filter((c) => c.system && !pinnedIds.includes(c.id));

    $: if ($activeCollection?.id && oldCollectionId != $activeCollection.id) {
        oldCollectionId = $activeCollection.id;
        if ($activeCollection.system && !pinnedCollections.find((c) => c.id == $activeCollection.id)) {
            showSystemSection = true;
        } else {
            showSystemSection = false;
        }
    }

    function scrollIntoView() {
        setTimeout(() => {
            const activeItem = document.querySelector(".collection-sidebar .sidebar-list-item.active");
            if (activeItem) {
                activeItem?.scrollIntoView({ block: "nearest" });
            }
        }, 0);
    }

    function loadPinned() {
        pinnedIds = [];

        try {
            const encoded = localStorage.getItem(pinnedStorageKey);
            if (encoded) {
                pinnedIds = JSON.parse(encoded) || [];
            }
        } catch (_) {}
    }

    function syncPinned() {
        pinnedIds = pinnedIds.filter((id) => !!$collections.find((c) => c.id == id));
    }
</script>

<PageSidebar class="collection-sidebar">
    <header class="sidebar-header">
        <div class="form-field search" class:active={hasSearch}>
            <div class="form-field-addon">
                <button
                    type="button"
                    class="btn btn-xs btn-transparent btn-circle btn-clear"
                    class:hidden={!hasSearch}
                    on:click={() => (searchTerm = "")}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
            <input
                type="text"
                placeholder="Search collections..."
                name="collections-search"
                bind:value={searchTerm}
            />
        </div>
    </header>

    <hr class="m-t-5 m-b-xs" />

    <div
        class="sidebar-content"
        class:fade={$isCollectionsLoading}
        class:sidebar-content-compact={filtered.length > 20}
    >
        {#if pinnedCollections.length}
            <div class="sidebar-title">Pinned</div>
            {#each pinnedCollections as collection (collection.id)}
                <CollectionSidebarItem {collection} bind:pinnedIds />
            {/each}
        {/if}

        {#if unpinnedRegularCollections.length}
            {#if pinnedCollections.length}
                <div class="sidebar-title">Others</div>
            {/if}
            {#each unpinnedRegularCollections as collection (collection.id)}
                <CollectionSidebarItem {collection} bind:pinnedIds />
            {/each}
        {/if}

        {#if unpinnedSystemCollections.length}
            <button
                type="button"
                class="sidebar-title m-b-xs"
                class:link-hint={!normalizedSearch.length}
                aria-label={showSystemSection ? "Expand system collections" : "Collapse system collections"}
                aria-expanded={showSystemSection || normalizedSearch.length}
                disabled={normalizedSearch.length}
                on:click={() => {
                    if (!normalizedSearch.length) {
                        showSystemSection = !showSystemSection;
                    }
                }}
            >
                <span class="txt">System</span>
                {#if !normalizedSearch.length}
                    <i class="ri-arrow-{showSystemSection ? 'up' : 'down'}-s-line" aria-hidden="true" />
                {/if}
            </button>
            {#if showSystemSection || normalizedSearch.length}
                {#each unpinnedSystemCollections as collection (collection.id)}
                    <CollectionSidebarItem {collection} bind:pinnedIds />
                {/each}
            {/if}
        {/if}

        {#if normalizedSearch.length && !filtered.length}
            <p class="txt-hint m-t-10 m-b-10 txt-center">No collections found.</p>
        {/if}
    </div>

    {#if !$hideControls}
        <footer class="sidebar-footer">
            <button type="button" class="btn btn-block btn-outline" on:click={() => collectionPanel?.show()}>
                <i class="ri-add-line" />
                <span class="txt">New collection</span>
            </button>
        </footer>
    {/if}
</PageSidebar>

<CollectionUpsertPanel bind:this={collectionPanel} />
