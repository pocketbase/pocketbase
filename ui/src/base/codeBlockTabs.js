window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Code highlighted tabs component.
 *
 * @example
 * ```js
 * app.components.codeBlockTabs({
 *     tabs: [
 *         {
 *             title: "Tab 1",
 *             language: "js",
 *             value: "console.log(123)","
 *             // other codeBlock props...
 *         },
 *         ...
 *     ],
 *     historyKey: "myTabs"
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.codeBlockTabs = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        activeTabIndex: 0,
        historyKey: "",
        tabs: [], // {title, ...codeBlockProps}
        get activeTab() {
            return props.tabs[props.activeTabIndex] || props.tabs[0];
        },
    });

    const watchers = app.utils.extendStore(props, propsArg);

    watchers.push(
        watch(() => props.activeTabIndex, (newIndex, oldIndex) => {
            if (oldIndex != undefined && props.historyKey) {
                localStorage.setItem(props.historyKey, newIndex);
            }
        }),
    );

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden || !props.tabs.length,
            inert: () => props.inert,
            className: () => `code-block-tabs ${props.className}`,
            onmount: () => {
                if (props.historyKey) {
                    props.activeTabIndex = localStorage.getItem(props.historyKey) << 0;
                }
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.header(
            { className: "tabs-header" },
            () => {
                return props.tabs.map((tab, i) => {
                    return t.button(
                        {
                            type: "button",
                            className: () => `tab-item ${props.activeTabIndex == i ? "active" : ""}`,
                            onclick: () => props.activeTabIndex = i,
                        },
                        (el) => {
                            if (typeof tab.title == "function") {
                                return tab.title(el);
                            }
                            return tab.title;
                        },
                    );
                });
            },
        ),
        t.div(
            { className: "code-block-tabs-content" },
            () => {
                if (props.activeTab) {
                    return app.components.codeBlock(props.activeTab);
                }
            },
        ),
    );
};
