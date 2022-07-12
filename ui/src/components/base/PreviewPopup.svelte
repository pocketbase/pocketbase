<script>
    import CommonHelper from "@/utils/CommonHelper";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let panel;
    let url = "";
    let fileType = false;

    export function show(newUrl) {
        if (newUrl === "") {
            return;
        }

        fileType = CommonHelper.getFileType(newUrl);

        if(fileType === "image") {
            CommonHelper.checkImageUrl(newUrl)
                .then(() => {
                    url = newUrl;
                    panel?.show();
                })
                .catch(() => {
                    console.warn("Invalid image preview url: ", newUrl);
                    hide();
                });
        } else if(fileType === "video") {
            CommonHelper.checkVideoUrl(newUrl)
                .then(() => {
                    url = newUrl;
                    panel?.show();
                })
                .catch(() => {
                    console.warn("Invalid video preview url: ", newUrl);
                    hide();
                });
        }
    }

    export function hide() {
        return panel?.hide();
    }
</script>

<OverlayPanel bind:this={panel} class="image-preview" popup on:show on:hide>
    {#if fileType === "image"}
        <img src={url} alt="Preview" />
    {:else if fileType === "video"}
        <video controls width="100%" alt="Preview"><source src={url}><track kind="captions"></video>
    {/if}

    <svelte:fragment slot="footer">
        <a href={url} class="link-hint txt-ellipsis">/../{url.substring(url.lastIndexOf("/") + 1)}</a>
        <div class="flex-fill" />
        <button type="button" class="btn btn-secondary" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>
