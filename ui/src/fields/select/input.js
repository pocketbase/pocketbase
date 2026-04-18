// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(data) {
    const uniqueId = "select_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-select" },
        t.div(
            { className: "field" },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.select.icon }),
                t.span({ className: "txt" }, () => data.field.name),
            ),
            app.components.select({
                id: uniqueId,
                max: () => data.field.maxSelect || 1,
                required: () => data.field.required,
                options: () => {
                    return data.field.values.map((v) => {
                        return { value: v };
                    });
                },
                value: () => {
                    return app.utils.toArray(data.record[data.field.name]);
                },
                onchange: (opts) => {
                    if (data.field.maxSelect <= 1) {
                        data.record[data.field.name] = opts?.[0]?.value || "";
                        return;
                    }

                    data.record[data.field.name] = opts.map((o) => o.value);
                },
            }),
        ),
        () => {
            if (data.field.help) {
                return t.div({ className: "field-help" }, data.field.help);
            }
        },
    );
}
