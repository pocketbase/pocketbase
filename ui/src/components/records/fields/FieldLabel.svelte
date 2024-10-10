<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let uniqueId;
    export let field;
    export let icon = true;
</script>

<label for={uniqueId}>
    {#if icon}
        {#if field.primaryKey}
            <i class={CommonHelper.getFieldTypeIcon("primary")} />
        {:else}
            <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        {/if}
    {/if}

    <span class="txt">{field.name}</span>

    {#if field.hidden}
        <small class="label label-sm label-danger">Hidden</small>
    {/if}

    <slot />
</label>
