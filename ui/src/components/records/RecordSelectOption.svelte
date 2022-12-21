<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";

    const excludedMetaProps = ["id", "created", "updated", "collectionId", "collectionName"];

    export let item = {}; // model

    $: meta = extractMeta(item);

    function extractMeta(model) {
        model = model || {};

        const props = [
            // prioritized common displayable props
            "title",
            "name",
            "email",
            "username",
            "label",
            "key",
            "heading",
            "content",
            "description",
            // fallback to the available props
            ...Object.keys(model),
        ];

        for (const prop of props) {
            if (
                typeof model[prop] === "string" &&
                !CommonHelper.isEmpty(model[prop]) &&
                !excludedMetaProps.includes(prop)
            ) {
                return prop + ": " + model[prop];
            }
        }

        return "";
    }
</script>

<i
    class="ri-information-line link-hint"
    use:tooltip={{ text: JSON.stringify(item, null, 2), position: "left", class: "code" }}
/>

<div class="content">
    <div class="block txt-ellipsis">{item.id}</div>
    {#if meta !== "" && meta !== item.id}
        <small class="block txt-hint txt-ellipsis">{meta}</small>
    {/if}
</div>

<style>
    .content {
        flex-shrink: 1;
        flex-grow: 0;
        width: auto;
        min-width: 0;
    }
</style>
