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

        for (let v of newKeyOfSelected) {
            const item = CommonHelper.findByKey(items, selectionKey, v);
            if (item) {
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
</script>

<Select bind:selected {items} {multiple} {labelComponent} {optionComponent} on:show on:hide {...$$restProps}>
    <svelte:fragment slot="afterOptions">
        <slot name="afterOptions" />
    </svelte:fragment>
</Select>
