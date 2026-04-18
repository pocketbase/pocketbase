// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div({ className: "record-field-view field-type-text" }, () => {
        const value = props.record[props.field.name] || "";

        if (value == "") {
            return t.span({ className: "missing-value" });
        }

        if (props.field.primaryKey) {
            let superuserYou = null;
            if (
                props.record?.collectionName == "_superusers"
                && app.store.superuser?.id == props.record?.id
            ) {
                superuserYou = t.strong({
                    className: "txt",
                    textContent: " (you)",
                });
            }

            return t.span(
                { className: "label" },
                app.components.copyButton(value),
                t.span({ className: "txt-ellipsis" }, app.utils.truncate(value), superuserYou),
            );
        }

        if (props.short) {
            return t.span({
                className: "txt txt-ellipsis",
                textContent: app.utils.truncate(value),
            });
        }

        return value;
    });
}
