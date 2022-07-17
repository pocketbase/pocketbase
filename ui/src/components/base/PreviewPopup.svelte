<script>
    import CommonHelper from "@/utils/CommonHelper";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let panel;
    let url = "";

    export function show(newUrl) {
        if (newUrl === "") {
            return;
        }

        CommonHelper.checkImageUrl(newUrl)
            .then(() => {
                url = newUrl;
                panel?.show();
            })
            .catch(() => {
                console.warn("Invalid image preview url: ", newUrl);
                hide();
            });
    }

    export function hide() {
        return panel?.hide();
    }
</script>

<OverlayPanel bind:this={panel} class="image-preview" btnClose={false} popup on:show on:hide>
    <svelte:fragment slot="header">
        <div class="overlay-close" on:click|preventDefault={hide}>
            <i class="ri-close-line" />
        </div>
    </svelte:fragment>

    <img src={url} alt="Preview" />

    <svelte:fragment slot="footer">
        <a href={url} title="Download" class="link-hint txt-ellipsis">
            /../{url.substring(url.lastIndexOf("/") + 1)}
        </a>
        <div class="flex-fill" />
        <button type="button" class="btn btn-secondary" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>
