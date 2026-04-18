import { toDeleteProp } from "@/base/fieldSettings";
import { collectionAuthOptionsTab } from "./collectionAuthOptionsTab";
import { collectionFieldsTab } from "./collectionFieldsTab";
import { collectionRulesTab } from "./collectionRulesTab";
import { collectionViewQueryTab } from "./collectionViewQueryTab";

window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openCollectionUpsert = function(collection = {}, modalSettings = {
    // base modal events
    onbeforeopen: null, // function(el) {},
    onafteropen: null, // function(el) {},
    onbeforeclose: null, // function(el) {},
    onafterclose: null, // function(el) {},
    // collection specific events
    onsave: null, // function(collection, isNew) {},
    ondelete: null, // function(collection) {},
    onduplicate: null, // function(collection) {},
    ontruncate: null, // function(collection) {},
}) {
    app.store.errors = null; // reset

    const modal = collectionUpsertModal(collection || {}, modalSettings || {});

    document.body.appendChild(modal);

    app.modals.open(modal);
};

window.app.collectionTypes = {
    "base": {
        "icon": "ri-folder-2-line",
        "tabs": {
            "Fields": collectionFieldsTab,
            "API rules": collectionRulesTab,
        },
    },
    "view": {
        "icon": "ri-table-line",
        "tabs": {
            "Query": collectionViewQueryTab,
            "API rules": collectionRulesTab,
        },
    },
    "auth": {
        "icon": "ri-group-line",
        "tabs": {
            "Fields": collectionFieldsTab,
            "API rules": collectionRulesTab,
            "Options": collectionAuthOptionsTab,
        },
    },
};

