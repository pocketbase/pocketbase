<script>
    import FullPage from "@/components/base/FullPage.svelte";
    import ApiClient from "@/utils/ApiClient";
    import PocketBase, { getTokenPayload, isTokenExpired } from "pocketbase";

    export let params;

    let successConfirm = false;
    let isConfirming = false;

    let successResend = false;
    let isResending = false;

    send();

    async function send() {
        if (isConfirming) {
            return;
        }

        isConfirming = true;

        // init a custom client to avoid interfering with the superuser state
        const client = new PocketBase(import.meta.env.PB_BACKEND_URL);

        try {
            const payload = getTokenPayload(params?.token);
            await client.collection(payload.collectionId).confirmVerification(params?.token);
            successConfirm = true;
        } catch (err) {
            successConfirm = false;
        }

        isConfirming = false;
    }

    $: canResend = params?.token && isTokenExpired(params.token);

    async function resend() {
        const payload = getTokenPayload(params?.token);

        if (isResending || !payload.collectionId || !payload.email) {
            return;
        }

        isResending = true;

        // init a custom client to avoid interfering with the superuser state
        const client = new PocketBase(import.meta.env.PB_BACKEND_URL);

        try {
            const payload = getTokenPayload(params?.token);
            await client.collection(payload.collectionId).requestVerification(payload.email);
            successResend = true;
        } catch (err) {
            ApiClient.error(err);
            successResend = false;
        }

        isResending = false;
    }
</script>

<FullPage nobranding>
    {#if isConfirming}
        <div class="txt-center">
            <div class="loader loader-lg">
                <em>Please wait...</em>
            </div>
        </div>
    {:else if successConfirm}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content txt-bold">
                <p>Successfully verified email address.</p>
            </div>
        </div>

        <button type="button" class="btn btn-transparent btn-block" on:click={() => window.close()}>
            Close
        </button>
    {:else if successResend}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content txt-bold">
                <p>Please check your email for the new verification link.</p>
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

        {#if canResend}
            <button
                type="button"
                class="btn btn-transparent btn-block"
                class:btn-loading={isResending}
                disabled={isResending}
                on:click={resend}
            >
                <span class="txt">Resend</span>
            </button>
        {:else}
            <button type="button" class="btn btn-transparent btn-block" on:click={() => window.close()}>
                Close
            </button>
        {/if}
    {/if}
</FullPage>
