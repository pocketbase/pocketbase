<script>
    import { scale } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { errors } from "@/stores/errors";
    import Accordion from "@/components/base/Accordion.svelte";
    import TokenField from "@/components/collections/TokenField.svelte";

    export let collection;

    let tokensList = [];

    $: isSuperusers = collection?.system && collection?.name === "_superusers";

    $: tokensList = isSuperusers
        ? [
              { key: "authToken", label: "Auth" },
              { key: "passwordResetToken", label: "Password reset" },
              { key: "fileToken", label: "Protected file access" },
          ]
        : [
              { key: "authToken", label: "Auth" },
              { key: "verificationToken", label: "Email verification" },
              { key: "passwordResetToken", label: "Password reset" },
              { key: "emailChangeToken", label: "Email change" },
              { key: "fileToken", label: "Protected file access" },
          ];

    $: hasErrors = hasTokenError($errors);

    function hasTokenError(errors) {
        if (CommonHelper.isEmpty(errors)) {
            return false;
        }

        for (let token of tokensList) {
            if (errors[token.key]) {
                return true;
            }
        }

        return false;
    }
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-key-2-line"></i>
            <span class="txt">Tokens options (invalidate, duration)</span>
        </div>

        <div class="flex-fill" />

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "Has errors", position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <div class="grid">
        {#each tokensList as token (token.key)}
            <div class="col-sm-6">
                <TokenField
                    key={token.key}
                    label={token.label}
                    bind:duration={collection[token.key].duration}
                    bind:secret={collection[token.key].secret}
                />
            </div>
        {/each}
    </div>
</Accordion>
