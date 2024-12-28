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
    let isSingle = field.maxSelect <= 1;
    let oldIsSingle = isSingle;

    $: if (typeof field.maxSelect == "undefined") {
        loadDefaults();
    } else {
        appendMissingMimeTypes();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.maxSelect = 1;
        } else {
            field.maxSelect = 99;
        }
    }

    function loadDefaults() {
        field.maxSelect = 1;
        field.thumbs = [];
        field.mimeTypes = [];

        isSingle = true;
        oldIsSingle = isSingle;
    }

    // append any previously set custom mime types to the predefined
    // list for backward compatibility
    function appendMissingMimeTypes() {
        if (CommonHelper.isEmpty(field.mimeTypes)) {
            return;
        }

        const missing = [];

        for (const v of field.mimeTypes) {
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
                <Field class="form-field" name="fields.{key}.mimeTypes" let:uniqueId>
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
                        bind:keyOfSelected={field.mimeTypes}
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
                                        field.mimeTypes = [
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
                                        field.mimeTypes = [
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
                                        field.mimeTypes = [
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
                                        field.mimeTypes = [
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
                <Field class="form-field" name="fields.{key}.thumbs" let:uniqueId>
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
                        placeholder="e.g. 50x50, 480x720"
                        bind:value={field.thumbs}
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
                                        (e.g. 100x50) - crop to WxH viewbox (from center)
                                    </li>
                                    <li>
                                        <strong>WxHt</strong>
                                        (e.g. 100x50t) - crop to WxH viewbox (from top)
                                    </li>
                                    <li>
                                        <strong>WxHb</strong>
                                        (e.g. 100x50b) - crop to WxH viewbox (from bottom)
                                    </li>
                                    <li>
                                        <strong>WxHf</strong>
                                        (e.g. 100x50f) - fit inside a WxH viewbox (without cropping)
                                    </li>
                                    <li>
                                        <strong>0xH</strong>
                                        (e.g. 0x50) - resize to H height preserving the aspect ratio
                                    </li>
                                    <li>
                                        <strong>Wx0</strong>
                                        (e.g. 100x0) - resize to W width preserving the aspect ratio
                                    </li>
                                </ul>
                            </Toggler>
                        </button>
                    </div>
                </Field>
            </div>

            <div class={!isSingle ? "col-sm-3" : "col-sm-4"}>
                <Field class="form-field" name="fields.{key}.maxSize" let:uniqueId>
                    <label for={uniqueId}>Max file size</label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        min="0"
                        max={Number.MAX_SAFE_INTEGER}
                        value={field.maxSize || ""}
                        on:input={(e) => (field.maxSize = parseInt(e.target.value, 10))}
                        placeholder="Default to max ~5MB"
                    />
                    <div class="help-block">Must be in bytes.</div>
                </Field>
            </div>

            {#if !isSingle}
                <div class="col-sm-3">
                    <Field class="form-field" name="fields.{key}.maxSelect" let:uniqueId>
                        <label for={uniqueId}>Max select</label>
                        <input
                            id={uniqueId}
                            type="number"
                            step="1"
                            min="2"
                            max={Number.MAX_SAFE_INTEGER}
                            required
                            placeholder="Default to single"
                            bind:value={field.maxSelect}
                        />
                    </Field>
                </div>
            {/if}

            <Field class="form-field form-field-toggle" name="fields.{key}.protected" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={field.protected} />
                <label for={uniqueId}>
                    <span class="txt">Protected</span>
                </label>
                <small class="txt-hint">
                    it will require View API rule permissions and file token to be accessible
                    <a
                        href={import.meta.env.PB_PROTECTED_FILE_DOCS}
                        class="toggle-info"
                        target="_blank"
                        rel="noopener"
                    >
                        (Learn more)
                    </a>
                </small>
            </Field>
        </div>
    </svelte:fragment>
</SchemaField>

<style>
    :global(.form-field-file-max-select) {
        width: 100px;
        flex-shrink: 0;
    }
</style>
