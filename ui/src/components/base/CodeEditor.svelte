<script>
    /**
     * This component uses Codemirror editor under the hood and its a "little heavy".
     * To allow manuall chunking it is recommended to load the component lazily!
     *
     * Example usage:
     * ```
     * <script>
     * import { onMount } from "svelte";
     *
     * let editorComponent;
     *
     * onMount(async () => {
     *     try {
     *         editorComponent = (await import("@/components/base/CodeEditor.svelte")).default;
     *     } catch (err) {
     *         console.warn(err);
     *     }
     * });
     * <//script>
     *
     * ...
     *
     * <svelte:component
     *     this={editorComponent}
     *     bind:value={value}
     *     disabled={disabled}
     *     language="html"
     * />
     * ```
     */
    import { onMount, createEventDispatcher } from "svelte";
    // code mirror imports
    // ---
    import {
        keymap,
        highlightSpecialChars,
        drawSelection,
        dropCursor,
        rectangularSelection,
        highlightActiveLineGutter,
        EditorView,
        placeholder as placeholderExt,
    } from "@codemirror/view";
    import { EditorState, Compartment } from "@codemirror/state";
    import { defaultHighlightStyle, syntaxHighlighting, bracketMatching } from "@codemirror/language";
    import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
    import { searchKeymap, highlightSelectionMatches } from "@codemirror/search";
    import {
        autocompletion,
        completionKeymap,
        closeBrackets,
        closeBracketsKeymap,
    } from "@codemirror/autocomplete";
    import { html as htmlLang } from "@codemirror/lang-html";
    import { javascript as javascriptLang } from "@codemirror/lang-javascript";
    // ---

    const dispatch = createEventDispatcher();

    export let id = "";
    export let value = "";
    export let maxHeight = null;
    export let disabled = false;
    export let placeholder = "";
    export let language = "javascript";
    export let singleLine = false;

    let editor;
    let container;
    let langCompartment = new Compartment();
    let editableCompartment = new Compartment();
    let readOnlyCompartment = new Compartment();
    let placeholderCompartment = new Compartment();

    $: if (id) {
        addLabelListeners();
    }

    $: if (editor && language) {
        editor.dispatch({
            effects: [langCompartment.reconfigure(getEditorLang())],
        });
    }

    $: if (editor && typeof disabled !== "undefined") {
        editor.dispatch({
            effects: [
                editableCompartment.reconfigure(EditorView.editable.of(!disabled)),
                readOnlyCompartment.reconfigure(EditorState.readOnly.of(disabled)),
            ],
        });

        triggerNativeChange();
    }

    $: if (editor && value != editor.state.doc.toString()) {
        editor.dispatch({
            changes: {
                from: 0,
                to: editor.state.doc.length,
                insert: value,
            },
        });
    }

    $: if (editor && typeof placeholder !== "undefined") {
        editor.dispatch({
            effects: [placeholderCompartment.reconfigure(placeholderExt(placeholder))],
        });
    }

    // Focus the editor (if inited).
    export function focus() {
        editor?.focus();
    }

    // Emulate native change event for the editor container element.
    function triggerNativeChange() {
        container?.dispatchEvent(
            new CustomEvent("change", {
                detail: { value },
                bubbles: true,
            })
        );
    }

    // Remove any attached label listeners.
    function removeLabelListeners() {
        if (!id) {
            return;
        }

        const labels = document.querySelectorAll('[for="' + id + '"]');
        for (let label of labels) {
            label.removeEventListener("click", focus);
        }
    }

    // Add `<label for="ID">...</label>` focus support.
    function addLabelListeners() {
        if (!id) {
            return;
        }

        removeLabelListeners();

        const labels = document.querySelectorAll('[for="' + id + '"]');
        for (let label of labels) {
            label.addEventListener("click", focus);
        }
    }

    // Returns the current active editor language.
    function getEditorLang() {
        return language === "html" ? htmlLang() : javascriptLang();
    }

    onMount(() => {
        const submitShortcut = {
            key: "Enter",
            run: (_) => {
                // trigger submit on enter for singleline input
                if (singleLine) {
                    dispatch("submit", value);
                }
            },
        };

        addLabelListeners();

        editor = new EditorView({
            parent: container,
            state: EditorState.create({
                doc: value,
                extensions: [
                    highlightActiveLineGutter(),
                    highlightSpecialChars(),
                    history(),
                    drawSelection(),
                    dropCursor(),
                    EditorState.allowMultipleSelections.of(true),
                    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
                    bracketMatching(),
                    closeBrackets(),
                    rectangularSelection(),
                    highlightSelectionMatches(),
                    keymap.of([
                        submitShortcut,
                        ...closeBracketsKeymap,
                        ...defaultKeymap,
                        searchKeymap.find((item) => item.key === "Mod-d"),
                        ...historyKeymap,
                        ...completionKeymap,
                    ]),
                    EditorView.lineWrapping,
                    autocompletion({
                        icons: false,
                    }),
                    langCompartment.of(getEditorLang()),
                    placeholderCompartment.of(placeholderExt(placeholder)),
                    editableCompartment.of(EditorView.editable.of(true)),
                    readOnlyCompartment.of(EditorState.readOnly.of(false)),
                    EditorState.transactionFilter.of((tr) => {
                        return singleLine && tr.newDoc.lines > 1 ? [] : tr;
                    }),
                    EditorView.updateListener.of((v) => {
                        if (!v.docChanged || disabled) {
                            return;
                        }
                        value = v.state.doc.toString();
                        triggerNativeChange();
                    }),
                ],
            }),
        });

        return () => {
            removeLabelListeners();
            editor?.destroy();
        };
    });
</script>

<div bind:this={container} class="code-editor" style:max-height={maxHeight ? maxHeight + "px" : "auto"} />
