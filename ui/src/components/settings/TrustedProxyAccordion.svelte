<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    const commonProxyHeaders = ["X-Forwarded-For", "Fly-Client-IP", "CF-Connecting-IP"];

    export let formSettings;
    export let healthData;

    let initialSettingsHash = "";

    $: settingsHash = JSON.stringify(formSettings);

    $: if (initialSettingsHash != settingsHash) {
        initialSettingsHash = settingsHash;
    }

    $: hasChanges = initialSettingsHash != settingsHash;

    $: hasErrors = !CommonHelper.isEmpty($errors?.trustedProxy);

    $: isEnabled = !CommonHelper.isEmpty(formSettings.trustedProxy.headers);

    $: suggestedProxyHeaders = !healthData.possibleProxyHeader
        ? commonProxyHeaders
        : [healthData.possibleProxyHeader].concat(
              commonProxyHeaders.filter((h) => h != healthData.possibleProxyHeader),
          );

    function setHeader(val) {
        formSettings.trustedProxy.headers = [val];
    }

    const ipOptions = [
        { label: "Use leftmost IP", value: true },
        { label: "Use rightmost IP", value: false },
    ];
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-route-line"></i>
            <span class="txt">User IP proxy headers</span>
            {#if !isEnabled && healthData.possibleProxyHeader}
                <i
                    class="ri-alert-line txt-sm txt-warning"
                    use:tooltip={"Detected proxy header.\nIt is recommend to list it as trusted."}
                />
            {:else if isEnabled && !hasChanges && !formSettings.trustedProxy.headers.includes(healthData.possibleProxyHeader)}
                <i
                    class="ri-alert-line txt-sm txt-hint"
                    use:tooltip={"The configured proxy header doesn't match with the detected one."}
                />
            {/if}
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

    <div class="alert alert-info m-b-sm">
        <div class="content">
            <div class="inline-flex flex-gap-5">
                <span>Resolved user IP:</span>
                <strong>{healthData.realIP || "N/A"}</strong>
                <i
                    class="ri-information-line txt-sm link-hint"
                    use:tooltip={"Must show your actual IP.\nIf not, set the correct proxy header."}
                />
            </div>
            <br />
            <div class="inline-flex flex-gap-5">
                <span>Detected proxy header:</span>
                <strong>{healthData.possibleProxyHeader || "N/A"}</strong>
            </div>
        </div>
    </div>

    <div class="content m-b-sm">
        <p>
            When PocketBase is deployed on platforms like Fly or it is accessible through proxies such as
            NGINX, requests from different users will originate from the same IP address (the IP of the proxy
            connecting to your PocketBase app).
        </p>
        <p>
            In this case to retrieve the actual user IP (used for rate limiting, logging, etc.) you need to
            properly configure your proxy and list below the trusted headers that PocketBase could use to
            extract the user IP.
        </p>
        <p class="txt-bold">When using such proxy, to avoid spoofing it is recommended to:</p>
        <ul class="m-t-0 txt-bold">
            <li>use headers that are controlled only by the proxy and cannot be manually set by the users</li>
            <li>make sure that the PocketBase server can be accessed only through the proxy</li>
        </ul>
        <p>You can clear the headers field if PocketBase is not deployed behind a proxy.</p>
    </div>

    <div class="grid grid-sm">
        <div class="col-lg-9">
            <Field class="form-field m-b-0" name="trustedProxy.headers" let:uniqueId>
                <label for={uniqueId}>Trusted proxy headers</label>
                <MultipleValueInput
                    id={uniqueId}
                    placeholder="Leave empty to disable"
                    bind:value={formSettings.trustedProxy.headers}
                />
                <div class="form-field-addon">
                    <button
                        type="button"
                        class="btn btn-sm btn-hint btn-transparent btn-clear"
                        class:hidden={CommonHelper.isEmpty(formSettings.trustedProxy.headers)}
                        on:click={() => (formSettings.trustedProxy.headers = [])}
                    >
                        Clear
                    </button>
                </div>
                <div class="help-block">
                    <p>
                        Comma separated list of headers such as:
                        {#each suggestedProxyHeaders as header}
                            <button
                                type="button"
                                class="label label-sm link-primary txt-mono"
                                on:click={() => setHeader(header)}
                            >
                                {header}
                            </button>&nbsp;
                        {/each}
                    </p>
                </div>
            </Field>
        </div>
        <div class="col-lg-3">
            <Field class="form-field m-0" name="trustedProxy.useLeftmostIP" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">IP priority selection</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "This is in case the proxy returns more than 1 IP as header value. The rightmost IP is usually considered to be the more trustworthy but this could vary depending on the proxy.",
                            position: "right",
                        }}
                    />
                </label>
                <ObjectSelect
                    items={ipOptions}
                    bind:keyOfSelected={formSettings.trustedProxy.useLeftmostIP}
                />
            </Field>
        </div>
    </div>
</Accordion>
