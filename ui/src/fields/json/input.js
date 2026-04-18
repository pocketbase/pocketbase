// {
//     collection: undefined,
//     originalRecord: undefined,
//     record: undefined,
//     field: undefined,
// }
export function input(props) {
    const uniqueId = "json_" + app.utils.randomString();

    const local = store({
        value: "",
    });

    const watchers = [
        watch(
            () => props.record[props.field.name],
            (newVal, oldVal) => {
                if (newVal !== "" && newVal === local.value) {
                    return;
                }

                // quote string values if not already
                if (typeof newVal == "string" && !newVal.startsWith("\"") && !newVal.endsWith("\"")) {
                    local.value = JSON.stringify(typeof newVal === "undefined" ? null : newVal);
                    props.record[props.field.name] = local.value;
                    return;
                }

                if (typeof newVal == "string" && newVal.startsWith("\"") && newVal.endsWith("\"")) {
                    local.value = newVal; // already double quoted
                } else if (newVal === null) {
                    local.value = "null";
                } else {
                    local.value = JSON.stringify(typeof newVal === "undefined" ? null : newVal, null, 2);
                }
            },
        ),
    ];

    function updateRecordValue() {
        const trimmed = local.value.trim();

        if (trimmed === "") {
            props.record[props.field.name] = null;
            return;
        }

        try {
            let parsed = JSON.parse(trimmed);
            if (typeof parsed == "string") {
                props.record[props.field.name] = JSON.stringify(parsed);
            } else {
                props.record[props.field.name] = parsed;
            }
        } catch (_) {
            props.record[props.field.name] = trimmed;
        }
    }

    let updateRecordValueTimeoutId;

    return t.div(
        { className: "record-field-input field-type-json" },
        t.div(
            {
                className: "field",
                onunmount: () => {
                    clearTimeout(updateRecordValueTimeoutId);
                    watchers.forEach((w) => w?.unwatch());
                },
            },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.json.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => props.field.name),
                t.span(
                    {
                        hidden: () => isValidStringifiedJSON(local.value.trim()),
                        className: "json-state",
                        ariaDescription: app.attrs.tooltip("Invalid JSON", "left"),
                    },
                    t.i({ className: "ri-error-warning-fill txt-danger", ariaHidden: true }),
                ),
                t.span(
                    {
                        hidden: () => !isValidStringifiedJSON(local.value.trim()),
                        className: "json-state",
                        ariaDescription: app.attrs.tooltip("Valid JSON", "left"),
                    },
                    t.i({ className: "ri-checkbox-circle-fill txt-success", ariaHidden: true }),
                ),
            ),
            app.components.codeEditor({
                language: "js",
                id: uniqueId,
                name: () => props.field.name,
                required: () => props.field.required,
                value: () => local.value,
                oninput: (val) => (local.value = val),
                onblur: () => updateRecordValue(),
            }),
        ),
        () => {
            if (props.field.help) {
                return t.div({ className: "field-help" }, props.field.help);
            }
        },
    );
}

function isValidStringifiedJSON(val) {
    if (val === "") {
        return true;
    }

    try {
        JSON.parse(val);
        return true;
    } catch (_) {
        return false;
    }
}
