// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(props) {
    const uniqueId = "f_" + app.utils.randomString();

    return app.components.fieldSettings(props, {
        content: () =>
            t.div(
                { className: "grid sm" },
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".min" }, "Min"),
                        t.input({
                            type: "text",
                            id: uniqueId + ".min",
                            name: () => `fields.${props.fieldIndex}.min`,
                            value: () => typeof props.field.min == "number" ? props.field.min : "",
                            oninput: (e) => {
                                if (!e.target.value) {
                                    props.field.min = null;
                                } else {
                                    props.field.min = Number(e.target.value);
                                }
                            },
                        }),
                    ),
                ),
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".max" }, "Max"),
                        t.input({
                            type: "text",
                            id: uniqueId + ".max",
                            min: () => props.field.min,
                            name: () => `fields.${props.fieldIndex}.max`,
                            value: () => typeof props.field.max == "number" ? props.field.max : "",
                            oninput: (e) => {
                                if (!e.target.value) {
                                    props.field.max = null;
                                } else {
                                    props.field.max = Number(e.target.value);
                                }
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
                    t.small({ className: "txt-hint" }, "(!=0)"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Requires the field value to be not 0."),
                    }),
                ),
            ),
            t.div(
                { className: "field" },
                t.input({
                    className: "sm",
                    type: "checkbox",
                    id: uniqueId + ".onlyInt",
                    name: () => `fields.${props.fieldIndex}.onlyInt`,
                    checked: () => !!props.field.onlyInt,
                    onchange: (e) => (props.field.onlyInt = e.target.checked),
                }),
                t.label(
                    { htmlFor: uniqueId + ".onlyInt" },
                    t.span({ className: "txt" }, "No decimals"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Existing decimal numbers will not be affected."),
                    }),
                ),
            ),
        ],
    });
}
