<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: isAuth = collection?.type === "auth";

    $: superusersOnly = collection?.updateRule === null;

    $: excludedTableFields = isAuth ? ["id", "password", "verified", "email", "emailVisibility"] : ["id"];

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
                  "message": "Failed to update record.",
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
        {
            code: 404,
            body: `
                {
                  "status": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `,
        },
    ];

    function getPayload(collection) {
        let payload = CommonHelper.dummyCollectionSchemaData(collection, true);

        if (isAuth) {
            payload.oldPassword = "12345678";
            payload.password = "87654321";
            payload.passwordConfirm = "87654321";
            delete payload.verified;
            delete payload.email;
        }

        return payload;
    }
</script>

<h3 class="m-b-sm">Update ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>Update a single <strong>{collection.name}</strong> record.</p>
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
    {#if isAuth}
        <p>
            <em>
                Note that in case of a password change all previously issued tokens for the current record
                will be automatically invalidated and if you want your user to remain signed in you need to
                reauthenticate manually after the update call.
            </em>
        </p>
    {/if}
</div>

<!-- prettier-ignore -->
<SdkTabs
    js={`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${backendAbsUrl}');

...

// example update data
const data = ${JSON.stringify(getPayload(collection), null, 4)};

const record = await pb.collection('${collection?.name}').update('RECORD_ID', data);
    `}
    dart={`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${backendAbsUrl}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(getPayload(collection), null, 2)};

final record = await pb.collection('${collection?.name}').update('RECORD_ID', body: body);
    `}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-warning">
    <strong class="label label-primary">PATCH</strong>
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
            <td>ID of the record to update.</td>
        </tr>
    </tbody>
</table>

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
                        <span class="label label-warning">Optional</span>
                        <span>email</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>
                    The auth record email address.
                    <br />
                    This field can be updated only by superusers or auth records with "Manage" access.
                    <br />
                    Regular accounts can update their email by calling "Request email change".
                </td>
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
                        <span class="label label-warning">Optional</span>
                        <span>oldPassword</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>
                    Old auth record password.
                    <br />
                    This field is required only when changing the record password. Superusers and auth records
                    with "Manage" access can skip this field.
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">Optional</span>
                        <span>password</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>New auth record password.</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">Optional</span>
                        <span>passwordConfirm</span>
                    </div>
                </td>
                <td>
                    <span class="label">String</span>
                </td>
                <td>New auth record password confirmation.</td>
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
                    {:else if field.type === "geoPoint"}
                        <code>{`{"lon":x,"lat":y}`}</code> object.
                    {:else if field.type === "file"}
                        File object.<br />
                        Set to <code>null</code> to delete already uploaded file(s).
                    {:else if field.type === "relation"}
                        Relation record {field.maxSelect == 1 ? "id" : "ids"}.
                    {/if}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

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
            <td>expand</td>
            <td>
                <span class="label">String</span>
            </td>
            <td>
                Auto expand relations when returning the updated record. Ex.:
                <CodeBlock content={`?expand=relField1,relField2.subRelField21`} />
                Supports up to 6-levels depth nested relations expansion. <br />
                The expanded relations will be appended to the record under the
                <code>expand</code> property (eg. <code>{`"expand": {"relField1": {...}, ...}`}</code>). Only
                the relations that the user has permissions to <strong>view</strong> will be expanded.
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
