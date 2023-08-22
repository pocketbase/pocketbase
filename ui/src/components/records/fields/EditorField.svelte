<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import TinyMCE from "@tinymce/tinymce-svelte";

    export let field;
    export let value = undefined;

    $: conf = Object.assign(CommonHelper.defaultEditorOptions(), {
        convert_urls: field.options?.convertUrls,
        relative_urls: false,
    });
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    <TinyMCE
        id={uniqueId}
        scriptSrc="{import.meta.env.BASE_URL}libs/tinymce/tinymce.min.js"
        {conf}
        bind:value
    />
</Field>
