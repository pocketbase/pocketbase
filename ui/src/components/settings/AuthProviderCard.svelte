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
                })
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

<style lang="scss">
    .provider-logo {
        $boxSize: 32px;
        $imgSize: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        width: $boxSize;
        height: $boxSize;
        border-radius: var(--baseRadius);
        background: var(--bodyColor);
        padding: 0;
        gap: 0;
        img {
            max-width: $imgSize;
            max-height: $imgSize;
            height: auto;
            flex-shrink: 0;
        }
    }
    .provider-card {
        display: flex;
        align-items: center;
        width: 100%;
        height: 100%;
        gap: 10px;
        padding: 10px;
        border-radius: var(--baseRadius);
        border: 1px solid var(--baseAlt1Color);
    }
</style>
