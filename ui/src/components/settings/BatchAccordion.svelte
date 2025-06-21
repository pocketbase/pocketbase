<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    export let formSettings;

    $: hasErrors = !CommonHelper.isEmpty($errors?.batch);

    $: isEnabled = !!formSettings.batch?.enabled;
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-archive-stack-line"></i>
            <span class="txt">Batch API</span>
        </div>

        <div class="flex-fill" />

        {#if isEnabled}
            <span class="label label-success">Enabled</span>
        {:else}
            <span class="label">Disabled</span>
        {/if}

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "Has errors", position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle m-b-sm" name="batch.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={formSettings.batch.enabled} />
        <label for={uniqueId}>Enable <small class="txt-hint">(experimental)</small></label>
    </Field>

    <div class="grid">
        <div class="col-lg-4">
            <Field class="form-field {isEnabled ? 'required' : ''}" name="batch.maxRequests" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">Max requests in a batch</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "Rate limiting (if enabled) also applies for the batch create/update/upsert/delete requests.",
                            position: "right",
                        }}
                    />
                </label>
                <input
                    type="number"
                    id={uniqueId}
                    min="0"
                    step="1"
                    required={isEnabled}
                    bind:value={formSettings.batch.maxRequests}
                />
            </Field>
        </div>

        <div class="col-lg-4">
            <Field class="form-field {isEnabled ? 'required' : ''}" name="batch.timeout" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">Max processing time (in seconds)</span>
                </label>
                <input
                    type="number"
                    id={uniqueId}
                    min="0"
                    step="1"
                    required={isEnabled}
                    bind:value={formSettings.batch.timeout}
                />
            </Field>
        </div>

        <div class="col-lg-4">
            <Field class="form-field" name="batch.maxBodySize" let:uniqueId>
                <label for={uniqueId}>Max body size (in bytes)</label>
                <input
                    type="number"
                    id={uniqueId}
                    min="0"
                    step="1"
                    placeholder="Default to 128MB"
                    value={formSettings.batch.maxBodySize || ""}
                    on:input={(e) => (formSettings.batch.maxBodySize = e.target.value << 0)}
                />
            </Field>
        </div>
    </div>
</Accordion>
