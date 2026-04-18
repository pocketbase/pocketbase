export function pageSuperuserLogin(route) {
    app.store.title = "Superuser login";

    const data = store({
        authMethods: {},

        identity: route.query.demoEmail?.[0] || "",
        password: route.query.demoPassword?.[0] || "",
        showPassword: false,

        otpId: "",
        lastOTPId: "",
        otpEmail: "",
        otpPassword: "",

        mfaId: "",
        totalSteps: 1,
        get currentStep() {
            return 1 + !!data.mfaId + !!data.otpId;
        },

        isAuthMethodsLoading: true,
        isPasswordAuthSubmitting: false,
        isOTPRequestSubmitting: false,
        isOTPAuthSubmitting: false,
    });

    async function loadAuthMethods() {
        data.isAuthMethodsLoading = true;

        try {
            data.authMethods = await app.pb.collection("_superusers").listAuthMethods();

            data.totalSteps = 1;
            data.mfaId = "";
            data.otpId = "";

            if (data.authMethods.mfa?.enabled && data.authMethods.otp?.enabled) {
                data.totalSteps += 2; // otp request + auth
            }
        } catch (err) {
            app.checkApiError(err);
        }

        data.isAuthMethodsLoading = false;
    }

    loadAuthMethods();

    return t.div(
        {
            pbEvent: "pageSuperuserLogin",
            className: "wrapper sm m-auto p-b-base",
        },
        t.header(
            { className: "txt-center m-b-base" },
            t.img({ className: "main-logo", src: () => app.store.mainLogo, ariaHidden: true, alt: "App logo" }),
            t.h5(
                { className: "m-t-10" },
                t.span(null, () => app.store.title),
                () => {
                    if (data.totalSteps > 1) {
                        return t.span(null, () => ` (${data.currentStep}/${data.totalSteps})`);
                    }
                },
            ),
        ),
        () => {
            if (data.isAuthMethodsLoading) {
                return t.div({ className: "block txt-center" }, t.span({ className: "loader lg" }));
            }

            if (data.authMethods.password?.enabled && !data.mfaId) {
                return authWithPasswordForm(data);
            }

            if (data.authMethods.otp?.enabled) {
                if (!data.otpId) {
                    return requestOTPForm(data);
                }

                return authWithOTPForm(data);
            }
        },
    );
}

// Auth with password
// -------------------------------------------------------------------

async function authWithPassword(data) {
    if (data.isPasswordAuthSubmitting) {
        return;
    }

    data.isPasswordAuthSubmitting = true;

    try {
        await app.pb.collection("_superusers").authWithPassword(data.identity, data.password);
        app.toasts.removeAll();
        app.store.errors = null;
        app.utils.toRememberedPath();
    } catch (err) {
        if (err.status == 401) {
            data.mfaId = err.response.mfaId;

            // show the otp forms
            if (
                // if the identity field is just the email use it to directly send an otp request
                data.authMethods?.password?.identityFields?.length == 1
                && data.authMethods.password.identityFields[0] == "email"
            ) {
                // prefill and request
                data.otpEmail = data.identity;
                await requestOTP(data);
            } else if (/^[^@\s]+@[^@\s]+$/.test(data.identity)) {
                // only prefill
                data.otpEmail = data.identity;
            }
        } else if (err.status != 400) {
            app.checkApiError(err);
        } else {
            app.toasts.error("Invalid login credentials.");
        }
    }

    data.isPasswordAuthSubmitting = false;
}

