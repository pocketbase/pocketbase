<script>
    import { createEventDispatcher } from "svelte";
    import { slide } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import { setErrors } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import ModelDateIcon from "@/components/base/ModelDateIcon.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import SecretGeneratorButton from "@/components/base/SecretGeneratorButton.svelte";

    const dispatch = createEventDispatcher();
    const formId = "admin_" + CommonHelper.randomString(5);

    let panel;
    let admin = {};
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let avatar = 0;
    let email = "";
    let password = "";
    let passwordConfirm = "";
    let changePasswordToggle = false;

    $: isNew = !admin?.id;

    $: hasChanges =
        (isNew && email != "") || changePasswordToggle || email !== admin.email || avatar !== admin.avatar;

    export function show(model) {
        load(model);

        confirmClose = true;

        return panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function load(model) {
        admin = structuredClone(model || {});
        reset(); // reset form
    }

    function reset() {
        changePasswordToggle = false;
        email = admin?.email || "";
        avatar = admin?.avatar || 0;
        password = "";
        passwordConfirm = "";
        setErrors({}); // reset errors
    }

    function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        const data = { email, avatar };
        if (isNew || changePasswordToggle) {
            data["password"] = password;
            data["passwordConfirm"] = passwordConfirm;
        }

        let request;
        if (isNew) {
            request = ApiClient.admins.create(data);
        } else {
            request = ApiClient.admins.update(admin.id, data);
        }

        request
            .then(async (result) => {
                confirmClose = false;
                hide();
                addSuccessToast(isNew ? "Successfully created admin." : "Successfully updated admin.");

                if (ApiClient.authStore.model?.id === result.id) {
                    ApiClient.authStore.save(ApiClient.authStore.token, result);
                }

                dispatch("save", result);
            })
            .catch((err) => {
                ApiClient.error(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function deleteConfirm() {
        if (!admin?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete the selected admin?`, () => {
            return ApiClient.admins
                .delete(admin.id)
                .then(() => {
                    confirmClose = false;
                    hide();
                    addSuccessToast("Successfully deleted admin.");
                    dispatch("delete", admin);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }
</script>

<OverlayPanel
    bind:this={panel}
    popup
    class="admin-panel"
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
        <h4>
            {isNew ? "New admin" : "Edit admin"}
        </h4>
    </svelte:fragment>

    <form id={formId} class="grid" autocomplete="off" on:submit|preventDefault={save}>
        {#if !isNew}
            <Field class="form-field readonly" name="id" let:uniqueId>
                <label for={uniqueId}>
                    <i class={CommonHelper.getFieldTypeIcon("primary")} />
                    <span class="txt">id</span>
                </label>
                <div class="form-field-addon">
                    <ModelDateIcon model={admin} />
                </div>
                <input type="text" id={uniqueId} value={admin.id} readonly />
            </Field>
        {/if}

        <div class="content">
            <p class="section-title">Avatar</p>
            <div class="flex flex-gap-xs flex-wrap">
                {#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as index}
                    <button
                        type="button"
                        class="link-fade thumb thumb-circle {index == avatar ? 'thumb-primary' : 'thumb-sm'}"
                        on:click={() => (avatar = index)}
                    >
                        <img
                            src="{import.meta.env.BASE_URL}images/avatars/avatar{index}.svg"
                            alt="Avatar {index}"
                        />
                    </button>
                {/each}
            </div>
        </div>

        <Field class="form-field required" name="email" let:uniqueId>
            <label for={uniqueId}>
                <i class={CommonHelper.getFieldTypeIcon("email")} />
                <span class="txt">Email</span>
            </label>
            <input type="email" autocomplete="off" id={uniqueId} required bind:value={email} />
        </Field>

        {#if !isNew}
            <Field class="form-field form-field-toggle" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={changePasswordToggle} />
                <label for={uniqueId}>Change password</label>
            </Field>
        {/if}

        {#if isNew || changePasswordToggle}
            <div class="col-12">
                <div class="grid" transition:slide={{ duration: 150 }}>
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
                            <div class="form-field-addon">
                                <SecretGeneratorButton />
                            </div>
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
    </form>

    <svelte:fragment slot="footer">
        {#if !isNew}
            <div
                tabindex="0"
                role="button"
                aria-label="More admin options"
                class="btn btn-sm btn-circle btn-transparent"
            >
                <!-- empty span for alignment -->
                <span aria-hidden="true" />
                <i class="ri-more-line" aria-hidden="true" />
                <Toggler class="dropdown dropdown-upside dropdown-left dropdown-nowrap">
                    <button
                        type="button"
                        class="dropdown-item txt-danger"
                        role="menuitem"
                        on:click={() => deleteConfirm()}
                    >
                        <i class="ri-delete-bin-7-line" aria-hidden="true" />
                        <span class="txt">Delete</span>
                    </button>
                </Toggler>
            </div>
            <div class="flex-fill" />
        {/if}

        <button type="button" class="btn btn-transparent" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!hasChanges || isSaving}
        >
            <span class="txt">{isNew ? "Create" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
