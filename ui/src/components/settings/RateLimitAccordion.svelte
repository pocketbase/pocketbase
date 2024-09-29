<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import AutocompleteInput from "@/components/base/AutocompleteInput.svelte";
    import { errors, setErrors } from "@/stores/errors";
    import { collections, loadCollections } from "@/stores/collections";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    export let formSettings;

    const basePredefinedTags = [
        { value: "*:list" },
        { value: "*:view" },
        { value: "*:create" },
        { value: "*:update" },
        { value: "*:delete" },
        { value: "*:file" },
        { value: "*:listAuthMethods" },
        { value: "*:authRefresh" },
        { value: "*:auth" },
        { value: "*:authWithPassword" },
        { value: "*:authWithOAuth2" },
        { value: "*:authWithOTP" },
        { value: "*:requestOTP" },
        { value: "*:requestPasswordReset" },
        { value: "*:confirmPasswordReset" },
        { value: "*:requestVerification" },
        { value: "*:confirmVerification" },
        { value: "*:requestEmailChange" },
        { value: "*:confirmEmailChange" },
    ];

    let predefinedTags = basePredefinedTags;

    $: hasErrors = !CommonHelper.isEmpty($errors?.rateLimits);

    loadPredefinedTags();

    async function loadPredefinedTags() {
        await loadCollections();

        predefinedTags = [];

        for (let collection of $collections) {
            if (collection.system) {
                continue;
            }

            predefinedTags.push({ value: collection.name + ":list" });
            predefinedTags.push({ value: collection.name + ":view" });

            if (collection.type != "view") {
                predefinedTags.push({ value: collection.name + ":create" });
                predefinedTags.push({ value: collection.name + ":update" });
                predefinedTags.push({ value: collection.name + ":delete" });
            }

            if (collection.type == "auth") {
                predefinedTags.push({ value: collection.name + ":listAuthMethods" });
                predefinedTags.push({ value: collection.name + ":authRefresh" });
                predefinedTags.push({ value: collection.name + ":auth" });
                predefinedTags.push({ value: collection.name + ":authWithPassword" });
                predefinedTags.push({ value: collection.name + ":authWithOAuth2" });
                predefinedTags.push({ value: collection.name + ":authWithOTP" });
                predefinedTags.push({ value: collection.name + ":requestOTP" });
                predefinedTags.push({ value: collection.name + ":requestPasswordReset" });
                predefinedTags.push({ value: collection.name + ":confirmPasswordReset" });
                predefinedTags.push({ value: collection.name + ":requestVerification" });
                predefinedTags.push({ value: collection.name + ":confirmVerification" });
                predefinedTags.push({ value: collection.name + ":requestEmailChange" });
                predefinedTags.push({ value: collection.name + ":confirmEmailChange" });
            }

            if (collection.fields.find((f) => f.type == "file")) {
                predefinedTags.push({ value: collection.name + ":file" });
            }
        }

        predefinedTags = predefinedTags.concat(basePredefinedTags);
    }

    function newRule() {
        setErrors({}); // reset

        if (!Array.isArray(formSettings.rateLimits.rules)) {
            formSettings.rateLimits.rules = [];
        }

        formSettings.rateLimits.rules.push({
            label: "",
            maxRequests: 300,
            duration: 10,
        });

        formSettings.rateLimits.rules = formSettings.rateLimits.rules;

        if (formSettings.rateLimits.rules.length == 1) {
            formSettings.rateLimits.enabled = true;
        }
    }

    function removeRule(i) {
        setErrors({}); // reset

        formSettings.rateLimits.rules.splice(i, 1);
        formSettings.rateLimits.rules = formSettings.rateLimits.rules;

        if (!formSettings.rateLimits.rules.length) {
            formSettings.rateLimits.enabled = false;
        }
    }
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-pulse-fill"></i>
            <span class="txt">Rate limiting</span>
        </div>

        <div class="flex-fill" />

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "Has errors", position: "left" }}
            />
        {/if}

        {#if formSettings.rateLimits.enabled}
            <span class="label label-success">Enabled</span>
        {:else}
            <span class="label">Disabled</span>
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle m-b-xs" name="rateLimits.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={formSettings.rateLimits.enabled} />
        <label for={uniqueId}>Enable</label>
    </Field>

    {#if !CommonHelper.isEmpty(formSettings.rateLimits.rules)}
        <table class="rate-limit-table">
            <thead>
                <tr>
                    <th>Rate limit label</th>
                    <th>Max requests (per IP)</th>
                    <th>Interval (in seconds)</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                {#each formSettings.rateLimits.rules || [] as rule, i}
                    <tr class="rate-limit-row">
                        <td class="col-tag">
                            <Field class="form-field" name={"rateLimits.rules." + i + ".label"} inlineError>
                                <AutocompleteInput
                                    required
                                    placeholder="tag (users:create) or path (/api/)"
                                    options={predefinedTags}
                                    bind:value={rule.label}
                                />
                            </Field>
                        </td>
                        <td class="col-requests">
                            <Field
                                class="form-field"
                                name={"rateLimits.rules." + i + ".maxRequests"}
                                inlineError
                            >
                                <input
                                    type="number"
                                    required
                                    placeholder="Max requests*"
                                    min="1"
                                    step="1"
                                    bind:value={rule.maxRequests}
                                />
                            </Field>
                        </td>
                        <td class="col-burst">
                            <Field
                                class="form-field"
                                name={"rateLimits.rules." + i + ".duration"}
                                inlineError
                            >
                                <input
                                    type="number"
                                    required
                                    placeholder="Interval*"
                                    min="1"
                                    step="1"
                                    bind:value={rule.duration}
                                />
                            </Field>
                        </td>
                        <td class="col-action">
                            <button
                                type="button"
                                title="Remove rule"
                                aria-label="Remove rule"
                                class="btn btn-xs btn-circle btn-hint btn-transparent"
                                on:click={() => removeRule(i)}
                            >
                                <i class="ri-close-line"></i>
                            </button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}

    <div class="flex m-t-sm">
        <button
            type="button"
            class="btn btn-sm btn-secondary m-r-auto"
            class:btn-danger={$errors?.rateLimits?.rules?.message}
            on:click={() => newRule()}
        >
            <i class="ri-add-line"></i>
            <span class="txt">Add rate limit rule</span>
        </button>

        <a
            href={import.meta.env.PB_RATE_LIMIT_DOCS}
            class="txt-nowrap txt-sm link-hint"
            target="_blank"
            rel="noopener noreferrer"
        >
            <em>Learn more about the rate limit rules</em>
        </a>
    </div>
</Accordion>
