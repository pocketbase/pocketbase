const tolerance = 5;

const tooltip = t.div({
    popover: "manual",
    className: "pb-tooltip",
});
document.body.appendChild(tooltip);

function updateTooltipPosition(node, position) {
    let nodeRect = node.getBoundingClientRect();

    tooltip.setAttribute("data-position", position);

    // reset tooltip position
    tooltip.style.top = "0px";
    tooltip.style.left = "0px";

    // note: doesn't use getBoundingClientRect() here because the
    // tooltip could be animated/scaled/transformed and we need the real size
    let tooltipHeight = tooltip.offsetHeight;
    let tooltipWidth = tooltip.offsetWidth;

    let top = 0;
    let left = 0;

    // calculate tooltip position based
    if (position == "left") {
        top = nodeRect.top + nodeRect.height / 2 - tooltipHeight / 2;
        left = nodeRect.left - tooltipWidth - tolerance;
    } else if (position == "right") {
        top = nodeRect.top + nodeRect.height / 2 - tooltipHeight / 2;
        left = nodeRect.right + tolerance;
    } else if (position == "top") {
        top = nodeRect.top - tooltipHeight - tolerance;
        left = nodeRect.left + nodeRect.width / 2 - tooltipWidth / 2;
    } else if (position == "top-left") {
        top = nodeRect.top - tooltipHeight - tolerance;
        left = nodeRect.left;
    } else if (position == "top-right") {
        top = nodeRect.top - tooltipHeight - tolerance;
        left = nodeRect.right - tooltipWidth;
    } else if (position == "bottom-left") {
        top = nodeRect.top + nodeRect.height + tolerance;
        left = nodeRect.left;
    } else if (position == "bottom-right") {
        top = nodeRect.top + nodeRect.height + tolerance;
        left = nodeRect.right - tooltipWidth;
    } else {
        // bottom
        top = nodeRect.top + nodeRect.height + tolerance;
        left = nodeRect.left + nodeRect.width / 2 - tooltipWidth / 2;
    }

    // right edge boundary
    if (left + tooltipWidth > document.documentElement.clientWidth) {
        left = document.documentElement.clientWidth - tooltipWidth;
    }

    // left edge boundary
    left = left >= 0 ? left : 0;

    // bottom edge boundary
    if (top + tooltipHeight > document.documentElement.clientHeight) {
        top = document.documentElement.clientHeight - tooltipHeight;
    }

    // top edge boundary
    top = top >= 0 ? top : 0;

    tooltip.style.top = top + "px";
    tooltip.style.left = left + "px";
}

function hideTooltip() {
    tooltip.hidePopover();
}

function showTooltip(node, text, position) {
    if (!node || !text) {
        hideTooltip();
        return;
    }

    tooltip.showPopover();
    tooltip.textContent = text;
    updateTooltipPosition(node, position);
}

document.body.addEventListener("mouseleave", () => {
    hideTooltip();
});

function tooltipAction(textOrFunc, position = "top") {
    return (el) => {
        if (!el._tooltipText) {
            el._tooltipText = store({
                value: "",
            });

            let tooltipTextWatcher;

            function showEventHandler() {
                tooltipTextWatcher?.unwatch();
                tooltipTextWatcher = watch(
                    () => el._tooltipText.value,
                    async (result) => {
                        showTooltip(el, result, position);
                    },
                );
            }

            async function hideEventHandler() {
                tooltipTextWatcher?.unwatch();
                tooltipTextWatcher = null;

                hideTooltip();
            }

            el.addEventListener("mouseenter", showEventHandler);
            el.addEventListener("focusin", showEventHandler);
            el.addEventListener("mouseleave", hideEventHandler);
            el.addEventListener("focusout", hideEventHandler);
            el.addEventListener("blur", hideEventHandler);

            const originalOnunmount = el.onunmount;
            el.onunmount = (el) => {
                tooltipTextWatcher?.unwatch();

                el._tooltipText = null;

                el?.removeEventListener("mouseenter", showEventHandler);
                el?.removeEventListener("focusin", showEventHandler);
                el?.removeEventListener("mouseleave", hideEventHandler);
                el?.removeEventListener("focusout", hideEventHandler);
                el?.removeEventListener("blur", hideEventHandler);

                originalOnunmount(el);
            };
        }

        if (typeof textOrFunc == "function") {
            el._tooltipText.value = textOrFunc();
        } else {
            el._tooltipText.value = textOrFunc;
        }

        return el._tooltipText.value;
    };
}

window.app = window.app || {};
window.app.attrs = window.app.attrs || {};
window.app.attrs.tooltip = tooltipAction;
