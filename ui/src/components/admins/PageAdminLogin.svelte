<script>
    import { link, replace, querystring } from "svelte-spa-router";
    import FullPage from "@/components/base/FullPage.svelte";
    import ApiClient from "@/utils/ApiClient";
    import Field from "@/components/base/Field.svelte";
    import { addErrorToast } from "@/stores/toasts";
    import EyeOpen from "@/components/base/svg/EyeOpen.svelte";
    import EyeClosed from "@/components/base/svg/EyeClosed.svelte";

    const queryParams = new URLSearchParams($querystring);

    let email = queryParams.get("demoEmail") || "";
    let password = queryParams.get("demoPassword") || "";
    let isLoading = false;
    let showPassword = false;
    const setPassword = (e) => {
        password = e.target.value;
    };

    function login() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        return ApiClient.admins
            .authViaEmail(email, password)
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
            <label for={uniqueId}>
                <button id="eye" on:click|preventDefault={() => (showPassword = !showPassword)}>
                    {#if showPassword}
                        <EyeClosed />
                    {:else}
                        <EyeOpen />
                    {/if}
                </button>
                Password
            </label>
            <input
                type={showPassword ? "text" : "password"}
                id={uniqueId}
                value={password}
                on:change={setPassword}
                required
            />
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
