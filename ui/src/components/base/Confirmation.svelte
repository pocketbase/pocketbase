<script>
    import { tick } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { confirmation, resetConfirmation } from "@/stores/confirmation";
    import { _ } from '@/services/i18n';

    let confirmationPopup;
    let isConfirmationBusy = false;

    $: if ($confirmation?.text) {
        confirmationPopup?.show();
    }
</script>

<OverlayPanel
    bind:this={confirmationPopup}
    class="confirm-popup hide-content overlay-panel-sm"
    overlayClose={!isConfirmationBusy}
    escClose={!isConfirmationBusy}
    btnClose={false}
    popup
    on:hide={async () => {
        if ($confirmation?.noCallback) {
            $confirmation.noCallback();
        }
        await tick();
        resetConfirmation();
    }}
>
    <h4 class="block center txt-break" slot="header">{$confirmation?.text}</h4>

    <svelte:fragment slot="footer">
        <!-- svelte-ignore a11y-autofocus -->
        <button
            autofocus
            type="button"
            class="btn btn-secondary btn-expanded-sm"
            disabled={isConfirmationBusy}
            on:click={() => {
                if ($confirmation?.noCallback) {
                    $confirmation.noCallback();
                }
                confirmationPopup?.hide();
            }}
        >
            <span class="txt">{$_("app.base.no")}</span>
        </button>
        <button
            type="button"
            class="btn btn-danger btn-expanded"
            class:btn-loading={isConfirmationBusy}
            disabled={isConfirmationBusy}
            on:click={async () => {
                if ($confirmation?.yesCallback) {
                    isConfirmationBusy = true;
                    await Promise.resolve($confirmation.yesCallback());
                    isConfirmationBusy = false;
                }
                confirmationPopup?.hide();
            }}
        >
            <span class="txt">{$_("app.base.yes")}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
