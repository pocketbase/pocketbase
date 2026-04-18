import { settings } from "./settings";

window.app = window.app || {};
window.app.fieldTypes = window.app.fieldTypes || {};
window.app.fieldTypes.password = {
    icon: "ri-lock-password-line",
    label: "Password",
    settings,
};
