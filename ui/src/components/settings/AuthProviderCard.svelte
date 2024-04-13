<script>
    import AuthProviderPanel from "@/components/settings/AuthProviderPanel.svelte";

    export let provider = {};
    export let config = {};

    let providerPanel;
</script>

<div class="provider-card">
    <figure class="provider-logo">
        {#if provider.logo}
            <img src="{import.meta.env.BASE_URL}images/oauth2/{provider.logo}" alt="{provider.title} logo" />
        {/if}
    </figure>
    <div class="title">{provider.title}</div>
    <em class="txt-hint txt-sm m-r-auto">({provider.key.slice(0, -4)})</em>
    {#if config.enabled}
        <div class="label label-success">Enabled</div>
    {/if}
    <button
        type="button"
        class="btn btn-circle btn-hint btn-transparent"
        aria-label="Provider settings"
        on:click={() => {
            providerPanel?.show(
                provider,
                Object.assign({}, config, {
                    enabled: config.clientId ? config.enabled : true,
                    pkce: config.clientId ? config.pkce : null,
                }),
            );
        }}
    >
        <i class="ri-settings-4-line" />
    </button>
</div>

<AuthProviderPanel
    bind:this={providerPanel}
    on:submit={(e) => {
        if (e.detail[provider.key]) {
            config = e.detail[provider.key];
        }
    }}
/>
