<script>
    import Draggable from "@/components/base/Draggable.svelte";
    import IndexesList from "@/components/collections/IndexesList.svelte";
    import NewField from "@/components/collections/schema/NewField.svelte";
    import SchemaFieldAutodate from "@/components/collections/schema/SchemaFieldAutodate.svelte";
    import SchemaFieldBool from "@/components/collections/schema/SchemaFieldBool.svelte";
    import SchemaFieldDate from "@/components/collections/schema/SchemaFieldDate.svelte";
    import SchemaFieldEditor from "@/components/collections/schema/SchemaFieldEditor.svelte";
    import SchemaFieldEmail from "@/components/collections/schema/SchemaFieldEmail.svelte";
    import SchemaFieldFile from "@/components/collections/schema/SchemaFieldFile.svelte";
    import SchemaFieldJson from "@/components/collections/schema/SchemaFieldJson.svelte";
    import SchemaFieldNumber from "@/components/collections/schema/SchemaFieldNumber.svelte";
    import SchemaFieldPassword from "@/components/collections/schema/SchemaFieldPassword.svelte";
    import SchemaFieldRelation from "@/components/collections/schema/SchemaFieldRelation.svelte";
    import SchemaFieldSelect from "@/components/collections/schema/SchemaFieldSelect.svelte";
    import SchemaFieldText from "@/components/collections/schema/SchemaFieldText.svelte";
    import SchemaFieldUrl from "@/components/collections/schema/SchemaFieldUrl.svelte";
    import SchemaFieldGeoPoint from "@/components/collections/schema/SchemaFieldGeoPoint.svelte";
    import { scaffolds } from "@/stores/collections";
    import { setErrors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let oldCollectionType;

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
        password: SchemaFieldPassword,
        autodate: SchemaFieldAutodate,
        geoPoint: SchemaFieldGeoPoint,
    };

    $: if (!collection.id && oldCollectionType != collection.type) {
        onTypeCanged();
        oldCollectionType = collection.type;
    }

    $: if (typeof collection.fields === "undefined") {
        collection.fields = [];
    }

    $: nonDeletedFields = collection.fields.filter((f) => !f._toDelete);

    function removeField(fieldIndex) {
        if (collection.fields[fieldIndex]) {
            collection.fields.splice(fieldIndex, 1);
            collection.fields = collection.fields;
        }
    }

    function duplicateField(fieldIndex) {
        const field = collection.fields[fieldIndex];
        if (!field) {
            return; // nothing to duplicate
        }

        field.onMountSelect = false;

        const clone = structuredClone(field);
        clone.id = "";
        clone.system = false;
        clone.name = getUniqueFieldName(clone.name + "_copy");
        clone.onMountSelect = true;

        collection.fields.splice(fieldIndex + 1, 0, clone);
        collection.fields = collection.fields;
    }

    function newField(fieldType = "text") {
        const field = CommonHelper.initSchemaField({
            name: getUniqueFieldName(),
            type: fieldType,
        });

        field.onMountSelect = true;

        // if the collection has created/updated last fields,
        // insert before the first autodate field, otherwise - append
        const idx = collection.fields.findLastIndex((f) => f.type != "autodate");
        if (field.type != "autodate" && idx >= 0) {
            collection.fields.splice(idx + 1, 0, field);
        } else {
            collection.fields.push(field);
        }

        collection.fields = collection.fields;
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
        return !!collection?.fields?.find((field) => field.name.toLowerCase() === name.toLowerCase());
    }

    function getSchemaFieldIndex(field) {
        return nonDeletedFields.findIndex((f) => f === field);
    }

    function replaceIndexesColumn(oldName, newName) {
        if (!collection?.fields?.length || oldName === newName || !newName) {
            return;
        }

        // field with the old name exists so there is no need to rename index columns
        if (!!collection?.fields?.find((f) => f.name == oldName && !f._toDelete)) {
            return;
        }

        // update indexes on renamed fields
        collection.indexes = collection.indexes.map((idx) =>
            CommonHelper.replaceIndexColumn(idx, oldName, newName),
        );
    }

    function onTypeCanged() {
        const newScaffold = structuredClone($scaffolds[collection.type]);

        // merge fields
        // -----------------------------------------------------------
        const oldFields = collection.fields || [];
        const nonSystemFields = oldFields.filter((f) => !f.system);

        collection.fields = newScaffold.fields;

        for (const oldField of oldFields) {
            if (!oldField.system) {
                continue;
            }

            const idx = collection.fields.findIndex((f) => f.name == oldField.name);
            if (idx < 0) {
                continue;
            }

            // merge the default field with the existing one
            collection.fields[idx] = Object.assign(collection.fields[idx], oldField);
        }

        for (const field of nonSystemFields) {
            collection.fields.push(field);
        }

        // merge indexes
        // -----------------------------------------------------------
        collection.indexes = collection.indexes || [];

        if (collection.indexes.length) {
            const oldScaffoldIndexes = $scaffolds[oldCollectionType]?.indexes || [];

            indexesLoop: for (let i = collection.indexes.length - 1; i >= 0; i--) {
                const parsed = CommonHelper.parseIndex(collection.indexes[i]);
                const parsedName = parsed.indexName.toLowerCase();

                // remove old scaffold indexes
                for (const idx of oldScaffoldIndexes) {
                    const oldScaffoldName = CommonHelper.parseIndex(idx).indexName.toLowerCase();
                    if (parsedName == oldScaffoldName) {
                        collection.indexes.splice(i, 1);
                        continue indexesLoop;
                    }
                }

                // remove indexes to nonexisting fields
                for (const column of parsed.columns) {
                    if (!hasFieldWithName(column.name)) {
                        collection.indexes.splice(i, 1);
                        continue indexesLoop;
                    }
                }
            }
        }

        // merge new scaffold indexes
        CommonHelper.mergeUnique(collection.indexes, newScaffold.indexes);
    }

    function replaceIdentityFields(oldName, newName) {
        if (oldName === newName || !newName) {
            return;
        }

        let identityFields = collection.passwordAuth?.identityFields || [];

        for (let i = 0; i < identityFields.length; i++) {
            if (identityFields[i] == oldName) {
                identityFields[i] = newName;
            }
        }
    }

    function onFieldRename(oldName, newName) {
        replaceIndexesColumn(oldName, newName);
        replaceIdentityFields(oldName, newName);
    }
</script>

<div class="schema-fields total-{collection.fields.length}">
    {#each collection.fields as field, i (field)}
        <Draggable
            bind:list={collection.fields}
            index={i}
            disabled={field._toDelete}
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
                {collection}
                bind:field
                on:remove={() => removeField(i)}
                on:duplicate={() => duplicateField(i)}
                on:rename={(e) => onFieldRename(e.detail.oldName, e.detail.newName)}
            />
        </Draggable>
    {/each}
</div>

<NewField class="btn btn-block btn-outline" on:select={(e) => newField(e.detail)} />

<hr />

<IndexesList bind:collection />
