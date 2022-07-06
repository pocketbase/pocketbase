<script>
    import { SchemaField } from "pocketbase";
    import FieldAccordion from "@/components/collections/FieldAccordion.svelte";

    const reservedNames = ["id", "created", "updated"];

    export let collection = {};

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

        for (let field of collection.schema) {
            if (field === currentField) {
                continue; // skip current
            }

            result.push(field.name);

            if (field.id && field.originalName !== "" && field.originalName !== field.name) {
                result.push(field.originalName);
            }
        }

        return result;
    }
</script>

<div class="accordions">
    {#each collection.schema as field, i (i)}
        <FieldAccordion
            bind:field
            key={i}
            excludeNames={reservedNames.concat(getSiblingsFieldNames(field))}
            on:remove={() => removeField(i)}
        />
    {/each}
</div>

<div class="clearfix m-t-xs" />

<button
    type="button"
    class="btn btn-block {collection.schema?.length ? 'btn-secondary' : 'btn-success'}"
    on:click={newField}
>
    <i class="ri-add-line" />
    <span class="txt">New field</span>
</button>
