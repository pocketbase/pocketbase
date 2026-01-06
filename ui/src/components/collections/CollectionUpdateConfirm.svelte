<script>
    import { createEventDispatcher, tick } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollection;
    let newCollection;
    let hideAfterSave;
    let conflictingOIDCs = [];
    let changedRules = [];

    $: isCollectionRenamed = oldCollection?.name != newCollection?.name;

    $: isNewCollectionView = newCollection?.type === "view";

    $: isNewCollectionAuth = newCollection?.type === "auth";

    $: renamedFields =
        (!isNewCollectionView &&
            newCollection?.fields?.filter(
                (field) => field.id && !field._toDelete && field._originalName != field.name,
            )) ||
        [];

    $: deletedFields =
        (!isNewCollectionView && newCollection?.fields?.filter((field) => field.id && field._toDelete)) || [];

    $: multipleToSingleFields =
        newCollection?.fields?.filter((field) => {
            const old = oldCollection?.fields?.find((f) => f.id == field.id);
            if (!old) {
                return false;
            }
            return old.maxSelect > 1 && field.maxSelect <= 1;
        }) || [];

    $: showChanges = !isNewCollectionView || isCollectionRenamed || changedRules.length;

    export async function show(original, changed, hideAfterSaveArg = true) {
        oldCollection = original;
        newCollection = changed;
        hideAfterSave = hideAfterSaveArg;

        await detectConflictingOIDCs();

        detectRulesChange();

        await tick();

        if (
            isCollectionRenamed ||
            renamedFields.length ||
            deletedFields.length ||
            multipleToSingleFields.length ||
            conflictingOIDCs.length ||
            changedRules.length
        ) {
            panel?.show();
        } else {
            // no changes to review -> confirm directly
            confirm();
        }
    }

    export function hide() {
        panel?.hide();
    }

    function confirm() {
        hide();
        dispatch("confirm", hideAfterSave);
    }

    const oidcProviders = ["oidc", "oidc2", "oidc3"];

    async function detectConflictingOIDCs() {
        conflictingOIDCs = [];

        for (let name of oidcProviders) {
            let oldProvider = oldCollection?.oauth2?.providers?.find((p) => p.name == name);
            let newProvider = newCollection?.oauth2?.providers?.find((p) => p.name == name);

            if (!oldProvider || !newProvider) {
                continue;
            }

            let oldHost = new URL(oldProvider.authURL).host;
            let newHost = new URL(newProvider.authURL).host;
            if (oldHost == newHost) {
                continue;
            }

            // check if there are existing externalAuths
            if (await haveExternalAuths(name)) {
                conflictingOIDCs.push({ name, oldHost, newHost });
            }
        }
    }

    async function haveExternalAuths(provider) {
        try {
            await ApiClient.collection("_externalAuths").getFirstListItem(
                ApiClient.filter("collectionRef={:collectionId} && provider={:provider}", {
                    collectionId: newCollection?.id,
                    provider: provider,
                }),
            );
            return true;
        } catch {}

        return false;
    }

    function getExternalAuthsFilterLink(provider) {
        return `#/collections?collection=_externalAuths&filter=collectionRef%3D%22${newCollection?.id}%22+%26%26+provider%3D%22${provider}%22`;
    }

    function detectRulesChange() {
        changedRules = [];

        // for now enable only for "production"
        if (window.location.protocol != "https:") {
            return;
        }

        const ruleProps = ["listRule", "viewRule"];
        if (!isNewCollectionView) {
            ruleProps.push("createRule", "updateRule", "deleteRule");
        }
        if (isNewCollectionAuth) {
            ruleProps.push("manageRule", "authRule");
        }

        let oldRule, newRule;
        for (let prop of ruleProps) {
            oldRule = oldCollection?.[prop];
            newRule = newCollection?.[prop];
            if (oldRule === newRule) {
                continue;
            }

            changedRules.push({ prop, oldRule, newRule });
        }
    }
</script>

