window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * A vertical dragline to allow width resizing an element.
 *
 * @example
 * ```js
 * return app.components.dragline({
 *     ondragstart: (e) => {
 *         el._startWidth = el.offsetWidth;
 *     },
 *     ondragging: (e, diffX, diffY) => {
 *         el.style.width = el._startWidth + diffX + "px";
 *     },
 * });
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.dragline = function(propsArg = {}) {
    let elem;

    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        tolerance: 0,
        ondragstart: function(e) {},
        ondragstop: function(e) {},
        ondragging: function(e, diffX, diffY) {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        dragStarted: false,
    });

    let startX, startY, shiftX, shiftY;

    function addDocumentEvents() {
        document.addEventListener("touchmove", onMove);
        document.addEventListener("mousemove", onMove);
        document.addEventListener("touchend", onStop);
        document.addEventListener("mouseup", onStop);
    }

    function removeDocumentEvents() {
        document.removeEventListener("touchmove", onMove);
        document.removeEventListener("mousemove", onMove);
        document.removeEventListener("touchend", onStop);
        document.removeEventListener("mouseup", onStop);
    }

    function dragInit(e) {
        e.stopPropagation();

        startX = e.clientX;
        startY = e.clientY;
        shiftX = e.clientX - elem.offsetLeft;
        shiftY = e.clientY - elem.offsetTop;

        addDocumentEvents();
    }

    function onStop(e) {
        if (data.dragStarted) {
            e.preventDefault();
            data.dragStarted = false;
            props.ondragstop?.(e);
        }

        removeDocumentEvents();
    }

    function onMove(e) {
        let diffX = e.clientX - startX;
        let diffY = e.clientY - startY;
        let left = e.clientX - shiftX;
        let top = e.clientY - shiftY;

        if (
            !data.dragStarted
            && Math.abs(left - elem.offsetLeft) < props.tolerance
            && Math.abs(top - elem.offsetTop) < props.tolerance
        ) {
            return;
        }

        e.preventDefault();

        if (!data.dragStarted) {
            data.dragStarted = true;
            props.ondragstart?.(e);
        }

        props.ondragging?.(e, diffX, diffY);
    }

    elem = t.div({
        rid: props.rid,
        id: () => props.id,
        hidden: () => props.hidden,
        inert: () => props.inert,
        className: () => `dragline ${data.dragStarted ? "dragging" : ""} ${props.className}`,
        onmousedown: (e) => {
            if (e.button == 0) {
                dragInit(e);
            }
        },
        ontouchstart: dragInit,
        onunmount: () => {
            removeDocumentEvents();
            watchers.forEach((w) => w?.unwatch());
        },
    });

    return elem;
};
