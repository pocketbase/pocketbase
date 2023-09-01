<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";

    export let collection;

    let responseTab = 200;
    let responses = [];
    let baseData = {};

    $: isAuth = collection.type === "auth";

    $: adminsOnly = collection?.createRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(CommonHelper.dummyCollectionRecord(collection), null, 2),
        },
        {
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${collection?.schema?.[0]?.name}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
        {
            code: 403,
            body: `
                {
                  "code": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `,
        },
    ];

    $: if (isAuth) {
        baseData = {
            username: "test_username",
            email: "test@example.com",
            emailVisibility: true,
            password: "12345678",
            passwordConfirm: "12345678",
        };
    } else {
        baseData = {};
    }
</script>

<h3 class="m-b-sm">Create ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>Create a new <strong>{collection.name}</strong> record.</p>
    <p>
        Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.
    </p>
    <p>
        File upload is supported only via <code>multipart/form-data</code>.
        <br />
        For more info and examples you could check the detailed
        <a href={import.meta.env.PB_FILE_UPLOAD_DOCS} target="_blank" rel="noopener noreferrer">
            Files upload and handling docs
        </a>.
    </p>
</div>

<!-- prettier-ignore -->
<SdkTabs
    js={`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${backendAbsUrl}');

...

// example create data
const data = ${JSON.stringify(Object.assign({}, baseData, CommonHelper.dummyCollectionSchemaData(collection)), null, 4)};

const record = await pb.collection('${collection?.name}').create(data);
` + (isAuth ?
`
// (optional) send an email verification request
await pb.collection('${collection?.name}').requestVerification('test@example.com');
` : ""
)}
    dart={`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${backendAbsUrl}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({}, baseData, CommonHelper.dummyCollectionSchemaData(collection)), null, 2)};

final record = await pb.collection('${collection?.name}').create(body: body);
` + (isAuth ?
`
// (optional) send an email verification request
await pb.collection('${collection?.name}').requestVerification('test@example.com');
` : ""
)}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">Requires admin <code>Authorization:TOKEN</code> header</p>
    {/if}
</div>

<div class="section-title">Body Parameters</div>
<table class="table-compact table-border m-b-base">
    <thead>
        <tr>
            <th>Param</th>
            <th>Type</th>
            <th width="50%">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-warning">Optional</span>
                    <span>id</span>
                </div>
            </td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                <strong>15 characters string</strong> to store as record ID.
                <br />
                If not set, it will be auto generated.
            </td>
        </tr>

        {#if isAuth}
            <tr>
                <td colspan="3" class="txt-hint">Auth fields</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">Optional</span>
                        <span>username</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>
                    The username of the auth record.
                    <br />
                    If not set, it will be auto generated.
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if collection?.options?.requireEmail}
                            <span class="label label-success">Required</span>
                        {:else}
                            <span class="label label-warning">Optional</span>
                        {/if}
                        <span>email</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>Auth record email address.</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">Optional</span>
                        <span>emailVisibility</span>
                    </div>
                </td>
                <td>
                    <span class="label">Boolean</span>
                </td>
                <td>Whether to show/hide the auth record email when fetching the record data.</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-success">Required</span>
                        <span>password</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>Auth record password.</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-success">Required</span>
                        <span>passwordConfirm</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>Auth record password confirmation.</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">Optional</span>
                        <span>verified</span>
                    </div>
                </td>
                <td>
                    <span class="label">Boolean</span>
                </td>
                <td>
                    Indicates whether the auth record is verified or not.
                    <br />
                    This field can be set only by admins or auth records with "Manage" access.
                </td>
            </tr>
            <tr>
                <td colspan="3" class="txt-hint">Schema fields</td>
            </tr>
        {/if}

        {#each collection?.schema as field (field.name)}
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if field.required}
                            <span class="label label-success">Required</span>
                        {:else}
                            <span class="label label-warning">Optional</span>
                        {/if}
                        <span>{field.name}</span>
                    </div>
                </td>
                <td>
                    <span class="label">{CommonHelper.getFieldValueType(field)}</span>
                </td>
                <td>
                    {#if field.type === "text"}
                        Plain text value.
                    {:else if field.type === "number"}
                        Number value.
                    {:else if field.type === "json"}
                        JSON array or object.
                    {:else if field.type === "email"}
                        Email address.
                    {:else if field.type === "url"}
                        URL address.
                    {:else if field.type === "file"}
                        File object.<br />
                        Set to <code>null</code> to delete already uploaded file(s).
                    {:else if field.type === "relation"}
                        Relation record {field.options?.maxSelect === 1 ? "id" : "ids"}.
                    {/if}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

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
            <td>expand</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                Auto expand relations when returning the created record. Ex.:
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                Supports up to 6-levels depth nested relations expansion. <br />
                The expanded relations will be appended to the record under the
                <code>expand</code> property (eg. <code>{`"expand": {"relField1": {...}, ...}`}</code>).
                <br />
                Only the relations to which the request user has permissions to <strong>view</strong> will be expanded.
            </td>
        </tr>
        <FieldsQueryParam />
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
