window.app = window.app || {};
window.app.oauth2 = window.app.oauth2 || {};

// note: data is the providerSettingsModal form store
window.app.oauth2.lark = function(providerInfo, namePrefix, data) {
    const uniqueId = "lark_" + app.utils.randomString();

    const domainOptions = [
        { label: "Feishu (China)", value: "feishu.cn" },
        { label: "Lark (International)", value: "larksuite.com" },
    ];

    const local = store({
        domain: data.config.authURL?.includes(domainOptions[1].value)
            ? domainOptions[1].value
            : domainOptions[0].value,
    });

    const watchers = [
        watch(() => local.domain, (domain) => {
            if (domain) {
                data.config.authURL = `https://accounts.${domain}/open-apis/authen/v1/authorize`;
                data.config.tokenURL = `https://open.${domain}/open-apis/authen/v2/oauth/token`;
                data.config.userInfoURL = `https://open.${domain}/open-apis/authen/v1/user_info`;
            }
        }),
    ];

    return t.div(
        {
            pbEvent: "oauth2LarkOptions",
            className: "oauth2-lark-options",
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: "grid" },
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "field" },
                    t.label({ htmlFor: uniqueId + ".site" }, "Site"),
                    app.components.select({
                        options: domainOptions,
                        required: true,
                        value: () => local.domain || "",
                        onchange: (selectedOpts) => {
                            local.domain = selectedOpts?.[0]?.value;
                        },
                    }),
                ),
            ),
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "alert info" },
                    "Note that the Lark user's ",
                    t.strong(null, "Union ID"),
                    " will be used for the association with the PocketBase user (see ",
                    t.a({
                        href:
                            "https://open.feishu.cn/document/platform-overveiw/basic-concepts/user-identity-introduction/introduction#3f2d4b63",
                        target: "_blank",
                        rel: "noopener noreferrer",
                        textContent: "Different Types of Lark User IDs",
                    }),
                    ").",
                ),
            ),
        ),
    );
};
