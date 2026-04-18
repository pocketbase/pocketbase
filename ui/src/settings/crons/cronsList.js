export function cronsList(propsArg = {}) {
    const props = store({
        reset: null,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        isLoading: false,
        isRunning: {},
        crons: [],
    });

    async function loadCrons() {
        data.isLoading = true;

        try {
            data.crons = await app.pb.crons.getFullList();
            data.isLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                data.isLoading = false;
            }
        }
    }

    async function runCron(jobId) {
        if (!jobId || data.isRunning[jobId]) {
            return;
        }

        data.isRunning[jobId] = true;

        try {
            await app.pb.crons.run(jobId);
            app.toasts.success(`Successfully triggered "${jobId}".`);
            data.isRunning[jobId] = false;
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
                data.isRunning[jobId] = false;
            }
        }
    }

    return t.div(
        {
            pbEvent: "cronsList",
            className: "list",
            onmount: () => {
                watchers.push(
                    watch(() => props.reset, () => {
                        loadCrons();
                    }),
                );
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        () => {
            if (!data.isLoading || data.crons.length) {
                return;
            }

            const skeletons = [];
            for (let i = 0; i < 4; i++) {
                skeletons.push(
                    t.div({ rid: "skeleton_" + i, className: "list-item" }, t.div({ className: "skeleton-loader" })),
                );
            }
            return skeletons;
        },
        t.div(
            {
                hidden: () => data.isLoading || data.crons.length,
                className: "list-item",
            },
            t.div({ className: "content block txt-hint" }, "No registered crons found."),
        ),
        () => {
            return data.crons.map((cron) => {
                return t.div(
                    { className: () => `list-item ${data.isLoading ? "faded" : ""}` },
                    t.div(
                        { className: "content" },
                        t.span({
                            className: "cron-id txt-code txt-ellipsis",
                            title: () => cron.id,
                            textContent: () => cron.id,
                        }),
                    ),
                    t.small({ className: "cron-expression txt-hint txt-nowrap txt-code" }, () => cron.expression),
                    t.nav(
                        { hidden: () => data.isLoading, className: "actions" },
                        t.button(
                            {
                                type: "button",
                                ariaDescription: app.attrs.tooltip("Run"),
                                className: () =>
                                    `btn sm circle secondary transparent ${data.isRunning[cron.id] ? "loading" : ""}`,
                                disabled: () => data.isRunning[cron.id],
                                onclick: () => runCron(cron.id),
                            },
                            t.i({ className: "ri-play-large-line" }),
                        ),
                    ),
                );
            });
        },
    );
}
