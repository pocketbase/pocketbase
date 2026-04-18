window.app = window.app || {};
window.app.modals = window.app.modals || {};

/**
 * Opens a new file preview popup.
 *
 * @example
 * ```js
 * app.modals.openFilePreview(url)
 * ```
 *
 * @param {string|Promise|Function} urlOrFactory
 */
window.app.modals.openFilePreview = function(urlOrFactory) {
    const modal = filePreviewModal(urlOrFactory);

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function filePreviewModal(urlOrFactory) {
    const data = store({
        url: "",
        get filename() {
            const url = data.url;
            const queryParamsIdx = url.indexOf("?");

            return url.substring(url.lastIndexOf("/") + 1, queryParamsIdx > 0 ? queryParamsIdx : undefined);
        },
        get fileType() {
            return app.utils.getFileType(data.filename);
        },
    });

    async function resolveUrlOrFactory() {
        let url = "";

        try {
            if (typeof urlOrFactory == "function") {
                url = await urlOrFactory();
            } else {
                // string or Promise
                url = await urlOrFactory;
            }
        } catch (err) {
            if (!err.isAbort) {
                console.warn("resolveUrlOrFactory file preview failure:", err);
            }
        }

        data.url = url;

        return url;
    }

    async function openInNewTab() {
        // resolve again because it may have expired
        let url = await resolveUrlOrFactory();
        if (!url) {
            return;
        }

        window.open(url, "_blank", "noreferrer,noopener");
    }

    return t.div(
        {
            pbEvent: "filePreviewModal",
            className: () => `modal preview preview-${data.fileType}`,
            onbeforeopen: () => {
                resolveUrlOrFactory();
            },
            onafterclose: (el) => {
                el.remove();
            },
        },
        t.div({ className: "modal-content" }, () => {
            if (!data.url) {
                return t.span({ className: "loader" });
            }

            if (data.fileType == "image") {
                return t.img({
                    src: () => data.url,
                    alt: () => `Preview ${data.filename}`,
                });
            }

            return t.object(
                {
                    data: data.url, // note: the reactive value doesn't trigger reload of the object
                    title: () => data.filename,
                },
                "Cannot preview the file.",
            );
        }),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "link-hint filename-link",
                    ariaDescription: app.attrs.tooltip("Open in new tab"),
                    onclick: () => openInNewTab(),
                },
                t.span({ className: "txt" }, () => data.filename),
            ),
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-l-auto",
                    onclick: () => app.modals.close(),
                },
                t.span({ className: "txt" }, "Close"),
            ),
        ),
    );
}
