<script>
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let panel;
    let url = "";

    export let type;

    export function show(newUrl) {
        if (newUrl === "") {
            return;
        }

        url = newUrl;

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    $: filename = url.substring(url.lastIndexOf("/") + 1);
</script>

<OverlayPanel bind:this={panel} class="preview {type}-preview" btnClose={false} popup on:show on:hide>
    <svelte:fragment slot="header">
        <button type="button" class="overlay-close" on:click|preventDefault={hide}>
            <i class="ri-close-line" />
        </button>
    </svelte:fragment>

    {#if type === "image"}
        <img src={url} alt="Preview {url}" />
    {:else if type === "pdf"}
        <object title={filename} data={url} type="application/pdf"> PDF embed not loaded. </object>
    {/if}

    <svelte:fragment slot="footer">
        <a
            href={url}
            title="Download"
            target="_blank"
            rel="noreferrer noopener"
            class="link-hint txt-ellipsis"
        >
            <i class="ri-file-download-line" />
            {filename}
        </a>
        <div class="flex-fill" />
        <button type="button" class="btn btn-secondary" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>
