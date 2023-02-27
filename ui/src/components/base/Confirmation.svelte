<script>
    import { tick } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { confirmation, resetConfirmation } from "@/stores/confirmation";

    let confirmationPopup;
    let isConfirmationBusy = false;
    let confirmed = false;

    $: if ($confirmation?.text) {
        confirmed = false;
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
        if (!confirmed && $confirmation?.noCallback) {
            $confirmation.noCallback();
        }
        await tick();
        confirmed = false;
        resetConfirmation();
    }}
>
    <h4 class="block center txt-break" slot="header">{$confirmation?.text}</h4>

    <svelte:fragment slot="footer">
        <!-- svelte-ignore a11y-autofocus -->
        <button
            autofocus
            type="button"
            class="btn btn-transparent btn-expanded-sm"
            disabled={isConfirmationBusy}
            on:click={() => {
                confirmed = false;
                confirmationPopup?.hide();
            }}
        >
            <span class="txt">No</span>
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
                confirmed = true;
                confirmationPopup?.hide();
            }}
        >
            <span class="txt">Yes</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
