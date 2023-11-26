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
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">Confirm password reset ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>Confirms <strong>{collection.name}</strong> password reset request and sets a new password.</p>
    <p>
        After this request all previously issued tokens for the specific record will be automatically
        invalidated.
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        let oldAuth = pb.authStore.model;

        await pb.collection('${collection?.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );

        // reauthenticate if needed
        // (after the above call all previously issued tokens are invalidated)
        await pb.collection('${collection?.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        final oldAuth = pb.authStore.model;

        await pb.collection('${collection?.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );

        // reauthenticate if needed
        // (after the above call all previously issued tokens are invalidated)
        await pb.collection('${collection?.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/confirm-password-reset
        </p>
    </div>
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
                    <span class="label label-success">Required</span>
                    <span>token</span>
                </div>
            </td>
            <td>
                <span class="label">String</span>
            </td>
            <td>The token from the password reset request email.</td>
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
            <td>The new password to set.</td>
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
            <td>The new password confirmation.</td>
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
