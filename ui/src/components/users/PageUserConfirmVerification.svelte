<script>
    import ApiClient from "@/utils/ApiClient";
    import FullPage from "@/components/base/FullPage.svelte";

    export let params;

    let success = false;
    let isLoading = false;

    send();

    async function send() {
        isLoading = true;

        try {
            await ApiClient.users.confirmVerification(params?.token);
            success = true;
        } catch (err) {
            console.warn(err);
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

        <button type="button" class="btn btn-secondary btn-block" on:click={() => window.close()}>
            Close
        </button>
    {:else}
        <div class="alert alert-danger">
            <div class="icon"><i class="ri-error-warning-line" /></div>
            <div class="content txt-bold">
                <p>Invalid or expired verification token.</p>
            </div>
        </div>

        <button type="button" class="btn btn-secondary btn-block" on:click={() => window.close()}>
            Close
        </button>
    {/if}
</FullPage>
