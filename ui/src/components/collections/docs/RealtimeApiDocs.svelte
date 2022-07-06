<script>
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";

    export let collection = new Collection();

    let sdkTab = "JavaScript";
    let sdkExamples = [];

    $: sdkExamples = [
        {
            lang: "JavaScript",
            code: `
                import PocketBase from 'pocketbase';

                const client = new PocketBase("${ApiClient.baseUrl}");

                // (Optionally) authenticate
                client.Users.authViaEmail("test@example.com", "123456");

                // Subscribe to changes in any record from the collection
                client.Realtime.subscribe("${collection?.name}", function (e) {
                    console.log(e.data);
                });

                // Subscribe to changes in a single record
                client.Realtime.subscribe("${collection?.name}/RECORD_ID", function (e) {
                    console.log(e.data);
                });

                // Unsubscribe
                client.Realtime.unsubscribe() // remove all subscriptions
                client.Realtime.unsubscribe("${collection?.name}") // remove the collection subscription
                client.Realtime.unsubscribe("${collection?.name}/RECORD_ID") // remove the record subscription
            `,
        },
    ];
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
<div class="tabs m-b-base">
    <div class="tabs-header compact left">
        {#each sdkExamples as example (example.lang)}
            <button
                class="tab-item"
                class:active={sdkTab === example.lang}
                on:click={() => (sdkTab = example.lang)}
            >
                {example.lang}
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each sdkExamples as example (example.lang)}
            <div class="tab-item" class:active={sdkTab === example.lang}>
                <CodeBlock content={example.code} />
            </div>
        {/each}
    </div>
</div>

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
