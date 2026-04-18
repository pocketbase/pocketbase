import { input } from "./input";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.editor = {
    icon: "ri-edit-2-line",
    label: "Rich editor",
    settings,
    input,
    view,
    filterModifiers: (f) => {
        return ["lower"];
    },
    dummyData: (f, forSubmit = false) => {
        return "Lorem ipsum dolor sit amet...";
    },
};
