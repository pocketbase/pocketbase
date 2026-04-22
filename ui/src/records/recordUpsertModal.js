window.app = window.app || {};
window.app.modals = window.app.modals || {};

/**
 * Opens a record upsert modal.
 *
 * @example
 * ```js
 * // create
 * app.modals.openRecordUpsert(collection)
 *
 * // update
 * app.modals.openRecordUpsert(collection, record)
 * ```
 *
 * @param {Object} collection
 * @param {Object} [record]
 * @param {Object} [modalSettings]
 */
window.app.modals.openRecordUpsert = function(collection, record = null, modalSettings = {
    // base modal events
    onbeforeopen: null, // function(el) {},
    onafteropen: null, // function(el) {},
    onbeforeclose: null, // function(el) {},
    onafterclose: null, // function(el) {},
    // record specific events
    onsave: null, // function(record, isNew) {},
    ondelete: null, // function(record) {},
    onduplicate: null, // function(record) {},
    ontokensreset: null, // function(record) {},
    onpasswordresetsend: null, // function(record) {},
    onverificationsend: null, // function(record) {},
}) {
    app.store.errors = null; // reset

    const modal = recordUpsertModal(collection, record, modalSettings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);

    app.modals.open(modal);
};

const defaultRedactFields = ["expand"];

function redacted(record, redactFields = defaultRedactFields) {
    // create redacted clone only if necessery
    if (redactFields.find((f) => typeof record[f] !== "undefined")) {
        record = Object.assign({}, record);
        for (let f of redactFields) {
            delete record[f];
        }
    }

    return record;
}

function downloadJSON(record) {
    record = redacted(record);
    app.utils.downloadJSON(record, record.collectionName + "_" + record.id + ".json");
}

function copyJSON(record) {
    record = redacted(record);
    app.utils.copyToClipboard(JSON.stringify(record, null, 2));
    app.toasts.success("Record copied to clipboard!");
}

function serializeRecord(record) {
    if (!record) {
        return "";
    }

    return JSON.stringify(redacted(record));
}

const TAB_MAIN = "main";
const TAB_AUTH_PROVIDERS = "authProviders";

