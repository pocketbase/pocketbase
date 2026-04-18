const toasts = new Map();

const toastsContainer = t.div({ className: "toasts-container" });

/**
 * Removes a single notification by its reference (key or content).
 *
 * @param {string|Node} toastRef
 * @param {boolean}     [animate]
 */
function removeToast(toastRef, animate = true) {
    const toast = toasts.get(toastRef);
    if (!toast || !toast.isConnected) {
        return;
    }

    toasts.delete(toastRef);
    clearTimeout(toast._removeTimeout);

    if (animate) {
        toast.classList.add("removing");
        setTimeout(() => {
            toast.remove();
        }, 300);
    } else {
        toast.remove();
    }
}

/**
 * Removes all registered notifications.
 *
 * @param {boolean} [animate]
 */
function removeAllToasts(animate = true) {
    toasts.forEach((_, key) => {
        window.app.toasts.remove(key, animate);
    });
}

/**
 * Adds "info" notification.
 *
 * @see {@link addToast}
 */
function infoToast(textOrElem, options = {}) {
    options.type = "info";
    options.duration = options.duration || 3000;

    addToast(textOrElem, options);
}

/**
 * Adds "success" notification.
 *
 * @see {@link addToast}
 */
function successToast(textOrElem, options = {}) {
    options.type = "success";
    options.duration = options.duration || 3000;

    addToast(textOrElem, options);
}

/**
 * Adds "error" notification.
 *
 * @see {@link addToast}
 */
function errorToast(textOrElem, options = {}) {
    options.type = "error";
    options.duration = options.duration || 3500;

    addToast(textOrElem, options);
}

/**
 * Creates and registers a new toast notification.
 *
 * @param {string|Node} textOrElem  The content of the notification as plain text or DOM node.
 * @param {object} [options] Toast options.
 * @param {number} [options.duration] Duration time in ms the notifaction would be active.
 * @param {string} [options.key]      Optional identifier that could be used to manually remove the specific notification (default to `textOrElem`).
 * @param {string} [options.type]     The CSS class type of the notification.
 */
function addToast(textOrElem, options = {}) {
    options = Object.assign({ duration: 3000, key: undefined, type: "info" }, options);

    if (!toastsContainer.isConnected) {
        document.body.appendChild(toastsContainer);
    }

    const toastRef = options.key || textOrElem;

    if (toasts.has(toastRef)) {
        removeToast(toastRef, false);
    }

    function initRemoveTimer(el) {
        if (el?._removeTimeout) {
            clearTimeout(el?._removeTimeout);
        }

        el._removeTimeout = setTimeout(() => {
            removeToast(toastRef);
        }, options.duration);
    }

    let newToast = t.div(
        {
            className: `toast ${options.type || ""}`,
            onmount: (el) => {
                initRemoveTimer(el);
            },
            onunmount: (el) => {
                if (el?._removeTimeout) {
                    clearTimeout(el?._removeTimeout);
                    newToast = null;
                }
            },
            onmouseover: () => {
                clearTimeout(newToast?._removeTimeout);
            },
            onmouseout: () => {
                initRemoveTimer(newToast);
            },
        },
        t.div(
            { className: "toast-container" },
            t.div({ className: "toast-icon" }),
            t.div(
                { className: "toast-content" },
                textOrElem,
                t.button(
                    {
                        className: "m-l-auto btn circle sm transparent secondary toast-remove",
                        title: "Clear",
                        onclick: () => removeToast(toastRef),
                    },
                    t.i({ className: "ri-close-line", ariaHidden: true }),
                ),
            ),
        ),
    );

    toasts.set(toastRef, newToast);
    toastsContainer.prepend(newToast);
}

// -------------------------------------------------------------------

window.app = window.app || {};
window.app.toasts = window.app.toasts || {};
window.app.toasts.info = infoToast;
window.app.toasts.error = errorToast;
window.app.toasts.success = successToast;
window.app.toasts.remove = removeToast;
window.app.toasts.removeAll = removeAllToasts;
