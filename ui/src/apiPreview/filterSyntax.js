export function filterSyntax() {
    const data = store({
        show: false,
    });

    return t.div(
        { className: "filter-details m-t-10" },
        t.button(
            {
                type: "button",
                className: "btn secondary sm",
                onclick: () => data.show = !data.show,
            },
            () => {
                if (data.show) {
                    return [
                        t.span({ className: "txt" }, "Hide details"),
                        t.i({ className: "ri-arrow-up-s-line" }),
                    ];
                }

                return [
                    t.span({ className: "txt" }, "Show details"),
                    t.i({ className: "ri-arrow-down-s-line" }),
                ];
            },
        ),
        app.components.slide(
            () => data.show,
            t.div(
                { className: "block p-t-5" },
                t.p(
                    null,
                    "The filter syntax follows the format ",
                    t.code(
                        null,
                        t.span({ className: "txt-success" }, "OPERAND"),
                        t.span({ className: "txt-danger" }, " OPERATOR "),
                        t.span({ className: "txt-success" }, "OPERAND"),
                    ),
                    ", where:",
                ),
                t.ul(
                    null,
                    t.li(
                        null,
                        t.span({ className: "txt-code txt-success" }, "OPERAND"),
                        " could be any of the above field literal, function, string (single or double quoted), number, null, true, false",
                    ),
                    t.li(
                        null,
                        t.span({ className: "txt-code txt-danger" }, "OPERATOR"),
                        " is one of:",
                        t.ul(
                            null,
                            t.li(null, t.code({ className: "filter-op" }, "="), " Equal"),
                            t.li(null, t.code({ className: "filter-op" }, "!="), " Not equal"),
                            t.li(null, t.code({ className: "filter-op" }, ">"), " Greater than"),
                            t.li(null, t.code({ className: "filter-op" }, ">="), " Greater than or equal"),
                            t.li(null, t.code({ className: "filter-op" }, "<"), " Less than"),
                            t.li(null, t.code({ className: "filter-op" }, "<="), " Less than or equal"),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "~"),
                                " Like/Contains",
                                t.div(
                                    { className: "txt-sm txt-hint" },
                                    t.em(
                                        null,
                                        "(auto wraps the right string OPERAND in a \"%\" for wildcard match if not explicitly set)",
                                    ),
                                ),
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "!~"),
                                " NOT Like/Contains",
                                t.div(
                                    { className: "txt-sm txt-hint" },
                                    t.em(
                                        null,
                                        "(auto wraps the right string OPERAND in a \"%\" for wildcard match if not explicitly set)",
                                    ),
                                ),
                            ),
                            // any/at-least-one-of
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?="),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Equal",
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?!="),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Not equal",
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?>"),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Greater than",
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?>="),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Greater than or equal",
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?<"),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Less than",
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?<="),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Less than or equal",
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?~"),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " Like/Contains",
                                t.div(
                                    { className: "txt-sm txt-hint" },
                                    t.em(
                                        null,
                                        "(auto wraps the right string OPERAND in a \"%\" for wildcard match if not explicitly set)",
                                    ),
                                ),
                            ),
                            t.li(
                                null,
                                t.code({ className: "filter-op" }, "?!~"),
                                t.span({ className: "txt-hint" }, " Any/At-least-one-of"),
                                " NOT Like/Contains",
                                t.div(
                                    { className: "txt-sm txt-hint" },
                                    t.em(
                                        null,
                                        "(auto wraps the right string OPERAND in a \"%\" for wildcard match if not explicitly set)",
                                    ),
                                ),
                            ),
                        ),
                        t.p(
                            null,
                            "To group and combine several expressions you could use brackets ",
                            t.code(null, "(...)"),
                            ", ",
                            t.code(null, "&&"),
                            ", (AND) and ",
                            t.code(null, "||"),
                            " (OR) tokens.",
                        ),
                    ),
                ),
            ),
        ),
    );
}
