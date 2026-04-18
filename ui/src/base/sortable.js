window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Creates a reactive single level sortable container (nested sortables are not supported).
 *
 * @example
 * ```js
 * app.components.sortable({
 *     data: () => data.list,
 *     dataItem: (item) => t.strong(null, "ID:", () => item.id),
 * })
 * ```
 *
 * @param  {Object>} [propsArg]
 * @return {Element}
 */
window.app.components.sortable = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        data: [],
        dataItem: function(item, i, parent) {
            return t.span(null, "Item " + i);
        },
        onchange: function(sortedList, fromIndex, toIndex) {},
        handle: "", // specific handle selector (if not set attached to the entire list item)
        before: undefined,
        after: undefined,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    function initSortEvents(listEl) {
        function clearDragData() {
            listEl.querySelectorAll(":scope > [data-dragstart=\"true\"]")?.forEach((item) => {
                item.dataset.dragstart = false;
            });

            listEl.querySelectorAll(":scope > [data-dragover=\"true\"]")?.forEach((item) => {
                item.dataset.dragover = false;
            });
        }

        // drag
        // ---
        listEl.addEventListener("dragstart", (e) => {
            if (props.handle && !e.target.closest(props.handle)) {
                e.preventDefault();
                return;
            }

            const child = closestChild(listEl, e.target);
            if (child) {
                child.dataset.dragstart = true;
            }
        });

        listEl.addEventListener("dragenter", (e) => {
            for (let child of listEl.children) {
                if (child.dataset.dragover) {
                    child.dataset.dragover = false;
                }
            }

            const to = closestChild(listEl, e.target);
            if (to) {
                to.dataset.dragover = true;
            }
        });

        listEl.addEventListener("dragend", (e) => {
            clearDragData();
        });

        // drop
        // ---
        // prevent default to allow drop
        // (https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/drop_event)
        listEl.addEventListener("dragover", (e) => {
            e.preventDefault();
        });
        listEl.addEventListener("drop", (e) => {
            if (!props.onchange) {
                clearDragData();
                return;
            }

            const from = listEl.querySelector(":scope > [data-dragstart=\"true\"]");
            const to = closestChild(listEl, e.target);

            clearDragData();

            if (!from || !to || to == from) {
                return;
            }

            const fromIndex = childIndex(from);
            const toIndex = childIndex(to);

            const clone = props.data.slice();
            const deleted = clone.splice(fromIndex, 1);
            clone.splice(toIndex, 0, deleted[0]);

            props.onchange(clone, fromIndex, toIndex);
        });
    }

    function childIndex(node) {
        if (!node?.parentNode) {
            return -1;
        }

        for (let i = 0; i < node.parentNode.children.length; i++) {
            if (node.parentNode.children[i] == node) {
                return i;
            }
        }

        return -1;
    }

    function closestChild(parent, node) {
        if (!node || !node.parentNode) {
            return null;
        }

        if (node.parentNode == parent) {
            return node;
        }

        return closestChild(parent, node.parentNode);
    }

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => props.className,
            onmount: (listEl) => {
                initSortEvents(listEl);
            },
            onunmount: (listEl) => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        (el) => {
            if (typeof props.before == "function") {
                return props.before(el);
            }
            return props.before;
        },
        (el) => {
            const children = [];

            for (let i = 0; i < props.data.length; i++) {
                let child = props.dataItem(props.data[i], i, el);
                if (!child) {
                    continue;
                }

                if (props.handle) {
                    const handle = child.querySelector(props.handle);
                    if (handle) {
                        handle.draggable = true;
                    }
                } else {
                    child.draggable = true;
                }

                children.push(child);
            }

            return children;
        },
        (el) => {
            if (typeof props.after == "function") {
                return props.after(el);
            }
            return props.after;
        },
    );
};
