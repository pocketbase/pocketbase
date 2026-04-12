<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    export let collection;

    $: if (CommonHelper.isEmpty(collection.webauthn)) {
        collection.webauthn = {
            enabled: false,
        };
    }

    $: hasErrors = !CommonHelper.isEmpty($errors?.webauthn);
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-fingerprint-line"></i>
            <span class="txt">WebAuthn / Passkeys</span>
        </div>

        <div class="flex-fill" />

        {#if collection.webauthn.enabled}
            <span class="label label-success">Enabled</span>
        {:else}
            <span class="label">Disabled</span>
        {/if}

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "Has errors", position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle" name="webauthn.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={collection.webauthn.enabled} />
        <label for={uniqueId}>Enable</label>
    </Field>

    {#if collection.webauthn.enabled}
        <div class="help-block m-t-5">
            <p>
                <i class="ri-information-line" />
                WebAuthn allows users to authenticate using hardware security keys or platform authenticators
                (fingerprint, face recognition, etc.). Users must first sign in with another method to register
                their passkey.
            </p>
        </div>
    {/if}
</Accordion>
