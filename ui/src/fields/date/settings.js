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
                        t.label({ htmlFor: uniqueId + ".min" }, t.span({ className: "txt" }, "Min date (Local)")),
                        t.input({
                            type: "datetime-local",
                            id: uniqueId + ".min",
                            step: 1,
                            name: () => `fields.${data.fieldIndex}.min`,
                            value: () => app.utils.toDatetimeLocalInputValue(data.field.min),
                            onchange: (e) => {
                                data.field.min = app.utils.toRFC3339Datetime(e.target.value);
                            },
                        }),
                    ),
                ),
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".max" }, t.span({ className: "txt" }, "Max date (Local)")),
                        t.input({
                            type: "datetime-local",
                            id: uniqueId + ".max",
                            step: 1,
                            name: () => `fields.${data.fieldIndex}.max`,
                            value: () => app.utils.toDatetimeLocalInputValue(data.field.max),
                            onchange: (e) => {
                                data.field.max = app.utils.toRFC3339Datetime(e.target.value);
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
