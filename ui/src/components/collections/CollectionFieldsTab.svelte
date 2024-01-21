<script>
    import { setErrors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import IndexesList from "@/components/collections/IndexesList.svelte";
    import NewField from "@/components/collections/schema/NewField.svelte";
    import SchemaFieldText from "@/components/collections/schema/SchemaFieldText.svelte";
    import SchemaFieldNumber from "@/components/collections/schema/SchemaFieldNumber.svelte";
    import SchemaFieldBool from "@/components/collections/schema/SchemaFieldBool.svelte";
    import SchemaFieldEmail from "@/components/collections/schema/SchemaFieldEmail.svelte";
    import SchemaFieldUrl from "@/components/collections/schema/SchemaFieldUrl.svelte";
    import SchemaFieldEditor from "@/components/collections/schema/SchemaFieldEditor.svelte";
    import SchemaFieldDate from "@/components/collections/schema/SchemaFieldDate.svelte";
    import SchemaFieldSelect from "@/components/collections/schema/SchemaFieldSelect.svelte";
    import SchemaFieldJson from "@/components/collections/schema/SchemaFieldJson.svelte";
    import SchemaFieldFile from "@/components/collections/schema/SchemaFieldFile.svelte";
    import SchemaFieldRelation from "@/components/collections/schema/SchemaFieldRelation.svelte";
    import Draggable from "@/components/base/Draggable.svelte";

    export let collection;

    const fieldComponents = {
        text: SchemaFieldText,
        number: SchemaFieldNumber,
        bool: SchemaFieldBool,
        email: SchemaFieldEmail,
        url: SchemaFieldUrl,
        editor: SchemaFieldEditor,
        date: SchemaFieldDate,
        select: SchemaFieldSelect,
        json: SchemaFieldJson,
        file: SchemaFieldFile,
        relation: SchemaFieldRelation,
    };

    $: if (typeof collection.schema === "undefined") {
        collection.schema = [];
    }

    $: nonDeletedFields = collection.schema.filter((f) => !f.toDelete) || [];

    function removeField(fieldIndex) {
        if (collection.schema[fieldIndex]) {
            collection.schema.splice(fieldIndex, 1);
            collection.schema = collection.schema;
        }
    }

    function duplicateField(fieldIndex) {
        const field = collection.schema[fieldIndex];
        if (!field) {
            return; // nothing to duplicate
        }

        field.onMountSelect = false;

        const clone = structuredClone(field);
        clone.id = "";
        clone.name = getUniqueFieldName(clone.name + "_copy");
        clone.onMountSelect = true;

        collection.schema.splice(fieldIndex + 1, 0, clone);
        collection.schema = collection.schema;
    }

    function newField(fieldType = "text") {
        const field = CommonHelper.initSchemaField({
            name: getUniqueFieldName(),
            type: fieldType,
        });

        field.onMountSelect = true;

        collection.schema.push(field);
        collection.schema = collection.schema;
    }

    function getUniqueFieldName(name = "field") {
        let result = name;
        let counter = 2;

        let suffix = name.match(/\d+$/)?.[0] || ""; // extract numeric suffix

        // name without the suffix
        let base = suffix ? name.substring(0, name.length - suffix.length) : name;

        while (hasFieldWithName(result)) {
            result = base + ((suffix << 0) + counter);
            counter++;
        }

        return result;
    }

    function hasFieldWithName(name) {
        return !!collection?.schema?.find((field) => field.name === name);
    }

    function getSchemaFieldIndex(field) {
        return nonDeletedFields.findIndex((f) => f === field);
    }

    function replaceIndexesColumn(oldName, newName) {
        if (!collection?.schema?.length || oldName === newName || !newName) {
            return;
        }

        // field with the old name exists so there is no need to rename index columns
        if (!!collection?.schema?.find((f) => f.name == oldName && !f.toDelete)) {
            return;
        }

        // update indexes on renamed fields
        collection.indexes = collection.indexes.map((idx) =>
            CommonHelper.replaceIndexColumn(idx, oldName, newName),
        );
    }
</script>

<div class="block m-b-25">
    <p class="txt-sm">
        System fields:
        <code class="txt-sm">id</code> ,
        <code class="txt-sm">created</code> ,
        <code class="txt-sm">updated</code>
        {#if collection.type === "auth"}
            ,
            <code class="txt-sm">username</code> ,
            <code class="txt-sm">email</code> ,
            <code class="txt-sm">emailVisibility</code> ,
            <code class="txt-sm">verified</code>
        {/if}
        .
    </p>
</div>

<div class="schema-fields">
    {#each collection.schema as field, i (field)}
        <Draggable
            bind:list={collection.schema}
            index={i}
            disabled={field.toDelete || (field.id && field.system)}
            dragHandleClass="drag-handle-wrapper"
            on:drag={(e) => {
                // blank drag placeholder
                if (!e.detail) {
                    return;
                }
                const ghost = e.detail.target;
                ghost.style.opacity = 0;
                setTimeout(() => {
                    ghost?.style?.removeProperty("opacity"); // restore
                }, 0);
                e.detail.dataTransfer.setDragImage(ghost, 0, 0);
            }}
            on:sort={() => {
                // reset errors since the schema keys index has changed
                setErrors({});
            }}
        >
            <svelte:component
                this={fieldComponents[field.type]}
                key={getSchemaFieldIndex(field)}
                bind:field
                on:remove={() => removeField(i)}
                on:duplicate={() => duplicateField(i)}
                on:rename={(e) => replaceIndexesColumn(e.detail.oldName, e.detail.newName)}
            />
        </Draggable>
    {/each}
</div>

<NewField class="btn btn-block btn-outline" on:select={(e) => newField(e.detail)} />

<hr />

<IndexesList bind:collection />
