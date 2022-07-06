<script>
    import { Collection } from "pocketbase";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import ListApiDocs from "@/components/collections/docs/ListApiDocs.svelte";
    import ViewApiDocs from "@/components/collections/docs/ViewApiDocs.svelte";
    import CreateApiDocs from "@/components/collections/docs/CreateApiDocs.svelte";
    import UpdateApiDocs from "@/components/collections/docs/UpdateApiDocs.svelte";
    import DeleteApiDocs from "@/components/collections/docs/DeleteApiDocs.svelte";
    import RealtimeApiDocs from "@/components/collections/docs/RealtimeApiDocs.svelte";

    const tabs = [
        {
            id: "list",
            label: "List",
            component: ListApiDocs,
        },
        {
            id: "view",
            label: "View",
            component: ViewApiDocs,
        },
        {
            id: "create",
            label: "Create",
            component: CreateApiDocs,
        },
        {
            id: "update",
            label: "Update",
            component: UpdateApiDocs,
        },
        {
            id: "delete",
            label: "Delete",
            component: DeleteApiDocs,
        },
        {
            id: "realtime",
            label: "Realtime",
            component: RealtimeApiDocs,
        },
    ];

    let collectionPanel;
    let collection = new Collection();
    let activeTab = tabs[0].id;

    export function show(model) {
        collection = model;

        changeTab(tabs[0].id);

        return collectionPanel?.show();
    }

    export function hide() {
        return collectionPanel?.hide();
    }

    export function changeTab(newTab) {
        activeTab = newTab;
    }

    function changeTabViaKey(e, newTab) {
        if (e.code === "Enter" || e.code === "Space") {
            e.preventDefault();
            changeTab(newTab);
        }
    }
</script>

<OverlayPanel
    bind:this={collectionPanel}
    on:hide
    on:show
    class="overlay-panel-xl colored-header collection-panel"
>
    <svelte:fragment slot="header">
        <h4><strong>{collection.name}</strong> records API</h4>

        <div class="tabs-header stretched">
            {#each tabs as tab (tab.id)}
                <button
                    tabindex="0"
                    class="tab-item"
                    class:active={activeTab === tab.id}
                    on:click={() => changeTab(tab.id)}
                    on:keydown|self={(e) => changeTabViaKey(e, tab.id)}
                >
                    <span class="txt">{tab.label}</span>
                </button>
            {/each}
        </div>
    </svelte:fragment>

    <div class="tabs-content">
        {#each tabs as tab (tab.id)}
            {#if activeTab === tab.id}
                <div class="tab-item active">
                    <svelte:component this={tab.component} {collection} />
                </div>
            {/if}
        {/each}
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" on:click={() => hide()}>
            <span class="txt">Close</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
