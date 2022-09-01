<script>
    import { createEventDispatcher } from "svelte";
    import { slide } from "svelte/transition";
    import { User } from "pocketbase";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { setErrors } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import ExternalAuthsList from "./ExternalAuthsList.svelte";

    const dispatch = createEventDispatcher();

    const formId = "user_" + CommonHelper.randomString(5);
    const accountTab = "account";
    const providersTab = "providers";

    let panel;
    let user = new User();
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let email = "";
    let password = "";
    let passwordConfirm = "";
    let changePasswordToggle = false;
    let verificationEmailToggle = true;
    let activeTab = accountTab;

    $: hasChanges = (user.isNew && email != "") || changePasswordToggle || email !== user.email;

    export function show(model) {
        load(model);

        confirmClose = true;

        activeTab = user.isNew || user.email ? accountTab : providersTab;

        return panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function load(model) {
        setErrors({}); // reset errors

        user = model?.clone ? model.clone() : new User();

        reset(); // reset form
    }

    function reset() {
        changePasswordToggle = false;
        verificationEmailToggle = true;
        email = user?.email || "";
        password = "";
        passwordConfirm = "";
    }

    function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        const data = { email: email };
        if (user.isNew || changePasswordToggle) {
            data["password"] = password;
            data["passwordConfirm"] = passwordConfirm;
        }

        let request;
        if (user.isNew) {
            request = ApiClient.users.create(data);
        } else {
            request = ApiClient.users.update(user.id, data);
        }

        request
            .then(async (result) => {
                user = result;

                if (verificationEmailToggle) {
                    sendVerificationEmail(false);
                }

                confirmClose = false;
                hide();
                addSuccessToast(user.isNew ? "Successfully created user." : "Successfully updated user.");
                dispatch("save", result);
            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function deleteConfirm() {
        if (!user?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete the selected user?`, () => {
            return ApiClient.users
                .delete(user.id)
                .then(() => {
                    confirmClose = false;
                    hide();
                    addSuccessToast("Successfully deleted user.");
                    dispatch("delete", user);
                })
                .catch((err) => {
                    ApiClient.errorResponseHandler(err);
                });
        });
    }

    function sendVerificationEmail(notify = true) {
        return ApiClient.users
            .requestVerification(user.email || email)
            .then(() => {
                confirmClose = false;
                hide();
                if (notify) {
                    addSuccessToast(`Successfully sent verification email to ${user.email}.`);
                }
            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
            });
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="user-panel"
    popup={user.isNew}
    beforeHide={() => {
        if (hasChanges && confirmClose) {
            confirm("You have unsaved changes. Do you really want to close the panel?", () => {
                confirmClose = false;
                hide();
            });
            return false;
        }
        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4>{user.isNew ? "New user" : "Edit user"}</h4>

        {#if !user.isNew}
            <button type="button" class="btn btn-sm btn-circle btn-secondary m-l-auto">
                <!-- empty span for alignment -->
                <span />
                <i class="ri-more-line" />
                <Toggler class="dropdown dropdown-right dropdown-nowrap">
                    {#if !user.verified && user.email}
                        <button type="button" class="dropdown-item" on:click={() => sendVerificationEmail()}>
                            <i class="ri-mail-check-line" />
                            <span class="txt">Send verification email</span>
                        </button>
                    {/if}
                    <button type="button" class="dropdown-item" on:click={() => deleteConfirm()}>
                        <i class="ri-delete-bin-7-line" />
                        <span class="txt">Delete</span>
                    </button>
                </Toggler>
            </button>
        {/if}
    </svelte:fragment>

    <div class="tabs user-tabs">
        {#if !user.isNew}
            <div class="tabs-header stretched">
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === accountTab}
                    on:click={() => (activeTab = accountTab)}
                >
                    Account
                </button>
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === providersTab}
                    on:click={() => (activeTab = providersTab)}
                >
                    Authorized providers
                </button>
            </div>
        {/if}
        <div class="tabs-content">
            <div class="tab-item" class:active={activeTab === accountTab}>
                <form id={formId} class="grid" autocomplete="off" on:submit|preventDefault={save}>
                    {#if !user.isNew}
                        <Field class="form-field disabled" name="id" let:uniqueId>
                            <label for={uniqueId}>
                                <i class={CommonHelper.getFieldTypeIcon("primary")} />
                                <span class="txt">ID</span>
                            </label>
                            <input type="text" id={uniqueId} value={user.id} disabled />
                        </Field>
                    {/if}

                    <Field class="form-field required" name="email" let:uniqueId>
                        <label for={uniqueId}>
                            <i class={CommonHelper.getFieldTypeIcon("email")} />
                            <span class="txt">Email</span>
                        </label>
                        {#if user.verified && user.email}
                            <div class="form-field-addon txt-success" use:tooltip={"Verified"}>
                                <i class="ri-shield-check-line" />
                            </div>
                        {/if}
                        <input type="email" autocomplete="off" id={uniqueId} required bind:value={email} />
                    </Field>

                    {#if !user.isNew && user.email}
                        <Field class="form-field form-field-toggle" let:uniqueId>
                            <input type="checkbox" id={uniqueId} bind:checked={changePasswordToggle} />
                            <label for={uniqueId}>Change password</label>
                        </Field>
                    {/if}

                    {#if user.isNew || !user.email || changePasswordToggle}
                        <div class="col-12">
                            <div class="grid" transition:slide|local={{ duration: 150 }}>
                                <div class="col-sm-6">
                                    <Field class="form-field required" name="password" let:uniqueId>
                                        <label for={uniqueId}>
                                            <i class="ri-lock-line" />
                                            <span class="txt">Password</span>
                                        </label>
                                        <input
                                            type="password"
                                            autocomplete="new-password"
                                            id={uniqueId}
                                            required
                                            bind:value={password}
                                        />
                                    </Field>
                                </div>
                                <div class="col-sm-6">
                                    <Field class="form-field required" name="passwordConfirm" let:uniqueId>
                                        <label for={uniqueId}>
                                            <i class="ri-lock-line" />
                                            <span class="txt">Password confirm</span>
                                        </label>
                                        <input
                                            type="password"
                                            autocomplete="new-password"
                                            id={uniqueId}
                                            required
                                            bind:value={passwordConfirm}
                                        />
                                    </Field>
                                </div>
                            </div>
                        </div>
                    {/if}

                    {#if user.isNew || !user.email}
                        <Field class="form-field form-field-toggle" let:uniqueId>
                            <input type="checkbox" id={uniqueId} bind:checked={verificationEmailToggle} />
                            <label for={uniqueId}>Send verification email</label>
                        </Field>
                    {/if}
                </form>
            </div>

            {#if !user.isNew}
                <div class="tab-item" class:active={activeTab === providersTab}>
                    <ExternalAuthsList {user} />
                </div>
            {/if}
        </div>
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>

        {#if activeTab === accountTab}
            <button
                type="submit"
                form={formId}
                class="btn btn-expanded"
                class:btn-loading={isSaving}
                disabled={!hasChanges || isSaving}
            >
                <span class="txt">{user.isNew ? "Create" : "Save changes"}</span>
            </button>
        {/if}
    </svelte:fragment>
</OverlayPanel>
