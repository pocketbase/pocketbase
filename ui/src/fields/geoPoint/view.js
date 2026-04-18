// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(data) {
    return t.div(
        { className: "record-field-view field-type-geoPoint" },
        t.span({ className: "label" }, () => {
            const coords = data.record[data.field.name];
            return `${coords?.lon || 0}, ${coords?.lat || 0}`;
        }),
    );
}
