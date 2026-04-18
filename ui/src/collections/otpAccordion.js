export function otpAccordion(collection) {
    const uniqueId = "otp_" + app.utils.randomString();

    const data = store({
        get config() {
            if (!collection.otp) {
                collection.otp = {
                    enabled: false,
                    duration: 300,
                    length: 8,
                };
            }

            return collection.otp;
        },
    });

    return t.details(
        {
            pbEvent: "otpAccordion",
            name: "auth-methods",
            className: "accordion otp-accordion",
        },
        t.summary(
            null,
            t.i({ className: "ri-time-line", ariaHidden: true }),
            t.span({ className: "txt", textContent: "One-time password (OTP)" }),
            t.span({
                className: () => `label m-l-auto ${data.config.enabled ? "success" : ""}`,
                textContent: () => (data.config.enabled ? "Enabled" : "Disabled"),
            }),
            () => {
                if (!app.store.errors?.otp) {
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
                        name: "otp.enabled",
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
                { className: "col-sm-6" },
                t.div(
                    { className: "field" },
                    t.label({
                        htmlFor: uniqueId + ".duration",
                        textContent: "Duration (in seconds)",
                    }),
                    t.input({
                        type: "number",
                        id: uniqueId + ".duration",
                        name: "otp.duration",
                        min: 1,
                        step: 1,
                        required: true,
                        value: () => data.config.duration || "",
                        oninput: (e) => (data.config.duration = parseInt(e.target.value, 10)),
                    }),
                ),
            ),
            t.div(
                { className: "col-sm-6" },
                t.div(
                    { className: "field" },
                    t.label({
                        htmlFor: uniqueId + ".length",
                        textContent: "Generated password length",
                    }),
                    t.input({
                        type: "number",
                        id: uniqueId + ".length",
                        name: "otp.length",
                        min: 1,
                        step: 1,
                        required: true,
                        value: () => data.config.length || "",
                        oninput: (e) => (data.config.length = parseInt(e.target.value, 10)),
                    }),
                ),
            ),
        ),
    );
}
