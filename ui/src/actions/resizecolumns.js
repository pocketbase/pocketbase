function getWidthInfo(el) {
    let { width, boxSizing, paddingLeft, paddingRight, borderLeftWidth, borderRightWidth } =
        window.getComputedStyle(el, null);

    width = parseFloat(width);

    let nonContentWidth =
        parseFloat(paddingLeft) +
        parseFloat(paddingRight) +
        parseFloat(borderLeftWidth) +
        parseFloat(borderRightWidth);

    let contentWidth = boxSizing === "content-box" ? width : width - nonContentWidth;

    return { boxSizing, contentWidth, nonContentWidth, width };
}

export default function resizecolumns(table, { id }) {
    let tables = {};
    let initialTableWidth;
    let initialLastThWidth;
    let lastTh;
    let lastCol;
    let dragging = false;
    let preDragWidth;
    let preDragMinWidth;
    let preDragX;

    function keepMinWidth(id) {
        let widthWithoutLastCol = 0;
        for (let col of Object.values(tables[id])) {
            if (col !== lastCol && col.active) {
                widthWithoutLastCol +=
                    col.boxSizing === "border-box" ? col.width : col.width + col.nonContentWidth;
            }
        }

        if (widthWithoutLastCol < initialTableWidth - initialLastThWidth) {
            lastCol.width = initialTableWidth - widthWithoutLastCol;
            if (lastCol.boxSizing === "content-box") lastCol.width -= lastCol.nonContentWidth;
            lastTh.style.width = Math.max(initialLastThWidth, lastCol.width) + "px";
        }
    }

    function setup(id) {
        if (!tables[id]) tables[id] = {};

        initialTableWidth = getWidthInfo(table).width;
        lastTh = table.querySelector("th:last-child");
        lastCol = tables[id][lastTh.getAttribute("column-key")] ?? getWidthInfo(lastTh);
        initialLastThWidth = lastCol.width;

        for (let th of table.getElementsByTagName("th")) {
            let key = th.getAttribute("column-key");

            tables[id][key] ??= getWidthInfo(th);

            th.style.width = tables[id][key].width + "px";
            th.style.minWidth = "auto";
            th.style.maxWidth = "none";

            tables[id][key].active = true;

            if (th.hasAttribute("column-noresize")) continue;

            function click(e) {
                e.stopPropagation();
            }

            function down(e) {
                e.currentTarget.setPointerCapture(e.pointerId);
                dragging = true;
                preDragWidth = tables[id][key].width;
                preDragMinWidth =
                    tables[id][key].boxSizing === "content-box" ? 0 : tables[id][key].nonContentWidth;
                preDragX = e.pageX;
                resizeHandler.classList.add("resize-handler-active");
            }

            function move(e) {
                if (!dragging) return;
                let width = Math.max(preDragMinWidth, preDragWidth + e.pageX - preDragX);
                tables[id][key].width = width;
                th.style.width = width + "px";

                keepMinWidth(id);
            }

            function up(e) {
                if (!dragging) return;
                dragging = false;
                e.currentTarget.releasePointerCapture(e.pointerId);
                resizeHandler.classList.remove("resize-handler-active");
            }

            let resizeHandler = document.createElement("div");
            resizeHandler.className = "resize-handler";
            let resizeIndicator = document.createElement("div");
            resizeIndicator.className = "resize-indicator";

            resizeHandler.appendChild(resizeIndicator);
            resizeHandler.addEventListener("click", click);
            resizeHandler.addEventListener("pointerdown", down);
            resizeHandler.addEventListener("pointermove", move);
            resizeHandler.addEventListener("pointerup", up);

            th.appendChild(resizeHandler);
        }

        table.style.tableLayout = "fixed";
        table.style.minWidth = "1px";
        table.style.width = "1px";

        keepMinWidth(id);
    }

    function teardown() {
        for (let th of table.getElementsByTagName("th")) {
            th.style.removeProperty("width");
            th.querySelector(".resize-handler")?.remove();
        }
        for (let table of Object.values(tables)) {
            for (let col of Object.values(table)) {
                col.active = false;
            }
        }
        table.style.removeProperty("table-layout");
        table.style.removeProperty("min-width");
        table.style.removeProperty("width");
    }

    setup(id);

    return {
        update({ id }) {
            teardown();
            setup(id);
        },
        destroy() {
            teardown();
        },
    };
}
