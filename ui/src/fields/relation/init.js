import { input } from "./input";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.relation = {
    icon: "ri-mind-map",
    label: "Relation",
    settings,
    input,
    view,
    filterModifiers: (f) => {
        return f.maxSelect > 1 ? ["each", "length"] : [];
    },
    dummyData: (f, forSubmit = false) => {
        return f.maxSelect > 1 ? ["RECORD_ID1", "RECORD_ID2"] : "RECORD_ID";
    },
};
