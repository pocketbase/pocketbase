// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "text_" + app.utils.randomString();

    const data = store({
        get hasAutogenerate() {
            return !app.utils.isEmpty(props.field.autogeneratePattern) && app.utils.isEmpty(props.originalRecord?.id);
        },
        get isDisabled() {
            return !app.utils.isEmpty(props.originalRecord?.id) && props.field.primaryKey;
        },
        get isRequired() {
            return props.field.required && !data.hasAutogenerate && !data.isDisabled;
        },
    });

    return t.div(
        { className: "record-field-input field-type-text" },
        t.div(
            { className: "fields" },
            t.div(
                { className: "field" },
                t.label(
                    { htmlFor: uniqueId },
                    t.i({
                        ariaHidden: true,
                        className: () => (props.field.primaryKey ? "ri-key-line" : app.fieldTypes.text.icon),
                    }),
                    t.span({ className: "txt" }, () => props.field.name),
                ),
                // @todo remove after Firefox add support for "field-sizing:content"
                //
                // note1: not using contenteditable because it requires keeping
                // track of the cursor offset when replacing the content in oninput
                //
                // note2: based on https://chriscoyier.net/2023/09/29/css-solves-auto-expanding-textareas-probably-eventually/
                t.div(
                    { className: "autoexpand-wrapper" },
                    t.textarea({
                        id: uniqueId,
                        // className: "autoexpand",
                        rows: 1,
                        name: () => props.field.name,
                        required: () => data.isRequired,
                        disabled: () => data.isDisabled,
                        placeholder: () => (data.hasAutogenerate ? "Leave empty to autogenerate..." : ""),
                        value: () => props.record[props.field.name] || "",
                        oninput: (e) => (props.record[props.field.name] = e.target.value || ""),
                    }),
                    t.div(
                        { className: "input" },
                        () => props.record[props.field.name],
                        // the empty space is necessary to prevent jumpy behavior
                        " ",
                    ),
                ),
            ),
            // list the autodate field values in a tooltip next to the primary key
            () => {
                if (!props.field.primaryKey || !props.originalRecord?.id) {
                    return;
                }

                const autodateFields = props.collection?.fields?.filter((f) => f.type == "autodate") || [];
                if (!autodateFields.length) {
                    return;
                }

                const autodateValues = [];
                for (let f of autodateFields) {
                    autodateValues.push(`${f.name}: ${app.utils.stringifyValue(props.record[f.name])}`);
                }

                return t.div(
                    { className: "field addon" },
                    t.i({
                        className: "ri-information-line txt-hint link-faded",
                        ariaDescription: app.attrs.tooltip(autodateValues.join("\n"), "left"),
                    }),
                );
            },
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}
