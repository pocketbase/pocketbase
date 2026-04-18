import { expandInfo } from "./expandInfo";
import { fieldsInfo } from "./fieldsInfo";

export function docsView(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const isSuperusersOnly = collection.viewRule === null;

    const baseDummyRecord = {
        collectionId: collection.id,
        collectionName: collection.name,
    };

    const responses = [
        {
            title: 200,
            value: JSON.stringify(
                Object.assign(baseDummyRecord, app.utils.getDummyFieldsData(collection)),
                null,
                2,
            ),
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
        { pbEvent: "apiPreviewView", className: "content" },
        // description
        t.p(null, `Fetch a single ${collection.name} record.`),
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

                        const record = await pb.collection('${collection.name}').getOne('RECORD_ID', {
                            expand: 'relField1,relField2.subRelField',
                        });
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

                        final record = await pb.collection('${collection.name}').getOne('RECORD_ID',
                          expand: 'relField1,relField2.subRelField',
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
                        curl \\
                          -H 'Authorization:TOKEN' \\
                          '${baseURL}/api/collections/${collection.name}/records/RECORD_ID'
                    `,
                },
            ],
        }),
        // api
        t.div({ className: "block m-t-base" }, t.strong(null, "API details")),
        t.div(
            { className: "alert info api-preview-alert" },
            t.span({ className: "label method" }, "GET"),
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
                    t.td(null, "ID of the record to view."),
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
