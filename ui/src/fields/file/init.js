import { input } from "./input";
import { onrecordduplicate } from "./onrecordduplicate";
import { onrecordsave } from "./onrecordsave";
import { settings } from "./settings";
import { view } from "./view";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.file = {
    icon: "ri-image-line",
    label: "File",
    settings,
    input,
    view,
    summaryPriority: -1,
    onrecordsave,
    onrecordduplicate,
    filterModifiers: (f) => {
        return f.maxSelect > 1 ? ["each", "length"] : [];
    },
    dummyData: (f, forSubmit = false) => {
        if (forSubmit) {
            if (f.maxSelect > 1) {
                return [dummyFileObject("test1.txt"), dummyFileObject("test2.txt")];
            }

            return dummyFileObject("test1.txt");
        }

        if (f.maxSelect > 1) {
            return [
                "test1_" + app.utils.randomString(10) + ".txt",
                "test2_" + app.utils.randomString(10) + ".txt",
            ];
        }

        return "test_" + app.utils.randomString(10) + ".txt";
    },
};

function dummyFileObject(name) {
    return {
        toString() {
            return `new File([...], '${name}')`;
        },
        toJSON() {
            // "[[ and ]]" will have to be manualy replaced after JSON.stringify
            return `[[new File([...], '${name}')]]`;
        },
    };
}
