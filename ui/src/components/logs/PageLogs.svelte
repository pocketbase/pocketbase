<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import LogsList from "@/components/logs/LogsList.svelte";
    import LogsChart from "@/components/logs/LogsChart.svelte";
    import LogViewPanel from "@/components/logs/LogViewPanel.svelte";
    import { _ } from '@/services/i18n';
import { tooltips } from "@codemirror/view";

    const ADMIN_LOGS_LOCAL_STORAGE_KEY = "includeAdminLogs";

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

    CommonHelper.setDocumentTitle($_("logs.pagetitle"));
</script>

<main class="page-wrapper">
    <div class="page-header-wrapper m-b-0">
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">{$_("logs.pagetitle")}</div>
            </nav>

            <button
                type="button"
                class="btn btn-circle btn-secondary"
                use:tooltip={{ text: $_("logs.tips.refresh"), position: "right" }}
                on:click={refresh}
            >
                <i class="ri-refresh-line" />
            </button>

            <div class="flex-fill" />

            <div class="inline-flex">
                <Field class="form-field form-field-toggle m-0" let:uniqueId>
                    <input type="checkbox" id={uniqueId} bind:checked={includeAdminLogs} />
                    <label for={uniqueId}>{$_("logs.tips.includeAdminLogs")}</label>
                </Field>
            </div>
        </header>

        <Searchbar
            value={filter}
            placeholder="{$_("logs.tips.search.placeholder")}"
            extraAutocompleteKeys={["method", "url", "ip", "referer", "status", "auth", "userAgent"]}
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
</main>

<LogViewPanel bind:this={logPanel} />
