window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * API rule input element.
 *
 * @example
 * ```js
 * app.components.ruleField({
 *     name: "listRule",
 *     autocomplete: (word) => {
 *         return app.utils.collectionAutocompleteKeys(someCollection, word);
 *     },
 *     value: () => someCollection.listRule,
 *     oninput: (newVal) => someCollection.listRule = newVal,
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.ruleField = function(propsArg = {}) {
    const uniqueId = "rule_" + app.utils.randomString();

    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        required: false,
        disabled: false,
        name: undefined,
        label: undefined,
        help: undefined,
        value: null,
        nullable: true,
        placeholder: "Leave empty to grant everyone access...",
        autocomplete: (word) => [],
        oninput: (newVal) => {},
        onmount: (el) => {},
        onunmount: (el) => {},
        // ---
        get isLocked() {
            return props.value == null;
        },
    });

    const watchers = app.utils.extendStore(props, propsArg, "isLocked");

    let ruleField;
    let _prevValue = "";

    function updateValue(newValue) {
        props.value = newValue;
        props.oninput?.(newValue);
        ruleField?.dispatchEvent(new CustomEvent("change", { detail: newValue }));
    }

    function lock() {
        if (props.value === null) {
            return;
        }

        _prevValue = props.value;
        updateValue(null);
    }

    function unlock() {
        if (_prevValue != null) {
            updateValue(_prevValue);
        } else {
            updateValue("");
        }

        setTimeout(() => {
            document.getElementById(uniqueId)?.focus();
        }, 0);
    }

    ruleField = t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            "html-name": () => props.name, // used for the error reset
            className: () =>
                [
                    "field",
                    "rule-field",
                    props.required ? "required" : null,
                    props.value === null ? "locked" : null,
                    props.disabled ? "disabled" : null,
                ].filter(Boolean).join(" "),
            onmount: (el) => {
                props.onmount?.(el);
            },
            onunmount: (el) => {
                props.onunmount?.(el);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.label(
            { htmlFor: uniqueId },
            (el) => {
                if (!props.label) {
                    return t.span({ className: "txt" }, "Rule");
                }

                if (typeof props.label == "function") {
                    return props.label(el);
                }

                if (typeof props.label == "string") {
                    return t.span({ className: "txt" }, props.label);
                }

                return props.label;
            },
            t.span({ hidden: () => !props.isLocked, className: "txt superusers-label" }, "(Superusers only)"),
        ),
        (el) => {
            if (props.isLocked) {
                return t.button(
                    {
                        type: "button",
                        className: "unlock-overlay",
                        disabled: () => props.disabled,
                        onclick: unlock,
                    },
                    t.span({ className: "txt" }, "Unlock and set custom rule"),
                    t.i({ className: "ri-lock-unlock-line", ariaHidden: true }),
                );
            }

            return [
                app.components.codeEditor({
                    id: uniqueId,
                    language: "pbrule",
                    required: () => props.required,
                    disabled: () => props.disabled,
                    value: () => props.value,
                    oninput: updateValue,
                    placeholder: () => props.placeholder,
                    autocomplete: props.autocomplete,
                    autocompleteContainer: el,
                }),
                t.button(
                    {
                        hidden: () => !props.nullable,
                        type: "button",
                        className: "superuser-toggle",
                        disabled: () => props.disabled,
                        onclick: lock,
                    },
                    t.i({ className: "ri-lock-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Set superusers only"),
                ),
            ];
        },
    );

    return ruleField;
};
