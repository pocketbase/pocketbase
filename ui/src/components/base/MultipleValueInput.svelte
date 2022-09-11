<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let value = [];
    export let separator = ",";

    $: valueStr = (value || []).join(", ");
</script>

<input
    type={$$restProps.type || "text"}
    value={valueStr}
    on:input={(e) => {
        value = CommonHelper.splitNonEmpty(e.target.value, separator);
    }}
    {...$$restProps}
/>
