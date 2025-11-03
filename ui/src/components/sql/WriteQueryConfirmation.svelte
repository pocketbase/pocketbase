<script>
    import { createEventDispatcher } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    export let active = false;

    function confirm() {
        dispatch("confirm");
        active = false;
    }

    function cancel() {
        dispatch("cancel");
        active = false;
    }
</script>

<OverlayPanel bind:active popup class="overlay-panel-sm" on:hide={cancel}>
    <svelte:fragment slot="header">
        <h4>Enable Write Queries</h4>
    </svelte:fragment>

    <div class="content">
        <div class="alert alert-warning m-b-base">
            <div class="icon">
                <i class="ri-alert-line" />
            </div>
            <div class="content">
                <p>
                    <strong>Warning:</strong> Enabling write queries allows you to execute
                    <code>INSERT</code>, <code>UPDATE</code>, <code>DELETE</code>, and other
                    database modification statements.
                </p>
                <p class="m-t-10">
                    This can permanently modify or delete data in your database. Please
                    proceed with caution.
                </p>
            </div>
        </div>

        <p class="txt-hint">
            Write mode will automatically be disabled after each query execution for safety.
        </p>
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={cancel}>
            <span class="txt">Cancel</span>
        </button>
        <button type="button" class="btn btn-expanded btn-warning" on:click={confirm}>
            <i class="ri-lock-unlock-line" />
            <span class="txt">Enable Writes</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

