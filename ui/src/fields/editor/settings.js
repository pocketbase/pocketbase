// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(props) {
    const uniqueId = "f_" + app.utils.randomString();

    const local = store({
        showInfo: false,
    });

    return app.components.fieldSettings(props, {
        content: () =>
            t.div(
                { className: "grid sm" },
                t.div(
                    { className: "col-sm-12" },
                    t.div(
                        { className: "field" },
                        t.label(
                            { htmlFor: uniqueId + ".maxSize" },
                            t.span(null, "Max size "),
                            t.small(null, "(bytes)"),
                        ),
                        t.input({
                            type: "number",
                            id: uniqueId + ".maxSize",
                            name: () => `fields.${props.fieldIndex}.maxSize`,
                            min: 0,
                            step: 1,
                            max: Number.MAX_SAFE_INTEGER,
                            placeholder: "Default to max ~5MB",
                            value: () => props.field.maxSize || "",
                            oninput: (e) => {
                                props.field.maxSize = parseInt(e.target.value, 10);
                            },
                        }),
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
                            name: () => `fields.${props.fieldIndex}.help`,
                            value: () => props.field.help || "",
                            oninput: (e) => (props.field.help = e.target.value),
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
                    name: () => `fields.${props.fieldIndex}.required`,
                    checked: () => !!props.field.required,
                    onchange: (e) => (props.field.required = e.target.checked),
                }),
                t.label(
                    { htmlFor: uniqueId + ".required" },
                    t.span({ className: "txt" }, "Required"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Requires the field value to be nonempty string"),
                    }),
                ),
            ),
            t.div(
                { className: "field" },
                t.input({
                    className: "sm",
                    type: "checkbox",
                    id: uniqueId + ".convertURLs",
                    name: () => `fields.${props.fieldIndex}.convertURLs`,
                    checked: () => !!props.field.convertURLs,
                    onchange: (e) => (props.field.convertURLs = e.target.checked),
                }),
                t.label(
                    { htmlFor: uniqueId + ".convertURLs" },
                    t.span({ className: "txt" }, "Strip URLs domain"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip(
                            "This could help making the editor content more portable between environments since there will be no local base url to replace.",
                        ),
                    }),
                ),
            ),
        ],
    });
}
