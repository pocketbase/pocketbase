// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(data) {
    const uniqueId = "f_" + app.utils.randomString();

    return app.components.fieldSettings(data, {
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
                            name: () => `fields.${data.fieldIndex}.min`,
                            value: () => typeof data.field.min == "number" ? data.field.min : "",
                            oninput: (e) => {
                                if (!e.target.value) {
                                    data.field.min = null;
                                } else {
                                    data.field.min = Number(e.target.value);
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
                            min: () => data.field.min,
                            name: () => `fields.${data.fieldIndex}.max`,
                            value: () => typeof data.field.max == "number" ? data.field.max : "",
                            oninput: (e) => {
                                if (!e.target.value) {
                                    data.field.max = null;
                                } else {
                                    data.field.max = Number(e.target.value);
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
                    id: uniqueId + ".onlyInt",
                    name: () => `fields.${data.fieldIndex}.onlyInt`,
                    checked: () => !!data.field.onlyInt,
                    onchange: (e) => (data.field.onlyInt = e.target.checked),
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
                    t.small({ className: "txt-hint" }, "(!=0)"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Requires the field value to be not 0."),
                    }),
                ),
            ),
        ],
    });
}
