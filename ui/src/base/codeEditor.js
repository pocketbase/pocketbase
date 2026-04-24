window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Basic code editor element with syntax highlight support.
 * For static code visualization use `app.components.codeBlock({ ... })`.
 *
 * @example
 * ```js
 * app.components.codeEditor({
 *     language: "html",
 *     value: () => data.myCode, // data is some store() instance
 *     singleLine: true,
 *     placeholder: "Type your html here...",
 *     oninput: (val) => {
 *         data.myCode = val
 *     },
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.codeEditor = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        name: undefined,
        className: "",
        value: "",
        language: "js", // see Prism.languages
        placeholder: "",
        disabled: false,
        required: false,
        singleLine: false,
        // [
        //     "value",
        //     {value, label},
        // ]
        autocomplete: undefined, // Array<string|Object> | function(word): Array<string|Object>,
        // ---
        oninput: function(val) {},
        onfocus: function(val) {},
        onblur: function(val) {},
    });

    const extendWatchers = app.utils.extendStore(props, propsArg, "autocomplete");

    let dropdown;
    let visibilityObserver;
    let isFieldVisible = true;

    function openAutocompleteDropdown(items) {
        closeAutocompleteDropdown();

        dropdown = t.div(
            {
                className: "dropdown autocomplete code-editor-dropdown",
                onmount: (el) => {
                    el._updatePosition = () => {
                        if (!isFieldVisible) {
                            closeAutocompleteDropdown();
                        } else {
                            updateDropdownPosition(dropdown);
                        }
                    };
                    el._closeOnEsc = (e) => {
                        if (e.key == "Escape") {
                            e.preventDefault();
                            closeAutocompleteDropdown();
                        }
                    };
                    window.addEventListener("scroll", el._updatePosition, true);
                    window.addEventListener("resize", el._updatePosition);
                    window.addEventListener("keydown", el._closeOnEsc);
                    el._updatePosition();
                },
                onunmount: (el) => {
                    if (el) {
                        window.removeEventListener("scroll", el._updatePosition, true);
                        window.removeEventListener("resize", el._updatePosition);
                        window.removeEventListener("keydown", el._closeOnEsc);
                    }
                },
            },
            items,
        );

        document.body.appendChild(dropdown);

        // track editor field visibility to hide the dropdown when
        // not in the view port to avoid overflow issues
        if (editorContent) {
            visibilityObserver?.disconnect();
            visibilityObserver = new IntersectionObserver(
                ([entry]) => {
                    isFieldVisible = entry.isIntersecting;
                },
                {
                    root: null,
                    threshold: 0.1,
                },
            );
            visibilityObserver.observe(editorContent);
        }
    }

    function closeAutocompleteDropdown() {
        if (dropdown) {
            dropdown.remove();
            dropdown = null;
        }

        if (visibilityObserver) {
            visibilityObserver.disconnect();
            visibilityObserver = null;
        }

        isFieldVisible = true;
    }

    function updateValue(newVal) {
        props.value = newVal;
        props.oninput?.(newVal);
        editorContent.dispatchEvent(new CustomEvent("change", { detail: newVal }));
    }

    let isCtrlOrCmdKey = false;

    let valueWatcher;

    // note1: use contenteditable so that we can call getBoundingClientRect on the selected text
    // note2: getSelection also doesn't seem to work in Firefox for textarea and inputs
    const editorContent = t.div({
        contentEditable: () => (props.disabled ? false : "plaintext-only"),
        tabIndex: 0,
        spellcheck: false,
        autocorrect: false,
        autocomplete: "off",
        autocapitalize: "off",
        role: "textbox",
        className: "editor-content",
        "html-data-placeholder": () => props.placeholder,
        onmount: (el) => {
            // auto change change textContent only if it props.value was
            // changed externally to preserve the focus and caret position
            valueWatcher?.unwatch();
            valueWatcher = watch(
                () => props.value,
                (value) => {
                    if (value != editorContent.textContent) {
                        editorContent.textContent = value;
                        closeAutocompleteDropdown();
                    }
                },
            );
        },
        onunmount: (el) => {
            valueWatcher?.unwatch();
            closeAutocompleteDropdown();
        },
        onfocus: () => {
            props.onfocus?.(props.value);
        },
        onblur: (e) => {
            // not blurred because of dropdown click
            if (dropdown && !dropdown.contains(e.relatedTarget)) {
                closeAutocompleteDropdown();
            }

            props.onblur?.(props.value);
        },
        oninput: (e) => {
            closeAutocompleteDropdown();

            updateValue(editorContent.textContent);

            if (!props.value?.length) {
                editorContent.textContent = ""; // ensure that no comments, br, etc. tags are left
                return;
            }

            if (!editorContent?.isConnected) {
                return;
            }

            const pos = getCaretPos(editorContent);

            const match = getWord(props.value, pos);

            if (
                !match.word.length
                // don't show suggestions in case the cursor is at the
                // beginning of an already typed word
                || pos == match.start
            ) {
                return;
            }

            let suggestions = [];
            if (typeof props.autocomplete == "function") {
                suggestions = props.autocomplete(match.word) || [];
            } else if (!app.utils.isEmpty(props.autocomplete)) {
                const wordLowercased = match.word.toLowerCase();
                suggestions = props.autocomplete.filter((item) => {
                    if (typeof item == "object") {
                        item = item?.value;
                    }

                    item = item?.toLowerCase();

                    return item && item != wordLowercased && item.includes(wordLowercased);
                });
            }

            if (!suggestions?.length) {
                return;
            }

            openAutocompleteDropdown(() => {
                return suggestions.map((suggestion, i) => {
                    return t.button({
                        type: "button",
                        className: `dropdown-item ${i == 0 ? "active" : ""}`,
                        textContent: suggestion.label || suggestion.value || suggestion,
                        onclick: (e) => {
                            e.preventDefault();

                            editorContent.focus();

                            const word = suggestion.value || suggestion;

                            // note: replacing the text doesn't preserve the native "undo" history
                            // (document.execCommand is being deprecated)
                            editorContent.textContent = editorContent.textContent.substring(0, match.start)
                                + word
                                + editorContent.textContent.substring(match.end + 1);

                            updateValue(editorContent.textContent);

                            try {
                                window
                                    .getSelection()
                                    .setPosition(editorContent.childNodes[0], match.start + word.length);
                            } catch (err) {
                                console.warn("failed to set caret position", err);
                            }

                            closeAutocompleteDropdown();
                        },
                    });
                });
            });
        },
        onkeydown: (e) => {
            isCtrlOrCmdKey = e.ctrlKey || e.metaKey;

            // autocomplete nav
            // -------------------------------------------------------

            if ((e.key == "Enter" || e.key == "Tab") && dropdown?.isConnected) {
                e.preventDefault();
                dropdown.querySelector(".dropdown-item.active")?.click();
                return;
            }

            if (e.key == "ArrowUp" && dropdown?.isConnected) {
                e.preventDefault();

                const currentActive = dropdown.querySelector(".dropdown-item.active");
                if (currentActive?.previousElementSibling) {
                    currentActive.classList.remove("active");
                    currentActive.previousElementSibling.classList.add("active");
                    currentActive.previousElementSibling.scrollIntoView(false);
                }

                return;
            }

            if (e.key == "ArrowDown" && dropdown?.isConnected) {
                e.preventDefault();

                const currentActive = dropdown.querySelector(".dropdown-item.active");
                if (currentActive?.nextElementSibling) {
                    currentActive.classList.remove("active");
                    currentActive.nextElementSibling.classList.add("active");
                    currentActive.nextElementSibling.scrollIntoView(false);
                }

                return;
            }

            // editor shortcuts
            // -------------------------------------------------------

            if (isCtrlOrCmdKey && e.key.toLowerCase() == "l") {
                e.preventDefault();
                selectLine(editorContent);
                return;
            }

            if (isCtrlOrCmdKey && e.key.toLowerCase() == "d") {
                e.preventDefault();
                selectWord(editorContent);
                return;
            }

            if (!props.singleLine && e.key == "Tab") {
                e.preventDefault();
                const selection = window.getSelection();
                if (!selection) {
                    return;
                }

                // -1 tab level
                if (e.shiftKey) {
                    selection.modify("extend", "backward", "character");
                    if (selection.toString()[0] == "\t") {
                        selection.deleteFromDocument();
                        updateValue(editorContent.textContent);
                    } else {
                        // check ahead and restore
                        selection.modify("extend", "forward", "character");
                        if (selection.toString()[0] == "\t") {
                            selection.deleteFromDocument();
                            updateValue(editorContent.textContent);
                        }
                    }

                    return;
                }

                // +1 tab level
                const range = selection.getRangeAt(0);
                if (range) {
                    range.deleteContents();
                    range.insertNode(document.createTextNode("\t"));
                    range.collapse();
                    updateValue(editorContent.textContent);
                }

                return;
            }

            // simulate single-line enter press
            if (props.singleLine && e.key == "Enter") {
                e.preventDefault();
                hiddenSubmit.click();
                return;
            }
        },
        onscroll: () => {
            closeAutocompleteDropdown();

            if (highlightOverlay) {
                highlightOverlay.scrollLeft = editorContent.scrollLeft;
                highlightOverlay.scrollTop = editorContent.scrollTop;
            }
        },
    });

    const highlightOverlay = t.div({
        className: "highlight-overlay",
        innerHTML: () => highlight(props.value, props.language),
        onscroll: () => {
            if (editorContent) {
                editorContent.scrollLeft = highlightOverlay.scrollLeft;
                editorContent.scrollTop = highlightOverlay.scrollTop;
            }
        },
    });

    const hiddenSubmit = t.button({
        type: "submit",
        className: "hidden",
    });

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            inert: () => props.inert,
            hidden: () => props.hidden,
            "html-name": () => props.name,
            "html-required": () => props.required || undefined,
            // dprint-ignore
            className: () => `input code-editor ${props.className} ${props.disabled ? "disabled" : ""} ${props.singleLine ? "single-line" : ""}`,
            onclick: () => {
                editorContent?.focus();
            },
            onunmount: () => {
                extendWatchers?.forEach((w) => w?.unwatch());
            },
        },
        t.div({ className: "code-editor-container" }, editorContent, highlightOverlay, hiddenSubmit),
    );
};

const highlightThreshold = 500;

function highlight(content, language) {
    content = typeof content == "string" ? content : "";
    if (!content) {
        return "";
    }

    if (
        !Prism.languages[language]
        // fallback to plain to avoid performance issues with large text blocks
        || content.length > highlightThreshold
    ) {
        language = "plain";
    }

    return Prism.highlight(content, Prism.languages[language], language);
}

const wordCharRegex = new RegExp(/[\p{Alphabetic}\p{Number}_@:\."'{}]/, "u");

function getWord(value, caretPos) {
    let start = caretPos;
    for (let i = caretPos - 1; i >= 0; i--) {
        if (!wordCharRegex.test(value[i])) {
            break;
        }
        start = i;
    }

    let end = start;
    for (let i = caretPos - 1; i < value.length; i++) {
        if (!wordCharRegex.test(value[i])) {
            break;
        }
        end = i;
    }

    return {
        word: value.substring(start, end + 1),
        start: start,
        end: end,
    };
}

function selectLine() {
    const selection = window.getSelection();
    selection?.modify("move", "forward", "lineboundary");
    selection?.modify("extend", "backward", "lineboundary");
}

function selectWord() {
    const selection = window.getSelection();
    selection?.modify("move", "forward", "word");
    selection?.modify("extend", "backward", "word");
}

function getCaretPos(editorContent) {
    const selection = window.getSelection();

    // new line adds a new text node which resets the selection counter
    // so we have to add them as offset
    let offset = 0;
    for (let node of editorContent.childNodes) {
        if (node == selection.focusNode) {
            break;
        } else {
            offset += node.length;
        }
    }

    return offset + selection.focusOffset;
}

function updateDropdownPosition(dropdown) {
    const targetRect = window.getSelection()?.getRangeAt(0)?.getBoundingClientRect();
    if (!targetRect || !dropdown) {
        return false;
    }

    if (targetRect.top < 0) {
        dropdown.classList.add("hidden");
        return;
    }

    dropdown.classList.remove("hidden");

    // reset
    dropdown.style.left = "0px";
    dropdown.style.top = "0px";

    const dropdownHeight = dropdown.offsetHeight;
    const dropdownWidth = dropdown.offsetWidth;

    let left = targetRect.left - 5;
    let top = targetRect.top + targetRect.height;

    // show on top if it cannot fit below the parent
    if (top + dropdownHeight > document.documentElement.clientHeight) {
        top = Math.max(targetRect.top - dropdownHeight, 0);
    }

    // align from the right edge if overflow
    if (left + dropdownWidth > document.documentElement.clientWidth) {
        left = Math.max(document.documentElement.clientWidth - dropdownWidth, 0);
    }

    dropdown.style.left = left + "px";
    dropdown.style.top = top + "px";
}
