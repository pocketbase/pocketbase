<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addSuccessToast } from "@/stores/toasts";
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
    import { Record } from "pocketbase";

    const recordId = "default";

    export let collection;

    let record = null;
    let uploadedFilesMap = {};
    let deletedFileNamesMap = {};
    let originalSerializedData = JSON.stringify(null);
    let serializedData = originalSerializedData;
    let isLoading = false;
    let hasChanges = false;
    let isSaving = false;

    $: if (collection?.id) {
        load();
    }

    $: hasFileChanges =
        CommonHelper.hasNonEmptyProps(uploadedFilesMap) || CommonHelper.hasNonEmptyProps(deletedFileNamesMap);

    $: serializedData = JSON.stringify(record);

    $: hasChanges = hasFileChanges || originalSerializedData != serializedData;

    function load() {
        record = null;
        isLoading = true;

        ApiClient.collection(collection.id)
            .getOne(recordId, {
                $cancelKey: "records_list",
            })
            .then(async (result) => {
                setRecord(result);
                isLoading = false;
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoading = false;
                    console.warn(err);
                    ApiClient.error(err, false);
                }
            });
    }

    function save() {
        if (isSaving || !hasChanges || !collection?.id) {
            return;
        }

        isSaving = true;

        const data = exportFormData();

        let request = ApiClient.collection(collection.id).update(recordId, data);

        request
            .then((result) => {
                setRecord(result);
                addSuccessToast("Successfully updated record.");
            })
            .catch((err) => {
                ApiClient.error(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function reset() {
        record = new Record(JSON.parse(originalSerializedData));
        uploadedFilesMap = {};
        deletedFileNamesMap = {};
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

    function setRecord(original) {
        record = original.$clone();
        uploadedFilesMap = {};
        deletedFileNamesMap = {};
        originalSerializedData = JSON.stringify(record);
    }
</script>

<div class="wrapper">
    <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
        {#if isLoading}
            <div class="loader" />
        {:else if record}
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
            <div class="flex">
                <div class="flex-fill" />
                {#if hasChanges}
                    <button
                        type="button"
                        class="btn btn-transparent btn-hint"
                        disabled={isSaving}
                        on:click={() => reset()}
                    >
                        <span class="txt">Cancel</span>
                    </button>
                {/if}
                <button
                    type="submit"
                    class="btn btn-expanded"
                    class:btn-loading={isSaving}
                    disabled={!hasChanges || isSaving}
                    on:click={() => save()}
                >
                    <span class="txt">Save changes</span>
                </button>
            </div>
        {:else}
            <div class="alert alert-danger m-0">
                <div class="icon">
                    <i class="ri-alert-line" />
                </div>
                <div class="content">Default record was not found.</div>
            </div>
        {/if}
    </form>
</div>
