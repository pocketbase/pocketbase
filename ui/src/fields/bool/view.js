// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div(
        { className: "record-field-view field-type-bool" },
        t.span(
            { className: () => `label ${props.record[props.field.name] ? "success" : ""}` },
            () => props.record[props.field.name] ? "True" : "False",
        ),
    );
}
