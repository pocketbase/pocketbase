<script>
    import { logLevels } from "@/utils/CommonHelper";

    export let level;

    $: label = logLevels.find((l) => l.level == level)?.label;
</script>

<div class="label log-level-label level-{level}">
    <span class="txt">
        {label || "UNKN"} ({level})
    </span>
</div>

<style lang="scss">
    .log-level-label {
        min-width: 75px;
        font-weight: 600;
        font-size: var(--xsFontSize);
        &:before {
            content: "";
            width: 5px;
            height: 5px;
            border-radius: 5px;
            background: var(--baseAlt4Color);
        }
        &.level--8:before {
            background: var(--primaryColor);
        }
        &.level-0:before {
            background: var(--infoColor);
        }
        &.level-4:before {
            background: var(--warningColor);
        }
        &.level-8:before {
            background: var(--dangerColor);
        }
    }
</style>
