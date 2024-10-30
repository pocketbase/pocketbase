import ApiClient from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";
import { get, writable } from "svelte/store";

export const collections = writable([]);
export const activeCollection = writable({});
export const isCollectionsLoading = writable(false);
export const protectedFilesCollectionsCache = writable({});
export const scaffolds = writable({});

let notifyChannel;

if (typeof BroadcastChannel != "undefined") {
    notifyChannel = new BroadcastChannel("collections");

    notifyChannel.onmessage = () => {
        loadCollections(get(activeCollection)?.id)
    }
}

function notifyOtherTabs() {
    notifyChannel?.postMessage("reload");
}

export function changeActiveCollectionByIdOrName(collectionIdOrName) {
    collections.update((list) => {
        const found = list.find((c) => c.id == collectionIdOrName || c.name == collectionIdOrName);
        if (found) {
            activeCollection.set(found);
        } else if (list.length) {
            activeCollection.set(list.find((c) => !c.system) || list[0]);
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

        notifyOtherTabs();

        return CommonHelper.sortCollections(list);
    });
}

export function removeCollection(collection) {
    collections.update((list) => {
        CommonHelper.removeByKey(list, "id", collection.id);

        activeCollection.update((current) => {
            if (current.id === collection.id) {
                return list.find((c) => !c.system) || list[0];
            }
            return current;
        });

        refreshProtectedFilesCollectionsCache();

        notifyOtherTabs();

        return list;
    });
}

// load all collections
export async function loadCollections(activeIdOrName = null) {
    isCollectionsLoading.set(true);

    try {
        let items = await ApiClient.collections.getFullList(200, {
            "sort": "+name",
        })

        items = CommonHelper.sortCollections(items);

        collections.set(items);

        const item = activeIdOrName && items.find((c) => c.id == activeIdOrName || c.name == activeIdOrName);
        if (item) {
            activeCollection.set(item);
        } else if (items.length) {
            activeCollection.set(items.find((c) => !c.system) || items[0]);
        }

        refreshProtectedFilesCollectionsCache();

        scaffolds.set(await ApiClient.collections.getScaffolds());
    } catch (err) {
        ApiClient.error(err);
    }

    isCollectionsLoading.set(false);
}

function refreshProtectedFilesCollectionsCache() {
    protectedFilesCollectionsCache.update((cache) => {
        collections.update((current) => {
            for (let c of current) {
                cache[c.id] = !!c.fields?.find((f) => f.type == "file" && f.protected);
            }

            return current;
        });

        return cache;
    });
}

