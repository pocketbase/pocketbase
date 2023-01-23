// Simple Svelte scrollend detection action.
// ===================================================================
//
// ### Example usage
//
// Simple form (with default 100px threshold):
// ```html
// <div class="list" use:scrollend={() => { console.log("end reached") }}>
//     ...
// </div>
// ```
//
// With custom threshold:
// ```html
// <div class="list" use:scrollend={{
//     threshold: 10,
//     callback:  () => { console.log("end reached") }
// }}>
//     ...
// </div>
// ```
// ===================================================================

function normalize(rawData) {
    if (typeof rawData === "function") {
        return {
            threshold: 100,
            callback: rawData,
        }
    }

    return rawData || {};
}

export default function scrollend(node, options) {
    options = normalize(options);

    options?.callback && options.callback();

    function onScroll(e) {
        if (!options?.callback) {
            return;
        }

        const offset = e.target.scrollHeight - e.target.clientHeight - e.target.scrollTop;

        if (offset <= options.threshold) {
            options.callback();
        }
    }

    node.addEventListener("scroll", onScroll);
    node.addEventListener("resize", onScroll);

    return {
        update(newOptions) {
            options = normalize(newOptions);
        },
        destroy() {
            node.removeEventListener("scroll", onScroll);
            node.removeEventListener("resize", onScroll);
        },
    };
}
