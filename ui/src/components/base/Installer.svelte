<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import Field from "@/components/base/Field.svelte";
    import EyeOpen from "@/components/base/svg/EyeOpen.svelte";
    import EyeClosed from "@/components/base/svg/EyeClosed.svelte";

    const dispatch = createEventDispatcher();

    let email = "";
    let password = "";
    let passwordConfirm = "";
    let isLoading = false;
    let showPassword = false;
    const setPassword = (e) => {
        password = e.target.value;
    };
    const setPasswordConfirm = (e) => {
        passwordConfirm = e.target.value;
    };

    async function submit() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.admins.create({
                email,
                password,
                passwordConfirm,
            });

            await ApiClient.admins.authViaEmail(email, password);

            dispatch("submit");
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoading = false;
    }
</script>

<!-- TODO -->

<form class="block" autocomplete="off" on:submit|preventDefault={submit}>
    <div class="content txt-center m-b-base">
        <h4>Create your first admin account in order to continue</h4>
    </div>

    <Field class="form-field required" name="email" let:uniqueId>
        <label for={uniqueId}>Email</label>
        <!-- svelte-ignore a11y-autofocus -->
        <input type="email" autocomplete="off" id={uniqueId} bind:value={email} required autofocus />
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
            autocomplete="new-password"
            minlength="10"
            id={uniqueId}
            value={password}
            on:change={setPassword}
            required
        />
        <div class="help-block">Minimum 10 characters.</div>
    </Field>

    <Field class="form-field required" name="passwordConfirm" let:uniqueId>
        <label for={uniqueId}>
            <button id="eye" on:click|preventDefault={() => (showPassword = !showPassword)}>
                {#if showPassword}
                    <EyeClosed />
                {:else}
                    <EyeOpen />
                {/if}
            </button>
            Password confirm
        </label>
        <input
            type={showPassword ? "text" : "password"}
            minlength="10"
            id={uniqueId}
            value={passwordConfirm}
            on:change={setPasswordConfirm}
            required
        />
    </Field>

    <button
        type="submit"
        class="btn btn-lg btn-block btn-next"
        class:btn-disabled={isLoading}
        class:btn-loading={isLoading}
    >
        <span class="txt">Create and login</span>
        <i class="ri-arrow-right-line" />
    </button>
</form>
