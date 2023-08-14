<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let collectionA = CommonHelper.initCollection();
    export let collectionB = CommonHelper.initCollection();
    export let deleteMissing = false;
    let schemaA = [];
    let schemaB = [];
    let removedFields = [];
    let sharedFields = [];
    let addedFields = [];

    $: isDeleteDiff = !collectionB?.id && !collectionB?.name;

    $: isCreateDiff = !isDeleteDiff && !collectionA?.id;

    $: schemaA = Array.isArray(collectionA?.schema) ? collectionA?.schema.concat() : [];

    $: if (
        typeof collectionA?.schema !== "undefined" ||
        typeof collectionB?.schema !== "undefined" ||
        typeof deleteMissing !== "undefined"
    ) {
        setSchemaB();
    }

    $: removedFields = schemaA.filter((fieldA) => {
        return !schemaB.find((fieldB) => fieldA.id == fieldB.id);
    });

    $: sharedFields = schemaB.filter((fieldB) => {
        return schemaA.find((fieldA) => fieldA.id == fieldB.id);
    });

    $: addedFields = schemaB.filter((fieldB) => {
        return !schemaA.find((fieldA) => fieldA.id == fieldB.id);
    });

    $: hasAnyChange = CommonHelper.hasCollectionChanges(collectionA, collectionB, deleteMissing);

    const mainModelProps = Object.keys(CommonHelper.initCollection()).filter(
        (key) => !["schema", "created", "updated"].includes(key)
    );

    function setSchemaB() {
        schemaB = Array.isArray(collectionB?.schema) ? collectionB?.schema.concat() : [];

        if (!deleteMissing) {
            schemaB = schemaB.concat(
                schemaA.filter((fieldA) => {
                    return !schemaB.find((fieldB) => fieldA.id == fieldB.id);
                })
            );
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

    function hasChanges(valA, valB) {
        // direct match
        if (valA === valB) {
            return false;
        }

        return JSON.stringify(valA) !== JSON.stringify(valB);
    }

    function displayValue(value) {
        if (typeof value === "undefined") {
            return "";
        }

        return CommonHelper.isObject(value) ? JSON.stringify(value, null, 4) : value;
    }
</script>

<div class="section-title">
    {#if !collectionA?.id}
        <span class="label label-success">Added</span>
        <strong>{collectionB?.name}</strong>
    {:else if !collectionB?.id}
        <span class="label label-danger">Deleted</span>
        <strong>{collectionA?.name}</strong>
    {:else}
        <div class="inline-flex fleg-gap-5">
            {#if hasAnyChange}
                <span class="label label-warning">Changed</span>
            {/if}
            {#if collectionA.name !== collectionB.name}
                <strong class="txt-strikethrough txt-hint">{collectionA.name}</strong>
                <i class="ri-arrow-right-line txt-sm" />
            {/if}
            <strong class="txt">{collectionB.name}</strong>
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
                        <span class="txt">field: {field.name}</span>
                        <span class="label label-danger m-l-5">
                            Deleted - <small>
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
                    <span class="txt">field: {field.name}</span>
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
                    <span class="txt">field: {field.name}</span>
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
