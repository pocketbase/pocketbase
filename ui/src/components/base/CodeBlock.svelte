<script>
    export let content = "";
    export let language = "javascript"; // javascript, html, dart, go, sql

    let classes = "";
    export { classes as class }; // export reserved keyword

    let formattedContent = "";

    $: if (typeof Prism !== "undefined" && content) {
        formattedContent = highlight(content);
    }

    function highlight(code) {
        code = typeof code === "string" ? code : "";

        // @see https://prismjs.com/plugins/normalize-whitespace
        code = Prism.plugins.NormalizeWhitespace.normalize(code, {
            "remove-trailing": true,
            "remove-indent": true,
            "left-trim": true,
            "right-trim": true,
        });

        return Prism.highlight(code, Prism.languages[language] || Prism.languages.javascript, language);
    }
</script>

<div class="code-wrapper prism-light {classes}">
    <code>{@html formattedContent}</code>
</div>

<style>
    code {
        display: block;
        width: 100%;
        padding: 10px 15px;
        white-space: pre-wrap;
        word-break: break-word;
    }
    .code-wrapper {
        display: block;
        width: 100%;
    }
    .prism-light code {
        color: var(--txtPrimaryColor);
        background: var(--baseAlt1Color);
    }
</style>
