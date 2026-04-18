import { input } from "./input";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.url = {
    icon: "ri-link",
    label: "URL",
    settings,
    input,
    view,
    dummyData: (f, forSubmit = false) => {
        return "https://example.com";
    },
};
