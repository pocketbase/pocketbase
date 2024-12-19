<script>
    import tooltip from "@/actions/tooltip";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import Draggable from "@/components/base/Draggable.svelte";
    import Field from "@/components/base/Field.svelte";
    import UploadedFilePreview from "@/components/base/UploadedFilePreview.svelte";
    import RecordFileThumb from "@/components/records/RecordFileThumb.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let record;
    export let field;
    export let value = "";
    export let uploadedFiles = []; // Array<File> array
    export let deletedFileNames = []; // Array<string> array

    let fileInput;
    let filesListElem;
    let isDragOver = false;

    // normalize uploadedFiles type
    $: if (!Array.isArray(uploadedFiles)) {
        uploadedFiles = CommonHelper.toArray(uploadedFiles);
    }

    // normalize deleted files
    $: if (!Array.isArray(deletedFileNames)) {
        deletedFileNames = CommonHelper.toArray(deletedFileNames);
    }

    $: isMultiple = field.maxSelect > 1;

    $: if (CommonHelper.isEmpty(value)) {
        value = isMultiple ? [] : "";
    }

    $: valueAsArray = CommonHelper.toArray(value);

    $: maxReached =
        (valueAsArray.length || uploadedFiles.length) &&
        field.maxSelect <= valueAsArray.length + uploadedFiles.length - deletedFileNames.length;

    $: if (uploadedFiles !== -1 || deletedFileNames !== -1) {
        triggerListChange();
    }

    function restoreExistingFile(name) {
        CommonHelper.removeByValue(deletedFileNames, name);
        deletedFileNames = deletedFileNames;
    }

    function removeExistingFile(name) {
        CommonHelper.pushUnique(deletedFileNames, name);
        deletedFileNames = deletedFileNames;
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
                detail: { value, uploadedFiles, deletedFileNames },
                bubbles: true,
            }),
        );
    }

    function dropHandler(e) {
        e.preventDefault();

        isDragOver = false;

        const files = e.dataTransfer?.files || [];

        if (maxReached || !files.length) {
            return;
        }

        for (const file of files) {
            const currentTotal = valueAsArray.length + uploadedFiles.length - deletedFileNames.length;

            if (field.maxSelect <= currentTotal) {
                break;
            }

            uploadedFiles.push(file);
        }

        uploadedFiles = uploadedFiles;
    }

    async function openInNewTab(filename) {
        try {
            let token = await ApiClient.getSuperuserFileToken(record.collectionId);
            let url = ApiClient.files.getURL(record, filename, { token });
            window.open(url, "_blank", "noreferrer, noopener");
        } catch (err) {
            console.warn("openInNewTab file token failure:", err);
        }
    }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    class="block"
    on:dragover|preventDefault={() => {
        isDragOver = true;
    }}
    on:dragleave={() => {
        isDragOver = false;
    }}
    on:drop={dropHandler}
>
    <Field
        class="
            form-field form-field-list form-field-file
            {field.required ? 'required' : ''}
            {isDragOver ? 'dragover' : ''}
        "
        name={field.name}
        let:uniqueId
    >
        <FieldLabel {uniqueId} {field} />

        <div bind:this={filesListElem} class="list">
            {#each valueAsArray as filename, i (filename + record.id)}
                {@const isDeleted = deletedFileNames.includes(filename)}
                <Draggable
                    bind:list={value}
                    group={field.name + "_uploaded"}
                    index={i}
                    disabled={!isMultiple}
                    let:dragging
                    let:dragover
                >
                    <div class="list-item" class:dragging class:dragover>
                        <div class:fade={isDeleted}>
                            <RecordFileThumb {record} {filename} />
                        </div>

                        <div class="content">
                            <button
                                type="button"
                                draggable={false}
                                class="txt-ellipsis {isDeleted
                                    ? 'txt-strikethrough link-hint'
                                    : 'link-primary'}"
                                title="Download"
                                on:auxclick={() => openInNewTab(filename)}
                                on:click={() => openInNewTab(filename)}
                            >
                                {filename}
                            </button>
                        </div>

                        <div class="actions">
                            {#if isDeleted}
                                <button
                                    type="button"
                                    class="btn btn-sm btn-danger btn-transparent"
                                    on:click={() => restoreExistingFile(filename)}
                                >
                                    <span class="txt">Restore</span>
                                </button>
                            {:else}
                                <button
                                    type="button"
                                    class="btn btn-transparent btn-hint btn-sm btn-circle btn-remove"
                                    use:tooltip={"Remove file"}
                                    on:click={() => removeExistingFile(filename)}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            {/if}
                        </div>
                    </div>
                </Draggable>
            {/each}

            {#each uploadedFiles as file, i (file.name + i)}
                <Draggable
                    bind:list={uploadedFiles}
                    group={field.name + "_new"}
                    index={i}
                    disabled={!isMultiple}
                    let:dragging
                    let:dragover
                >
                    <div class="list-item" class:dragging class:dragover>
                        <figure class="thumb">
                            <UploadedFilePreview {file} />
                        </figure>
                        <div class="filename m-r-auto" title={file.name}>
                            <small class="label label-success m-r-5">New</small>
                            <span class="txt">{file.name}</span>
                        </div>
                        <button
                            type="button"
                            class="btn btn-transparent btn-hint btn-sm btn-circle btn-remove"
                            use:tooltip={"Remove file"}
                            on:click={() => removeNewFile(i)}
                        >
                            <i class="ri-close-line" />
                        </button>
                    </div>
                </Draggable>
            {/each}

            <div class="list-item list-item-btn">
                <input
                    bind:this={fileInput}
                    type="file"
                    class="hidden"
                    accept={field.mimeTypes?.join(",") || null}
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
</div>
