<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: isAuth = collection.type === "auth";

    $: superusersOnly = collection?.createRule === null;

    $: excludedTableFields = isAuth ? ["password", "verified", "email", "emailVisibility"] : [];

    $: tableFields =
        collection?.fields?.filter((f) => {
            return !f.hidden && f.type != "autodate" && !excludedTableFields.includes(f.name);
        }) || [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(CommonHelper.dummyCollectionRecord(collection), null, 2),
        },
        {
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${collection?.fields?.[0]?.name}": {
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
                  "status": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `,
        },
    ];

    function getPayload(collection) {
        let payload = CommonHelper.dummyCollectionSchemaData(collection, true);

        if (isAuth) {
            payload.password = "12345678";
            payload.passwordConfirm = "12345678";
            delete payload.verified;
        }

        return payload;
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
const data = ${JSON.stringify(getPayload(collection), null, 4)};

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
final body = <String, dynamic>${JSON.stringify(getPayload(collection), null, 2)};

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
    {#if superusersOnly}
        <p class="txt-hint txt-sm txt-right">Requires superuser <code>Authorization:TOKEN</code> header</p>
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
        {#if isAuth}
            <tr>
                <td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if collection?.fields?.find((f) => f.name == "email")?.required}
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
                        {#if collection?.fields?.find((f) => f.name == "emailVisibility")?.required}
                            <span class="label label-success">Required</span>
                        {:else}
                            <span class="label label-warning">Optional</span>
                        {/if}
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
                    This field can be set only by superusers or auth records with "Manage" access.
                </td>
            </tr>
            <tr>
                <td colspan="3" class="txt-hint txt-bold">Other fields</td>
            </tr>
        {/if}

        {#each tableFields as field (field.name)}
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if !field.required || (field.type == "text" && field.autogeneratePattern)}
                            <span class="label label-warning">Optional</span>
                        {:else}
                            <span class="label label-success">Required</span>
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
                        {#if field.autogeneratePattern}
                            It is autogenerated if not set.
                        {/if}
                    {:else if field.type === "number"}
                        Number value.
                    {:else if field.type === "json"}
                        JSON array or object.
                    {:else if field.type === "email"}
                        Email address.
                    {:else if field.type === "url"}
                        URL address.
                    {:else if field.type === "geoPoint"}
                        <code>{`{"lon":x,"lat":y}`}</code> object.
                    {:else if field.type === "file"}
                        File object.<br />
                        Set to empty value (<code>null</code>, <code>""</code> or <code>[]</code>) to delete
                        already uploaded file(s).
                    {:else if field.type === "relation"}
                        Relation record {field.maxSelect === 1 ? "id" : "ids"}.
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
