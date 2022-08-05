<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addInfoToast, addErrorToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";

    $pageTitle = "Import collections";

    let uniquePageId = "import_" + CommonHelper.randomString(5);

    let fileInput;

    let schema = "";
    let isImporting = false;
    let isLoadingFile = false;
    let newCollections = [];
    let oldCollections = [];
    let isLoadingOldCollections = false;

    $: if (typeof schema !== "undefined") {
        loadNewCollections(schema);
    }

    $: isValid =
        !!schema &&
        newCollections.length &&
        newCollections.length === newCollections.filter((item) => !!item.id && !!item.name).length;

    $: canImport = isValid && !isLoadingOldCollections;

    $: collectionsToDelete = oldCollections.filter((collection) => {
        return !CommonHelper.findByKey(newCollections, "id", collection.id);
    });

    $: collectionsToAdd = newCollections.filter((collection) => {
        return !CommonHelper.findByKey(oldCollections, "id", collection.id);
    });

    $: collectionsToModify = newCollections.filter((newCollection) => {
        const oldCollection = CommonHelper.findByKey(oldCollections, "id", newCollection.id);
        if (!oldCollection?.id) {
            return false;
        }

        return JSON.stringify(oldCollection) !== JSON.stringify(newCollection);
    });

    loadOldCollections();

    async function loadOldCollections() {
        isLoadingOldCollections = true;

        try {
            oldCollections = await ApiClient.collections.getFullList(100, {
                $cancelKey: uniquePageId,
            });
            // delete timestamps
            for (let collection of oldCollections) {
                delete collection.created;
                delete collection.updated;
            }
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingOldCollections = false;
    }

    function loadNewCollections() {
        newCollections = [];

        try {
            newCollections = JSON.parse(schema);
        } catch (_) {}

        if (!Array.isArray(newCollections)) {
            newCollections = [];
        }

        // delete timestamps
        for (let collection of newCollections) {
            delete collection.created;
            delete collection.updated;
        }
    }

    function loadFile(file) {
        isLoadingFile = true;

        const reader = new FileReader();

        reader.onload = (event) => {
            schema = event.target.result;

            isLoadingFile = false;
            fileInput.value = ""; // reset
        };

        reader.onerror = (err) => {
            console.log(err);
            addErrorToast("Failed to load the imported JSON.");

            isLoadingFile = false;
            fileInput.value = ""; // reset
        };

        reader.readAsText(file);
    }

    function submitImport() {
        isImporting = true;

        try {
            const newCollections = JSON.parse(schema);
            ApiClient.collections.import(newCollections);
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isImporting = false;
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
            <div class="content txt-xl m-b-base">
                <input
                    bind:this={fileInput}
                    type="file"
                    class="hidden"
                    accept=".json"
                    on:change={() => {
                        if (fileInput.files.length) {
                            loadFile(fileInput.files[0]);
                        }
                    }}
                />

                <p>
                    Paste below the collections schema you want to import or
                    <button
                        class="btn btn-outline btn-sm"
                        class:btn-loading={isLoadingFile}
                        on:click={() => {
                            fileInput.click();
                        }}
                    >
                        <span class="txt">Import from JSON file</span>
                    </button>
                </p>
            </div>

            <Field class="form-field {!isValid ? 'field-error' : ''}" name="collections" let:uniqueId>
                <label for={uniqueId}>Collections schema</label>
                <textarea
                    id={uniqueId}
                    class="json-editor"
                    spellcheck="false"
                    rows="20"
                    required
                    bind:value={schema}
                />
                {#if !!schema && !isValid}
                    <div class="help-block help-block-error">Invalid collections schema.</div>
                {/if}
            </Field>

            <div class="section-title">Detected changes</div>
            <p>No changes to your current collections schema were found.</p>

            {#each collectionsToDelete as collection (collection.id)}
                Delete {collection.name}
                <br />
            {/each}

            {#each collectionsToModify as collection (collection.id)}
                Modify {collection.name}
                <br />
            {/each}

            {#each collectionsToAdd as collection (collection.id)}
                Add {collection.name}
                <br />
            {/each}

            <div class="flex m-t-base">
                <div class="flex-fill" />
                <button
                    type="button"
                    class="btn btn-expanded"
                    class:btn-loading={isImporting}
                    disabled={!canImport}
                    on:click={() => submitImport()}
                >
                    <span class="txt">Import</span>
                </button>
            </div>
        </div>
    </div>
</main>

<style>
    .json-editor {
        font-size: 15px;
        line-height: 1.379rem;
        font-family: var(--monospaceFontFamily);
    }
</style>
