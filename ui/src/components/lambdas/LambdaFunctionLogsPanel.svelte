<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addErrorToast } from "@/stores/toasts";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let func = {};
    let logs = [];
    let isLoading = false;

    export function show(functionToShow) {
        func = functionToShow;
        logs = [];
        loadLogs();
        panel?.show();
    }

    export function hide() {
        panel?.hide();
    }

    async function loadLogs() {
        if (!func?.id) return;

        isLoading = true;

        try {
            const result = await ApiClient.send(`/api/lambdas/${func.id}/logs`, {
                method: "GET",
            });

            logs = result.logs || [];
        } catch (err) {
            if (!err?.isAbort) {
                addErrorToast(err);
            }
        }

        isLoading = false;
    }

    function getStatusClass(success) {
        return success ? "label-success" : "label-danger";
    }

    function formatDuration(durationMs) {
        if (durationMs < 1000) {
            return `${durationMs}ms`;
        }
        return `${(durationMs / 1000).toFixed(2)}s`;
    }

    function formatOutput(output) {
        if (typeof output === "string") {
            return output;
        }
        return JSON.stringify(output, null, 2);
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="overlay-panel-xl logs-panel"
    on:hide
>
    <header class="panel-header">
        <h4>Logs: {func?.name}</h4>
        <div class="flex flex-gap-10">
            <RefreshButton class="btn-sm" on:refresh={loadLogs} />
        </div>
    </header>

    <div class="panel-content">
        <div class="grid">
            <div class="col-12">
                <div class="content txt-hint m-b-sm">
                    Execution logs and history for this lambda function.
                </div>
            </div>

            {#if isLoading}
                <div class="col-12">
                    <div class="block txt-center">
                        <span class="loader" />
                    </div>
                </div>
            {:else if !logs.length}
                <div class="col-12">
                    <div class="block txt-center txt-hint">
                        <h6>No execution logs found.</h6>
                        <p>Logs will appear here after the function is executed.</p>
                    </div>
                </div>
            {:else}
                <div class="col-12">
                    <div class="table-wrapper">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>Status</th>
                                    <th>Trigger</th>
                                    <th>Duration</th>
                                    <th>Timestamp</th>
                                    <th>Output/Error</th>
                                </tr>
                            </thead>
                            <tbody>
                                {#each logs as log (log.id)}
                                    <tr>
                                        <td>
                                            <span class="label {getStatusClass(log.success)}">
                                                {log.success ? "Success" : "Error"}
                                            </span>
                                        </td>
                                        <td>
                                            <span class="txt">{log.trigger_type}</span>
                                        </td>
                                        <td>
                                            <span class="txt">{formatDuration(log.duration_ms)}</span>
                                        </td>
                                        <td>
                                            <span class="txt-hint">
                                                {CommonHelper.formatToLocalDate(log.timestamp)}
                                            </span>
                                        </td>
                                        <td>
                                            <div class="log-content">
                                                {#if log.success}
                                                    {#if log.output}
                                                        <pre class="output-content">{formatOutput(log.output)}</pre>
                                                    {:else}
                                                        <span class="txt-hint">No output</span>
                                                    {/if}
                                                {:else}
                                                    <pre class="error-content">{log.error}</pre>
                                                {/if}
                                            </div>
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                </div>
            {/if}
        </div>
    </div>

    <footer class="panel-footer">
        <button type="button" class="btn btn-transparent" on:click={hide}>
            Close
        </button>
    </footer>
</OverlayPanel>

<style>
    :global(.logs-panel .log-content) {
        max-width: 400px;
    }
    
    :global(.logs-panel .output-content),
    :global(.logs-panel .error-content) {
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        padding: 8px;
        margin: 0;
        font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
        font-size: 12px;
        line-height: 1.4;
        max-height: 120px;
        overflow-y: auto;
        white-space: pre-wrap;
        word-break: break-word;
    }
    
    :global(.logs-panel .error-content) {
        background: var(--dangerAltColor);
        color: var(--dangerColor);
    }
    
    :global(.logs-panel .table td) {
        vertical-align: top;
    }
</style>