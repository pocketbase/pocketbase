<script>
    import { scale } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let identityFieldsOptions = [];
    let oldIndexes = "";

    $: isSuperusers = collection?.system && collection?.name === "_superusers";

    $: if (CommonHelper.isEmpty(collection?.passwordAuth)) {
        collection.passwordAuth = {
            enabled: true,
            identityFields: ["email"],
        };
    }

    $: hasErrors = !CommonHelper.isEmpty($errors?.passwordAuth);

    $: if (collection && oldIndexes != collection.indexes.join("")) {
        refreshIdentityFieldsOptions();
    }

    function refreshIdentityFieldsOptions() {
        // email is always available in auth collections
        identityFieldsOptions = [{ value: "email" }];

        const fields = collection?.fields || [];
        const indexes = collection?.indexes || [];

        oldIndexes = indexes.join("");

        for (let idx of indexes) {
            const parsed = CommonHelper.parseIndex(idx);
            if (!parsed.unique || parsed.columns.length != 1 || parsed.columns[0].name == "email") {
                continue;
            }

            const field = fields.find((f) => {
                return !f.hidden && f.name.toLowerCase() == parsed.columns[0].name.toLowerCase();
            });
            if (field) {
                identityFieldsOptions.push({ value: field.name });
            }
        }
    }
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-lock-password-line"></i>
            <span class="txt">Identity/Password</span>
        </div>

        <div class="flex-fill" />

        {#if collection.passwordAuth.enabled}
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

    <Field class="form-field form-field-toggle" name="passwordAuth.enabled" let:uniqueId>
        <input
            type="checkbox"
            id={uniqueId}
            bind:checked={collection.passwordAuth.enabled}
            disabled={isSuperusers}
        />
        <label for={uniqueId}>Enable</label>
        {#if isSuperusers}
            <i
                class="ri-information-line link-hint"
                use:tooltip={{
                    text: "Superusers are required to have password auth enabled.",
                    position: "right",
                }}
            />
        {/if}
    </Field>

    <Field class="form-field required m-0" name="passwordAuth.identityFields" let:uniqueId>
        <label for={uniqueId}>
            <span class="txt">Unique identity fields</span>
        </label>
        <ObjectSelect
            items={identityFieldsOptions}
            multiple
            bind:keyOfSelected={collection.passwordAuth.identityFields}
        />
    </Field>
</Accordion>
