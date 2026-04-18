// {
//     record: undefined,
//     field: undefined,
//     short: false,
// }
export function view(props) {
    return t.div({ className: "record-field-view field-type-file" }, () => {
        const filenames = app.utils.toArray(props.record[props.field.name]);
        if (!filenames.length) {
            return t.span({ className: "missing-value" });
        }

        const result = [];

        // truncate "full" view too to prevent freezing the browser tab
        const maxIndex = props.short ? 5 : 100;

        for (let i = 0; i < filenames.length; i++) {
            if (i >= maxIndex) {
                result.push(t.span({ className: "thumb sm" }, "+", filenames.length - maxIndex));
                break;
            }

            result.push(
                app.components.recordFileThumb({
                    record: props.record,
                    filename: filenames[i],
                }),
            );
        }

        return result;
    });
}
