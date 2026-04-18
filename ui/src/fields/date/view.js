// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div(
        { className: "record-field-view field-type-date" },
        app.components.formattedDate({
            value: () => props.record[props.field.name],
            short: () => props.short,
        }),
    );
}