// @todo consider exporting the "tabs" with the final version
function recordUpsertModal(collection, rawRecord, modalSettings) {
    if (!collection?.id) {
        console.warn("[recordUpsertModal] missing required collection");
        return;
    }

    let modal;

    const uniqueId = "record_upsert_" + app.utils.randomString();

    const listingColumnsPreferences = app.utils.getLocalHistory(app.consts.COLUMNS_STORAGE_PREFIX + collection.id, {});

    const data = store({
        isLoading: true,
        isSaving: false,
        isLocked: false,
        originalRecord: {},
        record: {},
        initialDraft: null,
        activeTab: TAB_MAIN,
        get isNew() {
            return app.utils.isEmpty(data.originalRecord?.id);
        },
        get isAuthCollection() {
            return collection.type == "auth";
        },
        get isSuperusersCollection() {
            return collection.name == "_superusers";
        },
        get showTabs() {
            return !data.isNew && data.isAuthCollection && !data.isSuperusersCollection;
        },
        get excludedFields() {
            const result = ["id"];

            if (data.isAuthCollection) {
                result.push("email", "emailVisibility", "verified", "password", "tokenKey");
            }

            return result;
        },
        get initialDraftHash() {
            return serializeRecord(data.initialDraft);
        },
        get recordHash() {
            return serializeRecord(data.record);
        },
        get originalRecordHash() {
            return serializeRecord(data.originalRecord);
        },
        get hasChanges() {
            return data.originalRecordHash != data.recordHash;
        },
        get isFormDisabled() {
            return data.isLoading || data.isSaving || (!data.isNew && !data.hasChanges);
        },
    });

    // note: not a getter to avoid the microtask batching
    function draftKey() {
        return "draft_" + collection.id + "_" + (data.originalRecord?.id || "");
    }

    function getDraftHash() {
        return window.localStorage.getItem(draftKey()) || "";
    }

    function getDraft() {
        try {
            const raw = getDraftHash();
            if (raw) {
                return JSON.parse(raw);
            }
        } catch (err) {
            console.warn("getDraft failure:", err);
            deleteDraft();
        }

        return null;
    }

    function saveDraft(serializedJSON) {
        try {
            window.localStorage.setItem(draftKey(), serializedJSON);
        } catch (e) {
            // ignore local storage errors in case the serialized data
            // exceed the browser localStorage single value quota
            console.warn("saveDraft failure:", e);
            deleteDraft();
        }
    }

    function deleteDraft() {
        window.localStorage.removeItem(draftKey());
        data.initialDraft = null;
    }

    function restoreDraft() {
        if (!data.initialDraft) {
            return;
        }

        if (app.store.errors) {
            app.store.errors = null;
        }

        const draftClone = JSON.parse(JSON.stringify(data.initialDraft));

        deleteDraft();

        data.record = draftClone;
    }

    let draftWatcher;

    function initDraftWatcher() {
        data.initialDraft = getDraft();

        draftWatcher?.unwatch();
        draftWatcher = watch(() => data.recordHash, (newVal, oldVal) => {
            if (typeof oldVal == "undefined") {
                return;
            }

            if (data.hasChanges) {
                saveDraft(data.recordHash);
            } else {
                deleteDraft();
            }
        });
    }

    async function initRecord(rawRecord) {
        data.isLoading = true;

        draftWatcher?.unwatch();

        // normalize rawRecord (could be plain id string)
        rawRecord = app.utils.isObject(rawRecord) ? rawRecord : { id: rawRecord || "" };

        // new record
        if (!rawRecord.id) {
            data.originalRecord = JSON.parse(JSON.stringify(rawRecord));
            data.record = JSON.parse(JSON.stringify(rawRecord));

            data.isLoading = false;
            data.isLocked = false;
            initDraftWatcher();
            return;
        }

        try {
            data.isLocked = !!app.store.settings?.meta?.hideControls;

            // preload to minimize content jumps
            data.originalRecord = JSON.parse(JSON.stringify(rawRecord));
            data.record = JSON.parse(JSON.stringify(rawRecord));

            // fetch to ensure that the main record fields are up-to-date
            let record = await app.pb.collection(collection.name).getOne(rawRecord.id, {
                requestKey: "upsert_load_" + rawRecord.id,
            });

            // preload existing expands (if any)
            if (rawRecord.expand) {
                record.expand = JSON.parse(JSON.stringify(rawRecord.expand));
            }

            // extend, not overwrite, to prevent reseting the reference passed down to the inputs
            Object.assign(data.originalRecord, JSON.parse(JSON.stringify(record)));
            Object.assign(data.record, JSON.parse(JSON.stringify(record)));

            data.isLoading = false;
            initDraftWatcher();
        } catch (err) {
            if (!err?.isAbort) {
                app.checkApiError(err);
                data.isLoading = false;
                setTimeout(() => app.modals.close(modal, true), 0);
            }
        }
    }

    async function exportPayload() {
        const payload = {};

        // shallow copy of the record fields
        for (const prop in data.record) {
            // skip expand and internal dynamic enumerable props
            if (prop == "expand" || prop.startsWith("@@")) {
                continue;
            }

            let val = data.record[prop]?.__raw || data.record[prop];

            // normalize undefined values
            if (typeof val == "undefined") {
                val = null;
            }

            payload[prop] = val;
        }

        // apply fields save normalization funcs
        for (const field of collection.fields) {
            const saveHook = app.fieldTypes[field.type]?.onrecordsave;
            if (!saveHook) {
                continue;
            }

            await saveHook({
                collection: collection,
                originalRecord: data.originalRecord,
                record: data.record,
                field: field,
                payload: payload,
            });
        }

        return payload;
    }

    async function save(close = true) {
        if (data.isLocked || data.isSaving || (!data.isNew && !data.hasChanges)) {
            return;
        }

        data.isSaving = true;

        try {
            const payload = await exportPayload();

            const isNew = app.utils.isEmpty(data.originalRecord?.id);

            let record;
            if (isNew) {
                record = await app.pb.collection(collection.name).create(payload);
            } else {
                record = await app.pb.collection(collection.name).update(data.originalRecord.id, payload);
            }

            deleteDraft();

            if (isNew) {
                // replace to ensure the same keys order and to force inputs rerender
                data.originalRecord = structuredClone(record);
                data.record = structuredClone(record);
            } else {
                // extend, not overwrite, to prevent reseting the reference passed down to the inputs
                Object.assign(data.originalRecord, structuredClone(record));
                Object.assign(data.record, structuredClone(record));
            }

            modalSettings.onsave?.(record, isNew);

            // reset all errors
            app.store.errors = null;

            let msg;
            if (isNew) {
                msg = `Successfully created ${collection.name} "${record.id}".`;
            } else {
                msg = `Successfully updated ${collection.name} "${record.id}".`;
            }
            app.toasts.success(msg, { key: "recordSave" });

            data.isSaving = false;

            if (close) {
                app.modals.close(modal, true);
            }
        } catch (err) {
            if (!err?.isAbort) {
                data.isSaving = false;
                app.checkApiError(err, false);
                app.toasts.error(err.message || "Failed to save record.", { key: "recordSave" });
            }
        }
    }

    function resetForm() {
        deleteDraft();
        data.record = JSON.parse(JSON.stringify(data.originalRecord));
    }

    async function duplicate() {
        const clone = data.originalRecord ? JSON.parse(JSON.stringify(data.originalRecord)) : {};
        clone.id = "";

        // apply fields duplicate hook
        for (const field of collection.fields) {
            const duplicateHook = app.fieldTypes[field.type]?.onrecordduplicate;
            if (!duplicateHook) {
                continue;
            }

            await duplicateHook({
                collection: collection,
                field: field,
                originalRecord: data.originalRecord,
                clone: clone,
            });
        }

        deleteDraft();

        modalSettings.onduplicate?.(clone);

        initRecord(clone);
    }

    function mainTab() {
        return [
            t.div(
                { className: "modal-content" },
                t.form(
                    {
                        id: uniqueId + "form",
                        className: "grid",
                        inert: () => data.isLoading || data.isSaving,
                        onsubmit: (e) => {
                            e.preventDefault();
                            // save(); // don't allow to prevent accidental save on input enter
                        },
                        onmount: (el) => {
                            el._quickSaveHandler = (e) => {
                                if ((e.ctrlKey || e.metaKey) && e.code == "KeyS") {
                                    e.preventDefault();
                                    save(false);
                                }
                            };
                            window.addEventListener("keydown", el._quickSaveHandler);
                        },
                        onunmount: (el) => {
                            if (el?._quickSaveHandler) {
                                window.removeEventListener("keydown", el?._quickSaveHandler);
                            }
                        },
                    },
                    // draft alert
                    () => {
                        if (
                            data.isLoading
                            || data.hasChanges
                            || app.utils.isEmpty(data.initialDraft)
                            || data.initialDraftHash == data.recordHash
                        ) {
                            return;
                        }

                        return t.div(
                            { className: "col-12" },
                            t.div(
                                { className: "alert warning flex gap-sm" },
                                t.div({ className: "content" }, "The record has previous unsaved changes."),
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn sm outline",
                                        onclick: () => restoreDraft(),
                                    },
                                    t.span({ className: "txt" }, "Restore draft"),
                                ),
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn sm secondary transparent circle m-l-auto",
                                        ariaLabel: app.attrs.tooltip("Discard draft", "left"),
                                        onclick: () => {
                                            deleteDraft();
                                        },
                                    },
                                    t.i({ className: "ri-close-line", ariaHidden: true }),
                                ),
                            ),
                        );
                    },
                    // primary key
                    () => {
                        const pkField = collection.fields?.find((f) => f.primaryKey);

                        return t.div(
                            { className: "col-12" },
                            app.fieldTypes[pkField.type].input({
                                get collection() {
                                    return collection;
                                },
                                get originalRecord() {
                                    return data.originalRecord;
                                },
                                get record() {
                                    return data.record;
                                },
                                get field() {
                                    return pkField;
                                },
                            }),
                        );
                    },
                    // special collection fields
                    () => {
                        if (!data.isAuthCollection) {
                            return;
                        }

                        const result = [
                            t.div({ className: "col-12" }, authFieldEmail(collection, data)),
                            t.div({ className: "col-12" }, authFieldPassword(collection, data)),
                        ];

                        // superusers are always verified
                        if (!data.isSuperusersCollection) {
                            result.push(
                                t.div({ className: "col-12" }, authFieldVerified(collection, data)),
                            );
                        }

                        return result;
                    },
                    // regular collection fields
                    () => {
                        const rows = [];

                        const excludedFields = data.excludedFields;

                        for (const field of collection.fields) {
                            if (!app.fieldTypes[field.type]?.input || excludedFields.includes(field.name)) {
                                continue;
                            }

                            rows.push(
                                t.div(
                                    // blur if not hidden and not explicitly toggle-on
                                    {
                                        className: () =>
                                            `col-12 ${
                                                field.hidden && !listingColumnsPreferences[field.id]
                                                    ? "hidden-field-blur"
                                                    : ""
                                            }`,
                                    },
                                    () => {
                                        return app.fieldTypes[field.type].input({
                                            get collection() {
                                                return collection;
                                            },
                                            get originalRecord() {
                                                return data.originalRecord;
                                            },
                                            get record() {
                                                return data.record;
                                            },
                                            get field() {
                                                return field;
                                            },
                                        });
                                    },
                                ),
                            );
                        }

                        if (rows.length && data.isAuthCollection) {
                            rows.unshift(t.div({ className: "col-12" }, t.hr({ className: "m-0" })));
                        }

                        return rows;
                    },
                ),
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
                t.button(
                    {
                        hidden: () => !data.isLocked,
                        type: "button",
                        className: "btn outline",
                        disabled: () => data.isFormDisabled,
                        onclick: () => data.isLocked = false,
                    },
                    t.i({ className: "ri-lock-unlock-line", ariaHidden: true }),
                    t.span({ className: "txt" }, "Unlock to save"),
                ),
                t.div(
                    {
                        hidden: () => data.isLocked,
                        className: "btns",
                    },
                    t.button(
                        {
                            type: "button",
                            className: () => `btn expanded-lg ${data.isLoading || data.isSaving ? "loading" : ""}`,
                            disabled: () => data.isLocked || data.isFormDisabled,
                            onclick: () => save(),
                        },
                        t.span({ className: "txt" }, () => (data.isNew ? "Create" : "Save changes")),
                    ),
                    t.button(
                        {
                            type: "button",
                            className: () => `btn p-5`,
                            title: "Save options",
                            disabled: () => data.isLocked || data.isFormDisabled,
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
                                    save(false);
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
        ];
    }

    modal = t.div(
        {
            pbEvent: "recordUpsertModal",
            className: "modal record-upsert-modal",
            onbeforeopen: () => {
                initRecord(rawRecord);

                return modalSettings.onbeforeopen?.(el);
            },
            onafteropen: (el) => {
                modalSettings.onafteropen?.(el);
            },
            onbeforeclose: (el, forceClosed) => {
                if (forceClosed) {
                    return modalSettings.onbeforeclose?.(el);
                }

                if (data.isLoading || data.isSaving) {
                    return false;
                }

                if (!data.hasChanges) {
                    return modalSettings.onbeforeclose?.(el);
                }

                return new Promise((r) => {
                    app.modals.confirm(
                        "You have unsaved changes. Do you really want to discard them?",
                        () => {
                            deleteDraft();
                            return r(modalSettings.onbeforeclose?.(el));
                        },
                        () => r(false),
                        { yesButton: "Yes, discard" },
                    );
                });
            },
            onafterclose: (el) => {
                modalSettings.onafterclose?.(el);
                el?.remove();
            },
            onunmount: () => {
                draftWatcher?.unwatch();
            },
        },
        t.header(
            { className: () => `modal-header ${data.showTabs ? "isolated" : ""}` },
            t.div(
                { className: "grid" },
                t.div(
                    { className: "col-12 flex" },
                    t.h6({ className: "modal-title" }, () => {
                        if (data.isLoading) {
                            return t.span({ className: "loader sm" });
                        }

                        return [
                            t.span(null, () => (data.isNew ? "Create " : "Edit ")),
                            t.strong(
                                { className: "txt-ellipsis collection-name", style: "max-width: 220px" },
                                () => collection.name,
                            ),
                            t.span(null, " record"),
                        ];
                    }),
                    t.div({ className: "flex-fill" }),
                    () => {
                        if (app.utils.isEmpty(data.originalRecord?.id)) {
                            return;
                        }

                        return [
                            t.button(
                                {
                                    type: "button",
                                    className: "btn sm circle transparent",
                                    title: "More options",
                                    disabled: () => data.isLoading,
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
                                // auth only options
                                () => {
                                    if (!data.isAuthCollection) {
                                        return;
                                    }

                                    const options = [];

                                    if (
                                        !data.originalRecord.verified
                                        && data.originalRecord.email
                                        // superusers are always verified
                                        && !data.isSuperusersCollection
                                    ) {
                                        options.push(sendVerificationDropdownItem(collection, data, modalSettings));
                                    }

                                    if (data.originalRecord.email) {
                                        options.push(
                                            sendPasswordResetEmailDropdownItem(collection, data, modalSettings),
                                        );
                                    }

                                    options.push(impersonateDropdownItem(collection, data, modalSettings));
                                    options.push(resetTokenKeyDropdownItem(collection, data, modalSettings));
                                    options.push(t.hr());

                                    return options;
                                },
                                t.button(
                                    {
                                        type: "button",
                                        className: "dropdown-item",
                                        onclick: (e) => {
                                            e.target.closest(".dropdown").hidePopover();
                                            copyJSON(data.originalRecord);
                                        },
                                    },
                                    t.i({ className: "ri-braces-line", ariaHidden: true }),
                                    t.span({ className: "txt" }, "Copy JSON"),
                                ),
                                () => {
                                    if (collection.type == "view") {
                                        return;
                                    }

                                    return [
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
                                        deleteDropdownItem(collection, data, modalSettings),
                                    ];
                                },
                            ),
                        ];
                    },
                ),
                () => {
                    if (!data.showTabs) {
                        return;
                    }

                    return t.div(
                        { className: "col-12" },
                        t.nav(
                            { className: "tabs-header equal-width" },
                            t.button(
                                {
                                    type: "button",
                                    disabled: () => data.isLoading || data.isSaving,
                                    className: () =>
                                        `tab-item ${
                                            data.activeTab == TAB_MAIN ? "active" : data.hasChanges ? "txt-warning" : ""
                                        }`,
                                    ariaDescription: app.attrs.tooltip(() =>
                                        data.hasChanges && data.activeTab != TAB_MAIN ? "Has unsaved changes" : ""
                                    ),
                                    onclick: () => data.activeTab = TAB_MAIN,
                                },
                                t.span({ className: "txt" }, () => (data.isAuthCollection ? "Account" : "Main")),
                            ),
                            t.button(
                                {
                                    type: "button",
                                    disabled: () => data.isLoading || data.isSaving,
                                    className: () => `tab-item ${data.activeTab == TAB_AUTH_PROVIDERS ? "active" : ""}`,
                                    onclick: () => data.activeTab = TAB_AUTH_PROVIDERS,
                                },
                                t.span({ className: "txt" }, "Auth providers"),
                            ),
                        ),
                    );
                },
            ),
        ),
        () => {
            if (
                !data.isNew
                && !data.isSuperusersCollection
                && data.activeTab == TAB_AUTH_PROVIDERS
            ) {
                return authProvidersTab(collection, data);
            }

            return mainTab();
        },
    );

    return modal;
}

// dropdown options
// -------------------------------------------------------------------

function resetTokenKeyDropdownItem(collection, data, modalSettings) {
    const local = store({
        isSubmitting: false,
    });

    async function resetTokenKey() {
        if (local.isSubmitting || !data.record.id) {
            return;
        }

        local.isSubmitting = true;

        try {
            const payload = {};

            const tokenKeyField = collection.fields.find((f) => f.name == "tokenKey");
            if (tokenKeyField.autogeneratePattern) {
                // leave the server to create the random string
                payload["tokenKey:autogenerate"] = "";
            } else {
                // set it manually
                payload["tokenKey"] = app.utils.randomSecret(
                    tokenKeyField.max << 0 || Math.max(2 * tokenKeyField.min << 0, 50),
                );
            }

            const updatedRecord = await app.pb.collection(collection.name).update(data.record.id, payload);

            modalSettings.ontokensreset?.(updatedRecord);

            // refresh all autodate fields
            const fields = collection.fields?.filter((f) => f.type == "autodate") || [];
            for (const field of fields) {
                const val = updatedRecord[field.name];
                if (data.initialDraft) {
                    data.initialDraft[field.name] = val;
                }
                data.originalRecord[field.name] = val;
                data.record[field.name] = val;
            }

            app.toasts.success("Successfully reset all tokens for the selected record.");
        } catch (err) {
            app.checkApiError(err);
        }

        local.isSubmitting = false;
    }

    return t.button(
        {
            type: "button",
            className: "dropdown-item",
            disabled: () => local.isSubmitting,
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();
                app.modals.confirm(
                    "Do you really want to reset all issued tokens for the selected auth record?",
                    resetTokenKey,
                    null,
                    { yesButton: "Reset all tokens" },
                );
            },
        },
        t.i({ className: "ri-reset-left-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Reset issued tokens"),
    );
}

function sendPasswordResetEmailDropdownItem(collection, data, modalSettings) {
    const local = store({
        isSubmitting: false,
    });

    async function sendPasswordResetEmail() {
        if (local.isSubmitting || !data.originalRecord?.email) {
            return;
        }

        local.isSubmitting = true;

        try {
            await app.pb.collection(collection.name).requestPasswordReset(data.originalRecord.email);

            modalSettings.onpasswordresetsend?.(JSON.parse(JSON.stringify(data.originalRecord)));

            app.toasts.success(`Successfully sent password reset email to ${data.originalRecord.email}.`);
        } catch (err) {
            app.checkApiError(err);
        }

        local.isSubmitting = false;
    }

    return t.button(
        {
            type: "button",
            className: "dropdown-item",
            disabled: () => local.isSubmitting,
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();
                app.modals.confirm(
                    `Do you really want to send password reset email to ${data.originalRecord?.email}?`,
                    sendPasswordResetEmail,
                    null,
                    { yesButton: "Send" },
                );
            },
        },
        t.i({ className: "ri-mail-lock-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Send password reset email"),
    );
}

function sendVerificationDropdownItem(collection, data, modalSettings) {
    const local = store({
        isSubmitting: false,
    });

    async function sendVerificationEmail() {
        if (local.isSubmitting || !data.originalRecord?.email || data.originalRecord?.verified) {
            return;
        }

        local.isSubmitting = true;

        try {
            await app.pb.collection(collection.name).requestVerification(data.originalRecord.email);

            modalSettings.onverificationsend?.(JSON.parse(JSON.stringify(data.originalRecord)));

            app.toasts.success(`Successfully sent verification email to ${data.originalRecord.email}.`);
        } catch (err) {
            app.checkApiError(err);
        }

        local.isSubmitting = false;
    }

    return t.button(
        {
            type: "button",
            className: "dropdown-item",
            disabled: () => local.isSubmitting,
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();
                app.modals.confirm(
                    `Do you really want to send verification email to ${data.originalRecord?.email}?`,
                    sendVerificationEmail,
                    null,
                    { yesButton: "Send" },
                );
            },
        },
        t.i({ className: "ri-mail-check-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Send verification email"),
    );
}

function impersonateDropdownItem(collection, data, modalSettings) {
    return t.button(
        {
            type: "button",
            className: "dropdown-item",
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();
                app.modals.openRecordImpersontate(data.originalRecord);
            },
        },
        t.i({ className: "ri-id-card-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Impersonate"),
    );
}

function deleteDropdownItem(collection, data, modalSettings) {
    const local = store({
        isSubmitting: false,
    });

    async function deleteRecord() {
        if (local.isSubmitting || !data.originalRecord?.id) {
            return;
        }

        local.isSubmitting = true;

        try {
            await app.pb.collection(collection.name).delete(data.originalRecord.id);

            modalSettings.ondelete?.(JSON.parse(JSON.stringify(data.originalRecord)));

            app.toasts.success(`Successfully deleted record "${data.originalRecord.id}".`);
        } catch (err) {
            app.checkApiError(err);
        }

        local.isSubmitting = false;
    }

    return t.button(
        {
            type: "button",
            className: "dropdown-item txt-danger",
            disabled: () => local.isSubmitting,
            onclick: (e) => {
                e.target.closest(".dropdown").hidePopover();
                app.modals.confirm(
                    `Do you really want to delete the selected record?`,
                    async () => {
                        await deleteRecord();
                        app.modals.close(e.target.closest(".modal"));
                    },
                    null,
                    { yesButton: "Delete record" },
                );
            },
        },
        t.i({ className: "ri-delete-bin-7-line", ariaHidden: true }),
        t.span({ className: "txt" }, "Delete"),
    );
}

// auth specific fields
// -------------------------------------------------------------------

function authFieldEmail(collection, data) {
    const emailField = collection.fields.find((f) => f.name == "email");
    if (!emailField) {
        console.warn("missing expected email field");
        return;
    }

    const uniqueId = "auth_email_" + app.utils.randomString();

    return t.div(
        { className: "record-field-input field-type-email field-type-auth-email" },
        t.div(
            { className: "fields" },
            t.div(
                { className: "field" },
                t.label(
                    { htmlFor: uniqueId },
                    t.i({ className: app.fieldTypes.email.icon, ariaHidden: true }),
                    t.span({ className: "txt" }, () => emailField.name),
                ),
                t.input({
                    type: "email",
                    id: uniqueId,
                    spellcheck: false,
                    name: () => emailField.name,
                    required: () => emailField.required,
                    value: () => data.record[emailField.name] || "",
                    oninput: (e) => (data.record[emailField.name] = e.target.value),
                }),
            ),
            t.div(
                { className: "field addon" },
                t.button(
                    {
                        type: "button",
                        className: () => `btn sm transparent ${data.record.emailVisibility ? "success" : "secondary"}`,
                        ariaDescription: app.attrs.tooltip("Make email public or private", "top-right"),
                        onclick: () => {
                            data.record.emailVisibility = !data.record.emailVisibility;
                        },
                    },
                    t.span({ className: "txt" }, "Public: ", () => (data.record.emailVisibility ? "On" : "Off")),
                ),
            ),
        ),
        () => {
            if (emailField.help) {
                return t.div({ className: "field-help" }, emailField.help);
            }
        },
    );
}

function authFieldVerified(collection, data) {
    const verifiedField = collection.fields.find((f) => f.name == "verified");
    if (!verifiedField) {
        console.warn("missing expected verified field");
        return;
    }

    const elem = app.fieldTypes.bool.input({
        get field() {
            return verifiedField;
        },
        get collection() {
            return collection;
        },
        get record() {
            return data.record;
        },
        get originalRecord() {
            return data.originalRecord;
        },
    });

    elem.addEventListener("change", (e) => {
        if (data.originalRecord.verified == data.record.verified) {
            return;
        }

        app.modals.confirm(
            `Do you really want to manually change the verified account state from "${!data.record
                .verified}" to "${data.record.verified}"?`,
            null,
            () => {
                data.record.verified = !data.record.verified;
            },
            { yesButton: "Yes, " + (data.record.verified ? "verify" : "unverify") },
        );
    });

    return elem;
}

function authFieldPassword(collection, data) {
    const uniqueId = "auth_pass_" + app.utils.randomString();

    const local = store({
        changePassword: false,
        get isNew() {
            return app.utils.isEmpty(data.originalRecord?.id);
        },
    });

    function clearPasswords() {
        delete data.record.password;
        delete data.record.passwordConfirm;

        if (app.store.errors) {
            delete app.store.errors.password;
            delete app.store.errors.passwordConfirm;
        }
    }

    return t.div(
        {
            className: "record-field-input field-type-password field-type-auth-password",
            onmount: (el) => {
                el._watchers?.forEach((w) => w?.unwatch());
                el._watchers = [
                    // force the toggle if any of the fields are populated
                    // (e.g. record update from outside) or there is an error
                    watch(() => {
                        if (local.changePassword) {
                            return; // already enabled
                        }

                        if (
                            app.store.errors?.password
                            || app.store.errors?.passwordConfirm
                            || data.record.password
                            || data.record.passwordConfirm
                        ) {
                            local.changePassword = true;
                        }
                    }),
                ];
            },
            onunmount: (el) => {
                el._watchers?.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            {
                hidden: () => local.isNew,
                className: "field",
            },
            t.input({
                type: "checkbox",
                id: uniqueId + "_change",
                className: "switch",
                checked: () => local.changePassword,
                onchange: (e) => {
                    local.changePassword = e.target.checked;
                    if (!e.target.checked) {
                        clearPasswords();
                    }
                },
            }),
            t.label({ htmlFor: uniqueId + "_change" }, t.span({ className: "txt" }, "change password")),
        ),
        app.components.slide(
            () => local.isNew || local.changePassword,
            t.div(
                { className: () => `fields ${local.isNew ? "" : "m-t-sm"}` },
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: uniqueId + "_password" },
                        t.i({ className: "ri-lock-line", ariaHidden: true }),
                        t.span({ className: "txt" }, "Password"),
                    ),
                    t.input({
                        type: "password",
                        id: uniqueId + "_password",
                        spellcheck: false,
                        name: "password",
                        className: "inline-error",
                        autocomplete: "new-password",
                        required: () => local.isNew || local.changePassword,
                        value: () => data.record.password || "",
                        oninput: (e) => {
                            if (!e.target.value) {
                                // delete to ensure that it is not submitted
                                delete data.record.password;
                            } else {
                                data.record.password = e.target.value;
                            }
                        },
                    }),
                ),
                t.div({ className: "delimiter" }),
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: uniqueId + "_password_confirm" },
                        t.i({ className: "ri-lock-line", ariaHidden: true }),
                        t.span({ className: "txt" }, "Confirm"),
                    ),
                    t.input({
                        type: "password",
                        id: uniqueId + "_password_confirm",
                        spellcheck: false,
                        name: "passwordConfirm",
                        className: "inline-error",
                        autocomplete: "new-password",
                        required: () => local.isNew || local.changePassword,
                        value: () => data.record.passwordConfirm || "",
                        oninput: (e) => {
                            if (!e.target.value) {
                                // delete to ensure that it is not submitted
                                delete data.record.passwordConfirm;
                            } else {
                                data.record.passwordConfirm = e.target.value;
                            }
                        },
                    }),
                ),
            ),
            () => {
                const helpText = collection.fields?.find((f) => f.name == "password")?.help || "";
                if (!helpText) {
                    return;
                }

                return t.div({ className: "field-help" }, helpText);
            },
            t.div(
                { className: "field-help" },
                t.span(
                    {
                        className: "txt link-hint",
                        role: "button",
                        onclick: (e) => {
                            e.preventDefault();
                            const random = app.utils.randomSecret(20);
                            data.record.password = random;
                            data.record.passwordConfirm = random;
                            app.utils.copyToClipboard(random);
                            app.toasts.info("Generated and copied random password to clipboard.");
                        },
                    },
                    "Generate and set random password",
                ),
            ),
        ),
    );
}

