<script>
    import { scale } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { errors } from "@/stores/errors";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import RuleField from "@/components/collections/RuleField.svelte";

    export let collection;

    $: hasErrors = !CommonHelper.isEmpty($errors?.mfa);
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-shield-check-line"></i>
            <span class="txt"> Multi-factor authentication (MFA) </span>
        </div>

        <div class="flex-fill" />

        {#if collection.mfa.enabled}
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

    <div class="content m-b-sm">
        <p class="txt-bold">This feature is experimental and may change in the future.</p>
        <p>
            Multi-factor authentication (MFA) requires the user to authenticate with any 2 different auth
            methods (otp, identity/password, oauth2) before issuing an auth token.
            <a
                href={import.meta.env.PB_MFA_DOCS}
                target="_blank"
                rel="noopener noreferrer"
                class="txt-sm link-hint"
                title="Learn more"
            >
                <em>(Learn more)</em>
            </a>.
        </p>
    </div>

    <div class="grid">
        <Field class="form-field form-field-toggle" name="mfa.enabled" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.mfa.enabled} />
            <label for={uniqueId}>
                <span class="txt">Enable</span>
            </label>
        </Field>

        <div class="content" class:fade={!collection.mfa.enabled}>
            <RuleField
                label="MFA rule"
                formKey="mfa.rule"
                superuserToggle={false}
                disabled={!collection.mfa.enabled}
                placeholder="Leave empty to require MFA for everyone"
                {collection}
                bind:rule={collection.mfa.rule}
            >
                <svelte:fragment>
                    <p>This optional rule could be used to enable/disable MFA per account basis.</p>
                    <p>
                        For example, to require MFA only for accounts with non-empty email you can set it to
                        <code>email != ''</code>.
                    </p>
                    <p>Leave the rule empty to require MFA for everyone.</p>
                </svelte:fragment>
            </RuleField>
        </div>
    </div>
</Accordion>
