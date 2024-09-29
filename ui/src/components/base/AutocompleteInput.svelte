<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let value = "";
    export let options = []; // [{label: "Option 1", value: "opt1"}, {label: "Option 2", value: "opt2"}, ...]

    const uniqueId = "list_" + CommonHelper.randomString(5);
</script>

<input
    type={$$restProps.type || "text"}
    list={uniqueId}
    {value}
    on:input={(e) => {
        value = e.target.value;
    }}
    {...$$restProps}
/>

<datalist id={uniqueId}>
    {#each options as opt}
        <option value={opt.value}>{opt.label || ""}</option>
    {/each}
</datalist>
