<script>
    import { tick, createEventDispatcher } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollection;
    let newCollection;

    $: isCollectionRenamed = oldCollection?.name != newCollection?.name;

    $: isNewCollectionView = newCollection?.type === "view";

    $: renamedFields =
        newCollection?.schema?.filter(
            (field) => field.id && !field.toDelete && field.originalName != field.name
        ) || [];

    $: deletedFields = newCollection?.schema?.filter((field) => field.id && field.toDelete) || [];

    $: multipleToSingleFields =
        newCollection?.schema?.filter((field) => {
            const old = oldCollection?.schema?.find((f) => f.id == field.id);
            if (!old) {
                return false;
            }
            return old.options?.maxSelect != 1 && field.options?.maxSelect == 1;
        }) || [];

    $: showChanges = !isNewCollectionView || isCollectionRenamed;

    export async function show(original, changed) {
        oldCollection = original;
        newCollection = changed;

        await tick();

        if (
            isCollectionRenamed ||
            renamedFields.length ||
            deletedFields.length ||
            multipleToSingleFields.length
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
        dispatch("confirm");
    }
</script>

<OverlayPanel bind:this={panel} class="confirm-changes-panel" popup on:hide on:show>
    <svelte:fragment slot="header">
        <h4>Confirm collection changes</h4>
    </svelte:fragment>

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
                            <strong class="txt-strikethrough txt-hint">{field.originalName}</strong>
                            <i class="ri-arrow-right-line txt-sm" />
                            <strong class="txt"> {field.name}</strong>
                        </div>
                    </li>
                {/each}

                {#each deletedFields as field}
                    <li class="txt-danger">
                        Removed field <span class="txt-bold">{field.name}</span>
                    </li>
                {/each}
            {/if}
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
