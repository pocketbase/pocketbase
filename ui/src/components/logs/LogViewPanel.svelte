<script>
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
    import LogLevel from "@/components/logs/LogLevel.svelte";
    import LogDate from "@/components/logs/LogDate.svelte";

    let logPanel;
    let log = {};

    $: hasData = !CommonHelper.isEmpty(log.data);

    export function show(model) {
        if (CommonHelper.isEmpty(model)) {
            return;
        }

        log = model;

        return logPanel?.show();
    }

    export function hide() {
        return logPanel?.hide();
    }

    const priotizedKeys = [
        "execTime",
        "type",
        "auth",
        "status",
        "method",
        "url",
        "referer",
        "remoteIp",
        "userIp",
        "error",
        "details",
        //
    ];

    function extractKeys(data) {
        if (!data) {
            return [];
        }

        let keys = [];

        for (let key of priotizedKeys) {
            if (typeof data[key] !== "undefined") {
                keys.push(key);
            }
        }

        // append the rest
        const original = Object.keys(data);
        for (let key of original) {
            if (!keys.includes(key)) {
                keys.push(key);
            }
        }

        return keys;
    }

    function downloadJson() {
        CommonHelper.downloadJson(log, "log_" + log.created.replaceAll(/[-:\. ]/gi, "") + ".json");
    }
</script>

<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<OverlayPanel bind:this={logPanel} class="overlay-panel-lg log-panel" on:hide on:show>
    <svelte:fragment slot="header">
        <h4>Request log</h4>
    </svelte:fragment>

    <table class="table-border">
        <tbody>
            <tr>
                <td class="min-width txt-hint txt-bold">id</td>
                <td>
                    <div class="label">
                        <CopyIcon value={log.id} />
                        <div class="txt">{log.id}</div>
                    </div>
                </td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">level</td>
                <td><LogLevel level={log.level} /></td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">created</td>
                <td><LogDate date={log.created} /></td>
            </tr>
            {#if log.data?.type != "request"}
                <tr>
                    <td class="min-width txt-hint txt-bold">message</td>
                    <td>
                        {#if log.message}
                            <span class="txt">{log.message}</span>
                        {:else}
                            <span class="txt txt-hint">N/A</span>
                        {/if}
                    </td>
                </tr>
            {/if}
            {#each extractKeys(log.data) as key}
                {@const value = log.data[key]}
                <tr>
                    <td class="min-width txt-hint txt-bold" class:v-align-top={hasData}>
                        data.{key}
                    </td>
                    <td>
                        {#if value !== null && typeof value == "object"}
                            <CodeBlock content={JSON.stringify(value, null, 2)} />
                        {:else if CommonHelper.isEmpty(value)}
                            <span class="txt txt-hint">N/A</span>
                        {:else}
                            <span class="txt">
                                {value}{key == "execTime" ? "ms" : ""}
                            </span>
                        {/if}
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Close</span>
        </button>

        <button type="button" class="btn btn-primary" on:click={() => downloadJson()}>
            <i class="ri-download-line" />
            <span class="txt">Download as JSON</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
