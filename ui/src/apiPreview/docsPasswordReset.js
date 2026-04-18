export function docsPasswordReset(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const actionTabs = [
        { title: "Request password reset", content: request },
        { title: "Confirm password reset", content: confirm },
    ];

    const data = store({
        activeActionIndex: 0,
    });

    return t.div(
        {
            pbEvent: "apiPreviewPasswordReset",
            className: "content",
        },
        // description
        t.p(null, `Sends ${collection.name} password reset email request.`),
        t.p(
            null,
            "On successful password reset all previously issued auth tokens for the specific record will be automatically invalidated.",
        ),
        app.components.codeBlockTabs({
            className: "sdk-examples m-t-sm",
            historyKey: "pbLastSDK",
            tabs: [
                {
                    title: "JS SDK",
                    language: "js",
                    value: `
                        import PocketBase from 'pocketbase';

                        const pb = new PocketBase('${baseURL}');

                        ...

                        await pb.collection('${collection.name}').requestPasswordReset('test@example.com');

                        // ---
                        // (optional) in your custom confirmation page:
                        // ---

                        // note: after this call all previously issued auth tokens are invalidated
                        await pb.collection('${collection.name}').confirmPasswordReset(
                            'RESET_TOKEN',
                            'NEW_PASSWORD',
                            'NEW_PASSWORD_CONFIRM',
                        );
                    `,
                    footnote: t.div(
                        { className: "txt-right" },
                        t.a({
                            href: import.meta.env.PB_JS_SDK_URL,
                            target: "_blank",
                            rel: "noopener noreferrer",
                            textContent: "JS SDK docs",
                        }),
                    ),
                },
                {
                    title: "Dart SDK",
                    language: "dart",
                    value: `
                        import 'package:pocketbase/pocketbase.dart';

                        final pb = PocketBase('${baseURL}');

                        ...

                        await pb.collection('${collection.name}').requestPasswordReset('test@example.com');

                        // ---
                        // (optional) in your custom confirmation page:
                        // ---

                        // note: after this call all previously issued auth tokens are invalidated
                        await pb.collection('${collection.name}').confirmPasswordReset(
                          'RESET_TOKEN',
                          'NEW_PASSWORD',
                          'NEW_PASSWORD_CONFIRM',
                        );
                    `,
                    footnote: t.div(
                        { className: "txt-right" },
                        t.a({
                            href: import.meta.env.PB_DART_SDK_URL,
                            target: "_blank",
                            rel: "noopener noreferrer",
                            textContent: "Dart SDK docs",
                        }),
                    ),
                },
                {
                    title: "curl",
                    language: "bash",
                    value: `
                        # Request password reset
                        curl -X POST \\
                          -H 'Content-Type:application/json' \\
                          -d '{ "email":"..." }' \\
                          '${baseURL}/api/collections/${collection.name}/request-password-reset'

                        # Confirm password reset
                        curl -X POST \\
                          -H 'Content-Type:application/json' \\
                          -d '{ "token":"...", "password":"", "passwordConfirm":"" }' \\
                          '${baseURL}/api/collections/${collection.name}/confirm-password-reset'
                    `,
                },
            ],
        }),
        t.nav(
            { className: "btns m-t-base m-b-sm" },
            () => {
                return actionTabs.map((tab, i) => {
                    return t.button({
                        type: "button",
                        className: () => `btn sm expanded ${data.activeActionIndex == i ? "active" : "secondary"}`,
                        textContent: () => tab.title,
                        onclick: () => data.activeActionIndex = i,
                    });
                });
            },
        ),
        () => actionTabs[data.activeActionIndex]?.content?.(collection),
    );
}

function request(collection) {
    const responses = [
        {
            title: 204,
            value: "null",
        },
        {
            title: 400,
            value: `
                {
                  "status": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "email": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
    ];

    return [
        // api
        t.div({ className: "block" }, t.strong(null, "API details")),
        t.div(
            { className: "alert success api-preview-alert" },
            t.span({ className: "label method" }, "POST"),
            t.span({ className: "path" }, `/api/collections/${collection.name}/request-password-reset`),
        ),
        t.table(
            { className: "api-preview-table body-params" },
            t.thead(
                null,
                t.tr(
                    null,
                    t.th({ className: "min-width txt-primary" }, "Body params"),
                    t.th({ className: "min-width" }, "Type"),
                    t.th(null, "Description"),
                ),
            ),
            t.tbody(
                null,
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "email ", t.em(null, "(required)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, "The auth record email address to send the password reset request (if exists)."),
                ),
            ),
        ),
        // responses
        t.div({ className: "block m-t-base m-b-sm" }, t.strong(null, "Example responses")),
        app.components.codeBlockTabs({
            tabs: responses,
        }),
    ];
}

function confirm(collection) {
    const responses = [
        {
            title: 204,
            value: "null",
        },
        {
            title: 400,
            value: `
                {
                  "status": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
    ];

    return [
        // api
        t.div({ className: "block" }, t.strong(null, "API details")),
        t.div(
            { className: "alert success api-preview-alert" },
            t.span({ className: "label method" }, "POST"),
            t.span({ className: "path" }, `/api/collections/${collection.name}/confirm-password-reset`),
        ),
        t.table(
            { className: "api-preview-table body-params" },
            t.thead(
                null,
                t.tr(
                    null,
                    t.th({ className: "min-width txt-primary" }, "Body params"),
                    t.th({ className: "min-width" }, "Type"),
                    t.th(null, "Description"),
                ),
            ),
            t.tbody(
                null,
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "token ", t.em(null, "(required)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, "The token from the password reset request email."),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "password ", t.em(null, "(required)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, "The new password to set."),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "passwordConfirm ", t.em(null, "(required)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, "Confirmation of the new password."),
                ),
            ),
        ),
        // responses
        t.div({ className: "block m-t-base m-b-sm" }, t.strong(null, "Example responses")),
        app.components.codeBlockTabs({
            tabs: responses,
        }),
    ];
}
