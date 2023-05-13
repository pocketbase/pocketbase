import { writable } from "svelte/store";
import ApiClient    from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";

export const collections                    = writable([]);
export const activeCollection               = writable({});
export const isCollectionsLoading           = writable(false);
export const protectedFilesCollectionsCache = writable({});

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

        refreshProtectedFilesCollectionsCache();

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

        refreshProtectedFilesCollectionsCache();

        return list;
    });
}

// load all collections (excluding the user profile)
export async function loadCollections(activeId = null) {
    isCollectionsLoading.set(true);

    try {
        let items = await ApiClient.collections.getFullList(200, {
            "sort": "+name",
        })

        items = CommonHelper.sortCollections(items);

        collections.set(items);

        const item = activeId && CommonHelper.findByKey(items, "id", activeId);
        if (item) {
            activeCollection.set(item);
        } else if (items.length) {
            activeCollection.set(items[0]);
        }

        refreshProtectedFilesCollectionsCache();
    } catch (err) {
        ApiClient.error(err);
    }

    isCollectionsLoading.set(false);
}

function refreshProtectedFilesCollectionsCache() {
    protectedFilesCollectionsCache.update((cache) => {
        collections.update((current) => {
            for (let c of current) {
                cache[c.id] = !!c.schema?.find((f) => f.type == "file" && f.options?.protected);
            }

            return current;
        });

        return cache;
    });
}
