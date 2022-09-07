<script>
    import { scale, slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";

    export let config = {}; // EmailAuthConfig

    let accordion;

    $: hasErrors = !CommonHelper.isEmpty($errors?.emailAuth);

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
            <i class="ri-mail-lock-line" />
            <span class="txt">Email/Password</span>
        </div>

        {#if config.enabled}
            <span class="label label-success">Enabled</span>
        {:else}
            <span class="label">Disabled</span>
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

    <Field class="form-field form-field-toggle m-b-0" name="emailAuth.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={config.enabled} />
        <label for={uniqueId}>Enable</label>
    </Field>

    {#if config.enabled}
        <div class="grid" transition:slide|local={{ duration: 150 }}>
            <div class="col-sm-12 m-t-sm">
                <Field class="form-field required" name="emailAuth.minPasswordLength" let:uniqueId>
                    <label for={uniqueId}>Minimum password length</label>
                    <input
                        type="number"
                        id={uniqueId}
                        required
                        min="5"
                        max="200"
                        bind:value={config.minPasswordLength}
                    />
                </Field>
            </div>
            <div class="col-lg-6">
                <Field
                    class="form-field {!CommonHelper.isEmpty(config.onlyDomains) ? 'disabled' : ''}"
                    name="emailAuth.exceptDomains"
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
                        disabled={!CommonHelper.isEmpty(config.onlyDomains)}
                        bind:value={config.exceptDomains}
                    />
                    <div class="help-block">Use comma as separator.</div>
                </Field>
            </div>
            <div class="col-lg-6">
                <Field
                    class="form-field {!CommonHelper.isEmpty(config.exceptDomains) ? 'disabled' : ''}"
                    name="emailAuth.onlyDomains"
                    let:uniqueId
                >
                    <label for="{uniqueId}.config.onlyDomains">
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
                        id="{uniqueId}.config.onlyDomains"
                        disabled={!CommonHelper.isEmpty(config.exceptDomains)}
                        bind:value={config.onlyDomains}
                    />
                    <div class="help-block">Use comma as separator.</div>
                </Field>
            </div>
        </div>
    {/if}
</Accordion>
