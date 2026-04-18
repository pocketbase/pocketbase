// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "editor_" + app.utils.randomString();

    const local = store({
        lazyEditor: null,
    });

    return t.div(
        {
            className: "record-field-input field-type-editor large-modal",
            onmount: () => {
                requestAnimationFrame(() => {
                    local.lazyEditor = app.components.tinymce({
                        id: uniqueId,
                        required: () => props.field.required,
                        convertURLs: () => props.field.convertURLs,
                        name: () => props.field.name,
                        value: () => props.record[props.field.name] || "",
                        onchange: (val) => {
                            props.record[props.field.name] = val;
                        },
                    });
                });
            },
        },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.editor.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            () => local.lazyEditor,
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
