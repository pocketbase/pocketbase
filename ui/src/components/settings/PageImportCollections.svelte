<script>
    import { tick } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addErrorToast } from "@/stores/toasts";
    import { setErrors } from "@/stores/errors";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Field from "@/components/base/Field.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import ImportPopup from "@/components/settings/ImportPopup.svelte";

    $pageTitle = "Import collections";

    let fileInput;
    let importPopup;

    let schemas = "";
    let isLoadingFile = false;
    let newCollections = [];
    let oldCollections = [];
    let collectionsToModify = [];
    let isLoadingOldCollections = false;

    $: if (typeof schemas !== "undefined") {
        loadNewCollections(schemas);
    }

    $: isValid =
        !!schemas &&
        newCollections.length &&
        newCollections.length === newCollections.filter((item) => !!item.id && !!item.name).length;

    $: collectionsToDelete = oldCollections.filter((collection) => {
        return isValid && !CommonHelper.findByKey(newCollections, "id", collection.id);
    });

    $: collectionsToAdd = newCollections.filter((collection) => {
        return isValid && !CommonHelper.findByKey(oldCollections, "id", collection.id);
    });

    $: if (typeof newCollections !== "undefined") {
        loadCollectionsToModify();
    }

    $: hasChanges =
        !!schemas && (collectionsToDelete.length || collectionsToAdd.length || collectionsToModify.length);

    $: canImport = !isLoadingOldCollections && isValid && hasChanges;

    loadOldCollections();

    async function loadOldCollections() {
        isLoadingOldCollections = true;

        try {
            oldCollections = await ApiClient.collections.getFullList(200);
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

    function loadCollectionsToModify() {
        collectionsToModify = [];

        if (!isValid) {
            return;
        }

        for (let newCollection of newCollections) {
            const oldCollection = CommonHelper.findByKey(oldCollections, "id", newCollection.id);
            if (
                // no old collection
                !oldCollection?.id ||
                // no changes
                JSON.stringify(oldCollection) === JSON.stringify(newCollection)
            ) {
                continue;
            }

            collectionsToModify.push({
                new: newCollection,
                old: oldCollection,
            });
        }
    }

    function loadNewCollections() {
        newCollections = [];

        try {
            newCollections = JSON.parse(schemas);
        } catch (_) {}

        if (!Array.isArray(newCollections)) {
            newCollections = [];
        } else {
            newCollections = CommonHelper.filterDuplicatesByKey(newCollections);
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

        reader.onload = async (event) => {
            isLoadingFile = false;
            fileInput.value = ""; // reset

            schemas = event.target.result;

            await tick();

            if (!newCollections.length) {
                addErrorToast("Invalid collections configuration.");
                clear();
            }
        };

        reader.onerror = (err) => {
            console.log(err);
            addErrorToast("Failed to load the imported JSON.");

            isLoadingFile = false;
            fileInput.value = ""; // reset
        };

        reader.readAsText(file);
    }

    function clear() {
        schemas = "";
        fileInput.value = "";
        setErrors({});
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
            {#if isLoadingOldCollections}
                <div class="loader" />
            {:else}
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

                <div class="content txt-xl m-b-base">
                    <p>
                        Paste below the collections configuration you want to import or
                        <button
                            class="btn btn-outline btn-sm m-l-5"
                            class:btn-loading={isLoadingFile}
                            on:click={() => {
                                fileInput.click();
                            }}
                        >
                            <span class="txt">Load from JSON file</span>
                        </button>
                    </p>
                </div>

                <Field class="form-field {!isValid ? 'field-error' : ''}" name="collections" let:uniqueId>
                    <label for={uniqueId} class="p-b-10">Collections</label>
                    <textarea
                        id={uniqueId}
                        class="code"
                        spellcheck="false"
                        rows="15"
                        required
                        bind:value={schemas}
                    />

                    {#if !!schemas && !isValid}
                        <div class="help-block help-block-error">Invalid collections configuration.</div>
                    {/if}
                </Field>

                {#if isValid && newCollections.length && !hasChanges}
                    <div class="alert alert-info">
                        <div class="icon">
                            <i class="ri-information-line" />
                        </div>
                        <div class="content">
                            <string>Your collections configuration is already up-to-date!</string>
                        </div>
                    </div>
                {/if}

                {#if isValid && newCollections.length && hasChanges}
                    <h5 class="section-title">Detected changes</h5>

                    <div class="list">
                        {#if collectionsToDelete.length}
                            {#each collectionsToDelete as collection (collection.id)}
                                <div class="list-item">
                                    <span class="label label-danger list-label">Deleted</span>
                                    <strong>{collection.name}</strong>
                                    {#if collection.id}
                                        <small class="txt-hint">({collection.id})</small>
                                    {/if}
                                </div>
                            {/each}
                        {/if}

                        {#if collectionsToModify.length}
                            {#each collectionsToModify as pair (pair.old.id + pair.new.id)}
                                <div class="list-item">
                                    <span class="label label-warning list-label">Modified</span>
                                    <strong>
                                        {#if pair.old.name !== pair.new.name}
                                            <span class="txt-strikethrough txt-hint">{pair.old.name}</span> -
                                        {/if}
                                        {pair.new.name}
                                    </strong>
                                    {#if pair.new.id}
                                        <small class="txt-hint">({pair.new.id})</small>
                                    {/if}
                                </div>
                            {/each}
                        {/if}

                        {#if collectionsToAdd.length}
                            {#each collectionsToAdd as collection (collection.id)}
                                <div class="list-item">
                                    <span class="label label-success list-label">New</span>
                                    <strong>{collection.name}</strong>
                                    {#if collection.id}
                                        <small class="txt-hint">({collection.id})</small>
                                    {/if}
                                </div>
                            {/each}
                        {/if}
                    </div>
                {/if}

                <div class="flex m-t-base">
                    {#if !!schemas}
                        <button type="button" class="btn btn-secondary link-hint" on:click={() => clear()}>
                            <span class="txt">Clear</span>
                        </button>
                    {/if}
                    <div class="flex-fill" />
                    <button
                        type="button"
                        class="btn btn-expanded btn-warning m-l-auto"
                        disabled={!canImport}
                        on:click={() => importPopup?.show(oldCollections, newCollections)}
                    >
                        <span class="txt">Review</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</PageWrapper>

<ImportPopup bind:this={importPopup} on:submit={() => clear()} />

<style>
    .list-label {
        min-width: 65px;
    }
</style>
