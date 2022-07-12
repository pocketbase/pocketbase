<script>
    import { scale, slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import { _ } from '@/services/i18n';

    export let config = {}; // EmailAuthConfig

    let accordion;

    $: hasErrors = !CommonHelper.isEmpty($errors?.emailPassword);

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
            <span class="txt">{$_("settings.auth.tips.emailPassword")}</span>
        </div>

        {#if config.enabled}
            <span class="label label-success">{$_("settings.auth.tips.enabled")}</span>
        {:else}
            <span class="label">{$_("settings.auth.tips.disabled")}</span>
        {/if}

        <div class="flex-fill" />

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: $_("settings.auth.tips.hasErrors"), position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle m-b-0" name="emailPassword.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={config.enabled} />
        <label for={uniqueId}>{$_("settings.auth.tips.enable")}</label>
    </Field>

    {#if config.enabled}
        <div class="grid" transition:slide|local={{ duration: 150 }}>
            <div class="col-sm-12 m-t-sm">
                <Field class="form-field required" name="emailPassword.minPasswordLength" let:uniqueId>
                    <label for={uniqueId}>{$_("settings.auth.tips.minPasswordLength")}</label>
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
                    name="emailPassword.exceptDomains"
                    let:uniqueId
                >
                    <label for={uniqueId}>
                        <span class="txt">{$_("settings.auth.tips.exceptDomains")}</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: $_("settings.auth.help.exceptDomains"),
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id={uniqueId}
                        disabled={!CommonHelper.isEmpty(config.onlyDomains)}
                        bind:value={config.exceptDomains}
                    />
                    <div class="help-block">{$_("settings.auth.help.comma")}</div>
                </Field>
            </div>
            <div class="col-lg-6">
                <Field
                    class="form-field {!CommonHelper.isEmpty(config.exceptDomains) ? 'disabled' : ''}"
                    name="emailPassword.onlyDomains"
                    let:uniqueId
                >
                    <label for="{uniqueId}.config.onlyDomains">
                        <span class="txt">{$_("settings.auth.tips.onlyDomains")}</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: $_("settings.auth.help.onlyDomains"),
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id="{uniqueId}.config.onlyDomains"
                        disabled={!CommonHelper.isEmpty(config.exceptDomains)}
                        bind:value={config.onlyDomains}
                    />
                    <div class="help-block">{$_("settings.auth.help.comma")}</div>
                </Field>
            </div>
        </div>
    {/if}
</Accordion>
