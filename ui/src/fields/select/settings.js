// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(props) {
    const uniqueId = "f_" + app.utils.randomString();

    const isMultipleOptions = [
        { label: "Single", value: false },
        { label: "Multiple", value: true },
    ];

    const optionsDropdown = t.div(
        {
            popover: "manual",
            className: "dropdown field-select-choices-dropdown",
        },
        t.div({ className: "field-help m-t-0", style: "font-size: 0.9em" }, "New-line separated choices:"),
        t.div(
            { className: "field" },
            t.textarea({
                className: "autoexpand",
                required: true,
                value: () => {
                    const vals = app.utils.toArray(props.field.values, false);
                    return vals.join("\n");
                },
                oninput: (e) => {
                    const vals = e.target.value.trimStart().replaceAll("\n\n", "\n").split("\n");
                    props.field.values = vals;

                    // clear previous errors
                    app.utils.deleteByPath(app.store.errors, `fields.${props.fieldIndex}.values`);
                },
                onchange: (e) => {
                    // filter duplicates and empty values
                    const unique = new Set();
                    const vals = e.target.value.split("\n");
                    for (let val of vals) {
                        if (val == "") {
                            continue;
                        }
                        unique.add(val);
                    }

                    props.field.values = Array.from(unique);
                },
                onblur: (e) => {
                    if (!e.relatedTarget || !optionsDropdown.contains(e.relatedTarget)) {
                        optionsDropdown.hidePopover();
                    }
                },
            }),
        ),
    );

    const watchers = [
        // cap maxSelect value
        watch(() => {
            if (props.field.values?.length && props.field.maxSelect > props.field.values.length) {
                props.field.maxSelect = props.field.values.length;
            }
        }),
    ];

    return app.components.fieldSettings(props, {
        header: [
            t.div(
                {
                    className: "field header-select field-select-choices-input",
                    onunmount: () => {
                        watchers.forEach((w) => w?.unwatch());
                    },
                },
                t.input({
                    type: "text",
                    placeholder: "Add choices*",
                    className: "txt-left inline-error",
                    value: () => props.field.values?.join(" • ") || "",
                    name: () => `fields.${props.fieldIndex}.values`,
                    onfocus: (e) => {
                        optionsDropdown?.showPopover({ source: e.target });
                        optionsDropdown.querySelector("textarea")?.focus();
                        return false;
                    },
                }),
                optionsDropdown,
            ),
            t.div(
                {
                    className: "field header-select single-multiple-select",
                },
                app.components.select({
                    required: true,
                    options: isMultipleOptions,
                    value: () => {
                        return props.field.maxSelect > 1;
                    },
                    onchange: (opts) => {
                        if (opts?.[0]?.value) {
                            props.field.maxSelect = props.field.values.length || 2;
                        } else {
                            props.field.maxSelect = 1;
                        }
                    },
                }),
            ),
        ],
        content: () =>
            t.div(
                { className: "grid sm" },
                () => {
                    if (props.field.maxSelect > 1) {
                        return t.div(
                            { className: "col-sm-12" },
                            t.div(
                                { className: "field" },
                                t.label({ htmlFor: uniqueId + ".maxSelect" }, "Max select"),
                                t.input({
                                    type: "number",
                                    id: uniqueId + ".maxSelect",
                                    placeholder: "Default to single",
                                    step: 1,
                                    min: 2,
                                    max: () => props.field.values?.length || 2,
                                    name: () => `fields.${props.fieldIndex}.maxSelect`,
                                    value: () => props.field.maxSelect || "",
                                    onchange: (e) => {
                                        const maxSelect = parseInt(e.target.value, 10);
                                        if (maxSelect > 1) {
                                            props.field.maxSelect = maxSelect;
                                        } else {
                                            props.field.maxSelect = 1;
                                        }
                                    },
                                }),
                            ),
                        );
                    }
                },
                t.div(
                    { className: "col-sm-12" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".help" }, "Help text"),
                        t.input({
                            type: "text",
                            id: uniqueId + ".help",
                            name: () => `fields.${props.fieldIndex}.help`,
                            value: () => props.field.help || "",
                            oninput: (e) => (props.field.help = e.target.value),
                        }),
                    ),
                ),
            ),
        footer: () => [
            t.div(
                { className: "field" },
                t.input({
                    className: "sm",
                    type: "checkbox",
                    id: uniqueId + ".required",
                    name: () => `fields.${props.fieldIndex}.required`,
                    checked: () => !!props.field.required,
                    onchange: (e) => (props.field.required = e.target.checked),
                }),
                t.label(
                    { htmlFor: uniqueId + ".required" },
                    t.span({ className: "txt" }, "Required"),
                    t.small({
                        className: "txt-hint",
                        textContent: () => (props.field.maxSelect > 1 ? "(!=[])" : "(!='')"),
                    }),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip(() => {
                            return `Requires the field value to be nonempty ${
                                props.field.maxSelect > 1 ? "array" : "string"
                            }.`;
                        }),
                    }),
                ),
            ),
        ],
    });
}
