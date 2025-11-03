<script>
    export let results = null;
    export let error = null;
    export let isWrite = false;
    export let rowsAffected = 0;
    export let executionMs = 0;

    function formatValue(value) {
        if (value === null || value === undefined) {
            return "NULL";
        }
        if (typeof value === "string" && value.length > 100) {
            return value.substring(0, 100) + "...";
        }
        return String(value);
    }

    function isNullValue(value) {
        return value === null || value === undefined;
    }
</script>

{#if error}
    <div class="alert alert-danger m-b-base">
        <div class="icon">
            <i class="ri-error-warning-line" />
        </div>
        <div class="content">
            <p><strong>Error:</strong></p>
            <p>{error}</p>
        </div>
    </div>
{/if}

{#if results !== null && results.length === 0 && !isWrite}
    <div class="alert alert-info m-b-base">
        <div class="icon">
            <i class="ri-information-line" />
        </div>
        <div class="content">
            <p>Query executed successfully but returned no results.</p>
            <p class="txt-hint">Execution time: {executionMs}ms</p>
        </div>
    </div>
{/if}

{#if results !== null && results.length > 0}
    <div class="wrapper">
        <div class="flex m-b-sm">
            <h5>Results</h5>
            <div class="flex-fill" />
            <span class="txt-hint">
                {results.length} row(s) • {executionMs}ms
            </span>
        </div>

        <div class="table-wrapper">
            <table class="table" style="table-layout: auto;">
                <thead>
                    <tr>
                        {#each Object.keys(results[0] || {}) as column}
                            <th>{column}</th>
                        {/each}
                    </tr>
                </thead>
                <tbody>
                    {#each results as row}
                        <tr>
                            {#each Object.keys(results[0] || {}) as column}
                                <td
                                    class:txt-hint={isNullValue(row[column])}
                                    title={String(row[column])}
                                >
                                    {formatValue(row[column])}
                                </td>
                            {/each}
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
{/if}

{#if isWrite && !error}
    <div class="alert alert-success m-b-base">
        <div class="icon">
            <i class="ri-checkbox-circle-line" />
        </div>
        <div class="content">
            <p>Query executed successfully. {rowsAffected} row(s) affected.</p>
            <p class="txt-hint">Execution time: {executionMs}ms</p>
        </div>
    </div>
{/if}

<style>
    .wrapper {
        width: 100%;
    }

    .table-wrapper {
        width: 100%;
        overflow-x: auto;
        max-height: 600px;
        overflow-y: auto;
        border: 1px solid var(--baseAlt2Color);
        border-radius: var(--baseRadius);
    }

    .table {
        font-size: 0.875rem;
    }

    .table td {
        max-width: 300px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
</style>

