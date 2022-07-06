<script>
    // example supported format: {label: "...", value: "...", icon: ""}
    export let item = {};
</script>

{#if item.icon}
    <i class="icon {item.icon}" />
{/if}

<span class="txt">{item.label || item.name || item.title || item.id || item.value}</span>
