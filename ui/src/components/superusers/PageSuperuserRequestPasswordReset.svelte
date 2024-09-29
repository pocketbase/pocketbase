<script>
    import { link } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";

    let email = "";
    let isLoading = false;
    let success = false;

    async function submit() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.collection("_superusers").requestPasswordReset(email);
            success = true;
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }
</script>

<FullPage>
    {#if success}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content">
                <p>Check <strong class="txt-nowrap">{email}</strong> for the recovery link.</p>
            </div>
        </div>
    {:else}
        <form class="m-b-base" on:submit|preventDefault={submit}>
            <div class="content txt-center m-b-sm">
                <h4 class="m-b-xs">Forgotten superuser password</h4>
                <p>Enter the email associated with your account and we’ll send you a recovery link:</p>
            </div>

            <Field class="form-field required" name="email" let:uniqueId>
                <label for={uniqueId}>Email</label>
                <!-- svelte-ignore a11y-autofocus -->
                <input type="email" id={uniqueId} required autofocus bind:value={email} />
            </Field>

            <button
                type="submit"
                class="btn btn-lg btn-block"
                class:btn-loading={isLoading}
                disabled={isLoading}
            >
                <i class="ri-mail-send-line" />
                <span class="txt">Send recovery link</span>
            </button>
        </form>
    {/if}

    <div class="content txt-center">
        <a href="/login" class="link-hint" use:link>Back to login</a>
    </div>
</FullPage>
