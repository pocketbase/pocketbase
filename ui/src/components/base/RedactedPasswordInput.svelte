<script>
    import { tick } from "svelte";
    import tooltip from "@/actions/tooltip";

    export let value = undefined; // note must be undefined so that it can be skipped from the submit data
    export let mask = false;

    let inputElem;

    async function unlock() {
        value = "";
        mask = false;
        await tick();
        inputElem?.focus();
    }
</script>

{#if mask}
    <div class="form-field-addon">
        <button
            type="button"
            class="btn btn-transparent btn-circle"
            use:tooltip={{ position: "left", text: "Set new value" }}
            on:click|preventDefault={unlock}
        >
            <i class="ri-key-line" />
        </button>
    </div>

    <input disabled type="text" placeholder="******" {...$$restProps} />
{:else}
    <input bind:this={inputElem} bind:value type="password" autocomplete="new-password" {...$$restProps} />
{/if}
