import { writable } from "svelte/store";

// logged app superuser
export const superuser = writable({});

export function setSuperuser(model) {
    superuser.set(model || {});
}
