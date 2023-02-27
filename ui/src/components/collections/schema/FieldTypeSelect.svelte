<script>
    import CommonHelper from "@/utils/CommonHelper";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";

    export let value = "text";

    let classes = "";
    export { classes as class }; // export reserved keyword

    const types = [
        {
            label: "Plain text",
            value: "text",
            icon: CommonHelper.getFieldTypeIcon("text"),
        },
        {
            label: "Rich editor",
            value: "editor",
            icon: CommonHelper.getFieldTypeIcon("editor"),
        },
        {
            label: "Number",
            value: "number",
            icon: CommonHelper.getFieldTypeIcon("number"),
        },
        {
            label: "Bool",
            value: "bool",
            icon: CommonHelper.getFieldTypeIcon("bool"),
        },
        {
            label: "Email",
            value: "email",
            icon: CommonHelper.getFieldTypeIcon("email"),
        },
        {
            label: "Url",
            value: "url",
            icon: CommonHelper.getFieldTypeIcon("url"),
        },
        {
            label: "DateTime",
            value: "date",
            icon: CommonHelper.getFieldTypeIcon("date"),
        },
        {
            label: "Select",
            value: "select",
            icon: CommonHelper.getFieldTypeIcon("select"),
        },
        {
            label: "File",
            value: "file",
            icon: CommonHelper.getFieldTypeIcon("file"),
        },
        {
            label: "Relation",
            value: "relation",
            icon: CommonHelper.getFieldTypeIcon("relation"),
        },
        {
            label: "JSON",
            value: "json",
            icon: CommonHelper.getFieldTypeIcon("json"),
        },
    ];
</script>

<ObjectSelect class="field-type-select {classes}" items={types} bind:keyOfSelected={value} {...$$restProps} />
