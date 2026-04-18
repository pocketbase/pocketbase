window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Thumb preview element for an uploaded File.
 * For non-image file the thumb is an icon representing the file type.
 *
 * @example
 * ```js
 * app.components.uploadedFileThumb({
 *     file: new File(...),
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.uploadedFileThumb = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        file: undefined, // File
        imageWidth: 100, // image thumb width
        imageHeight: 100, // image thumb height
        extraClasses: "sm", // any .thumb related classes
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        thumbSrc: undefined,
    });

    watchers.push(
        watch(
            () => [props.file, props.imageWidth, props.imageHeight],
            () => {
                if (app.utils.hasImageExtension(props.file?.name)) {
                    app.utils
                        .generateThumb(props.file, props.imageWidth, props.imageHeight)
                        .then((url) => {
                            data.thumbSrc = url;
                        })
                        .catch((err) => {
                            console.warn("unable to generate thumb:", err);
                            data.thumbSrc = undefined;
                        });
                } else {
                    data.thumbSrc = undefined;
                }
            },
        ),
    );

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `thumb ${props.extraClasses}`,
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        () => {
            const fileType = app.utils.getFileType(props.file?.name);

            if (fileType == "image" && data.thumbSrc) {
                return t.img({
                    draggable: false,
                    loading: "lazy",
                    alt: () => "Thumb of " + props.file.name,
                    src: data.thumbSrc,
                });
            }

            return t.i({
                className: app.utils.fileTypeIcons[fileType] || "ri-file-line",
                ariaHidden: true,
            });
        },
    );
};
