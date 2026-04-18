window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openMailTest = function(preselectedCollectionIdOrName = "", template = "") {
    const modal = mailTestModal(preselectedCollectionIdOrName, template);

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function mailTestModal(preselectedCollectionIdOrName = "", template = "") {
    const uniqueId = "mail_test_" + app.utils.randomString();
    const emailStorageKey = "pbLastTestEmail";
    const testRequestKey = "email_test_request";

    const templateOptions = [
        { label: "Verification", value: "verification" },
        { label: "Password reset", value: "password-reset" },
        { label: "Confirm email change", value: "email-change" },
        { label: "OTP", value: "otp" },
        { label: "Login alert", value: "login-alert" },
    ];

    const data = store({
        email: localStorage.getItem(emailStorageKey) || app.store.superuser?.email || "",
        template: template || templateOptions[0].value,
        isSending: false,
        collectionIdOrName: preselectedCollectionIdOrName,
        get isAuthCollectionsLoading() {
            return app.store.isCollectionsLoading;
        },
        get authCollections() {
            return app.utils.sortedCollections(
                app.store.collections.filter((c) => c.type == "auth"),
            );
        },
        get canSubmit() {
            return !!data.email && !!data.template && !!data.collectionIdOrName;
        },
    });

    let testTimeoutId;

    async function send() {
        if (data.isSending || !data.canSubmit) {
            return;
        }

        data.isSending = true;

        // auto cancel the test request after 15sec
        clearTimeout(testTimeoutId);
        testTimeoutId = setTimeout(() => {
            data.isSending = false;
            app.pb.cancelRequest(testRequestKey);
            app.modals.close();
            app.toasts.error("Test email send timeout.");
        }, 15000);

        try {
            // store as preset
            if (data.email != app.pb.authStore.record?.email) {
                localStorage.setItem(emailStorageKey, data.email);
            }

            await app.pb.settings.testEmail(data.collectionIdOrName, data.email, data.template, {
                requestKey: testRequestKey,
            });

            app.toasts.success("Successfully sent test email.");

            app.modals.close();
        } catch (err) {
            app.checkApiError(err);
        }

        data.isSending = false;
        clearTimeout(testTimeoutId);
    }

    const watchers = [];

    return t.div(
        {
            className: "modal popup sm",
            onbeforeopen: (el) => {
                // preselect the first auth collection as fallback
                watchers.push(
                    watch(() => data.isAuthCollectionsLoading, (isLoading) => {
                        if (!isLoading && !data.collectionIdOrName) {
                            data.collectionIdOrName = data.authCollections[0]?.id || "";
                        }
                    }),
                );
            },
            onafterclose: (el) => {
                clearTimeout(testTimeoutId);
                el?.remove();
            },
            onunmount: () => {
                clearTimeout(testTimeoutId);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.header({ className: "modal-header" }, t.h5({ className: "m-auto" }, "Send test email")),
        t.form(
            {
                id: uniqueId,
                className: "modal-content mail-settings-test-form",
                onsubmit: (e) => {
                    e.preventDefault();
                    send();
                },
            },
            t.div(
                { className: "grid" },
                t.div({ className: "col-lg-12" }, () => {
                    return templateOptions.map((opt, i) => {
                        return t.field(
                            { className: () => `field ${i > 0 ? "m-t-10" : ""}` },
                            t.input({
                                type: "radio",
                                id: uniqueId + ".template." + opt.value,
                                name: "template",
                                checked: () => data.template == opt.value,
                                onchange: (e) => (data.template = opt.value),
                            }),
                            t.label({ htmlFor: uniqueId + ".template." + opt.value }, opt.label || opt.value),
                        );
                    });
                }),
                () => {
                    if (preselectedCollectionIdOrName) {
                        return;
                    }

                    return t.div(
                        { className: "col-lg-12" },
                        t.div(
                            { className: "field" },
                            t.label({ htmlFor: uniqueId + ".collection" }, "Auth collection"),
                            app.components.select({
                                id: uniqueId + ".collection",
                                name: "collection",
                                required: true,
                                placeholder: () =>
                                    data.isAuthCollectionsLoading
                                        ? "Loading auth collections..."
                                        : "Select auth collection",
                                options: () =>
                                    data.authCollections.map((c) => {
                                        return { value: c.id, label: c.name };
                                    }),
                                value: () => data.collectionIdOrName || "",
                                onchange: (selected) => {
                                    data.collectionIdOrName = selected?.[0]?.value;
                                },
                            }),
                        ),
                    );
                },
                t.div(
                    { className: "col-lg-12" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: uniqueId + ".email" }, "To email address"),
                        t.input({
                            id: uniqueId + ".email",
                            name: "email",
                            type: "email",
                            required: true,
                            value: () => data.email || "",
                            oninput: (e) => (data.email = e.target.value),
                        }),
                    ),
                ),
            ),
        ),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    onclick: () => app.modals.close(),
                    disabled: () => data.isSending,
                },
                t.span({ className: "txt" }, "Close"),
            ),
            t.button(
                {
                    "html-form": uniqueId,
                    type: "submit",
                    className: () => `btn expanded ${data.isSending ? "loading" : ""}`,
                    disabled: () => data.isSending || !data.canSubmit,
                },
                t.i({ className: "ri-mail-send-line", ariaHidden: true }),
                t.span({ className: "txt" }, "Send"),
            ),
        ),
    );
}
