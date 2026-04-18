window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Helper button that shows a dropdown with previous search attempts.
 *
 * @example
 * ```js
 * app.components.searchHistoryButton({
 *     historyKey: "anything", // localStorage history key
 *     value: () => data.search,
 *     onselect: (historyVal) => data.search = historyVal,
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.searchHistoryButton = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        value: undefined,
        historyKey: "default",
        max: 15,
        openInNewTabParam: "filter",
        btnClassName: "btn sm pill secondary transparent p-r-5",
        onselect: function(val) {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const history = store({
        items: app.utils.getLocalHistory(props.historyKey, []),
    });

    function addToHistory(val) {
        removeFromHistory(val);
        history.items.unshift(val);
    }

    function removeFromHistory(val) {
        app.utils.removeByValue(history.items, val);
    }

    const uniqueId = "history_dropdown_" + app.utils.randomString();

    watchers.push(
        watch(
            () => props.value,
            (val) => {
                if (val) {
                    addToHistory(val);
                }
            },
        ),
    );

    watchers.push(
        watch(() => {
            if (history.items.length > props.max) {
                history.items = history.items.slice(0, props.max);
            }
            app.utils.saveLocalHistory(props.historyKey, history.items);
        }),
    );

    const dropdown = t.div(
        {
            id: uniqueId,
            className: "dropdown sm left nowrap history-searchbar-dropdown",
            popover: "hint",
            onclick: (e) => {
                e.stopPropagation();
                return false;
            },
        },
        t.div({ className: "block p-5" }, t.small({ className: "txt-hint" }, "Search history")),
        () => {
            if (!history.items?.length) {
                return t.div(
                    { rid: "no-history", className: "block p-5" },
                    t.span(null, "Your recent searches will show up here."),
                );
            }

            return history.items.slice(0, props.max).map((h) => {
                return t.button(
                    {
                        type: "button",
                        className: "dropdown-item txt-code",
                        onclick: () => {
                            dropdown.hidePopover();
                            props.onselect?.(h);
                            addToHistory(h);
                        },
                        onauxclick: () => {
                            if (props.openInNewTabParam) {
                                addToHistory(h);
                                dropdown.hidePopover();

                                const url = app.utils.replaceHashQueryParams(
                                    {
                                        [props.openInNewTabParam]: h,
                                    },
                                    false,
                                );
                                window.open(url, "_blank");
                            }
                        },
                    },
                    t.span({ className: "txt-ellipsis", title: h, textContent: h }),
                    t.small(
                        {
                            role: "button",
                            className: "remove-btn link-hint m-l-auto p-l-5 p-r-5",
                            title: "Clear",
                            onauxclick: (e) => {
                                e.stopPropagation();
                                return false;
                            },
                            onclick: (e) => {
                                e.stopPropagation();
                                removeFromHistory(h);
                                return false;
                            },
                        },
                        t.i({ className: "ri-close-line", ariaHidden: true }),
                    ),
                );
            });
        },
    );

    return t.button(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            type: "button",
            title: "Search history",
            className: () => props.btnClassName,
            "html-popovertarget": uniqueId,
            onunmount: () => {
                watchers?.forEach((w) => w?.unwatch());
            },
        },
        t.i({ className: "ri-search-line", ariaHidden: true }),
        t.i({ className: "ri-arrow-drop-down-line", ariaHidden: true }),
        dropdown,
    );
};
