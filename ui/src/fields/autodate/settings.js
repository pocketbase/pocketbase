// {
//     originalCollection: undefined,
//     collection: undefined,
//     field
//     get fieldIndex: int/-1,
//     get originalField: undefined
// }
export function settings(props) {
    const ON_CREATE = 1;
    const ON_UPDATE = 2;
    const ON_CREATE_UPDATE = 3;

    const options = [
        { label: "Create", value: ON_CREATE },
        { label: "Update", value: ON_UPDATE },
        { label: "Create/Update", value: ON_CREATE_UPDATE },
    ];

    function getOptionFromField(field) {
        if (field.onCreate && field.onUpdate) {
            return ON_CREATE_UPDATE;
        }

        if (field.onUpdate) {
            return ON_UPDATE;
        }

        return ON_CREATE;
    }

    function updateField(option) {
        switch (option) {
            case ON_CREATE:
                props.field.onCreate = true;
                props.field.onUpdate = false;
                break;
            case ON_UPDATE:
                props.field.onCreate = false;
                props.field.onUpdate = true;
                break;
            case ON_CREATE_UPDATE:
                props.field.onCreate = true;
                props.field.onUpdate = true;
                break;
        }
    }

    const local = store({
        isDropdownOpen: false,
    });

    return app.components.fieldSettings(props, {
        header: t.div(
            {
                className: "field header-select autodate-select",
                ariaDescription: app.attrs.tooltip("Auto set on", "left"),
                onmount: () => {
                    // init default value
                    updateField(getOptionFromField(props.field));
                },
            },
            app.components.select({
                required: true,
                options: options,
                disabled: () => props.originalCollection?.system,
                value: () => getOptionFromField(props.field),
                onchange: (opts) => updateField(opts?.[0]?.value),
                ondropdowntoggle: (e) => {
                    local.isDropdownOpen = e.newState == "open";
                },
            }),
        ),
    });
}
