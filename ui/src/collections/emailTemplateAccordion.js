export function emailTemplateAccordion(collection, key, propsArg = {}) {
    const uniqueId = "emailTemplate" + app.utils.randomString();

    const props = store({
        title: "Email template",
        placeholders: [],
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        get config() {
            let val = app.utils.getByPath(collection, key);

            if (!val) {
                val = { subject: "", body: "" };
                app.utils.setByPath(collection, key, val);
            }

            return val;
        },
        get tokensList() {
            return [];
        },
    });

    const placeholdersList = () => {
        if (!props.placeholders?.length) {
            return;
        }

        return t.div(
            { className: "field-help" },
            t.div({ className: "flex flex-wrap gap-5" }, t.span({ className: "txt" }, "Placeholders:"), () => {
                return props.placeholders.map((p) => {
                    return t.span({ className: "label sm" }, app.components.copyButton(p, p));
                });
            }),
        );
    };

    return t.details(
        {
            pbEvent: "emailTemplateAccordion",
            className: "accordion email-template-accordion",
            name: "email-template",
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.summary(
            null,
            t.i({ className: "ri-draft-line" }),
            t.span({ className: "txt", textContent: () => props.title }),
            () => {
                if (!app.utils.getByPath(app.store.errors, key)) {
                    return;
                }

                return t.i({
                    className: "ri-error-warning-fill txt-danger m-l-auto",
                    ariaDescription: app.attrs.tooltip("Has errors", "left"),
                });
            },
        ),
        t.div(
            { className: "grid" },
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "field" },
                    t.label({
                        htmlFor: uniqueId + ".subject",
                        textContent: "Subject",
                    }),
                    app.components.codeEditor({
                        id: uniqueId + ".subject",
                        name: key + ".subject",
                        required: true,
                        singleLine: true,
                        language: "text",
                        autocomplete: props.placeholders,
                        value: () => data.config.subject || "",
                        oninput: (val) => (data.config.subject = val),
                    }),
                ),
                placeholdersList,
            ),
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "field" },
                    t.label({
                        htmlFor: uniqueId + ".body",
                        textContent: "Body (HTML)",
                    }),
                    app.components.codeEditor({
                        id: uniqueId + ".body",
                        name: key + ".body",
                        required: true,
                        language: "html",
                        className: "pre-wrap",
                        autocomplete: props.placeholders,
                        value: () => data.config.body || "",
                        oninput: (val) => (data.config.body = val),
                    }),
                ),
                placeholdersList,
            ),
        ),
    );
}
