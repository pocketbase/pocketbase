<script>
    import tooltip from "@/actions/tooltip";
    import RuleField from "@/components/collections/RuleField.svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import { slide } from "svelte/transition";

    export let collection;

    $: fieldNames = CommonHelper.getAllCollectionIdentifiers(collection);

    $: hiddenFieldNames = collection.fields?.filter((f) => f.hidden).map((f) => f.name);

    let showFiltersInfo = false;

    let showExtraRules = collection.manageRule !== null || collection.authRule !== "";
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
                        {#each fieldNames as name}
                            {#if !hiddenFieldNames.includes(name)}
                                <code>{name}</code>
                            {/if}
                        {/each}
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        The request fields could be accessed with the special <em>@request</em> filter:
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@request.headers.*</code>
                        <code>@request.query.*</code>
                        <code>@request.body.*</code>
                        <code>@request.auth.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        You could also add constraints and query other collections using the
                        <em>@collection</em> filter:
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
        <svelte:fragment slot="afterLabel" let:isSuperuserOnly>
            {#if !isSuperuserOnly}
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `The main record fields hold the values that are going to be inserted in the database.`,
                        position: "top",
                    }}
                />
            {/if}
        </svelte:fragment>
    </RuleField>

    <RuleField label="Update rule" formKey="updateRule" {collection} bind:rule={collection.updateRule}>
        <svelte:fragment slot="afterLabel" let:isSuperuserOnly>
            {#if !isSuperuserOnly}
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `The main record fields represent the old/existing record field values.\nTo target the newly submitted ones you can use @request.body.*`,
                        position: "top",
                    }}
                />
            {/if}
        </svelte:fragment>
    </RuleField>

    <RuleField label="Delete rule" formKey="deleteRule" {collection} bind:rule={collection.deleteRule} />
{/if}

{#if collection?.type === "auth"}
    <hr />

    <button
        type="button"
        class="btn btn-sm m-b-sm {showExtraRules ? 'btn-secondary' : 'btn-hint btn-transparent'}"
        on:click={() => {
            showExtraRules = !showExtraRules;
        }}
    >
        <strong class="txt">Additional auth collection rules</strong>
        {#if showExtraRules}
            <i class="ri-arrow-up-s-line txt-sm" />
        {:else}
            <i class="ri-arrow-down-s-line txt-sm" />
        {/if}
    </button>

    {#if showExtraRules}
        <div class="block" transition:slide={{ duration: 150 }}>
            <RuleField
                label="Authentication rule"
                formKey="authRule"
                placeholder=""
                {collection}
                bind:rule={collection.authRule}
            >
                <svelte:fragment>
                    <p>
                        This rule is executed every time before authentication allowing you to restrict who
                        can authenticate.
                    </p>
                    <p>
                        For example, to allow only verified users you can set it to
                        <code>verified = true</code>.
                    </p>
                    <p>Leave it empty to allow anyone with an account to authenticate.</p>
                    <p>To disable authentication entirely you can change it to "Set superusers only".</p>
                </svelte:fragment>
            </RuleField>

            <RuleField
                label="Manage rule"
                formKey="manageRule"
                placeholder=""
                required={collection.manageRule !== null}
                {collection}
                bind:rule={collection.manageRule}
            >
                <svelte:fragment>
                    <p>
                        This rule is executed in addition to the <code>create</code> and <code>update</code> API
                        rules.
                    </p>
                    <p>
                        It enables superuser-like permissions to allow fully managing the auth record(s), eg.
                        changing the password without requiring to enter the old one, directly updating the
                        verified state or email, etc.
                    </p>
                </svelte:fragment>
            </RuleField>
        </div>
    {/if}
{/if}
