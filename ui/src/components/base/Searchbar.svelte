<script>
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher, onMount } from "svelte";
    import { fly } from "svelte/transition";

    const dispatch = createEventDispatcher();
    const uniqueId = "search_" + CommonHelper.randomString(7);

    export let value = "";
    export let placeholder = 'Search term or filter like created > "2022-01-01"...';

    // autocomplete filter component fields
    export let autocompleteCollection = null;
    export let extraAutocompleteKeys = [];

    let filterComponent;
    let isFilterComponentLoading = false;

    let searchInput;
    let tempValue = "";

    $: if (typeof value === "string") {
        tempValue = value;
    }

    function clear(focusInput = true) {
        tempValue = "";
        if (focusInput) {
            searchInput?.focus();
        }
        dispatch("clear");
    }

    function submit() {
        value = tempValue;
        dispatch("submit", value);
    }

    async function loadFilterComponent() {
        if (filterComponent || isFilterComponentLoading) {
            return; // already loaded or in the process
        }

        isFilterComponentLoading = true;

        filterComponent = (await import("@/components/base/FilterAutocompleteInput.svelte")).default;

        isFilterComponentLoading = false;
    }

    onMount(() => {
        loadFilterComponent();
    });
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<form class="searchbar" on:click|stopPropagation on:submit|preventDefault={submit}>
    <label for={uniqueId} class="m-l-10 txt-xl">
        <i class="ri-search-line" />
    </label>

    {#if filterComponent && !isFilterComponentLoading}
        <svelte:component
            this={filterComponent}
            id={uniqueId}
            singleLine
            disableRequestKeys
            disableCollectionJoinKeys
            {extraAutocompleteKeys}
            baseCollection={autocompleteCollection}
            placeholder={value || placeholder}
            bind:value={tempValue}
            on:submit={submit}
        />
    {:else}
        <input
            bind:this={searchInput}
            type="text"
            id={uniqueId}
            placeholder={value || placeholder}
            bind:value={tempValue}
        />
    {/if}

    {#if (value.length || tempValue.length) && tempValue != value}
        <button
            type="submit"
            class="btn btn-expanded-sm btn-sm btn-warning"
            transition:fly={{ duration: 150, x: 5 }}
        >
            <span class="txt">Search</span>
        </button>
    {/if}

    {#if value.length || tempValue.length}
        <button
            type="button"
            class="btn btn-transparent btn-sm btn-hint p-l-xs p-r-xs m-l-10"
            transition:fly={{ duration: 150, x: 5 }}
            on:click={() => {
                clear(false);
                submit();
            }}
        >
            <span class="txt">Clear</span>
        </button>
    {/if}
</form>
