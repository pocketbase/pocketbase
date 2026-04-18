// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "email_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-email" },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.email.icon }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            t.input({
                type: "email",
                id: uniqueId,
                spellcheck: false,
                autocomplete: false,
                name: () => props.field.name,
                required: () => props.field.required,
                value: () => props.record[props.field.name] || "",
                oninput: (e) => (props.record[props.field.name] = e.target.value),
            }),
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
