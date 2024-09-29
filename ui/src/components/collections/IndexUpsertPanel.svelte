<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher, onMount } from "svelte";

    const dispatch = createEventDispatcher();

    export let collection;

    let panel;
    let original = "";
    let index = "";
    let key = "";
    let codeEditorComponent;
    let isCodeEditorComponentLoading = false;

    $: presetColumns =
        collection?.fields?.filter((f) => !f.toDelete && f.name != "id")?.map((f) => f.name) || [];

    $: indexParts = CommonHelper.parseIndex(index);

    $: indexColumns = indexParts.columns?.map((c) => c.name) || [];

    export function show(showIndex, showKey) {
        key = !CommonHelper.isEmpty(showKey) ? showKey : "";
        original = showIndex || blankIndex();
        index = original;

        return panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function blankIndex() {
        const parsed = CommonHelper.parseIndex("");
        parsed.tableName = collection?.name || "";

        return CommonHelper.buildIndex(parsed);
    }

    function remove() {
        dispatch("remove", original);

        hide();
    }

    function submit() {
        if (!indexColumns.length) {
            return;
        }

        dispatch("submit", {
            old: original,
            new: index,
        });

        hide();
    }

    function toggleColumn(column) {
        const clone = CommonHelper.clone(indexParts);

        const col = clone.columns.find((c) => c.name == column);
        if (col) {
            CommonHelper.removeByValue(clone.columns, col);
        } else {
            CommonHelper.pushUnique(clone.columns, { name: column });
        }

        index = CommonHelper.buildIndex(clone);
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

<OverlayPanel bind:this={panel} popup on:hide on:show {...$$restProps}>
    <svelte:fragment slot="header">
        <h4>{original ? "Update" : "Create"} index</h4>
    </svelte:fragment>

    <Field class="form-field form-field-toggle m-b-sm" let:uniqueId>
        <input
            type="checkbox"
            id={uniqueId}
            checked={indexParts.unique}
            on:change={(e) => {
                indexParts.unique = e.target.checked;
                indexParts.tableName = indexParts.tableName || collection?.name;
                index = CommonHelper.buildIndex(indexParts);
            }}
        />
        <label for={uniqueId}>Unique</label>
    </Field>

    <Field class="form-field required m-b-sm" name={`indexes.${key || ""}`} let:uniqueId>
        {#if isCodeEditorComponentLoading}
            <textarea disabled rows="7" placeholder="Loading..." />
        {:else}
            <svelte:component
                this={codeEditorComponent}
                id={uniqueId}
                placeholder={`eg. CREATE INDEX idx_test on ${collection?.name} (created)`}
                language="sql-create-index"
                minHeight="85"
                bind:value={index}
            />
        {/if}
    </Field>

    {#if presetColumns.length > 0}
        <div class="inline-flex gap-10">
            <span class="txt txt-hint">Presets</span>
            {#each presetColumns as column}
                <button
                    type="button"
                    class="label link-primary"
                    class:label-info={indexColumns.includes(column)}
                    on:click={() => toggleColumn(column)}
                >
                    {column}
                </button>
            {/each}
        </div>
    {/if}

    <svelte:fragment slot="footer">
        {#if original != ""}
            <button
                type="button"
                class="btn btn-sm btn-circle btn-hint btn-transparent m-r-auto"
                use:tooltip={{ text: "Delete", position: "top" }}
                on:click={() => remove()}
            >
                <i class="ri-delete-bin-7-line" />
            </button>
        {/if}
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
            type="button"
            class="btn"
            class:btn-disabled={indexColumns.length <= 0}
            on:click={() => submit()}
        >
            <span class="txt">Set index</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
