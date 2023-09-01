<script>
    import { slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import CommonHelper from "@/utils/CommonHelper";
    import RuleField from "@/components/collections/RuleField.svelte";

    export let collection;

    $: fields = CommonHelper.getAllCollectionIdentifiers(collection);

    let showFiltersInfo = false;
</script>

<div class="block m-b-sm handle">
    <div class="flex txt-sm txt-hint m-b-5">
        <p>
            All rules follow the
            <a href={import.meta.env.PB_RULES_SYNTAX_DOCS} target="_blank" rel="noopener noreferrer">
                PocketBase filter syntax and operators
            </a>.
        </p>
        <button
            type="button"
            class="expand-handle txt-sm txt-bold txt-nowrap link-hint"
            on:click={() => (showFiltersInfo = !showFiltersInfo)}
        >
            {showFiltersInfo ? "Hide available fields" : "Show available fields"}
        </button>
    </div>

    {#if showFiltersInfo}
        <div transition:slide={{ duration: 150 }}>
            <div class="alert alert-warning m-0">
                <div class="content">
                    <p class="m-b-0">The following record fields are available:</p>
                    <div class="inline-flex flex-gap-5">
                        {#each fields as name}
                            <code>{name}</code>
                        {/each}
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        The request fields could be accessed with the special <em>@request</em> filter:
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@request.headers.*</code>
                        <code>@request.query.*</code>
                        <code>@request.data.*</code>
                        <code>@request.auth.*</code>
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
                        <code>@request.auth.id != "" && created > "2022-01-01 00:00:00"</code>
                    </p>
                </div>
            </div>
        </div>
    {/if}
</div>

<RuleField label="List/Search rule" formKey="listRule" {collection} bind:rule={collection.listRule} />

<RuleField label="View rule" formKey="viewRule" {collection} bind:rule={collection.viewRule} />

{#if collection?.type !== "view"}
    <RuleField label="Create rule" formKey="createRule" {collection} bind:rule={collection.createRule}>
        <svelte:fragment slot="afterLabel" let:isAdminOnly>
            {#if !isAdminOnly}
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `The Create rule is executed after a "dry save" of the submitted data, giving you access to the main record fields as in every other rule.`,
                        position: "top",
                    }}
                />
            {/if}
        </svelte:fragment>
    </RuleField>

    <RuleField label="Update rule" formKey="updateRule" {collection} bind:rule={collection.updateRule} />

    <RuleField label="Delete rule" formKey="deleteRule" {collection} bind:rule={collection.deleteRule} />
{/if}

{#if collection?.type === "auth"}
    <RuleField
        label="Manage rule"
        formKey="options.manageRule"
        placeholder=""
        required={collection.options.manageRule !== null}
        {collection}
        bind:rule={collection.options.manageRule}
    >
        <svelte:fragment>
            <p>
                This API rule gives admin-like permissions to allow fully managing the auth record(s), eg.
                changing the password without requiring to enter the old one, directly updating the verified
                state or email, etc.
            </p>
            <p>
                This rule is executed in addition to the <code>create</code> and <code>update</code> API rules.
            </p>
        </svelte:fragment>
    </RuleField>
{/if}