function authWithPasswordForm(data) {
    return t.form(
        {
            pbEvent: "authWithPasswordForm",
            className: "grid auth-with-password-form",
            onsubmit: (e) => {
                e.preventDefault();
                authWithPassword(data);
            },
        },
        t.div(
            { className: "col-12" },
            t.div(
                { className: "field" },
                t.label(
                    { htmlFor: "login_identity" },
                    () => app.utils.sentenize(data.authMethods.password.identityFields.join(" or "), false),
                ),
                t.input({
                    id: "login_identity",
                    name: "identity",
                    type: () => {
                        if (
                            data.authMethods.password.identityFields.length == 1
                            && data.authMethods.password.identityFields[0] == "email"
                        ) {
                            return "email";
                        }
                        return "text";
                    },
                    required: true,
                    autofocus: true,
                    value: () => data.identity,
                    oninput: (e) => (data.identity = e.target.value),
                }),
            ),
        ),
        t.div(
            { className: "col-12" },
            t.div(
                { className: "fields" },
                t.div(
                    { className: "field" },
                    t.label({ htmlFor: "login_pass" }, "Password"),
                    t.input({
                        id: "login_pass",
                        name: "password",
                        required: true,
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
            t.a(
                {
                    href: "#/request-password-reset",
                    className: "link-hint m-t-5",
                    onclick: (e) => e.stopPropagation(),
                },
                t.small(null, "Forgotten password"),
            ),
        ),
        t.div(
            { className: "col-12" },
            t.button(
                {
                    className: () => `btn lg block next ${data.isPasswordAuthSubmitting ? "loading" : ""}`,
                    disabled: () => data.isPasswordAuthSubmitting,
                },
                t.span({ className: "txt" }, () => (data.totalSteps > 1 ? "Next" : "Login")),
                t.i({ className: "ri-arrow-right-line" }),
            ),
        ),
    );
}

// Request OTP
// -------------------------------------------------------------------

async function requestOTP(data) {
    if (data.isOTPRequestSubmitting) {
        return;
    }

    data.isOTPRequestSubmitting = true;

    try {
        const result = await app.pb.collection("_superusers").requestOTP(data.otpEmail);
        data.otpId = result.otpId;
        data.lastOTPId = data.otpId;
        app.toasts.removeAll();
        app.store.errors = null;
    } catch (err) {
        // reset the form
        if (err.status == 429) {
            data.otpId = data.lastOTPId;
        }

        app.checkApiError(err);
    }

    data.isOTPRequestSubmitting = false;
}

function requestOTPForm(data) {
    return t.form(
        {
            pbEvent: "requestOTPForm",
            className: "grid request-otp-form",
            onsubmit: (e) => {
                e.preventDefault();
                requestOTP(data);
            },
        },
        t.div(
            { className: "col-12" },
            t.div(
                { className: "field" },
                t.label({ htmlFor: "otp_email" }, "Email"),
                t.input({
                    id: "otp_email",
                    name: "email",
                    type: "email",
                    required: true,
                    autofocus: true,
                    value: () => data.otpEmail,
                    oninput: (e) => (data.otpEmail = e.target.value),
                }),
            ),
        ),
        t.div(
            { className: "col-12" },
            t.button(
                {
                    className: () => `btn lg block ${data.isOTPRequestSubmitting ? "loading" : ""}`,
                    disabled: () => data.isOTPRequestSubmitting,
                },
                t.i({ className: "ri-mail-send-line" }),
                t.span({ className: "txt" }, "Send OTP"),
            ),
        ),
    );
}

// Auth with OTP
// -------------------------------------------------------------------

async function authWithOTP(data) {
    if (data.isOTPAuthSubmitting) {
        return;
    }

    data.isOTPAuthSubmitting = true;

    try {
        await app.pb.collection("_superusers").authWithOTP(data.otpId || data.lastOTPId, data.otpPassword, {
            mfaId: data.mfaId,
        });
        app.toasts.removeAll();
        app.store.errors = null;
        app.utils.toRememberedPath();
    } catch (err) {
        app.checkApiError(err);
    }

    data.isOTPAuthSubmitting = false;
}

function authWithOTPForm(data) {
    return t.form(
        {
            pbEvent: "authWithOTPForm",
            className: "grid auht-with-otp-form",
            onsubmit: (e) => {
                e.preventDefault();
                authWithOTP(data);
            },
        },
        () => {
            if (data.otpEmail) {
                return t.div(
                    { className: "col-12" },
                    t.div(
                        { className: "content txt-center" },
                        "Check your ",
                        t.strong(null, data.otpEmail),
                        " inbox and enter below the received One-time password (OTP).",
                    ),
                );
            }
        },
        t.div(
            { className: "col-12" },
            t.div(
                { className: "field" },
                t.label({ htmlFor: "otp_id" }, "Id"),
                t.input({
                    id: "otp_id",
                    name: "otpId",
                    type: "text",
                    required: true,
                    placeholder: () => data.lastOTPId,
                    value: () => data.otpId,
                    onchange: (e) => {
                        data.otpId = e.target.value || data.lastOTPId;
                        e.target.value = data.otpId;
                    },
                }),
            ),
        ),
        t.div(
            { className: "col-12" },
            t.div(
                { className: "fields" },
                t.div(
                    { className: "field" },
                    t.label({ htmlFor: "otp_password" }, "One-time password"),
                    t.input({
                        id: "otp_password",
                        name: "password",
                        required: true,
                        type: () => (data.showPassword ? "text" : "password"),
                        value: () => data.otpPassword,
                        oninput: (e) => (data.otpPassword = e.target.value),
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
                    className: () => `btn lg block next ${data.isOTPAuthSubmitting ? "loading" : ""}`,
                    disabled: () => data.isOTPAuthSubmitting,
                },
                t.span({ className: "txt" }, "Login"),
                t.i({ className: "ri-arrow-right-line" }),
            ),
            t.div(
                { className: "block m-t-sm txt-center" },
                t.button(
                    {
                        type: "button",
                        className: "link-hint txt-sm",
                        disabled: () => data.isOTPAuthSubmitting,
                        onclick: () => {
                            data.otpId = "";
                            data.otpPassword = "";
                        },
                    },
                    "Request another OTP",
                ),
            ),
        ),
    );
}
