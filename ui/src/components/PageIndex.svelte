<script>
    import { tick } from "svelte";
    import { replace } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import FullPage from "@/components/base/FullPage.svelte";
    import Installer from "@/components/base/Installer.svelte";

    let showInstaller = false;

    handler();

    function handler() {
        showInstaller = false;

        const realQueryParams = new URLSearchParams(window.location.search);
        if (realQueryParams.has(import.meta.env.PB_INSTALLER_PARAM)) {
            ApiClient.logout(false);
            showInstaller = true;
            return;
        }

        if (ApiClient.authStore.isValid) {
            replace("/collections");
        } else {
            ApiClient.logout();
        }
    }
</script>

{#if showInstaller}
    <FullPage>
        <Installer
            on:submit={async () => {
                showInstaller = false;

                await tick();

                // clear the installer param
                window.location.search = "";
            }}
        />
    </FullPage>
{/if}
