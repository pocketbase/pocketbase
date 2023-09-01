<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";

    export let collection;

    let responseTab = 204;
    let responses = [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: responses = [
        {
            code: 204,
            body: "null",
        },
        {
            code: 401,
            body: `
                {
                  "code": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `,
        },
        {
            code: 403,
            body: `
                {
                  "code": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `,
        },
        {
            code: 404,
            body: `
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">Unlink OAuth2 account ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>
        Unlink a single external OAuth2 provider from <strong>{collection.name}</strong> record.
    </p>
    <p>Only admins and the account owner can access this action.</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${collection?.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${collection?.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-danger">
    <strong class="label label-primary">DELETE</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong
            >/external-auths/<strong>:provider</strong>
        </p>
    </div>
    <p class="txt-hint txt-sm txt-right">Requires <code>Authorization:TOKEN</code> header</p>
</div>

<div class="section-title">Path Parameters</div>
<table class="table-compact table-border m-b-base">
    <thead>
        <tr>
            <th>Param</th>
            <th>Type</th>
            <th width="60%">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>id</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>ID of the auth record.</td>
        </tr>
        <tr>
            <td>provider</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                The name of the auth provider to unlink, eg. <code>google</code>, <code>twitter</code>,
                <code>github</code>, etc.
            </td>
        </tr>
    </tbody>
</table>

<div class="section-title">Responses</div>
<div class="tabs">
    <div class="tabs-header compact combined left">
        {#each responses as response (response.code)}
            <button
                class="tab-item"
                class:active={responseTab === response.code}
                on:click={() => (responseTab = response.code)}
            >
                {response.code}
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each responses as response (response.code)}
            <div class="tab-item" class:active={responseTab === response.code}>
                <CodeBlock content={response.body} />
            </div>
        {/each}
    </div>
</div>
