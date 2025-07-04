<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import LambdaFunctionUpsertPanel from "@/components/lambdas/LambdaFunctionUpsertPanel.svelte";
    import LambdaFunctionExecutePanel from "@/components/lambdas/LambdaFunctionExecutePanel.svelte";
    import LambdaFunctionLogsPanel from "@/components/lambdas/LambdaFunctionLogsPanel.svelte";

    $pageTitle = "Lambda Functions";

    let functions = [];
    let isLoading = false;
    let upsertPanel;
    let executePanel;
    let logsPanel;

    export const queryParams = {};

    async function loadFunctions() {
        isLoading = true;

        try {
            console.log("Loading lambda functions...");
            const response = await ApiClient.send("/api/lambdas", {
                method: "GET",
            });
            console.log("Lambda functions response:", response);
            functions = response || [];
        } catch (err) {
            console.error("Error loading lambda functions:", err);
            if (!err?.isAbort) {
                addErrorToast(err);
            }
        }

        isLoading = false;
    }

    function create() {
        upsertPanel?.show();
    }

    async function update(func) {
        try {
            // Fetch the full function data including code
            const fullFunc = await ApiClient.send(`/api/lambdas/${func.id}`, {
                method: "GET",
            });
            upsertPanel?.show(fullFunc);
        } catch (err) {
            addErrorToast(err);
        }
    }

    function execute(func) {
        executePanel?.show(func);
    }

    function viewLogs(func) {
        logsPanel?.show(func);
    }

    async function deleteConfirm(func) {
        return confirm(`Are you sure you want to delete "${func.name}"?`, {
            title: "Delete lambda function",
            yesText: "Delete",
        });
    }

    async function deleteFunction(func) {
        if (!(await deleteConfirm(func))) {
            return;
        }

        try {
            await ApiClient.send(`/api/lambdas/${func.id}`, {
                method: "DELETE",
            });

            addSuccessToast(`Lambda function "${func.name}" was deleted.`);
            loadFunctions();
        } catch (err) {
            addErrorToast(err);
        }
    }

    async function toggleFunction(func) {
        try {
            await ApiClient.send(`/api/lambdas/${func.id}`, {
                method: "PATCH",
                body: {
                    enabled: !func.enabled
                }
            });

            const action = func.enabled ? "disabled" : "enabled";
            addSuccessToast(`Lambda function "${func.name}" was ${action}.`);
            loadFunctions();
        } catch (err) {
            addErrorToast(err);
        }
    }

    loadFunctions();

</script>

<LambdaFunctionUpsertPanel bind:this={upsertPanel} on:save={loadFunctions} />
<LambdaFunctionExecutePanel bind:this={executePanel} />
<LambdaFunctionLogsPanel bind:this={logsPanel} />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>

    <div class="wrapper">
        <div class="panel">
            <div class="flex m-b-sm flex-gap-10">
                <span class="txt-xl">Lambda Functions</span>
                <RefreshButton class="btn-sm" on:refresh={loadFunctions} />
                <button
                    type="button"
                    class="btn btn-sm btn-outline"
                    on:click={() => create()}
                >
                    <i class="ri-add-line" />
                    <span class="txt">New function</span>
                </button>
            </div>

            {#if isLoading}
                <div class="block txt-center">
                    <span class="loader" />
                </div>
            {:else if !functions.length}
                <div class="block txt-center txt-hint">
                    <h6>No lambda functions found.</h6>
                    <button
                        type="button"
                        class="btn btn-sm btn-outline"
                        on:click={() => create()}
                    >
                        <i class="ri-add-line" />
                        <span class="txt">Create your first function</span>
                    </button>
                </div>
            {:else}
                <div class="table-wrapper">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Status</th>
                                <th>Timeout</th>
                                <th>Created</th>
                                <th>Updated</th>
                                <th class="min-width" />
                            </tr>
                        </thead>
                        <tbody>
                            {#each functions as func (func.id)}
                                <tr class="row-handle">
                                    <td>
                                        <div class="flex flex-gap-10">
                                            <strong class="txt">{func.name}</strong>
                                            {#if func.description}
                                                <span class="txt-hint">- {func.description}</span>
                                            {/if}
                                        </div>
                                    </td>
                                    <td>
                                        <span
                                            class="label"
                                            class:label-success={func.enabled}
                                            class:label-hint={!func.enabled}
                                        >
                                            {func.enabled ? "Enabled" : "Disabled"}
                                        </span>
                                    </td>
                                    <td>
                                        <span class="txt">{func.timeout}s</span>
                                    </td>
                                    <td>
                                        <span class="txt-hint">
                                            {CommonHelper.formatToLocalDate(func.created)}
                                        </span>
                                    </td>
                                    <td>
                                        <span class="txt-hint">
                                            {CommonHelper.formatToLocalDate(func.updated)}
                                        </span>
                                    </td>
                                    <td class="min-width nowrap">
                                        <div class="inline-flex">
                                            <button
                                                type="button"
                                                class="btn btn-xs btn-circle btn-hint"
                                                title="Execute function"
                                                on:click={() => execute(func)}
                                            >
                                                <i class="ri-play-circle-line" />
                                            </button>
                                            <button
                                                type="button"
                                                class="btn btn-xs btn-circle btn-hint"
                                                title="View logs"
                                                on:click={() => viewLogs(func)}
                                            >
                                                <i class="ri-file-list-3-line" />
                                            </button>
                                            <button
                                                type="button"
                                                class="btn btn-xs btn-circle btn-hint"
                                                title={func.enabled ? "Disable" : "Enable"}
                                                on:click={() => toggleFunction(func)}
                                            >
                                                <i class="ri-{func.enabled ? 'pause' : 'play'}-circle-line" />
                                            </button>
                                            <button
                                                type="button"
                                                class="btn btn-xs btn-circle btn-hint"
                                                title="Edit function"
                                                on:click={() => update(func)}
                                            >
                                                <i class="ri-edit-line" />
                                            </button>
                                            <button
                                                type="button"
                                                class="btn btn-xs btn-circle btn-hint"
                                                title="Delete function"
                                                on:click={() => deleteFunction(func)}
                                            >
                                                <i class="ri-delete-bin-7-line" />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            {/if}
        </div>
    </div>
</PageWrapper>