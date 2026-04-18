// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div(
        { className: "record-field-view field-type-url" },
        () => {
            const value = props.record[props.field.name] || "";

            if (!value) {
                return t.span({ className: "missing-value" });
            }

            return t.a({
                href: () => value,
                className: "txt txt-ellipsis",
                rel: "noopener noreferrer",
                target: "_blank",
                textContent: app.utils.truncate(value),
                ariaDescription: app.attrs.tooltip("Open in new tab"),
                onclick: (e) => {
                    e.stopPropagation();
                },
            });
        },
    );
}
