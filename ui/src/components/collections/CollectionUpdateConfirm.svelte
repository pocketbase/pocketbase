<script>
    import { tick, createEventDispatcher } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let collection;

    $: isCollectionRenamed = collection?.originalName != collection?.name;

    $: renamedFields =
        collection?.schema.filter(
            (field) => field.id && !field.toDelete && field.originalName != field.name
        ) || [];

    $: deletedFields = collection?.schema.filter((field) => field.id && field.toDelete) || [];

    export async function show(collectionToCheck) {
        collection = collectionToCheck;

        await tick();

        if (!isCollectionRenamed && !renamedFields.length && !deletedFields.length) {
            // no confirm required changes
            confirm();
        } else {
            panel?.show();
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
                If any of the following changes is part of another collection rule or filter, you'll have to
                update it manually!
            </p>
            {#if deletedFields.length}
                <p>All data associated with the removed fields will be permanently deleted!</p>
            {/if}
        </div>
    </div>

    <h6>Changes:</h6>
    <ul class="changes-list">
        {#if isCollectionRenamed}
            <li>
                <div class="inline-flex">
                    Renamed collection
                    <strong class="txt-strikethrough txt-hint">{collection.originalName}</strong>
                    <i class="ri-arrow-right-line txt-sm" />
                    <strong class="txt"> {collection.name}</strong>
                </div>
            </li>
        {/if}

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
    </ul>

    <svelte:fragment slot="footer">
        <!-- svelte-ignore a11y-autofocus -->
        <button autofocus type="button" class="btn btn-secondary" on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button type="button" class="btn btn-expanded" on:click={() => confirm()}>
            <span class="txt">Confirm</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style>
    .changes-list {
        word-break: break-all;
    }
</style>
