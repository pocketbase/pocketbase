import { writable } from "svelte/store";

// eg.
// {
//   "text":        "Do you really want to delete the selectedItem",
//   "yesCallback": function() {...},
//   "noCallback":  function() {...},
// }
export const confirmation = writable({});

/**
 * @param {String}   text
 * @param {Function} [yesCallback]
 * @param {Function} [noCallback]
 */
export function confirm(text, yesCallback, noCallback) {
    confirmation.set({
        text:        text,
        yesCallback: yesCallback,
        noCallback:  noCallback,
    });
}

export function resetConfirmation() {
    confirmation.set({});
}
