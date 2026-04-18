export function docsRealtime(collection) {
    const baseURL = app.utils.getApiExampleURL();

    const dummyRecord = Object.assign({
        collectionId: collection.id,
        collectionName: collection.name,
    }, app.utils.getDummyFieldsData(collection));

    return t.div(
        { pbEvent: "apiPreviewRealtime", className: "content" },
        // description
        t.p(null, `Subscribe to realtime changes via Server-Sent Events (SSE).`),
        t.p(
            null,
            "Events are sent for ",
            t.strong(null, "create"),
            ", ",
            t.strong(null, "update"),
            " and ",
            t.strong(null, "delete"),
            ` record operations (see "Event data format" below).`,
        ),
        t.div(
            { className: "alert info" },
            t.p({ className: "txt-bold" }, "You could subscribe to a single record or to an entire collection."),
            t.p(
                null,
                "When you subscribe to a ",
                t.strong(null, "single record"),
                ", the collection's ",
                t.strong(null, "View rule"),
                " will be used to determine whether the subscriber is allowed to receive the event message.",
            ),
            t.p(
                null,
                "When you subscribe to an ",
                t.strong(null, "entire collection"),
                ", the collection's ",
                t.strong(null, "List/Search rule"),
                " will be used to determine whether the subscriber is allowed to receive the event message.",
            ),
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

                        // (optionally) authenticate
                        await pb.collection('users').authWithPassword('test@example.com', '123456');

                        // subscribe to changes in any ${baseURL} record
                        pb.collection('${baseURL}').subscribe('*', function (e) {
                            console.log(e.action);
                            console.log(e.record);
                        }, { /* other options like: filter, expand, custom headers, etc. */ });

                        // subscribe to changes only in the specified record
                        pb.collection('${baseURL}').subscribe('RECORD_ID', function (e) {
                            console.log(e.action);
                            console.log(e.record);
                        }, { /* other options like: filter, expand, custom headers, etc. */ });

                        ...

                        // unsubscribe - remove all 'RECORD_ID' subscriptions
                        pb.collection('${baseURL}').unsubscribe('RECORD_ID');

                        // unsubscribe - remove all '*' topic subscriptions
                        pb.collection('${baseURL}').unsubscribe('*');

                        // unsubscribe - remove all collection subscriptions
                        pb.collection('${baseURL}').unsubscribe();
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

                        // (optionally) authenticate
                        await pb.collection('users').authWithPassword('test@example.com', '123456');

                        // subscribe to changes in any ${baseURL} record
                        pb.collection('${baseURL}').subscribe('*', (e) {
                            print(e.action);
                            print(e.record);
                        }, /* other options like: filter, expand, custom headers, etc. */);

                        // subscribe to changes only in the specified record
                        pb.collection('${baseURL}').subscribe('RECORD_ID', (e) {
                            print(e.action);
                            print(e.record);
                        }, /* other options like: filter, expand, custom headers, etc. */);

                        ...

                        // unsubscribe - remove all 'RECORD_ID' subscriptions
                        pb.collection('${baseURL}').unsubscribe('RECORD_ID');

                        // unsubscribe - remove all '*' topic subscriptions
                        pb.collection('${baseURL}').unsubscribe('*');

                        // unsubscribe - remove all collection subscriptions
                        pb.collection('${baseURL}').unsubscribe();
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
                        # init an SSE connection and start listening for messages
                        # (the first message is always PB_CONNECT with the connection "clientId")
                        curl -N '${baseURL}/api/realtime'

                        # open a new terminal and submit the subscription topic(s)
                        # with the "clientId" from the initial PB_CONNECT message
                        curl -X POST \\
                          -H 'Authorization:TOKEN' \\
                          -H 'Content-Type:application/json' \\
                          -d '{ "clientId": "YOUR_CLIENT_ID", "subscriptions": ["${collection.name}/*"] }' \\
                          '${baseURL}/api/realtime'

                        # create/update/delete a record in the ${collection.name} collection and
                        # you should see the event message(s) in the first terminal
                        # (as long as your client satisfies the topic API rule)
                    `,
                },
            ],
        }),
        // api
        t.div({ className: "block m-t-base" }, t.strong(null, "API details")),
        t.div(
            { className: "alert api-preview-alert" },
            t.span({ className: "label method" }, "GET/POST"),
            t.span({ className: "path" }, "/api/realtime"),
            t.div(
                { className: "extra" },
                t.a({
                    href: import.meta.env.PB_REALTIME_DOCS,
                    target: "_blank",
                    rel: "noopener noreferrer",
                    textContent: "Realtime docs",
                }),
            ),
        ),
        t.div({ className: "block m-t-base m-b-sm" }, t.strong(null, "Event data format")),
        app.components.codeBlock({
            value: JSON.stringify(
                {
                    "action": "create",
                    "record": dummyRecord,
                },
                null,
                2,
            ).replace(`"action": "create",`, "\"action\": \"create\", // create, update or delete"),
        }),
    );
}
