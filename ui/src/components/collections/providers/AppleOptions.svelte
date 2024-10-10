<script>
    import AppleSecretPopup from "@/components/collections/providers/AppleSecretPopup.svelte";

    export let key = "";
    export let config = {};

    let generatorPopup;
</script>

<button
    type="button"
    class="btn btn-sm btn-secondary btn-provider-{key}"
    on:click={() => generatorPopup?.show({ clientId: config.clientId })}
>
    <i class="ri-key-line" />
    <span class="txt">Generate secret</span>
</button>

<AppleSecretPopup
    bind:this={generatorPopup}
    on:submit={(e) => {
        config.clientSecret = e.detail?.secret || "";
    }}
/>