<OverlayPanel bind:this={panel} class="confirm-changes-panel" popup on:hide on:show>
    <svelte:fragment slot="header">
        <h4>Confirm collection changes</h4>
    </svelte:fragment>

    {#if isCollectionRenamed || deletedFields.length || renamedFields.length}
        <div class="alert alert-warning">
            <div class="icon">
                <i class="ri-error-warning-line" />
            </div>
            <div class="content txt-bold">
                <p>
                    If any of the collection changes is part of another collection rule, filter or view query,
                    you'll have to update it manually!
                </p>
                {#if deletedFields.length}
                    <p>All data associated with the removed fields will be permanently deleted!</p>
                {/if}
            </div>
        </div>
    {/if}

    {#if showChanges}
        <h6>Changes:</h6>
        <ul class="changes-list">
            {#if isCollectionRenamed}
                <li>
                    <div class="inline-flex">
                        Renamed collection
                        <strong class="txt-strikethrough txt-hint">{oldCollection?.name}</strong>
                        <i class="ri-arrow-right-line txt-sm" />
                        <strong class="txt"> {newCollection?.name}</strong>
                    </div>
                </li>
            {/if}

            {#if !isNewCollectionView}
                {#each multipleToSingleFields as field}
                    <li>
                        Multiple to single value conversion of field
                        <strong>{field.name}</strong>
                        <em class="txt-sm">(will keep only the last array item)</em>
                    </li>
                {/each}

                {#each renamedFields as field}
                    <li>
                        <div class="inline-flex">
                            Renamed field
                            <strong class="txt-strikethrough txt-hint">{field._originalName}</strong>
                            <i class="ri-arrow-right-line txt-sm" />
                            <strong class="txt">{field.name}</strong>
                        </div>
                    </li>
                {/each}

                {#each deletedFields as field}
                    <li class="txt-danger">
                        Removed field <span class="txt-bold">{field.name}</span>
                    </li>
                {/each}
            {/if}

            {#each changedRules as ruleChange}
                <li>
                    Changed API rule <code class="txt-sm">{ruleChange.prop}</code>:
                    <br />
                    <small class="txt-mono txt-hint">
                        <strong>Old</strong>:
                        <span class="txt-preline">
                            {ruleChange.oldRule === null
                                ? "null (superusers only)"
                                : ruleChange.oldRule || '""'}
                        </span>
                    </small>
                    <br />
                    <small class="txt-mono txt-success">
                        <strong>New</strong>:
                        <span class="txt-preline">
                            {ruleChange.newRule === null
                                ? "null (superusers only)"
                                : ruleChange.newRule || '""'}
                        </span>
                    </small>
                </li>
            {/each}

            {#each conflictingOIDCs as oidc}
                <li>
                    Changed <code>{oidc.name}</code> host
                    <div class="inline-flex m-l-5">
                        <strong class="txt-strikethrough txt-hint">{oidc.oldHost}</strong>
                        <i class="ri-arrow-right-line txt-sm" />
                        <strong class="txt">{oidc.newHost}</strong>
                    </div>
                    <br />
                    <em class="txt-hint">
                        If the old and new OIDC configuration is not for the same provider consider deleting
                        all old <code class="txt-sm">_externalAuths</code> records associated to the current
                        collection and provider, otherwise it may result in account linking errors.
                        <a href={getExternalAuthsFilterLink(oidc.name)} target="_blank">
                            Review existing <code class="txt-sm">_externalAuths</code> records
                            <i class="ri-external-link-line txt-sm"></i>
                        </a>.
                    </em>
                </li>
            {/each}
        </ul>
    {/if}

    <svelte:fragment slot="footer">
        <!-- svelte-ignore a11y-autofocus -->
        <button autofocus type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button type="button" class="btn btn-expanded" on:click={() => confirm()}>
            <span class="txt">Confirm</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style lang="scss">
    .changes-list {
        word-break: break-word;
        line-height: var(--smLineHeight);
        li {
            margin-top: 10px;
            margin-bottom: 10px;
        }
    }
</style>
