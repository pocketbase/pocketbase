<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";

    const SDK_PREFERENCE_KEY = "pb_sdk_preference";

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
        },
        {
            title: "Dart",
            language: "dart",
            content: dart,
        },
    ];
</script>

<div class="tabs sdk-tabs m-b-lg">
    <div class="tabs-header compact left">
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
            </div>
        {/each}
    </div>
</div>

<style>
    .sdk-tabs .tabs-header .tab-item {
        min-width: 100px;
    }
</style>
