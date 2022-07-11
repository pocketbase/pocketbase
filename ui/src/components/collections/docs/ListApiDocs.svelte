<script>
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FilterSyntax from "@/components/collections/docs/FilterSyntax.svelte";

    export let collection = new Collection();

    let responseTab = 200;
    let sdkTab = "JavaScript";
    let responses = [];
    let sdkExamples = [];

    $: adminsOnly = collection?.listRule === null;

    $: if (collection?.id) {
        responses.push({
            code: 200,
            body: JSON.stringify(
                {
                    page: 1,
                    perPage: 30,
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

        responses.push({
            code: 404,
            body: `
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `,
        });
    }

    $: sdkExamples = [
        {
            lang: "JavaScript",
            code: `
                import PocketBase from 'pocketbase';

                const client = new PocketBase("${ApiClient.baseUrl}");

                client.Records.getList("${collection?.name}", { page: 2 })
                    .then(function (list) {
                        // success...
                    }).catch(function (error) {
                        // error...
                    });

                // alternatively you can also fetch all records at once via getFullList:
                client.Records.getFullList("${collection?.name}", 200 /* batch size */)
                    .then(function (records) {
                        // success...
                    }).catch(function (error) {
                        // error...
                    });
            `,
        },
    ];
</script>

<div class="alert alert-info">
    <strong class="label label-primary">GET</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">Requires <code>Authorization: Admin TOKEN</code> header</p>
    {/if}
</div>

<div class="content m-b-base">
    <p>Fetch a paginated <strong>{collection.name}</strong> records list.</p>
</div>

<div class="section-title">Client SDKs example</div>
<div class="tabs m-b-lg">
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

<div class="section-title">Query parameters</div>
<table class="table-compact table-border m-b-lg">
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
                Auto expand nested record relations. Ex.:
                <CodeBlock
                    content={`
                        ?expand=rel1,rel2.subrel21.subrel22
                    `}
                />
                Supports up to 6-levels depth nested relations expansion. <br />
                The expanded relations will be appended to each individual record under the
                <code>@expand</code> property (eg. <code>{`"@expand": {"rel1": {...}, ...}`}</code>). Only the
                relations that the user has permissions to <strong>view</strong> will be expanded.
            </td>
        </tr>
    </tbody>
</table>

<div class="section-title">Responses</div>
<div class="tabs">
    <div class="tabs-header compact left">
        {#each responses as response (response.code)}
            <div
                class="tab-item"
                class:active={responseTab === response.code}
                on:click={() => (responseTab = response.code)}
            >
                {response.code}
            </div>
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
