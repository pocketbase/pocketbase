<script>
    import CommonHelper from "../../utils/CommonHelper";

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
                return model[prop];
            }
        }

        return "";
    }
</script>

{#if meta !== "" && meta !== item.id}
    <span class="label txt-base txt-mono" title={meta}>{meta}</span>
{:else}
    <span class="txt txt-hint">N/A</span>
{/if}