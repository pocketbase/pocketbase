<script>
    import { onMount } from "svelte";
    import { errors, removeError } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let collection;

    let codeEditorComponent;
    let isCodeEditorComponentLoading = false;
    let schemaErrors = [];

    $: checkSchemaErrors($errors);

    function checkSchemaErrors(errs) {
        schemaErrors = [];

        const raw = CommonHelper.getNestedVal(errs, "schema", null);

        if (CommonHelper.isEmpty(raw)) {
            return;
        }

        // generic schema error
        // ---
        if (raw?.message) {
            schemaErrors.push(raw?.message);
            return;
        }

        // schema fields error
        // ---
        const columns = CommonHelper.extractColumnsFromQuery(collection?.options?.query);
        // remove base system fields
        CommonHelper.removeByValue(columns, "id");
        CommonHelper.removeByValue(columns, "created");
        CommonHelper.removeByValue(columns, "updated");

        for (let idx in raw) {
            for (let key in raw[idx]) {
                const message = raw[idx][key].message;
                const fieldName = columns[idx] || idx;

                schemaErrors.push(CommonHelper.sentenize(fieldName + ": " + message));
            }
        }
    }

    onMount(async () => {
        isCodeEditorComponentLoading = true;

        try {
            codeEditorComponent = (await import("@/components/base/CodeEditor.svelte")).default;
        } catch (err) {
            console.warn(err);
        }

        isCodeEditorComponentLoading = false;
    });
</script>

<Field class="form-field required {schemaErrors.length ? 'error' : ''}" name="options.query" let:uniqueId>
    <label for={uniqueId}>
        <span class="txt">Select query</span>
    </label>

    {#if isCodeEditorComponentLoading}
        <textarea disabled rows="7" placeholder="Loading..." />
    {:else}
        <svelte:component
            this={codeEditorComponent}
            id={uniqueId}
            placeholder="eg. SELECT id, name from posts"
            language="sql-select"
            minHeight="150"
            on:change={() => {
                if (schemaErrors.length) {
                    removeError("schema");
                }
            }}
            bind:value={collection.options.query}
        />
    {/if}

    <div class="help-block">
        <ul>
            <li>Wildcard columns (<code>*</code>) are not supported.</li>
            <li>
                The query must have a unique <code>id</code> column.
                <br />
                If your query doesn't have a suitable one, you can use the universal
                <code>(ROW_NUMBER() OVER()) as id</code>.
            </li>
            <li>
                Expressions must be aliased with a valid formatted field name (eg.
                <code>MAX(balance) as maxBalance</code>).
            </li>
        </ul>
    </div>

    {#if schemaErrors.length}
        <div class="help-block help-block-error">
            <div class="content">
                {#each schemaErrors as err}
                    <p>{err}</p>
                {/each}
            </div>
        </div>
    {/if}
</Field>
