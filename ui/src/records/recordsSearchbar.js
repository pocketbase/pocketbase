window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Creates a new records searchbar element with builtin autocomplete.
 * The searchbar is based on `app.components.search`.
 *
 * Note that the created element doesn't do any search. It is responsible
 * only for binding a reactive search input value.
 *
 * @example
 * ```js
 * app.components.recordsSearchbar({
 *     value: () => data.searchTerm,
 *     onsubmit: (newVal) => data.searchTerm = newVal,
 * })
 * ```
 *
 * @param  {Object} propsArg
 * @return {Element}
 */
window.app.components.recordsSearchbar = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        disabled: undefined,
        value: "",
        className: "",
        collection: undefined,
        onsubmit: (newValue) => {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `full-width records-searchbar-wrapper ${props.className}`,
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        app.components.searchbar({
            placeholder: () => (!props.disabled && !props.collection?.id ? "Loading..." : "Search term or filter..."),
            historyKey: () => "pbRecordsSearchHistory_" + props.collection?.id,
            disabled: () => props.disabled || !props.collection,
            value: () => props.value,
            autocomplete: (word) => {
                return app.utils.collectionAutocompleteKeys(props.collection, word, {
                    requestKeys: false,
                    collectionJoinKeys: false,
                });
            },
            onsubmit: props.onsubmit,
        }),
    );
};
