window.app = window.app || {};
window.app.utils = window.app.utils || {};

const defaultOptions = {
    maxKeys: 30,
    requestKeys: true,
    collectionJoinKeys: true,
};

/**
 * Generates an array with the suitable autocomplete words for the targeted collection.
 *
 * @param  {string|Object} targetCollection Collection model or identifier.
 * @param  {string}  word The autocomplete triggered "word".
 * @param  {Object}  [options]
 * @param  {number}  [options.maxKeys] The max number of returned autocomplete keys (default to 30).
 * @param  {boolean} [options.requestKeys] Whether to include the `@request.*` keys (default to true).
 * @param  {boolean} [options.collectionJoinKeys] Whether to include the `@collection.*` keys (default to true).
 * @return {Array}
 */
window.app.utils.collectionAutocompleteKeys = function(targetCollection, word, options = {}) {
    if (!targetCollection || !word || !app.store.collections?.length) {
        return [];
    }

    options = Object.assign({}, defaultOptions, options);

    let result = collectionFieldsAutocomplete(word, app.store.collections, targetCollection).sort(keysSort);

    if (options.requestKeys) {
        const keys = requestFieldsAutocomplete(word, app.store.collections, targetCollection).sort(keysSort);
        for (let k of keys) {
            result.push(k);
        }
    }

    if (options.collectionJoinKeys) {
        const keys = collectionJoinAutocomplete(word, app.store.collections).sort(keysSort);
        for (let k of keys) {
            result.push(k);
        }
    }

    if (result.length > options.maxKeys) {
        return result.slice(0, options.maxKeys);
    }

    return result;
};

// sort shorter keys first
function keysSort(a, b) {
    return a.length - b.length;
}

/**
 * Generates recursively a list with all the autocomplete field keys
 * for the collectionNameOrId collection.
 *
 * @param  {string}        word
 * @param  {Array}         collections
 * @param  {string|object} collection
 * @param  {string}        [prefix]
 * @param  {number}        [level]
 * @return {Array}
 */
function collectionFieldsAutocomplete(word, collections, collection, prefix = "", level = 0) {
    if (!word || level >= 4) {
        return [];
    }

    if (typeof collection == "string") {
        collection = collections.find((c) => c.name == collection || c.id == collection);
    }
    if (!collection) {
        return [];
    }

    word = word.toLowerCase();

    const isAuth = collection.type == "auth";

    const result = app.utils
        .getAllCollectionIdentifiers(collection, prefix)
        .filter((item) => item.toLowerCase().includes(word));

    const fields = collection.fields || [];
    for (const field of fields) {
        if (field.type == "password" || (isAuth && field.name == "tokenKey")) {
            continue;
        }

        const keys = [];

        // special @request.body modifiers
        if (prefix == "@request.body.") {
            keys.push(prefix + field.name + ":changed");
            keys.push(prefix + field.name + ":isset");
        }

        if (typeof app.fieldTypes[field.type]?.filterModifiers == "function") {
            const modifiers = app.fieldTypes[field.type]?.filterModifiers(field) || [];
            for (const m of modifiers) {
                keys.push(prefix + field.name + ":" + m);
            }
        }

        for (const key of keys) {
            if (key.toLowerCase().includes(word)) {
                result.push(key);
            }
        }

        // add relation fields
        if (field.type == "relation" && field.collectionId) {
            const subKeys = collectionFieldsAutocomplete(
                word,
                collections,
                field.collectionId,
                prefix + field.name + ".",
                level + 1,
            );
            for (const k of subKeys) {
                result.push(k);
            }
        }
    }

    // add back relations
    for (const ref of collections) {
        const refFields = ref.fields || [];
        for (const field of refFields) {
            if (field.type != "relation" || field.collectionId != collection.id) {
                continue;
            }

            const key = prefix + ref.name + "_via_" + field.name;
            const subKeys = collectionFieldsAutocomplete(word, collections, ref, key + ".", level + 2); // +2 to reduce the recursive results
            for (const k of subKeys) {
                result.push(k);
            }
        }
    }

    return result;
}

/**
 * Generates a list with all @request.* autocomplete field keys.
 *
 * @param  {string}        word
 * @param  {Array}         collections
 * @param  {string|object} baseCollection (used for the `@request.body.*` fields)
 * @return {Array}
 */
function requestFieldsAutocomplete(word, collections, baseCollection) {
    if (!word) {
        return [];
    }

    word = word.toLowerCase();

    const result = [];

    const common = [
        "@request.context",
        "@request.method",
        "@request.query.",
        "@request.body.",
        "@request.headers.",
        "@request.auth.collectionId",
        "@request.auth.collectionName",
    ];
    for (const w of common) {
        if (!w.toLowerCase().includes(word)) {
            continue;
        }
        result.push(w);
    }

    // load auth collection fields
    const authCollections = collections.filter((collection) => collection.type === "auth");
    for (const collection of authCollections) {
        if (collection.system) {
            continue; // skip system collections for now
        }
        const authKeys = collectionFieldsAutocomplete(word, collections, collection, "@request.auth.");
        for (const k of authKeys) {
            app.utils.pushUnique(result, k);
        }
    }

    if (typeof baseCollection == "string") {
        baseCollection = collections.find((c) => c.name == baseCollection || c.id == baseCollection);
    }
    if (!baseCollection) {
        return result;
    }

    // load base collection fields into @request.body.*
    const keys = collectionFieldsAutocomplete(word, collections, baseCollection, "@request.body.");
    for (const key of keys) {
        result.push(key);
    }

    return result;
}

/**
 * Generates a list with all @collection.* autocomplete field keys.
 *
 * @param  {string} word
 * @param  {Array}  collections
 * @return {Array}
 */
function collectionJoinAutocomplete(word, collections) {
    const result = [];

    let basePrefix = "@collection.";

    // to avoid unnecessary loading all @collection.* keys match with the word first
    let base, search;
    if (basePrefix.length < word.length) {
        base = word;
        search = basePrefix;
    } else {
        base = basePrefix;
        search = word;
    }
    if (!base.includes(search)) {
        return result;
    }

    for (const collection of collections) {
        if (collection.system) {
            continue; // skip system collections for now
        }

        const keys = collectionFieldsAutocomplete(word, collections, collection, basePrefix + collection.name + ".");
        for (const key of keys) {
            result.push(key);
        }
    }

    return result;
}
