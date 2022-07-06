<script>
    import { SchemaField } from "pocketbase";
    import Flatpickr from "svelte-flatpickr";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let field = new SchemaField();
    export let value = undefined;
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