function collectionUpsertModal(rawCollection, modalSettings) {
    let modal;

    const uniqueId = "collection_upsert_" + app.utils.randomString();

    const data = store({
        isSaving: false,
        originalCollection: {},
        collection: {},
        selectedTab: "",
        get activeTab() {
            if (!app.collectionTypes[data.collection.type]?.tabs) {
                return data.selectedTab;
            }

            if (
                !data.selectedTab
                || !app.collectionTypes[data.collection.type].tabs?.[data.selectedTab]
            ) {
                return Object.keys(app.collectionTypes[data.collection.type].tabs)?.[0] || "";
            }

            return data.selectedTab;
        },
        get isNew() {
            return app.utils.isEmpty(data.originalCollection?.id);
        },
        get collectionTypeOptions() {
            return Object.keys(app.collectionTypes).map((type) => {
                return {
                    value: type,
                    label: app.utils.sentenize(type, false) + " collection",
                };
            });
        },
        get collectionHash() {
            Object.keys(data.collection).length;
            return JSON.stringify(data.collection);
        },
        get originalCollectionHash() {
            return JSON.stringify(data.originalCollection);
        },
        get hasChanges() {
            return data.originalCollectionHash != data.collectionHash;
        },
        get canSave() {
            return !data.isSaving && (data.isNew || data.hasChanges);
        },
    });

    async function initCollection(collection) {
        if (app.utils.isEmpty(collection)) {
            collection = JSON.parse(JSON.stringify(app.store.collectionScaffolds.base)) || {
                type: "base",
                fields: [],
            };

            // add commonly used timestamp fields
            collection.fields.push({
                type: "autodate",
                name: "created",
                onCreate: true,
            });
            collection.fields.push({
                type: "autodate",
                name: "updated",
                onCreate: true,
                onUpdate: true,
            });
        }

        data.originalCollection = JSON.parse(JSON.stringify(collection));
        data.collection = JSON.parse(JSON.stringify(collection));
    }

    async function confirmSave(close = true) {
        if (!data.canSave) {
            return;
        }

        data.isSaving = true;

        app.modals.openCollectionChangesConfirmation(
            data.originalCollection,
            data.collection,
            () => save(close),
            () => {
                data.isSaving = false;
            },
        );
    }

    function exportPayload() {
        const payload = JSON.parse(JSON.stringify(data.collection));
        payload.fields = payload.fields || [];

        // remove fields marked for deletion
        for (let i = payload.fields.length - 1; i >= 0; i--) {
            const field = payload.fields[i];

            if (field[toDeleteProp]) {
                payload.fields.splice(i, 1);
                continue;
            }
        }

        return payload;
    }

    async function save(close = true) {
        data.isSaving = true;

        try {
            const payload = exportPayload();

            const isNew = app.utils.isEmpty(data.originalCollection?.id);

            let request;
            if (isNew) {
                request = app.pb.collections.create(payload);
            } else {
                request = app.pb.collections.update(data.originalCollection.id, payload);
            }

            const rawCollection = JSON.stringify(await request);

            data.originalCollection = JSON.parse(rawCollection);
            data.collection = JSON.parse(rawCollection);
            app.store.addOrUpdateCollection(JSON.parse(rawCollection));

            modalSettings?.onsave?.(JSON.parse(rawCollection), isNew);

            data.isSaving = false;

            app.toasts.success(
                isNew
                    ? `Successfully created collection "${data.collection.name}".`
                    : `Successfully updated collection "${data.collection.name}".`,
                { key: "collectionSave" },
            );

            // reset all errors
            app.store.errors = null;

            if (close) {
                app.modals.close(modal, true);
            }
        } catch (err) {
            if (!err?.isAbort) {
                data.isSaving = false;
                app.checkApiError(err, false);
                app.toasts.error(err.message || "Failed to save collection.", { key: "collectionSave" });
            }
        }
    }

    function resetForm() {
        data.collection = JSON.parse(JSON.stringify(data.originalCollection));
    }

    async function duplicate() {
        const clone = data.originalCollection ? JSON.parse(JSON.stringify(data.originalCollection)) : {};
        clone.id = "";
        clone.system = false;
        clone.name = clone.name + "_duplicate";
        clone.created = "";
        clone.updated = "";

        // updated indexes ids
        clone.indexes = clone.indexes?.map((idx) => {
            return app.utils.replaceIndexFields(idx, (parsed) => {
                return {
                    indexName: parsed.indexName + app.utils.randomString(3),
                    tableName: clone.name,
                };
            });
        }) || [];

        await modalSettings.onduplicate?.(clone);

        return initCollection(clone);
    }

    async function changeTab(tabName) {
        data.selectedTab = tabName;

        // ui tick
        await new Promise((r) => setTimeout(r, 0));

        // refresh errors in case to retrigger validations
        if (app.store.errors) {
            app.store.errors = JSON.parse(JSON.stringify(app.store.errors));
        }
    }

    modal = t.div(
        {
            pbEvent: "collectionUpsertModal",
            "html-data-collectionId": () => data.originalCollection?.id,
            "html-data-collectionName": () => data.originalCollection?.name,
            className: "modal collection-upsert-modal",
            inert: () => data.isSaving,
            onkeydown: (e) => {
                if ((e.ctrlKey || e.metaKey) && e.code == "KeyS") {
                    e.preventDefault();
                    // temp blur any active input to make sure that onchange/blur events are fired
                    const input = document.activeElement;
                    input?.blur();

                    confirmSave(false);

                    // restore previous active input
                    input?.focus();
                }
            },
            onbeforeopen: () => {
                initCollection(rawCollection);

                return modalSettings.onbeforeopen?.(el);
            },
            onafteropen: (el) => {
                modalSettings.onafteropen?.(el);
            },
            onbeforeclose: (el, forceClosed) => {
                if (forceClosed) {
                    return modalSettings.onbeforeclose?.(el);
                }

                if (data.isSaving) {
                    return false;
                }

                if (!data.hasChanges) {
                    return modalSettings.onbeforeclose?.(el);
                }

                return new Promise((r) => {
                    app.modals.confirm(
                        "You have unsaved changes. Do you really want to discard them?",
                        () => r(modalSettings.onbeforeclose?.(el)),
                        () => r(false),
                    );
                });
            },
            onafterclose: (el) => {
                modalSettings.onafterclose?.(el);
                el?.remove();
            },
            onmount: (el) => {
                el._watchers?.forEach((w) => w?.unwatch());
                el._watchers = [
                    watch(
                        () => data.collection.type,
                        (newType, oldType) => {
                            if (!oldType || newType == oldType || !app.store.collectionScaffolds[newType]) {
                                return;
                            }

                            // reset fields list errors on type change
                            app.utils.deleteByPath(app.store.errors, "fields");

                            // merge with the scaffold to ensure that the minimal props are set
                            const scaffold = JSON.parse(JSON.stringify(app.store.collectionScaffolds[newType]));
                            data.collection = Object.assign(
                                structuredClone(scaffold),
                                JSON.parse(JSON.stringify(data.collection)),
                            );
                            data.originalCollection = scaffold;
                            syncFieldsAndIndexesWithScaffold(data.collection);
                        },
                    ),

                    // collection rename
                    watch(
                        () => data.collection.name,
                        (newName, oldName) => {
                            newName = app.utils.slugify(newName);
                            data.collection.name = newName;

                            if (typeof oldName == "undefined" || !newName || newName == oldName) {
                                return;
                            }

                            // update indexes with the latest collection name as table name
                            clearTimeout(el.__collectionRenameTimeoutId);
                            el.__collectionRenameTimeoutId = setTimeout(() => {
                                data.collection.indexes = data.collection.indexes?.map((idx) => {
                                    return app.utils.replaceIndexFields(idx, { tableName: data.collection.name });
                                });
                            }, 150);
                        },
                    ),
                ];
            },
            onunmount: (el) => {
                clearTimeout(el?.__collectionRenameTimeoutId);
                el?._watchers?.forEach((w) => w?.unwatch());
            },
        },
        t.header(
            { className: "modal-header isolated" },
            t.div(
                { className: "grid sm" },
                t.div(
                    { className: "col-12 flex" },
                    t.h6(
                        { className: "modal-title" },
                        t.span(null, () => (data.isNew ? "Create " : "Edit ")),
                        t.strong(
                            {
                                hidden: () => data.isNew,
                                className: "txt-ellipsis collection-name",
                            },
                            () => data.originalCollection?.name,
                        ),
                        t.span(null, " collection"),
                    ),
                    t.div({ className: "flex-fill" }),
                    () => {
                        if (app.utils.isEmpty(data.originalCollection?.id)) {
                            return;
                        }

                        return [
                            t.button(
                                {
                                    type: "button",
                                    className: "btn sm circle transparent",
                                    title: "More options",
                                    "html-popovertarget": uniqueId + "modal-header-dropdown",
                                },
                                t.i({ className: "ri-more-line", ariaHidden: true }),
                            ),
                            t.div(
                                {
                                    id: uniqueId + "modal-header-dropdown",
                                    className: "dropdown nowrap modal-header-dropdown",
                                    popover: "auto",
                                },
                                t.button(
                                    {
                                        type: "button",
                                        className: "dropdown-item",
                                        onclick: (e) => {
                                            e.target.closest(".dropdown").hidePopover();
                                            app.utils.copyToClipboard(
                                                JSON.stringify(data.originalCollection, null, 2),
                                            );
                                            app.toasts.success("Collection copied to clipboard!");
                                        },
                                    },
                                    t.i({ className: "ri-braces-line", ariaHidden: true }),
                                    t.span({ className: "txt" }, "Copy JSON"),
                                ),
                                t.button(
                                    {
                                        type: "button",
                                        className: "dropdown-item",
                                        onclick: (e) => {
                                            e.target.closest(".dropdown").hidePopover();

                                            if (data.hasChanges) {
                                                app.modals.confirm(
                                                    "You have unsaved changes. Do you really want to discard them?",
                                                    duplicate,
                                                    null,
                                                    { yesButton: "Yes, discard" },
                                                );
                                            } else {
                                                duplicate();
                                            }
                                        },
                                    },
                                    t.i({ className: "ri-file-copy-line", ariaHidden: true }),
                                    t.span({ className: "txt" }, "Duplicate"),
                                ),
                                t.hr(),
                                () => {
                                    if (data.collection.type == "view") {
                                        return; // view don't have their own records
                                    }
                                    return truncateDropdownItem(data, modalSettings);
                                },
                                () => {
                                    if (data.collection.system) {
                                        return; // system collections cannot be deleted
                                    }
                                    return deleteDropdownItem(data, modalSettings);
                                },
                            ),
                        ];
                    },
                ),
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "fields" },
                        t.div(
                            { className: "field" },
                            t.label({
                                htmlFor: uniqueId + "col_name",
                                textContent: () => {
                                    return `Name${data.collection?.system ? " (system)" : ""}`;
                                },
                            }),
                            t.input({
                                id: uniqueId + "col_name",
                                type: "text",
                                name: "name",
                                required: true,
                                spellcheck: false,
                                placeholder: "e.g. posts",
                                autofocus: () => data.isNew,
                                disabled: () => !data.isNew && data.collection?.system,
                                value: () => data.collection.name || "",
                                onmount: (el) => {
                                    el.addEventListener("compositionend", (e) => {
                                        data.collection.name = e.target.value;
                                    });
                                },
                                oninput: (e) => {
                                    if (e.isComposing) {
                                        return;
                                    }
                                    data.collection.name = e.target.value;
                                },
                            }),
                        ),
                        t.div(
                            { className: "field addon" },
                            t.button(
                                {
                                    type: "button",
                                    disabled: () => !data.isNew,
                                    className: () =>
                                        `btn sm collection-type-select ${data.isNew ? "outline" : "transparent"}`,
                                    "html-popovertarget": uniqueId + "col_type_dropdown",
                                },
                                t.span(
                                    { className: "txt" },
                                    "Type: ",
                                    () => app.utils.sentenize(data.collection.type, false) || "N/A",
                                ),
                                t.i({
                                    hidden: () => !data.isNew,
                                    ariaHidden: true,
                                    className: "ri-arrow-drop-down-line m-l-auto",
                                }),
                            ),
                            t.div(
                                {
                                    id: uniqueId + "col_type_dropdown",
                                    className: "dropdown nowrap collection-type-dropdown",
                                    popover: "auto",
                                },
                                () => {
                                    let options = [];

                                    for (const opt of data.collectionTypeOptions) {
                                        options.push(
                                            t.button(
                                                {
                                                    type: "button",
                                                    className: () =>
                                                        `dropdown-item ${
                                                            opt.value == data.collection.type ? "active" : ""
                                                        }`,
                                                    onclick: (e) => {
                                                        e.target.closest(".dropdown").hidePopover();
                                                        data.collection.type = opt.value;
                                                    },
                                                },
                                                t.i({
                                                    ariaHidden: true,
                                                    className: app.collectionTypes[opt.value]?.icon
                                                        || app.utils.fallbackCollectionIcon,
                                                }),
                                                t.span({ className: "txt" }, opt.label || opt.value),
                                            ),
                                        );
                                    }

                                    return options;
                                },
                            ),
                        ),
                    ),
                ),
                t.div(
                    { className: "col-12" },
                    t.nav(
                        { className: "tabs-header equal-width" },
                        () => {
                            const tabItems = [];

                            const tabs = app.collectionTypes[data.collection.type]?.tabs || {};
                            for (let tabName in tabs) {
                                tabItems.push(
                                    t.button(
                                        {
                                            type: "button",
                                            disabled: () => data.isSaving,
                                            className: () => `tab-item ${data.activeTab == tabName ? "active" : ""}`,
                                            onclick: () => changeTab(tabName),
                                        },
                                        t.span({ className: "txt" }, tabName),
                                    ),
                                );
                            }

                            return tabItems;
                        },
                    ),
                ),
            ),
        ),
        t.div(
            { className: "modal-content" },
            () => app.collectionTypes[data.collection.type]?.tabs?.[data.activeTab]?.(data),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    disabled: () => data.isSaving,
                    onclick: () => app.modals.close(modal),
                },
                t.span({ className: "txt" }, "Close"),
            ),
            () => {
                const rawErrors = JSON.stringify(app.store.errors);
                if (rawErrors == "" || rawErrors == "null" || rawErrors == "{}" || rawErrors == "[]") {
                    return;
                }

                return t.i({
                    className: "ri-alert-line txt-danger",
                    ariaDescription: app.attrs.tooltip(() => "Raw error:\n" + rawErrors),
                });
            },
            t.div(
                { className: "btns" },
                t.button(
                    {
                        type: "button",
                        className: () => `btn expanded-lg ${data.isSaving ? "loading" : ""}`,
                        disabled: () => !data.canSave,
                        onclick: () => confirmSave(true),
                    },
                    t.span({ className: "txt" }, () => (data.isNew ? "Create" : "Save changes")),
                ),
                t.button(
                    {
                        type: "button",
                        title: "Save options",
                        className: () => `btn p-5`,
                        disabled: () => !data.canSave,
                        "html-popovertarget": uniqueId + "save_options",
                    },
                    t.i({ className: "ri-arrow-up-s-line", ariaHidden: true }),
                ),
                t.div(
                    { id: uniqueId + "save_options", className: "dropdown nowrap", popover: "auto" },
                    t.button(
                        {
                            type: "button",
                            className: "dropdown-item",
                            onclick: (e) => {
                                e.target.closest(".dropdown").hidePopover();
                                confirmSave(false);
                            },
                        },
                        t.span({ className: "txt" }, "Save and continue"),
                        t.small({ className: "txt-hint" }, "(Ctrl+S)"),
                    ),
                    t.hr(),
                    t.button(
                        {
                            type: "button",
                            className: "dropdown-item",
                            onclick: (e) => {
                                e.target.closest(".dropdown").hidePopover();
                                resetForm();
                            },
                        },
                        t.span({ className: "txt" }, "Reset form"),
                    ),
                ),
            ),
        ),
    );

    return modal;
}

