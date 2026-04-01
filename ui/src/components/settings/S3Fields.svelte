<script>
    import { onMount } from "svelte";
    import { slide } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { removeError } from "@/stores/errors";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";

    const testRequestKey = "s3_test_request";

    const PROVIDERS = [
        {
            label: "AWS S3",
            value: "s3",
            endpointPlaceholder: "https://s3.amazonaws.com",
            regionPlaceholder: "us-east-1",
            accessKeyLabel: "Access key",
            secretLabel: "Secret",
            defaults: {},
        },
        {
            label: "Backblaze B2",
            value: "b2",
            endpointPlaceholder: "https://s3.{region}.backblazeb2.com",
            regionPlaceholder: "us-west-004",
            accessKeyLabel: "Application Key ID",
            secretLabel: "Application Key",
            defaults: {
                forcePathStyle: true,
            },
        },
    ];

    export let originalConfig = {};
    export let config = {};
    export let configKey = "s3";
    export let toggleLabel = "Enable S3";
    export let testFilesystem = "storage"; // storage or backups
    export let testError = null;
    export let isTesting = false;

    let testTimeoutId = null;
    let testDebounceId = null;
    let maskSecret = false;
    let selectedProvider = "s3";
    let prevProvider = "s3";
    let prevRegion = "";

    $: currentPreset = PROVIDERS.find((p) => p.value === selectedProvider) || PROVIDERS[0];

    // react to provider changes
    $: if (selectedProvider !== prevProvider) {
        onProviderChange(selectedProvider);
        prevProvider = selectedProvider;
    }

    // auto-fill endpoint when region changes for B2
    $: if (selectedProvider === "b2" && config.region && config.region !== prevRegion) {
        if (!config.endpoint || /^https:\/\/s3\.[^.]*\.backblazeb2\.com$/.test(config.endpoint)) {
            config.endpoint = `https://s3.${config.region}.backblazeb2.com`;
        }
        prevRegion = config.region;
    }

    function onProviderChange(newProvider) {
        const preset = PROVIDERS.find((p) => p.value === newProvider);
        if (!preset) return;

        for (const [key, value] of Object.entries(preset.defaults)) {
            config[key] = value;
        }

        if (newProvider === "b2" && config.region && !config.endpoint) {
            config.endpoint = `https://s3.${config.region}.backblazeb2.com`;
        }
    }

    $: if (originalConfig?.enabled) {
        refreshMaskSecret();
        testConnectionWithDebounce(100);
    }

    // clear s3 errors on disable
    $: if (!config.enabled) {
        removeError(configKey);
    }

    function refreshMaskSecret() {
        maskSecret = !!originalConfig?.accessKey;
    }

    function testConnectionWithDebounce(timeout) {
        isTesting = true;
        clearTimeout(testDebounceId);
        testDebounceId = setTimeout(() => {
            testConnection();
        }, timeout);
    }

    async function testConnection() {
        testError = null;

        if (!config.enabled) {
            isTesting = false;
            return testError; // nothing to test
        }

        // auto cancel the test request after 30sec
        ApiClient.cancelRequest(testRequestKey);
        clearTimeout(testTimeoutId);
        testTimeoutId = setTimeout(() => {
            ApiClient.cancelRequest(testRequestKey);
            testError = new Error("S3 test connection timeout.");
            isTesting = false;
        }, 30000);

        isTesting = true;

        let err;

        try {
            await ApiClient.settings.testS3(testFilesystem, {
                $cancelKey: testRequestKey,
            });
        } catch (e) {
            err = e;
        }

        if (!err?.isAbort) {
            testError = err;
            isTesting = false;
            clearTimeout(testTimeoutId);
        }

        return testError;
    }

    onMount(() => {
        if (config.endpoint?.includes("backblazeb2.com")) {
            selectedProvider = "b2";
        }
        prevProvider = selectedProvider;
        prevRegion = config.region || "";

        return () => {
            clearTimeout(testTimeoutId);
            clearTimeout(testDebounceId);
        };
    });
</script>

<Field class="form-field form-field-toggle" let:uniqueId>
    <input type="checkbox" id={uniqueId} required bind:checked={config.enabled} />
    <label for={uniqueId}>{toggleLabel}</label>
</Field>

<slot {isTesting} {testError} enabled={config.enabled} />

{#if config.enabled}
    <div class="grid" transition:slide={{ duration: 150 }}>
        <div class="col-lg-12">
            <Field class="form-field" let:uniqueId>
                <label for={uniqueId}>Provider</label>
                <ObjectSelect id={uniqueId} items={PROVIDERS} bind:keyOfSelected={selectedProvider} />
            </Field>
        </div>
        <div class="col-lg-6">
            <Field class="form-field required" name="{configKey}.endpoint" let:uniqueId>
                <label for={uniqueId}>Endpoint</label>
                <input type="text" id={uniqueId} required placeholder={currentPreset.endpointPlaceholder} bind:value={config.endpoint} />
            </Field>
        </div>
        <div class="col-lg-3">
            <Field class="form-field required" name="{configKey}.bucket" let:uniqueId>
                <label for={uniqueId}>Bucket</label>
                <input type="text" id={uniqueId} required bind:value={config.bucket} />
            </Field>
        </div>
        <div class="col-lg-3">
            <Field class="form-field required" name="{configKey}.region" let:uniqueId>
                <label for={uniqueId}>Region</label>
                <input type="text" id={uniqueId} required placeholder={currentPreset.regionPlaceholder} bind:value={config.region} />
            </Field>
        </div>
        <div class="col-lg-6">
            <Field class="form-field required" name="{configKey}.accessKey" let:uniqueId>
                <label for={uniqueId}>{currentPreset.accessKeyLabel}</label>
                <input type="text" id={uniqueId} required bind:value={config.accessKey} />
            </Field>
        </div>
        <div class="col-lg-6">
            <Field class="form-field required" name="{configKey}.secret" let:uniqueId>
                <label for={uniqueId}>{currentPreset.secretLabel}</label>
                <RedactedPasswordInput
                    required
                    id={uniqueId}
                    bind:mask={maskSecret}
                    bind:value={config.secret}
                />
            </Field>
        </div>
        <div class="col-lg-12">
            <Field class="form-field" name="{configKey}.forcePathStyle" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={config.forcePathStyle} />
                <label for={uniqueId}>
                    <span class="txt">Force path-style addressing</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: 'Forces the request to use path-style addressing, eg. "https://s3.amazonaws.com/BUCKET/KEY" instead of the default "https://BUCKET.s3.amazonaws.com/KEY".',
                            position: "top",
                        }}
                    />
                </label>
            </Field>
        </div>
        <!-- margin helper -->
        <div class="col-lg-12" />
    </div>
{/if}
