<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Field from "@/components/base/Field.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let func = {};
    let inputData = "{}";
    let outputData = "";
    let isExecuting = false;
    let hasExecuted = false;

    export function show(functionToExecute) {
        func = functionToExecute;
        inputData = "{}";
        outputData = "";
        hasExecuted = false;
        panel?.show();
    }

    export function hide() {
        panel?.hide();
    }

    async function execute() {
        if (isExecuting) {
            return;
        }

        isExecuting = true;

        try {
            let input = {};
            
            if (inputData.trim()) {
                try {
                    input = JSON.parse(inputData);
                } catch (err) {
                    throw new Error("Invalid JSON input: " + err.message);
                }
            }

            const result = await ApiClient.send(`/api/lambdas/${func.id}/execute`, {
                method: "POST",
                body: { input },
            });

            outputData = JSON.stringify(result, null, 2);
            hasExecuted = true;
            addSuccessToast("Function executed successfully");
        } catch (err) {
            outputData = JSON.stringify({
                error: err.message || "Execution failed",
                timestamp: new Date().toISOString()
            }, null, 2);
            hasExecuted = true;
            addErrorToast(err);
        }

        isExecuting = false;
    }

    function resetExecution() {
        outputData = "";
        hasExecuted = false;
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="overlay-panel-lg execute-panel"
    beforeHide={() => !isExecuting}
    on:hide
>
    <header class="panel-header">
        <h4>Execute: {func?.name}</h4>
    </header>

    <div class="panel-content">
        <div class="grid">
            <div class="col-12">
                <div class="content txt-hint m-b-sm">
                    Test your lambda function by providing input data and viewing the execution result.
                </div>
            </div>

            <!-- Function Info -->
            <div class="col-12">
                <div class="field-group">
                    <label>Function Details</label>
                </div>
                <div class="list">
                    <div class="list-item">
                        <strong>Name:</strong> {func?.name}
                    </div>
                    <div class="list-item">
                        <strong>Status:</strong>
                        <span class="label" class:label-success={func?.enabled} class:label-hint={!func?.enabled}>
                            {func?.enabled ? "Enabled" : "Disabled"}
                        </span>
                    </div>
                    <div class="list-item">
                        <strong>Timeout:</strong> {func?.timeout}s
                    </div>
                </div>
            </div>

            <!-- Input Data -->
            <div class="col-12">
                <Field class="form-field" name="input" let:uniqueId>
                    <label for={uniqueId}>Input Data (JSON)</label>
                    <textarea
                        id={uniqueId}
                        bind:value={inputData}
                        placeholder='&#123;"key": "value"&#125;'
                        rows="6"
                        disabled={isExecuting}
                    />
                    <div class="help-text">
                        Provide input data as JSON. This will be available as <code>$input</code> in your function.
                    </div>
                </Field>
            </div>

            <!-- Execute Button -->
            <div class="col-12">
                <div class="flex flex-gap-10">
                    <button
                        type="button"
                        class="btn btn-expanded"
                        class:btn-loading={isExecuting}
                        disabled={isExecuting || !func?.enabled}
                        on:click={execute}
                    >
                        <i class="ri-play-circle-line" />
                        <span class="txt">{isExecuting ? "Executing..." : "Execute Function"}</span>
                    </button>
                    
                    {#if hasExecuted}
                        <button
                            type="button"
                            class="btn btn-outline"
                            on:click={resetExecution}
                        >
                            <i class="ri-refresh-line" />
                            <span class="txt">Reset</span>
                        </button>
                    {/if}
                </div>
                
                {#if !func?.enabled}
                    <div class="help-text txt-warning">
                        ⚠️ Function is disabled and cannot be executed.
                    </div>
                {/if}
            </div>

            <!-- Output Data -->
            {#if hasExecuted}
                <div class="col-12">
                    <Field class="form-field" name="output" let:uniqueId>
                        <label for={uniqueId}>Execution Result</label>
                        <textarea
                            id={uniqueId}
                            bind:value={outputData}
                            readonly
                            rows="10"
                            class="font-mono"
                        />
                    </Field>
                </div>
            {/if}
        </div>
    </div>

    <footer class="panel-footer">
        <button type="button" class="btn btn-transparent" disabled={isExecuting} on:click={hide}>
            Close
        </button>
    </footer>
</OverlayPanel>

<style>
    :global(.execute-panel .list) {
        background: var(--baseAlt1Color);
        border-radius: var(--baseRadius);
        padding: 15px;
    }
    
    :global(.execute-panel .list-item) {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 5px 0;
    }
    
    :global(.execute-panel .list-item:not(:last-child)) {
        border-bottom: 1px solid var(--baseAlt2Color);
        margin-bottom: 5px;
        padding-bottom: 10px;
    }
    
    .field-group {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 10px;
        font-weight: 600;
    }
    
    .font-mono {
        font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
    }
</style>