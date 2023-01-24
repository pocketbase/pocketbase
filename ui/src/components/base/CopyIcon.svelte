<script>
    import { onMount } from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";

    export let value = "";
    export let idleClasses = "ri-file-copy-line txt-sm link-hint";
    export let successClasses = "ri-check-line txt-sm txt-success";
    export let successDuration = 500; // ms

    let copyTimeout;

    function copy() {
        if (!value) {
            return;
        }

        CommonHelper.copyToClipboard(value);

        clearTimeout(copyTimeout);
        copyTimeout = setTimeout(() => {
            clearTimeout(copyTimeout);
            copyTimeout = null;
        }, successDuration);
    }

    onMount(() => {
        return () => {
            if (copyTimeout) {
                clearTimeout(copyTimeout);
            }
        };
    });
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<i
    class={copyTimeout ? successClasses : idleClasses}
    use:tooltip={!copyTimeout ? "Copy" : ""}
    on:click|stopPropagation={copy}
/>
