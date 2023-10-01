<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import TinyMCE from "@tinymce/tinymce-svelte";
    import { onMount } from "svelte";

    export let field;
    export let value = undefined;

    let mounted = false;
    let mountedTimeoutId = null;

    $: conf = Object.assign(CommonHelper.defaultEditorOptions(), {
        convert_urls: field.options?.convertUrls,
        relative_urls: false,
    });

    // normalize value
    // (depending on the editor plugins, `undefined` may throw an error in case the TinyMCE text functions are used)
    $: if (typeof value == "undefined") {
        value = "";
    }

    onMount(() => {
        mountedTimeoutId = setTimeout(() => {
            mounted = true;
        }, 100);

        return () => {
            clearTimeout(mountedTimeoutId);
        };
    });
</script>

<Field class="form-field form-field-editor {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name}</span>
    </label>
    {#if mounted}
        <TinyMCE
            id={uniqueId}
            scriptSrc="{import.meta.env.BASE_URL}libs/tinymce/tinymce.min.js"
            {conf}
            bind:value
        />
    {:else}
        <div class="tinymce-wrapper" />
    {/if}
</Field>
