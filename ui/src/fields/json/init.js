import { input } from "./input";
import { onrecordsave } from "./onrecordsave";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.json = {
    icon: "ri-braces-line",
    label: "JSON",
    settings,
    input,
    view,
    onrecordsave,
    dummyData: (f, forSubmit = false) => {
        return { "example": 123 };
    },
};
