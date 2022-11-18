<script>
    import { replace, querystring } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { admin as loggedAdmin } from "@/stores/admin";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import SortHeader from "@/components/base/SortHeader.svelte";
    import IdLabel from "@/components/base/IdLabel.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import HorizontalScroller from "@/components/base/HorizontalScroller.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import AdminUpsertPanel from "@/components/admins/AdminUpsertPanel.svelte";

    $pageTitle = "Admins";

    const queryParams = new URLSearchParams($querystring);

    let adminUpsertPanel;
    let admins = [];
    let isLoading = false;
    let filter = queryParams.get("filter") || "";
    let sort = queryParams.get("sort") || "-created";

    $: if (sort !== -1 && filter !== -1) {
        // keep listing params in sync
        const query = new URLSearchParams({ filter, sort }).toString();
        replace("/settings/admins?" + query);

        loadAdmins();
    }

    export function loadAdmins() {
        isLoading = true;

        admins = []; // reset

        return ApiClient.admins
            .getFullList(100, {
                sort: sort || "-created",
                filter: filter,
            })
            .then((result) => {
                admins = result;
                isLoading = false;
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoading = false;
                    console.warn(err);
                    clearList();
                    ApiClient.errorResponseHandler(err, false);
                }
            });
    }

    function clearList() {
        admins = [];
    }
</script>

<SettingsSidebar />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>

        <RefreshButton on:refresh={() => loadAdmins()} />

        <div class="flex-fill" />

        <button type="button" class="btn btn-expanded" on:click={() => adminUpsertPanel?.show()}>
            <i class="ri-add-line" />
            <span class="txt">New admin</span>
        </button>
    </header>

    <Searchbar
        value={filter}
        placeholder={"Search filter, eg. email='test@example.com'"}
        extraAutocompleteKeys={["email"]}
        on:submit={(e) => (filter = e.detail)}
    />

    <HorizontalScroller class="table-wrapper">
        <table class="table" class:table-loading={isLoading}>
            <thead>
                <tr>
                    <th class="min-width" />

                    <SortHeader class="col-type-text" name="id" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("primary")} />
                            <span class="txt">id</span>
                        </div>
                    </SortHeader>

                    <SortHeader class="col-type-email col-field-email" name="email" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("email")} />
                            <span class="txt">email</span>
                        </div>
                    </SortHeader>

                    <SortHeader class="col-type-date col-field-created" name="created" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("date")} />
                            <span class="txt">created</span>
                        </div>
                    </SortHeader>

                    <SortHeader class="col-type-date col-field-updated" name="updated" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("date")} />
                            <span class="txt">updated</span>
                        </div>
                    </SortHeader>

                    <th class="col-type-action min-width" />
                </tr>
            </thead>
            <tbody>
                {#each admins as admin (admin.id)}
                    <tr
                        tabindex="0"
                        class="row-handle"
                        on:click={() => adminUpsertPanel?.show(admin)}
                        on:keydown={(e) => {
                            if (e.code === "Enter" || e.code === "Space") {
                                e.preventDefault();
                                adminUpsertPanel?.show(admin);
                            }
                        }}
                    >
                        <td class="min-width">
                            <figure class="thumb thumb-sm thumb-circle">
                                <img
                                    src="{import.meta.env.BASE_URL}images/avatars/avatar{admin.avatar ||
                                        0}.svg"
                                    alt="Admin avatar"
                                />
                            </figure>
                        </td>

                        <td class="col-type-text col-field-id">
                            <IdLabel id={admin.id} />
                            {#if admin.id === $loggedAdmin.id}
                                <span class="label label-warning m-l-5">You</span>
                            {/if}
                        </td>

                        <td class="col-type-email col-field-email">
                            <span class="txt txt-ellipsis" title={admin.email}>
                                {admin.email}
                            </span>
                        </td>

                        <td class="col-type-date col-field-created">
                            <FormattedDate date={admin.created} />
                        </td>

                        <td class="col-type-date col-field-updated">
                            <FormattedDate date={admin.updated} />
                        </td>

                        <td class="col-type-action min-width">
                            <i class="ri-arrow-right-line" />
                        </td>
                    </tr>
                {:else}
                    {#if isLoading}
                        <tr>
                            <td colspan="99" class="p-xs">
                                <span class="skeleton-loader" />
                            </td>
                        </tr>
                    {:else}
                        <tr>
                            <td colspan="99" class="txt-center txt-hint p-xs">
                                <h6>No admins found.</h6>
                                {#if filter?.length}
                                    <button
                                        type="button"
                                        class="btn btn-hint btn-expanded m-t-sm"
                                        on:click={() => (filter = "")}
                                    >
                                        <span class="txt">Clear filters</span>
                                    </button>
                                {/if}
                            </td>
                        </tr>
                    {/if}
                {/each}
            </tbody>
        </table>
    </HorizontalScroller>

    {#if admins.length}
        <small class="block txt-hint txt-right m-t-sm">Showing {admins.length} of {admins.length}</small>
    {/if}
</PageWrapper>

<AdminUpsertPanel bind:this={adminUpsertPanel} on:save={() => loadAdmins()} on:delete={() => loadAdmins()} />
