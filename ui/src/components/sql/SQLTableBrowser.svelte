<script>
    import { createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import Accordion from "@/components/base/Accordion.svelte";

    const dispatch = createEventDispatcher();

    export let tables = [];
    export let isLoading = false;

    // Separate regular tables from system tables (starting with _)
    // Also exclude tables containing "sqlite" in their name
    $: regularTables = tables.filter(
        (t) => !t.name.startsWith("_") && !t.name.toLowerCase().includes("sqlite")
    );
    $: systemTables = tables.filter(
        (t) => t.name.startsWith("_") && !t.name.toLowerCase().includes("sqlite")
    );

    function handleRefresh() {
        dispatch("refresh");
    }

    function handleQuickQuery(tableName) {
        dispatch("quickQuery", { tableName });
    }
</script>

<PageSidebar>
    <header class="sidebar-header">
        <div class="sidebar-title flex flex-gap-10">
            <span class="txt">Database Schema</span>
            <button
                type="button"
                class="btn btn-xs btn-circle btn-transparent"
                on:click={handleRefresh}
                disabled={isLoading}
                aria-label="Refresh schema"
                use:tooltip={{ text: "Refresh schema", position: "right" }}
            >
                <i class="ri-refresh-line" class:rotating={isLoading} />
            </button>
        </div>
    </header>

    <hr class="m-t-5 m-b-xs" />

    <div class="sidebar-content">
        {#if isLoading}
            <div class="txt-center p-base txt-hint">
                <span class="loader loader-sm" />
            </div>
        {:else if tables.length === 0}
            <div class="txt-center p-base txt-hint">No tables found</div>
        {:else}
            <!-- Regular Tables -->
            {#if regularTables.length > 0}
                {#each regularTables as table}
                    <Accordion single class="table-accordion">
                        <svelte:fragment slot="header">
                            <div class="flex flex-gap-5">
                                <i
                                    class={table.type === "view"
                                        ? "ri-eye-line"
                                        : "ri-table-line"}
                                />
                                <span class="txt">{table.name}</span>
                                <div class="flex-fill" />
                                <button
                                    type="button"
                                    class="btn btn-xs btn-circle btn-secondary btn-outline"
                                    on:click|stopPropagation={() =>
                                        handleQuickQuery(table.name)}
                                    aria-label="Quick query"
                                    use:tooltip={{
                                        text: "SELECT * LIMIT 100",
                                        position: "left",
                                    }}
                                >
                                    <i class="ri-play-line" />
                                </button>
                            </div>
                        </svelte:fragment>

                        <div class="columns-list">
                            {#each table.columns as column}
                                <div class="column-item" title={column.type}>
                                    <i
                                        class="ri-key-line txt-sm"
                                        class:txt-primary={column.primaryKey}
                                        class:txt-disabled={!column.primaryKey}
                                    />
                                    <span class="column-name">{column.name}</span>
                                    <span class="column-type txt-hint txt-sm"
                                        >{column.type}</span
                                    >
                                </div>
                            {/each}
                        </div>
                    </Accordion>
                {/each}
            {/if}

            <!-- System Tables Section -->
            {#if systemTables.length > 0}
                <Accordion single class="system-tables-accordion">
                    <svelte:fragment slot="header">
                        <div class="system-tables-header">
                            <i class="ri-shield-line" />
                            <span class="txt">System Tables ({systemTables.length})</span>
                        </div>
                    </svelte:fragment>

                    <div class="system-tables-content">
                        {#each systemTables as table}
                            <Accordion single class="table-accordion nested-accordion">
                                <svelte:fragment slot="header">
                                    <div class="flex flex-gap-5">
                                        <i
                                            class={table.type === "view"
                                                ? "ri-eye-line"
                                                : "ri-table-line"}
                                        />
                                        <span class="txt">{table.name}</span>
                                        <div class="flex-fill" />
                                        <button
                                            type="button"
                                            class="btn btn-xs btn-circle btn-secondary btn-outline"
                                            on:click|stopPropagation={() =>
                                                handleQuickQuery(table.name)}
                                            aria-label="Quick query"
                                            use:tooltip={{
                                                text: "SELECT * LIMIT 100",
                                                position: "left",
                                            }}
                                        >
                                            <i class="ri-play-line" />
                                        </button>
                                    </div>
                                </svelte:fragment>

                                <div class="columns-list">
                                    {#each table.columns as column}
                                        <div class="column-item" title={column.type}>
                                            <i
                                                class="ri-key-line txt-sm"
                                                class:txt-primary={column.primaryKey}
                                                class:txt-disabled={!column.primaryKey}
                                            />
                                            <span class="column-name">{column.name}</span>
                                            <span class="column-type txt-hint txt-sm"
                                                >{column.type}</span
                                            >
                                        </div>
                                    {/each}
                                </div>
                            </Accordion>
                        {/each}
                    </div>
                </Accordion>
            {/if}
        {/if}
    </div>
</PageSidebar>

<style>
    :global(.table-accordion) {
        margin: 0 !important;
        border-radius: 0 !important;
        border: none !important;
        border-bottom: 1px solid var(--baseAlt2Color) !important;
    }

    :global(.table-accordion:last-child) {
        border-bottom: none !important;
    }

    :global(.system-tables-accordion) {
        margin: 0 !important;
        border-radius: 0 !important;
        border: none !important;
        border-top: 2px solid var(--baseAlt3Color) !important;
        background: var(--baseAlt2Color) !important;
    }

    .system-tables-header {
        display: flex;
        align-items: center;
        gap: 8px;
        font-weight: 600;
        opacity: 0.8;
    }

    .system-tables-content {
        background: var(--baseAlt1Color);
    }

    :global(.nested-accordion) {
        background: var(--baseAlt1Color) !important;
    }

    .columns-list {
        padding: 5px 10px 10px;
        font-size: 0.875rem;
    }

    .column-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 4px 0;
    }

    .column-name {
        font-weight: 500;
    }

    .column-type {
        margin-left: auto;
        font-size: 0.8125rem;
    }

    .rotating {
        animation: rotate 1s linear infinite;
    }

    @keyframes rotate {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }
</style>

