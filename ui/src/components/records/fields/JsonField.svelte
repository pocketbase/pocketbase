<script>
    import { onMount } from "svelte";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    export let field;
    export let value = undefined;

    let editorComponent;

    let serialized = serialize(value);

    $: if (value !== serialized?.trim()) {
        serialized = serialize(value);
        value = serialized;
    }

    $: isValid = isValidJson(serialized);

    function serialize(val) {
        if (typeof val == "string" && isValidJson(val)) {
            return val; // already serlialized
        }

        return JSON.stringify(typeof val === "undefined" ? null : val, null, 2);
    }

    function isValidJson(val) {
        try {
            JSON.parse(val === "" ? null : val);
            return true;
        } catch (_) {}

        return false;
    }

    onMount(async () => {
        try {
            editorComponent = (await import("@/components/base/CodeEditor.svelte")).default;
        } catch (err) {
            console.warn(err);
        }
    });
</script>

<Field class="form-field {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <FieldLabel {uniqueId} {field}>
        <span
            class="json-state"
            use:tooltip={{ position: "left", text: isValid ? "Valid JSON" : "Invalid JSON" }}
        >
            {#if isValid}
                <i class="ri-checkbox-circle-fill txt-success" />
            {:else}
                <i class="ri-error-warning-fill txt-danger" />
            {/if}
        </span>
    </FieldLabel>

    {#if editorComponent}
        <svelte:component
            this={editorComponent}
            id={uniqueId}
            maxHeight="500"
            language="json"
            value={serialized}
            on:change={(e) => {
                serialized = e.detail;
                value = serialized.trim();
            }}
        />
    {:else}
        <input type="text" class="txt-mono" value="Loading..." disabled />
    {/if}
</Field>

<style>
    .json-state {
        position: absolute;
        right: 10px;
    }
</style>
