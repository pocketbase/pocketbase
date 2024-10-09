<script>
    import { onMount } from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltipAction from "@/actions/tooltip";

    export let value = "";
    export let tooltip = "Copy";
    export let idleClasses = "ri-file-copy-line txt-sm link-hint";
    export let successClasses = "ri-check-line txt-sm txt-success";
    export let successDuration = 500; // ms

    let copyTimeout;

    function copy() {
        if (CommonHelper.isEmpty(value)) {
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
    tabindex="-1"
    role="button"
    class={copyTimeout ? successClasses : idleClasses}
    aria-label={"Copy to clipboard"}
    use:tooltipAction={!copyTimeout ? tooltip : undefined}
    on:click|stopPropagation={copy}
/>
