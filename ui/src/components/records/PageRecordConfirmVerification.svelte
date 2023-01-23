<script>
    import PocketBase, { getTokenPayload } from "pocketbase";
    import FullPage from "@/components/base/FullPage.svelte";

    export let params;

    let success = false;
    let isLoading = false;

    send();

    async function send() {
        isLoading = true;

        // init a custom client to avoid interfering with the admin state
        const client = new PocketBase(import.meta.env.PB_BACKEND_URL);

        try {
            const payload = getTokenPayload(params?.token);
            await client.collection(payload.collectionId).confirmVerification(params?.token);
            success = true;
        } catch (err) {
            success = false;
        }

        isLoading = false;
    }
</script>

<FullPage nobranding>
    {#if isLoading}
        <div class="txt-center">
            <div class="loader loader-lg">
                <em>Please wait...</em>
            </div>
        </div>
    {:else if success}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content txt-bold">
                <p>Successfully verified email address.</p>
            </div>
        </div>

        <button type="button" class="btn btn-transparent btn-block" on:click={() => window.close()}>
            Close
        </button>
    {:else}
        <div class="alert alert-danger">
            <div class="icon"><i class="ri-error-warning-line" /></div>
            <div class="content txt-bold">
                <p>Invalid or expired verification token.</p>
            </div>
        </div>

        <button type="button" class="btn btn-transparent btn-block" on:click={() => window.close()}>
            Close
        </button>
    {/if}
</FullPage>
