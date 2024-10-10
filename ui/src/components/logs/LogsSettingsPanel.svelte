<script>
    import { createEventDispatcher } from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import LogsLevelsInfo from "@/components/logs/LogsLevelsInfo.svelte";

    const dispatch = createEventDispatcher();
    const formId = "logs_settings_" + CommonHelper.randomString(3);

    let panel;
    let isSaving = false;
    let isLoading = false;
    let originalFormSettings = {};
    let formSettings = {};

    $: initialHash = JSON.stringify(originalFormSettings);

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    export function show() {
        reset();

        loadSettings();

        return panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function reset() {
        setErrors();
        originalFormSettings = {};
        formSettings = JSON.parse(JSON.stringify(originalFormSettings || {}));
    }

    async function loadSettings() {
        isLoading = true;

        try {
            const settings = (await ApiClient.settings.getAll()) || {};
            init(settings);
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    async function save() {
        if (!hasChanges) {
            return;
        }

        isSaving = true;

        try {
            const settings = await ApiClient.settings.update(CommonHelper.filterRedactedProps(formSettings));
            init(settings);

            isSaving = false;

            hide();

            addSuccessToast("Successfully saved logs settings.");

            dispatch("save", settings);
        } catch (err) {
            isSaving = false;
            ApiClient.error(err);
        }
    }

    function init(settings = {}) {
        formSettings = {
            logs: settings?.logs || {},
        };

        originalFormSettings = JSON.parse(JSON.stringify(formSettings));
    }
</script>

<OverlayPanel bind:this={panel} popup class="superuser-panel" beforeHide={() => !isSaving} on:hide on:show>
    <svelte:fragment slot="header">
        <h4>Logs settings</h4>
    </svelte:fragment>

    {#if isLoading}
        <div class="block txt-center">
            <div class="loader" />
        </div>
    {:else}
        <form id={formId} class="grid" autocomplete="off" on:submit|preventDefault={save}>
            <Field class="form-field required" name="logs.maxDays" let:uniqueId>
                <label for={uniqueId}>Max days retention</label>
                <input type="number" id={uniqueId} required bind:value={formSettings.logs.maxDays} />
                <div class="help-block">
                    Set to <code>0</code> to disable logs persistence.
                </div>
            </Field>

            <Field class="form-field" name="logs.minLevel" let:uniqueId>
                <label for={uniqueId}>Min log level</label>
                <input type="number" required bind:value={formSettings.logs.minLevel} min="-100" max="100" />
                <div class="help-block">
                    <p>Logs with level below the minimum will be ignored.</p>
                    <LogsLevelsInfo />
                </div>
            </Field>

            <Field class="form-field form-field-toggle" name="logs.logIP" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={formSettings.logs.logIP} />
                <label for={uniqueId}>Enable IP logging</label>
            </Field>

            <Field class="form-field form-field-toggle" name="logs.logAuthId" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={formSettings.logs.logAuthId} />
                <label for={uniqueId}>Enable Auth Id logging</label>
            </Field>
        </form>
    {/if}

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" disabled={isSaving} on:click={hide}>
            <span class="txt">Cancel</span>
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!hasChanges || isSaving}
        >
            <span class="txt">Save changes</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
