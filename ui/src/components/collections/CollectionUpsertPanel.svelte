<script>
    import { createEventDispatcher, tick } from "svelte";
    import { scale } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { confirm } from "@/stores/confirmation";
    import { errors, removeError, setErrors } from "@/stores/errors";
    import { addSuccessToast, removeAllToasts } from "@/stores/toasts";
    import {
        addCollection,
        removeCollection,
        scaffolds,
        activeCollection,
        refreshScaffolds,
    } from "@/stores/collections";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import CollectionAuthOptionsTab from "@/components/collections/CollectionAuthOptionsTab.svelte";
    import CollectionFieldsTab from "@/components/collections/CollectionFieldsTab.svelte";
    import CollectionQueryTab from "@/components/collections/CollectionQueryTab.svelte";
    import CollectionRulesTab from "@/components/collections/CollectionRulesTab.svelte";
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
    let collection = {};
    let isSaving = false;
    let isLoadingConfirmation = false;
    let confirmClose = false; // prevent close recursion
    let activeTab = TAB_SCHEMA;
    let initialFormHash = calculateFormHash(collection);
    let fieldsTabError = "";
    let baseCollectionKeys = [];

    $: baseCollectionKeys = Object.keys($scaffolds["base"] || {});

    $: isAuth = collection.type === TYPE_AUTH;

    $: isView = collection.type === TYPE_VIEW;

    $: if ($errors.fields || $errors.viewQuery || $errors.indexes) {
        // extract the direct fields list error, otherwise - return a generic message
        fieldsTabError = CommonHelper.getNestedVal($errors, "fields.message") || "Has errors";
    } else {
        fieldsTabError = "";
    }

    $: isSystemUpdate = !!collection.id && collection.system;

    $: isSuperusers = !!collection.id && collection.system && collection.name == "_superusers";

    $: hasChanges = initialFormHash != calculateFormHash(collection);

    $: canSave = !collection.id || hasChanges;

    $: if (activeTab === TAB_OPTIONS && collection.type !== "auth") {
        // reset selected tab
        changeTab(TAB_SCHEMA);
    }

    $: if (collection.type === "view") {
        // reset non-view fields
        collection.createRule = null;
        collection.updateRule = null;
        collection.deleteRule = null;
        collection.indexes = [];
    }

    // update indexes on collection rename
    $: if (collection.name && original?.name != collection.name && collection.indexes.length > 0) {
        collection.indexes = collection.indexes?.map((idx) =>
            CommonHelper.replaceIndexTableName(idx, collection.name),
        );
    }

    export function changeTab(newTab) {
        activeTab = newTab;
    }

    export function show(model) {
        load(model);

        confirmClose = true;
        isLoadingConfirmation = false;
        isSaving = false;

        changeTab(TAB_SCHEMA);

        return collectionPanel?.show();
    }

    export function hide() {
        return collectionPanel?.hide();
    }

    export function forceHide() {
        confirmClose = false;
        hide();
    }

    async function load(model) {
        setErrors({}); // reset errors

        if (typeof model !== "undefined") {
            original = model;
            collection = structuredClone(model);
        } else {
            original = null;
            collection = structuredClone($scaffolds["base"]);

            // add default timestamp fields
            collection.fields.push({
                type: "autodate",
                name: "created",
                onCreate: true,
            });
            collection.fields.push({
                type: "autodate",
                name: "updated",
                onCreate: true,
                onUpdate: true,
            });
        }

        // normalize
        collection.fields = collection.fields || [];
        collection._originalName = collection.name || "";

        await tick();

        initialFormHash = calculateFormHash(collection);
    }

    async function saveConfirm(hideAfterSave = true) {
        if (isLoadingConfirmation) {
            return;
        }

        isLoadingConfirmation = true;

        try {
            if (!collection.id) {
                await save(hideAfterSave);
            } else {
                await confirmChangesPanel?.show(original, collection, hideAfterSave);
            }
        } catch {}

        isLoadingConfirmation = false;
    }

    async function save(hideAfterSave = true) {
        if (isSaving) {
            return;
        }

        isSaving = true;

        const data = exportFormData();
        const isNew = !collection.id;

        try {
            let result;
            if (isNew) {
                result = await ApiClient.collections.create(data);
            } else {
                result = await ApiClient.collections.update(collection.id, data);
            }

            removeAllToasts();

            addCollection(result);

            if (hideAfterSave) {
                confirmClose = false;
                hide();
            } else {
                load(result);
            }

            addSuccessToast(
                !collection.id ? "Successfully created collection." : "Successfully updated collection.",
            );

            dispatch("save", {
                isNew: isNew,
                collection: result,
            });

            if (isNew) {
                $activeCollection = result;

                await refreshScaffolds();
            }
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function exportFormData() {
        const data = Object.assign({}, collection);
        data.fields = data.fields.slice(0);

        // remove deleted fields
        for (let i = data.fields.length - 1; i >= 0; i--) {
            const field = data.fields[i];
            if (field._toDelete) {
                data.fields.splice(i, 1);
            }
        }

        return data;
    }

    function truncateConfirm() {
        if (!original?.id) {
            return; // nothing to truncate
        }

        confirm(
            `Do you really want to delete all "${original.name}" records, including their cascade delete references and files?`,
            () => {
                return ApiClient.collections
                    .truncate(original.id)
                    .then(() => {
                        forceHide();
                        addSuccessToast(`Successfully truncated collection "${original.name}".`);
                        dispatch("truncate");
                    })
                    .catch((err) => {
                        ApiClient.error(err);
                    });
            },
        );
    }

    function deleteConfirm() {
        if (!original?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete collection "${original.name}" and all its records?`, () => {
            return ApiClient.collections
                .delete(original.id)
                .then(() => {
                    forceHide();
                    addSuccessToast(`Successfully deleted collection "${original.name}".`);
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

        // merge with the scaffold to ensure that the minimal props are set
        collection = Object.assign(structuredClone($scaffolds[t]), collection);

        // reset fields list errors on type change
        removeError("fields");
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
        const clone = original ? structuredClone(original) : null;

        if (clone) {
            clone.id = "";
            clone.created = "";
            clone.updated = "";
            clone.name += "_duplicate";

            // reset the fields list
            if (!CommonHelper.isEmpty(clone.fields)) {
                for (const field of clone.fields) {
                    field.id = "";
                }
            }

            // update indexes with the new table name
            if (!CommonHelper.isEmpty(clone.indexes)) {
                for (let i = 0; i < clone.indexes.length; i++) {
                    const parsed = CommonHelper.parseIndex(clone.indexes[i]);
                    parsed.indexName = "idx_" + CommonHelper.randomString(10);
                    parsed.tableName = clone.name;
                    clone.indexes[i] = CommonHelper.buildIndex(parsed);
                }
            }
        }

        show(clone);

        await tick();

        initialFormHash = "";
    }

    function hasOtherKeys(obj, excludes = []) {
        if (CommonHelper.isEmpty(obj)) {
            return false;
        }

        const errorKeys = Object.keys(obj);
        for (let key of errorKeys) {
            if (!excludes.includes(key)) {
                return true;
            }
        }

        return false;
    }
</script>

<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
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
            {!collection.id ? "New collection" : "Edit collection"}
        </h4>

        {#if !!collection.id && (!collection.system || !isView)}
            <div class="flex-fill" />
            <div
                tabindex="0"
                role="button"
                aria-label="More collection options"
                class="btn btn-sm btn-circle btn-transparent flex-gap-0"
            >
                <i class="ri-more-line" aria-hidden="true" />
                <Toggler class="dropdown dropdown-right m-t-5">
                    {#if !collection.system}
                        <button
                            type="button"
                            class="dropdown-item"
                            role="menuitem"
                            on:click={() => duplicateConfirm()}
                        >
                            <i class="ri-file-copy-line" aria-hidden="true" />
                            <span class="txt">Duplicate</span>
                        </button>
                        <hr />
                    {/if}
                    {#if !isView}
                        <button
                            type="button"
                            class="dropdown-item txt-danger"
                            role="menuitem"
                            on:click={() => truncateConfirm()}
                        >
                            <i class="ri-eraser-line" aria-hidden="true"></i>
                            <span class="txt">Truncate</span>
                        </button>
                    {/if}
                    {#if !collection.system}
                        <button
                            type="button"
                            class="dropdown-item txt-danger"
                            role="menuitem"
                            on:click|preventDefault|stopPropagation={() => deleteConfirm()}
                        >
                            <i class="ri-delete-bin-7-line" aria-hidden="true" />
                            <span class="txt">Delete</span>
                        </button>
                    {/if}
                </Toggler>
            </div>
        {/if}

        <form
            class="block"
            on:submit|preventDefault={() => {
                canSave && saveConfirm();
            }}
        >
            <Field class="form-field collection-field-name required m-b-0" name="name" let:uniqueId>
                <label for={uniqueId}>Name</label>

                <!-- svelte-ignore a11y-autofocus -->
                <input
                    type="text"
                    id={uniqueId}
                    required
                    disabled={isSystemUpdate}
                    spellcheck="false"
                    class:txt-bold={collection.system}
                    autofocus={!collection.id}
                    placeholder={isAuth ? `eg. "users"` : `eg. "posts"`}
                    value={collection.name}
                    on:input={(e) => {
                        collection.name = CommonHelper.slugify(e.target.value);
                        e.target.value = collection.name;
                    }}
                />

                <div class="form-field-addon">
                    <div
                        tabindex={!collection.id ? 0 : -1}
                        role={!collection.id ? "button" : ""}
                        aria-label="View types"
                        class="btn btn-sm p-r-10 p-l-10 {!collection.id ? 'btn-outline' : 'btn-transparent'}"
                        class:btn-disabled={!!collection.id}
                    >
                        <!-- empty span for alignment -->
                        <span aria-hidden="true" />
                        <span class="txt">Type: {collectionTypes[collection.type] || "N/A"}</span>
                        {#if !collection.id}
                            <i class="ri-arrow-down-s-fill" aria-hidden="true" />
                            <Toggler class="dropdown dropdown-right dropdown-nowrap m-t-5">
                                {#each Object.entries(collectionTypes) as [type, label]}
                                    <button
                                        type="button"
                                        role="menuitem"
                                        class="dropdown-item closable"
                                        class:selected={type == collection.type}
                                        on:click={() => setCollectionType(type)}
                                    >
                                        <i
                                            class={CommonHelper.getCollectionTypeIcon(type)}
                                            aria-hidden="true"
                                        />
                                        <span class="txt">{label} collection</span>
                                    </button>
                                {/each}
                            </Toggler>
                        {/if}
                    </div>
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
                <span class="txt">{isView ? "Query" : "Fields"}</span>
                {#if !CommonHelper.isEmpty(fieldsTabError)}
                    <i
                        class="ri-error-warning-fill txt-danger"
                        transition:scale={{ duration: 150, start: 0.7 }}
                        use:tooltip={fieldsTabError}
                    />
                {/if}
            </button>

            {#if !isSuperusers}
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === TAB_RULES}
                    on:click={() => changeTab(TAB_RULES)}
                >
                    <span class="txt">API Rules</span>
                    {#if !CommonHelper.isEmpty($errors?.listRule) || !CommonHelper.isEmpty($errors?.viewRule) || !CommonHelper.isEmpty($errors?.createRule) || !CommonHelper.isEmpty($errors?.updateRule) || !CommonHelper.isEmpty($errors?.deleteRule) || !CommonHelper.isEmpty($errors?.authRule) || !CommonHelper.isEmpty($errors?.manageRule)}
                        <i
                            class="ri-error-warning-fill txt-danger"
                            transition:scale={{ duration: 150, start: 0.7 }}
                            use:tooltip={"Has errors"}
                        />
                    {/if}
                </button>
            {/if}

            {#if isAuth}
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === TAB_OPTIONS}
                    on:click={() => changeTab(TAB_OPTIONS)}
                >
                    <span class="txt">Options</span>
                    {#if $errors && hasOtherKeys($errors, baseCollectionKeys.concat( ["manageRule", "authRule"], ))}
                        <i
                            class="ri-error-warning-fill txt-danger"
                            transition:scale={{ duration: 150, start: 0.7 }}
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
            {#if isView}
                <CollectionQueryTab bind:collection />
            {:else}
                <CollectionFieldsTab bind:collection />
            {/if}
        </div>

        {#if !isSuperusers && activeTab === TAB_RULES}
            <div class="tab-item active">
                <CollectionRulesTab bind:collection />
            </div>
        {/if}

        {#if isAuth}
            <div class="tab-item" class:active={activeTab === TAB_OPTIONS}>
                <CollectionAuthOptionsTab bind:collection />
            </div>
        {/if}
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>

        <div class="btns-group no-gap">
            <button
                type="button"
                title="Save and close"
                class="btn"
                class:btn-expanded={!collection.id}
                class:btn-expanded-sm={!!collection.id}
                class:btn-loading={isSaving || isLoadingConfirmation}
                disabled={!canSave || isSaving || isLoadingConfirmation}
                on:click={() => saveConfirm()}
            >
                <span class="txt">{!collection.id ? "Create" : "Save changes"}</span>
            </button>

            {#if collection.id}
                <button
                    type="button"
                    class="btn p-l-5 p-r-5 flex-gap-0"
                    disabled={!canSave || isSaving || isLoadingConfirmation}
                >
                    <i class="ri-arrow-down-s-line" aria-hidden="true"></i>

                    <Toggler class="dropdown dropdown-upside dropdown-right dropdown-nowrap m-b-5">
                        <button
                            type="button"
                            class="dropdown-item closable"
                            role="menuitem"
                            on:click={() => saveConfirm(false)}
                        >
                            <span class="txt">Save and continue</span>
                        </button>
                    </Toggler>
                </button>
            {/if}
        </div>
    </svelte:fragment>
</OverlayPanel>

<CollectionUpdateConfirm bind:this={confirmChangesPanel} on:confirm={(e) => save(e.detail)} />

<style>
    .upsert-panel-title {
        display: inline-flex;
        align-items: center;
        min-height: var(--smBtnHeight);
    }
    .tabs-content:focus-within {
        z-index: 9; /* autocomplete dropdown overlay fix */
    }
    :global(.collection-panel .panel-content) {
        scrollbar-gutter: stable;
        padding-right: calc(var(--baseSpacing) - 5px);
    }
</style>
