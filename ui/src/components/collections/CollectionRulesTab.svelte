<script>
    import { slide } from "svelte/transition";
    import { Collection } from "pocketbase";
    import RuleField from "@/components/collections/RuleField.svelte";

    export let collection = new Collection();

    let showFiltersInfo = false;
</script>

<div class="block m-b-base">
    <div class="flex txt-sm m-b-5">
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

<RuleField label="List/Search action" formKey="listRule" {collection} bind:rule={collection.listRule} />

<hr class="m-t-sm m-b-sm" />
<RuleField label="View action" formKey="viewRule" {collection} bind:rule={collection.viewRule} />

<hr class="m-t-sm m-b-sm" />
<RuleField label="Create action" formKey="createRule" {collection} bind:rule={collection.createRule} />

<hr class="m-t-sm m-b-sm" />
<RuleField label="Update action" formKey="updateRule" {collection} bind:rule={collection.updateRule} />

<hr class="m-t-sm m-b-sm" />
<RuleField label="Delete action" formKey="deleteRule" {collection} bind:rule={collection.deleteRule} />

{#if collection?.isAuth}
    <hr class="m-t-sm m-b-sm" />
    <RuleField
        label="Manage action"
        formKey="options.manageRule"
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
