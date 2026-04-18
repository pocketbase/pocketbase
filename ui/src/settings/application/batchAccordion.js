export function batchAccordion(pageData) {
    return t.details(
        {
            pbEvent: "batchApiAccordion",
            className: "accordion batch-api-accordion",
            name: "settingsAccordion",
        },
        t.summary(
            null,
            t.i({ className: "ri-archive-stack-line", ariaHidden: true }),
            t.span({ className: "txt" }, "Batch API"),
            t.div({ className: "flex-fill" }),
            () => {
                if (pageData.formSettings.batch.enabled) {
                    return t.span({ className: "label success" }, "Enabled");
                }
                return t.span({ className: "label" }, "Disabled");
            },
            () => {
                if (!app.utils.isEmpty(app.store.errors?.batch)) {
                    return t.i({
                        className: "ri-error-warning-fill txt-danger",
                        ariaDescription: app.attrs.tooltip("Has errors", "left"),
                    });
                }
            },
        ),
        t.div(
            { className: "grid sm" },
            t.div(
                { className: "col-lg-12" },
                t.div(
                    { className: "field" },
                    t.input({
                        id: "batch.enabled",
                        name: "batch.enabled",
                        type: "checkbox",
                        className: "switch",
                        checked: () => pageData.formSettings.batch.enabled || false,
                        onchange: (e) => (pageData.formSettings.batch.enabled = e.target.checked),
                    }),
                    t.label(
                        { htmlFor: "batch.enabled" },
                        t.span({ className: "txt" }, "Enable"),
                        t.small({ className: "txt-hint" }, " (experimental)"),
                    ),
                ),
            ),
            t.div(
                { className: "col-lg-4" },
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: "batch.maxRequests" },
                        t.span({ className: "txt" }, "Max requests in a batch"),
                        t.i({
                            className: "ri-information-line link-faded",
                            ariaDescription: app.attrs.tooltip(
                                "Rate limiting (if enabled) also applies for the batch create/update/upsert/delete requests.",
                                "right",
                            ),
                        }),
                    ),
                    t.input({
                        id: "batch.maxRequests",
                        name: "batch.maxRequests",
                        type: "number",
                        min: 1,
                        step: 1,
                        required: () => pageData.formSettings.batch.enabled,
                        disabled: () => !pageData.formSettings.batch.enabled,
                        value: () => pageData.formSettings.batch.maxRequests,
                        oninput: (e) => (pageData.formSettings.batch.maxRequests = e.target.value << 0),
                    }),
                ),
            ),
            t.div(
                { className: "col-lg-4" },
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: "batch.timeout" },
                        t.span({ className: "txt" }, "Max processing time (in seconds)"),
                    ),
                    t.input({
                        id: "batch.timeout",
                        name: "batch.timeout",
                        type: "number",
                        min: 1,
                        step: 1,
                        required: () => pageData.formSettings.batch.enabled,
                        disabled: () => !pageData.formSettings.batch.enabled,
                        value: () => pageData.formSettings.batch.timeout,
                        oninput: (e) => pageData.formSettings.batch.timeout = parseInt(e.target.value, 10),
                    }),
                ),
            ),
            t.div(
                { className: "col-lg-4" },
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: "batch.maxBodySize" },
                        t.span({ className: "txt" }, "Max body size (in bytes)"),
                    ),
                    t.input({
                        id: "batch.maxBodySize",
                        name: "batch.maxBodySize",
                        type: "number",
                        min: 0,
                        step: 1,
                        placeholder: "Default to 128MB",
                        disabled: () => !pageData.formSettings.batch.enabled,
                        value: () => pageData.formSettings.batch.maxBodySize || "",
                        oninput: (e) => pageData.formSettings.batch.maxBodySize = parseInt(e.target.value, 10),
                    }),
                ),
            ),
        ),
    );
}
