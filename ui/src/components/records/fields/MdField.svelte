<script>
    import 'bytemd/dist/index.css'
    import "highlight.js/styles/default.css";

    import {SchemaField} from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import {Editor} from 'bytemd'
    import gfm from '@bytemd/plugin-gfm'
    import frontmatter from '@bytemd/plugin-frontmatter'
    import highlight from '@bytemd/plugin-highlight'

    export let field = new SchemaField();
    export let value = undefined;

    const plugins = [
        gfm(),
        frontmatter(),
        highlight()
    ]

    function handleChange(e) {
        value = e.detail.value
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label>
        <i class={CommonHelper.getFieldTypeIcon(field.type)}/>
        <span class="txt">{field.name}</span>
    </label>
    <Editor {value} {plugins} on:change={handleChange}/>
</Field>

<style>
    :global(.bytemd code) {
        white-space: inherit;
    }

    :global(.bytemd-fullscreen.bytemd) {
        z-index: 1;
    }

    :global(.bytemd label) {
        display: initial;
        padding-top: 0;
        padding-bottom: 0;
        background: none;
        font-weight: initial;
        text-transform: initial;
    }

    :global(.bytemd input[type=checkbox]) {
        display: initial;
        padding-top: 0;
        padding-bottom: 0;
        background: none;
        position: initial;
        opacity: 1;
        width: initial;
        height: initial;
    }

    /* Change editor and viewer background color */
    /*:global(.bytemd .CodeMirror-scroll), :global(.bytemd .bytemd-preview) {*/
    /*    background: var(--baseAlt1Color);*/
    /*}*/

    :global(.bytemd .markdown-body .task-list-item) {
        list-style-type: none;
    }

    :global(.bytemd .form-field:focus-within label) {
        color: initial;
        background: none;
    }

    label {
        padding-bottom: 12px;
    }
</style>