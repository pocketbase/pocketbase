// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "date_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-date" },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.date.icon }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            t.input({
                id: uniqueId,
                step: 1,
                type: "datetime-local",
                name: () => props.field.name,
                required: () => props.field.required,
                value: () => app.utils.toDatetimeLocalInputValue(props.record[props.field.name]),
                onchange: (e) => {
                    props.record[props.field.name] = app.utils.toRFC3339Datetime(e.target.value);
                },
            }),
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
