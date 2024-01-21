<script>
    import { scale, slide } from "svelte/transition";
    import { errors } from "@/stores/errors";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import Accordion from "@/components/base/Accordion.svelte";

    export let collection;

    $: if (collection.type === "auth" && CommonHelper.isEmpty(collection.options)) {
        collection.options = {
            allowEmailAuth: true,
            allowUsernameAuth: true,
            allowOAuth2Auth: true,
            minPasswordLength: 8,
        };
    }

    $: hasUsernameErrors = false;

    $: hasEmailErrors =
        !CommonHelper.isEmpty($errors?.options?.allowEmailAuth) ||
        !CommonHelper.isEmpty($errors?.options?.onlyEmailDomains) ||
        !CommonHelper.isEmpty($errors?.options?.exceptEmailDomains);

    $: hasOAuth2Errors = !CommonHelper.isEmpty($errors?.options?.allowOAuth2Auth);
</script>

<h4 class="section-title">Auth methods</h4>

<div class="accordions">
    <Accordion single>
        <svelte:fragment slot="header">
            <div class="inline-flex">
                <i class="ri-user-star-line" />
                <span class="txt">Username/Password</span>
            </div>

            <div class="flex-fill" />

            {#if collection.options.allowUsernameAuth}
                <span class="label label-success">Enabled</span>
            {:else}
                <span class="label">Disabled</span>
            {/if}

            {#if hasUsernameErrors}
                <i
                    class="ri-error-warning-fill txt-danger"
                    transition:scale={{ duration: 150, start: 0.7 }}
                    use:tooltip={{ text: "Has errors", position: "left" }}
                />
            {/if}
        </svelte:fragment>

        <Field class="form-field form-field-toggle m-b-0" name="options.allowUsernameAuth" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.options.allowUsernameAuth} />
            <label for={uniqueId}>Enable</label>
        </Field>
    </Accordion>

    <Accordion single>
        <svelte:fragment slot="header">
            <div class="inline-flex">
                <i class="ri-mail-star-line" />
                <span class="txt">Email/Password</span>
            </div>

            <div class="flex-fill" />

            {#if collection.options.allowEmailAuth}
                <span class="label label-success">Enabled</span>
            {:else}
                <span class="label">Disabled</span>
            {/if}

            {#if hasEmailErrors}
                <i
                    class="ri-error-warning-fill txt-danger"
                    transition:scale={{ duration: 150, start: 0.7 }}
                    use:tooltip={{ text: "Has errors", position: "left" }}
                />
            {/if}
        </svelte:fragment>

        <Field class="form-field form-field-toggle m-0" name="options.allowEmailAuth" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.options.allowEmailAuth} />
            <label for={uniqueId}>Enable</label>
        </Field>

        {#if collection.options.allowEmailAuth}
            <div class="grid grid-sm p-t-sm" transition:slide={{ duration: 150 }}>
                <div class="col-lg-6">
                    <Field
                        class="form-field {!CommonHelper.isEmpty(collection.options.onlyEmailDomains)
                            ? 'disabled'
                            : ''}"
                        name="options.exceptEmailDomains"
                        let:uniqueId
                    >
                        <label for={uniqueId}>
                            <span class="txt">Except domains</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: 'Email domains that are NOT allowed to sign up. \n This field is disabled if "Only domains" is set.',
                                    position: "top",
                                }}
                            />
                        </label>
                        <MultipleValueInput
                            id={uniqueId}
                            disabled={!CommonHelper.isEmpty(collection.options.onlyEmailDomains)}
                            bind:value={collection.options.exceptEmailDomains}
                        />
                        <div class="help-block">Use comma as separator.</div>
                    </Field>
                </div>
                <div class="col-lg-6">
                    <Field
                        class="form-field {!CommonHelper.isEmpty(collection.options.exceptEmailDomains)
                            ? 'disabled'
                            : ''}"
                        name="options.onlyEmailDomains"
                        let:uniqueId
                    >
                        <label for={uniqueId}>
                            <span class="txt">Only domains</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: 'Email domains that are ONLY allowed to sign up. \n This field is disabled if "Except domains" is set.',
                                    position: "top",
                                }}
                            />
                        </label>
                        <MultipleValueInput
                            id={uniqueId}
                            disabled={!CommonHelper.isEmpty(collection.options.exceptEmailDomains)}
                            bind:value={collection.options.onlyEmailDomains}
                        />
                        <div class="help-block">Use comma as separator.</div>
                    </Field>
                </div>
            </div>
        {/if}
    </Accordion>

    <Accordion single>
        <svelte:fragment slot="header">
            <div class="inline-flex">
                <i class="ri-shield-star-line" />
                <span class="txt">OAuth2</span>
            </div>

            <div class="flex-fill" />

            {#if collection.options.allowOAuth2Auth}
                <span class="label label-success">Enabled</span>
            {:else}
                <span class="label">Disabled</span>
            {/if}

            {#if hasOAuth2Errors}
                <i
                    class="ri-error-warning-fill txt-danger"
                    transition:scale={{ duration: 150, start: 0.7 }}
                    use:tooltip={{ text: "Has errors", position: "left" }}
                />
            {/if}
        </svelte:fragment>

        <Field class="form-field form-field-toggle m-b-0" name="options.allowOAuth2Auth" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.options.allowOAuth2Auth} />
            <label for={uniqueId}>Enable</label>
        </Field>

        {#if collection.options.allowOAuth2Auth}
            <div class="block" transition:slide={{ duration: 150 }}>
                <div class="flex p-t-base">
                    <a href="#/settings/auth-providers" target="_blank" class="btn btn-sm btn-outline">
                        <span class="txt">Manage OAuth2 providers</span>
                    </a>
                </div>
            </div>
        {/if}
    </Accordion>
</div>

<hr />

<h4 class="section-title">General</h4>

<Field class="form-field required" name="options.minPasswordLength" let:uniqueId>
    <label for={uniqueId}>Minimum password length</label>
    <input
        type="number"
        id={uniqueId}
        required
        min="6"
        max="72"
        bind:value={collection.options.minPasswordLength}
    />
</Field>

<Field class="form-field form-field-toggle m-b-sm" name="options.requireEmail" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={collection.options.requireEmail} />
    <label for={uniqueId}>
        <span class="txt">Always require email</span>
        <i
            class="ri-information-line txt-sm link-hint"
            use:tooltip={{
                text: "The constraint is applied only for new records.\nAlso note that some OAuth2 providers (like Twitter), don't return an email and the authentication may fail if the email field is required.",
                position: "right",
            }}
        />
    </label>
</Field>

<Field class="form-field form-field-toggle m-b-sm" name="options.onlyVerified" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={collection.options.onlyVerified} />
    <label for={uniqueId}>
        <span class="txt">Forbid authentication for unverified users</span>
        <i
            class="ri-information-line txt-sm link-hint"
            use:tooltip={{
                text: [
                    "If enabled, it returns 403 for new unverified user authentication requests.",
                    "If you need more granular control, don't enable this option and instead use the `@request.auth.verified = true` rule in the specific collection(s) you are targeting.",
                ].join("\n"),
                position: "right",
            }}
        />
    </label>
</Field>
