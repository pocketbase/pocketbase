<script>
    import Field from "@/components/base/Field.svelte";
    import EmailTemplateAccordion from "@/components/collections/EmailTemplateAccordion.svelte";
    import TokenOptionsAccordion from "@/components/collections/TokenOptionsAccordion.svelte";
    import MFAAccordion from "@/components/collections/MFAAccordion.svelte";
    import OAuth2Accordion from "@/components/collections/OAuth2Accordion.svelte";
    import OTPAccordion from "@/components/collections/OTPAccordion.svelte";
    import PasswordAuthAccordion from "@/components/collections/PasswordAuthAccordion.svelte";
    import EmailTestPopup from "@/components/settings/EmailTestPopup.svelte";

    export let collection;

    let emailTemplatesList = [];
    let emailTestPopup;

    $: isSuperusers = collection.system && collection.name === "_superusers";

    // nested email template normalizations
    $: if (typeof collection.otp?.emailTemplate == "undefined") {
        collection.otp = collection.otp || {};
        collection.otp.emailTemplate = {};
    }
    $: if (typeof collection.authAlert?.emailTemplate == "undefined") {
        collection.authAlert = collection.authAlert || {};
        collection.authAlert.emailTemplate = {};
    }

    // predefined email template configs
    $: resetPasswordTemplate = {
        key: "resetPasswordTemplate",
        label: "Default Password reset email template",
        placeholders: ["APP_NAME", "APP_URL", "RECORD:*", "TOKEN"],
        config: collection.resetPasswordTemplate,
    };
    $: verificationTemplate = {
        key: "verificationTemplate",
        label: "Default Verification email template",
        placeholders: ["APP_NAME", "APP_URL", "RECORD:*", "TOKEN"],
        config: collection.verificationTemplate,
    };
    $: confirmEmailChangeTemplate = {
        key: "confirmEmailChangeTemplate",
        label: "Default Confirm email change email template",
        placeholders: ["APP_NAME", "APP_URL", "RECORD:*", "TOKEN"],
        config: collection.confirmEmailChangeTemplate,
    };
    $: otpTemplate = {
        key: "otp.emailTemplate",
        label: "Default OTP email template",
        placeholders: ["APP_NAME", "APP_URL", "RECORD:*", "OTP", "OTP_ID"],
        config: collection.otp.emailTemplate,
    };
    $: authAlertTemplate = {
        key: "authAlert.emailTemplate",
        label: "Default Login alert email template",
        placeholders: ["APP_NAME", "APP_URL", "RECORD:*", "ALERT_INFO"],
        config: collection.authAlert.emailTemplate,
    };
    $: emailTemplatesList = isSuperusers
        ? [resetPasswordTemplate, otpTemplate, authAlertTemplate]
        : [
              verificationTemplate,
              resetPasswordTemplate,
              confirmEmailChangeTemplate,
              otpTemplate,
              authAlertTemplate,
          ];
</script>

<h4 class="section-title">
    <div class="flex">
        <span class="txt">Auth methods</span>
        <div class="m-l-auto handle">
            <Field
                class="form-field form-field-sm form-field-toggle m-0"
                name="authAlert.enabled"
                inlineError={true}
                let:uniqueId
            >
                <input type="checkbox" id={uniqueId} bind:checked={collection.authAlert.enabled} />
                <label for={uniqueId}>Send email alert for new logins</label>
            </Field>
        </div>
    </div>
</h4>
<div class="accordions m-b-35">
    <PasswordAuthAccordion bind:collection />

    {#if !isSuperusers}
        <OAuth2Accordion bind:collection />
    {/if}

    <OTPAccordion bind:collection />

    <MFAAccordion bind:collection />
</div>

<h4 class="section-title">
    <span class="txt">Mail templates</span>
    <button
        type="button"
        class="btn btn-xs m-l-auto btn-secondary"
        on:click={() => emailTestPopup?.show(collection.id)}
    >
        Send test email
    </button>
</h4>
<div class="accordions m-b-35">
    <div class="accordions">
        {#each emailTemplatesList as template (template.key)}
            <EmailTemplateAccordion
                single
                key={template.key}
                title={template.label}
                placeholders={template?.placeholders}
                bind:config={template.config}
            />
        {/each}
    </div>
</div>

<h4 class="section-title">Other</h4>
<div class="accordions m-b-base">
    <TokenOptionsAccordion bind:collection />
</div>

<EmailTestPopup bind:this={emailTestPopup} />
