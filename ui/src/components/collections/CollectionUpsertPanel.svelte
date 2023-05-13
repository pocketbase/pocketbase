<script>
    import { createEventDispatcher, tick } from "svelte";
    import { scale } from "svelte/transition";
    import { Collection } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import { errors, setErrors, removeError } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { removeAllToasts, addSuccessToast } from "@/stores/toasts";
    import { addCollection, removeCollection } from "@/stores/collections";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import CollectionFieldsTab from "@/components/collections/CollectionFieldsTab.svelte";
    import CollectionRulesTab from "@/components/collections/CollectionRulesTab.svelte";
    import CollectionQueryTab from "@/components/collections/CollectionQueryTab.svelte";
    import CollectionAuthOptionsTab from "@/components/collections/CollectionAuthOptionsTab.svelte";
    import CollectionUpdateConfirm from "@/components/collections/CollectionUpdateConfirm.svelte";

    const TAB_SCHEMA = "schema";
    const TAB_RULES = "api_rules";
    const TAB_OPTIONS = "options";

    const TYPE_BASE = "base";
    const TYPE_AUTH = "auth";
    const TYPE_VIEW = "view";

    const collectionTypes = {};
    collectionTypes[TYPE_BASE] = "Base";
    collectionTypes[TYPE_VIEW] = "View";
    collectionTypes[TYPE_AUTH] = "Auth";

    const dispatch = createEventDispatcher();

    let collectionPanel;
    let confirmChangesPanel;
    let original = null;
    let collection = new Collection();
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let activeTab = TAB_SCHEMA;
    let initialFormHash = calculateFormHash(collection);
    let schemaTabError = "";

    $: if ($errors.schema || $errors.options?.query) {
        // extract the direct schema field error, otherwise - return a generic message
        schemaTabError = CommonHelper.getNestedVal($errors, "schema.message") || "Has errors";
    } else {
        schemaTabError = "";
    }

    $: isSystemUpdate = !collection.$isNew && collection.system;

    $: hasChanges = initialFormHash != calculateFormHash(collection);

    $: canSave = collection.$isNew || hasChanges;

    $: if (activeTab === TAB_OPTIONS && collection.type !== TYPE_AUTH) {
        // reset selected tab
        changeTab(TAB_SCHEMA);
    }

    $: if (collection.type === TYPE_VIEW) {
        // reset non-view fields
        collection.createRule = null;
        collection.updateRule = null;
        collection.deleteRule = null;
        collection.indexes = [];
    }

    // update indexes on collection rename
    $: if (collection?.name && original?.name != collection?.name) {
        collection.indexes = collection.indexes?.map((idx) =>
            CommonHelper.replaceIndexTableName(idx, collection.name)
        );
    }

    export function changeTab(newTab) {
        activeTab = newTab;
    }

    export function show(model) {
        load(model);

        confirmClose = true;

        changeTab(TAB_SCHEMA);

        return collectionPanel?.show();
    }

    export function hide() {
        return collectionPanel?.hide();
    }

    async function load(model) {
        setErrors({}); // reset errors

        if (typeof model !== "undefined") {
            original = model;
            collection = model.$clone();
        } else {
            original = null;
            collection = new Collection();
        }

        // normalize
        collection.schema = collection.schema || [];
        collection.originalName = collection.name || "";

        await tick();

        initialFormHash = calculateFormHash(collection);
    }

    function saveConfirm() {
        if (collection.$isNew) {
            save();
        } else {
            confirmChangesPanel?.show(original, collection);
        }
    }

    function save() {
        if (isSaving) {
            return;
        }

        isSaving = true;

        const data = exportFormData();

        let request;
        if (collection.$isNew) {
            request = ApiClient.collections.create(data);
        } else {
            request = ApiClient.collections.update(collection.id, data);
        }

        request
            .then((result) => {
                removeAllToasts();

                addCollection(result);

                confirmClose = false;
                hide();

                addSuccessToast(
                    collection.$isNew
                        ? "Successfully created collection."
                        : "Successfully updated collection."
                );

                dispatch("save", {
                    isNew: collection.$isNew,
                    collection: result,
                });
            })
            .catch((err) => {
                ApiClient.error(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function exportFormData() {
        const data = collection.$export();
        data.schema = data.schema.slice(0);

        // remove deleted fields
        for (let i = data.schema.length - 1; i >= 0; i--) {
            const field = data.schema[i];
            if (field.toDelete) {
                data.schema.splice(i, 1);
            }
        }

        return data;
    }

    function deleteConfirm() {
        if (!original?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete collection "${original?.name}" and all its records?`, () => {
            return ApiClient.collections
                .delete(original?.id)
                .then(() => {
                    hide();
                    addSuccessToast(`Successfully deleted collection "${original?.name}".`);
                    dispatch("delete", original);
                    removeCollection(original);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function calculateFormHash(m) {
        return JSON.stringify(m);
    }

    function setCollectionType(t) {
        collection.type = t;

        // reset schema errors on type change
        removeError("schema");
    }

    function duplicateConfirm() {
        if (hasChanges) {
            confirm("You have unsaved changes. Do you really want to discard them?", () => {
                duplicate();
            });
        } else {
            duplicate();
        }
    }

    async function duplicate() {
        const clone = original?.$clone();

        if (clone) {
            clone.id = "";
            clone.created = "";
            clone.updated = "";
            clone.name += "_duplicate";

            // reset the schema
            if (!CommonHelper.isEmpty(clone.schema)) {
                for (const field of clone.schema) {
                    field.id = "";
                }
            }

            // update indexes with the new table name
            if (!CommonHelper.isEmpty(clone.indexes)) {
                for (let i = 0; i < clone.indexes.length; i++) {
                    const parsed = CommonHelper.parseIndex(clone.indexes[i]);
                    parsed.indexName = "idx_" + CommonHelper.randomString(7);
                    parsed.tableName = clone.name;
                    clone.indexes[i] = CommonHelper.buildIndex(parsed);
                }
            }
        }

        show(clone);

        await tick();

        initialFormHash = "";
    }
</script>

<OverlayPanel
    bind:this={collectionPanel}
    class="overlay-panel-lg colored-header collection-panel"
    escClose={false}
    overlayClose={!isSaving}
    beforeHide={() => {
        if (hasChanges && confirmClose) {
            confirm("You have unsaved changes. Do you really want to close the panel?", () => {
                confirmClose = false;
                hide();
            });
            return false;
        }
        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4 class="upsert-panel-title">
            {collection.$isNew ? "New collection" : "Edit collection"}
        </h4>

        {#if !collection.$isNew && !collection.system}
            <div class="flex-fill" />
            <button type="button" aria-label="More" class="btn btn-sm btn-circle btn-transparent flex-gap-0">
                <i class="ri-more-line" />
                <Toggler class="dropdown dropdown-right m-t-5">
                    <button type="button" class="dropdown-item closable" on:click={() => duplicateConfirm()}>
                        <i class="ri-file-copy-line" />
                        <span class="txt">Duplicate</span>
                    </button>
                    <button
                        type="button"
                        class="dropdown-item txt-danger closable"
                        on:click|preventDefault|stopPropagation={() => deleteConfirm()}
                    >
                        <i class="ri-delete-bin-7-line" />
                        <span class="txt">Delete</span>
                    </button>
                </Toggler>
            </button>
        {/if}

        <form
            class="block"
            on:submit|preventDefault={() => {
                canSave && saveConfirm();
            }}
        >
            <Field
                class="form-field collection-field-name required m-b-0 {isSystemUpdate ? 'disabled' : ''}"
                name="name"
                let:uniqueId
            >
                <label for={uniqueId}>Name</label>

                <!-- svelte-ignore a11y-autofocus -->
                <input
                    type="text"
                    id={uniqueId}
                    required
                    disabled={isSystemUpdate}
                    spellcheck="false"
                    autofocus={collection.$isNew}
                    placeholder={collection.$isAuth ? `eg. "users"` : `eg. "posts"`}
                    value={collection.name}
                    on:input={(e) => {
                        collection.name = CommonHelper.slugify(e.target.value);
                        e.target.value = collection.name;
                    }}
                />

                <div class="form-field-addon">
                    <button
                        type="button"
                        class="btn btn-sm p-r-10 p-l-10 {collection.$isNew
                            ? 'btn-outline'
                            : 'btn-transparent'}"
                        disabled={!collection.$isNew}
                    >
                        <!-- empty span for alignment -->
                        <span />
                        <span class="txt">Type: {collectionTypes[collection.type] || "N/A"}</span>
                        {#if collection.$isNew}
                            <i class="ri-arrow-down-s-fill" />
                            <Toggler class="dropdown dropdown-right dropdown-nowrap m-t-5">
                                {#each Object.entries(collectionTypes) as [type, label]}
                                    <button
                                        type="button"
                                        class="dropdown-item closable"
                                        class:selected={type == collection.type}
                                        on:click={() => setCollectionType(type)}
                                    >
                                        <i class={CommonHelper.getCollectionTypeIcon(type)} />
                                        <span class="txt">{label} collection</span>
                                    </button>
                                {/each}
                            </Toggler>
                        {/if}
                    </button>
                </div>

                {#if collection.system}
                    <div class="help-block">System collection</div>
                {/if}
            </Field>

            <input type="submit" class="hidden" tabindex="-1" />
        </form>

        <div class="tabs-header stretched">
            <button
                type="button"
                class="tab-item"
                class:active={activeTab === TAB_SCHEMA}
                on:click={() => changeTab(TAB_SCHEMA)}
            >
                <span class="txt">{collection?.$isView ? "Query" : "Fields"}</span>
                {#if !CommonHelper.isEmpty(schemaTabError)}
                    <i
                        class="ri-error-warning-fill txt-danger"
                        transition:scale|local={{ duration: 150, start: 0.7 }}
                        use:tooltip={schemaTabError}
                    />
                {/if}
            </button>

            <button
                type="button"
                class="tab-item"
                class:active={activeTab === TAB_RULES}
                on:click={() => changeTab(TAB_RULES)}
            >
                <span class="txt">API Rules</span>
                {#if !CommonHelper.isEmpty($errors?.listRule) || !CommonHelper.isEmpty($errors?.viewRule) || !CommonHelper.isEmpty($errors?.createRule) || !CommonHelper.isEmpty($errors?.updateRule) || !CommonHelper.isEmpty($errors?.deleteRule) || !CommonHelper.isEmpty($errors?.options?.manageRule)}
                    <i
                        class="ri-error-warning-fill txt-danger"
                        transition:scale|local={{ duration: 150, start: 0.7 }}
                        use:tooltip={"Has errors"}
                    />
                {/if}
            </button>

            {#if collection.$isAuth}
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === TAB_OPTIONS}
                    on:click={() => changeTab(TAB_OPTIONS)}
                >
                    <span class="txt">Options</span>
                    {#if !CommonHelper.isEmpty($errors?.options) && !$errors?.options?.manageRule}
                        <i
                            class="ri-error-warning-fill txt-danger"
                            transition:scale|local={{ duration: 150, start: 0.7 }}
                            use:tooltip={"Has errors"}
                        />
                    {/if}
                </button>
            {/if}
        </div>
    </svelte:fragment>

    <div class="tabs-content">
        <!-- avoid rerendering the fields tab -->
        <div class="tab-item" class:active={activeTab === TAB_SCHEMA}>
            {#if collection.$isView}
                <CollectionQueryTab bind:collection />
            {:else}
                <CollectionFieldsTab bind:collection />
            {/if}
        </div>

        {#if activeTab === TAB_RULES}
            <div class="tab-item active">
                <CollectionRulesTab bind:collection />
            </div>
        {/if}

        {#if collection.$isAuth}
            <div class="tab-item" class:active={activeTab === TAB_OPTIONS}>
                <CollectionAuthOptionsTab bind:collection />
            </div>
        {/if}
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
            type="button"
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!canSave || isSaving}
            on:click={() => saveConfirm()}
        >
            <span class="txt">{collection.$isNew ? "Create" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<CollectionUpdateConfirm bind:this={confirmChangesPanel} on:confirm={() => save()} />

<style>
    .upsert-panel-title {
        display: inline-flex;
        align-items: center;
        min-height: var(--smBtnHeight);
    }
    .tabs-content:focus-within {
        z-index: 9; /* autocomplete dropdown overlay fix */
    }
</style>
