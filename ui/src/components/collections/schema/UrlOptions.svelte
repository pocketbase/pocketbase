<script>
    import EmailOptions from "./EmailOptions.svelte";

    export let key = "";
    export let options = {};
</script>

<!-- shares the same options with the email field -->
<EmailOptions bind:key bind:options />
