<script>
    import BaseSelectOption from "@/components/base/BaseSelectOption.svelte";
    import Select from "@/components/base/Select.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    // original select props
    export let items = [];
    export let multiple = false;
    export let selected = multiple ? [] : undefined;
    export let labelComponent = BaseSelectOption; // custom component to use for each selected option label
    export let optionComponent = BaseSelectOption; // custom component to use for each dropdown option item

    // custom props
    export let selectionKey = "value";
    export let keyOfSelected = multiple ? [] : undefined;

    let oldKeyOfSelectedHash = JSON.stringify(keyOfSelected);

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
        if (!items.length) {
            return; // options are still loading...
        }

        let extractedKeys = CommonHelper.toArray(newSelected, true).map((item) => item[selectionKey]);
        let newKeyOfSelected = multiple ? extractedKeys : extractedKeys[0];

        if (JSON.stringify(newKeyOfSelected) != oldKeyOfSelectedHash) {
            keyOfSelected = newKeyOfSelected;
            oldKeyOfSelectedHash = JSON.stringify(keyOfSelected);
        }
    }
</script>

<Select
    bind:selected
    {items}
    {multiple}
    {labelComponent}
    {optionComponent}
    on:show
    on:hide
    on:change
    {...$$restProps}
>
    <svelte:fragment slot="afterOptions">
        <slot name="afterOptions" />
    </svelte:fragment>
</Select>
