<script>
    import SchemaFieldEmail from "./SchemaFieldEmail.svelte";
    import Field from "@/components/base/Field.svelte";

    export let field;
    export let key = "";
</script>

<!-- shares the same options with the email field -->
<SchemaFieldEmail bind:field {key} on:rename on:remove {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-12">
                <Field class="form-field" name="schema.{key}.title" let:title>
                    <label for={title}>Display Name</label>
                    <input type="text" id={title} bind:value={field.title} />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaFieldEmail>
