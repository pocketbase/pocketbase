<script>
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let panel;
    let url = "";

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
</script>

<OverlayPanel bind:this={panel} class="image-preview" btnClose={false} popup on:show on:hide>
    <svelte:fragment slot="header">
        <button type="button" class="overlay-close" on:click|preventDefault={hide}>
            <i class="ri-close-line" />
        </button>
    </svelte:fragment>

    <img src={url} alt="Preview {url}" />

    <svelte:fragment slot="footer">
        <a
            href={url}
            title="Download"
            target="_blank"
            rel="noreferrer noopener"
            class="link-hint txt-ellipsis"
        >
            {url.substring(url.lastIndexOf("/") + 1)}
        </a>
        <div class="flex-fill" />
        <button type="button" class="btn btn-secondary" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>
