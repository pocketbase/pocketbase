window.app = window.app || {};
window.app.modals = window.app.modals || {};

const modalAttr = "data-modal-state";
const modalManualClass = "manual";

let oldActiveElem;

/**
 * Initializes and opens `el` as modal.
 *
 * You can also make use of the following custom `el` function properties:
 *
 *  - onbeforeopen(el)  - triggered before the opening sequence (return `false` to stop)
 *  - onafteropen(el)   - triggered after the opening sequence
 *  - onbeforeclose(el, forceClosed) - triggered before the closing sequence (return `false` to stop);
 *                        (note that when force closing the result of this callback is ignored)
 *  - onafterclose(el, forceClosed)  - triggered after the closing sequence
 *
 * @example
 * ```js
 * modal = t.div(
 *     {
 *         className: "modal popup sm",
 *         onbeforeclose: (el) => {
 *             if (!someImportantCheck) {
 *                 return false
 *             }
 *         },
 *     },
 *     t.header({ className: "modal-header" },
 *         t.h5({ className: "m-auto" }, "Logs settings"),
 *     ),
 *     t.form(
 *         {
 *             id: "myForm",
 *             className: "modal-content",
 *             onsubmit: (e) => {
 *                 e.preventDefault();
 *                 // ...
 *             },
 *         },
 *         // ... content ...
 *     ),
 *     t.footer({ className: "modal-footer" },
 *         t.button(
 *             {
 *                 type: "button",
 *                 className: "btn transparent m-r-auto",
 *                 onclick: () => app.modals.close(),
 *             },
 *             t.span({ className: "txt" }, "Close"),
 *         ),
 *         t.button(
 *             {
 *                 "html-form": "myForm",
 *                 type: "submit",
 *                 className: "btn",
 *             },
 *             t.span({ className: "txt" }, "Save changes"),
 *         ),
 *     ),
 * );
 *
 * app.modals.open(modal)
 * ```
 *
 * @param {Element} el
 */
window.app.modals.open = async function(el) {
    if (!el?.isConnected) {
        console.error("modals.open requies an active DOM element", el);
        return;
    }

    let beforeopen;
    if (el.onbeforeopen) {
        beforeopen = await el.onbeforeopen(el);
    }
    if (beforeopen === false) {
        return;
    }

    // note: currently doesn't wait for the entrance animation but this may change in the future
    if (el.onafteropen) {
        let resizeObserver = new ResizeObserver((entries) => {
            for (const entry of entries) {
                if (entry.contentRect.height > 0 && el?.onafteropen) {
                    el.onafteropen(el);
                    resizeObserver.disconnect();
                    resizeObserver = null; // observers uses weak  but clear it explicitly nonetheless
                }
            }
        });
        resizeObserver.observe(el);
    }

    oldActiveElem = document.activeElement;

    // init (if not already)
    initModal(el);

    // force close on history change
    el._forceClose = () => app.modals.close(el, true);
    window.addEventListener("popstate", el._forceClose);

    const largestZIndex = Math.max(findLargestZIndexModal(el)?.style.zIndex << 0, 1000);
    el.style.zIndex = largestZIndex + 1;

    // note: use data attribute to avoid conflict with reactive className
    el.setAttribute(modalAttr, "open");
};

/**
 * Closes the specified modal (or the last/top open one if not explicitly set).
 *
 * Example:
 * ```js
 * app.modals.close()
 * app.modals.close(myModal)
 *
 * // force close
 * app.modals.close(null, true)
 * app.modals.close(myModal, true)
 * ```
 *
 * @param {Element} [el]
 * @param {boolean} [forceClose]
 */
window.app.modals.close = async function(el = null, forceClose = false) {
    el = el || findLargestZIndexModal();
    if (!el) {
        return;
    }

    window.removeEventListener("popstate", el._forceClose);

    if (forceClose) {
        el.onbeforeclose?.(el, true);
        el.setAttribute(modalAttr, "close");
        el.onafterclose?.(el, true);
    } else {
        if (
            el.onbeforeclose
            && (await el.onbeforeclose(el, false)) === false
        ) {
            return;
        }

        if (el.onafterclose) {
            let resizeObserver = new ResizeObserver((entries) => {
                for (const entry of entries) {
                    if (entry.contentRect.height <= 0 && el?.onafterclose) {
                        el.onafterclose(el, false);
                        resizeObserver.disconnect();
                        resizeObserver = null; // observers uses weak ref but clear it explicitly nonetheless
                    }
                }
            });
            resizeObserver.observe(el);
        }

        el.setAttribute(modalAttr, "close");
    }

    // restore original focus position without making the focus visible
    if (oldActiveElem) {
        oldActiveElem.focus?.();
        setTimeout(() => {
            oldActiveElem?.blur?.();
            oldActiveElem = null;
        }, 0);
    }
};

function initModal(el) {
    if (!el.getAttribute("tabindex")) {
        el.setAttribute("tabindex", "-1");
    }

    // focus to capture key events and to change the tab navigation to the modal
    // (execute after the element rendering task)
    setTimeout(() => {
        const autofocusEl = el?.querySelector("[autofocus]");
        if (autofocusEl) {
            autofocusEl.focus();
        } else {
            el?.focus();
        }
    }, 0);

    // already initialized
    if (el.getAttribute(modalAttr)) {
        return;
    }

    el.setAttribute(modalAttr, "");

    // dismiss handlers
    el.addEventListener("keydown", (e) => {
        if (
            e.key != "Escape"
            || el.classList.contains(modalManualClass)
            || (e.target !== el && el.contains(e.target))
        ) {
            return;
        }

        window.app.modals.close(el);
    });

    let startedInside = false;
    const startedInsideFunc = (e) => {
        startedInside = e.target !== el && el.contains(e.target);
    };
    el.addEventListener("mousedown", startedInsideFunc);
    el.addEventListener("touchstart", startedInsideFunc);

    let endedInside = false;
    const endedInsideFunc = (e) => {
        endedInside = e.target !== el && el.contains(e.target);
    };
    el.addEventListener("mouseup", endedInsideFunc);
    el.addEventListener("touchend", endedInsideFunc);

    el.addEventListener("click", (e) => {
        if (
            startedInside
            || endedInside
            || el.classList.contains(modalManualClass)
            // e.g. in case a btn is clicked with the keyboard (Enter/Space/etc.)
            || (e.target !== el && el.contains(e.target))
        ) {
            return;
        }

        window.app.modals.close(el);
    });
}

function findLargestZIndexModal(excludeEl) {
    const opened = document.querySelectorAll(`[${modalAttr}="open"]`);

    let z = 0;
    let max = 0;
    let largest;

    for (const m of opened) {
        if (excludeEl && m == excludeEl) {
            continue;
        }

        z = m.style.zIndex << 0;
        if (z > max) {
            max = z;
            largest = m;
        }
    }

    return largest;
}
