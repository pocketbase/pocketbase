<script>
    import CommonHelper from "@/utils/CommonHelper";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    let panel;
    let url = "";
    let urlOrFactory;

    $: queryParamsIndex = url.indexOf("?");

    $: filename = url.substring(
        url.lastIndexOf("/") + 1,
        queryParamsIndex > 0 ? queryParamsIndex : undefined,
    );

    $: type = CommonHelper.getFileType(filename);

    export async function show(urlOrFactoryArg) {
        urlOrFactory = urlOrFactoryArg;
        if (!urlOrFactory) {
            return;
        }

        url = await resolveUrlOrFactory();

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    async function resolveUrlOrFactory() {
        if (typeof urlOrFactory == "function") {
            return await urlOrFactory();
        }

        // string or Promise
        return await urlOrFactory;
    }

    async function openInNewTab() {
        try {
            // resolve again because it may have expired
            url = await resolveUrlOrFactory();
            window.open(url, "_blank", "noreferrer,noopener");
        } catch (err) {
            if (!err.isAbort) {
                console.warn("openInNewTab file token failure:", err);
            }
        }
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
        <button
            type="button"
            title={filename}
            class="link-hint txt-ellipsis inline-flex"
            on:auxclick={openInNewTab}
            on:click={openInNewTab}
        >
            {filename}
            <i class="ri-external-link-line" />
        </button>
        <div class="flex-fill" />
        <button type="button" class="btn btn-transparent" on:click={hide}>Close</button>
    </svelte:fragment>
</OverlayPanel>
