window.app = window.app || {};
window.app.components = window.app.components || {};

window.app.components.addCollectionFieldButton = function(collection) {
    const uniqueId = "new_field_" + app.utils.randomString();

    function addNewField(fieldType) {
        const field = {
            id: "",
            name: getUniqueFieldName(fieldType),
            type: fieldType,
            system: false,
            hidden: false,
            presentable: false,
            required: false,

            __focus: true, // see fieldSettings
        };

        collection.fields = collection.fields || [];

        // if the collection has created/updated last fields,
        // insert before the first autodate field, otherwise - append
        const idx = collection.fields.findLastIndex((f) => f.type != "autodate");
        if (field.type != "autodate" && idx >= 0) {
            collection.fields.splice(idx + 1, 0, field);
        } else {
            collection.fields.push(field);
        }
    }

    function getUniqueFieldName(baseName = "") {
        let result = baseName;
        let counter = 2;

        let suffix = baseName.match(/\d+$/)?.[0] || ""; // extract numeric suffix

        // name without the suffix
        let base = suffix ? baseName.substring(0, baseName.length - suffix.length) : baseName;

        while (hasFieldWithName(result)) {
            result = base + ((suffix << 0) + counter);
            counter++;
        }

        return result;
    }

    function hasFieldWithName(name) {
        return !!collection.fields?.find((f) => f.name.toLowerCase() === name.toLowerCase());
    }

    return t.div(
        { className: "new-collection-field-btn-wrapper" },
        t.button(
            {
                type: "button",
                className: "btn block outline",
                "html-popovertarget": uniqueId + "_dropdown",
            },
            t.i({ className: "ri-add-line", ariaHidden: true }),
            t.span({ className: "txt" }, "New field"),
        ),
        t.div(
            {
                id: uniqueId + "_dropdown",
                className: "dropdown field-types-dropdown",
                popover: "auto",
            },
            () => {
                const options = [];

                for (const type in app.fieldTypes) {
                    // for now skip password field types
                    if (type == "password") {
                        continue;
                    }

                    const def = app.fieldTypes[type];
                    if (!def.settings) {
                        continue;
                    }

                    options.push(
                        t.button(
                            {
                                type: "button",
                                className: "dropdown-item",
                                onclick: (e) => {
                                    e.target.closest(".dropdown")?.hidePopover();
                                    addNewField(type);
                                },
                            },
                            t.i({ className: def.icon || app.utils.fallbackFieldIcon, ariaHidden: true }),
                            t.span({ className: "txt" }, def.label || type),
                        ),
                    );
                }

                return options;
            },
        ),
    );
};
