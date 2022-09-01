<script>
    import { replace, querystring } from "svelte-spa-router";
    import { Collection } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { pageTitle, hideControls } from "@/stores/app";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import SortHeader from "@/components/base/SortHeader.svelte";
    import IdLabel from "@/components/base/IdLabel.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import UserUpsertPanel from "@/components/users/UserUpsertPanel.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";
    import RecordFieldCell from "@/components/records/RecordFieldCell.svelte";

    $pageTitle = "Users";

    const queryParams = new URLSearchParams($querystring);
    const excludedProfileFields = ["id", "userId", "created", "updated"];

    let userUpsertPanel;
    let collectionUpsertPanel;
    let recordUpsertPanel;
    let users = [];
    let currentPage = 1;
    let totalItems = 0;
    let isLoadingUsers = false;
    let filter = queryParams.get("filter") || "";
    let sort = queryParams.get("sort") || "-created";
    let profileCollection = new Collection();
    let isLoadingProfileCollection = false;

    $: if (sort !== -1 && filter !== -1) {
        // keep query params
        const query = new URLSearchParams({ filter, sort }).toString();
        replace("/users?" + query);

        loadUsers();
    }

    $: canLoadMore = totalItems > users.length;

    $: profileFields = profileCollection?.schema?.filter(
        (field) => !excludedProfileFields.includes(field.name)
    );

    loadProfilesCollection();

    export async function loadUsers(page = 1) {
        isLoadingUsers = true;

        if (page <= 1) {
            clearList();
        }

        return ApiClient.users
            .getList(page, 50, {
                sort: sort || "-created",
                filter: filter,
            })
            .then((result) => {
                isLoadingUsers = false;
                users = users.concat(result.items);
                currentPage = result.page;
                totalItems = result.totalItems;
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoadingUsers = false;
                    console.warn(err);
                    clearList();
                    ApiClient.errorResponseHandler(err, false);
                }
            });
    }

    function clearList() {
        users = [];
        currentPage = 1;
        totalItems = 0;
    }

    function setUserProfile(profile) {
        const user = users.find((u) => u.id === profile?.userId);
        if (user) {
            user.profile = profile;
        }
        users = users;
    }

    async function loadProfilesCollection() {
        isLoadingProfileCollection = true;

        try {
            profileCollection = await ApiClient.collections.getOne(import.meta.env.PB_PROFILE_COLLECTION);
        } catch (err) {
            ApiClient.errorResponseHandler(err);
        }

        isLoadingProfileCollection = false;
    }
</script>

<PageWrapper>
    {#if isLoadingProfileCollection}
        <div class="placeholder-section m-b-base">
            <span class="loader loader-lg" />
            <h1>Loading users...</h1>
        </div>
    {:else}
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">{$pageTitle}</div>
            </nav>

            {#if !$hideControls}
                <button
                    type="button"
                    class="btn btn-secondary btn-circle"
                    use:tooltip={{ text: "Edit profile collection", position: "right" }}
                    on:click={() => collectionUpsertPanel?.show(profileCollection)}
                >
                    <i class="ri-settings-4-line" />
                </button>
            {/if}

            <RefreshButton on:refresh={() => loadUsers()} />

            <div class="flex-fill" />

            <button type="button" class="btn btn-expanded" on:click={() => userUpsertPanel?.show()}>
                <i class="ri-add-line" />
                <span class="txt">New user</span>
            </button>
        </header>

        <Searchbar
            value={filter}
            placeholder={"Search filter, eg. verified=1"}
            extraAutocompleteKeys={["verified", "email"]}
            on:submit={(e) => (filter = e.detail)}
        />

        <div class="table-wrapper">
            <table class="table" class:table-loading={isLoadingUsers}>
                <thead>
                    <tr>
                        <SortHeader class="col-type-text col-field-id" name="id" bind:sort>
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

                        {#each profileFields as field (field.name)}
                            <th class="col-type-{field.type} col-field-{field.name}" name={field.name}>
                                <div class="col-header-content">
                                    <i class={CommonHelper.getFieldTypeIcon(field.type)} />
                                    <span class="txt">profile.{field.name}</span>
                                </div>
                            </th>
                        {/each}

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
                    {#each users as user (user.id)}
                        <tr>
                            <td class="col-type-text col-field-id">
                                <IdLabel id={user.id} />
                            </td>

                            <td class="col-type-email col-field-email">
                                <div class="inline-flex">
                                    {#if user.email}
                                        <span class="txt" title={user.email}>{user.email}</span>
                                        <span
                                            class="label"
                                            class:label-success={user.verified}
                                            class:label-warning={!user.verified}
                                        >
                                            {user.verified ? "Verified" : "Unverified"}
                                        </span>
                                    {:else}
                                        <div class="txt-hint">N/A</div>
                                        {#if user.verified}
                                            <span class="label label-success">OAuth2 verified</span>
                                        {/if}
                                    {/if}
                                </div>
                            </td>

                            {#each profileFields as field (field.name)}
                                <RecordFieldCell {field} record={user.profile || {}} />
                            {/each}

                            <td class="col-type-date col-field-created">
                                <FormattedDate date={user.created} />
                            </td>

                            <td class="col-type-date col-field-updated">
                                <FormattedDate date={user.updated} />
                            </td>

                            <td class="col-type-action min-width">
                                <button
                                    type="button"
                                    class="btn btn-sm btn-outline"
                                    on:click|stopPropagation={() => userUpsertPanel?.show(user)}
                                >
                                    <i class="ri-user-settings-line" />
                                    <span class="txt">Edit user</span>
                                </button>
                                <button
                                    type="button"
                                    class="btn btn-sm m-l-10"
                                    on:click|stopPropagation={() => recordUpsertPanel?.show(user.profile)}
                                >
                                    <i class="ri-profile-line" />
                                    <span class="txt">Edit profile</span>
                                </button>
                            </td>
                        </tr>
                    {:else}
                        {#if isLoadingUsers}
                            <tr>
                                <td colspan="99" class="p-xs">
                                    <span class="skeleton-loader" />
                                </td>
                            </tr>
                        {:else}
                            <tr>
                                <td colspan="99" class="txt-center txt-hint p-xs">
                                    <h6>No users found.</h6>
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
        </div>

        {#if users.length}
            <small class="block txt-hint txt-right m-t-sm">Showing {users.length} of {totalItems}</small>
        {/if}

        {#if users.length && canLoadMore}
            <div class="block txt-center m-t-xs">
                <button
                    type="button"
                    class="btn btn-lg btn-secondary btn-expanded"
                    class:btn-loading={isLoadingUsers}
                    class:btn-disabled={isLoadingUsers}
                    on:click={() => loadUsers(currentPage + 1)}
                >
                    <span class="txt">Load more ({totalItems - users.length})</span>
                </button>
            </div>
        {/if}
    {/if}
</PageWrapper>

<UserUpsertPanel bind:this={userUpsertPanel} on:save={() => loadUsers()} on:delete={() => loadUsers()} />

<CollectionUpsertPanel bind:this={collectionUpsertPanel} on:save={(e) => (profileCollection = e.detail)} />

<RecordUpsertPanel
    bind:this={recordUpsertPanel}
    collection={profileCollection}
    on:save={(e) => setUserProfile(e.detail)}
/>
