window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Generic page searchbar element.
 *
 * @example
 * ```js
 * app.components.searchbar({
 *     value: () => data.search,
 *     onsubmit: (newValue) => data.search = newValue,
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.searchbar = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        value: "",
        className: "",
        placeholder: "Search...",
        disabled: false,
        historyKey: "",
        autocomplete: undefined, // Array<string|Object> | function(word): Array<string|Object>,
        onsubmit: (newValue) => {},
    });

    const watchers = app.utils.extendStore(props, propsArg, "autocomplete");

    const local = store({
        value: "",
    });

    function submit() {
        props.value = local.value;
        props.onsubmit?.(local.value);
    }

    function clear() {
        local.value = "";
        submit();
    }

    watchers.push(
        // init and local sync changes
        watch(
            () => props.value,
            (searchTerm) => {
                local.value = searchTerm;
            },
        ),
    );

    return t.form(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `fields searchbar ${props.className}`,
            onsubmit: (e) => {
                e.preventDefault();
                submit();
            },
            onunmount: (el) => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        () => {
            if (!props.historyKey) {
                return;
            }

            return t.div(
                { className: "field addon p-l-5" },
                app.components.searchHistoryButton({
                    historyKey: () => props.historyKey,
                    value: () => props.value,
                    onselect: (val) => {
                        local.value = val;
                        submit();
                    },
                }),
            );
        },
        t.div(
            { className: "field" },
            app.components.codeEditor({
                singleLine: true,
                language: "pbrule",
                className: () => props.historyKey ? "p-l-5" : "p-l-20",
                placeholder: () => props.placeholder,
                disabled: () => props.disabled,
                value: () => local.value,
                oninput: (val) => (local.value = val),
                autocomplete: props.autocomplete,
            }),
        ),
        () => {
            if (props.value.length > 0 || local.value.length > 0) {
                return t.div(
                    { rid: "search-ctrls", className: "field addon p-r-5" },
                    t.button(
                        {
                            type: "submit",
                            className: "btn sm pill warning",
                            hidden: () => props.value == local.value,
                        },
                        "Search",
                    ),
                    t.button(
                        {
                            type: "button",
                            className: "btn sm pill secondary transparent",
                            onclick: () => clear(),
                        },
                        "Clear",
                    ),
                );
            }
        },
    );
};
