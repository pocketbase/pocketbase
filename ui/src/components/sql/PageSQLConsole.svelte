<script>
    import { onMount } from "svelte";
    import tooltip from "@/actions/tooltip";
    import ApiClient from "@/utils/ApiClient";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Field from "@/components/base/Field.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import SQLTableBrowser from "@/components/sql/SQLTableBrowser.svelte";
    import SQLQueryHistory from "@/components/sql/SQLQueryHistory.svelte";
    import SQLQueryResults from "@/components/sql/SQLQueryResults.svelte";
    import WriteQueryConfirmation from "@/components/sql/WriteQueryConfirmation.svelte";

    $pageTitle = "SQL Console";

    const HISTORY_STORAGE_KEY = "pb_sql_console_history";
    const HISTORY_VISIBLE_KEY = "pb_sql_console_history_visible";
    const TABS_STORAGE_KEY = "pb_sql_console_tabs";
    const MAX_HISTORY = 50;

    let editorComponent;
    let queryHistory = [];
    let readOnlyMode = true;
    let showConfirmation = false;
    let tables = [];
    let isLoadingSchema = false;
    let showHistory = true;
    
    // Tab management
    let tabs = [];
    let activeTabIndex = 0;
    let nextTabId = 1;

    onMount(async () => {
        try {
            editorComponent = (await import("@/components/base/CodeEditor.svelte")).default;
        } catch (err) {
            console.warn("Failed to load CodeEditor component:", err);
        }

        loadQueryHistory();
        loadHistoryVisibility();
        loadTabs();
        loadSchema();
    });

    function createNewTab() {
        const newTab = {
            id: nextTabId++,
            name: `Query ${tabs.length + 1}`,
            query: "",
            results: null,
            error: null,
            isExecuting: false,
            executionMs: 0,
            rowsAffected: 0,
            isWrite: false,
        };
        tabs = [...tabs, newTab];
        activeTabIndex = tabs.length - 1;
        saveTabs();
    }

    function closeTab(index) {
        if (tabs.length === 1) {
            // Don't close the last tab, just reset it
            tabs[0] = {
                ...tabs[0],
                query: "",
                results: null,
                error: null,
                executionMs: 0,
                rowsAffected: 0,
                isWrite: false,
            };
        } else {
            tabs = tabs.filter((_, i) => i !== index);
            // Adjust active tab if necessary
            if (activeTabIndex >= tabs.length) {
                activeTabIndex = tabs.length - 1;
            } else if (activeTabIndex > index) {
                activeTabIndex--;
            }
        }
        saveTabs();
    }

    function selectTab(index) {
        activeTabIndex = index;
        saveTabs();
    }

    function renameTab(index, newName) {
        tabs[index].name = newName;
        tabs = tabs;
        saveTabs();
    }

    function loadTabs() {
        try {
            const stored = localStorage.getItem(TABS_STORAGE_KEY);
            if (stored) {
                const data = JSON.parse(stored);
                tabs = data.tabs || [];
                activeTabIndex = data.activeTabIndex || 0;
                nextTabId = data.nextTabId || 1;
                
                // Ensure all tabs have the required properties
                tabs = tabs.map(tab => ({
                    id: tab.id || nextTabId++,
                    name: tab.name || "Query",
                    query: tab.query || "",
                    results: null, // Don't persist results
                    error: null,
                    isExecuting: false,
                    executionMs: 0,
                    rowsAffected: 0,
                    isWrite: false,
                }));
            }
            
            // Create default tab if none exist
            if (tabs.length === 0) {
                createNewTab();
            }
        } catch (err) {
            console.warn("Failed to load tabs:", err);
            createNewTab();
        }
    }

    function saveTabs() {
        try {
            // Only save tab structure and queries, not results
            const dataToSave = {
                tabs: tabs.map(tab => ({
                    id: tab.id,
                    name: tab.name,
                    query: tab.query,
                })),
                activeTabIndex,
                nextTabId,
            };
            localStorage.setItem(TABS_STORAGE_KEY, JSON.stringify(dataToSave));
        } catch (err) {
            console.warn("Failed to save tabs:", err);
        }
    }

    $: activeTab = tabs[activeTabIndex];

    function loadHistoryVisibility() {
        try {
            const stored = localStorage.getItem(HISTORY_VISIBLE_KEY);
            if (stored !== null) {
                showHistory = JSON.parse(stored);
            }
        } catch (err) {
            console.warn("Failed to load history visibility:", err);
        }
    }

    function toggleHistory() {
        showHistory = !showHistory;
        try {
            localStorage.setItem(HISTORY_VISIBLE_KEY, JSON.stringify(showHistory));
        } catch (err) {
            console.warn("Failed to save history visibility:", err);
        }
    }

    function loadQueryHistory() {
        try {
            const stored = localStorage.getItem(HISTORY_STORAGE_KEY);
            if (stored) {
                queryHistory = JSON.parse(stored);
            }
        } catch (err) {
            console.warn("Failed to load query history:", err);
        }
    }

    function saveQueryHistory() {
        try {
            localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(queryHistory));
        } catch (err) {
            console.warn("Failed to save query history:", err);
        }
    }

    function clearQueryHistory() {
        queryHistory = [];
        localStorage.removeItem(HISTORY_STORAGE_KEY);
    }

    async function loadSchema() {
        isLoadingSchema = true;
        try {
            tables = await ApiClient.send("/api/sql/schema", {
                method: "GET",
            });
        } catch (err) {
            console.warn("Failed to load schema:", err);
            ApiClient.error(err, false);
        }
        isLoadingSchema = false;
    }

    function handleKeydown(e) {
        // Ctrl+Enter or Cmd+Enter to execute
        if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
            e.preventDefault();
            executeQuery();
        }
        // Ctrl+T or Cmd+T to create new tab
        if ((e.ctrlKey || e.metaKey) && e.key === "t") {
            e.preventDefault();
            createNewTab();
        }
    }

    function toggleWriteMode() {
        if (readOnlyMode) {
            // Show confirmation when enabling write mode
            showConfirmation = true;
        } else {
            // Can disable write mode directly
            readOnlyMode = true;
        }
    }

    function onConfirmWrite() {
        readOnlyMode = false;
    }

    function onCancelWrite() {
        // showConfirmation will be automatically set to false by the modal
    }

    async function executeQuery() {
        if (!activeTab || !activeTab.query.trim() || activeTab.isExecuting) {
            return;
        }

        // Clear previous results and errors for this tab
        activeTab.results = null;
        activeTab.error = null;
        activeTab.executionMs = 0;
        activeTab.rowsAffected = 0;
        activeTab.isWrite = false;
        activeTab.isExecuting = true;
        tabs = tabs; // Trigger reactivity

        try {
            const response = await ApiClient.send("/api/sql", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    sql: activeTab.query,
                    allowWrite: !readOnlyMode,
                }),
            });

            // Add to history (max 50, newest first)
            const trimmedQuery = activeTab.query.trim();
            // Remove duplicates and add to front
            queryHistory = [
                trimmedQuery,
                ...queryHistory.filter((q) => q !== trimmedQuery),
            ].slice(0, MAX_HISTORY);
            saveQueryHistory();

            // Store results in the tab
            activeTab.results = response.results || [];
            activeTab.rowsAffected = response.rowsAffected || 0;
            activeTab.isWrite = response.isWrite || false;
            activeTab.executionMs = response.executionMs || 0;

            if (activeTab.isWrite) {
                addSuccessToast(
                    `Query executed successfully. ${activeTab.rowsAffected} row(s) affected.`
                );
                // Reset to read-only mode after write query
                readOnlyMode = true;
            }
        } catch (err) {
            if (err?.data?.message) {
                activeTab.error = err.data.message;
            } else if (err?.message) {
                activeTab.error = err.message;
            } else {
                activeTab.error = "An unexpected error occurred.";
            }
        }

        activeTab.isExecuting = false;
        tabs = tabs; // Trigger reactivity
        saveTabs();
    }

    function selectFromHistory(query) {
        if (activeTab) {
            activeTab.query = query;
            tabs = tabs;
            saveTabs();
        }
    }

    function clearResults() {
        if (activeTab) {
            activeTab.results = null;
            activeTab.error = null;
            activeTab.executionMs = 0;
            activeTab.rowsAffected = 0;
            activeTab.isWrite = false;
            tabs = tabs;
        }
    }

    function quickQuery(tableName) {
        if (activeTab) {
            activeTab.query = `SELECT * FROM ${tableName} LIMIT 100`;
            tabs = tabs;
            saveTabs();
            executeQuery();
        }
    }

    function updateTabQuery(query) {
        if (activeTab) {
            activeTab.query = query;
            tabs = tabs;
            saveTabs();
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

<SQLTableBrowser
    {tables}
    isLoading={isLoadingSchema}
    on:refresh={loadSchema}
    on:quickQuery={(e) => quickQuery(e.detail.tableName)}
/>

<div class="sql-main-content">
    <PageWrapper class="flex-content">
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">SQL Console</div>
            </nav>

            <RefreshButton on:refresh={loadSchema} />

            <div class="btns-group">
                <span
                    class="label"
                    class:label-success={!readOnlyMode}
                    class:label-primary={readOnlyMode}
                >
                    {readOnlyMode ? "READ-ONLY" : "WRITE ENABLED"}
                </span>
                <button
                    type="button"
                    class="btn btn-outline"
                    class:btn-warning={readOnlyMode}
                    on:click={toggleWriteMode}
                >
                    <i class="ri-lock-line" aria-hidden="true" />
                    <span class="txt"
                        >{readOnlyMode ? "Enable writes" : "Disable writes"}</span
                    >
                </button>
                <button
                    type="button"
                    class="btn btn-outline"
                    class:btn-hint={!showHistory}
                    aria-label={showHistory ? "Hide history" : "Show history"}
                    use:tooltip={{
                        text: showHistory ? "Hide history" : "Show history",
                        position: "left",
                    }}
                    on:click={toggleHistory}
                >
                    <i class="ri-history-line" />
                    <span class="txt">{showHistory ? "Hide history" : "Show history"}</span>
                </button>
            </div>
        </header>

        <!-- Tab Bar -->
        <div class="tabs-bar">
            <div class="tabs-list">
                {#each tabs as tab, index}
                    <button
                        type="button"
                        class="tab"
                        class:active={activeTabIndex === index}
                        on:click={() => selectTab(index)}
                        aria-label={`Tab: ${tab.name}`}
                    >
                        <span class="tab-name">{tab.name}</span>
                        {#if tab.isExecuting}
                            <i class="ri-loader-4-line rotating txt-hint" />
                        {/if}
                        {#if tabs.length > 1}
                            <button
                                type="button"
                                class="tab-close"
                                on:click|stopPropagation={() => closeTab(index)}
                                aria-label="Close tab"
                            >
                                <i class="ri-close-line" />
                            </button>
                        {/if}
                    </button>
                {/each}
                <button
                    type="button"
                    class="tab-new"
                    on:click={createNewTab}
                    aria-label="New tab"
                    use:tooltip={{ text: "New tab (Ctrl+T)", position: "bottom" }}
                >
                    <i class="ri-add-line" />
                </button>
            </div>
        </div>

        {#if activeTab}
            <div class="wrapper">
                <Field class="form-field m-b-base" name="sql" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">SQL Query</span>
                        <i
                            class="ri-information-line link-hint"
                            aria-label="Press Ctrl+Enter (Cmd+Enter on Mac) to execute the query"
                            use:tooltip={{ position: "right" }}
                        />
                    </label>
                    <svelte:component
                        this={editorComponent}
                        id={uniqueId}
                        value={activeTab.query}
                        on:change={(e) => updateTabQuery(e.detail)}
                        language="sql-select"
                        placeholder="Enter your SQL query here..."
                        minHeight={150}
                        maxHeight={400}
                    />
                </Field>

                <div class="flex gap-10 m-b-base">
                    <button
                        type="button"
                        class="btn btn-expanded"
                        disabled={!activeTab.query.trim() || activeTab.isExecuting}
                        on:click={executeQuery}
                    >
                        <span class="txt"
                            >{activeTab.isExecuting ? "Executing..." : "Execute Query"}</span
                        >
                        {#if !activeTab.isExecuting}
                            <span class="txt txt-hint">(Ctrl+Enter)</span>
                        {/if}
                    </button>

                    {#if activeTab.results !== null || activeTab.error !== null}
                        <button
                            type="button"
                            class="btn btn-secondary btn-outline"
                            on:click={clearResults}
                        >
                            <i class="ri-close-line" />
                            <span class="txt">Clear</span>
                        </button>
                    {/if}
                </div>

                <SQLQueryResults
                    results={activeTab.results}
                    error={activeTab.error}
                    isWrite={activeTab.isWrite}
                    rowsAffected={activeTab.rowsAffected}
                    executionMs={activeTab.executionMs}
                />
            </div>
        {/if}
    </PageWrapper>

    {#if showHistory}
        <SQLQueryHistory
            visible={showHistory}
            history={queryHistory}
            on:select={(e) => selectFromHistory(e.detail.query)}
            on:clear={clearQueryHistory}
        />
    {/if}
</div>

<style>
    .sql-main-content {
        display: flex;
        flex: 1;
        min-width: 0;
        height: 100%;
    }

    .sql-main-content :global(.page-wrapper) {
        flex: 1;
        min-width: 0;
    }

    .tabs-bar {
        width: 100%;
        border-bottom: 1px solid var(--baseAlt2Color);
        background: var(--baseAlt1Color);
        padding: 0 var(--baseSpacing);
        overflow-x: auto;
        overflow-y: hidden;
    }

    .tabs-list {
        display: flex;
        gap: 2px;
        min-width: min-content;
    }

    .wrapper {
        width: 100%;
    }

    .tab {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 10px 16px;
        background: transparent;
        border: none;
        border-top: 2px solid transparent;
        cursor: pointer;
        user-select: none;
        white-space: nowrap;
        transition: all 0.15s ease;
        position: relative;
    }

    .tab:hover {
        background: var(--baseAlt2Color);
    }

    .tab.active {
        background: var(--baseColor);
        border-top-color: var(--primaryColor);
    }

    .tab-name {
        font-size: 14px;
        color: var(--txtPrimaryColor);
    }

    .tab:not(.active) .tab-name {
        color: var(--txtHintColor);
    }

    .tab-close {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 18px;
        height: 18px;
        padding: 0;
        background: transparent;
        border: none;
        border-radius: var(--baseRadius);
        cursor: pointer;
        opacity: 0.6;
        transition: all 0.15s ease;
    }

    .tab-close:hover {
        opacity: 1;
        background: var(--dangerAlt1Color);
        color: var(--dangerColor);
    }

    .tab-close i {
        font-size: 14px;
    }

    .tab-new {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 10px 12px;
        background: transparent;
        border: none;
        cursor: pointer;
        opacity: 0.6;
        transition: all 0.15s ease;
    }

    .tab-new:hover {
        opacity: 1;
        background: var(--baseAlt2Color);
    }

    .tab-new i {
        font-size: 18px;
    }

    @keyframes rotating {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    .rotating {
        animation: rotating 1s linear infinite;
    }
</style>

<WriteQueryConfirmation
    bind:active={showConfirmation}
    on:confirm={onConfirmWrite}
    on:cancel={onCancelWrite}
/>
