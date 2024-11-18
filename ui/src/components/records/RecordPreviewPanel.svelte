<script>
    import { addErrorToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import RecordFieldValue from "@/components/records/RecordFieldValue.svelte";

    export let collection;

    let recordPanel;
    let record = {};
    let isLoading = false;

    $: hasEditorField = !!collection?.fields?.find((f) => f.type === "editor");

    export function show(model) {
        load(model);

        return recordPanel?.show();
    }

    export function hide() {
        isLoading = false;
        return recordPanel?.hide();
    }

    async function load(model) {
        record = {}; // reset

        isLoading = true;

        record = (await resolveModel(model)) || {};

        isLoading = false;
    }

    async function resolveModel(model) {
        if (model && typeof model === "string") {
            // load from id
            try {
                return await ApiClient.collection(collection.id).getOne(model);
            } catch (err) {
                if (!err.isAbort) {
                    hide();
                    console.warn("resolveModel:", err);
                    addErrorToast(`Unable to load record with id "${model}"`);
                }
            }

            return null;
        }

        return model;
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

    <table class="table-border preview-table" class:table-loading={isLoading}>
        <tbody>
            {#each collection?.fields as field}
                <tr>
                    <td class="min-width txt-hint txt-bold">{field.name}</td>
                    <td class="col-field">
                        <RecordFieldValue {field} {record} />
                    </td>
                </tr>
            {/each}
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
