// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div(
        { className: "record-field-view field-type-editor" },
        () => {
            if (props.short) {
                const value = props.record[props.field.name];
                if (!value) {
                    return t.span({ className: "missing-value" });
                }

                return t.span({
                    className: "txt",
                    textContent: app.utils.truncate(app.utils.plainText(value), 200),
                });
            }

            return app.components.tinymce({
                readonly: true,
                className: "large-modal",
                value: () => props.record[props.field.name] || "",
            });
        },
    );
}
