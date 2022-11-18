<script context="module">
    let cachedRuleComponent;
</script>

<script>
    import { tick } from "svelte";
    import tooltip from "@/actions/tooltip";
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

    async function loadEditorComponent() {
        if (ruleInputComponent || isRuleComponentLoading) {
            return; // already loaded or in the process
        }

        isRuleComponentLoading = true;

        ruleInputComponent = (await import("@/components/base/FilterAutocompleteInput.svelte")).default;

        cachedRuleComponent = ruleInputComponent;

        isRuleComponentLoading = false;
    }

    loadEditorComponent();
</script>

{#if isRuleComponentLoading}
    <div class="txt-center">
        <span class="loader" />
    </div>
{:else}
    <div class="rule-block">
        {#if isAdminOnly}
            <button
                type="button"
                class="rule-toggle-btn btn btn-circle btn-outline btn-success"
                use:tooltip={{
                    text: "Unlock and set custom rule",
                    position: "left",
                }}
                on:click={async () => {
                    rule = tempValue || "";
                    await tick();
                    editorRef?.focus();
                }}
            >
                <i class="ri-lock-unlock-line" />
            </button>
        {:else}
            <button
                type="button"
                class="rule-toggle-btn btn btn-circle btn-outline"
                use:tooltip={{
                    text: "Lock and set to Admins only",
                    position: "left",
                }}
                on:click={() => {
                    tempValue = rule;
                    rule = null;
                }}
            >
                <i class="ri-lock-line" />
            </button>
        {/if}

        <Field
            class="form-field rule-field m-0 {required ? 'requied' : ''} {isAdminOnly ? 'disabled' : ''}"
            name={formKey}
            let:uniqueId
        >
            <label for={uniqueId}>
                {label} - {isAdminOnly ? "Admins only" : "Custom rule"}
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
                            Only admins will be able to perform this action (unlock to change)
                        {:else}
                            Leave empty to grant everyone access
                        {/if}
                    </p>
                </slot>
            </div>
        </Field>
    </div>
{/if}

<style>
    .rule-block {
        display: flex;
        align-items: flex-start;
        gap: var(--xsSpacing);
    }
    .rule-toggle-btn {
        margin-top: 15px;
    }
</style>
