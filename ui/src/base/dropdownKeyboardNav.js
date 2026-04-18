const optionClass = "dropdown-item";

document.addEventListener(
    "toggle",
    (e) => {
        if (e.newState != "open" || !e.target?.matches(".dropdown") || e.target.__keyboardNavRegistered) {
            return;
        }

        e.target.__keyboardNavRegistered = true;

        const dropdown = e.target;

        function onKeydown(e) {
            // remove keydown listener in case the element was removed while still open
            if (!dropdown.isConnected) {
                document.removeEventListener("keydown", onKeydown);
                return;
            }

            let optElem;
            if (document.activeElement && document.activeElement.classList.contains(optionClass)) {
                optElem = document.activeElement;
            } else {
                optElem = dropdown.querySelector("." + optionClass + ":not([hidden]):not(.disabled)");
            }

            if (!optElem) {
                return;
            }

            if (e.key == "ArrowUp") {
                e.preventDefault();

                const prevElem = firstActiveSibling(optElem, -1);
                if (optElem == document.activeElement && prevElem?.classList?.contains(optionClass)) {
                    prevElem?.focus();
                } else {
                    optElem.focus();
                }
            } else if (e.key == "ArrowDown") {
                e.preventDefault();

                const nextElem = firstActiveSibling(optElem, 1);
                if (optElem == document.activeElement && nextElem?.classList?.contains(optionClass)) {
                    nextElem?.focus();
                } else {
                    optElem.focus();
                }
            } else if (
                // a-z
                (e.keyCode >= 65 && e.keyCode <= 90)
                // 0-9
                || (e.keyCode >= 48 && e.keyCode <= 57)
            ) {
                // autofocus the only available input when start typing
                // (e.g. for search)
                const inputs = dropdown.querySelectorAll("input,textare,select");
                if (inputs.length == 1) {
                    inputs[0].focus();
                }
            }
        }

        dropdown.addEventListener("toggle", (e) => {
            if (e.newState == "open") {
                updatePopovertargetsData(e.target.id, true);
                document.addEventListener("keydown", onKeydown);
            } else {
                updatePopovertargetsData(e.target.id, false);
                document.removeEventListener("keydown", onKeydown);
            }
        });
    },
    true,
);

function updatePopovertargetsData(popoverId, state = false) {
    if (!popoverId) {
        return;
    }

    document.querySelectorAll("[popovertarget='" + popoverId + "']")
        ?.forEach((el) => el.setAttribute("data-popover-state", state));
}

function firstActiveSibling(el, dir = -1) {
    const sibling = dir < 0 ? el.previousElementSibling : el.nextElementSibling;

    if (
        sibling
        && (sibling.hidden || sibling.classList.contains("disabled") || !sibling.classList.contains(optionClass))
    ) {
        return firstActiveSibling(sibling, dir);
    }

    return sibling;
}
