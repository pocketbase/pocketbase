// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "editor_" + app.utils.randomString();

    return t.div(
        {
            className: "record-field-input field-type-editor large-modal",
        },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.editor.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            () => {
                return app.components.tinymce({
                    id: uniqueId,
                    name: () => props.field.name,
                    required: () => props.field.required,
                    convertURLs: () => props.field.convertURLs,
                    value: () => props.record[props.field.name] || "",
                    onchange: (val) => {
                        props.record[props.field.name] = val;
                    },
                });
            },
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
