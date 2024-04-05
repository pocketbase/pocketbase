<script context="module">
    let siblings = [];
</script>

<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import { errors, setErrors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher, onMount } from "svelte";
    import { slide } from "svelte/transition";

    const componentId = "f_" + CommonHelper.randomString(8);

    const dispatch = createEventDispatcher();

    const customRequiredLabels = {
        // type => label
        bool: "Nonfalsey",
        number: "Nonzero",
    };

    export let key = "";
    export let field = CommonHelper.initSchemaField();

    let nameInput;
    let showOptions = false;

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
            collapse();
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

    function duplicate() {
        if (!field.toDelete) {
            collapse();
            dispatch("duplicate");
        }
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
    class="schema-field"
    class:required={field.required}
    class:expanded={interactive && showOptions}
    class:deleted={field.toDelete}
    transition:slide={{ duration: 150 }}
>
    <div class="schema-field-header">
        {#if interactive}
            <div class="drag-handle-wrapper" draggable={true} aria-label="Sort">
                <span class="drag-handle" />
            </div>
        {/if}
        <Field
            class="form-field required m-0 {!interactive ? 'disabled' : ''}"
            name="schema.{key}.name"
            inlineError
        >
            {#if field.required}
                <div class="field-labels">
                    <span class="label label-success">{requiredLabel}</span>
                </div>
            {/if}

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
                aria-expanded={showOptions}
            >
                <i class="ri-settings-3-line" />
            </button>
        {/if}
    </div>

    {#if interactive && showOptions}
        <div class="schema-field-options" transition:slide={{ duration: 150 }}>
            <div class="hidden-empty m-b-sm">
                <slot name="options" {interactive} {hasErrors} />
            </div>

            <div class="schema-field-options-footer">
                <Field class="form-field form-field-toggle" name="requried" let:uniqueId>
                    <input type="checkbox" id={uniqueId} bind:checked={field.required} />
                    <label for={uniqueId}>
                        <span class="txt">{requiredLabel}</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: `Requires the field value NOT to be ${CommonHelper.zeroDefaultStr(
                                    field,
                                )}.`,
                            }}
                        />
                    </label>
                </Field>

                <Field class="form-field form-field-toggle" name="presentable" let:uniqueId>
                    <input type="checkbox" id={uniqueId} bind:checked={field.presentable} />
                    <label for={uniqueId}>
                        <span class="txt">Presentable</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: `Whether the field should be preferred in the Admin UI relation listings (default to auto).`,
                            }}
                        />
                    </label>
                </Field>

                <slot name="optionsFooter" {interactive} {hasErrors} />

                {#if !field.toDelete}
                    <div class="m-l-auto txt-right">
                        <div class="inline-flex flex-gap-sm flex-nowrap">
                            <div
                                tabindex="0"
                                role="button"
                                aria-label="More"
                                class="btn btn-circle btn-sm btn-transparent"
                            >
                                <i class="ri-more-line" aria-hidden="true" />
                                <Toggler
                                    class="dropdown dropdown-sm dropdown-upside dropdown-right dropdown-nowrap no-min-width"
                                >
                                    <button
                                        type="button"
                                        class="dropdown-item"
                                        role="menuitem"
                                        on:click|preventDefault={duplicate}
                                    >
                                        <span class="txt">Duplicate</span>
                                    </button>
                                    <button
                                        type="button"
                                        class="dropdown-item"
                                        role="menuitem"
                                        on:click|preventDefault={remove}
                                    >
                                        <span class="txt">Remove</span>
                                    </button>
                                </Toggler>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>
