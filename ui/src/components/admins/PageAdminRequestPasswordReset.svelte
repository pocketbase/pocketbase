<script>
    import { link } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";
    import { _ } from '@/services/i18n';
    

    let email = "";
    let isLoading = false;
    let success = false;

    async function submit() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.Admins.requestPasswordReset(email);
            success = true;
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }
</script>

<FullPage>
    {#if success}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content">
                <p>{@html $_("admins.reset.checkemail",{values:{email:email}})}</p>
            </div>
        </div>
    {:else}
        <form class="m-b-base" on:submit|preventDefault={submit}>
            <div class="content txt-center m-b-sm">
                <h4 class="m-b-xs">{$_("admins.reset.title")}</h4>
                <p>{$_("admins.reset.desc")}</p>
            </div>

            <Field class="form-field required" name="email" let:uniqueId>
                <label for={uniqueId}>{$_("admins.reset.email")}</label>
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
                <span class="txt">{$_("admins.reset.send")}</span>
            </button>
        </form>
    {/if}

    <div class="content txt-center">
        <a href="/login" class="link-hint" use:link>{$_("admins.reset.back")}</a>
    </div>
</FullPage>
