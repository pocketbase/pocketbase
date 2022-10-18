<script>
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { setErrors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { Record } from "pocketbase";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let view;

    let recordPanel;
    let record;

    export function show(model) {
        load(model);

        return recordPanel?.show();
    }

    export function hide() {
        return recordPanel?.hide();
    }

    function load(model) {
        setErrors({}); // reset errors
        record = model?.clone ? model.clone() : new Record();
    }
</script>

<OverlayPanel
    bind:this={recordPanel}
    class="overlay-panel-lg record-panel"
    beforeHide={() => {
        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4>
            {view.name} record
        </h4>
    </svelte:fragment>

    <div class="block">
        {#each view?.schema || [] as field (field.name)}
            <Field class="form-field disabled" name={field.name} let:uniqueId>
                <label for={uniqueId}>
                    <i class={CommonHelper.getFieldTypeIcon(field.type)} />
                    <span class="txt">{field.name}</span>
                    <span class="flex-fill" />
                </label>
                {#if field.type === "text" || field.type === "number"}
                    <input type={field.type} value={record[field.name]} disabled />
                {:else if field.type === "bool"}
                    <input type="checkbox" checked={record[field.name]} disabled />
                {:else if field.type === "json"}
                    <textarea class="txt-mono" value={record[field.name]} disabled />
                {/if}
            </Field>
        {/each}
    </div>
</OverlayPanel>
