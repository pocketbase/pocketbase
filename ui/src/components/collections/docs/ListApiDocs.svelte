<script>
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FilterSyntax from "@/components/collections/docs/FilterSyntax.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";

    export let collection = new Collection();

    let responseTab = 200;
    let responses = [];

    $: fieldNames = CommonHelper.getAllCollectionIdentifiers(collection);

    $: adminsOnly = collection?.listRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: if (collection?.id) {
        responses.push({
            code: 200,
            body: JSON.stringify(
                {
                    page: 1,
                    perPage: 30,
                    totalPages: 1,
                    totalItems: 2,
                    items: [
                        CommonHelper.dummyCollectionRecord(collection),
                        CommonHelper.dummyCollectionRecord(collection),
                    ],
                },
                null,
                2
            ),
        });

        responses.push({
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `,
        });

        if (adminsOnly) {
            responses.push({
                code: 403,
                body: `
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `,
            });
        }
    }
</script>

<h3 class="m-b-sm">List/Search ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>
        Fetch a paginated <strong>{collection.name}</strong> records list, supporting sorting and filtering.
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${collection?.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${collection?.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${collection?.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${collection?.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${collection?.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${collection?.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-info">
    <strong class="label label-primary">GET</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">Requires admin <code>Authorization:TOKEN</code> header</p>
    {/if}
</div>

<div class="section-title">Query parameters</div>
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
            <td>page</td>
            <td>
                <span class="label">Number</span>
            </td>
            <td>The page (aka. offset) of the paginated list (default to 1).</td>
        </tr>
        <tr>
            <td>perPage</td>
            <td>
                <span class="label">Number</span>
            </td>
            <td>Specify the max returned records per page (default to 30).</td>
        </tr>
        <tr>
            <td>sort</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                Specify the records order attribute(s). <br />
                Add <code>-</code> / <code>+</code> (default) in front of the attribute for DESC / ASC order.
                Ex.:
                <CodeBlock
                    content={`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}
                />
                <p>
                    <strong>Supported record sort fields:</strong> <br />
                    <code>@random</code>,
                    {#each fieldNames as name, i}
                        <code>{name}</code>{i < fieldNames.length - 1 ? ", " : ""}
                    {/each}
                </p>
            </td>
        </tr>
        <tr>
            <td>filter</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                Filter the returned records. Ex.:
                <CodeBlock
                    content={`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}
                />
                <FilterSyntax />
            </td>
        </tr>
        <tr>
            <td>expand</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                Auto expand record relations. Ex.:
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                Supports up to 6-levels depth nested relations expansion. <br />
                The expanded relations will be appended to each individual record under the
                <code>expand</code> property (eg. <code>{`"expand": {"relField1": {...}, ...}`}</code>).
                <br />
                Only the relations to which the request user has permissions to <strong>view</strong> will be expanded.
            </td>
        </tr>
        <tr>
            <td id="query-page">fields</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                Comma separated string of the fields to return in the JSON response
                <em>(by default returns all fields)</em>.
            </td>
        </tr>
    </tbody>
</table>

<div class="section-title">Responses</div>
<div class="tabs">
    <div class="tabs-header compact left">
        {#each responses as response (response.code)}
            <button
                type="button"
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
