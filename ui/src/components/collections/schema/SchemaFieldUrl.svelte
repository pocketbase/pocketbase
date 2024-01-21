<script>
    import SchemaFieldEmail from "./SchemaFieldEmail.svelte";

    export let field;
    export let key = "";
</script>

<!-- shares the same options with the email field -->
<SchemaFieldEmail bind:field {key} on:rename on:remove on:duplicate {...$$restProps} />
