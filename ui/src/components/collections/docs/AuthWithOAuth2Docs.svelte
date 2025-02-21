<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(
                {
                    token: "JWT_AUTH_TOKEN",
                    record: CommonHelper.dummyCollectionRecord(collection),
                    meta: {
                        id: "abc123",
                        name: "John Doe",
                        username: "john.doe",
                        email: "test@example.com",
                        avatarURL: "https://example.com/avatar.png",
                        accessToken: "...",
                        refreshToken: "...",
                        expiry: "2022-01-01 10:00:00.123Z",
                        isNew: false,
                        rawUser: {},
                    },
                },
                null,
                2,
            ),
        },
        {
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "An error occurred while submitting the form.",
                  "data": {
                    "provider": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">Auth with OAuth2 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p>
    <p>
        For more details please check the
        <a href={import.meta.env.PB_OAUTH2_EXAMPLE} target="_blank" rel="noopener noreferrer">
            OAuth2 integration documentation
        </a>.
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${backendAbsUrl}/api/oauth2-redirect as redirect url.
        const authData = await pb.collection('${collection.name}').authWithOAuth2({ provider: 'google' });

        // OR authenticate with manual OAuth2 code exchange
        // const authData = await pb.collection('${collection.name}').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        // OAuth2 authentication with a single realtime call.
        //
        // Make sure to register ${backendAbsUrl}/api/oauth2-redirect as redirect url.
        final authData = await pb.collection('${collection.name}').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // OR authenticate with manual OAuth2 code exchange
        // final authData = await pb.collection('${collection.name}').authWithOAuth2Code(...);

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}
/>

<h6 class="m-b-xs">API details</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/auth-with-oauth2
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
                    <span>provider</span>
                </div>
            </td>
            <td>
                <span class="label">String</span>
            </td>
            <td>The name of the OAuth2 client provider (eg. "google").</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">Required</span>
                    <span>code</span>
                </div>
            </td>
            <td>
                <span class="label">String</span>
            </td>
            <td>The authorization code returned from the initial request.</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">Required</span>
                    <span>codeVerifier</span>
                </div>
            </td>
            <td>
                <span class="label">String</span>
            </td>
            <td>The code verifier sent with the initial request as part of the code_challenge.</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">Required</span>
                    <span>redirectURL</span>
                </div>
            </td>
            <td>
                <span class="label">String</span>
            </td>
            <td>The redirect url sent with the initial request.</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-warning">Optional</span>
                    <span>createData</span>
                </div>
            </td>
            <td>
                <span class="label">Object</span>
            </td>
            <td>
                <p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p>
                <p>
                    The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br />
                    <em>
                        The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.
                    </em>
                </p>
            </td>
        </tr>
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
                Auto expand record relations. Ex.:
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                Supports up to 6-levels depth nested relations expansion. <br />
                The expanded relations will be appended to the record under the
                <code>expand</code> property (eg. <code>{`"expand": {"relField1": {...}, ...}`}</code>).
                <br />
                Only the relations to which the request user has permissions to <strong>view</strong> will be expanded.
            </td>
        </tr>
        <FieldsQueryParam prefix="record." />
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
