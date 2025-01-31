<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";

    export let collection;

    let responseTab = 204;
    let responses = [];

    $: superusersOnly = collection?.deleteRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: if (collection?.id) {
        responses.push({
            code: 204,
            body: `
                null
            `,
        });

        responses.push({
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `,
        });

        if (superusersOnly) {
            responses.push({
                code: 403,
                body: `
                    {
                      "status": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `,
            });
        }

        responses.push({
            code: 404,
            body: `
                {
                  "status": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `,
        });
    }
</script>

<h3 class="m-b-sm">Delete ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>Delete a single <strong>{collection.name}</strong> record.</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').delete('RECORD_ID');
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').delete('RECORD_ID');
    `}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-danger">
    <strong class="label label-primary">DELETE</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong>
        </p>
    </div>
    {#if superusersOnly}
        <p class="txt-hint txt-sm txt-right">Requires superuser <code>Authorization:TOKEN</code> header</p>
    {/if}
</div>

<div class="section-title">Path parameters</div>
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
            <td>ID of the record to delete.</td>
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
