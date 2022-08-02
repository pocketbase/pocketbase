<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";

    export let params;

    let newPassword = "";
    let newPasswordConfirm = "";
    let isLoading = false;
    let success = false;

    $: email = CommonHelper.getJWTPayload(params?.token).email || "";

    async function submit() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.users.confirmPasswordReset(params?.token, newPassword, newPasswordConfirm);
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
                <p>Password changed</p>
                <p>You can now sign in with your new password.</p>
            </div>
        </div>

        <button type="button" class="btn btn-secondary btn-block" on:click={() => window.close()}>
            Close
        </button>
    {:else}
        <form on:submit|preventDefault={submit}>
            <div class="content txt-center m-b-sm">
                <h4 class="m-b-xs">
                    Reset your user password
                    {#if email}
                        for <strong>{email}</strong>
                    {/if}
                </h4>
            </div>

            <Field class="form-field required" name="password" let:uniqueId>
                <label for={uniqueId}>New password</label>
                <!-- svelte-ignore a11y-autofocus -->
                <input type="password" id={uniqueId} required autofocus bind:value={newPassword} />
            </Field>

            <Field class="form-field required" name="passwordConfirm" let:uniqueId>
                <label for={uniqueId}>New password confirm</label>
                <input type="password" id={uniqueId} required bind:value={newPasswordConfirm} />
            </Field>

            <button
                type="submit"
                class="btn btn-lg btn-block"
                class:btn-loading={isLoading}
                disabled={isLoading}
            >
                <span class="txt">Set new password</span>
            </button>
        </form>
    {/if}
</FullPage>
