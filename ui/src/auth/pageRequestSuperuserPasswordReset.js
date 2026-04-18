export function pageRequestSuperuserPasswordReset(route) {
    app.store.title = "Forgotten superuser password";

    const data = store({
        email: "",
        isSubmitting: false,
        success: false,
    });

    async function submit() {
        if (data.isSubmitting) {
            return;
        }

        data.isSubmitting = true;

        try {
            await app.pb.collection("_superusers").requestPasswordReset(data.email);
            data.success = true;
        } catch (err) {
            app.checkApiError(err);
        }

        data.isSubmitting = false;
    }

    return t.div(
        {
            pbEvent: "pageSuperuserPasswordReset",
            className: "wrapper sm m-auto p-b-base",
        },
        t.header(
            { className: "txt-center m-b-base" },
            t.img({ className: "main-logo", src: () => app.store.mainLogo, ariaHidden: true, alt: "App logo" }),
            t.h5({ className: "m-t-10" }, () => app.store.title),
        ),
        () => {
            if (data.success) {
                return t.div(
                    { pbEvent: "superuserPasswordResetAlert", className: "alert success txt-center" },
                    t.p(null, "Check ", t.strong(null, data.email), " for the recovery link!"),
                );
            }

            return t.form(
                {
                    pbEvent: "superuserPasswordResetForm",
                    className: "grid request-password-reset-form",
                    onsubmit: (e) => {
                        e.preventDefault();
                        submit();
                    },
                },
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "content txt-center m-b-sm" },
                        t.p(null, "Enter the email associated with your account and we'll send you a recovery link:"),
                    ),
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: "password_reset_email" }, "Email"),
                        t.input({
                            id: "password_reset_email",
                            name: "email",
                            type: "email",
                            required: true,
                            autofocus: true,
                            value: () => data.email,
                            oninput: (e) => (data.email = e.target.value),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-12" },
                    t.button(
                        {
                            className: () => `btn lg block ${data.isSubmitting ? "loading" : ""}`,
                            disabled: () => data.isSubmitting,
                        },
                        t.i({ className: "ri-mail-send-line" }),
                        t.span({ className: "txt" }, "Send recovery link"),
                    ),
                ),
            );
        },
        t.div(
            { className: "block m-t-sm txt-center" },
            t.a({ href: "#/login", className: "link-hint" }, "Back to login"),
        ),
    );
}
