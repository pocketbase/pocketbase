<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addInfoToast } from "@/stores/toasts";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    $pageTitle = "Export collections";

    const uniqueId = "export_" + CommonHelper.randomString(5);

    let collections = [];
    let isLoadingCollections = false;

    $: schema = JSON.stringify(collections, null, 2);

    loadCollections();

    async function loadCollections() {
        isLoadingCollections = true;

        try {
            collections = await ApiClient.collections.getFullList(100, {
                $cancelKey: uniqueId,
            });
            // delete timestamps
            for (let collection of collections) {
                delete collection.created;
                delete collection.updated;
            }
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingCollections = false;
    }

    function download() {
        CommonHelper.downloadJson(collections, "pb_schema");
    }

    function copy() {
        CommonHelper.copyToClipboard(schema);
        addInfoToast("The schema was copied to your clipboard!", 3000);
    }
</script>

<SettingsSidebar />

<main class="page-wrapper">
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>

    <div class="wrapper">
        <div class="panel">
            {#if isLoadingCollections}
                <div class="loader" />
            {:else}
                <div class="content txt-xl m-b-base">
                    <p>
                        Below you'll find your current collections schema that you could import later in
                        another PocketBase environment.
                    </p>
                </div>

                <div class="export-preview">
                    <button
                        type="button"
                        class="btn btn-sm btn-secondary fade copy-schema"
                        on:click={() => copy()}
                    >
                        <span class="txt">Copy</span>
                    </button>

                    <CodeBlock content={schema} />
                </div>

                <div class="flex m-t-base">
                    <div class="flex-fill" />
                    <button type="button" class="btn btn-expanded" on:click={() => download()}>
                        <i class="ri-download-line" />
                        <span class="txt">Download as JSON</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</main>

<style>
    .export-preview {
        position: relative;
        height: 500px;
    }
    .export-preview .copy-schema {
        position: absolute;
        right: 15px;
        top: 15px;
    }
</style>
