import { writable } from "svelte/store";

function getPreferenceFromLocalStorage() {
    const isDarkTheme = window.localStorage.getItem('pb_darkTheme');
    if (isDarkTheme === null)
        return null;
    if (isDarkTheme === "true")
        return true;
    return false;
}

function getPreferenceFromSystem() {
    return window.matchMedia("(prefers-color-scheme: dark)").matches === true;
}

function getPreference() {
    const preferenceLocalStorage = getPreferenceFromLocalStorage();

    if (preferenceLocalStorage === null)
        return getPreferenceFromSystem();
    return preferenceLocalStorage;
}

export const darkTheme = writable(getPreference());

export function setDarkTheme(value) {
    window.localStorage.setItem("pb_darkTheme", value === true);
    darkTheme.set(value === true);
}