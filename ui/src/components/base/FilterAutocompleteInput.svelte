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
    } from "@codemirror/language";
    import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
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
    export let disableIndirectCollectionsKeys = false;

    let editor;
    let container;
    let oldDisabledState = disabled;
    let langCompartment = new Compartment();
    let editableCompartment = new Compartment();
    let readOnlyCompartment = new Compartment();
    let placeholderCompartment = new Compartment();

    let cachedCollections = [];
    let cachedRequestKeys = [];
    let cachedIndirectCollectionKeys = [];
    let cachedBaseKeys = [];
    let baseKeysChangeHash = "";
    let oldBaseKeysChangeHash = "";

    $: baseKeysChangeHash = getCollectionKeysChangeHash(baseCollection);

    $: if (
        !disabled &&
        (oldBaseKeysChangeHash != baseKeysChangeHash ||
            disableRequestKeys !== -1 ||
            disableIndirectCollectionsKeys !== -1)
    ) {
        oldBaseKeysChangeHash = baseKeysChangeHash;
        refreshCachedKeys();
    }

    $: if (id) {
        addLabelListeners();
    }

    $: if (editor && baseCollection?.schema) {
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

    let refreshDebounceId = null;

    // Refresh the cached autocomplete keys.
    function refreshCachedKeys() {
        clearTimeout(refreshDebounceId);
        refreshDebounceId = setTimeout(() => {
            cachedCollections = concatWithBaseCollection($collections);
            cachedBaseKeys = getBaseKeys();
            cachedRequestKeys = !disableRequestKeys ? getRequestKeys() : [];
            cachedIndirectCollectionKeys = !disableIndirectCollectionsKeys ? getIndirectCollectionKeys() : [];
        }, 300);
    }

    // Return a collection keys hash string that can be used to compare with previous states.
    function getCollectionKeysChangeHash(collection) {
        return JSON.stringify([collection?.name, collection?.type, collection?.schema]);
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

    // Returns a list with all collection field keys recursively.
    function getCollectionFieldKeys(nameOrId, prefix = "", level = 0) {
        let collection = cachedCollections.find((item) => item.name == nameOrId || item.id == nameOrId);
        if (!collection || level >= 4) {
            return [];
        }

        let result = CommonHelper.getAllCollectionIdentifiers(collection, prefix);

        for (const field of collection.schema) {
            const key = prefix + field.name;

            // add relation fields
            if (field.type === "relation" && field.options?.collectionId) {
                const subKeys = getCollectionFieldKeys(field.options.collectionId, key + ".", level + 1);
                if (subKeys.length) {
                    result = result.concat(subKeys);
                }
            }

            // add ":each" field modifier
            if (field.type === "select" && field.options?.maxSelect != 1) {
                result.push(key + ":each");
            }

            // add ":length" field modifier to arrayble fields
            if (field.options?.maxSelect != 1 && ["select", "file", "relation"].includes(field.type)) {
                result.push(key + ":length");
            }
        }

        return result;
    }

    // Returns baseCollection keys.
    function getBaseKeys() {
        return getCollectionFieldKeys(baseCollection?.name);
    }

    // Returns @request.* keys.
    function getRequestKeys() {
        const result = [];

        result.push("@request.method");
        result.push("@request.query.");
        result.push("@request.data.");
        result.push("@request.headers.");
        result.push("@request.auth.id");
        result.push("@request.auth.collectionId");
        result.push("@request.auth.collectionName");
        result.push("@request.auth.verified");
        result.push("@request.auth.username");
        result.push("@request.auth.email");
        result.push("@request.auth.emailVisibility");
        result.push("@request.auth.created");
        result.push("@request.auth.updated");

        // load auth collection fields
        const authCollections = cachedCollections.filter((collection) => collection.$isAuth);
        for (const collection of authCollections) {
            const authKeys = getCollectionFieldKeys(collection.id, "@request.auth.");
            for (const k of authKeys) {
                CommonHelper.pushUnique(result, k);
            }
        }

        // load base collection fields into @request.data.*
        const issetExcludeList = ["created", "updated"];
        if (baseCollection?.id) {
            const keys = getCollectionFieldKeys(baseCollection.name, "@request.data.");
            for (const key of keys) {
                result.push(key);

                // add ":isset" modifier to non-base keys
                const parts = key.split(".");
                if (
                    parts.length === 3 &&
                    // doesn't contain another modifier
                    parts[2].indexOf(":") === -1 &&
                    // is not from the exclude list
                    !issetExcludeList.includes(parts[2])
                ) {
                    result.push(key + ":isset");
                }
            }
        }

        return result;
    }

    // Returns @collection.* keys.
    function getIndirectCollectionKeys() {
        const result = [];

        for (const collection of cachedCollections) {
            const prefix = "@collection." + collection.name + ".";
            const keys = getCollectionFieldKeys(collection.name, prefix);
            for (const key of keys) {
                result.push(key);
            }
        }

        return result;
    }

    // Returns an array with all the supported keys.
    function getAllKeys(includeRequestKeys = true, includeIndirectCollectionsKeys = true) {
        let result = [].concat(extraAutocompleteKeys);

        // add base keys
        result = result.concat(cachedBaseKeys || []);

        // add @request.* keys
        if (includeRequestKeys) {
            result = result.concat(cachedRequestKeys || []);
        }

        // add @collections.* keys
        if (includeIndirectCollectionsKeys) {
            result = result.concat(cachedIndirectCollectionKeys || []);
        }

        // sort longer keys first because the highlighter will highlight
        // the first match and stops until an operator is found
        result.sort(function (a, b) {
            return b.length - a.length;
        });

        return result;
    }

    // Returns object with all the completions matching the context.
    function completions(context) {
        let word = context.matchBefore(/[\'\"\@\w\.]*/);
        if (word && word.from == word.to && !context.explicit) {
            return null;
        }

        let options = [{ label: "false" }, { label: "true" }, { label: "@now" }];

        if (!disableIndirectCollectionsKeys) {
            options.push({ label: "@collection.*", apply: "@collection." });
        }

        const keys = getAllKeys(!disableRequestKeys, !disableRequestKeys && word.text.startsWith("@c"));
        for (const key of keys) {
            options.push({
                label: key.endsWith(".") ? key + "*" : key,
                apply: key,
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
                    { regex: CommonHelper.escapeRegExp("@request.method"), token: "keyword" },
                ],
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
                        override: [completions],
                        icons: false,
                    }),
                    placeholderCompartment.of(placeholderExt(placeholder)),
                    editableCompartment.of(EditorView.editable.of(!disabled)),
                    readOnlyCompartment.of(EditorState.readOnly.of(disabled)),
                    langCompartment.of(ruleLang()),
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
            clearTimeout(refreshDebounceId);
            removeLabelListeners();
            editor?.destroy();
        };
    });
</script>

<div bind:this={container} class="code-editor" />
