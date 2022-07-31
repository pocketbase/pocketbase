<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Select from "@/components/base/Select.svelte";
    import BaseSelectOption from "@/components/base/BaseSelectOption.svelte";

    // original select props
    export let items = [];
    export let multiple = false;
    export let selected = multiple ? [] : undefined;
    export let labelComponent = BaseSelectOption; // custom component to use for each selected option label
    export let optionComponent = BaseSelectOption; // custom component to use for each dropdown option item

    // custom props
    export let selectionKey = "value";
    export let keyOfSelected = multiple ? [] : undefined;

    $: if (items) {
        handleKeyOfSelectedChange(keyOfSelected);
    }

    $: handleSelectedChange(selected);

    function handleKeyOfSelectedChange(newKeyOfSelected) {
        newKeyOfSelected = CommonHelper.toArray(newKeyOfSelected, true);

        let newSelected = [];
        let allItems = getFlattenItems();

        for (let item of allItems) {
            if (CommonHelper.inArray(newKeyOfSelected, item[selectionKey])) {
                newSelected.push(item);
            }
        }

        if (newKeyOfSelected.length && !newSelected.length) {
            return; // options are still loading...
        }

        selected = multiple ? newSelected : newSelected[0];
    }

    async function handleSelectedChange(newSelected) {
        let extractedKeys = CommonHelper.toArray(newSelected, true).map((item) => item[selectionKey]);

        if (!items.length) {
            return; // options are still loading...
        }

        keyOfSelected = multiple ? extractedKeys : extractedKeys[0];
    }

    function getFlattenItems() {
        if (!CommonHelper.isObjectArrayWithKeys(items, ["group", "items"])) {
            return items; // already flatten
        }

        // extract items from groups
        let result = [];
        for (const group of items) {
            result = result.concat(group.items);
        }

        return result;
    }
</script>

<Select bind:selected {items} {multiple} {labelComponent} {optionComponent} on:show on:hide {...$$restProps}>
    <svelte:fragment slot="afterOptions">
        <slot name="afterOptions" />
    </svelte:fragment>
</Select>
