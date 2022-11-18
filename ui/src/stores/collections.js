import { writable } from "svelte/store";
import ApiClient    from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";

export const collections          = writable([]);
export const activeCollection     = writable({});
export const isCollectionsLoading = writable(false);

export function changeActiveCollectionById(collectionId) {
    collections.update((list) => {
        const found = CommonHelper.findByKey(list, "id", collectionId);

        if (found) {
            activeCollection.set(found);
        } else if (list.length) {
            activeCollection.set(list[0]);
        }

        return list;
    });
}

// add or update collection
export function addCollection(collection) {
    activeCollection.update((current) => {
        return CommonHelper.isEmpty(current?.id) || current.id === collection.id ? collection : current;
    });

    collections.update((list) => {
        CommonHelper.pushOrReplaceByKey(list, collection, "id");
        return CommonHelper.sortCollections(list);
    });
}

export function removeCollection(collection) {
    collections.update((list) => {
        CommonHelper.removeByKey(list, "id", collection.id);

        activeCollection.update((current) => {
            if (current.id === collection.id) {
                return list[0];
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
            collections.set(CommonHelper.sortCollections(items));

            const item = activeId && CommonHelper.findByKey(items, "id", activeId);
            if (item) {
                activeCollection.set(item);
            } else if (items.length) {
                activeCollection.set(items[0]);
            }
        })
        .catch((err) => {
            ApiClient.errorResponseHandler(err);
        })
        .finally(() => {
            isCollectionsLoading.set(false);
        });
}
