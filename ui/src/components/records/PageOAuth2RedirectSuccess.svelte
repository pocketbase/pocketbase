<script>
    import { onMount } from "svelte";
    import { pageTitle } from "@/stores/app";

    $pageTitle = "OAuth2 auth completed";

    onMount(() => {
        window.close();
    });
</script>

<div class="content txt-hint txt-center p-base">
    <h3 class="m-b-sm">Auth completed.</h3>
    <h5>You can close this window and go back to the app.</h5>
</div>
