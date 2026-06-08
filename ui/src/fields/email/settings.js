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
                            disabled: () => !app.utils.isEmpty(props.field.onlyDomains),
                            name: () => `fields.${props.fieldIndex}.exceptDomains`,
                            value: () => app.utils.joinNonEmpty(props.field.exceptDomains),
                            onchange: (
                                e,
                            ) => (props.field.exceptDomains = app.utils.splitNonEmpty(e.target.value, ",")),
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
                            disabled: () => !app.utils.isEmpty(props.field.exceptDomains),
                            name: () => `fields.${props.fieldIndex}.onlyDomains`,
                            value: () => app.utils.joinNonEmpty(props.field.onlyDomains),
                            onchange: (e) => (props.field.onlyDomains = app.utils.splitNonEmpty(e.target.value, ",")),
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
                    t.small({ className: "txt-hint" }, "(!='')"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Requires the field value to be nonempty string."),
                    }),
                ),
            ),
        ],
    });
}
