<script>
    import CommonHelper from "@/utils/CommonHelper";

    export let value;

    let classes = undefined;
    export { classes as class }; // export reserved keyword
</script>

<figure class="thumb thumb-circle {classes}">
    <span class="txt initials">{CommonHelper.getInitials(value)}</span>
    <slot />
</figure>

<style>
    .initials {
        font-weight: bold;
        font-size: 0.9em;
        text-transform: uppercase;
        letter-spacing: 2px;
        text-align: center;
        line-height: 1;
        margin-right: -1px;
        color: var(--bodyColor);
    }
    .thumb {
        background: var(--primaryColor);
    }
</style>
