<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let key = "";
    export let config = {};

    if (CommonHelper.isEmpty(config.pkce)) {
        config.pkce = true;
    }

    if (!config.displayName) {
        config.displayName = "OIDC";
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

<Field class="form-field required" name="{key}.userInfoURL" let:uniqueId>
    <label for={uniqueId}>User info URL</label>
    <input type="url" id={uniqueId} bind:value={config.userInfoURL} required />
</Field>

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
