<script>
    import { link } from "svelte-spa-router";
    import { replace } from "svelte-spa-router";
    import FullPage from "@/components/base/FullPage.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import { addErrorToast } from "@/stores/toasts";

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
                addErrorToast("Invalid login credentials.");
            })
            .finally(() => {
                isLoading = false;
            });
    }
</script>

<FullPage>
    <form class="block" on:submit|preventDefault={login}>
        <div class="content txt-center m-b-base">
            <h4>Admin sign in</h4>
        </div>

        <Field class="form-field required" name="email" let:uniqueId>
            <label for={uniqueId}>Email</label>
            <!-- svelte-ignore a11y-autofocus -->
            <input type="email" id={uniqueId} bind:value={email} required autofocus />
        </Field>

        <Field class="form-field required" name="password" let:uniqueId>
            <label for={uniqueId}>Password</label>
            <input type="password" id={uniqueId} bind:value={password} required />
            <div class="help-block">
                <a href="/request-password-reset" class="link-hint" use:link>Forgotten password?</a>
            </div>
        </Field>

        <button
            type="submit"
            class="btn btn-lg btn-block btn-next"
            class:btn-disabled={isLoading}
            class:btn-loading={isLoading}
        >
            <span class="txt">Login</span>
            <i class="ri-arrow-right-line" />
        </button>
    </form>
</FullPage>
