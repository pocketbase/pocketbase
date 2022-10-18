<script>
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import { confirm } from "@/stores/confirmation";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import { activeView, addView, removeView } from "@/stores/views";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { View } from "pocketbase";
    import { createEventDispatcher, tick } from "svelte";
    import ViewUpdateConfirm from "../views/ViewUpdateConfirm.svelte";
    import ViewForm from "./ViewForm.svelte";

    const dispatch = createEventDispatcher();

    let viewPanel;
    let confirmViewPanel;

    let original = null;
    let view = new View();
    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let initialFormHash = calculateFormHash(view);

    $: hasChanges = initialFormHash != calculateFormHash(view);

    $: canSave = view.isNew || hasChanges;

    export function changeTab(newTab) {
        activeTab = newTab;
    }

    export function show(model) {
        load(model);

        confirmClose = true;

        return viewPanel?.show();
    }

    export function hide() {
        return viewPanel?.hide();
    }

    async function load(model) {
        setErrors({}); // reset errors
        if (typeof model !== "undefined") {
            original = model;
            view = model?.clone();
        } else {
            original = null;
            view = new View();
        }
        // normalize
        view.originalSql = view.sql || "";
        view.originalName = view.name || "";
        view.originalListRule = view.listRule || "";

        await tick();

        initialFormHash = calculateFormHash(view);
    }

    function saveWithConfirm() {
        if (view.isNew) {
            return save();
        } else {
            confirmViewPanel?.show(view);
        }
    }

    function save() {
        if (isSaving) {
            return;
        }

        isSaving = true;

        const data = exportFormData();
        console.log(data);

        let request;
        if (view.isNew) {
            request = ApiClient.views.create(data);
        } else {
            request = ApiClient.views.update(view.id, data);
        }

        request
            .then((result) => {
                confirmClose = false;
                hide();
                addSuccessToast(view.isNew ? "Successfully created views." : "Successfully updated view.");
                addView(result);

                if (view.isNew) {
                    $activeView = result;
                }

                dispatch("save", result);
            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function exportFormData() {
        return view.export();
    }

    function deleteConfirm() {
        if (!original?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete View "${original?.name}"`, () => {
            return ApiClient.views
                .delete(original?.id)
                .then(() => {
                    hide();
                    addSuccessToast(`Successfully deleted view "${original?.name}".`);
                    dispatch("delete", original);
                    removeView(original);
                })
                .catch((err) => {
                    ApiClient.errorResponseHandler(err);
                });
        });
    }

    function calculateFormHash(m) {
        return JSON.stringify(m);
    }
</script>

<OverlayPanel
    bind:this={viewPanel}
    class="overlay-panel-lg colored-header compact-header collection-panel"
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
            {view.isNew ? "New view" : "Edit view"}
        </h4>

        {#if !view.isNew}
            <div class="flex-fill" />
            <button type="button" class="btn btn-sm btn-circle btn-secondary flex-gap-0">
                <i class="ri-more-line" />
                <Toggler class="dropdown dropdown-right m-t-5">
                    <button type="button" class="dropdown-item closable" on:click={() => deleteConfirm()}>
                        <i class="ri-delete-bin-7-line" />
                        <span class="txt">Delete</span>
                    </button>
                </Toggler>
            </button>
        {/if}

        <form
            class="block"
            on:submit|preventDefault={() => {
                canSave && saveWithConfirm();
            }}
        >
            <Field class="form-field required m-b-0" name="name" let:uniqueId>
                <label for={uniqueId}>Name</label>
                <!-- svelte-ignore a11y-autofocus -->
                <input
                    type="text"
                    id={uniqueId}
                    required
                    spellcheck="false"
                    autofocus={view.isNew}
                    placeholder={`eg. "posts"`}
                    value={view.name}
                    on:input={(e) => {
                        view.name = CommonHelper.slugify(e.target.value);
                        e.target.value = view.name;
                    }}
                />
            </Field>
        </form>
    </svelte:fragment>

    <ViewForm bind:view />

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
            type="button"
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!canSave || isSaving}
            on:click={() => saveWithConfirm()}
        >
            <span class="txt">{view.isNew ? "Create" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<ViewUpdateConfirm bind:this={confirmViewPanel} on:confirm={() => save()} />
