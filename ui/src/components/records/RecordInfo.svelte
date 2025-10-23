<script>
    import tooltip from "@/actions/tooltip";
    import RecordInfoContent from "@/components/records/RecordInfoContent.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let record;

    function excludeProps(item, ...props) {
        const result = Object.assign({}, item);
        for (let prop of props) {
            delete result[prop];
        }
        return result;
    }
</script>

<div class="record-info-excerpt">
    <div class="info-content">
        <RecordInfoContent {record} />
    </div>

    <a
        href="#/collections?collection={record.collectionId}&recordId={record.id}"
        target="_blank"
        class="record-link link-hint"
        rel="noopener noreferrer"
        use:tooltip={{
            text:
                "Open relation record in new tab:\n" +
                CommonHelper.truncate(
                    JSON.stringify(CommonHelper.truncateObject(excludeProps(record, "expand")), null, 2),
                    800,
                    true,
                ),
            class: "code",
            position: "left",
        }}
        on:click|stopPropagation
        on:keydown|stopPropagation
    >
        <i class="ri-external-link-line txt-sm"></i>
    </a>
</div>
