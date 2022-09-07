<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";

    export let key = "";
    export let options = {};

    $: if (CommonHelper.isEmpty(options)) {
        // load defaults
        options = {
            maxSelect: 1,
            maxSize: 5242880,
            thumbs: [],
            mimeTypes: [],
        };
    }
</script>

<div class="grid">
    <div class="col-sm-6">
        <Field class="form-field required" name="schema.{key}.options.maxSize" let:uniqueId>
            <label for={uniqueId}>Max file size (bytes)</label>
            <input type="number" id={uniqueId} step="1" min="0" bind:value={options.maxSize} />
        </Field>
    </div>

    <div class="col-sm-6">
        <Field class="form-field required" name="schema.{key}.options.maxSelect" let:uniqueId>
            <label for={uniqueId}>Max files</label>
            <input type="number" id={uniqueId} step="1" min="" required bind:value={options.maxSelect} />
        </Field>
    </div>

    <div class="col-sm-12">
        <Field class="form-field" name="schema.{key}.options.mimeTypes" let:uniqueId>
            <label for={uniqueId}>
                <span class="txt">Mime types</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: "Allow files ONLY with the listed mime types. \n Leave empty for no restriction.",
                        position: "top",
                    }}
                />
            </label>
            <MultipleValueInput
                id={uniqueId}
                placeholder="eg. image/png, application/pdf..."
                bind:value={options.mimeTypes}
            />
            <div class="help-block">
                <span class="txt">Use comma as separator.</span>
                <button type="button" class="inline-flex flex-gap-0">
                    <span class="txt link-primary">Choose presets</span>
                    <i class="ri-arrow-drop-down-fill" />
                    <Toggler class="dropdown dropdown-sm dropdown-nowrap">
                        <div
                            tabindex="0"
                            class="dropdown-item closable"
                            on:click={() => {
                                options.mimeTypes = [
                                    "application/pdf",
                                    "application/msword",
                                    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
                                    "application/vnd.ms-excel",
                                    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
                                ];
                            }}
                        >
                            <span class="txt">Documents (pdf, doc/docx, xls/xlsx)</span>
                        </div>
                        <div
                            tabindex="0"
                            class="dropdown-item closable"
                            on:click={() => {
                                options.mimeTypes = [
                                    "image/jpg",
                                    "image/jpeg",
                                    "image/png",
                                    "image/svg+xml",
                                    "image/gif",
                                ];
                            }}
                        >
                            <span class="txt">Images (jpg, png, svg, gif)</span>
                        </div>
                        <div
                            tabindex="0"
                            class="dropdown-item closable"
                            on:click={() => {
                                options.mimeTypes = [
                                    "video/mp4",
                                    "video/x-ms-wmv",
                                    "video/quicktime",
                                    "video/3gpp",
                                ];
                            }}
                        >
                            <span class="txt">Videos (mp4, avi, mov, 3gp)</span>
                        </div>
                        <div
                            tabindex="0"
                            class="dropdown-item closable"
                            on:click={() => {
                                options.mimeTypes = [
                                    "application/zip",
                                    "application/x-7z-compressed",
                                    "application/x-rar-compressed",
                                ];
                            }}
                        >
                            <span class="txt">Archives (zip, 7zip, rar)</span>
                        </div>
                    </Toggler>
                </button>
            </div>
        </Field>
    </div>

    <div class="col-sm-12">
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
            <MultipleValueInput id={uniqueId} placeholder="eg. 50x50, 480x720" bind:value={options.thumbs} />
            <div class="help-block">
                <span class="txt">Use comma as separator.</span>
                <button type="button" class="inline-flex flex-gap-0">
                    <span class="txt link-primary">Supported formats</span>
                    <i class="ri-arrow-drop-down-fill" />
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
</div>
