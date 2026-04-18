// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "bool_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-bool" },
        t.div(
            { className: "field" },
            t.input({
                type: "checkbox",
                id: uniqueId,
                className: "switch",
                name: () => props.field.name,
                required: () => props.field.required,
                checked: () => props.record[props.field.name] || false,
                onchange: (e) => (props.record[props.field.name] = e.target.checked || false),
            }),
            t.label({ htmlFor: uniqueId }, () => props.field.name),
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
