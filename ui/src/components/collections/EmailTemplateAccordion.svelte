<script context="module">
    let cachedEditorComponent;
</script>

<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import { errors, removeError } from "@/stores/errors";
    import { addInfoToast } from "@/stores/toasts";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    export let key;
    export let title;
    export let config = {};
    export let placeholders = [];

    let accordion;
    let editorComponent = cachedEditorComponent;
    let isEditorComponentLoading = false;

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

    async function loadEditorComponent() {
        if (editorComponent || isEditorComponentLoading) {
            return; // already loaded or in the process
        }

        isEditorComponentLoading = true;

        editorComponent = (await import("@/components/base/CodeEditor.svelte")).default;

        cachedEditorComponent = editorComponent;

        isEditorComponentLoading = false;
    }

    function copy(param) {
        param = param.replace("*", ""); // strip wildcard
        CommonHelper.copyToClipboard(param);
        addInfoToast(`Copied ${param} to clipboard`, 2000);
    }

    loadEditorComponent();
</script>

<Accordion bind:this={accordion} on:expand on:collapse on:toggle {...$$restProps}>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-draft-line" />
            <span class="txt">{title}</span>
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

    <Field class="form-field required" name="{key}.subject" let:uniqueId>
        <label for={uniqueId}>Subject</label>
        <input type="text" id={uniqueId} bind:value={config.subject} spellcheck="false" required />
        {#if placeholders?.length > 0}
            <div class="help-block">
                Available placeholder parameters:
                {#each placeholders as placeholder}
                    <button
                        type="button"
                        class="label label-sm link-primary txt-mono"
                        on:click={() => copy("{" + placeholder + "}")}
                    >
                        {"{" + placeholder + "}"}
                    </button>&nbsp;
                {/each}
            </div>
        {/if}
    </Field>

    <Field class="form-field m-0 required" name="{key}.body" let:uniqueId>
        <label for={uniqueId}>Body (HTML)</label>

        {#if editorComponent && !isEditorComponentLoading}
            <svelte:component this={editorComponent} id={uniqueId} language="html" bind:value={config.body} />
        {:else}
            <textarea
                id={uniqueId}
                class="txt-mono"
                spellcheck="false"
                rows="14"
                required
                bind:value={config.body}
            />
        {/if}

        {#if placeholders?.length > 0}
            <div class="help-block">
                Available placeholder parameters:
                {#each placeholders as placeholder}
                    <button
                        type="button"
                        class="label label-sm link-primary txt-mono"
                        on:click={() => copy("{" + placeholder + "}")}
                    >
                        {"{" + placeholder + "}"}
                    </button>&nbsp;
                {/each}
            </div>
        {/if}
    </Field>
</Accordion>
