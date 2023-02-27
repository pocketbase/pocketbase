<script context="module">
    let cachedRuleComponent;
</script>

<script>
    import { tick } from "svelte";
    import Field from "@/components/base/Field.svelte";

    export let collection = null;
    export let rule = null;
    export let label = "Rule";
    export let formKey = "rule";
    export let required = false;

    let editorRef = null;
    let tempValue = null;
    let ruleInputComponent = cachedRuleComponent;
    let isRuleComponentLoading = false;

    $: isAdminOnly = rule === null;

    loadEditorComponent();

    async function loadEditorComponent() {
        if (ruleInputComponent || isRuleComponentLoading) {
            return; // already loaded or in the process
        }

        isRuleComponentLoading = true;

        ruleInputComponent = (await import("@/components/base/FilterAutocompleteInput.svelte")).default;

        cachedRuleComponent = ruleInputComponent;

        isRuleComponentLoading = false;
    }

    async function unlock() {
        rule = tempValue || "";
        await tick();
        editorRef?.focus();
    }

    async function lock() {
        tempValue = rule;
        rule = null;
    }
</script>

{#if isRuleComponentLoading}
    <div class="txt-center">
        <span class="loader" />
    </div>
{:else}
    <Field
        class="form-field rule-field m-0 {required ? 'requied' : ''} {isAdminOnly ? 'disabled' : ''}"
        name={formKey}
        let:uniqueId
    >
        <label for={uniqueId}>
            <span class="txt" class:txt-hint={isAdminOnly}>
                {label}
            </span>
            <span class="label label-sm">
                {isAdminOnly ? "Admins only" : "Custom rule"}
            </span>

            {#if isAdminOnly}
                <button
                    type="button"
                    class="btn btn-sm btn-transparent btn-success lock-toggle"
                    on:click={unlock}
                >
                    <i class="ri-lock-unlock-line" />
                    <span class="txt">Enable custom rule</span>
                </button>
            {:else}
                <button
                    type="button"
                    class="btn  btn-sm btn-transparent btn-hint lock-toggle"
                    on:click={lock}
                >
                    <i class="ri-lock-line" />
                    <span class="txt">Set Admins only</span>
                </button>
            {/if}
        </label>

        <svelte:component
            this={ruleInputComponent}
            id={uniqueId}
            bind:this={editorRef}
            bind:value={rule}
            baseCollection={collection}
            disabled={isAdminOnly}
        />

        <div class="help-block">
            <slot {isAdminOnly}>
                <p>
                    {#if isAdminOnly}
                        Only admins will be able to perform this action (
                        <button type="button" class="link-primary" on:click={unlock}>unlock to change</button>
                        ).
                    {:else}
                        Leave empty to grant everyone access.
                    {/if}
                </p>
            </slot>
        </div>
    </Field>
{/if}

<style>
    label .label {
        margin: -5px 0;
        background: rgba(53, 71, 104, 0.12);
    }
    .lock-toggle {
        position: absolute;
        right: 0px;
        top: 0px;
        min-width: 135px;
        padding: 10px;
        border-top-left-radius: 0;
        border-bottom-right-radius: 0;
        background: rgba(53, 71, 104, 0.09);
    }
</style>
