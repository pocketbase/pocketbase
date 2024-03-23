<script>
    import { link } from "svelte-spa-router";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { activeCollection } from "@/stores/collections";

    export let collection;
    export let pinnedIds;

    $: isPinned = pinnedIds.includes(collection.id);

    function toggleCollectionPin(collection) {
        if (pinnedIds.includes(collection.id)) {
            CommonHelper.removeByValue(pinnedIds, collection.id);
        } else {
            pinnedIds.push(collection.id);
        }

        pinnedIds = pinnedIds;
    }
</script>

<a
    href="/collections?collectionId={collection.id}"
    class="sidebar-list-item"
    title={collection.name}
    class:active={$activeCollection?.id === collection.id}
    use:link
>
    <i class={CommonHelper.getCollectionTypeIcon(collection.type)} aria-hidden="true" />
    <span class="txt m-r-auto">{collection.name}</span>
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <span
        class="btn btn-xs btn-circle btn-hint btn-transparent pin-collection"
        aria-label={"Pin collection"}
        aria-hidden="true"
        use:tooltip={{ position: "right", text: (isPinned ? "Unpin" : "Pin") + " collection" }}
        on:click|preventDefault|stopPropagation={() => toggleCollectionPin(collection)}
    >
        {#if isPinned}
            <i class="ri-unpin-line" />
        {:else}
            <i class="ri-pushpin-line m-l-auto" />
        {/if}
    </span>
</a>

<style lang="scss">
    .pin-collection {
        margin: 0 -5px 0 -15px;
        opacity: 0;
        transition: opacity var(--baseAnimationSpeed);
        i {
            font-size: inherit;
        }
    }
    a:hover .pin-collection {
        opacity: 0.4;
        &:hover {
            opacity: 1;
        }
    }
</style>
