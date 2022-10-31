<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import providersList from "@/providers.js";

    const dispatch = createEventDispatcher();

    export let record;

    let externalAuths = [];
    let isLoading = false;

    function getProviderTitle(provider) {
        return providersList[provider + "Auth"]?.title || CommonHelper.sentenize(provider, false);
    }

    function getProviderIcon(provider) {
        return providersList[provider + "Auth"]?.icon || `ri-${provider}-line`;
    }

    async function loadExternalAuths() {
        if (!record?.id) {
            externalAuths = [];
            isLoading = false;
            return;
        }

        isLoading = true;

        try {
            externalAuths = await ApiClient.collection(record.collectionId).listExternalAuths(record.id);
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }

    function unlinkExternalAuth(provider) {
        if (!record?.id || !provider) {
            return; // nothing to unlink
        }

        confirm(`Do you really want to unlink the ${getProviderTitle(provider)} provider?`, () => {
            return ApiClient.collection(record.collectionId)
                .unlinkExternalAuth(record.id, provider)
                .then(() => {
                    addSuccessToast(`Successfully unlinked the ${getProviderTitle(provider)} provider.`);
                    dispatch("unlink", provider);
                    loadExternalAuths(); // reload list
                })
                .catch((err) => {
                    ApiClient.errorResponseHandler(err);
                });
        });
    }

    loadExternalAuths();
</script>

{#if isLoading}
    <div class="block txt-center">
        <span class="loader" />
    </div>
{:else if record?.id && externalAuths.length}
    <div class="list">
        {#each externalAuths as auth}
            <div class="list-item">
                <i class={getProviderIcon(auth.provider)} />
                <span class="txt">{getProviderTitle(auth.provider)}</span>
                <div class="txt-hint">ID: {auth.providerId}</div>
                <button
                    type="button"
                    class="btn btn-secondary link-hint btn-circle btn-sm m-l-auto"
                    on:click={() => unlinkExternalAuth(auth.provider)}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
        {/each}
    </div>
{:else}
    <p class="txt-hint txt-center">No linked OAuth2 providers.</p>
{/if}
