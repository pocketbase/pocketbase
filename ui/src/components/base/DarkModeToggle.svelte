<script>
    import { darkMode } from "@/stores/darkMode";
    import tooltip from "@/actions/tooltip";
</script>

<button
    type="button"
    class="btn btn-circle btn-transparent btn-hint"
    aria-label="Toggle dark mode"
    use:tooltip={{ text: $darkMode ? "Switch to light mode" : "Switch to dark mode", position: "bottom" }}
    on:click={darkMode.toggle}
>
    {#if $darkMode}
        <i class="ri-sun-line" />
    {:else}
        <i class="ri-moon-line" />
    {/if}
</button>