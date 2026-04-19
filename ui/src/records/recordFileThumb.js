window.app = window.app || {};
window.app.components = window.app.components || {};

// rudimentary semaphore to resolve the images url on batches to prevent
// overhelming and freezing the browser tab with rendering too many images at once
const semaphore = {
    max: 10,
    pending: new Set(),
    processing: new Set(),
};

function semaphoreAdd(fn) {
    semaphore.pending.add(fn);

    if (semaphore.processing.size <= semaphore.max) {
        semaphoreProcess();
    }

    // release func that must be called manually after done with the loading
    return () => {
        semaphore.pending.delete(fn);
        semaphore.processing.delete(fn);

        if (semaphore.processing.size < semaphore.max) {
            semaphoreProcess();
        }
    };
}

function semaphoreProcess() {
    for (const fn of semaphore.pending) {
        semaphore.pending.delete(fn);
        semaphore.processing.add(fn);
        fn();
        return;
    }
}

/**
 * Creates a record file thumb element.
 *
 * @example
 * ```js
 * app.components.recordFileThumb({
 *     record: () => data.record,
 *     filename: () => data.record.myFile,
 * })
 * ```
 *
 * @param  {Object} propsArg
 * @return {Element}
 */
window.app.components.recordFileThumb = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        record: {},
        filename: "",
        extraClasses: "sm", // any .thumb related classes
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        isPreviewLoading: false,
        previewToken: "",
        get fileType() {
            return app.utils.getFileType(props.filename);
        },
        get hasPreview() {
            return ["image", "audio", "video"].includes(data.fileType) || props.filename.endsWith(".pdf");
        },
        previewURL: undefined,
    });

    return t.button(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            type: "button",
            draggable: false,
            className: () => `thumb ${props.extraClasses} ${data.isPreviewLoading ? "loading" : ""}`,
            title: () => (data.hasPreview ? "Preview" : "Download") + " " + props.filename,
            onclick: async (e) => {
                e.stopPropagation();

                async function resolveURL() {
                    const token = await app.getFileToken(props.record.collectionId);
                    return app.pb.files.getURL(props.record, props.filename, { token });
                }

                if (data.hasPreview) {
                    app.modals.openFilePreview(resolveURL);
                } else {
                    const url = await resolveURL();
                    app.utils.download(url, props.filename);
                }
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        () => {
            if (data.fileType == "image") {
                const img = t.img({
                    draggable: false,
                    alt: () => "Thumb of " + props.filename,
                    src: () => data.previewURL,
                    onerror: (err) => {
                        console.warn("[recordFileThumb] load err:", err);
                        data.isPreviewLoading = false;
                        img?._semaphoreRelease?.();
                    },
                    onload: () => {
                        data.isPreviewLoading = false;
                        img?._semaphoreRelease?.();
                    },
                    onmount: (el) => {
                        data.isPreviewLoading = true;

                        el._semaphoreRelease = semaphoreAdd(async () => {
                            try {
                                data.previewToken = await app.getFileToken(props.record.collectionId);

                                data.previewURL = app.pb.files.getURL(props.record, props.filename, {
                                    thumb: "100x100",
                                    token: data.previewToken,
                                });
                            } catch (err) {
                                console.warn(err);
                            }
                        });
                    },
                    onunmount: (el) => {
                        data.isPreviewLoading = false;
                        el._semaphoreRelease?.();
                    },
                });

                return img;
            }

            return t.i({ className: app.utils.fileTypeIcons[data.fileType] || "ri-file-line", ariaHidden: true });
        },
    );
};
