<script>
    import { LogRequest } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let logPanel;
    let item = new LogRequest();

    export function show(model) {
        item = model;

        return logPanel?.show();
    }

    export function hide() {
        return logPanel?.hide();
    }
</script>

<OverlayPanel bind:this={logPanel} class="overlay-panel-lg log-panel" on:hide on:show>
    <svelte:fragment slot="header">
        <h4>Request log</h4>
    </svelte:fragment>

    <table class="table-border">
        <tbody>
            <tr>
                <td class="min-width txt-hint txt-bold">ID</td>
                <td>{item.id}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Status</td>
                <td>
                    <span class="label" class:label-danger={item.status >= 400}>
                        {item.status}
                    </span>
                </td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Method</td>
                <td>{item.method?.toUpperCase()}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Auth</td>
                <td>{item.auth}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">URL</td>
                <td>{item.url}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Referer</td>
                <td>
                    {#if item.referer}
                        <a href={item.referer} target="_blank" rel="noopener noreferrer">
                            {item.referer}
                        </a>
                    {:else}
                        <span class="txt-hint">N/A</span>
                    {/if}
                </td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Remote IP</td>
                <td>{item.remoteIp}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">User IP</td>
                <td>{item.userIp}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">UserAgent</td>
                <td>{item.userAgent}</td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Meta</td>
                <td>
                    {#if !CommonHelper.isEmpty(item.meta)}
                        <div class="block">
                            <CodeBlock content={JSON.stringify(item.meta, null, 2)} />
                        </div>
                    {:else}
                        <span class="txt-hint">N/A</span>
                    {/if}
                </td>
            </tr>
            <tr>
                <td class="min-width txt-hint txt-bold">Created</td>
                <td><FormattedDate date={item.created} /></td>
            </tr>
        </tbody>
    </table>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Close</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
