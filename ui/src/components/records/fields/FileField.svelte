<script>
    import { SchemaField } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import UploadedFilePreview from "@/components/base/UploadedFilePreview.svelte";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";

    export let record;
    export let value = "";
    export let uploadedFiles = []; // Array<File> array
    export let deletedFileIndexes = []; // Array<int> array
    export let field = new SchemaField();

    let fileInput;
    let filesListElem;

    // normalize uploadedFiles type
    $: if (!Array.isArray(uploadedFiles)) {
        uploadedFiles = CommonHelper.toArray(uploadedFiles);
    }

    // normalize delited file indexes
    $: if (!Array.isArray(deletedFileIndexes)) {
        deletedFileIndexes = CommonHelper.toArray(deletedFileIndexes);
    }

    $: isMultiple = field.options?.maxSelect > 1;

    $: if (CommonHelper.isEmpty(value)) {
        value = isMultiple ? [] : "";
    }

    $: valueAsArray = CommonHelper.toArray(value);

    $: maxReached =
        (valueAsArray.length || uploadedFiles.length) &&
        field.options?.maxSelect <= valueAsArray.length + uploadedFiles.length - deletedFileIndexes.length;

    $: if (uploadedFiles !== -1 || deletedFileIndexes !== -1) {
        triggerListChange();
    }

    function restoreExistingFile(valueIndex) {
        CommonHelper.removeByValue(deletedFileIndexes, valueIndex);
        deletedFileIndexes = deletedFileIndexes;
    }

    function removeExistingFile(valueIndex) {
        CommonHelper.pushUnique(deletedFileIndexes, valueIndex);
        deletedFileIndexes = deletedFileIndexes;
    }

    function removeNewFile(index) {
        if (!CommonHelper.isEmpty(uploadedFiles[index])) {
            uploadedFiles.splice(index, 1);
        }
        uploadedFiles = uploadedFiles;
    }

    // emulate native change event
    function triggerListChange() {
        filesListElem?.dispatchEvent(
            new CustomEvent("change", {
                detail: { value, uploadedFiles, deletedFileIndexes },
                bubbles: true,
            })
        );
    }
</script>

<Field
    class="form-field form-field-list form-field-file {field.required ? 'required' : ''}"
    name={field.name}
    let:uniqueId
>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>

    <div bind:this={filesListElem} class="list">
        {#each valueAsArray as filename, i (filename + record.id)}
            {@const isDeleted = deletedFileIndexes.includes(i)}
            <div class="list-item">
                <div class:fade={deletedFileIndexes.includes(i)}>
                    <RecordFileThumb {record} {filename} />
                </div>

                <div class="content">
                    <a
                        href={ApiClient.getFileUrl(record, filename)}
                        class="txt-ellipsis {isDeleted ? 'txt-strikethrough txt-hint' : 'link-primary'}"
                        title="Download"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        {filename}
                    </a>
                </div>

                <div class="actions">
                    {#if deletedFileIndexes.includes(i)}
                        <button
                            type="button"
                            class="btn btn-sm btn-danger btn-transparent"
                            on:click={() => restoreExistingFile(i)}
                        >
                            <span class="txt">Restore</span>
                        </button>
                    {:else}
                        <button
                            type="button"
                            class="btn btn-transparent btn-hint btn-sm btn-circle btn-remove"
                            use:tooltip={"Remove file"}
                            on:click={() => removeExistingFile(i)}
                        >
                            <i class="ri-close-line" />
                        </button>
                    {/if}
                </div>
            </div>
        {/each}

        {#each uploadedFiles as file, i}
            <div class="list-item">
                <figure class="thumb">
                    <UploadedFilePreview {file} />
                </figure>
                <div class="filename m-r-auto" title={file.name}>
                    <small class="label label-success m-r-5">New</small>
                    <span class="txt">{file.name}</span>
                </div>
                <button
                    type="button"
                    class="btn btn-transparent btn-sm btn-circle btn-remove"
                    use:tooltip={"Remove file"}
                    on:click={() => removeNewFile(i)}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
        {/each}

        <div class="list-item list-item-btn">
            <input
                bind:this={fileInput}
                type="file"
                class="hidden"
                multiple={isMultiple}
                on:change={() => {
                    for (let file of fileInput.files) {
                        uploadedFiles.push(file);
                    }
                    uploadedFiles = uploadedFiles;
                    fileInput.value = null; // reset
                }}
            />
            <button
                type="button"
                class="btn btn-transparent btn-sm btn-block"
                disabled={maxReached}
                on:click={() => fileInput?.click()}
            >
                <i class="ri-upload-cloud-line" />
                <span class="txt">Upload new file</span>
            </button>
        </div>
    </div>
</Field>
