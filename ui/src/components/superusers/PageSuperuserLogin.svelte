<script>
    import { link, replace, querystring } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";
    import { setErrors } from "@/stores/errors";
    import { addErrorToast, removeAllToasts } from "@/stores/toasts";

    const queryParams = new URLSearchParams($querystring);

    let identity = queryParams.get("demoEmail") || "";
    let password = queryParams.get("demoPassword") || "";

    let authMethods = {};
    let currentStep = 1;
    let totalSteps = 1;

    let passwordAuthSubmitting = false;
    let otpRequestSubmitting = false;
    let otpAuthSubmitting = false;
    let isLoading = false;

    let mfaId = "";
    let otpId = "";
    let lastOTPId = "";
    let otpEmail = "";
    let otpPassword = "";

    $: {
        totalSteps = 1;
        currentStep = 1;

        if (authMethods?.mfa?.enabled) {
            totalSteps++;
        }

        if (authMethods?.otp?.enabled) {
            totalSteps++;
        }

        if (mfaId != "") {
            currentStep++;
        }

        if (otpId != "") {
            currentStep++;
        }
    }

    loadAuthMethods();

    async function loadAuthMethods() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            authMethods = await ApiClient.collection("_superusers").listAuthMethods();
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    async function authWithPassword() {
        if (passwordAuthSubmitting) {
            return;
        }

        passwordAuthSubmitting = true;

        try {
            await ApiClient.collection("_superusers").authWithPassword(identity, password);
            removeAllToasts();
            setErrors({});
            replace("/");
        } catch (err) {
            if (err.status == 401) {
                mfaId = err.response.mfaId;

                // show the otp forms
                if (
                    // if the identity field is just the email use it to directly send an otp request
                    authMethods?.password?.identityFields?.length == 1 &&
                    authMethods.password.identityFields[0] == "email"
                ) {
                    // prefill and request
                    otpEmail = identity;
                    await requestOTP();
                } else if (/^[^@\s]+@[^@\s]+$/.test(identity)) {
                    // only prefill
                    otpEmail = identity;
                }
            } else if (err.status != 400) {
                ApiClient.error(err);
            } else {
                addErrorToast("Invalid login credentials.");
            }
        }

        passwordAuthSubmitting = false;
    }

    async function requestOTP() {
        if (otpRequestSubmitting) {
            return;
        }

        otpRequestSubmitting = true;

        try {
            const result = await ApiClient.collection("_superusers").requestOTP(otpEmail);
            otpId = result.otpId;
            lastOTPId = otpId;
            removeAllToasts();
            setErrors({});
        } catch (err) {
            // reset the form
            if (err.status == 429) {
                otpId = lastOTPId;
            }

            ApiClient.error(err);
        }

        otpRequestSubmitting = false;
    }

    async function authWithOTP() {
        if (otpAuthSubmitting) {
            return;
        }

        otpAuthSubmitting = true;

        try {
            await ApiClient.collection("_superusers").authWithOTP(otpId || lastOTPId, otpPassword, { mfaId });
            removeAllToasts();
            setErrors({});
            replace("/");
        } catch (err) {
            ApiClient.error(err);
        }

        otpAuthSubmitting = false;
    }
</script>

<FullPage>
    <div class="content txt-center m-b-base">
        <h4>
            Superuser login
            {#if totalSteps > 1}
                ({currentStep}/{totalSteps})
            {/if}
        </h4>
    </div>

    {#if isLoading}
        <div class="block txt-center">
            <span class="loader" />
        </div>
    {:else if authMethods.password.enabled && !mfaId}
        <!-- auth with password -->
        <form class="block" on:submit|preventDefault={authWithPassword}>
            <Field class="form-field required" name="identity" let:uniqueId>
                <label for={uniqueId}>
                    {CommonHelper.sentenize(authMethods.password.identityFields.join(" or "), false)}
                </label>
                <!-- svelte-ignore a11y-autofocus -->
                <input
                    id={uniqueId}
                    type={authMethods.password.identityFields.length == 1 &&
                    authMethods.password.identityFields[0] == "email"
                        ? "email"
                        : "text"}
                    value={identity}
                    on:input={(e) => {
                        identity = e.target.value;
                    }}
                    required
                    autofocus
                />
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
                class:btn-disabled={passwordAuthSubmitting}
                class:btn-loading={passwordAuthSubmitting}
            >
                <span class="txt">{totalSteps > 1 ? "Next" : "Login"}</span>
                <i class="ri-arrow-right-line" />
            </button>
        </form>
    {:else if authMethods.otp.enabled}
        {#if !otpId}
            <!-- request otp -->
            <form class="block" on:submit|preventDefault={requestOTP}>
                <Field class="form-field required" name="email" let:uniqueId>
                    <label for={uniqueId}>Email</label>
                    <input type="email" id={uniqueId} bind:value={otpEmail} required />
                </Field>

                <button
                    type="submit"
                    class="btn btn-lg btn-block btn-next"
                    class:btn-disabled={otpRequestSubmitting}
                    class:btn-loading={otpRequestSubmitting}
                >
                    <i class="ri-mail-send-line" />
                    <span class="txt">Send OTP</span>
                </button>
            </form>
        {:else}
            {#if otpEmail}
                <div class="content txt-center m-b-sm">
                    <p>
                        Check your <strong>{otpEmail}</strong> inbox and enter in the input below the received
                        One-time password (OTP).
                    </p>
                </div>
            {/if}

            <!-- auth with otp -->
            <form class="block" on:submit|preventDefault={authWithOTP}>
                <Field class="form-field required" name="otpId" let:uniqueId>
                    <label for={uniqueId}>Id</label>
                    <input
                        type="text"
                        id={uniqueId}
                        value={otpId}
                        placeholder={lastOTPId}
                        on:change={(e) => {
                            otpId = e.target.value || lastOTPId;
                            e.target.value = otpId;
                        }}
                        required
                    />
                </Field>

                <Field class="form-field required" name="password" let:uniqueId>
                    <label for={uniqueId}>One-time password</label>
                    <!-- svelte-ignore a11y-autofocus -->
                    <input type="password" id={uniqueId} bind:value={otpPassword} required autofocus />
                </Field>

                <button
                    type="submit"
                    class="btn btn-lg btn-block btn-next"
                    class:btn-disabled={otpAuthSubmitting}
                    class:btn-loading={otpAuthSubmitting}
                >
                    <span class="txt">Login</span>
                    <i class="ri-arrow-right-line" />
                </button>
            </form>

            <div class="content txt-center m-t-sm">
                <button
                    type="button"
                    class="link-hint"
                    disabled={otpAuthSubmitting}
                    on:click={() => {
                        otpId = "";
                    }}
                >
                    Request another OTP
                </button>
            </div>
        {/if}
    {/if}
</FullPage>
