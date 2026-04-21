// {
//     get collection: undefined,
//     get originalRecord: undefined,
//     get record: undefined,
//     get field: undefined,
// }
export function input(data) {
    const uniqueId = "geo_" + app.utils.randomString();

    const local = store({
        showMap: false,
    });

    return t.div(
        { className: "record-field-input field-type-geoPoint" },
        t.div(
            { className: () => `field-list ${data.field.required ? "required" : ""}` },
            t.label(
                { htmlFor: uniqueId },
                t.i({ className: app.fieldTypes.geoPoint.icon, ariaHidden: true }),
                t.span({ className: "txt" }, () => data.field.name),
            ),
            t.div(
                { className: "field-list-content" },
                t.div(
                    { className: "field-list-item p-0" },
                    t.div(
                        { className: "fields" },
                        t.div({ className: "field addon" }, t.label({ htmlFor: uniqueId + ".lon" }, "Longitude:")),
                        t.div(
                            { className: "field" },
                            t.input({
                                id: uniqueId + ".lon",
                                type: "number",
                                step: "any",
                                min: -180,
                                max: 180,
                                placeholder: 0,
                                name: () => data.field.name,
                                required: () => data.field.required,
                                value: () => data.record[data.field.name]?.lon || "",
                                onchange: (e) => {
                                    data.record[data.field.name] = data.record[data.field.name] || {};
                                    data.record[data.field.name].lon = Number(e.target.value);
                                },
                            }),
                        ),
                        t.span({ className: "delimiter" }),
                        t.div({ className: "field addon" }, t.label({ htmlFor: uniqueId + ".lat" }, "Latitude:")),
                        t.div(
                            { className: "field" },
                            t.input({
                                id: uniqueId + ".lat",
                                type: "number",
                                step: "any",
                                min: -90,
                                max: 90,
                                placeholder: 0,
                                name: () => data.field.name,
                                required: () => data.field.required,
                                value: () => data.record[data.field.name]?.lat || "",
                                onchange: (e) => {
                                    data.record[data.field.name] = data.record[data.field.name] || {};
                                    data.record[data.field.name].lat = Number(e.target.value);
                                },
                            }),
                        ),
                        t.span({ className: "delimiter" }),
                        t.div(
                            { className: "field addon p-5" },
                            t.button(
                                {
                                    type: "button",
                                    className: () => `btn sm circle secondary ${local.showMap ? "" : "transparent"}`,
                                    onclick: () => (local.showMap = !local.showMap),
                                },
                                t.i({ className: "ri-map-2-line" }),
                            ),
                        ),
                    ),
                ),
                () => {
                    if (!local.showMap) {
                        return;
                    }

                    return t.div(
                        { className: "field-list-item p-0", style: "height: 250px" },
                        app.components.leaflet({
                            point: () => data.record[data.field.name] || { lat: 0, lon: 0 },
                            onchange: (newPoint) => {
                                data.record[data.field.name] = structuredClone(newPoint);
                            },
                        }),
                    );
                },
            ),
        ),
        () => {
            if (data.field.help) {
                return t.div({ className: "field-help" }, data.field.help);
            }
        },
    );
}
