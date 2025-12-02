<script>
    import { onMount } from "svelte";
    import { isDarkMode } from "@/stores/app";
    import ApiClient from "@/utils/ApiClient";

    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import TinyMCE from "@/components/base/TinyMCE.svelte";
    import RecordFilePicker from "@/components/records/RecordFilePicker.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let field;
    export let value = "";

    let picker;
    let editor;
    let mounted = false;
    let mountedTimeoutId = null;

    $: conf = Object.assign(CommonHelper.defaultEditorOptions(), {
        convert_urls: field.convertURLs,
        relative_urls: false,
        content_style: `
            body {
                font-size: 14px;
                font-family: 'Source Sans 3', sans-serif, emoji;
                color: ${$isDarkMode ? "#dcdcdc" : "#1a1a24"};
                background: ${$isDarkMode ? "#1c1c21" : "#ffffff"};
            }
            a { color: #5499e8; }
            code {
                background-color: ${$isDarkMode ? "#303038" : "#e8e8e8"};
                color: inherit;
                border-radius: 3px;
                padding: 0.1rem 0.2rem;
            }
            blockquote {
                border-left: 2px solid ${$isDarkMode ? "#42424e" : "#ccc"};
                margin-left: 1.5rem;
                padding-left: 1rem;
            }
            table th, table td {
                border-color: ${$isDarkMode ? "#42424e" : "#ccc"};
            }
            hr {
                border-color: ${$isDarkMode ? "#42424e" : "#ccc"};
            }
        `,
    });

    // normalize value
    // (depending on the editor plugins, `undefined` may throw an error in case the TinyMCE text functions are used)
    $: if (typeof value == "undefined") {
        value = "";
    }

    onMount(async () => {
        if (typeof value == "undefined") {
            value = "";
        }

        // slight "offset" the editor mount to avoid blocking the rendering of the other fields
        mountedTimeoutId = setTimeout(() => {
            mounted = true;
        }, 100);

        return () => {
            clearTimeout(mountedTimeoutId);
        };
    });
</script>

<Field class="form-field form-field-editor {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <FieldLabel {uniqueId} {field} />

    {#if mounted}
        {#key $isDarkMode}
            <TinyMCE
                id={uniqueId}
                {conf}
                bind:value
                on:init={(initEvent) => {
                    editor = initEvent.detail.editor;
                    editor.on("collections_file_picker", () => {
                        picker?.show();
                    });
                }}
            />
        {/key}
    {:else}
        <div class="tinymce-wrapper" />
    {/if}
</Field>

<RecordFilePicker
    bind:this={picker}
    title="Select an image"
    fileTypes={["image"]}
    on:submit={(e) => {
        editor?.execCommand(
            "InsertImage",
            false,
            ApiClient.files.getURL(e.detail.record, e.detail.name, {
                thumb: e.detail.size,
            }),
        );
    }}
/>
