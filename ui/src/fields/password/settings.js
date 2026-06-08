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
        showHidden: false,
        showPresentable: false,
        showDuplicate: false,
        content: t.div(
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
                        name: () => `fields.${props.fieldIndex}.min`,
                        step: 1,
                        min: 0,
                        max: 71,
                        placeholder: "No min limit",
                        value: () => props.field.min || "",
                        oninput: (e) => {
                            props.field.min = parseInt(e.target.value, 10);
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
                                "Clear the field or set it to 0 to fallback to the default limit (71).",
                            ),
                        }),
                    ),
                    t.input({
                        type: "number",
                        id: uniqueId + ".max",
                        name: () => `fields.${props.fieldIndex}.max`,
                        step: 1,
                        min: () => props.field.min || 0,
                        max: 71,
                        placeholder: "Up to 71 chars",
                        value: () => props.field.max || "",
                        oninput: (e) => {
                            props.field.max = parseInt(e.target.value, 10);
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
                        t.span({ className: "txt" }, "Bcrypt cost"),
                        t.i({
                            className: "ri-information-line link-hint",
                            ariaDescription: app.attrs.tooltip(
                                "Clear the field or set it to 0 to fallback to the default (10).",
                            ),
                        }),
                    ),
                    t.input({
                        type: "number",
                        id: uniqueId + ".cost",
                        name: () => `fields.${props.fieldIndex}.cost`,
                        step: 1,
                        // https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-constants
                        min: 4,
                        max: 31,
                        placeholder: "Default to 10",
                        value: () => props.field.cost || "",
                        oninput: (e) => {
                            props.field.cost = parseInt(e.target.value, 10);
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
                    ),
                    t.input({
                        type: "text",
                        id: uniqueId + ".pattern",
                        placeholder: "ex. ^\\w+$",
                        name: () => `fields.${props.fieldIndex}.pattern`,
                        value: () => props.field.pattern || "",
                        oninput: (e) => (props.field.pattern = e.target.value),
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
        footer: () => {
            // the system password auth field is always required
            if (props.collection?.type == "auth" && props.field.name == "password") {
                return;
            }

            return [
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
                            ariaDescription: app.attrs.tooltip("Requires the field value to be nonempty string"),
                        }),
                    ),
                ),
            ];
        },
    });
}
