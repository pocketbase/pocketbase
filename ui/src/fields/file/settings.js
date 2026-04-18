// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(data) {
    const uniqueId = "f_" + app.utils.randomString();

    const isMultipleOptions = [
        { label: "Single", value: false },
        { label: "Multiple", value: true },
    ];

    return app.components.fieldSettings(data, {
        header: [
            t.div(
                {
                    className: "field header-select single-multiple-select",
                },
                app.components.select({
                    required: true,
                    options: isMultipleOptions,
                    value: () => {
                        return data.field.maxSelect > 1;
                    },
                    onchange: (opts) => {
                        if (opts?.[0]?.value) {
                            if (!data.field.maxSelect || data.field.maxSelect < 2) {
                                data.field.maxSelect = 10;
                            }
                        } else {
                            data.field.maxSelect = 1;
                        }
                    },
                }),
            ),
        ],
        content: () =>
            t.div(
                { className: "grid sm" },
                t.div(
                    { className: "col-sm-12" },
                    t.div(
                        { className: "field" },
                        t.label(
                            { htmlFor: uniqueId + ".mimeTypes" },
                            t.span({ className: "txt" }, "Allowed mime types"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    "Allow files ONLY with the listed mime types.\n Leave empty for no restriction.",
                                ),
                            }),
                        ),
                        app.components.select({
                            max: 99,
                            placeholder: "No restriction",
                            options: app.utils.mimeTypes.map((opt) => {
                                return {
                                    value: opt.mimeType,
                                    label: () =>
                                        t.div(
                                            { className: "inline-flex gap-10" },
                                            t.span({ className: "txt" }, opt.ext || "-"),
                                            t.small({ className: "txt-hint" }, opt.mimeType),
                                        ),
                                };
                            }),
                            name: () => `fields.${data.fieldIndex}.mimeTypes`,
                            value: () => app.utils.toArray(data.field.mimeTypes),
                            onchange: (opts) => (data.field.mimeTypes = opts.map((opt) => opt.value)),
                        }),
                    ),
                    t.div(
                        { className: "field-help" },
                        t.button(
                            {
                                "type": "button",
                                "className": "link-hint gap-0",
                                "html-popovertarget": uniqueId + "mimeTypesDropdown",
                            },
                            t.span({ className: "txt" }, "Choose presets"),
                            t.i({ className: "ri-arrow-drop-down-fill", ariaHidden: true }),
                        ),
                        t.div(
                            {
                                id: uniqueId + "mimeTypesDropdown",
                                className: "dropdown sm nowrap left p-10",
                                popover: "auto",
                            },
                            t.button({
                                type: "button",
                                className: "dropdown-item",
                                role: "menuitem",
                                onclick: (e) => {
                                    data.field.mimeTypes = [
                                        "image/jpeg",
                                        "image/png",
                                        "image/svg+xml",
                                        "image/gif",
                                        "image/webp",
                                    ];

                                    e.target.closest(".dropdown").hidePopover();
                                },
                                textContent: "Images (jpg, png, svg, gif, webp)",
                            }),
                            t.button({
                                type: "button",
                                className: "dropdown-item",
                                role: "menuitem",
                                onclick: (e) => {
                                    data.field.mimeTypes = [
                                        "application/pdf",
                                        "application/msword",
                                        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
                                        "application/vnd.ms-excel",
                                        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
                                    ];

                                    e.target.closest(".dropdown").hidePopover();
                                },
                                textContent: "Documents (pdf, doc/docx, xls/xlsx)",
                            }),
                            t.button({
                                type: "button",
                                className: "dropdown-item",
                                role: "menuitem",
                                onclick: (e) => {
                                    data.field.mimeTypes = [
                                        "video/mp4",
                                        "video/mpeg",
                                        "video/x-msvideo",
                                        "video/quicktime",
                                        "video/3gpp",
                                    ];

                                    e.target.closest(".dropdown").hidePopover();
                                },
                                textContent: "Videos (mp4, mpeg, avi, mov, 3gp)",
                            }),
                            t.button({
                                type: "button",
                                className: "dropdown-item",
                                role: "menuitem",
                                onclick: (e) => {
                                    data.field.mimeTypes = [
                                        "application/zip",
                                        "application/x-7z-compressed",
                                        "application/x-rar-compressed",
                                    ];

                                    e.target.closest(".dropdown").hidePopover();
                                },
                                textContent: "Archives (zip, 7zip, rar)",
                            }),
                        ),
                    ),
                ),
                t.div(
                    { className: () => (data.field.maxSelect > 1 ? "col-sm-6" : "col-sm-9") },
                    t.div(
                        { className: "field" },
                        t.label(
                            {
                                htmlFor: uniqueId + ".thumbs",
                            },
                            t.span({ className: "txt" }, "Thumb sizes"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    "List of additional thumb sizes for image files, along with the default thumb size of 100x100. The thumbs are generated lazily on first access.",
                                ),
                            }),
                        ),
                        t.input({
                            type: "text",
                            id: uniqueId + ".thumbs",
                            placeholder: "e.g. 50x50, 480x720",
                            name: () => `fields.${data.fieldIndex}.thumbs`,
                            value: () => app.utils.joinNonEmpty(data.field.thumbs),
                            onchange: (e) => (data.field.thumbs = app.utils.splitNonEmpty(e.target.value, ",")),
                        }),
                    ),
                    t.div(
                        { className: "field-help" },
                        t.span({ className: "txt m-r-5" }, "Use comma as separator."),
                        t.button(
                            {
                                "type": "button",
                                "className": "link-hint gap-0",
                                "html-popovertarget": uniqueId + "thumbFormatsDropdown",
                            },
                            t.span({ className: "txt" }, "Supported formats"),
                            t.i({ className: "ri-arrow-drop-down-fill", ariaHidden: true }),
                        ),
                        t.div(
                            {
                                id: uniqueId + "thumbFormatsDropdown",
                                className: "dropdown sm nowrap left p-10",
                                popover: "auto",
                            },
                            t.ul(
                                { className: "m-0 p-l-sm" },
                                t.li(
                                    null,
                                    t.strong(null, "WxH"),
                                    t.span(null, " (e.g. 100x50) - crop to WxH viewbox (from center)"),
                                ),
                                t.li(
                                    null,
                                    t.strong(null, "WxHt"),
                                    t.span(null, " (e.g. 100x50t) - crop to WxH viewbox (from top)"),
                                ),
                                t.li(
                                    null,
                                    t.strong(null, "WxHb"),
                                    t.span(null, " (e.g. 100x50b) - crop to WxH viewbox (from bottom)"),
                                ),
                                t.li(
                                    null,
                                    t.strong(null, "WxHf"),
                                    t.span(null, " (e.g. 100x50f) - fit inside a WxH viewbox (without cropping)"),
                                ),
                                t.li(
                                    null,
                                    t.strong(null, "0xH"),
                                    t.span(null, " (e.g. 0x50) - resize to H height preserving the aspect ratio"),
                                ),
                                t.li(
                                    null,
                                    t.strong(null, "Wx0"),
                                    t.span(null, " (e.g. 100x0) - resize to W width preserving the aspect ratio"),
                                ),
                            ),
                        ),
                    ),
                ),
                t.div(
                    { className: "col-sm-3" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".maxSize" }, "Max size"),
                        t.input({
                            type: "number",
                            id: uniqueId + ".maxSize",
                            step: 1,
                            min: 0,
                            max: Number.MAX_SAFE_INTEGER,
                            placeholder: "~5MB default",
                            name: () => `fields.${data.fieldIndex}.maxSize`,
                            value: () => data.field.maxSize || "",
                            oninput: (e) => (data.field.maxSize = parseInt(e.target.value, 10)),
                        }),
                    ),
                    t.div({ className: "field-help" }, "In bytes."),
                ),
                () => {
                    if (data.field.maxSelect > 1) {
                        return t.div(
                            { className: "col-sm-3" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: uniqueId + ".maxSelect" }, "Max select"),
                                t.input({
                                    type: "number",
                                    id: uniqueId + ".maxSelect",
                                    placeholder: "Default to single",
                                    step: 1,
                                    min: 2,
                                    required: true,
                                    max: Number.MAX_SAFE_INTEGER,
                                    name: () => `fields.${data.fieldIndex}.maxSelect`,
                                    value: () => data.field.maxSelect || "",
                                    onchange: (e) => {
                                        const maxSelect = parseInt(e.target.value, 10);
                                        if (maxSelect > 1) {
                                            props.field.maxSelect = maxSelect;
                                        } else {
                                            props.field.maxSelect = 1;
                                        }
                                    },
                                }),
                            ),
                        );
                    }
                },
                t.div(
                    { className: "col-sm-12" },
                    t.div(
                        { className: "field m-t-5 m-b-5" },
                        t.input({
                            className: "switch",
                            type: "checkbox",
                            id: uniqueId + ".protected",
                            name: () => `fields.${data.fieldIndex}.protected`,
                            checked: () => !!data.field.protected,
                            onchange: (e) => (data.field.protected = e.target.checked),
                        }),
                        t.label(
                            { htmlFor: uniqueId + ".protected" },
                            t.span({ className: "txt" }, "Protected"),
                            t.small(
                                { className: "txt-hint" },
                                "Files will require View API rule permissions and file token (",
                                t.a({
                                    href: import.meta.env.PB_PROTECTED_FILE_DOCS,
                                    target: "_blank",
                                    rel: "noopener noreferrer",
                                    textContent: "Learn more",
                                }),
                                ").",
                            ),
                        ),
                    ),
                ),
                t.div(
                    { className: "col-sm-12" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".help" }, "Help text"),
                        t.input({
                            type: "text",
                            id: uniqueId + ".help",
                            name: () => `fields.${data.fieldIndex}.help`,
                            value: () => data.field.help || "",
                            oninput: (e) => (data.field.help = e.target.value),
                        }),
                    ),
                ),
            ),
        footer: () => [
            t.div(
                { className: "field" },
                t.input({
                    className: "sm",
                    type: "checkbox",
                    id: uniqueId + ".required",
                    name: () => `fields.${data.fieldIndex}.required`,
                    checked: () => !!data.field.required,
                    onchange: (e) => (data.field.required = e.target.checked),
                }),
                t.label(
                    { htmlFor: uniqueId + ".required" },
                    t.span({ className: "txt" }, "Required"),
                    t.small({ className: "txt-hint" }, "(!='')"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Requires the field value to be nonempty string"),
                    }),
                ),
            ),
        ],
    });
}
