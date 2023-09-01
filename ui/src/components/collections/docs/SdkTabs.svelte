<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";

    const SDK_PREFERENCE_KEY = "pb_sdk_preference";

    let classes = "m-b-sm";
    export { classes as class }; // export reserved keyword

    export let js = "";
    export let dart = "";

    let activeTab = localStorage.getItem(SDK_PREFERENCE_KEY) || "javascript";

    $: if (activeTab) {
        // store user preference
        localStorage.setItem(SDK_PREFERENCE_KEY, activeTab);
    }

    $: sdkExamples = [
        {
            title: "JavaScript",
            language: "javascript",
            content: js,
            url: import.meta.env.PB_JS_SDK_URL,
        },
        {
            title: "Dart",
            language: "dart",
            content: dart,
            url: import.meta.env.PB_DART_SDK_URL,
        },
    ];
</script>

<div class="tabs sdk-tabs {classes}">
    <div class="tabs-header compact combined left">
        {#each sdkExamples as example (example.language)}
            <button
                class="tab-item"
                class:active={activeTab === example.language}
                on:click={() => (activeTab = example.language)}
            >
                <div class="txt">{example.title}</div>
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each sdkExamples as example (example.language)}
            <div class="tab-item" class:active={activeTab === example.language}>
                <CodeBlock language={example.language} content={example.content} />
                <div class="txt-right">
                    <em class="txt-sm txt-hint">
                        <a href={example.url} target="_blank" rel="noopener noreferrer">
                            {example.title} SDK
                        </a>
                    </em>
                </div>
            </div>
        {/each}
    </div>
</div>

<style>
    .sdk-tabs .tabs-header .tab-item {
        min-width: 100px;
    }
</style>
