<script>
    import { Collection } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";

    export let collectionA = new Collection();
    export let collectionB = new Collection();
    export let deleteMissing = false;

    $: isDeleteDiff = !collectionB?.id && !collectionB?.name;

    $: isCreateDiff = !isDeleteDiff && !collectionA?.id;

    $: schemaA = Array.isArray(collectionA?.schema) ? collectionA?.schema : [];

    $: schemaB = Array.isArray(collectionB?.schema) ? collectionB?.schema : [];

    $: removedFields = schemaA.filter((fieldA) => {
        return !schemaB.find((fieldB) => fieldA.id == fieldB.id);
    });

    $: sharedFields = schemaB.filter((fieldB) => {
        return schemaA.find((fieldA) => fieldA.id == fieldB.id);
    });

    $: addedFields = schemaB.filter((fieldB) => {
        return !schemaA.find((fieldA) => fieldA.id == fieldB.id);
    });

    $: if (typeof deleteMissing !== "undefined") {
        normalizeSchemaB();
    }

    $: hasAnyChange = detectChanges(collectionA, collectionB);

    const mainModelProps = Object.keys(new Collection().export()).filter(
        (key) => !["schema", "created", "updated"].includes(key)
    );

    function normalizeSchemaB() {
        schemaB = Array.isArray(collectionB?.schema) ? collectionB?.schema : [];
        if (!deleteMissing) {
            schemaB = schemaB.concat(removedFields);
        }
    }

    function getFieldById(schema, id) {
        schema = schema || [];

        for (let field of schema) {
            if (field.id == id) {
                return field;
            }
        }
        return null;
    }

    function detectChanges() {
        // added or removed fields
        if (addedFields?.length || (deleteMissing && removedFields?.length)) {
            return true;
        }

        // changes in the main model props
        for (let prop of mainModelProps) {
            if (hasChanges(collectionA?.[prop], collectionB?.[prop])) {
                return true;
            }
        }

        // changes in the schema fields
        for (let field of sharedFields) {
            if (hasChanges(field, CommonHelper.findByKey(schemaA, "id", field.id))) {
                return true;
            }
        }

        return false;
    }

    function hasChanges(valA, valB) {
        // direct match
        if (valA === valB) {
            return false;
        }

        return JSON.stringify(valA) !== JSON.stringify(valB);
    }

    function displayValue(value) {
        if (typeof value === "undefined") {
            return "N/A";
        }

        return CommonHelper.isObject(value) ? JSON.stringify(value, null, 4) : value;
    }
</script>

<div class="section-title">
    {#if !collectionA?.id}
        <strong>{collectionB?.name}</strong>
        <span class="label label-success">Added</span>
    {:else if !collectionB?.id}
        <strong>{collectionA?.name}</strong>
        <span class="label label-danger">Removed</span>
    {:else}
        <div class="inline-flex fleg-gap-5">
            {#if collectionA.name !== collectionB.name}
                <strong class="txt-strikethrough txt-hint">{collectionA.name}</strong>
                <i class="ri-arrow-right-line txt-sm" />
            {/if}
            <strong class="txt">{collectionB.name}</strong>
            {#if hasAnyChange}
                <span class="label label-warning">Changed</span>
            {/if}
        </div>
    {/if}
</div>

<table class="table collections-diff-table m-b-base">
    <thead>
        <tr>
            <th>Props</th>
            <th width="10%">Old</th>
            <th width="10%">New</th>
        </tr>
    </thead>

    <tbody>
        {#each mainModelProps as prop}
            <tr class:txt-primary={hasChanges(collectionA?.[prop], collectionB?.[prop])}>
                <td class="min-width">
                    <span>{prop}</span>
                </td>
                <td
                    class:changed-old-col={!isCreateDiff &&
                        hasChanges(collectionA?.[prop], collectionB?.[prop])}
                    class:changed-none-col={isCreateDiff}
                >
                    <pre class="txt">{displayValue(collectionA?.[prop])}</pre>
                </td>
                <td
                    class:changed-new-col={!isDeleteDiff &&
                        hasChanges(collectionA?.[prop], collectionB?.[prop])}
                    class:changed-none-col={isDeleteDiff}
                >
                    <pre class="txt">{displayValue(collectionB?.[prop])}</pre>
                </td>
            </tr>
        {/each}

        {#if deleteMissing || isDeleteDiff}
            {#each removedFields as field}
                <tr>
                    <th class="min-width" colspan="3">
                        <span class="txt">schema.{field.name}</span>
                        <span class="label label-danger m-l-5">
                            Removed - <small>
                                All stored data related to <strong>{field.name}</strong> will be deleted!
                            </small>
                        </span>
                    </th>
                </tr>

                {#each Object.entries(field) as [key, value]}
                    <tr class="txt-primary">
                        <td class="min-width field-key-col">{key}</td>
                        <td class="changed-old-col">
                            <pre class="txt">{displayValue(value)}</pre>
                        </td>
                        <td class="changed-none-col" />
                    </tr>
                {/each}
            {/each}
        {/if}

        {#each sharedFields as field}
            <tr>
                <th class="min-width" colspan="3">
                    <span class="txt">schema.{field.name}</span>
                    {#if hasChanges(getFieldById(schemaA, field.id), getFieldById(schemaB, field.id))}
                        <span class="label label-warning m-l-5">Changed</span>
                    {/if}
                </th>
            </tr>

            {#each Object.entries(field) as [key, newValue]}
                <tr class:txt-primary={hasChanges(getFieldById(schemaA, field.id)?.[key], newValue)}>
                    <td class="min-width field-key-col">{key}</td>
                    <td class:changed-old-col={hasChanges(getFieldById(schemaA, field.id)?.[key], newValue)}>
                        <pre class="txt">{displayValue(getFieldById(schemaA, field.id)?.[key])}</pre>
                    </td>
                    <td class:changed-new-col={hasChanges(getFieldById(schemaA, field.id)?.[key], newValue)}>
                        <pre class="txt">{displayValue(newValue)}</pre>
                    </td>
                </tr>
            {/each}
        {/each}

        {#each addedFields as field}
            <tr>
                <th class="min-width" colspan="3">
                    <span class="txt">schema.{field.name}</span>
                    <span class="label label-success m-l-5">Added</span>
                </th>
            </tr>

            {#each Object.entries(field) as [key, value]}
                <tr class="txt-primary">
                    <td class="min-width field-key-col">{key}</td>
                    <td class="changed-none-col" />
                    <td class="changed-new-col">
                        <pre class="txt">{displayValue(value)}</pre>
                    </td>
                </tr>
            {/each}
        {/each}
    </tbody>
</table>

<style lang="scss">
    .collections-diff-table {
        color: var(--txtHintColor);
        border: 2px solid var(--primaryColor);
        tr {
            background: none;
        }
        th,
        td {
            height: auto;
            padding: 2px 15px;
            border-bottom: 1px solid rgba(#000, 0.07);
        }
        th {
            height: 35px;
            padding: 4px 15px;
            color: var(--txtPrimaryColor);
        }
        thead tr {
            background: var(--primaryColor);
            th {
                color: var(--baseColor);
                background: none;
            }
        }
        .label {
            font-weight: normal;
        }
        .changed-none-col {
            color: var(--txtDisabledColor);
            background: var(--baseAlt1Color);
        }
        .changed-old-col {
            color: var(--txtPrimaryColor);
            background: var(--dangerAltColor);
        }
        .changed-new-col {
            color: var(--txtPrimaryColor);
            background: var(--successAltColor);
        }
        .field-key-col {
            padding-left: 30px;
        }
    }
</style>
