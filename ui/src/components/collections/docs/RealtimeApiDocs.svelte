<script>
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";

    export let collection = new Collection();

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);
</script>

<div class="alert">
    <strong class="label label-primary">SSE</strong>
    <div class="content">
        <p>/api/realtime</p>
    </div>
</div>

<div class="content m-b-base">
    <p>Subscribe to realtime changes via Server-Sent Events (SSE).</p>
    <p>
        Events are send for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see "Event data format" section below).
    </p>
    <div class="alert alert-info m-t-10">
        <div class="icon">
            <i class="ri-information-line" />
        </div>
        <div class="contet">
            <p>
                <strong>You could subscribe to a single record or to an entire collection.</strong>
            </p>
            <p>
                When you subscribe to a <strong>single record</strong>, the collection's
                <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive
                the event message.
            </p>
            <p>
                When you subscribe to an <strong>entire collection</strong>, the collection's
                <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive
                the event message.
            </p>
        </div>
    </div>
</div>

<div class="section-title">Client SDKs example</div>
<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const client = new PocketBase('${backendAbsUrl}');

        ...

        // (Optionally) authenticate
        client.users.authViaEmail('test@example.com', '123456');

        // Subscribe to changes in any record from the collection
        client.realtime.subscribe('${collection?.name}', function (e) {
            console.log(e.record);
        });

        // Subscribe to changes in a single record
        client.realtime.subscribe('${collection?.name}/RECORD_ID', function (e) {
            console.log(e.record);
        });

        // Unsubscribe
        client.realtime.unsubscribe() // remove all subscriptions
        client.realtime.unsubscribe('${collection?.name}') // remove only the collection subscription
        client.realtime.unsubscribe('${collection?.name}/RECORD_ID') // remove only the record subscription
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final client = PocketBase('${backendAbsUrl}');

        ...

        // (Optionally) authenticate
        client.users.authViaEmail('test@example.com', '123456');

        // Subscribe to changes in any record from the collection
        client.realtime.subscribe('${collection?.name}', (e) {
          print(e.record);
        });

        // Subscribe to changes in a single record
        client.realtime.subscribe('${collection?.name}/RECORD_ID', (e) {
          print(e.record);
        });

        // Unsubscribe
        client.realtime.unsubscribe() // remove all subscriptions
        client.realtime.unsubscribe('${collection?.name}') // remove only the collection subscription
        client.realtime.unsubscribe('${collection?.name}/RECORD_ID') // remove only the record subscription
    `}
/>

<div class="section-title">Event data format</div>
<CodeBlock
    content={JSON.stringify(
        {
            action: "create",
            record: CommonHelper.dummyCollectionRecord(collection),
        },
        null,
        2
    ).replace('"action": "create"', '"action": "create" // create, update or delete')}
/>
