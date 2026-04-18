// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(props) {
    const uniqueId = "f_" + app.utils.randomString();

    const cascadeOptions = [
        { label: "False", value: false },
        { label: "True", value: true },
    ];

    const isMultipleOptions = [
        { label: "Single", value: false },
        { label: "Multiple", value: true },
    ];

    const watchers = [
        // reset minSelect
        watch(
            () => props.field.maxSelect,
            (maxSelect) => {
                maxSelect = maxSelect || 1;
                if (maxSelect <= 1) {
                    props.field.minSelect = 0;
                }
            },
        ),
    ];

    return app.components.fieldSettings(props, {
        header: [
            t.div(
                {
                    className: "field header-select collections-select",
                    onunmount: () => {
                        watchers.forEach((w) => w?.unwatch());
                    },
                },
                app.components.select({
                    required: true,
                    className: "inline-error",
                    placeholder: "Select collection*",
                    name: () => `fields.${props.fieldIndex}.collectionId`,
                    disabled: () => !!props.originalField?.id,
                    options: () =>
                        app.utils.sortedCollections(app.store.collections.filter((c) => c.type != "view")).map(
                            (c) => {
                                return { value: c.id, label: c.name };
                            },
                        ),
                    value: () => props.field.collectionId,
                    onchange: (opts) => {
                        props.field.collectionId = opts?.[0]?.value || "";
                    },
                    after: () => {
                        return [
                            t.hr({ className: "m-t-5 m-b-5" }),
                            t.button(
                                {
                                    type: "button",
                                    className: "btn sm outline",
                                    onclick: () => {
                                        app.modals.openCollectionUpsert({}, {
                                            onsave: (newCollection) => {
                                                props.field.collectionId = newCollection.id;
                                            },
                                        });
                                    },
                                },
                                t.i({ className: "ri-add-line", ariaHidden: true }),
                                t.span({ className: "txt" }, "New collection"),
                            ),
                        ];
                    },
                }),
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
                            if (props.field.maxSelect << 0 < 2) {
                                props.field.maxSelect = 10;
                            }
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
                t.div(
                    { className: "col-sm-6", hidden: () => props.field.maxSelect << 0 < 2 },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".minSelect" }, "Min select"),
                        t.input({
                            type: "number",
                            id: uniqueId + ".minSelect",
                            step: 1,
                            min: 0,
                            max: Number.MAX_SAFE_INTEGER,
                            placeholder: "No min limit",
                            name: () => `fields.${props.fieldIndex}.minSelect`,
                            value: () => props.field.minSelect || "",
                            onchange: (e) => (props.field.minSelect = parseInt(e.target.value, 10)),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-sm-6", hidden: () => props.field.maxSelect << 0 < 2 },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".maxSelect" }, "Max select"),
                        t.input({
                            type: "number",
                            id: uniqueId + ".maxSelect",
                            step: 1,
                            min: () => props.field.minSelect || 2,
                            max: Number.MAX_SAFE_INTEGER,
                            placeholder: "Default to single",
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
                ),
                t.div(
                    { className: "col-sm-12" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".cascadeDelete" }, "Cascade delete"),
                        app.components.select({
                            required: true,
                            id: uniqueId + ".cascadeDelete",
                            name: () => `fields.${props.fieldIndex}.cascadeDelete`,
                            options: cascadeOptions,
                            value: () => props.field.cascadeDelete || false,
                            onchange: (opts) => {
                                props.field.cascadeDelete = !!opts?.[0].value;
                            },
                        }),
                    ),
                ),
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
                    t.small({ className: "txt-hint" }, "(!='')"),
                    t.i({
                        className: "ri-information-line link-hint",
                        ariaDescription: app.attrs.tooltip("Requires the field value to be nonempty string"),
                    }),
                ),
            ),
        ],
    });
}
