<script>
    import CommonHelper from "@/utils/CommonHelper";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let panel;
    let url = "";

    $: queryParamsIndex = url.indexOf("?");

    $: filename = url.substring(
        url.lastIndexOf("/") + 1,
        queryParamsIndex > 0 ? queryParamsIndex : undefined
    );

    $: type = CommonHelper.getFileType(filename);

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

<OverlayPanel bind:this={panel} class="preview preview-{type}" btnClose={false} popup on:show on:hide>
    <svelte:fragment slot="header">
        <button type="button" class="overlay-close" on:click|preventDefault={hide}>
            <i class="ri-close-line" />
        </button>
    </svelte:fragment>

    {#if panel?.isActive()}
        {#if type === "image"}
            <img src={url} alt="Preview {filename}" />
        {:else}
            <object title={filename} data={url}>Cannot preview the file.</object>
        {/if}
    {/if}

    <svelte:fragment slot="footer">
        <a
            href={url}
            title={filename}
            target="_blank"
            rel="noreferrer noopener"
            class="link-hint txt-ellipsis inline-flex"
        >
            {filename}
            <i class="ri-external-link-line" />
        </a>
        <div class="flex-fill" />
        <button type="button" class="btn btn-transparent" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>
