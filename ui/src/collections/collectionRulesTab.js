export function collectionRulesTab(upsertData) {
    const local = store({
        showRulesInfo: false,
        showAuthRules: false,
    });

    const systemRuleTooltip = () =>
        app.attrs.tooltip(
            upsertData.originalCollection?.system ? "System collection rule cannot be changed." : null,
            "top-left",
        );

    function autocomplete(word) {
        return app.utils.collectionAutocompleteKeys(upsertData.collection, word);
    }

    return t.div(
        { className: "collection-tab-content collection-rules-tab-content" },
        t.div(
            { className: "grid" },
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "flex txt-hint txt-sm" },
                    t.span(
                        { className: "txt" },
                        "All rules follow the ",
                        t.a({
                            target: "_blank",
                            rel: "noopener noreferrer",
                            href: import.meta.env.PB_RULES_SYNTAX_DOCS,
                            textContent: "PocketBase filter syntax and operators",
                        }),
                        ".",
                    ),
                    t.strong({
                        tabIndex: -1,
                        className: "m-l-auto link-hint",
                        textContent: () => (local.showRulesInfo ? "Hide available fields" : "Show available fields"),
                        onclick: () => (local.showRulesInfo = !local.showRulesInfo),
                    }),
                ),
                app.components.slide(
                    () => local.showRulesInfo,
                    t.div(
                        { className: "alert warning m-t-sm" },
                        t.div(
                            { className: "content" },
                            t.p(null, "The following record fields are available:"),
                            t.div({ className: "flex flex-wrap gap-5" }, () => {
                                const identifiers = app.utils.getAllCollectionIdentifiers(upsertData.collection);
                                return identifiers.map((f) => {
                                    return t.code(null, f);
                                });
                            }),
                            t.hr({ className: "m-t-10 m-b-10" }),
                            t.p(
                                null,
                                "The request fields could be accessed with the special ",
                                t.strong(null, "@request"),
                                " fields:",
                            ),
                            t.div(
                                { className: "flex flex-wrap gap-5" },
                                t.code(null, "@request.headers.*"),
                                t.code(null, "@request.query.*"),
                                t.code(null, "@request.body.*"),
                                t.code(null, "@request.auth.*"),
                            ),
                            t.hr({ className: "m-t-10 m-b-10" }),
                            t.p(
                                null,
                                "You could also add constraints and query other collections using the ",
                                t.strong(null, "@collection"),
                                " field:",
                            ),
                            t.div(
                                { className: "flex flex-wrap gap-5" },
                                t.code(null, "@collection.ANY_COLLECTION_NAME.*"),
                            ),
                            t.hr({ className: "m-t-10 m-b-10" }),
                            t.p(null, "Example rule:"),
                            () => {
                                const dateField = upsertData.collection.fields?.find(
                                    (f) => f.type == "date" || f.type == "autodate",
                                );
                                if (dateField) {
                                    return t.code(
                                        null,
                                        `@request.auth.id != "" && ${dateField.name} > "2022-01-01 00:00:00.000Z"`,
                                    );
                                }
                                return t.code(null, `@request.auth.id != ""`);
                            },
                        ),
                    ),
                ),
            ),
            t.div(
                { className: "col-12", ariaDescription: systemRuleTooltip() },
                app.components.ruleField({
                    label: "List/Search rule",
                    name: "listRule",
                    autocomplete: autocomplete,
                    disabled: () => upsertData.originalCollection?.system,
                    value: () => upsertData.collection.listRule,
                    oninput: (val) => (upsertData.collection.listRule = val),
                }),
            ),
            t.div(
                { className: "col-12", ariaDescription: systemRuleTooltip() },
                app.components.ruleField({
                    label: "View rule",
                    name: "viewRule",
                    autocomplete: autocomplete,
                    disabled: () => upsertData.originalCollection?.system,
                    value: () => upsertData.collection.viewRule,
                    oninput: (val) => (upsertData.collection.viewRule = val),
                }),
            ),
            () => {
                // view collections has only List and View API rules
                if (upsertData.collection.type == "view") {
                    return;
                }

                return [
                    t.div(
                        { className: "col-12", ariaDescription: systemRuleTooltip() },
                        app.components.ruleField({
                            label: [
                                t.span({ className: "txt", textContent: "Create rule" }),
                                t.i({
                                    hidden: () => upsertData.collection.createRule == null,
                                    className: "ri-information-line link-hint",
                                    ariaDescription: app.attrs.tooltip(
                                        "The main record fields hold the values that are going to be inserted in the database.",
                                    ),
                                }),
                            ],
                            name: "createRule",
                            autocomplete: autocomplete,
                            disabled: () => upsertData.originalCollection?.system,
                            value: () => upsertData.collection.createRule,
                            oninput: (val) => (upsertData.collection.createRule = val),
                        }),
                    ),
                    t.div(
                        { className: "col-12", ariaDescription: systemRuleTooltip() },
                        app.components.ruleField({
                            label: [
                                t.span({ className: "txt", textContent: "Update rule" }),
                                t.i({
                                    hidden: () => upsertData.collection.updateRule == null,
                                    className: "ri-information-line link-hint",
                                    ariaDescription: app.attrs.tooltip(
                                        "The main record fields hold the old/existing record field values.\nTo target the newly submitted ones you can use @request.body.*.",
                                    ),
                                }),
                            ],
                            name: "updateRule",
                            autocomplete: autocomplete,
                            disabled: () => upsertData.originalCollection?.system,
                            value: () => upsertData.collection.updateRule,
                            oninput: (val) => (upsertData.collection.updateRule = val),
                        }),
                    ),
                    t.div(
                        { className: "col-12", ariaDescription: systemRuleTooltip() },
                        app.components.ruleField({
                            label: "Delete rule",
                            name: "deleteRule",
                            autocomplete: autocomplete,
                            disabled: () => upsertData.originalCollection?.system,
                            value: () => upsertData.collection.deleteRule,
                            oninput: (val) => (upsertData.collection.deleteRule = val),
                        }),
                    ),
                ];
            },
        ),
        // auth specific fields
        () => {
            if (upsertData.collection.type != "auth") {
                return;
            }

            return [
                t.hr({ className: "m-t-base m-b-base" }),
                t.button(
                    {
                        type: "button",
                        onmount: () => {
                            local.showAuthRules = upsertData.collection.manageRule !== null
                                || upsertData.collection.authRule !== "";
                        },
                        className: () => `btn secondary sm ${local.showAuthRules ? "" : "transparent"}`,
                        onclick: () => {
                            local.showAuthRules = !local.showAuthRules;
                        },
                    },
                    t.span({ className: "txt" }, "Additional auth collection rules"),
                    t.i({
                        className: () => (local.showAuthRules ? "ri-arrow-drop-up-line" : "ri-arrow-drop-down-line"),
                    }),
                ),
                app.components.slide(
                    () => local.showAuthRules,
                    t.div(
                        { className: "grid sm m-t-sm" },
                        t.div(
                            { className: "col-12", ariaDescription: systemRuleTooltip() },
                            app.components.ruleField({
                                label: "Authentication rule",
                                name: "authRule",
                                placeholder: "",
                                autocomplete: autocomplete,
                                disabled: () => upsertData.originalCollection?.system,
                                value: () => upsertData.collection.authRule,
                                oninput: (val) => (upsertData.collection.authRule = val),
                            }),
                            t.div(
                                { className: "field-help" },
                                t.p(
                                    null,
                                    "This rule is executed every time ",
                                    t.strong(null, "before authentication"),
                                    " allowing you to restrict who can authenticate.",
                                ),
                                t.p(
                                    null,
                                    "For example, to allow only verified users you can set it to ",
                                    t.code(null, "verified = true"),
                                    ".",
                                ),
                                t.p(null, "Leave it empty to allow anyone with an account to authenticate."),
                                t.p(
                                    null,
                                    `To disable authentication entirely you can change it to "Set superusers only".`,
                                ),
                            ),
                        ),
                        t.div(
                            { className: "col-12", ariaDescription: systemRuleTooltip() },
                            app.components.ruleField({
                                label: "Manage rule",
                                name: "manageRule",
                                autocomplete: autocomplete,
                                disabled: () => upsertData.originalCollection?.system,
                                value: () => upsertData.collection.manageRule,
                                oninput: (val) => (upsertData.collection.manageRule = val),
                            }),
                            t.div(
                                { className: "field-help" },
                                t.p(
                                    null,
                                    "This rule is executed in addition to the ",
                                    t.strong(null, "create"),
                                    " and ",
                                    t.strong(null, "update"),
                                    " API rules.",
                                ),
                                t.p(
                                    null,
                                    "It enables superuser-like permissions to allow fully managing the auth record(s), eg. changing the password without requiring to enter the old one, directly updating the verified state or email, etc.",
                                ),
                            ),
                        ),
                    ),
                ),
            ];
        },
    );
}
