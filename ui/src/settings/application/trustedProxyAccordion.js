export function trustedProxyAccordion(pageData) {
    const commonProxyHeaders = ["X-Forwarded-For", "Fly-Client-IP", "CF-Connecting-IP"];

    const ipOptions = [
        { label: "Use leftmost IP", value: true },
        { label: "Use rightmost IP", value: false },
    ];

    const proxyInfo = store({
        isLoading: false,
        realIP: "",
        possibleProxyHeader: "",
        get suggestedProxyHeaders() {
            if (!proxyInfo.possibleProxyHeader) {
                return commonProxyHeaders;
            }

            return [proxyInfo.possibleProxyHeader].concat(
                commonProxyHeaders.filter((h) => h != proxyInfo.possibleProxyHeader),
            );
        },
        get isEnabled() {
            return !app.utils.isEmpty(pageData.formSettings.trustedProxy?.headers);
        },
    });

    loadProxyInfo();

    async function loadProxyInfo() {
        proxyInfo.isLoading = true;

        try {
            const health = await app.pb.health.check({ requestKey: "loadProxyInfo" });

            proxyInfo.realIP = health.data?.realIP || "";
            proxyInfo.possibleProxyHeader = health.data?.possibleProxyHeader || "";
            proxyInfo.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                proxyInfo.isLoading = false;
            }
        }
    }

    return t.details(
        {
            pbEvent: "trustedProxyAccordion",
            className: "accordion trusted-proxy-accordion",
            name: "settingsAccordion",
            open: () => (proxyInfo.isLoading ? false : null),
        },
        t.summary(
            null,
            t.i({ className: "ri-route-line" }),
            t.span({ className: "txt" }, "User IP proxy headers"),
            () => {
                if (proxyInfo.isLoading) {
                    return t.span({ className: "loader sm" });
                }

                if (!proxyInfo.isEnabled && proxyInfo.possibleProxyHeader) {
                    return t.i({
                        className: "ri-alert-line txt-warning",
                        ariaDescription: app.attrs.tooltip(
                            "Detected proxy header.\nIt is recommend to list it as trusted.",
                            "right",
                        ),
                    });
                }

                if (
                    proxyInfo.isEnabled
                    && proxyInfo.possibleProxyHeader
                    && !pageData.formSettings.trustedProxy.headers.includes(proxyInfo.possibleProxyHeader)
                ) {
                    return t.i({
                        className: "ri-alert-line txt-hint",
                        ariaDescription: app.attrs.tooltip(
                            "The configured proxy header doesn't match with the detected one.",
                            "right",
                        ),
                    });
                }
            },
            t.div({ className: "flex-fill" }),
            () => {
                if (proxyInfo.isEnabled) {
                    return t.span({ className: "label success" }, "Enabled");
                }
                return t.span({ className: "label" }, "Disabled");
            },
            () => {
                if (!app.utils.isEmpty(app.store.errors?.trustedProxy)) {
                    return t.i({
                        className: "ri-error-warning-fill txt-danger",
                        ariaDescription: app.attrs.tooltip("Has errors", "left"),
                    });
                }
            },
        ),
        t.p(
            { className: "m-t-0" },
            "Below you should see your real IP. If not - configure the correct proxy header for your environment.",
        ),
        t.div(
            {
                hidden: () => proxyInfo.isLoading,
                className: "alert info m-b-sm",
            },
            t.div(
                { className: "flex gap-5" },
                t.span(null, "Resolved user IP:"),
                t.strong(null, () => proxyInfo.realIP || "N/A"),
            ),
            t.div(
                { className: "flex gap-5" },
                t.span(null, "Detected proxy header:"),
                t.strong(null, () => proxyInfo.possibleProxyHeader || "N/A"),
            ),
        ),
        t.div(
            { className: "content m-b-sm" },
            t.p(
                null,
                `
                When PocketBase is deployed on platforms like Fly or it is accessible through proxies such as
                NGINX, requests from different users will originate from the same IP address (the IP of the proxy
                connecting to your PocketBase app).
            `,
            ),
            t.p(
                null,
                `
                In this case to retrieve the actual user IP (used for rate limiting, logging, etc.) you need to
                properly configure your proxy and list below the trusted headers that PocketBase could use to
                extract the user IP.
            `,
            ),
            t.p({ className: "txt-bold" }, `When using such proxy, to avoid spoofing it is recommended to:`),
            t.ul(
                { className: "txt-bold" },
                t.li(
                    null,
                    "use headers that are controlled only by the proxy and cannot be manually set by the users",
                ),
                t.li(null, "make sure that the PocketBase server can be accessed ONLY through the proxy"),
            ),
            t.p(null, "You can clear the headers field if PocketBase is not deployed behind a proxy."),
        ),
        t.div(
            { className: "grid sm" },
            t.div(
                { className: "col-lg-9" },
                t.div(
                    { className: "fields" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: "trustedProxy.headers" }, "Trusted IP proxy headers"),
                        t.input({
                            type: "text",
                            id: "trustedProxy.headers",
                            name: "trustedProxy.headers",
                            placeholder: "Leave empty to disable",
                            value: () => app.utils.joinNonEmpty(pageData.formSettings.trustedProxy.headers),
                            oninput: (e) => {
                                const newValue = app.utils.splitNonEmpty(e.target.value, ",");
                                const newStr = app.utils.joinNonEmpty(newValue);
                                const oldStr = app.utils.joinNonEmpty(pageData.formSettings.trustedProxy.headers);

                                // has an actual change
                                if (oldStr != newStr) {
                                    pageData.formSettings.trustedProxy.headers = newValue;
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
                                        app.utils.isEmpty(pageData.formSettings.trustedProxy.headers) ? "hidden" : ""
                                    }`,
                                onclick: () => {
                                    pageData.formSettings.trustedProxy.headers = [];
                                },
                            },
                            t.span({ className: "txt" }, "Clear"),
                        ),
                    ),
                ),
                t.div(
                    { className: "field-help" },
                    "Comma separated list of headers such as: ",
                    t.div({ className: "inline-flex gap-5" }, () => {
                        return proxyInfo.suggestedProxyHeaders.map((header) => {
                            return t.div({
                                type: "button",
                                className: "label sm link-hint",
                                onclick: () => {
                                    pageData.formSettings.trustedProxy.headers = [header];
                                },
                                textContent: header,
                            });
                        });
                    }),
                ),
            ),
            t.div(
                { className: "col-lg-3" },
                t.div(
                    { className: "field" },
                    t.label(
                        { htmlFor: "trustedProxy.useLeftmostIP" },
                        t.span({ className: "txt" }, "IP priority"),
                        t.i({
                            className: "ri-information-line tooltip-right",
                            ariaDescription: app.attrs.tooltip(
                                "This is in case the proxy returns more than 1 IP as header value. The rightmost IP is usually considered to be the more trustworthy but this could vary depending on the proxy.",
                            ),
                        }),
                    ),
                    app.components.select({
                        id: "trustedProxy.useLeftmostIP",
                        name: "trustedProxy.useLeftmostIP",
                        options: ipOptions,
                        required: true,
                        value: () => pageData.formSettings.trustedProxy.useLeftmostIP || false,
                        onchange: (selected) => {
                            pageData.formSettings.trustedProxy.useLeftmostIP = selected?.[0]?.value;
                        },
                    }),
                ),
            ),
        ),
    );
}
