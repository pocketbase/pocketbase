function getWidth(el) {
    let { width, boxSizing } = window.getComputedStyle(el, null);

    switch (boxSizing) {
        case "content-box":
            return parseFloat(width);
        case "border-box":
            return el.offsetWidth;
    }
}

export default function resizecolumns(table, { id }) {
    let tables = {};
    let dragging = false;
    let preDragWidth;
    let preDragX;

    function setup(id) {
        if (!tables[id]) tables[id] = {};

        for (let th of table.getElementsByTagName("th")) {
            let key = th.getAttribute("column-key");

            let width = tables[id][key]?.width ?? getWidth(th);

            th.style.width = width + "px";
            th.style.minWidth = "auto";
            th.style.maxWidth = "none";

            if (!key) continue;

            tables[id][key] = { width };

            function click(e) {
                e.stopPropagation();
            }

            function down(e) {
                e.currentTarget.setPointerCapture(e.pointerId);
                dragging = true;
                preDragWidth = getWidth(th);
                preDragX = e.pageX;
                resizeHandler.classList.add("resize-handler-active");
            }

            function move(e) {
                if (!dragging) return;
                let width = preDragWidth + e.pageX - preDragX;
                tables[id][key].width = width;
                th.style.width = width + "px";
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
    }

    function teardown() {
        for (let th of table.getElementsByTagName("th")) {
            th.style.removeProperty("width");
            th.querySelector(".resize-handler")?.remove();
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
