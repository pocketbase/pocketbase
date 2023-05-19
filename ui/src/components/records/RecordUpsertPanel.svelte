<script>
    import { createEventDispatcher, tick } from "svelte";
    import { slide } from "svelte/transition";
    import { Record } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { setErrors } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import AuthFields from "@/components/records/fields/AuthFields.svelte";
    import TextField from "@/components/records/fields/TextField.svelte";
    import NumberField from "@/components/records/fields/NumberField.svelte";
    import BoolField from "@/components/records/fields/BoolField.svelte";
    import EmailField from "@/components/records/fields/EmailField.svelte";
    import UrlField from "@/components/records/fields/UrlField.svelte";
    import DateField from "@/components/records/fields/DateField.svelte";
    import SelectField from "@/components/records/fields/SelectField.svelte";
    import JsonField from "@/components/records/fields/JsonField.svelte";
    import FileField from "@/components/records/fields/FileField.svelte";
    import RelationField from "@/components/records/fields/RelationField.svelte";
    import EditorField from "@/components/records/fields/EditorField.svelte";
    import ExternalAuthsList from "@/components/records/ExternalAuthsList.svelte";

    const dispatch = createEventDispatcher();
    const formId = "record_" + CommonHelper.randomString(5);
    const tabFormKey = "form";
    const tabProviderKey = "providers";

    export let collection;

    let recordPanel;
    let original = null;
    let record = null;
    let initialDraft = null;
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let uploadedFilesMap = {}; // eg.: {"field1":[File1, File2], ...}
    let deletedFileNamesMap = {}; // eg.: {"field1":[0, 1], ...}
    let originalSerializedData = JSON.stringify(null);
    let serializedData = originalSerializedData;
    let activeTab = tabFormKey;
    let isNew = true;
    let isLoaded = false;

    $: hasEditorField = !!collection?.schema?.find((f) => f.type === "editor");

    $: hasFileChanges =
        CommonHelper.hasNonEmptyProps(uploadedFilesMap) || CommonHelper.hasNonEmptyProps(deletedFileNamesMap);

    $: serializedData = JSON.stringify(record);

    $: hasChanges = hasFileChanges || originalSerializedData != serializedData;

    $: isNew = !original || original.$isNew;

    $: canSave = isNew || hasChanges;

    $: if (isLoaded) {
        updateDraft(serializedData);
    }

    export function show(model) {
        load(model);

        confirmClose = true;

        activeTab = tabFormKey;

        return recordPanel?.show();
    }

    export function hide() {
        return recordPanel?.hide();
    }

    async function load(model) {
        isLoaded = false;
        setErrors({}); // reset errors
        original = model || new Record();
        record = original.$clone();
        uploadedFilesMap = {};
        deletedFileNamesMap = {};

        // wait to populate the fields to get the normalized values
        await tick();

        initialDraft = getDraft();
        if (!initialDraft || areRecordsEqual(record, initialDraft)) {
            initialDraft = null;
        } else {
            delete initialDraft.password;
            delete initialDraft.passwordConfirm;
        }

        originalSerializedData = JSON.stringify(record);
        isLoaded = true;
    }

    async function replaceOriginal(newOriginal) {
        setErrors({}); // reset errors
        original = newOriginal || new Record();
        uploadedFilesMap = {};
        deletedFileNamesMap = {};

        // to avoid layout shifts we replace only the file and non-schema fields
        const skipFields = collection?.schema?.filter((f) => f.type != "file")?.map((f) => f.name) || [];
        for (let k in newOriginal.$export()) {
            if (skipFields.includes(k)) {
                continue;
            }
            record[k] = newOriginal[k];
        }

        // wait to populate the fields to get the normalized values
        await tick();

        originalSerializedData = JSON.stringify(record);

        deleteDraft();
    }

    function draftKey() {
        return "record_draft_" + (collection?.id || "") + "_" + (original?.id || "");
    }

    function getDraft(fallbackRecord) {
        try {
            const raw = window.localStorage.getItem(draftKey());
            if (raw) {
                return new Record(JSON.parse(raw));
            }
        } catch (_) {}

        return fallbackRecord;
    }

    function updateDraft(newSerializedData) {
        window.localStorage.setItem(draftKey(), newSerializedData);
    }

    function restoreDraft() {
        if (initialDraft) {
            record = initialDraft;
            initialDraft = null;
        }
    }

    function areRecordsEqual(recordA, recordB) {
        const cloneA = recordA?.$clone();
        const cloneB = recordB?.$clone();

        const fileFields = collection?.schema?.filter((f) => f.type === "file");
        for (let field of fileFields) {
            delete cloneA?.[field.name];
            delete cloneB?.[field.name];
        }

        // delete password props
        delete cloneA?.password;
        delete cloneA?.passwordConfirm;
        delete cloneB?.password;
        delete cloneB?.passwordConfirm;

        return JSON.stringify(cloneA) == JSON.stringify(cloneB);
    }

    function deleteDraft() {
        initialDraft = null;
        window.localStorage.removeItem(draftKey());
    }

    function save(hidePanel = true) {
        if (isSaving || !canSave || !collection?.id) {
            return;
        }

        isSaving = true;

        const data = exportFormData();

        let request;
        if (isNew) {
            request = ApiClient.collection(collection.id).create(data);
        } else {
            request = ApiClient.collection(collection.id).update(record.id, data);
        }

        request
            .then((result) => {
                addSuccessToast(isNew ? "Successfully created record." : "Successfully updated record.");

                deleteDraft();

                if (hidePanel) {
                    confirmClose = false;
                    hide();
                } else {
                    replaceOriginal(result);
                }

                dispatch("save", result);
            })
            .catch((err) => {
                ApiClient.error(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function deleteConfirm() {
        if (!original?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete the selected record?`, () => {
            return ApiClient.collection(original.collectionId)
                .delete(original.id)
                .then(() => {
                    hide();
                    addSuccessToast("Successfully deleted record.");
                    dispatch("delete", original);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function exportFormData() {
        const data = record?.$export() || {};
        const formData = new FormData();

        const exportableFields = {
            id: data.id,
        };

        for (const field of collection?.schema || []) {
            exportableFields[field.name] = true;
        }

        if (collection?.isAuth) {
            exportableFields["username"] = true;
            exportableFields["email"] = true;
            exportableFields["emailVisibility"] = true;
            exportableFields["password"] = true;
            exportableFields["passwordConfirm"] = true;
            exportableFields["verified"] = true;
        }

        // export base fields
        for (const key in data) {
            // skip non-schema fields
            if (!exportableFields[key]) {
                continue;
            }

            // normalize nullable values
            if (typeof data[key] === "undefined") {
                data[key] = null;
            }

            CommonHelper.addValueToFormData(formData, key, data[key]);
        }

        // add uploaded files  (if any)
        for (const key in uploadedFilesMap) {
            const files = CommonHelper.toArray(uploadedFilesMap[key]);
            for (const file of files) {
                formData.append(key, file);
            }
        }

        // unset deleted files (if any)
        for (const key in deletedFileNamesMap) {
            const names = CommonHelper.toArray(deletedFileNamesMap[key]);
            for (const name of names) {
                formData.append(key + "." + name, "");
            }
        }

        return formData;
    }

    function sendVerificationEmail() {
        if (!collection?.id || !original?.email) {
            return;
        }

        confirm(`Do you really want to sent verification email to ${original.email}?`, () => {
            return ApiClient.collection(collection.id)
                .requestVerification(original.email)
                .then(() => {
                    addSuccessToast(`Successfully sent verification email to ${original.email}.`);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function sendPasswordResetEmail() {
        if (!collection?.id || !original?.email) {
            return;
        }

        confirm(`Do you really want to sent password reset email to ${original.email}?`, () => {
            return ApiClient.collection(collection.id)
                .requestPasswordReset(original.email)
                .then(() => {
                    addSuccessToast(`Successfully sent password reset email to ${original.email}.`);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
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

            // reset file fields
            const fields = collection?.schema || [];
            for (const field of fields) {
                if (field.type === "file") {
                    delete clone[field.name];
                }
            }
        }

        deleteDraft();
        show(clone);

        await tick();

        originalSerializedData = "";
    }

    function handleFormKeydown(e) {
        if ((e.ctrlKey || e.metaKey) && e.code == "KeyS") {
            e.preventDefault();
            e.stopPropagation();
            save(false);
        }
    }
</script>

<OverlayPanel
    bind:this={recordPanel}
    class="
        record-panel
        {hasEditorField ? 'overlay-panel-xl' : 'overlay-panel-lg'}
        {collection?.$isAuth && !isNew ? 'colored-header' : ''}
    "
    beforeHide={() => {
        if (hasChanges && confirmClose) {
            confirm("You have unsaved changes. Do you really want to close the panel?", () => {
                confirmClose = false;
                hide();
            });

            return false;
        }

        setErrors({});
        deleteDraft();

        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4 class="panel-title">
            {isNew ? "New" : "Edit"}
            <strong>{collection?.name}</strong> record
        </h4>

        {#if !isNew}
            <div class="flex-fill" />
            <button type="button" aria-label="More" class="btn btn-sm btn-circle btn-transparent flex-gap-0">
                <i class="ri-more-line" />
                <Toggler class="dropdown dropdown-right dropdown-nowrap">
                    {#if collection.$isAuth && !original.verified && original.email}
                        <button
                            type="button"
                            class="dropdown-item closable"
                            on:click={() => sendVerificationEmail()}
                        >
                            <i class="ri-mail-check-line" />
                            <span class="txt">Send verification email</span>
                        </button>
                    {/if}
                    {#if collection.$isAuth && original.email}
                        <button
                            type="button"
                            class="dropdown-item closable"
                            on:click={() => sendPasswordResetEmail()}
                        >
                            <i class="ri-mail-lock-line" />
                            <span class="txt">Send password reset email</span>
                        </button>
                    {/if}
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

        {#if collection.$isAuth && !isNew}
            <div class="tabs-header stretched">
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === tabFormKey}
                    on:click={() => (activeTab = tabFormKey)}
                >
                    Account
                </button>
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === tabProviderKey}
                    on:click={() => (activeTab = tabProviderKey)}
                >
                    Authorized providers
                </button>
            </div>
        {/if}
    </svelte:fragment>

    <div class="tabs-content">
        <form
            id={formId}
            class="tab-item"
            class:active={activeTab === tabFormKey}
            on:submit|preventDefault={save}
            on:keydown={handleFormKeydown}
        >
            {#if !hasChanges && initialDraft}
                <div class="block" out:slide={{ duration: 150 }}>
                    <div class="alert alert-info m-0">
                        <div class="icon">
                            <i class="ri-information-line" />
                        </div>
                        <div class="flex flex-gap-xs">
                            The record has previous unsaved changes.
                            <button
                                type="button"
                                class="btn btn-sm btn-secondary"
                                on:click={() => restoreDraft()}
                            >
                                Restore draft
                            </button>
                        </div>
                        <button
                            type="button"
                            class="close"
                            aria-label="Discard draft"
                            use:tooltip={"Discard draft"}
                            on:click|preventDefault={() => deleteDraft()}
                        >
                            <i class="ri-close-line" />
                        </button>
                    </div>
                    <div class="clearfix p-b-base" />
                </div>
            {/if}

            <Field class="form-field {!isNew ? 'readonly' : ''}" name="id" let:uniqueId>
                <label for={uniqueId}>
                    <i class={CommonHelper.getFieldTypeIcon("primary")} />
                    <span class="txt">id</span>
                    <span class="flex-fill" />
                </label>
                {#if !isNew}
                    <div class="form-field-addon">
                        <i
                            class="ri-calendar-event-line txt-disabled"
                            use:tooltip={{
                                text: `Created: ${record.created}\nUpdated: ${record.updated}`,
                                position: "left",
                            }}
                        />
                    </div>
                {/if}
                <input
                    type="text"
                    id={uniqueId}
                    placeholder="Leave empty to auto generate..."
                    minlength="15"
                    readonly={!isNew}
                    bind:value={record.id}
                />
            </Field>

            {#if collection?.isAuth}
                <AuthFields bind:record {isNew} {collection} />

                {#if collection?.schema?.length}
                    <hr />
                {/if}
            {/if}

            {#each collection?.schema || [] as field (field.name)}
                {#if field.type === "text"}
                    <TextField {field} bind:value={record[field.name]} />
                {:else if field.type === "number"}
                    <NumberField {field} bind:value={record[field.name]} />
                {:else if field.type === "bool"}
                    <BoolField {field} bind:value={record[field.name]} />
                {:else if field.type === "email"}
                    <EmailField {field} bind:value={record[field.name]} />
                {:else if field.type === "url"}
                    <UrlField {field} bind:value={record[field.name]} />
                {:else if field.type === "editor"}
                    <EditorField {field} bind:value={record[field.name]} />
                {:else if field.type === "date"}
                    <DateField {field} bind:value={record[field.name]} />
                {:else if field.type === "select"}
                    <SelectField {field} bind:value={record[field.name]} />
                {:else if field.type === "json"}
                    <JsonField {field} bind:value={record[field.name]} />
                {:else if field.type === "file"}
                    <FileField
                        {field}
                        {record}
                        bind:value={record[field.name]}
                        bind:uploadedFiles={uploadedFilesMap[field.name]}
                        bind:deletedFileNames={deletedFileNamesMap[field.name]}
                    />
                {:else if field.type === "relation"}
                    <RelationField {field} bind:value={record[field.name]} />
                {/if}
            {/each}
        </form>

        {#if collection.$isAuth && !isNew}
            <div class="tab-item" class:active={activeTab === tabProviderKey}>
                <ExternalAuthsList {record} />
            </div>
        {/if}
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>

        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!canSave || isSaving}
        >
            <span class="txt">{isNew ? "Create" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style>
    .panel-title {
        line-height: var(--smBtnHeight);
    }
</style>
