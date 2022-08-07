<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { addSuccessToast } from "@/stores/toasts";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollections = [];
    let newCollections = [];
    let changes = [];
    let isImporting = false;

    $: if (Array.isArray(oldCollections) && Array.isArray(newCollections)) {
        loadChanges();
    }

    export function show(a, b) {
        oldCollections = a;
        newCollections = b;

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function loadChanges() {
        changes = [];

        // add deleted and modified collections
        for (const oldCollection of oldCollections) {
            const newCollection = CommonHelper.findByKey(newCollections, "id", oldCollection.id) || null;
            if (!newCollection?.id || JSON.stringify(oldCollection) != JSON.stringify(newCollection)) {
                changes.push({
                    old: oldCollection,
                    new: newCollection,
                });
            }
        }

        // add only new collections
        for (const newCollection of newCollections) {
            const oldCollection = CommonHelper.findByKey(oldCollections, "id", newCollection.id) || null;
            if (!oldCollection?.id) {
                changes.push({
                    old: oldCollection,
                    new: newCollection,
                });
            }
        }
    }

    function diffsToHtml(diffs, ops = [window.DIFF_INSERT, window.DIFF_DELETE, window.DIFF_EQUAL]) {
        const html = [];
        const pattern_amp = /&/g;
        const pattern_lt = /</g;
        const pattern_gt = />/g;
        const pattern_para = /\n/g;

        for (let i = 0; i < diffs.length; i++) {
            const op = diffs[i][0]; // operation (insert, delete, equal)

            if (!ops.includes(op)) {
                continue;
            }

            const text = diffs[i][1]
                .replace(pattern_amp, "&amp;")
                .replace(pattern_lt, "&lt;")
                .replace(pattern_gt, "&gt;")
                .replace(pattern_para, "<br>");

            switch (op) {
                case DIFF_INSERT:
                    html[i] = '<ins class="block">' + text + "</ins>";
                    break;
                case DIFF_DELETE:
                    html[i] = '<del class="block">' + text + "</del>";
                    break;
                case DIFF_EQUAL:
                    html[i] = text;
                    break;
            }
        }

        return html.join("");
    }

    function diff(obj1, obj2, ops = [window.DIFF_INSERT, window.DIFF_DELETE, window.DIFF_EQUAL]) {
        const dmp = new diff_match_patch();
        const lines = dmp.diff_linesToChars_(
            obj1 ? JSON.stringify(obj1, null, 4) : "",
            obj2 ? JSON.stringify(obj2, null, 4) : ""
        );
        const diffs = dmp.diff_main(lines.chars1, lines.chars2, false);

        dmp.diff_charsToLines_(diffs, lines.lineArray);

        return diffsToHtml(diffs, ops);
    }

    async function submitImport() {
        if (isImporting) {
            return;
        }

        isImporting = true;

        try {
            await ApiClient.collections.import(newCollections);
            addSuccessToast("Successfully imported the collections configuration.");
            dispatch("submit");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isImporting = false;

        hide();
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="full-width-popup import-popup"
    overlayClose={false}
    escClose={!isImporting}
    beforeHide={() => !isImporting}
    popup
    on:show
    on:hide
>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">Side-by-side diff</h4>
    </svelte:fragment>

    <div class="grid grid-sm m-b-sm">
        {#each changes as pair (pair.old?.id + pair.new?.id)}
            <div class="col-12">
                <div class="flex flex-gap-10">
                    {#if !pair.old?.id}
                        <span class="label label-success">New</span>
                        <strong>{pair.new?.name}</strong>
                    {:else if !pair.new?.id}
                        <span class="label label-danger">Deleted</span>
                        <strong>{pair.old?.name}</strong>
                    {:else}
                        <span class="label label-warning">Modified</span>
                        <div class="inline-flex fleg-gap-5">
                            {#if pair.old.name !== pair.new.name}
                                <strong class="txt-strikethrough txt-hint">{pair.old.name}</strong>
                                <i class="ri-arrow-right-line txt-sm" />
                            {/if}
                            <strong class="txt">{pair.new.name}</strong>
                        </div>
                    {/if}
                </div>
            </div>
            <div class="col-6 p-b-10">
                <code class="code-block">
                    {@html diff(pair.old, pair.new, [window.DIFF_DELETE, window.DIFF_EQUAL]) || "N/A"}
                </code>
            </div>
            <div class="col-6 p-b-10">
                <code class="code-block">
                    {@html diff(pair.old, pair.new, [window.DIFF_INSERT, window.DIFF_EQUAL]) || "N/A"}
                </code>
            </div>
        {/each}
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" on:click={hide} disabled={isImporting}>Close</button>
        <button
            type="button"
            class="btn btn-expanded"
            class:btn-loading={isImporting}
            disabled={isImporting}
            on:click={() => submitImport()}
        >
            <span class="txt">Confirm and import</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style>
    code {
        color: var(--txtHintColor);
        min-height: 100%;
    }
</style>
