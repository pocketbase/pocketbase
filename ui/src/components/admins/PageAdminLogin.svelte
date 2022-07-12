<script>
    import { link } from "svelte-spa-router";
    import { replace } from "svelte-spa-router";
    import FullPage from "@/components/base/FullPage.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import { addErrorToast } from "@/stores/toasts";
    import { _ } from '@/services/i18n';

    const queryParams = CommonHelper.getQueryParams(window.location?.href);

    let email = queryParams.demoEmail || "";
    let password = queryParams.demoPassword || "";
    let isLoading = false;

    function login() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        return ApiClient.Admins.authViaEmail(email, password)
            .then(() => {
                replace("/");
            })
            .catch(() => {
                addErrorToast($_("admins.login.invalid"));
            })
            .finally(() => {
                isLoading = false;
            });
    }
</script>

<FullPage>
    <form class="block" on:submit|preventDefault={login}>
        <div class="content txt-center m-b-base">
            <h4>{$_("admins.login.title")}</h4>
        </div>

        <Field class="form-field required" name="email" let:uniqueId>
            <label for={uniqueId}>{$_("admins.login.email")}</label>
            <!-- svelte-ignore a11y-autofocus -->
            <input type="email" id={uniqueId} bind:value={email} required autofocus />
        </Field>

        <Field class="form-field required" name="password" let:uniqueId>
            <label for={uniqueId}>{$_("admins.login.password")}</label>
            <input type="password" id={uniqueId} bind:value={password} required />
            <div class="help-block">
                <a href="/request-password-reset" class="link-hint" use:link>{$_("admins.login.forgotten_password")}</a>
            </div>
        </Field>

        <button
            type="submit"
            class="btn btn-lg btn-block btn-next"
            class:btn-disabled={isLoading}
            class:btn-loading={isLoading}
        >
            <span class="txt">{$_("admins.login.login")}</span>
            <i class="ri-arrow-right-line" />
        </button>
    </form>
</FullPage>
