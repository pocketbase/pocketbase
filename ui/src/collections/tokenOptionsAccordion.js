export function tokenOptionsAccordion(collection) {
    const uniqueId = "token_" + app.utils.randomString();

    const data = store({
        get tokensList() {
            if (collection?.name === "_superusers") {
                return [
                    { key: "authToken", label: "Auth" },
                    { key: "passwordResetToken", label: "Password reset" },
                    { key: "fileToken", label: "Protected file" },
                ];
            }

            return [
                { key: "authToken", label: "Auth" },
                { key: "verificationToken", label: "Email verification" },
                { key: "passwordResetToken", label: "Password reset" },
                { key: "emailChangeToken", label: "Email change" },
                { key: "fileToken", label: "Protected file" },
            ];
        },
    });

    return t.details(
        {
            pbEvent: "tokenOptionsAccordion",
            name: "other",
            className: "accordion token-options-accordion",
        },
        t.summary(
            null,
            t.i({ className: "ri-key-2-line" }),
            t.span({ className: "txt", textContent: "Token options (invalidate, duration)" }),
        ),
        t.div({ className: "grid sm" }, () => {
            return data.tokensList.map((token) => {
                const fieldId = uniqueId + token.key;

                return t.div(
                    { className: "col-sm-6" },
                    t.div(
                        { className: "field token-field" },
                        t.label({
                            htmlFor: fieldId,
                            textContent: () => token.label + " duration (in seconds)",
                        }),
                        t.input({
                            id: fieldId,
                            type: "number",
                            min: 1,
                            step: 1,
                            required: true,
                            name: () => token.key + ".duration",
                            value: () => collection[token.key].duration,
                            oninput: (e) => (collection[token.key].duration = parseInt(e.target.value, 10)),
                        }),
                    ),
                    t.div(
                        { className: "field-help m-b-10" },
                        t.button({
                            type: "button",
                            className: () => `link-hint ${collection[token.key].secret ? "txt-success" : ""}`,
                            textContent: "Invalidate all previously issued tokens",
                            onclick: () => {
                                // toggle
                                if (collection[token.key].secret) {
                                    delete collection[token.key].secret;
                                } else {
                                    collection[token.key].secret = app.utils.randomSecret(50);
                                }
                            },
                        }),
                    ),
                );
            });
        }),
    );
}
