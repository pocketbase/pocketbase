<script>
    import { Collection, SchemaField } from "pocketbase";
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

    export let collection = new Collection();

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
        collection = collection || new Collection();
        collection.schema = [];
    }

    $: nonDeletedFields = collection.schema.filter((f) => !f.toDelete) || [];

    function removeField(fieldIndex) {
        if (collection.schema[fieldIndex]) {
            collection.schema.splice(fieldIndex, 1);
            collection.schema = collection.schema;
        }
    }

    function newField(fieldType = "text") {
        const field = new SchemaField({
            name: getUniqueFieldName(),
            type: fieldType,
        });

        field.onMountSelect = true;

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

    function getSchemaFieldIndex(field) {
        return nonDeletedFields.findIndex((f) => f === field);
    }

    function replaceIndexesColumn(oldName, newName) {
        if (!collection?.schema?.length || oldName === newName || !newName) {
            return;
        }

        // update indexes on renamed fields
        collection.indexes = collection.indexes.map((idx) =>
            CommonHelper.replaceIndexColumn(idx, oldName, newName)
        );
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

        // reset errors since the schema keys index has changed
        setErrors({});
    }
</script>

<div class="block m-b-25">
    <p class="txt-sm">
        System fields:
        <code class="txt-sm">id</code> ,
        <code class="txt-sm">created</code> ,
        <code class="txt-sm">updated</code>
        {#if collection.$isAuth}
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
        <svelte:component
            this={fieldComponents[field.type]}
            key={getSchemaFieldIndex(field)}
            bind:field
            on:remove={() => removeField(i)}
            on:rename={(e) => replaceIndexesColumn(e.detail.oldName, e.detail.newName)}
            on:dragstart={(e) => onFieldDrag(e.detail, i)}
            on:drop={(e) => onFieldDrop(e.detail, i)}
        />
    {/each}
</div>

<NewField class="btn btn-block btn-outline" on:select={(e) => newField(e.detail)} />

<hr />

<IndexesList bind:collection />
