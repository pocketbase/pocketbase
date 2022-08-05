<script>
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    export let title = "Side-by-side diff";
    export let contentATitle = "Old state";
    export let contentBTitle = "New state";

    let panel;
    let contentA = "";
    let contentB = "";

    export function show(a, b) {
        contentA = a;
        contentB = b;

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function diffsToHtml(diffs, ops = [DIFF_INSERT, DIFF_DELETE, DIFF_EQUAL]) {
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

    function diff(ops = [DIFF_INSERT, DIFF_DELETE, DIFF_EQUAL]) {
        const dmp = new diff_match_patch();
        const lines = dmp.diff_linesToChars_(contentA, contentB);
        const diffs = dmp.diff_main(lines.chars1, lines.chars2, false);

        dmp.diff_charsToLines_(diffs, lines.lineArray);

        return diffsToHtml(diffs, ops);
    }
</script>

<OverlayPanel bind:this={panel} class="full-width-popup diff-popup" popup on:show on:hide>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">{title}</h4>
    </svelte:fragment>

    <div class="grid m-b-base">
        <div class="col-6">
            <div class="section-title">{contentATitle}</div>
            <code class="code-block">{@html diff([DIFF_DELETE, DIFF_EQUAL])}</code>
        </div>
        <div class="col-6">
            <div class="section-title">{contentBTitle}</div>
            <code class="code-block">{@html diff([DIFF_INSERT, DIFF_EQUAL])}</code>
        </div>
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>

<style>
    code {
        color: var(--txtHintColor);
        min-height: 100%;
    }
</style>
