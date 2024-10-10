<script>
    import Field from "@/components/base/Field.svelte";
    import { errors, removeError } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { onMount } from "svelte";

    export let collection;

    let codeEditorComponent;
    let isCodeEditorComponentLoading = false;
    let fieldsErrors = [];

    $: checkFieldsErrors($errors);

    function checkFieldsErrors(errs) {
        fieldsErrors = [];

        const raw = CommonHelper.getNestedVal(errs, "fields", null);

        if (CommonHelper.isEmpty(raw)) {
            return;
        }

        // generic fields list error
        // ---
        if (raw?.message) {
            fieldsErrors.push(raw?.message);
            return;
        }

        // individual field error
        // ---
        const columns = CommonHelper.extractColumnsFromQuery(collection?.viewQuery);
        // remove base system fields
        CommonHelper.removeByValue(columns, "id");
        CommonHelper.removeByValue(columns, "created");
        CommonHelper.removeByValue(columns, "updated");

        for (let idx in raw) {
            for (let key in raw[idx]) {
                const message = raw[idx][key].message;
                const fieldName = columns[idx] || idx;

                fieldsErrors.push(CommonHelper.sentenize(fieldName + ": " + message));
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

<Field class="form-field required {fieldsErrors.length ? 'error' : ''}" name="viewQuery" let:uniqueId>
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
                if (fieldsErrors.length) {
                    removeError("fields");
                }
            }}
            bind:value={collection.viewQuery}
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
                Expressions must be aliased with a valid formatted field name, e.g.
                <code>MAX(balance) as maxBalance</code>.
            </li>
            <li>
                Combined/multi-spaced expressions must be wrapped in parenthesis, e.g.
                <code>(MAX(balance) + 1) as maxBalance</code>.
            </li>
        </ul>
    </div>

    {#if fieldsErrors.length}
        <div class="help-block help-block-error">
            <div class="content">
                {#each fieldsErrors as err}
                    <p>{err}</p>
                {/each}
            </div>
        </div>
    {/if}
</Field>
