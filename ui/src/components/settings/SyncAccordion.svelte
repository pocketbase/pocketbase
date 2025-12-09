<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";
    export let formSettings;

    $: hasErrors = !CommonHelper.isEmpty($errors?.sync);

    $: isEnabled = !!formSettings.sync?.enabled;
    $: nodeType = formSettings.sync?.nodeType || 'add';

    // Generate instance ID manually
    function generateInstanceID() {
        const bytes = new Uint8Array(16);
        crypto.getRandomValues(bytes);
        bytes[6] = (bytes[6] & 0x0f) | 0x40;
        bytes[8] = (bytes[8] & 0x3f) | 0x80;
        const uuid = Array.from(bytes)
            .map((b, i) => {
                if (i === 4 || i === 6 || i === 8 || i === 10) {
                    return "-" + b.toString(16).padStart(2, "0");
                }
                return b.toString(16).padStart(2, "0");
            })
            .join("");
        formSettings.sync.instanceID = uuid;
    }

</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-sync-line"></i>
            <span class="txt">Sync</span>
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

    <Field class="form-field form-field-toggle m-b-sm" name="sync.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={formSettings.sync.enabled} />
        <label for={uniqueId}>Enable horizontal scaling via NATS JetStream</label>
    </Field>

    {#if isEnabled}
        <div class="m-b-lg">
            <Field class="form-field" name="sync.nodeType" let:uniqueId>
                <label for={uniqueId}>Sync Action</label>
                <div class="form-field-block m-b-sm">
                    <input
                        type="radio"
                        id={uniqueId + 'add'}
                        name={uniqueId}
                        value="add"
                        bind:group={formSettings.sync.nodeType}
                    />
                    <label for={uniqueId + 'add'}>
                        <i class="ri-add-circle-line"></i>
                        <span class="txt">Add To Sync Chain</span>
                    </label>
                </div>
                <div class="form-field-block">
                    <input
                        type="radio"
                        id={uniqueId + 'initial'}
                        name={uniqueId}
                        value="initial"
                        bind:group={formSettings.sync.nodeType}
                    />
                    <label for={uniqueId + 'initial'}>
                        <i class="ri-server-line"></i>
                        <span class="txt">Initial Node</span>
                    </label>
                </div>
                {#if nodeType === 'add'}
                    <div class="form-help">
                        <i class="ri-information-line"></i>
                        This instance will discover and sync with an existing instance in the sync chain.
                    </div>
                {:else}
                    <div class="form-help">
                        <i class="ri-information-line"></i>
                        This instance will be the first node in the sync chain and will publish schema and record changes.
                    </div>
                {/if}
            </Field>
        </div>

        <div class="grid">
            {#if nodeType === 'add'}
                <div class="col-lg-6">
                    <Field class="form-field required" name="sync.targetInstanceAddress" let:uniqueId>
                        <label for={uniqueId}>
                            <span class="txt">Remote Instance Address</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: "IP address and port of the remote instance (e.g., 192.168.1.100:4222)",
                                    position: "right",
                                }}
                            />
                        </label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            placeholder="192.168.1.100:4222"
                            bind:value={formSettings.sync.targetInstanceAddress}
                        />
                    </Field>
                </div>

                <div class="col-lg-6">
                    <Field class="form-field required" name="sync.targetInstanceID" let:uniqueId>
                        <label for={uniqueId}>
                            <span class="txt">Remote Instance ID</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: "Instance ID of the remote node to connect to",
                                    position: "right",
                                }}
                            />
                        </label>
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            placeholder="xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxx"
                            bind:value={formSettings.sync.targetInstanceID}
                        />
                    </Field>
                </div>
            {/if}

            <div class="col-lg-6">
                <Field class="form-field required" name="sync.instanceID" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Local Instance ID</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: "Unique identifier for this local instance",
                                position: "right",
                            }}
                        />
                    </label>
                    <div class="input-group">
                        <input
                            type="text"
                            id={uniqueId}
                            required
                            bind:value={formSettings.sync.instanceID}
                        />
                        <button
                            type="button"
                            class="btn btn-transparent"
                            on:click={() => generateInstanceID()}
                            title="Generate new Local Instance ID"
                        >
                            <i class="ri-refresh-line"></i>
                        </button>
                    </div>
                </Field>
            </div>

            <div class="col-lg-6">
                <Field class="form-field required" name="sync.natsURL" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Local NATS URL</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: "Local NATS server address in format host:port (e.g., 0.0.0.0:4222)",
                                position: "right",
                            }}
                        />
                    </label>
                    <input
                        type="text"
                        id={uniqueId}
                        required
                        placeholder="0.0.0.0:4222"
                        bind:value={formSettings.sync.natsURL}
                    />
                </Field>
            </div>

            <div class="col-lg-6">
                <Field class="form-field required" name="sync.snapshotInterval" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Snapshot Interval (hours)</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: "How often to take snapshots of schema and data for new instance synchronization",
                                position: "right",
                            }}
                        />
                    </label>
                    <input
                        type="number"
                        id={uniqueId}
                        min="1"
                        step="1"
                        required
                        bind:value={formSettings.sync.snapshotInterval}
                    />
                </Field>
            </div>
        </div>


    {/if}
</Accordion>
