<script>
    import { fly } from "svelte/transition";

    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import providersList from "@/providers.js";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let disabled = [];

    let panel;
    let searchTerm = "";
    let filteredProviders = [];

    $: if (searchTerm !== -1 || disabled !== -1) {
        filteredProviders = filterProviders();
    }

    export function show() {
        clearSearch();
        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function select(provider) {
        dispatch("select", provider);
        hide();
    }

    function filterProviders() {
        const search = (searchTerm || "").toLowerCase();

        return providersList.filter(
            (p) =>
                !disabled.includes(p.key) &&
                (search == "" ||
                    p.key.toLowerCase().includes(search) ||
                    p.title.toLowerCase().includes(search)),
        );
    }

    function clearSearch() {
        searchTerm = "";
    }
</script>

<OverlayPanel bind:this={panel} on:show on:hide btnClose={false}>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">Add OAuth2 provider</h4>
    </svelte:fragment>

    <Field class="searchbar m-b-sm" let:uniqueId>
        <label for={uniqueId} class="m-l-10 txt-xl">
            <i class="ri-search-line" />
        </label>
        <input id={uniqueId} type="text" placeholder="Search provider" bind:value={searchTerm} />
        {#if searchTerm != ""}
            <button
                type="button"
                class="btn btn-transparent btn-sm btn-hint p-l-xs p-r-xs m-l-10"
                transition:fly={{ duration: 150, x: 5 }}
                on:click={() => (searchTerm = "")}
            >
                <span class="txt">Clear</span>
            </button>
        {/if}
    </Field>

    <div class="grid grid-sm">
        {#each filteredProviders as provider (provider.key)}
            <div class="col-6">
                <button type="button" class="provider-card handle" on:click={() => select(provider)}>
                    <figure class="provider-logo">
                        {#if provider.logo}
                            <img
                                src="{import.meta.env.BASE_URL}images/oauth2/{provider.logo}"
                                alt="{provider.title} logo"
                            />
                        {/if}
                    </figure>
                    <div class="content">
                        <div class="title">{provider.title}</div>
                        <em class="txt-hint txt-sm m-r-auto">{provider.key}</em>
                    </div>
                </button>
            </div>
        {:else}
            <div class="flex inline-flex">
                <span class="txt-hint">No providers to select.</span>
                {#if searchTerm != ""}
                    <button type="button" class="btn btn-sm btn-secondary" on:click={clearSearch}>
                        Clear filter
                    </button>
                {/if}
            </div>
        {/each}
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide}>Cancel</button>
    </svelte:fragment>
</OverlayPanel>
