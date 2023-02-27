<script>
    import { onMount } from "svelte";
    import { errors, removeError } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";

    const uniqueId = "field_" + CommonHelper.randomString(7);
    const defaultError = "Invalid value";

    export let name = "";

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
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div bind:this={container} class={classes} class:error={fieldErrors.length} on:click>
    <slot {uniqueId} />

    {#each fieldErrors as error}
        <div class="help-block help-block-error">
            {#if typeof error === "object"}
                <pre>{error?.message || error?.code || defaultError}</pre>
            {:else}
                {error || defaultError}
            {/if}
        </div>
    {/each}
</div>
