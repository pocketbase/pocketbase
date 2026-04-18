window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Color picker component (with predefined colors support).
 *
 * @example
 * ```js
 * app.components.colorPicker({
 *     value: () => data.color,
 *     predefinedColors: ["#ff0000", "#123456"],
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.colorPicker = function(propsArg = {}) {
    const uniqueId = "picker_" + app.utils.randomString();

    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        name: "",
        required: false,
        disabled: false,
        value: "",
        predefinedColors: [],
        // ---
        onchange: (newColor) => {},
        onmount: (el) => {},
        onunmount: (el) => {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const local = store({
        inputValue: "#ffffff",
    });

    function handleOnchange(val) {
        val = val?.toLowerCase() || "";
        props.onchange?.(val);
    }

    let inputTimeoutId;

    let input = t.input({
        type: "color",
        className: "color-picker-input",
        id: () => props.id,
        name: () => props.name,
        required: () => props.required,
        disabled: () => props.disabled,
        value: () => {
            local.inputValue = props.value || "#ffffff";
            return props.value || undefined;
        },
        oninput: (e) => {
            local.inputValue = e.target.value;

            clearTimeout(inputTimeoutId);
            inputTimeoutId = setTimeout(() => {
                handleOnchange(e.target.value);
            }, 50);
        },
    });

    return t.div(
        {
            rid: props.rid,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `color-picker ${props.className}`,
            onmount: (el) => {
                props.onmount?.(el);
            },
            onunmount: (el) => {
                clearTimeout(inputTimeoutId);

                props.onunmount?.(el);

                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: "color-picker-input-wrapper" },
            input,
            t.output({
                className: "result",
                // black or white text (https://developer.chrome.com/blog/css-relative-color-syntax)
                // @todo replace with contrast-color once there is better support?
                style: () => `color: lch(from ${local.inputValue || "#ffffff"} calc((49 - l) * infinity) 0 0);`,
                textContent: () => local.inputValue,
            }),
        ),
        t.button(
            {
                hidden: () => !props.predefinedColors.length,
                type: "button",
                title: "Predefined colors",
                className: "link-hint predefined-colors-btn",
                "html-popovertarget": uniqueId + "predefined-colors-dropdown",
            },
            t.i({ className: "ri-arrow-down-s-line", roleHidden: true }),
        ),
        t.div(
            {
                pbEvent: "predefinedColorsDropdown",
                id: uniqueId + "predefined-colors-dropdown",
                className: "dropdown predefined-colors-dropdown",
                popover: "auto",
            },
            t.div(
                {
                    className: "predefined-colors-list",
                },
                () => {
                    return props.predefinedColors?.map((color) => {
                        return t.button({
                            type: "button",
                            className: () => `color ${props.value == color ? "active" : ""}`,
                            style: `background:${color}`,
                            onclick: (e) => {
                                if (!input) {
                                    return;
                                }

                                e.target.closest(".dropdown")?.hidePopover();
                                input.value = color || undefined;
                                input.dispatchEvent(new Event("input", { bubbles: true }));
                            },
                        });
                    });
                },
            ),
        ),
    );
};
