import { emailTemplateAccordion } from "./emailTemplateAccordion";
import { mfaAccordion } from "./mfaAccordion";
import { oauth2Accordion } from "./oauth2Accordion";
import { otpAccordion } from "./otpAccordion";
import { passwordAuthAccordion } from "./passwordAuthAccordion";
import { tokenOptionsAccordion } from "./tokenOptionsAccordion";

export function collectionAuthOptionsTab(upsertData) {
    const uniqueId = "options_" + app.utils.randomString();

    return t.div(
        { className: "collection-tab-content collection-options-tab-content" },
        t.div(
            { className: "grid" },
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "section-heading" },
                    t.strong(null, "Auth methods"),
                    t.div({ className: "flex-fill" }),
                    t.div(
                        { className: "field" },
                        t.input({
                            id: uniqueId + ".authAlert",
                            name: "authAlert.enabled",
                            type: "checkbox",
                            className: "switch sm",
                            checked: () => !!upsertData.collection.authAlert?.enabled,
                            onchange: (e) => {
                                upsertData.collection.authAlert = upsertData.collection.authAlert || {};
                                upsertData.collection.authAlert.enabled = e.target.checked;
                            },
                        }),
                        t.label({ htmlFor: uniqueId + ".authAlert" }, "Send email alert for new logins"),
                    ),
                ),
                passwordAuthAccordion(upsertData.collection),
                () => {
                    if (upsertData.originalCollection?.name == "_superusers") {
                        return;
                    }

                    return oauth2Accordion(upsertData.collection);
                },
                otpAccordion(upsertData.collection),
                mfaAccordion(upsertData.collection),
            ),
            t.div(
                { className: "col-12" },
                t.div(
                    { className: "section-heading" },
                    t.strong(null, "Mail templates"),
                    t.button({
                        tabIndex: -1,
                        type: "buttton",
                        className: "m-l-auto label handle txt-bold",
                        textContent: "Send test email",
                        onclick: () => app.modals.openMailTest(upsertData.collection?.name),
                    }),
                ),
                emailTemplateAccordion(upsertData.collection, "verificationTemplate", {
                    title: "Default Verification email template",
                    placeholders: ["{APP_NAME}", "{APP_URL}", "{RECORD:*}", "{TOKEN}"],
                }),
                emailTemplateAccordion(upsertData.collection, "resetPasswordTemplate", {
                    title: "Default Password reset email template",
                    placeholders: ["{APP_NAME}", "{APP_URL}", "{RECORD:*}", "{TOKEN}"],
                }),
                emailTemplateAccordion(upsertData.collection, "confirmEmailChangeTemplate", {
                    title: "Default Confirm email change email template",
                    placeholders: ["{APP_NAME}", "{APP_URL}", "{RECORD:*}", "{TOKEN}"],
                }),
                emailTemplateAccordion(upsertData.collection, "otp.emailTemplate", {
                    title: "Default OTP email template",
                    placeholders: ["{APP_NAME}", "{APP_URL}", "{RECORD:*}", "{OTP}", "{OTP_ID}"],
                }),
                emailTemplateAccordion(upsertData.collection, "authAlert.emailTemplate", {
                    title: "Default Login alert email template",
                    placeholders: ["{APP_NAME}", "{APP_URL}", "{RECORD:*}", "{ALERT_INFO}"],
                }),
            ),
            t.div(
                { className: "col-12" },
                t.div({ className: "section-heading" }, t.strong(null, "Other")),
                tokenOptionsAccordion(upsertData.collection),
            ),
        ),
    );
}
