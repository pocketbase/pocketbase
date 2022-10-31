<script>
    import { scale, slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import { errors, removeError } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";

    export let key;
    export let title;
    export let icon = "";
    export let config = {};
    export let optionsComponent;

    let accordion;

    $: hasErrors = !CommonHelper.isEmpty(CommonHelper.getNestedVal($errors, key));

    $: if (!config.enabled) {
        removeError(key);
    }

    export function expand() {
        accordion?.expand();
    }

    export function collapse() {
        accordion?.collapse();
    }

    export function collapseSiblings() {
        accordion?.collapseSiblings();
    }
</script>

<Accordion bind:this={accordion} on:expand on:collapse on:toggle {...$$restProps}>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            {#if icon}
                <i class={icon} />
            {/if}
            <span class="txt">{title}</span>
        </div>

        {#if config.enabled}
            <span class="label label-success">Enabled</span>
        {:else}
            <span class="label label-hint">Disabled</span>
        {/if}

        <div class="flex-fill" />

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "Has errors", position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle m-b-0" name="{key}.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={config.enabled} />
        <label for={uniqueId}>Enable</label>
    </Field>

    {#if config.enabled}
        <div class="grid" transition:slide|local={{ duration: 200 }}>
            <div class="col-12 spacing" />
            <div class="col-lg-6">
                <Field class="form-field required" name="{key}.clientId" let:uniqueId>
                    <label for={uniqueId}>Client ID</label>
                    <input type="text" id={uniqueId} bind:value={config.clientId} required />
                </Field>
            </div>

            <div class="col-lg-6">
                <Field class="form-field required" name="{key}.clientSecret" let:uniqueId>
                    <label for={uniqueId}>Client Secret</label>
                    <RedactedPasswordInput bind:value={config.clientSecret} id={uniqueId} required />
                </Field>
            </div>

            {#if optionsComponent}
                <div class="col-lg-12">
                    <svelte:component this={optionsComponent} {key} bind:config />
                </div>
            {/if}
        </div>
    {/if}
</Accordion>
