<script>
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let record;
    export let collection;

    let credentials = [];
    let isLoading = false;
    let isClearingAll = false;

    async function loadCredentials() {
        if (!record?.id || !collection?.id) {
            credentials = [];
            isLoading = false;
            return;
        }

        isLoading = true;

        try {
            credentials = await ApiClient.send(
                `/api/collections/${encodeURIComponent(collection.name)}/auth-with-webauthn/credentials-by-record/${encodeURIComponent(record.id)}`,
                { method: "GET" },
            ).catch(() => {
                // fallback: query the system collection directly
                return ApiClient.collection("_webauthnCredentials").getFullList({
                    filter: ApiClient.filter("collectionRef = {:collectionId} && recordRef = {:recordId}", {
                        collectionId: collection.id,
                        recordId: record.id,
                    }),
                });
            });
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    function deleteCredential(credential) {
        if (!record?.id || !credential) {
            return;
        }

        const name = credential.name || credential.id;
        confirm(`Do you really want to delete the passkey "${name}"?`, () => {
            return ApiClient.collection("_webauthnCredentials")
                .delete(credential.id)
                .then(() => {
                    addSuccessToast(`Successfully deleted passkey "${name}".`);
                    dispatch("delete", credential.id);
                    loadCredentials();
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function clearAllCredentials() {
        if (!record?.id) {
            return;
        }

        confirm(
            `Do you really want to remove ALL passkeys for this user? They will need to use another auth method to sign in.`,
            async () => {
                isClearingAll = true;

                try {
                    const result = await ApiClient.send(
                        `/api/collections/${encodeURIComponent(collection.name)}/auth-with-webauthn/credentials-by-record/${encodeURIComponent(record.id)}`,
                        { method: "DELETE" },
                    );
                    addSuccessToast(`Successfully removed ${result.deleted || 0} passkey(s).`);
                    dispatch("clear");
                    loadCredentials();
                } catch (err) {
                    ApiClient.error(err);
                }

                isClearingAll = false;
            },
        );
    }

    loadCredentials();
</script>

{#if isLoading}
    <div class="block txt-center">
        <span class="loader" />
    </div>
{:else if record?.id && credentials.length}
    <div class="list">
        {#each credentials as cred}
            <div class="list-item">
                <i class="ri-fingerprint-line txt-hint" />
                <span class="txt">{cred.name || "Unnamed passkey"}</span>
                <div class="txt-hint txt-sm">
                    {#if cred.created}
                        Registered: {CommonHelper.formatToLocalDate(cred.created)}
                    {/if}
                    {#if cred.signCount !== undefined}
                        <span class="m-l-5">· Uses: {cred.signCount}</span>
                    {/if}
                </div>
                <button
                    type="button"
                    class="btn btn-transparent link-hint btn-circle btn-sm m-l-auto"
                    on:click={() => deleteCredential(cred)}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
        {/each}
    </div>

    <div class="m-t-sm txt-right">
        <button
            type="button"
            class="btn btn-xs btn-transparent btn-danger"
            class:btn-loading={isClearingAll}
            disabled={isClearingAll}
            on:click={clearAllCredentials}
        >
            <i class="ri-delete-bin-line" />
            <span class="txt">Clear all passkeys</span>
        </button>
    </div>
{:else}
    <h6 class="txt-hint txt-center m-t-sm m-b-sm">No registered passkeys.</h6>
{/if}
