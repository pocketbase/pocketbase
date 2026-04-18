import { input } from "./input";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.bool = {
    icon: "ri-toggle-line",
    label: "Bool",
    settings,
    input,
    view,
    dummyData: (f, forSubmit = false) => {
        return [true, false][Math.floor(Math.random() * 2)];
    },
};
