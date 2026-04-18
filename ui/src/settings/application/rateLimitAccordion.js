import { basePredefinedTags, openRateLimitInfoModal } from "./rateLimitInfoModal";

// sort the specified rules list in place
export function sortRules(rules) {
    if (!rules) {
        return;
    }

    let compare = [{}, {}];

    rules.sort((a, b) => {
        compare[0].length = a.label.length;
        compare[0].isTag = a.label.includes(":") || !a.label.includes("/");
        compare[0].isWildcardTag = compare[0].isTag && a.label.startsWith("*");
        compare[0].isExactTag = compare[0].isTag && !compare[0].isWildcardTag;
        compare[0].isPrefix = !compare[0].isTag && a.label.endsWith("/");
        compare[0].hasMethod = !compare[0].isTag && a.label.includes(" /");

        compare[1].length = b.label.length;
        compare[1].isTag = b.label.includes(":") || !b.label.includes("/");
        compare[1].isWildcardTag = compare[1].isTag && b.label.startsWith("*");
        compare[1].isExactTag = compare[1].isTag && !compare[1].isWildcardTag;
        compare[1].isPrefix = !compare[1].isTag && b.label.endsWith("/");
        compare[1].hasMethod = !compare[1].isTag && b.label.includes(" /");

        for (let item of compare) {
            item.priority = 0; // reset

            if (item.isTag) {
                item.priority += 1000;

                if (item.isExactTag) {
                    item.priority += 10;
                } else {
                    item.priority += 5;
                }
            } else {
                if (item.hasMethod) {
                    item.priority += 10;
                }

                if (!item.isPrefix) {
                    item.priority += 5;
                }
            }
        }
        // sort additionally prefix paths based on their length
        if (
            compare[0].isPrefix
            && compare[1].isPrefix
            && ((compare[0].hasMethod && compare[1].hasMethod) || (!compare[0].hasMethod && !compare[1].hasMethod))
        ) {
            if (compare[0].length > compare[1].length) {
                compare[0].priority += 1;
            } else if (compare[0].length < compare[1].length) {
                compare[1].priority += 1;
            }
        }

        if (compare[0].priority > compare[1].priority) {
            return -1;
        }

        if (compare[0].priority < compare[1].priority) {
            return 1;
        }

        return 0;
    });

    return rules;
}

