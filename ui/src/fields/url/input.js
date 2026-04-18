// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "url_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-url" },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.url.icon }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            t.input({
                type: "url",
                id: uniqueId,
                spellcheck: false,
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
