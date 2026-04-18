import { fullDummyPayload, primitivesDummyPayload, replaceDummyPayloadPlaceholder } from "./docsCreate";
import { expandInfo } from "./expandInfo";
import { fieldsInfo } from "./fieldsInfo";

export function docsUpdate(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const isSuperusersOnly = collection.updateRule === null;

    const isAuth = collection.type === "auth";

    const excludedTableFields = isAuth ? ["id", "password", "verified", "email", "emailVisibility"] : ["id"];

    const tableFields =
        collection.fields?.filter((f) => !f.hidden && f.type != "autodate" && !excludedTableFields.includes(f.name))
        || [];

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
        {
            title: 400,
            value: `
                {
                  "status": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${tableFields.find((f) => !f.primaryKey)?.name || "someField"}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
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
                  "message": "Only superusers can perform this action.",
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
        { pbEvent: "apiPreviewUpdate", className: "content" },
        // description
        t.p(null, `Updates an existing ${collection.name} record.`),
        t.p(
            null,
            "Body parameters could be sent as ",
            t.code(null, "application/json"),
            " or ",
            t.code(null, "multipart/form-data"),
            ".",
        ),
        t.p(
            null,
            "File upload is supported only via ",
            t.code(null, "multipart/form-data"),
            ". For more info and examples you could check the detailed ",
            t.a({
                href: import.meta.env.PB_FILE_UPLOAD_DOCS,
                target: "_blank",
                rel: "noopener noreferrer",
                textContent: "Files upload and handling docs",
            }),
            ".",
        ),
        t.p(
            null,
            t.em(
                null,
                "Note that in case of a password change all previously issued tokens for the current record will be automatically invalidated and if you want your user to remain signed in you need to reauthenticate manually after the update call.",
            ),
        ),
        app.components.codeBlockTabs({
            className: "sdk-examples m-t-sm",
            historyKey: "pbLastSDK",
            tabs: [
                {
                    title: "JS SDK",
                    language: "js",
                    // dprint-ignore
                    value: `
import PocketBase from 'pocketbase';

const pb = new PocketBase('${baseURL}');

...

// example update body
const body = ${replaceDummyPayloadPlaceholder(JSON.stringify(fullDummyPayload(collection, true), null, 2))};

const record = await pb.collection('${collection.name}').update('RECORD_ID', body);
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
                    // dprint-ignore
                    value: `
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${baseURL}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(primitivesDummyPayload(collection, true), null, 2)};

final record = await pb.collection('${collection.name}').update(
  'RECORD_ID',
  body: body,
  files: [],
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
                        curl -X PATCH \\
                          -H 'Authorization:TOKEN' \\
                          -H 'Content-Type:application/json' \\
                          -d '{ ... }' \\
                          '${baseURL}/api/collections/${collection.name}/records/RECORD_ID'
                    `,
                },
            ],
        }),
        // api
        t.div({ className: "block m-t-base" }, t.strong(null, "API details")),
        t.div(
            { className: "alert warning api-preview-alert" },
            t.span({ className: "label method" }, "PATCH"),
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
                    t.td(null, "ID of the record to update."),
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
