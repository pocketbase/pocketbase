<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Toggler from "@/components/base/Toggler.svelte";
    import Draggable from "@/components/base/Draggable.svelte";

    export let id = null;
    export let items = [];
    export let disabled = false;
    export let emptyPlaceholder = "Add items";
    export let newPlaceholder = "e.g. optionA";

    let newInput;
    let newInputVal = "";

    $: formattedValue = items.join(" • ");

    function remove(item) {
        items = items || [];
        CommonHelper.removeByValue(items, item);

        if (!items.length) {
            newInput?.focus();
        }
    }

    function add(item) {
        const val = item.trim();
        if (!val.length) {
            return;
        }

        items = items || [];
        CommonHelper.pushUnique(items, val);

        // reset input
        newInputVal = "";
    }
</script>

<div class="block">
    <input
        {id}
        readonly
        type="text"
        class="formatted-value-input"
        {disabled}
        placeholder={emptyPlaceholder}
        value={formattedValue}
        title={formattedValue}
    />

    {#if !disabled}
        <Toggler
            class="dropdown dropdown-block options-dropdown dropdown-left m-t-0 p-0"
            on:hide={() => (newInputVal = "")}
        >
            <!-- svelte-ignore a11y-no-static-element-interactions -->
            <div
                class="block"
                on:dragover|stopPropagation
                on:dragleave|stopPropagation
                on:dragend|stopPropagation
                on:dragstart|stopPropagation
                on:drop|stopPropagation
            >
                {#each items as item, i (item)}
                    <Draggable bind:list={items} index={i} group={"options_" + id}>
                        <div class="dropdown-item plain">
                            <span class="txt">{item}</span>
                            <div class="flex-fill"></div>
                            <button
                                type="button"
                                class="btn btn-circle btn-transparent btn-hint btn-xs"
                                title="Remove"
                                on:click|stopPropagation={() => remove(item)}
                            >
                                <i class="ri-close-line" aria-hidden="true"></i>
                            </button>
                        </div>
                    </Draggable>
                {/each}

                <div class="new-item-form">
                    <div class="form-field form-field-sm m-0">
                        <div class="input-group">
                            <!-- svelte-ignore a11y-autofocus -->
                            <input
                                bind:this={newInput}
                                autofocus
                                type="text"
                                class="new-item-input"
                                placeholder={newPlaceholder}
                                bind:value={newInputVal}
                                on:keydown={(e) => {
                                    if (e.code === "Enter") {
                                        e.preventDefault();
                                        add(e.target.value);
                                    }
                                }}
                            />
                            <div class="form-field-addon suffix">
                                <button
                                    type="button"
                                    class="btn btn-transparent btn-xs btn-circle new-item-btn"
                                    title="Add new"
                                    class:btn-disabled={!newInputVal.length}
                                    on:click={() => add(newInputVal)}
                                >
                                    <i class="ri-add-line" aria-hidden="true"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Toggler>
    {/if}
</div>

<style lang="scss">
    .formatted-value-input {
        padding-left: 10px;
        padding-right: 10px;
        cursor: pointer;
        color: var(--txtPrimaryColor);
    }
    .dropdown-item {
        padding-top: 5px;
        padding-bottom: 5px;
    }
    .new-item-form {
        position: sticky;
        z-index: 99;
        bottom: 0;
        padding: 10px;
        background: var(--baseColor);
        border-bottom-left-radius: var(--baseRadius);
        border-bottom-right-radius: var(--baseRadius);
        &:not(:first-child) {
            margin-top: 5px;
            border-top: 1px solid var(--baseAlt1Color);
        }
    }
    .new-item-input {
        padding-right: 40px;
        padding-left: 10px;
    }
    .new-item-btn {
        right: -5px;
    }
</style>
