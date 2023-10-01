<script>
    import { onMount, createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    let classes = "";
    export { classes as class }; // export reserved keyword

    export let vThreshold = 0;
    export let hThreshold = 0;
    export let dispatchOnNoScroll = true;

    let wrapper = null;
    let scrollClasses = "";
    let throttleTimeoutId = null;
    let hDiff;
    let vDiff;
    let wrapperWidth;
    let wrapperHeight;
    let observer;

    export function resetVerticalScroll() {
        if (!wrapper) {
            return;
        }

        wrapper.scrollTop = 0;
    }

    export function resetHorizontalScroll() {
        if (!wrapper) {
            return;
        }

        wrapper.scrollLeft = 0;
    }

    export function refresh() {
        if (!wrapper) {
            return;
        }

        scrollClasses = "";

        // +2 extra threshold as a workaround for the lack of subpixel precision
        wrapperWidth = wrapper.clientWidth + 2;
        wrapperHeight = wrapper.clientHeight + 2;

        hDiff = wrapper.scrollWidth - wrapperWidth;
        vDiff = wrapper.scrollHeight - wrapperHeight;

        // vertical scroller
        if (vDiff > 0) {
            scrollClasses += " v-scroll";

            if (vThreshold >= wrapperHeight) {
                vThreshold = 0;
            }

            if (wrapper.scrollTop - vThreshold <= 0) {
                scrollClasses += " v-scroll-start";
                dispatch("vScrollStart");
            }

            if (wrapper.scrollTop + vThreshold >= vDiff) {
                scrollClasses += " v-scroll-end";
                dispatch("vScrollEnd");
            }
        } else if (dispatchOnNoScroll) {
            dispatch("vScrollEnd");
        }

        // horizontal scroller
        if (hDiff > 0) {
            scrollClasses += " h-scroll";

            if (hThreshold >= wrapperWidth) {
                hThreshold = 0;
            }

            if (wrapper.scrollLeft - hThreshold <= 0) {
                scrollClasses += " h-scroll-start";
                dispatch("hScrollStart");
            }
            if (wrapper.scrollLeft + hThreshold >= hDiff) {
                scrollClasses += " h-scroll-end";
                dispatch("hScrollEnd");
            }
        } else if (dispatchOnNoScroll) {
            dispatch("hScrollEnd");
        }
    }

    export function throttleRefresh() {
        if (throttleTimeoutId) {
            return;
        }

        throttleTimeoutId = setTimeout(() => {
            refresh();

            throttleTimeoutId = null;
        }, 150);
    }

    onMount(() => {
        throttleRefresh();

        observer = new MutationObserver(throttleRefresh);

        observer.observe(wrapper, {
            attributeFilter: ["width", "height"],
            childList: true,
            subtree: true,
        });

        return () => {
            observer?.disconnect();
            clearTimeout(throttleTimeoutId);
        };
    });
</script>

<svelte:window on:resize={throttleRefresh} />

<div class="scroller-wrapper">
    <slot name="before" />

    <div bind:this={wrapper} class="scroller {classes} {scrollClasses}" on:scroll={throttleRefresh}>
        <slot />
    </div>

    <slot name="after" />
</div>

<style>
    .scroller {
        width: auto;
        min-height: 0;
        overflow: auto;
    }
    .scroller-wrapper {
        position: relative;
        min-height: 0;
    }
    :global(.scroller-wrapper .columns-dropdown) {
        top: 40px;
        z-index: 101;
        max-height: 340px;
    }
</style>
