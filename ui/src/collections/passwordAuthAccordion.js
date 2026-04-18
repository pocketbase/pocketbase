export function passwordAuthAccordion(collection) {
    const uniqueId = "passwordAuth_" + app.utils.randomString();

    const data = store({
        get config() {
            if (!collection.passwordAuth) {
                collection.passwordAuth = {
                    enabled: true,
                    identityFields: ["email"],
                };
            }

            return collection.passwordAuth;
        },
        get identityFieldOptions() {
            // email is always available in auth collections
            const options = [{ value: "email" }];

            const fields = collection?.fields || [];
            const indexes = collection?.indexes || [];

            for (let index of indexes) {
                const parsed = app.utils.parseIndex(index);
                if (!parsed.unique || parsed.columns.length != 1 || parsed.columns[0].name == "email") {
                    continue;
                }

                const field = fields.find((f) => {
                    return !f.hidden && f.name.toLowerCase() == parsed.columns[0].name.toLowerCase();
                });
                if (field) {
                    options.push({ value: field.name });
                }
            }

            return options;
        },
    });

    return t.details(
        {
            pbEvent: "passwordAuthAccordion",
            name: "auth-methods",
            className: "accordion password-auth-accordion",
        },
        t.summary(
            null,
            t.i({ className: "ri-lock-password-line", ariaHidden: true }),
            t.span({ className: "txt", textContent: "Identity/Password" }),
            t.span({
                className: () => `label m-l-auto ${data.config.enabled ? "success" : ""}`,
                textContent: () => (data.config.enabled ? "Enabled" : "Disabled"),
            }),
            () => {
                if (!app.store.errors?.passwordAuth) {
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
                    { className: "field" },
                    t.input({
                        type: "checkbox",
                        id: uniqueId + ".enabled",
                        name: "passwordAuth.enabled",
                        className: "switch",
                        checked: () => data.config.enabled,
                        onchange: (e) => (data.config.enabled = e.target.checked),
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
                        htmlFor: uniqueId + ".identityFields",
                        textContent: "Identity fields",
                    }),
                    app.components.select({
                        id: uniqueId + ".identityFields",
                        name: "passwordAuth.identityFields",
                        max: 99,
                        required: true,
                        options: () => data.identityFieldOptions,
                        value: () => data.config.identityFields,
                        onchange: (selectedOpts) => {
                            data.config.identityFields = selectedOpts.map((opt) => opt.value);
                        },
                    }),
                ),
                t.div(
                    { className: "field-help" },
                    "Only non-hidden fields with UNIQUE index constraint can be selected.",
                ),
            ),
        ),
    );
}
