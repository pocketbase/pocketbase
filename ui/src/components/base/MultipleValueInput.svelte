<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let value = [];
    export let separator = ",";
    export let readonly = null;
    export let disabled = null;

    $: valueStr = CommonHelper.joinNonEmpty(value, separator + " ");
</script>

<input
    type={$$restProps.type || "text"}
    value={valueStr}
    {disabled}
    {readonly}
    on:input={(e) => {
        value = CommonHelper.splitNonEmpty(e.target.value, separator);
    }}
    {...$$restProps}
/>
