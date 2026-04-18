import { expandInfo } from "./expandInfo";
import { fieldsInfo } from "./fieldsInfo";

export function docsAuthWithPassword(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const identityFields = collection.passwordAuth?.identityFields || [];

    const exampleIdentityLabel = identityFields.length == 0
        ? "NONE"
        : "YOUR_" + identityFields.join("_OR_").toUpperCase();

    const baseDummyRecord = {
        collectionId: collection.id,
        collectionName: collection.name,
    };

    const responses = [
        {
            title: 200,
            value: JSON.stringify(
                {
                    token: "...JWT...",
                    record: Object.assign(baseDummyRecord, app.utils.getDummyFieldsData(collection)),
                },
                null,
                2,
            ),
        },
        {
            title: 400,
            value: `
                {
                  "status": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "identity": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
    ];

    return t.div(
        {
            pbEvent: "apiPreviewAuthWithPassword",
            className: "content",
        },
        // description
        t.p(
            null,
            "Authenticate with combination of ",
            t.strong(null, identityFields.join("/")),
            " and ",
            t.strong(null, "password"),
            ".",
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

                        const authData = await pb.collection('${collection.name}').authWithPassword(
                          '${exampleIdentityLabel}',
                          'YOUR_PASSWORD',
                        );

                        // after the above you can also access the auth data from the authStore
                        console.log(pb.authStore.isValid);
                        console.log(pb.authStore.token);
                        console.log(pb.authStore.record.id);

                        // "logout"
                        pb.authStore.clear();
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

                        final authData = await pb.collection('${collection.name}').authWithPassword(
                          '${exampleIdentityLabel}',
                          'YOUR_PASSWORD',
                        );

                        // after the above you can also access the auth data from the authStore
                        print(pb.authStore.isValid);
                        print(pb.authStore.token);
                        print(pb.authStore.record.id);

                        // "logout"
                        pb.authStore.clear();
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
                        curl -X POST \\
                          -H 'Content-Type:application/json' \\
                          -d '{ "identity":"${exampleIdentityLabel}", "password":"YOUR_PASSWORD" }' \\
                          '${baseURL}/api/collections/${collection.name}/auth-with-password'
                    `,
                },
            ],
        }),
        // api
        t.div({ className: "block m-t-base" }, t.strong(null, "API details")),
        t.div(
            { className: "alert success api-preview-alert" },
            t.span({ className: "label method" }, "POST"),
            t.span({ className: "path" }, `/api/collections/${collection.name}/auth-with-password`),
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
                    t.td({ className: "min-width" }, "identity ", t.em(null, "(required)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(
                        null,
                        app.utils.sentenize(identityFields.join(" or "), false),
                        " of the record to authenticate.",
                    ),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "identityField ", t.em(null, "(optional)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(
                        null,
                        "In case of multiple identity fields, explicitly set the field name to use when searching for the auth record.",
                        t.br(),
                        "Leave it empty for auto detection.",
                    ),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "password ", t.em(null, "(required)")),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, "The auth record password."),
                ),
            ),
        ),
        t.table(
            { className: "api-preview-table query-params" },
            t.thead(
                null,
                t.tr(
                    null,
                    t.th({ className: "min-width txt-primary" }, "?query params"),
                    t.th({ className: "min-width" }, "Type"),
                    t.th(null, "Description"),
                ),
            ),
            t.tbody(
                null,
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "expand"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, expandInfo()),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "fields"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, fieldsInfo()),
                ),
            ),
        ),
        // responses
        t.div({ className: "block m-t-base m-b-sm" }, t.strong(null, "Example responses")),
        app.components.codeBlockTabs({
            tabs: responses,
        }),
    );
}
