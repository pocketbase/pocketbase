<script>
    import { createEventDispatcher, tick } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addErrorToast, addSuccessToast } from "@/stores/toasts";
    import { setErrors } from "@/stores/errors";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Field from "@/components/base/Field.svelte";

    const dispatch = createEventDispatcher();

    const formId = "email_test_" + CommonHelper.randomString(5);
    const emailStorageKey = "last_email_test";
    const testRequestKey = "email_test_request";

    const templateOptions = [
        { label: '"Verification" template', value: "verification" },
        { label: '"Password reset" template', value: "password-reset" },
        { label: '"Confirm email change" template', value: "email-change" },
    ];

    let panel;
    let email = localStorage.getItem(emailStorageKey);
    let template = templateOptions[0].value;
    let isSubmitting = false;
    let testTimeoutId = null;

    $: canSubmit = !!email && !!template;

    export function show(emailArg = "", templateArg = "") {
        email = emailArg || localStorage.getItem(emailStorageKey);
        template = templateArg || templateOptions[0].value;

        setErrors({}); // reset any previous errors

        panel?.show();
    }

    export function hide() {
        clearTimeout(testTimeoutId);
        return panel?.hide();
    }

    async function submit() {
        if (!canSubmit || isSubmitting) {
            return;
        }

        isSubmitting = true;

        // store in local storage for later use
        localStorage?.setItem(emailStorageKey, email);

        // auto cancel the test request after 30sec
        clearTimeout(testTimeoutId);
        testTimeoutId = setTimeout(() => {
            ApiClient.cancelRequest(testRequestKey);
            addErrorToast("Test email send timeout.");
        }, 30000);

        try {
            await ApiClient.settings.testEmail(email, template, {
                $cancelKey: testRequestKey,
            });

            addSuccessToast("Successfully sent test email.");
            dispatch("submit");
            isSubmitting = false;

            await tick();

            hide();
        } catch (err) {
            isSubmitting = false;
            ApiClient.error(err);
        }

        clearTimeout(testTimeoutId);
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="overlay-panel-sm email-test-popup"
    overlayClose={!isSubmitting}
    escClose={!isSubmitting}
    beforeHide={() => !isSubmitting}
    popup
    on:show
    on:hide
>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">Send test email</h4>
    </svelte:fragment>

    <form id={formId} autocomplete="off" on:submit|preventDefault={() => submit()}>
        <Field class="form-field required" name="template" let:uniqueId>
            {#each templateOptions as option (option.value)}
                <div class="form-field-block">
                    <input
                        type="radio"
                        name="template"
                        id={uniqueId + option.value}
                        value={option.value}
                        bind:group={template}
                    />
                    <label for={uniqueId + option.value}>{option.label}</label>
                </div>
            {/each}
        </Field>

        <Field class="form-field required m-0" name="email" let:uniqueId>
            <label for={uniqueId}>To email address</label>
            <!-- svelte-ignore a11y-autofocus -->
            <input type="email" id={uniqueId} autofocus required bind:value={email} />
        </Field>
    </form>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide} disabled={isSubmitting}
            >Close</button
        >
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSubmitting}
            disabled={!canSubmit || isSubmitting}
        >
            <i class="ri-mail-send-line" />
            <span class="txt">Send</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
