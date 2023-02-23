<script>
    import { tick } from "svelte";
    import tooltip from "@/actions/tooltip";

    export let value = "";
    export let mask = "******";

    let inputElem;
    let locked = false;

    $: locked = value === mask;

    async function unlock() {
        value = "";
        locked = false;
        await tick();
        inputElem?.focus();
    }
</script>

{#if locked}
    <div class="form-field-addon">
        <button
            type="button"
            class="btn btn-transparent btn-circle"
            use:tooltip={{ position: "left", text: "Set new value" }}
            on:click={() => unlock()}
        >
            <i class="ri-key-line" />
        </button>
    </div>

    <input readonly type="text" placeholder={mask} {...$$restProps} />
{:else}
    <input bind:this={inputElem} bind:value type="password" autocomplete="new-password" {...$$restProps} />
{/if}
