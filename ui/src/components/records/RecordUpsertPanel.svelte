<script>
    import { createEventDispatcher, tick } from "svelte";
    import { Record } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import { setErrors } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
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
    import UserField from "@/components/records/fields/UserField.svelte";

    const dispatch = createEventDispatcher();
    const formId = "record_" + CommonHelper.randomString(5);

    export let collection;

    let recordPanel;
    let original = null;
    let record = new Record();
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let uploadedFilesMap = {}; // eg.: {"field1":[File1, File2], ...}
    let deletedFileIndexesMap = {}; // eg.: {"field1":[0, 1], ...}
    let initialFormHash = "";

    $: hasFileChanges =
        CommonHelper.hasNonEmptyProps(uploadedFilesMap) ||
        CommonHelper.hasNonEmptyProps(deletedFileIndexesMap);

    $: hasChanges = hasFileChanges || initialFormHash != calculateFormHash(record);

    $: canSave = record.isNew || hasChanges;

    $: isProfileCollection = collection?.name !== import.meta.env.PB_PROFILE_COLLECTION;

    export function show(model) {
        load(model);

        confirmClose = true;

        return recordPanel?.show();
    }

    export function hide() {
        return recordPanel?.hide();
    }

    async function load(model) {
        setErrors({}); // reset errors
        original = model || {};
        record = model?.clone ? model.clone() : new Record();
        uploadedFilesMap = {};
        deletedFileIndexesMap = {};
        await tick(); // wait to populate the fields to get the normalized values
        initialFormHash = calculateFormHash(record);
    }

    function calculateFormHash(m) {
        return JSON.stringify(m);
    }

    function save() {
        if (isSaving || !canSave) {
            return;
        }

        isSaving = true;

        const data = exportFormData();

        let request;
        if (record.isNew) {
            request = ApiClient.records.create(collection?.id, data);
        } else {
            request = ApiClient.records.update(collection?.id, record.id, data);
        }

        request
            .then(async (result) => {
                addSuccessToast(
                    record.isNew ? "Successfully created record." : "Successfully updated record."
                );
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
            return ApiClient.records.delete(original["@collectionId"], original.id)
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
        const data = record?.export() || {};
        const formData = new FormData();

        const schemaMap = {};
        for (const field of collection?.schema || []) {
            schemaMap[field.name] = field;
        }

        // export base fields
        for (const key in data) {
            // skip non-schema fields
            if (!schemaMap[key]) {
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
</script>

<OverlayPanel
    bind:this={recordPanel}
    class="overlay-panel-lg record-panel"
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
        <h4>
            {record.isNew ? "New" : "Edit"}
            {collection.name} record
        </h4>

        {#if !record.isNew && isProfileCollection}
            <div class="flex-fill" />
            <button type="button" class="btn btn-sm btn-circle btn-secondary">
                <div class="content">
                    <i class="ri-more-line" />
                    <Toggler class="dropdown dropdown-right m-t-5">
                        <div tabindex="0" class="dropdown-item closable" on:click={() => deleteConfirm()}>
                            <i class="ri-delete-bin-7-line" />
                            <span class="txt">Delete</span>
                        </div>
                    </Toggler>
                </div>
            </button>
        {/if}
    </svelte:fragment>

    <form id={formId} class="block" on:submit|preventDefault={save}>
        {#if !record.isNew}
            <Field class="form-field disabled" name="id" let:uniqueId>
                <label for={uniqueId}>
                    <i class={CommonHelper.getFieldTypeIcon("primary")} />
                    <span class="txt">id</span>
                    <span class="flex-fill" />
                </label>
                <input type="text" id={uniqueId} value={record.id} disabled />
            </Field>
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
            {:else if field.type === "user"}
                <UserField {field} bind:value={record[field.name]} />
            {/if}
        {:else}
            <div class="block txt-center txt-disabled">
                <h5>No custom fields to be set</h5>
            </div>
        {/each}
    </form>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!canSave || isSaving}
        >
            <span class="txt">{record.isNew ? "Create" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
