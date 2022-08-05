<script>
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import CollectionsExportForm from "@/components/collections/CollectionsExportForm.svelte";

    let overlayPanel;

    export function show() {
        return overlayPanel?.show();
    }

    export function hide() {
        return overlayPanel?.hide();
    }
</script>

<OverlayPanel
    bind:this={overlayPanel}
    class="overlay-panel-xl collections-export-panel"
    on:hide
    on:show
    popup
    active
>
    <svelte:fragment slot="header">
        <h4>Export collections schema</h4>
    </svelte:fragment>

    <CollectionsExportForm />
</OverlayPanel>
