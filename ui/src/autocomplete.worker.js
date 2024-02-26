import CommonHelper from "@/utils/CommonHelper";

const maxKeys = 11000;

onmessage = (e) => {
    if (!e.data.collections) {
        return;
    }

    const result = {};

    result.baseKeys = CommonHelper.getCollectionAutocompleteKeys(e.data.collections, e.data.baseCollection?.name);
    result.baseKeys = limitArray(result.baseKeys.sort(keysSort), maxKeys);

    if (!e.data.disableRequestKeys) {
        result.requestKeys = CommonHelper.getRequestAutocompleteKeys(e.data.collections, e.data.baseCollection?.name);
        result.requestKeys = limitArray(result.requestKeys.sort(keysSort), maxKeys);
    }

    if (!e.data.disableCollectionJoinKeys) {
        result.collectionJoinKeys = CommonHelper.getCollectionJoinAutocompleteKeys(e.data.collections);
        result.collectionJoinKeys = limitArray(result.collectionJoinKeys.sort(keysSort), maxKeys);
    }

    postMessage(result);
};

// sort shorter keys first
function keysSort(a, b) {
    return a.length - b.length;
}

function limitArray(arr, max) {
    if (arr.length > max)  {
        return arr.slice(0, max);
    }

    return arr;
}
