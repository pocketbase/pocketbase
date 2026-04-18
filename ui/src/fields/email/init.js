import { input } from "./input";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.email = {
    icon: "ri-mail-line",
    label: "Email",
    settings,
    input,
    view,
    filterModifiers: (f) => {
        return ["lower"];
    },
    dummyData: (f, forSubmit = false) => {
        return `test_${app.utils.randomString(3, "123567890")}@example.com`;
    },
};