function syncFieldsAndIndexesWithScaffold(collection) {
    const newScaffold = JSON.parse(JSON.stringify(app.store.collectionScaffolds[collection.type]));

    // merge fields
    // -----------------------------------------------------------
    const oldFields = JSON.parse(JSON.stringify(collection.fields)) || [];
    const nonSystemFields = oldFields.filter((f) => !f.system);

    collection.fields = newScaffold.fields || [];

    for (const oldField of oldFields) {
        if (!oldField.system) {
            continue;
        }

        const field = collection.fields.find((f) => f.name == oldField.name);
        if (!field) {
            continue;
        }

        // merge the default field with the existing one
        Object.assign(field, oldField);
    }

    for (const field of nonSystemFields) {
        collection.fields.push(field);
    }

    // merge indexes
    // -----------------------------------------------------------
    collection.indexes = collection.indexes || [];

    if (collection.indexes.length) {
        const scaffoldIndexes = newScaffold?.indexes || [];

        indexesLoop: for (let i = collection.indexes.length - 1; i >= 0; i--) {
            const parsed = app.utils.parseIndex(collection.indexes[i]);
            const parsedName = parsed.indexName.toLowerCase();

            // remove old scaffold indexes
            for (const idx of scaffoldIndexes) {
                const oldScaffoldName = app.utils.parseIndex(idx).indexName.toLowerCase();
                if (parsedName == oldScaffoldName) {
                    collection.indexes.splice(i, 1);
                    continue indexesLoop;
                }
            }

            // remove indexes to nonexisting fields
            for (const column of parsed.columns) {
                const hasFieldWithName = !!collection.fields.find(
                    (f) => f.name.toLowerCase() == column.name.toLowerCase(),
                );
                if (!hasFieldWithName) {
                    collection.indexes.splice(i, 1);
                    continue indexesLoop;
                }
            }
        }
    }

    // merge new scaffold indexes
    app.utils.mergeUnique(collection.indexes, newScaffold.indexes);
}