function authProvidersTab(collection, data) {
    const local = store({
        isLoading: false,
        externalAuths: [],
    });

    async function loadExternalAuths() {
        local.isLoading = true;

        try {
            local.externalAuths = await app.pb.collection("_externalAuths").getFullList({
                filter: app.pb.filter("collectionRef={:collectionId} && recordRef={:recordId}", {
                    collectionId: data.record.collectionId,
                    recordId: data.record.id,
                }),
            });

            local.isLoading = false;
        } catch (err) {
            if (err?.isAbort) {
                app.pb.checkApiError(err);
                local.isLoading = false;
            }
        }
    }

    function confirmAndUnlink(externalAuth) {
        const providerInfo = app.store.oauth2Providers?.find((p) => p.name == externalAuth.provider) || {};
        const name = providerInfo.displayName || externalAuth.provider;

        app.modals.confirm(
            `Do you really want to unlink the ${name} provider?`,
            () => {
                return app.pb
                    .collection("_externalAuths")
                    .delete(externalAuth.id)
                    .then(() => {
                        app.toasts.success(`Successfully unlinked ${name}.`);
                        loadExternalAuths(); // reload list
                    })
                    .catch((err) => {
                        app.checkApiError(err);
                    });
            },
            null,
            { yesButton: "Unlink" },
        );
    }

    return [
        t.div(
            { className: "modal-content" },
            t.div(
                {
                    className: "list",
                    onmount: () => {
                        loadExternalAuths();
                    },
                },
                () => {
                    if (local.isLoading) {
                        return t.div({ className: "list-item" }, t.div({ className: "skeleton-loader" }));
                    }

                    if (!local.externalAuths.length) {
                        return t.div(
                            { className: "list-item" },
                            t.div({ className: "block txt-hint txt-center" }, "No external auth providers found."),
                        );
                    }

                    return local.externalAuths.map((externalAuth) => {
                        const providerInfo = app.store.oauth2Providers?.find((p) => p.name == externalAuth.provider)
                            || {};

                        return t.div(
                            { className: "list-item" },
                            t.figure(
                                { className: "provider-logo" },
                                () => {
                                    if (providerInfo.logo) {
                                        return t.img({
                                            src: "data:image/svg+xml;base64," + btoa(providerInfo.logo),
                                            alt: externalAuth.provider + " logo",
                                        });
                                    }

                                    return t.i({ className: app.utils.fallbackProviderIcon, ariaHidden: true });
                                },
                            ),
                            t.div(
                                { className: "content" },
                                t.span(
                                    { className: "txt-nowrap" },
                                    () => providerInfo.displayName || externalAuth.provider,
                                ),
                                t.small({ className: "txt-hint" }, "ID: ", () => externalAuth.providerId),
                            ),
                            t.div(
                                { className: "actions" },
                                t.button(
                                    {
                                        type: "button",
                                        className: "btn sm secondary transparent circle",
                                        ariaLabel: app.attrs.tooltip("Unlink", "left"),
                                        onclick: () => confirmAndUnlink(externalAuth),
                                    },
                                    t.i({ className: "ri-close-line", ariaHidden: true }),
                                ),
                            ),
                        );
                    });
                },
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => app.modals.close(),
                },
                t.span({ className: "txt" }, "Close"),
            ),
        ),
    ];
}
