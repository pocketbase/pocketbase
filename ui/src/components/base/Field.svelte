<script>
    import { onMount } from "svelte";
    import { slide, scale } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import { errors, removeError } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";

    const uniqueId = "field_" + CommonHelper.randomString(7);
    const defaultError = "Invalid value";

    export let name = "";
    export let inlineError = false;

    let classes = undefined;
    export { classes as class }; // export reserved keyword

    let container;
    let fieldErrors = [];

    $: {
        fieldErrors = CommonHelper.toArray(CommonHelper.getNestedVal($errors, name));
    }

    export function changed() {
        removeError(name);
    }

    onMount(() => {
        container.addEventListener("input", changed);
        container.addEventListener("change", changed);

        return () => {
            container.removeEventListener("input", changed);
            container.removeEventListener("change", changed);
        };
    });

    function getErrorMessage(err) {
        if (typeof err === "object") {
            return err?.message || err?.code || defaultError;
        }

        return err || defaultError;
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div bind:this={container} class={classes} class:error={fieldErrors.length} on:click>
    <slot {uniqueId} />

    {#if inlineError && fieldErrors.length}
        <div class="form-field-addon inline-error-icon">
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{
                    position: "left",
                    text: fieldErrors.map(getErrorMessage).join("\n"),
                }}
            />
        </div>
    {:else}
        {#each fieldErrors as error}
            <div class="help-block help-block-error" transition:slide={{ duration: 150 }}>
                <pre>{getErrorMessage(error)}</pre>
            </div>
        {/each}
    {/if}
</div>