function truncateDropdownItem(data, modalSettings) {
    const uniqueId = "truncate_" + app.utils.randomString();

    const local = store({
        isSubmitting: false,
        nameConfirm: "",
    });

    async function truncateCollection() {
        if (
            local.isSubmitting
            || !data.originalCollection?.name
            || data.originalCollection.name != local.nameConfirm
        ) {
            return false;
        }

        local.isSubmitting = true;

        try {
            await app.pb.collections.truncate(data.originalCollection.name);

            modalSettings.ontruncate?.(JSON.parse(JSON.stringify(data.originalCollection)));

            app.toasts.success(`Successfully truncated collection "${data.originalCollection.name}".`);

            local.isSubmitting = false;

            return true;
        } catch (err) {
            local.isSubmitting = false;
            app.checkApiError(err);
        }

        return false;
    }

    return t.button(
        {
            type: "button",
            className: "dropdown-item txt-danger",
            disabled: () => local.isSubmitting,
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();
                app.modals.confirm(
                    t.div(
                        null,
                        t.h6(
                            { className: "block txt-center" },
                            "Do you really want to delete all records of the collection?",
                        ),
                        t.div(
                            { className: "confirm-collection-label txt-bold m-t-sm m-b-sm" },
                            "Type the collection name ",
                            t.div(
                                { className: "label" },
                                () => data.originalCollection.name,
                                app.components.copyButton(() => data.originalCollection?.name),
                            ),
                            " to confirm:",
                        ),
                        t.div(
                            { className: "field" },
                            t.label({ htmlFor: uniqueId + ".confirm_name" }, "Collection name"),
                            t.input({
                                id: uniqueId + ".confirm_name",
                                type: "text",
                                required: true,
                                pattern: () => RegExp.escape(data.originalCollection.name),
                                value: () => local.nameConfirm,
                                oninput: (e) => local.nameConfirm = e.target.value,
                            }),
                        ),
                    ),
                    async () => {
                        document.getElementById(uniqueId + ".confirm_name")?.reportValidity();

                        const truncated = await truncateCollection();
                        if (!truncated) {
                            return false;
                        }

                        app.modals.close(e.target.closest(".modal"));
                    },
                    () => {
                        local.nameConfirm = "";
                    },
                );
            },
        },
        t.i({ className: "ri-eraser-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Truncate"),
    );
}

