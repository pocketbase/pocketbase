<script>
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CodeBlock from "@/components/base/CodeBlock.svelte";

    export let collection = new Collection();

    let responseTab = 204;
    let sdkTab = "JavaScript";
    let responses = [];
    let sdkExamples = [];

    $: adminsOnly = collection?.deleteRule === null;

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
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
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

                ...

                await client.Records.delete("${collection?.name}", "RECORD_ID");
            `,
        },
    ];
</script>

<div class="alert alert-danger">
    <strong class="label label-primary">DELETE</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong>
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">Requires <code>Authorization: Admin TOKEN</code> header</p>
    {/if}
</div>

<div class="content m-b-base">
    <p>Delete a single <strong>{collection.name}</strong> record.</p>
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

<div class="section-title">Path parameters</div>
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
    <div class="tabs-header compact left">
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
