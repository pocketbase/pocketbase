import { expandInfo } from "./expandInfo";
import { fieldsInfo } from "./fieldsInfo";
import { filterSyntax } from "./filterSyntax";

export function docsList(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const isSuperusersOnly = collection.listRule === null;

    const baseDummyRecord = {
        collectionId: collection.id,
        collectionName: collection.name,
    };

    const responses = [
        {
            title: 200,
            value: JSON.stringify(
                {
                    page: 1,
                    perPage: 30,
                    totalPages: 1,
                    totalItems: 2,
                    items: [
                        Object.assign(baseDummyRecord, app.utils.getDummyFieldsData(collection)),
                        Object.assign(baseDummyRecord, app.utils.getDummyFieldsData(collection)),
                    ],
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
                  "message": "Something went wrong while processing your request.",
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

    return t.div(
        { pbEvent: "apiPreviewList", className: "content" },
        // description
        t.p(null, `Fetch a paginated ${collection.name} records list, supporting sorting and filtering.`),
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

                        // fetch a paginated records list
                        const resultList = await pb.collection('${collection.name}').getList(1, 50, {
                          filter: 'someField1 != someField2',
                        });

                        // you can also fetch all records at once via getFullList
                        const records = await pb.collection('${collection.name}').getFullList({
                          sort: '-someField',
                        });

                        // or fetch only the first record that matches the specified filter
                        const record = await pb.collection('${collection.name}').getFirstListItem(
                          'someField="test"',
                          { expand: 'relField1,relField2.subRelField' },
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

                        // fetch a paginated records list
                        final resultList = await pb.collection('${collection.name}').getList(
                          page: 1,
                          perPage: 50,
                          filter: 'someField1 != someField2',
                        );

                        // you can also fetch all records at once via getFullList
                        final records = await pb.collection('${collection.name}').getFullList(
                          sort: '-someField',
                        );

                        // or fetch only the first record that matches the specified filter
                        final record = await pb.collection('${collection.name}').getFirstListItem(
                          'someField="test"',
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
                          '${baseURL}/api/collections/${collection.name}/records?perPage=50'
                    `,
                },
            ],
        }),
        // api
        t.div({ className: "block m-t-base" }, t.strong(null, "API details")),
        t.div(
            { className: "alert info api-preview-alert" },
            t.span({ className: "label method" }, "GET"),
            t.span({ className: "path" }, `/api/collections/${collection.name}/records`),
            () => {
                if (isSuperusersOnly) {
                    return t.small({ className: "extra" }, "Requires superuser Authorization:TOKEN header");
                }
            },
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
                    t.td({ className: "min-width" }, "page"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "Number")),
                    t.td(null, "The page (aka. offset) of the paginated list (default to 1)."),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "perPage"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "Number")),
                    t.td(null, "Specify the max returned records per page (default to 30)."),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "sort"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(
                        null,
                        t.p(
                            null,
                            "Specify the records order attribute(s).",
                            t.br(),
                            "Add -/+ (default) in front of the attribute for DESC / ASC order.",
                        ),
                        t.p(
                            null,
                            "For example:",
                            app.components.codeBlock({
                                value: `// DESC by created and ASC by id\n?sort=-created,id`,
                            }),
                        ),
                        t.p(
                            null,
                            "In addition to the collection non-hidden fields, the following special sort fields could be also used: ",
                            t.code(null, "@random"),
                            " ",
                            t.code({ hidden: () => collection.type == "view" }, "@rowid"),
                            ".",
                        ),
                    ),
                ),
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "filter"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "String")),
                    t.td(
                        null,
                        t.p(null, "Filter the returned records. For example:"),
                        app.components.codeBlock({
                            value: `?filter=(id='abc' && created>'2022-01-01')`,
                            footnote: "All query params must be properly URL encoded (the SDKs do this automatically).",
                        }),
                        filterSyntax(),
                    ),
                ),
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
                t.tr(
                    null,
                    t.td({ className: "min-width" }, "skipTotal"),
                    t.td({ className: "min-width" }, t.span({ className: "label" }, "Boolean")),
                    t.td(
                        null,
                        t.p(
                            null,
                            "If set to ",
                            t.code(null, "1/true"),
                            " the total counts query will be skipped and the response fields ",
                            t.code(null, "totalItems"),
                            " and ",
                            t.code(null, "totalPages"),
                            " will have -1 value.",
                        ),
                        t.p(
                            null,
                            "This could drastically speed up the search queries when the total counters are not needed or cursor based pagination is used.",
                            " For optimization purposes, it is set by default in the ",
                            t.code(null, "getFirstListItem()"),
                            " and ",
                            t.code(null, "getFullList()"),
                            " SDKs methods.",
                        ),
                    ),
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
