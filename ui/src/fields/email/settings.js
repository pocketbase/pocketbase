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
                            { htmlFor: uniqueId + ".exceptDomains" },
                            t.span({ className: "txt" }, "Except domains"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    `List of domains that are NOT allowed.\nThis field is disabled if "Only domains" is set.`,
                                ),
                            }),
                        ),
                        t.input({
                            type: "text",
                            id: uniqueId + ".exceptDomains",
                            disabled: () => !app.utils.isEmpty(data.field.onlyDomains),
                            name: () => `fields.${data.fieldIndex}.exceptDomains`,
                            value: () => app.utils.joinNonEmpty(data.field.exceptDomains),
                            onchange: (
                                e,
                            ) => (data.field.exceptDomains = app.utils.splitNonEmpty(e.target.value, ",")),
                        }),
                    ),
                    t.div({ className: "field-help" }, "Use comma as separator."),
                ),
                t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field" },
                        t.label(
                            { htmlFor: uniqueId + ".onlyDomains" },
                            t.span({ className: "txt" }, "Only domains"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    `List of domains that are ONLY allowed.\nThis field is disabled if "Except domains" is set.`,
                                ),
                            }),
                        ),
                        t.input({
                            type: "text",
                            id: uniqueId + ".onlyDomains",
                            disabled: () => !app.utils.isEmpty(data.field.exceptDomains),
                            name: () => `fields.${data.fieldIndex}.onlyDomains`,
                            value: () => app.utils.joinNonEmpty(data.field.onlyDomains),
                            onchange: (e) => (data.field.onlyDomains = app.utils.splitNonEmpty(e.target.value, ",")),
                        }),
                    ),
                    t.div({ className: "field-help" }, "Use comma as separator."),
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
