<script>
    import { Record } from "pocketbase";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import RecordFieldValue from "./RecordFieldValue.svelte";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";

    export let collection;

    let recordPanel;
    let record = new Record();

    $: hasEditorField = !!collection?.schema?.find((f) => f.type === "editor");

    export function show(model) {
        record = model;

        return recordPanel?.show();
    }

    export function hide() {
        return recordPanel?.hide();
    }
</script>

<OverlayPanel
    bind:this={recordPanel}
    class="record-preview-panel {hasEditorField ? 'overlay-panel-xl' : 'overlay-panel-lg'}"
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4><strong>{collection?.name}</strong> record preview</h4>
    </svelte:fragment>

    <table class="table-border preview-table">
        <tbody>
            <tr>
                <td class="min-width txt-hint txt-bold">id</td>
                <td class="col-field">
                    <div class="label">
                        <CopyIcon value={record.id} />
                        <span class="txt">{record.id}</span>
                    </div>
                </td>
            </tr>

            {#each collection?.schema as field}
                <tr>
                    <td class="min-width txt-hint txt-bold">{field.name}</td>
                    <td class="col-field">
                        <RecordFieldValue {field} {record} />
                    </td>
                </tr>
            {/each}

            {#if record.created}
                <tr>
                    <td class="min-width txt-hint txt-bold">created</td>
                    <td class="col-field"><FormattedDate date={record.created} /></td>
                </tr>
            {/if}

            {#if record.updated}
                <tr>
                    <td class="min-width txt-hint txt-bold">updated</td>
                    <td class="col-field"><FormattedDate date={record.updated} /></td>
                </tr>
            {/if}
        </tbody>
    </table>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Close</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style lang="scss">
    .col-field {
        max-width: 1px; // text overflow workaround
    }
</style>
