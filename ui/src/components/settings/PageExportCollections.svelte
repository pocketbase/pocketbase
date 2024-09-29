<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addInfoToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Field from "@/components/base/Field.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    $pageTitle = "Export collections";

    const uniqueId = "export_" + CommonHelper.randomString(5);

    let previewContainer;
    let collections = [];
    let bulkSelected = {};
    let isLoadingCollections = false;

    $: schema = JSON.stringify(Object.values(bulkSelected), null, 4);

    $: totalBulkSelected = Object.keys(bulkSelected).length;

    $: areAllSelected = collections.length && totalBulkSelected === collections.length;

    loadCollections();

    async function loadCollections() {
        isLoadingCollections = true;

        try {
            collections = await ApiClient.collections.getFullList({
                batch: 100,
                $cancelKey: uniqueId,
            });

            collections = CommonHelper.sortCollections(collections);

            for (let collection of collections) {
                // delete timestamps
                delete collection.created;
                delete collection.updated;

                // unset oauth2 providers
                delete collection.oauth2?.providers;
            }

            selectAll();
        } catch (err) {
            ApiClient.error(err);
        }

        isLoadingCollections = false;
    }

    function download() {
        CommonHelper.downloadJson(Object.values(bulkSelected), "pb_schema");
    }

    function copy() {
        CommonHelper.copyToClipboard(schema);
        addInfoToast("The configuration was copied to your clipboard!", 3000);
    }

    function toggleSelectAll() {
        if (areAllSelected) {
            deselectAll();
        } else {
            selectAll();
        }
    }

    function deselectAll() {
        bulkSelected = {};
    }

    function selectAll() {
        bulkSelected = {};

        for (const collection of collections) {
            bulkSelected[collection.id] = collection;
        }
    }

    function toggleSelectCollection(collection) {
        if (!bulkSelected[collection.id]) {
            bulkSelected[collection.id] = collection;
        } else {
            delete bulkSelected[collection.id];
        }

        bulkSelected = bulkSelected; // trigger reactivity
    }
</script>

<SettingsSidebar />

<PageWrapper>
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
                        Below you'll find your current collections configuration that you could import in
                        another PocketBase environment.
                    </p>
                </div>

                <div class="export-panel">
                    <div class="export-list">
                        <div class="list-item list-item-section">
                            <Field class="form-field" let:uniqueId>
                                <input
                                    type="checkbox"
                                    id={uniqueId}
                                    disabled={!collections.length}
                                    checked={areAllSelected}
                                    on:change={() => toggleSelectAll()}
                                />
                                <label for={uniqueId}>Select all</label>
                            </Field>
                        </div>
                        {#each collections as collection (collection.id)}
                            <div class="list-item list-item-collection">
                                <Field class="form-field" let:uniqueId>
                                    <input
                                        type="checkbox"
                                        id={uniqueId}
                                        checked={bulkSelected[collection.id]}
                                        on:change={() => toggleSelectCollection(collection)}
                                    />
                                    <label for={uniqueId} title={collection.name}>{collection.name}</label>
                                </Field>
                            </div>
                        {/each}
                    </div>

                    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
                    <!-- svelte-ignore a11y-no-static-element-interactions -->
                    <div
                        bind:this={previewContainer}
                        tabindex="0"
                        class="export-preview"
                        on:keydown={(e) => {
                            // select all
                            if (e.ctrlKey && e.code === "KeyA") {
                                e.preventDefault();
                                const selection = window.getSelection();
                                const range = document.createRange();
                                range.selectNodeContents(previewContainer);
                                selection.removeAllRanges();
                                selection.addRange(range);
                            }
                        }}
                    >
                        <button
                            type="button"
                            class="btn btn-sm btn-transparent fade copy-schema"
                            disabled={!totalBulkSelected}
                            on:click={() => copy()}
                        >
                            <span class="txt">Copy</span>
                        </button>

                        <pre class="code-wrapper">{schema}</pre>
                    </div>
                </div>

                <div class="flex m-t-base">
                    <div class="flex-fill" />
                    <button
                        type="button"
                        class="btn btn-expanded"
                        disabled={!totalBulkSelected}
                        on:click={() => download()}
                    >
                        <i class="ri-download-line" />
                        <span class="txt">Download as JSON</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</PageWrapper>
