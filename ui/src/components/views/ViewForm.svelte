<script>
    import { onMount, tick } from "svelte";
    import { slide } from "svelte/transition";
    import { View } from "pocketbase";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";

    let tempValues = {};
    let editorRefs = {};
    let showFiltersInfo = false;
    let sqlInputComponent;
    let isSQLeditorLoading = false;
    let ruleInputComponent;
    let isRuleComponentLoading = false;

    export let view = new View();
    function isAdminOnly(propVal) {
        return propVal === null;
    }

    async function loadEditorComponent() {
        isSQLeditorLoading = true;
        try {
            sqlInputComponent = (await import("@/components/views/SQLeditor.svelte")).default;
        } catch (err) {
            console.warn(err);
            sqlInputComponent = null;
        }
        isSQLeditorLoading = false;
        isRuleComponentLoading = true;
        try {
            ruleInputComponent = (await import("@/components/base/FilterAutocompleteInput.svelte")).default;
        } catch (err) {
            console.warn(err);
            ruleInputComponent = null;
        }
        isRuleComponentLoading = false;
    }

    onMount(() => {
        loadEditorComponent();
    });
</script>

{#if isSQLeditorLoading}
    <div class="txt-center">
        <span class="loader" />
    </div>
{:else}
    <div>
        <Field let:uniqueId name="sql" class="form-field required m-b-0">
            <label for={uniqueId}>Query</label>
            <svelte:component this={sqlInputComponent} bind:value={view.sql} id={uniqueId} />
        </Field>
    </div>
    <hr />
{/if}
<div class="block m-b-base">
    <div class="flex">
        <p>
            All rules follow the
            <a href={import.meta.env.PB_RULES_SYNTAX_DOCS} target="_blank" rel="noopener">
                PocketBase filter syntax and operators
            </a>.
        </p>
        <span
            class="expand-handle txt-sm txt-bold txt-nowrap link-hint"
            on:click={() => (showFiltersInfo = !showFiltersInfo)}
        >
            {showFiltersInfo ? "Hide available fields" : "Show available fields"}
        </span>
    </div>

    {#if showFiltersInfo}
        <div transition:slide|local={{ duration: 150 }}>
            <div class="alert alert-warning m-0">
                <div class="content">
                    <p class="m-b-1 ">
                        The fields used in the rule should be returned from the sql query.<br /> duplicated
                        column names should be renamed to be able to use them in the rule.<br />
                    </p>
                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        The request fields could be accessed with the special <em>@request</em> filter:
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@request.method</code>
                        <code>@request.query.*</code>
                        <code>@request.data.*</code>
                        <code>@request.user.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        You could also add constraints and query other collections using the <em
                            >@collection</em
                        > filter:
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@collection.ANY_COLLECTION_NAME.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p>
                        Example rule:
                        <br />
                        <code>@request.user.id!="" && created>"2022-01-01 00:00:00"</code>
                    </p>
                </div>
            </div>
        </div>
    {/if}
</div>
{#if isRuleComponentLoading}
    <div class="txt-center">
        <span class="loader" />
    </div>
{:else}
    <div class="rule-block">
        {#if isAdminOnly(view.listRule)}
            <button
                type="button"
                class="rule-toggle-btn btn btn-circle btn-outline btn-success"
                use:tooltip={"Unlock and set custom rule"}
                on:click={async () => {
                    view.listRule = tempValues.listRule || "";
                    await tick();
                    editorRefs.listRule?.focus();
                }}
            >
                <i class="ri-lock-unlock-line" />
            </button>
        {:else}
            <button
                type="button"
                class="rule-toggle-btn btn btn-circle btn-outline"
                use:tooltip={"Lock and set to Admins only"}
                on:click={() => {
                    tempValues.listRule = view.listRule;
                    view.listRule = null;
                }}
            >
                <i class="ri-lock-line" />
            </button>
        {/if}

        <Field
            class="form-field rule-field m-0 {isAdminOnly(view.listRule) ? 'disabled' : ''}"
            name="listRule"
            let:uniqueId
        >
            <label for={uniqueId}>
                List Rule - {isAdminOnly(view.listRule) ? "Admins only" : "Custom rule"}
            </label>

            <svelte:component
                this={ruleInputComponent}
                id={uniqueId}
                bind:this={editorRefs.listRule}
                bind:value={view.listRule}
                baseCollection={(() => {
                    // filter schema fields with ':' cause its invalid rule character
                    const schema = view.schema;
                    return { schema: schema.filter((i) => !i.name.includes(":")) };
                })()}
                disabled={isAdminOnly(view.listRule)}
            />

            <div class="help-block">
                {#if isAdminOnly(view.listRule)}
                    Only admins will be able to access (unlock to change)
                {:else}
                    Leave empty to grant everyone access
                {/if}
            </div>
        </Field>
    </div>
{/if}

<style>
    .rule-block {
        display: flex;
        align-items: flex-start;
        gap: var(--xsSpacing);
    }
    .rule-toggle-btn {
        margin-top: 15px;
    }
</style>
