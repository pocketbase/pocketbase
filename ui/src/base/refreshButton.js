window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Animated refresh button element.
 *
 * @example
 * ```js
 * app.components.refreshButton({
 *     onclick: () => { console.log("clicked...") },
 * })
 * ```
 *
 * @param  {Object} propsArg
 * @return {Element}
 */
window.app.components.refreshButton = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        tooltip: "Refresh",
        className: "btn transparent secondary circle rotate-btn",
        disabled: false,
        onclick: function(e) {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    let refreshTimeoutId;

    const btn = t.button(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            type: "button",
            ariaLabel: app.attrs.tooltip(() => props.tooltip),
            disabled: () => props.disabled,
            className: () => props.className,
            onunmount: () => {
                clearTimeout(refreshTimeoutId);
                watchers.forEach((w) => w?.unwatch());
            },
            onclick: (e) => {
                e.preventDefault();

                if (props.onclick) {
                    props.onclick(e);
                }

                btn.classList.add("rotate");
                btn.addEventListener("animationend", () => {
                    btn.classList.remove("rotate");
                });

                // fallback
                clearTimeout(refreshTimeoutId);
                refreshTimeoutId = setTimeout(() => {
                    clearTimeout(refreshTimeoutId);
                    btn.classList.remove("rotate");
                }, 500);
            },
        },
        t.i({ className: "ri-refresh-line", ariaHidden: true }),
    );

    return btn;
};
