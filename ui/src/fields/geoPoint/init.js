import { input } from "./input";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.geoPoint = {
    icon: "ri-map-pin-2-line",
    label: "Geo Point",
    settings,
    input,
    view,
    identifierExtractor: function(field, prefix = "") {
        return [prefix + field.name + ".lon", prefix + field.name + ".lat"];
    },
    dummyData: (f, forSubmit = false) => {
        return { lon: 0, lat: 0 };
    },
};
