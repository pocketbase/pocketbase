window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Generic custom select element.
 * If max > 1 the select is multiple, otherwise - single (default).
 *
 * Note that if label or selected are custom DOM elements they need to be
 * wrapped in a function to allow recreation when toggling the select options.
 *
 * @example
 * ```js
 * app.components.select({
 *     options: [
 *         { value: "opt1", label: "Opt 1" },
 *         { value: "opt2", label: "Opt 2", selected: "Opt 2 selected label" },
 *         { value: "opt2", label: () => t.div(null, "Custom element") },
 *     ],
 *     value: () => data.selected,
 *     onchange: (opts) => data.selected = opts.map((opt) => opt.value),
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.select = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined, // used for associating with a label
        name: undefined, // used for error matching
        hidden: undefined,
        inert: undefined,
        className: "",
        value: undefined,
        options: [], // [{value, label?, selected?}, ...]
        before: null,
        after: null,
        max: 1,
        searchThreshold: 6,
        required: false,
        disabled: false,
        placeholder: "- Select -",
        noItemsFoundText: "No items found",
        onchange: function(selectedOpts) {},
        ondropdowntoggle: function(e) {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    if (props.max <= 0) {
        props.max = 1;
    }

    const internalData = store({
        selected: [],
        search: "",
        get hasSearch() {
            return internalData.search?.length > 0;
        },
        get allowRemove() {
            return !props.disabled && (!props.required || props.max > 1);
        },
    });

    function syncSelected() {
        if (typeof props.value === "undefined") {
            return; // nothing to sync
        }

        const optVals = app.utils.toArray(props.value, true);
        const cappedOptVals = optVals.slice(0, props.max || 1);

        if (optVals.length != cappedOptVals.length) {
            console.warn(
                `[select] the provided select values (${optVals.length}) are more than the allowed max selected options (${cappedOptVals.length}):`,
                optVals,
            );
            props.value = props.max > 1 ? cappedOptVals : cappedOptVals[0];
        }

        internalData.selected = optVals
            .map((value) => {
                return props.options.find((opt) => opt.value === value);
            })
            .filter(Boolean);
    }

    watchers.push(
        watch(
            () => props.value,
            () => syncSelected(),
        ),
    );

    async function toggle(opt) {
        const idx = internalData.selected.findIndex((o) => o.value === opt.value);
        if (idx >= 0) {
            if (!internalData.allowRemove) {
                dropdown?.hidePopover();
                return; // no change
            }

            internalData.selected.splice(idx, 1);
        } else {
            // clear last redundant elements (leaving place for the new selected)
            let toRemove = internalData.selected.length - props.max;
            while (toRemove >= 0) {
                internalData.selected.pop();
                toRemove--;
            }

            internalData.selected.push(opt);
        }

        if (props.max <= 1) {
            dropdown?.hidePopover();
        }

        if (props.onchange) {
            await props.onchange(internalData.selected);
            syncSelected(); // manually sync in case in the onchange handler the value didn't change
        }

        // trigger custom change event for clearing field errors
        if (selectedContainer?.isConnected) {
            selectedContainer.dispatchEvent(
                new CustomEvent("change", {
                    detail: internalData.selected,
                    bubbles: true,
                }),
            );
        }
    }

    function isSelected(opt) {
        return internalData.selected.findIndex((o) => o.value === opt.value) >= 0;
    }

    const searchInput = t.input({
        type: "text",
        placeholder: "Search...",
        value: () => internalData.search,
        oninput: (e) => (internalData.search = e.target.value),
    });

    function clearSearch(focus = false) {
        internalData.search = "";

        if (focus) {
            searchInput?.focus();
        }
    }

    const noItemsFoundElem = t.div({ className: "txt-hint txt-center m-0 p-5", hidden: true }, props.noItemsFoundText);

    async function toggleNoItemsFoundElem() {
        if (!dropdown) {
            return;
        }

        await new Promise((r) => setTimeout(r, 0));

        if (dropdown.querySelector(".select-option:not([hidden])")) {
            noItemsFoundElem.hidden = true;
        } else {
            noItemsFoundElem.hidden = false;
        }
    }

    const dropdown = t.div(
        {
            tabIndex: -1,
            popover: "auto",
            className: "dropdown",
            onbeforetoggle: (e) => {
                if (e.newState == "closed") {
                    clearSearch();
                }
                return props.ondropdowntoggle?.(e);
            },
        },
        t.div(
            {
                className: "fields dropdown-search",
                hidden: () => props.options.length < props.searchThreshold,
            },
            t.div({ className: "field" }, searchInput),
            t.div(
                {
                    className: "field addon p-r-5",
                    hidden: () => !internalData.hasSearch,
                },
                t.button(
                    {
                        type: "button",
                        title: "Clear",
                        className: "btn sm secondary transparent circle",
                        onclick: () => clearSearch(true),
                    },
                    t.i({ className: "ri-close-line", ariaHidden: true }),
                ),
            ),
        ),
        () => props.before?.__raw || props.before,
        () => {
            return props.options.map((opt) => {
                return t.button(
                    {
                        type: "button",
                        className: () => `dropdown-item select-option ${isSelected(opt) ? "active" : ""}`,
                        onclick: () => {
                            toggle(opt);
                            return false;
                        },
                    },
                    opt.label || opt.value,
                );
            });
        },
        noItemsFoundElem,
        () => props.after?.__raw || props.after,
    );

    const selectedContainer = t.button(
        {
            type: "button",
            id: () => props.id,
            name: () => props.name,
            disabled: () => props.disabled,
            className: () => `selected-container ${props.className}`,
            popoverTargetElement: dropdown,
            onclick: (e) => {
                e.stopPropagation();
            },
        },
        () => {
            if (!internalData.selected.length) {
                return t.span({ rid: "selected-placeholder", className: "placeholder" }, () => props.placeholder);
            }

            return internalData.selected.map((opt) => {
                return t.div({ className: "selected-item" }, opt.selected || opt.label || opt.value, () => {
                    if (!internalData.allowRemove) {
                        return;
                    }

                    return t.i({
                        tabIndex: -1,
                        role: "button",
                        className: "ri-close-line link-hint btn-option-unset",
                        ariaLabel: app.attrs.tooltip("Unset"),
                        onclick: () => {
                            toggle(opt);
                            return false;
                        },
                    });
                });
            });
        },
    );

    watchers.push(
        watch(
            () => props.options,
            () => {
                toggleNoItemsFoundElem();
            },
        ),
    );

    // search watcher
    let searchDebounce;
    watchers.push(
        watch(
            () => internalData.search,
            () => {
                const normalizedSearch = internalData.search.toLowerCase().replaceAll(" ", "");

                clearTimeout(searchDebounce);
                searchDebounce = setTimeout(() => {
                    const options = dropdown.querySelectorAll(".select-option");

                    if (!normalizedSearch.length) {
                        options.forEach((opt) => (opt.hidden = false));
                    } else {
                        options.forEach((opt) => {
                            const txt = opt.textContent.toLowerCase().replaceAll(" ", "");
                            if (!txt.includes(normalizedSearch)) {
                                opt.hidden = true;
                            } else {
                                opt.hidden = false;
                            }
                        });
                    }

                    toggleNoItemsFoundElem();
                }, 100);
            },
        ),
    );

    return t.div(
        {
            rid: props.rid,
            hidden: () => props.hidden,
            inert: () => props.inert,
            onmount: (el) => {
                el.addEventListener("focusout", function(e) {
                    if (!e.relatedTarget || !el.contains(e.relatedTarget)) {
                        dropdown?.hidePopover();
                    }
                });
            },
            onunmount: () => {
                clearTimeout(searchDebounce);
                watchers.forEach((w) => w.unwatch());
            },
            className: () => {
                return [
                    "input",
                    "select",
                    props.max > 1 ? "multiple" : "single",
                    props.disabled ? "disabled" : "",
                    props.required ? "required" : "",
                ].join(" ");
            },
        },
        selectedContainer,
        dropdown,
    );
};
