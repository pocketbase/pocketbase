<script>
    import { onMount, createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";
    import Toggler from "@/components/base/Toggler.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let id = "";
    export let noOptionsText = "No options found";
    export let selectPlaceholder = "- Select -";
    export let searchPlaceholder = "Search...";
    export let items = [];
    export let multiple = false;
    export let disabled = false;
    export let readonly = false;
    export let upside = false;
    export let zeroFunc = () => (multiple ? [] : undefined);
    export let selected = zeroFunc();
    export let toggle = multiple; // toggle option on click
    export let closable = true; // close the dropdown on option select/deselect
    export let labelComponent = undefined; // custom component to use for each selected option label
    export let labelComponentProps = {}; // props to pass to the custom option component
    export let optionComponent = undefined; // custom component to use for each dropdown option item
    export let optionComponentProps = {}; // props to pass to the custom option component
    export let searchable = false; // whether to show the dropdown options search input
    export let searchFunc = undefined; // custom search option filter: `function(item, searchTerm):boolean`

    const dispatch = createEventDispatcher();

    let classes = "";
    export { classes as class }; // export reserved keyword

    let toggler;
    let searchTerm = "";
    let container = undefined;
    let labelDiv = undefined;

    $: if (items) {
        ensureSelectedExist();
        resetSearch();
    }

    $: filteredItems = filterItems(items, searchTerm);

    $: isSelected = function (item) {
        const normalized = CommonHelper.toArray(selected);

        return CommonHelper.inArray(normalized, item);
    };

    // Selection handlers
    // ---------------------------------------------------------------
    export function deselectItem(item) {
        if (CommonHelper.isEmpty(selected)) {
            return; // nothing to deselect
        }

        let normalized = CommonHelper.toArray(selected);
        if (CommonHelper.inArray(normalized, item)) {
            CommonHelper.removeByValue(normalized, item);
            selected = multiple ? normalized : normalized?.[0] || zeroFunc();
        }

        dispatch("change", { selected });

        // emulate native change event
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function selectItem(item) {
        if (multiple) {
            let normalized = CommonHelper.toArray(selected);
            if (!CommonHelper.inArray(normalized, item)) {
                selected = [...normalized, item];
            }
        } else {
            selected = item;
        }

        dispatch("change", { selected });

        // emulate native change event
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function toggleItem(item) {
        return isSelected(item) ? deselectItem(item) : selectItem(item);
    }

    export function reset() {
        selected = zeroFunc();

        dispatch("change", { selected });

        // emulate native change event
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function showDropdown() {
        toggler?.show && toggler?.show();
    }

    export function hideDropdown() {
        toggler?.hide && toggler?.hide();
    }

    function ensureSelectedExist() {
        if (CommonHelper.isEmpty(selected) || CommonHelper.isEmpty(items)) {
            return; // nothing to check
        }

        let selectedArray = CommonHelper.toArray(selected);
        let unselectedArray = [];

        // find missing
        for (const selectedItem of selectedArray) {
            if (!CommonHelper.inArray(items, selectedItem)) {
                unselectedArray.push(selectedItem);
            }
        }

        // trigger reactivity
        if (unselectedArray.length) {
            for (const item of unselectedArray) {
                CommonHelper.removeByValue(selectedArray, item);
            }

            selected = multiple ? selectedArray : selectedArray[0];
        }
    }

    // Search handlers
    // ---------------------------------------------------------------
    function defaultSearchFunc(item, search) {
        let normalizedSearch = ("" + search).replace(/\s+/g, "").toLowerCase();
        let normalizedItem = item;

        try {
            if (typeof item === "object" && item !== null) {
                normalizedItem = JSON.stringify(item);
            }
        } catch (e) {}

        return ("" + normalizedItem).replace(/\s+/g, "").toLowerCase().includes(normalizedSearch);
    }

    function resetSearch() {
        searchTerm = "";
    }

    function filterItems(items, search) {
        items = items || [];

        const filterFunc = searchFunc || defaultSearchFunc;

        return items.filter((item) => filterFunc(item, search)) || [];
    }

    // Option actions
    // ---------------------------------------------------------------
    function handleOptionSelect(e, item) {
        e.preventDefault();

        if (toggle && multiple) {
            toggleItem(item);
        } else {
            selectItem(item);
        }
    }

    function handleOptionKeypress(e, item) {
        if (e.code === "Enter" || e.code === "Space") {
            handleOptionSelect(e, item);
            if (closable) {
                hideDropdown();
            }
        }
    }

    function onDropdownShow() {
        resetSearch();

        // ensure that the first selected option is visible
        setTimeout(() => {
            const selected = container?.querySelector(".dropdown-item.option.selected");
            if (selected) {
                selected.focus();
                selected.scrollIntoView({ block: "nearest" });
            }
        }, 0);
    }

    // Label(s) activation
    // ---------------------------------------------------------------
    function onLabelClick(e) {
        e.stopPropagation();

        !readonly && !disabled && toggler?.toggle();
    }

    onMount(() => {
        const labels = document.querySelectorAll(`label[for="${id}"]`);

        for (const label of labels) {
            label.addEventListener("click", onLabelClick);
        }

        return () => {
            for (const label of labels) {
                label.removeEventListener("click", onLabelClick);
            }
        };
    });
</script>

<div bind:this={container} class="select {classes}" class:upside class:multiple class:disabled class:readonly>
    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
    <div
        bind:this={labelDiv}
        tabindex={disabled || readonly ? "-1" : "0"}
        class="selected-container"
        class:disabled
        class:readonly
        role="button"
    >
        {#each CommonHelper.toArray(selected) as item, i}
            <div class="option">
                {#if labelComponent}
                    <svelte:component this={labelComponent} {item} {...labelComponentProps} />
                {:else}
                    <span class="txt">{item}</span>
                {/if}

                {#if multiple || toggle}
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <!-- svelte-ignore a11y-no-static-element-interactions -->
                    <span
                        class="clear"
                        use:tooltip={"Clear"}
                        on:click|preventDefault|stopPropagation={() => deselectItem(item)}
                    >
                        <i class="ri-close-line" />
                    </span>
                {/if}
            </div>
        {:else}
            <div class="block txt-placeholder" class:link-hint={!disabled && !readonly}>
                {selectPlaceholder}
            </div>
        {/each}
    </div>

    {#if !disabled && !readonly}
        <Toggler
            bind:this={toggler}
            class="dropdown dropdown-block options-dropdown dropdown-left {upside ? 'dropdown-upside' : ''}"
            trigger={labelDiv}
            on:show={onDropdownShow}
            on:hide
        >
            {#if searchable}
                <div class="form-field form-field-sm options-search">
                    <label class="input-group">
                        <div class="addon p-r-0">
                            <i class="ri-search-line" />
                        </div>
                        <!-- svelte-ignore a11y-autofocus -->
                        <input
                            autofocus
                            type="text"
                            placeholder={searchPlaceholder}
                            bind:value={searchTerm}
                        />

                        {#if searchTerm.length}
                            <div class="addon suffix p-r-5">
                                <button
                                    type="button"
                                    class="btn btn-sm btn-circle btn-transparent clear"
                                    on:click|preventDefault|stopPropagation={resetSearch}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            </div>
                        {/if}
                    </label>
                </div>
            {/if}

            <slot name="beforeOptions" />

            <div class="options-list">
                {#each filteredItems as item}
                    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
                    <!-- svelte-ignore a11y-no-static-element-interactions -->
                    <div
                        tabindex="0"
                        class="dropdown-item option"
                        role="menuitem"
                        class:closable
                        class:selected={isSelected(item)}
                        on:click={(e) => handleOptionSelect(e, item)}
                        on:keydown={(e) => handleOptionKeypress(e, item)}
                    >
                        {#if optionComponent}
                            <svelte:component this={optionComponent} {item} {...optionComponentProps} />
                        {:else}{item}{/if}
                    </div>
                {:else}
                    {#if noOptionsText}
                        <div class="txt-missing">{noOptionsText}</div>
                    {/if}
                {/each}
            </div>

            <slot name="afterOptions" />
        </Toggler>
    {/if}
</div>
