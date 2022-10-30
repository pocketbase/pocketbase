<script>
    import { SchemaField } from "pocketbase";
    import Flatpickr from "svelte-flatpickr";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;

    // strip ms and zone for backwards compatibility with the older format
    // and because flatpickr currently doesn't have integrated
    // zones support and requires manual parsing and formatting
    $: if (value && value.length > 19) {
        value = value.substring(0, 19);
    }
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <label for={uniqueId}>
        <i class={CommonHelper.getFieldTypeIcon(field.type)} />
        <span class="txt">{field.name} (UTC)</span>
    </label>
    <Flatpickr
        id={uniqueId}
        options={CommonHelper.defaultFlatpickrOptions()}
        {value}
        bind:formattedValue={value}
    />
</Field>
