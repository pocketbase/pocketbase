<script>
    import { createEventDispatcher, tick } from "svelte";
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
    const TAB_FORM = "form";
    const TAB_PROVIDERS = "providers";

    export let collection;

    let recordPanel;
    let original = null;
    let record = new Record();
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let uploadedFilesMap = {}; // eg.: {"field1":[File1, File2], ...}
    let deletedFileIndexesMap = {}; // eg.: {"field1":[0, 1], ...}
    let initialFormHash = "";
    let activeTab = TAB_FORM;
    let isNew = true;

    $: hasEditorField = !!collection?.schema?.find((f) => f.type === "editor");

    $: hasFileChanges =
        CommonHelper.hasNonEmptyProps(uploadedFilesMap) ||
        CommonHelper.hasNonEmptyProps(deletedFileIndexesMap);

    $: hasChanges = hasFileChanges || initialFormHash != calculateFormHash(record);

    $: isNew = !original || original.isNew;

    $: canSave = isNew || hasChanges;

    export function show(model) {
        load(model);

        confirmClose = true;

        activeTab = TAB_FORM;

        return recordPanel?.show();
    }

    export function hide() {
        return recordPanel?.hide();
    }

    async function load(model) {
        setErrors({}); // reset errors
        original = model || new Record();
        if (model?.$clone) {
            record = model.$clone();
        } else {
            record = new Record();
        }
        uploadedFilesMap = {};
        deletedFileIndexesMap = {};
        await tick(); // wait to populate the fields to get the normalized values
        initialFormHash = calculateFormHash(record);
    }

    function calculateFormHash(m) {
        return JSON.stringify(m);
    }

    function save() {
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
                confirmClose = false;
                hide();
                dispatch("save", result);
            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
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
                    ApiClient.errorResponseHandler(err);
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
        for (const key in deletedFileIndexesMap) {
            const indexes = CommonHelper.toArray(deletedFileIndexesMap[key]);
            for (const index of indexes) {
                formData.append(key + "." + index, "");
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
                    ApiClient.errorResponseHandler(err);
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
                    ApiClient.errorResponseHandler(err);
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

        show(clone);

        await tick();

        initialFormHash = "";
    }
</script>

<OverlayPanel
    bind:this={recordPanel}
    class="
        record-panel
        {hasEditorField ? 'overlay-panel-xl' : 'overlay-panel-lg'}
        {collection?.isAuth && !isNew ? 'colored-header' : ''}
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
        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4>
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
                    class:active={activeTab === TAB_FORM}
                    on:click={() => (activeTab = TAB_FORM)}
                >
                    Account
                </button>
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === TAB_PROVIDERS}
                    on:click={() => (activeTab = TAB_PROVIDERS)}
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
            class:active={activeTab === TAB_FORM}
            on:submit|preventDefault={save}
        >
            <Field class="form-field readonly" name="id" let:uniqueId>
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
                    bind:value={record.id}
                    readonly={!isNew}
                />
            </Field>

            {#if collection?.isAuth}
                <AuthFields bind:record {collection} />

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
                        bind:deletedFileIndexes={deletedFileIndexesMap[field.name]}
                    />
                {:else if field.type === "relation"}
                    <RelationField {field} bind:value={record[field.name]} />
                {/if}
            {/each}
        </form>

        {#if collection.$isAuth && !isNew}
            <div class="tab-item" class:active={activeTab === TAB_PROVIDERS}>
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
