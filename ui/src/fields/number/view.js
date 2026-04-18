// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div(
        { className: "record-field-view field-type-number" },
        t.span({ className: "txt" }, () => props.record[props.field.name]),
    );
}
