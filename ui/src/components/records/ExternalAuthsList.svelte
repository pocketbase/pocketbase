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

    function getProviderConfig(provider) {
        return providersList.find((p) => p.key == provider + "Auth") || {};
    }

    function getProviderTitle(provider) {
        return getProviderConfig(provider)?.title || CommonHelper.sentenize(provider, false);
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
            ApiClient.error(err);
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
                    ApiClient.error(err);
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
                <figure class="provider-logo">
                    <img
                        src="{import.meta.env.BASE_URL}images/oauth2/{getProviderConfig(auth.provider)?.logo}"
                        alt="Provider logo"
                    />
                </figure>
                <span class="txt">{getProviderTitle(auth.provider)}</span>
                <div class="txt-hint">ID: {auth.providerId}</div>
                <button
                    type="button"
                    class="btn btn-transparent link-hint btn-circle btn-sm m-l-auto"
                    on:click={() => unlinkExternalAuth(auth.provider)}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
        {/each}
    </div>
{:else}
    <h6 class="txt-hint txt-center m-t-sm m-b-sm">No linked OAuth2 providers.</h6>
{/if}
