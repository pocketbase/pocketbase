<script>
    import { createEventDispatcher } from "svelte";
    import { slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import PageSidebar from "@/components/base/PageSidebar.svelte";

    const dispatch = createEventDispatcher();

    export let history = [];
    export let visible = true;

    function selectQuery(query) {
        dispatch("select", { query });
    }

    function clearHistory() {
        dispatch("clear");
    }
</script>

{#if visible}
    <div class="history-wrapper" transition:slide={{ duration: 200, axis: "x" }}>
        <PageSidebar>
            <header class="sidebar-header">
                <div class="sidebar-title flex">
                    <span class="txt">Query History</span>
                    <div class="flex-fill" />
                    {#if history.length > 0}
                        <button
                            type="button"
                            class="btn btn-xs btn-circle btn-transparent txt-danger"
                            on:click={clearHistory}
                            aria-label="Clear history"
                            use:tooltip={{ text: "Clear history", position: "left" }}
                        >
                            <i class="ri-delete-bin-line" />
                        </button>
                    {/if}
                </div>
            </header>

            <hr class="m-t-5 m-b-xs" />

            <div class="sidebar-content">
                {#if history.length === 0}
                    <div class="txt-center p-base txt-hint">No query history yet</div>
                {:else}
                    {#each history as query}
                        <button
                            type="button"
                            class="sidebar-list-item history-item"
                            on:click={() => selectQuery(query)}
                            title={query}
                        >
                            <div class="history-query">{query}</div>
                        </button>
                    {/each}
                {/if}
            </div>
        </PageSidebar>
    </div>
{/if}

<style>
    .history-wrapper {
        height: 100vh;
        display: flex;
        flex-direction: column;
    }

    .history-wrapper :global(.page-sidebar) {
        height: 100%;
    }

    .history-item {
        width: 100%;
        text-align: left;
        font-family: var(--monospaceFontFamily);
        font-size: 0.8125rem;
        padding: 10px 15px;
        display: block;
        color: inherit;
    }

    .history-query {
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        line-clamp: 3;
        -webkit-box-orient: vertical;
        word-break: break-all;
    }
</style>

