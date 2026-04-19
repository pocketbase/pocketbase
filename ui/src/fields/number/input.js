// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "number_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-number" },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.number.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            t.input({
                type: "number",
                id: uniqueId,
                step: "any",
                name: () => props.field.name,
                required: () => props.field.required,
                min: () => props.field.min,
                max: () => props.field.max,
                value: () => props.record[props.field.name] || "",
                oninput: (e) => props.record[props.field.name] = Number(e.target.value),
            }),
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
