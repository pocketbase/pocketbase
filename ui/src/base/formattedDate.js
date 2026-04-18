window.app = window.app || {};
window.app.components = window.app.components || {};

const tzName = Intl.DateTimeFormat().resolvedOptions().timeZone;

window.app.components.formattedDate = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        value: "",
        short: false,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            ariaDescription: app.attrs.tooltip(() => {
                if (props.short && !!props.value) {
                    return app.utils.toLocalDatetime(props.value) + "\n" + tzName;
                }

                return null;
            }),
            "html-class": "formatted-date",
            className: () => `formatted-date ${props.short ? "short" : "full"}`,
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        () => {
            if (!props.value) {
                return t.span({ className: "missing-value" });
            }

            if (props.short) {
                const parts = props.value.split(" ");
                return [
                    t.span({ className: "primary-date" }, parts[0]),
                    t.span({ className: "secondary-date" }, parts[1]),
                ];
            }

            return [
                t.span({ className: "primary-date" }, app.utils.toLocalDatetime(props.value)),
                t.span({ className: "secondary-date" }, props.value),
            ];
        },
    );
};
