window.app = window.app || {};
window.app.modals = window.app.modals || {};

/**
 * Opens a confirmation dialog and executes `yesCallback` or `noCallback`
 * depending on the user's choice.
 *
 * The callbacks can return a `Promise`` that will be awaited before
 * closing the confirmation modal.
 *
 * @example
 * ```js
 * app.modals.confirm("Are you sure?", () => console.log("confirmed"))
 * ```
 *
 * @param {string|Element} textOrElem The confirmation message.
 * @param {function} [yesCallback]
 * @param {function} [noCallback]
 * @param {Object}   [settings]
 */
window.app.modals.confirm = function(textOrElem, yesCallback, noCallback, settings = {
    className: undefined,
    yesButton: "",
    noButton: "",
}) {
    data.textOrElem = textOrElem;
    data.yesCallback = yesCallback;
    data.yesCallbackWaiting = false;
    data.noCallback = noCallback;
    data.noCallbackWaiting = false;
    data.className = typeof settings.className == "string" ? settings.className : "sm";
    data.yesButton = settings.yesButton || "Yes";
    data.noButton = settings.noButton || "No";

    if (!confirmElem.isConnected) {
        document.body.appendChild(confirmElem);
    }

    window.app.modals.open(confirmElem);
};

const data = store({
    className: "",
    textOrElem: null,
    // ---
    yesButton: "",
    yesCallback: null,
    yesCallbackWaiting: false,
    // ---
    noButton: "",
    noCallback: null,
    noCallbackWaiting: false,
    // ---
    get isBusy() {
        return data.yesCallbackWaiting || data.noCallbackWaiting;
    },
});

const confirmElem = t.div(
    { className: () => `modal popup manual ${data.className || ""}` },
    t.div(
        { className: "modal-content" },
        (el) => {
            if (typeof data.textOrElem === "string") {
                return t.h6({ className: "block txt-center" }, data.textOrElem);
            }

            if (typeof data.textOrElem === "function") {
                return data.textOrElem(el);
            }

            return data.textOrElem;
        },
    ),
    t.footer(
        { className: "modal-footer p-sm" },
        t.div(
            { className: "grid sm" },
            t.div(
                { className: "col-sm-6" },
                t.button(
                    {
                        type: "button",
                        className: () => `btn lg block secondary ${data.noCallbackWaiting ? "loading" : ""}`,
                        disabled: () => data.isBusy,
                        onclick: async () => {
                            if (data.noCallback) {
                                data.noCallbackWaiting = true;
                                try {
                                    const result = await data.noCallback();
                                    if (result === false) {
                                        return;
                                    }
                                } catch (err) {
                                    console.log("confirm noCallback error:", err);
                                } finally {
                                    data.noCallbackWaiting = false;
                                }
                            }

                            window.app.modals.close(confirmElem);
                        },
                    },
                    () => {
                        if (typeof data.noButton === "string") {
                            return t.span({ className: "txt" }, data.noButton);
                        }

                        if (typeof data.noButton === "function") {
                            return data.noButton(el);
                        }

                        return data.noButton;
                    },
                ),
            ),
            t.div(
                { className: "col-sm-6" },
                t.button(
                    {
                        type: "button",
                        className: () => `btn lg block warning ${data.yesCallbackWaiting ? "loading" : ""}`,
                        disabled: () => data.isBusy,
                        onclick: async () => {
                            if (data.yesCallback) {
                                data.yesCallbackWaiting = true;
                                try {
                                    const result = await data.yesCallback();
                                    if (result === false) {
                                        return;
                                    }
                                } catch (err) {
                                    console.log("confirm yesCallback error:", err);
                                } finally {
                                    data.yesCallbackWaiting = false;
                                }
                            }

                            window.app.modals.close(confirmElem);
                        },
                    },
                    () => {
                        if (typeof data.yesButton === "string") {
                            return t.span({ className: "txt" }, data.yesButton);
                        }

                        if (typeof data.yesButton === "function") {
                            return data.yesButton(el);
                        }

                        return data.yesButton;
                    },
                ),
            ),
        ),
    ),
);
