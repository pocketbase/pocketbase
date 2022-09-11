<script>
    import { onMount, tick } from "svelte";
    import { slide } from "svelte/transition";
    import { Collection } from "pocketbase";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";

    export let collection = new Collection();

    let tempValues = {};
    let showFiltersInfo = false;
    let editorRefs = {};
    let ruleInputComponent;
    let isRuleComponentLoading = false;

    // all supported collection rules in "collection_rule_prop: label" format
    const ruleProps = {
        listRule: "List Action",
        viewRule: "View Action",
        createRule: "Create Action",
        updateRule: "Update Action",
        deleteRule: "Delete Action",
    };

    function isAdminOnly(propVal) {
        return propVal === null;
    }

    async function loadEditorComponent() {
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
                    <p class="m-b-0">The following record fields are available:</p>
                    <div class="inline-flex flex-gap-5">
                        <code>id</code>
                        <code>created</code>
                        <code>updated</code>
                        {#each collection.schema as field}
                            {#if field.type === "relation" || field.type === "user"}
                                <code>{field.name}.*</code>
                            {:else}
                                <code>{field.name}</code>
                            {/if}
                        {/each}
                    </div>

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
    {#each Object.entries(ruleProps) as [prop, label] (prop)}
        <hr class="m-t-sm m-b-sm" />
        <div class="rule-block">
            {#if isAdminOnly(collection[prop])}
                <button
                    type="button"
                    class="rule-toggle-btn btn btn-circle btn-outline btn-success"
                    use:tooltip={"Unlock and set custom rule"}
                    on:click={async () => {
                        collection[prop] = tempValues[prop] || "";
                        await tick();
                        editorRefs[prop]?.focus();
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
                        tempValues[prop] = collection[prop];
                        collection[prop] = null;
                    }}
                >
                    <i class="ri-lock-line" />
                </button>
            {/if}

            <Field
                class="form-field rule-field m-0 {isAdminOnly(collection[prop]) ? 'disabled' : ''}"
                name={prop}
                let:uniqueId
            >
                <label for={uniqueId}>
                    {label} - {isAdminOnly(collection[prop]) ? "Admins only" : "Custom rule"}
                </label>

                <svelte:component
                    this={ruleInputComponent}
                    id={uniqueId}
                    bind:this={editorRefs[prop]}
                    bind:value={collection[prop]}
                    baseCollection={collection}
                    disabled={isAdminOnly(collection[prop])}
                />

                <div class="help-block">
                    {#if isAdminOnly(collection[prop])}
                        Only admins will be able to access (unlock to change)
                    {:else}
                        Leave empty to grant everyone access
                    {/if}
                </div>
            </Field>
        </div>
    {/each}
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
