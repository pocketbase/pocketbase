import cssVars from "@/css/vars.css?inline";

window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Creates a new TinyMCE editor element.
 *
 * @example
 * ```js
 * const data = store({ value: "" })
 *
 * app.components.tinymce({
 *     value: () => data.value,
 *     onchange: (val) => data.value = val,
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.tinymce = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        name: undefined,
        className: "",
        value: "",
        readonly: false,
        disabled: false,
        required: false,
        convertURLs: false,
        onchange: function(val) {},
        onbeforeinit: function(opts) {},
        onafterinit: function(editor) {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    let editorRef;
    let textarea;
    let oldChange;

    watchers.push(watch(() => props.value, setEditorContentValue));
    watchers.push(watch(() => props.disabled || props.readonly, setDisabled));
    watchers.push(watch(() => props.convertURLs, setConvertURLs));
    watchers.push(watch(() => app.store.activeColorScheme, setEditorBodyColorScheme));

    // generic error handling wrapper to prevent throws from tinymce API calls from crashing the UI
    function catchError(fn) {
        try {
            fn();
        } catch (err) {
            console.warn("tinymce error:", err);
        }
    }

    function setEditorContentValue() {
        if (oldChange != props.value) {
            catchError(() => {
                editorRef?.setContent("" + (props.value || "")); // stringify and normalize
            });
        }
    }

    function setDisabled() {
        catchError(() => {
            // https://www.tiny.cloud/docs/tinymce/6/editor-important-options/#readonly
            editorRef?.mode?.set(props.disabled || props.readonly ? "readonly" : "design");
        });
    }

    function setConvertURLs() {
        catchError(() => {
            editorRef?.options?.set("convert_urls", !!props.convertURLs);
        });
    }

    function setEditorBodyColorScheme() {
        catchError(() => {
            editorRef?.getBody()?.setAttribute("data-color-scheme", app.store.activeColorScheme);
        });
    }

    let changeTimeoutId;
    function triggerOnchangeWithDebounce(debounce = 150) {
        clearTimeout(changeTimeoutId);
        changeTimeoutId = setTimeout(triggerOnchange, debounce);
    }

    function triggerOnchange() {
        if (!editorRef) {
            return;
        }

        clearTimeout(changeTimeoutId);

        let content;
        catchError(() => {
            content = editorRef.getContent();
        });
        if (content == oldChange) {
            return; // no change
        }

        oldChange = content;
        props.onchange?.(content);

        // trigger custom change event for clearing field errors
        textarea?.dispatchEvent(
            new CustomEvent("change", {
                detail: { editor: editorRef, content: content },
                bubbles: true,
            }),
        );
    }

    function destroyEditor() {
        if (!editorRef) {
            return; // already removed or not initialized yet
        }

        clearTimeout(changeTimeoutId);

        // workaround for https://github.com/tinymce/tinymce/issues/9377
        editorRef.dom?.unbind(document);

        catchError(() => {
            window.tinymce?.remove(editorRef);
        });
        editorRef = null;
        oldChange = null;
    }

    async function initEditor(el) {
        await loadTinyMCE();

        destroyEditor();

        // removed while loading
        if (!el.isConnected) {
            return;
        }

        const opts = {
            target: el,
            content_style: cssVars,
            branding: false,
            promotion: false,
            menubar: false,
            resize: false,
            min_height: 265,
            height: 265,
            max_height: 600,
            sandbox_iframes: true,
            convert_unsafe_embeds: true, // GHSA-5359
            codesample_global_prismjs: true,
            convert_urls: false,
            relative_urls: false,
            autoresize_bottom_margin: 30,
            media_poster: false,
            media_alt_source: false,
            ui_mode: "split",
            codesample_languages: [
                { text: "HTML/XML", value: "markup" },
                { text: "CSS", value: "css" },
                { text: "SQL", value: "sql" },
                { text: "JavaScript", value: "javascript" },
                { text: "Go", value: "go" },
                { text: "Dart", value: "dart" },
                { text: "Zig", value: "zig" },
                { text: "Rust", value: "rust" },
                { text: "Lua", value: "lua" },
                { text: "PHP", value: "php" },
                { text: "Ruby", value: "ruby" },
                { text: "Python", value: "python" },
                { text: "Java", value: "java" },
                { text: "C", value: "c" },
                { text: "C#", value: "csharp" },
                { text: "C++", value: "cpp" },
                // other non-highlighted languages
                { text: "Markdown", value: "markdown" },
                { text: "Swift", value: "swift" },
                { text: "Kotlin", value: "kotlin" },
                { text: "Elixir", value: "elixir" },
                { text: "Scala", value: "scala" },
                { text: "Julia", value: "julia" },
                { text: "Haskell", value: "haskell" },
            ],
            plugins: [
                "autolink",
                "autoresize",
                "code",
                "codesample",
                "directionality",
                "image",
                "link",
                "lists",
                "media",
                "table",
                "wordcount",
            ],
            toolbar:
                "styles | alignleft aligncenter alignright | bold italic forecolor backcolor | bullist numlist | link table media_picker codesample | direction code",
            paste_postprocess: (editor, args) => {
                cleanupPastedNode(args.node);
            },
            // @see https://www.tiny.cloud/docs/tinymce/6/file-image-upload/#interactive-example
            file_picker_types: "image",
            file_picker_callback: (callback, value, meta) => {
                const input = document.createElement("input");
                input.setAttribute("type", "file");
                input.setAttribute("accept", "image/*");

                input.addEventListener("change", (e) => {
                    const file = e.target.files[0];
                    const reader = new FileReader();

                    reader.addEventListener("load", () => {
                        if (!tinymce) {
                            return;
                        }

                        // We need to register the blob in TinyMCEs image blob registry.
                        // In future TinyMCE version this part will be handled internally.
                        const id = "blobid" + new Date().getTime();
                        const blobCache = tinymce.activeEditor.editorUpload.blobCache;
                        const base64 = reader.result.split(",")[1];
                        const blobInfo = blobCache.create(id, file, base64);
                        blobCache.add(blobInfo);

                        // call the callback and populate the Title field with the file name
                        callback(blobInfo.blobUri(), { title: file.name });
                    });

                    reader.readAsDataURL(file);
                });

                input.click();
            },
            setup: (editor) => {
                editorRef = editor;

                editor.on("init", (e) => {
                    props.onafterinit?.(editorRef);

                    setConvertURLs();
                    setDisabled();
                    setEditorBodyColorScheme();
                    setEditorContentValue();
                });

                // propagate save shortcut to the parent
                editor.on("keydown", (e) => {
                    if ((e.ctrlKey || e.metaKey) && e.code == "KeyS" && editor.formElement) {
                        e.preventDefault();
                        e.stopPropagation();
                        editor.formElement.dispatchEvent(new KeyboardEvent("keydown", e));
                    }
                });

                editor.on("input", (e) => {
                    triggerOnchangeWithDebounce();
                });

                editor.on("change", (e) => {
                    triggerOnchange();
                });

                registerDirectionButton(editor);
                registerMediaButton(editor);
            },
        };

        if (props.readonly) {
            opts.statusbar = false;
            opts.min_height = 30;
            opts.height = 30;
            opts.max_height = 500;
            opts.autoresize_bottom_margin = 5;
            opts.resize = false;
            opts.toolbar = false;
            opts.plugins = ["autoresize", "codesample", "directionality"];
        }

        if (props.onbeforeinit) {
            props.onbeforeinit(opts);
        }

        window.tinymce.init(opts);
    }

    textarea = t.textarea({
        name: () => props.name,
        onmount: (el) => {
            initEditor(el).catch((err) => {
                console.warn("tinymce init error:", err);
            });
        },
        onunmount: destroyEditor,
    });

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `pb-tinymce ${props.className}`,
            "html-required": () => props.required || undefined, // set on the parent because the textarea will be hidden
            onunmount: (el) => {
                clearTimeout(changeTimeoutId);
                watchers.forEach((w) => w?.unwatch());
                textarea = null;
            },
        },
        textarea,
    );
};

function registerDirectionButton(editor) {
    const lastDirectionKey = "pbTinymceLastDirection";

    // load last used text direction for blank editors
    editor.on("init", () => {
        const lastDirection = window.localStorage.getItem(lastDirectionKey);
        if (!editor.isDirty() && editor.getContent() == "" && lastDirection == "rtl") {
            editor.execCommand("mceDirectionRTL");
        }
    });

    // text direction dropdown
    editor.ui.registry.addMenuButton("direction", {
        icon: "visualchars",
        tooltip: "Direction",
        fetch: (callback) => {
            const items = [
                {
                    type: "menuitem",
                    text: "LTR content",
                    icon: "ltr",
                    onAction: () => {
                        window?.localStorage?.setItem(lastDirectionKey, "ltr");
                        editor.execCommand("mceDirectionLTR");
                    },
                },
                {
                    type: "menuitem",
                    text: "RTL content",
                    icon: "rtl",
                    onAction: () => {
                        window?.localStorage?.setItem(lastDirectionKey, "rtl");
                        editor.execCommand("mceDirectionRTL");
                    },
                },
            ];

            callback(items);
        },
    });
}

function registerMediaButton(editor) {
    editor.ui.registry.addMenuButton("media_picker", {
        tooltip: "Insert media",
        icon: "embed",
        fetch: (callback) => {
            const items = [
                {
                    type: "menuitem",
                    text: "Inline image (Base64)",
                    onAction: () => {
                        editor.execCommand("mceImage");
                    },
                },
                {
                    type: "menuitem",
                    text: "Media from collection",
                    onAction: () => {
                        app.modals.openRecordFilePicker({
                            fileTypes: ["image", "audio", "video"],
                            onselect: (selected) => {
                                const url = app.pb.files.getURL(selected.record, selected.name, {
                                    thumb: selected.thumb || undefined,
                                });

                                // just an extra precaution in case the editor fail for whatever reason to sanitize the inserted raw htmls
                                const escapedName = app.utils.encodeEntities(selected.name);
                                const escapedUrl = app.utils.encodeEntities(url);

                                if (app.utils.hasImageExtension(selected.name)) {
                                    editor?.execCommand("InsertImage", false, url);
                                } else if (app.utils.hasAudioExtension(selected.name)) {
                                    editor?.execCommand(
                                        "InsertHTML",
                                        false,
                                        `<audio controls src="${escapedUrl}"></audio>`,
                                    );
                                } else if (app.utils.hasVideoExtension(escapedName)) {
                                    editor?.execCommand(
                                        "InsertHTML",
                                        false,
                                        `
                                        <video controls width="300">
                                            <source src="${escapedUrl}" />
                                            <p>Download: <a href="${escapedUrl}" download="${escapedName}">${escapedName}</a>.</p>
                                        </video>
                                    `,
                                    );
                                }
                            },
                        });
                    },
                },
                {
                    type: "menuitem",
                    text: "Manual embed",
                    onAction: () => {
                        tinymce.activeEditor.execCommand("mceMedia");
                    },
                },
            ];

            callback(items);
        },
    });
}

const allowedPasteNodes = [
    "DIV",
    "P",
    "A",
    "EM",
    "B",
    "STRONG",
    "H1",
    "H2",
    "H3",
    "H4",
    "H5",
    "H6",
    "TABLE",
    "TR",
    "TD",
    "TH",
    "TBODY",
    "THEAD",
    "TFOOT",
    "BR",
    "HR",
    "Q",
    "SUP",
    "SUB",
    "DEL",
    "IMG",
    "OL",
    "UL",
    "LI",
    "CODE",
];

function cleanupPastedNode(node) {
    if (!node) {
        return; // nothing to cleanup
    }

    for (const child of node.children) {
        cleanupPastedNode(child);
    }

    if (!allowedPasteNodes.includes(node.tagName)) {
        unwrap(node);
    } else {
        node.removeAttribute("style");
        node.removeAttribute("class");
    }
}

function unwrap(node) {
    let parent = node.parentNode;

    // move children outside of the parent node
    while (node.firstChild) {
        parent.insertBefore(node.firstChild, node);
    }

    // remove the now empty parent element
    parent.removeChild(node);
}

async function loadTinyMCE() {
    // already loaded
    if (typeof window.tinymce != "undefined") {
        return;
    }

    const scriptId = "lazy-tinymce-js";

    // in the process of being loaded
    if (document.getElementById(scriptId)) {
        return new Promise((resolve, reject) => {
            function cleanup() {
                document.removeEventListener("tinymceLoadSuccess", successHandler);
                document.removeEventListener("tinymceLoadError", errorHandler);
            }

            const successHandler = function() {
                cleanup();
                resolve();
            };

            const errorHandler = function(e) {
                cleanup();
                reject(e?.details);
            };

            document.addEventListener("tinymceLoadSuccess", successHandler);
            document.addEventListener("tinymceLoadError", errorHandler);
        });
    }

    return new Promise((resolve, reject) => {
        document.head.querySelector("#shablon-script").after(
            t.script({
                id: scriptId,
                src: import.meta.env.BASE_URL + "libs/tinymce/tinymce.min.js",
                onload: () => {
                    resolve();
                },
                onerror: (err) => {
                    console.warn("failed to load tinymce.min.js:", err);
                    reject(err);
                },
            }),
        );
    }).then(() => {
        document.dispatchEvent(new CustomEvent("tinymceLoadSuccess"));
    }).catch((err) => {
        document.dispatchEvent(new CustomEvent("tinymceLoadError", { detail: err }));
    });
}
