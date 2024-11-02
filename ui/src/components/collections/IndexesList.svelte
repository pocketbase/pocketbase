<script>
    import { scale } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { errors, removeError } from "@/stores/errors";
    import IndexUpsertPanel from "@/components/collections/IndexUpsertPanel.svelte";

    export let collection;

    let upsertPanel;

    function pushOrReplace(oldIndex, newIndex) {
        for (let i = 0; i < collection.indexes.length; i++) {
            // replace
            if (collection.indexes[i] == oldIndex) {
                collection.indexes[i] = newIndex;
                removeError("indexes." + i);
                return;
            }
        }

        // push missing
        collection.indexes.push(newIndex);
        collection.indexes = collection.indexes;
    }
</script>

<div class="section-title">
    Unique constraints and indexes ({collection?.indexes?.length || 0})
    {#if $errors?.indexes?.message}
        <i
            class="ri-error-warning-fill txt-danger"
            transition:scale={{ duration: 150 }}
            use:tooltip={$errors?.indexes.message}
        />
    {/if}
</div>
<div class="indexes-list">
    {#each collection?.indexes || [] as rawIndex, i}
        {@const parsed = CommonHelper.parseIndex(rawIndex)}
        <button
            type="button"
            class="label link-primary {$errors.indexes?.[i]?.message ? 'label-danger' : ''}"
            use:tooltip={$errors.indexes?.[i]?.message || ""}
            on:click={() => upsertPanel?.show(rawIndex, i)}
        >
            {#if parsed.unique}
                <strong>Unique:</strong>
            {/if}
            <span class="txt">
                {parsed.columns?.map((c) => c.name).join(", ")}
            </span>
        </button>
    {/each}
    <button
        type="button"
        class="btn btn-xs btn-transparent btn-pill btn-outline"
        on:click={() => upsertPanel?.show()}
    >
        <span class="txt">+</span>
        <span class="txt">New index</span>
    </button>
</div>

<IndexUpsertPanel
    bind:this={upsertPanel}
    bind:collection
    on:remove={(e) => {
        for (let i = 0; i < collection.indexes.length; i++) {
            if (collection.indexes[i] == e.detail) {
                collection.indexes.splice(i, 1);
                removeError("indexes." + i);
                break;
            }
        }
        collection.indexes = collection.indexes;
    }}
    on:submit={(e) => {
        // clear generic error
        if ($errors.indexes?.message) {
            removeError("indexes");
        }

        pushOrReplace(e.detail.old, e.detail.new);
    }}
/>

<style lang="scss">
    .indexes-list {
        display: flex;
        flex-wrap: wrap;
        width: 100%;
        gap: 10px;
    }
    .label {
        overflow: hidden;
        min-width: 50px;
    }
</style>
