<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";

    export let params;

    let password = "";
    let isLoading = false;
    let success = false;

    $: newEmail = CommonHelper.getJWTPayload(params?.token).newEmail || "";

    async function submit() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.users.confirmEmailChange(params?.token, password);
            success = true;
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }
</script>

<FullPage nobranding>
    {#if success}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content txt-bold">
                <p>Email address changed</p>
                <p>You can now sign in with your new email address.</p>
            </div>
        </div>

        <button type="button" class="btn btn-secondary btn-block" on:click={() => window.close()}>
            Close
        </button>
    {:else}
        <form on:submit|preventDefault={submit}>
            <div class="content txt-center m-b-sm">
                <h4 class="m-b-xs">
                    Type your password to confirm changing your email address
                    {#if newEmail}
                        to <strong class="txt-nowrap">{newEmail}</strong>
                    {/if}
                </h4>
            </div>

            <Field class="form-field required" name="password" let:uniqueId>
                <label for={uniqueId}>Password</label>
                <!-- svelte-ignore a11y-autofocus -->
                <input type="password" id={uniqueId} required autofocus bind:value={password} />
            </Field>

            <button
                type="submit"
                class="btn btn-lg btn-block"
                class:btn-loading={isLoading}
                disabled={isLoading}
            >
                <span class="txt">Confirm new email</span>
            </button>
        </form>
    {/if}
</FullPage>
