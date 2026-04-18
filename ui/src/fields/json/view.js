// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div({ className: "record-field-view field-type-json" }, () => {
        const rawValue = props.record[props.field.name];

        if (props.short) {
            return t.span({
                className: "txt-code txt-ellipsis",
                textContent: app.utils.truncate(app.utils.trimQuotedValue(JSON.stringify(rawValue)) || ""),
            });
        }

        return app.components.codeBlock({
            value: () => JSON.stringify(rawValue, null, 2),
        });
    });
}
