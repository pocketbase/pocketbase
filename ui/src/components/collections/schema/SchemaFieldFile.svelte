<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import MimeTypeSelectOption from "@/components/base/MimeTypeSelectOption.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import baseMimeTypesList from "@/mimes.js";
    import CommonHelper from "@/utils/CommonHelper";

    export let field;
    export let key = "";

    const isSingleOptions = [
        { label: "Single", value: true },
        { label: "Multiple", value: false },
    ];

    let mimeTypesList = baseMimeTypesList.slice();
    let isSingle = field.options?.maxSelect <= 1;
    let oldIsSingle = isSingle;

    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    } else {
        appendMissingMimeTypes();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.options.maxSelect = 1;
        } else {
            field.options.maxSelect = field.options?.values?.length || 99;
        }
    }

    function loadDefaults() {
        field.options = {
            maxSelect: 1,
            maxSize: 5242880,
            thumbs: [],
            mimeTypes: [],
        };
        isSingle = true;
        oldIsSingle = isSingle;
    }

    // append any previously set custom mime types to the predefined
    // list for backward compatibility
    function appendMissingMimeTypes() {
        if (CommonHelper.isEmpty(field.options.mimeTypes)) {
            return;
        }

        const missing = [];

        for (const v of field.options.mimeTypes) {
            if (!!mimeTypesList.find((item) => item.mimeType === v)) {
                continue; // exist
            }

            missing.push({ mimeType: v });
        }

        if (missing.length) {
            mimeTypesList = mimeTypesList.concat(missing);
        }
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment let:interactive>
        <div class="separator" />

        <Field
            class="form-field form-field-single-multiple-select {!interactive ? 'readonly' : ''}"
            inlineError
            let:uniqueId
        >
            <ObjectSelect
                id={uniqueId}
                items={isSingleOptions}
                readonly={!interactive}
                bind:keyOfSelected={isSingle}
            />
        </Field>

        <div class="separator" />
    </svelte:fragment>

    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-12">
                <Field class="form-field" name="schema.{key}.options.mimeTypes" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Allowed mime types</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: "Allow files ONLY with the listed mime types. \n Leave empty for no restriction.",
                                position: "top",
                            }}
                        />
                    </label>
                    <ObjectSelect
                        id={uniqueId}
                        multiple
                        searchable
                        closable={false}
                        selectionKey="mimeType"
                        selectPlaceholder="No restriction"
                        items={mimeTypesList}
                        labelComponent={MimeTypeSelectOption}
                        optionComponent={MimeTypeSelectOption}
                        bind:keyOfSelected={field.options.mimeTypes}
                    />
                    <div class="help-block">
                        <div tabindex="0" role="button" class="inline-flex flex-gap-0">
                            <span class="txt link-primary">Choose presets</span>
                            <i class="ri-arrow-drop-down-fill" aria-hidden="true" />
                            <Toggler class="dropdown dropdown-sm dropdown-nowrap dropdown-left">
                                <button
                                    type="button"
                                    class="dropdown-item closable"
                                    role="menuitem"
                                    on:click={() => {
                                        field.options.mimeTypes = [
                                            "image/jpeg",
                                            "image/png",
                                            "image/svg+xml",
                                            "image/gif",
                                            "image/webp",
                                        ];
                                    }}
                                >
                                    <span class="txt">Images (jpg, png, svg, gif, webp)</span>
                                </button>
                                <button
                                    type="button"
                                    class="dropdown-item closable"
                                    role="menuitem"
                                    on:click={() => {
                                        field.options.mimeTypes = [
                                            "application/pdf",
                                            "application/msword",
                                            "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
                                            "application/vnd.ms-excel",
                                            "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
                                        ];
                                    }}
                                >
                                    <span class="txt">Documents (pdf, doc/docx, xls/xlsx)</span>
                                </button>
                                <button
                                    type="button"
                                    class="dropdown-item closable"
                                    role="menuitem"
                                    on:click={() => {
                                        field.options.mimeTypes = [
                                            "video/mp4",
                                            "video/x-ms-wmv",
                                            "video/quicktime",
                                            "video/3gpp",
                                        ];
                                    }}
                                >
                                    <span class="txt">Videos (mp4, avi, mov, 3gp)</span>
                                </button>
                                <button
                                    type="button"
                                    class="dropdown-item closable"
                                    role="menuitem"
                                    on:click={() => {
                                        field.options.mimeTypes = [
                                            "application/zip",
                                            "application/x-7z-compressed",
                                            "application/x-rar-compressed",
                                        ];
                                    }}
                                >
                                    <span class="txt">Archives (zip, 7zip, rar)</span>
                                </button>
                            </Toggler>
                        </div>
                    </div>
                </Field>
            </div>

            <div class={!isSingle ? "col-sm-6" : "col-sm-8"}>
                <Field class="form-field" name="schema.{key}.options.thumbs" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">Thumb sizes</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: "List of additional thumb sizes for image files, along with the default thumb size of 100x100. The thumbs are generated lazily on first access.",
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id={uniqueId}
                        placeholder="eg. 50x50, 480x720"
                        bind:value={field.options.thumbs}
                    />
                    <div class="help-block">
                        <span class="txt">Use comma as separator.</span>
                        <button type="button" class="inline-flex flex-gap-0">
                            <span class="txt link-primary">Supported formats</span>
                            <i class="ri-arrow-drop-down-fill" aria-hidden="true" />
                            <Toggler class="dropdown dropdown-sm dropdown-center dropdown-nowrap p-r-10">
                                <ul class="m-0">
                                    <li>
                                        <strong>WxH</strong>
                                        (eg. 100x50) - crop to WxH viewbox (from center)
                                    </li>
                                    <li>
                                        <strong>WxHt</strong>
                                        (eg. 100x50t) - crop to WxH viewbox (from top)
                                    </li>
                                    <li>
                                        <strong>WxHb</strong>
                                        (eg. 100x50b) - crop to WxH viewbox (from bottom)
                                    </li>
                                    <li>
                                        <strong>WxHf</strong>
                                        (eg. 100x50f) - fit inside a WxH viewbox (without cropping)
                                    </li>
                                    <li>
                                        <strong>0xH</strong>
                                        (eg. 0x50) - resize to H height preserving the aspect ratio
                                    </li>
                                    <li>
                                        <strong>Wx0</strong>
                                        (eg. 100x0) - resize to W width preserving the aspect ratio
                                    </li>
                                </ul>
                            </Toggler>
                        </button>
                    </div>
                </Field>
            </div>

            <div class={!isSingle ? "col-sm-3" : "col-sm-4"}>
                <Field class="form-field required" name="schema.{key}.options.maxSize" let:uniqueId>
                    <label for={uniqueId}>Max file size</label>
                    <input type="number" id={uniqueId} step="1" min="0" bind:value={field.options.maxSize} />
                    <div class="help-block">Must be in bytes.</div>
                </Field>
            </div>

            {#if !isSingle}
                <div class="col-sm-3">
                    <Field class="form-field required" name="schema.{key}.options.maxSelect" let:uniqueId>
                        <label for={uniqueId}>Max select</label>
                        <input
                            id={uniqueId}
                            type="number"
                            step="1"
                            min="2"
                            required
                            bind:value={field.options.maxSelect}
                        />
                    </Field>
                </div>
            {/if}
        </div>
    </svelte:fragment>

    <svelte:fragment slot="optionsFooter">
        <Field class="form-field form-field-toggle" name="schema.{key}.options.protected" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={field.options.protected} />
            <label for={uniqueId}>
                <span class="txt">Protected</span>
            </label>
            <a
                href={import.meta.env.PB_PROTECTED_FILE_DOCS}
                class="toggle-info txt-sm txt-hint m-l-5"
                target="_blank"
                rel="noopener"
            >
                (Learn more)
            </a>
        </Field>
    </svelte:fragment>
</SchemaField>

<style>
    :global(.form-field-file-max-select) {
        width: 100px;
        flex-shrink: 0;
    }
</style>
