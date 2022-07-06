<script>
    import { onMount } from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Toggler from "@/components/base/Toggler.svelte";

    const baseGroup = "_base_"; // reserved items group name

    export let id = "";
    export let noOptionsText = "No options found";
    export let selectPlaceholder = "- Select -";
    export let searchPlaceholder = "Search...";
    export let items = []; // for groups support wrap in `[{group: 'My group', items: [...]}]`
    export let multiple = false;
    export let disabled = false;
    export let selected = multiple ? [] : undefined;
    export let toggle = false; // toggle option on click
    export let labelComponent = undefined; // custom component to use for each selected option label
    export let labelComponentProps = {}; // props to pass to the custom option component
    export let optionComponent = undefined; // custom component to use for each dropdown option item
    export let optionComponentProps = {}; // props to pass to the custom option component
    export let searchable = false; // whether to show the dropdown options search input
    export let searchFunc = undefined; // custom search option filter: `function(item, searchTerm):boolean`

    let classes = "";
    export { classes as class }; // export reserved keyword

    let toggler;
    let searchTerm = "";
    let container = undefined;
    let labelDiv = undefined;

    $: groupedItems = CommonHelper.isObjectArrayWithKeys(items, ["group"])
        ? items
        : [{ group: baseGroup, items: items }];

    $: if (items) {
        ensureSelectedExist();
        resetSearch();
    }

    $: filteredGroups = filterGroups(groupedItems, searchTerm);

    $: isSelected = function (item) {
        let normalized = CommonHelper.toArray(selected);

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
            selected = normalized;
        }

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

        // emulate native change event
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function toggleItem(item) {
        return isSelected(item) ? deselectItem(item) : selectItem(item);
    }

    export function reset() {
        selected = multiple ? [] : undefined;
    }

    export function showDropdown() {
        toggler?.show && toggler?.show();
    }

    export function hideDropdown() {
        toggler?.hide && toggler?.hide();
    }

    function ensureSelectedExist() {
        if (CommonHelper.isEmpty(selected) || CommonHelper.isEmpty(groupedItems)) {
            return; // nothing to check
        }

        let selectedArray = CommonHelper.toArray(selected);
        let unselectedArray = [];

        // find missing
        for (const selectedItem of selectedArray) {
            let exist = false;

            for (const group of groupedItems) {
                if (CommonHelper.inArray(group.items, selectedItem)) {
                    exist = true;
                    break;
                }
            }

            if (!exist) {
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

    function filterGroups(groups, search) {
        const result = [];
        const filterFunc = searchFunc || defaultSearchFunc;

        for (const group of groups) {
            let groupItems;
            if (typeof search === "string" && search.length) {
                groupItems = group.items?.filter((item) => filterFunc(item, search)) || [];
            } else {
                groupItems = group.items || [];
            }

            if (groupItems.length) {
                result.push({ group: group.group, items: groupItems });
            }
        }

        return result;
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

        !disabled && toggler?.toggle();
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

<div class="select {classes}" class:multiple class:disabled bind:this={container}>
    <div tabindex={disabled ? "-1" : "0"} class="selected-container" class:disabled bind:this={labelDiv}>
        {#each CommonHelper.toArray(selected) as item}
            <div class="option">
                {#if labelComponent}
                    <svelte:component this={labelComponent} {item} {...labelComponentProps} />
                {:else}<span class="txt">{item}</span>{/if}

                {#if multiple || toggle}
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
            <div class="txt-placeholder">{selectPlaceholder}</div>
        {/each}
    </div>

    {#if !disabled}
        <Toggler
            class="dropdown dropdown-block options-dropdown dropdown-left"
            trigger={labelDiv}
            on:show={onDropdownShow}
            on:hide
            bind:this={toggler}
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
                                    class="btn btn-sm btn-circle btn-secondary clear"
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
                {#each filteredGroups as group}
                    {#if group.group != baseGroup}
                        <div class="dropdown-item separator">{group.group}</div>
                    {/if}

                    {#each group.items as item}
                        <div
                            tabindex="0"
                            class="dropdown-item option closable"
                            class:selected={isSelected(item)}
                            on:click={(e) => handleOptionSelect(e, item)}
                            on:keydown={(e) => handleOptionKeypress(e, item)}
                        >
                            {#if optionComponent}
                                <svelte:component this={optionComponent} {item} {...optionComponentProps} />
                            {:else}{item}{/if}
                        </div>
                    {/each}
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
