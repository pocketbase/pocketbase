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
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import Scroller from "@/components/base/Scroller.svelte";
    import CopyIcon from "@/components/base/CopyIcon.svelte";
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

        const normalizedFilter = CommonHelper.normalizeSearchFilter(filter, [
            "id",
            "email",
            "created",
            "updated",
        ]);

        return ApiClient.admins
            .getFullList(100, {
                sort: sort || "-created",
                filter: normalizedFilter,
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
                    ApiClient.error(err, !normalizedFilter || err?.status != 400); // silence filter errors
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

        <div class="btns-group">
            <button type="button" class="btn btn-expanded" on:click={() => adminUpsertPanel?.show()}>
                <i class="ri-add-line" />
                <span class="txt">New admin</span>
            </button>
        </div>
    </header>

    <Searchbar
        value={filter}
        placeholder={"Search term or filter like email='test@example.com'"}
        extraAutocompleteKeys={["email"]}
        on:submit={(e) => (filter = e.detail)}
    />
    <div class="clearfix m-b-base" />

    <Scroller class="table-wrapper">
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
                            <div class="label">
                                <CopyIcon value={admin.id} />
                                <span class="txt">{admin.id}</span>
                            </div>
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
                                <span class="skeleton-loader m-0" />
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
    </Scroller>

    <svelte:fragment slot="footer">
        <div class="m-r-auto txt-sm txt-hint">Total found: {admins.length}</div>
    </svelte:fragment>
</PageWrapper>

<AdminUpsertPanel bind:this={adminUpsertPanel} on:save={() => loadAdmins()} on:delete={() => loadAdmins()} />
