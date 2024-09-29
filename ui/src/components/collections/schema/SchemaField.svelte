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

    // @todo refactor once the UI is dynamic
    const authHideNonemptyToggle = ["password", "tokenKey", "id", "autodate"];
    const authHideHiddenToggle = ["password", "tokenKey", "id", "email"];
    const authHidePresentableToggle = ["password", "tokenKey"];

    export let key = "";
    export let field = CommonHelper.initSchemaField();
    export let draggable = true;
    export let collection = {};

    let nameInput;
    let showOptions = false;

    $: isAuthCollection = collection?.type == "auth";

    $: if (field._toDelete) {
        // reset the name if it was previously deleted
        if (field._originalName && field.name !== field._originalName) {
            field.name = field._originalName;
        }
    }

    $: if (!field._originalName && field.name) {
        field._originalName = field.name;
    }

    $: if (typeof field._toDelete === "undefined") {
        field._toDelete = false; // normalize
    }

    $: if (field.required) {
        field.nullable = false;
    }

    $: interactive = !field._toDelete;

    $: hasErrors = !CommonHelper.isEmpty(CommonHelper.getNestedVal($errors, `fields.${key}`));

    $: requiredLabel = customRequiredLabels[field?.type] || "Nonempty";

    function remove() {
        if (!field.id) {
            collapse();
            dispatch("remove");
        } else {
            field._toDelete = true;
        }
    }

    function restore() {
        field._toDelete = false;

        // reset all errors since the error index key would have been changed
        setErrors({});
    }

    function duplicate() {
        if (!field._toDelete) {
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
    class:deleted={field._toDelete}
    transition:slide={{ duration: 150 }}
>
    <div class="schema-field-header">
        {#if interactive && draggable}
            <div class="drag-handle-wrapper" draggable={true} aria-label="Sort">
                <span class="drag-handle" />
            </div>
        {/if}
        <Field
            class="form-field required m-0 {!interactive ? 'disabled' : ''}"
            name="fields.{key}.name"
            inlineError
        >
            <div class="field-labels">
                {#if field.required}
                    <span class="label label-success">{requiredLabel}</span>
                {/if}
                {#if field.hidden}
                    <span class="label label-danger">Hidden</span>
                {/if}
            </div>

            <!-- svelte-ignore a11y-no-static-element-interactions -->
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div
                class="form-field-addon prefix field-type-icon"
                class:txt-disabled={!interactive || field.system}
                use:tooltip={field.type + (field.system ? " (system)" : "")}
                on:click={() => nameInput?.focus()}
            >
                <i class={CommonHelper.getFieldTypeIcon(field.type)} />
            </div>

            <input
                bind:this={nameInput}
                type="text"
                required
                disabled={!interactive || field.system}
                spellcheck="false"
                placeholder="Field name"
                value={field.name}
                title="System field"
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

        {#if field._toDelete}
            <button
                type="button"
                class="btn btn-sm btn-circle btn-success btn-transparent options-trigger"
                aria-label="Restore"
                use:tooltip={"Restore"}
                on:click={restore}
            >
                <i class="ri-restart-line" />
            </button>
        {:else if interactive}
            <button
                type="button"
                aria-label="Toggle {field.name} field options"
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
        <div class="schema-field-options" transition:slide={{ delay: 10, duration: 150 }}>
            <div class="hidden-empty m-b-sm">
                <slot name="options" {interactive} {hasErrors} />
            </div>

            <div class="schema-field-options-footer">
                <!-- @todo move to each field after the refactoring -->
                {#if !field.primaryKey && field.type != "autodate" && (!isAuthCollection || !authHideNonemptyToggle.includes(field.name))}
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
                {/if}

                {#if !field.primaryKey && (!isAuthCollection || !authHideHiddenToggle.includes(field.name))}
                    <Field class="form-field form-field-toggle" name="hidden" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            bind:checked={field.hidden}
                            on:change={(e) => {
                                if (e.target.checked) {
                                    field.presentable = false;
                                }
                            }}
                        />
                        <label for={uniqueId}>
                            <span class="txt">Hidden</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: `Hide from the JSON API response and filters.`,
                                }}
                            />
                        </label>
                    </Field>
                {/if}

                {#if !isAuthCollection || !authHidePresentableToggle.includes(field.name)}
                    <Field class="form-field form-field-toggle m-0" name="presentable" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            bind:checked={field.presentable}
                            disabled={field.hidden}
                        />
                        <label for={uniqueId}>
                            <span class="txt">Presentable</span>
                            <i
                                class="ri-information-line {field.hidden ? 'txt-disabled' : 'link-hint'}"
                                use:tooltip={{
                                    text: `Whether the field should be preferred in the Superuser UI relation listings (default to auto).`,
                                }}
                            />
                        </label>
                    </Field>
                {/if}

                <slot name="optionsFooter" {interactive} {hasErrors} />

                {#if !field._toDelete && !field.primaryKey}
                    <div class="m-l-auto txt-right">
                        <div class="inline-flex flex-gap-sm flex-nowrap">
                            <div
                                tabindex="0"
                                role="button"
                                title="More field options"
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
                                    {#if !field.system}
                                        <button
                                            type="button"
                                            class="dropdown-item"
                                            role="menuitem"
                                            on:click|preventDefault={remove}
                                        >
                                            <span class="txt">Remove</span>
                                        </button>
                                    {/if}
                                </Toggler>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>
