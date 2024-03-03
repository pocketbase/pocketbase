<script>
    import { pageTitle } from "@/stores/app";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import Graph from "@/components/graph/Graph.svelte";

    $pageTitle = "Graph";


    let refreshKey = 1;

    function refresh() {
        refreshKey++;
    }
</script>

<PageWrapper>
    <div class="page-header-wrapper m-b-0">
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">{$pageTitle}</div>
            </nav>

            <RefreshButton on:refresh={() => refresh()} />

            <div class="flex-fill" />
        </header>
        <div>
            {#key refreshKey}
                <Graph />
            {/key}
        </div>
    </div>
</PageWrapper>