function deleteDropdownItem(data, modalSettings) {
    const uniqueId = "delete_" + app.utils.randomString();

    const local = store({
        isSubmitting: false,
        nameConfirm: "",
    });

    async function deleteCollection() {
        if (
            local.isSubmitting
            || !data.originalCollection?.name
            || data.originalCollection.name != local.nameConfirm
        ) {
            return false;
        }

        local.isSubmitting = true;

        try {
            await app.pb.collections.delete(data.originalCollection.name);

            modalSettings.ondelete?.(JSON.parse(JSON.stringify(data.originalCollection)));

            app.utils.removeByKey(app.store.collections, "id", data.originalCollection.id);

            app.toasts.success(`Successfully deleted collection "${data.originalCollection.name}".`);

            local.isSubmitting = false;

            return true;
        } catch (err) {
            local.isSubmitting = false;
            app.checkApiError(err);
        }

        return false;
    }

    return t.button(
        {
            type: "button",
            className: "dropdown-item txt-danger",
            disabled: () => local.isSubmitting,
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();

                const collectionModal = e.target.closest(".modal");

                app.modals.confirm(
                    t.div(
                        { className: "block" },
                        t.h6(
                            { className: "block txt-center" },
                            () => {
                                if (data.originalCollection.type == "view") {
                                    return "Do you really want to delete the selected collection?";
                                }

                                return "Do you really want to delete the selected collection and all its records";
                            },
                        ),
                        t.div(
                            { className: "confirm-collection-label txt-bold m-t-sm m-b-sm" },
                            "Type the collection name ",
                            t.div(
                                { className: "label" },
                                () => data.originalCollection.name,
                                app.components.copyButton(() => data.originalCollection?.name),
                            ),
                            " to confirm:",
                        ),
                        t.div(
                            { className: "field" },
                            t.label({ htmlFor: uniqueId + ".confirm_name" }, "Collection name"),
                            t.input({
                                id: uniqueId + ".confirm_name",
                                type: "text",
                                required: true,
                                pattern: () => RegExp.escape(data.originalCollection.name),
                                value: () => local.nameConfirm,
                                oninput: (e) => local.nameConfirm = e.target.value,
                            }),
                        ),
                    ),
                    async () => {
                        document.getElementById(uniqueId + ".confirm_name")?.reportValidity();

                        const deleted = await deleteCollection();
                        if (!deleted) {
                            return false;
                        }

                        app.modals.close(collectionModal);
                    },
                    () => {
                        local.nameConfirm = "";
                    },
                );
            },
        },
        t.i({ className: "ri-delete-bin-7-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Delete"),
    );
}
