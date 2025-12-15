import { writable, get } from "svelte/store";

const STORAGE_KEY = "pocketbase_saved_filters";

/**
 * Load saved filters from localStorage.
 * Structure: { collectionId: [{ id, name, filter }] }
 */
function loadFromStorage() {
    try {
        const stored = localStorage.getItem(STORAGE_KEY);
        return stored ? JSON.parse(stored) : {};
    } catch (e) {
        console.warn("Failed to load saved filters from localStorage", e);
        return {};
    }
}

/**
 * Save filters to localStorage.
 */
function saveToStorage(data) {
    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
    } catch (e) {
        console.warn("Failed to save filters to localStorage", e);
    }
}

/**
 * Generate a unique ID for a saved filter.
 */
function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substring(2);
}

export const savedFilters = writable(loadFromStorage());

/**
 * Get saved filters for a specific collection.
 * @param {string} collectionId
 * @returns {Array<{id: string, name: string, filter: string}>}
 */
export function getSavedFilters(collectionId) {
    const all = get(savedFilters);
    return all[collectionId] || [];
}

/**
 * Add a new saved filter for a collection.
 * @param {string} collectionId
 * @param {string} name - Display name for the filter
 * @param {string} filter - The filter string
 */
export function addSavedFilter(collectionId, name, filter) {
    savedFilters.update((all) => {
        if (!all[collectionId]) {
            all[collectionId] = [];
        }
        all[collectionId].push({
            id: generateId(),
            name: name.trim(),
            filter: filter.trim(),
        });
        saveToStorage(all);
        return all;
    });
}

/**
 * Remove a saved filter from a collection.
 * @param {string} collectionId
 * @param {string} filterId
 */
export function removeSavedFilter(collectionId, filterId) {
    savedFilters.update((all) => {
        if (all[collectionId]) {
            all[collectionId] = all[collectionId].filter((f) => f.id !== filterId);
            if (all[collectionId].length === 0) {
                delete all[collectionId];
            }
        }
        saveToStorage(all);
        return all;
    });
}

/**
 * Update an existing saved filter.
 * @param {string} collectionId
 * @param {string} filterId
 * @param {string} name
 * @param {string} filter
 */
export function updateSavedFilter(collectionId, filterId, name, filter) {
    savedFilters.update((all) => {
        if (all[collectionId]) {
            const idx = all[collectionId].findIndex((f) => f.id === filterId);
            if (idx !== -1) {
                all[collectionId][idx] = {
                    id: filterId,
                    name: name.trim(),
                    filter: filter.trim(),
                };
            }
        }
        saveToStorage(all);
        return all;
    });
}
