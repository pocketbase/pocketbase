<script>
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";

    export let key;
    export let label;
    export let duration;
    export let secret;
</script>

<Field class="form-field required" name="{key}.duration" let:uniqueId>
    <label for={uniqueId}>{label} duration (in seconds)</label>
    <input type="number" id={uniqueId} required bind:value={duration} />
    <div class="help-block">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <span
            class="link-primary"
            class:txt-success={!!secret}
            on:click={() => {
                // toggle
                if (secret) {
                    secret = undefined;
                } else {
                    secret = CommonHelper.randomSecret(50);
                }
            }}
        >
            Invalidate all previously issued tokens
        </span>
    </div>
</Field>
