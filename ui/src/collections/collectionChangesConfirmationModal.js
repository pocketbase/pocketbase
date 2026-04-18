import { toDeleteProp } from "@/base/fieldSettings";

window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openCollectionChangesConfirmation = async function(
    oldCollection,
    newCollection,
    yesCallback,
    noCallback,
) {
    const data = store({
        isLoadingConflictingOIDCProviders: false,
        conflictingOIDCProviders: [],
        // ---
        get isCollectionRenamed() {
            return oldCollection?.name != newCollection?.name;
        },
        get isNewCollectionAuth() {
            return newCollection?.type === "auth";
        },
        get isNewCollectionView() {
            return newCollection?.type === "view";
        },
        get renamedFields() {
            if (data.isNewCollectionView) {
                return [];
            }

            return newCollection?.fields?.filter?.((f) => {
                let oldField;
                if (f.id && !f[toDeleteProp]) {
                    oldField = oldCollection.fields?.find?.((old) => old.id == f.id);
                }

                return oldField && oldField.name != f.name;
            }) || [];
        },
        get deletedFields() {
            if (data.isNewCollectionView) {
                return [];
            }

            return newCollection?.fields?.filter?.((f) => {
                return f.id && f[toDeleteProp];
            }) || [];
        },
        get multipleToSingleFields() {
            if (data.isNewCollectionView) {
                return [];
            }

            return newCollection?.fields?.filter?.((newField) => {
                const oldField = oldCollection?.fields?.find?.((f) => f.id == newField.id);
                if (!oldField || typeof oldField.maxSelect == "undefined") {
                    return false;
                }

                // normalize
                const oldMaxSelect = oldField.maxSelect || 1;
                const newMaxSelect = newField.maxSelect || 1;

                return oldMaxSelect > 1 && newMaxSelect == 1;
            }) || [];
        },
        get changedRules() {
            // for now enable only for "production"
            if (window.location.protocol != "https:") {
                return [];
            }

            const result = [];

            const ruleProps = ["listRule", "viewRule"];
            if (!data.isNewCollectionView) {
                ruleProps.push("createRule", "updateRule", "deleteRule");
            }
            if (data.isNewCollectionAuth) {
                ruleProps.push("manageRule", "authRule");
            }

            let oldRule, newRule;
            for (let prop of ruleProps) {
                oldRule = oldCollection?.[prop];
                newRule = newCollection?.[prop];
                if (oldRule === newRule) {
                    continue;
                }

                result.push({ prop, oldRule, newRule });
            }

            return result;
        },
        get needConfirmation() {
            return !app.utils.isEmpty(oldCollection?.id) && (
                data.isCollectionRenamed
                || data.renamedFields.length
                || data.deletedFields.length
                || data.multipleToSingleFields.length
                || data.changedRules.length
                || data.conflictingOIDCProviders.length
            );
        },
    });

    const knownOIDCProviders = ["oidc", "oidc2", "oidc3"];

    async function detectConflictingOIDCProviders() {
        if (app.utils.isEmpty(oldCollection?.id) || !data.isNewCollectionAuth) {
            return;
        }

        data.isLoadingConflictingOIDCProviders = true;

        try {
            data.conflictingOIDCProviders = [];

            for (const name of knownOIDCProviders) {
                const oldProvider = oldCollection?.oauth2?.providers?.find?.((p) => p.name == name);
                const newProvider = newCollection?.oauth2?.providers?.find?.((p) => p.name == name);

                if (!oldProvider || !newProvider) {
                    continue;
                }

                const oldHost = new URL(oldProvider.authURL).host;
                const newHost = new URL(newProvider.authURL).host;
                if (oldHost == newHost) {
                    continue;
                }

                // check if there are existing externalAuths
                const haveExternalAuths = await app.pb.collection("_externalAuths").getFirstListItem(
                    app.pb.filter("collectionRef={:collectionId} && provider={:provider}", {
                        collectionId: newCollection?.id,
                        provider: name,
                    }),
                    {
                        requestKey: null,
                    },
                );
                if (haveExternalAuths) {
                    data.conflictingOIDCProviders.push({ name, oldHost, newHost });
                }
            }

            data.isLoadingConflictingOIDCProviders = false;
        } catch (err) {
            if (err.isAbort) {
                data.isLoadingConflictingOIDCProviders = false;
                app.checkApiError(err);
            }
        }
    }

    await detectConflictingOIDCProviders();

    if (!data.needConfirmation) {
        return yesCallback();
    }

    app.modals.confirm(
        t.div(
            { className: "dangerous-collection-changes-list" },
            t.h5({ className: "block txt-center m-b-base" }, "Do you really want to save the collection changes?"),
            // general collection warning
            () => {
                if (!data.isCollectionRenamed && !data.deletedFields.length && !data.renamedFields.length) {
                    return;
                }

                return t.div(
                    { className: "alert warning m-b-base" },
                    t.p(
                        null,
                        "If the collection participate in another collection rule, filter or view query, you'll have to update it manually!",
                    ),
                    () => {
                        if (data.deletedFields.length) {
                            return t.p(
                                null,
                                "All data associated with the removed fields will be permanently deleted!",
                            );
                        }
                    },
                );
            },
            // renamed collection
            () => {
                if (!data.isCollectionRenamed) {
                    return;
                }

                return t.ul(
                    { className: "collection-changes-list changes-renamed-collection" },
                    t.li(
                        { className: "list-item" },
                        "Renamed collection ",
                        t.strong({ className: "label warning" }, oldCollection?.name),
                        t.i({ className: "ri-arrow-right-line txt-sm", ariaHidden: true }),
                        t.strong({ className: "label success" }, newCollection?.name || "N/A"),
                    ),
                );
            },
            // renamed fields
            () => {
                if (!data.renamedFields.length) {
                    return;
                }

                return t.ul(
                    { className: "collection-changes-list changes-renamed-fields" },
                    () => {
                        return data.renamedFields.map((newField) => {
                            const oldField = oldCollection?.fields?.find?.((f) => f.id == newField.id);
                            return t.li(
                                { className: "list-item" },
                                "Renamed field ",
                                t.strong({ className: "label warning" }, oldField?.name),
                                t.i({ className: "ri-arrow-right-line txt-sm", ariaHidden: true }),
                                t.strong({ className: "label success" }, newField.name || "N/A"),
                            );
                        });
                    },
                );
            },
            // deleted fields
            () => {
                if (!data.deletedFields.length) {
                    return;
                }

                return t.ul(
                    { className: "collection-changes-list changes-deleted-fields" },
                    () => {
                        return data.deletedFields.map((field) => {
                            return t.li(
                                { className: "list-item" },
                                "Deleted field ",
                                t.strong({ className: "label danger" }, field.name || "N/A"),
                            );
                        });
                    },
                );
            },
            // multiple->single fields
            () => {
                if (!data.multipleToSingleFields.length) {
                    return;
                }

                return t.ul(
                    { className: "collection-changes-list changes-multiple-to-single-fields" },
                    () => {
                        return data.multipleToSingleFields.map((field) => {
                            return t.li(
                                { className: "list-item" },
                                "Multiple to single value conversion of field ",
                                t.strong({ className: "label warning" }, field.name || field.id),
                                t.em({ className: "txt-sm" }, " (will keep only the last array item)"),
                            );
                        });
                    },
                );
            },
            // API rule changes
            () => {
                if (!data.changedRules.length) {
                    return;
                }

                return t.ul(
                    { className: "collection-changes-list changes-api-rules" },
                    () => {
                        return data.changedRules.map((ruleChange) => {
                            return t.li(
                                { className: "list-item" },
                                t.div(
                                    { className: "content" },
                                    t.span({ className: "txt" }, "Changed API rule for "),
                                    t.code(null, ruleChange.prop),
                                ),
                                t.small({ className: "txt-bold" }, "Old:"),
                                t.div(
                                    { className: "rule-content old-rule" },
                                    ruleChange.oldRule === null
                                        ? "null (superusers only)"
                                        : (ruleChange.oldRule || "\"\""),
                                ),
                                t.small({ className: "txt-bold" }, "New:"),
                                t.div(
                                    { className: "rule-content new-rule" },
                                    ruleChange.newRule === null
                                        ? "null (superusers only)"
                                        : (ruleChange.newRule || "\"\""),
                                ),
                            );
                        });
                    },
                );
            },
            // Conflicting OIDC changes
            () => {
                if (!data.conflictingOIDCProviders.length) {
                    return;
                }

                return t.ul(
                    { className: "collection-changes-list changes-api-rules" },
                    () => {
                        return data.conflictingOIDCProviders.map((oidc) => {
                            return t.li(
                                { className: "list-item" },
                                "Changed OIDC ",
                                oidc.name,
                                " host ",
                                t.strong({ className: "label warning" }, oidc.oldHost),
                                t.i({ className: "ri-arrow-right-line txt-sm", ariaHidden: true }),
                                t.strong({ className: "label success" }, oidc.newHost),
                                t.br(),
                                t.span(
                                    { className: "txt-hint" },
                                    "If the old and new OIDC configuration is not for the same provider consider deleting",
                                    " all old _externalAuths records associated to the current collection and provider,",
                                    " otherwise it may result in account linking errors.",
                                ),
                                " ",
                                t.a({
                                    rel: "noopenener noreferrer",
                                    target: "_blank",
                                    href: () => {
                                        return `#/collections?collection=_externalAuths&filter=collectionRef%3D%22${newCollection?.id}%22+%26%26+provider%3D%22${oidc.name}%22`;
                                    },
                                    textContent: "Review existing _externalAuths records",
                                }),
                            );
                        });
                    },
                );
            },
        ),
        yesCallback,
        noCallback,
        { className: "collection-changes-confirm-modal" },
    );
};
