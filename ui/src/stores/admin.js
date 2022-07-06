import { writable } from "svelte/store";

// logged app admin
export const admin = writable({});

export function setAdmin(model) {
    admin.set(model || {});
}
