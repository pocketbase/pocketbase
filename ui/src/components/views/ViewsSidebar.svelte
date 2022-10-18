<script>
    import { hideControls } from "@/stores/app";
    import { views, activeView } from "@/stores/views";
    import ViewUpsertPanel from "@/components/views/ViewUpsertPanel.svelte";

    let viewPanel;
    let searchTerm = "";

    $: normalizedSearch = searchTerm.replace(/\s+/g, "").toLowerCase();

    $: hasSearch = searchTerm !== "";

    $: filteredViews = $views.filter((view) => {
        return (
            view.id == searchTerm || view.name.replace(/\s+/g, "").toLowerCase().includes(normalizedSearch)
        );
    });

    function selectView(view) {
        $activeView = view;
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
            <input type="text" placeholder="Search views..." bind:value={searchTerm} />
        </div>
    </header>

    <hr class="m-t-5 m-b-xs" />

    <div class="sidebar-content">
        {#each filteredViews as view (view.id)}
            <div
                tabindex="0"
                class="sidebar-list-item"
                class:active={$activeView?.id === view.id}
                on:click={() => selectView(view)}
            >
                {#if $activeView?.id === view.id}
                    <i class="ri-terminal-box-fill" />
                {:else}
                    <!-- <i class="ri-folder-2-line" /> -->
                    <i class="ri-terminal-box-line" />
                {/if}
                <span class="txt">{view.name}</span>
            </div>
        {:else}
            {#if normalizedSearch.length}
                <p class="txt-hint m-t-10 m-b-10 txt-center">No views found.</p>
            {/if}
        {/each}
    </div>

    {#if !$hideControls}
        <footer class="sidebar-footer">
            <button type="button" class="btn btn-block btn-outline" on:click={() => viewPanel?.show()}>
                <i class="ri-add-line" />
                <span class="txt">New view</span>
            </button>
        </footer>
    {/if}
</aside>

<ViewUpsertPanel bind:this={viewPanel} />
