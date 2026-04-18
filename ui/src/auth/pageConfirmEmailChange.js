import PocketBase, { getTokenPayload } from "pocketbase";

export function pageConfirmEmailChange(route) {
    const token = route.params?.token || "";
    const tokenPayload = getTokenPayload(token);

    if (!tokenPayload.newEmail || !tokenPayload.collectionId) {
        app.toasts.error("Invalid or expired email change token.");
        window.location.hash = "#/";
        return;
    }

    app.store.title = "Confirm email change";

    const data = store({
        password: "",
        isSubmitting: false,
        isSuccess: false,
        showPassword: false,
    });

    async function submit() {
        if (data.isSubmitting) {
            return;
        }

        data.isSubmitting = true;

        // init a custom client to avoid interfering with the superuser state
        const client = new PocketBase(app.pb.baseURL);

        try {
            await client.collection(tokenPayload.collectionId).confirmEmailChange(token, data.password);
            data.isSuccess = true;
        } catch (err) {
            app.checkApiError(err);
            data.isSuccess = false;
        }

        data.isSubmitting = false;
    }

    return t.div(
        {
            pbEvent: "pageConfirmEmailChange",
            className: "wrapper sm m-auto p-b-base",
        },
        t.header(
            { className: "txt-center m-b-base" },
            t.img({ className: "main-logo", src: () => app.store.mainLogo, ariaHidden: true, alt: "App logo" }),
            t.h5({ className: "m-t-10" }, () => app.store.title),
        ),
        () => {
            if (data.isSuccess) {
                return t.div(
                    {
                        pbEvent: "confirmEmailChangeAlert",
                        className: "alert success txt-center",
                    },
                    t.p(null, "The email was successfully changed."),
                    t.p(null, "You can go back and sign in with your new email address."),
                );
            }

            return t.form(
                {
                    pbEvent: "confirmEmailChangeForm",
                    className: "grid confirm-email-change-form",
                    onsubmit: (e) => {
                        e.preventDefault();
                        submit();
                    },
                },
                t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "content txt-center m-b-sm" },
                        "Type your password to confirm changing your email address to ",
                        t.strong(null, tokenPayload.newEmail),
                        ":",
                    ),
                    t.div(
                        { className: "fields" },
                        t.div(
                            { className: "field" },
                            t.label({ htmlFor: "password_confirm" }, "Password"),
                            t.input({
                                id: "password_confirm",
                                name: "password",
                                required: true,
                                autofocus: true,
                                type: () => (data.showPassword ? "text" : "password"),
                                value: () => data.password,
                                oninput: (e) => (data.password = e.target.value),
                            }),
                        ),
                        t.div(
                            { className: "field addon" },
                            t.button(
                                {
                                    type: "button",
                                    tabIndex: -1,
                                    className: "btn sm transparent secondary circle tooltip-right",
                                    ariaDescription: app.attrs.tooltip(() =>
                                        data.showPassword ? "Hide password" : "Show password"
                                    ),
                                    onclick: () => (data.showPassword = !data.showPassword),
                                },
                                t.i({
                                    className: () => (data.showPassword ? "ri-eye-off-line" : "ri-eye-line"),
                                }),
                            ),
                        ),
                    ),
                ),
                t.div(
                    { className: "col-12" },
                    t.button(
                        {
                            className: () => `btn lg block ${data.isSubmitting ? "loading" : ""}`,
                            disabled: () => data.isSubmitting,
                        },
                        t.span({ className: "txt" }, "Confirm new email"),
                    ),
                ),
            );
        },
    );
}
