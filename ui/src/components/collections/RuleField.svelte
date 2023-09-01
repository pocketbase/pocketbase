<script context="module">
    let cachedRuleComponent;
</script>

<script>
    import { tick } from "svelte";
    import { scale } from "svelte/transition";
    import Field from "@/components/base/Field.svelte";

    export let collection = null;
    export let rule = null;
    export let label = "Rule";
    export let formKey = "rule";
    export let required = false;
    export let placeholder = "Leave empty to grant everyone access...";

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
        class="form-field rule-field {required ? 'requied' : ''} {isAdminOnly ? 'disabled' : ''}"
        name={formKey}
        let:uniqueId
    >
        <div class="input-wrapper">
            <label for={uniqueId}>
                <slot name="beforeLabel" {isAdminOnly} />

                <span class="txt" class:txt-hint={isAdminOnly}>
                    {label}
                    {isAdminOnly ? "- Admins only" : ""}
                </span>

                <slot name="afterLabel" {isAdminOnly} />

                {#if !isAdminOnly}
                    <button
                        type="button"
                        class="btn btn-sm btn-transparent btn-hint lock-toggle"
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
                placeholder={!isAdminOnly ? placeholder : ""}
            />

            {#if isAdminOnly}
                <button
                    type="button"
                    class="unlock-overlay"
                    aria-label="Unlock and set custom rule"
                    transition:scale={{ duration: 150, start: 0.98 }}
                    on:click={unlock}
                >
                    <small class="txt">Unlock and set custom rule</small>
                    <div class="icon">
                        <i class="ri-lock-unlock-line" />
                    </div>
                </button>
            {/if}
        </div>

        <div class="help-block">
            <slot {isAdminOnly} />
        </div>
    </Field>
{/if}

<style lang="scss">
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
    :global(.rule-field .code-editor .cm-placeholder) {
        font-family: var(--baseFontFamily);
    }
    .input-wrapper {
        position: relative;
    }
    .unlock-overlay {
        --hoverAnimationSpeed: 0.2s;
        position: absolute;
        z-index: 1;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
        display: flex;
        padding: 20px;
        gap: 10px;
        align-items: center;
        justify-content: end;
        text-align: center;
        border-radius: var(--baseRadius);
        background: rgba(255, 255, 255, 0.2);
        outline: 0;
        cursor: pointer;
        text-decoration: none;
        color: var(--successColor);
        border: 2px solid var(--baseAlt1Color);
        transition: border-color var(--baseAnimationSpeed);
        i {
            font-size: inherit;
        }
        .icon {
            color: var(--successColor);
            font-size: 1.15rem;
            line-height: 1;
            font-weight: normal;
            transition: transform var(--hoverAnimationSpeed);
        }
        .txt {
            opacity: 0;
            font-size: var(--xsFontSize);
            font-weight: 600;
            line-height: var(--smLineHeight);
            transform: translateX(5px);
            transition: transform var(--hoverAnimationSpeed), opacity var(--hoverAnimationSpeed);
        }
        &:hover,
        &:focus-visible,
        &:active {
            border-color: var(--baseAlt3Color);
            .icon {
                transform: scale(1.1);
            }
            .txt {
                opacity: 1;
                transform: scale(1);
            }
        }
        &:active {
            transition-duration: var(--activeAnimationSpeed);
            border-color: var(--baseAlt3Color);
        }
    }
</style>
