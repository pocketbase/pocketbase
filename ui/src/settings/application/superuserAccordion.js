export function superuserAccordion(pageData) {
    const info = store({
        isLoading: false,
        realIP: "",
    });

    async function loadInfo() {
        info.isLoading = true;

        try {
            const health = await app.pb.health.check({ requestKey: "loadSuperuserIPsInfo" });

            info.realIP = health.data?.realIP || "";
            info.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                info.isLoading = false;
            }
        }
    }

    return t.details(
        {
            pbEvent: "superuserAccordion",
            className: "accordion superuser-accordion",
            name: "settingsAccordion",
            onmount: (el) => {
                el._ipwatcher?.unwatch();
                el._ipwatcher = watch(
                    () => JSON.stringify(app.store.settings?.trustedProxy?.headers),
                    (newHash, oldHash) => {
                        if (newHash != oldHash) {
                            loadInfo();
                        }
                    },
                );
            },
            onunmount: (el) => {
                el._ipwatcher?.unwatch();
            },
        },
        t.summary(
            null,
            t.i({ className: "ri-fingerprint-2-line", ariaHidden: true }),
            t.span({ className: "txt" }, "Superuser IPs"),
            t.div({ className: "flex-fill" }),
            () => {
                if (pageData.formSettings?.superuserIPs?.length) {
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
            { className: "content m-b-sm" },
            t.p(null, "A comma separated list of superusers allowed IPs and subnets."),
            t.p(
                null,
                "Enabling this option greatly helps harden the security of your application because even if someone manages to get their hands on a superuser auth token they will not be able to use it.",
            ),
            t.p(
                null,
                "In case your IP changes, you can always reset the field value with the ",
                t.a(
                    {
                        href: import.meta.env.PB_SUPERUSER_IPS_RESET_DOCS,
                        target: "_blank",
                        rel: "noopener noreferrer",
                        className: "link-primary txt-bold txt-sm",
                    },
                    t.code(
                        null,
                        "superuser ips",
                        t.i({ ariaHidden: true, className: "ri-arrow-right-up-line txt-sm" }),
                    ),
                ),
                " console command.",
            ),
        ),
        t.div(
            { className: "fields" },
            t.div(
                { className: "field" },
                t.label(
                    { htmlFor: "superuserIPs" },
                    t.span({ className: "txt" }, "Superuser IPs and subnets"),
                ),
                t.input({
                    id: "superuserIPs",
                    name: "superuserIPs",
                    type: "text",
                    placeholder: "Leave empty for no restriction",
                    value: () => app.utils.joinNonEmpty(pageData.formSettings.superuserIPs),
                    oninput: (e) => {
                        const newValue = app.utils.splitNonEmpty(e.target.value, ",");
                        const newStr = app.utils.joinNonEmpty(newValue);
                        const oldStr = app.utils.joinNonEmpty(pageData.formSettings.superuserIPs);

                        // has an actual change
                        if (oldStr != newStr) {
                            pageData.formSettings.superuserIPs = newValue;
                        }
                    },
                }),
            ),
            t.div(
                { className: "field addon" },
                t.button(
                    {
                        type: "button",
                        className: () =>
                            `btn sm secondary transparent ${
                                app.utils.isEmpty(pageData.formSettings.superuserIPs) ? "hidden" : ""
                            }`,
                        onclick: () => {
                            pageData.formSettings.superuserIPs = [];

                            if (app.store.errors?.superuserIPs) {
                                delete app.store.errors.superuserIPs;
                            }
                        },
                    },
                    t.span({ className: "txt" }, "Clear"),
                ),
            ),
        ),
        t.div(
            { className: "field-help" },
            "Comma separated list of IPs and subnets such as: ",
            t.div(
                { className: "inline-flex gap-5" },
                t.div({
                    role: "button",
                    className: "label sm link-primary",
                    onclick: () => {
                        if (info.isLoading) {
                            return;
                        }

                        const ips = app.utils.toArray(pageData.formSettings.superuserIPs);
                        app.utils.pushUnique(ips, info.realIP);
                        pageData.formSettings.superuserIPs = ips;
                    },
                    textContent: () => info.isLoading ? "..." : (info.realIP + " (you)"),
                }),
            ),
        ),
    );
}
