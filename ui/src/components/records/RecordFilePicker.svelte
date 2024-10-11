<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Scroller from "@/components/base/Scroller.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";
    import { collections } from "@/stores/collections";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();
    const uniqueId = "file_picker_" + CommonHelper.randomString(5);
    const batchSize = 50;

    export let title = "Select a file";
    export let submitText = "Insert";
    export let fileTypes = ["image", "document", "video", "audio", "file"];

    let pickerPanel;
    let upsertPanel;
    let filter = "";
    let list = [];
    let currentPage = 1;
    let lastItemsCount = 0;
    let isLoading = false;
    let fileCollections = [];
    let fileFields = [];
    let sizeOptions = [];
    let selectedCollection = {};
    let selectedFile = {};
    let selectedSize = "";

    // find all collections with at least one non-protected file field
    $: fileCollections = $collections.filter((c) => {
        return (
            c.type !== "view" &&
            !!CommonHelper.toArray(c.fields).find((f) => {
                return (
                    // is file field
                    f.type === "file" &&
                    // is public (aka. doesn't require file token)
                    !f.protected &&
                    // allow any MIME type OR image/*
                    (!f.mimeTypes?.length || !!f.mimeTypes?.find((t) => t.startsWith("image/")))
                );
            })
        );
    });

    // auto select the first collection from the list
    $: if (!selectedCollection?.id && fileCollections.length > 0) {
        selectedCollection = fileCollections[0];
    }

    $: fileFields = selectedCollection?.fields?.filter((f) => f.type === "file" && !f.protected);

    // reset filter on collection change
    $: if (selectedCollection?.id) {
        clearFilter();
        refreshSizeOptions();
    }

    // refresh the size options on selected file change
    $: if (selectedFile?.name) {
        refreshSizeOptions();
    }

    // reset list on filter or collection change
    $: if (typeof filter !== "undefined" && selectedCollection?.id && pickerPanel?.isActive()) {
        loadList(true);
    }

    $: isSelected = (record, name) => {
        return selectedFile?.name == name && selectedFile?.record?.id == record.id;
    };

    $: hasAtleastOneFile = list.find((r) => extractFiles(r).length > 0);

    $: canLoadMore = !isLoading && lastItemsCount == batchSize;

    $: canSubmit = !isLoading && !!selectedFile?.name;

    export function show() {
        loadList(true);

        return pickerPanel?.show();
    }

    export function hide() {
        return pickerPanel?.hide();
    }

    function clearList() {
        list = [];
        selectedFile = {};
        selectedSize = "";
    }

    function clearFilter() {
        filter = "";
    }

    async function loadList(reset = false) {
        if (!selectedCollection?.id) {
            return;
        }

        isLoading = true;

        if (reset) {
            clearList();
        }

        try {
            const page = reset ? 1 : currentPage + 1;

            const fallbackSearchFields = CommonHelper.getAllCollectionIdentifiers(selectedCollection);

            let normalizedFilter = CommonHelper.normalizeSearchFilter(filter, fallbackSearchFields) || "";

            if (normalizedFilter) {
                normalizedFilter += " && ";
            }
            normalizedFilter += "(" + fileFields.map((f) => `${f.name}:length>0`).join("||") + ")";

            let sort = "";
            if (selectedCollection.type != "view") {
                sort = "-@rowid"; // all collections with exception to the view has this field
            }

            const result = await ApiClient.collection(selectedCollection.id).getList(page, batchSize, {
                filter: normalizedFilter,
                sort: sort,
                fields: "*:excerpt(100)",
                skipTotal: 1,
                requestKey: uniqueId + "loadImagePicker",
            });

            list = CommonHelper.filterDuplicatesByKey(list.concat(result.items));
            currentPage = result.page;
            lastItemsCount = result.items.length;

            isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
                isLoading = false;
            }
        }
    }

    function refreshSizeOptions() {
        let sizes = ["100x100"]; // default Superuser UI thumb

        // extract the thumb sizes of the selected file field
        if (selectedFile?.record?.id) {
            for (const field of fileFields) {
                if (CommonHelper.toArray(selectedFile.record[field.name]).includes(selectedFile.name)) {
                    sizes = sizes.concat(CommonHelper.toArray(field.thumbs));
                    break;
                }
            }
        }

        // construct the dropdown options
        sizeOptions = [{ label: "Original size", value: "" }];
        for (const size of sizes) {
            sizeOptions.push({
                label: `${size} thumb`,
                value: size,
            });
        }

        // reset selected size if missing
        if (selectedSize && !sizes.includes(selectedSize)) {
            selectedSize = "";
        }
    }

    function extractFiles(record) {
        let result = [];

        for (const field of fileFields) {
            const names = CommonHelper.toArray(record[field.name]);
            for (const name of names) {
                if (fileTypes.includes(CommonHelper.getFileType(name))) {
                    result.push(name);
                }
            }
        }

        return result;
    }

    function select(record, name) {
        selectedFile = { record, name };
    }

    function submit() {
        if (!canSubmit) {
            return;
        }

        dispatch(
            "submit",
            Object.assign(
                {
                    size: selectedSize,
                },
                selectedFile,
            ),
        );

        hide();
    }
