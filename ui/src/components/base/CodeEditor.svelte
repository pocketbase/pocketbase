<script>
    import { onMount, createEventDispatcher } from "svelte";
    // code mirror imports
    // ---
    import {
        keymap,
        highlightSpecialChars,
        drawSelection,
        dropCursor,
        rectangularSelection,
        EditorView,
        placeholder as placeholderExt,
    } from "@codemirror/view";
    import { EditorState, Compartment } from "@codemirror/state";
    import { defaultHighlightStyle, syntaxHighlighting, bracketMatching } from "@codemirror/language";
    import { defaultKeymap, history, historyKeymap, indentWithTab } from "@codemirror/commands";
    import { searchKeymap, highlightSelectionMatches } from "@codemirror/search";
    import { closeBrackets, closeBracketsKeymap } from "@codemirror/autocomplete";
    import { javascript } from "@codemirror/lang-javascript";
    // ---

    const dispatch = createEventDispatcher();

    export let value = "";
    export let disabled = false;
    export let placeholder = "";
    export let singleLine = false;

    let editor;
    let container;
    let editableCompartment = new Compartment();
    let readOnlyCompartment = new Compartment();
    let placeholderCompartment = new Compartment();

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

        editor = new EditorView({
            parent: container,
            state: EditorState.create({
                doc: value,
                extensions: [
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
                        indentWithTab,
                        ...closeBracketsKeymap,
                        ...defaultKeymap,
                        ...searchKeymap,
                        ...historyKeymap,
                    ]),
                    EditorView.lineWrapping,
                    javascript(),
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

        return () => editor?.destroy();
    });
</script>

<div bind:this={container} class="code-editor" />
