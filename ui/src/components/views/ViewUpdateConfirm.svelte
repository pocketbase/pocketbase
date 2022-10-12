<script>
    import { tick, createEventDispatcher } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let view;

    $: isViewRenamed = view?.originalName != view?.name;
    $: isSqlChange = view?.originalSql != view?.sql;
    $: isListRuleChange = view?.originalListRule != view?.listRule;

    export async function show(viewToCheck) {
        view = viewToCheck;

        await tick();

        if (!isViewRenamed && !isSqlChange && !isListRuleChange) {
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
        <h4>Confirm view changes</h4>
    </svelte:fragment>

    <h6>Changes:</h6>
    <ul class="changes-list">
        {#if isSqlChange}
            <li>
                <div class="inline-flex">
                    <strong class="txt">SQL Query change</strong>
                </div>
            </li>
        {/if}
    </ul>
    <ul class="changes-list">
        {#if isViewRenamed}
            <li>
                <div class="inline-flex">
                    Renamed view
                    <strong class="txt-strikethrough txt-hint">{view.originalName}</strong>
                    <i class="ri-arrow-right-line txt-sm" />
                    <strong class="txt"> {view.name}</strong>
                </div>
            </li>
        {/if}
    </ul>
    <ul class="changes-list">
        {#if isListRuleChange}
            <li>
                <div class="inline-flex">
                    List rule changes
                    <code class="txt-strikethrough txt-hint txt-sm">{view.originalListRule}</code>
                    <i class="ri-arrow-right-line txt-sm" />
                    <code class="txt-sm"> {view.listRule}</code>
                </div>
            </li>
        {/if}
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
