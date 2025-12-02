import { writable } from "svelte/store";

export const pageTitle = writable('');

export const appName = writable('');

export const hideControls = writable(false);

const initialDark = typeof localStorage !== "undefined" && localStorage.getItem("pb_theme") === "dark";
if (initialDark) {
    document.documentElement.classList.add("theme-dark");
}

export const isDarkMode = writable(initialDark);

isDarkMode.subscribe((v) => {
    if (v) {
        document.documentElement.classList.add("theme-dark");
        if (typeof localStorage !== "undefined") {
            localStorage.setItem("pb_theme", "dark");
        }
    } else {
        document.documentElement.classList.remove("theme-dark");
        if (typeof localStorage !== "undefined") {
            localStorage.setItem("pb_theme", "light");
        }
    }
});
