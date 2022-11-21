<script>
    import { SchemaField } from "pocketbase";
    import FieldAccordion from "@/components/collections/FieldAccordion.svelte";

    export let collection = {};

    const baseReservedNames = [
        "id",
        "created",
        "updated",
        "collectionId",
        "collectionName",
        "expand",
        "true",
        "false",
        "null",
    ];

    let reservedNames = [];

    $: if (collection.isAuth) {
        reservedNames = baseReservedNames.concat([
            "username",
            "email",
            "emailVisibility",
            "verified",
            "tokenKey",
            "passwordHash",
            "lastResetSentAt",
            "lastVerificationSentAt",
            "password",
            "passwordConfirm",
            "oldPassword",
        ]);
    } else {
        reservedNames = baseReservedNames.slice(0);
    }

    $: if (typeof collection?.schema === "undefined") {
        collection = collection || {};
        collection.schema = [];
    }

    function removeField(fieldIndex) {
        if (collection.schema[fieldIndex]) {
            collection.schema.splice(fieldIndex, 1);
            collection.schema = collection.schema;
        }
    }

    function newField() {
        const field = new SchemaField({
            name: getUniqueFieldName(),
        });

        collection.schema.push(field);
        collection.schema = collection.schema;
    }

    function getUniqueFieldName(base = "field") {
        let counter = "";

        while (hasFieldWithName(base + counter)) {
            ++counter;
        }

        return base + counter;
    }

    function hasFieldWithName(name) {
        return !!collection.schema.find((field) => field.name === name);
    }

    function getSiblingsFieldNames(currentField) {
        let result = [];

        if (currentField.toDelete) {
            return result;
        }

        for (let field of collection.schema) {
            if (field === currentField || field.toDelete) {
                continue; // skip current and deleted fields
            }

            result.push(field.name);
        }

        return result;
    }

    // ---------------------------------------------------------------
    // fields drag&drop handling
    // ---------------------------------------------------------------

    function onFieldDrag(event, i) {
        if (!event) {
            return;
        }

        event.dataTransfer.effectAllowed = "move";
        event.dataTransfer.dropEffect = "move";
        event.dataTransfer.setData("text/plain", i);
    }

    function onFieldDrop(event, target) {
        if (!event) {
            return;
        }

        event.dataTransfer.dropEffect = "move";

        const start = parseInt(event.dataTransfer.getData("text/plain"));
        const newSchema = collection.schema;

        if (start < target) {
            newSchema.splice(target + 1, 0, newSchema[start]);
            newSchema.splice(start, 1);
        } else {
            newSchema.splice(target, 0, newSchema[start]);
            newSchema.splice(start + 1, 1);
        }

        collection.schema = newSchema;
    }
</script>

<div class="block m-b-25">
    <p class="txt-sm">
        System fields:
        <code class="txt-sm">id</code> ,
        <code class="txt-sm">created</code> ,
        <code class="txt-sm">updated</code>
        {#if collection.isAuth}
            ,
            <code class="txt-sm">username</code> ,
            <code class="txt-sm">email</code> ,
            <code class="txt-sm">emailVisibility</code> ,
            <code class="txt-sm">verified</code>
        {/if}
        .
    </p>
</div>

<div class="accordions">
    {#each collection.schema as field, i (field)}
        <FieldAccordion
            bind:field
            key={i}
            excludeNames={reservedNames.concat(getSiblingsFieldNames(field))}
            on:remove={() => removeField(i)}
            on:dragstart={(e) => onFieldDrag(e?.detail, i)}
            on:drop={(e) => onFieldDrop(e?.detail, i)}
        />
    {/each}
</div>

<div class="clearfix m-t-xs" />

<button
    type="button"
    class="btn btn-block {collection?.isAuth || collection.schema?.length ? 'btn-secondary' : 'btn-warning'}"
    on:click={newField}
>
    <i class="ri-add-line" />
    <span class="txt">New field</span>
</button>
