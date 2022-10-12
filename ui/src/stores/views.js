import { writable } from "svelte/store";
import ApiClient from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";

export const views = writable([]);
export const activeView = writable({});
export const isViewLoading = writable(false);

// add or update collection
export function addView(view) {
    activeView.update((current) => {
        return CommonHelper.isEmpty(current?.id) || current.id === view.id ? view : current;
    });

    views.update((list) => {
        CommonHelper.pushOrReplaceByKey(list, view, "id");
        return list;
    });
}

export function removeView(view) {
    views.update((list) => {
        CommonHelper.removeByKey(list, "id", view.id);

        activeView.update((current) => {
            return current;
        });

        return list;
    });
}

// load all views
export async function loadViews(activeId = null) {
    isViewLoading.set(true);

    activeView.set({});
    views.set([]);

    return ApiClient.views
        .getFullList(200)
        .then((items) => {
            views.set(items);

            const item = activeId && CommonHelper.findByKey(items, "id", activeId);
            if (item) {
                activeView.set(item);
            } else if (items.length) {
                // fallback to the first non-profile collection item
                const nonProfile = items.find((c) => c.name != import.meta.env.PB_PROFILE_COLLECTION);
                if (nonProfile) {
                    activeView.set(nonProfile);
                }
            }
        })
        .catch((err) => {
            ApiClient.errorResponseHandler(err);
        })
        .finally(() => {
            isViewLoading.set(false);
        });
}
