<script context="module">
    /*
     * ---------------------------------------------------------------
     * The below component is similar to https://github.com/tinymce/tinymce-svelte
     * but with removed unnecessary dependencies (eg. the TinyMCE cloud loading script)
     * and with extra error catching to handle the async edge-cases
     * when the init event is fired after the Svelte component was destroyed.
     * ---------------------------------------------------------------
     */

    function createScriptLoader() {
        let state = {
            listeners: [],
            scriptLoaded: false,
            injected: false,
        };

        function injectScript(doc, url, callback) {
            state.injected = true;
            const script = doc.createElement("script");
            script.referrerPolicy = "origin";
            script.type = "application/javascript";
            script.src = url;
            script.onload = () => {
                callback();
            };
            if (doc.head) {
                doc.head.appendChild(script);
            }
        }

        function load(doc, url, callback) {
            if (state.scriptLoaded) {
                callback();
            } else {
                state.listeners.push(callback);
                // check we can access doc
                if (!state.injected) {
                    injectScript(doc, url, () => {
                        state.listeners.forEach((fn) => fn());
                        state.scriptLoaded = true;
                    });
                }
            }
        }

        return { load };
    }

    let scriptLoader = createScriptLoader();
</script>

<script>
    import { onMount, createEventDispatcher } from "svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let id = "tinymce_svelte" + CommonHelper.randomString(7);
    export let inline = undefined;
    export let disabled = false;
    export let scriptSrc = `${import.meta.env.BASE_URL}libs/tinymce/tinymce.min.js`;
    export let conf = {};
    export let modelEvents = "change input undo redo";
    export let value = "";
    export let text = "";
    export let cssClass = "tinymce-wrapper";

    // Events
    // ---------------------------------------------------------------
    const validEvents = [
        "Activate",
        "AddUndo",
        "BeforeAddUndo",
        "BeforeExecCommand",
        "BeforeGetContent",
        "BeforeRenderUI",
        "BeforeSetContent",
        "BeforePaste",
        "Blur",
        "Change",
        "ClearUndos",
        "Click",
        "ContextMenu",
        "Copy",
        "Cut",
        "Dblclick",
        "Deactivate",
        "Dirty",
        "Drag",
        "DragDrop",
        "DragEnd",
        "DragGesture",
        "DragOver",
        "Drop",
        "ExecCommand",
        "Focus",
        "FocusIn",
        "FocusOut",
        "GetContent",
        "Hide",
        "Init",
        "KeyDown",
        "KeyPress",
        "KeyUp",
        "LoadContent",
        "MouseDown",
        "MouseEnter",
        "MouseLeave",
        "MouseMove",
        "MouseOut",
        "MouseOver",
        "MouseUp",
        "NodeChange",
        "ObjectResizeStart",
        "ObjectResized",
        "ObjectSelected",
        "Paste",
        "PostProcess",
        "PostRender",
        "PreProcess",
        "ProgressState",
        "Redo",
        "Remove",
        "Reset",
        "ResizeEditor",
        "SaveContent",
        "SelectionChange",
        "SetAttrib",
        "SetContent",
        "Show",
        "Submit",
        "Undo",
        "VisualAid",
    ];

    const bindHandlers = (editor, dispatch) => {
        validEvents.forEach((eventName) => {
            editor.on(eventName, (e) => {
                dispatch(eventName.toLowerCase(), {
                    eventName,
                    event: e,
                    editor,
                });
            });
        });
    };
    // ---------------------------------------------------------------

    let container;
    let element;
    let editorRef;

    let lastVal = value;
    let disablindCache = disabled;

    const dispatch = createEventDispatcher();

    $: {
        try {
            if (editorRef && lastVal !== value) {
                editorRef.setContent(value);
                text = editorRef.getContent({ format: "text" });
            }
            if (editorRef && disabled !== disablindCache) {
                disablindCache = disabled;
                if (typeof editorRef.mode?.set === "function") {
                    editorRef.mode.set(disabled ? "readonly" : "design");
                } else {
                    editorRef.setMode(disabled ? "readonly" : "design");
                }
            }
        } catch (err) {
            console.warn("TinyMCE reactive error:", err);
        }
    }

    function getTinymce() {
        return window && window.tinymce ? window.tinymce : null;
    }

    function init() {
        const finalInit = {
            ...conf,
            target: element,
            inline: inline !== undefined ? inline : conf.inline !== undefined ? conf.inline : false,
            readonly: disabled,
            setup: (editor) => {
                editorRef = editor;
                editor.on("init", () => {
                    editor.setContent(value);
                    // bind model events
                    editor.on(modelEvents, () => {
                        lastVal = editor.getContent();
                        if (lastVal !== value) {
                            value = lastVal;
                            text = editor.getContent({ format: "text" });
                        }
                    });
                });

                bindHandlers(editor, dispatch);

                if (typeof conf.setup === "function") {
                    conf.setup(editor);
                }
            },
        };

        element.style.visibility = "";

        getTinymce().init(finalInit);
    }

    onMount(() => {
        if (getTinymce() !== null) {
            init();
        } else {
            scriptLoader.load(container.ownerDocument, scriptSrc, () => {
                // init if the container is not removed from the DOM
                if (container) {
                    init();
                }
            });
        }

        return () => {
            try {
                if (editorRef) {
                    // temp workaround for https://github.com/tinymce/tinymce/issues/9377
                    editorRef.dom?.unbind(document);

                    getTinymce()?.remove(editorRef);
                }
            } catch (_) {}
        };
    });
</script>

<div bind:this={container} class={cssClass}>
    {#if inline}
        <div {id} bind:this={element} />
    {:else}
        <textarea {id} bind:this={element} style="visibility: hidden" />
    {/if}
</div>
