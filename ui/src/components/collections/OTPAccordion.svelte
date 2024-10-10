<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    export let collection;

    $: isSuperusers = collection?.system && collection?.name === "_superusers";

    $: if (CommonHelper.isEmpty(collection.otp)) {
        collection.otp = {
            enabled: true,
            duration: 300,
            length: 8,
        };
    }

    $: hasErrors = !CommonHelper.isEmpty($errors?.otp);
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-time-line"></i>
            <span class="txt">One-time password (OTP)</span>
        </div>

        <div class="flex-fill" />

        {#if collection.otp.enabled}
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

    <Field class="form-field form-field-toggle" name="otp.enabled" let:uniqueId>
        <input
            type="checkbox"
            id={uniqueId}
            bind:checked={collection.otp.enabled}
            on:change={(e) => {
                if (isSuperusers) {
                    collection.mfa.enabled = e.target.checked;
                }
            }}
        />
        <label for={uniqueId}>Enable</label>
        {#if isSuperusers}
            <i
                class="ri-information-line link-hint"
                use:tooltip={{
                    text: "Superusers can have OTP only as part of Two-factor authentication.",
                    position: "right",
                }}
            />
        {/if}
    </Field>

    <div class="grid grid-sm">
        <div class="col-sm-6">
            <Field class="form-field form-field-toggle required" name="otp.duration" let:uniqueId>
                <label for={uniqueId}>Duration (in seconds)</label>
                <input
                    type="number"
                    min="0"
                    step="1"
                    id={uniqueId}
                    bind:value={collection.otp.duration}
                    required
                />
            </Field>
        </div>
        <div class="col-sm-6">
            <Field class="form-field form-field-toggle required" name="otp.length" let:uniqueId>
                <label for={uniqueId}>Generated password length</label>
                <input
                    type="number"
                    min="0"
                    step="1"
                    id={uniqueId}
                    bind:value={collection.otp.length}
                    required
                />
            </Field>
        </div>
    </div>
</Accordion>