</script>

<OverlayPanel bind:this={pickerPanel} popup class="file-picker-popup" on:hide on:show {...$$restProps}>
    <svelte:fragment slot="header">
        <h4>{title}</h4>
    </svelte:fragment>

    {#if !fileCollections.length}
        <h6 class="txt-center txt-hint">
            You currently don't have any collection with <code>file</code> field.
        </h6>
    {:else}
        <div class="file-picker">
            <aside class="file-picker-sidebar">
                {#each fileCollections as collection (collection.id)}
                    <button
                        type="button"
                        class="sidebar-item"
                        class:active={selectedCollection?.id == collection.id}
                        on:click|preventDefault={() => {
                            selectedCollection = collection;
                        }}
                    >
                        {collection.name}
                    </button>
                {/each}
            </aside>

            <div class="file-picker-content">
                <div class="flex m-b-base flex-gap-10">
                    <Searchbar
                        value={filter}
                        placeholder="Record search term or filter..."
                        autocompleteCollection={selectedCollection}
                        on:submit={(e) => (filter = e.detail)}
                    />
                    <button
                        type="button"
                        class="btn btn-pill btn-transparent btn-hint p-l-xs p-r-xs"
                        on:click={() => upsertPanel?.show()}
                    >
                        <div class="txt">New record</div>
                    </button>
                </div>
                <Scroller
                    class="files-list"
                    vThreshold={100}
                    on:vScrollEnd={() => {
                        if (canLoadMore) {
                            loadList();
                        }
                    }}
                >
                    {#if hasAtleastOneFile}
                        {#each list as record (record.id)}
                            {@const names = extractFiles(record)}
                            {#each names as name}
                                <button
                                    type="button"
                                    class="thumb handle"
                                    use:tooltip={name + "\n(record: " + record.id + ")"}
                                    class:thumb-warning={isSelected(record, name)}
                                    on:click|preventDefault={select(record, name)}
                                >
                                    {#if CommonHelper.hasImageExtension(name)}
                                        <img
                                            loading="lazy"
                                            src={ApiClient.files.getURL(record, name, { thumb: "100x100" })}
                                            alt={name}
                                        />
                                    {:else}
                                        <i class="ri-file-3-line" />
                                    {/if}
                                </button>
                            {/each}
                        {/each}
                    {:else if !isLoading}
                        <div class="inline-flex">
                            <span class="txt txt-hint">No records with images found.</span>
                            {#if filter?.length}
                                <button
                                    type="button"
                                    class="btn btn-hint btn-sm"
                                    on:click|preventDefault={clearFilter}
                                >
                                    <span class="txt">Clear filter</span>
                                </button>
                            {/if}
                        </div>
                    {/if}

                    {#if isLoading}
                        <div class="block txt-center">
                            <span class="loader loader-sm active" />
                        </div>
                    {/if}
                </Scroller>
            </div>
        </div>
    {/if}

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent m-r-auto" disabled={isLoading} on:click={hide}>
            <span class="txt">Cancel</span>
        </button>

        {#if CommonHelper.hasImageExtension(selectedFile?.name)}
            <Field class="form-field file-picker-size-select" let:uniqueId>
                <ObjectSelect
                    upside
                    id={uniqueId}
                    items={sizeOptions}
                    disabled={!canSubmit}
                    selectPlaceholder="Select size"
                    bind:keyOfSelected={selectedSize}
                />
            </Field>
        {/if}

        <button type="button" class="btn btn-expanded" disabled={!canSubmit} on:click={submit}>
            <span class="txt">{submitText}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<RecordUpsertPanel
    bind:this={upsertPanel}
    collection={selectedCollection}
    on:save={(e) => {
        CommonHelper.removeByKey(list, "id", e.detail.record.id);
        list.unshift(e.detail.record);
        list = list;

        const names = extractFiles(e.detail.record);
        if (names.length > 0) {
            select(e.detail.record, names[0]);
        }
    }}
    on:delete={(e) => {
        if (selectedFile?.record?.id == e.detail.id) {
            selectedFile = {};
        }

        CommonHelper.removeByKey(list, "id", e.detail.id);
        list = list;
    }}
/>
