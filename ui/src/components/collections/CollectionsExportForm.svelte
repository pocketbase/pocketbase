<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import { onMount } from "svelte";

    const uniqueId = "exports_" + CommonHelper.randomString(5);

    let collections = [];
    let isLoadingCollections = false;

    loadCollections();

    async function loadCollections() {
        isLoadingCollections = true;

        try {
            collections = await ApiClient.collections.getFullList(100, {
                $cancelKey: uniqueId,
            });
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingCollections = false;
    }

    let oldSchema = "";
    let newSchema = "";

    function diff_prettyHtml(diffs, showInsert) {
        const html = [];
        const pattern_amp = /&/g;
        const pattern_lt = /</g;
        const pattern_gt = />/g;
        const pattern_para = /\n/g;
        for (let x = 0; x < diffs.length; x++) {
            let op = diffs[x][0]; // Operation (insert, delete, equal)
            let data = diffs[x][1]; // Text of change.
            let text = data
                .replace(pattern_amp, "&amp;")
                .replace(pattern_lt, "&lt;")
                .replace(pattern_gt, "&gt;")
                .replace(pattern_para, "<br>");
            // text = CommonHelper.stripTags(text);
            switch (op) {
                case DIFF_INSERT:
                    if (showInsert) {
                        html[x] = '<ins class="block">' + text + "</ins>";
                    }
                    break;
                case DIFF_DELETE:
                    if (!showInsert) {
                        html[x] = '<del class="block">' + text + "</del>";
                    }
                    break;
                case DIFF_EQUAL:
                    html[x] = "<span>" + text + "</span>";
                    break;
            }
        }
        return html.join("");
    }

    onMount(() => {
        var dmp = new diff_match_patch();
        const text1 = [
            {
                id: "zwWlxR46txtoAwx",
                created: "2022-08-01 17:32:24.329",
                updated: "2022-08-04 10:19:57.248",
                name: "profia sdles",
                system: true,
                listRule: "userId = @request.user.id",
                viewRule: "userId = @request.user.id",
                createRule: "userId = @request.user.id",
                updateRule: "userId = @request.user.id",
                deleteRule: null,
                schema: [
                    {
                        id: "nsght7oy",
                        name: "userId",
                        type: "user",
                        system: true,
                        required: true,
                        unique: true,
                        options: {
                            maxSelect: 1,
                            cascadeDelete: true,
                        },
                    },
                    {
                        id: "atpc4yjm",
                        name: "name",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                    {
                        id: "akb4s9de",
                        name: "avatar",
                        type: "file",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            maxSelect: 1,
                            maxSize: 5242880,
                            mimeTypes: ["image/jpg", "image/jpeg", "image/png", "image/svg+xml", "image/gif"],
                            thumbs: null,
                        },
                    },
                ],
            },
            {
                id: "IV8FbE78jmXF56d",
                created: "",
                updated: "2022-08-04 10:21:54.100",
                name: "abc",
                system: false,
                listRule: null,
                viewRule: null,
                createRule: null,
                updateRule: null,
                deleteRule: null,
                schema: [
                    {
                        id: "t2pukeas",
                        name: "demo",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                    {
                        id: "dddddd",
                        name: "aaaa",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                    {
                        id: "squmamtm",
                        name: "test",
                        type: "date",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: "",
                            max: "",
                        },
                    },
                ],
            },
        ];

        const text2 = [
            {
                id: "zwWlxR46txtoAwx",
                created: "2022-08-01 17:32:24.329",
                updated: "2022-08-04 10:19:57.248",
                name: "Demo",
                system: true,
                listRule: "userId = @request.user.id",
                viewRule: "userId = @request.user.id",
                createRule: "userId = @request.user.id",
                updateRule: "userId = @request.user.id",
                deleteRule: null,
                schema: [
                    {
                        id: "nsght7oy",
                        name: "userId",
                        type: "user",
                        system: true,
                        required: true,
                        unique: true,
                        options: {
                            maxSelect: 1,
                            cascadeDelete: true,
                        },
                    },
                    {
                        id: "atpc4yjm",
                        name: "name",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                    {
                        id: "akb4s9de",
                        name: "avatar",
                        type: "file",
                        system: false,
                        required: false,
                        unique: true,
                        options: {
                            maxSelect: 1,
                            maxSize: 5242880,
                            mimeTypes: ["image/jpg", "image/jpeg", "image/png", "image/svg+xml", "image/gif"],
                            thumbs: null,
                        },
                    },
                ],
            },
            {
                id: "IV8FbE78jmXF56d",
                created: "",
                updated: "2022-08-04 10:21:54.100",
                name: "abc",
                system: false,
                listRule: null,
                viewRule: null,
                createRule: null,
                updateRule: null,
                deleteRule: null,
                schema: [
                    {
                        id: "t2pukeas",
                        name: "demo",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                    {
                        id: "dddddd",
                        name: "aaaa",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                    {
                        id: "squmamtm",
                        name: "test",
                        type: "date",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: "",
                            max: "",
                        },
                    },
                ],
            },
            {
                id: "GGACt8sa1tcJp7T",
                created: "2022-08-04 10:22:15.871",
                updated: "2022-08-04 10:22:15.871",
                name: "asdasd",
                system: true,
                listRule: null,
                viewRule: null,
                createRule: null,
                updateRule: null,
                deleteRule: null,
                schema: [
                    {
                        id: "0eklwfvl",
                        name: "field",
                        type: "text",
                        system: false,
                        required: false,
                        unique: false,
                        options: {
                            min: null,
                            max: null,
                            pattern: "",
                        },
                    },
                ],
            },
        ];

        // var diffs = dmp.diff_main(JSON.stringify(text1, null, 2), JSON.stringify(text2, null, 2));

        var a = dmp.diff_linesToChars_(JSON.stringify(text1, null, 2), JSON.stringify(text2, null, 2));
        var lineText1 = a.chars1;
        var lineText2 = a.chars2;
        var lineArray = a.lineArray;
        var diffs = dmp.diff_main(lineText1, lineText2, false);
        dmp.diff_charsToLines_(diffs, lineArray);

        oldSchema = diff_prettyHtml(diffs, false);
        newSchema = diff_prettyHtml(diffs, true);
    });
</script>

<br />
<div class="grid">
    <div class="col-6">
        <code>
            {@html oldSchema}
        </code>
    </div>
    <div class="col-6">
        <code>
            {@html newSchema}
        </code>
    </div>
</div>

<style lang="scss">
    .collections-list {
        column-count: 2;
        column-gap: var(--baseSpacing);
    }
    code {
        display: block;
        width: 100%;
        overflow: auto;
        padding: var(--xsSpacing);
        white-space: pre;
        background: var(--baseAlt1Color);
    }
</style>
