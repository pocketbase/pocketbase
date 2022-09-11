import { writable } from "svelte/store";
import ApiClient    from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";

export const collections          = writable([]);
export const activeCollection     = writable({});
export const isCollectionsLoading = writable(false);

// add or update collection
export function addCollection(collection) {
    activeCollection.update((current) => {
        return CommonHelper.isEmpty(current?.id) || current.id === collection.id ? collection : current;
    });

    collections.update((list) => {
        CommonHelper.pushOrReplaceByKey(list, collection, "id");
        return list;
    });
}

export function removeCollection(collection) {
    collections.update((list) => {
        CommonHelper.removeByKey(list, "id", collection.id);

        activeCollection.update((current) => {
            if (current.id === collection.id) {
                // fallback to the first non-profile collection item
                return list.find((c) => c.name != import.meta.env.PB_PROFILE_COLLECTION) || {}
            }
            return current;
        });

        return list;
    });

}

// load all collections (excluding the user profile)
export async function loadCollections(activeId = null) {
    isCollectionsLoading.set(true);

    activeCollection.set({});
    collections.set([]);

    return ApiClient.collections.getFullList(200, {
        "sort": "+created",
    })
        .then((items) => {
            collections.set(items);

            const item = activeId && CommonHelper.findByKey(items, "id", activeId);
            if (item) {
                activeCollection.set(item);
            } else if (items.length) {
                // fallback to the first non-profile collection item
                const nonProfile = items.find((c) => c.name != import.meta.env.PB_PROFILE_COLLECTION)
                if (nonProfile) {
                    activeCollection.set(nonProfile);
                }
            }
        })
        .catch((err) => {
            ApiClient.errorResponseHandler(err);
        })
        .finally(() => {
            isCollectionsLoading.set(false);
        });
}
