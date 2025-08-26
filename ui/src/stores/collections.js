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

export async function refreshScaffolds() {
    scaffolds.set(await ApiClient.collections.getScaffolds());
}

// load all collections
export async function loadCollections(activeIdOrName = null) {
    isCollectionsLoading.set(true);

    try {
        const promises = [];
        promises.push(ApiClient.collections.getScaffolds());
        promises.push(ApiClient.collections.getFullList());

        let [resultScaffolds, resultCollections] = await Promise.all(promises);

        scaffolds.set(resultScaffolds);

        resultCollections = CommonHelper.sortCollections(resultCollections);

        collections.set(resultCollections);

        const found = activeIdOrName && resultCollections.find((c) => c.id == activeIdOrName || c.name == activeIdOrName);
        if (found) {
            activeCollection.set(found);
        } else if (resultCollections.length) {
            activeCollection.set(resultCollections.find((c) => !c.system) || resultCollections[0]);
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
                cache[c.id] = !!c.fields?.find((f) => f.type == "file" && f.protected);
            }

            return current;
        });

        return cache;
    });
}
