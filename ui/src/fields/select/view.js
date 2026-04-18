// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div(
        { className: "record-field-view field-type-select" },
        t.div({ className: "inline-flex gap-5" }, () => {
            const opts = app.utils.toArray(props.record[props.field.name], false);

            if (!opts.length) {
                return t.span({ className: "missing-value" });
            }

            return opts.map((opt) => {
                return t.span({
                    className: "label",
                    title: opt,
                    textContent: app.utils.truncate(opt, 100),
                });
            });
        }),
    );
}
