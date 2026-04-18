window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * A helper label element that performs live S3 connectivity tests.
 *
 * ```js
 * app.components.s3Test({
 *     config: () => data.settings.backups.s3,
 *     testFilesystem: "backups",
 * })
 * ```
 *
 * @param  {Object} propsArg
 * @return {Element}
 */
window.app.components.s3Test = function(propsArg = {}) {
    const testRequestKey = "s3_test_request";

    const props = store({
        rid: undefined,
        config: null, // S3 config store
        label: "Use S3 storage",
        testFilesystem: "storage", // "storage" or "backups"
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        isTesting: false,
        testError: null,
        get hasError() {
            return !app.utils.isEmpty(data.testError);
        },
    });

    let testDebounceId;
    let testTimeoutId;

    function testS3WithDebounce(timeout = 150) {
        if (!props.config.enabled) {
            clearTimeout(testDebounceId);
            return;
        }

        data.isTesting = true;

        clearTimeout(testDebounceId);
        testDebounceId = setTimeout(() => {
            testS3();
        }, timeout);
    }

    async function testS3() {
        data.isTesting = true;

        if (!props.config.enabled || !props.testFilesystem) {
            data.testError = null;
            data.isTesting = false;
            return; // nothing to test
        }

        // auto cancel the test request after 30sec
        app.pb.cancelRequest(testRequestKey);
        clearTimeout(testTimeoutId);
        testTimeoutId = setTimeout(() => {
            app.pb.cancelRequest(testRequestKey);
            data.testError = new Error("S3 test connection timeout.");
            data.isTesting = false;
        }, 30000);

        try {
            await app.pb.props.testS3(props.testFilesystem, {
                requestKey: testRequestKey,
            });
            data.testError = null;
            data.isTesting = false;
        } catch (err) {
            if (!err?.isAbort) {
                data.testError = err;
                data.isTesting = false;
                clearTimeout(testTimeoutId);
            }
        }
    }

    watchers.push(
        watch(
            () => props.testFilesystem && props.config,
            () => testS3WithDebounce(),
        ),
    );

    return t.div(
        {
            pbEvent: "s3Test",
            rid: props.rid,
            hidden: () => !props.testFilesystem,
            className: () => `label s3-test-label txt-nowrap ${data.hasError ? "warning" : "success"}`,
            ariaDescription: app.attrs.tooltip(() => data.testError?.data?.message),
            onunmount: () => {
                clearTimeout(testTimeoutId);
                clearTimeout(testDebounceId);
                watchers.forEach((w) => w?.unwatch());
            },
        },
        () => {
            if (data.isTesting) {
                return t.span({ className: "loader sm" });
            }

            if (data.hasError) {
                return [
                    t.i({ className: "ri-error-warning-line txt-warning" }),
                    t.span({ className: "txt" }, "Failed to establish S3 connection"),
                ];
            }

            return [
                t.i({ className: "ri-checkbox-circle-line txt-success" }),
                t.span({ className: "txt" }, "S3 connected successfully"),
            ];
        },
    );
};
