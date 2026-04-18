window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Readonly code highlight component.
 *
 * @example
 * ```js
 * app.components.codeBlock({
 *     value: () => data.myCode,
 *     language: "html",
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.codeBlock = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        language: "js", // see Prism.languages
        value: undefined,
        footnote: undefined,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `code-wrapper ${props.className}`,
            tabIndex: -1,
            onmount: (el) => {
                el.addEventListener("keydown", (e) => {
                    if ((e.ctrlKey || e.metaKey) && (e.key == "a" || e.key == "A")) {
                        e.preventDefault();
                        window.getSelection().selectAllChildren(el);
                    }
                });
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.code({
            className: "block",
            innerHTML: () => highlight(props.value, props.language),
        }),
        t.div({ className: "footnote" }, (el) => {
            if (typeof props.footnote == "function") {
                return props.footnote(el);
            }

            return props.footnote;
        }),
    );
};

function highlight(content, language) {
    content = typeof content == "string" ? content : "";

    // @see https://prismjs.com/plugins/normalize-whitespace
    content = Prism.plugins.NormalizeWhitespace.normalize(content, {
        "remove-trailing": true,
        "remove-indent": true,
        "left-trim": true,
        "right-trim": true,
    });

    return Prism.highlight(content, Prism.languages[language] || Prism.languages.js, language);
}
