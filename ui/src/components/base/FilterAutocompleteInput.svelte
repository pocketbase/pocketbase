<script>
    /**
     * @todo consider combining with the CodeEditor component.
     *
     * This component uses Codemirror editor under the hood and its a "little heavy".
     * To allow manuall chunking it is recommended to load the component lazily!
     *
     * Example usage:
     * ```
     * <script>
     * import { onMount } from "svelte";
     *
     * let inputComponent;
     *
     * onMount(async () => {
     *     try {
     *         inputComponent = (await import("@/components/base/FilterAutocompleteInput.svelte")).default;
     *     } catch (err) {
     *         console.warn(err);
     *     }
     * });
     * <//script>
     *
     * ...
     *
     * <svelte:component
     *     this={inputComponent}
     *     bind:value={value}
     *     baseCollection={baseCollection}
     *     disabled={disabled}
     * />
     * ```
     */
    import { onMount, createEventDispatcher } from "svelte";
    import { collections } from "@/stores/collections";
    import CommonHelper from "@/utils/CommonHelper";
    import AutocompleteWorker from "@/autocomplete.worker.js?worker";
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
    import {
        defaultHighlightStyle,
        syntaxHighlighting,
        bracketMatching,
        StreamLanguage,
        syntaxTree,
    } from "@codemirror/language";
    import { defaultKeymap, indentWithTab, history, historyKeymap } from "@codemirror/commands";
    import { searchKeymap, highlightSelectionMatches } from "@codemirror/search";
    import {
        autocompletion,
        completionKeymap,
        closeBrackets,
        closeBracketsKeymap,
    } from "@codemirror/autocomplete";
    import { simpleMode } from "@codemirror/legacy-modes/mode/simple-mode";
    // ---

    const dispatch = createEventDispatcher();

    export let id = "";
    export let value = "";
    export let disabled = false;
    export let placeholder = "";
    export let baseCollection = null;
    export let singleLine = false;
    export let extraAutocompleteKeys = []; // eg. ["test1", "test2"]
    export let disableRequestKeys = false;
    export let disableCollectionJoinKeys = false;

    let editor;
    let container;
    let oldDisabledState = disabled;
    let langCompartment = new Compartment();
    let editableCompartment = new Compartment();
    let readOnlyCompartment = new Compartment();
    let placeholderCompartment = new Compartment();
    let autocompleteWorker = new AutocompleteWorker();

    let cachedRequestKeys = [];
    let cachedCollectionJoinKeys = [];
    let cachedBaseKeys = [];
    let baseKeysChangeHash = "";
    let oldBaseKeysChangeHash = "";

    $: baseKeysChangeHash = getCollectionKeysChangeHash(baseCollection);

    $: if (
        !disabled &&
        (oldBaseKeysChangeHash != baseKeysChangeHash ||
            disableRequestKeys !== -1 ||
            disableCollectionJoinKeys !== -1)
    ) {
        oldBaseKeysChangeHash = baseKeysChangeHash;
        refreshCachedKeys();
    }

    $: if (id) {
        addLabelListeners();
    }

    $: if (editor && baseCollection?.fields) {
        editor.dispatch({
            effects: [langCompartment.reconfigure(ruleLang())],
        });
    }

    $: if (editor && oldDisabledState != disabled) {
        editor.dispatch({
            effects: [
                editableCompartment.reconfigure(EditorView.editable.of(!disabled)),
                readOnlyCompartment.reconfigure(EditorState.readOnly.of(disabled)),
            ],
        });
        oldDisabledState = disabled;
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

    // Refresh the cached autocomplete keys.
    // ---
    let refreshDebounceId = null;

    autocompleteWorker.onmessage = (e) => {
        cachedBaseKeys = e.data.baseKeys || [];
        cachedRequestKeys = e.data.requestKeys || [];
        cachedCollectionJoinKeys = e.data.collectionJoinKeys || [];
    };

    function refreshCachedKeys() {
        clearTimeout(refreshDebounceId);
        refreshDebounceId = setTimeout(() => {
            autocompleteWorker.postMessage({
                baseCollection: baseCollection,
                collections: concatWithBaseCollection($collections),
                disableRequestKeys: disableRequestKeys,
                disableCollectionJoinKeys: disableCollectionJoinKeys,
            });
        }, 250);
    }
    // ---

    // Return a collection keys hash string that can be used to compare with previous states.
    function getCollectionKeysChangeHash(collection) {
        return JSON.stringify([collection?.name, collection?.type, collection?.fields]);
    }

    // Merge the base collection in a new list with the provided collections.
    function concatWithBaseCollection(collections) {
        let copy = collections.slice();

        if (baseCollection) {
            CommonHelper.pushOrReplaceByKey(copy, baseCollection, "id");
        }

        return copy;
    }

    // Emulate native change event for the editor container element.
    function triggerNativeChange() {
        container?.dispatchEvent(
            new CustomEvent("change", {
                detail: { value },
                bubbles: true,
            }),
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

    // Returns an array with all the supported keys.
    function getAllKeys(includeRequestKeys = true, includeCollectionJoinKeys = true) {
        let result = [].concat(extraAutocompleteKeys);

        // add base keys
        result = result.concat(cachedBaseKeys || []);

        // add @request.* keys
        if (includeRequestKeys) {
            result = result.concat(cachedRequestKeys || []);
        }

        // add @collection.* keys
        if (includeCollectionJoinKeys) {
            result = result.concat(cachedCollectionJoinKeys || []);
        }

        return result;
    }

    // Returns object with all the completions matching the context.
    function completions(context) {
        let word = context.matchBefore(/[\'\"\@\w\.]*/);
        if (word && word.from == word.to && !context.explicit) {
            return null;
        }

        // skip for comments
        let nodeBefore = syntaxTree(context.state).resolveInner(context.pos, -1);
        if (nodeBefore?.type?.name == "comment") {
            return null;
        }

        let options = [
            { label: "false" },
            { label: "true" },
            { label: "@now" },
            { label: "@second" },
            { label: "@minute" },
            { label: "@hour" },
            { label: "@year" },
            { label: "@day" },
            { label: "@month" },
            { label: "@weekday" },
            { label: "@yesterday" },
            { label: "@tomorrow" },
            { label: "@todayStart" },
            { label: "@todayEnd" },
            { label: "@monthStart" },
            { label: "@monthEnd" },
            { label: "@yearStart" },
            { label: "@yearEnd" },
        ];

        if (!disableCollectionJoinKeys) {
            options.push({ label: "@collection.*", apply: "@collection." });
        }

        let keys = getAllKeys(
            !disableRequestKeys && word.text.startsWith("@r"),
            !disableCollectionJoinKeys && word.text.startsWith("@c"),
        );

        for (const key of keys) {
            options.push({
                label: key.endsWith(".") ? key + "*" : key,
                apply: key,
                boost: key.indexOf("_via_") > 0 ? -1 : 0, // deprioritize _via_ keys
            });
        }

        return {
            from: word.from,
            options: options,
        };
    }

    // Creates a new language mode.
    // @see https://codemirror.net/5/demo/simplemode.html
    function ruleLang() {
        return StreamLanguage.define(
            simpleMode({
                start: [
                    // base literals
                    {
                        regex: /true|false|null/,
                        token: "atom",
                    },
                    // comments
                    { regex: /\/\/.*/, token: "comment" },
                    // double quoted string
                    { regex: /"(?:[^\\]|\\.)*?(?:"|$)/, token: "string" },
                    // single quoted string
                    { regex: /'(?:[^\\]|\\.)*?(?:'|$)/, token: "string" },
                    // numbers
                    {
                        regex: /0x[a-f\d]+|[-+]?(?:\.\d+|\d+\.?\d*)(?:e[-+]?\d+)?/i,
                        token: "number",
                    },
                    // operators
                    {
                        regex: /\&\&|\|\||\=|\!\=|\~|\!\~|\>|\<|\>\=|\<\=/,
                        token: "operator",
                    },
                    // indent and dedent properties guide autoindentation
                    { regex: /[\{\[\(]/, indent: true },
                    { regex: /[\}\]\)]/, dedent: true },
                    // keywords
                    { regex: /\w+[\w\.]*\w+/, token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@now"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@second"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@minute"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@hour"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@year"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@day"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@month"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@weekday"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@todayStart"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@todayEnd"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@monthStart"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@monthEnd"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@yearStart"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@yearEnd"), token: "keyword" },
                    { regex: CommonHelper.escapeRegExp("@request.method"), token: "keyword" },
                ],
                meta: {
                    lineComment: "//",
                },
            }),
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

        addLabelListeners();

        let keybindings = [
            submitShortcut,
            ...closeBracketsKeymap,
            ...defaultKeymap,
            searchKeymap.find((item) => item.key === "Mod-d"),
            ...historyKeymap,
            ...completionKeymap,
        ];
        if (!singleLine) {
            keybindings.push(indentWithTab);
        }

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
                    keymap.of(keybindings),
                    EditorView.lineWrapping,
                    autocompletion({
                        override: [completions],
                        icons: false,
                    }),
                    placeholderCompartment.of(placeholderExt(placeholder)),
                    editableCompartment.of(EditorView.editable.of(!disabled)),
                    readOnlyCompartment.of(EditorState.readOnly.of(disabled)),
                    langCompartment.of(ruleLang()),
                    EditorState.transactionFilter.of((tr) => {
                        if (singleLine && tr.newDoc.lines > 1) {
                            if (!tr.changes?.inserted?.filter((i) => !!i.text.find((t) => t))?.length) {
                                return []; // only empty lines
                            }
                            // it is ok to mutate the current transaction as we don't change the doc length
                            tr.newDoc.text = [tr.newDoc.text.join(" ")];
                        }
                        return tr;
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
            clearTimeout(refreshDebounceId);
            removeLabelListeners();
            editor?.destroy();
            autocompleteWorker.terminate();
        };
    });
</script>

<div bind:this={container} class="code-editor" />
