window.app = window.app || {};
window.app.components = window.app.components || {};

const maxScale = 2;
const minScale = 0.4;
const areaPadding = 40;

/**
 * Component for rendering ERD-like draggable chart based on the specified collections.
 *
 * @example
 * ```js
 * return app.components.erd({
 *     collections: () => [collection1, collection2, collection3],
 * });
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.erd = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        collections: [],
        height: 500, // number or CSS string
        cols: 5,
        marginX: 90,
        marginY: 70,
        scale: 0.8,
        onscalechange: (newScale) => {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        activeCollection: null,
        positions: {},
        viewX: 0,
        viewY: 0,
        panStartX: 0,
        panStartY: 0,
        isPanning: false,
        isUpdating: true,
    });

    let erdEl;
    let initialScale;
    const uniqueId = "erd_" + app.utils.randomString();

    watchers.push(watch(() => props.scale, (newScale) => {
        if (newScale > maxScale) {
            props.scale = maxScale;
        }

        if (newScale < minScale) {
            props.scale = minScale;
        }

        if (!initialScale) {
            initialScale = newScale;
        }

        props.onscalechange?.(newScale);
    }));

    watchers.push(watch(() => JSON.stringify(props.collections), (_, oldHash) => {
        if (typeof oldHash == "undefined") {
            return; // initial
        }

        updateTablePositions();
    }));

    async function updateTablePositions(withRelations = true) {
        data.isUpdating = true;

        await new Promise((r) => setTimeout(r, 0));

        data.positions = {};

        const colYOffsets = new Array(props.cols).fill(areaPadding);

        erdEl.querySelectorAll(".erd-table")?.forEach((table, i) => {
            const colIndex = i % props.cols;

            data.positions[table.dataset.collectionId] = {
                x: areaPadding + colIndex * (table.clientWidth + props.marginX),
                y: colYOffsets[colIndex],
            };

            colYOffsets[colIndex] += table.clientHeight + props.marginY;
        });

        horizontalCenter();

        if (withRelations) {
            await updateRelations();
        } else {
            clearRelations();
        }

        data.isUpdating = false;
    }

    function clearRelations() {
        backPathsGroup.innerHTML = "";
        frontPathsGroup.innerHTML = "";
    }

    async function updateRelations() {
        await new Promise((r) => setTimeout(r, 0));

        clearRelations();

        const paths = [];

        for (const collection of props.collections) {
            const fields = collection.fields || [];
            for (const field of fields) {
                if (field.type != "relation") {
                    continue;
                }

                const fromEl = erdEl?.querySelector(`[data-collection-id="${collection.id}"]`)
                    ?.querySelector(`[data-field-name="${field.name}"]`);

                const toEl = erdEl?.querySelector(`[data-collection-id="${field.collectionId}"]`)
                    ?.querySelector(`[data-field-name="id"]`);

                paths.push(createPath(fromEl, toEl, props.scale, collection.id, field.collectionId));
            }
        }

        if (paths.length) {
            backPathsGroup.append(...paths);
        }
    }

    function horizontalCenter() {
        const containerWidth = erdEl.clientWidth || 0;
        const tableWidth = erdEl.querySelector(".erd-table")?.offsetWidth || 0;
        const areaWidth = 2 * areaPadding + (props.cols * (tableWidth + props.marginX)) - props.marginX;

        data.viewY = 0;
        data.viewX = (containerWidth - areaWidth * props.scale) / 2;
    }

    function resetScale() {
        props.scale = initialScale || 1;
        horizontalCenter();
    }

    function isHighlighted(collection) {
        if (!data.activeCollection) {
            return false;
        }

        if (data.activeCollection.id == collection.id) {
            return true;
        }

        const relFromField = data.activeCollection.fields?.find((f) =>
            f.type == "relation" && f.collectionId == collection.id
        );
        if (relFromField) {
            return true;
        }

        const relToField = collection.fields?.find((f) =>
            f.type == "relation" && f.collectionId == data.activeCollection.id
        );
        return !!relToField;
    }

    // maintain 2 svg layers because z-index on paths can't escape the svg
    //
    // (for the path commands see https://developer.mozilla.org/en-US/docs/Web/SVG/Reference/Attribute/d#path_commands)
    // ---
    let svgBack = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svgBack.classList.add("erd-paths", "back");
    svgBack.innerHTML = `
        <defs>
          <marker id="${uniqueId}_arrow1" class="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto" fill="context-stroke">
            <path d="M 0 0 L 10 5 L 0 10 z" />
          </marker>
        </defs>
        <g class="paths-group" marker-end="url(#${uniqueId}_arrow1)"></g>
    `;
    let backPathsGroup = svgBack.querySelector(".paths-group");

    let svgFront = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svgFront.classList.add("erd-paths", "front");
    svgFront.innerHTML = `
        <defs>
          <marker id="${uniqueId}_arrow2" class="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto" fill="context-stroke">
            <path d="M 0 0 L 10 5 L 0 10 z" />
          </marker>
        </defs>
        <g class="paths-group" marker-end="url(#${uniqueId}_arrow2)"></g>
    `;
    let frontPathsGroup = svgFront.querySelector(".paths-group");
    // ---

    erdEl = t.div(
        {
            tabIndex: -1,
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () =>
                // dprint-ignore
                `erd ${props.className} ${data.isUpdating ? "updating" : ""} ${data.isPanning ? "panning" : ""} ${data.activeCollection ? "active" : ""}`,
            onkeydown: (e) => {
                if ((e.ctrlKey || e.metaKey) && e.key == "0") {
                    resetScale();
                }
            },
            onmount: async (el) => {
                // zoom
                el.addEventListener("wheel", (e) => {
                    e.preventDefault();

                    const rect = el.getBoundingClientRect();
                    const mouseX = e.clientX - rect.left;
                    const mouseY = e.clientY - rect.top;

                    const layoutX = (mouseX - data.viewX) / props.scale;
                    const layoutY = (mouseY - data.viewY) / props.scale;

                    const newScale = Math.min(Math.max(-e.deltaY * 0.001 + props.scale, minScale), maxScale);

                    data.viewX = mouseX - layoutX * newScale;
                    data.viewY = mouseY - layoutY * newScale;
                    props.scale = newScale;
                });

                // pan
                el.addEventListener("pointerdown", (e) => {
                    if (e.buttons != 1) {
                        return;
                    }

                    data.isPanning = true;
                    data.panStartX = e.clientX - data.viewX;
                    data.panStartY = e.clientY - data.viewY;
                });
                el._ondragging = function(e) {
                    if (!data.isPanning) {
                        return;
                    }

                    data.viewX = e.clientX - data.panStartX;
                    data.viewY = e.clientY - data.panStartY;
                };
                el._ondragstop = function() {
                    data.isPanning = false;
                };
                // note: attach mouse end events to window to allow panning when outside the element
                window.addEventListener("pointermove", el._ondragging);
                window.addEventListener("pointerup", el._ondragstop);

                updateTablePositions();
            },
            onunmount: (el) => {
                window.removeEventListener("pointermove", el._ondragging);
                window.removeEventListener("pointerup", el._ondragstop);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            {
                className: "erd-area",
                style: () => `transform: translate(${data.viewX}px, ${data.viewY}px) scale(${props.scale});`,
            },
            svgBack,
            t.div(
                { className: "erd-tables" },
                () => {
                    return props.collections.map((collection) => {
                        return t.div(
                            {
                                // dpint-ignore
                                style: () =>
                                    `left:${data.positions[collection.id]?.x || 0}px;top:${
                                        data.positions[collection.id]?.y || 0
                                    }px`,
                                className: () =>
                                    `erd-table type-${collection.type} ${isHighlighted(collection) ? "active" : ""} ${
                                        collection.system ? "system" : ""
                                    }`,
                                "html-data-collection-id": () => collection.id,
                                "html-data-collection-name": () => collection.name,
                                onmouseenter: () => {
                                    backPathsGroup.querySelectorAll(`[data-to="${collection.id}"]`)?.forEach(
                                        (child) => {
                                            child.classList.add("active-to");
                                            frontPathsGroup.append(child);
                                        },
                                    );
                                    backPathsGroup.querySelectorAll(`[data-from="${collection.id}"]`)?.forEach(
                                        (child) => {
                                            child.classList.add("active-from");
                                            frontPathsGroup.append(child);
                                        },
                                    );

                                    data.activeCollection = collection;
                                },
                                onmouseleave: () => {
                                    for (const child of frontPathsGroup.children) {
                                        child.classList.remove("active-from", "active-to");
                                    }
                                    backPathsGroup.append(...frontPathsGroup.children);

                                    data.activeCollection = null;
                                },
                            },
                            t.div(
                                { className: "erd-table-row header" },
                                () => collection.name,
                            ),
                            () => {
                                return collection.fields?.map((field) => {
                                    return t.div(
                                        {
                                            className: `erd-table-row type-${field.type} ${
                                                field.primaryKey ? "primary-key" : ""
                                            }`,
                                            "html-data-field-id": () => field.id,
                                            "html-data-field-name": () => field.name,
                                        },
                                        t.i({
                                            ariaHidden: true,
                                            title: () => field.type,
                                            className: () =>
                                                `field-icon ${
                                                    app.fieldTypes[field.type].icon || app.utils.fallbackFieldIcon
                                                }`,
                                        }),
                                        t.span({ className: "field-name" }, () => field.name),
                                        () => {
                                            if (field.hidden) {
                                                return t.span(
                                                    { className: "label danger field-hidden-label" },
                                                    "Hidden",
                                                );
                                            }
                                        },
                                        () => {
                                            if (typeof field.maxSelect != "undefined") {
                                                return t.span(
                                                    { className: "meta" },
                                                    field.maxSelect > 1 ? "multiple" : "single",
                                                );
                                            }
                                        },
                                    );
                                });
                            },
                        );
                    });
                },
            ),
            svgFront,
        ),
        t.nav(
            {
                className: "erd-nav",
                onmousedown: (e) => {
                    e.stopImmediatePropagation();
                },
            },
            t.button(
                {
                    type: "button",
                    className: "btn sm circle secondary",
                    title: "Zoom in",
                    onclick: () => {
                        props.scale += 0.05;
                    },
                },
                t.i({ className: "ri-add-line", ariaHidden: true }),
            ),
            t.button(
                {
                    type: "button",
                    className: "btn sm circle secondary",
                    title: "Zoom out",
                    onclick: () => {
                        props.scale -= 0.05;
                    },
                },
                t.i({ className: "ri-subtract-line", ariaHidden: true }),
            ),
        ),
    );

    return erdEl;
};

// creates a new svg path line in svgEl to visually connect el1 and el2.
function createPath(
    el1,
    el2,
    scale = 1,
    fromCollectionId = "",
    toCollectionId = "",
    extraSpacing = 2,
) {
    if (!el1 || !el2) {
        return;
    }

    const workspaceRect = el1.closest(".erd-area").getBoundingClientRect();

    const r1 = el1.getBoundingClientRect();
    const r2 = el2.getBoundingClientRect();

    extraSpacing *= scale;

    let y1 = (r1.top - workspaceRect.top) + r1.height / 2;
    let y2 = (r2.top - workspaceRect.top) + r2.height / 2;

    let x1, x2;
    if (r1.left < r2.left) {
        // left -> right
        x1 = (r1.left - workspaceRect.left) + r1.width + extraSpacing;
        x2 = (r2.left - workspaceRect.left) - extraSpacing;
    } else if (r1.left > r2.left) {
        // right <- left
        x1 = (r1.left - workspaceRect.left) - extraSpacing;
        x2 = (r2.left - workspaceRect.left) + r2.width + extraSpacing;
    } else {
        // within the same column
        x1 = (r1.left - workspaceRect.left) - extraSpacing;
        x2 = (r2.left - workspaceRect.left) - extraSpacing;
    }

    // rescale
    x1 /= scale;
    x2 /= scale;
    y1 /= scale;
    y2 /= scale;

    // mid point for the line squirish/orthogonal break
    let midX = x1 + (x2 - x1) / 2;
    if (x1 == x2) {
        midX -= 20; // offset slightly to prevent overlap when the 2 elements are in the same column
    }

    const d = `M ${x1} ${y1}
        L ${midX} ${y1}
        L ${midX} ${y2}
        L ${x2} ${y2}`;

    const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
    path.setAttribute("class", "relation-path");
    path.setAttribute("data-from", fromCollectionId || "");
    path.setAttribute("data-to", toCollectionId || "");
    path.setAttribute("d", d);

    return path;
}
