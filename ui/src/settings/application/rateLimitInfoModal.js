export let basePredefinedTags = [
    { value: "*:list" },
    { value: "*:view" },
    { value: "*:create" },
    { value: "*:update" },
    { value: "*:delete" },
    { value: "*:file", description: "targets the files download endpoint" },
    { value: "*:listAuthMethods" },
    { value: "*:authRefresh" },
    { value: "*:auth", description: "targets all auth methods" },
    { value: "*:authWithPassword" },
    { value: "*:authWithOAuth2" },
    { value: "*:authWithOTP" },
    { value: "*:requestOTP" },
    { value: "*:requestPasswordReset" },
    { value: "*:confirmPasswordReset" },
    { value: "*:requestVerification" },
    { value: "*:confirmVerification" },
    { value: "*:requestEmailChange" },
    { value: "*:confirmEmailChange" },
];

export function openRateLimitInfoModal() {
    const modal = rateLimitInfoModal();

    document.body.appendChild(modal);

    app.modals.open(modal);
}

function rateLimitInfoModal() {
    return t.div(
        {
            pbEvent: "rateLimitInfoModal",
            className: "modal rate-limit-info-modal",
            onafterclose: (el) => {
                el?.remove();
            },
        },
        t.header({ className: "modal-header" }, t.h5(null, "Rate limit label format")),
        t.div(
            { className: "modal-content" },
            t.p(null, "The rate limit rules are resolved in the following order (stops on the first match):"),
            t.ol(
                null,
                t.li(null, "exact tag (e.g. ", t.code(null, "users:create")),
                t.li(null, "wildcard tag (e.g. ", t.code(null, "*:create")),
                t.li(null, "METHOD + exact path (e.g. ", t.code(null, "POST /a/b")),
                t.li(null, "METHOD + prefix path (e.g. ", t.code(null, "POST /a/b", t.strong(null, "/"))),
                t.li(null, "exact path (e.g. ", t.code(null, "/a/b")),
                t.li(null, "prefix path (e.g. ", t.code(null, "/a/b", t.strong(null, "/"))),
            ),
            t.p(
                null,
                `In case of multiple rules with the same label but different target user audience (e.g. "guest" vs "auth"), only the matching audience rule is taken in consideration.`,
            ),
            t.hr(),
            t.p(null, "The rate limit label could be in one of the following formats:"),
            t.ul(
                null,
                t.li(
                    { className: "m-b-sm" },
                    t.code(null, "[METHOD ]/my/path"),
                    " - full exact route match (",
                    t.strong(null, "must be without trailing slash"),
                    "; \"METHOD\" is optional).",
                    t.br(),
                    "For example:",
                    t.ul(
                        { className: "m-0" },
                        t.li(
                            null,
                            t.code(null, "/hello"),
                            " - matches ",
                            t.code(null, "GET /hello"),
                            ", ",
                            t.code(null, "POST /hello"),
                            ", etc.",
                        ),
                        t.li(null, t.code(null, "POST /hello"), " - matches only ", t.code(null, "POST /hello")),
                    ),
                ),
                t.li(
                    { className: "m-b-sm" },
                    t.code(null, "[METHOD ]/my/prefix", t.strong(null, "/")),
                    " - path prefix (",
                    t.strong(null, "must end with trailing slash;"),
                    "\"METHOD\" is optional). For example:",
                    t.ul(
                        { className: "m-0" },
                        t.li(
                            null,
                            t.code(null, "/hello/"),
                            " - matches ",
                            t.code(null, "GET /hello"),
                            ", ",
                            t.code(null, "POST /hello/a/b/c"),
                            ", etc.",
                        ),
                        t.li(
                            null,
                            t.code(null, "POST /hello/"),
                            " - matches ",
                            t.code(null, "POST /hello"),
                            ", ",
                            t.code(null, "POST /hello/a/b/c"),
                            ", etc.",
                        ),
                    ),
                ),
                t.li(
                    { className: "m-b-0" },
                    t.code(null, "collectionName:predefinedTag"),
                    " - targets a specific action of a single collection.",
                    " To apply the rule for all collections you can use the ",
                    t.code(null, "*"),
                    " wildcard. For example:",
                    t.code(null, "posts:create"),
                    ", ",
                    t.code(null, "users:listAuthMethods"),
                    ", ",
                    t.code(null, "*:auth"),
                    ".",
                    t.br(),
                    "The predifined collection tags are (",
                    t.em(null, "there should be autocomplete once you start typing"),
                    "):",
                    t.ul({ className: "m-0" }, () => {
                        return basePredefinedTags.map((tag) => {
                            return t.li(null, tag.value.replace("*:", ":"), () => {
                                if (tag.description) {
                                    return t.em({ className: "txt-hint" }, " (", tag.description, ")");
                                }
                            });
                        });
                    }),
                ),
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => app.modals.close(),
                },
                t.span({ className: "txt" }, "Close"),
            ),
        ),
    );
}
