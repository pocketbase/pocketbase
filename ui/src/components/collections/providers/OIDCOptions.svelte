<script>
    import { slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let key = "";
    export let config = {};

    const userInfoOptions = [
        { label: "User info URL", value: true },
        { label: "ID Token", value: false },
    ];

    let hasUserInfoURL = !!config.userInfoURL;

    if (CommonHelper.isEmpty(config.pkce)) {
        config.pkce = true;
    }

    if (!config.displayName) {
        config.displayName = "OIDC";
    }

    if (!config.extra) {
        config.extra = {};
        hasUserInfoURL = true;
    }

    $: if (typeof hasUserInfoURL !== undefined) {
        refreshUserInfoState();
    }

    function refreshUserInfoState() {
        if (!hasUserInfoURL) {
            config.userInfoURL = "";
            config.extra = config.extra || {};
        } else {
            config.extra = {};
        }
    }
</script>

<Field class="form-field required" name="{key}.displayName" let:uniqueId>
    <label for={uniqueId}>Display name</label>
    <input type="text" id={uniqueId} bind:value={config.displayName} required />
</Field>

<div class="section-title">Endpoints</div>

<Field class="form-field required" name="{key}.authURL" let:uniqueId>
    <label for={uniqueId}>Auth URL</label>
    <input type="url" id={uniqueId} bind:value={config.authURL} required />
</Field>

<Field class="form-field required" name="{key}.tokenURL" let:uniqueId>
    <label for={uniqueId}>Token URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenURL} required />
</Field>

<Field class="form-field m-b-xs" let:uniqueId>
    <label for={uniqueId}>Fetch user info from</label>
    <ObjectSelect id={uniqueId} items={userInfoOptions} bind:keyOfSelected={hasUserInfoURL} />
</Field>

<div class="sub-panel m-b-base">
    {#if hasUserInfoURL}
        <div class="content" transition:slide={{ delay: 10, duration: 150 }}>
            <Field class="form-field required" name="{key}.userInfoURL" let:uniqueId>
                <label for={uniqueId}>User info URL</label>
                <input type="url" id={uniqueId} bind:value={config.userInfoURL} required />
            </Field>
        </div>
    {:else}
        <div class="content" transition:slide={{ delay: 10, duration: 150 }}>
            <p class="txt-hint txt-sm m-b-xs">
                <em>
                    Both fields are considered optional because the parsed <code>id_token</code>
                    is a direct result of the trusted server code->token exchange response.
                </em>
            </p>
            <Field class="form-field m-b-xs" name="{key}.extra.jwksURL" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">JWKS verification URL</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "URL to the public token verification keys.",
                            position: "top",
                        }}
                    />
                </label>
                <input type="url" id={uniqueId} bind:value={config.extra.jwksURL} />
            </Field>
            <Field class="form-field" name="{key}.extra.issuers" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">Issuers</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "Comma separated list of accepted values for the iss token claim validation.",
                            position: "top",
                        }}
                    />
                </label>
                <MultipleValueInput id={uniqueId} bind:value={config.extra.issuers} />
            </Field>
        </div>
    {/if}
</div>

<Field class="form-field" name="{key}.pkce" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={config.pkce} />
    <label for={uniqueId}>
        <span class="txt">Support PKCE</span>
        <i
            class="ri-information-line link-hint"
            use:tooltip={{
                text: "Usually it should be safe to be always enabled as most providers will just ignore the extra query parameters if they don't support PKCE.",
                position: "right",
            }}
        />
    </label>
</Field>
