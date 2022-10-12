<script>
    import { loadCollections } from "@/stores/collections";
    import { createEventDispatcher, onMount } from "svelte";
    // code mirror imports
    // ---
    import ApiClient from "@/utils/ApiClient";
    import { autocompletion, closeBrackets } from "@codemirror/autocomplete";
    import { history } from "@codemirror/commands";
    import { sql, SQLite } from "@codemirror/lang-sql";
    import { bracketMatching, defaultHighlightStyle, syntaxHighlighting } from "@codemirror/language";
    import { highlightSelectionMatches } from "@codemirror/search";
    import { EditorState } from "@codemirror/state";
    import {
        drawSelection,
        dropCursor,
        EditorView,
        highlightActiveLineGutter,
        placeholder,
        rectangularSelection,
    } from "@codemirror/view";
    // ---

    const dispatch = createEventDispatcher();

    export let id = "";
    export let value = "";
    export let disabled = false;

    let editor;
    let container;

    // Focus the editor (if inited).
    export function focus() {
        editor?.focus();
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
    // Emulate native change event for the editor container element.
    function triggerNativeChange() {
        container?.dispatchEvent(
            new CustomEvent("change", {
                detail: { value },
                bubbles: true,
            })
        );
    }

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
    onMount(async () => {
        addLabelListeners();
        const schema = await ApiClient.collections.getFullList().then((collections) => {
            let kv = {};
            collections.map((collection) => {
                kv[collection.name] = [
                    "id",
                    "created",
                    "updated",
                    ...collection.schema.map((s) => ({ label: s.name })),
                ];
            });
            return kv;
        });
        editor = new EditorView({
            parent: container,
            state: EditorState.create({
                doc: value,
                extensions: [
                    sql({ dialect: SQLite, upperCaseKeywords: true, schema }),
                    autocompletion({ icons: false }),
                    highlightActiveLineGutter(),
                    history(),
                    drawSelection(),
                    dropCursor(),
                    placeholder("SELECT posts.* from posts"),
                    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
                    bracketMatching(),
                    closeBrackets(),
                    rectangularSelection(),
                    highlightSelectionMatches(),
                    EditorView.updateListener.of((v) => {
                        if (!v.docChanged || disabled) {
                            return;
                        }
                        value = v.state.doc.toString();
                        triggerNativeChange();
                    }),
                    EditorView.theme({
                        ".cm-content": {
                            minHeight: "120px",
                        },
                    }),
                ],
            }),
        });

        return () => {
            removeLabelListeners();
            editor?.destroy();
        };
    });
    loadCollections();
</script>

<div bind:this={container} class="code-editor" />
