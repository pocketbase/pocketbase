<script>
    import { createEventDispatcher, onMount } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Field from "@/components/base/Field.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let originalFunc = {};
    let func = {};
    let isLoading = false;
    let codeComponent = null;

    $: isNew = !originalFunc?.id;

    export function show(functionToEdit = {}) {
        load(functionToEdit);
        
        panel?.show();
    }

    export function hide() {
        panel?.hide();
    }

    function load(functionToEdit = {}) {
        isLoading = true;

        originalFunc = functionToEdit;
        func = {
            name: "",
            code: "// Enter your JavaScript code here\n\n",
            description: "",
            enabled: true,
            timeout: 30,
            triggers: {},
            env_vars: {},
            ...functionToEdit,
        };

        // Parse triggers and env_vars if they're strings
        if (typeof func.triggers === "string") {
            try {
                func.triggers = JSON.parse(func.triggers);
            } catch {
                func.triggers = {};
            }
        }

        if (typeof func.env_vars === "string") {
            try {
                func.env_vars = JSON.parse(func.env_vars);
            } catch {
                func.env_vars = {};
            }
        }

        isLoading = false;
    }

    async function save() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            let result;
            
            if (isNew) {
                result = await ApiClient.send("/api/lambdas", {
                    method: "POST",
                    body: func,
                });
            } else {
                result = await ApiClient.send(`/api/lambdas/${originalFunc.id}`, {
                    method: "PATCH",
                    body: func,
                });
            }

            addSuccessToast(
                isNew 
                    ? `Lambda function "${func.name}" was created.`
                    : `Lambda function "${func.name}" was updated.`
            );

            hide();
            dispatch("save", result);
        } catch (err) {
            addErrorToast(err);
        }

        isLoading = false;
    }

    // Load code editor component dynamically
    onMount(async () => {
        try {
            const module = await import("@/components/base/CodeEditor.svelte");
            codeComponent = module.default;
        } catch (err) {
            console.warn("Failed to load CodeEditor component:", err);
        }
    });

    // Trigger management
    let httpTriggers = [];
    let dbTriggers = [];
    let cronTriggers = [];

    $: {
        if (func.triggers) {
            httpTriggers = func.triggers.http || [];
            dbTriggers = func.triggers.database || [];
            cronTriggers = func.triggers.cron || [];
        }
    }

    function updateTriggers() {
        func.triggers = {
            ...(httpTriggers.length && { http: httpTriggers }),
            ...(dbTriggers.length && { database: dbTriggers }),
            ...(cronTriggers.length && { cron: cronTriggers }),
        };
    }

    function addHttpTrigger() {
        httpTriggers = [...httpTriggers, { method: "GET", path: "/" }];
        updateTriggers();
    }

    function removeHttpTrigger(index) {
        httpTriggers = httpTriggers.filter((_, i) => i !== index);
        updateTriggers();
    }

    function addDbTrigger() {
        dbTriggers = [...dbTriggers, { collection: "", event: "create" }];
        updateTriggers();
    }

    function removeDbTrigger(index) {
        dbTriggers = dbTriggers.filter((_, i) => i !== index);
        updateTriggers();
    }

    function addCronTrigger() {
        cronTriggers = [...cronTriggers, { schedule: "0 * * * *" }];
        updateTriggers();
    }

    function removeCronTrigger(index) {
        cronTriggers = cronTriggers.filter((_, i) => i !== index);
        updateTriggers();
    }

    // Environment variables management
    let envVarsList = [];

    $: {
        envVarsList = Object.entries(func.env_vars || {}).map(([key, value]) => ({ key, value }));
    }

    function updateEnvVars() {
        func.env_vars = {};
        envVarsList.forEach(({ key, value }) => {
            if (key.trim()) {
                func.env_vars[key.trim()] = value || "";
            }
        });
    }

    function addEnvVar() {
        envVarsList = [...envVarsList, { key: "", value: "" }];
    }

    function removeEnvVar(index) {
        envVarsList = envVarsList.filter((_, i) => i !== index);
        updateEnvVars();
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="overlay-panel-lg edge-function-panel"
    beforeHide={() => !isLoading}
    on:hide
>
    <header class="panel-header">
        <h4>
            {isNew ? "New" : "Edit"} lambda function
        </h4>
    </header>

    <div class="panel-content">
        {#if isLoading}
            <div class="block txt-center">
                <span class="loader" />
            </div>
        {:else}
            <div class="grid">
                <!-- Basic Information -->
                <div class="col-12">
                    <Field class="form-field required" name="name" let:uniqueId>
                        <label for={uniqueId}>Name</label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={func.name}
                            placeholder="my-function"
                        />
                    </Field>
                </div>

                <div class="col-12">
                    <Field class="form-field" name="description" let:uniqueId>
                        <label for={uniqueId}>Description</label>
                        <textarea
                            id={uniqueId}
                            bind:value={func.description}
                            placeholder="Function description..."
                        />
                    </Field>
                </div>

                <div class="col-6">
                    <Field class="form-field" name="enabled" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            bind:checked={func.enabled}
                        />
                        <label for={uniqueId}>Enabled</label>
                    </Field>
                </div>

                <div class="col-6">
                    <Field class="form-field required" name="timeout" let:uniqueId>
                        <label for={uniqueId}>Timeout (seconds)</label>
                        <input
                            type="number"
                            id={uniqueId}
                            required
                            min="1"
                            max="300"
                            bind:value={func.timeout}
                        />
                    </Field>
                </div>

                <!-- Code Editor -->
                <div class="col-12">
                    <Field class="form-field required" name="code" let:uniqueId>
                        <label for={uniqueId}>JavaScript Code</label>
                        {#if codeComponent}
                            <svelte:component
                                this={codeComponent}
                                bind:value={func.code}
                                language="javascript"
                                minHeight={300}
                                placeholder="// Enter your JavaScript code here..."
                            />
                        {:else}
                            <textarea
                                id={uniqueId}
                                required
                                bind:value={func.code}
                                placeholder="// Enter your JavaScript code here..."
                                rows="10"
                            />
                        {/if}
                    </Field>
                </div>

                <!-- HTTP Triggers -->
                <div class="col-12">
                    <label class="field-group">
                        HTTP Triggers
                        <button
                            type="button"
                            class="btn btn-xs btn-circle btn-outline"
                            title="Add HTTP trigger"
                            on:click={addHttpTrigger}
                        >
                            <i class="ri-add-line" />
                        </button>
                    </label>
                    
                    {#each httpTriggers as trigger, i}
                        <div class="form-field-group">
                            <div class="flex flex-gap-10">
                                <Field class="form-field flex-1" name="method_{i}">
                                    <select bind:value={trigger.method} on:change={updateTriggers}>
                                        <option value="GET">GET</option>
                                        <option value="POST">POST</option>
                                        <option value="PUT">PUT</option>
                                        <option value="PATCH">PATCH</option>
                                        <option value="DELETE">DELETE</option>
                                    </select>
                                </Field>
                                <Field class="form-field flex-3" name="path_{i}">
                                    <input
                                        type="text"
                                        bind:value={trigger.path}
                                        placeholder="/api/my-endpoint"
                                        on:input={updateTriggers}
                                    />
                                </Field>
                                <button
                                    type="button"
                                    class="btn btn-xs btn-circle btn-hint"
                                    title="Remove"
                                    on:click={() => removeHttpTrigger(i)}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>

                <!-- Database Triggers -->
                <div class="col-12">
                    <label class="field-group">
                        Database Triggers
                        <button
                            type="button"
                            class="btn btn-xs btn-circle btn-outline"
                            title="Add database trigger"
                            on:click={addDbTrigger}
                        >
                            <i class="ri-add-line" />
                        </button>
                    </label>
                    
                    {#each dbTriggers as trigger, i}
                        <div class="form-field-group">
                            <div class="flex flex-gap-10">
                                <Field class="form-field flex-1" name="collection_{i}">
                                    <input
                                        type="text"
                                        bind:value={trigger.collection}
                                        placeholder="collection_name"
                                        on:input={updateTriggers}
                                    />
                                </Field>
                                <Field class="form-field flex-1" name="event_{i}">
                                    <select bind:value={trigger.event} on:change={updateTriggers}>
                                        <option value="create">Create</option>
                                        <option value="update">Update</option>
                                        <option value="delete">Delete</option>
                                    </select>
                                </Field>
                                <button
                                    type="button"
                                    class="btn btn-xs btn-circle btn-hint"
                                    title="Remove"
                                    on:click={() => removeDbTrigger(i)}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>

                <!-- Cron Triggers -->
                <div class="col-12">
                    <label class="field-group">
                        Cron Triggers
                        <button
                            type="button"
                            class="btn btn-xs btn-circle btn-outline"
                            title="Add cron trigger"
                            on:click={addCronTrigger}
                        >
                            <i class="ri-add-line" />
                        </button>
                    </label>
                    
                    {#each cronTriggers as trigger, i}
                        <div class="form-field-group">
                            <div class="flex flex-gap-10">
                                <Field class="form-field flex-1" name="schedule_{i}">
                                    <input
                                        type="text"
                                        bind:value={trigger.schedule}
                                        placeholder="0 * * * *"
                                        on:input={updateTriggers}
                                    />
                                </Field>
                                <button
                                    type="button"
                                    class="btn btn-xs btn-circle btn-hint"
                                    title="Remove"
                                    on:click={() => removeCronTrigger(i)}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>

                <!-- Environment Variables -->
                <div class="col-12">
                    <label class="field-group">
                        Environment Variables
                        <button
                            type="button"
                            class="btn btn-xs btn-circle btn-outline"
                            title="Add environment variable"
                            on:click={addEnvVar}
                        >
                            <i class="ri-add-line" />
                        </button>
                    </label>
                    
                    {#each envVarsList as envVar, i}
                        <div class="form-field-group">
                            <div class="flex flex-gap-10">
                                <Field class="form-field flex-1" name="env_key_{i}">
                                    <input
                                        type="text"
                                        bind:value={envVar.key}
                                        placeholder="VARIABLE_NAME"
                                        on:input={updateEnvVars}
                                    />
                                </Field>
                                <Field class="form-field flex-1" name="env_value_{i}">
                                    <input
                                        type="text"
                                        bind:value={envVar.value}
                                        placeholder="variable_value"
                                        on:input={updateEnvVars}
                                    />
                                </Field>
                                <button
                                    type="button"
                                    class="btn btn-xs btn-circle btn-hint"
                                    title="Remove"
                                    on:click={() => removeEnvVar(i)}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </div>

    <footer class="panel-footer">
        <button type="button" class="btn btn-transparent" disabled={isLoading} on:click={hide}>
            Cancel
        </button>
        <button type="button" class="btn btn-expanded" disabled={isLoading} on:click={save}>
            <span class="txt">{isNew ? "Create" : "Save changes"}</span>
        </button>
    </footer>
</OverlayPanel>

<style>
    .field-group {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 10px;
        font-weight: 600;
    }
    
    :global(.edge-function-panel .form-field-group) {
        margin-bottom: 10px;
    }
</style>