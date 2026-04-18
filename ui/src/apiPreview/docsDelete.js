export function docsDelete(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const isSuperusersOnly = collection.deleteRule === null;

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
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `,
        },
    ];
    if (isSuperusersOnly) {
        responses.push({
            title: 403,
            value: `
                {
                  "status": 403,
                  "message": "Only superusers can access this action.",
                  "data": {}
                }
            `,
        });
    }
    responses.push({
        title: 404,
        value: `
            {
              "status": 404,
              "message": "The requested resource wasn't found.",
              "data": {}
            }
        `,
    });

    return t.div(
        { pbEvent: "apiPreviewDelete", className: "content" },
        // description
        t.p(null, `Delete a single ${collection.name} record.`),
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

                        await pb.collection('${collection.name}').delete('RECORD_ID');
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

                        await pb.collection('${collection.name}').delete('RECORD_ID');
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
                        curl -X DELETE \\
                          -H 'Authorization:TOKEN' \\
                          '${baseURL}/api/collections/${collection.name}/records/RECORD_ID'
                    `,
                },
            ],
        }),
        // api
        t.div({ className: "block m-t-base" }, t.strong(null, "API details")),
        t.div(
            { className: "alert danger api-preview-alert" },
            t.span({ className: "label method" }, "DELETE"),
            t.span({ className: "path" }, `/api/collections/${collection.name}/records/`, t.strong(null, ":id")),
            () => {
                if (isSuperusersOnly) {
                    return t.small({ className: "extra" }, "Requires superuser Authorization:TOKEN header");
                }
            },
        ),
        t.table(
            { className: "api-preview-table path-params" },
            t.thead(
                null,
                t.tr(
                    null,
                    t.th({ className: "min-width txt-primary" }, "Path params"),
                    t.th({ className: "min-width" }, "Type"),
                    t.th(null, "Description"),
                ),
            ),
            t.tbody(
                null,
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "id"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(null, "ID of the record to delete."),
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
