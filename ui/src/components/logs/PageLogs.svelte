<script>
    import { pageTitle } from "@/stores/app";
    import Field from "@/components/base/Field.svelte";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import LogsList from "@/components/logs/LogsList.svelte";
    import LogsChart from "@/components/logs/LogsChart.svelte";
    import LogViewPanel from "@/components/logs/LogViewPanel.svelte";

    const ADMIN_LOGS_LOCAL_STORAGE_KEY = "includeAdminLogs";

    $pageTitle = "Request logs";

    let logPanel;
    let filter = "";
    let includeAdminLogs = window.localStorage?.getItem(ADMIN_LOGS_LOCAL_STORAGE_KEY) << 0;
    let refreshToken = 1;

    $: presets = !includeAdminLogs ? 'auth!="admin"' : "";

    $: if (typeof includeAdminLogs !== "undefined" && window.localStorage) {
        window.localStorage.setItem(ADMIN_LOGS_LOCAL_STORAGE_KEY, includeAdminLogs << 0);
    }

    function refresh() {
        refreshToken++;
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

            <div class="inline-flex">
                <Field class="form-field form-field-toggle m-0" let:uniqueId>
                    <input type="checkbox" id={uniqueId} bind:checked={includeAdminLogs} />
                    <label for={uniqueId}>Include requests by admins</label>
                </Field>
            </div>
        </header>

        <Searchbar
            value={filter}
            placeholder="Search logs, ex. status > 200"
            extraAutocompleteKeys={[
                "method",
                "url",
                "remoteIp",
                "userIp",
                "referer",
                "status",
                "auth",
                "userAgent",
            ]}
            on:submit={(e) => (filter = e.detail)}
        />

        <div class="clearfix m-b-xs" />

        {#key refreshToken}
            <LogsChart bind:filter {presets} />
        {/key}
    </div>

    {#key refreshToken}
        <LogsList bind:filter {presets} on:select={(e) => logPanel?.show(e?.detail)} />
    {/key}
</PageWrapper>

<LogViewPanel bind:this={logPanel} />
