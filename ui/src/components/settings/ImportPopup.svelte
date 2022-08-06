<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { addSuccessToast } from "@/stores/toasts";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollections = [];
    let newCollections = [];
    let isImporting = false;

    export function show(a, b) {
        oldCollections = a;
        newCollections = b;

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
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

    function diff(ops = [window.DIFF_INSERT, window.DIFF_DELETE, window.DIFF_EQUAL]) {
        const dmp = new diff_match_patch();
        const lines = dmp.diff_linesToChars_(
            JSON.stringify(oldCollections, null, 4),
            JSON.stringify(newCollections, null, 4)
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
            addSuccessToast("Successfully imported the provided schema.");
            dispatch("submit");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        hide();

        isImporting = false;
    }
</script>

<OverlayPanel bind:this={panel} class="full-width-popup import-popup" popup on:show on:hide>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">Side-by-side diff</h4>
    </svelte:fragment>

    <div class="grid m-b-base">
        <div class="col-6">
            <div class="section-title">Old schema</div>
            <code class="code-block">{@html diff([window.DIFF_DELETE, window.DIFF_EQUAL])}</code>
        </div>
        <div class="col-6">
            <div class="section-title">New schema</div>
            <code class="code-block">{@html diff([window.DIFF_INSERT, window.DIFF_EQUAL])}</code>
        </div>
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" on:click={hide}>Close</button>
        <button
            type="button"
            class="btn btn-expanded m-l-auto"
            class:btn-loading={isImporting}
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
