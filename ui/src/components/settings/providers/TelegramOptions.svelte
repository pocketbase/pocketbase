<script>
    import Field from "@/components/base/Field.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";

    export let key = "";
    export let config = {
        botToken: ""
    };
</script>

<Field class="form-field required" name="{key}.botToken" let:uniqueId>
    <label for={uniqueId}>Bot Token</label>
    <RedactedPasswordInput type="text" id={uniqueId} bind:value={config.botToken} required />
    <div class="help-block">
        Check <a href="https://core.telegram.org/bots#how-do-i-create-a-bot" target="_blank">How Do I Create a Bot?</a>
    </div>
</Field>
