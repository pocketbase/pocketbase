<script>
    import { SchemaField } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import UploadedFilePreview from "@/components/base/UploadedFilePreview.svelte";
    import PreviewPopup from "@/components/base/PreviewPopup.svelte";
    import RecordFilePreview from "@/components/records/RecordFilePreview.svelte";

    export let record;
    export let value = null;
    export let uploadedFiles = []; // Array<File> array
    export let deletedFileIndexes = []; // Array<int> array
    export let field = new SchemaField();

    let fileInput;
    let previewPopup;
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

    $: if (typeof value === "undefined" || value === null) {
        value = isMultiple ? [] : null;
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

<Field class="form-field form-field-file {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>

    <div bind:this={filesListElem} class="files-list">
        {#each valueAsArray as filename, i (filename)}
            <div class="list-item">
                <figute
                    class="thumb"
                    class:fade={deletedFileIndexes.includes(i)}
                    class:link-fade={CommonHelper.hasImageExtension(filename)}
                    title={CommonHelper.hasImageExtension(filename) ? "Preview" : ""}
                    on:click={() =>
                        CommonHelper.hasImageExtension(filename)
                            ? previewPopup?.show(ApiClient.Records.getFileUrl(record, filename))
                            : false}
                >
                    <RecordFilePreview {record} {filename} />
                </figute>
                <a
                    href={ApiClient.Records.getFileUrl(record, filename)}
                    class="filename"
                    class:txt-strikethrough={deletedFileIndexes.includes(i)}
                    title={"Download " + filename}
                    target="_blank"
                    rel="noopener"
                    download
                >
                    /.../{filename}
                </a>

                {#if deletedFileIndexes.includes(i)}
                    <button
                        type="button"
                        class="btn btn-sm btn-danger btn-secondary"
                        on:click={() => restoreExistingFile(i)}
                    >
                        <span class="txt">Restore</span>
                    </button>
                {:else}
                    <button
                        type="button"
                        class="btn btn-secondary btn-sm btn-circle btn-remove txt-hint"
                        use:tooltip={"Remove file"}
                        on:click={() => removeExistingFile(i)}
                    >
                        <i class="ri-close-line" />
                    </button>
                {/if}
            </div>
        {/each}

        {#each uploadedFiles as file, i}
            <div class="list-item">
                <figute class="thumb">
                    <UploadedFilePreview {file} />
                </figute>
                <div class="filename" title={file.name}>
                    <small class="label label-success m-r-5">New</small>
                    <span class="txt">{file.name}</span>
                </div>
                <button
                    type="button"
                    class="btn btn-secondary btn-sm btn-circle btn-remove"
                    use:tooltip={"Remove file"}
                    on:click={() => removeNewFile(i)}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
        {/each}

        {#if !maxReached}
            <div class="list-item btn-list-item">
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
                    class="btn btn-secondary btn-sm btn-block"
                    on:click={() => fileInput?.click()}
                >
                    <i class="ri-upload-cloud-line" />
                    <span class="txt">Upload new file</span>
                </button>
            </div>
        {/if}
    </div>
</Field>

<PreviewPopup bind:this={previewPopup} />
