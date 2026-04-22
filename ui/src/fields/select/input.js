// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "select_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-select" },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.select.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
            ),
            app.components.select({
                id: uniqueId,
                name: () => props.field.name,
                max: () => props.field.maxSelect || 1,
                required: () => props.field.required,
                options: () => {
                    return props.field.values.map((v) => {
                        return { value: v };
                    });
                },
                value: () => {
                    return app.utils.toArray(props.record[props.field.name]);
                },
                onchange: (opts) => {
                    if (props.field.maxSelect <= 1) {
                        props.record[props.field.name] = opts?.[0]?.value || "";
                        return;
                    }

                    props.record[props.field.name] = opts.map((o) => o.value);
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
