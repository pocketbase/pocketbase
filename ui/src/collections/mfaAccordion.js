export function mfaAccordion(collection) {
    const uniqueId = "mfa_" + app.utils.randomString();

    const data = store({
        get config() {
            if (!collection.mfa) {
                collection.mfa = {
                    enabled: false,
                    duration: 900,
                    rule: "",
                };
            }

            return collection.mfa;
        },
        get isSuperusers() {
            return collection.system && collection.name == "_superusers";
        },
    });

    return t.details(
        {
            pbEvent: "mfaAccordion",
            name: "auth-methods",
            className: "accordion mfa-accordion",
        },
        t.summary(
            null,
            t.i({ className: "ri-shield-check-line", ariaHidden: true }),
            t.span({ className: "txt", textContent: "Multi-factor authentication (MFA)" }),
            t.span({
                className: () => `label m-l-auto ${data.config.enabled ? "success" : ""}`,
                textContent: () => (data.config.enabled ? "Enabled" : "Disabled"),
            }),
            () => {
                if (!app.store.errors?.mfa) {
                    return;
                }

                return t.i({
                    className: "ri-error-warning-fill txt-danger",
                    ariaDescription: app.attrs.tooltip("Has errors", "left"),
                });
            },
        ),
        t.div(
            { className: "grid sm" },
            t.div(
                { className: "col-sm-12" },
                t.div(
                    { className: "alert info" },
                    t.div(
                        { className: "content" },
                        t.p(
                            null,
                            "Multi-factor authentication (MFA) requires the user to authenticate with any 2 different auth methods (otp, identity/password, oauth2) before issuing an auth token. ",
                            t.a({
                                href: import.meta.env.PB_MFA_DOCS,
                                className: "link-hint",
                                target: "_blank",
                                rel: "noopener noreferrer",
                                textContent: "Learn more.",
                            }),
                        ),
                    ),
                ),
            ),
            t.div(
                { className: "col-sm-12" },
                t.div(
                    { className: "field" },
                    t.input({
                        type: "checkbox",
                        id: uniqueId + ".enabled",
                        name: "mfa.enabled",
                        className: "switch",
                        checked: () => data.config.enabled,
                        onchange: (e) => {
                            data.config.enabled = e.target.checked;

                            if (data.isSuperusers) {
                                collection.otp.enabled = e.target.checked;
                            }
                        },
                    }),
                    t.label({
                        htmlFor: uniqueId + ".enabled",
                        textContent: "Enable",
                    }),
                ),
            ),
            t.div(
                { className: "col-sm-12" },
                t.div(
                    { className: "field" },
                    t.label({
                        htmlFor: uniqueId + ".duration",
                        textContent: "Max duration between 2 authentications (in seconds)",
                    }),
                    t.input({
                        type: "number",
                        id: uniqueId + ".duration",
                        name: "mfa.duration",
                        min: 1,
                        step: 1,
                        required: true,
                        value: () => data.config.duration || "",
                        oninput: (e) => (data.config.duration = parseInt(e.target.value, 10)),
                    }),
                ),
            ),
            t.div(
                { className: "col-sm-12" },
                app.components.ruleField({
                    label: "MFA rule",
                    id: uniqueId + ".rule",
                    name: "mfa.rule",
                    nullable: false,
                    placeholder: "Leave empty to require MFA for everyone",
                    autocomplete: (word) => {
                        return app.utils.collectionAutocompleteKeys(collection, word);
                    },
                    value: () => data.config.rule || "",
                    oninput: (newVal) => (data.config.rule = newVal),
                }),
                t.div(
                    { className: "field-help" },
                    t.p(null, "This optional rule could be used to enable/disable MFA per account basis."),
                    t.p(
                        null,
                        "For example, to require MFA only for accounts with non-empty email you can set it to ",
                        t.code(null, "email != ''"),
                        ".",
                    ),
                    t.p(null, "Leave the rule empty to require MFA for everyone."),
                ),
            ),
        ),
    );
}
