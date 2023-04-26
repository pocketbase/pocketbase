<script context="module">
    let siblings = [];
</script>

<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { slide } from "svelte/transition";
    import { SchemaField } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { errors, setErrors } from "@/stores/errors";
    import Toggler from "@/components/base/Toggler.svelte";
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let field = new SchemaField();

    let nameInput;
    let isDragOver = false;
    let showOptions = false;

    const componentId = "f_" + CommonHelper.randomString(8);

    const dispatch = createEventDispatcher();

    const customRequiredLabels = {
        // type => label
        bool: "Nonfalsey",
        number: "Nonzero",
    };

    $: if (field.toDelete) {
        // reset the name if it was previously deleted
        if (field.originalName && field.name !== field.originalName) {
            field.name = field.originalName;
        }
    }

    $: if (!field.originalName && field.name) {
        field.originalName = field.name;
    }

    $: if (typeof field.toDelete === "undefined") {
        field.toDelete = false; // normalize
    }

    $: if (field.required) {
        field.nullable = false;
    }

    $: interactive = !field.toDelete && !(field.id && field.system);

    $: hasErrors = !CommonHelper.isEmpty(CommonHelper.getNestedVal($errors, `schema.${key}`));

    $: requiredLabel = customRequiredLabels[field?.type] || "Nonempty";

    function remove() {
        if (!field.id) {
            dispatch("remove");
        } else {
            field.toDelete = true;
        }
    }

    function restore() {
        field.toDelete = false;

        // reset all errors since the error index key would have been changed
        setErrors({});
    }

    function normalizeFieldName(name) {
        return CommonHelper.slugify(name);
    }

    function expand() {
        showOptions = true;
        collapseSiblings();
    }

    function collapse() {
        showOptions = false;
    }

    function toggle() {
        if (showOptions) {
            collapse();
        } else {
            expand();
        }
    }

    function collapseSiblings() {
        for (let f of siblings) {
            if (f.id == componentId) {
                continue;
            }
            f.collapse();
        }
    }

    onMount(() => {
        siblings.push({
            id: componentId,
            collapse: collapse,
        });

        if (field.onMountSelect) {
            field.onMountSelect = false;
            nameInput?.select();
        }

        return () => {
            CommonHelper.removeByKey(siblings, "id", componentId);
        };
    });
</script>

<div
    draggable={true}
    class="schema-field"
    class:required={field.required}
    class:expanded={interactive && showOptions}
    class:deleted={field.toDelete}
    class:drag-over={isDragOver}
    transition:slide|local={{ duration: 150 }}
    on:dragstart={(e) => {
        if (!e.target.classList.contains("drag-handle-wrapper")) {
            e.preventDefault();
            return;
        }

        const blank = document.createElement("div");
        e.dataTransfer.setDragImage(blank, 0, 0);
        interactive && dispatch("dragstart", e);
    }}
    on:dragenter={(e) => {
        if (interactive) {
            isDragOver = true;
            dispatch("dragenter", e);
        }
    }}
    on:drop|preventDefault={(e) => {
        if (interactive) {
            isDragOver = false;
            dispatch("drop", e);
        }
    }}
    on:dragleave={(e) => {
        if (interactive) {
            isDragOver = false;
            dispatch("dragleave", e);
        }
    }}
    on:dragover|preventDefault
>
    <div class="schema-field-header">
        {#if interactive}
            <div class="drag-handle-wrapper" draggable="true" aria-label="Sort">
                <span class="drag-handle" />
            </div>
        {/if}
        <Field
            class="form-field required m-0 {!interactive ? 'disabled' : ''}"
            name="schema.{key}.name"
            inlineError
        >
            <div class="markers">
                {#if field.required}
                    <span class="marker marker-required" use:tooltip={requiredLabel} />
                {/if}
            </div>

            <div
                class="form-field-addon prefix no-pointer-events field-type-icon"
                class:txt-disabled={!interactive}
            >
                <i class={CommonHelper.getFieldTypeIcon(field.type)} />
            </div>

            <!-- svelte-ignore a11y-autofocus -->
            <input
                bind:this={nameInput}
                type="text"
                required
                disabled={!interactive}
                readonly={field.id && field.system}
                spellcheck="false"
                autofocus={!field.id}
                placeholder="Field name"
                value={field.name}
                on:input={(e) => {
                    const oldName = field.name;
                    field.name = normalizeFieldName(e.target.value);
                    e.target.value = field.name;

                    dispatch("rename", { oldName: oldName, newName: field.name });
                }}
            />
        </Field>

        <slot {interactive} {hasErrors}>
            <span class="separator" />
        </slot>

        {#if field.toDelete}
            <button
                type="button"
                class="btn btn-sm btn-circle btn-warning btn-transparent options-trigger"
                aria-label="Restore"
                use:tooltip={"Restore"}
                on:click={restore}
            >
                <i class="ri-restart-line" />
            </button>
        {:else if interactive}
            <button
                type="button"
                aria-label="Toggle field options"
                class="btn btn-sm btn-circle options-trigger {showOptions
                    ? 'btn-secondary'
                    : 'btn-transparent'}"
                class:btn-hint={!showOptions && !hasErrors}
                class:btn-danger={hasErrors}
                on:click={toggle}
            >
                <i class="ri-settings-3-line" />
            </button>
        {/if}
    </div>

    {#if interactive && showOptions}
        <div class="schema-field-options" transition:slide|local={{ duration: 150 }}>
            <div class="grid grid-sm">
                <div class="col-sm-12 hidden-empty">
                    <slot name="options" {interactive} {hasErrors} />
                </div>

                <slot name="beforeNonempty" {interactive} {hasErrors} />

                <div class="col-sm-4">
                    <Field class="form-field form-field-toggle m-0" name="requried" let:uniqueId>
                        <input type="checkbox" id={uniqueId} bind:checked={field.required} />
                        <label for={uniqueId}>
                            <span class="txt">{requiredLabel}</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: `Requires the field value NOT to be ${CommonHelper.zeroDefaultStr(
                                        field
                                    )}.`,
                                    position: "right",
                                }}
                            />
                        </label>
                    </Field>
                </div>

                <slot name="afterNonempty" {interactive} {hasErrors} />

                {#if !field.toDelete}
                    <div class="col-sm-4 m-l-auto txt-right">
                        <div class="flex-fill" />
                        <div class="inline-flex flex-gap-sm flex-nowrap">
                            <button
                                type="button"
                                aria-label="More"
                                class="btn btn-circle btn-sm btn-transparent"
                            >
                                <i class="ri-more-line" />
                                <Toggler
                                    class="dropdown dropdown-sm dropdown-upside dropdown-right dropdown-nowrap no-min-width"
                                >
                                    <button type="button" class="dropdown-item txt-right" on:click={remove}>
                                        <span class="txt">Remove</span>
                                    </button>
                                </Toggler>
                            </button>
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>
