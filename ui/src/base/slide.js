window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Wraps the children elements in a collapsible container.
 *
 * @example
 * ```js
 * app.components.slide(
 *     () => data.showToggle,
 *     t.div(null, "child1..."),
 *     t.div(null, "child2..."),
 * )
 * ```
 *
 * @param {function} boolFunc Boolean function that indicates whether the container is visible or not.
 * @param  {Array<Element>} [children]
 * @return {Element}
 */
window.app.components.slide = function(boolFunc, ...children) {
    let initTimeoutId;

    return t.div(
        {
            className: (el) => `block slide-block ${boolFunc?.(el) ? "" : "hidden"}`,
            onmount: (el) => {
                // add a ready attribute with slight delay to avoid @starting-style flickering
                initTimeoutId = setTimeout(() => {
                    el?.setAttribute("data-slide", "1");
                }, 200);
            },
            onunmount: () => {
                clearTimeout(initTimeoutId);
            },
        },
        ...children,
    );
};
