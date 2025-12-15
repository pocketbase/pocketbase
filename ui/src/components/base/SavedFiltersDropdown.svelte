<script>
    import { createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";
    import { savedFilters, addSavedFilter, removeSavedFilter } from "@/stores/savedFilters";
    import Toggler from "@/components/base/Toggler.svelte";

    export let collectionId = "";
    export let currentFilter = "";

    const dispatch = createEventDispatcher();

    let showSaveForm = false;
    let newFilterName = "";
    let toggler;
    let nameInput;

    $: collectionFilters = $savedFilters[collectionId] || [];

    function selectFilter(filter) {
        dispatch("select", filter.filter);
        toggler?.hide();
    }

    function openSaveForm() {
        showSaveForm = true;
        // Focus input after DOM update
        setTimeout(() => nameInput?.focus(), 50);
    }

    function closeSaveForm() {
        showSaveForm = false;
        newFilterName = "";
    }

    function saveCurrentFilter() {
        if (!newFilterName.trim() || !currentFilter.trim()) return;
        addSavedFilter(collectionId, newFilterName.trim(), currentFilter);
        closeSaveForm();
    }

    function deleteFilter(e, filter) {
        e.stopPropagation();
        removeSavedFilter(collectionId, filter.id);
    }

    function handleKeydown(e) {
        if (e.key === "Escape") {
            closeSaveForm();
        }
    }
</script>

<Toggler bind:this={toggler} class="dropdown dropdown-right dropdown-nowrap saved-filters-dropdown">
    <button
        type="button"
        slot="trigger"
        class="btn btn-transparent btn-circle"
        class:btn-hint={!collectionFilters.length}
        use:tooltip={{ text: "Saved filters", position: "left" }}
    >
        <i class="ri-bookmark-line" />
    </button>

    <div class="dropdown-content">
        <div class="txt-hint txt-sm p-5 m-b-5">Saved Filters</div>

        {#if collectionFilters.length === 0 && !showSaveForm}
            <div class="txt-hint txt-center p-10 txt-sm">No saved filters yet</div>
        {/if}

        {#each collectionFilters as filter (filter.id)}
            <button
                type="button"
                class="dropdown-item closable saved-filter-item"
                on:click={() => selectFilter(filter)}
            >
                <span class="txt txt-ellipsis">{filter.name}</span>
                <span class="txt-hint txt-ellipsis filter-preview" title={filter.filter}>
                    {filter.filter}
                </span>
                <button
                    type="button"
                    class="btn btn-transparent btn-circle btn-hint btn-xs delete-btn"
                    on:click={(e) => deleteFilter(e, filter)}
                    use:tooltip={{ text: "Delete", position: "left" }}
                >
                    <i class="ri-close-line" />
                </button>
            </button>
        {/each}

        {#if currentFilter && !showSaveForm}
            <hr class="dropdown-divider m-t-5 m-b-5" />
            <button type="button" class="dropdown-item" on:click={openSaveForm}>
                <i class="ri-add-line" />
                <span class="txt">Save current filter</span>
            </button>
        {/if}

        {#if showSaveForm}
            <hr class="dropdown-divider m-t-5 m-b-5" />
            <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
            <form
                class="save-filter-form p-10"
                on:submit|preventDefault={saveCurrentFilter}
                on:keydown={handleKeydown}
            >
                <label class="txt-hint txt-sm m-b-5" for="filter-name-input">Filter name</label>
                <input
                    bind:this={nameInput}
                    type="text"
                    id="filter-name-input"
                    class="input-sm"
                    placeholder="e.g., Active users"
                    bind:value={newFilterName}
                />
                <div class="txt-hint txt-xs m-t-5 filter-preview-text">
                    <strong>Filter:</strong> {currentFilter}
                </div>
                <div class="flex gap-5 m-t-10">
                    <button type="button" class="btn btn-sm btn-transparent" on:click={closeSaveForm}>
                        Cancel
                    </button>
                    <button type="submit" class="btn btn-sm" disabled={!newFilterName.trim()}>
                        Save
                    </button>
                </div>
            </form>
        {/if}
    </div>
</Toggler>

<style>
    :global(.saved-filters-dropdown) {
        min-width: 280px;
    }
    .saved-filter-item {
        display: flex;
        align-items: center;
        gap: 10px;
        padding-right: 5px;
    }
    .saved-filter-item .txt {
        flex-shrink: 0;
        max-width: 120px;
    }
    .saved-filter-item .filter-preview {
        flex: 1;
        min-width: 0;
        text-align: right;
        font-size: var(--smFontSize);
    }
    .saved-filter-item .delete-btn {
        flex-shrink: 0;
        opacity: 0;
        transition: opacity var(--baseAnimationSpeed);
    }
    .saved-filter-item:hover .delete-btn,
    .saved-filter-item:focus-within .delete-btn {
        opacity: 1;
    }
    .save-filter-form {
        min-width: 220px;
    }
    .filter-preview-text {
        word-break: break-all;
        max-height: 60px;
        overflow: hidden;
        text-overflow: ellipsis;
    }
</style>
