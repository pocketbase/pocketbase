window.app = window.app || {};
window.app.oauth2 = window.app.oauth2 || {};

// note: data is the providerSettingsModal form store
window.app.oauth2.microsoft = function(providerInfo, namePrefix, data) {
    const uniqueId = "microsoft_" + app.utils.randomString();

    return t.div(
        { pbEvent: "oauth2MicrosoftOptions", className: "oauth2-microsoft-options" },
        t.p({ className: "txt-bold" }, "Azure AD endpoints"),
        t.div(
            { className: "grid" },
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "field" },
                    t.label({ htmlFor: uniqueId + ".authURL" }, "Auth URL"),
                    t.input({
                        id: uniqueId + ".authURL",
                        name: namePrefix + ".authURL",
                        type: "url",
                        required: true,
                        value: () => data.config.authURL || "",
                        oninput: (e) => data.config.authURL = e.target.value,
                    }),
                ),
                t.div(
                    { className: "field-help" },
                    "Ex. https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/authorize",
                ),
            ),
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "field" },
                    t.label({ htmlFor: uniqueId + ".tokenURL" }, "Token URL"),
                    t.input({
                        id: uniqueId + ".tokenURL",
                        name: namePrefix + ".tokenURL",
                        type: "url",
                        required: true,
                        value: () => data.config.tokenURL || "",
                        oninput: (e) => data.config.tokenURL = e.target.value,
                    }),
                ),
                t.div(
                    { className: "field-help" },
                    "Ex. https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/token",
                ),
            ),
        ),
    );
};
