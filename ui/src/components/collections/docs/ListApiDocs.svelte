<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FilterSyntax from "@/components/collections/docs/FilterSyntax.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: fieldNames = CommonHelper.getAllCollectionIdentifiers(collection);

    $: superusersOnly = collection?.listRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: dummyRecord = CommonHelper.dummyCollectionRecord(collection);

    $: if (collection?.id) {
        responses.push({
            code: 200,
            body: JSON.stringify(
                {
                    page: 1,
                    perPage: 30,
                    totalPages: 1,
                    totalItems: 2,
                    items: [dummyRecord, Object.assign({}, dummyRecord, { id: dummyRecord.id + "2" })],
                },
                null,
                2,
            ),
        });

        responses.push({
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
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
            filter: 'someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${collection?.name}').getFullList({
            sort: '-someField',
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
          filter: 'someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${collection?.name}').getFullList(
          sort: '-someField',
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
    {#if superusersOnly}
        <p class="txt-hint txt-sm txt-right">Requires superuser <code>Authorization:TOKEN</code> header</p>
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
                    <code>@rowid</code>,
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
        <FieldsQueryParam />
        <tr>
            <td id="query-page">skipTotal</td>
            <td>
                <span class="label">Boolean</span>
            </td>
            <td>
                If it is set the total counts query will be skipped and the response fields
                <code>totalItems</code> and <code>totalPages</code> will have <code>-1</code> value.
                <br />
                This could drastically speed up the search queries when the total counters are not needed or cursor
                based pagination is used.
                <br />
                For optimization purposes, it is set by default for the
                <code>getFirstListItem()</code>
                and
                <code>getFullList()</code> SDKs methods.
            </td>
        </tr>
    </tbody>
</table>

<div class="section-title">Responses</div>
<div class="tabs">
    <div class="tabs-header compact combined left">
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
