import { writable } from "svelte/store";

const createDarkModeStore = () => {
    const defaultValue = false;
    const initialValue = typeof window !== "undefined" ? localStorage.getItem("darkMode") === "true" : defaultValue;
    
    const { subscribe, set, update } = writable(initialValue);

    return {
        subscribe,
        set: (value) => {
            if (typeof window !== "undefined") {
                localStorage.setItem("darkMode", value.toString());
                document.documentElement.setAttribute("data-theme", value ? "dark" : "light");
            }
            set(value);
        },
        update,
        toggle: () => update(n => {
            const newValue = !n;
            if (typeof window !== "undefined") {
                localStorage.setItem("darkMode", newValue.toString());
                document.documentElement.setAttribute("data-theme", newValue ? "dark" : "light");
            }
            return newValue;
        })
    };
};

export const darkMode = createDarkModeStore();