export function rateLimitAccordion(pageData) {
    const audienceOptions = [
        { value: "", label: "All" },
        { value: "@guest", label: "Guest only" },
        { value: "@auth", label: "Auth only" },
    ];

    const accordionData = store({
        predefinedTags: basePredefinedTags,
    });

    loadPredefinedTags();

    async function loadPredefinedTags() {
        let collections = [];

        // fetch an up-to-date collections list
        try {
            collections = await app.pb.collections.getFullList();
        } catch (err) {
            console.warn("loadPredefinedTags: failed to load collections", err);
            return;
        }

        accordionData.predefinedTags = [];

        for (const collection of collections) {
            if (collection.system) {
                continue;
            }

            accordionData.predefinedTags.push({ value: collection.name + ":list" });
            accordionData.predefinedTags.push({ value: collection.name + ":view" });

            if (collection.type != "view") {
                accordionData.predefinedTags.push({ value: collection.name + ":create" });
                accordionData.predefinedTags.push({ value: collection.name + ":update" });
                accordionData.predefinedTags.push({ value: collection.name + ":delete" });
            }

            if (collection.type == "auth") {
                accordionData.predefinedTags.push({
                    value: collection.name + ":listAuthMethods",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":authRefresh",
                });
                accordionData.predefinedTags.push({ value: collection.name + ":auth" });
                accordionData.predefinedTags.push({
                    value: collection.name + ":authWithPassword",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":authWithOAuth2",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":authWithOTP",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":requestOTP",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":requestPasswordReset",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":confirmPasswordReset",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":requestVerification",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":confirmVerification",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":requestEmailChange",
                });
                accordionData.predefinedTags.push({
                    value: collection.name + ":confirmEmailChange",
                });
            }

            if (collection.fields.find((f) => f.type == "file")) {
                accordionData.predefinedTags.push({ value: collection.name + ":file" });
            }
        }

        accordionData.predefinedTags = accordionData.predefinedTags.concat(basePredefinedTags);
    }

    function newRule() {
        if (!Array.isArray(pageData.formSettings.rateLimits.rules)) {
            pageData.formSettings.rateLimits.rules = [];
        }

        pageData.formSettings.rateLimits.rules.push({
            label: "",
            maxRequests: 200,
            duration: 3,
            audience: "",
        });

        // enable the rate limiter if this is the first rule that is being added
        if (pageData.formSettings.rateLimits.rules.length == 1) {
            pageData.formSettings.rateLimits.enabled = true;
        }
    }

    function removeRule(i) {
        pageData.formSettings.rateLimits.rules.splice(i, 1);

        if (!pageData.formSettings.rateLimits.rules.length) {
            pageData.formSettings.rateLimits.enabled = false;
        }
    }

    const watchers = [];

    return t.details(
        {
            pbEvent: "rateLimitAccordion",
            className: "accordion rate-limit-accordion",
            name: "settingsAccordion",
            onmount: () => {
                watchers.push(
                    // clear rules errors on any rule change since an error could be
                    // for a duplicated tag that may have been updated in a different rule
                    watch(
                        () => JSON.stringify(pageData.formSettings.rateLimits.rules),
                        () => {
                            if (!app.store.errors?.rateLimits?.rules) {
                                return;
                            }

                            delete app.store.errors.rateLimits;
                        },
                    ),
                );
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.summary(
            null,
            t.i({ className: "ri-pulse-fill", ariaHidden: true }),
            t.span({ className: "txt" }, "Rate limiting"),
            t.div({ className: "flex-fill" }),
            () => {
                if (pageData.formSettings.rateLimits.enabled) {
                    return t.span({ className: "label success" }, "Enabled");
                }
                return t.span({ className: "label" }, "Disabled");
            },
            () => {
                if (!app.utils.isEmpty(app.store.errors?.rateLimits)) {
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
                        id: "rateLimits.enabled",
                        name: "rateLimits.enabled",
                        type: "checkbox",
                        className: "switch",
                        checked: () => pageData.formSettings.rateLimits.enabled || false,
                        onchange: (e) => (pageData.formSettings.rateLimits.enabled = e.target.checked),
                    }),
                    t.label(
                        { htmlFor: "rateLimits.enabled" },
                        t.span({ className: "txt" }, "Enable"),
                        t.small({ className: "txt-hint" }, " (experimental)"),
                    ),
                ),
            ),
            t.div(
                { className: "col-lg-12" },
                t.div(
                    { className: "rate-limit-table-wrapper" },
                    t.table(
                        { className: "rate-limit-table" },
                        t.thead(
                            {
                                hidden: () => !pageData.formSettings.rateLimits.rules?.length,
                            },
                            t.tr(
                                null,
                                t.th({ className: "col-label" }, "Rate limit label"),
                                t.th(
                                    { className: "col-requests" },
                                    "Max requests",
                                    t.br(),
                                    t.small(null, "(per IP)"),
                                ),
                                t.th(
                                    { className: "col-duration" },
                                    "Interval",
                                    t.br(),
                                    t.small(null, "(in seconds)"),
                                ),
                                t.th({ className: "col-audience" }, "Targeted users"),
                                t.th({ className: "col-action" }),
                            ),
                        ),
                        t.tbody(null, () => {
                            const rows = [];
                            const rules = pageData.formSettings.rateLimits.rules || [];

                            for (let i = 0; i < rules.length; i++) {
                                const rule = rules[i];

                                rows.push(
                                    t.tr(
                                        { className: "rate-limit-row" },
                                        t.td(
                                            { className: "col-label" },
                                            t.div(
                                                { className: "field" },
                                                t.input({
                                                    type: "text",
                                                    required: true,
                                                    className: "inline-error",
                                                    id: "rateLimits.rules." + i + ".label",
                                                    name: "rateLimits.rules." + i + ".label",
                                                    placeholder: "tag (users:create) or path (/api/)",
                                                    "html-list": "rateLimits.rules." + i + ".label_list",
                                                    value: () => rule.label,
                                                    oninput: (e) => (rule.label = e.target.value),
                                                }),
                                                t.datalist(
                                                    {
                                                        id: "rateLimits.rules." + i + ".label_list",
                                                    },
                                                    () => {
                                                        return accordionData.predefinedTags.map((tag) => {
                                                            return t.option({ value: tag.value }, tag.label || "");
                                                        });
                                                    },
                                                ),
                                            ),
                                        ),
                                        t.td(
                                            { className: "col-requests" },
                                            t.div(
                                                { className: "field" },
                                                t.input({
                                                    type: "number",
                                                    required: true,
                                                    placeholder: "Max requests*",
                                                    className: "inline-error",
                                                    min: 1,
                                                    step: 1,
                                                    name: "rateLimits.rules." + i + ".maxRequests",
                                                    value: () => rule.maxRequests || 0,
                                                    oninput: (e) => rule.maxRequests = parseInt(e.target.value, 10),
                                                }),
                                            ),
                                        ),
                                        t.td(
                                            { className: "col-duration" },
                                            t.div(
                                                { className: "field" },
                                                t.input({
                                                    type: "number",
                                                    required: true,
                                                    placeholder: "Interval*",
                                                    className: "inline-error",
                                                    min: 1,
                                                    step: 1,
                                                    name: "rateLimits.rules." + i + ".duration",
                                                    value: () => rule.duration,
                                                    oninput: (e) => rule.duration = parseInt(e.target.value, 10),
                                                }),
                                            ),
                                        ),
                                        t.td(
                                            { className: "col-audience" },
                                            t.div(
                                                { className: "field" },
                                                app.components.select({
                                                    name: "rateLimits.rules." + i + ".audience",
                                                    className: "inline-error",
                                                    options: audienceOptions,
                                                    required: true,
                                                    value: () => rule.audience || "",
                                                    onchange: (selected) => {
                                                        rule.audience = selected?.[0]?.value;
                                                    },
                                                }),
                                            ),
                                        ),
                                        t.td(
                                            { className: "col-action" },
                                            t.button(
                                                {
                                                    type: "button",
                                                    araiaDescription: app.attrs.tooltip("Remove rule"),
                                                    className: "btn sm secondary transparent circle",
                                                    onclick: () => removeRule(i),
                                                },
                                                t.i({ className: "ri-close-line" }),
                                            ),
                                        ),
                                    ),
                                );
                            }

                            return rows;
                        }),
                    ),
                ),
                t.div(
                    { className: "flex m-t-sm" },
                    t.button(
                        {
                            type: "button",
                            className: "btn secondary sm",
                            onclick: () => newRule(),
                        },
                        t.i({ className: "ri-add-line", ariaHidden: true }),
                        t.span({ className: "txt" }, "Add rate limit rule"),
                    ),
                    t.button(
                        {
                            type: "button",
                            className: "link-hint txt-sm m-l-auto",
                            onclick: () => openRateLimitInfoModal(),
                        },
                        t.em(null, "Learn more about the rate limit rules"),
                    ),
                ),
            ),
        ),
    );
}
