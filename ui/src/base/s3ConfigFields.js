window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Returns a wrapper form element with common S3 config fields.
 *
 * @example
 * ```js
 * app.components.s3ConfigFields({
 *     config: () => data.settings.storage,
 * })
 * ```
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.s3ConfigFields = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        config: {}, // S3 config store (pass as a function in case the object is being replaced)
        configKey: "s3", // used for the fields error matching
        toggleLabel: "Use S3 storage",
        testFilesystem: "storage",
        before: null,
        after: null,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    if (props.configKey.endsWith(".")) {
        props.configKey = props.configKey.substring(0, props.configKey.length - 1);
    }

    const data = store({
        originalHash: "",
        originalConfig: null,
    });

    watchers.push(
        watch(
            () => props.config,
            (c) => {
                data.originalHash = JSON.stringify(c);
                data.originalConfig = JSON.parse(data.originalHash);
            },
        ),
    );

    return t.div(
        {
            pbEvent: "s3ConfigFields",
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `block s3-fields s3-config-${props.configKey} ${props.className}`,
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.div(
            { className: "field" },
            t.input({
                id: () => `${props.configKey}.enabled`,
                name: () => `${props.configKey}.enabled`,
                type: "checkbox",
                className: "switch",
                checked: () => props.config.enabled,
                onchange: (e) => props.config.enabled = e.target.checked,
            }),
            t.label({ htmlFor: () => `${props.configKey}.enabled` }, () => props.toggleLabel),
        ),
        (el) => {
            if (typeof props.before == "function") {
                return props.before(el);
            }
            return props.before;
        },
        app.components.slide(
            () => props.config.enabled,
            t.div(
                { className: "grid m-t-base" },
                t.div(
                    { className: "col-lg-6" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: () => `${props.configKey}.endpoint` }, "Endpoint"),
                        t.input({
                            id: () => `${props.configKey}.endpoint`,
                            name: () => `${props.configKey}.endpoint`,
                            type: "text",
                            required: () => props.config.enabled,
                            value: () => props.config.endpoint || "",
                            oninput: (e) => (props.config.endpoint = e.target.value),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-lg-3" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: () => `${props.configKey}.bucket` }, "Bucket"),
                        t.input({
                            id: () => `${props.configKey}.bucket`,
                            name: () => `${props.configKey}.bucket`,
                            type: "text",
                            required: () => props.config.enabled,
                            value: () => props.config.bucket || "",
                            oninput: (e) => (props.config.bucket = e.target.value),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-lg-3" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: () => `${props.configKey}.region` }, "Region"),
                        t.input({
                            id: () => `${props.configKey}.region`,
                            name: () => `${props.configKey}.region`,
                            type: "text",
                            required: () => props.config.enabled,
                            value: () => props.config.region || "",
                            oninput: (e) => (props.config.region = e.target.value),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-lg-6" },
                    t.div(
                        { className: "field" },
                        t.label({ htmlFor: () => `${props.configKey}.accessKey` }, "Access key"),
                        t.input({
                            id: () => `${props.configKey}.accessKey`,
                            name: () => `${props.configKey}.accessKey`,
                            type: "text",
                            autocomplete: "off",
                            required: () => props.config.enabled,
                            value: () => props.config.accessKey || "",
                            oninput: (e) => (props.config.accessKey = e.target.value),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-lg-6" },
                    t.div(
                        {
                            className: () => `field ${props.config.enabled ? "" : "required"}`,
                        },
                        t.label({ htmlFor: () => `${props.configKey}.secret` }, "Secret"),
                        t.input({
                            id: () => `${props.configKey}.secret`,
                            name: () => `${props.configKey}.secret`,
                            type: "password",
                            autocomplete: "new-password",
                            value: () => props.config.secret || "",
                            oninput: (e) => (props.config.secret = e.target.value),
                            onkeyup: (e) => {
                                if (
                                    e.key == "Backspace"
                                    && typeof props.config.secret === "undefined"
                                ) {
                                    props.config.secret = "";
                                }
                            },
                            placeholder: () => (typeof props.config.secret !== "undefined" ? "" : "* * * * * *"),
                        }),
                    ),
                ),
                t.div(
                    { className: "col-lg-6", style: "min-height: 25px" },
                    t.div(
                        { className: "field" },
                        t.input({
                            id: () => `${props.configKey}.forcePathStyle`,
                            name: () => `${props.configKey}.forcePathStyle`,
                            type: "checkbox",
                            checked: () => props.config.forcePathStyle || false,
                            onchange: (e) => (props.config.forcePathStyle = e.target.checked),
                        }),
                        t.label(
                            { htmlFor: () => `${props.configKey}.forcePathStyle` },
                            t.span({ className: "txt" }, "Force path-style addressing"),
                            t.i({
                                className: "ri-information-line link-hint",
                                ariaDescription: app.attrs.tooltip(
                                    `Forces the request to use path-style addressing, eg. "https://s3.amazonaws.com/BUCKET/KEY" instead of the default "https://BUCKET.s3.amazonaws.com/KEY".`,
                                ),
                            }),
                        ),
                    ),
                ),
                t.div({ className: "col-lg-6 txt-right" }, () => {
                    if (!props.config?.enabled || data.originalHash != JSON.stringify(props.config)) {
                        return;
                    }

                    return app.components.s3Test({
                        config: () => props.config,
                        testFilesystem: () => props.testFilesystem,
                    });
                }),
            ),
        ),
        (el) => {
            if (typeof props.after == "function") {
                return props.after(el);
            }
            return props.after;
        },
    );
};
