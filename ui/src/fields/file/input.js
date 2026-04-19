import { filesToDeleteProp } from "./onrecordsave.js";

// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "file_" + app.utils.randomString();

    function isDeleted(nameOrFile) {
        if (typeof nameOrFile != "string") {
            return false;
        }

        return !!props.record[filesToDeleteProp]?.[props.field.name]?.includes(nameOrFile);
    }

    function toDelete(nameOrFile) {
        // existing files are just marked for delete to allow restore
        if (typeof nameOrFile == "string") {
            props.record[filesToDeleteProp] = props.record[filesToDeleteProp] || {};
            props.record[filesToDeleteProp][props.field.name] = props.record[filesToDeleteProp][props.field.name] || [];
            app.utils.pushUnique(props.record[filesToDeleteProp][props.field.name], nameOrFile);
            triggerChangeEvent();
            return;
        }

        // new files are directly removed
        const normalized = app.utils.toArray(props.record[props.field.name]);
        const index = normalized.indexOf(nameOrFile);
        if (index >= 0) {
            normalized.splice(index, 1);
            props.record[props.field.name] = normalized;
            triggerChangeEvent();
        }
    }

    function restoreDeleted(nameOrFile) {
        if (typeof nameOrFile != "string") {
            return;
        }

        app.utils.removeByValue(props.record[filesToDeleteProp]?.[props.field.name], nameOrFile);
        triggerChangeEvent();
    }

    function totalDeletedFiles() {
        return props.record[filesToDeleteProp]?.[props.field.name]?.length || 0;
    }

    function totalFiles() {
        const totalNormalized = app.utils.toArray(props.record[props.field.name]).length;
        return totalNormalized - totalDeletedFiles();
    }

    // trigger custom change event for clearing field errors
    function triggerChangeEvent() {
        fieldContentEl?.dispatchEvent(
            new CustomEvent("change", {
                detail: { data: props },
                bubbles: true,
            }),
        );
    }

    const local = store({
        get maxReached() {
            const maxSelect = props.field.maxSelect || 1;
            return totalFiles() >= maxSelect;
        },
    });

    function addFiles(files) {
        const normalized = app.utils.toArray(props.record[props.field.name]);

        for (let file of files) {
            if (local.maxReached) {
                console.warn("can't add more files - max allowed files reached");
                break;
            }

            normalized.push(file);
        }

        props.record[props.field.name] = normalized;

        triggerChangeEvent();
    }

    const fileInput = t.input({
        type: "file",
        hidden: true,
        multiple: () => props.field.maxSelect > 1,
        accept: () => props.field.mimeTypes?.join(",") || undefined,
        onchange: (e) => {
            addFiles(e.target.files);
            e.target.value = null; // reset
        },
    });

    const fieldContentEl = t.output(
        {
            className: "field-content",
            name: () => props.field.name,
        },
        // @todo enable ordering new files before/inbetween existing
        app.components.sortable({
            className: "list",
            data: () => {
                const vals = app.utils.toArray(props.record[props.field.name]);

                let hadInvalid = false;
                // filter empty or invalid values (e.g. from old serialized draft)
                for (let i = vals.length - 1; i >= 0; i--) {
                    if (typeof vals[i] == "string" || vals[i] instanceof Blob) {
                        continue; // valid
                    }

                    hadInvalid = true;
                    vals.splice(i, 1);
                }

                // update record model to prevent conflict with required and other validators
                if (hadInvalid) {
                    props.record[props.field.name] = vals;
                }

                return vals;
            },
            onchange: (sortedList) => {
                props.record[props.field.name] = sortedList;
                triggerChangeEvent();
            },
            dataItem: (nameOrFile, i) => {
                return t.div(
                    {
                        rid: nameOrFile,
                        className: () => `list-item highlight ${isDeleted(nameOrFile) ? "deleted" : ""}`,
                    },
                    t.div({ className: "content gap-10" }, () => {
                        if (typeof nameOrFile == "string") {
                            return [
                                app.components.recordFileThumb({
                                    record: props.record,
                                    filename: nameOrFile,
                                }),
                                t.button(
                                    {
                                        type: "button",
                                        ariaDescription: app.attrs.tooltip("Open in new tab"),
                                        onclick: async () => {
                                            const token = await app.getFileToken(props.record.collectionId);
                                            const url = app.pb.files.getURL(props.record, nameOrFile, {
                                                token,
                                            });
                                            window.open(url, "_blank", "noreferrer,noopener");
                                        },
                                    },
                                    t.span({ className: "txt link-primary" }, nameOrFile),
                                ),
                            ];
                        }

                        return [
                            app.components.uploadedFileThumb({
                                file: nameOrFile,
                            }),
                            t.span({ className: "label success" }, "New"),
                            t.span({ className: "txt" }, nameOrFile.name),
                        ];
                    }),
                    t.div(
                        { className: "actions" },
                        t.button(
                            {
                                type: "button",
                                className: "btn sm secondary transparent circle",
                                ariaLabel: app.attrs.tooltip("Remove file"),
                                hidden: () => isDeleted(nameOrFile),
                                onclick: () => toDelete(nameOrFile),
                            },
                            t.i({ className: "ri-close-line", ariaHidden: true }),
                        ),
                        t.button(
                            {
                                type: "button",
                                className: "btn sm warning transparent",
                                hidden: () => !isDeleted(nameOrFile),
                                onclick: () => restoreDeleted(nameOrFile),
                            },
                            t.span({ className: "txt" }, "Restore"),
                        ),
                    ),
                );
            },
        }),
        t.hr({
            className: "m-t-5 m-b-0",
            hidden: () => app.utils.toArray(props.record[props.field.name]).length > 0,
        }),
        t.button(
            {
                type: "button",
                className: "btn sm secondary block",
                title: () => local.maxReached ? "Max allowed files reached" : undefined,
                disabled: () => local.maxReached,
                onclick: (e) => {
                    if (!local.maxReached) {
                        fileInput?.click();
                    }
                    document.activeElement?.blur();
                },
            },
            t.i({ className: "ri-upload-cloud-line", ariaHidden: true }),
            t.span({ className: "txt" }, "Upload or drop new file"),
        ),
    );

    return t.div(
        {
            className: "record-field-input field-type-file",
            ondragover: (e) => {
                e.preventDefault(); // prevent default to allow drop
            },
            ondrop: (e) => {
                const files = e.dataTransfer?.files || [];
                if (!files.length) {
                    return; // not a file drop
                }

                e.preventDefault();

                if (local.maxReached) {
                    return;
                }

                addFiles(files);
            },
        },
        t.div(
            { className: () => `field ${props.field.required ? "required" : ""}` },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.file.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            fileInput,
            fieldContentEl,
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
