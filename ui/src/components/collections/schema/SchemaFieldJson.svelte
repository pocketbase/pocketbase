<script>
    import { slide } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    let showInfo = false;

    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    }

    function loadDefaults() {
        field.options = {
            maxSize: 2000000,
        };
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <Field class="form-field required m-b-sm" name="schema.{key}.options.maxSize" let:uniqueId>
            <label for={uniqueId}>Max size <small>(bytes)</small></label>
            <input type="number" id={uniqueId} step="1" min="0" bind:value={field.options.maxSize} />
        </Field>

        <button
            type="button"
            class="btn btn-sm {showInfo ? 'btn-secondary' : 'btn-hint btn-transparent'}"
            on:click={() => {
                showInfo = !showInfo;
            }}
        >
            <strong class="txt">String value normalizations</strong>
            {#if showInfo}
                <i class="ri-arrow-up-s-line txt-sm" />
            {:else}
                <i class="ri-arrow-down-s-line txt-sm" />
            {/if}
        </button>
        {#if showInfo}
            <div class="block" transition:slide={{ duration: 150 }}>
                <div class="alert alert-warning m-b-0 m-t-10">
                    <div class="content">
                        In order to support seamlessly both <code>application/json</code> and
                        <code>multipart/form-data</code>
                        requests, the following normalization rules are applied if the <code>json</code> field
                        is a
                        <strong>plain string</strong>:
                        <ul>
                            <li>"true" is converted to the json <code>true</code></li>
                            <li>"false" is converted to the json <code>false</code></li>
                            <li>"null" is converted to the json <code>null</code></li>
                            <li>"[1,2,3]" is converted to the json <code>[1,2,3]</code></li>
                            <li>
                                {'"{"a":1,"b":2}"'} is converted to the json <code>{'{"a":1,"b":2}'}</code>
                            </li>
                            <li>numeric strings are converted to json number</li>
                            <li>double quoted strings are left as they are (aka. without normalizations)</li>
                            <li>any other string (empty string too) is double quoted</li>
                        </ul>
                        Alternatively, if you want to avoid the string value normalizations, you can wrap your
                        data inside an object, eg.<code>{'{"data": anything}'}</code>
                    </div>
                </div>
            </div>
        {/if}
    </svelte:fragment>
</SchemaField>
