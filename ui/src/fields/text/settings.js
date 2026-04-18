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
                        t.label(
                            { htmlFor: uniqueId + ".min" },
                            t.span({ className: "txt" }, "Min length"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip("Clear the field or set it to 0 for no limit."),
                            }),
                        ),
                        t.input({
                            type: "number",
                            id: uniqueId + ".min",
                            name: () => `fields.${data.fieldIndex}.min`,
                            step: 1,
                            min: 0,
                            max: Number.MAX_SAFE_INTEGER,
                            placeholder: "No min limit",
                            value: () => data.field.min || "",
                            oninput: (e) => {
                                data.field.min = parseInt(e.target.value, 10);
                            },
                        }),
                    ),
                ),
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label(
                            { htmlFor: uniqueId + ".max" },
                            t.span({ className: "txt" }, "Max length"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    "Clear the field or set it to 0 to fallback to the default limit.",
                                ),
                            }),
                        ),
                        t.input({
                            type: "number",
                            id: uniqueId + ".max",
                            name: () => `fields.${data.fieldIndex}.max`,
                            step: 1,
                            min: () => data.field.min || 0,
                            max: Number.MAX_SAFE_INTEGER,
                            placeholder: "Default to max 5000 characters",
                            value: () => data.field.max || "",
                            oninput: (e) => {
                                data.field.max = parseInt(e.target.value, 10);
                            },
                        }),
                    ),
                ),
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label(
                            { htmlFor: uniqueId + ".pattern" },
                            t.span({ className: "txt" }, "Validation pattern"),
                            () => {
                                if (data.field.primaryKey) {
                                    return t.i({
                                        className: "ri-information-line link-hint",
                                        ariaDescription: app.attrs.tooltip(
                                            "All record ids have forbidden characters and unique case-insensitive (ASCII) validations in addition to the user defined regex pattern.",
                                        ),
                                    });
                                }
                            },
                        ),
                        t.input({
                            type: "text",
                            id: uniqueId + ".pattern",
                            name: () => `fields.${data.fieldIndex}.pattern`,
                            value: () => data.field.pattern || "",
                            oninput: (e) => (data.field.pattern = e.target.value),
                        }),
                    ),
                    t.div({ className: "field-help" }, "Ex. ", t.code(null, "^[a-z0-9]+$")),
                ),
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label(
                            { htmlFor: uniqueId + ".autogeneratePattern" },
                            t.span({ className: "txt" }, "Autogenerate pattern"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    "Set and autogenerate text matching the pattern on missing record create value.",
                                ),
                            }),
                        ),
                        t.input({
                            type: "text",
                            id: uniqueId + ".autogeneratePattern",
                            name: () => `fields.${data.fieldIndex}.autogeneratePattern`,
                            value: () => data.field.autogeneratePattern || "",
                            oninput: (e) => (data.field.autogeneratePattern = e.target.value),
                        }),
                    ),
                    t.div({ className: "field-help" }, "Ex. ", t.code(null, "[a-z0-9]{30}")),
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
