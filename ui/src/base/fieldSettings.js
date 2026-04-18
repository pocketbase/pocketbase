window.app = window.app || {};
window.app.components = window.app.components || {};

export const toDeleteProp = "@toDelete";

window.app.components.fieldSettings = function(data, settingsArg = {}) {
    const uniqueId = "base_" + app.utils.randomString();

    const settings = store({
        // options
        showHidden: true,
        showPresentable: true,
        showDuplicate: true,
        showRemove: true, // for system fields this is ignored
        // slots
        header: (el) => null,
        content: (el) => null,
        footer: (el) => null,
    });

    const watchers = app.utils.extendStore(settings, settingsArg);

    function duplicateField() {
        const clone = JSON.parse(JSON.stringify(data.field));

        clone.id = "";
        clone.system = false;
        clone.name = getUniqueFieldName(data.collection.fields, clone.name + "_copy");
        clone.__detailsOpen = true;

        if (clone.primaryKey) {
            clone.primaryKey = false;
        }

        if (clone[toDeleteProp]) {
            delete clone[toDeleteProp];
        }

        data.collection.fields.splice(data.fieldIndex + 1, 0, clone);
    }

    function removeField() {
        if (data.field.id) {
            // existing fields are only marked as deleted
            data.field[toDeleteProp] = true;
        } else {
            // new fields are immediately removed
            data.collection.fields.splice(data.fieldIndex, 1);
        }
    }

    return t.details(
        {
            // duplicate the class as raw attribute because the reactive state
            // is applied AFTER mount (MutationObserver related quirk)
            "html-class": "accordion record-field-settings",
            className: () =>
                `accordion record-field-settings field-type-${data.field.type} ${
                    data.field[toDeleteProp] ? "deleted" : ""
                }`,
            name: "collection_field",
            onmount: (el) => {
                if (data.field.__detailsOpen) {
                    delete data.field.__detailsOpen;
                    el.open = true;
                }

                // name normalizer
                watchers.push(
                    watch(
                        () => data.field.name,
                        (newName, oldName) => {
                            newName = app.utils.slugify(newName);
                            data.field.name = newName;

                            if (typeof oldName == "undefined") {
                                return;
                            }

                            replaceIndexesColumn(data.collection, oldName, newName);
                            replaceIdentityFields(data.collection, oldName, newName);
                        },
                    ),
                );

                // reset the name if it was previously deleted
                watchers.push(
                    watch(
                        () => data.field[toDeleteProp],
                        (deleted) => {
                            if (deleted && data.originalField?.name && data.field.name != data.originalField.name) {
                                data.field.name = data.originalField.name;
                            }
                        },
                    ),
                );

                // disable presentable
                watchers.push(
                    watch(() => {
                        if (data.field.presentable && data.field.hidden) {
                            data.field.presentable = false;
                            app.toasts.info("The field cannot be presentable if hidden.");
                        }
                    }),
                );

                // special cases for some system fields
                watchers.push(
                    watch(() => {
                        if (
                            (data.field.name == "id"
                                || (data.collection.type == "auth"
                                    && ["password", "tokenKey"].includes(data.field.name)))
                            && data.originalField
                            && data.field.required != data.originalField.required
                        ) {
                            data.field.required = data.originalField.required;
                            app.toasts.info(`The option cannot be changed for field "${data.field.name}".`);
                        }
                    }),
                );

                // prevent hidden prop change for special fields
                watchers.push(
                    watch(() => {
                        if (
                            (data.field.name == "id"
                                || (data.collection.type == "auth"
                                    && ["password", "tokenKey", "email"].includes(data.field.name)))
                            && data.originalField
                            && data.field.hidden != data.originalField.hidden
                        ) {
                            data.field.hidden = data.originalField.hidden;
                            app.toasts.info(`The option cannot be changed for field "${data.field.name}".`);
                        }
                    }),
                );
            },
            onunmount: (el) => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.summary(
            { tabIndex: -1, onfocusout: () => false, onclick: () => false, onkeyup: () => false },
            t.span({ className: "sort-handle" }, t.i({ className: "ri-draggable" })),
            t.header(
                {
                    className: "header-fields",
                    inert: () => data.field[toDeleteProp],
                    onclick: (e) => {
                        e.stopPropagation();
                        e.preventDefault();
                    },
                },
                t.div(
                    { className: "fields" },
                    t.label(
                        {
                            htmlFor: uniqueId + ".name",
                            className: () => `field addon ${data.field.system ? "txt-disabled" : ""}`,
                        },
                        t.i({
                            className: app.fieldTypes[data.field.type]?.icon || app.utils.fallbackFieldIcon,
                            ariaDescription: app.attrs.tooltip(() => {
                                if (data.field.system) {
                                    return data.field.type + " (system)";
                                }
                                return data.field.type;
                            }),
                        }),
                    ),
                    t.div(
                        { className: "field prop-name" },
                        t.input({
                            type: "text",
                            id: uniqueId + ".name",
                            name: () => `fields.${data.fieldIndex}.name`,
                            required: true,
                            spellcheck: false,
                            placeholder: "Field name*",
                            className: "inline-error",
                            disabled: () => data.field[toDeleteProp] || data.field.system,
                            value: () => data.field.name || "",
                            oninput: (e) => {
                                if (e.isComposing) {
                                    return;
                                }
                                data.field.name = e.target.value;
                            },
                            onmount: (nameInput) => {
                                nameInput.addEventListener("compositionend", (e) => {
                                    data.field.name = e.target.value;
                                });

                                setTimeout(() => {
                                    if (nameInput && data.field.__focus) {
                                        nameInput.select();
                                        delete data.field.__focus;
                                    }
                                }, 0);
                            },
                        }),
                        t.div({ className: "field-labels" }, () => {
                            const labels = [];

                            if (data.field.required) {
                                labels.push(t.span({ className: "label success" }, "Required"));
                            }

                            if (data.field.hidden) {
                                labels.push(t.span({ className: "label danger" }, "Hidden"));
                            } else if (data.field.presentable) {
                                labels.push(t.span({ className: "label info" }, "Presentable"));
                            }

                            return labels;
                        }),
                    ),
                ),
                (el) => {
                    if (typeof settings.header == "function") {
                        return settings.header(el);
                    }
                    return settings.header;
                },
            ),
            t.button(
                {
                    type: "button",
                    className: () => {
                        const hasError = !app.utils.isEmpty(
                            app.utils.getByPath(app.store.errors, `fields.${data.fieldIndex}`),
                        );
                        return `btn sm circle transparent secondary ${hasError ? "txt-danger" : ""}`;
                    },
                    title: "Field options",
                    hidden: () => data.field[toDeleteProp],
                    onclick: (e) => {
                        const details = e.target.closest("details");
                        if (details) {
                            details.open = !details.open;
                        }
                    },
                },
                t.i({ className: "ri-settings-3-line", ariaHidden: true }),
            ),
            t.button(
                {
                    type: "button",
                    className: "btn sm circle transparent warning",
                    hidden: () => !data.field[toDeleteProp],
                    onclick: () => delete data.field[toDeleteProp],
                    ariaLabel: app.attrs.tooltip("Restore"),
                },
                t.i({ className: "ri-restart-line", ariaHidden: true }),
            ),
        ),
        (el) => {
            if (typeof settings.content == "function") {
                return settings.content(el);
            }
            return settings.content;
        },
        t.footer(
            { className: "record-field-settings-footer" },
            (el) => {
                if (typeof settings.footer == "function") {
                    return settings.footer(el);
                }
                return settings.footer;
            },
            () => {
                if (!settings.showPresentable) {
                    return;
                }

                return t.div(
                    { className: "field prop-presentable" },
                    t.input({
                        type: "checkbox",
                        id: uniqueId + ".presentable",
                        name: () => `fields.${data.fieldIndex}.presentable`,
                        className: "sm",
                        disabled: () => data.field.hidden,
                        checked: () => !!data.field.presentable,
                        onchange: (e) => (data.field.presentable = e.target.checked),
                    }),
                    t.label(
                        { htmlFor: uniqueId + ".presentable" },
                        t.span({ className: "txt" }, "Presentable"),
                        t.i({
                            className: "ri-information-line link-hint",
                            ariaDescription: app.attrs.tooltip(
                                () => {
                                    let msg =
                                        "Whether the field should be preferred in the Superuser UI relation listings (default to auto).";
                                    if (data.field.hidden) {
                                        msg += "\nThe field cannot be presentable if hidden.";
                                    }
                                    return msg;
                                },
                            ),
                        }),
                    ),
                );
            },
            () => {
                if (!settings.showHidden) {
                    return;
                }

                return t.div(
                    { className: "field prop-hidden" },
                    t.input({
                        type: "checkbox",
                        id: uniqueId + ".hidden",
                        className: "sm",
                        name: () => `fields.${data.fieldIndex}.hidden`,
                        checked: () => !!data.field.hidden,
                        onchange: (e) => (data.field.hidden = e.target.checked),
                    }),
                    t.label(
                        { htmlFor: uniqueId + ".hidden" },
                        t.span({ className: "txt" }, "Hidden"),
                        t.i({
                            className: "ri-information-line link-hint",
                            ariaDescription: app.attrs.tooltip("Hide from the JSON API response and filters."),
                        }),
                    ),
                );
            },
            t.button(
                {
                    hidden: () => !settings.showDuplicate && (!settings.showRemove || data.field.system),
                    type: "button",
                    title: "More options",
                    className: "btn sm circle transparent secondary more-btn m-l-auto",
                    "html-popovertarget": uniqueId + "_options_dropdown",
                },
                t.i({ className: "ri-more-line", ariaHidden: true }),
            ),
            t.div(
                {
                    id: uniqueId + "_options_dropdown",
                    className: "dropdown sm field-options-dropdown",
                    popover: "auto",
                },
                () => {
                    if (!settings.showDuplicate) {
                        return;
                    }

                    return t.button({
                        type: "button",
                        className: "dropdown-item",
                        role: "menuitem",
                        textContent: "Duplicate",
                        onclick: (e) => {
                            duplicateField();
                            e.target.closest(".dropdown").hidePopover();
                        },
                    });
                },
                () => {
                    if (!settings.showRemove || data.field.system) {
                        return;
                    }

                    return t.button({
                        type: "button",
                        className: "dropdown-item",
                        role: "menuitem",
                        textContent: "Remove",
                        onclick: (e) => {
                            removeField();
                            e.target.closest(".dropdown").hidePopover();
                            e.target.closest("details").open = false;
                        },
                    });
                },
            ),
        ),
    );
};

function getUniqueFieldName(allFields, name = "field") {
    let result = name;
    let counter = 2;

    let suffix = name.match(/\d+$/)?.[0] || ""; // extract numeric suffix

    // name without the suffix
    let base = suffix ? name.substring(0, name.length - suffix.length) : name;

    while (!!allFields?.find((field) => field.name.toLowerCase() == result.toLowerCase())) {
        result = base + ((suffix << 0) + counter);
        counter++;
    }

    return result;
}

function replaceIdentityFields(collection, oldName, newName) {
    if (
        !newName
        || typeof oldName == "undefined"
        || oldName === newName
        || !collection?.passwordAuth?.identityFields?.length
    ) {
        return;
    }

    let identityFields = collection.passwordAuth.identityFields;
    for (let i = 0; i < identityFields.length; i++) {
        if (identityFields[i] == oldName) {
            identityFields[i] = newName;
        }
    }
}

function replaceIndexesColumn(collection, oldName, newName) {
    if (
        !newName
        || typeof oldName == "undefined"
        || oldName === newName
        || !collection?.indexes?.length
        || !collection?.fields?.length
    ) {
        return;
    }

    // field with the old name exists so there is no need to rename index columns
    if (!!collection.fields.find((f) => !f[toDeleteProp] && f.name == oldName)) {
        return;
    }

    // update indexes on renamed fields
    collection.indexes = collection.indexes.map((idx) => app.utils.replaceIndexColumn(idx, oldName, newName));
}
